package dto

import (
	"time"

	"github.com/srcndev/message-service/internal/domain"
)

// MessageResponse represents the response payload for a message
type MessageResponse struct {
	ID          uint                 `json:"id" example:"1"`
	PhoneNumber string               `json:"phoneNumber" example:"+905551111111"`
	Content     string               `json:"content" example:"Hello"`
	Status      domain.MessageStatus `json:"status" example:"pending"`
	MessageID   *string              `json:"messageId,omitempty" example:"67f2f8a8-ea58-4ed0-a6f9-ff217df4d849"`
	SentAt      *time.Time           `json:"sentAt,omitempty" example:"2025-11-09T10:30:00Z"`
	CreatedAt   time.Time            `json:"createdAt" example:"2025-11-09T10:00:00Z"`
	UpdatedAt   time.Time            `json:"updatedAt" example:"2025-11-09T10:00:00Z"`
}

// ToResponse converts domain model to response DTO
func ToResponse(m *domain.Message) MessageResponse {
	return MessageResponse{
		ID:          m.ID,
		PhoneNumber: m.PhoneNumber,
		Content:     m.Content,
		Status:      m.Status,
		MessageID:   m.MessageID,
		SentAt:      m.SentAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
