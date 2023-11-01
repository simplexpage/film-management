package migrate

import (
	"errors"
	"film-management/pkg/database/postgresql"
	"go.uber.org/zap"
)

var (
	ErrMigrateFilmDatabase = errors.New("error migrate film database")
	ErrConnectFilmDB       = errors.New("error connect to film database")
)

// PostgresDatabase migrate database.
func PostgresDatabase(sc *postgresql.Config, logger *zap.Logger) error {
	logger.Info("Run cron migrate database")
	clientDB, errDB := postgresql.Connect(sc, logger)

	if errDB != nil {
		logger.Error("Error connect to p2p database", zap.Error(errDB))

		return ErrConnectFilmDB
	}

	err := clientDB.AutoMigrate()

	if err != nil {
		logger.Error("Error migrate p2p database", zap.Error(err))

		return ErrMigrateFilmDatabase
	}

	logger.Info("Migrate p2p database success")

	return nil
}
