package routes
import (
	"strings"
	"media-server/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)
func Setup(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
	api.Post("/scan", handlers.ScanLibrary)
	api.Get("/movies", handlers.GetMovies)
	api.Get("/series", handlers.GetSeries)
	api.Get("/media/continue-watching", handlers.GetContinueWatching)
	api.Get("/media", handlers.GetMediaItems)
	api.Get("/media/:id", handlers.GetMediaItemByID)
	api.Get("/stream/:id", handlers.StreamMedia)
	api.Get("/progress/:id", handlers.GetProgress)
	api.Post("/progress/:id", handlers.UpdateProgress)
	api.Get("/series/:id/seasons", handlers.GetSeriesSeasons)
	api.Get("/seasons/:seasonId/episodes", handlers.GetSeasonEpisodes)
	api.Get("/transcode/status", handlers.GetTranscodeStatus)
	api.Get("/files/list", handlers.ListDirectory)
	api.Get("/settings", handlers.GetSettings)
	api.Post("/settings", handlers.SaveSettings)
	api.Get("/logs", handlers.GetLogs)

	api.Get("/apk/info", handlers.GetAPKInfo)
	api.Get("/apk/download", handlers.DownloadAPK)

	app.Static("/", "./frontend/dist")
	app.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api") {
			return c.Status(404).JSON(fiber.Map{
				"error": "API route not found",
			})
		}
		return c.SendFile("./frontend/dist/index.html")
	})
}
