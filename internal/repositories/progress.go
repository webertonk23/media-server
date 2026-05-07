package repositories

import (
	"time"

	"media-server/internal/database"
	"media-server/internal/models"
)

type ProgressRepository struct{}

func NewProgressRepository() *ProgressRepository {
	return &ProgressRepository{}
}

func (r *ProgressRepository) Save(
	progress *models.MediaProgress,
) error {

	var existing models.MediaProgress

	result := database.DB.
		Where("media_item_id = ?", progress.MediaItemID).
		First(&existing)

	if result.RowsAffected == 0 {
		return database.DB.Create(progress).Error
	}

	existing.Position = progress.Position
	existing.Duration = progress.Duration
	existing.Finished = progress.Finished
	existing.LastWatchedAt = time.Now()

	return database.DB.Save(&existing).Error
}
