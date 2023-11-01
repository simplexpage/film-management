package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) CreatePassword(passwordString string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordString), 14)
	return string(bytes), err
}

func (u *User) CheckPassword(passwordString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwordString))
	return err == nil
}

// Operation is a type for user operation.
type Operation string

const (
	// OperationAdd is an operation for adding user.
	OperationAdd Operation = "add"
)
