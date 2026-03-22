# Go DDD Scaffold 技术架构文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的技术架构设计，包括整体架构模式、分层设计、技术选型理由以及各组件间的交互关系。

## 整体架构设计

### 架构模式选择

项目采用 **Clean Architecture + DDD + 构造函数注入** 的混合架构模式，结合了三种架构思想的优势：

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│  (HTTP Controllers, CLI Commands, gRPC Services)        │
├─────────────────────────────────────────────────────────┤
│                    Application Layer                     │
│  (Use Cases, Application Services, DTOs)                │
├─────────────────────────────────────────────────────────┤
│                      Domain Layer                        │
│  (Entities, Value Objects, Aggregates, Domain Services) │
├─────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                     │
│  (Repositories, External Services, Persistence)         │
└─────────────────────────────────────────────────────────┘
                         ↑
              main.go (Composition Root)
              - 创建 Infra 结构体
              - 实例化 Module
              - 遍历注册 HTTP/Event
```

### Composition Root 设计

#### 核心理念

**Composition Root（组合根）** 是应用的启动入口，负责创建和组装所有依赖。在 `go-ddd-scaffold` 中，由 `cmd/api/main.go` 承担此职责，采用 **Infra + Module** 模式。

```
┌──────────────────────────────────────┐
│       main.go                        │
│    (Composition Root)                │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ Infra 结构体（基础设施）    │ │
│  │   - DB (PostgreSQL/GORM)       │ │
│  │   - Redis                      │ │
│  │   - Logger                     │ │
│  │   - Config                     │ │
│  │   - Snowflake                  │ │
│  │   - EventPublisher             │ │
│  │   - EventBus                   │ │
│  │   - AsynqClient                │ │
│  │   - ErrorMapper                │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ Module（模块组装层）        │ │
│  │   - UserModule                 │ │
│  │   - AuthModule                 │ │
│  │   - 内部构建完整依赖链         │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ 遍历注册                    │ │
│  │   - RegisterHTTP(api)          │ │
│  │   - RegisterSubscriptions(bus) │ │
│  └────────────────────────────────┘ │
└──────────────────────────────────────┘
```

#### 职责分离

| 组件 | 职责 | 管理内容 |
|------|------|----------|
| **main.go** | Composition Root | 加载配置、创建 Infra、实例化 Module、注册路由 |
| **Infra** | 基础设施容器 | DB/Redis/Logger/Config/Snowflake/EventPublisher/EventBus/AsynqClient/ErrorMapper |
| **Module** | 模块组装层（胶水层） | 构建模块内完整依赖链，实现 HTTPModule/EventModule 接口 |
| **Domain** | 业务逻辑 | 实体、值对象、领域服务 |
| **HTTP** | 接口层 | Handler、Routes、响应处理 |

### 核心设计原则

1. **依赖倒置原则**：高层模块不依赖低层模块，都依赖抽象
2. **单一职责原则**：每层都有明确的职责边界
3. **开闭原则**：对扩展开放，对修改封闭
4. **里氏替换原则**：子类型必须能够替换它们的基类型

## 分层详细设计

### 1. 领域层 (Domain Layer)

#### 核心职责
- 业务逻辑的核心表达
- 领域概念的精确建模
- 业务规则的封装实现

#### 目录结构
```
internal/domain/
├── user/                      # 用户领域
│   ├── aggregate/             # 聚合根
│   │   ├── user.go            # User聚合根
│   │   └── builder.go         # Builder 模式
│   ├── valueobject/           # 值对象
│   │   └── value_objects.go   # 值对象定义
│   ├── event/                 # 领域事件
│   │   └── events.go          # 事件定义
│   ├── service/               # 领域服务
│   │   ├── password_hasher.go # 密码哈希服务
│   │   └── password_policy.go # 密码策略服务
│   ├── repository/            # 仓储接口
│   │   └── user_repository.go # 用户仓储接口
│   └── types.go               # 类型别名统一导出
├── tenant/                    # 租户领域
│   └── ...                    # 类似的结构
└── shared/kernel/             # 领域共享内核
    ├── entity.go              # 聚合根基类
    ├── value_object.go        # 值对象接口
    ├── event.go               # 领域事件基类
    ├── event_bus.go           # 事件总线接口
    ├── repository.go          # 仓储接口
    ├── errors.go              # 领域错误定义
    └── response.go            # ErrorMapper（错误映射器）
