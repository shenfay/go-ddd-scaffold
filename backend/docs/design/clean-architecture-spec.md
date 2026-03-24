# Clean Architecture 规范

本文档详细定义了 Clean Architecture（整洁架构）在本项目中的实现规范。

## 📋 架构分层

### 四层结构

```
┌─────────────────────────────────────────┐
│   Interfaces Layer (接口层)             │  ← HTTP/gRPC/CLI
├─────────────────────────────────────────┤
│   Application Layer (应用层)            │  ← Use Cases + Ports
├─────────────────────────────────────────┤
│   Domain Layer (领域层)                 │  ← Business Logic
├─────────────────────────────────────────┤
│   Infrastructure Layer (基础设施层)     │  ← Adapters
└─────────────────────────────────────────┘
```

### 依赖规则

**核心原则：** 依赖指向内层，外层知道内层，内层不知道外层

```
允许依赖:
✓ Interfaces → Application
✓ Application → Domain  
✓ Infrastructure → Domain (通过适配器)
✓ Bootstrap → All

禁止依赖:
✗ Application → Infrastructure
✗ Domain → Application
✗ Infrastructure → Application
✗ Interfaces → Domain
```

---

## 🎯 各层详细规范

### 1. Domain Layer（领域层）

#### 职责

- 封装业务逻辑和状态变更规则
- 定义业务对象（聚合根、值对象）
- 定义领域服务
- 定义 Repository 接口（Domain 需要持久化）

#### 目录结构

```
domain/
├── shared/
│   ├── kernel/           # 核心抽象
│   │   ├── entity.go
│   │   ├── valueobject.go
│   │   ├── aggregate.go
│   │   └── errors.go
│   └── event/            # 共享事件
│
├── user/                 # 用户限界上下文
│   ├── aggregate/        # 聚合根
│   │   └── user.go
│   ├── valueobject/      # 值对象
│   │   ├── username.go
│   │   ├── email.go
│   │   └── password.go
│   ├── event/            # 领域事件
│   │   ├── user_registered.go
│   │   └── user_logged_in.go
│   ├── service/          # 领域服务
│   │   └── password_hasher.go
│   └── repository/       # 仓储接口
│       └── user_repository.go
│
└── tenant/               # 租户限界上下文
    └── ...
```

#### 实现规范

**聚合根（Aggregate）**

```go
// domain/user/aggregate/user.go
package aggregate

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// User 用户聚合根
type User struct {
    *kernel.Entity
    username              vo.Username
    email                 vo.Email
    password              vo.Password
    status                vo.UserStatus
    failedLoginAttempts   int
    lastLoginAt           *time.Time
    lastLoginIP           string
}

// NewUser 创建新用户（构造函数）
func NewUser(username, email, password string) (*User, error) {
    // 验证并创建值对象
    uName, err := vo.NewUsername(username)
    if err != nil {
        return nil, err
    }
    
    uEmail, err := vo.NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    uPassword, err := vo.NewPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Entity:   kernel.NewEntity(),
        username: uName,
        email:    uEmail,
        password: uPassword,
        status:   vo.UserStatusPending,
    }
    
    // 发布领域事件
    user.RecordEvent(&event.UserRegistered{
        UserID:  user.ID().Value(),
        Email:   user.email.String(),
        Created: time.Now(),
    })
    
    return user, nil
}

// Login 登录业务逻辑
func (u *User) Login(password string, ip string, userAgent string) error {
    // 验证状态
    if u.status != vo.UserStatusActive {
        return kernel.ErrUserInactive
    }
    
    // 验证密码
    if !u.password.Verify(password) {
        u.failedLoginAttempts++
        if u.failedLoginAttempts >= MaxLoginAttempts {
            u.Lock()
            return kernel.ErrUserLocked
        }
        return kernel.ErrInvalidCredentials
    }
    
    // 登录成功
    u.failedLoginAttempts = 0
    u.lastLoginAt = ptr.Time(time.Now())
    u.lastLoginIP = ip
    
    u.RecordEvent(&event.UserLoggedIn{
        UserID:    u.ID().Value(),
        Email:     u.email.String(),
        IP:        ip,
        UserAgent: userAgent,
    })
    
    return nil
}

// Lock 锁定账户
func (u *User) Lock() {
    u.status = vo.UserStatusLocked
    u.RecordEvent(&event.UserLocked{
        UserID: u.ID().Value(),
    })
}

// Activate 激活账户
func (u *User) Activate() {
    u.status = vo.UserStatusActive
    u.RecordEvent(&event.UserActivated{
        UserID: u.ID().Value(),
    })
}

// Getters - 只读访问
func (u *User) Username() vo.Username { return u.username }
func (u *User) Email() vo.Email       { return u.email }
func (u *User) Status() vo.UserStatus { return u.status }
```

