package services

import (
	"crypto/rsa"
	"film-management/pkg/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

// AuthService is a struct for Auth.
type AuthService struct {
	logger *zap.Logger
	cfg    auth.Config
}

// NewAuthService is a constructor for Service.
func NewAuthService(cfg auth.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		cfg:    cfg,
		logger: logger,
	}
}

// GenerateAuthToken is a function for generating auth token.
func (a AuthService) GenerateAuthToken(UUID string) (string, time.Time, error) {
	privateKeyBytes, err := getPrivateKeyFile(a.cfg.PathPrivateKeyFile)
	if err != nil {
		a.logger.Debug("during getPrivateKeyFile", zap.Error(err))

		return "", time.Time{}, auth.ErrGetPrivateKeyFile
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		a.logger.Debug("during jwt.ParseRSAPrivateKeyFromPEM", zap.Error(err))

		return "", time.Time{}, auth.ErrParseRSAPrivateKeyFromPEM
	}

	expirationTime := time.Now().Add(time.Duration(a.cfg.AuthDurationMin) * time.Minute)

	claims := &auth.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expirationTime),
			IssuedAt:  jwt.At(time.Now()),
		},
		UUID: UUID,
	}

	a.logger.Debug("claims", zap.Any("claims", claims))

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		a.logger.Debug("during jwt.NewWithClaims", zap.Error(err))

		return "", time.Time{}, auth.ErrCreateSignToken
	}

	return token, expirationTime, nil
}

// ParseAuthToken is a function for parsing token.
func (a AuthService) ParseAuthToken(accessToken string) (*auth.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			method, methodOk := token.Header["alg"].(string)
			if !methodOk {
				a.logger.Debug("unexpected signing method", zap.String("method", "unknown"))

				return nil, auth.ErrUnexpectedSigningMethod
			}
			a.logger.Debug("unexpected signing method", zap.String("method", method))

			return nil, auth.ErrUnexpectedSigningMethod
		}

		publicKey, err := getPublicKeyFromFile(a.cfg.PathPublicKeyFile)
		if err != nil {
			a.logger.Debug("during getPublicKeyFromFile", zap.Error(err))

			return nil, auth.ErrGetPublicKeyFromFile
		}

		return publicKey, nil
	})

	if err != nil {
		a.logger.Debug("during jwt.ParseWithClaims", zap.Error(err))

		return nil, auth.ErrParseWithClaims
	}

	if !token.Valid {
		return nil, auth.ErrValidAuthToken
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		a.logger.Debug("token claims are not of type JwtClaims")

		return nil, auth.ErrAuthTokenNoOfTypeJwt
	}

	return claims, nil
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