```

#### 主要组件

**聚合根 (Aggregates)**
```go
// User 用户聚合根 - 具有唯一标识的领域对象
type User struct {
    ddd.BaseEntity
    
    username       *UserName
    email          *Email
    password       *HashedPassword
    status         UserStatus
    displayName    string
    firstName      string
    lastName       string
    gender         UserGender
    phoneNumber    string
    avatarURL      string
    lastLoginAt    *time.Time
    loginCount     int
    lockedUntil    *time.Time
    failedAttempts int
    createdAt      time.Time
    updatedAt      time.Time
}

func (u *User) ChangePassword(oldPass, newPass string, ipAddress string) error {
    if !u.password.Matches(oldPass) {
        return ddd.NewBusinessError("INVALID_OLD_PASSWORD", "invalid old password")
    }
    u.password = NewHashedPassword(newPass)
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserPasswordChangedEvent(u.ID().(UserID), ipAddress)
    u.ApplyEvent(event)
    
    return nil
}
```

**值对象 (Value Objects)**
```go
// UserID 用户标识 - 不可变，基于值相等
type UserID struct {
    ddd.Int64Identity
}

func NewUserID(value int64) UserID {
    return UserID{Int64Identity: ddd.NewInt64Identity(value)}
}

func (uid UserID) String() string {
    return fmt.Sprintf("user-%d", uid.Int64())
}

// Email 邮箱值对象
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    email := &Email{value: strings.TrimSpace(strings.ToLower(value))}
    if err := email.Validate(); err != nil {
        return nil, err
    }
    return email, nil
}

func (e *Email) Value() string { return e.value }
func (e *Email) Validate() error { /* 验证逻辑 */ }
```

**聚合根 (Aggregate Roots)**
```go
// User 用户聚合根 - 管理用户相关的业务一致性
type User struct {
    ddd.BaseEntity  // 包含 ID、Version、UncommittedEvents
    
    // 业务属性
    username *UserName
    email    *Email
    password *HashedPassword
    status   UserStatus
    // ... 其他字段
}

// Activate 激活用户 - 业务方法示例
func (u *User) Activate() error {
    if u.status != UserStatusPending {
        return ddd.NewBusinessError("USER_NOT_PENDING", "user is not in pending status")
    }
    u.status = UserStatusActive
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserActivatedEvent(u.ID().(UserID))
    u.ApplyEvent(event)
    
    return nil
}

// GetUncommittedEvents 获取未提交事件（由仓储使用）
func (u *User) GetUncommittedEvents() []ddd.DomainEvent {
    return u.BaseEntity.GetUncommittedEvents()
}
```

**领域服务 (Domain Services)**
```go
// AuthenticationService 认证服务
type AuthenticationService struct {
    userRepo       UserRepository
    tokenService   TokenService
    passwordPolicy PasswordPolicyService
}

func (s *AuthenticationService) Authenticate(ctx context.Context, usernameOrEmail, password string, ipAddress, userAgent string) (*AuthenticateResult, error) {
    // 1. 查找用户
    u, err := s.findUser(ctx, usernameOrEmail)
    if err != nil {
        return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
    }

    // 2. 验证密码
    if !u.Password().Matches(password) {
        u.RecordFailedLogin(ipAddress, userAgent, "invalid_password")
        _ = s.userRepo.Save(ctx, u)
        return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
    }

    // 3. 检查账户状态
    if !u.CanLogin() {
        return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "account cannot login")
    }

    // 4. 记录成功登录并发布事件
    u.RecordLogin(ipAddress, userAgent)
    if err := s.userRepo.Save(ctx, u); err != nil {
        return nil, err
    }

    // 5. 生成令牌
    tokenPair, err := s.tokenService.GenerateTokenPair(u.ID().(UserID))
    if err != nil {
        return nil, err
    }

    return &AuthenticateResult{
        UserID:       u.ID().(UserID),
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
    }, nil
}
```

#### 2. 应用层 (Application Layer)

#### 核心职责
- 协调领域对象完成应用用例
- 定义和实现应用服务
- 处理跨领域的业务流程

#### 目录结构
```
internal/application/
├── user/
│   ├── service.go         # 用户应用服务
│   ├── dtos.go            # DTOs（Commands + Results）
│   └── event_handlers.go  # 领域事件处理器
└── auth/
    ├── service.go         # 认证应用服务
    └── dtos.go            # 认证 DTOs
