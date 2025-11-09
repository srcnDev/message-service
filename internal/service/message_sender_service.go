package service

import (
	"context"
	"fmt"
	"log"

	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/pkg/webhook"
)

// MessageSenderService defines the message sender service interface
type MessageSenderService interface {
	// SendPendingMessages fetches and sends pending messages
	SendPendingMessages(ctx context.Context) error
}

type messageSenderService struct {
	messageService MessageService
	webhookClient  webhook.Client
	batchSize      int
}

// Compile-time interface compliance check
var _ MessageSenderService = (*messageSenderService)(nil)

// NewMessageSenderService creates a new message sender service
func NewMessageSenderService(messageService MessageService, webhookClient webhook.Client, batchSize int) MessageSenderService {
	if batchSize <= 0 {
		batchSize = 2 // Default batch size from case study
	}

	return &messageSenderService{
		messageService: messageService,
		webhookClient:  webhookClient,
		batchSize:      batchSize,
	}
}

// SendPendingMessages fetches and sends pending messages in batches
func (s *messageSenderService) SendPendingMessages(ctx context.Context) error {
	// Get pending messages
	messages, err := s.messageService.GetPendingMessages(ctx, s.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending messages: %w", err)
	}

	// No pending messages
	if len(messages) == 0 {
		return nil
	}

	// Send each message
	for _, msg := range messages {
		if err := s.sendMessage(ctx, msg); err != nil {
			// Log error but continue with other messages
			log.Printf("Failed to send message %d: %v\n", msg.ID, err)
			continue
		}
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
		// Mark as failed
		if markErr := s.messageService.SetFailed(ctx, msg.ID); markErr != nil {
			log.Printf("Failed to mark message %d as failed: %v\n", msg.ID, markErr)
		}
		return fmt.Errorf("webhook send failed: %w", err)
	}

	// Mark as sent with messageID from webhook
	if err := s.messageService.SetSent(ctx, msg.ID, resp.MessageID); err != nil {
		return fmt.Errorf("failed to mark message as sent: %w", err)
	}

	return nil
}
