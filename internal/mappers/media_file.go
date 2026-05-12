package mappers
import (
	"media-server/internal/dto"
	"media-server/internal/models"
)
func ToMediaFileResponse(file models.MediaFile) dto.MediaFileResponse {
	return dto.MediaFileResponse{
		ID:           file.ID,
		ULID:         file.ULID,
		MediaItemID:  file.MediaItemID,
		Path:         file.Path,
		Size:         file.Size,
		Fingerprint:  file.Fingerprint,
		Quality:      file.Quality,
		Status:       file.Status,
		OriginalPath: file.OriginalPath,
		ErrorMessage: file.ErrorMessage,
		RetryCount:   file.RetryCount,
		Duration:     file.Duration,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}
func ToMediaFileResponses(files []models.MediaFile) []dto.MediaFileResponse {
	responses := make([]dto.MediaFileResponse, len(files))
	for i, f := range files {
		responses[i] = ToMediaFileResponse(f)
	}
	return responses
}
