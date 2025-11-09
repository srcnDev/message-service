package dto

import "github.com/srcndev/message-service/internal/domain"

// UpdateMessageRequest represents the request payload for updating a message
type UpdateMessageRequest struct {
	PhoneNumber *string               `json:"phoneNumber,omitempty" binding:"omitempty,e164"`
	Content     *string               `json:"content,omitempty" binding:"omitempty,max=160"`
	Status      *domain.MessageStatus `json:"status,omitempty" binding:"omitempty,oneof=pending sent"`
}
