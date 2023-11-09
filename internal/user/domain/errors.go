package domain

import "errors"

var (
	ErrUserCreate               = errors.New("failed to create user")
	ErrUserFindByUsername       = errors.New("failed to find user by username")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserExistsWithUsername   = errors.New("user already exists with the same username")
	ErrIncorrectLoginOrPassword = errors.New("incorrect username or password")
	ErrUserCheckExistence       = errors.New("failed to check user existence")
	ErrGeneratePasswordHash     = errors.New("failed to generate password hash")
	ErrGenerateAuthToken        = errors.New("failed to generate auth token")
)
