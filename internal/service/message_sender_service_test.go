package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/srcndev/message-service/internal/repository"
	"github.com/srcndev/message-service/pkg/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageService mocks MessageService interface
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) Create(ctx context.Context, req dto.CreateMessageRequest) (*domain.Message, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) GetByID(ctx context.Context, id uint) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) List(ctx context.Context, limit, offset int) ([]*domain.Message, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageService) GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageService) SetSent(ctx context.Context, id uint, messageID string) error {
	args := m.Called(ctx, id, messageID)
	return args.Error(0)
}

func (m *MockMessageService) Update(ctx context.Context, id uint, req dto.UpdateMessageRequest) (*domain.Message, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockWebhookClient mocks webhook.Client interface
type MockWebhookClient struct {
	mock.Mock
}

func (m *MockWebhookClient) SendMessage(ctx context.Context, req *webhook.SendMessageRequest) (*webhook.SendMessageResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*webhook.SendMessageResponse), args.Error(1)
}

// MockCacheRepository mocks MessageCacheRepository interface
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) CacheSentMessage(ctx context.Context, messageID string, sentAt time.Time) error {
	args := m.Called(ctx, messageID, sentAt)
	return args.Error(0)
}

func (m *MockCacheRepository) GetCachedMessage(ctx context.Context, messageID string) (*repository.CachedMessage, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CachedMessage), args.Error(1)
}

func (m *MockCacheRepository) IsCached(ctx context.Context, messageID string) (bool, error) {
	args := m.Called(ctx, messageID)
	return args.Bool(0), args.Error(1)
}

func TestMessageSenderService_SendPendingMessages_Success(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, true)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
		{ID: 2, PhoneNumber: "+905552222222", Content: "Message 2", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)

	// First message
	mockWebhook.On("SendMessage", mock.Anything, mock.MatchedBy(func(req *webhook.SendMessageRequest) bool {
		return req.To == "+905551111111" && req.Content == "Message 1"
	})).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-1",
	}, nil)
	mockMsgService.On("SetSent", mock.Anything, uint(1), "webhook-id-1").Return(nil)
	mockCache.On("CacheSentMessage", mock.Anything, "webhook-id-1", mock.Anything).Return(nil)

	// Second message
	mockWebhook.On("SendMessage", mock.Anything, mock.MatchedBy(func(req *webhook.SendMessageRequest) bool {
		return req.To == "+905552222222" && req.Content == "Message 2"
	})).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-2",
	}, nil)
	mockMsgService.On("SetSent", mock.Anything, uint(2), "webhook-id-2").Return(nil)
	mockCache.On("CacheSentMessage", mock.Anything, "webhook-id-2", mock.Anything).Return(nil)

	err := service.SendPendingMessages(context.Background())

	assert.NoError(t, err)
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestMessageSenderService_SendPendingMessages_NoPendingMessages(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, true)

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return([]*domain.Message{}, nil)

	err := service.SendPendingMessages(context.Background())

	assert.NoError(t, err)
	mockMsgService.AssertExpectations(t)
	// Webhook should not be called
	mockWebhook.AssertNotCalled(t, "SendMessage")
}

func TestMessageSenderService_SendPendingMessages_GetPendingError(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)

	dbError := errors.New("database error")
	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(nil, dbError)

	err := service.SendPendingMessages(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_LIST_FAILED")
	mockMsgService.AssertExpectations(t)
}

func TestMessageSenderService_SendPendingMessages_WebhookFailure(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
		{ID: 2, PhoneNumber: "+905552222222", Content: "Message 2", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)

	// Both messages fail webhook
	webhookError := errors.New("webhook connection error")
	mockWebhook.On("SendMessage", mock.Anything, mock.Anything).Return(nil, webhookError)

	err := service.SendPendingMessages(context.Background())

	// Should return error when ALL messages fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_SEND_FAILED")
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
	// SetSent should NOT be called for failed messages
	mockMsgService.AssertNotCalled(t, "SetSent")
}

