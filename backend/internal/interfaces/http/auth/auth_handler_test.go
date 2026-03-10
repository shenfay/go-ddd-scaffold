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
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MockAuthenticationService 模拟认证服务
type MockAuthenticationService struct{}

func (m *MockAuthenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	if req.Email == "existing@example.com" {
		return nil, errors.ErrUserExists
	}
	return &dto.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Nickname: req.Nickname,
	}, nil
}

func (m *MockAuthenticationService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if req.Email != "test@example.com" || req.Password != "password123" {
		return nil, errors.ErrInvalidPassword
	}
	return &dto.LoginResponse{
		AccessToken: "mock_jwt_token",
		User: &dto.User{
			ID:       uuid.New().String(),
			Email:    req.Email,
			Nickname: "Test User",
		},
	}, nil
}

func (m *MockAuthenticationService) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	return nil
}

func setupAuthRouter(authService service.AuthenticationService) *gin.Engine {
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

func TestAuthHandler_Register(t *testing.T) {
	router := setupAuthRouter(&MockAuthenticationService{})

	tests := []struct {
		name           string
		request        dto.RegisterRequest
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "成功注册",
			request: dto.RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				Nickname:  "New User",
				TenantID:  strPtr(uuid.New().String()),
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "Success",
		},
		{
			name: "邮箱已存在",
			request: dto.RegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				Nickname:  "Existing User",
				TenantID:  strPtr(uuid.New().String()),
			},
			expectedStatus: http.StatusConflict,
			expectedCode:   "User.Exists",
		},
		{
			name: "无效邮箱格式",
			request: dto.RegisterRequest{
				Email:     "invalid-email",
				Password:  "password123",
				Nickname:  "User",
				TenantID:  strPtr(uuid.New().String()),
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "ValidationFailed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var resp response.Response
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.expectedCode, resp.Code)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	router := setupAuthRouter(&MockAuthenticationService{})

	tests := []struct {
		name           string
		request        dto.LoginRequest
		expectedStatus int
		expectedCode   string
		expectToken    bool
	}{
		{
			name: "登录成功",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "Success",
			expectToken:    true,
		},
		{
			name: "密码错误",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "User.InvalidPassword",
			expectToken:    false,
		},
		{
			name: "参数验证失败",
			request: dto.LoginRequest{
				Email:    "",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "ValidationFailed",
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var resp response.Response
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.expectedCode, resp.Code)
			
			if tt.expectToken {
				assert.NotNil(t, resp.Data)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
