# Repository 指南

本文档详细介绍 Repository（仓储）模式在 Go DDD Scaffold 中的实现和使用。

## 📋 什么是 Repository？

### 核心概念

**Repository（仓储）** 是领域层定义的接口，用于抽象数据持久化操作。

```
┌─────────────────────────────────────┐
│     Domain Layer                    │
│                                     │
│  type UserRepository interface {    │  ← Domain 定义接口
│      FindByID(id UserID) (*User, error)
│      Save(user *User) error
│  }
└─────────────────────────────────────┘
                ↑
┌─────────────────────────────────────┤
│   Infrastructure Layer              │
│                                     │
│  type userRepositoryImpl struct {   │  ← Infrastructure 实现
│      db *gorm.DB
│  }
│                                     │
│  func (r *userRepositoryImpl)       │
│      FindByID(...) (*User, error) { │
│      // 具体实现
│  }
└─────────────────────────────────────┘
```

### Repository 的价值

1. **解耦** - Domain 不依赖具体持久化技术
2. **可测试** - 可以轻松 Mock
3. **可替换** - 可以从 MySQL 切换到 PostgreSQL
4. **统一** - 统一的访问接口

---

## 🏗️ Repository 架构

### 目录结构

```
backend/internal/
├── domain/
│   └── user/
│       └── repository/              # ⭐ Domain 层定义接口
│           └── user_repository.go
│
└── infrastructure/
    └── persistence/
        └── repository/              # ⭐ Infrastructure 层实现
            ├── user_repository.go
            └── tenant_repository.go
```

### 分层职责

| 层级 | 职责 | 示例 |
|------|------|------|
| **Domain 层** | 定义 Repository 接口 | `type UserRepository interface {}` |
| **Infrastructure 层** | 实现 Repository 接口 | `type userRepositoryImpl struct {}` |
| **Application 层** | 使用 Repository 接口 | `userRepo.FindByID(...)` |
| **Module 层** | 组装依赖 | `NewUserRepository(db, daoQuery)` |

---

## 🎯 Repository 设计模式

### 1. 标准 Repository 接口

```go
// domain/user/repository/user_repository.go
package repository

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserRepository 用户仓储接口
type UserRepository interface {
    // === 基本 CRUD ===
    
    // FindByID 根据 ID 查找用户
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    
    // FindByUsername 根据用户名查找用户
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    
    // FindByEmail 根据邮箱查找用户
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    
    // Save 保存用户（创建或更新）
    Save(ctx context.Context, user *aggregate.User) error
    
    // Delete 删除用户
    Delete(ctx context.Context, id vo.UserID) error
    
    // === 查询方法 ===
    
    // ExistsByEmail 检查邮箱是否已存在
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    
    // CountByStatus 统计指定状态的用户数
    CountByStatus(ctx context.Context, status vo.UserStatus) (int, error)
    
    // ListActiveUsers 获取活跃用户列表
    ListActiveUsers(ctx context.Context, limit int) ([]*aggregate.User, error)
}
```

### 2. 通用 Repository 模式

```go
// domain/shared/repository/base_repository.go
package repository

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// BaseRepository 通用仓储接口
type BaseRepository[T any, ID any] interface {
    FindByID(ctx context.Context, id ID) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    Save(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id ID) error
}

// 使用示例
type UserRepository interface {
    BaseRepository[aggregate.User, vo.UserID]
    // 添加用户特有的方法
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
}
```

**优点：** 减少重复代码  
**缺点：** Go 泛型较复杂，不如显式定义清晰

**推荐：** 为每个聚合根显式定义 Repository 接口

---

## 💻 Repository 实现详解

### 步骤 1：定义 Domain 接口

```go
// domain/user/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// 定义错误
var (
    ErrUserNotFound = errors.New("user not found")
)

// UserRepository 用户仓储接口
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
    Delete(ctx context.Context, id vo.UserID) error
}
```

### 步骤 2：实现 Infrastructure Repository

