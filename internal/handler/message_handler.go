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
	ListSent(c *gin.Context)
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
		messages.GET("/sent", h.ListSent)
		messages.PUT("/:id", h.Update)
		messages.DELETE("/:id", h.Delete)
	}
}

// Create godoc
// @Summary      Create a new message
// @Description  Create a new message to be sent via webhook
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        message  body      dto.CreateMessageRequest  true  "Message details"
// @Success      201      {object}  customresponse.CustomResponse{data=dto.MessageResponse}
// @Failure      400      {object}  customresponse.CustomResponse
// @Failure      500      {object}  customresponse.CustomResponse
// @Router       /messages [post]
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

// GetByID godoc
// @Summary      Get message by ID
// @Description  Get a single message by its ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Message ID"
// @Success      200  {object}  customresponse.CustomResponse{data=dto.MessageResponse}
// @Failure      400  {object}  customresponse.CustomResponse
// @Failure      404  {object}  customresponse.CustomResponse
// @Failure      500  {object}  customresponse.CustomResponse
// @Router       /messages/{id} [get]
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

// List godoc
// @Summary      List messages
// @Description  Get a list of messages with pagination
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit"   default(10)
// @Param        offset  query     int  false  "Offset"  default(0)
// @Success      200     {object}  customresponse.CustomResponse{data=[]dto.MessageResponse}
// @Failure      500     {object}  customresponse.CustomResponse
// @Router       /messages [get]
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

// ListSent godoc
// @Summary      List sent messages
// @Description  Get a list of sent messages with pagination
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit"   default(10)
// @Param        offset  query     int  false  "Offset"  default(0)
// @Success      200     {object}  customresponse.CustomResponse{data=[]dto.MessageResponse}
// @Failure      500     {object}  customresponse.CustomResponse
// @Router       /messages/sent [get]
func (h *messageHandler) ListSent(c *gin.Context) {
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

	messages, err := h.service.ListSentMessages(c.Request.Context(), limit, offset)
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

// Update godoc
// @Summary      Update message
// @Description  Update an existing message by ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Message ID"
// @Param        message  body      dto.UpdateMessageRequest  true  "Message details"
// @Success      200      {object}  customresponse.CustomResponse{data=dto.MessageResponse}
// @Failure      400      {object}  customresponse.CustomResponse
// @Failure      404      {object}  customresponse.CustomResponse
// @Failure      500      {object}  customresponse.CustomResponse
// @Router       /messages/{id} [put]
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

// Delete godoc
// @Summary      Delete message
// @Description  Soft delete a message by ID
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Message ID"
// @Success      204  {object}  customresponse.CustomResponse
// @Failure      400  {object}  customresponse.CustomResponse
// @Failure      404  {object}  customresponse.CustomResponse
// @Failure      500  {object}  customresponse.CustomResponse
// @Router       /messages/{id} [delete]
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
