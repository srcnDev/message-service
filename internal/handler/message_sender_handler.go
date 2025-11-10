package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/internal/job"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/customresponse"
)

// MessageSenderHandler interface defines message sender HTTP handlers
type MessageSenderHandler interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
	Status(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// messageSenderHandler is the private implementation of MessageSenderHandler interface
type messageSenderHandler struct {
	messageSenderJob job.MessageSenderJob
}

// Compile-time interface compliance check
var _ MessageSenderHandler = (*messageSenderHandler)(nil)

// NewMessageSenderHandler creates a new message sender handler
func NewMessageSenderHandler(messageSenderJob job.MessageSenderJob) MessageSenderHandler {
	return &messageSenderHandler{
		messageSenderJob: messageSenderJob,
	}
}

// RegisterRoutes registers message sender routes
func (h *messageSenderHandler) RegisterRoutes(router *gin.RouterGroup) {
	sender := router.Group("/sender")
	{
		sender.POST("/start", h.Start)
		sender.POST("/stop", h.Stop)
		sender.GET("/status", h.Status)
	}
}

// Start godoc
// @Summary      Start message sender
// @Description  Start the message sender job manually
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  customresponse.CustomResponse{data=map[string]string}
// @Failure      400  {object}  customresponse.CustomResponse
// @Failure      500  {object}  customresponse.CustomResponse
// @Router       /sender/start [post]
func (h *messageSenderHandler) Start(c *gin.Context) {
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

// Stop godoc
// @Summary      Stop message sender
// @Description  Stop the message sender job manually
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  customresponse.CustomResponse{data=map[string]string}
// @Failure      400  {object}  customresponse.CustomResponse
// @Failure      500  {object}  customresponse.CustomResponse
// @Router       /sender/stop [post]
func (h *messageSenderHandler) Stop(c *gin.Context) {
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

// Status godoc
// @Summary      Get sender status
// @Description  Check if the message sender job is running
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  customresponse.CustomResponse{data=map[string]bool}
// @Router       /sender/status [get]
func (h *messageSenderHandler) Status(c *gin.Context) {
	customresponse.Success(c, http.StatusOK, gin.H{
		"running": h.messageSenderJob.IsRunning(),
	})
}
