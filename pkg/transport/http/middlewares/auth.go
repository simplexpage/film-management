package middlewares

import (
	"film-management/pkg/auth"
	httpTransport "film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var (
	ErrMissingAuthToken = errors.New("missing auth token")
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrAuthTokenEmpty   = errors.New("auth token is empty")
	ErrWrongAuthToken   = errors.New("wrong auth token")
)

const (
	AuthorizationHeader string = "Authorization"
	AuthorizationPrefix string = "Bearer"
)

type AuthService interface {
	ParseAuthToken(token string) (*auth.JwtClaims, error)
}

// AuthMiddleware is a middleware for authentication.
func AuthMiddleware(notAuthUrls []string, authService AuthService) mux.MiddlewareFunc {
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
				httpTransport.EncodeError(r.Context(), http.StatusUnauthorized, ErrMissingAuthToken, w)

				return
			}

			// Check if token is valid
			tokenParts := strings.Split(tokenHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != AuthorizationPrefix {
				httpTransport.EncodeError(r.Context(), http.StatusUnauthorized, ErrInvalidAuthToken, w)

				return
			}

			// Check if token is empty
			if len(tokenParts[1]) == 0 {
				httpTransport.EncodeError(r.Context(), http.StatusUnauthorized, ErrAuthTokenEmpty, w)

				return
			}

			// Parse token
			token, err := authService.ParseAuthToken(tokenParts[1])
			if err != nil {
				httpTransport.EncodeError(r.Context(), http.StatusUnauthorized, ErrWrongAuthToken, w)

				return
			}

			// Add user id to context
			ctx := auth.SetUserIDToContext(r.Context(), token.UUID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
