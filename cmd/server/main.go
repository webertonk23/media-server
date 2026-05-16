package main

import (
	"log"
	"media-server/internal/api/routes"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/network"
	"media-server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()
	utils.InitLogger(500, "data/app.log")
	database.Connect()
	app := fiber.New()
	routes.Setup(app)
	network.StartDiscovery()
	log.Fatal(app.ListenTLS(
		":"+config.AppConfig.Port,
		"./certs/cert.pem",
		"./certs/key.pem",
	))
}
