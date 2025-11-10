package app

import (
	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/pkg/middleware"

	"github.com/srcndev/message-service/docs" // Swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// setupRouter configures all HTTP routes
func (a *App) setupRouter() {
	router := gin.Default()

	// Global error handler middleware
	router.Use(middleware.ErrorHandler())

	// Health check route (outside versioned API)
	a.container.HealthHandler.RegisterRoutes(&router.RouterGroup)

	// API Documentation - Swagger UI
	// Initialize swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "Message Service API"
	docs.SwaggerInfo.Version = "1.0"

	// Swagger route with custom handler to handle /swagger redirect
	router.GET("/swagger/*any", func(c *gin.Context) {
		// If path is exactly /swagger or /swagger/, redirect to /swagger/index.html
		if c.Request.URL.Path == "/swagger" || c.Request.URL.Path == "/swagger/" {
			c.Redirect(301, "/swagger/index.html")
			return
		}
		// Otherwise, use the swagger handler
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		a.container.MessageHandler.RegisterRoutes(v1)
		a.container.MessageSenderHandler.RegisterRoutes(v1)
	}

	a.router = router
}
