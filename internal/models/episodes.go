package models

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Episode struct {
	ID   uint   `gorm:"primaryKey"`
	ULID string `gorm:"column:ulid;uniqueIndex;size:26;not null"`

	MediaItemID uint      `gorm:"uniqueIndex;not null"`
	MediaItem   MediaItem `gorm:"foreignKey:MediaItemID"`

	SeasonID uint   `gorm:"not null;uniqueIndex:idx_season_episode"`
	Season   Season `gorm:"foreignKey:SeasonID"`

	Number int `gorm:"not null;uniqueIndex:idx_season_episode"`

	// Informações do episódio
	Name     string `gorm:"size:500"`
	Overview string `gorm:"type:text"`
	Still    string `gorm:"size:500"` // Imagem do episódio

	AirDate *time.Time
	Runtime int `gorm:"default:0"` // Duração em minutos

	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook para gerar ULID automaticamente
func (e *Episode) BeforeCreate(tx *gorm.DB) error {
	if e.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		e.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
