package job

import (
	"context"
	"time"

	"github.com/srcndev/message-service/internal/apperror"
	"github.com/srcndev/message-service/internal/service"
	"github.com/srcndev/message-service/pkg/logger"
	"github.com/srcndev/message-service/pkg/scheduler"
)

// MessageSenderJob defines the interface for scheduled message sending
type MessageSenderJob interface {
	// Start starts the scheduled job
	Start(ctx context.Context) error
	// Stop stops the scheduled job
	Stop(ctx context.Context) error
	// IsRunning returns whether the job is running
	IsRunning() bool
}

// messageSenderJob manages the scheduled message sending
type messageSenderJob struct {
	senderService service.MessageSenderService
	scheduler     scheduler.Scheduler
}

// Compile-time interface compliance check
var _ MessageSenderJob = (*messageSenderJob)(nil)

// NewMessageSenderJob creates a new message sender job with the sender service
func NewMessageSenderJob(senderService service.MessageSenderService, interval time.Duration) (MessageSenderJob, error) {
	j := &messageSenderJob{
		senderService: senderService,
	}

	// Create scheduler
	sch, err := scheduler.New(j.run, interval)
	if err != nil {
		return nil, apperror.ErrSchedulerInitFailed.WithError(err)
	}
	j.scheduler = sch

	return j, nil
}

// run is the job function that gets executed by scheduler
func (j *messageSenderJob) run(ctx context.Context) error {
	logger.Info("Starting message sending cycle")

	err := j.senderService.SendPendingMessages(ctx)
	if err != nil {
		logger.Error("Error sending messages: %v", err)
		return err
	}

	logger.Info("Message sending cycle completed")
	return nil
}

// Start starts the scheduled job
func (j *messageSenderJob) Start(ctx context.Context) error {
	logger.Info("Starting message sender job")
	return j.scheduler.Start(ctx)
}

// Stop stops the scheduled job
func (j *messageSenderJob) Stop(ctx context.Context) error {
	logger.Info("Stopping message sender job")
	return j.scheduler.Stop(ctx)
}

// IsRunning returns whether the job is running
func (j *messageSenderJob) IsRunning() bool {
	return j.scheduler.IsRunning()
}
