package models
import (
	"crypto/rand"
	"time"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)
type MediaItem struct {
	ID            uint        `gorm:"primaryKey"`
	ULID          string      `gorm:"column:ulid;uniqueIndex;size:26;not null"`
	Type          string      `gorm:"index;size:20;not null"`
	Title         string      `gorm:"size:500;not null"`
	OriginalTitle string      `gorm:"size:500"`
	Year          int         `gorm:"index"`
	Overview      string      `gorm:"type:text"`
	Poster        string      `gorm:"size:500"`
	Backdrop      string      `gorm:"size:500"`
	TMDBID        int         `gorm:"index"`
	Duration      float64     `gorm:"default:0"`
	Files         []MediaFile `gorm:"foreignKey:MediaItemID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
func (m *MediaItem) BeforeCreate(tx *gorm.DB) error {
	if m.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		m.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
