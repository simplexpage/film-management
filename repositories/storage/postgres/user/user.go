package user

import (
	"context"
	"errors"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository is a struct for User.
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewRepository is a constructor for Repository.
func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CreateUser is a method to create user.
func (r Repository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("failed to add a new user in db", zap.Error(err))

		return domain.ErrUserCreate
	}

	return nil
}

// FindOneUserByUUID is a method to find one user.
func (r Repository) FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error) {
	var user models.User

	if result := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, domain.ErrUserNotFound
		}

		r.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, domain.ErrUserFind
	}

	return user, nil
}

// FindOneUserByUsername is a method to find one user.
func (r Repository) FindOneUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	if result := r.db.WithContext(ctx).Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, domain.ErrUserNotFound
		}

		r.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, domain.ErrUserFind
	}

	return user, nil
}

// UserExistsWithUsername checks if an user with the given username exists.
func (r Repository) UserExistsWithUsername(ctx context.Context, username string) error {
	var count int64

	// Check if a user with the same username exists
	err := r.db.
		WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).
		Error

	if err != nil {
		r.logger.Error("failed to check user existence by username in db", zap.Error(err))

		return domain.ErrUserCheckExistence
	}

	if count > 0 {
		return domain.ErrUserExistsWithUsername
	}

	return nil
}
