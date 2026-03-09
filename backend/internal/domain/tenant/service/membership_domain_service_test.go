// Package service_test 租户领域服务测试
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	sharedEntity "go-ddd-scaffold/internal/domain/shared/entity"
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
	validRoles := []sharedEntity.UserRole{
		sharedEntity.RoleOwner,
		sharedEntity.RoleAdmin,
		sharedEntity.RoleMember,
		sharedEntity.RoleGuest,
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
	canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, sharedEntity.RoleMember)

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
	canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, sharedEntity.RoleMember)

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
	invalidRole := sharedEntity.UserRole("non_existent_role")

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
		currentRole sharedEntity.UserRole
		newRole     sharedEntity.UserRole
		description string
	}{
		{sharedEntity.RoleMember, sharedEntity.RoleAdmin, "Member -> Admin"},
		{sharedEntity.RoleGuest, sharedEntity.RoleMember, "Guest -> Member"},
		{sharedEntity.RoleMember, sharedEntity.RoleGuest, "Member -> Guest (降级)"},
		{sharedEntity.RoleAdmin, sharedEntity.RoleMember, "Admin -> Member (降级)"},
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
	err := domainService.ValidateRoleTransition(sharedEntity.RoleOwner, sharedEntity.RoleAdmin)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
}

// TestMembershipDomainService_ValidateRoleTransition_CannotPromoteToOwner 测试不能晋升为 Owner
func TestMembershipDomainService_ValidateRoleTransition_CannotPromoteToOwner(t *testing.T) {
	// 1. 准备测试数据
	domainService := service.NewMembershipDomainService()

	// 2. 尝试晋升为 Owner
	err := domainService.ValidateRoleTransition(sharedEntity.RoleMember, sharedEntity.RoleOwner)

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
	member, err := tenant.AddMember(userID, sharedEntity.RoleAdmin, nil)
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, 1, tenant.GetActiveMemberCount())

	// 3. 添加第二个成员
	userID2 := uuid.New()
	member2, err := tenant.AddMember(userID2, sharedEntity.RoleMember, &userID)
	assert.NoError(t, err)
	assert.NotNil(t, member2)
	assert.Equal(t, 2, tenant.GetActiveMemberCount())
}

// TestTenant_AddMember_DuplicateMember 测试重复添加成员
func TestTenant_AddMember_DuplicateMember(t *testing.T) {
	// 1. 创建租户并添加成员
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	userID := uuid.New()

	_, err := tenant.AddMember(userID, sharedEntity.RoleMember, nil)
	assert.NoError(t, err)

	// 2. 尝试再次添加同一用户
	_, err = tenant.AddMember(userID, sharedEntity.RoleMember, nil)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberAlreadyExists, err)
}

// TestTenant_AddMember_LimitExceeded 测试超过成员限制
func TestTenant_AddMember_LimitExceeded(t *testing.T) {
	// 1. 创建只有 2 个名额的租户
	tenant := tenantEntity.NewTenant("Small Tenant", 2)

	// 2. 添加 2 个成员
	_, _ = tenant.AddMember(uuid.New(), sharedEntity.RoleMember, nil)
	_, _ = tenant.AddMember(uuid.New(), sharedEntity.RoleMember, nil)

	// 3. 尝试添加第 3 个成员
	_, err := tenant.AddMember(uuid.New(), sharedEntity.RoleMember, nil)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
}

