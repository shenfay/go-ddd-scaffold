// Package service_test AuthenticationService Mock 测试
package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	eventBus "go-ddd-scaffold/internal/infrastructure/event"
	helper "go-ddd-scaffold/tests/helper"
)

// ==================== Mock 实现 ====================

// MockUserRepository UserRepository 的 Mock 实现（复用已有的）
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *user_entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user_entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *user_entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*user_entity.User, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user_entity.User), args.Error(1)
}

func (m *MockUserRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) WithTx(tx *gorm.DB) repository.BaseRepository[user_entity.User, uuid.UUID] {
	args := m.Called(tx)
	return args.Get(0).(repository.UserRepository)
}

// MockTenantRepository TenantRepository 的 Mock 实现
type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *user_entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*user_entity.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Update(ctx context.Context, tenant *user_entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTenantRepository) WithTx(tx *gorm.DB) repository.BaseRepository[user_entity.Tenant, uuid.UUID] {
	args := m.Called(tx)
	return args.Get(0).(repository.TenantRepository)
}

// MockTenantMemberRepository TenantMemberRepository 的 Mock 实现（复用已有的）
type MockTenantMemberRepository struct {
	mock.Mock
}

func (m *MockTenantMemberRepository) Create(ctx context.Context, member *user_entity.TenantMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockTenantMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*user_entity.TenantMember, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.TenantMember), args.Error(1)
}

func (m *MockTenantMemberRepository) Update(ctx context.Context, member *user_entity.TenantMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockTenantMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTenantMemberRepository) GetByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) (*user_entity.TenantMember, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.TenantMember), args.Error(1)
}

func (m *MockTenantMemberRepository) DeleteByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) error {
	args := m.Called(ctx, userID, tenantID)
	return args.Error(0)
}

func (m *MockTenantMemberRepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTenantMemberRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*user_entity.TenantMember, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user_entity.TenantMember), args.Error(1)
}

func (m *MockTenantMemberRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*user_entity.TenantMember, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user_entity.TenantMember), args.Error(1)
}

func (m *MockTenantMemberRepository) WithTx(tx *gorm.DB) repository.BaseRepository[user_entity.TenantMember, uuid.UUID] {
	args := m.Called(tx)
	return args.Get(0).(repository.TenantMemberRepository)
}

// MockJWTService JWTService 的 Mock 实现
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*user_entity.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_entity.TokenClaims), args.Error(1)
}

func (m *MockJWTService) GenerateTokenWithTenant(userID uuid.UUID, tenantID uuid.UUID) (string, error) {
	args := m.Called(userID, tenantID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// MockEventBus EventBus 的 Mock 实现
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event eventBus.DomainEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockTokenBlacklistService TokenBlacklistService 的 Mock 实现
type MockTokenBlacklistService struct {
	mock.Mock
}

func (m *MockTokenBlacklistService) AddToBlacklist(ctx context.Context, token string, expireAt time.Time) error {
	args := m.Called(ctx, token, expireAt)
	return args.Error(0)
}

func (m *MockTokenBlacklistService) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenBlacklistService) IsBlacklistedBatch(ctx context.Context, tokens []string) (map[string]bool, error) {
	args := m.Called(ctx, tokens)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockTokenBlacklistService) RemoveFromBlacklist(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenBlacklistService) CleanExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockPasswordHasher PasswordHasher 的 Mock 实现（复用已有的）
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

// ==================== 测试用例 ====================

// TestAuthenticationService_Register_Success 测试用户注册成功
func TestAuthenticationService_Register_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	req := &dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "TestUser",
	}

	// Mock 检查邮箱不存在
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return((*user_entity.User)(nil), nil)

	// Mock 密码哈希
	mockPasswordHasher.On("Hash", req.Password).Return("$2a$12$...", nil)

	// Mock 创建用户
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// Mock 发布事件
	mockEventBus.On("Publish", ctx, mock.AnythingOfType("*event.UserRegisteredEvent")).Return(nil)

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	user, err := authSvc.Register(ctx, req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestAuthenticationService_Register_EmailExists 测试邮箱已存在
func TestAuthenticationService_Register_EmailExists(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	req := &dto.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Nickname: "TestUser",
	}

	// Mock 检查邮箱存在
	factory := helper.NewUserFactory(t)
	existingUser := factory.CreateUser(helper.WithEmail(req.Email))
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	user, err := authSvc.Register(ctx, req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "用户已存在")

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	// 不应该调用创建和发布事件
	mockUserRepo.AssertNotCalled(t, "Create")
	mockEventBus.AssertNotCalled(t, "Publish")
}

