package httpclient

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeHTTPRequestFailed    = "HTTP_REQUEST_FAILED"
	ErrCodeHTTPTimeout          = "HTTP_TIMEOUT"
	ErrCodeInvalidHTTPRequest   = "INVALID_HTTP_REQUEST"
	ErrCodeUnexpectedHTTPStatus = "UNEXPECTED_HTTP_STATUS"
)

// Error messages
const (
	MsgHTTPRequestFailed    = "HTTP request failed"
	MsgHTTPTimeout          = "HTTP request timed out"
	MsgInvalidHTTPRequest   = "Invalid HTTP request"
	MsgUnexpectedHTTPStatus = "Unexpected HTTP status code"
)

// Predefined errors
var (
	ErrRequestFailed = customerror.NewCustomError(
		ErrCodeHTTPRequestFailed,
		MsgHTTPRequestFailed,
		http.StatusBadGateway,
	)

	ErrTimeout = customerror.NewCustomError(
		ErrCodeHTTPTimeout,
		MsgHTTPTimeout,
		http.StatusGatewayTimeout,
	)

	ErrInvalidRequest = customerror.NewCustomError(
		ErrCodeInvalidHTTPRequest,
		MsgInvalidHTTPRequest,
		http.StatusBadRequest,
	)

	ErrUnexpectedStatus = customerror.NewCustomError(
		ErrCodeUnexpectedHTTPStatus,
		MsgUnexpectedHTTPStatus,
		http.StatusBadGateway,
	)
)