**值对象（Value Object）**

```go
// domain/user/valueobject/email.go
package valueobject

import (
    "regexp"
    "strings"
)

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 创建邮箱（工厂方法）
func NewEmail(value string) (Email, error) {
    value = strings.TrimSpace(strings.ToLower(value))
    
    if !isValidEmail(value) {
        return Email{}, kernel.FieldError("email", "无效的邮箱格式", value)
    }
    
    return Email{value: value}, nil
}

func (e Email) String() string { return e.value }

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

func isValidEmail(email string) bool {
    pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
    return regexp.MustCompile(pattern).MatchString(email)
}
```

**Repository 接口**

```go
// domain/user/repository/user_repository.go
package repository

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserRepository 用户仓储接口（Domain 层定义）
type UserRepository interface {
    // 基本 CRUD
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
    Delete(ctx context.Context, id vo.UserID) error
    
    // 查询方法
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    CountByStatus(ctx context.Context, status vo.UserStatus) (int, error)
}
```

---

### 2. Application Layer（应用层）

#### 职责

- 编排业务流程（Use Cases）
- 定义 Ports（外部依赖接口）⭐
- 定义 DTO（数据传输对象）

#### 目录结构

```
application/
├── ports/                # ⭐ Ports 接口定义
│   ├── auth/
│   │   └── token_service.go
│   ├── cache/
│   │   └── user_cache.go
│   ├── idgen/
│   │   └── generator.go
│   ├── email/
│   │   └── email_service.go
│   └── repository/       # 或者统一放在这里
│       ├── user_repository.go
│       └── tenant_repository.go
│
├── user/                 # 用户应用服务
│   ├── service.go
│   └── dto.go
│
├── auth/                 # 认证应用服务
│   ├── service.go
│   └── dto.go
│
└── unit_of_work.go       # 工作单元
```

#### 实现规范

**Ports 定义**

```go
// application/ports/auth/token_service.go
package ports

import "time"

// TokenPair 令牌对
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresAt    int64
}

// TokenClaims 令牌声明
type TokenClaims struct {
    UserID    int64  `json:"user_id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    JTI       string `json:"jti"`
    IssuedAt  int64  `json:"issued_at"`
    ExpiresAt int64  `json:"expires_at"`
}

// TokenService 令牌服务端口（Application 层定义）
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
    BlacklistToken(token string, expiresAt time.Time) error
}
```

**应用服务**

```go
// application/auth/service.go
package auth

import (
    "context"
    "errors"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
    "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
    "github.com/shenfay/go-ddd-scaffold/internal/application/ports/idgen"
    "go.uber.org/zap"
)

// AuthServiceImpl 认证应用服务实现
type AuthServiceImpl struct {
    logger         *zap.Logger
    userRepo       repository.UserRepository
    tokenService   ports.TokenService
    idGenerator    idgen.Generator
    eventPublisher EventPublisher
}

