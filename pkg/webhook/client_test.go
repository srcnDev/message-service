package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/srcndev/message-service/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock for httpclient.Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(ctx context.Context, req *httpclient.Request) (*httpclient.Response, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func (m *MockHTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*httpclient.Response, error) {
	args := m.Called(ctx, url, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func (m *MockHTTPClient) Post(ctx context.Context, url string, body any, headers map[string]string) (*httpclient.Response, error) {
	args := m.Called(ctx, url, body, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func (m *MockHTTPClient) Put(ctx context.Context, url string, body any, headers map[string]string) (*httpclient.Response, error) {
	args := m.Called(ctx, url, body, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func (m *MockHTTPClient) Delete(ctx context.Context, url string, headers map[string]string) (*httpclient.Response, error) {
	args := m.Called(ctx, url, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func (m *MockHTTPClient) Patch(ctx context.Context, url string, body any, headers map[string]string) (*httpclient.Response, error) {
	args := m.Called(ctx, url, body, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*httpclient.Response), args.Error(1)
}

func TestClient_SendMessage_Success(t *testing.T) {
	tests := []struct {
		name           string
		request        *SendMessageRequest
		statusCode     int
		responseBody   SendMessageResponse
		expectedResult *SendMessageResponse
	}{
		{
			name: "successful send with 202 Accepted",
			request: &SendMessageRequest{
				To:      "+905551234567",
				Content: "Test message",
			},
			statusCode: http.StatusAccepted,
			responseBody: SendMessageResponse{
				Message:   "Accepted",
				MessageID: "test-message-id-123",
			},
			expectedResult: &SendMessageResponse{
				Message:   "Accepted",
				MessageID: "test-message-id-123",
			},
		},
		{
			name: "successful send with 200 OK",
			request: &SendMessageRequest{
				To:      "+905551234567",
				Content: "Test message",
			},
			statusCode: http.StatusOK,
			responseBody: SendMessageResponse{
				Message:   "OK",
				MessageID: "test-message-id-456",
			},
			expectedResult: &SendMessageResponse{
				Message:   "OK",
				MessageID: "test-message-id-456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := new(MockHTTPClient)

			responseBytes, _ := json.Marshal(tt.responseBody)
			mockHTTP.On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(&httpclient.Response{
					StatusCode: tt.statusCode,
					Body:       responseBytes,
				}, nil)

			client := &client{
				httpClient: mockHTTP,
				baseURL:    "https://webhook.test",
				authKey:    "test-key",
			}

			result, err := client.SendMessage(context.Background(), tt.request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult.Message, result.Message)
			assert.Equal(t, tt.expectedResult.MessageID, result.MessageID)
			mockHTTP.AssertExpectations(t)
		})
	}
}

func TestClient_SendMessage_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		request     *SendMessageRequest
		expectedErr error
	}{
		{
			name:        "nil request",
			request:     nil,
			expectedErr: ErrInvalidRequest,
		},
		{
			name: "empty phone number",
			request: &SendMessageRequest{
				To:      "",
				Content: "Test message",
			},
			expectedErr: ErrInvalidPhoneNumber,
		},
		{
			name: "empty content",
			request: &SendMessageRequest{
				To:      "+905551234567",
				Content: "",
			},
			expectedErr: ErrEmptyContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := new(MockHTTPClient)

			client := &client{
				httpClient: mockHTTP,
				baseURL:    "https://webhook.test",
				authKey:    "test-key",
			}

			result, err := client.SendMessage(context.Background(), tt.request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestClient_SendMessage_HTTPErrors(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		expectedErrCode string
	}{
		{
			name:            "401 Unauthorized",
			statusCode:      http.StatusUnauthorized,
			expectedErrCode: "WEBHOOK_UNAUTHORIZED",
		},
		{
			name:            "500 Internal Server Error",
			statusCode:      http.StatusInternalServerError,
			expectedErrCode: "WEBHOOK_SERVER_ERROR",
		},
		{
			name:            "503 Service Unavailable",
			statusCode:      http.StatusServiceUnavailable,
			expectedErrCode: "WEBHOOK_SERVER_ERROR",
		},
		{
			name:            "400 Bad Request",
			statusCode:      http.StatusBadRequest,
			expectedErrCode: "WEBHOOK_INVALID_REQUEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := new(MockHTTPClient)

			mockHTTP.On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(&httpclient.Response{
					StatusCode: tt.statusCode,
					Body:       []byte("{}"),
				}, nil)

			client := &client{
				httpClient: mockHTTP,
				baseURL:    "https://webhook.test",
				authKey:    "test-key",
			}

			request := &SendMessageRequest{
				To:      "+905551234567",
				Content: "Test message",
			}

			result, err := client.SendMessage(context.Background(), request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.expectedErrCode)
			mockHTTP.AssertExpectations(t)
		})
	}
}

func TestClient_SendMessage_ConnectionError(t *testing.T) {
	mockHTTP := new(MockHTTPClient)

	connectionErr := errors.New("connection refused")
	mockHTTP.On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, connectionErr)

	client := &client{
		httpClient: mockHTTP,
		baseURL:    "https://webhook.test",
		authKey:    "test-key",
	}

	request := &SendMessageRequest{
		To:      "+905551234567",
		Content: "Test message",
	}

	result, err := client.SendMessage(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "WEBHOOK_CONNECTION_FAILED")
	mockHTTP.AssertExpectations(t)
}

func TestClient_SendMessage_InvalidJSON(t *testing.T) {
	mockHTTP := new(MockHTTPClient)

	mockHTTP.On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&httpclient.Response{
			StatusCode: http.StatusAccepted,
			Body:       []byte("invalid json"),
		}, nil)

	client := &client{
		httpClient: mockHTTP,
		baseURL:    "https://webhook.test",
		authKey:    "test-key",
	}

	request := &SendMessageRequest{
		To:      "+905551234567",
		Content: "Test message",
	}

	result, err := client.SendMessage(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "WEBHOOK_PARSING_ERROR")
	mockHTTP.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	cfg := Config{
		URL:    "https://webhook.test",
		AuthKey:    "test-key",
		Timeout:    30,
		MaxRetries: 3,
	}

	webhookClient := NewWebhookClient(cfg)

	assert.NotNil(t, webhookClient)

	// Type assertion to access private fields
	c, ok := webhookClient.(*client)
	assert.True(t, ok)
	assert.Equal(t, cfg.URL, c.baseURL)
	assert.Equal(t, cfg.AuthKey, c.authKey)
	assert.NotNil(t, c.httpClient)
}

func TestClient_InterfaceCompliance(t *testing.T) {
	var _ Client = (*client)(nil) // Compile-time check

	cfg := Config{
		URL: "https://test.com",
		AuthKey: "key",
	}

	c := NewWebhookClient(cfg)
	assert.NotNil(t, c)

	// Verify it satisfies Client interface
	var client Client = c
	assert.NotNil(t, client)
}

func TestClient_SendMessage_All2xxSuccessCodes(t *testing.T) {
	successCodes := []int{200, 201, 202, 203, 204}

	for _, code := range successCodes {
		t.Run(fmt.Sprintf("status_%d", code), func(t *testing.T) {
			mockHTTP := new(MockHTTPClient)

			responseBody := SendMessageResponse{
				Message:   "Success",
				MessageID: "test-id",
			}
			responseBytes, _ := json.Marshal(responseBody)

			mockHTTP.On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(&httpclient.Response{
					StatusCode: code,
					Body:       responseBytes,
				}, nil)

			client := &client{
				httpClient: mockHTTP,
				baseURL:    "https://webhook.test",
				authKey:    "test-key",
			}

			request := &SendMessageRequest{
				To:      "+905551234567",
				Content: "Test",
			}

			result, err := client.SendMessage(context.Background(), request)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "test-id", result.MessageID)
			mockHTTP.AssertExpectations(t)
		})
	}
}
