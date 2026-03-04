// Package service 租户领域服务
package service

import (
	"context"

	"github.com/google/uuid"

	sharedEntity "go-ddd-scaffold/internal/domain/shared/entity"
	tenantEntity "go-ddd-scaffold/internal/domain/tenant/entity"
)

// MembershipDomainService 成员关系领域服务接口
type MembershipDomainService interface {
	// ValidateMemberLimit 验证租户成员数量限制
	ValidateMemberLimit(tenant *tenantEntity.Tenant, currentCount int) error
	// CanUserJoinTenant 检查用户是否可以加入租户
	CanUserJoinTenant(ctx context.Context, userID, tenantID uuid.UUID, role sharedEntity.UserRole) bool
	// ValidateRoleTransition 验证角色转换是否合法（如从 member 升级到 admin）
	ValidateRoleTransition(currentRole, newRole sharedEntity.UserRole) error
}

// membershipDomainService 成员关系领域服务实现
type membershipDomainService struct{}

// NewMembershipDomainService 创建成员关系领域服务实例
func NewMembershipDomainService() MembershipDomainService {
	return &membershipDomainService{}
}

// ValidateMemberLimit 验证租户成员数量限制
func (s *membershipDomainService) ValidateMemberLimit(tenant *tenantEntity.Tenant, currentCount int) error {
	if !tenant.IsValid() {
		return tenantEntity.ErrTenantInvalid
	}

	if !tenant.CanAddMoreMembers(currentCount) {
		return tenantEntity.ErrTenantMemberLimitExceeded
	}

	return nil
}

// CanUserJoinTenant 检查用户是否可以加入租户
func (s *membershipDomainService) CanUserJoinTenant(ctx context.Context, userID, tenantID uuid.UUID, role sharedEntity.UserRole) bool {
	// 基础验证：用户 ID 和租户 ID 不能为空
	if userID == uuid.Nil || tenantID == uuid.Nil {
		return false
	}

	// 角色验证：必须是有效的角色
	switch role {
	case sharedEntity.RoleOwner, sharedEntity.RoleAdmin, sharedEntity.RoleMember, sharedEntity.RoleGuest:
		// 有效角色
	default:
		return false
	}

	// 更复杂的业务规则可以在这里扩展
	// 例如：检查用户是否已被邀请、是否有黑名单记录等

	return true
}

// ValidateRoleTransition 验证角色转换是否合法
func (s *membershipDomainService) ValidateRoleTransition(currentRole, newRole sharedEntity.UserRole) error {
	// Owner 角色不能被移除或降级
	if currentRole == sharedEntity.RoleOwner {
		return ErrCannotChangeOwnerRole
	}

	// 不允许直接跳到 Owner
	if newRole == sharedEntity.RoleOwner {
		return ErrCannotPromoteToOwner
	}

	// 其他情况都是允许的
	return nil
}
