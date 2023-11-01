package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
	"time"
)

func NewLoggingMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Debug("endpoint", zap.Error(err), zap.Duration("took", time.Since(begin)))
			}(time.Now())

			return next(ctx, request)
		}
	}
}
