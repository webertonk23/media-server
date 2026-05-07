package dto

// MediaItemResponse representa a resposta da API para um item de mídia
type MediaItemResponse struct {
	ID       string `json:"id"` // ULID público
	Type     string `json:"type"`
	Title    string `json:"title"`
	Year     int    `json:"year,omitempty"`
	Overview string `json:"overview,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Backdrop string `json:"backdrop,omitempty"`

	// URL para streaming
	StreamURL string `json:"stream_url,omitempty"`
}
