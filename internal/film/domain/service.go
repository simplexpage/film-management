package domain

import (
	"context"
	modelsFilm "film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"film-management/pkg/validation"
	"github.com/google/uuid"
)

// service is a struct for domain service.
type service struct {
	Opts
}

// NewService is a constructor for domain service.
func NewService(filmRepository FilmRepository, opts ...OptFunc) Service {
	// Init default options
	o := defaultOpts(filmRepository)

	// Apply options
	for _, opt := range opts {
		opt(&o)
	}

	return &service{
		Opts: o,
	}
}

func (s service) AddFilm(ctx context.Context, model *modelsFilm.Film) error {
	// Check if a film with the same title already exists
	if err := s.filmRepository.FilmExistsWithTitle(ctx, model.Title); err != nil {
		switch err {
		case ErrFilmExistsWithTitle:
			return validation.CustomError{Field: "title", Err: ErrFilmExistsWithTitle}
		default:
			return err
		}
	}

	// Create a film in db
	return s.filmRepository.CreateFilm(ctx, model)
}

func (s service) UpdateFilm(ctx context.Context, model *modelsFilm.Film) error {
	// Get film from db
	filmFromDB, err := s.filmRepository.FindOneFilmByUUID(ctx, model.UUID)
	if err != nil {
		return err
	}

	// Check permission
	if s.checkFilmPermission(model.CreatorID, filmFromDB.CreatorID) {
		return ErrFilmNotPermission
	}

	// Set new film data
	filmFromDB.SetDataForUpdate(model)

	// Update a film in db
	return s.filmRepository.UpdateFilm(ctx, &filmFromDB)
}

func (s service) ViewFilm(ctx context.Context, filmID uuid.UUID) (modelsFilm.Film, error) {
	return s.filmRepository.FindOneFilmByUUIDWithCreator(ctx, filmID)
}

func (s service) ViewAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]modelsFilm.Film, pagination.Pagination, error) {
	return s.filmRepository.FindAllFilms(ctx, filterSortPagination)
}

func (s service) DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) error {
	// Get film from db
	filmFromDB, err := s.filmRepository.FindOneFilmByUUID(ctx, filmID)
	if err != nil {
		return err
	}

	// Check permission
	if s.checkFilmPermission(userID, filmFromDB.CreatorID) {
		return ErrFilmNotPermission
	}

	// Delete a film in db
	return s.filmRepository.DeleteFilm(ctx, filmID)
}

// checkFilmPermission Check if a film user is the same as the user who wants to do an action.
func (s service) checkFilmPermission(creatorID uuid.UUID, filmCreatorID uuid.UUID) bool {
	return filmCreatorID != creatorID
}
