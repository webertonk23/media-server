package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"media-server/internal/dto"
	"media-server/internal/mappers"
	"media-server/internal/models"
	"media-server/internal/repositories"
	"media-server/internal/services"
)

// GetSeries lista todas as séries
func GetSeries(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	mediaItemService := services.NewMediaItemService()
	seriesService := services.NewSeriesService()

	items, total, err := mediaItemService.GetMediaItems(
		page,
		limit,
		search,
		models.MediaTypeSeries,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []dto.SeriesResponse

	for _, item := range items {
		// Buscar informações da série
		series, err := seriesService.GetByMediaItemID(item.ID)
		if err != nil {
			// Se não encontrar, criar resposta básica
			response = append(response, dto.SeriesResponse{
				ID:       item.ULID,
				Type:     item.Type,
				Title:    item.Title,
				Year:     item.Year,
				Overview: item.Overview,
				Poster:   item.Poster,
				Backdrop: item.Backdrop,
			})
			continue
		}

		response = append(response, mappers.ToSeriesResponse(item, *series))
	}

	return c.JSON(dto.PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Items: response,
	})
}

// GetSeriesSeasons lista as temporadas de uma série
func GetSeriesSeasons(c *fiber.Ctx) error {
	ulid := c.Params("id")

	mediaItemService := services.NewMediaItemService()
	seriesService := services.NewSeriesService()

	// Buscar MediaItem da série
	mediaItem, err := mediaItemService.GetByULID(ulid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "series not found",
		})
	}

	// Buscar Series
	series, err := seriesService.GetByMediaItemID(mediaItem.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "series data not found",
		})
	}

	// Buscar temporadas
	seasons, err := seriesService.GetSeasons(series.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []dto.SeasonResponse
	for _, season := range seasons {
		response = append(response, mappers.ToSeasonResponse(season))
	}

	return c.JSON(response)
}

// GetSeasonEpisodes lista os episódios de uma temporada
func GetSeasonEpisodes(c *fiber.Ctx) error {
	seasonULID := c.Params("seasonId")

	seriesService := services.NewSeriesService()
	seasonRepo := repositories.NewSeasonRepository()

	// Buscar Season
	season, err := seasonRepo.FindByULID(seasonULID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "season not found",
		})
	}

	// Buscar episódios
	episodes, err := seriesService.GetEpisodes(season.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []dto.EpisodeResponse
	for _, episode := range episodes {
		response = append(response, mappers.ToEpisodeResponse(episode, season.Number))
	}

	return c.JSON(response)
}
