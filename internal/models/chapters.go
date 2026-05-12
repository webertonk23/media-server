package models
import (
	"time"
)
type Chapter struct {
	ID          uint      `gorm:"primaryKey"`
	MediaItemID uint      `gorm:"index"`
	Title       string    `gorm:"size:255"`
	StartTime   float64   `gorm:"not null"`
	EndTime     float64   `gorm:"not null"`
	CreatedAt   time.Time
}
