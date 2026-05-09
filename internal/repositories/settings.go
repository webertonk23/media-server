package repositories

import (
	"media-server/internal/database"
	"media-server/internal/models"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) Get() (*models.Settings, error) {
	var settings models.Settings
	err := database.DB.First(&settings).Error
	if err != nil {
		// Retornar default se não existir
		return &models.Settings{MoviePath: "/app/media/movies", SeriesPath: "/app/media/series", ScanInterval: 1}, nil
	}
	return &settings, nil
}

func (r *SettingsRepository) Save(settings *models.Settings) error {
	return database.DB.Save(settings).Error
}
