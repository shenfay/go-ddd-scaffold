package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-ddd-scaffold/internal/domain/tenant/factory"
	user_entity"go-ddd-scaffold/internal/domain/user/entity"
	user_spec "go-ddd-scaffold/internal/domain/user/specification"
	shared_spec "go-ddd-scaffold/internal/domain/shared/specification"
)

func TestTenantLifecycle(t *testing.T) {
	ownerID := uuid.New()

	t.Run("CreateTenantWithBuilder", func(t *testing.T) {
		tenant := factory.NewTenantBuilder(ownerID, "测试租户").
			WithDescription("描述信息").
			WithMaxMembers(50).
			MustBuild()

		require.NotNil(t, tenant)
		assert.Equal(t, "测试租户", tenant.Name)
		assert.Equal(t, 50, tenant.MaxMembers)
		assert.True(t, tenant.ExpiredAt.After(time.Now()))
	})

	t.Run("TenantCapacityCheck", func(t *testing.T) {
		tenant := factory.NewTenantBuilder(ownerID, "容量测试租户").MustBuild()

		capacitySpec := user_spec.NewTenantHasCapacitySpec(0)
		assert.True(t, capacitySpec.IsSatisfiedBy(tenant))

		capacitySpec = user_spec.NewTenantHasCapacitySpec(int64(tenant.MaxMembers - 1))
		assert.True(t, capacitySpec.IsSatisfiedBy(tenant))

		capacitySpec = user_spec.NewTenantHasCapacitySpec(int64(tenant.MaxMembers))
		assert.False(t, capacitySpec.IsSatisfiedBy(tenant))
	})

	t.Run("MemberStatusCheck", func(t *testing.T) {
		activeMember := &user_entity.TenantMember{Status: user_entity.MemberStatusActive}
		inactiveMember := &user_entity.TenantMember{Status: user_entity.MemberStatusInactive}

		activeSpec := &user_spec.ActiveMemberSpec{}
		assert.True(t, activeSpec.IsSatisfiedBy(activeMember))
		assert.False(t, activeSpec.IsSatisfiedBy(inactiveMember))
	})

	t.Run("CombinedSpecifications", func(t *testing.T) {
		activeOwner := &user_entity.TenantMember{
			Status: user_entity.MemberStatusActive,
			Role:   user_entity.RoleOwner,
		}

		activeMember := &user_entity.TenantMember{
			Status: user_entity.MemberStatusActive,
			Role:   user_entity.RoleMember,
		}

		combinedSpec := shared_spec.And(&user_spec.ActiveMemberSpec{}, &user_spec.OwnerRoleSpec{})
		assert.True(t, combinedSpec.IsSatisfiedBy(activeOwner))
		assert.False(t, combinedSpec.IsSatisfiedBy(activeMember))
	})
}
