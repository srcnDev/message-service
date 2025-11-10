package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	job := func(ctx context.Context) error {
		return nil
	}

	scheduler, err := NewScheduler(job, 100*time.Millisecond)

	assert.NoError(t, err)
	assert.NotNil(t, scheduler)
	assert.Equal(t, 100*time.Millisecond, scheduler.interval)
	assert.NotNil(t, scheduler.job)
	assert.False(t, scheduler.IsRunning())
}

func TestNew_InvalidInterval(t *testing.T) {
	job := func(ctx context.Context) error { return nil }

	tests := []struct {
		name     string
		interval time.Duration
	}{
		{"zero interval", 0},
		{"negative interval", -1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler, err := NewScheduler(job, tt.interval)

			assert.Error(t, err)
			assert.Nil(t, scheduler)
			assert.Contains(t, err.Error(), "SCHEDULER_INVALID_INTERVAL")
		})
	}
}

func TestNew_NilJob(t *testing.T) {
	scheduler, err := NewScheduler(nil, 1*time.Second)

	assert.Error(t, err)
	assert.Nil(t, scheduler)
	assert.Contains(t, err.Error(), "SCHEDULER_NIL_JOB")
}

func TestScheduler_Start_Success(t *testing.T) {
	callCount := 0
	job := func(ctx context.Context) error {
		callCount++
		return nil
	}

	scheduler, _ := NewScheduler(job, 50*time.Millisecond)

	err := scheduler.Start(context.Background())
	assert.NoError(t, err)
	assert.True(t, scheduler.IsRunning())

	// Wait for at least 2 ticks
	time.Sleep(150 * time.Millisecond)

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)

	// Job should be called at least 3 times (immediate + 2 ticks)
	assert.GreaterOrEqual(t, callCount, 3)
}

func TestScheduler_Start_AlreadyRunning(t *testing.T) {
	job := func(ctx context.Context) error { return nil }

	scheduler, _ := NewScheduler(job, 1*time.Second)

	// Start first time
	err := scheduler.Start(context.Background())
	assert.NoError(t, err)

	// Try to start again
	err = scheduler.Start(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SCHEDULER_ALREADY_RUNNING")

	// Cleanup
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)
}

func TestScheduler_Stop_Success(t *testing.T) {
	job := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	scheduler, _ := NewScheduler(job, 100*time.Millisecond)

	// Start scheduler
	_ = scheduler.Start(context.Background())
	assert.True(t, scheduler.IsRunning())

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := scheduler.Stop(stopCtx)

	assert.NoError(t, err)
	assert.False(t, scheduler.IsRunning())
}

func TestScheduler_Stop_NotRunning(t *testing.T) {
	job := func(ctx context.Context) error { return nil }

	scheduler, _ := NewScheduler(job, 1*time.Second)

	// Try to stop without starting
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := scheduler.Stop(stopCtx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SCHEDULER_NOT_RUNNING")
}

func TestScheduler_IsRunning(t *testing.T) {
	job := func(ctx context.Context) error { return nil }

	scheduler, _ := NewScheduler(job, 100*time.Millisecond)

	// Initially not running
	assert.False(t, scheduler.IsRunning())

	// Start it
	_ = scheduler.Start(context.Background())
	assert.True(t, scheduler.IsRunning())

	// Stop it
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)
	assert.False(t, scheduler.IsRunning())
}

func TestScheduler_JobError_ContinuesRunning(t *testing.T) {
	callCount := 0
	job := func(ctx context.Context) error {
		callCount++
		return errors.New("job error")
	}

	scheduler, _ := NewScheduler(job, 50*time.Millisecond)

	_ = scheduler.Start(context.Background())

	// Wait for multiple ticks
	time.Sleep(150 * time.Millisecond)

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)

	// Even with errors, job should continue running
	assert.GreaterOrEqual(t, callCount, 3)
}

func TestScheduler_JobPanic_Recovered(t *testing.T) {
	callCount := 0
	panicCount := 0

	job := func(ctx context.Context) error {
		callCount++
		if callCount == 1 {
			panicCount++
			panic("job panic")
		}
		return nil
	}

	scheduler, _ := NewScheduler(job, 50*time.Millisecond)

	_ = scheduler.Start(context.Background())

	// Wait for multiple ticks
	time.Sleep(150 * time.Millisecond)

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)

	// Job should continue even after panic
	assert.GreaterOrEqual(t, callCount, 3)
	assert.Equal(t, 1, panicCount)
}

func TestScheduler_ContextCancellation(t *testing.T) {
	callCount := 0
	job := func(ctx context.Context) error {
		callCount++
		return nil
	}

	scheduler, _ := NewScheduler(job, 50*time.Millisecond)

	// Start with cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	_ = scheduler.Start(ctx)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Cancel context via Stop
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer stopCancel()
	_ = scheduler.Stop(stopCtx)

	beforeCount := callCount
	time.Sleep(100 * time.Millisecond)
	afterCount := callCount

	// No new calls after stop
	assert.Equal(t, beforeCount, afterCount)
	cancel() // cleanup
}

func TestScheduler_ImmediateExecution(t *testing.T) {
	callCount := 0
	firstCallTime := time.Time{}

	job := func(ctx context.Context) error {
		callCount++
		if callCount == 1 {
			firstCallTime = time.Now()
		}
		return nil
	}

	scheduler, _ := NewScheduler(job, 1*time.Second)

	startTime := time.Now()
	_ = scheduler.Start(context.Background())

	// Wait briefly
	time.Sleep(50 * time.Millisecond)

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)

	// First call should be immediate (within 50ms of start)
	assert.Equal(t, 1, callCount)
	assert.WithinDuration(t, startTime, firstCallTime, 50*time.Millisecond)
}

func TestScheduler_ConcurrentStartStop(t *testing.T) {
	job := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	scheduler, _ := NewScheduler(job, 50*time.Millisecond)

	// Start scheduler
	_ = scheduler.Start(context.Background())

	// Try multiple concurrent starts (should all fail)
	errCount := 0
	for i := 0; i < 5; i++ {
		if err := scheduler.Start(context.Background()); err != nil {
			errCount++
		}
	}
	assert.Equal(t, 5, errCount)

	// Stop scheduler
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = scheduler.Stop(stopCtx)

	// Try multiple concurrent stops (should all fail)
	errCount = 0
	for i := 0; i < 5; i++ {
		stopCtx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel2()
		if err := scheduler.Stop(stopCtx2); err != nil {
			errCount++
		}
	}
	assert.Equal(t, 5, errCount)
}

func TestScheduler_InterfaceCompliance(t *testing.T) {
	var _ Scheduler = (*scheduler)(nil) // Compile-time check

	job := func(ctx context.Context) error { return nil }
	s, _ := NewScheduler(job, 1*time.Second)

	assert.NotNil(t, s)
}
