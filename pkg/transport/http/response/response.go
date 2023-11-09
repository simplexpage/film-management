package response

import (
	"context"
	customError "film-management/pkg/errors"
	transportHttp "film-management/pkg/transport/http"
	"film-management/pkg/validation"
	endpointKit "github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// SuccessResponse is the common struct for all success responses.
type SuccessResponse struct {
	Code    int         `json:"code" example:"200"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"OK"`
}

// ErrorResponse is the common struct for all error responses.
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

// ErrorResponseValidation is the common struct for all error responses with validation errors.
type ErrorResponseValidation struct {
	Code    int               `json:"code" example:"421"`
	Message string            `json:"message" example:"Validation Error"`
	Data    map[string]string `json:"data"`
}

// EncodeHTTPResponse is the common method to encode all success responses to the HTTP response.
func EncodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// Check if response is Failer interface and encode error
	if f, ok := response.(endpointKit.Failer); ok && f.Failed() != nil {
		EncodeError(ctx, f.Failed(), w)

		return nil
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Success response
	data := SuccessResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    response,
	}

	return jsoniter.NewEncoder(w).Encode(data)
}

// EncodeError is the default error handler. It encodes errors to the HTTP response.
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Error response
	data, code := createErrorResponse(err)
	w.WriteHeader(code)

	if errEncode := jsoniter.NewEncoder(w).Encode(data); errEncode != nil {
		return
	}
}

// createErrorResponse is the common method to create all error responses.
func createErrorResponse(err error) (data interface{}, code int) {
	// Check if error is Validation Errors and convert to map
	validationErr := errorsValidationMap(err)

	switch {
	case len(validationErr) > 0:
		data, code = handleValidationErrors(validationErr)
	case errors.Is(err, transportHttp.ErrBadRouting),
		errors.Is(err, transportHttp.ErrJSONDecode):
		data, code = handleBadRequestErrors(err)
	case errors.Is(err, transportHttp.ErrNotFound),
		errors.As(err, &customError.NotFoundError{}):
		data, code = handleNotFoundError(err)
	case errors.As(err, &customError.AuthError{}):
		data, code = handleAuthErrors(err)
	case errors.As(err, &customError.CorsError{}) ||
		errors.As(err, &customError.PermissionError{}):
		data, code = handlePermissionErrors(err)
	default:
		data, code = handleDefaultErrors(err)
	}

	return
}

// handleAuthErrors is the common method to handle all auth errors.
func handleAuthErrors(err error) (data interface{}, code int) {
	data = ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	}
	code = http.StatusUnauthorized

	return
}

// handlePermissionErrors is the common method to handle all permission errors.
func handlePermissionErrors(err error) (data interface{}, code int) {
	data = ErrorResponse{
		Code:    http.StatusForbidden,
		Message: err.Error(),
	}
	code = http.StatusForbidden

	return
}

// handleValidationErrors is the common method to handle all validation errors.
func handleValidationErrors(validationErr map[string]string) (data interface{}, code int) {
	data = ErrorResponseValidation{
		Code:    http.StatusUnprocessableEntity,
		Message: transportHttp.ErrDataValidation.Error(),
		Data:    validationErr,
	}
	code = http.StatusUnprocessableEntity

	return
}

// handleBadRequestErrors is the common method to handle all bad request errors.
func handleBadRequestErrors(err error) (data interface{}, code int) {
	data = ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
	code = http.StatusBadRequest

	return
}

// handleNotFoundError is the common method to handle all not found errors.
func handleNotFoundError(err error) (data interface{}, code int) {
	data = ErrorResponse{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
	code = http.StatusNotFound

	return
}

// handleDefaultErrors is the common method to handle all default errors.
func handleDefaultErrors(err error) (data interface{}, code int) {
	data = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	code = http.StatusInternalServerError

	return
}

// errorsValidationMap convert Validation Errors to map.
func errorsValidationMap(err error) map[string]string {
	result := make(map[string]string)

	var (
		validationErrors validator.ValidationErrors
		validationError  customError.ValidationError
	)

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			result[strings.ToLower(fieldError.Field())] = fieldError.Translate(validation.GetTranslator())
		}
	} else if errors.As(err, &validationError) {
		result[validationError.Field] = validationError.Err.Error()
	}

	return result
}

// SetErrorHandlers is the default error handler. It encodes errors to the HTTP response.
func SetErrorHandlers(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(NotFoundFunc)
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedFunc)
}

// NotFoundFunc is the default error handler. It encodes errors to the HTTP response.
func NotFoundFunc(w http.ResponseWriter, _ *http.Request) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	// Error response
	data := ErrorResponse{
		Code:    http.StatusNotFound,
		Message: transportHttp.ErrSystemActionNotFound.Error(),
	}

	if errEncode := jsoniter.NewEncoder(w).Encode(data); errEncode != nil {
		return
	}
}

// MethodNotAllowedFunc is the default error handler. It encodes errors to the HTTP response.
func MethodNotAllowedFunc(w http.ResponseWriter, _ *http.Request) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)

	// Error response
	data := ErrorResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: transportHttp.ErrSystemActionMethodNotAllowed.Error(),
	}

	if errEncode := jsoniter.NewEncoder(w).Encode(data); errEncode != nil {
		return
	}
}
