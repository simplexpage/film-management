package endpoints

import (
	"film-management/internal/user/domain"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
)

// SetEndpoints collects all the endpoints that compose an ad service.
type SetEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

// NewEndpoints returns a SetEndpoints that wraps the provided server, and wires in all the provided middlewares.
func NewEndpoints(s domain.Service, logger *zap.Logger) SetEndpoints {
	var registerEndpoint endpoint.Endpoint
	{
		registerEndpoint = MakeRegisterEndpoint(s)
		registerEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "Register")))(registerEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(s)
		loginEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "Login")))(loginEndpoint)
	}

	return SetEndpoints{
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
	}
}
