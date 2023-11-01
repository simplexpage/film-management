package http

import (
	"context"
	"go.uber.org/zap"
)

// LogErrorHandler is a error handler that logs the error.
type LogErrorHandler struct {
	logger *zap.Logger
}

// Handle logs the error.
func (l LogErrorHandler) Handle(_ context.Context, err error) {
	l.logger.Error(err.Error())
}

// NewLogErrorHandler returns a new LogErrorHandler.
func NewLogErrorHandler(logger *zap.Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}
