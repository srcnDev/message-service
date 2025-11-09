package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/internal/job"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/customresponse"
)

// SenderHandler interface defines message sender HTTP handlers
type SenderHandler interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
	Status(c *gin.Context)
	RegisterRoutes(rg *gin.RouterGroup)
}

// senderHandler is the private implementation of SenderHandler interface
type senderHandler struct {
	messageSenderJob job.MessageSenderJob
}

// Compile-time interface compliance check
var _ SenderHandler = (*senderHandler)(nil)

// NewSenderHandler creates a new message sender handler
func NewSenderHandler(messageSenderJob job.MessageSenderJob) SenderHandler {
	return &senderHandler{
		messageSenderJob: messageSenderJob,
	}
}

// RegisterRoutes registers message sender routes
func (h *senderHandler) RegisterRoutes(rg *gin.RouterGroup) {
	sender := rg.Group("/sender")
	{
		sender.POST("/start", h.Start)
		sender.POST("/stop", h.Stop)
		sender.GET("/status", h.Status)
	}
}

func (h *senderHandler) Start(c *gin.Context) {
	if err := h.messageSenderJob.Start(c.Request.Context()); err != nil {
		if customErr, ok := err.(*customerror.CustomError); ok {
			customresponse.Error(c, customErr.GetStatusCode(), customErr.Code, customErr.Message)
		} else {
			customresponse.Error(c, http.StatusInternalServerError, "START_FAILED", err.Error())
		}
		return
	}

	customresponse.Success(c, http.StatusOK, gin.H{"message": "Message sender started"})
}

func (h *senderHandler) Stop(c *gin.Context) {
	if err := h.messageSenderJob.Stop(c.Request.Context()); err != nil {
		if customErr, ok := err.(*customerror.CustomError); ok {
			customresponse.Error(c, customErr.GetStatusCode(), customErr.Code, customErr.Message)
		} else {
			customresponse.Error(c, http.StatusInternalServerError, "STOP_FAILED", err.Error())
		}
		return
	}

	customresponse.Success(c, http.StatusOK, gin.H{"message": "Message sender stopped"})
}

func (h *senderHandler) Status(c *gin.Context) {
	customresponse.Success(c, http.StatusOK, gin.H{
		"running": h.messageSenderJob.IsRunning(),
	})
}
