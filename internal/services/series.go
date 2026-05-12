package services
import (
	"media-server/internal/models"
	"media-server/internal/repositories"
)
type SeriesService struct {
	seriesRepository  *repositories.SeriesRepository
	seasonRepository  *repositories.SeasonRepository
	episodeRepository *repositories.EpisodeRepository
}
func NewSeriesService() *SeriesService {
	return &SeriesService{
		seriesRepository:  repositories.NewSeriesRepository(),
		seasonRepository:  repositories.NewSeasonRepository(),
		episodeRepository: repositories.NewEpisodeRepository(),
	}
}
func (s *SeriesService) GetByULID(ulid string) (*models.Series, error) {
	return s.seriesRepository.FindByULID(ulid)
}
func (s *SeriesService) GetByMediaItemID(mediaItemID uint) (*models.Series, error) {
	return s.seriesRepository.FindByMediaItemID(mediaItemID)
}
func (s *SeriesService) GetSeasons(seriesID uint) ([]models.Season, error) {
	return s.seasonRepository.FindBySeriesID(seriesID)
}
func (s *SeriesService) GetEpisodes(seasonID uint) ([]models.Episode, error) {
	return s.episodeRepository.FindBySeasonID(seasonID)
}
func (s *SeriesService) GetEpisodeByULID(ulid string) (*models.Episode, error) {
	return s.episodeRepository.FindByULID(ulid)
}
func (s *SeriesService) Create(series *models.Series) error {
	return s.seriesRepository.Create(series)
}
func (s *SeriesService) CreateSeason(season *models.Season) error {
	return s.seasonRepository.Create(season)
}
func (s *SeriesService) CreateEpisode(episode *models.Episode) error {
	return s.episodeRepository.Create(episode)
}
func (s *SeriesService) FindOrCreateSeason(seriesID uint, number int) (*models.Season, error) {
	season, err := s.seasonRepository.FindBySeriesAndNumber(seriesID, number)
	if err == nil {
		return season, nil
	}
	season = &models.Season{
		SeriesID: seriesID,
		Number:   number,
	}
	err = s.seasonRepository.Create(season)
	if err != nil {
		return nil, err
	}
	return season, nil
}
func (s *SeriesService) FindOrCreateEpisode(seasonID uint, number int, mediaItemID uint) (*models.Episode, error) {
	episode, err := s.episodeRepository.FindBySeasonAndNumber(seasonID, number)
	if err == nil {
		return episode, nil
	}
	episode = &models.Episode{
		SeasonID:    seasonID,
		Number:      number,
		MediaItemID: mediaItemID,
	}
	err = s.episodeRepository.Create(episode)
	if err != nil {
		return nil, err
	}
	return episode, nil
}
