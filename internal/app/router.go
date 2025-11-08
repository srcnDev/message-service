package app

import (
	"github.com/gin-gonic/gin"
)

// setupRouter configures all HTTP routes
func (a *App) setupRouter() {
	router := gin.Default()

	// Health check
	router.GET("/health", a.container.HealthHandler.Check)

	a.router = router
}
