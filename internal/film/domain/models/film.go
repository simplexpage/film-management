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
	DirectorID  uint      `json:"directorID" gorm:"not null"`
	ReleaseDate time.Time `json:"release_date" gorm:"type:date;not null;index"`
	Casts       []Cast    `json:"casts" gorm:"many2many:film_casts;constraint:OnDelete:CASCADE"`
	Genres      []Genre   `json:"genres" gorm:"many2many:film_genres;constraint:OnDelete:CASCADE"`
	Synopsis    string    `json:"synopsis" gorm:"type:text;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime"`

	Creator  models.User `json:"creator" gorm:"foreignKey:CreatorID;references:UUID;constraint:OnDelete:CASCADE"`
	Director Director    `json:"director" gorm:"foreignKey:DirectorID;references:ID;constraint:OnDelete:CASCADE"`
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
	f.Casts = data.Casts
	f.Synopsis = data.Synopsis
	f.Genres = data.Genres
}

// Operation is a type for film operation.
type Operation string

const (
	// OperationAdd is an operation for adding film.
	OperationAdd Operation = "add"
	// OperationUpdate is an operation for updating film.
	OperationUpdate Operation = "update"
)
