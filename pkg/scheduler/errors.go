package scheduler

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

// Error codes
const (
	ErrCodeSchedulerInvalidInterval = "SCHEDULER_INVALID_INTERVAL"
	ErrCodeSchedulerNilJob          = "SCHEDULER_NIL_JOB"
	ErrCodeSchedulerAlreadyRunning  = "SCHEDULER_ALREADY_RUNNING"
	ErrCodeSchedulerNotRunning      = "SCHEDULER_NOT_RUNNING"
)

// Error messages
const (
	MsgSchedulerInvalidInterval = "Interval must be positive"
	MsgSchedulerNilJob          = "Job cannot be nil"
	MsgSchedulerAlreadyRunning  = "Scheduler already running"
	MsgSchedulerNotRunning      = "Scheduler not running"
)

// Predefined errors
var (
	ErrInvalidInterval = customerror.New(
		ErrCodeSchedulerInvalidInterval,
		MsgSchedulerInvalidInterval,
		http.StatusBadRequest,
	)

	ErrNilJob = customerror.New(
		ErrCodeSchedulerNilJob,
		MsgSchedulerNilJob,
		http.StatusBadRequest,
	)

	ErrAlreadyRunning = customerror.New(
		ErrCodeSchedulerAlreadyRunning,
		MsgSchedulerAlreadyRunning,
		http.StatusConflict,
	)

	ErrNotRunning = customerror.New(
		ErrCodeSchedulerNotRunning,
		MsgSchedulerNotRunning,
		http.StatusConflict,
	)
)
