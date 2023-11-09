package recovery

import (
	httpTransport "film-management/pkg/transport/http/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrInvalidServer = errors.New("invalid server error")
)

// Middleware is a middleware for recovering from panic.
func Middleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic occurred", zap.Any("error", err))
					httpTransport.EncodeError(r.Context(), errors.Wrap(ErrInvalidServer, "middlewares.RecoveryMiddleware"), w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
