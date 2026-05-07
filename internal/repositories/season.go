package repositories

import (
	"media-server/internal/database"
	"media-server/internal/models"
)

type SeasonRepository struct{}

func NewSeasonRepository() *SeasonRepository {
	return &SeasonRepository{}
}

func (r *SeasonRepository) Create(season *models.Season) error {
	return database.DB.Create(season).Error
}

func (r *SeasonRepository) Update(season *models.Season) error {
	return database.DB.Save(season).Error
}

func (r *SeasonRepository) FindByID(id uint) (*models.Season, error) {
	var season models.Season
	result := database.DB.First(&season, id)
	return &season, result.Error
}

func (r *SeasonRepository) FindByULID(ulid string) (*models.Season, error) {
	var season models.Season
	result := database.DB.Where("ulid = ?", ulid).First(&season)
	return &season, result.Error
}

func (r *SeasonRepository) FindBySeriesAndNumber(seriesID uint, number int) (*models.Season, error) {
	var season models.Season
	result := database.DB.Where("series_id = ? AND number = ?", seriesID, number).First(&season)
	return &season, result.Error
}

func (r *SeasonRepository) FindBySeriesID(seriesID uint) ([]models.Season, error) {
	var seasons []models.Season
	result := database.DB.Where("series_id = ?", seriesID).Order("number ASC").Find(&seasons)
	return seasons, result.Error
}
