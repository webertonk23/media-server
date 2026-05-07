package models

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Season struct {
	ID   uint   `gorm:"primaryKey"`
	ULID string `gorm:"column:ulid;uniqueIndex;size:26;not null"`

	SeriesID uint   `gorm:"not null;uniqueIndex:idx_series_season"`
	Series   Series `gorm:"foreignKey:SeriesID"`

	Number int `gorm:"not null;uniqueIndex:idx_series_season"`

	// Informações da temporada
	Name     string `gorm:"size:200"`
	Overview string `gorm:"type:text"`
	Poster   string `gorm:"size:500"`

	AirDate       *time.Time
	EpisodeCount  int `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook para gerar ULID automaticamente
func (s *Season) BeforeCreate(tx *gorm.DB) error {
	if s.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		s.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
