# Repository-DAO 使用指南

## 概述

本文档说明了在 go-ddd-scaffold 项目中如何正确使用 Repository 层和 DAO 层的职责分工。

## 目录结构

```
backend/internal/infrastructure/persistence/
├── dao/                    # GORM/gen 自动生成的数据访问对象
│   ├── users.gen.go
│   ├── tenants.gen.go
│   └── ...
├── model/                  # GORM/gen 自动生成的数据库模型
│   ├── users.gen.go
│   ├── tenants.gen.go
│   └── ...
├── repository/             # 仓储实现（按领域分组）
│   ├── user_repository.go
│   ├── tenant_repository.go
│   └── ...
├── interfaces.go          # DB 接口抽象
└── postgres.go            # PostgreSQL 连接配置
```

**设计原则**：
- ✅ **职责分离**：DAO、Repository 各司其职
- ✅ **高内聚**：同一领域的仓储放在一起
- ✅ **易扩展**：新增领域时目录结构清晰

## 架构设计原则

### 分层职责划分

```
┌─────────────────────────────────────────────────────┐
│           Application Layer (Services)               │
│  - 应用服务                                          │
│  - 业务协调                                          │
└────────────────────┬────────────────────────────────┘
                     │ 使用领域对象
                     ▼
┌─────────────────────────────────────────────────────┐
│              Repository Layer (Domain)               │
│  - 领域对象转换 (DAO Model ↔ Domain Entity)          │
│  - 业务规则验证                                      │
│  - 领域事件管理                                      │
│  - 复杂查询封装（调用 DAO）                          │
└────────────────────┬────────────────────────────────┘
                     │ 调用
                     ▼
┌─────────────────────────────────────────────────────┐
│                DAO Layer (Generated)                 │
│  - 基础 CRUD（增删改查）                              │
│  - 类型安全查询                                     │
│  - 简单条件过滤                                      │
└─────────────────────────────────────────────────────┘
```

## 各层职责详解

### 1. DAO 层（数据访问对象）

**职责**：
- 提供类型安全的数据库操作接口
- 执行基础的 CRUD 操作
- 支持简单的条件查询和过滤

**特点**：
- ✅ 由 GORM/gen 自动生成
- ✅ 零手写 SQL
- ✅ 完全的字段名和类型检查
- ❌ 不包含任何业务逻辑

**示例**：
```go
// 创建用户
userModel := &model.User{
    ID:       1234567890,
    Username: "testuser",
    Email:    "test@example.com",
    Status:   1,
}
err := dao.User.WithContext(ctx).Create(userModel)

// 查询单个用户
user, err := dao.User.WithContext(ctx).
    Where(dao.User.ID.Eq(1234567890)).
    First()

// 条件查询
users, err := dao.User.WithContext(ctx).
    Where(dao.User.Status.Eq(1)).
    Where(dao.User.DeletedAt.IsNull()).
    Order(dao.User.CreatedAt.Desc()).
    Limit(10).
    Offset(0).
    Find()

// 统计
count, err := dao.User.WithContext(ctx).
    Where(dao.User.Status.Eq(1)).
    Count()
```

### 2. Repository 层（仓储）

**职责**：
- **领域对象转换**：DAO Model ↔ Domain Entity
- **业务规则验证**：确保领域对象的完整性
- **领域事件管理**：保存和发布领域事件
- **复杂查询封装**：组合多个 DAO 调用实现业务查询

**特点**：
- ✅ 包含领域知识
- ✅ 维护聚合根不变性
- ✅ 处理乐观锁和并发控制
- ✅ 领域事件持久化

#### 类型转换工具

项目使用 `pkg/util` 包提供统一的类型转换和指针操作函数：

```go
import "github.com/shenfay/go-ddd-scaffold/pkg/util"
```

**四类核心函数**：

```go
// 1. 类型转换（To 前缀，处理任意输入）
util.ToString(v interface{}) string
util.ToInt32(v interface{}) int32
util.ToBool(v interface{}) bool

// 2. 创建指针（类型名即函数名，类型必须匹配）
util.String(s string) *string
util.Int32(i int32) *int32
util.Bool(b bool) *bool

// 3. 获取值（Value 后缀，安全防护）
util.StringValue(s *string) string
util.Int32Value(i *int32) int32
util.BoolValue(b *bool) bool

// 4. 智能转换（根据值决定返回 nil 或指针）⭐
util.StringPtrNilIfEmpty(s string) *string   // "" → nil, "hello" → *string("hello")
util.Int16PtrNilIfZero(i int16) *int16       // 0 → nil, 10 → *int16(10)
util.Int32PtrNilIfZero(i int32) *int32       // 0 → nil, 10 → *int32(10)
util.Int64PtrNilIfZero(i int64) *int64       // 0 → nil, 10 → *int64(10)
```

