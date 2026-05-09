package dto

import "time"

type MediaFileResponse struct {
	ID           uint      `json:"id"`
	ULID         string    `json:"ulid"`
	MediaItemID  uint      `json:"media_item_id"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	Fingerprint  string    `json:"fingerprint"`
	Quality      string    `json:"quality"`
	Status       string    `json:"status"`
	OriginalPath string    `json:"original_path,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	RetryCount   int       `json:"retry_count"`
	Duration     float64   `json:"duration"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
