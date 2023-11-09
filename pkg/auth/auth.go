package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrValidAuthToken            = errors.New("token is not valid")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrAuthTokenNoOfTypeJwt      = errors.New("token claims are not of type JwtClaims")
	ErrCheckTokeHeaderSigningKey = errors.New("error checking token header signing key")
	ErrGetPublicKeyFromFile      = errors.New("error getting public key from file")
	ErrGetPrivateKeyFile         = errors.New("error getting private key from file")
	ErrParseRSAPrivateKeyFromPEM = errors.New("error parsing RSA private key from PEM")
	ErrNewClaims                 = errors.New("error creating new claims")
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

// Service is a struct for Auth.
type Service struct {
	logger *zap.Logger
	cfg    Config
}

// NewAuthService is a constructor for Service.
func NewAuthService(cfg Config, logger *zap.Logger) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger,
	}
}

// GenerateAuthToken is a function for generating auth token.
func (a Service) GenerateAuthToken(uuid string) (string, time.Time, error) {
	privateKeyBytes, err := getPrivateKeyFile(a.cfg.PathPrivateKeyFile)
	if err != nil {
		a.logger.Error("during getPrivateKeyFile", zap.Error(err))

		return "", time.Time{}, errors.Wrap(ErrGetPrivateKeyFile, "authService.GenerateAuthToken.getPrivateKeyFile")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		a.logger.Error("during jwt.ParseRSAPrivateKeyFromPEM", zap.Error(err))

		return "", time.Time{}, errors.Wrap(ErrParseRSAPrivateKeyFromPEM, "authService.GenerateAuthToken.ParseRSAPrivateKeyFromPEM")
	}

	expirationTime := time.Now().Add(time.Duration(a.cfg.AuthDurationMin) * time.Minute)

	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expirationTime),
			IssuedAt:  jwt.At(time.Now()),
		},
		UUID: uuid,
	}

	a.logger.Debug("claims", zap.Any("claims", claims))

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		a.logger.Error("during jwt.NewWithClaims", zap.Error(err))

		return "", time.Time{}, errors.Wrap(ErrNewClaims, "authService.GenerateAuthToken.NewWithClaims")
	}

	return token, expirationTime, nil
}

// ParseAuthToken is a function for parsing token.
func (a Service) ParseAuthToken(accessToken string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			method, methodOk := token.Header["alg"].(string)
			if !methodOk {
				a.logger.Error("check if token.Header[\"alg\"] is string", zap.String("token.Header[\"alg\"]", fmt.Sprintf("%v", token.Header["alg"])))

				return nil, errors.Wrap(ErrCheckTokeHeaderSigningKey, "authService.ParseAuthToken.checkIfTokenHeaderAlgIsString")
			}
			a.logger.Error("unexpected signing method", zap.String("method", method))

			return nil, errors.Wrap(ErrUnexpectedSigningMethod, "authService.ParseAuthToken.unexpectedSigningMethod")
		}

		publicKey, err := getPublicKeyFromFile(a.cfg.PathPublicKeyFile)
		if err != nil {
			a.logger.Error("during getPublicKeyFromFile", zap.Error(err))

			return nil, errors.Wrap(ErrGetPublicKeyFromFile, "authService.ParseAuthToken.getPublicKeyFromFile")
		}

		return publicKey, nil
	})

	if err != nil {
		a.logger.Debug("during jwt.ParseWithClaims", zap.Error(err))

		return nil, errors.Wrap(ErrValidAuthToken, "authService.ParseAuthToken.ParseWithClaims")
	}

	if !token.Valid {
		return nil, errors.Wrap(ErrValidAuthToken, "authService.ParseAuthToken.Valid")
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		a.logger.Error("token claims are not of type JwtClaims")

		return nil, errors.Wrap(ErrAuthTokenNoOfTypeJwt, "authService.ParseAuthToken.tokenClaimsAreNotOfTypeJwt")
	}

	return claims, nil
}

// getPrivateKeyFile is a function for getting private key from file.
func getPrivateKeyFile(pathPrivateKeyFile string) ([]byte, error) {
	path, err := filepath.Abs(pathPrivateKeyFile)

	if err != nil {
		return nil, errors.Wrap(err, "authService.getPrivateKeyFile.filepath.Abs")
	}

	verifyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "authService.getPrivateKeyFile.os.ReadFile")
	}

	return verifyBytes, nil
}

// getPublicKeyFromFile is a function for getting public key from file.
func getPublicKeyFromFile(pathPublicKeyFile string) (*rsa.PublicKey, error) {
	path, err := filepath.Abs(pathPublicKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "authService.getPublicKeyFromFile.filepath.Abs")
	}

	verifyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "authService.getPublicKeyFromFile.os.ReadFile")
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "authService.getPublicKeyFromFile.jwt.ParseRSAPublicKeyFromPEM")
	}

	return verifyKey, nil
}
