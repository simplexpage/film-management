package auth

import (
	"context"
	"film-management/pkg/auth"
	customError "film-management/pkg/errors"
	httpResponse "film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
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
func Middleware(notAuthUrls []string, authService Service) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for some urls
			requestPath := r.URL.Path
			for _, value := range notAuthUrls {
				if strings.Contains(requestPath, value) {
					next.ServeHTTP(w, r)

					return
				}
			}

			// Get authorization header
			tokenHeader := r.Header.Get(AuthorizationHeader)

			// Check if token is missing
			if tokenHeader == "" {
				httpResponse.EncodeError(r.Context(), customError.AuthError{Err: ErrMissingAuthToken}, w)

				return
			}

			// Check if token is valid
			tokenParts := strings.Split(tokenHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != AuthorizationPrefix {
				httpResponse.EncodeError(r.Context(), customError.AuthError{Err: ErrInvalidAuthToken}, w)

				return
			}

			// Check if token is empty
			if len(tokenParts[1]) == 0 {
				httpResponse.EncodeError(r.Context(), customError.AuthError{Err: ErrAuthTokenEmpty}, w)

				return
			}

			// Parse token
			token, err := authService.ParseAuthToken(tokenParts[1])
			if err != nil {
				httpResponse.EncodeError(r.Context(), customError.AuthError{Err: ErrWrongAuthToken}, w)

				return
			}

			// Add user id to context
			ctx := setUserIDToContext(r.Context(), token.UUID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// setUserIDToContext is a function for setting user ID to context.
func setUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}
