package endpoints

import (
	"context"
	"film-management/internal/user/domain"
	"film-management/pkg/validation"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"time"
)

// MakeLoginEndpoint is an endpoint for Login.
func MakeLoginEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqForm, ok := request.(LoginRequest)
		if !ok {
			return LoginResponse{}, ErrInvalidRequest
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

	// Register custom validator "username"
	err = customValidator.GetValidate().RegisterValidation("username", UsernameValidator)
	if err != nil {
		return err
	}

	// Add translation for "username"
	err = customValidator.AddTranslation("username", fmt.Sprintf("{0} must be valid (alphanumeric starting with letter)"))
	if err != nil {
		return err
	}

	// Validate form
	return customValidator.Validate(r)
}

// LoginResponse is a response for Login.
type LoginResponse struct {
	AuthToken string    `json:"auth_token"`
	ExpiredAt time.Time `json:"expired_at"`
	Err       error     `json:"err,omitempty" swaggerignore:"true"`
}

// Failed implements response.Failed.
func (r LoginResponse) Failed() error { return r.Err }
