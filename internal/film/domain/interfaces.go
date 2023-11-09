package domain

import (
	"context"
	"film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"github.com/google/uuid"
)

// Service is an interface for domain service.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_service.go -package=mocks
type Service interface {
	AddFilm(ctx context.Context, model *models.Film) error
	UpdateFilm(ctx context.Context, model *models.Film) error
	ViewFilm(ctx context.Context, filmID uuid.UUID) (models.Film, error)
	ViewAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]models.Film, pagination.Pagination, error)
	DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) error
}

// Repository is a repository for domain service
type Repository interface {
	FilmRepository
	GenreRepository
	CastRepository
}

// FilmRepository is a repository for film.
type FilmRepository interface {
	CreateFilm(ctx context.Context, model *models.Film) error
	UpdateFilm(ctx context.Context, model *models.Film) error
	FindOneFilmByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error)
	FindOneFilmForViewByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error)
	FindAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]models.Film, pagination.Pagination, error)
	DeleteFilm(ctx context.Context, uuid uuid.UUID) error
	FilmExistsWithTitle(ctx context.Context, title string, filmID uuid.UUID, operation models.Operation) error
}

// GenreRepository is a repository for genre.
type GenreRepository interface {
	CreateGenre(ctx context.Context, model *models.Genre) (*models.Genre, error)
	GetGenresByNames(ctx context.Context, names []string) ([]models.Genre, error)
}

// CastRepository is a repository for cast.
type CastRepository interface {
	CreateCast(ctx context.Context, model *models.Cast) (*models.Cast, error)
	GetCastsByNames(ctx context.Context, names []string) ([]models.Cast, error)
}
