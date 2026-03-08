# UnitOfWork 事务管理使用指南

## 📋 概述

本项目已实现完整的 UnitOfWork（工作单元）模式，用于保证跨仓储操作的原子性和一致性。

---

## 🏗️ 架构设计

### 接口定义层（Domain 层）

**位置**: `internal/domain/shared/transaction/unit_of_work.go`

```go
// Transaction 事务接口
type Transaction interface {
    Commit() error
    Rollback() error
    GetDB() *gorm.DB
}

// UnitOfWork 工作单元接口
type UnitOfWork interface {
    Begin(ctx context.Context) (Transaction, error)
}
```

### 基础设施实现层（Infrastructure 层）

**位置**: `internal/infrastructure/transaction/unit_of_work.go`

```go
// gormUnitOfWork GORM 工作单元实现
type gormUnitOfWork struct {
    db *gorm.DB
}

// NewGormUnitOfWork 创建 GORM 工作单元实例
func NewGormUnitOfWork(db *gorm.DB) UnitOfWork {
    return &gormUnitOfWork{db: db}
}
```

---

## 💡 使用场景

### 场景 1: 用户注册并加入租户（跨聚合根事务）

```go
// Application Service
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
        // 1. 创建用户
        user, err := entity.NewUser(email, password, nickname)
        if err != nil {
            return err
        }
        
        // 2. 保存到用户仓储
        if err := s.userRepo.Create(ctx, user); err != nil {
            return err
        }
        
        // 3. 获取租户信息
        tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
        if err != nil {
            return err
        }
        
        // 4. 验证成员限制
        currentCount, _ := s.memberRepo.CountByTenant(ctx, tenantID)
        if err := tenant.ValidateMemberLimit(currentCount + 1); err != nil {
            return err
        }
        
        // 5. 创建租户成员关系
        member, err := entity.NewTenantMember(tenantID, user.ID, role)
        if err != nil {
            return err
        }
        
        if err := s.memberRepo.Create(ctx, member); err != nil {
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

### 场景 2: 手动管理事务

```go
func (s *UserService) ComplexOperation(ctx context.Context) error {
    // 1. 开启事务
    tx, err := s.uow.Begin(ctx)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            // panic 时回滚
            _ = tx.Rollback()
            panic(p)
        }
    }()
    
    // 2. 在事务中执行操作
    // 使用 tx.GetDB() 代替原始 db
    userDB := s.userRepo.WithTx(tx.GetDB())
    
    user, err := userDB.GetByID(ctx, userID)
    if err != nil {
        _ = tx.Rollback()
        return err
    }
    
    user.UpdateNickname(newNickname)
    
    if err := userDB.Update(ctx, user); err != nil {
        _ = tx.Rollback()
        return err
    }
    
    // 3. 提交事务
    if err := tx.Commit(); err != nil {
        return err
    }
    
    return nil
}
```

---

### 场景 3: 仓储的 WithTx 支持

```go
// UserRepository 接口扩展
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // 新增：支持事务的变体
    WithTx(tx *gorm.DB) UserRepository
}

// 实现
type userDAORepository struct {
    db *gorm.DB
}

func (r *userDAORepository) WithTx(tx *gorm.DB) repository.UserRepository {
    return &userDAORepository{db: tx}
}

// 使用
func (s *UserService) UpdateUserWithAudit(ctx context.Context, userID uuid.UUID, newNickname string) error {
    return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 切换到事务 DB
        userRepo := s.userRepo.WithTx(tx.GetDB())
        auditRepo := s.auditRepo.WithTx(tx.GetDB())
        
        // 所有操作共享同一事务
        user, _ := userRepo.GetByID(ctx, userID)
        user.UpdateNickname(newNickname)
        userRepo.Update(ctx, user)
        
        auditRepo.Log(ctx, "用户昵称更新")
        
        return nil
    })
}
```

---

## 🔧 Wire 集成

### providers.go 配置

```go
// internal/infrastructure/wire/providers.go

import (
    "go-ddd-scaffold/internal/infrastructure/transaction"
)

var TransactionSet = wire.NewSet(
    transaction.NewGormUnitOfWork,
    wire.Bind(new(transaction.UnitOfWork), new(*transaction.gormUnitOfWork)),
)

// 在 injector 中注入
func InitializeApp(db *gorm.DB) (*App, error) {
    wire.Build(
        // ... 其他依赖
        TransactionSet,
        NewApplicationService,
    )
    return nil, nil
}
```

---

## ✅ 最佳实践

### 1. 事务边界

```go
// ✅ 推荐：在 Application Service 层管理事务
func (s *ApplicationService) CreateUser(ctx context.Context, dto CreateUserDTO) error {
    return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 业务逻辑
        return nil
    })
}

// ❌ 不推荐：在 Domain Service 或 Repository 层管理事务
func (d *DomainService) DoSomething(ctx context.Context) error {
    tx, _ := d.uow.Begin(ctx) // 错误的位置
    // ...
}
```

### 2. 错误处理

```go
// ✅ 推荐：使用 WithTransaction 自动处理 commit/rollback
err := s.uow.WithTransaction(ctx, func(ctx context.Context) error {
    // 任何错误都会自动回滚
    return errors.New("some error")
})

// ❌ 不推荐：忘记处理异常
func (s *Service) BadExample(ctx context.Context) error {
    tx, _ := s.uow.Begin(ctx)
    // 如果发生 panic，事务不会回滚
    doSomething() // panic!
    tx.Commit()
}
```

### 3. 嵌套事务

```go
// ⚠️ 注意：GORM 不支持真正的嵌套事务
// 应该使用 SavePoint（如果数据库支持）

func (s *Service) NestedExample(ctx context.Context) error {
    return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 外层事务
        
        // 内层逻辑 - 不要再次开启事务
        s.doInnerWork(ctx)
        
        return nil
    })
}
```

---

## 📊 测试示例

### UnitOfWork 单元测试

```go
func TestUnitOfWork_Commit(t *testing.T) {
    db := setupTestDB(t)
    uow := transaction.NewGormUnitOfWork(db)
    
    err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
        // 执行操作
        return nil
    })
    
    assert.NoError(t, err)
}

func TestUnitOfWork_Rollback(t *testing.T) {
    db := setupTestDB(t)
    uow := transaction.NewGormUnitOfWork(db)
    
    err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
        return errors.New("force rollback")
    })
    
    assert.Error(t, err)
    assert.Equal(t, "force rollback", err.Error())
}
```

---

## 📈 性能考虑

1. **事务粒度**: 保持事务尽可能短
2. **避免长事务**: 不要在事务中调用外部 API
3. **并发控制**: 使用乐观锁或悲观锁防止并发冲突

---

## 🔍 常见问题

### Q: 什么时候使用 UnitOfWork？
A: 当需要同时修改多个聚合根或跨多个仓储操作时。

### Q: 单个仓储操作需要事务吗？
A: 不需要，GORM 会自动为单个操作开启事务。

### Q: 如何处理分布式事务？
A: 本实现仅支持本地事务，分布式事务需要使用 Saga 或 TCC 模式。

---

## 📚 相关文档

- [Domain Driven Design](https://martinfowler.com/tags/domain_driven_design.html)
- [Unit of Work Pattern](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [GORM Transactions](https://gorm.io/docs/transactions.html)
