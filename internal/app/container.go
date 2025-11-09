package app

import (
	"log"

	"gorm.io/gorm"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/health"
	"github.com/srcndev/message-service/pkg/database"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Repositories

	// Services
	HealthService health.Service

	// Handlers
	HealthHandler health.Handler
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	container := &Container{
		Config: cfg,
		DB:     db,
	}

	// Wire dependencies
	container.setupRepositories()
	container.setupServices()
	container.setupHandlers()

	return container, nil
}

// setupRepositories initializes all repositories
func (c *Container) setupRepositories() {
	// Message repository
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
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			log.Printf("Warning: failed to get database instance: %v", err)
			return nil
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Warning: failed to close database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
	return nil
}
