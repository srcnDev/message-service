package apperror

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeMessageNotFound     = "MESSAGE_NOT_FOUND"
	ErrCodeMessageCreateFailed = "MESSAGE_CREATE_FAILED"
	ErrCodeMessageUpdateFailed = "MESSAGE_UPDATE_FAILED"
	ErrCodeMessageDeleteFailed = "MESSAGE_DELETE_FAILED"
	ErrCodeMessageListFailed   = "MESSAGE_LIST_FAILED"
)

// Error messages
const (
	MsgMessageNotFound     = "Message not found"
	MsgMessageCreateFailed = "Failed to create message"
	MsgMessageUpdateFailed = "Failed to update message"
	MsgMessageDeleteFailed = "Failed to delete message"
	MsgMessageListFailed   = "Failed to list messages"
)

// Predefined errors
var (
	ErrMessageNotFound = customerror.NewCustomError(
		ErrCodeMessageNotFound,
		MsgMessageNotFound,
		http.StatusNotFound,
	)

	ErrMessageCreateFailed = customerror.NewCustomError(
		ErrCodeMessageCreateFailed,
		MsgMessageCreateFailed,
		http.StatusInternalServerError,
	)

	ErrMessageUpdateFailed = customerror.NewCustomError(
		ErrCodeMessageUpdateFailed,
		MsgMessageUpdateFailed,
		http.StatusInternalServerError,
	)

	ErrMessageDeleteFailed = customerror.NewCustomError(
		ErrCodeMessageDeleteFailed,
		MsgMessageDeleteFailed,
		http.StatusInternalServerError,
	)

	ErrMessageListFailed = customerror.NewCustomError(
		ErrCodeMessageListFailed,
		MsgMessageListFailed,
		http.StatusInternalServerError,
	)
)