#### 智能转换函数的使用场景

**场景 1：数据库字段可 NULL**

当数据库字段允许为 NULL 时，使用智能转换函数区分"空值"和"零值"：

```go
// users.display_name (VARCHAR, 可 NULL)
DisplayName: util.StringPtrNilIfEmpty(displayName)
// "" → nil (数据库存储为 NULL)
// "张三" → *string("张三") (数据库存储为 '张三')

// users.gender (SMALLINT, 可 NULL)  
Gender: util.Int16PtrNilIfZero(int16(u.Gender()))
// 0 → nil (数据库存储为 NULL)
// 1 → *int16(1) (数据库存储为 1)
```

**场景 2：数据库字段 NOT NULL**

当数据库字段不允许为 NULL 时，直接使用简洁的指针创建函数：

```go
// login_logs.login_type (VARCHAR, NOT NULL)
LoginType: util.String(log.LoginType)
// 直接创建指针，不检查空值

// domain_events.processed (BOOLEAN, NOT NULL)
Processed: util.Bool(false)
// 直接创建指针
```

**场景 3：从 DAO Model 转换为 Domain Entity**

使用 Value 后缀函数安全解引用：

```go
// DAO Model → Domain Entity
domainUser := &user.User{
    DisplayName: util.StringValue(daoModel.DisplayName),
    // nil → "", *string("张三") → "张三"
    
    Gender: user.Gender(util.Int32Value(daoModel.Gender)),
    // nil → 0, *int32(1) → 1
}
```

**示例**：
```go
// 保存用户（包含创建和更新）
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
    if u.Version() == 0 {
        return r.insert(ctx, u)  // 创建新用户
    }
    return r.update(ctx, u)      // 更新现有用户（带乐观锁）
}

// 查找用户（返回领域对象）
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
    // 1. 使用 DAO 查询
    userModel, err := dao.User.WithContext(ctx).
        Where(dao.User.ID.Eq(id.Int64())).
        Where(dao.User.DeletedAt.IsNull()).
        First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ddd.ErrAggregateNotFound
        }
        return nil, err
    }
    
    // 2. 转换为领域对象（核心职责）
    return r.toDomain(userModel), nil
}

// 删除用户（软删除 + 版本号递增）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id user.UserID) error {
    // 使用 GORM.Delete 进行软删除
    result, err := dao.User.WithContext(ctx).Delete(userModel)
    if err != nil {
        return err
    }
    
    if result.RowsAffected == 0 {
        return ddd.ErrAggregateNotFound
    }
    
    return nil
}
```

### 3. Projector 层（投影器 - CQRS 读侧）

**职责**：
- 监听领域事件
- 将事件投影到读模型（user_read_model）
- 优化查询性能

**特点**：
- ✅ 使用原生 SQL（性能优先）
- ✅ 最终一致性
- ✅ 针对查询优化

**示例**：
```go
// 处理用户注册事件
func (p *UserProjectorImpl) handleUserRegistered(ctx context.Context, event *userdomain.UserRegisteredEvent) error {
    query := `
        INSERT INTO user_read_model (
            user_id, username, email, status, created_at, updated_at, login_count
        ) VALUES ($1, $2, $3, $4, NOW(), NOW(), 0)
    `
    
    _, err := p.db.ExecContext(ctx, query,
        event.UserID.Int64(),
        event.Username,
        event.Email,
        int(user.UserStatusPending),
    )
    
    return err
}
```

## 最佳实践

### 1. 领域对象转换

