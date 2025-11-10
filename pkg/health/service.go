package health

import (
	"time"
)

// Service defines health check operations
type Service interface {
	GetStatus() Status
}

// service handles health check logic
type service struct {
	startTime time.Time
}

// Compile-time interface compliance check
var _ Service = (*service)(nil)

// NewHealthService creates a new health check service
func NewHealthService() Service {
	return &service{
		startTime: time.Now(),
	}
}

// GetStatus returns current health status
func (s *service) GetStatus() Status {
	return Status{
		Status: "healthy",
		Uptime: time.Since(s.startTime).String(),
	}
}
