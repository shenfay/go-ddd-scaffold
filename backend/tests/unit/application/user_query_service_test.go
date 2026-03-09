// Package service_test UserQueryService Mock 测试
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/application/user/service"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	helper "go-ddd-scaffold/tests/helper"
)

// getNow 返回当前时间指针（测试辅助函数）
func getNow() time.Time {
	return time.Now()
}

// ==================== 测试用例 ====================

// TestUserQueryService_GetUser_Success 测试获取用户信息成功
func TestUserQueryService_GetUser_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	userID := uuid.New()

	// 创建测试数据
	user := factory.CreateUser(
		helper.WithID(userID),
		helper.WithEmail("test@example.com"),
		helper.WithNickname("TestUser"),
	)

	tenantID := uuid.New()
	member := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     user_entity.RoleMember,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: getNow(),
	}

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock 查询租户成员关系
	mockMemberRepo.On("ListByUser", ctx, userID).Return([]*user_entity.TenantMember{member}, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.GetUser(ctx, userID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "TestUser", result.Nickname)
	assert.Equal(t, string(user_entity.RoleMember), result.Role)
	assert.NotNil(t, result.TenantID)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockMemberRepo.AssertExpectations(t)
}

// TestUserQueryService_GetUser_NoTenant 测试获取用户信息（无租户）
func TestUserQueryService_GetUser_NoTenant(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	userID := uuid.New()

	// 创建测试数据
	user := factory.CreateUser(
		helper.WithID(userID),
		helper.WithEmail("test@example.com"),
		helper.WithNickname("TestUser"),
	)

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock 查询租户成员关系 - 返回空列表
	mockMemberRepo.On("ListByUser", ctx, userID).Return([]*user_entity.TenantMember{}, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.GetUser(ctx, userID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "TestUser", result.Nickname)
	assert.Empty(t, result.Role)
	assert.Nil(t, result.TenantID)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockMemberRepo.AssertExpectations(t)
}

// TestUserQueryService_GetUser_UserNotFound 测试获取用户信息 - 用户不存在
func TestUserQueryService_GetUser_UserNotFound(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()

	// Mock 查询用户 - 返回错误
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, assert.AnError)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.GetUser(ctx, userID)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
}

// TestUserQueryService_GetUserInfo_Success 测试获取当前用户信息成功
func TestUserQueryService_GetUserInfo_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	userID := uuid.New()

	// 创建测试数据
	user := factory.CreateUser(
		helper.WithID(userID),
		helper.WithEmail("test@example.com"),
		helper.WithNickname("TestUser"),
	)

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.GetUserInfo(ctx, userID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "TestUser", result.Nickname)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
}

