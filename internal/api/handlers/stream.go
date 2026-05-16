package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"context"
	"time"

	"media-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

var (
	hwCodec string
	once    sync.Once
	timeRegex  = regexp.MustCompile(`^\d{1,2}:\d{2}:\d{2}(\.\d+)?$`)
	audioRegex = regexp.MustCompile(`^\d+$`)
)

type limitReadCloser struct {
	io.Reader
	closer io.Closer
}

func (l limitReadCloser) Close() error {
	return l.closer.Close()
}

func StreamMedia(c *fiber.Ctx) error {
	ulid := c.Params("id")
	mediaItemService := services.NewMediaItemService()
	mediaFileService := services.NewMediaFileService()
	mediaItem, err := mediaItemService.GetByULID(ulid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "media item not found"})
	}
	files, err := mediaFileService.GetByMediaItemID(mediaItem.ID)
	if err != nil || len(files) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "media file not found"})
	}
	file := files[0]
	for _, f := range files {
		if f.Status == "completed" {
			file = f
			break
		}
	}
	if _, err := os.Stat(file.Path); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "file not found on disk"})
	}
	ext := strings.ToLower(filepath.Ext(file.Path))
	if ext == ".mp4" {
		return streamMP4(c, file.Path)
	}
	start := c.Query("start")
	audio := c.Query("audio")
	return streamViaFFmpeg(c, file.Path, start, audio)
}

func streamMP4(c *fiber.Ctx, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot open file"})
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return c.Status(500).JSON(fiber.Map{"error": "cannot stat file"})
	}
	fileSize := stat.Size()
	c.Set("Content-Type", "video/mp4")
	c.Set("Accept-Ranges", "bytes")
	c.Set("Content-Disposition", "inline")

	if c.Method() == "HEAD" {
		c.Set("Content-Length", fmt.Sprintf("%d", fileSize))
		return c.SendStatus(200)
	}

	rangeHeader := c.Get("Range")
	if rangeHeader == "" {
		c.Set("Content-Length", fmt.Sprintf("%d", fileSize))
		c.Status(200)
		c.Response().SetBodyStream(f, -1)
		return nil
	}

	var start, end int64
	end = fileSize - 1
	rangeVal := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeVal, "-")
	if len(parts) != 2 {
		f.Close()
		return c.Status(416).JSON(fiber.Map{"error": "invalid range"})
	}
	if parts[0] != "" {
		fmt.Sscanf(parts[0], "%d", &start)
	}
	if parts[1] != "" {
		fmt.Sscanf(parts[1], "%d", &end)
	}
	if start < 0 || end >= fileSize || start > end {
		f.Close()
		c.Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		return c.Status(416).JSON(fiber.Map{"error": "range not satisfiable"})
	}
	chunkSize := end - start + 1
	f.Seek(start, io.SeekStart)

	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Set("Content-Length", fmt.Sprintf("%d", chunkSize))
	c.Status(206)
	c.Response().SetBodyStream(limitReadCloser{io.LimitReader(f, chunkSize), f}, int(chunkSize))
	return nil
}

func streamViaFFmpeg(c *fiber.Ctx, inputPath string, start string, audio string) error {
	once.Do(func() {
		hwCodec = os.Getenv("TRANSCODE_CODEC")
		if hwCodec == "" {
			hwCodec = "libx264"
		}
		if hwCodec != "libx264" {
			args := []string{"-f", "lavfi", "-i", "nullsrc=s=16x16:d=0.1", "-c:v", hwCodec, "-f", "null", "-"}
			if strings.Contains(hwCodec, "vaapi") {
				args = append([]string{"-vaapi_device", "/dev/dri/renderD129"}, args...)
			}
			testCmd := exec.Command("ffmpeg", args...)
			if testCmd.Run() != nil {
				log.Printf("[Stream] Codec HW %s falhou, fallback para libx264\n", hwCodec)
				hwCodec = "libx264"
			}
		}
	})

	codec := hwCodec
	if shouldCopyVideo(inputPath) {
		codec = "copy"
	}

	args, err := buildFFmpegArgs(inputPath, codec, start, audio)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create ffmpeg pipe"})
	}

	if err := cmd.Start(); err != nil {
		stdout.Close()
		return c.Status(500).JSON(fiber.Map{"error": "failed to start ffmpeg"})
	}

	c.Response().Header.Set("Content-Type", "video/mp4")
	c.Response().Header.Set("Cache-Control", "no-cache")
	c.Status(200)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer stdout.Close()
		defer func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			cmd.Wait()
		}()
		
		io.Copy(w, stdout)
		
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("[Stream] Timeout de 2 horas excedido")
		}
	})
	return nil
}

func shouldCopyVideo(path string) bool {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "h264"
}

func buildFFmpegArgs(inputPath, codec, start, audio string) ([]string, error) {
	args := []string{}
	
	if strings.Contains(codec, "vaapi") {
		args = append(args, "-vaapi_device", "/dev/dri/renderD129")
	}

	if start != "" {
		if !timeRegex.MatchString(start) {
			return nil, fmt.Errorf("formato de tempo inválido: deve ser HH:MM:SS")
		}
		args = append(args, "-ss", start)
	}
	
	args = append(args, "-i", inputPath, "-map", "0:v:0")

	if audio != "" {
		if !audioRegex.MatchString(audio) {
			return nil, fmt.Errorf("índice de áudio inválido: deve ser um número")
		}
		var audioIdx int
		fmt.Sscanf(audio, "%d", &audioIdx)
		if audioIdx < 0 || audioIdx > 32 {
			return nil, fmt.Errorf("índice de áudio fora do range permitido (0-32)")
		}
		args = append(args, "-map", fmt.Sprintf("0:a:%d", audioIdx))
	} else {
		args = append(args, "-map", "0:a:0")
	}
	
	if strings.Contains(codec, "vaapi") {
		args = append(args, "-vf", "format=nv12,hwupload")
	}
	
	args = append(args, "-c:v", codec)
	
	if codec == "libx264" {
		args = append(args, "-preset", "veryfast", "-crf", "23", "-tune", "zerolatency", "-profile:v", "baseline", "-level", "3.0")
	}
	
	args = append(args, "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "128k", "-ac", "2", "-movflags", "frag_keyframe+empty_moov+default_base_moof+omit_tfhd_offset+frag_discont+faststart", "-f", "mp4", "pipe:1")
	
	return args, nil
}
