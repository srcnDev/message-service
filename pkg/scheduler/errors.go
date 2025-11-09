package scheduler

import (
	"net/http"

	"github.com/srcndev/message-service/pkg/customerror"
)

var (
	// ErrInvalidInterval is returned when interval is invalid
	ErrInvalidInterval = customerror.New("SCHEDULER_INVALID_INTERVAL", "interval must be positive", http.StatusBadRequest)

	// ErrNilJob is returned when job is nil
	ErrNilJob = customerror.New("SCHEDULER_NIL_JOB", "job cannot be nil", http.StatusBadRequest)

	// ErrAlreadyRunning is returned when scheduler is already running
	ErrAlreadyRunning = customerror.New("SCHEDULER_ALREADY_RUNNING", "scheduler already running", http.StatusConflict)

	// ErrNotRunning is returned when scheduler is not running
	ErrNotRunning = customerror.New("SCHEDULER_NOT_RUNNING", "scheduler not running", http.StatusConflict)
)

// Compile-time interface compliance check
var _ error = (*customerror.CustomError)(nil)
