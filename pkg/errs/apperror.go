package errs

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code       string // Error code (e.g., "NOT_FOUND", "VALIDATION_ERROR")
	Message    string // User-friendly error message
	StatusCode *int   // HTTP status code (nullable, defaults based on code if nil)
	Err        error  // Original error (nullable, for logging/debugging)
}

var _ error = (*AppError)(nil)

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error (for errors.Is and errors.As)
func (e *AppError) Unwrap() error {
	return e.Err
}

// GetStatusCode returns the HTTP status code, with smart defaults if nil
func (e *AppError) GetStatusCode() int {
	if e.StatusCode != nil {
		return *e.StatusCode
	}

	return http.StatusInternalServerError
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Err:        err,
	}
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	status := statusCode
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: &status,
	}
}

// NewWithDefaults creates a new AppError with automatic status code
func NewWithDefaults(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: nil, // Will use default based on code
	}
}

// Wrap wraps an error with AppError
func Wrap(err error, code, message string, statusCode int) *AppError {
	status := statusCode
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: &status,
		Err:        err,
	}
}
