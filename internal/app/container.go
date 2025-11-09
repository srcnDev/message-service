package app

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/health"
	"github.com/srcndev/message-service/internal/message"
	"github.com/srcndev/message-service/internal/messagesender"
	"github.com/srcndev/message-service/pkg/database"
	"github.com/srcndev/message-service/pkg/webhook"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Repositories
	MessageRepo message.Repository

	// Services
	HealthService        health.Service
	MessageService       message.Service
	MessageSenderService messagesender.Service

	// Jobs
	MessageSenderJob messagesender.Job

	// Handlers
	HealthHandler        health.Handler
	MessageHandler       message.Handler
	MessageSenderHandler messagesender.Handler

	// Clients
	WebhookClient webhook.Client
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
	container.setupClients()
	container.setupRepositories()
	container.setupServices()
	container.setupJobs()
	container.setupHandlers()

	// Run migrations
	if err := container.migrate(); err != nil {
		return nil, err
	}

	return container, nil
}

// setupClients initializes all external clients
func (c *Container) setupClients() {
	c.WebhookClient = webhook.New(webhook.Config{
		BaseURL:    c.Config.Webhook.BaseURL,
		AuthKey:    c.Config.Webhook.AuthKey,
		Timeout:    c.Config.Webhook.Timeout,
		MaxRetries: c.Config.Webhook.MaxRetries,
	})
}

// setupRepositories initializes all repositories
func (c *Container) setupRepositories() {
	c.MessageRepo = message.NewRepository(c.DB)
}

// setupServices initializes all services
func (c *Container) setupServices() {
	c.HealthService = health.NewService()
	c.MessageService = message.NewService(c.MessageRepo)
	c.MessageSenderService = messagesender.NewService(
		c.MessageService,
		c.WebhookClient,
		2, // Batch size: 2 messages per cycle
	)
}

// setupJobs initializes all scheduled jobs
func (c *Container) setupJobs() {
	job, err := messagesender.NewJob(
		c.MessageSenderService,
		5*time.Second, // Every 5 seconds
	)
	if err != nil {
		log.Fatalf("Failed to create message sender job: %v", err)
	}
	c.MessageSenderJob = job
}

// setupHandlers initializes all HTTP handlers
func (c *Container) setupHandlers() {
	c.HealthHandler = health.NewHandler(c.HealthService)
	c.MessageHandler = message.NewHandler(c.MessageService)
	c.MessageSenderHandler = messagesender.NewHandler(c.MessageSenderJob)
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
