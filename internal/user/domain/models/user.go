package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is a model for user.
type User struct {
	UUID      uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey"`
	Username  string    `json:"username" gorm:"size:40;unique;not null"`
	Password  string    `json:"password" gorm:"size:255;not null"`
	CreatedAt int64     `gorm:"autoCreateTime"`
	UpdatedAt int64     `gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.UUID = uuid.New()

	return nil
}
