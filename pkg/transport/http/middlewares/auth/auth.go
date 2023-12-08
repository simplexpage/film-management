package auth

import (
	"context"
	"film-management/pkg/auth"
	customError "film-management/pkg/errors"
	httpResponse "film-management/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

const (
	AuthorizationHeader string     = "Authorization"
	AuthorizationPrefix string     = "Bearer"
	ContextKeyUserID    ContextKey = "user_id"
)

var (
	ErrMissingAuthToken = errors.New("missing auth token")
	ErrInvalidAuthToken = errors.New("invalid auth token. Bearer token is expected")
	ErrAuthTokenEmpty   = errors.New("auth token is empty")
	ErrWrongAuthToken   = errors.New("wrong auth token")
)

type ContextKey string

type Service interface {
	ParseAuthToken(token string) (*auth.JwtClaims, error)
}

// Middleware is a middleware for authentication.
func Middleware(authService Service, notAuthUrls []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for some urls
		requestPath := c.Request.URL.Path
		for _, value := range notAuthUrls {
			if strings.Contains(requestPath, value) {
				c.Next()

				return
			}
		}

		// Get authorization header
		tokenHeader := c.Request.Header.Get(AuthorizationHeader)

		// Check if token is missing
		if tokenHeader == "" {
			httpResponse.EncodeError(c.Request.Context(), customError.AuthError{Err: ErrMissingAuthToken}, c.Writer)
			c.Abort()

			return
		}

		// Check if token is valid
		tokenParts := strings.Split(tokenHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != AuthorizationPrefix {
			httpResponse.EncodeError(c.Request.Context(), customError.AuthError{Err: ErrInvalidAuthToken}, c.Writer)
			c.Abort()

			return
		}

		// Check if token is empty
		if len(tokenParts[1]) == 0 {
			httpResponse.EncodeError(c.Request.Context(), customError.AuthError{Err: ErrAuthTokenEmpty}, c.Writer)
			c.Abort()

			return
		}

		// Parse token
		token, err := authService.ParseAuthToken(tokenParts[1])
		if err != nil {
			httpResponse.EncodeError(c.Request.Context(), customError.AuthError{Err: ErrWrongAuthToken}, c.Writer)
			c.Abort()

			return
		}

		// Add user id to context
		ctx := setUserIDToContext(c.Request.Context(), token.UUID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// setUserIDToContext is a function for setting user ID to context.
func setUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}
