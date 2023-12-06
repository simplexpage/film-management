package http

import (
	"context"
	"film-management/internal/user/endpoints"
	httpTransport "film-management/pkg/transport/http"
	"film-management/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
	httpKitTransport "github.com/go-kit/kit/transport/http"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
)

// SetHTTPRoutes is a function that makes a set of endpoints available on predefined paths.
func SetHTTPRoutes(router *gin.Engine, endpoints endpoints.SetEndpoints, logger *zap.Logger) {
	options := []httpKitTransport.ServerOption{
		httpKitTransport.ServerErrorHandler(httpTransport.NewLogErrorHandler(logger)),
		httpKitTransport.ServerErrorEncoder(response.EncodeError),
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

	// Routes
	//
	// User
	v1 := router.Group("/api/v1/user")
	{
		// Register
		v1.POST("/register", gin.WrapH(registerHandler))
		// Login
		v1.POST("/login", gin.WrapH(loginHandler))
	}
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
