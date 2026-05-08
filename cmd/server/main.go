package main

import (
	"log"

	"media-server/internal/api/routes"
	"media-server/internal/config"
	"media-server/internal/database"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.Load()

	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
