package cors_test

import (
	"film-management/pkg/transport/http/middlewares/cors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                string
		method              string
		origin              string
		expectedStatus      int
		expectedErrorString string
	}

	logger := zap.NewNop()
	corsAllowedOrigins := []string{"http://example.com", "http://test.com"}

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []testCase{
		{
			name:           "ValidOrigin",
			method:         http.MethodGet,
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:                "InvalidOrigin",
			method:              http.MethodGet,
			origin:              "http://invalid.com",
			expectedStatus:      http.StatusForbidden,
			expectedErrorString: cors.ErrInvalidOriginCORS.Error(),
		},
		{
			name:           "EmptyOrigin",
			method:         http.MethodGet,
			origin:         "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OptionsMethod",
			method:         http.MethodOptions,
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
		},
	}

	handler := cors.Middleware(corsAllowedOrigins, logger)(corsHandler)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(tc.method, "/", nil)
			request.Header.Set(cors.HeaderOrigin, tc.origin)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedErrorString != "" {
				// Check if response body contains the expected error string
				assert.Contains(t, recorder.Body.String(), tc.expectedErrorString)
			} else {
				if tc.method == http.MethodOptions {
					// Check if the response has status OK for OPTIONS method
					assert.Equal(t, http.StatusOK, recorder.Code)
				} else {
					// Check if the Access-Control-Allow-Origin header is set correctly
					assert.Equal(t, tc.origin, recorder.Header().Get("Access-Control-Allow-Origin"))
				}
			}
		})
	}
}
