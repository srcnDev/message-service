package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Run("creates service successfully", func(t *testing.T) {
		svc := NewService()

		assert.NotNil(t, svc)
		assert.IsType(t, &service{}, svc)
	})
}

func TestService_GetStatus(t *testing.T) {
	t.Run("returns healthy status", func(t *testing.T) {
		// Setup
		svc := NewService().(*service)

		// Execute
		status := svc.GetStatus()

		// Verify
		assert.Equal(t, "healthy", status.Status)
		assert.NotEmpty(t, status.Uptime)
	})

	t.Run("uptime increases over time", func(t *testing.T) {
		// Setup
		svc := NewService().(*service)

		// Execute - first check
		status1 := svc.GetStatus()
		time.Sleep(100 * time.Millisecond)
		status2 := svc.GetStatus()

		// Verify - uptime should be different
		assert.Equal(t, "healthy", status1.Status)
		assert.Equal(t, "healthy", status2.Status)
		assert.NotEqual(t, status1.Uptime, status2.Uptime)
	})

	t.Run("uptime format is valid duration string", func(t *testing.T) {
		// Setup
		svc := NewService().(*service)

		// Execute
		status := svc.GetStatus()

		// Verify - should be parseable as duration
		_, err := time.ParseDuration(status.Uptime)
		assert.NoError(t, err, "uptime should be a valid duration string")
	})

	t.Run("multiple calls return consistent status", func(t *testing.T) {
		// Setup
		svc := NewService().(*service)

		// Execute
		for i := 0; i < 5; i++ {
			status := svc.GetStatus()
			assert.Equal(t, "healthy", status.Status)
			assert.NotEmpty(t, status.Uptime)
		}
	})
}

func TestService_StartTime(t *testing.T) {
	t.Run("start time is initialized on creation", func(t *testing.T) {
		// Setup
		before := time.Now()
		svc := NewService().(*service)
		after := time.Now()

		// Verify
		assert.True(t, svc.startTime.After(before) || svc.startTime.Equal(before))
		assert.True(t, svc.startTime.Before(after) || svc.startTime.Equal(after))
	})
}

func TestService_InterfaceCompliance(t *testing.T) {
	t.Run("service implements Service interface", func(t *testing.T) {
		var _ Service = (*service)(nil)
		var _ Service = NewService()
	})
}

func TestStatus_Structure(t *testing.T) {
	t.Run("status has correct fields", func(t *testing.T) {
		status := Status{
			Status: "healthy",
			Uptime: "1m30s",
		}

		assert.Equal(t, "healthy", status.Status)
		assert.Equal(t, "1m30s", status.Uptime)
	})
}
