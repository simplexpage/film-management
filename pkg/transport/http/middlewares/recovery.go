package middlewares

import (
	"errors"
	httpTransport "film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrInvalidServer = errors.New("invalid server error")
)

// RecoveryMiddleware is a middleware for recovering from panic.
func RecoveryMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic occurred", zap.Any("error", err))
					httpTransport.EncodeError(r.Context(), http.StatusBadRequest, ErrInvalidServer, w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
