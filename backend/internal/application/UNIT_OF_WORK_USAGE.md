# Application Service UnitOfWork 集成指南

## 📋 概述

Application Service 现已注入 `transaction.UnitOfWork`，可以在跨聚合根操作中使用事务保证原子性。

---

## 🔧 Wire 配置

### providers.go 新增内容

```go
// InitializeUnitOfWork 初始化工作单元（事务管理）
func InitializeUnitOfWork(db *gorm.DB) transaction.UnitOfWork {
	return transaction.NewGormUnitOfWork(db)
}
```

### injector.wire 配置

```go
// wireinject/wire.go
func InitializeApp() (*App, error) {
	wire.Build(
		// ... 其他依赖
		wire.ProviderSet{
			transaction.InitializeUnitOfWork,
			// ...
		},
	)
	return nil, nil
}
```

---

## 💡 Application Service 改造

### UserCommandService 结构

```go
type userCommandService struct {
	userRepo         repository.UserRepository
	tenantMemberRepo repository.TenantMemberRepository
	passwordHasher   service.PasswordHasher
	uow              transaction.UnitOfWork // ✅ 新增
}

func NewUserCommandService(
	userRepo repository.UserRepository,
	tenantMemberRepo repository.TenantMemberRepository,
	passwordHasher service.PasswordHasher,
	uow transaction.UnitOfWork, // ✅ 新增参数
) UserCommandService {
	return &userCommandService{
		userRepo:         userRepo,
		tenantMemberRepo: tenantMemberRepo,
		passwordHasher:   passwordHasher,
		uow:              uow,
	}
}
```

---

## 🎯 使用场景

### 场景 1: 用户注册并自动加入租户（跨聚合根事务）

**需求**: 注册用户后，自动将其加入指定租户，两个操作必须原子性完成。

```go
type TenantApplicationService struct {
	uow              transaction.UnitOfWork
	userRepo         repository.UserRepository
	tenantRepo       repository.TenantRepository
	memberRepo       repository.TenantMemberRepository
	membershipDomain domain.MembershipDomainService
}

func (s *TenantApplicationService) RegisterUserAndJoinTenant(
	ctx context.Context,
	req *RegisterAndJoinRequest,
) (*entity.User, error) {
	
	var createdUser *entity.User
	
	// ✅ 使用 UnitOfWork 保证原子性
	err := s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 步骤 1: 创建用户
		hashedPassword, err := s.passwordHasher.Hash(req.Password)
		if err != nil {
			return err
		}
		
		user, err := entity.NewUser(req.Email, hashedPassword, req.Nickname)
		if err != nil {
			return err
		}
		
		// 切换到事务仓储
		tx := getTxFromContext(ctx)
		userRepo := s.userRepo.WithTx(tx)
		
		if err := userRepo.Create(ctx, user); err != nil {
			return err
		}
		
		// 步骤 2: 验证租户成员限制
		tenantRepo := s.tenantRepo.WithTx(tx)
		tenant, err := tenantRepo.GetByID(ctx, req.TenantID)
		if err != nil {
			return err
		}
		
		memberRepo := s.memberRepo.WithTx(tx)
		currentCount, _ := memberRepo.CountByTenant(ctx, req.TenantID)
		
		if err := s.membershipDomain.ValidateMemberLimit(tenant, currentCount+1); err != nil {
			return err
		}
		
		// 步骤 3: 创建租户成员关系
		member, err := entity.NewTenantMember(req.TenantID, user.ID, req.Role)
		if err != nil {
			return err
		}
		
		if err := memberRepo.Create(ctx, member); err != nil {
			return err
		}
		
		createdUser = user
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return createdUser, nil
}
```

**事务保证**:
- ✅ 要么用户和成员关系都创建成功
- ✅ 要么都回滚，不会出现脏数据
- ✅ 租户成员限制验证在事务中进行

---

### 场景 2: 批量更新用户信息

```go
func (s *UserCommandService) BatchUpdateUserInfo(
	ctx context.Context,
	userIDs []uuid.UUID,
	updateFunc func(*entity.User) error,
) error {
	return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		tx := getTxFromContext(ctx)
		userRepo := s.userRepo.WithTx(tx)
		
		for _, userID := range userIDs {
			user, err := userRepo.GetByID(ctx, userID)
			if err != nil {
				return err
			}
			
			if err := updateFunc(user); err != nil {
				return err
			}
			
			if err := userRepo.Update(ctx, user); err != nil {
				return err
			}
		}
		
		return nil
	})
}
```

---

### 场景 3: 删除租户及其所有成员

```go
func (s *TenantService) DeleteTenantAndMembers(
	ctx context.Context,
	tenantID uuid.UUID,
) error {
	return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		tx := getTxFromContext(ctx)
		
		tenantRepo := s.tenantRepo.WithTx(tx)
		memberRepo := s.memberRepo.WithTx(tx)
		
		// 步骤 1: 获取所有成员
		members, err := memberRepo.ListByTenant(ctx, tenantID)
		if err != nil {
			return err
		}
		
		// 步骤 2: 删除所有成员
		for _, member := range members {
			if err := memberRepo.Delete(ctx, member.ID); err != nil {
				return err
			}
		}
		
		// 步骤 3: 删除租户
		tenant, err := tenantRepo.GetByID(ctx, tenantID)
		if err != nil {
			return err
		}
		
		return tenantRepo.Delete(ctx, tenant.ID)
	})
}
```

---

## 🔍 辅助函数

### 从 Context 获取事务

