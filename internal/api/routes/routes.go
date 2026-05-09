package routes

import (
	"strings"

	"media-server/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// API routes group
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	api.Post("/scan", handlers.ScanLibrary)
	api.Post("/scan/movies", handlers.ScanMovies)
	api.Post("/scan/series", handlers.ScanSeries)

	api.Get("/media", handlers.GetMediaItems)
	api.Get("/media/:id", handlers.GetMediaItemByID)
	api.Get("/movies", handlers.GetMovies)
	api.Get("/series", handlers.GetSeries)
	api.Get("/stream/:id", handlers.StreamMedia)
	api.Post("/progress/:id", handlers.UpdateProgress)

	api.Get("/series/:id/seasons", handlers.GetSeriesSeasons)
	api.Get("/seasons/:seasonId/episodes", handlers.GetSeasonEpisodes)
	api.Get("/transcode/status", handlers.GetTranscodeStatus)
	api.Get("/files/list", handlers.ListDirectory)
	api.Get("/settings", handlers.GetSettings)
	api.Post("/settings", handlers.SaveSettings)
	api.Get("/logs", handlers.GetLogs)

	// Serve static files from frontend/dist
	app.Static("/", "./frontend/dist")

	// SPA fallback - serve index.html for all non-API routes
	app.Use(func(c *fiber.Ctx) error {
		// If the request starts with /api but didn't match any route, return 404
		if strings.HasPrefix(c.Path(), "/api") {
			return c.Status(404).JSON(fiber.Map{
				"error": "API route not found",
			})
		}
		// For all other routes, serve the frontend
		return c.SendFile("./frontend/dist/index.html")
	})

}
