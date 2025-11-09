package service

import (
	"context"
	"errors"
	"testing"

	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockMessageRepository mocks the MessageRepository interface
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByID(ctx context.Context, id uint) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Message, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Update(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestMessageService_Create_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	req := dto.CreateMessageRequest{
		PhoneNumber: "+905551234567",
		Content:     "Test message",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *domain.Message) bool {
		return msg.PhoneNumber == req.PhoneNumber &&
			msg.Content == req.Content &&
			msg.Status == domain.StatusPending
	})).Return(nil)

	result, err := service.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.PhoneNumber, result.PhoneNumber)
	assert.Equal(t, req.Content, result.Content)
	assert.Equal(t, domain.StatusPending, result.Status)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Create_Error(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	req := dto.CreateMessageRequest{
		PhoneNumber: "+905551234567",
		Content:     "Test message",
	}

	dbError := errors.New("database error")
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(dbError)

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MESSAGE_CREATE_FAILED")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	expectedMsg := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Test message",
		Status:      domain.StatusPending,
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedMsg, nil)

	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMsg.ID, result.ID)
	assert.Equal(t, expectedMsg.PhoneNumber, result.PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MESSAGE_NOT_FOUND")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetByID_Error(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	dbError := errors.New("database error")
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(nil, dbError)

	result, err := service.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MESSAGE_LIST_FAILED")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_List_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	expectedMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
		{ID: 2, PhoneNumber: "+905552222222", Content: "Message 2", Status: domain.StatusSent},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(expectedMessages, nil)

	result, err := service.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedMessages[0].ID, result[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_List_Error(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	dbError := errors.New("database error")
	mockRepo.On("List", mock.Anything, 10, 0).Return(nil, dbError)

	result, err := service.List(context.Background(), 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MESSAGE_LIST_FAILED")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetPendingMessages_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	expectedMessages := []*domain.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Message 1", Status: domain.StatusPending},
		{ID: 2, PhoneNumber: "+905552222222", Content: "Message 2", Status: domain.StatusPending},
	}

	mockRepo.On("GetPendingMessages", mock.Anything, 2).Return(expectedMessages, nil)

	result, err := service.GetPendingMessages(context.Background(), 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetPendingMessages_Error(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	dbError := errors.New("database error")
	mockRepo.On("GetPendingMessages", mock.Anything, 2).Return(nil, dbError)

	result, err := service.GetPendingMessages(context.Background(), 2)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_SetSent_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	existingMsg := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Test message",
		Status:      domain.StatusPending,
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(existingMsg, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(msg *domain.Message) bool {
		return msg.ID == 1 &&
			msg.Status == domain.StatusSent &&
			msg.MessageID != nil &&
			*msg.MessageID == "webhook-msg-id" &&
			msg.SentAt != nil
	})).Return(nil)

	err := service.SetSent(context.Background(), 1, "webhook-msg-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_SetSent_NotFound(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	err := service.SetSent(context.Background(), 999, "webhook-msg-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_NOT_FOUND")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_SetSent_UpdateError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	existingMsg := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Test",
		Status:      domain.StatusPending,
	}

	dbError := errors.New("database error")
	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(existingMsg, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(dbError)

	err := service.SetSent(context.Background(), 1, "webhook-msg-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_UPDATE_FAILED")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Update_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	existingMsg := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Old content",
		Status:      domain.StatusPending,
	}

	newPhone := "+905559999999"
	newContent := "New content"
	newStatus := domain.StatusSent

	updateReq := dto.UpdateMessageRequest{
		PhoneNumber: &newPhone,
		Content:     &newContent,
		Status:      &newStatus,
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(existingMsg, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(msg *domain.Message) bool {
		return msg.ID == 1 &&
			msg.PhoneNumber == newPhone &&
			msg.Content == newContent &&
			msg.Status == newStatus
	})).Return(nil)

	result, err := service.Update(context.Background(), 1, updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newPhone, result.PhoneNumber)
	assert.Equal(t, newContent, result.Content)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Update_PartialUpdate(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	existingMsg := &domain.Message{
		ID:          1,
		PhoneNumber: "+905551234567",
		Content:     "Old content",
		Status:      domain.StatusPending,
	}

	newContent := "New content"
	updateReq := dto.UpdateMessageRequest{
		Content: &newContent,
		// PhoneNumber and Status not provided
	}

	mockRepo.On("GetByID", mock.Anything, uint(1)).Return(existingMsg, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(msg *domain.Message) bool {
		return msg.ID == 1 &&
			msg.PhoneNumber == existingMsg.PhoneNumber && // Unchanged
			msg.Content == newContent && // Changed
			msg.Status == existingMsg.Status // Unchanged
	})).Return(nil)

	result, err := service.Update(context.Background(), 1, updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newContent, result.Content)
	assert.Equal(t, existingMsg.PhoneNumber, result.PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	newContent := "New content"
	updateReq := dto.UpdateMessageRequest{
		Content: &newContent,
	}

	mockRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.Update(context.Background(), 999, updateReq)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MESSAGE_NOT_FOUND")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Delete_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err := service.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(999)).Return(gorm.ErrRecordNotFound)

	err := service.Delete(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_NOT_FOUND")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Delete_Error(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	dbError := errors.New("database error")
	mockRepo.On("Delete", mock.Anything, uint(1)).Return(dbError)

	err := service.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_DELETE_FAILED")
	mockRepo.AssertExpectations(t)
}

func TestMessageService_InterfaceCompliance(t *testing.T) {
	var _ MessageService = (*messageService)(nil) // Compile-time check

	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)
	assert.NotNil(t, service)
}

func TestNewMessageService(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := NewMessageService(mockRepo)

	assert.NotNil(t, service)

	// Type assertion
	svc, ok := service.(*messageService)
	assert.True(t, ok)
	assert.NotNil(t, svc.repo)
}
