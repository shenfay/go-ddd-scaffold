package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	http2 "go-ddd-scaffold/internal/interfaces/http/auth"
	"go-ddd-scaffold/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MockAuthBenchmarkService 模拟认证服务（用于基准测试）
type MockAuthBenchmarkService struct{}

func (m *MockAuthBenchmarkService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	return &dto.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Nickname: req.Nickname,
	}, nil
}

func (m *MockAuthBenchmarkService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	return &dto.LoginResponse{
		AccessToken: "mock_jwt_token_" + uuid.New().String(),
		User: &dto.User{
			ID:       uuid.New().String(),
			Email:    req.Email,
			Nickname: "Test User",
		},
	}, nil
}

func (m *MockAuthBenchmarkService) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	return nil
}


func setupBenchmarkRouter(authService service.AuthenticationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	authHandler := http2.NewAuthHandler(authService, zap.NewNop(), nil)
	
	// 注册中间件
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.HandlePanic())
	
	// 注册路由
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}
	
	return router
}

func BenchmarkAuthHandler_Register(b *testing.B) {
	router := setupBenchmarkRouter(&MockAuthBenchmarkService{})
	
	reqBody := dto.RegisterRequest{
		Email:     "benchmark@example.com",
		Password:  "password123",
		Nickname:  "Benchmark User",
		TenantID:  strPtr(uuid.New().String()),
	}
	body, _ := json.Marshal(reqBody)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			b.Fatalf("期望状态码 200, 得到 %d", w.Code)
		}
	}
}

func BenchmarkAuthHandler_Login(b *testing.B) {
	router := setupBenchmarkRouter(&MockAuthBenchmarkService{})
	
	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			b.Fatalf("期望状态码 200, 得到 %d", w.Code)
		}
	}
}