```

---

#### 3. 基础设施层 (Infrastructure Layer)

#### 核心职责
- 实现技术细节（数据库、缓存、消息队列等）
- 提供持久化机制
- 实现外部服务集成

#### 目录结构
```
internal/infrastructure/
├── persistence/           # 数据持久化
│   ├── dao/               # DAO 层
│   ├── model/             # 数据库模型
│   └── repository/        # 仓储实现
├── cache/                 # 缓存实现
├── messaging/             # 消息队列
├── logging/               # 日志实现
└── config/                # 配置加载
```

---

#### 4. 接口层 (Interfaces Layer)

#### 核心职责
- 处理外部请求和响应
- 协议转换（HTTP/gRPC/CLI）
- 请求验证和响应格式化

#### 目录结构
```
internal/interfaces/
├── http/
│   ├── middleware/        # HTTP 中间件
│   ├── response.go        # HTTP 响应处理器
│   ├── router.go          # 路由配置
│   ├── user/              # 用户领域 HTTP Handler
│   │   ├── handler.go
│   │   ├── request.go
│   │   └── response.go
│   └── auth/              # 认证领域 HTTP Handler
│       └── handler.go
└── grpc/                  # gRPC 接口（预留）
```

---

#### 5. 公共包 (pkg/)

#### 核心职责
- 提供可跨层复用的公共组件
- 不包含业务逻辑
- 可被任何层引用

#### 目录结构
```
pkg/
├── response/              # HTTP 响应公共包
│   └── response.go        # Response, ErrorResponse, PageData
└── util/                  # 通用工具包
    ├── cast.go            # 类型转换
    └── time.go            # 时间处理
