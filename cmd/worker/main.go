package main

import (
	"log"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/services"
	"media-server/internal/utils"
	"os"
)

func main() {
	config.Load()
	utils.InitLogger(100, "data/app.log")
	database.Connect()
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run cmd/worker/main.go [scanner|transcoder]")
	}
	workerType := os.Args[1]
	switch workerType {
	case "scanner":
		libraryService := services.NewLibraryService()
		libraryService.StartScannerWorker()
	case "transcoder":
		transcoderService := services.NewTranscoderService()
		transcoderService.StartWorker()
	default:
		log.Fatalf("Tipo de worker desconhecido: %s", workerType)
	}
}
