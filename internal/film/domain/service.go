package domain

import (
	"context"
	modelsFilm "film-management/internal/film/domain/models"
	customError "film-management/pkg/errors"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

	// Create film in db
	if err := s.repository.CreateFilm(ctx, model); err != nil {
		return ErrFilmCreate
	}

	return nil
}

func (s service) UpdateFilm(ctx context.Context, model *modelsFilm.Film) error {
	// Get film from db
	filmFromDB, err := s.getFilmFromDB(ctx, model.UUID)
	if err != nil {
		return err
	}

	// Check permission
	if errPermission := s.checkFilmPermission(model.CreatorID, filmFromDB.CreatorID); errPermission != nil {
		return errPermission
	}

	// Check duplicate film with title
	if errDuplicate := s.checkDuplicateFilm(ctx, model.Title, filmFromDB.UUID, modelsFilm.OperationUpdate); errDuplicate != nil {
		return errDuplicate
	}

	// Set and validate film genres
	if errGenres := s.setAndValidateFilmGenres(ctx, model); errGenres != nil {
		return errGenres
	}

	// Set and create film casts
	if errCasts := s.setAndCreateFilmCasts(ctx, model); errCasts != nil {
		return errCasts
	}

	// Set new film data
	filmFromDB.SetDataForUpdate(model)

	// Update a film in db
	if errUpdate := s.repository.UpdateFilm(ctx, &filmFromDB); errUpdate != nil {
		return ErrFilmUpdate
	}

	return nil
}

// ViewFilm View a film.
func (s service) ViewFilm(ctx context.Context, filmID uuid.UUID) (modelsFilm.Film, error) {
	// Get film from db
	filmFromDB, err := s.repository.FindOneFilmForViewByUUID(ctx, filmID)
	if err != nil {
		switch {
		case errors.Is(err, ErrFilmNotFound):
			return modelsFilm.Film{}, customError.NotFoundError{Err: ErrFilmNotFound}
		default:
			return modelsFilm.Film{}, ErrFilmFind
		}
	}

	return filmFromDB, nil
}

// ViewAllFilms View all films.
func (s service) ViewAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) ([]modelsFilm.Film, pagination.Pagination, error) {
	filmsFromDB, p, err := s.repository.FindAllFilms(ctx, filterSortPagination)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}

	return filmsFromDB, p, nil
}

// DeleteFilm Delete a film.
func (s service) DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) error {
	// Get film from db
	filmFromDB, err := s.getFilmFromDB(ctx, filmID)
	if err != nil {
		return err
	}

	// Check permission
	if errPermission := s.checkFilmPermission(userID, filmFromDB.CreatorID); errPermission != nil {
		return errPermission
	}

	// Delete a film in db
	if errDelete := s.repository.DeleteFilm(ctx, filmID); errDelete != nil {
		return ErrFilmDelete
	}

	return nil
}

// getFilmFromDB Get film from db.
func (s service) getFilmFromDB(ctx context.Context, filmID uuid.UUID) (modelsFilm.Film, error) {
	filmFromDB, err := s.repository.FindOneFilmByUUID(ctx, filmID)
	if err != nil {
		switch {
		case errors.Is(err, ErrFilmNotFound):
			return modelsFilm.Film{}, customError.NotFoundError{Err: ErrFilmNotFound}
		default:
			return modelsFilm.Film{}, ErrFilmFind
		}
	}

	return filmFromDB, nil
}

// checkFilmPermission Check if a film user is the same as the user who wants to do an action.
func (s service) checkFilmPermission(firstUserID uuid.UUID, secondUserID uuid.UUID) error {
	if firstUserID != secondUserID {
		return customError.PermissionError{Err: ErrFilmNotPermission}
	}

	return nil
}

// checkDuplicateFilm Check if a film with the same title already exists.
func (s service) checkDuplicateFilm(ctx context.Context, title string, filmID uuid.UUID, operation modelsFilm.Operation) error {
	if err := s.repository.FilmExistsWithTitle(ctx, title, filmID, operation); err != nil {
		switch {
		case errors.Is(err, ErrFilmExistsWithTitle):
			return customError.ValidationError{Field: "title", Err: ErrFilmExistsWithTitle}
		default:
			return ErrFilmCheckExistence
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
			return customError.ValidationError{Field: "genres", Err: fmt.Errorf("genre %s does not exist in the database", genre.Name)}
		}
	}

	return nil
}

// getExistingGenres Get existing genres in db.
func (s service) getExistingGenres(ctx context.Context, genres []modelsFilm.Genre) (map[string]modelsFilm.Genre, error) {
	// Get genre names
	genreNames := make([]string, len(genres))
	for i, genre := range genres {
		genreNames[i] = genre.Name
	}

	// Get existing genres in db
	existingGenres, err := s.repository.GetGenresByNames(ctx, genreNames)
	if err != nil {
		return nil, ErrFilmGetGenresByNames
	}

	// Create map of existing genres
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
		// If cast exists in db, set it
		if existingCast, ok := existingCasts[cast.Name]; ok {
			model.Casts[i] = existingCast
		} else {
			// If cast does not exist in db, create it
			newCast, errCreate := s.repository.CreateCast(ctx, &cast)
			if errCreate != nil {
				return ErrFilmCreateCast
			}
			model.Casts[i] = *newCast
		}
	}

	return nil
}

// getExistingCasts Get existing casts in db.
func (s service) getExistingCasts(ctx context.Context, casts []modelsFilm.Cast) (map[string]modelsFilm.Cast, error) {
	// Get cast names
	castNames := make([]string, len(casts))
	for i, cast := range casts {
		castNames[i] = cast.Name
	}

	// Get existing casts in db
	existingCasts, err := s.repository.GetCastsByNames(ctx, castNames)
	if err != nil {
		return nil, ErrFilmGetCastsByNames
	}

	// Create map of existing casts
	existingCastsMap := make(map[string]modelsFilm.Cast)
	for _, cast := range existingCasts {
		existingCastsMap[cast.Name] = cast
	}

	return existingCastsMap, nil
}
