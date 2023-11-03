package http

import (
	"context"
	"film-management/config"
	httpCommon "film-management/internal/common/transport/http"
	"film-management/internal/film/endpoints"
	"film-management/pkg/auth"
	"film-management/pkg/query"
	httpTransport "film-management/pkg/transport/http"
	"film-management/pkg/transport/http/middlewares"
	"film-management/pkg/transport/http/response"
	httpKitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	APIPath = httpCommon.APIPath + "film/"

	AddFilmPath      = APIPath + "add"
	UpdateFilmPath   = APIPath + "update/{uuid}"
	ViewFilmPath     = APIPath + "view/{uuid}"
	ViewAllFilmsPath = APIPath + "view-all"
	DeleteFilmPath   = APIPath + "delete/{uuid}"
)

// NewHTTPHandlers is a function that returns a http.Handler that makes a set of endpoints available on predefined paths.
func NewHTTPHandlers(endpoints endpoints.SetEndpoints, authService middlewares.AuthService, cfg *config.Config, logger *zap.Logger) http.Handler {
	options := []httpKitTransport.ServerOption{
		httpKitTransport.ServerErrorHandler(httpTransport.NewLogErrorHandler(logger)),
		httpKitTransport.ServerErrorEncoder(response.EncodeServerError),
	}

	// Handlers
	// Add a film
	addAdHandler := httpKitTransport.NewServer(
		endpoints.AddFilmEndpoint,
		decodeHTTPAddFilmRequest,
		response.EncodeHTTPResponse,
		options...,
	)
	// Update the film
	updateAdHandler := httpKitTransport.NewServer(
		endpoints.UpdateFilmEndpoint,
		decodeHTTPUpdateFilmRequest,
		response.EncodeHTTPResponse,
		options...,
	)
	// View the film
	viewAdHandler := httpKitTransport.NewServer(
		endpoints.ViewFilmEndpoint,
		decodeHTTPViewFilmRequest,
		response.EncodeHTTPResponse,
		options...,
	)
	// View all films
	viewAllFilmsHandler := httpKitTransport.NewServer(
		endpoints.ViewAllFilmsEndpoint,
		decodeHTTPViewAllFilmsRequest,
		response.EncodeHTTPResponse,
		options...,
	)
	// Delete the film
	deleteAdHandler := httpKitTransport.NewServer(
		endpoints.DeleteFilmEndpoint,
		decodeHTTPDeleteFilmRequest,
		response.EncodeHTTPResponse,
		options...,
	)

	r := mux.NewRouter()

	// CORS
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middlewares.CORSMiddleware(cfg.HTTP.CorsAllowedOrigins, logger))

	// Recovery
	r.Use(middlewares.RecoveryMiddleware(logger))

	// AUTH
	r.Use(middlewares.AuthMiddleware(cfg.HTTP.NotAuthUrls, authService))

	// Routes

	// Film
	//
	// Add a film
	r.Handle(AddFilmPath, addAdHandler).Methods(http.MethodPost, http.MethodOptions)
	// Update a film
	r.Handle(UpdateFilmPath, updateAdHandler).Methods(http.MethodPut, http.MethodOptions)
	// View a film
	r.Handle(ViewFilmPath, viewAdHandler).Methods(http.MethodGet, http.MethodOptions)
	// View all films
	r.Handle(ViewAllFilmsPath, viewAllFilmsHandler).Methods(http.MethodGet, http.MethodOptions)
	// Delete a film
	r.Handle(DeleteFilmPath, deleteAdHandler).Methods(http.MethodDelete, http.MethodOptions)

	// Set custom error handlers
	response.SetErrorHandlers(r)

	return r
}

// AddFilm godoc
// @Summary Add a film
// @Description Add a film
// @Tags Film
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param form body endpoints.AddFilmRequest true "Add Film Form"
// @Success 200 {object} response.SuccessResponse{data=endpoints.AddFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/film/add [post] .
func decodeHTTPAddFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.AddFilmRequest

	// Get UserID from context
	userID, errUserID := auth.GetUserIDFromContext(r.Context())
	if errUserID != nil {
		return nil, httpTransport.ErrContextUserID
	}

	// Decode JSON
	if e := jsoniter.NewDecoder(r.Body).Decode(&reqForm); e != nil {
		return nil, httpTransport.ErrJSONDecode
	}

	// Set CreatorID
	reqForm.CreatorID = userID

	return reqForm, nil
}

// UpdateFilm godoc
// @Summary Update a film
// @Description Update a film
// @Tags Film
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param uuid path string true "Film UUID"
// @Param form body endpoints.UpdateFilmRequest true "Update film form"
// @Success 200 {object} response.SuccessResponse{data=endpoints.UpdateFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/film/update/{uuid} [put] .
func decodeHTTPUpdateFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.UpdateFilmRequest

	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "uuid")
	if err != nil {
		return nil, err
	}

	// Get UserID from context
	userID, errUserID := auth.GetUserIDFromContext(r.Context())
	if errUserID != nil {
		return nil, httpTransport.ErrContextUserID
	}

	// Decode JSON
	if e := jsoniter.NewDecoder(r.Body).Decode(&reqForm); e != nil {
		return nil, httpTransport.ErrJSONDecode
	}

	// Set UUID and CreatorID
	reqForm.UUID = uuidFromPath
	reqForm.CreatorID = userID

	return reqForm, nil
}

// ViewFilm godoc
// @Summary View a film
// @Description View a film
// @Tags Film
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param uuid path string true "Film UUID"
// @Success 200 {object} response.SuccessResponse{data=endpoints.ViewFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/film/view/{uuid} [get] .
func decodeHTTPViewFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "uuid")
	if err != nil {
		return nil, err
	}

	return endpoints.ViewFilmRequest{UUID: uuidFromPath}, nil
}

// ViewAllFilms godoc
// @Summary View all films
// @Description View all films
// @Tags Film
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param title query string false "title" example(Title)
// @Param release_date query string false "date" example(2023-12-11 or 2023-10-11:2023-12-11)
// @Param sort query string false "sort" example(title.asc, title.desc, release_date.asc, release_date.desc)
// @Param limit query string false "limit" example(10)
// @Param offset query string false "offset" example(1)
// @Success 200 {object} response.SuccessResponse{data=endpoints.ViewAllFilmsResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/film/view-all [get] .
func decodeHTTPViewAllFilmsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	filterSortLimit, err := query.NewFilterSortLimitFromHTTPRequest(r, FilmAvailableFilterFields{}, "release_date.desc")
	if err != nil {
		return nil, err
	}

	return endpoints.ViewAllFilmsRequest{FilterSortLimit: filterSortLimit}, nil
}

// FilmAvailableFilterFields is a struct that contains available fields for filtering.
type FilmAvailableFilterFields struct {
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
}

// DeleteFilm godoc
// @Summary Delete a film
// @Description Delete a film
// @Tags Film
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param uuid path string true "Film UUID"
// @Success 200 {object} response.SuccessResponse{data=endpoints.DeleteFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/film/delete/{uuid} [delete] .
func decodeHTTPDeleteFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "uuid")
	if err != nil {
		return nil, err
	}

	// Get UserID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		return nil, httpTransport.ErrContextUserID
	}

	return endpoints.DeleteFilmRequest{UUID: uuidFromPath, CreatorID: userID}, nil
}
