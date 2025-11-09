package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/customresponse"
	"github.com/srcndev/message-service/pkg/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageSenderJob is a mock implementation of MessageSenderJob
type MockMessageSenderJob struct {
	mock.Mock
}

func (m *MockMessageSenderJob) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageSenderJob) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageSenderJob) IsRunning() bool {
	args := m.Called()
	return args.Bool(0)
}

// Error handler middleware for tests
func senderErrorHandlerMiddleware() gin.HandlerFunc {
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
func setupSenderRouter(handler SenderHandler) *gin.Engine {
	router := gin.New()
	router.Use(senderErrorHandlerMiddleware())
	handler.RegisterRoutes(router.Group("/api"))
	return router
}

func TestNewSenderHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		mockJob := new(MockMessageSenderJob)
		handler := NewSenderHandler(mockJob)

		assert.NotNil(t, handler)
		assert.Implements(t, (*SenderHandler)(nil), handler)
	})
}

func TestSenderHandler_Start(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*MockMessageSenderJob)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "success - starts sender",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Start", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "Message sender started", data["message"])
			},
		},
		{
			name: "error - already running",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Start", mock.Anything).Return(scheduler.ErrAlreadyRunning)
			},
			expectedStatus: http.StatusConflict,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "SCHEDULER_ALREADY_RUNNING", resp.Error.Code)
			},
		},
		{
			name: "error - generic error",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Start", mock.Anything).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "START_FAILED", resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJob := new(MockMessageSenderJob)
			tt.mockSetup(mockJob)

			handler := NewSenderHandler(mockJob)
			router := setupSenderRouter(handler)

			req := httptest.NewRequest(http.MethodPost, "/api/sender/start", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockJob.AssertExpectations(t)
		})
	}
}

func TestSenderHandler_Stop(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*MockMessageSenderJob)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "success - stops sender",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Stop", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "Message sender stopped", data["message"])
			},
		},
		{
			name: "error - not running",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Stop", mock.Anything).Return(scheduler.ErrNotRunning)
			},
			expectedStatus: http.StatusConflict,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "SCHEDULER_NOT_RUNNING", resp.Error.Code)
			},
		},
		{
			name: "error - generic error",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("Stop", mock.Anything).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "STOP_FAILED", resp.Error.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJob := new(MockMessageSenderJob)
			tt.mockSetup(mockJob)

			handler := NewSenderHandler(mockJob)
			router := setupSenderRouter(handler)

			req := httptest.NewRequest(http.MethodPost, "/api/sender/stop", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockJob.AssertExpectations(t)
		})
	}
}

func TestSenderHandler_Status(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*MockMessageSenderJob)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "success - sender is running",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("IsRunning").Return(true)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, true, data["running"])
			},
		},
		{
			name: "success - sender is not running",
			mockSetup: func(m *MockMessageSenderJob) {
				m.On("IsRunning").Return(false)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp customresponse.CustomResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, false, data["running"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJob := new(MockMessageSenderJob)
			tt.mockSetup(mockJob)

			handler := NewSenderHandler(mockJob)
			router := setupSenderRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/sender/status", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
			mockJob.AssertExpectations(t)
		})
	}
}

func TestSenderHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("registers all routes", func(t *testing.T) {
		mockJob := new(MockMessageSenderJob)
		handler := NewSenderHandler(mockJob)
		router := gin.New()
		handler.RegisterRoutes(router.Group("/api"))

		routes := router.Routes()
		expectedRoutes := map[string]bool{
			"POST /api/sender/start": true,
			"POST /api/sender/stop":  true,
			"GET /api/sender/status": true,
		}

		foundRoutes := make(map[string]bool)
		for _, route := range routes {
			key := route.Method + " " + route.Path
			if expectedRoutes[key] {
				foundRoutes[key] = true
			}
		}

		assert.Equal(t, len(expectedRoutes), len(foundRoutes), "All routes should be registered")
	})
}

func TestSenderHandler_InterfaceCompliance(t *testing.T) {
	t.Run("handler implements SenderHandler interface", func(t *testing.T) {
		mockJob := new(MockMessageSenderJob)
		var _ SenderHandler = NewSenderHandler(mockJob)
	})
}

func TestSenderHandler_MultipleStartStopCycles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("can start and stop multiple times", func(t *testing.T) {
		mockJob := new(MockMessageSenderJob)
		handler := NewSenderHandler(mockJob)
		router := setupSenderRouter(handler)

		// First start
		mockJob.On("Start", mock.Anything).Return(nil).Once()
		req := httptest.NewRequest(http.MethodPost, "/api/sender/start", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// First stop
		mockJob.On("Stop", mock.Anything).Return(nil).Once()
		req = httptest.NewRequest(http.MethodPost, "/api/sender/stop", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Second start
		mockJob.On("Start", mock.Anything).Return(nil).Once()
		req = httptest.NewRequest(http.MethodPost, "/api/sender/start", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		mockJob.AssertExpectations(t)
	})
}

func TestSenderHandler_StatusCheckWhileRunning(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("status changes correctly", func(t *testing.T) {
		mockJob := new(MockMessageSenderJob)
		handler := NewSenderHandler(mockJob)
		router := setupSenderRouter(handler)

		// Check status - not running
		mockJob.On("IsRunning").Return(false).Once()
		req := httptest.NewRequest(http.MethodGet, "/api/sender/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp customresponse.CustomResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, false, data["running"])

		// Start
		mockJob.On("Start", mock.Anything).Return(nil).Once()
		req = httptest.NewRequest(http.MethodPost, "/api/sender/start", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Check status - running
		mockJob.On("IsRunning").Return(true).Once()
		req = httptest.NewRequest(http.MethodGet, "/api/sender/status", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &resp)
		data = resp.Data.(map[string]interface{})
		assert.Equal(t, true, data["running"])

		mockJob.AssertExpectations(t)
	})
}
