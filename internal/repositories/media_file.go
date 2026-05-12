package repositories

import (
	"media-server/internal/database"
	"media-server/internal/models"
)

type MediaFileRepository struct{}

func NewMediaFileRepository() *MediaFileRepository {
	return &MediaFileRepository{}
}
func (r *MediaFileRepository) Create(file *models.MediaFile) error {
	return database.DB.Create(file).Error
}
func (r *MediaFileRepository) Update(file *models.MediaFile) error {
	return database.DB.Save(file).Error
}
func (r *MediaFileRepository) FindByID(id uint) (*models.MediaFile, error) {
	var file models.MediaFile
	result := database.DB.First(&file, id)
	return &file, result.Error
}
func (r *MediaFileRepository) FindByULID(ulid string) (*models.MediaFile, error) {
	var file models.MediaFile
	result := database.DB.Where("ulid = ?", ulid).First(&file)
	return &file, result.Error
}
func (r *MediaFileRepository) FindByPath(path string) (*models.MediaFile, error) {
	var file models.MediaFile
	result := database.DB.Where("path = ?", path).First(&file)
	return &file, result.Error
}
func (r *MediaFileRepository) FindByMediaItemID(mediaItemID uint) ([]models.MediaFile, error) {
	var files []models.MediaFile
	result := database.DB.Where("media_item_id = ?", mediaItemID).Find(&files)
	return files, result.Error
}

func (r *MediaFileRepository) DeleteByMediaItemID(mediaItemID uint) ([]models.MediaFile, error) {
	var files []models.MediaFile
	err := database.DB.Where("media_item_id = ?", mediaItemID).Find(&files).Error
	if err != nil {
		return nil, err
	}
	err = database.DB.Where("media_item_id = ?", mediaItemID).Delete(&models.MediaFile{}).Error
	return files, err
}

func (r *MediaFileRepository) FindAllFingerprints() (map[string]bool, error) {
	var files []models.MediaFile
	result := database.DB.
		Select("fingerprint").
		Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	fingerprints := make(map[string]bool)
	for _, file := range files {
		fingerprints[file.Fingerprint] = true
	}
	return fingerprints, nil
}
func (r *MediaFileRepository) FindAllPaths() (map[string]bool, error) {
	var files []models.MediaFile
	result := database.DB.Select("path").Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	paths := make(map[string]bool)
	for _, file := range files {
		paths[file.Path] = true
	}
	return paths, nil
}
func (r *MediaFileRepository) FindByStatus(status string) ([]models.MediaFile, error) {
	var files []models.MediaFile
	result := database.DB.Where("status = ?", status).Find(&files)
	return files, result.Error
}
func (r *MediaFileRepository) FindNextPending() (*models.MediaFile, error) {
	var file models.MediaFile
	result := database.DB.Where("status = ?", models.FileStatusPending).Order("created_at asc").First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}
func (r *MediaFileRepository) UpdateStatus(id uint, status string, errorMessage string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMessage != "" {
		updates["error_message"] = errorMessage
	}
	return database.DB.Model(&models.MediaFile{}).Where("id = ?", id).Updates(updates).Error
}
