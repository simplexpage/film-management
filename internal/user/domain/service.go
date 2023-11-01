package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"film-management/pkg/auth"
	"film-management/pkg/validation"
)

// service is a struct for domain service.
type service struct {
	Opts
}

// NewService is a constructor for domain service.
func NewService(userRepository UserRepository, opts ...OptFunc) Service {
	// Init default options
	o := defaultOpts(userRepository)

	// Apply options
	for _, opt := range opts {
		opt(&o)
	}

	return &service{
		Opts: o,
	}
}

// Register is a method to register user.
func (s service) Register(ctx context.Context, model *models.User) (err error) {
	// Check if a user with the same username already exists
	if err := s.userRepository.UserExists(ctx, model.Username, models.OperationAdd); err != nil {
		switch err {
		case ErrUserExistsWithUsername:
			return validation.CustomError{Field: "username", Err: err}
		default:
			return ErrUserExists
		}
	}

	// Create hashed password
	hashPassword, errPassword := model.CreatePassword(model.Password)
	if errPassword != nil {
		return errPassword
	}
	model.Password = hashPassword

	// Create user
	if err := s.userRepository.CreateUser(ctx, model); err != nil {
		return ErrUserCreate
	}

	return nil
}

// Login is a method to login user.
func (s service) Login(ctx context.Context, username string, password string) (authToken string, err error) {
	// Find user by username
	user, err := s.userRepository.FindOneUserByUsername(ctx, username)
	if err != nil {
		return "", validation.CustomError{Field: "password", Err: ErrIncorrectLoginOrPassword}
	}

	// Check password
	if ok := user.CheckPassword(password); !ok {
		return "", validation.CustomError{Field: "password", Err: ErrIncorrectLoginOrPassword}
	}

	// Generate auth token
	authToken, errAuthToken := auth.GenerateAuthToken(
		user.UUID.String(),
		s.cfg.Services.User.PathPrivateKeyFile,
		s.cfg.Services.User.AuthDurationMin)
	if errAuthToken != nil {
		return "", errAuthToken
	}

	return authToken, nil
}
