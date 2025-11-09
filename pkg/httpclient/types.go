package httpclient

import (
	"time"
)

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// Config holds HTTP client configuration
type Config struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	DefaultHeaders map[string]string
}
