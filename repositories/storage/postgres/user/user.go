package user

import (
	"context"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository is a struct for User.
type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserRepository is a constructor for Repository.
func NewUserRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CreateUser is a method to create user.
func (r Repository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("userRepo.CreateUser.Create", zap.Error(err))

		return errors.Wrap(err, "userRepo.CreateUser.Create")
	}

	return nil
}

// FindOneUserByUUID is a method to find one user.
func (r Repository) FindOneUserByUUID(ctx context.Context, uuid uuid.UUID) (models.User, error) {
	var user models.User

	if result := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.Wrap(domain.ErrUserNotFound, "userRepo.FindOneUserByUUID.First")
		}
		r.logger.Error("userRepo.FindOneUserByUUID.First", zap.Error(result.Error))

		return models.User{}, errors.Wrap(result.Error, "userRepo.FindOneUserByUUID.First")
	}

	return user, nil
}

// FindOneUserByUsername is a method to find one user.
func (r Repository) FindOneUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	if result := r.db.WithContext(ctx).Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.Wrap(domain.ErrUserNotFound, "userRepo.FindOneUserByUsername.First")
		}
		r.logger.Error("userRepo.FindOneUserByUsername.First", zap.Error(result.Error))

		return models.User{}, errors.Wrap(result.Error, "userRepo.FindOneUserByUsername.First")
	}

	return user, nil
}

// UserExistsWithUsername checks if a user with the given username exists.
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
		r.logger.Error("userRepo.UserExistsWithUsername.Count", zap.Error(err))

		return errors.Wrap(err, "userRepo.UserExistsWithUsername.Count")
	}

	// If count > 0, then a user with the same username exists
	if count > 0 {
		return errors.Wrap(domain.ErrUserExistsWithUsername, "userRepo.UserExistsWithUsername.Count")
	}

	return nil
}
