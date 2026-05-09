package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"media-server/internal/services"
)

// StreamMedia faz streaming de um MediaItem pelo ULID
// Busca o melhor MediaFile associado (preferindo transcoded)
func StreamMedia(c *fiber.Ctx) error {
	ulid := c.Params("id")

	mediaItemService := services.NewMediaItemService()
	mediaFileService := services.NewMediaFileService()

	// Verificar se o MediaItem existe
	mediaItem, err := mediaItemService.GetByULID(ulid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "media item not found",
		})
	}

	// Buscar os arquivos associados ao MediaItem
	files, err := mediaFileService.GetByMediaItemID(mediaItem.ID)
	if err != nil || len(files) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "media file not found",
		})
	}

	// Selecionar o melhor arquivo: preferir status 'completed'
	var file = files[0]
	for _, f := range files {
		if f.Status == "completed" {
			file = f
			break
		}
	}

	// Verificar se o arquivo existe no sistema
	_, err = os.Stat(file.Path)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "file not found on disk",
		})
	}

	// Definir Content-Type baseado na extensão
	ext := strings.ToLower(filepath.Ext(file.Path))

	switch ext {
	case ".mp4":
		c.Set("Content-Type", "video/mp4")
	case ".mkv":
		c.Set("Content-Type", "video/x-matroska")
	case ".avi":
		c.Set("Content-Type", "video/x-msvideo")
	case ".webm":
		c.Set("Content-Type", "video/webm")
	default:
		c.Set("Content-Type", "application/octet-stream")
	}

	c.Set("Content-Disposition", "inline")

	return c.SendFile(file.Path, false)
}
