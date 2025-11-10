package apperror

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes for message sender
const (
	ErrCodeMessageSendFailed = "MESSAGE_SEND_FAILED"
	ErrCodeWebhookCallFailed = "WEBHOOK_CALL_FAILED"
	ErrCodeMarkSentFailed    = "MARK_SENT_FAILED"
	ErrCodeMarkFailedFailed  = "MARK_FAILED_FAILED"
)

// Error messages
const (
	MsgMessageSendFailed = "Failed to send message"
	MsgWebhookCallFailed = "Webhook call failed"
	MsgMarkSentFailed    = "Failed to mark message as sent"
	MsgMarkFailedFailed  = "Failed to mark message as failed"
)

// Predefined errors
var (
	ErrMessageSendFailed = customerror.NewCustomError(
		ErrCodeMessageSendFailed,
		MsgMessageSendFailed,
		http.StatusInternalServerError,
	)

	ErrWebhookCallFailed = customerror.NewCustomError(
		ErrCodeWebhookCallFailed,
		MsgWebhookCallFailed,
		http.StatusBadGateway,
	)

	ErrMarkSentFailed = customerror.NewCustomError(
		ErrCodeMarkSentFailed,
		MsgMarkSentFailed,
		http.StatusInternalServerError,
	)

	ErrMarkFailedFailed = customerror.NewCustomError(
		ErrCodeMarkFailedFailed,
		MsgMarkFailedFailed,
		http.StatusInternalServerError,
	)
)