```go
// infrastructure/persistence/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "fmt"
    "gorm.io/gorm"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// userRepositoryImpl 用户仓储实现
type userRepositoryImpl struct {
    db       *gorm.DB
    daoQuery *dao_query.Query
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB, daoQuery *dao_query.Query) repository.UserRepository {
    return &userRepositoryImpl{
        db:       db,
        daoQuery: daoQuery,
    }
}

// FindByID 根据 ID 查找用户
func (r *userRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    // 1. 使用 GORM Gen 查询
    dao, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.ID.Eq(id.Value())).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, kernel.ErrAggregateNotFound
        }
        return nil, fmt.Errorf("find user by id failed: %w", err)
    }
    
    // 2. DAO → Domain 转换
    return r.toDomain(dao)
}

// FindByUsername 根据用户名查找用户
func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
    dao, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.Username.Eq(username)).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrUserNotFound
        }
        return nil, err
    }
    
    return r.toDomain(dao)
}

// FindByEmail 根据邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
    dao, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.Email.Eq(email)).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrUserNotFound
        }
        return nil, err
    }
    
    return r.toDomain(dao)
}

// Save 保存用户
func (r *userRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 1. 保存或更新用户
    err := r.saveUser(tx, user)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("save user failed: %w", err)
    }
    
    // 2. 保存领域事件（Outbox Pattern）
    events := user.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("save domain event failed: %w", err)
        }
    }
    
    return tx.Commit()
}

// Delete 删除用户
func (r *userRepositoryImpl) Delete(ctx context.Context, id vo.UserID) error {
    _, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.ID.Eq(id.Value())).Delete()
    if err != nil {
        return fmt.Errorf("delete user failed: %w", err)
    }
    return nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    count, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.Email.Eq(email)).Count()
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// toDomain DAO 转换为 Domain
func (r *userRepositoryImpl) toDomain(dao *dao.User) (*aggregate.User, error) {
    // 创建值对象
    username, err := vo.NewUsername(dao.Username)
    if err != nil {
        return nil, err
    }
    
    email, err := vo.NewEmail(dao.Email)
    if err != nil {
        return nil, err
    }
    
    // 重建聚合根
    user := &aggregate.User{
        Entity: kernel.ReconstructEntity(dao.ID, dao.CreatedAt, dao.UpdatedAt),
        username: username,
        email:    email,
        status:   vo.UserStatus(dao.Status),
    }
    
    return user, nil
}

// saveUser Domain 转换为 DAO 并保存
func (r *userRepositoryImpl) saveUser(tx *gorm.DB, user *aggregate.User) error {
    dao := &dao.User{
        ID:        user.ID().Value(),
        Username:  user.Username().String(),
        Email:     user.Email().String(),
        Status:    string(user.Status()),
        CreatedAt: user.CreatedAt(),
        UpdatedAt: user.UpdatedAt(),
    }
    
    // 如果 ID 为 0，则是新建
    if user.ID().Value() == 0 {
        return tx.Create(dao).Error
    }
    
    // 否则更新
    return tx.Save(dao).Error
}

// saveEvent 保存领域事件
func (r *userRepositoryImpl) saveEvent(tx *gorm.DB, event kernel.DomainEvent) error {
    eventDAO := &dao.DomainEvent{
        ID:            snowflake.Generate(),
        EventType:     event.Type(),
        AggregateType: event.AggregateType(),
        AggregateID:   fmt.Sprintf("%d", event.AggregateID()),
        EventData:     toJSON(event),
        OccurredAt:    event.Timestamp(),
    }
    
    return tx.Create(eventDAO).Error
}

// toJSON 序列化为 JSON
func toJSON(v interface{}) []byte {
    data, _ := json.Marshal(v)
    return data
}
```

---

## 🔄 Unit of Work 模式

### 什么是 Unit of Work？

Unit of Work（工作单元）用于维护业务操作的事务一致性。

```go
// application/unit_of_work.go
package application

import (
    "context"
    "gorm.io/gorm"
)

// UnitOfWork 工作单元
type UnitOfWork struct {
    db *gorm.DB
}

// NewUnitOfWork 创建工作单元
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
    return &UnitOfWork{db: db}
}

// Transaction 执行事务
func (uow *UnitOfWork) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(ctx)
    })
}

// WithTransaction 返回带事务的 DB
func (uow *UnitOfWork) WithTransaction(ctx context.Context) (*gorm.DB, error) {
    tx := uow.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }
    return tx.WithContext(ctx), nil
}
```

### 使用示例

