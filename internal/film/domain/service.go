package domain

import (
	"context"
	modelsFilm "film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"film-management/pkg/validation"
	"fmt"
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
	if err := s.checkDuplicateFilm(ctx, model.Title, model.UUID, modelsFilm.OperationAdd); err != nil {
		return err
	}

	// Set and validate film genres
	if err := s.setAndValidateFilmGenres(ctx, model); err != nil {
		return err
	}

	// Set and create film casts
	if err := s.setAndCreateFilmCasts(ctx, model); err != nil {
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

	// Check permission
	if s.checkFilmPermission(model.CreatorID, filmFromDB.CreatorID) {
		return ErrFilmNotPermission
	}

	// Check duplicate film with title
	if err := s.checkDuplicateFilm(ctx, model.Title, filmFromDB.UUID, modelsFilm.OperationUpdate); err != nil {
		return err
	}

	// Set and validate film genres
	if err := s.setAndValidateFilmGenres(ctx, model); err != nil {
		return err
	}

	// Set and create film casts
	if err := s.setAndCreateFilmCasts(ctx, model); err != nil {
		return err
	}

	// Set new film data
	filmFromDB.SetDataForUpdate(model)

	// Update a film in db
	return s.repository.UpdateFilm(ctx, &filmFromDB)
}

func (s service) ViewFilm(ctx context.Context, filmID uuid.UUID) (modelsFilm.Film, error) {
	return s.repository.FindOneFilmForViewByUUID(ctx, filmID)
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
func (s service) checkDuplicateFilm(ctx context.Context, title string, filmID uuid.UUID, operation modelsFilm.Operation) error {
	if err := s.repository.FilmExists(ctx, title, filmID, operation); err != nil {
		switch err {
		case ErrFilmExistsWithTitle:
			return validation.CustomError{Field: "title", Err: ErrFilmExistsWithTitle}
		default:
			return err
		}
	}

	return nil
}

// setAndValidateFilmGenres Set and validate film genres.
func (s service) setAndValidateFilmGenres(ctx context.Context, model *modelsFilm.Film) error {
	// Get existing genres in db
	existingGenres, err := s.getExistingGenres(ctx, model.Genres)
	if err != nil {
		return err
	}

	// Set and validate film genres
	for i, genre := range model.Genres {
		if existingGenre, ok := existingGenres[genre.Name]; ok {
			model.Genres[i] = existingGenre
		} else {
			return validation.CustomError{Field: "genres", Err: fmt.Errorf("genre %s does not exist in the database", genre.Name)}
		}
	}

	return nil
}

// getExistingGenres Get existing genres in db.
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

// setAndCreateFilmCasts Set and create film casts.
func (s service) setAndCreateFilmCasts(ctx context.Context, model *modelsFilm.Film) error {
	// Get existing casts in db
	existingCasts, err := s.getExistingCasts(ctx, model.Casts)
	if err != nil {
		return err
	}

	// Set film casts
	for i, cast := range model.Casts {
		if existingCast, ok := existingCasts[cast.Name]; ok {
			model.Casts[i] = existingCast
		} else {
			newCast, err := s.repository.CreateCast(ctx, &cast)
			if err != nil {
				return err
			}
			model.Casts[i] = *newCast
		}
	}

	return nil
}

// getExistingCasts Get existing casts in db.
func (s service) getExistingCasts(ctx context.Context, casts []modelsFilm.Cast) (map[string]modelsFilm.Cast, error) {
	castNames := make([]string, len(casts))
	for i, cast := range casts {
		castNames[i] = cast.Name
	}

	existingCasts, err := s.repository.GetCastsByNames(ctx, castNames)
	if err != nil {
		return nil, err
	}

	existingCastsMap := make(map[string]modelsFilm.Cast)
	for _, cast := range existingCasts {
		existingCastsMap[cast.Name] = cast
	}

	return existingCastsMap, nil
}