**推荐做法**：
```go
// toDomain: DAO Model → Domain Entity
func (r *UserRepositoryImpl) toDomain(model *model.User) *user.User {
    // 处理可能的 nil 值
    loginCount := 0
    if model.LoginCount != nil {
        loginCount = int(*model.LoginCount)
    }
    
    gender := user.UserGender(0)
    if model.Gender != nil {
        gender = user.UserGender(*model.Gender)
    }
    
    // 使用 Builder 模式构建领域对象
    return user.NewUserBuilder().
        WithID(model.ID).
        WithUsername(model.Username).
        WithEmail(model.Email).
        WithPasswordHash(model.PasswordHash).
        WithStatus(user.UserStatus(model.Status)).
        WithGender(gender).
        WithDisplayName(displayName).
        WithPhoneNumber(phoneNumber).
        WithAvatarURL(avatarURL).
        WithLastLoginAt(model.LastLoginAt).
        WithLoginCount(loginCount).
        WithFailedAttempts(failedAttempts).
        WithLockedUntil(model.LockedUntil).
        WithVersion(version).
        WithTimestamps(createdAt, updatedAt).
        Build()
}

// fromDomain: Domain Entity → DAO Model
func (r *UserRepositoryImpl) fromDomain(u *user.User) *model.User {
	displayName := u.DisplayName()
	phoneNumber := u.PhoneNumber()
	avatarURL := u.AvatarURL()
	loginCount := int(u.LoginCount())
	failedAttempts := int(u.FailedAttempts())
	version := int(u.Version())
	
	return &model.User{
		ID:             u.ID().(user.UserID).Int64(),
		Username:       u.Username().Value(),
		Email:          u.Email().Value(),
		PasswordHash:   u.Password().Value(),
		Status:         int16(u.Status()),
		DisplayName:    util.StringPtrNilIfEmpty(displayName),
		Gender:         util.Int16PtrNilIfZero(int16(u.Gender())),
		PhoneNumber:    util.StringPtrNilIfEmpty(phoneNumber),
		AvatarURL:      util.StringPtrNilIfEmpty(avatarURL),
		LastLoginAt:    u.LastLoginAt(),
		LoginCount:     util.Int32PtrNilIfZero(int32(loginCount)),
		FailedAttempts: util.Int32PtrNilIfZero(int32(failedAttempts)),
		LockedUntil:    u.LockedUntil(),
		Version:        util.Int32PtrNilIfZero(int32(version)),
		CreatedAt:      util.Time(u.CreatedAt()),
		UpdatedAt:      util.Time(u.UpdatedAt()),
	}
}

```

### 2. 乐观锁实现

GORM 自动支持乐观锁，只需确保模型包含 `version` 字段：

```go
type User struct {
    Version *int32 `gorm:"column:version;type:integer"`
    // ... 其他字段
}

// 更新时自动检查版本号
result, err := dao.User.WithContext(ctx).
    Where(dao.User.ID.Eq(userID)).
    Updates(userModel)  // GORM 会自动生成 WHERE version = ?

if result.RowsAffected == 0 {
    // 版本号为 0 表示被其他事务更新了
    return ddd.NewConcurrencyError(...)
}
```

### 3. 领域事件保存

**重要**：领域事件必须保存在同一事务中：

```go
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
    // 1. 保存聚合根
    if u.Version() == 0 {
        err := dao.User.WithContext(ctx).Create(r.fromDomain(u))
        if err != nil {
            return err
        }
    } else {
        result, err := dao.User.WithContext(ctx).
            Where(dao.User.ID.Eq(u.ID().(user.UserID).Int64())).
            Updates(r.fromDomain(u))
        if err != nil {
            return err
        }
        
        // 乐观锁检查
        if result.RowsAffected == 0 {
            return ddd.NewConcurrencyError(...)
        }
    }
    
    // 2. 保存领域事件（在同一事务中）
    return r.saveEvents(ctx, u)
}

func (r *UserRepositoryImpl) saveEvents(ctx context.Context, u *user.User) error {
    events := u.GetUncommittedEvents()
    if len(events) == 0 {
        return nil
    }
    
    for _, event := range events {
        eventData, err := json.Marshal(event)
        if err != nil {
            return err
        }
        
        domainEvent := &model.DomainEvent{
            AggregateID:   u.ID().(user.UserID).String(),
            AggregateType: "user",
            EventType:     event.EventName(),
            EventVersion:  int32(event.Version()),
            EventData:     string(eventData),
            OccurredOn:    event.OccurredOn(),
            Processed:     util.Bool(false),
        }
        
        err = dao.DomainEvent.WithContext(ctx).Create(domainEvent)
        if err != nil {
            return err
        }
    }
    
    // 清除已保存的事件
    u.ClearUncommittedEvents()
    return nil
}
```

