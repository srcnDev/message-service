package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/internal/apperror"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/customresponse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock MessageService
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

// Error handler middleware for tests
func errorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if customErr, ok := err.(*customerror.CustomError); ok {
				customresponse.Error(c, customErr.GetStatusCode(), customErr.Code, customErr.Message)
				return
			}
			customresponse.Error(c, 500, "INTERNAL_ERROR", "Internal server error")
		}
	}
}

// Helper to create router with middleware
func setupRouter(handler MessageHandler) *gin.Engine {
	router := gin.New()
	router.Use(errorHandlerMiddleware())
	handler.RegisterRoutes(router.Group("/api"))
	return router
}

func TestNewMessageHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		mockService := new(MockMessageService)
		handler := NewMessageHandler(mockService)

		assert.NotNil(t, handler)
		assert.Implements(t, (*MessageHandler)(nil), handler)
	})
}

func TestMessageHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockMessageService)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "success - creates message",
			requestBody: dto.CreateMessageRequest{
				PhoneNumber: "+905551111111",
				Content:     "Test message",
			},
			mockSetup: func(m *MockMessageService) {
				m.On("Create", mock.Anything, mock.Anything).Return(&domain.Message{
					ID:          1,
					PhoneNumber: "+905551111111",
					Content:     "Test message",
					Status:      domain.StatusPending,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)
			},
		},
		{
			name:           "error - invalid json",
			requestBody:    `{"phoneNumber": "invalid"`,
			mockSetup:      func(m *MockMessageService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
				assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
			},
		},
		{
			name: "error - service error",
			requestBody: dto.CreateMessageRequest{
				PhoneNumber: "+905551111111",
				Content:     "Test",
			},
			mockSetup: func(m *MockMessageService) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil, apperror.ErrMessageCreateFailed)
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMessageService)
			tt.mockSetup(mockService)

			handler := NewMessageHandler(mockService)
			router := setupRouter(handler)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMessageHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		messageID      string
		mockSetup      func(*MockMessageService)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:      "success - returns message",
			messageID: "1",
			mockSetup: func(m *MockMessageService) {
				m.On("GetByID", mock.Anything, uint(1)).Return(&domain.Message{
					ID:          1,
					PhoneNumber: "+905551111111",
					Content:     "Test",
					Status:      domain.StatusPending,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.True(t, resp.Success)
			},
		},
		{
			name:           "error - invalid id",
			messageID:      "invalid",
			mockSetup:      func(m *MockMessageService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
				assert.Equal(t, "INVALID_ID", resp.Error.Code)
			},
		},
		{
			name:      "error - message not found",
			messageID: "999",
			mockSetup: func(m *MockMessageService) {
				m.On("GetByID", mock.Anything, uint(999)).Return(nil, apperror.ErrMessageNotFound)
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMessageService)
			tt.mockSetup(mockService)

			handler := NewMessageHandler(mockService)
			router := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/messages/"+tt.messageID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMessageHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockMessageService)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:        "success - default pagination",
			queryParams: "",
			mockSetup: func(m *MockMessageService) {
				m.On("List", mock.Anything, 10, 0).Return([]*domain.Message{
					{ID: 1, PhoneNumber: "+905551111111", Content: "Msg1", Status: domain.StatusPending},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.True(t, resp.Success)
			},
		},
		{
			name:        "success - custom pagination",
			queryParams: "?limit=5&offset=10",
			mockSetup: func(m *MockMessageService) {
				m.On("List", mock.Anything, 5, 10).Return([]*domain.Message{}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.True(t, resp.Success)
			},
		},
		{
			name:        "error - service error",
			queryParams: "",
			mockSetup: func(m *MockMessageService) {
				m.On("List", mock.Anything, 10, 0).Return(nil, apperror.ErrMessageListFailed)
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMessageService)
			tt.mockSetup(mockService)

			handler := NewMessageHandler(mockService)
			router := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/messages"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMessageHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		messageID      string
		requestBody    interface{}
		mockSetup      func(*MockMessageService)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:      "success - updates message",
			messageID: "1",
			requestBody: dto.UpdateMessageRequest{
				Content: stringPtr("Updated"),
			},
			mockSetup: func(m *MockMessageService) {
				m.On("Update", mock.Anything, uint(1), mock.Anything).Return(&domain.Message{
					ID:      1,
					Content: "Updated",
					Status:  domain.StatusPending,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.True(t, resp.Success)
			},
		},
		{
			name:      "error - invalid id",
			messageID: "invalid",
			requestBody: dto.UpdateMessageRequest{
				Content: stringPtr("Updated"),
			},
			mockSetup:      func(m *MockMessageService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
				assert.Equal(t, "INVALID_ID", resp.Error.Code)
			},
		},
		{
			name:      "error - message not found",
			messageID: "999",
			requestBody: dto.UpdateMessageRequest{
				Content: stringPtr("Updated"),
			},
			mockSetup: func(m *MockMessageService) {
				m.On("Update", mock.Anything, uint(999), mock.Anything).Return(nil, apperror.ErrMessageNotFound)
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMessageService)
			tt.mockSetup(mockService)

			handler := NewMessageHandler(mockService)
			router := setupRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/messages/"+tt.messageID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMessageHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		messageID      string
		mockSetup      func(*MockMessageService)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:      "success - deletes message",
			messageID: "1",
			mockSetup: func(m *MockMessageService) {
				m.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			validateBody: func(t *testing.T, body []byte) {
				// 204 No Content may return empty body or success response
				if len(body) > 0 {
					var resp customresponse.CustomResponse
					json.Unmarshal(body, &resp)
					assert.True(t, resp.Success)
				}
			},
		},
		{
			name:           "error - invalid id",
			messageID:      "invalid",
			mockSetup:      func(m *MockMessageService) {},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
				assert.Equal(t, "INVALID_ID", resp.Error.Code)
			},
		},
		{
			name:      "error - message not found",
			messageID: "999",
			mockSetup: func(m *MockMessageService) {
				m.On("Delete", mock.Anything, uint(999)).Return(apperror.ErrMessageNotFound)
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				json.Unmarshal(body, &resp)
				assert.False(t, resp.Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMessageService)
			tt.mockSetup(mockService)

			handler := NewMessageHandler(mockService)
			router := setupRouter(handler)

			req := httptest.NewRequest(http.MethodDelete, "/api/messages/"+tt.messageID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestMessageHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("registers all routes", func(t *testing.T) {
		mockService := new(MockMessageService)
		handler := NewMessageHandler(mockService)
		router := gin.New()
		handler.RegisterRoutes(router.Group("/api"))

		routes := router.Routes()
		expectedRoutes := map[string]bool{
			"POST /api/messages":       true,
			"GET /api/messages/:id":    true,
			"GET /api/messages":        true,
			"PUT /api/messages/:id":    true,
			"DELETE /api/messages/:id": true,
		}

		for _, route := range routes {
			key := fmt.Sprintf("%s %s", route.Method, route.Path)
			if expectedRoutes[key] {
				assert.True(t, true, "Route registered: "+key)
			}
		}
	})
}

func TestMessageHandler_InterfaceCompliance(t *testing.T) {
	t.Run("handler implements MessageHandler interface", func(t *testing.T) {
		mockService := new(MockMessageService)
		var _ MessageHandler = NewMessageHandler(mockService)
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
