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

// FilmRepository is a repository for film.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_film_repository.go -package=mocks
type FilmRepository interface {
	CreateFilm(ctx context.Context, model *models.Film) error
	UpdateFilm(ctx context.Context, model *models.Film) error
	FindOneFilmByUUID(ctx context.Context, uuid uuid.UUID) (models.Film, error)
	FindOneFilmByUUIDWithCreator(ctx context.Context, uuid uuid.UUID) (models.Film, error)
	FindAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]models.Film, pagination.Pagination, error)
	DeleteFilm(ctx context.Context, uuid uuid.UUID) error
	FilmExistsWithTitle(ctx context.Context, title string) error
}
