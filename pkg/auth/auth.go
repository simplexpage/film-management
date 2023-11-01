package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

type ContextKey string

const (
	ContextKeyUserID ContextKey = "user_uuid"
)

var (
	ErrContextUserID           = errors.New("user uuid not found in context")
	ErrValidAuthToken          = errors.New("token is not valid")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrAuthTokenNoOfTypeJwt    = errors.New("token claims are not of type JwtClaims")
)

// JwtClaims is a struct for JWT auth.
type JwtClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

func GetUserUUIDFromContext(ctx context.Context) (string, error) {
	userUUID, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok {
		return "", ErrContextUserID
	}

	return userUUID, nil
}

// SetUserUUIDToContext is a function for setting user UUID to context.
func SetUserUUIDToContext(ctx context.Context, UUID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, UUID)
}

// GenerateAuthToken is a function for generating auth token.
func GenerateAuthToken(UUID string, pathPrivateKeyFile string, authDurationMin time.Duration) (string, error) {
	privateKeyBytes, err := getPrivateKeyFile(pathPrivateKeyFile)
	if err != nil {
		return "", fmt.Errorf("error getting key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("parse key: %w", err)
	}

	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(authDurationMin)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UUID: UUID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("create sign token: %w", err)
	}

	return token, nil
}

// getPrivateKeyFile is a function for getting private key from file.
func getPrivateKeyFile(pathPrivateKeyFile string) ([]byte, error) {
	path, err := filepath.Abs(pathPrivateKeyFile)

	if err != nil {
		return nil, fmt.Errorf("during filepath.Abs: %w", err)
	}

	verifyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("during jwt.ParseRSAPrivateKeyFromPEM: %w", err)
	}

	return verifyBytes, nil
}

// getPublicKeyFromFile is a function for getting public key from file.
func getPublicKeyFromFile(pathPublicKeyFile string) (*rsa.PublicKey, error) {
	path, err := filepath.Abs(pathPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("during filepath.Abs: %w", err)
	}

	verifyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("during os.ReadFile: %w", err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, fmt.Errorf("during jwt.ParseRSAPublicKeyFromPEM: %w", err)
	}

	return verifyKey, nil
}

// ParseToken is a function for parsing token.
func ParseToken(accessToken, pathPublicKeyFile string, logger *zap.Logger) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			method, methodOk := token.Header["alg"].(string)
			if !methodOk {
				logger.Error("unexpected signing method", zap.String("method", "unknown"))

				return nil, ErrUnexpectedSigningMethod
			}
			logger.Error("unexpected signing method", zap.String("method", method))

			return nil, ErrUnexpectedSigningMethod
		}

		publicKey, err := getPublicKeyFromFile(pathPublicKeyFile)
		if err != nil {
			logger.Error("during getPublicKeyFromFile", zap.Error(err))

			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrValidAuthToken
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		logger.Error("token claims are not of type JwtClaims")

		return nil, ErrAuthTokenNoOfTypeJwt
	}

	return claims, nil
}
