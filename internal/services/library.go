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

// scanDirectory é o método genérico que escaneia um diretório
func (s *LibraryService) scanDirectory(path string, mediaType string) error {
	// 1. Scanner: descobre arquivos
	files, err := scanner.ScanDirectory(path)
	if err != nil {
		return err
	}

	// Buscar fingerprints existentes para evitar duplicados
	existingFingerprints, err := s.mediaFileRepository.FindAllFingerprints()
	if err != nil {
		return err
	}

	for _, file := range files {
		// Pular se já existe
		if existingFingerprints[file.Fingerprint] {
			continue
		}

		// 2. Parser: detectar se é filme ou série
		if utils.IsSeriesPath(file.Path) {
			// É uma série
			err := s.processSeries(file)
			if err != nil {
				log.Printf("Erro ao processar série %s: %v\n", file.Path, err)
			}
		} else {
			// É um filme
			err := s.processMovie(file)
			if err != nil {
				log.Printf("Erro ao processar filme %s: %v\n", file.Path, err)
			}
		}
	}

	return nil
}

// Scan é um alias para ScanAll (escaneia filmes e séries)
func (s *LibraryService) Scan() error {
	return s.ScanAll()
}

// StartScannerWorker inicia o loop de scan da biblioteca
func (s *LibraryService) StartScannerWorker() {
	log.Println("Scanner Worker iniciado. Monitorando pastas...")
	for {
		log.Println("Iniciando varredura da biblioteca...")
		err := s.ScanAll()
		if err != nil {
			log.Printf("Erro durante varredura: %v\n", err)
		}
		log.Println("Varredura concluída. Aguardando 1 minuto para a próxima...")
		time.Sleep(1 * time.Minute)
	}
}

// processMovie processa um arquivo de filme
func (s *LibraryService) processMovie(file scanner.ScannedFile) error {
	// Parser: extrai informações do nome do arquivo
	parsed := utils.ParseMovieFilename(file.Path)

	// Metadata provider: busca informações no TMDB
	metadata, err := s.metadataService.SearchMovie(parsed.Title, parsed.Year)
	if err != nil {
		log.Printf("Erro ao buscar metadata para %s: %v\n", parsed.Title, err)
	}

	// Criar ou buscar MediaItem
	var mediaItem *models.MediaItem

	if metadata != nil {
		// Verificar se já existe um MediaItem com esse TMDB ID
		existing, err := s.mediaItemRepository.FindByTMDBID(metadata.TMDBID, models.MediaTypeMovie)
		if err == nil && existing != nil {
			mediaItem = existing
			log.Printf("MediaItem já existe: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
		}
	}

	// Se não existe, criar novo MediaItem
	if mediaItem == nil {
		mediaItem = &models.MediaItem{
			Type:          models.MediaTypeMovie,
			Title:         parsed.Title,
			OriginalTitle: parsed.Title,
			Year:          parsed.Year,
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

	// Criar MediaFile associado ao MediaItem
	mediaFile := &models.MediaFile{
		MediaItemID: mediaItem.ID,
		Path:        file.Path,
		Size:        file.Size,
		Fingerprint: file.Fingerprint,
		Quality:     parsed.Quality,
		Status:      models.FileStatusPending, // Inicia como pendente para o transcoder
	}

	err = s.mediaFileRepository.Create(mediaFile)
	if err != nil {
		return err
	}

	log.Printf("MediaFile criado: %s\n", file.Path)
	return nil
}

// processSeries processa um arquivo de série/episódio
func (s *LibraryService) processSeries(file scanner.ScannedFile) error {
	// Parser: extrai informações de série
	parsed := utils.ParseSeriesFilename(file.Path)

	if !parsed.IsSeries {
		return nil
	}

	log.Printf("Processando série: %s S%02dE%02d\n", parsed.Title, parsed.Season, parsed.Episode)

	// Buscar metadata no TMDB TV para unificar nomes (ex: Foundation vs Fundação)
	metadata, err := s.metadataService.SearchSeries(parsed.Title, parsed.Year)
	if err != nil {
		log.Printf("Erro ao buscar metadata para série %s: %v\n", parsed.Title, err)
	}

	// Buscar ou criar Series (MediaItem)
	_, seriesModel, err := s.findOrCreateSeries(parsed.Title, parsed.Year, metadata)
	if err != nil {
		return err
	}

	// Buscar ou criar Season
	season, err := s.findOrCreateSeason(seriesModel.ID, parsed.Season)
	if err != nil {
		return err
	}

	log.Printf("Season encontrada/criada: S%02d (ID: %d)\n", season.Number, season.ID)

	// Criar MediaItem para o episódio
	episodeMediaItem := &models.MediaItem{
		Type:          models.MediaTypeEpisode,
		Title:         parsed.Title, // Será atualizado com metadata
		OriginalTitle: parsed.Title,
		Year:          parsed.Year,
	}

	err = s.mediaItemRepository.Create(episodeMediaItem)
	if err != nil {
		return err
	}

	log.Printf("MediaItem do episódio criado (ID: %d)\n", episodeMediaItem.ID)

	// Criar Episode
	episode, err := s.findOrCreateEpisode(season.ID, parsed.Episode, episodeMediaItem.ID)
	if err != nil {
		return err
	}

	log.Printf("Episode criado: S%02dE%02d (ID: %d)\n", season.Number, episode.Number, episode.ID)

	// Criar MediaFile associado ao episódio
	mediaFile := &models.MediaFile{
		MediaItemID: episodeMediaItem.ID,
		Path:        file.Path,
		Size:        file.Size,
		Fingerprint: file.Fingerprint,
		Quality:     parsed.Quality,
		Status:      models.FileStatusPending,
	}

	err = s.mediaFileRepository.Create(mediaFile)
	if err != nil {
		return err
	}

	log.Printf("MediaFile criado: %s\n", file.Path)
	return nil
}

// findOrCreateSeries busca ou cria uma série
func (s *LibraryService) findOrCreateSeries(title string, year int, metadata *SeriesMetadata) (*models.MediaItem, *models.Series, error) {
	var mediaItem *models.MediaItem

	// 1. Tentar buscar por TMDB ID se tiver metadata
	if metadata != nil {
		existing, err := s.mediaItemRepository.FindByTMDBID(metadata.TMDBID, models.MediaTypeSeries)
		if err == nil && existing != nil {
			mediaItem = existing
			log.Printf("Série encontrada por TMDB ID: %s (ID: %d)\n", mediaItem.Title, mediaItem.ID)
		}
	}

	// 2. Se não achou por ID, tentar por título exato (fallback)
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

	// 3. Se ainda não existe, criar nova
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

	// Buscar ou criar registro na tabela Series
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

// findOrCreateSeason busca ou cria uma temporada
func (s *LibraryService) findOrCreateSeason(seriesID uint, number int) (*models.Season, error) {
	season, err := s.seasonRepository.FindBySeriesAndNumber(seriesID, number)
	if err == nil {
		return season, nil
	}

	// Criar nova temporada
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

// findOrCreateEpisode busca ou cria um episódio
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
