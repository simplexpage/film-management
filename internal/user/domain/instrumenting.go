package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"film-management/pkg/instrumenting"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingMiddleware struct {
	requestCount    metrics.Counter
	requestDuration metrics.Histogram
	next            Service
}

// NewInstrumentingMiddleware returns an instance of the instrumenting middleware.
func NewInstrumentingMiddleware(requestCount metrics.Counter,
	requestDuration metrics.Histogram) Middleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{
			requestCount,
			requestDuration,
			next,
		}
	}
}

func (i instrumentingMiddleware) Register(ctx context.Context, model *models.User) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Register(ctx, model)
}

func (i instrumentingMiddleware) Login(ctx context.Context, username string, password string) (authToken string, expirationTime time.Time, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Login(ctx, username, password)
}
