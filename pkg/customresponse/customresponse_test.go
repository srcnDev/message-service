package customresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		want       CustomResponse
	}{
		{
			name:       "success with map data",
			statusCode: http.StatusOK,
			data: map[string]interface{}{
				"id":   1,
				"name": "test",
			},
			want: CustomResponse{
				Success: true,
				Data: map[string]interface{}{
					"id":   float64(1), // JSON unmarshals numbers as float64
					"name": "test",
				},
			},
		},
		{
			name:       "success with string data",
			statusCode: http.StatusCreated,
			data:       "created successfully",
			want: CustomResponse{
				Success: true,
				Data:    "created successfully",
			},
		},
		{
			name:       "success with nil data",
			statusCode: http.StatusOK,
			data:       nil,
			want: CustomResponse{
				Success: true,
				Data:    nil,
			},
		},
		{
			name:       "success with slice data",
			statusCode: http.StatusOK,
			data:       []string{"item1", "item2"},
			want: CustomResponse{
				Success: true,
				Data:    []interface{}{"item1", "item2"},
			},
		},
		{
			name:       "success with struct data",
			statusCode: http.StatusOK,
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   42,
				Name: "test-struct",
			},
			want: CustomResponse{
				Success: true,
				Data: map[string]interface{}{
					"id":   float64(42),
					"name": "test-struct",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Execute
			Success(c, tt.statusCode, tt.data)

			// Verify
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var got CustomResponse
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Success, got.Success)
			assert.Equal(t, tt.want.Data, got.Data)
			assert.Nil(t, got.Error)
		})
	}
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		statusCode int
		code       string
		message    string
		want       CustomResponse
	}{
		{
			name:       "bad request error",
			statusCode: http.StatusBadRequest,
			code:       "INVALID_INPUT",
			message:    "Invalid input provided",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "INVALID_INPUT",
					Message: "Invalid input provided",
				},
			},
		},
		{
			name:       "not found error",
			statusCode: http.StatusNotFound,
			code:       "NOT_FOUND",
			message:    "Resource not found",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "NOT_FOUND",
					Message: "Resource not found",
				},
			},
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			code:       "INTERNAL_ERROR",
			message:    "An internal error occurred",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "INTERNAL_ERROR",
					Message: "An internal error occurred",
				},
			},
		},
		{
			name:       "empty error code and message",
			statusCode: http.StatusBadRequest,
			code:       "",
			message:    "",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "",
					Message: "",
				},
			},
		},
		{
			name:       "unauthorized error",
			statusCode: http.StatusUnauthorized,
			code:       "UNAUTHORIZED",
			message:    "Authentication required",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "UNAUTHORIZED",
					Message: "Authentication required",
				},
			},
		},
		{
			name:       "conflict error",
			statusCode: http.StatusConflict,
			code:       "DUPLICATE_ENTRY",
			message:    "Resource already exists",
			want: CustomResponse{
				Success: false,
				Error: &ErrorInfo{
					Code:    "DUPLICATE_ENTRY",
					Message: "Resource already exists",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Execute
			Error(c, tt.statusCode, tt.code, tt.message)

			// Verify
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var got CustomResponse
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Success, got.Success)
			assert.Nil(t, got.Data)
			assert.NotNil(t, got.Error)
			assert.Equal(t, tt.want.Error.Code, got.Error.Code)
			assert.Equal(t, tt.want.Error.Message, got.Error.Message)
		})
	}
}

func TestCustomResponse_JSONStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success response omits error field", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Success(c, http.StatusOK, "test data")

		var raw map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &raw)
		assert.NoError(t, err)

		// Verify structure
		assert.Contains(t, raw, "success")
		assert.Contains(t, raw, "data")
		assert.NotContains(t, raw, "error")
	})

	t.Run("error response omits data field", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Error(c, http.StatusBadRequest, "ERROR_CODE", "error message")

		var raw map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &raw)
		assert.NoError(t, err)

		// Verify structure
		assert.Contains(t, raw, "success")
		assert.Contains(t, raw, "error")
		assert.NotContains(t, raw, "data")
	})
}

func TestErrorInfo_Structure(t *testing.T) {
	t.Run("error info has correct fields", func(t *testing.T) {
		errorInfo := ErrorInfo{
			Code:    "TEST_CODE",
			Message: "test message",
		}

		data, err := json.Marshal(errorInfo)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		assert.Equal(t, "TEST_CODE", result["code"])
		assert.Equal(t, "test message", result["message"])
	})
}

func TestCustomResponse_Structure(t *testing.T) {
	t.Run("response has correct structure", func(t *testing.T) {
		resp := CustomResponse{
			Success: true,
			Data:    "test",
		}

		data, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		assert.Equal(t, true, result["success"])
		assert.Equal(t, "test", result["data"])
	})

	t.Run("response with error has correct structure", func(t *testing.T) {
		resp := CustomResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "ERROR",
				Message: "error occurred",
			},
		}

		data, err := json.Marshal(resp)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		assert.Equal(t, false, result["success"])
		assert.NotNil(t, result["error"])

		errorMap := result["error"].(map[string]interface{})
		assert.Equal(t, "ERROR", errorMap["code"])
		assert.Equal(t, "error occurred", errorMap["message"])
	})
}

func TestSuccess_WithComplexData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("nested struct data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		type Address struct {
			City    string `json:"city"`
			Country string `json:"country"`
		}
		type User struct {
			ID      int     `json:"id"`
			Name    string  `json:"name"`
			Address Address `json:"address"`
		}

		data := User{
			ID:   1,
			Name: "John",
			Address: Address{
				City:    "Istanbul",
				Country: "Turkey",
			},
		}

		Success(c, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)

		var got CustomResponse
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.True(t, got.Success)
		assert.NotNil(t, got.Data)
	})
}

func TestError_WithVariousStatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	statusCodes := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}

	for _, code := range statusCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, code, "ERROR_CODE", "error message")

			assert.Equal(t, code, w.Code)

			var got CustomResponse
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.False(t, got.Success)
			assert.NotNil(t, got.Error)
		})
	}
}