```

#### 使用示例
```go
import (
    "github.com/shenfay/go-ddd-scaffold/pkg/response"
    "github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// HTTP 响应
resp := response.NewResponse(data)
c.JSON(http.StatusOK, resp)

// 类型转换
userID := util.ToInt64(c.Query("user_id"))

#### 目录结构
```
internal/application/
├── user/
│   ├── service.go           # UserService 接口和实现
│   ├── dtos.go              # 所有 DTOs（Commands + Results）
│   └── event_handlers.go    # 领域事件处理器
├── auth/
│   ├── service.go           # AuthService 接口和实现
│   └── dtos.go              # Auth 相关 DTOs
└── shared/dto/
    └── page.go              # 分页 DTO
```

#### 核心职责
- 协调领域对象完成业务用例
- 处理跨聚合的业务逻辑
- 定义应用服务接口
- **统一管理 DTOs**（Input Commands + Output Results）

### Application Service 设计

每个领域一个统一的 Service，负责协调领域对象完成用例。

```go
// user/service.go
type UserService interface {
    RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error)
    AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error)
    GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error)
    UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
    ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error
}
```

**DTO 组织方式**

所有 DTOs 合并在 `dtos.go` 文件中，按功能分组排列。

#### 主要用例实现
```go
// 用户管理用例
type UserUseCase struct {
    userRepo     UserRepository
    tenantRepo   TenantRepository
    eventBus     EventBus
    passwordSvc  PasswordPolicyService
}

func (uc *UserUseCase) CreateUser(req CreateUserRequest) (*User, error) {
    // 1. 验证业务规则
    if err := uc.passwordSvc.Validate(req.Password); err != nil {
        return nil, err
    }
    
    // 2. 创建领域对象
    user := NewUser(req.Username, req.Email, req.Password)
    
    // 3. 持久化
    if err := uc.userRepo.Save(user); err != nil {
        return nil, err
    }
    
    // 4. 发布领域事件
    uc.eventBus.Publish(UserCreatedEvent{UserID: user.ID()})
    
    return user, nil
}
```

### 3. 接口层 (Interfaces Layer)

#### 核心职责
- 处理外部请求和响应
- 协议转换（HTTP/gRPC/CLI）
- 输入验证和错误处理

#### HTTP控制器示例
```go
type UserController struct {
    userUseCase UserUseCase
    logger      *zap.Logger
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ctrl.handleValidationError(c, err)
        return
    }
    
    user, err := ctrl.userUseCase.CreateUser(req)
    if err != nil {
        ctrl.handleBusinessError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, UserResponseFromDomain(user))
}
```

### 4. 基础设施层 (Infrastructure Layer)

#### 核心职责
- 实现领域层定义的接口
- 处理技术细节（数据库、缓存、外部服务）
- 提供基础设施服务

#### 仓储实现示例
```go
// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
    db DB
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u *User) error {
    model := r.toModel(u)

    return r.db.WithContext(ctx).Transaction(func(tx Tx) error {
        // 检查用户是否存在
        var existingVersion int
        err := tx.QueryRow(ctx, "SELECT version FROM users WHERE id = ?", model.ID).Scan(&existingVersion)

        if err != nil { // 不存在，创建新用户
            _, err = tx.Exec(ctx,
                `INSERT INTO users (...) VALUES (...)`,
                model.ID, model.Username, model.Email, ...,
            )
        } else {
            // 乐观锁检查
            if existingVersion != u.Version()-1 {
                return ddd.NewConcurrencyError(u.ID(), u.Version()-1, existingVersion, "version conflict")
            }
            // 更新用户
            _, err = tx.Exec(ctx, `UPDATE users SET ... WHERE id = ?`, ..., model.ID)
        }
        
        // 保存未提交的领域事件到事件存储
        events := u.GetUncommittedEvents()
        if len(events) > 0 {
            if err := r.eventStore.AppendEvents(ctx, u.ID(), events); err != nil {
                return err
            }
            u.ClearUncommittedEvents()
        }
        
        return err
    })
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id UserID) (*User, error) {
    var model UserModel
    err := r.db.QueryRow(ctx, `SELECT ... FROM users WHERE id = ?`, id.Int64()).Scan(...)
    if err != nil {
        return nil, ddd.ErrAggregateNotFound
    }
    return r.toDomain(&model), nil
}
```

## 技术选型详解

### 核心技术栈

#### 后端技术选型

**Go 1.25.6**
- **选择理由**：高性能、静态类型、并发支持好、生态成熟
- **优势**：编译速度快、部署简单、内存占用低
- **适用场景**：高并发Web服务、微服务、API网关

**标准库 database/sql + 驱动**
- **选择理由**：轻量级、无额外依赖、完全控制SQL执行
- **关键特性**：连接池管理、预处理语句、事务支持
- **优势**：避免ORM复杂性，更好的性能控制

**PostgreSQL**
- **选择理由**：ACID特性、JSON支持、扩展性好
- **优势**：数据完整性保障、复杂查询支持、地理信息系统
- **适用场景**：OLTP系统、数据分析、复杂业务逻辑

**Redis**
- **选择理由**：高性能、丰富的数据结构、持久化支持
- **使用场景**：会话存储、缓存、分布式锁、消息队列
- **集群支持**：主从复制、哨兵模式、集群模式

#### DDD 基础设施

**领域事件存储**
- **实现方式**：PostgreSQL 表存储（domain_events）
- **特点**：支持事件回放、审计追踪、异步处理
- **优势**：与业务数据同一事务，保证一致性

**事件总线**
- **实现方式**：内存事件总线（支持同步/异步）
- **特点**：发布订阅模式、多处理器支持
- **扩展**：可替换为消息队列（RabbitMQ/Kafka）

**CQRS 读模型**
- **实现方式**：独立读模型表 + 投影器
- **特点**：最终一致性、查询优化、读写分离
- **优势**：读性能优化、灵活的数据结构

#### 主键生成方案

**Snowflake ID**
```
时间戳(41位) + 机器ID(10位) + 序列号(12位)
```
- **优势**：全局唯一、趋势递增、分布式友好
- **性能**：纳秒级生成速度，无网络依赖
- **存储**：64位整数，索引效率高

#### 认证授权方案

**JWT双Token机制**
```
Access Token (短期) + Refresh Token (长期)
```
- **安全性**：签名验证、过期控制、撤销机制
- **性能**：无状态认证，服务端无需存储会话
- **扩展性**：支持负载均衡和水平扩展

### 前端技术选型（预留）

**React 18**
- **选择理由**：组件化、虚拟DOM、生态丰富
- **核心特性**：Hooks、Suspense、并发渲染
- **开发体验**：热重载、TypeScript支持、调试工具

**Tailwind CSS**
- **选择理由**：实用优先、原子化CSS、开发效率高
- **优势**：减少CSS文件体积、响应式设计简单
- **定制性**：可通过配置文件深度定制

## 领域事件与最终一致性

### 领域事件的用途

领域事件主要用于触发副作用（side effects），如发送邮件、记录审计日志、更新统计信息等。

### 事件驱动的最终一致性

```
┌──────────────────┐
│  Application     │
│   Service        │
└────────┬─────────┘
         │ 1. 执行业务逻辑
         ▼
┌──────────────────┐
│  Domain Model    │
│  (Aggregate)     │
└────────┬─────────┘
         │ 2. 发布领域事件
         ▼
┌──────────────────┐
│  Event Publisher │
└────────┬─────────┘
         │ 3. 分发给订阅者
         ├──► Email Handler (发送邮件)
         ├──► Audit Handler (记录日志)
         └──► Stats Handler (更新统计)
```

## DDD 战术模式应用

### 聚合设计原则
```
用户聚合边界
├── User (聚合根)
├── UserProfile (实体)
├── UserSettings (值对象)
└── UserTenant (值对象)
```

**设计要点**：
- 聚合根负责维护一致性边界
- 跨聚合通过领域事件通信
- 聚合尽量保持小巧专注

### 领域事件模式
```go
// 领域事件接口定义
type DomainEvent interface {
    EventName() string
    OccurredOn() time.Time
    AggregateID() interface{}
    Version() int
}

// 用户领域事件体系
type (
    // 用户生命周期事件
    UserCreatedEvent struct {
        UserID    UserID    `json:"user_id"`
        Username  string    `json:"username"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
    }
    
    UserActivatedEvent struct {
        UserID     UserID    `json:"user_id"`
        ActivatedAt time.Time `json:"activated_at"`
    }
    
    UserDeactivatedEvent struct {
        UserID       UserID    `json:"user_id"`
        Reason       string    `json:"reason"`
        DeactivatedAt time.Time `json:"deactivated_at"`
    }
    
    // 用户属性变更事件
    UserEmailChangedEvent struct {
        UserID    UserID `json:"user_id"`
        OldEmail  string `json:"old_email"`
        NewEmail  string `json:"new_email"`
        ChangedAt time.Time `json:"changed_at"`
    }
)

