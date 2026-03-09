// Package service_test 租户领域服务测试
package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/domain/shared/entity"
	tenantEntity "go-ddd-scaffold/internal/domain/tenant/entity"
	"go-ddd-scaffold/internal/domain/tenant/service"
)

// TestMembershipDomainService_ValidateMemberLimit_Success 测试验证成员限制成功
func TestMembershipDomainService_ValidateMemberLimit_Success(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	tenant := tenantEntity.NewTenant("Test Tenant", 10) // 最大 10 个成员
	currentCount := 5                                   // 当前 5 个成员

	// 2. 执行验证
	err := domainService.ValidateMemberLimit(tenant, currentCount)

	// 3. 验证结果
	assert.NoError(t, err)
}

// TestMembershipDomainService_ValidateMemberLimit_Exceeded 测试成员数超限
func TestMembershipDomainService_ValidateMemberLimit_Exceeded(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	tenant := tenantEntity.NewTenant("Test Tenant", 10)
	currentCount := 10 // 已达到上限

	// 2. 执行验证
	err := domainService.ValidateMemberLimit(tenant, currentCount)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
}

// TestMembershipDomainService_ValidateMemberLimit_InvalidTenant 测试无效租户
func TestMembershipDomainService_ValidateMemberLimit_InvalidTenant(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	// 创建已过期的租户
	tenant := &tenantEntity.Tenant{
		ID:         uuid.New(),
		Name:       "Expired Tenant",
		MaxMembers: 10,
		ExpiredAt:  tenantEntity.NewTenant("Test", 10).ExpiredAt.AddDate(-2, 0, 0), // 2 年前过期
	}

	// 2. 执行验证
	err := domainService.ValidateMemberLimit(tenant, 0)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantInvalid, err)
}

// TestMembershipDomainService_CanUserJoinTenant_Success 测试用户可以加入租户
func TestMembershipDomainService_CanUserJoinTenant_Success(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.New()

	// 2. 测试各种有效角色
	validRoles := []entity.UserRole{
		entity.RoleOwner,
		entity.RoleAdmin,
		entity.RoleMember,
		entity.RoleGuest,
	}

	for _, role := range validRoles {
		canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, role)
		assert.True(t, canJoin, "Role %s should be able to join", role)
	}
}

// TestMembershipDomainService_CanUserJoinTenant_InvalidUserID 测试无效用户 ID
func TestMembershipDomainService_CanUserJoinTenant_InvalidUserID(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	userID := uuid.Nil // 无效的用户 ID
	tenantID := uuid.New()

	// 2. 执行验证
	canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, entity.RoleMember)

	// 3. 验证结果
	assert.False(t, canJoin)
}

// TestMembershipDomainService_CanUserJoinTenant_InvalidTenantID 测试无效租户 ID
func TestMembershipDomainService_CanUserJoinTenant_InvalidTenantID(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.Nil // 无效的租户 ID

	// 2. 执行验证
	canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, entity.RoleMember)

	// 3. 验证结果
	assert.False(t, canJoin)
}

// TestMembershipDomainService_CanUserJoinTenant_InvalidRole 测试无效角色
func TestMembershipDomainService_CanUserJoinTenant_InvalidRole(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.New()

	// 2. 使用不存在的角色
	invalidRole := entity.UserRole("non_existent_role")

	// 3. 执行验证
	canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, invalidRole)

	// 4. 验证结果
	assert.False(t, canJoin)
}

// TestMembershipDomainService_ValidateRoleTransition_Success 测试角色转换成功
func TestMembershipDomainService_ValidateRoleTransition_Success(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	// 2. 测试有效的角色转换
	testCases := []struct {
		currentRole entity.UserRole
		newRole     entity.UserRole
		description string
	}{
		{entity.RoleMember, entity.RoleAdmin, "Member -> Admin"},
		{entity.RoleGuest, entity.RoleMember, "Guest -> Member"},
		{entity.RoleMember, entity.RoleGuest, "Member -> Guest (降级)"},
		{entity.RoleAdmin, entity.RoleMember, "Admin -> Member (降级)"},
	}

	for _, tc := range testCases {
		err := domainService.ValidateRoleTransition(tc.currentRole, tc.newRole)
		assert.NoError(t, err, tc.description)
	}
}

