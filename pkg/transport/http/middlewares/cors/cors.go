package cors

import (
	customError "film-management/pkg/errors"
	httpTransport "film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

const (
	HeaderOrigin = "Origin"
)

var (
	ErrInvalidOriginCORS = errors.New("forbidden - Origin is not allowed")
)

// Middleware is a middleware for CORS.
func Middleware(corsAllowedOrigins []string, logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			origin := r.Header.Get(HeaderOrigin)
			if origin != "" {
				logger.Debug("origin", zap.String("origin", origin))
				if !contains(corsAllowedOrigins, origin) {
					httpTransport.EncodeError(r.Context(), customError.CorsError{Err: ErrInvalidOriginCORS}, w)

					return
				}

				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
