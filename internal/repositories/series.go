package repositories

import (
	"media-server/internal/database"
	"media-server/internal/models"
)

type SeriesRepository struct{}

func NewSeriesRepository() *SeriesRepository {
	return &SeriesRepository{}
}

func (r *SeriesRepository) Create(series *models.Series) error {
	return database.DB.Create(series).Error
}

func (r *SeriesRepository) Update(series *models.Series) error {
	return database.DB.Save(series).Error
}

func (r *SeriesRepository) FindByID(id uint) (*models.Series, error) {
	var series models.Series
	result := database.DB.Preload("MediaItem").First(&series, id)
	return &series, result.Error
}

func (r *SeriesRepository) FindByULID(ulid string) (*models.Series, error) {
	var series models.Series
	result := database.DB.Preload("MediaItem").Where("ulid = ?", ulid).First(&series)
	return &series, result.Error
}

func (r *SeriesRepository) FindByMediaItemID(mediaItemID uint) (*models.Series, error) {
	var series models.Series
	result := database.DB.Where("media_item_id = ?", mediaItemID).First(&series)
	return &series, result.Error
}
