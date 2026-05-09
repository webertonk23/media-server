package handlers

import (
	"github.com/gofiber/fiber/v2"
	"media-server/internal/mappers"
	"media-server/internal/models"
	"media-server/internal/repositories"
)

func GetTranscodeStatus(c *fiber.Ctx) error {
	repo := repositories.NewMediaFileRepository()

	pending, err := repo.FindByStatus(models.FileStatusPending)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar pendentes"})
	}

	processing, err := repo.FindByStatus(models.FileStatusProcessing)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar processando"})
	}

	return c.JSON(fiber.Map{
		"pending":    mappers.ToMediaFileResponses(pending),
		"processing": mappers.ToMediaFileResponses(processing),
	})
}
