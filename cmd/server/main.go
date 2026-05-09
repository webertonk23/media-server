package main

import (
	"log"

	"media-server/internal/api/routes"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.Load()
	utils.InitLogger(500, "data/app.log")

	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
