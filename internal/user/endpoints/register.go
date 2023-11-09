package endpoints

import (
	"context"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/models"
	"film-management/pkg/errors"
	"film-management/pkg/validation"
	"github.com/go-kit/kit/endpoint"
)

// MakeRegisterEndpoint is an endpoint for Register.
func MakeRegisterEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(RegisterRequest)
		if !ok {
			return RegisterResponse{}, errors.ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return RegisterResponse{Err: errValidate}, nil
		}

		// Prepare a User model
		model := &models.User{
			Username: reqForm.Username,
			Password: reqForm.Password,
		}

		// Register user
		if errRegister := s.Register(ctx, model); errRegister != nil {
			return RegisterResponse{Err: errRegister}, nil
		}

		return RegisterResponse{UUID: model.UUID.String(), Username: model.Username}, nil
	}
}

// RegisterRequest is a request for Register.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,username,min=5,max=40" example:"test123"`
	Password string `json:"password" validate:"required,min=8,max=30" example:"12345678"`
}

// Validate is a method to validate form.
func (r *RegisterRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// RegisterResponse is a response for Register.
type RegisterResponse struct {
	UUID     string `json:"uuid,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string `json:"username,omitempty" example:"test123"`
	Err      error  `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r RegisterResponse) Failed() error { return r.Err }
