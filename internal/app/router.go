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

	a.router = router
}
