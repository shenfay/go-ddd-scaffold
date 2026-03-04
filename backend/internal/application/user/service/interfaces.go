// Package service 用户应用服务接口定义
package service

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
	// Logout 用户登出
	Logout(ctx context.Context, userID uuid.UUID) error
}

// UserQueryService 用户查询服务接口
type UserQueryService interface {
	// GetUser 根据 ID 获取用户
	GetUser(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	// GetUserInfo 获取当前用户信息
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	// ListUsersByTenant 列出租户下所有用户
	ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error)
	// ListMembersByTenant 列出租户下的所有成员
	ListMembersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error)
}

// UserCommandService 用户命令服务接口
type UserCommandService interface {
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) error
	// UpdateProfile 更新用户个人资料
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

// TenantService 租户服务接口
type TenantService interface {
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.Tenant, error)
	// GetTenant 获取租户信息
	GetTenant(ctx context.Context, tenantID uuid.UUID) (*dto.Tenant, error)
}
