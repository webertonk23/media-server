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
}
var AppConfig Config
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env não encontrado")
	}
	AppConfig = Config{
		Port:         getEnv("PORT", "9000"),
		MediaPath:    getEnv("MEDIA_PATH", "/media"),
		IncomingPath: getEnv("INCOMING_PATH", "/media/incoming"),
		LibraryPath:  getEnv("LIBRARY_PATH", "/media/library"),
		DBPath:       getEnv("DB_PATH", "./data/media.db"),
		TMDBApiKey:   getEnv("TMDB_API_KEY", ""),
		DBDriver:     getEnv("DB_DRIVER", "sqlite"),
	}
}
func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
