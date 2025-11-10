package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/srcndev/message-service/pkg/httpclient"
)

// Client defines the webhook client interface
type Client interface {
	// SendMessage sends a message via webhook
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
}

// client is the private implementation
type client struct {
	httpClient httpclient.Client
	baseURL    string
	authKey    string
}

// Compile-time interface compliance check
var _ Client = (*client)(nil)

// Config holds webhook client configuration
type Config struct {
	URL        string
	AuthKey    string
	Timeout    time.Duration
	MaxRetries int
}

// SendMessageRequest represents the webhook request payload
type SendMessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// SendMessageResponse represents the webhook response
type SendMessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

// New creates a new webhook client
func New(cfg Config) Client {
	httpCfg := httpclient.Config{
		Timeout:    cfg.Timeout,
		MaxRetries: cfg.MaxRetries,
		DefaultHeaders: map[string]string{
			"Content-Type":   "application/json",
			"x-ins-auth-key": cfg.AuthKey,
		},
	}

	return &client{
		httpClient: httpclient.New(httpCfg),
		baseURL:    cfg.URL,
		authKey:    cfg.AuthKey,
	}
}

// SendMessage sends a message via webhook
func (c *client) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.To == "" {
		return nil, ErrInvalidPhoneNumber
	}

	if req.Content == "" {
		return nil, ErrEmptyContent
	}

	// Send HTTP request
	resp, err := c.httpClient.Post(ctx, c.baseURL, req, nil)
	if err != nil {
		return nil, ErrConnectionFailed.WithError(err)
	}

	// Check response status
	if resp.StatusCode == 401 {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode >= 500 {
		return nil, ErrServerError.WithError(fmt.Errorf("status: %d", resp.StatusCode))
	}

	// Accept any 2xx success status (200-299)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ErrInvalidRequest.WithError(fmt.Errorf("unexpected status: %d", resp.StatusCode))
	}

	// Parse response
	var webhookResp SendMessageResponse
	if err := json.Unmarshal(resp.Body, &webhookResp); err != nil {
		return nil, ErrParsingResponse.WithError(err)
	}

	return &webhookResp, nil
}
