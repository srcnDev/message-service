package message

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for message data operations
type Repository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id uint) (*Message, error)
	List(ctx context.Context, limit, offset int) ([]*Message, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

// Compile-time interface compliance check
var _ Repository = (*repository)(nil)

// NewRepository creates a new message repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new message into the database
func (r *repository) Create(ctx context.Context, message *Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetByID retrieves a message by its ID
func (r *repository) GetByID(ctx context.Context, id uint) (*Message, error) {
	var message Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// List retrieves all messages with pagination
func (r *repository) List(ctx context.Context, limit, offset int) ([]*Message, error) {
	var messages []*Message
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// Update updates an existing message
func (r *repository) Update(ctx context.Context, message *Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

// Delete soft deletes a message
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Message{}, id).Error
}