### 4. 事务处理

在 Service 层统一处理事务：

```go
func (s *UserService) CreateUser(cmd CreateUserCommand) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 传递事务中的 DB
        repo := persistence.NewUserRepository(tx)
        
        // 创建用户
        user, _ := user.NewUser(...)
        if err := repo.Save(ctx, user); err != nil {
            return err
        }
        
        // 其他操作...
        
        return nil
    })
}
```

## 常见错误

### ❌ 错误：在 DAO 层混入业务逻辑

``go
// 错误示例
func CreateUserWithBusinessLogic(ctx context.Context, cmd CreateUserCommand) error {
    // 不应该在 DAO 层做业务验证
    if exists, _ := dao.User.Where(...).First(); exists {
        return errors.New("user already exists")
    }
    // ...
}
```

### ✅ 正确：在 Repository 层处理业务逻辑

``go
// 正确示例
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
    // 业务规则验证
    if u.Username().Value() == "admin" {
        return errors.New("username 'admin' is reserved")
    }
    
    // 使用 DAO 保存
    return dao.User.WithContext(ctx).Create(r.fromDomain(u))
}
```

### ❌ 错误：直接暴露 DAO Model

``go
// 错误示例
func GetUser(ctx context.Context, id int64) (*model.User, error) {
    return dao.User.WithContext(ctx).Where(dao.User.ID.Eq(id)).First()
}
```

### ✅ 正确：返回 Domain Entity

``go
// 正确示例
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
    userModel, err := dao.User.WithContext(ctx).Where(dao.User.ID.Eq(id.Int64())).First()
    if err != nil {
        return nil, err
    }
    return r.toDomain(userModel), nil  // 转换为领域对象
}
```

## 总结

### Repository 层最佳实践

#### 1. 使用统一的类型转换工具

✅ **推荐**：使用 `pkg/util` 包
```go
import "github.com/shenfay/go-ddd-scaffold/pkg/util"

// 智能转换（可 NULL 字段）
DisplayName: util.StringPtrNilIfEmpty(displayName)
Gender: util.Int16PtrNilIfZero(int16(u.Gender()))

// 简洁转换（NOT NULL 字段）
LoginType: util.String(log.LoginType)
Processed: util.Bool(false)

// 安全解引用
domainName := util.StringValue(daoModel.DisplayName)
domainGender := util.Int32Value(daoModel.Gender)
```

❌ **不推荐**：内联函数
```go
// 代码冗长，不易维护
DisplayName: func() *string {
    if displayName != "" {
        return &displayName
    }
    return nil
}()
```

#### 2. 函数选择指南

| 场景 | 推荐函数 | 示例 |
|------|---------|------|
| **数据库字段可 NULL** | `StringPtrNilIfEmpty`, `Int16PtrNilIfZero` | `util.StringPtrNilIfEmpty("")` → `nil` |
| **数据库字段 NOT NULL** | `String`, `Int32`, `Bool` | `util.String("hello")` → `*string("hello")` |
| **从 DAO Model 读取** | `StringValue`, `Int32Value` | `util.StringValue(nil)` → `""` |
| **任意类型转换** | `ToString`, `ToInt32` | `util.ToInt32("123")` → `123` |

#### 3. 已更新的文件

| 文件 | 改进 | 代码行数变化 |
|------|------|-------------|
| `user_repository.go` | ✅ 使用智能转换函数 | ~35 行 → ~20 行 (-43%) |
| `audit_log_repository.go` | ✅ 使用智能转换函数 | ~50 行 → ~20 行 (-60%) |
| `login_log_repository.go` | ✅ 使用简洁转换函数 | 保持不变 |

---

### 架构分层总结

| 层级 | 职责 | 技术栈 | 是否包含业务逻辑 |
|------|------|--------|-----------------|
| **DAO** | 基础 CRUD、类型安全查询 | GORM/gen | ❌ 否 |
| **Repository** | 领域转换、业务规则、事件管理 | GORM + 原生 SQL | ✅ 是 |
| **Projector** | CQRS 读模型投影 | 原生 SQL | ❌ 否（仅数据同步） |

遵循这个分层架构，可以确保：
- ✅ 领域模型的纯净性
- ✅ 代码的可维护性
- ✅ 职责的清晰分离
- ✅ 便于单元测试
