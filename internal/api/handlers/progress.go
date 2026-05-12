package handlers
import (
	"github.com/gofiber/fiber/v2"
	"media-server/internal/dto"
	"media-server/internal/mappers"
	"media-server/internal/services"
)
func GetProgress(c *fiber.Ctx) error {
	ulid := c.Params("id")
	mediaItemService := services.NewMediaItemService()
	mediaItem, err := mediaItemService.GetByULID(ulid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "media item not found",
		})
	}
	progressService := services.NewProgressService()
	progress, err := progressService.GetByMediaItemID(mediaItem.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "progress not found",
		})
	}
	return c.JSON(fiber.Map{
		"position": progress.Position,
		"duration": progress.Duration,
		"finished": progress.Finished,
	})
}
func GetContinueWatching(c *fiber.Ctx) error {
	progressService := services.NewProgressService()
	progressList, err := progressService.GetContinueWatching()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var response []dto.ContinueWatchingResponse
	for _, p := range progressList {
		response = append(response, dto.ContinueWatchingResponse{
			Media:    mappers.ToMediaItemResponse(p.MediaItem),
			Position: p.Position,
			Duration: p.Duration,
			Finished: p.Finished,
		})
	}
	return c.JSON(response)
}
func UpdateProgress(c *fiber.Ctx) error {
	ulid := c.Params("id")
	var body dto.UpdateProgressRequest
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}
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
