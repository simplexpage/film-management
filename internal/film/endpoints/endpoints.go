package endpoints

import (
	"film-management/internal/film/domain"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
)

// SetEndpoints collects all the endpoints that compose an ad service.
type SetEndpoints struct {
	AddFilmEndpoint      endpoint.Endpoint
	UpdateFilmEndpoint   endpoint.Endpoint
	ViewFilmEndpoint     endpoint.Endpoint
	ViewAllFilmsEndpoint endpoint.Endpoint
	DeleteFilmEndpoint   endpoint.Endpoint
}

// NewEndpoints returns a SetEndpoints that wraps the provided server, and wires in all the provided middlewares.
func NewEndpoints(s domain.Service, logger *zap.Logger) SetEndpoints {
	var addFilmEndpoint endpoint.Endpoint
	{
		addFilmEndpoint = MakeAddFilmEndpoint(s)
		addFilmEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "AddFilm")))(addFilmEndpoint)
	}

	var updateFilmEndpoint endpoint.Endpoint
	{
		updateFilmEndpoint = MakeUpdateFilmEndpoint(s)
		updateFilmEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "UpdateFilm")))(updateFilmEndpoint)
	}

	var viewFilmEndpoint endpoint.Endpoint
	{
		viewFilmEndpoint = MakeViewFilmEndpoint(s)
		viewFilmEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "ViewFilm")))(viewFilmEndpoint)
	}

	var viewAllFilmsEndpoint endpoint.Endpoint
	{
		viewAllFilmsEndpoint = MakeViewAllFilmsEndpoint(s)
		viewAllFilmsEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "ViewAllFilms")))(viewAllFilmsEndpoint)
	}

	var deleteFilmEndpoint endpoint.Endpoint
	{
		deleteFilmEndpoint = MakeDeleteFilmEndpoint(s)
		deleteFilmEndpoint = NewLoggingMiddleware(logger.With(zap.String("method", "DeleteFilm")))(deleteFilmEndpoint)
	}

	return SetEndpoints{
		AddFilmEndpoint:      addFilmEndpoint,
		UpdateFilmEndpoint:   updateFilmEndpoint,
		ViewFilmEndpoint:     viewFilmEndpoint,
		ViewAllFilmsEndpoint: viewAllFilmsEndpoint,
		DeleteFilmEndpoint:   deleteFilmEndpoint,
	}
}
