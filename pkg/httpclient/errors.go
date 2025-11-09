package httpclient

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

var (
	// ErrRequestFailed is returned when HTTP request fails
	ErrRequestFailed = customerror.New(
		"HTTP_REQUEST_FAILED",
		"HTTP request failed",
		http.StatusBadGateway,
	)

	// ErrTimeout is returned when request times out
	ErrTimeout = customerror.New(
		"HTTP_TIMEOUT",
		"HTTP request timed out",
		http.StatusGatewayTimeout,
	)

	// ErrInvalidRequest is returned when request is invalid
	ErrInvalidRequest = customerror.New(
		"INVALID_HTTP_REQUEST",
		"Invalid HTTP request",
		http.StatusBadRequest,
	)

	// ErrUnexpectedStatus is returned when response status is unexpected
	ErrUnexpectedStatus = customerror.New(
		"UNEXPECTED_HTTP_STATUS",
		"Unexpected HTTP status code",
		http.StatusBadGateway,
	)
)

// Compile-time interface compliance check
var _ error = (*customerror.CustomError)(nil)
