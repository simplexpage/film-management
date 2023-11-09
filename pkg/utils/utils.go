package utils

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrGetValueFromContext = errors.New("failed to get value from context")
)

// GetValueFromContext is a function for getting value from context.
func GetValueFromContext(ctx context.Context, value interface{}) (string, error) {
	// Get user ID from context
	userID, ok := ctx.Value(value).(string)
	if !ok {
		return "", ErrGetValueFromContext
	}

	return userID, nil
}
