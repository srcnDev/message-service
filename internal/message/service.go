package message

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Service defines the business logic interface for messages
type Service interface {
	Create(ctx context.Context, req CreateMessageRequest) (*Message, error)
	GetByID(ctx context.Context, id uint) (*Message, error)
	List(ctx context.Context, limit, offset int) ([]*Message, error)
	Update(ctx context.Context, id uint, req UpdateMessageRequest) (*Message, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

// Compile-time interface compliance check
var _ Service = (*service)(nil)

// NewService creates a new message service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new message
func (s *service) Create(ctx context.Context, req CreateMessageRequest) (*Message, error) {
	message := &Message{
		PhoneNumber: req.PhoneNumber,
		Content:     req.Content,
		Status:      StatusPending,
	}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, ErrMessageCreateFailed.WithError(err)
	}

	return message, nil
}

// GetByID retrieves a message by ID
func (s *service) GetByID(ctx context.Context, id uint) (*Message, error) {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMessageNotFound
		}
		return nil, ErrMessageListFailed.WithError(err)
	}

	return message, nil
}

// List retrieves all messages with pagination
func (s *service) List(ctx context.Context, limit, offset int) ([]*Message, error) {
	messages, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, ErrMessageListFailed.WithError(err)
	}

	return messages, nil
}

// Update updates an existing message
func (s *service) Update(ctx context.Context, id uint, req UpdateMessageRequest) (*Message, error) {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMessageNotFound
		}
		return nil, ErrMessageUpdateFailed.WithError(err)
	}

	// Update only provided fields
	if req.PhoneNumber != nil {
		message.PhoneNumber = *req.PhoneNumber
	}
	if req.Content != nil {
		message.Content = *req.Content
	}
	if req.Status != nil {
		message.Status = *req.Status
	}

	if err := s.repo.Update(ctx, message); err != nil {
		return nil, ErrMessageUpdateFailed.WithError(err)
	}

	return message, nil
}

// Delete deletes a message
func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMessageNotFound
		}
		return ErrMessageDeleteFailed.WithError(err)
	}
	return nil
}
