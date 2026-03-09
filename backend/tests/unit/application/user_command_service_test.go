package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/infrastructure/transaction"
	helper "go-ddd-scaffold/tests/helper"
)

// MockUnitOfWork UnitOfWork 的 Mock 实现（本文件特有）
type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Begin(ctx context.Context) (transaction.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockUnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	// 模拟成功场景，直接执行函数
	return fn(ctx)
}

// ==================== 测试用例 ====================

// TestUserCommandService_UpdateUser_Success 测试更新用户信息成功
func TestUserCommandService_UpdateUser_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockPasswordHasher := new(MockPasswordHasher)
	mockUoW := new(MockUnitOfWork)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()
	newEmail := "newemail@example.com"

	// 使用工厂创建测试用户
	factory := helper.NewUserFactory(t)
	existingUser := factory.CreateUser(
		helper.WithID(userID),
		helper.WithEmail("old@example.com"),
	)

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)

	// Mock 密码哈希
	mockPasswordHasher.On("Hash", "newpassword123").Return("$2a$12$...", nil)

	// Mock 更新用户
	mockUserRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// 3. 创建服务实例
	userSvc := service.NewUserCommandService(mockUserRepo, mockMemberRepo, mockPasswordHasher, mockUoW)

	// 4. 执行测试
	req := &dto.UpdateUserRequest{
		Email:    &newEmail,
		Password: stringPtr("newpassword123"),
	}
	
	err := userSvc.UpdateUser(ctx, userID, req)

	// 5. 验证结果
	assert.NoError(t, err)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
	mockUoW.AssertExpectations(t)
}

// TestUserCommandService_UpdateUser_UserNotFound 测试用户不存在
func TestUserCommandService_UpdateUser_UserNotFound(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockPasswordHasher := new(MockPasswordHasher)
	mockUoW := new(MockUnitOfWork)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()

	// Mock 查询用户（不存在）
	mockUserRepo.On("GetByID", ctx, userID).Return((*user_entity.User)(nil), errors.New("user not found"))

	// 3. 创建服务实例
	userSvc := service.NewUserCommandService(mockUserRepo, mockMemberRepo, mockPasswordHasher, mockUoW)

	// 4. 执行测试
	req := &dto.UpdateUserRequest{
		Email: stringPtr("new@example.com"),
	}
	
	err := userSvc.UpdateUser(ctx, userID, req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	// 不应该调用更新
	mockUserRepo.AssertNotCalled(t, "Update")
}

// TestUserCommandService_UpdateProfile_Success 测试更新个人资料成功
func TestUserCommandService_UpdateProfile_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockPasswordHasher := new(MockPasswordHasher)
	mockUoW := new(MockUnitOfWork)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()
	newNickname := "TestNickname"
	newPhone := "13800138000"
	newBio := "This is a bio"

	// 使用工厂创建测试用户
	factory := helper.NewUserFactory(t)
	existingUser := factory.CreateUser(
		helper.WithID(userID),
	)

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)

	// Mock 更新用户
	mockUserRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// 3. 创建服务实例
	userSvc := service.NewUserCommandService(mockUserRepo, mockMemberRepo, mockPasswordHasher, mockUoW)

	// 4. 执行测试
	req := &dto.UpdateProfileRequest{
		Nickname: &newNickname,
		Phone:    &newPhone,
		Bio:      &newBio,
	}
	
	err := userSvc.UpdateProfile(ctx, userID, req)

	// 5. 验证结果
	assert.NoError(t, err)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockUoW.AssertExpectations(t)
}

// TestUserCommandService_DeleteUser_Success 测试删除用户成功
func TestUserCommandService_DeleteUser_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockPasswordHasher := new(MockPasswordHasher)
	mockUoW := new(MockUnitOfWork)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()

	// Mock 删除用户
	mockUserRepo.On("Delete", ctx, userID).Return(nil)

	// 3. 创建服务实例
	userSvc := service.NewUserCommandService(mockUserRepo, mockMemberRepo, mockPasswordHasher, mockUoW)

	// 4. 执行测试
	err := userSvc.DeleteUser(ctx, userID)

	// 5. 验证结果
	assert.NoError(t, err)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
}

// TestUserCommandService_UpdateUser_UnitOfWorkError 测试事务失败
func TestUserCommandService_UpdateUser_UnitOfWorkError(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockPasswordHasher := new(MockPasswordHasher)
	mockUoW := new(MockUnitOfWork)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()
	newEmail := "new@example.com"

	// 使用工厂创建测试用户
	factory := helper.NewUserFactory(t)
	existingUser := factory.CreateUser(helper.WithID(userID))

	// Mock 查询用户
	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)

	// Mock 更新用户（会失败）
	mockUserRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("update error"))

	// 3. 创建服务实例
	userSvc := service.NewUserCommandService(mockUserRepo, mockMemberRepo, mockPasswordHasher, mockUoW)

	// 4. 执行测试
	req := &dto.UpdateUserRequest{
		Email: &newEmail,
	}
	
	err := userSvc.UpdateUser(ctx, userID, req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
}

// stringPtr 辅助函数
func stringPtr(s string) *string {
	return &s
}
