package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/apperror"
	"github.com/srcndev/message-service/pkg/logger"
)

// App represents the application with all dependencies
type App struct {
	container *Container
	router    *gin.Engine
	server    *http.Server
}

// New creates and initializes the application
func New(cfg *config.Config) (*App, error) {
	// Create DI container
	container, err := NewContainer(cfg)
	if err != nil {
		return nil, apperror.ErrContainerInitFailed.WithError(err)
	}

	app := &App{
		container: container,
	}

	// Setup router
	app.setupRouter()

	// Setup HTTP server
	app.server = &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      app.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return app, nil
}

// Run starts the HTTP server
func (a *App) Run() error {
	// Start background jobs before starting HTTP server
	if err := a.container.StartJobs(); err != nil {
		logger.Error("Failed to start background jobs: %v", err)
	}

	logger.Info("Starting server on %s", a.server.Addr)
	logger.Info("Health check: http://localhost%s/health", a.server.Addr)
	logger.Info("API base URL: http://localhost%s/api/v1", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return apperror.ErrServerStartFailed.WithError(err)
	}

	return nil
}

// Shutdown gracefully stops the application
func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return apperror.ErrServerStopFailed.WithError(err)
	}

	// Close container resources
	if err := a.container.Close(); err != nil {
		logger.Error("Container close error: %v", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}
