package postgresql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
	"time"
)

// Connect to postgres database
func Connect(config *Config, logger *zap.Logger) (*gorm.DB, error) {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	gormLogger := NewGormLogger(logger)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logger.Debug("Failed to connect to database", zap.String("url", dbURL))

		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Connected to postgres database")
	logger.Debug("PostgresUrl", zap.String("url", dbURL))

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
