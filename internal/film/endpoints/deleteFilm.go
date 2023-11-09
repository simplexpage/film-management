package endpoints

import (
	"context"
	"film-management/internal/film/domain"
	"film-management/pkg/errors"
	"film-management/pkg/validation"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

func MakeDeleteFilmEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(DeleteFilmRequest)
		if !ok {
			return DeleteFilmResponse{}, errors.ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return DeleteFilmResponse{Err: errValidate}, nil
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

		// Delete a film
		if errDeleteFilm := s.DeleteFilm(ctx, parseUUID, parseCreatorUUID); errDeleteFilm != nil {
			return DeleteFilmResponse{Err: errDeleteFilm}, nil
		}

		return DeleteFilmResponse{}, nil
	}
}

// DeleteFilmRequest is a request for Delete Film.
type DeleteFilmRequest struct {
	UUID      string `json:"uuid" validate:"required,uuid4" swaggerignore:"true"`
	CreatorID string `json:"creatorID" validate:"required,uuid4" swaggerignore:"true"`
}

// Validate is a method to validate form.
func (r *DeleteFilmRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// DeleteFilmResponse is a response for UpdateAd.
type DeleteFilmResponse struct {
	Err error `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r DeleteFilmResponse) Failed() error { return r.Err }
