package specification

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// ActiveMemberSpec 活跃成员规格
type ActiveMemberSpec struct{}

func (s *ActiveMemberSpec) IsSatisfiedBy(candidate *entity.TenantMember) bool {
	return candidate.Status == entity.MemberStatusActive
}

// OwnerRoleSpec 所有者角色规格
type OwnerRoleSpec struct{}

func (s *OwnerRoleSpec) IsSatisfiedBy(candidate *entity.TenantMember) bool {
	return candidate.Role == entity.RoleOwner
}

// MemberInStatusSpec 指定状态成员规格
type MemberInStatusSpec struct{
	status entity.MemberStatus
}

func NewMemberInStatusSpec(status entity.MemberStatus) *MemberInStatusSpec {
	return &MemberInStatusSpec{status: status}
}

func (s *MemberInStatusSpec) IsSatisfiedBy(candidate *entity.TenantMember) bool {
	return candidate.Status == s.status
}

// MemberHasRoleSpec 指定角色成员规格
type MemberHasRoleSpec struct{
	role entity.UserRole
}

func NewMemberHasRoleSpec(role entity.UserRole) *MemberHasRoleSpec {
	return &MemberHasRoleSpec{role: role}
}

func (s *MemberHasRoleSpec) IsSatisfiedBy(candidate *entity.TenantMember) bool {
	return candidate.Role == s.role
}

// TenantNotExpiredSpec 租户未过期规格
type TenantNotExpiredSpec struct{}

func (s *TenantNotExpiredSpec) IsSatisfiedBy(candidate *entity.Tenant) bool {
	return candidate.ExpiredAt.After(time.Now())
}

// TenantHasCapacitySpec 租户有容量规格
type TenantHasCapacitySpec struct {
	currentMembers int64
}

func NewTenantHasCapacitySpec(currentMembers int64) *TenantHasCapacitySpec {
	return &TenantHasCapacitySpec{currentMembers: currentMembers}
}

func (s *TenantHasCapacitySpec) IsSatisfiedBy(candidate *entity.Tenant) bool {
	return int64(candidate.MaxMembers) > s.currentMembers
}

// AdminRoleSpec 管理员角色规格
type AdminRoleSpec struct{}

func (s *AdminRoleSpec) IsSatisfiedBy(candidate *entity.TenantMember) bool {
	return candidate.Role == entity.RoleAdmin
}
