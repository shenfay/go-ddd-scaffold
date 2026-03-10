// Package dto 租户应用层 DTO 定义
package dto

import (
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/tenant/factory"
	user_entity"go-ddd-scaffold/internal/domain/user/entity"
	cast "go-ddd-scaffold/pkg/uitl"
)

// CreateTenantRequest 创建租户请求 DTO
type CreateTenantRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description,omitempty" binding:"max=500"`
	MaxMembers  int    `json:"maxMembers" binding:"required,min=1"`
}

// TenantResponse 租户响应 DTO
type TenantResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	MaxMembers  int      `json:"maxMembers"`
	MemberCount int64     `json:"memberCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TenantWithRoleResponse 租户带角色响应 DTO（用于租户列表）
type TenantWithRoleResponse struct {
	TenantResponse
	Role     string    `json:"role"`     // 用户在该租户的角色
	JoinedAt time.Time `json:"joinedAt"` // 加入时间
}

// ToTenantDTO 将领域实体转换为租户响应 DTO
func ToTenantDTO(entity *user_entity.Tenant, memberCount int64) *TenantResponse {
	if entity == nil {
		return nil
	}

	return &TenantResponse{
		ID:          entity.ID.String(),
		Name:        entity.Name,
		Description: &entity.Description,
		MaxMembers:  entity.MaxMembers,
		MemberCount: memberCount,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToTenantWithRoleDTO 将租户和成员关系转换为带角色的 DTO
func ToTenantWithRoleDTO(tenant *user_entity.Tenant, membership *user_entity.TenantMember, memberCount int64) *TenantWithRoleResponse {
	if tenant == nil || membership == nil {
		return nil
	}

	base := ToTenantDTO(tenant, memberCount)
	return &TenantWithRoleResponse{
		TenantResponse: *base,
		Role:           string(membership.Role),
		JoinedAt:       membership.JoinedAt,
	}
}

// FromCreateRequest 从创建请求转换为领域实体
func FromCreateRequest(req *CreateTenantRequest, ownerID uuid.UUID) (*user_entity.Tenant, error) {
	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	return factory.NewTenantBuilder(ownerID, req.Name).
		WithDescription(description).
		WithMaxMembers(req.MaxMembers).
		Build()
}

// UserFromDTO 将 DTO 转换为实体（用于更新，仅包含基础字段）
func UserFromDTO(dto *TenantResponse) *user_entity.Tenant {
	if dto == nil {
		return nil
	}

	id, _ := cast.ToUUID(dto.ID)
	description := ""
	if dto.Description != nil {
		description = *dto.Description
	}

	tenant := &user_entity.Tenant{
		ID:          id,
		Name:        dto.Name,
		Description: description,
		MaxMembers:  dto.MaxMembers,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}

	return tenant
}
