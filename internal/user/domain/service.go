package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"film-management/pkg/validation"
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
	if err := s.userRepository.UserExistsWithUsername(ctx, model.Username); err != nil {
		switch err {
		case ErrUserExistsWithUsername:
			return validation.CustomError{Field: "username", Err: ErrUserExistsWithUsername}
		default:
			return err
		}
	}

	// Generated hashed password
	hashPassword, errPassword := s.passwordService.GeneratePasswordHash(model.Password)
	if errPassword != nil {
		return ErrGeneratePasswordHash
	}
	// Set hashed password
	model.Password = hashPassword

	// Create user in db
	return s.userRepository.CreateUser(ctx, model)
}

// Login is a method to login user.
func (s service) Login(ctx context.Context, username string, password string) (string, time.Time, error) {
	// Find user by username in db
	user, err := s.userRepository.FindOneUserByUsername(ctx, username)
	if err != nil {
		return "", time.Time{}, validation.CustomError{Field: "password", Err: ErrIncorrectLoginOrPassword}
	}

	// Compare password hash
	if ok := s.passwordService.ComparePasswordHash(password, user.Password); !ok {
		return "", time.Time{}, validation.CustomError{Field: "password", Err: ErrIncorrectLoginOrPassword}
	}

	// Generate auth token
	if authToken, expirationTime, errAuthToken := s.authService.GenerateAuthToken(user.UUID.String()); errAuthToken != nil {
		return "", time.Time{}, errAuthToken
	} else {
		return authToken, expirationTime, nil
	}
}