// TestAuthenticationService_Login_Success 测试登录成功
func TestAuthenticationService_Login_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	req := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock 查询用户
	factory := helper.NewUserFactory(t)
	user := factory.CreateUser(
		helper.WithEmail(req.Email),
	)
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Mock 密码验证（注意参数顺序：Verify(hash, plainPassword)）
	mockPasswordHasher.On("Verify", string(user.Password), req.Password).Return(true)

	// Mock 生成 Token
	mockJWTService.On("GenerateToken", user.ID).Return("mock_token_123", nil)

	// Mock 发布登录事件
	mockEventBus.On("Publish", ctx, mock.AnythingOfType("*event.UserLoggedInEvent")).Return(nil)

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	response, err := authSvc.Login(ctx, req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "mock_token_123", response.AccessToken)
	assert.Equal(t, req.Email, response.User.Email)

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
}

// TestAuthenticationService_Login_UserNotFound 测试用户不存在
func TestAuthenticationService_Login_UserNotFound(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	req := &dto.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	// Mock 查询用户不存在
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return((*user_entity.User)(nil), errors.New("user not found"))

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	response, err := authSvc.Login(ctx, req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "用户不存在")

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
}

// TestAuthenticationService_Login_WrongPassword 测试密码错误
func TestAuthenticationService_Login_WrongPassword(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	req := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock 查询用户
	factory := helper.NewUserFactory(t)
	user := factory.CreateUser(helper.WithEmail(req.Email))
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Mock 密码验证失败（注意参数顺序）
	mockPasswordHasher.On("Verify", string(user.Password), req.Password).Return(false)

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	response, err := authSvc.Login(ctx, req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "密码错误")

	// 6. 验证 Mock 期望
	mockUserRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
}

// TestAuthenticationService_Logout_Success 测试登出成功
func TestAuthenticationService_Logout_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUserRepo := new(MockUserRepository)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockTenantMemberRepository)
	mockJWTService := new(MockJWTService)
	mockEventBus := new(MockEventBus)
	mockTokenBlacklist := new(MockTokenBlacklistService)
	mockPasswordHasher := new(MockPasswordHasher)

	// 2. 设置期望
	ctx := context.Background()
	userID := uuid.New()
	token := "valid_token_123"

	// Mock Token 验证成功
	claims := &user_entity.TokenClaims{UserID: userID}
	mockJWTService.On("ValidateToken", token).Return(claims, nil)

	// Mock 加入黑名单
	mockTokenBlacklist.On("AddToBlacklist", ctx, token, mock.AnythingOfType("time.Time")).Return(nil)

	// 3. 创建服务实例
	authSvc := service.NewAuthenticationService(
		mockUserRepo,
		mockTenantRepo,
		mockMemberRepo,
		mockJWTService,
		mockEventBus,
		mockTokenBlacklist,
		mockPasswordHasher,
	)

	// 4. 执行测试
	err := authSvc.Logout(ctx, userID, token)

	// 5. 验证结果
	assert.NoError(t, err)

	// 6. 验证 Mock 期望
	mockJWTService.AssertExpectations(t)
	mockTokenBlacklist.AssertExpectations(t)
}
