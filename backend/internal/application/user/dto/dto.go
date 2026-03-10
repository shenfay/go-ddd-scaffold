// Package dto 用户应用层 DTO 定义
package dto

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	cast "go-ddd-scaffold/pkg/uitl"
)

// RegisterRequest 注册请求DTO
type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Nickname string  `json:"nickname" binding:"required"`
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=member guest"` // 可选，默认为 member
	TenantID *string `json:"tenantId,omitempty"`
}

// LoginRequest 登录请求DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应DTO
type LoginResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"accessToken"`
}

// UpdateUserRequest 更新用户请求DTO
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Status   *string `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED"`
}

// UpdateProfileRequest 更新个人资料请求DTO（用户自己更新资料）
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}


// User 用户 DTO
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Phone     *string   `json:"phone"`
	Bio       *string   `json:"bio"`
	Avatar    *string   `json:"avatar"`
	Role      string    `json:"role"`
	TenantID  *string   `json:"tenantId,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Tenant 租户 DTO
type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	MaxMembers  int       `json:"maxMembers"`
	MemberCount int64     `json:"memberCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToUserDTO 将实体转换为 User DTO（基础信息，不含租户和角色）
func ToUserDTO(entity *entity.User) *User {
	if entity == nil {
		return nil
	}

	// 注意：不再返回 Role 和 TenantID
	// - 角色通过租户上下文动态获取（用户在每个租户可能有不同角色）
	// - 租户信息通过单独的接口获取

	return &User{
		ID:        entity.ID.String(),
		Email:     entity.Email.String(),
		Nickname:  entity.Nickname.String(),
		Phone:     entity.Phone,
		Bio:       entity.Bio,
		Avatar:    entity.Avatar,
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToTenantDTO 将实体转换为 Tenant DTO
func ToTenantDTO(entity *entity.Tenant, userCount int64) *Tenant {
	if entity == nil {
		return nil
	}

	return &Tenant{
		ID:          entity.ID.String(),
		Name:        entity.Name,
		Description: &entity.Description,
		MaxMembers:  entity.MaxMembers,
		MemberCount: userCount,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// UserFromDTO 将 DTO 转换为实体（用于更新，仅包含基础字段）
func UserFromDTO(dto *User) *entity.User {
	if dto == nil {
		return nil
	}

	id, _ := cast.ToUUID(dto.ID)
	email := valueobject.NewEmailFromString(dto.Email)
	nickname := valueobject.NewNicknameFromString(dto.Nickname)

	user := &entity.User{
		ID:        id,
		Email:     email,
		Nickname:  nickname,
		Status:    entity.UserStatus(dto.Status),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	// 注意：不再设置 Role 和 TenantID
	// - 角色通过 tenant_members 表管理
	// - 租户关系通过单独的接口处理

	return user
}
