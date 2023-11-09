package password

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service is a struct for Password.
type Service struct {
	logger *zap.Logger
}

// NewPasswordService is a constructor for Service.
func NewPasswordService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s Service) GeneratePasswordHash(password string) (string, error) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14); err != nil {
		s.logger.Error("failed to generate password hash", zap.Error(err))

		return "", errors.Wrap(err, "passwordService.GeneratePasswordHash")
	} else {
		return string(bytes), nil
	}
}

func (s Service) ComparePasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		s.logger.Debug("failed to compare password hash", zap.Error(err))

		return errors.Wrap(err, "passwordService.ComparePasswordHash")
	}

	return nil
}
