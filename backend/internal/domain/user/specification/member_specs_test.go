package specification_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/domain/user/entity"
	user_spec "go-ddd-scaffold/internal/domain/user/specification"
	shared_spec "go-ddd-scaffold/internal/domain/shared/specification"
)

func TestActiveMemberSpec(t *testing.T) {
	spec := &user_spec.ActiveMemberSpec{}

	activeMember := &entity.TenantMember{Status: entity.MemberStatusActive}
	inactiveMember := &entity.TenantMember{Status: entity.MemberStatusInactive}
	removedMember := &entity.TenantMember{Status: entity.MemberStatusRemoved}

	assert.True(t, spec.IsSatisfiedBy(activeMember), "活跃成员应满足规格")
	assert.False(t, spec.IsSatisfiedBy(inactiveMember), "非活跃成员不应满足规格")
	assert.False(t, spec.IsSatisfiedBy(removedMember), "已移除成员不应满足规格")
}

func TestOwnerRoleSpec(t *testing.T) {
	spec := &user_spec.OwnerRoleSpec{}

	owner := &entity.TenantMember{Role: entity.RoleOwner}
	admin := &entity.TenantMember{Role: entity.RoleAdmin}
	member := &entity.TenantMember{Role: entity.RoleMember}

	assert.True(t, spec.IsSatisfiedBy(owner), "所有者应满足规格")
	assert.False(t, spec.IsSatisfiedBy(admin), "管理员不应满足规格")
	assert.False(t, spec.IsSatisfiedBy(member), "普通成员不应满足规格")
}

func TestMemberInStatusSpec(t *testing.T) {
	spec := user_spec.NewMemberInStatusSpec(entity.MemberStatusActive)

	activeMember := &entity.TenantMember{Status: entity.MemberStatusActive}
	inactiveMember := &entity.TenantMember{Status: entity.MemberStatusInactive}

	assert.True(t, spec.IsSatisfiedBy(activeMember), "状态匹配应满足规格")
	assert.False(t, spec.IsSatisfiedBy(inactiveMember), "状态不匹配不应满足规格")
}

func TestMemberHasRoleSpec(t *testing.T) {
	spec := user_spec.NewMemberHasRoleSpec(entity.RoleAdmin)

	admin := &entity.TenantMember{Role: entity.RoleAdmin}
	owner := &entity.TenantMember{Role: entity.RoleOwner}

	assert.True(t, spec.IsSatisfiedBy(admin), "角色匹配应满足规格")
	assert.False(t, spec.IsSatisfiedBy(owner), "角色不匹配不应满足规格")
}

func TestTenantNotExpiredSpec(t *testing.T) {
	spec := &user_spec.TenantNotExpiredSpec{}

	futureTenant := &entity.Tenant{ExpiredAt: time.Now().AddDate(1, 0, 0)}
	expiredTenant := &entity.Tenant{ExpiredAt: time.Now().AddDate(-1, 0, 0)}

	assert.True(t, spec.IsSatisfiedBy(futureTenant), "未过期租户应满足规格")
	assert.False(t, spec.IsSatisfiedBy(expiredTenant), "已过期租户不应满足规格")
}

func TestTenantHasCapacitySpec(t *testing.T) {
	spec := user_spec.NewTenantHasCapacitySpec(5)

	tenantWithCapacity := &entity.Tenant{MaxMembers: 10}
	tenantFull := &entity.Tenant{MaxMembers: 5}
	tenantOverCapacity := &entity.Tenant{MaxMembers: 3}

	assert.True(t, spec.IsSatisfiedBy(tenantWithCapacity), "有容量应满足规格")
	assert.False(t, spec.IsSatisfiedBy(tenantFull), "容量已满不应满足规格")
	assert.False(t, spec.IsSatisfiedBy(tenantOverCapacity), "超容量不应满足规格")
}

func TestAndSpec(t *testing.T) {
	activeOwnerSpec := shared_spec.And(&user_spec.ActiveMemberSpec{}, &user_spec.OwnerRoleSpec{})

	activeOwner := &entity.TenantMember{Status: entity.MemberStatusActive, Role: entity.RoleOwner}
	activeMember := &entity.TenantMember{Status: entity.MemberStatusActive, Role: entity.RoleMember}
	inactiveOwner := &entity.TenantMember{Status: entity.MemberStatusInactive, Role: entity.RoleOwner}

	assert.True(t, activeOwnerSpec.IsSatisfiedBy(activeOwner), "活跃的所有者应满足与规格")
	assert.False(t, activeOwnerSpec.IsSatisfiedBy(activeMember), "活跃的非所有者不应满足与规格")
	assert.False(t, activeOwnerSpec.IsSatisfiedBy(inactiveOwner), "非活跃的所有者不应满足与规格")
}

func TestOrSpec(t *testing.T) {
	adminOrOwnerSpec := shared_spec.Or(&user_spec.AdminRoleSpec{}, &user_spec.OwnerRoleSpec{})

	admin := &entity.TenantMember{Role: entity.RoleAdmin}
	owner := &entity.TenantMember{Role: entity.RoleOwner}
	member := &entity.TenantMember{Role: entity.RoleMember}

	assert.True(t, adminOrOwnerSpec.IsSatisfiedBy(admin), "管理员应满足或规格")
	assert.True(t, adminOrOwnerSpec.IsSatisfiedBy(owner), "所有者应满足或规格")
	assert.False(t, adminOrOwnerSpec.IsSatisfiedBy(member), "普通成员不应满足或规格")
}

func TestNotSpec(t *testing.T) {
	notActiveSpec := shared_spec.Not(&user_spec.ActiveMemberSpec{})

	activeMember := &entity.TenantMember{Status: entity.MemberStatusActive}
	inactiveMember := &entity.TenantMember{Status: entity.MemberStatusInactive}

	assert.False(t, notActiveSpec.IsSatisfiedBy(activeMember), "活跃成员不应满足非规格")
	assert.True(t, notActiveSpec.IsSatisfiedBy(inactiveMember), "非活跃成员应满足非规格")
}