```go
// transaction/context_tx.go
package transaction

import (
	"context"
	
	"gorm.io/gorm"
)

type txKeyType struct{}

var txKey = txKeyType{}

// ContextWithTx 将事务添加到 context
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// GetTxFromContext 从 context 获取事务
func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// getTxFromContext 简化版本（内部使用）
func getTxFromContext(ctx context.Context) *gorm.DB {
	return GetTxFromContext(ctx)
}
```

---

## ⚠️ 注意事项

### 1. 不要在事务外使用事务仓储

```go
// ❌ 错误示范
tx, _ := uow.Begin(ctx)
userRepo := s.userRepo.WithTx(tx.GetDB())

// 在另一个非事务上下文中使用
userRepo.GetByID(otherCtx, userID) // 可能失败！

// ✅ 正确做法
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
	tx := getTxFromContext(ctx)
	userRepo := s.userRepo.WithTx(tx)
	return userRepo.Create(ctx, user)
})
```

### 2. 事务仓储是临时的

```go
// ✅ 推荐：每次事务都创建新实例
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
	tx := getTxFromContext(ctx)
	userRepo := s.userRepo.WithTx(tx)
	// 使用...
	return nil
})

// ❌ 不推荐：缓存事务仓储
cachedRepo := s.userRepo.WithTx(someTx) // 不要这样做！
```

### 3. 原始仓储仍然可用

```go
// ✅ 正常操作使用原始仓储
user, err := s.userRepo.GetByID(ctx, userID)

// ✅ 事务中使用 WithTx
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
	tx := getTxFromContext(ctx)
	txRepo := s.userRepo.WithTx(tx)
	return txRepo.Update(ctx, user)
})
```

---

## 📊 测试示例

### 单元测试 - Mock UnitOfWork

```go
type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Begin(ctx context.Context) (transaction.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockUnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestTenantService_RegisterUserAndJoinTenant(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockUow := new(MockUnitOfWork)
	
	service := NewTenantApplicationService(mockUow, mockUserRepo, ...)
	
	mockUow.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	
	user, err := service.RegisterUserAndJoinTenant(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockUow.AssertExpectations(t)
}
```

### 集成测试 - 真实事务

```go
func TestRegisterUserAndJoinTenant_Integration(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)
	
	userRepo := repo.NewUserDAORepository(db)
	tenantRepo := repo.NewTenantDAORepository(db)
	memberRepo := repo.NewTenantMemberDAORepository(db)
	
	service := NewTenantApplicationService(uow, userRepo, tenantRepo, memberRepo)
	
	// 执行跨聚合根操作
	user, err := service.RegisterUserAndJoinTenant(
		context.Background(),
		&RegisterAndJoinRequest{
			Email:    "test@example.com",
			Password: "Password123!",
			Nickname: "测试用户",
			TenantID: tenantID,
			Role:     sharedEntity.RoleMember,
		},
	)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	
	// 验证数据一致性
	var userCount, memberCount int64
	db.Model(&model.User{}).Where("id = ?", user.ID.String()).Count(&userCount)
	db.Model(&model.TenantMember{}).Where("user_id = ?", user.ID.String()).Count(&memberCount)
	
	assert.Equal(t, int64(1), userCount)
	assert.Equal(t, int64(1), memberCount)
}
```

---

## 🎯 最佳实践

### 1. 在 Application Service 层管理事务

```go
// ✅ 推荐：Application Service 层
type ApplicationService struct {
	uow        transaction.UnitOfWork
	userRepo   repository.UserRepository
	// ...
}

func (s *ApplicationService) ComplexOperation(ctx context.Context) error {
	return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 使用 WithTx 切换仓储
		return nil
	})
}

// ❌ 不推荐：Domain Service 层管理事务
type DomainService struct {
	uow transaction.UnitOfWork // 不应该！
}
```

### 2. 保持事务简短

```go
// ✅ 推荐：短事务
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
	tx := getTxFromContext(ctx)
	userRepo := s.userRepo.WithTx(tx)
	memberRepo := s.memberRepo.WithTx(tx)
	
	// 快速完成
	user, _ := userRepo.GetByID(ctx, userID)
	member, _ := memberRepo.GetByUserAndTenant(ctx, userID, tenantID)
	return nil
})

// ❌ 不推荐：长事务（包含外部 API 调用）
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
	// 数据库操作
	tx := getTxFromContext(ctx)
	userRepo := s.userRepo.WithTx(tx)
	user, _ := userRepo.GetByID(ctx, userID)
	
	// ❌ 外部 API 调用（慢！）
	sendEmail(user.Email) // 耗时 2 秒
	
	// 更多数据库操作
	return nil
})
```

### 3. 错误处理

```go
// ✅ 自动回滚
err := s.uow.WithTransaction(ctx, func(ctx context.Context) error {
	tx := getTxFromContext(ctx)
	userRepo := s.userRepo.WithTx(tx)
	
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return err // 自动回滚
	}
	
	// 业务验证失败也会回滚
	if !user.IsActive() {
		return errors.New("user is not active")
	}
	
	return nil
})

if err != nil {
	// 处理错误（已自动回滚）
	log.Printf("Transaction failed: %v", err)
}
```

---

## 📚 相关文档

- [UnitOfWork 使用指南](../transaction/README.md)
- [仓储事务支持](../persistence/gorm/repo/README.md)
- [DDD 架构重构计划](../../../docs/DDD_ARCHITECTURE_RESTRUCTURE_PLAN.md)

---

**更新时间**: 2026-03-08  
**状态**: Complete ✅
