package service

import (
	"context"
	"errors"
	"time"

	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/srcndev/message-service/internal/repository"
	"gorm.io/gorm"
)

// MessageService defines the business logic interface for messages
type MessageService interface {
	Create(ctx context.Context, req dto.CreateMessageRequest) (*domain.Message, error)
	GetByID(ctx context.Context, id uint) (*domain.Message, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Message, error)
	GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error)
	SetSent(ctx context.Context, id uint, messageID string) error
	SetFailed(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, req dto.UpdateMessageRequest) (*domain.Message, error)
	Delete(ctx context.Context, id uint) error
}

type messageService struct {
	repo repository.MessageRepository
}

// Compile-time interface compliance check
var _ MessageService = (*messageService)(nil)

// NewMessageService creates a new message service
func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{
		repo: repo,
	}
}

// Create creates a new message
func (s *messageService) Create(ctx context.Context, req dto.CreateMessageRequest) (*domain.Message, error) {
	message := &domain.Message{
		PhoneNumber: req.PhoneNumber,
		Content:     req.Content,
		Status:      domain.StatusPending,
	}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, dto.ErrMessageCreateFailed.WithError(err)
	}

	return message, nil
}

// GetByID retrieves a message by ID
func (s *messageService) GetByID(ctx context.Context, id uint) (*domain.Message, error) {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrMessageNotFound
		}
		return nil, dto.ErrMessageListFailed.WithError(err)
	}

	return message, nil
}

// List retrieves all messages with pagination
func (s *messageService) List(ctx context.Context, limit, offset int) ([]*domain.Message, error) {
	messages, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, dto.ErrMessageListFailed.WithError(err)
	}

	return messages, nil
}

// GetPendingMessages retrieves pending messages
func (s *messageService) GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error) {
	messages, err := s.repo.GetPendingMessages(ctx, limit)
	if err != nil {
		return nil, dto.ErrMessageListFailed.WithError(err)
	}
	return messages, nil
}

// SetSent marks a message as sent
func (s *messageService) SetSent(ctx context.Context, id uint, messageID string) error {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMessageNotFound
		}
		return dto.ErrMessageUpdateFailed.WithError(err)
	}

	now := time.Now()
	message.Status = domain.StatusSent
	message.MessageID = &messageID
	message.SentAt = &now

	if err := s.repo.Update(ctx, message); err != nil {
		return dto.ErrMessageUpdateFailed.WithError(err)
	}

	return nil
}

// SetFailed marks a message as failed
func (s *messageService) SetFailed(ctx context.Context, id uint) error {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMessageNotFound
		}
		return dto.ErrMessageUpdateFailed.WithError(err)
	}

	message.Status = domain.StatusFailed

	if err := s.repo.Update(ctx, message); err != nil {
		return dto.ErrMessageUpdateFailed.WithError(err)
	}

	return nil
}

// Update updates an existing message
func (s *messageService) Update(ctx context.Context, id uint, req dto.UpdateMessageRequest) (*domain.Message, error) {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrMessageNotFound
		}
		return nil, dto.ErrMessageUpdateFailed.WithError(err)
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
		return nil, dto.ErrMessageUpdateFailed.WithError(err)
	}

	return message, nil
}

// Delete deletes a message
func (s *messageService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrMessageNotFound
		}
		return dto.ErrMessageDeleteFailed.WithError(err)
	}
	return nil
}
