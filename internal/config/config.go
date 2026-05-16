package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MediaPath    string
	IncomingPath string
	LibraryPath  string
	DBPath       string
	TMDBApiKey   string
	DBDriver     string
	AdminUser    string
	AdminPass    string
}

var AppConfig Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env não encontrado, usando valores padrão")
	}

	AppConfig = Config{
		Port:         getEnv("PORT", "9000"),
		MediaPath:    getEnv("MEDIA_PATH", "./media"),
		IncomingPath: getEnv("INCOMING_PATH", "./media/incoming"),
		LibraryPath:  getEnv("LIBRARY_PATH", "./media/library"),
		DBPath:       getEnv("DB_PATH", "./data/media.db"),
		TMDBApiKey:   getEnv("TMDB_API_KEY", ""),
		DBDriver:     getEnv("DB_DRIVER", "sqlite"),
		AdminUser:    getEnv("ADMIN_USER", "admin"),
		AdminPass:    getEnv("ADMIN_PASS", ""),
	}

	if AppConfig.AdminPass == "" {
		log.Println("⚠️  AVISO: ADMIN_PASS não configurado. Rotas administrativas estarão desprotegidas!")
		log.Println("⚠️  Configure ADMIN_USER e ADMIN_PASS no .env para proteger /scan, /settings e /logs")
	} else {
		log.Printf("✅ Autenticação configurada para usuário: %s", AppConfig.AdminUser)
	}

	if AppConfig.TMDBApiKey == "" {
		log.Println("AVISO: TMDB_API_KEY não configurado. Funcionalidades de metadata podem não funcionar.")
	}

	log.Printf("Configuração carregada: Port=%s, MediaPath=%s", AppConfig.Port, AppConfig.MediaPath)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
