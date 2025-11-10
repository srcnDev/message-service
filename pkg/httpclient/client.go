package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	// Do executes an HTTP request
	Do(ctx context.Context, req *Request) (*Response, error)

	// Get executes a GET request
	Get(ctx context.Context, url string, headers map[string]string) (*Response, error)

	// Post executes a POST request
	Post(ctx context.Context, url string, body any, headers map[string]string) (*Response, error)

	// Put executes a PUT request
	Put(ctx context.Context, url string, body any, headers map[string]string) (*Response, error)

	// Delete executes a DELETE request
	Delete(ctx context.Context, url string, headers map[string]string) (*Response, error)

	// Patch executes a PATCH request
	Patch(ctx context.Context, url string, body any, headers map[string]string) (*Response, error)
}

// client is the private implementation
type client struct {
	httpClient     *http.Client
	defaultHeaders map[string]string
	maxRetries     int
	retryDelay     time.Duration
}

// Compile-time interface compliance check
var _ Client = (*client)(nil)

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(cfg Config) Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 0
	}

	if cfg.RetryDelay == 0 {
		cfg.RetryDelay = 1 * time.Second
	}

	return &client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		defaultHeaders: cfg.DefaultHeaders,
		maxRetries:     cfg.MaxRetries,
		retryDelay:     cfg.RetryDelay,
	}
}

// Do executes an HTTP request with retry logic
func (c *client) Do(ctx context.Context, req *Request) (*Response, error) {
	if err := c.validateRequest(req); err != nil {
		return nil, err
	}

	var lastErr error
	attempts := c.maxRetries + 1

	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-time.After(c.retryDelay):
			case <-ctx.Done():
				return nil, ErrTimeout.WithError(ctx.Err())
			}
		}

		resp, err := c.doRequest(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

// doRequest executes a single HTTP request
func (c *client) doRequest(ctx context.Context, req *Request) (*Response, error) {
	// Marshal request body
	bodyBytes, err := c.marshalBody(req.Body)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, ErrInvalidRequest.WithError(err)
	}

	// Set default headers
	for key, value := range c.defaultHeaders {
		httpReq.Header.Set(key, value)
	}

	// Set request-specific headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, ErrTimeout.WithError(err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, ErrRequestFailed.WithError(err)
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       body,
	}, nil
}

// validateRequest validates the HTTP request
func (c *client) validateRequest(req *Request) error {
	if req == nil {
		return ErrInvalidRequest.WithError(fmt.Errorf("request is nil"))
	}

	if req.Method == "" {
		return ErrInvalidRequest.WithError(fmt.Errorf("method is required"))
	}

	if req.URL == "" {
		return ErrInvalidRequest.WithError(fmt.Errorf("URL is required"))
	}

	return nil
}

// Get executes a GET request
func (c *client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	})
}

// Post executes a POST request
func (c *client) Post(ctx context.Context, url string, body any, headers map[string]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method:  http.MethodPost,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Put executes a PUT request
func (c *client) Put(ctx context.Context, url string, body any, headers map[string]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method:  http.MethodPut,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Delete executes a DELETE request
func (c *client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
	})
}

// Patch executes a PATCH request
func (c *client) Patch(ctx context.Context, url string, body any, headers map[string]string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method:  http.MethodPatch,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// marshalBody marshals body to JSON bytes
func (c *client) marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	// If already []byte, return as is
	if bytes, ok := body.([]byte); ok {
		return bytes, nil
	}

	// If string, convert to []byte
	if str, ok := body.(string); ok {
		return []byte(str), nil
	}

	// Otherwise, marshal to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, ErrInvalidRequest.WithError(fmt.Errorf("failed to marshal body: %w", err))
	}

	return bodyBytes, nil
}
