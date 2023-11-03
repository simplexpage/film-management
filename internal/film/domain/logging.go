package domain

import (
	"context"
	modelsFilm "film-management/internal/film/domain/models"
	"film-management/pkg/query"
	"film-management/pkg/query/pagination"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	next   Service
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (l loggingMiddleware) AddFilm(ctx context.Context, model *modelsFilm.Film) (err error) {
	defer func() {
		l.logger.With(zap.String("method", "AddFilm")).
			Debug("domain",
				zap.Any("film", model),
				zap.Error(err))
	}()

	return l.next.AddFilm(ctx, model)
}

func (l loggingMiddleware) UpdateFilm(ctx context.Context, model *modelsFilm.Film) (err error) {
	defer func() {
		l.logger.With(zap.String("method", "UpdateFilm")).
			Debug("domain",
				zap.Any("film", model),
				zap.Error(err))
	}()

	return l.next.UpdateFilm(ctx, model)
}

func (l loggingMiddleware) ViewFilm(ctx context.Context, filmID uuid.UUID) (model modelsFilm.Film, err error) {
	defer func() {
		l.logger.With(zap.String("method", "ViewFilm")).
			Debug("domain",
				zap.Any("film", model),
				zap.Error(err))
	}()

	return l.next.ViewFilm(ctx, filmID)
}

func (l loggingMiddleware) ViewAllFilms(ctx context.Context, filterSortLimit query.FilterSortLimit) (models []modelsFilm.Film, p pagination.Pagination, err error) {
	defer func() {
		l.logger.With(zap.String("method", "ViewAllFilms")).
			Debug("domain",
				zap.String("sort_field", filterSortLimit.Sort.Field()),
				zap.String("sort_order", filterSortLimit.Sort.Order()),
				zap.Int("limit", filterSortLimit.Limit),
				zap.Int("offset", filterSortLimit.Offset),
				zap.Int("page", p.Page),
				zap.Int("page-size", p.PageSize),
				zap.Int("total-count", p.TotalCount),
				zap.Error(err))
	}()

	return l.next.ViewAllFilms(ctx, filterSortLimit)
}

func (l loggingMiddleware) DeleteFilm(ctx context.Context, filmID uuid.UUID, userID uuid.UUID) (err error) {
	defer func() {
		l.logger.With(zap.String("method", "DeleteFilm")).
			Debug("domain",
				zap.Any("filmID", filmID),
				zap.Any("userID", userID),
				zap.Error(err))
	}()

	return l.next.DeleteFilm(ctx, filmID, userID)
}