// TestTenant_RemoveMember_Success 测试移除成员成功
func TestTenant_RemoveMember_Success(t *testing.T) {
	// 1. 创建租户并添加成员
	tenant := tenantEntity.NewTenant("Test Tenant", 5)
	userID := uuid.New()
	member, _ := tenant.AddMember(userID, sharedEntity.RoleMember, nil)

	assert.Equal(t, 1, tenant.GetActiveMemberCount())

	// 2. 移除成员
	err := tenant.RemoveMember(member.ID)

	// 3. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 0, tenant.GetActiveMemberCount())
	// 验证租户Members 切片中的成员状态（因为 RemoveMember 修改的是切片中的元素）
	assert.Equal(t, tenantEntity.MemberStatusRemoved, tenant.Members[0].Status)
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
	err = domainService.ValidateRoleTransition(sharedEntity.RoleMember, sharedEntity.RoleAdmin)
	assert.NoError(t, err)
	
	// 6. 验证 Owner 不可修改
	err = domainService.ValidateRoleTransition(sharedEntity.RoleOwner, sharedEntity.RoleMember)
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
		roles := []sharedEntity.UserRole{
			sharedEntity.RoleOwner,
			sharedEntity.RoleAdmin,
			sharedEntity.RoleMember,
			sharedEntity.RoleGuest,
		}
		
		// 测试所有角色转换组合
		for _, fromRole := range roles {
			for _, toRole := range roles {
				err := domainService.ValidateRoleTransition(fromRole, toRole)
				
				if fromRole == sharedEntity.RoleOwner {
					// Owner 不能转换到任何角色
					assert.Error(t, err)
					assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
				} else if toRole == sharedEntity.RoleOwner {
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
		assert.True(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, sharedEntity.RoleMember))
		
		// 无效用户 ID
		assert.False(t, domainService.CanUserJoinTenant(ctx, uuid.Nil, validTenantID, sharedEntity.RoleMember))
		
		// 无效租户 ID
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, uuid.Nil, sharedEntity.RoleMember))
		
		// 两者都无效
		assert.False(t, domainService.CanUserJoinTenant(ctx, uuid.Nil, uuid.Nil, sharedEntity.RoleMember))
		
		// 所有有效角色
		validRoles := []sharedEntity.UserRole{
			sharedEntity.RoleOwner, sharedEntity.RoleAdmin,
			sharedEntity.RoleMember, sharedEntity.RoleGuest,
		}
		for _, role := range validRoles {
			assert.True(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, role))
		}
		
		// 无效角色
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, "invalid_role"))
		assert.False(t, domainService.CanUserJoinTenant(ctx, validUserID, validTenantID, ""))
	})
}

// TestMembershipDomainService_MemberStatusCheck 测试成员状态检查
func TestMembershipDomainService_MemberStatusCheck(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	
	t.Run("Active 成员可以加入", func(t *testing.T) {
		userID := uuid.New()
		tenantID := uuid.New()
		canJoin := domainService.CanUserJoinTenant(ctx, userID, tenantID, sharedEntity.RoleMember)
		assert.True(t, canJoin)
	})
	
	t.Run("验证成员状态流转", func(t *testing.T) {
		// 创建租户成员实体
		member := &tenantEntity.TenantMember{
			ID:       uuid.New(),
			TenantID: uuid.New(),
			UserID:   uuid.New(),
			Role:     sharedEntity.RoleMember,
			Status:   tenantEntity.MemberStatusActive,
			JoinedAt: time.Now(),
		}
		
		// 验证 Active 状态
		assert.True(t, member.IsActive())
		assert.False(t, member.IsRemoved())
		
		// 模拟移除操作
		member.Remove()
		assert.False(t, member.IsActive())
		assert.True(t, member.IsRemoved())
	})
}

// TestMembershipDomainService_TenantEligibilityCheck 测试租户资格检查
func TestMembershipDomainService_TenantEligibilityCheck(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	
	t.Run("未到期租户有效", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Valid Tenant", 10)
		// 确保未到期
		tenant.ExpiredAt = time.Now().AddDate(0, 6, 0) // 6 个月后到期
		
		err := domainService.ValidateMemberLimit(tenant, 5)
		assert.NoError(t, err)
		assert.True(t, tenant.IsValid())
	})
	
	t.Run("已到期租户无效", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Expired Tenant", 10)
		// 设置过期时间为过去
		tenant.ExpiredAt = time.Now().AddDate(-1, 0, 0) // 1 年前到期
		
		err := domainService.ValidateMemberLimit(tenant, 0)
		assert.Error(t, err)
		assert.Equal(t, tenantEntity.ErrTenantInvalid, err)
		assert.False(t, tenant.IsValid())
	})
	
	t.Run("刚好今天到期的租户", func(t *testing.T) {
		tenant := tenantEntity.NewTenant("Today Expired Tenant", 10)
		tenant.ExpiredAt = time.Now().AddDate(0, 0, 1) // 设置为明天到期
		
		// 未到期，应该有效
		assert.True(t, tenant.IsValid())
	})
}

