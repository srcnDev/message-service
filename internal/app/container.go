package app

import (
	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/health"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config

	// Services
	HealthService *health.Service

	// Handlers
	HealthHandler *health.Handler
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	// Wire dependencies
	container.setupServices()
	container.setupHandlers()

	return container, nil
}

// setupServices initializes all services
func (c *Container) setupServices() {
	// Health service
	c.HealthService = health.NewService()

}

// setupHandlers initializes all HTTP handlers
func (c *Container) setupHandlers() {
	// Health handler
	c.HealthHandler = health.NewHandler(c.HealthService)

}

// Close gracefully closes all resources
func (c *Container) Close() error {

	return nil
}
