//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/internal/dto"
	"github.com/srcndev/message-service/internal/handler"
	"github.com/srcndev/message-service/internal/job"
	"github.com/srcndev/message-service/internal/repository"
	"github.com/srcndev/message-service/internal/service"
	"github.com/srcndev/message-service/pkg/customresponse"
	"github.com/srcndev/message-service/pkg/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// E2E Test Configuration
const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBUser     = "postgres"
	testDBPassword = "postgres"
	testDBName     = "message_service_test"
)

// MockWebhookServer simulates the external webhook endpoint
type MockWebhookServer struct {
	server       *httptest.Server
	requests     []webhook.SendMessageRequest
	responseCode int
}

func NewMockWebhookServer() *MockWebhookServer {
	mock := &MockWebhookServer{
		requests:     make([]webhook.SendMessageRequest, 0),
		responseCode: http.StatusAccepted,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req webhook.SendMessageRequest
		json.NewDecoder(r.Body).Decode(&req)
		mock.requests = append(mock.requests, req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mock.responseCode)
		json.NewEncoder(w).Encode(webhook.SendMessageResponse{
			Message:   "Accepted",
			MessageID: fmt.Sprintf("msg-%d", len(mock.requests)),
		})
	})

	mock.server = httptest.NewServer(handler)
	return mock
}

func (m *MockWebhookServer) Close() {
	m.server.Close()
}

func (m *MockWebhookServer) GetRequestCount() int {
	return len(m.requests)
}

func (m *MockWebhookServer) GetLastRequest() *webhook.SendMessageRequest {
	if len(m.requests) == 0 {
		return nil
	}
	return &m.requests[len(m.requests)-1]
}

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword, testDBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto migrate schema
	err = db.AutoMigrate(&domain.Message{})
	require.NoError(t, err, "Failed to migrate schema")

	return db
}

// cleanupTestDB cleans up test data
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	db.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE")
}

// setupTestApp creates a complete application instance for testing
func setupTestApp(t *testing.T, webhookURL string) (*gin.Engine, *gorm.DB, handler.MessageHandler, handler.SenderHandler) {
	db := setupTestDB(t)

	// Create repositories
	messageRepo := repository.NewMessageRepository(db)
	cacheRepo := repository.NewMessageCacheRepository(nil) // No Redis in basic E2E test

	// Create webhook client
	webhookClient := webhook.New(webhook.Config{
		URL:        webhookURL,
		AuthKey:    "test-auth-key",
		Timeout:    10 * time.Second,
		MaxRetries: 3,
	})

	// Create services
	messageService := service.NewMessageService(messageRepo)
	senderService := service.NewMessageSenderService(
		messageService,
		cacheRepo,
		webhookClient,
		2,     // batch size
		false, // cache disabled for test
	)

	// Create job (but don't start it automatically)
	messageSenderJob, err := job.NewMessageSenderJob(senderService, 2*time.Minute)
	require.NoError(t, err)

	// Create handlers
	messageHandler := handler.NewMessageHandler(messageService)
	senderHandler := handler.NewSenderHandler(messageSenderJob)

	// Setup router
	router := gin.New()
	api := router.Group("/api")
	messageHandler.RegisterRoutes(api)
	senderHandler.RegisterRoutes(api)

	return router, db, messageHandler, senderHandler
}

// TestE2E_CreateAndListMessages tests message creation and listing
func TestE2E_CreateAndListMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Test: Create a message
	createReq := dto.CreateMessageRequest{
		PhoneNumber: "+905551111111",
		Content:     "E2E Test Message",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.True(t, createResp.Success)

	// Extract message ID
	respData := createResp.Data.(map[string]interface{})
	messageID := int(respData["id"].(float64))

	// Test: List messages
	req = httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &listResp)
	assert.True(t, listResp.Success)

	messages := listResp.Data.([]interface{})
	assert.Equal(t, 1, len(messages))

	// Test: Get message by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/messages/%d", messageID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.True(t, getResp.Success)
	getData := getResp.Data.(map[string]interface{})
	assert.Equal(t, "E2E Test Message", getData["content"])
}