func (e UserCreatedEvent) EventName() string { return "UserCreated" }
func (e UserCreatedEvent) OccurredOn() time.Time { return e.CreatedAt }
func (e UserCreatedEvent) AggregateID() interface{} { return e.UserID }
func (e UserCreatedEvent) Version() int { return 1 }

// 事件总线接口
type EventBus interface {
    Publish(event DomainEvent) error
    Subscribe(eventType string, handler EventHandler) error
}

// 事件处理器类型
type EventHandler func(event DomainEvent) error

// 异步事件处理示例
func (h *EmailNotificationHandler) Handle(event DomainEvent) error {
    switch evt := event.(type) {
    case UserCreatedEvent:
        return h.sendWelcomeEmail(evt.Email, evt.Username)
    case UserEmailChangedEvent:
        return h.sendEmailChangeNotification(evt.OldEmail, evt.NewEmail)
    }
    return nil
}
```

## Application Service 实现方案

### Application Service 职责

Application Service 负责协调领域对象完成业务用例。它直接调用 Repository、领域服务和其他基础设施组件。

**注意**：本项目不使用 Command Handler 模式，Command 仅作为参数对象在 Application Service 方法中传递。

```go
// Application Service 实现示例
type UserServiceImpl struct {
    userRepo       UserRepository
    eventPublisher EventPublisher
    passwordHasher PasswordHasher
    tokenService   TokenService
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error) {
    // 1. 验证业务规则
    if err := s.passwordHasher.Validate(cmd.Password); err != nil {
        return nil, err
    }
    
    // 2. 检查用户名/邮箱是否已存在
    if exists, _ := s.userRepo.ExistsByUsername(cmd.Username); exists {
        return nil, ddd.NewBusinessError("USERNAME_EXISTS", "username already exists")
    }
    
    // 3. 创建领域对象
    hashedPassword := s.passwordHasher.Hash(cmd.Password)
    user, err := model.NewUser(cmd.Username, cmd.Email, hashedPassword, s.idGenerator)
    if err != nil {
        return nil, err
    }
    
    // 4. 持久化
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. 发布领域事件
    events := user.GetUncommittedEvents()
    for _, event := range events {
        if err := s.eventPublisher.Publish(ctx, event); err != nil {
            // 记录错误但不中断主流程
        }
    }
    
    // 6. 返回结果 DTO
    return &RegisterUserResult{
        UserID:   user.ID().(model.UserID).Int64(),
        Username: user.Username().Value(),
        Email:    user.Email().Value(),
    }, nil
}
```

### 查询处理

查询直接从聚合根读取数据，保证强一致性。对于复杂查询，可以通过 Repository 的自定义方法实现。

```go
// 查询直接通过 Repository 获取聚合根
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error) {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 转换为结果 DTO
    return &GetUserResult{
        ID:          user.ID().(model.UserID).Int64(),
        Username:    user.Username().Value(),
        Email:       user.Email().Value(),
        DisplayName: user.DisplayName(),
        Status:      int(user.Status()),
        CreatedAt:   user.CreatedAt(),
    }, nil
}

