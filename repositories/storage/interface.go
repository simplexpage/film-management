package storage

import (
	"context"
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
)

// UserRepositorier is a repository for User.
type UserRepositorier interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error)
}
