package webhook

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeWebhookConnectionFailed = "WEBHOOK_CONNECTION_FAILED"
	ErrCodeWebhookTimeout          = "WEBHOOK_TIMEOUT"
	ErrCodeWebhookInvalidURL       = "WEBHOOK_INVALID_URL"
	ErrCodeWebhookInvalidRequest   = "WEBHOOK_INVALID_REQUEST"
	ErrCodeWebhookUnauthorized     = "WEBHOOK_UNAUTHORIZED"
	ErrCodeWebhookServerError      = "WEBHOOK_SERVER_ERROR"
	ErrCodeWebhookParsingResponse  = "WEBHOOK_PARSING_ERROR"
	ErrCodeInvalidPhoneNumber      = "INVALID_PHONE_NUMBER"
	ErrCodeEmptyContent            = "EMPTY_CONTENT"
)

// Error messages
const (
	MsgWebhookConnectionFailed = "Failed to connect to webhook"
	MsgWebhookTimeout          = "Webhook request timed out"
	MsgWebhookInvalidURL       = "Invalid webhook URL"
	MsgWebhookInvalidRequest   = "Invalid webhook request"
	MsgWebhookUnauthorized     = "Webhook authentication failed"
	MsgWebhookServerError      = "Webhook server error"
	MsgWebhookParsingResponse  = "Failed to parse webhook response"
	MsgInvalidPhoneNumber      = "Invalid phone number format"
	MsgEmptyContent            = "Message content cannot be empty"
)

// Predefined errors
var (
	ErrConnectionFailed = customerror.NewCustomError(
		ErrCodeWebhookConnectionFailed,
		MsgWebhookConnectionFailed,
		http.StatusServiceUnavailable,
	)

	ErrTimeout = customerror.NewCustomError(
		ErrCodeWebhookTimeout,
		MsgWebhookTimeout,
		http.StatusGatewayTimeout,
	)

	ErrInvalidURL = customerror.NewCustomError(
		ErrCodeWebhookInvalidURL,
		MsgWebhookInvalidURL,
		http.StatusInternalServerError,
	)

	ErrInvalidRequest = customerror.NewCustomError(
		ErrCodeWebhookInvalidRequest,
		MsgWebhookInvalidRequest,
		http.StatusBadRequest,
	)

	ErrUnauthorized = customerror.NewCustomError(
		ErrCodeWebhookUnauthorized,
		MsgWebhookUnauthorized,
		http.StatusUnauthorized,
	)

	ErrServerError = customerror.NewCustomError(
		ErrCodeWebhookServerError,
		MsgWebhookServerError,
		http.StatusBadGateway,
	)

	ErrParsingResponse = customerror.NewCustomError(
		ErrCodeWebhookParsingResponse,
		MsgWebhookParsingResponse,
		http.StatusInternalServerError,
	)

	ErrInvalidPhoneNumber = customerror.NewCustomError(
		ErrCodeInvalidPhoneNumber,
		MsgInvalidPhoneNumber,
		http.StatusBadRequest,
	)

	ErrEmptyContent = customerror.NewCustomError(
		ErrCodeEmptyContent,
		MsgEmptyContent,
		http.StatusBadRequest,
	)
)
