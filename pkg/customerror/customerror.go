package customerror

import (
	"fmt"
	"net/http"
)

// CustomError represents a structured application error
type CustomError struct {
	Code       string // Error code (e.g., "NOT_FOUND", "VALIDATION_ERROR")
	Message    string // User-friendly error message
	StatusCode *int   // HTTP status code (nullable, defaults based on code if nil)
	Err        error  // Original error (nullable, for logging/debugging)
}

var _ error = (*CustomError)(nil)

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error (for errors.Is and errors.As)
func (e *CustomError) Unwrap() error {
	return e.Err
}

// GetStatusCode returns the HTTP status code, with smart defaults if nil
func (e *CustomError) GetStatusCode() int {
	if e.StatusCode != nil {
		return *e.StatusCode
	}

	return http.StatusInternalServerError
}

// WithError wraps an underlying error
func (e *CustomError) WithError(err error) *CustomError {
	return &CustomError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Err:        err,
	}
}

// NewCustomError creates a new CustomError
func NewCustomError(code, message string, statusCode int) *CustomError {
	status := statusCode
	return &CustomError{
		Code:       code,
		Message:    message,
		StatusCode: &status,
	}
}

// NewWithDefaults creates a new CustomError with automatic status code
func NewWithDefaults(code, message string) *CustomError {
	return &CustomError{
		Code:       code,
		Message:    message,
		StatusCode: nil, // Will use default based on code
	}
}

// Wrap wraps an error with CustomError
func Wrap(err error, code, message string, statusCode int) *CustomError {
	status := statusCode
	return &CustomError{
		Code:       code,
		Message:    message,
		StatusCode: &status,
		Err:        err,
	}
}
