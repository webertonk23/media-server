package handlers
import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"media-server/internal/dto"
	"media-server/internal/mappers"
	"media-server/internal/models"
	"media-server/internal/services"
)
func GetMediaItems(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")
	mediaType := c.Query("type", "") 
	mediaItemService := services.NewMediaItemService()
	items, total, err := mediaItemService.GetMediaItems(
		page,
		limit,
		search,
		mediaType,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var response []dto.MediaItemResponse
	for _, item := range items {
		response = append(
			response,
			mappers.ToMediaItemResponse(item),
		)
	}
	return c.JSON(dto.PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Items: response,
	})
}
func GetMediaItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "ID é obrigatório",
		})
	}
	mediaItemService := services.NewMediaItemService()
	item, err := mediaItemService.GetByULID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Mídia não encontrada",
		})
	}
	return c.JSON(mappers.ToMediaItemResponse(*item))
}
func GetMovies(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")
	mediaItemService := services.NewMediaItemService()
	items, total, err := mediaItemService.GetMediaItems(
		page,
		limit,
		search,
		models.MediaTypeMovie,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var response []dto.MediaItemResponse
	for _, item := range items {
		response = append(
			response,
			mappers.ToMediaItemResponse(item),
		)
	}
	return c.JSON(dto.PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Items: response,
	})
}
