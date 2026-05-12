package services

import (
	"media-server/internal/metadata"
)

type MetadataService struct{}

func NewMetadataService() *MetadataService {
	return &MetadataService{}
}

type MovieMetadata struct {
	Title         string
	OriginalTitle string
	Year          int
	Overview      string
	Poster        string
	Backdrop      string
	TMDBID        int
}

type SeriesMetadata struct {
	Title         string
	OriginalTitle string
	Year          int
	Overview      string
	Poster        string
	Backdrop      string
	TMDBID        int
}

func (s *MetadataService) SearchMovie(title string, year int) (*MovieMetadata, error) {
	result, err := metadata.SearchMovie(title, year)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	meta := &MovieMetadata{
		Title:         result.Title,
		OriginalTitle: result.OriginalTitle,
		Year:          result.Year,
		Overview:      result.Overview,
		TMDBID:        result.ID,
	}
	if result.PosterPath != "" {
		meta.Poster = "https://image.tmdb.org/t/p/w500" + result.PosterPath
	}
	if result.Backdrop != "" {
		meta.Backdrop = "https://image.tmdb.org/t/p/original" + result.Backdrop
	}
	return meta, nil
}

func (s *MetadataService) SearchSeries(title string, year int) (*SeriesMetadata, error) {
	result, err := metadata.SearchSeries(title, year)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	meta := &SeriesMetadata{
		Title:         result.Name,
		OriginalTitle: result.OriginalName,
		Year:          result.Year,
		Overview:      result.Overview,
		TMDBID:        result.ID,
	}
	if result.PosterPath != "" {
		meta.Poster = "https://image.tmdb.org/t/p/w500" + result.PosterPath
	}
	if result.Backdrop != "" {
		meta.Backdrop = "https://image.tmdb.org/t/p/original" + result.Backdrop
	}
	return meta, nil
}
