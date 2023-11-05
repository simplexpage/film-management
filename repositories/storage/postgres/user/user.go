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

// UserRepository is a struct for User.
type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserRepository is a constructor for userRepository.
func NewUserRepository(db *gorm.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser is a method to create user.
func (a UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := a.db.WithContext(ctx).Create(user).Error; err != nil {
		a.logger.Error("failed to add a new user in db", zap.Error(err))

		return domain.ErrUserCreate
	}

	return nil
}

// FindOneUserByUUID is a method to find one user.
func (a UserRepository) FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error) {
	var user models.User

	if result := a.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, domain.ErrUserNotFound
		}

		a.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, domain.ErrUserFind
	}

	return user, nil
}

// FindOneUserByUsername is a method to find one user.
func (a UserRepository) FindOneUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	if result := a.db.WithContext(ctx).Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, domain.ErrUserNotFound
		}

		a.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, domain.ErrUserFind
	}

	return user, nil
}

// UserExistsWithUsername checks if an user with the given username exists.
func (a UserRepository) UserExistsWithUsername(ctx context.Context, username string) error {
	var count int64

	// Check if a user with the same username exists
	err := a.db.
		WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).
		Error

	if err != nil {
		a.logger.Error("failed to check user existence by username in db", zap.Error(err))

		return domain.ErrUserCheckExistence
	}

	if count > 0 {
		return domain.ErrUserExistsWithUsername
	}

	return nil
}