```go
// application/user/service.go
func (s *UserService) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*User, error) {
    // 使用 Unit of Work 保证事务一致性
    err := s.uow.Transaction(ctx, func(ctx context.Context) error {
        // 1. 创建用户
        user, err := aggregate.NewUser(cmd.Username, cmd.Email, cmd.Password)
        if err != nil {
            return err
        }
        
        // 2. 保存用户
        err = s.userRepo.Save(ctx, user)
        if err != nil {
            return err
        }
        
        // 3. 初始化租户配置
        err = s.tenantService.InitializeTenant(ctx, user.ID())
        if err != nil {
            return err
        }
        
        // 4. 发送欢迎邮件
        err = s.emailService.SendWelcomeEmail(ctx, user.Email())
        if err != nil {
            s.logger.Warn("send welcome email failed", zap.Error(err))
            // 不返回错误，继续执行
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

---

## 🎯 Repository 最佳实践

### 1. 返回领域对象，而非 DAO

```go
// ✅ 正确
func (r *userRepositoryImpl) FindByID(id int64) (*aggregate.User, error) {
    // 返回 Domain 对象
}

// ❌ 错误
func (r *userRepositoryImpl) FindByID(id int64) (*dao.User, error) {
    // 不应该暴露 DAO
}
```

### 2. 使用值对象作为参数

```go
// ✅ 正确
func (r *userRepositoryImpl) FindByID(id vo.UserID) (*aggregate.User, error) {
    // 使用值对象
}

// ❌ 错误
func (r *userRepositoryImpl) FindByID(id int64) (*aggregate.User, error) {
    // 原始类型不安全
}
```

### 3. 处理并发更新

```go
// 使用乐观锁
func (r *userRepositoryImpl) Save(user *aggregate.User) error {
    result := r.db.Model(&dao.User{}).
        Where("id = ? AND version = ?", user.ID(), user.Version()).
        UpdateColumn("version", user.Version()+1)
    
    if result.RowsAffected == 0 {
        return kernel.ErrOptimisticLockFailed
    }
    
    return nil
}
```

### 4. 实现缓存 Repository

```go
// 装饰器模式实现缓存
type CachedUserRepository struct {
    base   repository.UserRepository
    cache  *redis.Client
}

func (r *CachedUserRepository) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    // 1. 尝试从缓存读取
    key := fmt.Sprintf("user:%d", id.Value())
    cached, err := r.cache.Get(ctx, key).Bytes()
    if err == nil {
        user := &aggregate.User{}
        json.Unmarshal(cached, user)
        return user, nil
    }
    
    // 2. 缓存未命中，从数据库读取
    user, err := r.base.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    data, _ := json.Marshal(user)
    r.cache.Set(ctx, key, data, 5*time.Minute)
    
    return user, nil
}
```

### 5. 实现软删除

```go
// 支持软删除的 Repository
func (r *userRepositoryImpl) Delete(ctx context.Context, id vo.UserID) error {
    // 使用 DeletedAt 字段
    err := r.daoQuery.User.WithContext(ctx).
        Where(r.daoQuery.User.ID.Eq(id.Value())).
        Update("deleted_at", gorm.NowFunc()).Error
    
    if err != nil {
        return err
    }
    
    return nil
}

// 查询时自动过滤已删除
func (r *userRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    dao, err := r.daoQuery.User.WithContext(ctx).
        Where(r.daoQuery.User.ID.Eq(id.Value())).
        Where(r.daoQuery.User.DeletedAt.IsNull()).  // ← 过滤已删除
        First()
    // ...
}
```

---

## ✅ Repository 检查清单

### 设计检查

- [ ] Repository 接口定义在 Domain 层
- [ ] Repository 实现在 Infrastructure 层
- [ ] 返回领域对象，而非 DAO
- [ ] 使用值对象作为参数
- [ ] 包含基本的 CRUD 方法
- [ ] 包含业务查询方法

### 实现检查

- [ ] 正确处理记录不存在的情况
- [ ] 使用事务保证一致性
- [ ] 保存领域事件（Outbox Pattern）
- [ ] DAO ↔ Domain 转换正确
- [ ] 错误信息清晰明确

### 编译检查

```bash
# 检查 Repository 实现是否正确
go build ./...

# 运行 Repository 相关测试
go test ./infrastructure/persistence/repository/... -v
```

---

## 📚 参考资源

- [Domain-Driven Design](https://martinfowler.com/bliki/DDD.html)
- [Implementing Repositories](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer)
- [Unit of Work 模式](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
