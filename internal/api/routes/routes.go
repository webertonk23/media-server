package routes

import (
	"log"
	"os"
	"strings"
	"time"

	"media-server/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Setup(app *fiber.App) {
	app.Use(helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         31536000,
		ReferrerPolicy:     "no-referrer",
	}))

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:9000,https://localhost:9000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Muitas requisições. Tente novamente mais tarde.",
			})
		},
	}))

	api := app.Group("/api")

	// Rotas públicas (sem autenticação) - para frontend e app funcionarem
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Rotas de leitura - públicas para o app e frontend funcionarem
	api.Get("/movies", handlers.GetMovies)
	api.Get("/series", handlers.GetSeries)
	api.Get("/media/continue-watching", handlers.GetContinueWatching)
	api.Get("/media", handlers.GetMediaItems)
	api.Get("/media/:id", handlers.GetMediaItemByID)
	api.Get("/stream/:id", handlers.StreamMedia)
	api.Get("/progress/:id", handlers.GetProgress)
	api.Get("/series/:id/seasons", handlers.GetSeriesSeasons)
	api.Get("/seasons/:seasonId/episodes", handlers.GetSeasonEpisodes)
	api.Get("/transcode/status", handlers.GetTranscodeStatus)
	api.Get("/apk/info", handlers.GetAPKInfo)
	api.Get("/apk/download", handlers.DownloadAPK)

	// Rotas de escrita - públicas para o app funcionar
	api.Post("/progress/:id", handlers.UpdateProgress)

	// Middleware de autenticação para rotas administrativas
	// Se ADMIN_PASS não estiver configurado, permite acesso (modo desenvolvimento)
	adminPass := os.Getenv("ADMIN_PASS")
	
	var authMiddleware fiber.Handler
	
	if adminPass == "" {
		// Modo sem autenticação - apenas loga um aviso
		authMiddleware = func(c *fiber.Ctx) error {
			log.Println("⚠️  Acesso administrativo sem autenticação - configure ADMIN_PASS no .env")
			return c.Next()
		}
	} else {
		// Modo com autenticação
		authMiddleware = basicauth.New(basicauth.Config{
			Users: map[string]string{
				os.Getenv("ADMIN_USER"): adminPass,
			},
			Realm: "Restricted",
			Unauthorized: func(c *fiber.Ctx) error {
				return c.Status(401).JSON(fiber.Map{
					"error": "Autenticação necessária",
				})
			},
		})
	}

	// Rotas protegidas - apenas operações administrativas sensíveis
	protected := api.Group("", authMiddleware)

	scanLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})

	protected.Post("/scan", scanLimiter, handlers.ScanLibrary)
	protected.Get("/files/list", handlers.ListDirectory)
	protected.Get("/settings", handlers.GetSettings)
	protected.Post("/settings", handlers.SaveSettings)
	protected.Get("/logs", handlers.GetLogs)

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
