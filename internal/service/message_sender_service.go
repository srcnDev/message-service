package service

import (
	"context"
	"time"

	"github.com/srcndev/message-service/internal/apperror"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/repository"
	"github.com/srcndev/message-service/pkg/logger"
	"github.com/srcndev/message-service/pkg/webhook"
)

// MessageSenderService defines the message sender service interface
type MessageSenderService interface {
	// SendPendingMessages fetches and sends pending messages
	SendPendingMessages(ctx context.Context) error
}

type messageSenderService struct {
	messageService MessageService
	cacheRepo      repository.MessageCacheRepository
	webhookClient  webhook.Client
	batchSize      int
	cacheEnabled   bool
}

// Compile-time interface compliance check
var _ MessageSenderService = (*messageSenderService)(nil)

// NewMessageSenderService creates a new message sender service
func NewMessageSenderService(
	messageService MessageService,
	cacheRepo repository.MessageCacheRepository,
	webhookClient webhook.Client,
	batchSize int,
	cacheEnabled bool,
) MessageSenderService {
	if batchSize <= 0 {
		batchSize = 2 // Default batch size from case study
	}

	return &messageSenderService{
		messageService: messageService,
		cacheRepo:      cacheRepo,
		webhookClient:  webhookClient,
		batchSize:      batchSize,
		cacheEnabled:   cacheEnabled,
	}
}

// SendPendingMessages fetches and sends pending messages in batches
func (s *messageSenderService) SendPendingMessages(ctx context.Context) error {
	// Get pending messages
	messages, err := s.messageService.GetPendingMessages(ctx, s.batchSize)
	if err != nil {
		return apperror.ErrMessageListFailed.WithError(err)
	}

	// No pending messages
	if len(messages) == 0 {
		return nil
	}

	// Send each message
	var sendErrors []error
	for _, msg := range messages {
		if err := s.sendMessage(ctx, msg); err != nil {
			// Log error but continue with other messages
			logger.Error("Failed to send message %d: %v", msg.ID, err)
			sendErrors = append(sendErrors, err)
			continue
		}
	}

	// If all messages failed, return error
	if len(sendErrors) > 0 && len(sendErrors) == len(messages) {
		return apperror.ErrMessageSendFailed
	}

	return nil
}

// sendMessage sends a single message via webhook
func (s *messageSenderService) sendMessage(ctx context.Context, msg *domain.Message) error {
	// Prepare webhook request
	req := &webhook.SendMessageRequest{
		To:      msg.PhoneNumber,
		Content: msg.Content,
	}

	// Send via webhook
	resp, err := s.webhookClient.SendMessage(ctx, req)
	if err != nil {
		// Don't mark as failed - leave it pending for retry in next cycle
		logger.Error("Failed to send message %d: %v (will retry in next cycle)", msg.ID, err)
		return apperror.ErrWebhookCallFailed.WithError(err)
	}

	// Mark as sent with messageID from webhook
	if err := s.messageService.SetSent(ctx, msg.ID, resp.MessageID); err != nil {
		return apperror.ErrMarkSentFailed.WithError(err)
	}

	// Cache to Redis if enabled (Bonus feature)
	if s.cacheEnabled && s.cacheRepo != nil {
		sentAt := time.Now()
		if cacheErr := s.cacheRepo.CacheSentMessage(ctx, resp.MessageID, sentAt); cacheErr != nil {
			// Log but don't fail the operation
			logger.Error("Failed to cache message %s to Redis: %v", resp.MessageID, cacheErr)
		} else {
			logger.Debug("Message %s cached to Redis successfully", resp.MessageID)
		}
	}

	logger.Info("Message %d sent successfully (webhook messageId: %s)", msg.ID, resp.MessageID)
	return nil
}
