package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/srcndev/message-service/internal/service"
	"github.com/srcndev/message-service/pkg/customresponse"
)

// MessageHandler interface defines message HTTP handlers
type MessageHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// messageHandler is the private implementation of MessageHandler interface
type messageHandler struct {
	service service.MessageService
}

// Compile-time interface compliance check
var _ MessageHandler = (*messageHandler)(nil)

// NewMessageHandler creates a new message handler
func NewMessageHandler(service service.MessageService) MessageHandler {
	return &messageHandler{
		service: service,
	}
}

// RegisterRoutes registers all message routes
func (h *messageHandler) RegisterRoutes(router *gin.RouterGroup) {
	messages := router.Group("/messages")
	{
		messages.POST("", h.Create)
		messages.GET("/:id", h.GetByID)
		messages.GET("", h.List)
		messages.PUT("/:id", h.Update)
		messages.DELETE("/:id", h.Delete)
	}
}

func (h *messageHandler) Create(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		customresponse.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	message, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	customresponse.Success(c, http.StatusCreated, dto.ToResponse(message))
}

func (h *messageHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		customresponse.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid message ID")
		return
	}

	message, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	customresponse.Success(c, http.StatusOK, dto.ToResponse(message))
}

func (h *messageHandler) List(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	messages, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.Error(err)
		return
	}

	responses := make([]dto.MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = dto.ToResponse(message)
	}

	customresponse.Success(c, http.StatusOK, responses)
}

func (h *messageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		customresponse.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid message ID")
		return
	}

	var req dto.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		customresponse.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	message, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		c.Error(err)
		return
	}

	customresponse.Success(c, http.StatusOK, dto.ToResponse(message))
}

func (h *messageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		customresponse.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid message ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	customresponse.Success(c, http.StatusNoContent, map[string]interface{}(nil))
}
