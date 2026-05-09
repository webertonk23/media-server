package dto

type MediaItemResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Year      int    `json:"year,omitempty"`
	Overview  string `json:"overview,omitempty"`
	Poster    string `json:"poster,omitempty"`
	Backdrop  string `json:"backdrop,omitempty"`
	StreamURL string `json:"stream_url,omitempty"`
	Quality   string `json:"quality,omitempty"`
}
