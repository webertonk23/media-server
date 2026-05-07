package services

import (
	"media-server/internal/models"
	"media-server/internal/repositories"
)

type MediaItemService struct {
	repository *repositories.MediaItemRepository
}

func NewMediaItemService() *MediaItemService {
	return &MediaItemService{
		repository: repositories.NewMediaItemRepository(),
	}
}

func (s *MediaItemService) GetMediaItems(
	page int,
	limit int,
	search string,
	mediaType string,
) ([]models.MediaItem, int64, error) {
	return s.repository.Paginate(page, limit, search, mediaType)
}

func (s *MediaItemService) GetByID(id uint) (*models.MediaItem, error) {
	return s.repository.FindByID(id)
}

func (s *MediaItemService) GetByULID(ulid string) (*models.MediaItem, error) {
	return s.repository.FindByULID(ulid)
}

func (s *MediaItemService) Create(item *models.MediaItem) error {
	return s.repository.Create(item)
}

func (s *MediaItemService) Update(item *models.MediaItem) error {
	return s.repository.Update(item)
}
