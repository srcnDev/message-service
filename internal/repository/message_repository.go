package repository

import (
	"context"

	"github.com/srcndev/message-service/internal/domain"
	"gorm.io/gorm"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, id uint) (*domain.Message, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Message, error)
	GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error)
	Update(ctx context.Context, message *domain.Message) error
	Delete(ctx context.Context, id uint) error
}

type messageRepository struct {
	db *gorm.DB
}

// Compile-time interface compliance check
var _ MessageRepository = (*messageRepository)(nil)

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Create inserts a new message into the database
func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetByID retrieves a message by its ID
func (r *messageRepository) GetByID(ctx context.Context, id uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// List retrieves all messages with pagination
func (r *messageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// GetPendingMessages retrieves pending messages with limit
func (r *messageRepository) GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).
		Where("status = ?", domain.StatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// Update updates an existing message
func (r *messageRepository) Update(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

// Delete soft deletes a message
func (r *messageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Message{}, id).Error
}
