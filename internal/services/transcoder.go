package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"media-server/internal/models"
	"media-server/internal/repositories"
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
	log.Println("Transcoder Worker iniciado. Aguardando arquivos pendentes...")
	for {
		file, err := s.mediaFileRepo.FindNextPending()
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Novo arquivo para transcodificar encontrado: %s\n", file.Path)
		err = s.ProcessFile(file)
		if err != nil {
			log.Printf("Erro ao processar arquivo %d (%s): %v\n", file.ID, file.Path, err)
			s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusError, err.Error())
		}
	}
}

func (s *TranscoderService) ProcessFile(file *models.MediaFile) error {
	err := s.mediaFileRepo.UpdateStatus(file.ID, models.FileStatusProcessing, "")
	if err != nil {
		return err
	}

	dir := filepath.Dir(file.Path)
	tempOutputPath := filepath.Join(dir, fmt.Sprintf("transcoding_%s.mp4", file.ULID))

	codecEnv := os.Getenv("TRANSCODE_CODEC")
	presetEnv := os.Getenv("TRANSCODE_PRESET")
	if presetEnv == "" {
		presetEnv = "fast"
	}

	var codecs []string
	if codecEnv != "" && codecEnv != "libx264" {
		codecs = append(codecs, codecEnv)
	}
	codecs = append(codecs, "libx264")

	var cmd *exec.Cmd
	var codecUsed string
	var stat os.FileInfo

	for _, codec := range codecs {
		log.Printf("Tentando transcodificar com codec: %s\n", codec)
		
		args := []string{"-i", file.Path, "-c:v", codec}
		if codec == "libx264" {
			args = append(args, "-preset", presetEnv, "-crf", "23")
		}
		args = append(args, "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "128k", "-movflags", "+faststart", "-y", tempOutputPath)

		cmd = exec.Command("ffmpeg", args...)
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			stat, err = os.Stat(tempOutputPath)
			if err == nil && stat.Size() > 0 {
				codecUsed = codec
				break
			} else {
				log.Printf("Arquivo de saída vazio ou não gerado com codec %s\n", codec)
			}
		} else {
			log.Printf("Falha ao usar %s: %s\n", codec, string(output))
		}
		
		os.Remove(tempOutputPath)
	}

	if codecUsed == "" {
		return fmt.Errorf("todos os codecs falharam")
	}

	originalPath := file.Path
	ext := strings.ToLower(filepath.Ext(file.Path))
	finalPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(originalPath), ext)+".mp4")

	if finalPath == originalPath {
		os.Remove(originalPath)
	}

	err = os.Rename(tempOutputPath, finalPath)
	if err != nil {
		return fmt.Errorf("falha ao renomear: %v", err)
	}

	if finalPath != originalPath {
		os.Remove(originalPath)
	}

	file.Path = finalPath
	file.OriginalPath = originalPath
	file.Status = models.FileStatusCompleted
	file.Size = stat.Size()
	
	return s.mediaFileRepo.Update(file)
}
