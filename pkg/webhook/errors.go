package webhook

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

var (
	// Connection errors
	ErrConnectionFailed = customerror.New("WEBHOOK_CONNECTION_FAILED", "failed to connect to webhook", http.StatusServiceUnavailable)
	ErrTimeout          = customerror.New("WEBHOOK_TIMEOUT", "webhook request timed out", http.StatusGatewayTimeout)

	// Request errors
	ErrInvalidURL     = customerror.New("WEBHOOK_INVALID_URL", "invalid webhook URL", http.StatusInternalServerError)
	ErrInvalidRequest = customerror.New("WEBHOOK_INVALID_REQUEST", "invalid webhook request", http.StatusBadRequest)

	// Response errors
	ErrUnauthorized    = customerror.New("WEBHOOK_UNAUTHORIZED", "webhook authentication failed", http.StatusUnauthorized)
	ErrServerError     = customerror.New("WEBHOOK_SERVER_ERROR", "webhook server error", http.StatusBadGateway)
	ErrParsingResponse = customerror.New("WEBHOOK_PARSING_ERROR", "failed to parse webhook response", http.StatusInternalServerError)

	// Validation errors
	ErrInvalidPhoneNumber = customerror.New("INVALID_PHONE_NUMBER", "invalid phone number format", http.StatusBadRequest)
	ErrEmptyContent       = customerror.New("EMPTY_CONTENT", "message content cannot be empty", http.StatusBadRequest)
)
