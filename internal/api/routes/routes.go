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

	// Rotas de scan
	api.Post("/scan", handlers.ScanLibrary)           // Escaneia tudo
	api.Post("/scan/movies", handlers.ScanMovies)     // Escaneia apenas filmes
	api.Post("/scan/series", handlers.ScanSeries)     // Escaneia apenas séries

	// Rotas de mídia (usando ULID)
	api.Get("/media", handlers.GetMediaItems)           // Lista todos os media items (filtrável por tipo)
	api.Get("/media/:id", handlers.GetMediaItemByID)    // Busca media item por ID
	api.Get("/movies", handlers.GetMovies)              // Lista apenas filmes
	api.Get("/series", handlers.GetSeries)              // Lista apenas séries
	api.Get("/stream/:id", handlers.StreamMedia)        // Stream de qualquer media item
	api.Post("/progress/:id", handlers.UpdateProgress)  // Atualizar progresso

	// Rotas de séries
	api.Get("/series/:id/seasons", handlers.GetSeriesSeasons)         // Lista temporadas de uma série
	api.Get("/seasons/:seasonId/episodes", handlers.GetSeasonEpisodes) // Lista episódios de uma temporada

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
