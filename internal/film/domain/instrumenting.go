package domain

import (
	"context"
	modelsFilm "film-management/internal/film/domain/models"
	"film-management/pkg/instrumenting"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
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

func (i instrumentingMiddleware) AddFilm(ctx context.Context, model *modelsFilm.Film) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddFilm", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.AddFilm(ctx, model)
}

func (i instrumentingMiddleware) UpdateFilm(ctx context.Context, model *modelsFilm.Film) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateFilm", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.UpdateFilm(ctx, model)
}

func (i instrumentingMiddleware) ViewFilm(ctx context.Context, filmID uuid.UUID) (model modelsFilm.Film, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ViewFilm", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.ViewFilm(ctx, filmID)
}

func (i instrumentingMiddleware) ViewAllFilms(ctx context.Context, filterSortPagination query.FilterSortLimit) (models []modelsFilm.Film, p pagination.Pagination, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ViewAllFilms", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.ViewAllFilms(ctx, filterSortPagination)
}

func (i instrumentingMiddleware) DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteFilm", "error", instrumenting.PrintErr(err)}
		i.requestCount.With(lvs...).Add(1)
		i.requestDuration.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.DeleteFilm(ctx, filmID, userID)
}