func TestMessageSenderService_SendPendingMessages_PartialSuccess(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
		{ID: 2, PhoneNumber: "+905552222222", Content: "Message 2", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)

	// First message succeeds
	mockWebhook.On("SendMessage", mock.Anything, mock.MatchedBy(func(req *webhook.SendMessageRequest) bool {
		return req.To == "+905551111111"
	})).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-1",
	}, nil)
	mockMsgService.On("SetSent", mock.Anything, uint(1), "webhook-id-1").Return(nil)

	// Second message fails
	mockWebhook.On("SendMessage", mock.Anything, mock.MatchedBy(func(req *webhook.SendMessageRequest) bool {
		return req.To == "+905552222222"
	})).Return(nil, errors.New("webhook error"))

	err := service.SendPendingMessages(context.Background())

	// Should NOT error because at least one succeeded
	assert.NoError(t, err)
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestMessageSenderService_SendPendingMessages_SetSentFailure(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)

	mockWebhook.On("SendMessage", mock.Anything, mock.Anything).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-1",
	}, nil)

	// SetSent fails
	mockMsgService.On("SetSent", mock.Anything, uint(1), "webhook-id-1").Return(errors.New("db error"))

	err := service.SendPendingMessages(context.Background())

	// Should return error because SetSent failed
	assert.Error(t, err)
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
}

func TestMessageSenderService_SendPendingMessages_CacheDisabled(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	// Cache disabled
	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)
	mockWebhook.On("SendMessage", mock.Anything, mock.Anything).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-1",
	}, nil)
	mockMsgService.On("SetSent", mock.Anything, uint(1), "webhook-id-1").Return(nil)

	err := service.SendPendingMessages(context.Background())

	assert.NoError(t, err)
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
	// Cache should NOT be called when disabled
	mockCache.AssertNotCalled(t, "CacheSentMessage")
}

func TestMessageSenderService_SendPendingMessages_CacheFailureNonBlocking(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, true)

	pendingMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
	}

	mockMsgService.On("GetPendingMessages", mock.Anything, 2).Return(pendingMessages, nil)
	mockWebhook.On("SendMessage", mock.Anything, mock.Anything).Return(&webhook.SendMessageResponse{
		Message:   "Accepted",
		MessageID: "webhook-id-1",
	}, nil)
	mockMsgService.On("SetSent", mock.Anything, uint(1), "webhook-id-1").Return(nil)

	// Cache fails but should not block operation
	mockCache.On("CacheSentMessage", mock.Anything, "webhook-id-1", mock.Anything).Return(errors.New("redis error"))

	err := service.SendPendingMessages(context.Background())

	// Should still succeed even if cache fails
	assert.NoError(t, err)
	mockMsgService.AssertExpectations(t)
	mockWebhook.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestNewMessageSenderService_DefaultBatchSize(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	// Invalid batch size (0 or negative)
	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 0, false)

	assert.NotNil(t, service)

	// Verify it uses default batch size (2)
	svc, ok := service.(*messageSenderService)
	assert.True(t, ok)
	assert.Equal(t, 2, svc.batchSize)
}

func TestNewMessageSenderService_CustomBatchSize(t *testing.T) {
	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 5, false)

	assert.NotNil(t, service)

	svc, ok := service.(*messageSenderService)
	assert.True(t, ok)
	assert.Equal(t, 5, svc.batchSize)
}

func TestMessageSenderService_InterfaceCompliance(t *testing.T) {
	var _ MessageSenderService = (*messageSenderService)(nil) // Compile-time check

	mockMsgService := new(MockMessageService)
	mockWebhook := new(MockWebhookClient)
	mockCache := new(MockCacheRepository)

	service := NewMessageSenderService(mockMsgService, mockCache, mockWebhook, 2, false)
	assert.NotNil(t, service)
}
