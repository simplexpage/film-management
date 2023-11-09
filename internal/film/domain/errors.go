package domain

import "errors"

var (
	ErrFilmCreate           = errors.New("failed to create film")
	ErrFilmUpdate           = errors.New("failed to update film")
	ErrFilmDelete           = errors.New("failed to delete film")
	ErrFilmFind             = errors.New("failed to find film")
	ErrFilmFindAll          = errors.New("failed to find all films")
	ErrFilmNotPermission    = errors.New("access denied, you do not have permission to edit this film")
	ErrFilmNotFound         = errors.New("film not found")
	ErrFilmExistsWithTitle  = errors.New("film already exists with the same title")
	ErrFilmCheckExistence   = errors.New("failed to check film existence")
	ErrFilmCreateCast       = errors.New("failed to create cast")
	ErrFilmGetCastsByNames  = errors.New("failed to get casts by names")
	ErrFilmGetGenresByNames = errors.New("failed to get genres by names")
	ErrFilmFindGenres       = errors.New("failed to find genres")
	ErrFilmGenresNotFound   = errors.New("genres do not exist in the database")
	ErrFilmFilterWrong      = errors.New("filter wrong")
	ErrFilmUnknownField     = errors.New("unknown field")
)
