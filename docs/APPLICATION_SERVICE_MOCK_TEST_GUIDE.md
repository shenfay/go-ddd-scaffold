# Application Service Mock 测试指南

## 📋 概述

Application Service Mock 测试用于验证应用服务层的业务逻辑，通过 Mock 依赖项实现快速、隔离的单元测试。

---

## 🎯 测试策略

### 测试层次

1. **单元测试** - Mock 所有外部依赖
2. **集成测试** - 使用真实数据库和依赖
3. **端到端测试** - 完整的 HTTP 请求测试

### Mock 重点

- ✅ UnitOfWork 事务管理
- ✅ Repository 仓储层
- ✅ CasbinService 权限服务
- ✅ 领域服务（可选）

---

## 🔧 Mock 实现示例

### 1. Mock UnitOfWork

```go
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
```

### 2. Mock Repository

```go
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

func (m *MockTenantRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*user_entity.TenantWithRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user_entity.TenantWithRole), args.Error(1)
}
```

### 3. Mock CasbinService

```go
type MockCasbinService struct{}

func (m *MockCasbinService) Enforce(sub, dom, obj, act string) (bool, error) {
	return true, nil // 总是允许
}

func (m *MockCasbinService) AddRoleForUser(userID uuid.UUID, tenantID uuid.UUID, role string) error {
	return nil
}

func (m *MockCasbinService) GetRolesForUser(userID, tenantID uuid.UUID) []string {
	return []string{"owner"}
}
```

---

## 📝 典型测试场景

### 场景 1: 成功创建租户

```go
func TestTenantService_CreateTenant_Success(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUoW := new(MockUnitOfWork)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockMemberRepository)
	mockCasbinService := new(MockCasbinService)

	// 2. 设置期望
	ctx := context.Background()
	ownerID := uuid.New()

	mockUoW.On("WithTransaction", ctx).Return(nil)
	mockTenantRepo.On("Create", ctx, mock.AnythingOfType("*entity.Tenant")).Return(nil)
	mockMemberRepo.On("Create", ctx, mock.AnythingOfType("*entity.TenantMember")).Return(nil)

	// 3. 创建服务实例
	tenantSvc := service.NewTenantService(mockTenantRepo, mockMemberRepo, mockCasbinService, mockUoW)

	// 4. 执行测试
	createdTenant, err := tenantSvc.CreateTenant(ctx, "Test", "Test", ownerID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, createdTenant)

	// 6. 验证 Mock 期望
	mockUoW.AssertExpectations(t)
	mockTenantRepo.AssertExpectations(t)
	mockMemberRepo.AssertExpectations(t)
}
```

### 场景 2: UnitOfWork 失败

```go
func TestTenantService_CreateTenant_UnitOfWorkError(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUoW := new(MockUnitOfWork)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockMemberRepository)
	mockCasbinService := new(MockCasbinService)

	// 2. 设置期望
	ctx := context.Background()
	ownerID := uuid.New()

	// Mock UnitOfWork 失败
	mockUoW.On("WithTransaction", ctx).Return(errors.New("transaction error"))

	// 3. 创建服务实例
	tenantSvc := service.NewTenantService(mockTenantRepo, mockMemberRepo, mockCasbinService, mockUoW)

	// 4. 执行测试
	createdTenant, err := tenantSvc.CreateTenant(ctx, "Test", "Test", ownerID)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, createdTenant)
	assert.Contains(t, err.Error(), "transaction error")

	// 6. 验证 Mock 期望
	mockUoW.AssertExpectations(t)
}
```

### 场景 3: Repository 异常

```go
func TestTenantService_CreateTenant_TenantCreateError(t *testing.T) {
	// 1. 准备 Mock 依赖
	mockUoW := new(MockUnitOfWork)
	mockTenantRepo := new(MockTenantRepository)
	mockMemberRepo := new(MockMemberRepository)
	mockCasbinService := new(MockCasbinService)

	// 2. 设置期望
	ctx := context.Background()
	ownerID := uuid.New()

	// Mock UnitOfWork 执行
	mockUoW.On("WithTransaction", ctx).Return(nil)

	// Mock TenantRepository.Create 失败
	mockTenantRepo.On("Create", ctx, mock.AnythingOfType("*entity.Tenant")).Return(errors.New("create tenant error"))

	// 3. 创建服务实例
	tenantSvc := service.NewTenantService(mockTenantRepo, mockMemberRepo, mockCasbinService, mockUoW)

	// 4. 执行测试
	createdTenant, err := tenantSvc.CreateTenant(ctx, "Test", "Test", ownerID)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, createdTenant)
	assert.Contains(t, err.Error(), "create tenant error")

	// 6. 验证 Mock 期望
	mockUoW.AssertExpectations(t)
	mockTenantRepo.AssertExpectations(t)
}
```

---

## 🚀 集成测试

### 使用真实数据库

