package mappers

import (
	"fmt"
	"media-server/internal/dto"
	"media-server/internal/models"
)

func ToMediaItemResponse(item models.MediaItem) dto.MediaItemResponse {
	return dto.MediaItemResponse{
		ID:       item.ULID, // Usar ULID público
		Type:     item.Type,
		Title:    item.Title,
		Year:     item.Year,
		Overview: item.Overview,
		Poster:   item.Poster,
		Backdrop: item.Backdrop,
		// StreamURL usa ULID
		StreamURL: fmt.Sprintf("/stream/%s", item.ULID),
	}
}