func NewAuthService(
    logger *zap.Logger,
    userRepo repository.UserRepository,
    tokenService ports.TokenService,
    idGenerator idgen.Generator,
    eventPublisher EventPublisher,
) *AuthServiceImpl {
    return &AuthServiceImpl{
        logger:         logger,
        userRepo:       userRepo,
        tokenService:   tokenService,
        idGenerator:    idGenerator,
        eventPublisher: eventPublisher,
    }
}

// AuthenticateUser 认证用户（用例）
func (s *AuthServiceImpl) AuthenticateUser(
    ctx context.Context, 
    cmd *AuthenticateUserCommand,
) (*AuthResult, error) {
    
    // 1. 查找用户
    user, err := s.findUserByIdentifier(ctx, cmd.Identifier)
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("find user failed: %w", err)
    }
    
    // 2. 调用领域方法
    err = user.Login(cmd.Password, cmd.IP, cmd.UserAgent)
    if err != nil {
        var bizErr *kernel.BusinessError
        if errors.As(err, &bizErr) {
            return nil, err
        }
        return nil, fmt.Errorf("user login failed: %w", err)
    }
    
    // 3. 保存用户状态
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("save user failed: %w", err)
    }
    
    // 4. 生成令牌
    pair, err := s.tokenService.GenerateTokenPair(
        user.ID().Value(),
        user.Username().String(),
        user.Email().String(),
    )
    if err != nil {
        return nil, fmt.Errorf("generate token failed: %w", err)
    }
    
    // 5. 发布事件
    s.eventPublisher.Publish(&event.UserLoggedIn{
        UserID:    user.ID().Value(),
        Email:     user.Email().String(),
        IP:        cmd.IP,
        UserAgent: cmd.UserAgent,
    })
    
    return &AuthResult{
        UserID:       user.ID().Value(),
        Username:     user.Username().String(),
        Email:        user.Email().String(),
        AccessToken:  pair.AccessToken,
        RefreshToken: pair.RefreshToken,
    }, nil
}
```

**DTO 定义**

```go
// application/auth/dto.go
package auth

// AuthenticateUserCommand 认证命令
type AuthenticateUserCommand struct {
    Identifier string // 邮箱或用户名
    Password   string
    IP         string
    UserAgent  string
}

// AuthResult 认证结果
type AuthResult struct {
    UserID       int64
    Username     string
    Email        string
    AccessToken  string
    RefreshToken string
}
```

---

### 3. Infrastructure Layer（基础设施层）

#### 职责

- 实现 Ports 接口（Driven Adapters）
- 数据持久化
- 外部服务集成

#### 目录结构

```
infrastructure/
├── persistence/          # 数据持久化
│   ├── dao/              # GORM 生成的 DAO
│   │   ├── user.gen.go
│   │   └── tenant.gen.go
│   └── repository/       # Repository 实现
│       ├── user_repository.go
│       └── tenant_repository.go
│
├── platform/             # 平台服务
│   ├── auth/             # 认证基础设施
│   │   ├── jwt_service.go
│   │   └── token_service_adapter.go  # ⭐ 适配器
│   ├── snowflake/        # Snowflake ID 生成器
│   │   └── node.go
│   └── notification/     # 通知服务
│       └── email_service_impl.go
│
├── cache/redis/          # Redis 缓存
│   └── redis.go
│
├── logging/              # 日志
│   └── zap.go
│
└── config/               # 配置加载
    ├── loader.go
    └── model.go
```

#### 实现规范

**Repository 实现**

```go
// infrastructure/persistence/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
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

func (r *userRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    dao, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.ID.Eq(id.Value())).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, kernel.ErrAggregateNotFound
        }
        return nil, err
    }
    
    return r.toDomain(dao)
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 保存或更新
    err := r.saveUser(tx, user)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 保存领域事件（Outbox Pattern）
    events := user.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}

func (r *userRepositoryImpl) toDomain(dao *dao.User) (*aggregate.User, error) {
    // DAO → Domain 转换
    username, _ := vo.NewUsername(dao.Username)
    email, _ := vo.NewEmail(dao.Email)
    
    user := &aggregate.User{
        // 重建聚合根
    }
    
    return user, nil
}

