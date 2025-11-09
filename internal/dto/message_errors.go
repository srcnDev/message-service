package dto

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
	ErrMessageNotFound = customerror.New(
		ErrCodeMessageNotFound,
		MsgMessageNotFound,
		http.StatusNotFound,
	)

	ErrMessageCreateFailed = customerror.New(
		ErrCodeMessageCreateFailed,
		MsgMessageCreateFailed,
		http.StatusInternalServerError,
	)

	ErrMessageUpdateFailed = customerror.New(
		ErrCodeMessageUpdateFailed,
		MsgMessageUpdateFailed,
		http.StatusInternalServerError,
	)

	ErrMessageDeleteFailed = customerror.New(
		ErrCodeMessageDeleteFailed,
		MsgMessageDeleteFailed,
		http.StatusInternalServerError,
	)

	ErrMessageListFailed = customerror.New(
		ErrCodeMessageListFailed,
		MsgMessageListFailed,
		http.StatusInternalServerError,
	)
)
