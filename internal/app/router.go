package app

import (
	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/pkg/middleware"
)

// setupRouter configures all HTTP routes
func (a *App) setupRouter() {
	router := gin.Default()

	// Global error handler middleware
	router.Use(middleware.ErrorHandler())

	// Health check route (outside versioned API)
	a.container.HealthHandler.RegisterRoutes(&router.RouterGroup)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		a.container.MessageHandler.RegisterRoutes(v1)
		a.container.SenderHandler.RegisterRoutes(v1)
	}

	a.router = router
}
