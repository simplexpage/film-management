package migrate

import (
	"errors"
	modelsFilm "film-management/internal/film/domain/models"
	"film-management/internal/user/domain/models"
	"film-management/pkg/database/postgresql"
	"go.uber.org/zap"
	"time"
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

	// Migrate database
	if err := clientDB.AutoMigrate(
		&models.User{},
		&modelsFilm.Film{},
		&modelsFilm.Genre{},
		&modelsFilm.Director{},
		&modelsFilm.Cast{}); err != nil {
		logger.Error("Error migrate p2p database", zap.Error(err))

		return ErrMigrateFilmDatabase
	}

	logger.Info("Migrate p2p database success")

	return nil
}

// SeedTestData seeds test data.
func SeedTestData(sc *postgresql.Config, logger *zap.Logger) error {
	logger.Info("Run cron migrate database")
	clientDB, errDB := postgresql.Connect(sc, logger)

	if errDB != nil {
		logger.Error("Error connect to p2p database", zap.Error(errDB))

		return ErrConnectFilmDB
	}

	// Add users
	users := []models.User{
		{Username: "user1", Password: "$2a$14$ETk8B0Jndb3mrJauT3Ns1OPSAgR.RnfKqTQhLzGLoaTiFIODum7XC"},
		{Username: "user2", Password: "$2a$14$ETk8B0Jndb3mrJauT3Ns1OPSAgR.RnfKqTQhLzGLoaTiFIODum7XC"},
	}
	clientDB.Create(&users)

	// Add genres
	genres := []modelsFilm.Genre{
		{Name: "action"},
		{Name: "comedy"},
		{Name: "drama"},
		{Name: "thriller"},
		{Name: "horror"},
		{Name: "romance"},
		{Name: "sci-fi"},
		{Name: "fantasy"},
		{Name: "adventure"},
		{Name: "animation"},
		{Name: "crime"},
		{Name: "documentary"},
		{Name: "family"},
		{Name: "history"},
		{Name: "mystery"},
		{Name: "war"},
	}
	clientDB.Create(&genres)

	directors := []modelsFilm.Director{
		{Name: "Director A"},
		{Name: "Director B"},
		{Name: "Director C"},
		{Name: "Director D"},
		{Name: "Director F"},
		{Name: "Director G"},
	}
	clientDB.Create(&directors)

	casts := []modelsFilm.Cast{
		{Name: "Cast A"},
		{Name: "Cast B"},
		{Name: "Cast C"},
		{Name: "Cast D"},
		{Name: "Cast F"},
		{Name: "Cast G"},
	}
	clientDB.Create(&casts)

	// Add films
	films := []modelsFilm.Film{
		{CreatorID: users[0].UUID, Title: "Film 1", Director: directors[0], Genres: []modelsFilm.Genre{genres[0], genres[1]}, ReleaseDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[0], casts[1]}, Synopsis: "Synopsis A"},
		{CreatorID: users[1].UUID, Title: "Film 2", Director: directors[1], Genres: []modelsFilm.Genre{genres[1], genres[2]}, ReleaseDate: time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[1], casts[2]}, Synopsis: "Synopsis B"},
		{CreatorID: users[0].UUID, Title: "Film 3", Director: directors[2], Genres: []modelsFilm.Genre{genres[3], genres[4]}, ReleaseDate: time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[2], casts[3]}, Synopsis: "Synopsis C"},
		{CreatorID: users[1].UUID, Title: "Film 4", Director: directors[3], Genres: []modelsFilm.Genre{genres[5], genres[6]}, ReleaseDate: time.Date(2017, 4, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[3], casts[4]}, Synopsis: "Synopsis D"},
		{CreatorID: users[0].UUID, Title: "Film 5", Director: directors[4], Genres: []modelsFilm.Genre{genres[7], genres[8]}, ReleaseDate: time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[5], casts[1]}, Synopsis: "Synopsis E"},
		{CreatorID: users[1].UUID, Title: "Film 6", Director: directors[5], Genres: []modelsFilm.Genre{genres[9], genres[10]}, ReleaseDate: time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC), Casts: []modelsFilm.Cast{casts[2], casts[4]}, Synopsis: "Synopsis F"},
	}
	clientDB.Create(&films)

	logger.Info("Test data seeded successfully")
	return nil
}
