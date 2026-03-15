# Go DDD 设计指南

## 文档概述

本文档详细阐述了在 go-ddd-scaffold 项目中如何实现标准的 **DDD（领域驱动设计）+ Clean Architecture（洁净架构）** 架构模式。

**重要说明：** 本项目采用**纯 DDD 架构**，而非 CQRS 模式。我们移除了读写模型分离、Projector 等 CQRS 组件，专注于领域模型的统一实现。

---

## DDD 架构模式详解

### 核心设计理念

DDD（Domain-Driven Design）将业务领域作为软件设计的核心，通过以下模式来组织代码：

```
┌─────────────────────────────────────┐
│   Interfaces Layer (接口层)          │
│   - HTTP Handler / Controller       │
│   - Request/Response DTOs           │
│   - Mapper (DTO ↔ Command/Query)    │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Application Layer (应用层)         │
│   - Application Services            │
│   - Commands & Queries              │
│   - 协调领域对象完成用例             │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Domain Layer (领域层) ⭐           │
│   - Aggregates (聚合根)             │
│   - Entities (实体)                 │
│   - Value Objects (值对象)          │
│   - Domain Events (领域事件)        │
│   - Domain Services (领域服务)      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Infrastructure Layer (基础设施层)    │
│ - Repository Implementations        │
│ - Database (GORM)                   │
│ - External Services (JWT, Email...) │
└─────────────────────────────────────┘
```

---

## 领域建模核心要素

### 1. 聚合根设计原则

#### 什么是聚合根？

聚合根是领域对象的根实体，它负责维护聚合内的一致性规则。

**示例：User 聚合根**
```go
package model

type User struct {
    ddd.BaseEntity
    
    // 值对象
    username    *UserName
    email       *Email
    password    *HashedPassword
    
    // 属性
    status      UserStatus
    displayName string
    createdAt   time.Time
    updatedAt   time.Time
    
    // 领域事件队列
    uncommittedEvents []ddd.DomainEvent
}
```

#### 聚合边界确定

**✅ 正确示例：User 聚合**
```go
// User 聚合根包含所有强相关的业务概念
type User struct {
    // 身份标识
    ID UserID
    
    // 核心属性（必须一起变化）
    Username UserName      // 用户名
    Email    Email         // 邮箱
    Password HashedPassword // 密码
    
    // 行为方法（维护一致性）
    func (u *User) ChangePassword(oldPwd, newPwd string) error {
        // 1. 验证旧密码
        // 2. 验证新密码强度
        // 3. 更新密码
        // 4. 发布 UserPasswordChangedEvent
    }
}
```

**❌ 错误示例：过度设计**
```go
// 不要把不相关的对象放在一个聚合里
type User struct {
    // ... user fields
    
    // ❌ 订单不应该在 User 聚合中
    orders []*Order
    
    // ❌ 日志不应该在 User 聚合中
    logs []*AuditLog
}
```

---

### 2. 值对象（Value Objects）

#### 什么是值对象？

值对象是没有身份标识的对象，它们通过属性值来定义相等性。

**示例：UserName**
```go
type UserName struct {
    value string
}

func NewUserName(value string) (*UserName, error) {
    if len(value) < 3 || len(value) > 50 {
        return nil, errors.New("用户名长度必须在 3-50 之间")
    }
    return &UserName{value: value}, nil
}

func (n *UserName) Value() string {
    return n.value
}
```

#### 值对象的优势

1. **自动验证** - 创建时即保证合法性
2. **不可变性** - 线程安全
3. **明确语义** - `UserName` 比 `string` 更清晰
4. **封装逻辑** - 可包含相关行为（如格式化）

---

### 3. 领域事件（Domain Events）

#### 什么是领域事件？

领域事件记录领域中发生的重要事情，用于触发副作用或与其他系统集成。

**示例：UserRegisteredEvent**
```go
type UserRegisteredEvent struct {
    *ddd.BaseEvent
    UserID       UserID    `json:"user_id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    RegisteredAt time.Time `json:"registered_at"`
}

func NewUserRegisteredEvent(userID UserID, username, email string) *UserRegisteredEvent {
    return &UserRegisteredEvent{
        BaseEvent:    ddd.NewBaseEvent("UserRegistered", userID, 1),
        UserID:       userID,
        Username:     username,
        Email:        email,
        RegisteredAt: time.Now(),
    }
}
```

#### 领域事件的用途

**✅ 正确使用场景：**
1. 发送通知邮件
2. 初始化统计数据
3. 记录审计日志
4. 触发工作流

**❌ 不使用场景：**
1. 同步更新读模型（这是 CQRS 的做法）
2. 强制业务验证（应该用领域服务）

---

### 4. 领域服务（Domain Services）

#### 什么是领域服务？

当某些业务逻辑不属于单个实体或值对象时，使用领域服务。

**示例：PasswordHasher**
```go
type PasswordHasher interface {
    Hash(password string) (string, error)
    Verify(password, hash string) bool
}

type BcryptPasswordHasher struct {
    cost int
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}
```

---

## Clean Architecture 分层实现

### 1. Domain Layer（领域层）⭐

**职责：** 包含核心业务逻辑，不依赖任何外部库

**目录结构：**
```
internal/domain/user/
├── model/
│   ├── user.go              # 聚合根
│   ├── valueobjects.go      # 值对象
│   └── builder.go           # Builder
├── event/
│   └── events.go            # 领域事件
├── service/
│   └── password_hasher.go   # 领域服务
└── repository/
    └── user_repository.go   # 仓储接口
