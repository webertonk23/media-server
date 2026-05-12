package models
import "time"
type MediaProgress struct {
	ID uint `gorm:"primaryKey"`
	MediaItemID uint      `gorm:"uniqueIndex;not null"`
	MediaItem   MediaItem `gorm:"foreignKey:MediaItemID"`
	Position int64 `gorm:"not null;default:0"`
	Duration int64 `gorm:"not null;default:0"`
	Finished bool `gorm:"not null;default:false"`
	LastWatchedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
