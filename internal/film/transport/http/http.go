package http

import (
	"context"
	"film-management/config"
	httpCommon "film-management/internal/common/transport/http"
	"film-management/internal/film/endpoints"
	"film-management/pkg/auth"
	httpTransport "film-management/pkg/transport/http"
	"film-management/pkg/transport/http/middlewares"
	"film-management/pkg/transport/http/response"
	httpKitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	APIPath = httpCommon.APIPath + "films/"
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
	r.Handle(APIPath, addAdHandler).Methods(http.MethodPost)
	// Update a film
	r.Handle(APIPath+"{id}", updateAdHandler).Methods(http.MethodPut)
	// View a film
	r.Handle(APIPath+"{id}", viewAdHandler).Methods(http.MethodGet)
	// View all films
	r.Handle(APIPath, viewAllFilmsHandler).Methods(http.MethodGet)
	// Delete a film
	r.Handle(APIPath+"{id}", deleteAdHandler).Methods(http.MethodDelete)

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
// @Router /api/v1/films/ [post] .
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
// @Router /api/v1/films/{id} [put] .
func decodeHTTPUpdateFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.UpdateFilmRequest

	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "id")
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
// @Router /api/v1/films/{id} [get] .
func decodeHTTPViewFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "id")
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
// @Param title query string false "title" example(Star Wars)
// @Param release_date query string false "date" example(2023-12-11 or 2023-10-11:2023-12-11)
// @Param genres query string false "genres" example(action,adventure)
// @Param sort query string false "sort" example(title.asc or title.desc or release_date.asc or release_date.desc)
// @Param limit query string false "limit" example(10)
// @Param offset query string false "offset" example(1)
// @Success 200 {object} response.SuccessResponse{data=endpoints.ViewAllFilmsResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/films/ [get] .
func decodeHTTPViewAllFilmsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.ViewAllFilmsRequest

	// Get sort from HTTP request
	req.Sort = r.URL.Query().Get("sort")

	// Get limit from HTTP request
	if err := httpTransport.GetIntParamFromHTTPRequest("limit", r, &req.Limit); err != nil {
		return nil, err
	}

	// Get offset from HTTP request
	if err := httpTransport.GetIntParamFromHTTPRequest("offset", r, &req.Offset); err != nil {
		return nil, err
	}

	// Get filters from HTTP request
	req.Title = r.URL.Query().Get("title")
	req.ReleaseDate = r.URL.Query().Get("release_date")

	if genres := r.URL.Query().Get("genres"); genres != "" {
		req.Genres = strings.Split(genres, ",")
	}

	return req, nil
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
// @Router /api/v1/films/{id} [delete] .
func decodeHTTPDeleteFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "id")
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
