package domain

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a message to be sent
type Message struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PhoneNumber string         `gorm:"type:varchar(20);not null;index" json:"phoneNumber"`
	Content     string         `gorm:"type:varchar(160);not null" json:"content"`
	Status      MessageStatus  `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	MessageID   *string        `gorm:"type:varchar(100);uniqueIndex" json:"messageId,omitempty"`
	SentAt      *time.Time     `json:"sentAt,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Message) TableName() string {
	return "messages"
}
