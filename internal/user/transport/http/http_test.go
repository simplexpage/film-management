package http_test

import (
	"context"
	"film-management/config"
	"film-management/internal/user/domain"
	"film-management/internal/user/domain/mocks"
	"film-management/internal/user/endpoints"
	userHttp "film-management/internal/user/transport/http"
	customError "film-management/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	ConfigPath = "../../../../config"
)

// TestRegisterHandler tests the register user handler.
func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	type mockServiceBehavior func(r *mocks.MockService)

	cfg := config.GetConfig(ConfigPath)
	log := zap.NewNop()

	type request struct {
		body string
	}

	tests := []struct {
		name string
		request
		mockServiceBehavior mockServiceBehavior
		expectedStatusCode  int
		expectedResponse    string
	}{
		{
			name: "success",
			request: request{
				body: `{"username":"user1","password":"12345678"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {
				r.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"code":200,"data":{"uuid":"00000000-0000-0000-0000-000000000000","username":"user1"},"message":"OK"}`,
		},
		{
			name: "json decode error",
			request: request{
				body: `{"`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponse:    `{"code":400,"message":"json decode failed"}`,
		},
		{
			name: "validation error username required",
			request: request{
				body: `{"password":"12345678"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {},
			expectedStatusCode:  http.StatusUnprocessableEntity,
			expectedResponse:    `{"code":422,"message":"data validation error","data":{"username":"Username is required"}}`,
		},
		{
			name: "validation error password required",
			request: request{
				body: `{"username":"user1"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {},
			expectedStatusCode:  http.StatusUnprocessableEntity,
			expectedResponse:    `{"code":422,"message":"data validation error","data":{"password":"Password is required"}}`,
		},
		{
			name: "validation error user username exists",
			request: request{
				body: `{"username":"user1","password":"12345678"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {
				r.EXPECT().Register(gomock.Any(), gomock.Any()).Return(customError.ValidationError{Field: "username", Err: domain.ErrUserExistsWithUsername}).AnyTimes()
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   `{"code":422,"message":"data validation error","data":{"username":"user already exists with the same username"}}`,
		},
		{
			name: "failed to check username exists",
			request: request{
				body: `{"username":"user1","password":"12345678"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {
				r.EXPECT().Register(gomock.Any(), gomock.Any()).Return(domain.ErrUserCheckExistence).AnyTimes()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":500,"message":"failed to check user existence"}`,
		},
		{
			name: "failed to add user",
			request: request{
				body: `{"username":"user1","password":"12345678"}`,
			},
			mockServiceBehavior: func(r *mocks.MockService) {
				r.EXPECT().Register(gomock.Any(), gomock.Any()).Return(domain.ErrUserCreate).AnyTimes()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":500,"message":"failed to create user"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mocks.NewMockService(ctrl)
			test.mockServiceBehavior(serviceMock)

			serviceEndpoints := endpoints.NewEndpoints(serviceMock, log)
			serviceHTTPHandler := userHttp.NewHTTPHandlers(serviceEndpoints, cfg, log)

			srv := httptest.NewServer(serviceHTTPHandler)
			defer srv.Close()

			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, userHttp.RegisterPath, strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			serviceHTTPHandler.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponse, strings.TrimSpace(w.Body.String()))
		})
	}
}