// TestMembershipDomainService_ValidateRoleTransition_CannotChangeOwner 测试不能修改 Owner 角色
func TestMembershipDomainService_ValidateRoleTransition_CannotChangeOwner(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	// 2. 尝试修改 Owner 角色
	err := domainService.ValidateRoleTransition(entity.RoleOwner, entity.RoleAdmin)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
}

// TestMembershipDomainService_ValidateRoleTransition_CannotPromoteToOwner 测试不能晋升为 Owner
func TestMembershipDomainService_ValidateRoleTransition_CannotPromoteToOwner(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	// 2. 尝试晋升为 Owner
	err := domainService.ValidateRoleTransition(entity.RoleMember, entity.RoleOwner)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, service.ErrCannotPromoteToOwner, err)
}

// TestTenant_AddMember_AggregateRootMethod 测试聚合根方法 AddMember
func TestTenant_AddMember_AggregateRootMethod(t *testing.T) {
	// 1. 创建租户
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	assert.Equal(t, 0, tenant.GetActiveMemberCount())

	// 2. 添加第一个成员
	userID := uuid.New()
	member, err := tenant.AddMember(userID, entity.RoleAdmin, nil)
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, 1, tenant.GetActiveMemberCount())

	// 3. 添加第二个成员
	userID2 := uuid.New()
	member2, err := tenant.AddMember(userID2, entity.RoleMember, &userID)
	assert.NoError(t, err)
	assert.NotNil(t, member2)
	assert.Equal(t, 2, tenant.GetActiveMemberCount())
}

// TestTenant_AddMember_DuplicateMember 测试重复添加成员
func TestTenant_AddMember_DuplicateMember(t *testing.T) {
	// 1. 创建租户并添加成员
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	userID := uuid.New()

	_, err := tenant.AddMember(userID, entity.RoleMember, nil)
	assert.NoError(t, err)

	// 2. 尝试再次添加同一用户
	_, err = tenant.AddMember(userID, entity.RoleMember, nil)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberAlreadyExists, err)
}

// TestTenant_AddMember_LimitExceeded 测试超过成员限制
func TestTenant_AddMember_LimitExceeded(t *testing.T) {
	// 1. 创建只有 2 个名额的租户
	tenant := tenantEntity.NewTenant("Small Tenant", 2)

	// 2. 添加 2 个成员
	_, _ = tenant.AddMember(uuid.New(), entity.RoleMember, nil)
	_, _ = tenant.AddMember(uuid.New(), entity.RoleMember, nil)

	// 3. 尝试添加第 3 个成员
	_, err := tenant.AddMember(uuid.New(), entity.RoleMember, nil)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
}

// TestTenant_RemoveMember_Success 测试移除成员成功
func TestTenant_RemoveMember_Success(t *testing.T) {
	// 1. 创建租户并添加成员
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	userID := uuid.New()
	member, _ := tenant.AddMember(userID, entity.RoleMember, nil)

	assert.Equal(t, 1, tenant.GetActiveMemberCount())

	// 2. 移除成员
	err := tenant.RemoveMember(member.ID)

	// 3. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 0, tenant.GetActiveMemberCount())
	assert.True(t, member.IsRemoved())
}

// TestTenant_RemoveMember_NotFound 测试移除不存在的成员
func TestTenant_RemoveMember_NotFound(t *testing.T) {
	// 1. 创建租户
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	nonExistentID := uuid.New()

	// 2. 尝试移除不存在的成员
	err := tenant.RemoveMember(nonExistentID)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberNotFound, err)
}

