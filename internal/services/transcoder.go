package services

import (
	"context"
	"fmt"
	"log"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/models"
	"media-server/internal/repositories"
	"media-server/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	safeLanguageCode     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type TranscoderService struct {
	mediaFileRepo *repositories.MediaFileRepository
}

func NewTranscoderService() *TranscoderService {
	return &TranscoderService{
		mediaFileRepo: repositories.NewMediaFileRepository(),
	}
}
func (s *TranscoderService) StartWorker() {
	log.Println("Pipeline de Transcodificação Avançado iniciado...")
	os.MkdirAll(config.AppConfig.IncomingPath, 0755)
	os.MkdirAll(config.AppConfig.LibraryPath, 0755)
	for {
		file, err := s.mediaFileRepo.FindNextPending()
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Processando novo arquivo: %s\n", file.Path)
		err = s.ProcessPipeline(file)
		if err != nil {
			log.Printf("Erro na pipeline para %s: %v\n", file.Path, err)
			s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusError, err.Error())
		}
	}
}
func (s *TranscoderService) ProcessPipeline(file *models.MediaFile) error {
	s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusProcessing, "")
	info, err := utils.Probe(file.Path)
	if err != nil {
		return fmt.Errorf("falha ao analisar arquivo: %v", err)
	}
	needsVideoTranscode := !info.HasH264()
	audioStreams := info.GetAudioStreams()
	dir := filepath.Dir(file.Path)
	tempOutput := filepath.Join(dir, "proc_"+file.ULID+".mp4")
	err = s.runFFmpegMultiAudio(file.Path, tempOutput, needsVideoTranscode, audioStreams)
	if err != nil {
		return err
	}
	finalPath, err := s.determineFinalPath(file.Path, info)
	if err != nil {
		return err
	}
	finalDir := filepath.Dir(finalPath)
	err = os.MkdirAll(finalDir, 0755)
	if err != nil {
		return err
	}
	s.extractSubtitles(file.Path, finalPath, info)
	s.extractThumbnail(file.Path, file.MediaItemID)
	s.saveChapters(file.MediaItemID, info)
	originalPath := file.Path
	err = os.Rename(tempOutput, finalPath)
	if err != nil {
		return err
	}
	file.Path = finalPath
	file.Status = models.FileStatusCompleted
	stat, _ := os.Stat(finalPath)
	file.Size = stat.Size()
	err = s.mediaFileRepo.Update(file)
	if err == nil {
		log.Printf("[Pipeline] Sucesso! Removendo original: %s\n", originalPath)
		os.Remove(originalPath)
	}
	return err
}
func (s *TranscoderService) runFFmpegMultiAudio(input, output string, transcodeVideo bool, audioStreams []utils.FFStream) error {
	var args []string
	codec := os.Getenv("TRANSCODE_CODEC")
	if codec == "" {
		codec = "h264_vaapi"
	}
	if strings.Contains(codec, "vaapi") {
		args = append(args, "-vaapi_device", "/dev/dri/renderD129")
	}
	args = append(args, "-i", input)
	args = append(args, "-map", "0:v:0")
	if transcodeVideo {
		log.Println("[Pipeline] Transcodificando vídeo...")
		if strings.Contains(codec, "vaapi") {
			args = append(args, "-vf", "format=nv12|vaapi,hwupload", "-c:v", codec)
		} else {
			args = append(args, "-c:v", "libx264", "-preset", "fast", "-crf", "23")
		}
	} else {
		log.Println("[Pipeline] Remuxing vídeo (Copy)...")
		args = append(args, "-c:v", "copy")
	}
	for i, stream := range audioStreams {
		args = append(args, "-map", fmt.Sprintf("0:a:%d", i))
		needsAudioTranscode := !(stream.CodecName == "aac" || stream.CodecName == "mp3" || stream.CodecName == "ac3")
		if needsAudioTranscode {
			log.Printf("[Pipeline] Transcodificando áudio %d (%s) para AAC...\n", i, stream.CodecName)
			args = append(args, fmt.Sprintf("-c:a:%d", i), "aac", fmt.Sprintf("-b:a:%d", i), "128k")
		} else {
			log.Printf("[Pipeline] Remuxing áudio %d (%s) (Copy)...\n", i, stream.CodecName)
			args = append(args, fmt.Sprintf("-c:a:%d", i), "copy")
		}
	}
	args = append(args, "-pix_fmt", "yuv420p", "-async", "1", "-fps_mode", "cfr", "-movflags", "+faststart", "-y", output)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("transcodificação excedeu o tempo limite de 30 minutos")
	}
	
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v, output: %s", err, string(out))
	}
	return nil
}
func (s *TranscoderService) extractThumbnail(input string, mediaItemID uint) {
	output := filepath.Join(config.AppConfig.MediaPath, "thumbnails", fmt.Sprintf("%d.jpg", mediaItemID))
	os.MkdirAll(filepath.Dir(output), 0755)
	log.Printf("[Pipeline] Extraindo thumbnail para MediaItem %d...\n", mediaItemID)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "ffmpeg", "-ss", "00:00:10", "-i", input, "-vframes", "1", "-q:v", "2", "-y", output)
	err := cmd.Run()
	
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[Pipeline] Timeout ao extrair thumbnail, tentando posição alternativa\n")
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel2()
		exec.CommandContext(ctx2, "ffmpeg", "-ss", "00:00:01", "-i", input, "-vframes", "1", "-q:v", "2", "-y", output).Run()
	}
}
func (s *TranscoderService) saveChapters(mediaItemID uint, info *utils.FFProbeOutput) {
	if len(info.Chapters) == 0 {
		return
	}
	log.Printf("[Pipeline] Salvando %d capítulos para MediaItem %d...\n", len(info.Chapters), mediaItemID)
	database.DB.Where("media_item_id = ?", mediaItemID).Delete(&models.Chapter{})
	for _, c := range info.Chapters {
		title := c.Tags.Title
		if title == "" {
			title = fmt.Sprintf("Capítulo %d", c.ID+1)
		}
		chapter := models.Chapter{
			MediaItemID: mediaItemID,
			Title:       title,
			StartTime:   c.GetStartTime(),
			EndTime:     c.GetEndTime(),
		}
		database.DB.Create(&chapter)
	}
}
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	
	name = invalidFilenameChars.ReplaceAllString(name, "_")
	
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")
	
	if len(name) > 255 {
		name = name[:255]
	}
	
	name = strings.Trim(name, "_")
	
	if name == "" {
		name = "unnamed"
	}
	
	return name
}

