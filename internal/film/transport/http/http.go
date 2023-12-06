package http

import (
	"context"
	"film-management/internal/film/endpoints"
	httpTransport "film-management/pkg/transport/http"
	"film-management/pkg/transport/http/middlewares/auth"
	"film-management/pkg/transport/http/response"
	"film-management/pkg/utils"
	"github.com/gin-gonic/gin"
	httpKitTransport "github.com/go-kit/kit/transport/http"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
)

// SetHTTPRoutes is a function that makes a set of endpoints available on predefined paths.
func SetHTTPRoutes(router *gin.Engine, endpoints endpoints.SetEndpoints) {
	// Handlers
	//
	// Add a film
	addAdHandler := httpKitTransport.NewServer(
		endpoints.AddFilmEndpoint,
		decodeHTTPAddFilmRequest,
		response.EncodeHTTPResponse,
	)
	// Update the film
	updateAdHandler := httpKitTransport.NewServer(
		endpoints.UpdateFilmEndpoint,
		decodeHTTPUpdateFilmRequest,
		response.EncodeHTTPResponse,
	)
	// View the film
	viewAdHandler := httpKitTransport.NewServer(
		endpoints.ViewFilmEndpoint,
		decodeHTTPViewFilmRequest,
		response.EncodeHTTPResponse,
	)
	// View all films
	viewAllFilmsHandler := httpKitTransport.NewServer(
		endpoints.ViewAllFilmsEndpoint,
		decodeHTTPViewAllFilmsRequest,
		response.EncodeHTTPResponse,
	)
	// Delete the film
	deleteAdHandler := httpKitTransport.NewServer(
		endpoints.DeleteFilmEndpoint,
		decodeHTTPDeleteFilmRequest,
		response.EncodeHTTPResponse,
	)

	// Routes
	//
	// Film
	v1 := router.Group("/api/v1/films")
	{
		// Add a film
		v1.POST("/", gin.WrapH(addAdHandler))
		// Update a film
		v1.PUT("/{id}", gin.WrapH(updateAdHandler))
		// View a film
		v1.GET("/{id}", gin.WrapH(viewAdHandler))
		// View all films
		v1.GET("/", gin.WrapH(viewAllFilmsHandler))
		// Delete a film
		v1.DELETE("/{id}", gin.WrapH(deleteAdHandler))
	}
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
// @Router /films/ [post] .
func decodeHTTPAddFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.AddFilmRequest

	// Get UserID from context
	userID, errUserID := utils.GetValueFromContext(r.Context(), auth.ContextKeyUserID)
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
// @Param id path string true "Film UUID"
// @Param form body endpoints.UpdateFilmRequest true "Update film form"
// @Success 200 {object} response.SuccessResponse{data=endpoints.UpdateFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /films/{id} [put] .
func decodeHTTPUpdateFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.UpdateFilmRequest

	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "id")
	if err != nil {
		return nil, err
	}

	// Get UserID from context
	userID, errUserID := utils.GetValueFromContext(r.Context(), auth.ContextKeyUserID)
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
// @Param id path string true "Film UUID"
// @Success 200 {object} response.SuccessResponse{data=endpoints.ViewFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /films/{id} [get] .
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
// @Router /films [get] .
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
// @Param id path string true "Film UUID"
// @Success 200 {object} response.SuccessResponse{data=endpoints.DeleteFilmResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /films/{id} [delete] .
func decodeHTTPDeleteFilmRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// Get UUID from path
	uuidFromPath, err := httpTransport.GetValueFromPath(r, "id")
	if err != nil {
		return nil, err
	}

	// Get UserID from context
	userID, err := utils.GetValueFromContext(r.Context(), auth.ContextKeyUserID)
	if err != nil {
		return nil, httpTransport.ErrContextUserID
	}

	return endpoints.DeleteFilmRequest{UUID: uuidFromPath, CreatorID: userID}, nil
}
