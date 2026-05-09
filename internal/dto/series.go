package dto

// SeriesResponse representa a resposta da API para uma série
type SeriesResponse struct {
	ID       string `json:"id"` // ULID
	Type     string `json:"type"`
	Title    string `json:"title"`
	Year     int    `json:"year,omitempty"`
	Overview string `json:"overview,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Backdrop string `json:"backdrop,omitempty"`

	// Informações específicas de séries
	Status         string `json:"status,omitempty"`
	NumberSeasons  int    `json:"number_seasons,omitempty"`
	NumberEpisodes int    `json:"number_episodes,omitempty"`
}

// SeasonResponse representa a resposta da API para uma temporada
type SeasonResponse struct {
	ID     string `json:"id"` // ULID
	Number int    `json:"number"`
	Name   string `json:"name,omitempty"`

	Overview     string `json:"overview,omitempty"`
	Poster       string `json:"poster,omitempty"`
	EpisodeCount int    `json:"episode_count,omitempty"`
}

// EpisodeResponse representa a resposta da API para um episódio
type EpisodeResponse struct {
	ID       string `json:"id"` // ULID do MediaItem
	Type     string `json:"type"`
	Title    string `json:"title"`
	Overview string `json:"overview,omitempty"`

	// Informações do episódio
	SeasonNumber  int    `json:"season_number"`
	EpisodeNumber int    `json:"episode_number"`
	Still         string `json:"still,omitempty"`
	Runtime       int    `json:"runtime,omitempty"`

	// URL para streaming
	StreamURL string `json:"stream_url,omitempty"`

	Quality string `json:"quality,omitempty"`
}
