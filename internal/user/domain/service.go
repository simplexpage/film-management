package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	customError "film-management/pkg/errors"
	"github.com/pkg/errors"
	"time"
)

// service is a struct for domain service.
type service struct {
	Opts
}

// NewService is a constructor for domain service.
func NewService(userRepository UserRepository, authService AuthService, passwordService PasswordService, opts ...OptFunc) Service {
	// Init default options
	o := defaultOpts(userRepository, authService, passwordService)

	// Apply options
	for _, opt := range opts {
		opt(&o)
	}

	return &service{
		Opts: o,
	}
}

// Register is a method to register new user.
func (s service) Register(ctx context.Context, model *models.User) error {
	// Check if a user with the same username already exists in db
	if err := s.checkDuplicateUser(ctx, model.Username); err != nil {
		return err
	}

	// Generated hashed password
	hashPassword, errPassword := s.passwordService.GeneratePasswordHash(model.Password)
	if errPassword != nil {
		return ErrGeneratePasswordHash
	}

	// Set hashed password
	model.Password = hashPassword

	// Create user in db
	if err := s.userRepository.CreateUser(ctx, model); err != nil {
		return ErrUserCreate
	}

	return nil
}

// checkDuplicateUser Check if a user with the same username already exists in db.
func (s service) checkDuplicateUser(ctx context.Context, username string) error {
	if err := s.userRepository.UserExistsWithUsername(ctx, username); err != nil {
		switch {
		case errors.Is(err, ErrUserExistsWithUsername):
			return customError.ValidationError{Field: "username", Err: ErrUserExistsWithUsername}
		default:
			return ErrUserCheckExistence
		}
	}

	return nil
}

// Login is a method to login user.
func (s service) Login(ctx context.Context, username string, password string) (string, time.Time, error) {
	// Find user by username in db
	user, err := s.userRepository.FindOneUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return "", time.Time{}, customError.ValidationError{Field: "username", Err: ErrIncorrectLoginOrPassword}
		default:
			return "", time.Time{}, ErrUserFindByUsername
		}
	}

	// Compare password hash
	if err := s.passwordService.ComparePasswordHash(password, user.Password); err != nil {
		return "", time.Time{}, customError.ValidationError{Field: "username", Err: ErrIncorrectLoginOrPassword}
	}

	// Generate auth token
	if authToken, expirationTime, errAuthToken := s.authService.GenerateAuthToken(user.UUID.String()); errAuthToken != nil {
		return "", time.Time{}, ErrGenerateAuthToken
	} else {
		return authToken, expirationTime, nil
	}
}
