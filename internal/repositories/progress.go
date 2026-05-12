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
func (r *ProgressRepository) GetByMediaItemID(
	mediaItemID uint,
) (*models.MediaProgress, error) {
	var progress models.MediaProgress
	err := database.DB.
		Where("media_item_id = ?", mediaItemID).
		First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}
func (r *ProgressRepository) GetContinueWatching() ([]models.MediaProgress, error) {
	var progressList []models.MediaProgress
	err := database.DB.
		Preload("MediaItem").
		Preload("MediaItem.Files").
		Where("finished = ?", false).
		Order("last_watched_at DESC").
		Find(&progressList).Error
	return progressList, err
}

func (r *ProgressRepository) DeleteByMediaItemID(mediaItemID uint) error {
	return database.DB.Where("media_item_id = ?", mediaItemID).Delete(&models.MediaProgress{}).Error
}


