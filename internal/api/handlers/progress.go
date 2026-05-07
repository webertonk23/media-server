package handlers

import (
	"github.com/gofiber/fiber/v2"

	"media-server/internal/dto"
	"media-server/internal/services"
)

func UpdateProgress(c *fiber.Ctx) error {
	ulid := c.Params("id")

	var body dto.UpdateProgressRequest

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	// Buscar MediaItem pelo ULID
	mediaItemService := services.NewMediaItemService()
	mediaItem, err := mediaItemService.GetByULID(ulid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "media item not found",
		})
	}

	progressService := services.NewProgressService()

	err = progressService.Save(
		mediaItem.ID,
		body.Position,
		body.Duration,
		body.Finished,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
