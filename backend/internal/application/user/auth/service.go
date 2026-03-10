package auth

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
)

// AuthenticationService 认证服务接口
type AuthenticationService interface {
	// Register 用户注册
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error)
	
	// Login 用户登录
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	
	// Logout 用户登出（将 token 加入黑名单）
	Logout(ctx context.Context, userID uuid.UUID, token string) error
}
