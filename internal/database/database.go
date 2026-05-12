package database
import (
	"log"
	"media-server/internal/config"
	"media-server/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
var DB *gorm.DB
func Connect() {
	var err error
	driver := config.AppConfig.DBDriver
	switch driver {
	case "sqlite":
		DB, err = gorm.Open(
			sqlite.Open(config.AppConfig.DBPath),
			&gorm.Config{},
		)
	default:
		log.Fatal("Driver de banco de dados não suportado")
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Banco conectado")
	err = DB.AutoMigrate(
		&models.MediaItem{},
		&models.MediaFile{},
		&models.MediaProgress{},
		&models.Series{},
		&models.Season{},
		&models.Episode{},
		&models.Settings{},
		&models.Chapter{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
