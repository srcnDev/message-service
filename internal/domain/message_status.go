package domain

// MessageStatus represents the current state of a message
type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSent    MessageStatus = "sent"
)
