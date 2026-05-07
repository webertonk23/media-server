package handlers

import (
	"github.com/gofiber/fiber/v2"
	"media-server/internal/services"
)

// ScanLibrary executa o scan de toda a biblioteca
func ScanLibrary(c *fiber.Ctx) error {
	libraryService := services.NewLibraryService()

	err := libraryService.ScanAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Library scan completed",
	})
}

// ScanMovies executa o scan apenas de filmes
func ScanMovies(c *fiber.Ctx) error {
	libraryService := services.NewLibraryService()

	err := libraryService.ScanMovies()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Movies scan completed",
	})
}

// ScanSeries executa o scan apenas de séries
func ScanSeries(c *fiber.Ctx) error {
	libraryService := services.NewLibraryService()

	err := libraryService.ScanSeries()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Series scan completed",
	})
}
