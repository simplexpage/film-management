package endpoints

import (
	"context"
	"film-management/internal/user/domain"
	"film-management/pkg/errors"
	"film-management/pkg/validation"
	"github.com/go-kit/kit/endpoint"
	"time"
)

// MakeLoginEndpoint is an endpoint for Login.
func MakeLoginEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(LoginRequest)
		if !ok {
			return LoginResponse{}, errors.ErrInvalidRequest
		}

		// Validate form
		if errValidate := reqForm.Validate(); errValidate != nil {
			return LoginResponse{Err: errValidate}, nil
		}

		// Login
		authToken, expirationTime, err := s.Login(ctx, reqForm.Username, reqForm.Password)

		return LoginResponse{AuthToken: authToken, ExpiredAt: expirationTime, Err: err}, nil
	}
}

// LoginRequest is a request for Login.
type LoginRequest struct {
	Username string `json:"username" validate:"required,username,min=5,max=40" example:"test123"`
	Password string `json:"password" validate:"required,min=8,max=30" example:"12345678"`
}

// Validate is a method to validate form.
func (r *LoginRequest) Validate() error {
	// Get custom validator
	customValidator, err := validation.GetValidator()
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// LoginResponse is a response for Login.
type LoginResponse struct {
	AuthToken string    `json:"auth_token,omitempty" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIj"`
	ExpiredAt time.Time `json:"expired_at,omitempty" example:"2023-11-09T15:21:15.973955426Z"`
	Err       error     `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r LoginResponse) Failed() error { return r.Err }