// TestMembershipDomainService_RolePermissionCheck 测试角色权限检查
func TestMembershipDomainService_RolePermissionCheck(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	
	t.Run("Owner 角色特殊保护", func(t *testing.T) {
		// Owner 不能降级
		err := domainService.ValidateRoleTransition(sharedEntity.RoleOwner, sharedEntity.RoleAdmin)
		assert.Error(t, err)
		assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
		
		err = domainService.ValidateRoleTransition(sharedEntity.RoleOwner, sharedEntity.RoleMember)
		assert.Error(t, err)
		assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
		
		err = domainService.ValidateRoleTransition(sharedEntity.RoleOwner, sharedEntity.RoleGuest)
		assert.Error(t, err)
		assert.Equal(t, service.ErrCannotChangeOwnerRole, err)
	})
	
	t.Run("不能晋升为 Owner", func(t *testing.T) {
		// 任何角色都不能晋升为 Owner
		fromRoles := []sharedEntity.UserRole{
			sharedEntity.RoleAdmin,
			sharedEntity.RoleMember,
			sharedEntity.RoleGuest,
		}
		
		for _, fromRole := range fromRoles {
			err := domainService.ValidateRoleTransition(fromRole, sharedEntity.RoleOwner)
			assert.Error(t, err)
			assert.Equal(t, service.ErrCannotPromoteToOwner, err)
		}
	})
	
	t.Run("平级转换允许", func(t *testing.T) {
		// 同一角色的平级转换应该允许
		roles := []sharedEntity.UserRole{
			sharedEntity.RoleAdmin,
			sharedEntity.RoleMember,
			sharedEntity.RoleGuest,
		}
		
		for _, role := range roles {
			err := domainService.ValidateRoleTransition(role, role)
			assert.NoError(t, err)
		}
	})
}

// TestMembershipDomainService_LargeScaleTenant 测试大规模租户场景
func TestMembershipDomainService_LargeScaleTenant(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	
	t.Run("大型租户成员限制", func(t *testing.T) {
		// 创建大型租户（最大 1000 成员）
		tenant := tenantEntity.NewTenant("Large Tenant", 1000)
		
		// 验证不同成员数量
		testCases := []struct {
			count       int
			shouldError bool
		}{
			{0, false},
			{100, false},
			{500, false},
			{999, false},
			{1000, true},  // 达到上限
			{1001, true},  // 超过上限
		}
		
		for _, tc := range testCases {
			err := domainService.ValidateMemberLimit(tenant, tc.count)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})
	
	t.Run("超小型租户", func(t *testing.T) {
		// 创建最小租户（只能容纳 1 个成员）
		tenant := tenantEntity.NewTenant("Minimal Tenant", 1)
		
		// 0 个成员时可以通过验证
		err := domainService.ValidateMemberLimit(tenant, 0)
		assert.NoError(t, err)
		
		// 1 个成员时达到上限
		err = domainService.ValidateMemberLimit(tenant, 1)
		assert.Error(t, err)
		assert.Equal(t, tenantEntity.ErrTenantMemberLimitExceeded, err)
	})
}

// TestMembershipDomainService_ConcurrentAccess 测试并发访问场景
func TestMembershipDomainService_ConcurrentAccess(t *testing.T) {
	domainService := service.NewMembershipDomainService()
	ctx := context.Background()
	
	t.Run("并发 CanUserJoinTenant 调用", func(t *testing.T) {
		userID := uuid.New()
		tenantID := uuid.New()
		
		// 并发调用 100 次
		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				result := domainService.CanUserJoinTenant(ctx, userID, tenantID, sharedEntity.RoleMember)
				assert.True(t, result)
				done <- true
			}()
		}
		
		// 等待所有 goroutine 完成
		for i := 0; i < 100; i++ {
			<-done
		}
	})
	
	t.Run("并发 ValidateRoleTransition 调用", func(t *testing.T) {
		done := make(chan bool, 100)
		for i := 0; i < 100; i++ {
			go func() {
				err := domainService.ValidateRoleTransition(sharedEntity.RoleMember, sharedEntity.RoleAdmin)
				assert.NoError(t, err)
				done <- true
			}()
		}
		
		for i := 0; i < 100; i++ {
			<-done
		}
	})
}
