package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/srcndev/message-service/pkg/customerror"
	"github.com/srcndev/message-service/pkg/response"
)

// ErrorHandler is a middleware that handles errors from handlers
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *customerror.CustomError
			if errors.As(err, &appErr) {
				response.Error(c, appErr.GetStatusCode(), appErr.Code, appErr.Message)
				return
			}

			// Fallback for unknown errors
			response.Error(c, 500, "INTERNAL_ERROR", "Internal server error")
		}
	}
}
