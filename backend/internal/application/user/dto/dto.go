// Package dto 用户应用层DTO定义
package dto

import (
	"time"

	"github.com/google/uuid"

	user_entity "go-ddd-scaffold/internal/domain/user/entity"
)

// RegisterRequest 注册请求DTO
type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required,oneof=MEMBER GUEST TEACHER ADMIN"`
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

// CreateTenantRequest 创建租户请求DTO
type CreateTenantRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	MaxMembers  int     `json:"maxMembers" binding:"min=0"`
}

// User 用户 DTO
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
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

// ToUserDTO 将实体转换为User DTO
func ToUserDTO(entity *user_entity.User) *User {
	if entity == nil {
		return nil
	}

	var tenantID *string
	if entity.TenantID != nil {
		id := entity.TenantID.String()
		tenantID = &id
	}

	return &User{
		ID:        entity.ID.String(),
		Email:     entity.Email,
		Role:      string(entity.Role),
		TenantID:  tenantID,
		Status:    string(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToTenantDTO 将实体转换为Tenant DTO
func ToTenantDTO(entity *user_entity.Tenant, userCount int64) *Tenant {
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

// UserFromDTO 将DTO转换为实体（用于更新）
func UserFromDTO(dto *User) *user_entity.User {
	if dto == nil {
		return nil
	}

	id, _ := uuid.Parse(dto.ID)
	user := &user_entity.User{
		ID:        id,
		Email:     dto.Email,
		Role:      user_entity.UserRole(dto.Role),
		Status:    user_entity.UserStatus(dto.Status),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	if dto.TenantID != nil {
		tenantID, _ := uuid.Parse(*dto.TenantID)
		user.TenantID = &tenantID
	}

	return user
}
