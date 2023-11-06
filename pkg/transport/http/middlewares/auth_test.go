package middlewares_test

import (
	"film-management/config"
	"film-management/pkg/transport/http/middlewares"
	"film-management/repositories/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	ConfigPath = "../../../../config"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		url            string
		authToken      string
		expectedStatus int
	}

	cfg := config.GetConfig(ConfigPath)

	var (
		logger      = zap.NewNop()
		notAuthURLs = []string{"/public", "/login"}
		authService = services.NewAuthService(cfg.Services.Auth, logger)
	)

	token, _, err := authService.GenerateAuthToken("d83d97ab-ff68-4de2-b2a9-7cd5f0fc9a5e")
	if err != nil {
		return
	}

	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middlewares.AuthMiddleware(notAuthURLs, authService)(authHandler)

	testCases := []testCase{
		{
			name:           "ValidToken",
			url:            "/protected",
			authToken:      middlewares.AuthorizationPrefix + " " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "MissingToken",
			url:            "/protected",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InvalidToken",
			url:            "/protected",
			authToken:      "invalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "EmptyToken",
			url:            "/protected",
			authToken:      "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "NotAuthURL",
			url:            "/public",
			authToken:      "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(http.MethodGet, tc.url, nil)
			request.Header.Set("Authorization", tc.authToken)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
		})
	}
}
