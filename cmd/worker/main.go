package main
import (
	"log"
	"os"
	"media-server/internal/config"
	"media-server/internal/database"
	"media-server/internal/repositories"
	"media-server/internal/services"
	"media-server/internal/utils"
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
	case "cleanup":
		repo := repositories.NewMediaFileRepository()
		count, err := repo.DeleteOrphanedMediaFiles()
		if err != nil {
			log.Fatalf("Erro ao limpar registros orfãos: %v", err)
		}
		log.Printf("Limpeza concluída: %d arquivos removidos.", count)
	default:
		log.Fatalf("Tipo de worker desconhecido: %s", workerType)
	}
}
