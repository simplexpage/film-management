package domain

import "errors"

var (
	ErrUserCreate               = errors.New("failed to create user")
	ErrUserExistsWithUsername   = errors.New("user already exists with the same username")
	ErrUserExists               = errors.New("failed check user exists")
	ErrIncorrectLoginOrPassword = errors.New("incorrect email or password")
)