```go
func TestTenantService_CreateTenant_Integration(t *testing.T) {
	// 1. 准备测试数据库（内存 SQLite）
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&user_entity.Tenant{}, &user_entity.TenantMember{})

	// 2. 初始化依赖
	uow := transaction.NewGormUnitOfWork(db)
	tenantRepo := repo.NewTenantDAORepository(db)
	memberRepo := repo.NewTenantMemberDAORepository(db)
	casbinService := &MockCasbinService{} // 使用 Mock Casbin

	// 3. 创建服务实例
	tenantSvc := service.NewTenantService(tenantRepo, memberRepo, casbinService, uow)

	// 4. 执行测试
	ctx := context.Background()
	ownerID := uuid.New()
	createdTenant, err := tenantSvc.CreateTenant(ctx, "Test", "Test", ownerID)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, createdTenant)

	// 6. 验证数据库中确实保存了数据
	var tenantCount int64
	db.Model(&user_entity.Tenant{}).Where("id = ?", createdTenant.ID).Count(&tenantCount)
	assert.Equal(t, int64(1), tenantCount)
}
```

---

## ✅ 最佳实践

### 1. Mock 职责分离

```go
// ✅ 好的做法：每个 Mock 只关注一个接口
type MockUnitOfWork struct{ mock.Mock }
type MockTenantRepository struct{ mock.Mock }
type MockMemberRepository struct{ mock.Mock }

// ❌ 不好的做法：一个大 Mock 实现多个接口
type MockEverything struct{ mock.Mock }
```

### 2. 清晰的测试命名

```go
// ✅ 清晰表达测试意图
TestTenantService_CreateTenant_Success
TestTenantService_CreateTenant_UnitOfWorkError
TestTenantService_CreateTenant_TenantCreateError
TestTenantService_GetUserTenants_Success

// ❌ 模糊的命名
TestTenant1
TestCreate
```

### 3. AAA 模式组织测试

```go
func TestExample(t *testing.T) {
	// Arrange - 准备阶段
	mockRepo := new(MockRepository)
	mockRepo.On("GetByID", id).Return(expectedEntity, nil)

	// Act - 执行阶段
	result, err := service.GetByID(id)

	// Assert - 验证阶段
	assert.NoError(t, err)
	assert.Equal(t, expectedEntity, result)
	mockRepo.AssertExpectations(t)
}
```

### 4. 使用表格驱动测试

```go
func TestTenantService_CreateTenant_MultipleScenarios(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func(*MockUnitOfWork, *MockTenantRepository)
		expectedError string
	}{
		{
			name: "成功创建",
			setupMocks: func(uow *MockUnitOfWork, repo *MockTenantRepository) {
				uow.On("WithTransaction", mock.Anything).Return(nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "UnitOfWork 失败",
			setupMocks: func(uow *MockUnitOfWork, repo *MockTenantRepository) {
				uow.On("WithTransaction", mock.Anything).Return(errors.New("tx error"))
			},
			expectedError: "tx error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试逻辑
		})
	}
}
```

---

## 📊 测试覆盖率要求

| 层级 | 最低覆盖率 | 推荐工具 |
|------|-----------|---------|
| **Application Service** | 80%+ | testify/mock |
| **Domain Service** | 90%+ | 原生 testing + assert |
| **Repository** | 70%+ | testify/mock |
| **Infrastructure** | 60%+ | 集成测试 |

---

## 🔍 常见问题

### Q1: 如何处理依赖注入？

**A**: 使用构造函数注入，便于在测试中替换为 Mock：

```go
// 生产代码
func NewTenantService(
	tenantRepo repository.TenantRepository,
	memberRepo repository.TenantMemberRepository,
	casbinService auth.CasbinService,
	uow transaction.UnitOfWork,
) TenantService {
	return &tenantService{...}
}

// 测试代码
mockRepo := new(MockTenantRepository)
service := NewTenantService(mockRepo, ...)
```

### Q2: Mock 期望不匹配怎么办？

**A**: 使用 `mock.Anything` 或自定义匹配器：

```go
// 宽松匹配
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Tenant")).Return(nil)

// 严格匹配
expectedID := uuid.New()
mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(t *entity.Tenant) bool {
	return t.ID == expectedID
})).Return(nil)
```

### Q3: 如何测试异步操作？

**A**: 使用 channel 等待异步完成：

```go
done := make(chan bool)
mockEventBus.On("Publish", mock.Anything).Run(func(args mock.Arguments) {
	done <- true
}).Return(nil)

<-done // 等待事件发布
```

---

## 📚 参考资源

- [testify/mock 官方文档](https://github.com/stretchr/testify)
- [Go Mock 最佳实践](https://go.dev/doc/tutorial/mocking-dependencies)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

## 🎯 下一步建议

1. ✅ **补充 UserCommandService Mock 测试**
2. ✅ **添加错误传播测试**
3. ✅ **编写并发安全测试**
4. ✅ **集成到 CI/CD 流水线**
