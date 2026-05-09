package mappers

import (
	"fmt"
	"media-server/internal/dto"
	"media-server/internal/models"
)

func ToMediaItemResponse(item models.MediaItem) dto.MediaItemResponse {
	quality := ""
	var files []dto.MediaFileResponse
	
	if len(item.Files) > 0 {
		quality = item.Files[0].Quality
		for _, f := range item.Files {
			files = append(files, ToMediaFileResponse(f))
		}
	}

	return dto.MediaItemResponse{
		ID:        item.ULID,
		Type:      item.Type,
		Title:     item.Title,
		Year:      item.Year,
		Overview:  item.Overview,
		Poster:    item.Poster,
		Backdrop:  item.Backdrop,
		StreamURL: fmt.Sprintf("/stream/%s", item.ULID),
		Quality:   quality,
		Duration:  item.Duration,
		Files:     files,
	}
}
