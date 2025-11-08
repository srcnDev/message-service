package health

import (
	"github.com/gin-gonic/gin"
)

// Handler interface defines health check HTTP handlers
type Handler interface {
	Check(c *gin.Context)
}

// handler is the private implementation of Handler interface
type handler struct {
	service Service
}

// Compile-time interface compliance check
var _ Handler = (*handler)(nil)

// NewHandler creates a new health check handler
func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

// Check handles GET /health endpoint
func (h *handler) Check(c *gin.Context) {
	status := h.service.GetStatus()
	c.JSON(200, status)
}
