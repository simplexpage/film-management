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
func NewService(repository Repository, opts ...OptFunc) Service {
	// Init default options
	o := defaultOpts(repository)

	// Apply options
	for _, opt := range opts {
		opt(&o)
	}

	return &service{
		Opts: o,
	}
}

func (s service) AddFilm(ctx context.Context, model *modelsFilm.Film) error {
	// Check duplicate film with title
	if err := s.checkDuplicateFilm(ctx, model.Title); err != nil {
		return err
	}

	// Update existing genres or create new genres
	if err := s.updateFilmGenres(ctx, model); err != nil {
		return err
	}

	// Create a film in db
	return s.repository.CreateFilm(ctx, model)
}

func (s service) UpdateFilm(ctx context.Context, model *modelsFilm.Film) error {
	// Get film from db
	filmFromDB, err := s.repository.FindOneFilmByUUID(ctx, model.UUID)
	if err != nil {
		return err
	}

	// Check duplicate film with title
	if err := s.checkDuplicateFilm(ctx, model.Title); err != nil {
		return err
	}

	// Check permission
	if s.checkFilmPermission(model.CreatorID, filmFromDB.CreatorID) {
		return ErrFilmNotPermission
	}

	// Set new film data
	filmFromDB.SetDataForUpdate(model)

	// Update a film in db
	return s.repository.UpdateFilm(ctx, &filmFromDB)
}

func (s service) ViewFilm(ctx context.Context, filmID uuid.UUID) (modelsFilm.Film, error) {
	return s.repository.FindOneFilmByUUIDWithCreator(ctx, filmID)
}

func (s service) ViewAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]modelsFilm.Film, pagination.Pagination, error) {
	return s.repository.FindAllFilms(ctx, filterSortPagination)
}

func (s service) DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) error {
	// Get film from db
	filmFromDB, err := s.repository.FindOneFilmByUUID(ctx, filmID)
	if err != nil {
		return err
	}

	// Check permission
	if s.checkFilmPermission(userID, filmFromDB.CreatorID) {
		return ErrFilmNotPermission
	}

	// Delete a film in db
	return s.repository.DeleteFilm(ctx, filmID)
}

// checkFilmPermission Check if a film user is the same as the user who wants to do an action.
func (s service) checkFilmPermission(creatorID uuid.UUID, filmCreatorID uuid.UUID) bool {
	return filmCreatorID != creatorID
}

// checkDuplicateFilm Check if a film with the same title already exists.
func (s service) checkDuplicateFilm(ctx context.Context, title string) error {
	if err := s.repository.FilmExistsWithTitle(ctx, title); err != nil {
		switch err {
		case ErrFilmExistsWithTitle:
			return validation.CustomError{Field: "title", Err: ErrFilmExistsWithTitle}
		default:
			return err
		}
	}
	return nil
}

// updateFilmGenres Update existing genres or create new genres.
func (s service) updateFilmGenres(ctx context.Context, model *modelsFilm.Film) error {
	existingGenres, err := s.getExistingGenres(ctx, model.Genres)
	if err != nil {
		return err
	}

	for i, genre := range model.Genres {
		if existingGenre, ok := existingGenres[genre.Name]; ok {
			model.Genres[i] = existingGenre
		} else {
			newGenre, err := s.repository.CreateGenre(ctx, &genre)
			if err != nil {
				return err
			}
			model.Genres[i] = *newGenre
		}
	}

	return nil
}

// getExistingGenres Get existing genres.
func (s service) getExistingGenres(ctx context.Context, genres []modelsFilm.Genre) (map[string]modelsFilm.Genre, error) {
	genreNames := make([]string, len(genres))
	for i, genre := range genres {
		genreNames[i] = genre.Name
	}

	existingGenres, err := s.repository.GetGenresByNames(ctx, genreNames)
	if err != nil {
		return nil, err
	}

	existingGenresMap := make(map[string]modelsFilm.Genre)
	for _, genre := range existingGenres {
		existingGenresMap[genre.Name] = genre
	}

	return existingGenresMap, nil
}
