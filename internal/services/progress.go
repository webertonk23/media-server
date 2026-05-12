package services
import (
	"time"
	"media-server/internal/models"
	"media-server/internal/repositories"
)
type ProgressService struct {
	progressRepository *repositories.ProgressRepository
}
func NewProgressService() *ProgressService {
	return &ProgressService{
		progressRepository: repositories.NewProgressRepository(),
	}
}
func (s *ProgressService) Save(
	mediaItemID uint,
	position int64,
	duration int64,
	finished bool,
) error {
	progress := models.MediaProgress{
		MediaItemID:   mediaItemID,
		Position:      position,
		Duration:      duration,
		Finished:      finished,
		LastWatchedAt: time.Now(),
	}
	return s.progressRepository.Save(&progress)
}
func (s *ProgressService) GetByMediaItemID(
	mediaItemID uint,
) (*models.MediaProgress, error) {
	return s.progressRepository.GetByMediaItemID(mediaItemID)
}
func (s *ProgressService) GetContinueWatching() ([]models.MediaProgress, error) {
	return s.progressRepository.GetContinueWatching()
}
