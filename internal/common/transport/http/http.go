package http

import (
	"film-management/config"
	"film-management/pkg/transport/http/middlewares"
	"film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"io"
	"net/http"

	// Import for generating documentation.
	_ "film-management/docs"
)

const (
	APIPath = "/api/v1/"

	HealthCheckPath = APIPath + "health"
	SwaggerPath     = APIPath + "swagger"
)

// @title Film management service API
// @version 1.0
// @description This is a film management service.
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// NewHTTPHandlers is a function that returns a http.Handler that makes a set of endpoints available on predefined paths.
func NewHTTPHandlers(cfg *config.Config, logger *zap.Logger) http.Handler {
	r := mux.NewRouter()

	// CORS
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middlewares.CORSMiddleware(cfg.HTTP.CorsAllowedOrigins, logger))

	// Recovery
	r.Use(middlewares.RecoveryMiddleware(logger))

	// Routes
	//
	// Health Check
	r.HandleFunc(HealthCheckPath, HealthCheckHandler)
	// Swagger
	r.PathPrefix(SwaggerPath).Handler(httpSwagger.WrapHandler)
	// Not found handler
	r.NotFoundHandler = http.HandlerFunc(response.NotFoundFunc)

	return r
}

// HealthCheckHandler Health Check godoc
// @Summary Health Check
// @Description Health Check
// @Tags Common
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthCheckResponse "OK"
// @Router /health [get] .
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, `{"alive": true}`)

	if err != nil {
		return
	}
}

// HealthCheckResponse is a response struct for Health Check.
type HealthCheckResponse struct {
	Alive bool `json:"alive" example:"true"`
}
