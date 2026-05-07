package mappers

import (
	"fmt"
	"media-server/internal/dto"
	"media-server/internal/models"
)

func ToSeriesResponse(mediaItem models.MediaItem, series models.Series) dto.SeriesResponse {
	return dto.SeriesResponse{
		ID:             mediaItem.ULID,
		Type:           mediaItem.Type,
		Title:          mediaItem.Title,
		Year:           mediaItem.Year,
		Overview:       mediaItem.Overview,
		Poster:         mediaItem.Poster,
		Backdrop:       mediaItem.Backdrop,
		Status:         series.Status,
		NumberSeasons:  series.NumberSeasons,
		NumberEpisodes: series.NumberEpisodes,
	}
}

func ToSeasonResponse(season models.Season) dto.SeasonResponse {
	return dto.SeasonResponse{
		ID:           season.ULID,
		Number:       season.Number,
		Name:         season.Name,
		Overview:     season.Overview,
		Poster:       season.Poster,
		EpisodeCount: season.EpisodeCount,
	}
}

func ToEpisodeResponse(episode models.Episode, seasonNumber int) dto.EpisodeResponse {
	return dto.EpisodeResponse{
		ID:            episode.MediaItem.ULID,
		Type:          episode.MediaItem.Type,
		Title:         episode.Name,
		Overview:      episode.Overview,
		SeasonNumber:  seasonNumber,
		EpisodeNumber: episode.Number,
		Still:         episode.Still,
		Runtime:       episode.Runtime,
		StreamURL:     fmt.Sprintf("/stream/%s", episode.MediaItem.ULID),
	}
}
