package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/config"
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
		return nil, fmt.Errorf("failed to create container: %w", err)
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
	log.Printf("Starting server on %s", a.server.Addr)
	log.Printf("Health check: http://localhost%s/health", a.server.Addr)
	log.Printf("API base URL: http://localhost%s/api/v1", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close container resources
	if err := a.container.Close(); err != nil {
		log.Printf("Warning: container close error: %v", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
