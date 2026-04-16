package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/test/factory"
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

func TestLoginRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should validate email format", func(t *testing.T) {
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
		reqBody := map[string]interface{}{
			"email":    "not-an-email",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should require password", func(t *testing.T) {
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
		reqBody := map[string]interface{}{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestRefreshTokenValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should require refresh token", func(t *testing.T) {
		engine := gin.New()
		engine.POST("/auth/refresh", func(c *gin.Context) {
			var req RefreshTokenRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		reqBody := map[string]interface{}{}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestRegisterRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	f := factory.NewUserFactory()

	t.Run("should accept valid registration request", func(t *testing.T) {
		engine := gin.New()
		engine.POST("/auth/register", func(c *gin.Context) {
			var req RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(201, gin.H{"status": "ok", "email": req.Email})
		})

		w := httptest.NewRecorder()
		reqBody := map[string]interface{}{
			"email":    f.CreateUser().Email,
			"password": "ValidPassword123!",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})
}

func TestPasswordResetRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	f := factory.NewUserFactory()

	t.Run("should validate email for password reset", func(t *testing.T) {
		engine := gin.New()
		engine.POST("/auth/password-reset", func(c *gin.Context) {
			var req struct {
				Email string `json:"email" binding:"required,email"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		reqBody := map[string]interface{}{
			"email": f.CreateUser().Email,
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/password-reset", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should reject invalid email format", func(t *testing.T) {
		engine := gin.New()
		engine.POST("/auth/password-reset", func(c *gin.Context) {
			var req struct {
				Email string `json:"email" binding:"required,email"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		reqBody := map[string]interface{}{
			"email": "not-an-email",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/password-reset", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}
