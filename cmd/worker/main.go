package main

import (
	"log"
	"os"

	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/services"
)

func main() {
	config.Load()
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
