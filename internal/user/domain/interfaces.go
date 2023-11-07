package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
	"time"
)

// Service is an interface for domain service.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	Register(ctx context.Context, model *models.User) error
	Login(ctx context.Context, username string, password string) (string, time.Time, error)
}

// UserRepository is a repository for user.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error)
	FindOneUserByUsername(ctx context.Context, username string) (models.User, error)
	UserExistsWithUsername(ctx context.Context, username string) error
}

type PasswordService interface {
	GeneratePasswordHash(password string) (string, error)
	ComparePasswordHash(password, hash string) bool
}

type AuthService interface {
	GenerateAuthToken(userID string) (string, time.Time, error)
}
