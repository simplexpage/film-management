package response

import (
	"context"
	"errors"
	transportHttp "film-management/pkg/transport/http"
	"film-management/pkg/validation"
	endpointKit "github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"net/http"
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
func EncodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// Check if response is Failer interface and encode error
	if f, ok := response.(endpointKit.Failer); ok && f.Failed() != nil {
		EncodeError(ctx, errorToCodeHTTPAnswer(f.Failed()), f.Failed(), w)

		return
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
func EncodeError(ctx context.Context, code int, err error, w http.ResponseWriter) {
	// Check if error is Validation Errors
	validationErr := errorsValidationMap(err)
	if len(validationErr) > 0 {
		encodeValidation(ctx, validationErr, w)

		return
	}

	// Check if error is NotFoundError
	if errors.As(err, &validation.NotFoundError{}) {
		code = http.StatusNotFound
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	// Error response
	data := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	if errEncode := jsoniter.NewEncoder(w).Encode(data); errEncode != nil {
		return
	}
}

// errorsValidationMap convert Validation Errors to map.
func errorsValidationMap(err error) map[string]string {
	result := make(map[string]string)

	var (
		validationErrors validator.ValidationErrors
		customError      validation.CustomError
	)

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			result[fieldError.Field()] = fieldError.Translate(validation.GetTranslator())
		}
	} else if errors.As(err, &customError) {
		result[customError.Field] = customError.Err.Error()
	}

	return result
}

// encodeValidation is the default error handler. It encodes errors to the HTTP response.
func encodeValidation(_ context.Context, validationErr map[string]string, w http.ResponseWriter) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)

	// Error response validation
	data := ErrorResponseValidation{
		Code:    http.StatusUnprocessableEntity,
		Message: transportHttp.ErrDataValidation.Error(),
		Data:    validationErr,
	}

	errEncode := jsoniter.NewEncoder(w).Encode(data)

	if errEncode != nil {
		return
	}
}

// EncodeServerError is the default error handler. It encodes errors to the HTTP response.
func EncodeServerError(_ context.Context, err error, w http.ResponseWriter) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorToCodeHTTPAnswer(err))

	// Error response
	data := ErrorResponse{
		Code:    errorToCodeHTTPAnswer(err),
		Message: err.Error(),
	}

	if errEncode := jsoniter.NewEncoder(w).Encode(data); errEncode != nil {
		return
	}
}

// errorToCodeHTTPAnswer convert error to HTTP code.
func errorToCodeHTTPAnswer(err error) int {
	switch {
	case errors.Is(err, transportHttp.ErrBadRouting) ||
		errors.Is(err, transportHttp.ErrContextUserID) ||
		errors.Is(err, transportHttp.ErrJSONDecode):
		return http.StatusBadRequest
	case errors.Is(err, transportHttp.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
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
