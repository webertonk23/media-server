package services
import (
	"media-server/internal/models"
	"media-server/internal/repositories"
)
type MediaFileService struct {
	repository *repositories.MediaFileRepository
}
func NewMediaFileService() *MediaFileService {
	return &MediaFileService{
		repository: repositories.NewMediaFileRepository(),
	}
}
func (s *MediaFileService) GetByID(id uint) (*models.MediaFile, error) {
	return s.repository.FindByID(id)
}
func (s *MediaFileService) GetByULID(ulid string) (*models.MediaFile, error) {
	return s.repository.FindByULID(ulid)
}
func (s *MediaFileService) GetByPath(path string) (*models.MediaFile, error) {
	return s.repository.FindByPath(path)
}
func (s *MediaFileService) GetByMediaItemID(mediaItemID uint) ([]models.MediaFile, error) {
	return s.repository.FindByMediaItemID(mediaItemID)
}
func (s *MediaFileService) Create(file *models.MediaFile) error {
	return s.repository.Create(file)
}
func (s *MediaFileService) Update(file *models.MediaFile) error {
	return s.repository.Update(file)
}
