package messagesender

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/srcndev/message-service/pkg/scheduler"
)

// Job defines the interface for scheduled message sending
type Job interface {
	// Start starts the scheduled job
	Start(ctx context.Context) error
	// Stop stops the scheduled job
	Stop(ctx context.Context) error
	// IsRunning returns whether the job is running
	IsRunning() bool
}

// job manages the scheduled message sending
type job struct {
	service   Service
	scheduler scheduler.Scheduler
}

// Compile-time interface compliance check
var _ Job = (*job)(nil)

// NewJob creates a new message sender job with scheduler
func NewJob(service Service, interval time.Duration) (Job, error) {
	j := &job{
		service: service,
	}

	// Create scheduler
	sch, err := scheduler.New(j.run, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}
	j.scheduler = sch

	return j, nil
}

// run is the job function that gets executed by scheduler
func (j *job) run(ctx context.Context) error {
	log.Println("[Job] Starting message sending cycle...")

	err := j.service.SendPendingMessages(ctx)
	if err != nil {
		log.Printf("[Job] Error sending messages: %v\n", err)
		return err
	}

	log.Println("[Job] Message sending cycle completed")
	return nil
}

// Start starts the scheduled job
func (j *job) Start(ctx context.Context) error {
	log.Println("[Job] Starting scheduler...")
	return j.scheduler.Start(ctx)
}

// Stop stops the scheduled job
func (j *job) Stop(ctx context.Context) error {
	log.Println("[Job] Stopping scheduler...")
	return j.scheduler.Stop(ctx)
}

// IsRunning returns whether the job is running
func (j *job) IsRunning() bool {
	return j.scheduler.IsRunning()
}