// TestUserQueryService_ListUsersByTenant_Success 测试列出租户下所有用户成功
func TestUserQueryService_ListUsersByTenant_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	tenantID := uuid.New()

	// 创建测试数据
	user1 := factory.CreateUser(
		helper.WithID(uuid.New()),
		helper.WithEmail("user1@example.com"),
		helper.WithNickname("User1"),
	)

	user2 := factory.CreateUser(
		helper.WithID(uuid.New()),
		helper.WithEmail("user2@example.com"),
		helper.WithNickname("User2"),
	)

	member1 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   user1.ID,
		Role:     user_entity.RoleOwner,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: getNow(),
	}

	member2 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   user2.ID,
		Role:     user_entity.RoleMember,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: getNow(),
	}

	// Mock 查询成员列表
	mockMemberRepo.On("ListByTenant", ctx, tenantID).Return([]*user_entity.TenantMember{member1, member2}, nil)

	// Mock 查询用户详情
	mockUserRepo.On("GetByID", ctx, user1.ID).Return(user1, nil)
	mockUserRepo.On("GetByID", ctx, user2.ID).Return(user2, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.ListUsersByTenant(ctx, tenantID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1@example.com", result[0].Email)
	assert.Equal(t, string(user_entity.RoleOwner), result[0].Role)
	assert.Equal(t, "user2@example.com", result[1].Email)
	assert.Equal(t, string(user_entity.RoleMember), result[1].Role)

	// 6. 验证 Mock 期望
	mockMemberRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestUserQueryService_ListUsersByTenant_SomeUsersNotFound 测试部分用户不存在时跳过
func TestUserQueryService_ListUsersByTenant_SomeUsersNotFound(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	tenantID := uuid.New()

	// 创建测试数据
	user1 := factory.CreateUser(
		helper.WithID(uuid.New()),
		helper.WithEmail("user1@example.com"),
		helper.WithNickname("User1"),
	)

	member1 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   user1.ID,
		Role:     user_entity.RoleOwner,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: getNow(),
	}

	member2 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   uuid.New(), // 这个用户不存在
		Role:     user_entity.RoleMember,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: getNow(),
	}

	// Mock 查询成员列表
	mockMemberRepo.On("ListByTenant", ctx, tenantID).Return([]*user_entity.TenantMember{member1, member2}, nil)

	// Mock 查询用户详情 - 第一个成功，第二个失败
	mockUserRepo.On("GetByID", ctx, user1.ID).Return(user1, nil)
	mockUserRepo.On("GetByID", ctx, member2.UserID).Return(nil, assert.AnError)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.ListUsersByTenant(ctx, tenantID)

	// 5. 验证结果 - 只返回存在的用户
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "user1@example.com", result[0].Email)

	// 6. 验证 Mock 期望
	mockMemberRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestUserQueryService_ListMembersByTenant_Success 测试列出活跃成员成功
func TestUserQueryService_ListMembersByTenant_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	factory := helper.NewUserFactory(t)
	tenantID := uuid.New()

	// 创建测试数据
	user1 := factory.CreateUser(
		helper.WithID(uuid.New()),
		helper.WithEmail("user1@example.com"),
		helper.WithNickname("User1"),
	)

	user2 := factory.CreateUser(
		helper.WithID(uuid.New()),
		helper.WithEmail("user2@example.com"),
		helper.WithNickname("User2"),
	)

	member1 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   user1.ID,
		Role:     user_entity.RoleOwner,
		Status:   user_entity.MemberStatusActive, // 活跃状态
		JoinedAt: getNow(),
	}

	member2 := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   user2.ID,
		Role:     user_entity.RoleMember,
		Status:   user_entity.MemberStatusInactive, // 非活跃状态
		JoinedAt: getNow(),
	}

	// Mock 查询成员列表
	mockMemberRepo.On("ListByTenant", ctx, tenantID).Return([]*user_entity.TenantMember{member1, member2}, nil)

	// Mock 查询用户详情（只查询活跃的）
	mockUserRepo.On("GetByID", ctx, user1.ID).Return(user1, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.ListMembersByTenant(ctx, tenantID)

	// 5. 验证结果 - 只返回活跃成员
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "user1@example.com", result[0].Email)
	assert.Equal(t, string(user_entity.RoleOwner), result[0].Role)

	// 6. 验证 Mock 期望
	mockMemberRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestUserQueryService_ListMembersByTenant_Empty 测试空租户
func TestUserQueryService_ListMembersByTenant_Empty(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)

	// 2. 设置期望
	ctx := context.Background()
	tenantID := uuid.New()

	// Mock 查询成员列表 - 返回空
	mockMemberRepo.On("ListByTenant", ctx, tenantID).Return([]*user_entity.TenantMember{}, nil)

	// 3. 创建服务实例
	querySvc := service.NewUserQueryService(mockUserRepo, mockMemberRepo)

	// 4. 执行测试
	result, err := querySvc.ListMembersByTenant(ctx, tenantID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.Empty(t, result)

	// 6. 验证 Mock 期望
	mockMemberRepo.AssertExpectations(t)
}
