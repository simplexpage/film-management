package http

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	// Import for generating documentation.
	_ "film-management/docs"
	swaggerfiles "github.com/swaggo/files"
)

// @title Film management service API
// @version 1.0
// @description This is a film management service.
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// SetHTTPRoutes is a function that makes a set of endpoints available on predefined paths.
func SetHTTPRoutes(router *gin.Engine) {
	// Routes
	//
	v1 := router.Group("/api/v1")
	// Health Check
	v1.GET("/health", HealthCheckHandler)
	// Swagger
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// HealthCheckHandler Health Check godoc
// @Summary Health Check
// @Description Health Check
// @Tags Common
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthCheckResponse "OK"
// @Router /health [get] .
// HealthCheckHandler Health Check
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}

// HealthCheckResponse is a response struct for Health Check.
type HealthCheckResponse struct {
	Alive bool `json:"alive" example:"true"`
}
