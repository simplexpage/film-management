package services

import (
	"film-management/internal/user/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService is a struct for Password.
type PasswordService struct {
	logger *zap.Logger
}

// NewPasswordService is a constructor for PasswordService.
func NewPasswordService(logger *zap.Logger) *PasswordService {
	return &PasswordService{
		logger: logger,
	}
}

func (s PasswordService) GeneratePasswordHash(password string) (string, error) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14); err != nil {
		s.logger.Error("failed to generate password hash", zap.Error(err))

		return "", domain.ErrGeneratePasswordHash
	} else {
		return string(bytes), nil
	}
}

func (s PasswordService) ComparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