// 复杂查询通过 Repository 的自定义方法实现
func (s *UserServiceImpl) ListUsersByStatus(ctx context.Context, status user.UserStatus) ([]*GetUserResult, error) {
    users, err := s.userRepo.FindByStatus(ctx, status)
    if err != nil {
        return nil, err
    }
    
    results := make([]*GetUserResult, len(users))
    for i, user := range users {
        results[i] = toGetUserResult(user)
    }
    return results, nil
}
```

### 领域事件的发布与应用

领域事件用于触发副作用，如发送邮件、记录审计日志等。

```go
// main.go 中创建模块并注册事件订阅
func main() {
    // ... 加载配置、创建 Logger ...
    
    // 创建基础设施
    infra, cleanup, err := bootstrap.NewInfra(appConfig, logger)
    defer cleanup()
    
    // 创建模块（内部自行构建完整依赖链）
    userMod := module.NewUserModule(infra)
    authMod := module.NewAuthModule(infra)
    modules := []bootstrap.Module{authMod, userMod}
    
    // 注册事件订阅
    for _, m := range modules {
        if em, ok := m.(bootstrap.EventModule); ok {
            em.RegisterSubscriptions(infra.EventBus)
        }
    }
    
    // 注册 HTTP 路由
    api := router.Group("/api/v1")
    for _, m := range modules {
        if h, ok := m.(bootstrap.HTTPModule); ok {
            h.RegisterHTTP(api)
        }
    }
}
```
        return p.handleUserCreated(evt)
    case UserEmailChangedEvent:
        return p.handleUserEmailChanged(evt)
    case UserDeactivatedEvent:
        return p.handleUserDeactivated(evt)
    }
    return nil
}

func (p *UserReadModelProjector) handleUserCreated(event UserCreatedEvent) error {
    readModel := UserReadModel{
        ID:        int64(event.UserID),
        Username:  event.Username,
        Email:     event.Email,
        Status:    1,
        CreatedAt: event.CreatedAt,
        UpdatedAt: event.CreatedAt,
    }
    return p.db.Table("user_read_model").Create(&readModel).Error
}
```



## 部署架构

### 单体应用部署
```
┌─────────────────────────────────────────┐
│              Load Balancer               │
├─────────────────────────────────────────┤
│           Application Servers            │
│        [Go App Instance 1..N]           │
├─────────────┬───────────────────────────┤
│   Cache     │        Database           │
│  (Redis)    │     (PostgreSQL)          │
└─────────────┴───────────────────────────┘
```

### 配置管理策略

**环境变量配置**：
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=scaffold_db
DB_USER=postgres
DB_PASSWORD=${SECRET_DB_PASSWORD}

# Redis配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=${SECRET_REDIS_PASSWORD}

# JWT配置
JWT_SECRET=${SECRET_JWT_KEY}
JWT_ACCESS_EXPIRE=30m
JWT_REFRESH_EXPIRE=7d
```

**YAML配置文件**：
```yaml
server:
  port: ${SERVER_PORT:8080}
  mode: ${ENV_MODE:debug}

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  max_idle_conns: 10
  max_open_conns: 100
```

## 性能优化策略

### 数据库优化
- **连接池配置**：合理的最大连接数和空闲连接数
- **索引策略**：复合索引、部分索引、表达式索引
- **查询优化**：预加载、批处理、分页查询

### 缓存策略
- **多级缓存**：本地缓存 + Redis缓存
- **缓存穿透**：布隆过滤器、空值缓存
- **缓存雪崩**：随机过期时间、熔断机制

### 并发处理
- **Goroutine池**：限制并发数量，防止资源耗尽
- **限流降级**：令牌桶算法、漏桶算法
- **超时控制**：合理的超时设置和重试机制

这个技术架构文档为项目提供了完整的技术蓝图，确保团队成员对系统设计有一致的理解。

---

## 📊 图表索引

为了更直观地理解架构设计，参考以下图表：

**[architecture-diagrams-section.md](./architecture-diagrams-section.md)** - 架构图集
- Clean Architecture 分层架构图（mermaid graph）
- Composition Root 设计图
- 依赖关系可视化说明
