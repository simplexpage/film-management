package domain

import "errors"

var (
	ErrFilmCreate           = errors.New("failed to create film")
	ErrFilmUpdate           = errors.New("failed to update film")
	ErrFilmDelete           = errors.New("failed to delete film")
	ErrFilmFind             = errors.New("failed to find film")
	ErrFilmFindAll          = errors.New("failed to find all films")
	ErrFilmNotPermission    = errors.New("film not permission")
	ErrFilmNotFound         = errors.New("film not found")
	ErrFilmExistsWithTitle  = errors.New("film already exists with the same title")
	ErrFilmCheckExistence   = errors.New("failed to check film existence")
	ErrFilmCreateGenre      = errors.New("failed to create genre")
	ErrFilmGetGenresByNames = errors.New("failed to get genres by names")
	ErrFilmGetCount         = errors.New("failed to get count")
)
