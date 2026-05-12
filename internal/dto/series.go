package dto
type SeriesResponse struct {
	ID       string `json:"id"` 
	Type     string `json:"type"`
	Title    string `json:"title"`
	Year     int    `json:"year,omitempty"`
	Overview string `json:"overview,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Backdrop string `json:"backdrop,omitempty"`
	Status         string `json:"status,omitempty"`
	NumberSeasons  int    `json:"number_seasons,omitempty"`
	NumberEpisodes int    `json:"number_episodes,omitempty"`
}
type SeasonResponse struct {
	ID     string `json:"id"` 
	Number int    `json:"number"`
	Name   string `json:"name,omitempty"`
	Overview     string `json:"overview,omitempty"`
	Poster       string `json:"poster,omitempty"`
	EpisodeCount int    `json:"episode_count,omitempty"`
}
type EpisodeResponse struct {
	ID       string `json:"id"` 
	Type     string `json:"type"`
	Title    string `json:"title"`
	Overview string `json:"overview,omitempty"`
	SeasonNumber  int    `json:"season_number"`
	EpisodeNumber int    `json:"episode_number"`
	Still         string `json:"still,omitempty"`
	Runtime       int    `json:"runtime,omitempty"`
	StreamURL string `json:"stream_url,omitempty"`
	Quality string `json:"quality,omitempty"`
}
