package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/errors"
	"film-management/pkg/validation"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"time"
)

// MakeViewFilmEndpoint is an endpoint for ViewAd.
func MakeViewFilmEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(ViewFilmRequest)
		if !ok {
			return ViewFilmResponse{}, errors.ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return ViewFilmResponse{Err: errValidate}, nil
		}

		// Parse UUID
		parseUUID, err := uuid.Parse(reqForm.UUID)
		if err != nil {
			return ViewFilmResponse{Err: err}, nil
		}

		if item, errViewFilm := s.ViewFilm(ctx, parseUUID); errViewFilm != nil {
			return ViewFilmResponse{Err: errViewFilm}, nil
		} else {
			return ViewFilmResponse{
				Item: domainFilmToItemViewFilm(item),
			}, nil
		}
	}
}

// ViewFilmRequest is a request for ViewFilm.
type ViewFilmRequest struct {
	UUID string `json:"uuid" validate:"required,uuid4" swaggerignore:"true"`
}

// Validate is a method to validate form.
func (r *ViewFilmRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// ViewFilmResponse is a response for ViewFilm.
type ViewFilmResponse struct {
	Item ItemViewFilm `json:"item,omitempty"`
	Err  error        `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r ViewFilmResponse) Failed() error { return r.Err }

// ItemViewFilm is a response for ViewFilm.
type ItemViewFilm struct {
	UUID        uuid.UUID   `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string      `json:"title" example:"The Shawshank Redemption"`
	Director    string      `json:"director" example:"Frank Darabont"`
	Genres      []string    `json:"genres" example:"drama,crime"`
	ReleaseDate string      `json:"release_date" example:"1994-09-23"`
	Casts       []string    `json:"casts" example:"Tim Robbins,Morgan Freeman"`
	Synopsis    string      `json:"synopsis" example:"This is a synopsis."`
	CreatedAt   string      `json:"created_at" example:"2021-01-01 00:00:00"`
	UpdatedAt   string      `json:"updated_at" example:"2021-01-01 00:00:00"`
	Creator     ItemCreator `json:"creator"`
}

// ItemCreator is a response for ViewFilm.
type ItemCreator struct {
	UUID     uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string    `json:"username" example:"test"`
}

// domainFilmToItemViewFilm is a method to convert domain Film to Item Film.
func domainFilmToItemViewFilm(item models.Film) ItemViewFilm {
	return ItemViewFilm{
		UUID:        item.UUID,
		Title:       item.Title,
		Director:    item.Director.Name,
		Genres:      convertGenresToStrings(item.Genres),
		ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
		Casts:       convertCastsToStrings(item.Casts),
		Synopsis:    item.Synopsis,
		CreatedAt:   time.Unix(item.CreatedAt, 0).Format(time.DateTime),
		UpdatedAt:   time.Unix(item.UpdatedAt, 0).Format(time.DateTime),
		Creator: ItemCreator{
			UUID:     item.Creator.UUID,
			Username: item.Creator.Username,
		},
	}
}
