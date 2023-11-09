package recovery_test

import (
	"film-management/pkg/transport/http/middlewares/recovery"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler := recovery.Middleware(logger)(panicHandler)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), recovery.ErrInvalidServer.Error())
}