// TestMembershipDomainService_IntegrationWithUnitOfWork 测试领域服务与 UnitOfWork 集成
func TestMembershipDomainService_IntegrationWithUnitOfWork(t *testing.T) {
	// 1. 准备领域服务和 UnitOfWork
	domainService := service.NewMembershipDomainService()
	
	// 2. 创建租户
	tenant := tenantEntity.NewTenant("Integration Test Tenant", 3)
	
	// 3. 验证成员限制
	currentCount := 2
	err := domainService.ValidateMemberLimit(tenant, currentCount)
	assert.NoError(t, err)
	
	// 4. 验证超过限制
	err = domainService.ValidateMemberLimit(tenant, 3)
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
	
	// 5. 验证角色转换规则
	err = domainService.ValidateRoleTransition(entity.RoleMember, entity.RoleAdmin)
	assert.NoError(t, err)
	
	// 6. 验证 Owner 不可修改
	err = domainService.ValidateRoleTransition(entity.RoleOwner, entity.RoleMember)
	assert.Error(t, err)
	assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
}

// TestMembershipDomainService_EdgeCases 测试边界场景
func TestMembershipDomainService_EdgeCases(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	
	t.Run("零成员限制", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Zero Limit Tenant", 0)
		err := domainService.ValidateMemberLimit(tenant, 0)
		// 零限制的租户本身就是无效的
		assert.Error(t, err)
		assert.Equal(t, tenantEntity.ErrTenantInvalid, err)
	})
	
	t.Run("刚好达到限制", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Exact Limit Tenant", 5)
		err := domainService.ValidateMemberLimit(tenant, 5)
		assert.Error(t, err)
		assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
	})
	
	t.Run("负数成员数", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Negative Count Tenant", 10)
		err := domainService.ValidateMemberLimit(tenant, -1)
		assert.NoError(t, err) // 负数应该可以通过验证
	})
	
	t.Run("空角色检查", func(t *testing.T) {
		userID := uuid.New()
		tenantID := uuid.New()
		
		canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, "")
		assert.False(t, canJoin)
	})
}

// TestMembershipDomainService_ComplexScenarios 测试复杂场景
func TestMembershipDomainService_ComplexScenarios(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	
	t.Run("批量成员加入验证", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Batch Join Tenant", 5)
		
		// 模拟批量加入场景
		for i := 0; i < 5; i++ {
			currentCount := i
			if i < 4 {
				err := domainService.ValidateMemberLimit(tenant, currentCount)
				assert.NoError(t, err)
			} else {
				// i == 4 时，currentCount=4，还可以添加第 5 个成员
				// i == 5 时才会失败，但循环只到 4
				err := domainService.ValidateMemberLimit(tenant, currentCount+1)
				assert.Error(t, err)
				assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
			}
		}
	})
	
	t.Run("角色转换矩阵验证", func(t *testing.T) {
		roles := []entity.UserRole{
			entity.RoleOwner,
			entity.RoleAdmin,
			entity.RoleMember,
			entity.RoleGuest,
		}
		
		// 测试所有角色转换组合
		for _, fromRole := range roles {
			for _, toRole := range roles {
				err := domainService.ValidateRoleTransition(fromRole, toRole)
				
				if fromRole == entity.RoleOwner {
					// Owner 不能转换到任何角色
					assert.Error(t, err)
					assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
				} else if toRole == entity.RoleOwner {
					// 不能晋升为 Owner
					assert.Error(t, err)
					assert.Equal(t, service.ErrCannotPromoteToOwner, err)
				} else {
					// 其他转换都允许
					assert.NoError(t, err)
				}
			}
		}
	})
	
	t.Run("用户加入资格综合验证", func(t *testing.T) {
		validUserID := uuid.New()
		validTenantID := uuid.New()
		
		// 有效场景
		assert.True(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, entity.RoleMember))
		
		// 无效用户 ID
		assert.False(t, domainService.CanUserJoinTenant(ctx, uuid.Nil, validTenantID, entity.RoleMember))
		
		// 无效租户 ID
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, uuid.Nil, entity.RoleMember))
		
		// 两者都无效
		assert.False(t, domainService.CanUserJoinTenant(ctx, uuid.Nil, uuid.Nil, entity.RoleMember))
		
		// 所有有效角色
		validRoles := []entity.UserRole{
			entity.RoleOwner, entity.RoleAdmin,
			entity.RoleMember, entity.RoleGuest,
		}
		for _, role := range validRoles {
			assert.True(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, role))
		}
		
		// 无效角色
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, "invalid_role"))
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, ""))
	})
}
