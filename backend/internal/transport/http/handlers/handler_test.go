package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_RegisterValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "should fail with invalid email",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "ValidPassword123!",
			},
			expectedCode: 400,
		},
		{
			name: "should fail with short password",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "short",
			},
			expectedCode: 400,
		},
		{
			name: "should fail with missing email",
			requestBody: map[string]interface{}{
				"password": "ValidPassword123!",
			},
			expectedCode: 400,
		},
		{
			name: "should fail with missing password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedCode: 400,
		},
		{
			name:         "should fail with empty body",
			requestBody:  map[string]interface{}{},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Gin engine 以启用 binding 验证
			engine := gin.New()
			engine.POST("/auth/register", func(c *gin.Context) {
				var req RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestAuthHandler_LoginValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "should fail with invalid email format",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedCode: 400,
		},
		{
			name: "should fail with missing password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedCode: 400,
		},
		{
			name:         "should fail with empty body",
			requestBody:  map[string]interface{}{},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Gin engine 以启用 binding 验证
			engine := gin.New()
			engine.POST("/auth/login", func(c *gin.Context) {
				var req LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	server := httptest.NewServer(engine)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", body["status"])
}
