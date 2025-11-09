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
	ErrMessageSendFailed = customerror.New(
		ErrCodeMessageSendFailed,
		MsgMessageSendFailed,
		http.StatusInternalServerError,
	)

	ErrWebhookCallFailed = customerror.New(
		ErrCodeWebhookCallFailed,
		MsgWebhookCallFailed,
		http.StatusBadGateway,
	)

	ErrMarkSentFailed = customerror.New(
		ErrCodeMarkSentFailed,
		MsgMarkSentFailed,
		http.StatusInternalServerError,
	)

	ErrMarkFailedFailed = customerror.New(
		ErrCodeMarkFailedFailed,
		MsgMarkFailedFailed,
		http.StatusInternalServerError,
	)
)
