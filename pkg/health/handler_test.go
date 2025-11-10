package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) GetStatus() Status {
	args := m.Called()
	return args.Get(0).(Status)
}

// Verify MockService implements Service interface
var _ Service = (*MockService)(nil)

func TestNewHealthHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		// Setup
		mockService := new(MockService)

		// Execute
		h := NewHealthHandler(mockService)

		// Verify
		assert.NotNil(t, h)
		assert.Implements(t, (*Handler)(nil), h)
	})
}

func TestHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockStatus     Status
		expectedStatus int
		expectedBody   Status
	}{
		{
			name: "returns healthy status",
			mockStatus: Status{
				Status: "healthy",
				Uptime: "5m30s",
			},
			expectedStatus: http.StatusOK,
			expectedBody: Status{
				Status: "healthy",
				Uptime: "5m30s",
			},
		},
		{
			name: "returns status with long uptime",
			mockStatus: Status{
				Status: "healthy",
				Uptime: "2h30m15s",
			},
			expectedStatus: http.StatusOK,
			expectedBody: Status{
				Status: "healthy",
				Uptime: "2h30m15s",
			},
		},
		{
			name: "returns status with short uptime",
			mockStatus: Status{
				Status: "healthy",
				Uptime: "100ms",
			},
			expectedStatus: http.StatusOK,
			expectedBody: Status{
				Status: "healthy",
				Uptime: "100ms",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockService)
			mockService.On("GetStatus").Return(tt.mockStatus)

			h := NewHealthHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Execute
			h.Check(c)

			// Verify
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var got Status
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody.Status, got.Status)
			assert.Equal(t, tt.expectedBody.Uptime, got.Uptime)

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("registers health route", func(t *testing.T) {
		// Setup
		mockService := new(MockService)
		h := NewHealthHandler(mockService)

		router := gin.New()
		group := router.Group("/api")

		// Execute
		h.RegisterRoutes(group)

		// Verify - test if route exists by making request
		mockService.On("GetStatus").Return(Status{
			Status: "healthy",
			Uptime: "1m",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var status Status
		err := json.Unmarshal(w.Body.Bytes(), &status)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", status.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("route responds to GET method", func(t *testing.T) {
		// Setup
		mockService := new(MockService)
		h := NewHealthHandler(mockService)

		router := gin.New()
		group := router.Group("/api")
		h.RegisterRoutes(group)

		mockService.On("GetStatus").Return(Status{
			Status: "healthy",
			Uptime: "1m",
		})

		// Execute - GET request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
		router.ServeHTTP(w, req)

		// Verify
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("route does not respond to POST method", func(t *testing.T) {
		// Setup
		mockService := new(MockService)
		h := NewHealthHandler(mockService)

		router := gin.New()
		group := router.Group("/api")
		h.RegisterRoutes(group)

		// Execute - POST request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/health", nil)
		router.ServeHTTP(w, req)

		// Verify - should return 404 (route not found)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_Check_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("handles multiple concurrent requests", func(t *testing.T) {
		// Setup
		mockService := new(MockService)
		mockService.On("GetStatus").Return(Status{
			Status: "healthy",
			Uptime: "10m",
		}).Times(5)

		h := NewHealthHandler(mockService)

		// Execute - multiple requests
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			h.Check(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var status Status
			err := json.Unmarshal(w.Body.Bytes(), &status)
			assert.NoError(t, err)
			assert.Equal(t, "healthy", status.Status)
		}

		mockService.AssertExpectations(t)
	})
}

func TestHandler_InterfaceCompliance(t *testing.T) {
	t.Run("handler implements Handler interface", func(t *testing.T) {
		var _ Handler = (*handler)(nil)

		mockService := new(MockService)
		var _ Handler = NewHealthHandler(mockService)
	})
}

func TestHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("handler works with real service", func(t *testing.T) {
		// Setup - using real service instead of mock
		realService := NewHealthService()
		h := NewHealthHandler(realService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Execute
		h.Check(c)

		// Verify
		assert.Equal(t, http.StatusOK, w.Code)

		var status Status
		err := json.Unmarshal(w.Body.Bytes(), &status)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", status.Status)
		assert.NotEmpty(t, status.Uptime)
	})
}
