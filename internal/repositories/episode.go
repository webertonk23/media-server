package repositories
import (
	"media-server/internal/database"
	"media-server/internal/models"
)
type EpisodeRepository struct{}
func NewEpisodeRepository() *EpisodeRepository {
	return &EpisodeRepository{}
}
func (r *EpisodeRepository) Create(episode *models.Episode) error {
	return database.DB.Create(episode).Error
}
func (r *EpisodeRepository) Update(episode *models.Episode) error {
	return database.DB.Save(episode).Error
}
func (r *EpisodeRepository) FindByID(id uint) (*models.Episode, error) {
	var episode models.Episode
	result := database.DB.Preload("MediaItem.Files").Preload("Season").First(&episode, id)
	return &episode, result.Error
}
func (r *EpisodeRepository) FindByULID(ulid string) (*models.Episode, error) {
	var episode models.Episode
	result := database.DB.Preload("MediaItem.Files").Preload("Season").Where("ulid = ?", ulid).First(&episode)
	return &episode, result.Error
}
func (r *EpisodeRepository) FindByMediaItemID(mediaItemID uint) (*models.Episode, error) {
	var episode models.Episode
	result := database.DB.Preload("MediaItem.Files").Where("media_item_id = ?", mediaItemID).First(&episode)
	return &episode, result.Error
}
func (r *EpisodeRepository) FindBySeasonAndNumber(seasonID uint, number int) (*models.Episode, error) {
	var episode models.Episode
	result := database.DB.Preload("MediaItem.Files").Where("season_id = ? AND number = ?", seasonID, number).First(&episode)
	return &episode, result.Error
}
func (r *EpisodeRepository) FindBySeasonID(seasonID uint) ([]models.Episode, error) {
	var episodes []models.Episode
	result := database.DB.Preload("MediaItem.Files").Where("season_id = ?", seasonID).Order("number ASC").Find(&episodes)
	return episodes, result.Error
}