func (s *TranscoderService) extractSubtitles(input string, finalVideoPath string, info *utils.FFProbeOutput) {
	subs := info.GetSubtitles()
	basePath := strings.TrimSuffix(finalVideoPath, filepath.Ext(finalVideoPath))

	absLibraryPath, _ := filepath.Abs(config.AppConfig.LibraryPath)
	
	for i, sub := range subs {
		if sub.Index < 0 || sub.Index > 100 {
			log.Printf("[Pipeline] Índice de legenda inválido: %d\n", sub.Index)
			continue
		}

		lang := sub.Tags.Language
		if lang == "" || !safeLanguageCode.MatchString(lang) {
			lang = fmt.Sprintf("track%d", i)
		}
		
		outputSub := basePath + "." + lang + ".vtt"
		
		absOutputSub, _ := filepath.Abs(outputSub)
		if !strings.HasPrefix(absOutputSub, absLibraryPath) {
			log.Printf("[Pipeline] Tentativa de escrever legenda fora do diretório permitido: %s\n", outputSub)
			continue
		}
		
		log.Printf("[Pipeline] Extraindo legenda para destino final: %s\n", outputSub)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", input, "-map", fmt.Sprintf("0:%d", sub.Index), "-y", outputSub)
		err := cmd.Run()
		cancel()
		
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("[Pipeline] Timeout ao extrair legenda: %s\n", outputSub)
			} else {
				log.Printf("[Pipeline] Erro ao extrair legenda: %v\n", err)
			}
		}
	}
}
func (s *TranscoderService) determineFinalPath(originalPath string, info *utils.FFProbeOutput) (string, error) {
	metadata := info.GetMetadata()
	parsedSeries := utils.ParseSeriesFilename(originalPath)
	if parsedSeries.IsSeries {
		title := metadata.Show
		if title == "" {
			title = parsedSeries.Title
		}
		title = sanitizeFilename(title)
		
		season := parsedSeries.Season
		if season == 0 && metadata.Season != "" {
			fmt.Sscanf(metadata.Season, "%d", &season)
		}
		if season == 0 {
			season = 1
		}
		episode := parsedSeries.Episode
		if episode == 0 && metadata.Episode != "" {
			fmt.Sscanf(metadata.Episode, "%d", &episode)
		}
		folderName := fmt.Sprintf("%s/Season %02d", title, season)
		fileName := fmt.Sprintf("%s - S%02dE%02d.mp4", title, season, episode)
		return filepath.Join(config.AppConfig.LibraryPath, folderName, fileName), nil
	}
	parsedMovie := utils.ParseMovieFilename(originalPath)
	title := metadata.Title
	if title == "" {
		title = parsedMovie.Title
	}
	title = sanitizeFilename(title)
	
	year := parsedMovie.Year
	folderName := title
	if year > 0 {
		folderName = fmt.Sprintf("%s (%d)", title, year)
	}
	fileName := fmt.Sprintf("%s.mp4", folderName)
	return filepath.Join(config.AppConfig.LibraryPath, folderName, fileName), nil
}
