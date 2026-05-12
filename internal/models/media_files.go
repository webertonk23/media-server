package models
import (
	"crypto/rand"
	"time"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)
const (
	FileStatusPending    = "pending"
	FileStatusProcessing = "transcoding"
	FileStatusCompleted  = "completed"
	FileStatusError      = "error"
)
type MediaFile struct {
	ID           uint      `gorm:"primaryKey"`
	ULID         string    `gorm:"column:ulid;uniqueIndex;size:26;not null"`
	MediaItemID  uint      `gorm:"index;not null"`
	MediaItem    MediaItem `gorm:"foreignKey:MediaItemID"`
	Path         string    `gorm:"uniqueIndex;size:1000;not null"`
	Size         int64     `gorm:"not null"`
	Fingerprint  string    `gorm:"index;size:64;not null"`
	Quality      string    `gorm:"size:20"`
	Duration     float64   `gorm:"default:0"`
	Status       string    `gorm:"index;size:20;default:pending"`
	OriginalPath string    `gorm:"size:1000"`
	ErrorMessage string    `gorm:"type:text"`
	RetryCount   int       `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
func (m *MediaFile) BeforeCreate(tx *gorm.DB) error {
	if m.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		m.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
