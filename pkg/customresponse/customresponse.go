package customresponse

import "github.com/gin-gonic/gin"

// CustomResponse represents a standardized API response
type CustomResponse[T any] struct {
	Success bool       `json:"success"`
	Data    T          `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success sends a successful response
func Success[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, CustomResponse[T]{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, CustomResponse[any]{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}
