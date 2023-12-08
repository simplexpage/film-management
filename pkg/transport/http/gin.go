package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	httpKitTransport "github.com/go-kit/kit/transport/http"
)

// WithGinContext Middleware to pass gin context
func WithGinContext(handler *httpKitTransport.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request.WithContext(context.WithValue(c, "ginContext", c)))
	}
}

// GetGinContext Get gin context from context
func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginCtx, ok := ctx.Value("ginContext").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin context not found")
	}
	return ginCtx, nil
}
