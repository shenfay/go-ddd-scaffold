# 仓储事务支持使用指南

## 📋 概述

所有 Repository 现已支持 `WithTx` 方法，可以在 UnitOfWork 管理的事务中执行操作。

---

## 🎯 核心接口

### UserRepository

```go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error)
    CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
    
    // WithTx 返回使用指定事务的仓储实例
    WithTx(tx *gorm.DB) UserRepository
}
```

### TenantRepository

```go
type TenantRepository interface {
    Create(ctx context.Context, tenant *entity.Tenant) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
    Update(ctx context.Context, tenant *entity.Tenant) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    WithTx(tx *gorm.DB) TenantRepository
}
```

### TenantMemberRepository

```go
type TenantMemberRepository interface {
    Create(ctx context.Context, member *entity.TenantMember) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.TenantMember, error)
    GetByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) (*entity.TenantMember, error)
    Update(ctx context.Context, member *entity.TenantMember) error
    Delete(ctx context.Context, id uuid.UUID) error
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.TenantMember, error)
    ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.TenantMember, error)
    DeleteByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) error
    
    WithTx(tx *gorm.DB) TenantMemberRepository
}
```

---

## 💡 使用示例

### 场景 1: 用户注册并加入租户（跨聚合根事务）

```go
type TenantApplicationService struct {
    uow        transaction.UnitOfWork
    userRepo   repository.UserRepository
    tenantRepo repository.TenantRepository
    memberRepo repository.TenantMemberRepository
}

func (s *TenantApplicationService) RegisterUserAndJoinTenant(
    ctx context.Context,
    email, password, nickname string,
    tenantID uuid.UUID,
    role sharedEntity.UserRole,
) (*entity.User, error) {
    
    var createdUser *entity.User
    
    // 使用 UnitOfWork 保证原子性
    err := s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 步骤 1: 创建用户
        user, err := entity.NewUser(email, password, nickname)
        if err != nil {
            return err
        }
        
        // 获取当前事务的 DB
        tx := getTxFromContext(ctx)
        
        // 切换到事务仓储
        userRepo := s.userRepo.WithTx(tx)
        if err := userRepo.Create(ctx, user); err != nil {
            return err
        }
        
        // 步骤 2: 验证租户成员限制
        tenantRepo := s.tenantRepo.WithTx(tx)
        tenant, err := tenantRepo.GetByID(ctx, tenantID)
        if err != nil {
            return err
        }
        
        memberRepo := s.memberRepo.WithTx(tx)
        currentCount, _ := memberRepo.CountByTenant(ctx, tenantID)
        
        if err := tenant.ValidateMemberLimit(currentCount + 1); err != nil {
            return err
        }
        
        // 步骤 3: 创建租户成员关系
        member, err := entity.NewTenantMember(tenantID, user.ID, role)
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

---

### 场景 2: 批量更新用户信息

```go
func (s *UserService) BatchUpdateUserInfo(
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

## 🔧 实现细节

### WithTx 方法实现

```go
// UserDAORepository.WithTx
func (r *UserDAORepository) WithTx(tx *gorm.DB) repository.UserRepository {
    return &UserDAORepository{
        db:      tx,
        querier: dao.Use(tx),
    }
}
```

**关键点**:
- ✅ 返回新的仓储实例，使用传入的 tx
- ✅ 重新初始化 querier（Gen 生成的 DAO）
- ✅ 原始仓储不受影响（可重复使用）

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

### 单元测试 - Mock 仓储

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) WithTx(tx *gorm.DB) repository.UserRepository {
    // 测试中可以返回自身或另一个 mock
    return m
}

func TestApplicationService_WithTransaction(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockUow := new(transaction.MockUnitOfWork)
    
    service := NewApplicationService(mockUow, mockUserRepo, ...)
    
    mockUow.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
    mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    err := service.RegisterUser(ctx, email, password, nickname)
    
    assert.NoError(t, err)
    mockUserRepo.AssertExpectations(t)
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
    
    service := NewApplicationService(uow, userRepo, tenantRepo, memberRepo)
    
    // 执行跨聚合根操作
    user, err := service.RegisterUserAndJoinTenant(
        context.Background(),
        "test@example.com",
        "Password123!",
        "测试用户",
        tenantID,
        sharedEntity.RoleMember,
    )
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    
    // 验证数据一致性
    var count int64
    db.Model(&model.User{}).Where("id = ?", user.ID.String()).Count(&count)
    assert.Equal(t, int64(1), count)
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
    userRepo := s.userRepo.WithTx(getTxFromContext(ctx))
    memberRepo := s.memberRepo.WithTx(getTxFromContext(ctx))
    
    // 快速完成
    user, _ := userRepo.GetByID(ctx, userID)
    member, _ := memberRepo.GetByUserAndTenant(ctx, userID, tenantID)
    return nil
})

// ❌ 不推荐：长事务（包含外部 API 调用）
err := uow.WithTransaction(ctx, func(ctx context.Context) error {
    // 数据库操作
    userRepo := s.userRepo.WithTx(getTxFromContext(ctx))
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

- [UnitOfWork 使用指南](./transaction/README.md)
- [DDD 架构重构计划](../../../docs/DDD_ARCHITECTURE_RESTRUCTURE_PLAN.md)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

---

**更新时间**: 2026-03-08  
**状态**: Complete ✅
