package health

import (
	"time"
)

// Service handles health check logic
type Service struct {
	startTime time.Time
}

// NewService creates a new health check service
func NewService() *Service {
	return &Service{
		startTime: time.Now(),
	}
}

// GetStatus returns current health status
func (s *Service) GetStatus() Status {
	return Status{
		Status: "healthy",
		Uptime: time.Since(s.startTime).String(),
	}
}
