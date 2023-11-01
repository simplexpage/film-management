package middlewares_test

import (
	"film-management/config"
	"film-management/pkg/transport/http/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	ConfigPath = "../../../config/test"
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
		logger            = zap.NewNop()
		notAuthURLs       = []string{"/public", "/login"}
		pathPublicKeyFile = cfg.HTTP.PathPublicKeyFile
		authTokenForTest  = cfg.HTTP.AuthTokenForTest
	)

	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []testCase{
		{
			name:           "ValidToken",
			url:            "/protected",
			authToken:      middlewares.AuthorizationPrefix + " " + authTokenForTest,
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
			authToken:      authTokenForTest,
			expectedStatus: http.StatusOK,
		},
	}

	handler := middlewares.AuthMiddleware(notAuthURLs, pathPublicKeyFile, logger)(authHandler)

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
