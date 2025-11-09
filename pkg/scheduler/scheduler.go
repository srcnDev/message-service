package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

// Scheduler defines the interface for a job scheduler
type Scheduler interface {
	// Start begins executing the job at the specified interval
	Start(ctx context.Context) error

	// Stop gracefully stops the scheduler
	Stop(ctx context.Context) error

	// IsRunning returns whether the scheduler is currently running
	IsRunning() bool
}

// scheduler is the private implementation of Scheduler interface
type scheduler struct {
	job      Job
	interval time.Duration

	mu        sync.Mutex
	running   bool
	ticker    *time.Ticker
	stoppedCh chan struct{}
	cancel    context.CancelFunc
}

// Compile-time interface compliance check
var _ Scheduler = (*scheduler)(nil)

// New creates a new scheduler with the given job and interval
func New(job Job, interval time.Duration) (*scheduler, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if job == nil {
		return nil, ErrNilJob
	}

	return &scheduler{
		job:      job,
		interval: interval,
	}, nil
}

// Start begins executing the job at the specified interval
func (s *scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return ErrAlreadyRunning
	}

	// Use background context for long-running scheduler
	// Don't use the HTTP request context which will be cancelled
	jobCtx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.stoppedCh = make(chan struct{})
	s.ticker = time.NewTicker(s.interval)
	s.running = true

	go s.run(jobCtx)

	return nil
}

// Stop gracefully stops the scheduler
func (s *scheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return ErrNotRunning
	}

	s.cancel()
	s.ticker.Stop()

	// Wait for graceful shutdown with timeout
	select {
	case <-s.stoppedCh:
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
	}

	s.running = false
	return nil
}

// IsRunning returns whether the scheduler is currently running
func (s *scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// run is the main scheduler loop
func (s *scheduler) run(ctx context.Context) {
	defer close(s.stoppedCh)

	log.Println("[Scheduler] Starting, executing job immediately...")
	// Execute job immediately on start (for pending messages)
	s.executeJob(ctx)

	// Continue executing on interval
	log.Printf("[Scheduler] Will run every %v\n", s.interval)
	for {
		select {
		case <-s.ticker.C:
			log.Println("[Scheduler] Tick received, executing job...")
			s.executeJob(ctx)
		case <-ctx.Done():
			log.Println("[Scheduler] Context cancelled, stopping...")
			return
		}
	}
}

// executeJob executes the job safely
func (s *scheduler) executeJob(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Scheduler] Job panicked: %v\n", r)
		}
	}()

	if err := s.job(ctx); err != nil {
		log.Printf("[Scheduler] Job returned error: %v (will retry on next tick)\n", err)
	}
}