func (r *userRepositoryImpl) saveUser(tx *gorm.DB, user *aggregate.User) error {
    // Domain → DAO 转换并保存
}
```

**适配器模式**

```go
// infrastructure/platform/auth/token_service_adapter.go
package auth

import (
    "time"
    ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
)

// TokenServiceAdapter TokenService 端口适配器
type TokenServiceAdapter struct {
    service *JWTService  // 具体实现
}

func NewTokenServiceAdapter(service *JWTService) *TokenServiceAdapter {
    return &TokenServiceAdapter{service: service}
}

func (a *TokenServiceAdapter) GenerateTokenPair(...) (*ports_auth.TokenPair, error) {
    // 1. 调用基础设施
    pair, err := a.service.GenerateTokenPair(...)
    if err != nil {
        return nil, err
    }
    
    // 2. ⭐ 类型转换
    return &ports_auth.TokenPair{
        AccessToken:  pair.AccessToken,
        RefreshToken: pair.RefreshToken,
        ExpiresAt:    pair.ExpiresAt,
    }, nil
}

// 编译期检查
var _ ports_auth.TokenService = (*TokenServiceAdapter)(nil)
```

---

### 4. Interfaces Layer（接口层）

#### 职责

- 协议适配（HTTP/gRPC/CLI）
- 请求验证
- 响应格式化

#### 目录结构

```
interfaces/
├── http/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── routes.go
│   │   └── dto.go
│   └── user/
│       ├── handler.go
│       ├── routes.go
│       └── dto.go
│
├── grpc/
│   └── ...
│
└── messaging/
    └── ...
```

#### 实现规范

**Handler**

```go
// interfaces/http/auth/handler.go
package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
    app_auth "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
)

type Handler struct {
    authService *app_auth.AuthServiceImpl
    respHandler *http_shared.ResponseHandler
}

func NewHandler(
    authService *app_auth.AuthServiceImpl,
    respHandler *http_shared.ResponseHandler,
) *Handler {
    return &Handler{
        authService: authService,
        respHandler: respHandler,
    }
}

// Login 登录接口
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, http.StatusBadRequest, err)
        return
    }
    
    cmd := &app_auth.AuthenticateUserCommand{
        Identifier: req.Identifier,
        Password:   req.Password,
        IP:         c.ClientIP(),
        UserAgent:  c.Request.UserAgent(),
    }
    
    result, err := h.authService.AuthenticateUser(c.Request.Context(), cmd)
    if err != nil {
        h.respHandler.Error(c, http.StatusUnauthorized, err)
        return
    }
    
    h.respHandler.Success(c, LoginResponse{
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken,
    })
}
```

---

## 🔄 Module 组装

```go
// module/auth.go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. ⭐ 创建适配器
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 创建应用服务（使用 Port）
    authSvc := authApp.NewAuthService(
        infra.Logger.Named("auth"),
        userRepo,
        tokenServiceAdapter,  // ← 使用适配器
        idGenerator,
        eventPublisher,
    )
    
    // 4. 创建 Handler 和 Routes
    handler := authHTTP.NewHandler(authSvc, respHandler)
    routes := authHTTP.NewRoutes(handler, jwtSvc)
    
    return &AuthModule{routes: routes}
}
```

---

## ✅ 架构验证清单

### 代码审查

- [ ] Domain 层无任何 import（除了 shared/kernel）
- [ ] Application 层不 import Infrastructure
- [ ] Infrastructure 通过适配器实现 Ports
- [ ] Module 负责组装所有组件
- [ ] 依赖方向指向内层

### 编译验证

```bash
# 查看依赖图
goda graph ./... | dot -Tpng -o deps.png

# 检查循环依赖
importcycle ./...
```

---

## 📚 参考资源

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ports & Adapters](https://alistair.cockburn.us/hexagonal-architecture/)
- [架构总览](./architecture-overview.md)
- [Ports 模式详解](./ports-pattern-design.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
