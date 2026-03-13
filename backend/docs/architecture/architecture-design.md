# Go DDD Scaffold 技术架构文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的技术架构设计，包括整体架构模式、分层设计、技术选型理由以及各组件间的交互关系。

## 整体架构设计

### 架构模式选择

项目采用 **Clean Architecture + DDD + Composition Root** 的混合架构模式，结合了三种架构思想的优势：

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
              Bootstrap (Composition Root)
              - 创建所有依赖
              - 组装依赖关系
              - 初始化应用服务
```

### Composition Root 设计

#### 核心理念

**Composition Root（组合根）** 是应用的启动入口，负责创建和组装所有依赖。在 `go-ddd-scaffold` 中，由 `Bootstrap` 类承担此职责。

```
┌──────────────────────────────────────┐
│       Bootstrap                      │
│    (Composition Root)                │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ 领域组件（按领域分组）      │ │
│  │   - user.*Handler              │ │
│  │   - tenant.*Handler            │ │
│  │   - order.*Handler             │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ 基础设施组件               │ │
│  │   - DB (PostgreSQL)            │ │
│  │   - Redis                      │ │
│  │   - Logger                     │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │ ✅ 接口层组件                 │ │
│  │   - HTTP Handlers              │ │
│  │   - Router                     │ │
│  └────────────────────────────────┘ │
└──────────────────────────────────────┘
```

#### 职责分离

| 组件 | 职责 | 管理内容 |
|------|------|----------|
| **Bootstrap** | Composition Root | 创建领域组件、组装依赖关系 |
| **Container** | 基础设施 + 路由 | DB/Redis/Logger、HTTP 路由、生命周期 |
| **Domain** | 业务逻辑 | 实体、值对象、领域服务 |
| **HTTP** | 接口层 | 路由、Handler、响应处理 |

**详细说明请参考**: [Container 重构说明](./container-refactoring.md)

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

#### 主要组件

**实体 (Entities)**
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

### 2. 应用层 (Application Layer)

#### 核心职责
- 协调领域对象完成业务用例
- 处理跨聚合的业务逻辑
- 定义应用服务接口

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

## CQRS架构模式详解

### 核心设计理念

CQRS（Command Query Responsibility Segregation）将系统的读写操作完全分离，这种模式特别适合复杂的业务场景：

```
┌─────────────────┐    ┌─────────────────┐
│   Command Side  │    │   Query Side    │
│  (Write Model)  │    │  (Read Model)   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│  Command Bus    │    │  Query Service  │
│  & Handlers     │    │  & Projections  │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│  Domain Model   │◄──►│  Read Stores    │
│  (Aggregates)   │    │  (Optimized)    │
└─────────────────┘    └─────────────────┘
```

### 读写模型分离策略

#### 写模型（Command Model）
- **职责**：处理业务逻辑、维护数据一致性
- **特点**：规范化设计、强一致性、复杂业务规则
- **存储**：事务性数据库（PostgreSQL）

#### 读模型（Query Model）
- **职责**：优化查询性能、支持复杂展示需求
- **特点**：非规范化设计、最终一致性、高性能读取
- **存储**：可选用不同存储（Redis、Elasticsearch、专用读库）

## DDD战术模式应用

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

**CQRS视角下的聚合设计**：
```go
// 命令侧聚合根 - 关注业务逻辑和一致性
type UserAggregate struct {
    baseAggregate BaseAggregate
    id           UserID
    username     string
    email        Email
    password     HashedPassword
    status       UserStatus
    profile      UserProfile
    settings     UserSettings
    roles        []UserRole
}

// 查询侧读模型 - 关注查询性能和展示需求
type UserReadModel struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    Status       int       `json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    TenantCount  int       `json:"tenant_count"`
    LastLoginAt  time.Time `json:"last_login_at,omitempty"`
    DisplayName  string    `json:"display_name"`
}
```

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

// CQRS事件投影示例
type UserReadModelProjector struct {
    db *gorm.DB
}

func (p *UserReadModelProjector) HandleUserCreated(event UserCreatedEvent) error {
    readModel := UserReadModel{
        ID:        int64(event.UserID),
        Username:  event.Username,
        Email:     event.Email,
        Status:    1, // Active
        CreatedAt: event.CreatedAt,
        UpdatedAt: event.CreatedAt,
    }
    
    return p.db.Table("user_read_model").Create(&readModel).Error
}

