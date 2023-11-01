package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
)

// Service is an interface for domain service.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	Register(ctx context.Context, model *models.User) error
	Login(ctx context.Context, username string, password string) (string, error)
}

// UserRepository is a repository for user.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_user_repository.go -package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error)
	FindOneUserByUsername(ctx context.Context, username string) (models.User, error)
	UserExists(ctx context.Context, username string, operation models.Operation) error
}
