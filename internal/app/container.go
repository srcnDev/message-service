package app

import (
	"log"

	"gorm.io/gorm"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/health"
	"github.com/srcndev/message-service/internal/message"
	"github.com/srcndev/message-service/pkg/database"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Repositories
	MessageRepo message.Repository

	// Services
	HealthService  health.Service
	MessageService message.Service

	// Handlers
	HealthHandler  health.Handler
	MessageHandler message.Handler
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

	// Run migrations
	if err := container.migrate(); err != nil {
		return nil, err
	}

	return container, nil
}

// setupRepositories initializes all repositories
func (c *Container) setupRepositories() {
	c.MessageRepo = message.NewRepository(c.DB)
}

// setupServices initializes all services
func (c *Container) setupServices() {
	c.HealthService = health.NewService()
	c.MessageService = message.NewService(c.MessageRepo)
}

// setupHandlers initializes all HTTP handlers
func (c *Container) setupHandlers() {
	c.HealthHandler = health.NewHandler(c.HealthService)
	c.MessageHandler = message.NewHandler(c.MessageService)
}

// migrate runs database migrations
func (c *Container) migrate() error {
	return c.DB.AutoMigrate(
		&message.Message{},
	)
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
