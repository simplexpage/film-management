package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/internal/film/domain/models"
	"film-management/pkg/validation"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"time"
)

// MakeUpdateFilmEndpoint is an endpoint for AddFilm.
func MakeUpdateFilmEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(UpdateFilmRequest)
		if !ok {
			return UpdateFilmResponse{}, ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return UpdateFilmResponse{Err: errValidate}, nil
		}

		// Parse UUID
		parseUUID, err := uuid.Parse(reqForm.UUID)
		if err != nil {
			return ViewFilmResponse{Err: err}, nil
		}

		// Parse creator UUID
		parseCreatorUUID, err := uuid.Parse(reqForm.CreatorID)
		if err != nil {
			return AddFilmResponse{Err: err}, nil
		}

		// Parse date
		parseDate, err := time.Parse(time.DateOnly, reqForm.ReleaseDate)
		if err != nil {
			return AddFilmResponse{Err: err}, nil
		}

		// Prepare a film model
		model := &models.Film{
			UUID:        parseUUID,
			CreatorID:   parseCreatorUUID,
			Title:       reqForm.Title,
			Director:    reqForm.Director,
			ReleaseDate: parseDate,
			Cast:        reqForm.Cast,
			Synopsis:    reqForm.Synopsis,
		}

		if errUpdateFilm := s.UpdateFilm(ctx, model); errUpdateFilm != nil {
			return UpdateFilmResponse{Err: errUpdateFilm}, nil
		}

		return UpdateFilmResponse{}, nil
	}
}

// UpdateFilmRequest is a request for Update Film.
type UpdateFilmRequest struct {
	UUID      string `json:"uuid" validate:"required,uuid4" swaggerignore:"true"`
	CreatorID string `json:"creatorID" validate:"required,uuid4" swaggerignore:"true"`

	Title       string `json:"title" validate:"required,min=3,max=100" example:"Garry Potter"`
	Director    string `json:"director" validate:"required,min=3,max=40" example:"John Doe"`
	ReleaseDate string `json:"releaseDate" validate:"required,customDate" example:"2021-01-01"`
	Cast        string `json:"cast" validate:"required" example:"John Doe, Jane Doe, Foo Bar, Baz Quux"`
	Synopsis    string `json:"synopsis" validate:"required,min=10,max=1000" example:"This is a synopsis."`
}

// Validate is a method to validate form.
func (r *UpdateFilmRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Register custom validator "customDate"
	err = customValidator.GetValidate().RegisterValidation("customDate", DateValidator)
	if err != nil {
		return err
	}

	// Add translation for "customDate"
	err = customValidator.AddTranslation("customDate", fmt.Sprintf("{0} must be valid (YYYY-MM-DD)"))
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// UpdateFilmResponse is a response for UpdateAd.
type UpdateFilmResponse struct {
	Err error `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r UpdateFilmResponse) Failed() error { return r.Err }
