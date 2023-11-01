package postgres

import (
	"context"
	"errors"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/models"
	"fmt"
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
			return models.User{}, fmt.Errorf("user not found")
		}

		a.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, fmt.Errorf("failed to find user in db: %w", result.Error)
	}

	return user, nil
}

// FindOneUserByUsername is a method to find one user.
func (a UserRepository) FindOneUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	if result := a.db.WithContext(ctx).Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("user not found")
		}

		a.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, fmt.Errorf("failed to find user in db: %w", result.Error)
	}

	return user, nil
}

// FindOneUserByRefreshToken is a method to find one user.
func (a UserRepository) FindOneUserByRefreshToken(ctx context.Context, token string) (models.User, error) {
	var user models.User

	if result := a.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("user not found")
		}

		a.logger.Error("failed to find user in db", zap.Error(result.Error))

		return models.User{}, fmt.Errorf("failed to find user in db: %w", result.Error)
	}

	return user, nil
}

// UserExists checks if an user with the given username exists.
// The operation parameter specifies the type of operation: "add" or "another".
func (a UserRepository) UserExists(ctx context.Context, username string, operation models.Operation) error {
	var count int64

	switch operation {
	case models.OperationAdd:
		// Check if a user with the same email exists
		err := a.db.
			WithContext(ctx).
			Model(&models.User{}).
			Where("username = ?", username).
			Count(&count).
			Error

		if err != nil {
			a.logger.Error("failed to check user existence by username in db", zap.Error(err))

			return fmt.Errorf("failed to check user existence by username in db: %w", err)
		}

		if count > 0 {
			return domain.ErrUserExistsWithUsername
		}

	default:
		a.logger.Error("unknown operation", zap.String("operation", string(operation)))

		return fmt.Errorf("unknown operation: %s", operation)
	}

	return nil
}

// UpdateRefreshToken is a method for update refresh token.
func (a UserRepository) UpdateRefreshToken(ctx context.Context, uuid uuid.UUID, refreshToken string) error {
	err := a.db.
		WithContext(ctx).
		Table("users").
		Where("uuid = ?", uuid).
		Update("refresh_token", refreshToken).
		Error

	if err != nil {
		return err
	}

	return nil
}
