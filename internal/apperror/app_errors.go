package apperror

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes for application lifecycle
const (
	ErrCodeContainerInitFailed = "CONTAINER_INIT_FAILED"
	ErrCodeServerStartFailed   = "SERVER_START_FAILED"
	ErrCodeServerStopFailed    = "SERVER_STOP_FAILED"
	ErrCodeSchedulerInitFailed = "SCHEDULER_INIT_FAILED"
)

// Error messages
const (
	MsgContainerInitFailed = "Failed to initialize application container"
	MsgServerStartFailed   = "Failed to start HTTP server"
	MsgServerStopFailed    = "Failed to stop HTTP server gracefully"
	MsgSchedulerInitFailed = "Failed to initialize message sender job"
)

// Predefined errors
var (
	ErrContainerInitFailed = customerror.New(
		ErrCodeContainerInitFailed,
		MsgContainerInitFailed,
		http.StatusInternalServerError,
	)

	ErrServerStartFailed = customerror.New(
		ErrCodeServerStartFailed,
		MsgServerStartFailed,
		http.StatusInternalServerError,
	)

	ErrServerStopFailed = customerror.New(
		ErrCodeServerStopFailed,
		MsgServerStopFailed,
		http.StatusInternalServerError,
	)

	ErrSchedulerInitFailed = customerror.New(
		ErrCodeSchedulerInitFailed,
		MsgSchedulerInitFailed,
		http.StatusInternalServerError,
	)
)
