package app

import (
	"context"

	"gorm.io/gorm"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/handler"
	"github.com/srcndev/message-service/internal/job"
	"github.com/srcndev/message-service/internal/repository"
	"github.com/srcndev/message-service/internal/service"
	"github.com/srcndev/message-service/pkg/database"
	"github.com/srcndev/message-service/pkg/health"
	"github.com/srcndev/message-service/pkg/logger"
	"github.com/srcndev/message-service/pkg/redis"
	"github.com/srcndev/message-service/pkg/webhook"
)

// Container holds all application dependencies
type Container struct {
	Config      *config.Config
	DB          *gorm.DB
	RedisClient redis.Client

	// Repositories
	MessageRepo      repository.MessageRepository
	MessageCacheRepo repository.MessageCacheRepository

	// Services
	HealthService        health.Service
	MessageService       service.MessageService
	MessageSenderService service.MessageSenderService

	// Jobs
	MessageSenderJob job.MessageSenderJob

	// Handlers
	HealthHandler  health.Handler
	MessageHandler handler.MessageHandler
	SenderHandler  handler.SenderHandler

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

	// Initialize Redis if enabled
	if cfg.Redis.Enabled {
		redisClient, err := redis.NewClient(redis.Config{
			Host:     cfg.Redis.Host,
			Port:     cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err != nil {
			logger.Error("Failed to connect to Redis: %v (continuing without cache)", err)
			// Don't fail the app, just disable Redis
			cfg.Redis.Enabled = false
		} else {
			container.RedisClient = redisClient
			logger.Info("Redis connection established")
		}
	}

	// Wire dependencies
	container.setupClients()
	container.setupRepositories()
	container.setupServices()
	container.setupHandlers()

	// Note: Migrations and seeding should be done via migrate CLI tool
	// This ensures clean separation of concerns and better control

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
	c.MessageRepo = repository.NewMessageRepository(c.DB)

	// Initialize cache repository if Redis is enabled
	if c.Config.Redis.Enabled && c.RedisClient != nil {
		c.MessageCacheRepo = repository.NewMessageCacheRepository(c.RedisClient)
	}
}

// setupServices initializes all services
func (c *Container) setupServices() {
	c.HealthService = health.NewService()
	c.MessageService = service.NewMessageService(c.MessageRepo)
	c.MessageSenderService = service.NewMessageSenderService(
		c.MessageService,
		c.MessageCacheRepo,
		c.WebhookClient,
		c.Config.MessageSender.BatchSize,
		c.Config.Redis.Enabled,
	)

	// Create scheduler job
	messageSenderJob, err := job.NewMessageSenderJob(
		c.MessageSenderService,
		c.Config.MessageSender.Interval,
	)
	if err != nil {
		logger.Fatal("Failed to create message sender job: %v", err)
	}
	c.MessageSenderJob = messageSenderJob
}

// setupHandlers initializes all HTTP handlers
func (c *Container) setupHandlers() {
	c.HealthHandler = health.NewHandler(c.HealthService)
	c.MessageHandler = handler.NewMessageHandler(c.MessageService)
	c.SenderHandler = handler.NewSenderHandler(c.MessageSenderJob)
}

// StartJobs starts all background jobs
func (c *Container) StartJobs() error {
	logger.Info("Starting background jobs...")

	// Use background context for the job lifecycle
	ctx := context.Background()

	if err := c.MessageSenderJob.Start(ctx); err != nil {
		return err
	}

	logger.Info("Background jobs started successfully")
	return nil
}

// Close gracefully closes all resources
func (c *Container) Close() error {
	// Stop message sender job first
	if c.MessageSenderJob != nil && c.MessageSenderJob.IsRunning() {
		logger.Info("Stopping background jobs...")
		ctx := context.Background()
		if err := c.MessageSenderJob.Stop(ctx); err != nil {
			logger.Error("Failed to stop message sender job: %v", err)
		}
	}

	// Close Redis connection if exists
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection: %v", err)
		} else {
			logger.Info("Redis connection closed")
		}
	}

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			logger.Error("Failed to get database instance: %v", err)
			return nil
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database: %v", err)
		} else {
			logger.Info("Database connection closed")
		}
	}
	return nil
}
