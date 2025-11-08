package health

import (
	"github.com/gin-gonic/gin"
)

// Handler handles health check HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new health check handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Check handles GET /health endpoint
func (h *Handler) Check(c *gin.Context) {
	status := h.service.GetStatus()
	c.JSON(200, status)
}
