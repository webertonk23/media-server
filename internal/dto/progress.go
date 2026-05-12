package dto
type UpdateProgressRequest struct {
	Position int64 `json:"position"`
	Duration int64 `json:"duration"`
	Finished bool  `json:"finished"`
}
type ContinueWatchingResponse struct {
	Media    MediaItemResponse `json:"media"`
	Position int64             `json:"position"`
	Duration int64             `json:"duration"`
	Finished bool              `json:"finished"`
}
