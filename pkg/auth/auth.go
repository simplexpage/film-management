package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
)

type ContextKey string

const (
	ContextKeyUserID ContextKey = "user_id"
)

var (
	ErrContextUserID             = errors.New("user uuid not found in context")
	ErrValidAuthToken            = errors.New("token is not valid")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrAuthTokenNoOfTypeJwt      = errors.New("token claims are not of type JwtClaims")
	ErrGetPublicKeyFromFile      = errors.New("error getting public key")
	ErrParseWithClaims           = errors.New("error parsing with claims")
	ErrGetPrivateKeyFile         = errors.New("error getting private key")
	ErrParseRSAPrivateKeyFromPEM = errors.New("error parsing rsa private key from pem")
	ErrCreateSignToken           = errors.New("error creating sign token")
)

// Config is a struct for auth config.
type Config struct {
	AuthDurationMin    int64
	PathPublicKeyFile  string
	PathPrivateKeyFile string
}

// JwtClaims is a struct for JWT auth.
type JwtClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

// GetUserIDFromContext is a function for getting user ID from context.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	// Get user ID from context
	userID, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok {
		return "", ErrContextUserID
	}

	return userID, nil
}

// SetUserIDToContext is a function for setting user ID to context.
func SetUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}
