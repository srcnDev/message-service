package dto

// CreateMessageRequest represents the request payload for creating a message
type CreateMessageRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,e164" example:"+905551111111"`
	Content     string `json:"content" binding:"required,max=160" example:"Hello World"`
}