func (p *UserReadModelProjector) HandleUserEmailChanged(event UserEmailChangedEvent) error {
    return p.db.Table("user_read_model").
        Where("id = ?", int64(event.UserID)).
        Updates(map[string]interface{}{
            "email":      event.NewEmail,
            "updated_at": event.ChangedAt,
        }).Error
}
```

## CQRS完整实现方案

### 命令侧实现

**命令对象设计**：
```go
// 命令接口定义
type Command interface {
    CommandName() string
    Validate() error
}

// 具体命令实现
type CreateUserCommand struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    TenantID *int64 `json:"tenant_id,omitempty"`
}

func (cmd CreateUserCommand) CommandName() string {
    return "CreateUser"
}

func (cmd CreateUserCommand) Validate() error {
    // 基础验证逻辑
    return nil
}
```

**命令处理器实现**：
```go
type CommandHandler interface {
    Handle(command Command) (interface{}, error)
}

type UserCommandHandler struct {
    userRepo    UserRepository
    tenantRepo  TenantRepository
    eventBus    EventBus
    passwordSvc PasswordService
}

func (h *UserCommandHandler) HandleCreateUser(cmd CreateUserCommand) (UserID, error) {
    // 1. 命令验证
    if err := cmd.Validate(); err != nil {
        return 0, err
    }
    
    // 2. 业务规则验证
    if exists, _ := h.userRepo.ExistsByEmail(cmd.Email); exists {
        return 0, errors.New("email already registered")
    }
    
    // 3. 创建聚合根
    user, err := NewUser(cmd.Username, cmd.Email, cmd.Password)
    if err != nil {
        return 0, err
    }
    
    // 4. 处理租户关联
    if cmd.TenantID != nil {
        tenant, err := h.tenantRepo.GetByID(TenantID(*cmd.TenantID))
        if err != nil {
            return 0, err
        }
        user.AssignToTenant(tenant.ID())
    }
    
    // 5. 持久化聚合
    if err := h.userRepo.Save(user); err != nil {
        return 0, err
    }
    
    // 6. 发布领域事件
    events := user.GetUncommittedEvents()
    for _, event := range events {
        h.eventBus.Publish(event)
    }
    user.ClearUncommittedEvents()
    
    return user.ID(), nil
}
```

### 查询侧实现

**查询对象设计**：
```go
// 查询接口定义
type Query interface {
    QueryName() string
}

type GetUserProfileQuery struct {
    UserID UserID `json:"user_id"`
}

func (q GetUserProfileQuery) QueryName() string {
    return "GetUserProfile"
}

type ListUsersQuery struct {
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    Status   *int   `json:"status,omitempty"`
    Keyword  string `json:"keyword,omitempty"`
}

func (q ListUsersQuery) QueryName() string {
    return "ListUsers"
}
```

**查询服务实现**：
```go
type QueryService interface {
    Execute(query Query) (interface{}, error)
}

type UserQueryService struct {
    db *gorm.DB
}

// 优化的读模型查询
func (qs *UserQueryService) GetUserProfile(query GetUserProfileQuery) (*UserProfileDTO, error) {
    var profile UserProfileDTO
    
    err := qs.db.Table("user_read_model").
        Select(`
            id,
            username,
            email,
            status,
            created_at,
            updated_at,
            tenant_count,
            last_login_at
        `).
        Where("id = ?", query.UserID).
        First(&profile).Error
    
    return &profile, err
}

// 复杂列表查询
func (qs *UserQueryService) ListUsers(query ListUsersQuery) (*PagedResult[UserListItemDTO], error) {
    // 分页和筛选逻辑
    // ...
    return result, nil
}
```

### 读模型同步策略

**事件驱动的读模型更新**：
```go
// 读模型投影器
type ReadModelProjector interface {
    Project(event DomainEvent) error
}

type UserReadModelProjector struct {
    db *gorm.DB
}

func (p *UserReadModelProjector) Project(event DomainEvent) error {
    switch evt := event.(type) {
    case UserCreatedEvent:
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