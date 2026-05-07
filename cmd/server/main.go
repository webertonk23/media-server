package main

import (
	"log"

	"media-server/internal/api/routes"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.Load()

	database.Connect()

	libraryService := services.NewLibraryService()

	err := libraryService.Scan()

	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	routes.Setup(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
