package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
)

// TestServer 测试服务器封装
type TestServer struct {
	Engine       *gin.Engine
	AuthService  *authentication.Service
	TokenService authentication.TokenService
	Recorder     *httptest.ResponseRecorder
}

// NewTestServer 创建测试服务器
func NewTestServer(t *testing.T, authService *authentication.Service, tokenService authentication.TokenService) *TestServer {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	// 设置路由
	engine.POST("/api/v1/auth/register", authHandler.Register)
	engine.POST("/api/v1/auth/login", authHandler.Login)
	engine.POST("/api/v1/auth/logout", authHandler.Logout)
	engine.POST("/api/v1/auth/refresh", authHandler.RefreshToken)
	engine.POST("/api/v1/auth/forgot-password", authHandler.RequestPasswordReset)
	engine.POST("/api/v1/auth/reset-password", authHandler.ResetPassword)
	engine.POST("/api/v1/auth/verify-email", authHandler.VerifyEmail)
	engine.POST("/api/v1/auth/resend-verification", authHandler.ResendVerificationEmail)
	engine.GET("/api/v1/auth/me", middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: tokenService,
	}), authHandler.GetCurrentUser)

	return &TestServer{
		Engine:       engine,
		AuthService:  authService,
		TokenService: tokenService,
		Recorder:     httptest.NewRecorder(),
	}
}

// PerformRequest 执行 HTTP 请求
func (ts *TestServer) PerformRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	ts.Recorder = httptest.NewRecorder()
	req := httptest.NewRequest(method, url, nil)

	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	}

	// 添加自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ts.Engine.ServeHTTP(ts.Recorder, req)
	return ts.Recorder
}

// ParseResponse 解析响应
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), target)
	assert.NoError(t, err)
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// ExtractToken 从响应中提取 Token
func ExtractToken(w *httptest.ResponseRecorder, tokenType string) string {
	var resp AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		return ""
	}

	if tokenData, ok := resp.Data[tokenType].(string); ok {
		return tokenData
	}
	return ""
}

// ExtractField 从响应中提取字段
func ExtractField(w *httptest.ResponseRecorder, field string) string {
	var resp AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		return ""
	}

	if data, ok := resp.Data[field].(string); ok {
		return data
	}
	if data, ok := resp.Data[field].(map[string]interface{}); ok {
		jsonBytes, _ := json.Marshal(data)
		return string(jsonBytes)
	}
	return ""
}

// GenerateTestEmail 生成测试邮箱
func GenerateTestEmail(testName string, index int) string {
	return fmt.Sprintf("test.%s.%d@example.com", testName, index)
}

// GenerateTestPassword 生成测试密码
func GenerateTestPassword() string {
	return "TestPassword123!"
}

// TestMetrics 创建测试指标(空实现)
func TestMetrics() *metrics.Metrics {
	return nil
}
