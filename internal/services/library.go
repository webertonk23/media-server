package services

import (
	"log"
	"path/filepath"
	"strings"
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

// ScanMovies escaneia o diretório de filmes
func (s *LibraryService) ScanMovies() error {
	moviesPath := filepath.Join(config.AppConfig.MediaPath, "movies")
	return s.scanDirectory(moviesPath, models.MediaTypeMovie)
}

// ScanSeries escaneia o diretório de séries
func (s *LibraryService) ScanSeries() error {
	seriesPath := filepath.Join(config.AppConfig.MediaPath, "series")
	return s.scanDirectory(seriesPath, models.MediaTypeSeries)
}

// ScanAll escaneia todos os diretórios
func (s *LibraryService) ScanAll() error {
	if err := s.ScanMovies(); err != nil {
		log.Printf("Erro ao escanear filmes: %v\n", err)
	}

	if err := s.ScanSeries(); err != nil {
		log.Printf("Erro ao escanear séries: %v\n", err)
	}

	return nil
}

func (s *LibraryService) scanDirectory(path string, mediaType string) error {
	files, err := scanner.ScanDirectory(path)
	if err != nil {
		return err
	}

	existingFingerprints, err := s.mediaFileRepository.FindAllFingerprints()
	if err != nil {
		return err
	}

	for _, file := range files {
		if existingFingerprints[file.Fingerprint] {
			continue
		}
		if utils.IsSeriesPath(file.Path) {
			err := s.processSeries(file)
			if err != nil {
				log.Printf("Erro ao processar série %s: %v\n", file.Path, err)
			}
		} else {
			err := s.processMovie(file)
			if err != nil {
				log.Printf("Erro ao processar filme %s: %v\n", file.Path, err)
			}
		}
	}

	return nil
}

func (s *LibraryService) Scan() error {
	return s.ScanAll()
}
func (s *LibraryService) StartScannerWorker() {
	log.Println("Scanner Worker iniciado. Monitorando pastas...")
	settingsRepo := repositories.NewSettingsRepository()
	for {
		log.Println("Iniciando varredura da biblioteca...")
		err := s.ScanAll()
		if err != nil {
			log.Printf("Erro durante varredura: %v\n", err)
		}

		interval := 1
		settings, err := settingsRepo.Get()
		if err == nil && settings.ScanInterval > 0 {
			interval = settings.ScanInterval
		}

		log.Printf("Varredura concluída. Aguardando %d minuto(s) para a próxima...\n", interval)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func (s *LibraryService) processMovie(file scanner.ScannedFile) error {
	parsed := utils.ParseMovieFilename(file.Path)
	metadata, err := s.metadataService.SearchMovie(parsed.Title, parsed.Year)
	if err != nil {
		log.Printf("Erro ao buscar metadata para %s: %v\n", parsed.Title, err)
	}

	duration, _ := utils.GetVideoDuration(file.Path)

	var mediaItem *models.MediaItem

	if metadata != nil {
		existing, err := s.mediaItemRepository.FindByTMDBID(metadata.TMDBID, models.MediaTypeMovie)
		if err == nil && existing != nil {
			mediaItem = existing
			log.Printf("MediaItem já existe: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)

			if mediaItem.Duration == 0 && duration > 0 {
				mediaItem.Duration = duration
				s.mediaItemRepository.Update(mediaItem)
			}
		}
	}

	if mediaItem == nil {
		mediaItem = &models.MediaItem{
			Type:          models.MediaTypeMovie,
			Title:         parsed.Title,
			OriginalTitle: parsed.Title,
			Year:          parsed.Year,
			Duration:      duration,
		}

		if metadata != nil {
			mediaItem.Overview = metadata.Overview
			mediaItem.Poster = metadata.Poster
			mediaItem.Backdrop = metadata.Backdrop
			mediaItem.TMDBID = metadata.TMDBID
			if metadata.OriginalTitle != "" {
				mediaItem.OriginalTitle = metadata.OriginalTitle
			}
		}

		err = s.mediaItemRepository.Create(mediaItem)
		if err != nil {
			return err
		}

		log.Printf("MediaItem criado: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
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

	err = s.mediaFileRepository.Create(mediaFile)
	if err != nil {
		return err
	}

	log.Printf("MediaFile criado: %s\n", file.Path)
	return nil
}

func (s *LibraryService) processSeries(file scanner.ScannedFile) error {
	parsed := utils.ParseSeriesFilename(file.Path)

	if !parsed.IsSeries {
		return nil
	}

	log.Printf("Processando série: %s S%02dE%02d\n", parsed.Title, parsed.Season, parsed.Episode)

	metadata, err := s.metadataService.SearchSeries(parsed.Title, parsed.Year)
	if err != nil {
		log.Printf("Erro ao buscar metadata para série %s: %v\n", parsed.Title, err)
	}

	duration, _ := utils.GetVideoDuration(file.Path)

	_, seriesModel, err := s.findOrCreateSeries(parsed.Title, parsed.Year, metadata)
	if err != nil {
		return err
	}

	season, err := s.findOrCreateSeason(seriesModel.ID, parsed.Season)
	if err != nil {
		return err
	}

	log.Printf("Season encontrada/criada: S%02d (ID: %d)\n", season.Number, season.ID)

	episodeMediaItem := &models.MediaItem{
		Type:          models.MediaTypeEpisode,
		Title:         parsed.Title,
		OriginalTitle: parsed.Title,
		Year:          parsed.Year,
		Duration:      duration,
	}

	err = s.mediaItemRepository.Create(episodeMediaItem)
	if err != nil {
		return err
	}

	log.Printf("MediaItem do episódio criado (ID: %d)\n", episodeMediaItem.ID)

	episode, err := s.findOrCreateEpisode(season.ID, parsed.Episode, episodeMediaItem.ID)
	if err != nil {
		return err
	}

	log.Printf("Episode criado: S%02dE%02d (ID: %d)\n", season.Number, episode.Number, episode.ID)

	mediaFile := &models.MediaFile{
		MediaItemID: episodeMediaItem.ID,
		Path:        file.Path,
		Size:        file.Size,
		Fingerprint: file.Fingerprint,
		Quality:     parsed.Quality,
		Duration:    duration,
		Status:      models.FileStatusPending,
	}

	err = s.mediaFileRepository.Create(mediaFile)
	if err != nil {
		return err
	}

	log.Printf("MediaFile criado: %s\n", file.Path)
	return nil
}

func (s *LibraryService) findOrCreateSeries(title string, year int, metadata *SeriesMetadata) (*models.MediaItem, *models.Series, error) {
	var mediaItem *models.MediaItem

	if metadata != nil {
		existing, err := s.mediaItemRepository.FindByTMDBID(metadata.TMDBID, models.MediaTypeSeries)
		if err == nil && existing != nil {
			mediaItem = existing
			log.Printf("Série encontrada por TMDB ID: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
		}
	}

	if mediaItem == nil {
		items, _, err := s.mediaItemRepository.Paginate(1, 10, title, models.MediaTypeSeries)
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if strings.EqualFold(item.Title, title) {
					mediaItem = &item
					log.Printf("Série encontrada por Título: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
					break
				}
			}
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
		log.Printf("MediaItem da série criado: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
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
		log.Printf("Registro Series criado para MediaItem %d\n", mediaItem.ID)
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
	if err != nil {
		return nil, err
	}

	return season, nil
}

func (s *LibraryService) findOrCreateEpisode(seasonID uint, number int, mediaItemID uint) (*models.Episode, error) {
	episode, err := s.episodeRepository.FindBySeasonAndNumber(seasonID, number)
	if err == nil {
		return episode, nil
	}

	// Criar novo episódio
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
