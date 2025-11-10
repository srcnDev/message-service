package scheduler

import (
	"context"
	"time"
)

// Job represents a function to be executed by the scheduler
type Job func(ctx context.Context) error

// Config holds scheduler configuration
type Config struct {
	// Job is the function to execute
	Job Job

	// Interval between job executions
	Interval time.Duration

	// Name for the scheduler (optional, for logging)
	Name string

	// MaxRetries for job execution (0 = no retry)
	MaxRetries int

	// RetryDelay between retries
	RetryDelay time.Duration
}

// Option is a functional option for scheduler configuration
type Option func(*Config)

// WithName sets the scheduler name
func WithName(name string) Option {
	return func(c *Config) {
		c.Name = name
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}

// WithRetryDelay sets the delay between retries
func WithRetryDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.RetryDelay = delay
	}
}
