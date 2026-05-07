package models

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Series struct {
	ID   uint   `gorm:"primaryKey"`
	ULID string `gorm:"column:ulid;uniqueIndex;size:26;not null"`

	MediaItemID uint      `gorm:"uniqueIndex;not null"`
	MediaItem   MediaItem `gorm:"foreignKey:MediaItemID"`

	// Informações específicas de séries
	Status       string `gorm:"size:50"`  // Returning, Ended, Canceled
	NumberSeasons int    `gorm:"default:0"`
	NumberEpisodes int   `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook para gerar ULID automaticamente
func (s *Series) BeforeCreate(tx *gorm.DB) error {
	if s.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		s.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
