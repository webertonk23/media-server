package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"media-server/internal/services"
)

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
	return streamViaFFmpeg(c, file.Path, start)
}

func streamMP4(c *fiber.Ctx, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot open file"})
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot stat file"})
	}

	fileSize := stat.Size()

	c.Set("Content-Type", "video/mp4")
	c.Set("Accept-Ranges", "bytes")
	c.Set("Content-Disposition", "inline")

	rangeHeader := c.Get("Range")
	if rangeHeader == "" {
		c.Set("Content-Length", fmt.Sprintf("%d", fileSize))
		c.Status(200)
		c.Response().SetBodyStream(f, int(fileSize))
		return nil
	}

	var start, end int64
	end = fileSize - 1

	rangeVal := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeVal, "-")

	if len(parts) != 2 {
		return c.Status(416).JSON(fiber.Map{"error": "invalid range"})
	}

	if parts[0] != "" {
		fmt.Sscanf(parts[0], "%d", &start)
	}
	if parts[1] != "" {
		fmt.Sscanf(parts[1], "%d", &end)
	}

	if start < 0 || end >= fileSize || start > end {
		c.Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		return c.Status(416).JSON(fiber.Map{"error": "range not satisfiable"})
	}

	chunkSize := end - start + 1

	_, err = f.Seek(start, io.SeekStart)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "seek failed"})
	}

	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Set("Content-Length", fmt.Sprintf("%d", chunkSize))
	c.Status(206)

	c.Response().SetBodyStream(io.LimitReader(f, chunkSize), int(chunkSize))
	return nil
}

func streamViaFFmpeg(c *fiber.Ctx, inputPath string, start string) error {
	codecEnv := os.Getenv("TRANSCODE_CODEC")
	if codecEnv == "" {
		codecEnv = "libx264"
	}

	codec := selectCodec(codecEnv)

	// Smart Remux: check if we can just copy the video stream
	if shouldCopyVideo(inputPath) {
		log.Printf("[Stream] Smart Remux: %s is already H264, using copy\n", inputPath)
		codec = "copy"
	}

	args := buildFFmpegArgs(inputPath, codec, start)

	log.Printf("[Stream] FFmpeg on-the-fly: %s → codec=%s (start=%s)\n", inputPath, codec, start)

	cmd := exec.Command("ffmpeg", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create ffmpeg pipe"})
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create ffmpeg stderr pipe"})
	}

	if err := cmd.Start(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to start ffmpeg"})
	}

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("[FFmpeg] %s\n", scanner.Text())
		}
	}()

	c.Response().Header.Set("Content-Type", "video/mp4")
	c.Response().Header.Set("Accept-Ranges", "none")
	c.Response().Header.Set("Content-Disposition", "inline")
	c.Response().Header.Set("Cache-Control", "no-cache")
	c.Response().Header.Set("Connection", "keep-alive")
	c.Status(200)

	log.Printf("[Stream] Starting stream for: %s\n", inputPath)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		log.Printf("[Stream] Body stream writer started for: %s\n", inputPath)
		n, err := io.Copy(w, stdout)
		if err != nil {
			log.Printf("[Stream] Error copying to response: %v\n", err)
		}
		log.Printf("[Stream] Finished copying %d bytes from FFmpeg to response\n", n)

		if err := cmd.Wait(); err != nil {
			log.Printf("[Stream] FFmpeg exit status: %v\n", err)
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
	codec := strings.TrimSpace(string(output))
	return codec == "h264"
}

func selectCodec(requested string) string {
	if requested == "libx264" {
		return "libx264"
	}
	testCmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "nullsrc=s=16x16:d=0.1",
		"-c:v", requested, "-f", "null", "-")
	testCmd.Stdout = io.Discard
	testCmd.Stderr = io.Discard

	if testCmd.Run() == nil {
		log.Printf("[Stream] Codec de hardware disponível: %s\n", requested)
		return requested
	}

	log.Printf("[Stream] Codec %s indisponível, usando libx264\n", requested)
	return "libx264"
}

func buildFFmpegArgs(inputPath, codec, start string) []string {
	args := []string{}
	if start != "" {
		args = append(args, "-ss", start)
	}

	args = append(args, "-i", inputPath)
	args = append(args, "-c:v", codec)

	if codec == "libx264" {
		args = append(args,
			"-preset", "veryfast",
			"-crf", "23",
			"-tune", "zerolatency",
			"-profile:v", "baseline",
			"-level", "3.0",
		)
	}

	movflags := "frag_keyframe+empty_moov+default_base_moof+omit_tfhd_offset+frag_discont+faststart"

	args = append(args,
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2",
		"-movflags", movflags,
		"-f", "mp4",
		"pipe:1",
	)

	return args
}
