package http

import (
	"context"
	"film-management/config"
	httpCommon "film-management/internal/common/transport/http"
	"film-management/internal/user/endpoints"
	httpTransport "film-management/pkg/transport/http"
	"film-management/pkg/transport/http/middlewares"
	"film-management/pkg/transport/http/response"
	httpKitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
)

const (
	APIPath = httpCommon.APIPath + "user/"

	RegisterPath = APIPath + "register"
	LoginPath    = APIPath + "login"
)

// NewHTTPHandlers is a function that returns a http.Handler that makes a set of endpoints available on predefined paths.
func NewHTTPHandlers(endpoints endpoints.SetEndpoints, cfg *config.Config, logger *zap.Logger) http.Handler {
	options := []httpKitTransport.ServerOption{
		httpKitTransport.ServerErrorHandler(httpTransport.NewLogErrorHandler(logger)),
		httpKitTransport.ServerErrorEncoder(response.EncodeServerError),
	}

	// Handlers

	// Register User
	registerHandler := httpKitTransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeHTTPRegisterRequest,
		response.EncodeHTTPResponse,
		options...,
	)

	// Login User
	loginHandler := httpKitTransport.NewServer(
		endpoints.LoginEndpoint,
		decodeHTTPLoginRequest,
		response.EncodeHTTPResponse,
		options...,
	)

	r := mux.NewRouter()

	// CORS
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middlewares.CORSMiddleware(cfg.HTTP.CorsAllowedOrigins, logger))

	// Recovery
	r.Use(middlewares.RecoveryMiddleware(logger))

	// Routes

	// User
	//
	// Register
	r.Handle(RegisterPath, registerHandler).Methods(http.MethodPost, http.MethodOptions)
	// Login
	r.Handle(LoginPath, loginHandler).Methods(http.MethodPost, http.MethodOptions)
	// Set custom error handlers
	response.SetErrorHandlers(r)

	return r
}

// Registration godoc
// @Summary Registration
// @Description Registration
// @Tags User
// @Accept json
// @Produce json
// @Param form body endpoints.RegisterRequest true "Register form"
// @Success 200 {object} response.SuccessResponse{data=endpoints.RegisterResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /user/register [post] .
func decodeHTTPRegisterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.RegisterRequest

	if e := jsoniter.NewDecoder(r.Body).Decode(&reqForm); e != nil {
		return nil, httpTransport.ErrJSONDecode
	}

	return reqForm, nil
}

// Authentication godoc
// @Summary Login
// @Description Login
// @Tags User
// @Accept json
// @Produce json
// @Param form body endpoints.LoginRequest true "Login form"
// @Success 200 {object} response.SuccessResponse{data=endpoints.LoginResponse} "Success"
// @Failure 400 {object} response.ErrorResponse	"Bad Request"
// @Failure 422 {object} response.ErrorResponseValidation "Data Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /user/login [post] .
func decodeHTTPLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqForm endpoints.LoginRequest

	if e := jsoniter.NewDecoder(r.Body).Decode(&reqForm); e != nil {
		return nil, httpTransport.ErrJSONDecode
	}

	return reqForm, nil
}
