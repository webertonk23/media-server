package routes

import (
	"media-server/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API routes - these are registered first to take precedence over static file serving
	// Rotas de scan
	app.Post("/scan", handlers.ScanLibrary)           // Escaneia tudo
	app.Post("/scan/movies", handlers.ScanMovies)     // Escaneia apenas filmes
	app.Post("/scan/series", handlers.ScanSeries)     // Escaneia apenas séries

	// Rotas de mídia (usando ULID)
	app.Get("/media", handlers.GetMediaItems)           // Lista todos os media items (filtrável por tipo)
	app.Get("/media/:id", handlers.GetMediaItemByID)    // Busca media item por ID
	app.Get("/movies", handlers.GetMovies)              // Lista apenas filmes
	app.Get("/series", handlers.GetSeries)              // Lista apenas séries
	app.Get("/stream/:id", handlers.StreamMedia)        // Stream de qualquer media item
	app.Post("/progress/:id", handlers.UpdateProgress)  // Atualizar progresso

	// Rotas de séries
	app.Get("/series/:id/seasons", handlers.GetSeriesSeasons)         // Lista temporadas de uma série
	app.Get("/seasons/:seasonId/episodes", handlers.GetSeasonEpisodes) // Lista episódios de uma temporada

	// Serve static files from frontend/dist
	// This is registered last so API routes take precedence
	app.Static("/", "./frontend/dist")

	// SPA fallback - serve index.html for all non-API routes
	// This allows Vue Router to handle client-side routing
	app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/dist/index.html")
	})

}
