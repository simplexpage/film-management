package postgresql

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
	"time"
)

// Connect to postgres database.
func Connect(config *Config, logger *zap.Logger) (*gorm.DB, error) {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	logger.Debug("postgres db url", zap.String("url", dbURL))

	gormLogger := NewGormLogger(logger)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))

		return nil, errors.Wrap(err, "failed to connect to database")
	}

	logger.Info("Connected to postgres database")

	return db, nil
}

func NewGormLogger(zapLogger *zap.Logger) zapgorm2.Logger {
	return zapgorm2.Logger{
		ZapLogger:                 zapLogger,
		LogLevel:                  logger.Warn,
		SlowThreshold:             time.Second,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: true,
		Context:                   nil,
	}
}

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}