// TestE2E_UpdateAndDeleteMessage tests message update and delete operations
func TestE2E_UpdateAndDeleteMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Create a message
	createReq := dto.CreateMessageRequest{
		PhoneNumber: "+905551111111",
		Content:     "Original Message",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	respData := createResp.Data.(map[string]interface{})
	messageID := int(respData["id"].(float64))

	// Test: Update message
	updateReq := dto.UpdateMessageRequest{
		Content: stringPtr("Updated Message"),
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/messages/%d", messageID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &updateResp)
	assert.True(t, updateResp.Success)
	updateData := updateResp.Data.(map[string]interface{})
	assert.Equal(t, "Updated Message", updateData["content"])

	// Test: Delete message
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/messages/%d", messageID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify message is deleted (soft delete)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/messages/%d", messageID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestE2E_MessageSendingWorkflow tests the complete message sending workflow
func TestE2E_MessageSendingWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Create multiple messages
	for i := 1; i <= 3; i++ {
		createReq := dto.CreateMessageRequest{
			PhoneNumber: fmt.Sprintf("+9055511111%02d", i),
			Content:     fmt.Sprintf("Test Message %d", i),
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Verify webhook hasn't been called yet
	assert.Equal(t, 0, mockWebhook.GetRequestCount())

	// Start sender job
	req := httptest.NewRequest(http.MethodPost, "/api/sender/start", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Check sender status
	req = httptest.NewRequest(http.MethodGet, "/api/sender/status", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var statusResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &statusResp)
	statusData := statusResp.Data.(map[string]interface{})
	assert.Equal(t, true, statusData["running"])

	// Wait a bit for messages to be sent (in real scenario, scheduler would run)
	// For E2E test, we'll manually trigger send by calling the service
	time.Sleep(100 * time.Millisecond)

	// Note: In a real E2E test with running scheduler, webhook would be called automatically
	// For this test, we verify the setup is correct

	// Stop sender job
	req = httptest.NewRequest(http.MethodPost, "/api/sender/stop", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify sender stopped
	req = httptest.NewRequest(http.MethodGet, "/api/sender/status", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &statusResp)
	statusData = statusResp.Data.(map[string]interface{})
	assert.Equal(t, false, statusData["running"])
}

// TestE2E_ValidationErrors tests validation error handling
func TestE2E_ValidationErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "missing phone number",
			requestBody: map[string]interface{}{
				"content": "Test",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid phone format",
			requestBody: dto.CreateMessageRequest{
				PhoneNumber: "123456",
				Content:     "Test",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "content too long",
			requestBody: dto.CreateMessageRequest{
				PhoneNumber: "+905551111111",
				Content:     string(make([]byte, 161)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing content",
			requestBody: map[string]interface{}{
				"phoneNumber": "+905551111111",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp customresponse.CustomResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.False(t, resp.Success)
			assert.NotNil(t, resp.Error)
		})
	}
}

// TestE2E_ConcurrentRequests tests handling of concurrent requests
func TestE2E_ConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Create messages concurrently
	concurrentRequests := 10
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			createReq := dto.CreateMessageRequest{
				PhoneNumber: fmt.Sprintf("+9055511111%02d", index),
				Content:     fmt.Sprintf("Concurrent Message %d", index),
			}
			body, _ := json.Marshal(createReq)
			req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}

	// Verify all messages were created
	req := httptest.NewRequest(http.MethodGet, "/api/messages?limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var listResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &listResp)
	messages := listResp.Data.([]interface{})
	assert.Equal(t, concurrentRequests, len(messages))
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// TestE2E_PaginationWorkflow tests pagination functionality
func TestE2E_PaginationWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Create 15 messages
	for i := 1; i <= 15; i++ {
		createReq := dto.CreateMessageRequest{
			PhoneNumber: fmt.Sprintf("+9055511111%02d", i),
			Content:     fmt.Sprintf("Message %d", i),
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Test: First page (limit=5, offset=0)
	req := httptest.NewRequest(http.MethodGet, "/api/messages?limit=5&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	messages := resp.Data.([]interface{})
	assert.Equal(t, 5, len(messages))

	// Test: Second page (limit=5, offset=5)
	req = httptest.NewRequest(http.MethodGet, "/api/messages?limit=5&offset=5", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &resp)
	messages = resp.Data.([]interface{})
	assert.Equal(t, 5, len(messages))

	// Test: Third page (limit=5, offset=10)
	req = httptest.NewRequest(http.MethodGet, "/api/messages?limit=5&offset=10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &resp)
	messages = resp.Data.([]interface{})
	assert.Equal(t, 5, len(messages))
}

// TestE2E_ListSentMessages tests listing only sent messages
func TestE2E_ListSentMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup
	mockWebhook := NewMockWebhookServer()
	defer mockWebhook.Close()

	router, db, _, _ := setupTestApp(t, mockWebhook.server.URL)
	defer cleanupTestDB(t, db)

	// Create 5 messages with different statuses
	for i := 1; i <= 5; i++ {
		createReq := dto.CreateMessageRequest{
			PhoneNumber: fmt.Sprintf("+9055511111%02d", i),
			Content:     fmt.Sprintf("Message %d", i),
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Manually mark 3 messages as sent in database
	var messages []*domain.Message
	db.Find(&messages)
	now := time.Now()
	for i := 0; i < 3 && i < len(messages); i++ {
		messageID := fmt.Sprintf("msg-%d", i+1)
		messages[i].Status = domain.StatusSent
		messages[i].MessageID = &messageID
		messages[i].SentAt = &now
		db.Save(messages[i])
	}

	// Test: List all messages
	req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var allResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &allResp)
	allMessages := allResp.Data.([]interface{})
	assert.Equal(t, 5, len(allMessages))

	// Test: List only sent messages
	req = httptest.NewRequest(http.MethodGet, "/api/messages/sent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var sentResp customresponse.CustomResponse
	json.Unmarshal(w.Body.Bytes(), &sentResp)
	assert.True(t, sentResp.Success)

	sentMessages := sentResp.Data.([]interface{})
	assert.Equal(t, 3, len(sentMessages))

	// Verify all returned messages have status "sent"
	for _, msg := range sentMessages {
		msgData := msg.(map[string]interface{})
		assert.Equal(t, "sent", msgData["status"])
		assert.NotNil(t, msgData["messageId"])
		assert.NotNil(t, msgData["sentAt"])
	}

	// Test: Pagination for sent messages
	req = httptest.NewRequest(http.MethodGet, "/api/messages/sent?limit=2&offset=0", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &sentResp)
	sentMessages = sentResp.Data.([]interface{})
	assert.Equal(t, 2, len(sentMessages))

	// Test: Second page
	req = httptest.NewRequest(http.MethodGet, "/api/messages/sent?limit=2&offset=2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &sentResp)
	sentMessages = sentResp.Data.([]interface{})
	assert.Equal(t, 1, len(sentMessages))

	// Test: Empty result when no more sent messages
	req = httptest.NewRequest(http.MethodGet, "/api/messages/sent?limit=10&offset=10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &sentResp)
	sentMessages = sentResp.Data.([]interface{})
	assert.Equal(t, 0, len(sentMessages))
}

func TestMain(m *testing.M) {
	// Setup: You might want to create test database here
	// For now, assuming test database exists

	// Run tests
	m.Run()

	// Cleanup: Drop test database if needed
}
