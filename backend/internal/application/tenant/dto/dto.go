// Package dto 租户应用层 DTO 定义
package dto

import (
	"time"

	"github.com/google/uuid"

	user_dto "go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	cast "go-ddd-scaffold/pkg/uitl"
)

// TenantMemberResponse 租户成员响应 DTO
type TenantMemberResponse struct{
	ID    string  `json:"id"`
	UserID  string  `json:"userId"`
	User    *user_dto.User   `json:"user,omitempty"`
	TenantID  string  `json:"tenantId"`
	Role    string  `json:"role"`
	Status  string  `json:"status"`
	JoinedAt  time.Time `json:"joinedAt"`
	LeftAt    *time.Time `json:"leftAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TenantResponse 租户响应 DTO
type TenantResponse struct{
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	MaxMembers int       `json:"maxMembers"`
	MemberCount int64     `json:"memberCount"`
	ExpiredAt   time.Time `json:"expiredAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TenantDetailResponse 租户详情响应 DTO（包含成员列表）
type TenantDetailResponse struct{
	TenantResponse
	Members []*TenantMemberResponse `json:"members"`
}

// CreateTenantRequest 创建租户请求 DTO
type CreateTenantRequest struct{
	Name      string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	MaxMembers int     `json:"maxMembers" binding:"required,min=1"`
}

// UpdateTenantRequest 更新租户请求 DTO
type UpdateTenantRequest struct{
	Name       *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	MaxMembers  *int    `json:"maxMembers,omitempty" binding:"omitempty,min=1"`
}

// ToTenantDTO 将租户实体转换为 Tenant DTO
func ToTenantDTO(tenant *entity.Tenant, memberCount int64) *TenantResponse {
	if tenant == nil {
		return nil
	}

	return &TenantResponse{
		ID:          tenant.ID.String(),
		Name:        tenant.Name,
		Description: cast.ToStringPtr(tenant.Description),
		MaxMembers:  tenant.MaxMembers,
		MemberCount: memberCount,
		ExpiredAt:   tenant.ExpiredAt,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
}

// ToTenantMemberDTO 将租户成员实体转换为 DTO
func ToTenantMemberDTO(member *entity.TenantMember, user *entity.User) *TenantMemberResponse {
	if member == nil {
		return nil
	}

	dto := &TenantMemberResponse{
		ID:        member.ID.String(),
		UserID:    member.UserID.String(),
		TenantID:  member.TenantID.String(),
		Role:     string(member.Role),
		Status:   string(member.Status),
		JoinedAt:  member.JoinedAt,
		LeftAt:    member.LeftAt,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}

	// 如果有用户信息，填充
	if user != nil {
		dto.User = user_dto.ToUserDTO(user)
	}

	return dto
}

// ToTenantDetailDTO 将租户实体和成员列表转换为详情 DTO
func ToTenantDetailDTO(tenant *entity.Tenant, members []*entity.TenantMember, users map[uuid.UUID]*entity.User) *TenantDetailResponse {
	if tenant == nil {
		return nil
	}

	response := &TenantDetailResponse{
		TenantResponse: *ToTenantDTO(tenant, int64(len(members))),
		Members:        make([]*TenantMemberResponse, 0, len(members)),
	}

	for _, member := range members {
		user := users[member.UserID]
		response.Members = append(response.Members, ToTenantMemberDTO(member, user))
	}

	return response
}
