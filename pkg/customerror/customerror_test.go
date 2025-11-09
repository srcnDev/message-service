package customerror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CustomError
		expected string
	}{
		{
			name: "error without underlying error",
			err: &CustomError{
				Code:    "NOT_FOUND",
				Message: "Resource not found",
			},
			expected: "[NOT_FOUND] Resource not found",
		},
		{
			name: "error with underlying error",
			err: &CustomError{
				Code:    "DB_ERROR",
				Message: "Database query failed",
				Err:     errors.New("connection timeout"),
			},
			expected: "[DB_ERROR] Database query failed: connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		err      *CustomError
		expected error
	}{
		{
			name: "unwrap with underlying error",
			err: &CustomError{
				Code:    "TEST",
				Message: "test message",
				Err:     fmt.Errorf("original error"),
			},
			expected: fmt.Errorf("original error"),
		},
		{
			name: "unwrap without underlying error",
			err: &CustomError{
				Code:    "TEST",
				Message: "test message",
				Err:     nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Unwrap()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Error(), result.Error())
			}
		})
	}
}

func TestCustomError_GetStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      *CustomError
		expected int
	}{
		{
			name: "with explicit status code",
			err: &CustomError{
				Code:       "NOT_FOUND",
				Message:    "Not found",
				StatusCode: intPtr(http.StatusNotFound),
			},
			expected: http.StatusNotFound,
		},
		{
			name: "without status code - default to 500",
			err: &CustomError{
				Code:       "UNKNOWN",
				Message:    "Unknown error",
				StatusCode: nil,
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.GetStatusCode()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomError_WithError(t *testing.T) {
	originalErr := errors.New("database timeout")
	baseErr := &CustomError{
		Code:       "DB_ERROR",
		Message:    "Database operation failed",
		StatusCode: intPtr(http.StatusInternalServerError),
	}

	wrappedErr := baseErr.WithError(originalErr)

	assert.Equal(t, baseErr.Code, wrappedErr.Code)
	assert.Equal(t, baseErr.Message, wrappedErr.Message)
	assert.Equal(t, baseErr.StatusCode, wrappedErr.StatusCode)
	assert.Equal(t, originalErr, wrappedErr.Err)
	assert.NotSame(t, baseErr, wrappedErr) // Should create new instance
}

func TestNew(t *testing.T) {
	code := "VALIDATION_ERROR"
	message := "Invalid input"
	statusCode := http.StatusBadRequest

	err := New(code, message, statusCode)

	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.NotNil(t, err.StatusCode)
	assert.Equal(t, statusCode, *err.StatusCode)
	assert.Nil(t, err.Err)
}

func TestNewWithDefaults(t *testing.T) {
	code := "GENERIC_ERROR"
	message := "Something went wrong"

	err := NewWithDefaults(code, message)

	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Nil(t, err.StatusCode)
	assert.Nil(t, err.Err)
	assert.Equal(t, http.StatusInternalServerError, err.GetStatusCode())
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("connection refused")
	code := "NETWORK_ERROR"
	message := "Failed to connect"
	statusCode := http.StatusServiceUnavailable

	wrappedErr := Wrap(originalErr, code, message, statusCode)

	assert.Equal(t, code, wrappedErr.Code)
	assert.Equal(t, message, wrappedErr.Message)
	assert.NotNil(t, wrappedErr.StatusCode)
	assert.Equal(t, statusCode, *wrappedErr.StatusCode)
	assert.Equal(t, originalErr, wrappedErr.Err)
}

func TestCustomError_ErrorsIsCompatibility(t *testing.T) {
	baseErr := errors.New("base error")
	customErr := &CustomError{
		Code:    "TEST",
		Message: "test",
		Err:     baseErr,
	}

	// Test errors.Is compatibility
	assert.True(t, errors.Is(customErr, baseErr))
}

func TestCustomError_InterfaceCompliance(t *testing.T) {
	var err error = &CustomError{
		Code:    "TEST",
		Message: "test message",
	}

	// Should compile and be usable as error interface
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "TEST")
}

// Helper function
func intPtr(i int) *int {
	return &i
}
