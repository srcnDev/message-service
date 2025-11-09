package messagesender

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/customresponse"
)

// Handler interface defines message sender HTTP handlers
type Handler interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
	Status(c *gin.Context)
	RegisterRoutes(rg *gin.RouterGroup)
}

// handler is the private implementation of Handler interface
type handler struct {
	job Job
}

// Compile-time interface compliance check
var _ Handler = (*handler)(nil)

// NewHandler creates a new message sender handler
func NewHandler(job Job) Handler {
	return &handler{
		job: job,
	}
}

// RegisterRoutes registers message sender routes
func (h *handler) RegisterRoutes(rg *gin.RouterGroup) {
	sender := rg.Group("/sender")
	{
		sender.POST("/start", h.Start)
		sender.POST("/stop", h.Stop)
		sender.GET("/status", h.Status)
	}
}

func (h *handler) Start(c *gin.Context) {
	if err := h.job.Start(c.Request.Context()); err != nil {
		if customErr, ok := err.(*customerror.CustomError); ok {
			customresponse.Error(c, customErr.GetStatusCode(), customErr.Code, customErr.Message)
		} else {
			customresponse.Error(c, http.StatusInternalServerError, "START_FAILED", err.Error())
		}
		return
	}

	customresponse.Success(c, http.StatusOK, gin.H{"message": "Message sender started"})
}

func (h *handler) Stop(c *gin.Context) {
	if err := h.job.Stop(c.Request.Context()); err != nil {
		if customErr, ok := err.(*customerror.CustomError); ok {
			customresponse.Error(c, customErr.GetStatusCode(), customErr.Code, customErr.Message)
		} else {
			customresponse.Error(c, http.StatusInternalServerError, "STOP_FAILED", err.Error())
		}
		return
	}

	customresponse.Success(c, http.StatusOK, gin.H{"message": "Message sender stopped"})
}

func (h *handler) Status(c *gin.Context) {
	customresponse.Success(c, http.StatusOK, gin.H{
		"running": h.job.IsRunning(),
	})
}
