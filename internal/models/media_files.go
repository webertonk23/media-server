package models

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MediaFile struct {
	ID   uint   `gorm:"primaryKey"`
	ULID string `gorm:"column:ulid;uniqueIndex;size:26;not null"`

	MediaItemID uint   `gorm:"index;not null"`
	MediaItem   MediaItem `gorm:"foreignKey:MediaItemID"`

	Path string `gorm:"uniqueIndex;size:1000;not null"`

	Size int64 `gorm:"not null"`

	Fingerprint string `gorm:"index;size:64;not null"`

	Quality string `gorm:"size:20"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook para gerar ULID automaticamente
func (m *MediaFile) BeforeCreate(tx *gorm.DB) error {
	if m.ULID == "" {
		entropy := ulid.Monotonic(rand.Reader, 0)
		m.ULID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	}
	return nil
}
