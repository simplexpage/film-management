package models

import (
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Film struct {
	UUID        uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey"`
	CreatorID   uuid.UUID `json:"creatorID" gorm:"type:uuid;not null"`
	Title       string    `json:"title" gorm:"unique;not null;size:100"`
	Director    string    `json:"director" gorm:"not null"`
	ReleaseDate time.Time `json:"release_date" gorm:"type:date;not null;index"`
	Cast        string    `json:"cast" gorm:"not null"`
	Genres      []Genre   `gorm:"many2many:film_genres;"`
	Synopsis    string    `json:"synopsis" gorm:"type:text;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime"`

	Creator models.User `json:"creator" gorm:"foreignKey:CreatorID;references:UUID;constraint:OnDelete:CASCADE"`
}

func (f *Film) BeforeCreate(_ *gorm.DB) error {
	f.UUID = uuid.New()

	return nil
}

// SetDataForUpdate sets data for update.
func (f *Film) SetDataForUpdate(data *Film) {
	f.Title = data.Title
	f.Director = data.Director
	f.ReleaseDate = data.ReleaseDate
	f.Cast = data.Cast
	f.Synopsis = data.Synopsis
}
