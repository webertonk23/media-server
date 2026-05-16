package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"media-server/internal/models"
	"media-server/internal/repositories"
)

func GetSettings(c *fiber.Ctx) error {
	repo := repositories.NewSettingsRepository()
	settings, err := repo.Get()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar config"})
	}
	return c.JSON(settings)
}

func SaveSettings(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return c.Status(415).JSON(fiber.Map{
			"error": "Content-Type deve ser application/json",
		})
	}

	var settings models.Settings
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
	}

	repo := repositories.NewSettingsRepository()
	if err := repo.Save(&settings); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao salvar"})
	}

	return c.JSON(settings)
}
