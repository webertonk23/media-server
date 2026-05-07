package repositories

import (
	"media-server/internal/database"
	"media-server/internal/models"
	"strings"
)

type MediaItemRepository struct{}

func NewMediaItemRepository() *MediaItemRepository {
	return &MediaItemRepository{}
}

func (r *MediaItemRepository) Create(item *models.MediaItem) error {
	return database.DB.Create(item).Error
}

func (r *MediaItemRepository) Update(item *models.MediaItem) error {
	return database.DB.Save(item).Error
}

func (r *MediaItemRepository) FindByID(id uint) (*models.MediaItem, error) {
	var item models.MediaItem
	result := database.DB.First(&item, id)
	return &item, result.Error
}

func (r *MediaItemRepository) FindByULID(ulid string) (*models.MediaItem, error) {
	var item models.MediaItem
	result := database.DB.Where("ulid = ?", ulid).First(&item)
	return &item, result.Error
}

func (r *MediaItemRepository) FindByTMDBID(tmdbID int, mediaType string) (*models.MediaItem, error) {
	var item models.MediaItem
	result := database.DB.Where("tmdb_id = ? AND type = ?", tmdbID, mediaType).First(&item)
	return &item, result.Error
}

func (r *MediaItemRepository) Paginate(
	page int,
	limit int,
	search string,
	mediaType string,
) ([]models.MediaItem, int64, error) {
	var items []models.MediaItem
	var total int64

	query := database.DB.Model(&models.MediaItem{})

	// Filtrar por tipo de mídia
	if mediaType != "" {
		query = query.Where("type = ?", mediaType)
	}

	// Filtrar por busca
	if search != "" {
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(original_title) LIKE ?",
			"%"+strings.ToLower(search)+"%",
			"%"+strings.ToLower(search)+"%",
		)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&items).
		Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
