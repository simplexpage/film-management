package domain

import (
	"context"
	"film-management/internal/user/domain/models"
	"go.uber.org/zap"
	"time"
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

func (l loggingMiddleware) Register(ctx context.Context, model *models.User) (err error) {
	defer func() {
		l.logger.With(zap.String("method", "Register")).
			Debug("domain",
				zap.String("username", model.Username),
				zap.String("password", model.Password),
				zap.Error(err))
	}()

	return l.next.Register(ctx, model)
}

func (l loggingMiddleware) Login(ctx context.Context, username string, password string) (authToken string, expirationTime time.Time, err error) {
	defer func() {
		l.logger.With(zap.String("method", "Login")).
			Debug("domain",
				zap.String("username", username),
				zap.String("password", password),
				zap.String("authToken", authToken),
				zap.Time("expirationTime", expirationTime),
				zap.Error(err))
	}()

	return l.next.Login(ctx, username, password)
}