```

**依赖规则：** 只能依赖标准库和 shared/ddd 包

---

### 2. Application Layer（应用层）

**职责：** 协调领域对象完成具体用例

**示例：UserService**
```go
type UserService interface {
    RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error)
    AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticationResult, error)
    GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error)
}

type UserServiceImpl struct {
    userRepo       user.UserRepository
    eventPublisher ddd.EventPublisher
    passwordHasher user.PasswordHasher
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error) {
    // 1. 检查唯一性
    // 2. 哈希密码
    // 3. 创建聚合根
    // 4. 保存并发布事件
}
```

**依赖规则：** 只能依赖 Domain 层和标准库

---

### 3. Infrastructure Layer（基础设施层）

**职责：** 实现技术细节（数据库、缓存、消息队列等）

**示例：UserRepositoryImpl**
```go
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u *model.User) error {
    // GORM 实现
    return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    var userDB dao.User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&userDB).Error
    // ...
}
```

**依赖规则：** 可以依赖任何层（但实际只依赖 Domain 的接口）

---

### 4. Interfaces Layer（接口层）

**职责：** 处理外部请求（HTTP、gRPC、CLI 等）

**示例：HTTP Handler**
```go
type Handler struct {
    userService userApp.UserService
}

func (h *Handler) GetUser(c *gin.Context) {
    userID, _ := h.mapper.ParseUserID(c.Param("id"))
    result, _ := h.userService.GetUserByID(c.Request.Context(), userID)
    h.respHandler.Success(c, toUserDetailDTO(result))
}
```

**依赖规则：** 可以依赖 Application 层

---

## 依赖注入（Bootstrap 层）

### Composition Root 模式

所有依赖关系在 Bootstrap 层组装：

```go
func (b *Bootstrap) initUserDomain(ctx context.Context) error {
    // 1. 基础设施
    eventPublisher := NewInMemoryEventPublisher(logger)
    userRepo := repositoryPkg.NewUserRepository(db)
    passwordHasher := userDomain.NewBcryptPasswordHasher(12)
    
    // 2. 应用服务
    b.user.service = userApp.NewUserService(userRepo, eventPublisher, passwordHasher)
    
    // 3. 事件处理器
    b.user.eventHandler = userApp.NewUserEventHandler(logger)
    
    return nil
}
```

---

## 核心业务流程示例

### 1. 用户注册流程

```
HTTP POST /auth/register
    ↓
AuthHandler.Register()
    ↓
UserService.RegisterUser(cmd)
    ├─ 1. 检查用户名唯一性
    ├─ 2. 检查邮箱唯一性
    ├─ 3. PasswordHasher.Hash(password)
    ├─ 4. User.NewUser(...) → 创建聚合根
    ├─ 5. UserRepository.Save(user) → 发布 UserRegisteredEvent
    └─ 6. JWTService.GenerateTokenPair() → 返回令牌
```

### 2. 用户登录流程

```
HTTP POST /auth/login
    ↓
AuthHandler.Login()
    ↓
UserService.AuthenticateUser(cmd)
    ├─ 1. UserRepository.FindByUsername(username)
    ├─ 2. PasswordHasher.Verify(password, hash)
    ├─ 3. User.RecordLogin() → 记录登录
    ├─ 4. UserRepository.Save(user) → 发布 UserLoggedInEvent
    └─ 5. JWTService.GenerateTokenPair() → 返回令牌
```

---

## 与 CQRS 的区别

### 我们的选择：纯 DDD

| 特性 | CQRS | 我们的实现 |
|------|------|-----------|
| 读写模型 | 分离 | ✅ 统一模型 |
| Projector | 需要 | ❌ 不需要 |
| Read Model 表 | 需要 | ❌ 不需要 |
| Command Handlers | 多个独立 | ✅ 统一 Service |
| Query Handlers | 多个独立 | ✅ 直接查询 |
| 事件用途 | 更新读模型 | ✅ 触发副作用 |
| 复杂度 | 高 | ✅ 适中 |

### 为什么选择纯 DDD？

1. **项目规模适中** - 当前阶段不需要读写分离
2. **简化开发** - 减少样板代码
3. **性能足够** - 单表查询性能良好
4. **易于维护** - 代码量少，逻辑清晰

---

## 最佳实践建议

### ✅ DO（推荐做法）

1. **领域层要纯粹** - 不包含任何框架代码
2. **使用值对象** - 封装业务概念
3. **聚合根要小** - 只包含强相关对象
4. **领域事件命名** - 使用过去式（UserRegistered）
5. **Repository 接口在 Domain 层** - 实现在 Infrastructure 层

### ❌ DON'T（避免做法）

1. **不要在领域层导入 GORM** - 使用 Repository 接口
2. **不要贫血模型** - 聚合根要有行为方法
3. **不要滥用领域事件** - 仅用于触发副作用
4. **不要提前优化** - 等业务需要再引入 CQRS

---

## 总结

本项目采用**纯 DDD + Clean Architecture**架构，专注于：

1. **领域建模** - 聚合根、值对象、领域事件
2. **分层清晰** - Domain → Application → Infrastructure → Interfaces
3. **依赖倒置** - 高层模块不依赖低层模块
4. **简单实用** - 不过度设计，按需演进

这种架构在保证代码质量的同时，避免了 CQRS 的复杂性，适合当前项目阶段！🎉
