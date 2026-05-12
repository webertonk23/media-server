package services
import (
	"log"
	"os"
	"time"
	"media-server/internal/config"
	"media-server/internal/models"
	"media-server/internal/repositories"
	"media-server/internal/scanner"
	"media-server/internal/utils"
)
type LibraryService struct {
	mediaItemRepository *repositories.MediaItemRepository
	mediaFileRepository *repositories.MediaFileRepository
	metadataService     *MetadataService
	seriesRepository    *repositories.SeriesRepository
	seasonRepository    *repositories.SeasonRepository
	episodeRepository   *repositories.EpisodeRepository
}
func NewLibraryService() *LibraryService {
	return &LibraryService{
		mediaItemRepository: repositories.NewMediaItemRepository(),
		mediaFileRepository: repositories.NewMediaFileRepository(),
		metadataService:     NewMetadataService(),
		seriesRepository:    repositories.NewSeriesRepository(),
		seasonRepository:    repositories.NewSeasonRepository(),
		episodeRepository:   repositories.NewEpisodeRepository(),
	}
}
func (s *LibraryService) StartScannerWorker() {
	log.Printf("Scanner Worker iniciado. Monitorando entrada: %s\n", config.AppConfig.IncomingPath)
	os.MkdirAll(config.AppConfig.IncomingPath, 0755)
	for {
		log.Println("Verificando novos arquivos em 'incoming'...")
		err := s.scanIncoming()
		if err != nil {
			log.Printf("Erro ao escanear entrada: %v\n", err)
		}
		time.Sleep(1 * time.Minute)
	}
}
func (s *LibraryService) scanIncoming() error {
	files, err := scanner.ScanDirectory(config.AppConfig.IncomingPath)
	if err != nil {
		return err
	}
	existingPaths, err := s.mediaFileRepository.FindAllPaths()
	if err != nil {
		return err
	}
	for _, file := range files {
		if existingPaths[file.Path] {
			continue
		}
		log.Printf("Novo arquivo detectado em incoming: %s\n", file.Path)
		if utils.IsSeriesPath(file.Path) {
			s.processSeries(file)
		} else {
			s.processMovie(file)
		}
	}
	return nil
}
func (s *LibraryService) processMovie(file scanner.ScannedFile) error {
	parsed := utils.ParseMovieFilename(file.Path)
	duration, _ := utils.GetVideoDuration(file.Path)
	mediaItem := &models.MediaItem{
		Type:          models.MediaTypeMovie,
		Title:         parsed.Title,
		OriginalTitle: parsed.Title,
		Year:          parsed.Year,
		Duration:      duration,
	}
	metadata, _ := s.metadataService.SearchMovie(parsed.Title, parsed.Year)
	if metadata != nil {
		mediaItem.Overview = metadata.Overview
		mediaItem.Poster = metadata.Poster
		mediaItem.Backdrop = metadata.Backdrop
		mediaItem.TMDBID = metadata.TMDBID
	}
	err := s.mediaItemRepository.Create(mediaItem)
	if err != nil {
		return err
	}
	mediaFile := &models.MediaFile{
		MediaItemID: mediaItem.ID,
		Path:        file.Path,
		Size:        file.Size,
		Fingerprint: file.Fingerprint,
		Quality:     parsed.Quality,
		Duration:    duration,
		Status:      models.FileStatusPending,
	}
	return s.mediaFileRepository.Create(mediaFile)
}
func (s *LibraryService) processSeries(file scanner.ScannedFile) error {
	parsed := utils.ParseSeriesFilename(file.Path)
	if !parsed.IsSeries { return nil }
	duration, _ := utils.GetVideoDuration(file.Path)
	metadata, _ := s.metadataService.SearchSeries(parsed.Title, parsed.Year)
	_, seriesModel, err := s.findOrCreateSeries(parsed.Title, parsed.Year, metadata)
	if err != nil { return err }
	season, err := s.findOrCreateSeason(seriesModel.ID, parsed.Season)
	if err != nil { return err }
	episodeMediaItem := &models.MediaItem{
		Type:          models.MediaTypeEpisode,
		Title:         parsed.Title,
		OriginalTitle: parsed.Title,
		Year:          parsed.Year,
		Duration:      duration,
	}
	s.mediaItemRepository.Create(episodeMediaItem)
	s.findOrCreateEpisode(season.ID, parsed.Episode, episodeMediaItem.ID)
	mediaFile := &models.MediaFile{
		MediaItemID: episodeMediaItem.ID,
		Path:        file.Path,
		Size:        file.Size,
		Fingerprint: file.Fingerprint,
		Quality:     parsed.Quality,
		Duration:    duration,
		Status:      models.FileStatusPending,
	}
	return s.mediaFileRepository.Create(mediaFile)
}
func (s *LibraryService) findOrCreateSeries(title string, year int, metadata *SeriesMetadata) (*models.MediaItem, *models.Series, error) {
	var mediaItem *models.MediaItem
	if metadata != nil {
		existing, err := s.mediaItemRepository.FindByTMDBID(metadata.TMDBID, models.MediaTypeSeries)
		if err == nil && existing != nil {
			mediaItem = existing
		}
	}
	if mediaItem == nil {
		mediaItem = &models.MediaItem{
			Type:          models.MediaTypeSeries,
			Title:         title,
			OriginalTitle: title,
			Year:          year,
		}
		if metadata != nil {
			mediaItem.Title = metadata.Title
			mediaItem.OriginalTitle = metadata.OriginalTitle
			mediaItem.Overview = metadata.Overview
			mediaItem.Poster = metadata.Poster
			mediaItem.Backdrop = metadata.Backdrop
			mediaItem.TMDBID = metadata.TMDBID
			mediaItem.Year = metadata.Year
		}
		err := s.mediaItemRepository.Create(mediaItem)
		if err != nil {
			return nil, nil, err
		}
	}
	seriesModel, err := s.seriesRepository.FindByMediaItemID(mediaItem.ID)
	if err != nil {
		seriesModel = &models.Series{
			MediaItemID: mediaItem.ID,
		}
		err = s.seriesRepository.Create(seriesModel)
		if err != nil {
			return nil, nil, err
		}
	}
	return mediaItem, seriesModel, nil
}
func (s *LibraryService) findOrCreateSeason(seriesID uint, number int) (*models.Season, error) {
	season, err := s.seasonRepository.FindBySeriesAndNumber(seriesID, number)
	if err == nil {
		return season, nil
	}
	season = &models.Season{
		SeriesID: seriesID,
		Number:   number,
	}
	err = s.seasonRepository.Create(season)
	return season, err
}
func (s *LibraryService) findOrCreateEpisode(seasonID uint, number int, mediaItemID uint) (*models.Episode, error) {
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
	return episode, err
}
