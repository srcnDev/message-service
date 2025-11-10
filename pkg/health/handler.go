package health

import (
	"github.com/gin-gonic/gin"
)

// Handler interface defines health check HTTP handlers
type Handler interface {
	Check(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// handler is the private implementation of Handler interface
type handler struct {
	service Service
}

// Compile-time interface compliance check
var _ Handler = (*handler)(nil)

// NewHealthHandler creates a new health check handler
func NewHealthHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

// RegisterRoutes registers health check routes
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.Check)
}

// Check godoc
// @Summary      Health check
// @Description  Check if the service is healthy
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  health.Status
// @Router       /health [get]
func (h *handler) Check(c *gin.Context) {
	status := h.service.GetStatus()
	c.JSON(200, status)
}
