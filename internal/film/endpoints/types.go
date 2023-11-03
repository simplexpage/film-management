package endpoints

import (
	"film-management/internal/film/domain/models"
)

// Genre is a type for film genre.
type Genre string

const (
	// GenreUnknown is a constant for the "unknown" genre.
	GenreUnknown Genre = "unknown"
	// GenreAction is a constant for the "action" genre.
	GenreAction Genre = "action"
	// GenreComedy is a constant for the "comedy" genre.
	GenreComedy Genre = "comedy"
)

// GenreToEnum converts Genre to film.Genre.
func GenreToEnum(genreType Genre) models.Genre {
	switch genreType {
	case GenreAction:
		return models.GenreAction
	case GenreComedy:
		return models.GenreComedy
	default:
		return models.GenreUnknown
	}
}

// GenreFromEnum converts film.Genre to Genre.
func GenreFromEnum(filmGenre models.Genre) Genre {
	switch filmGenre {
	case models.GenreAction:
		return GenreAction
	case models.GenreComedy:
		return GenreComedy
	default:
		return GenreUnknown
	}
}
