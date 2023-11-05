package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
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
			return ViewFilmResponse{}, ErrInvalidRequest
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
				Item: domainFilmToItemFilm(item),
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

	return customValidator.Validate(r)
}

// ViewFilmResponse is a response for ViewFilm.
type ViewFilmResponse struct {
	Item ItemFilm `json:"item,omitempty"`
	Err  error    `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r ViewFilmResponse) Failed() error { return r.Err }

// ItemFilm is a response for ViewFilm.
type ItemFilm struct {
	UUID        uuid.UUID   `json:"uuid"`
	Title       string      `json:"title"`
	Director    string      `json:"director"`
	ReleaseDate string      `json:"release_date"`
	Cast        string      `json:"cast"`
	Synopsis    string      `json:"synopsis"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Creator     ItemCreator `json:"creator"`
}

// ItemCreator is a response for ViewFilm.
type ItemCreator struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
}

// domainFilmToItemFilm is a method to convert domain Film to Item Film.
func domainFilmToItemFilm(item models.Film) ItemFilm {
	return ItemFilm{
		UUID:        item.UUID,
		Title:       item.Title,
		Director:    item.Director,
		ReleaseDate: item.ReleaseDate.Format(time.DateOnly),
		Cast:        item.Cast,
		Synopsis:    item.Synopsis,
		CreatedAt:   time.Unix(item.CreatedAt, 0).Format(time.DateTime),
		UpdatedAt:   time.Unix(item.UpdatedAt, 0).Format(time.DateTime),
		Creator: ItemCreator{
			UUID:     item.Creator.UUID,
			Username: item.Creator.Username,
		},
	}
}
