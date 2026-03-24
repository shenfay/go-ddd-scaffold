# Ports 模式详解

本文档深入讲解 Ports & Adapters（六边形架构）模式在 Go DDD Scaffold 中的应用。

## 📋 什么是 Ports 模式？

### 核心概念

**Ports & Adapters模式**（也称为六边形架构）是一种架构模式，用于创建松耦合的应用程序结构。

```
┌─────────────────────────────────────────┐
│          Driving Adapters               │  ← HTTP/gRPC/CLI
│         (Primary Adapters)              │
└───────────────┬─────────────────────────┘
                │
┌───────────────▼─────────────────────────┐
│            Input Ports                  │  ← Application Services
├─────────────────────────────────────────┤
│          Core Business Logic            │  ← Domain Layer
├─────────────────────────────────────────┤
│           Output Ports                  │  ← Repository/Service Interfaces
└───────────────┬─────────────────────────┘
                │
┌───────────────▼─────────────────────────┐
│         Driven Adapters                 │  ← Database/Cache/External APIs
│        (Secondary Adapters)             │
└─────────────────────────────────────────┘
```

### 核心组件

| 组件 | 职责 | 位置 | 示例 |
|------|------|------|------|
| **Port** | 定义接口 | Application/Domain | `TokenService`, `UserRepository` |
| **Driving Adapter** | 调用应用 | Interfaces | HTTP Handler, gRPC Service |
| **Driven Adapter** | 被应用调用 | Infrastructure | Repository Impl, JWT Service |

---

## 🎯 为什么需要 Ports 模式？

### 解决的问题

#### ❌ 没有 Ports 的问题

```go
// ❌ 错误：Application 直接依赖 Infrastructure
package auth

import (
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
)

type AuthService struct {
    jwtService *auth.JWTService  // 依赖具体实现
    userDAO    *dao.UserDAO      // 依赖具体实现
}

// 问题：
// 1. 难以测试（需要真实的数据库和 JWT）
// 2. 无法替换实现（想换 Redis 缓存？修改代码！）
// 3. 循环依赖（Infra 想复用逻辑怎么办？）
```

#### ✅ 使用 Ports 的优势

```go
// ✅ 正确：Application 依赖 Port 接口
package auth

import (
    ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
    ports_repo "github.com/shenfay/go-ddd-scaffold/internal/application/ports/repository"
)

type AuthService struct {
    tokenService ports_auth.TokenService      // 依赖接口
    userRepo     ports_repo.UserRepository    // 依赖接口
}

// 优势：
// 1. 易于测试（可以 Mock）
// 2. 可替换实现（适配器模式）
// 3. 清晰的依赖关系
```

---

## 🏗️ Port 的分类

### 1. Input Ports（输入端口）

**定义：** 由外部调用者（Driving Adapters）使用的接口

**位置：** `application/service/`

**示例：**

```go
// application/ports/auth/token_service.go
package ports

// TokenService 令牌服务端口
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
    BlacklistToken(token string, expiresAt time.Time) error
}
```

**特点：**
- 由 Application 层定义
- 被 Interfaces 层调用
- 封装了应用的核心功能

### 2. Output Ports（输出端口）

**定义：** 被外部基础设施（Driven Adapters）实现的接口

**位置：** 
- Repository Ports: `domain/*/repository/`
- Service Ports: `application/ports/*/`

**示例 1 - Repository Port：**

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
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
    Delete(ctx context.Context, id vo.UserID) error
}
```

**示例 2 - Infrastructure Port：**

```go
// application/ports/cache/user_cache.go
package ports

import "context"

// UserCache 用户缓存端口
type UserCache interface {
    Get(ctx context.Context, userID int64) ([]byte, error)
    Set(ctx context.Context, userID int64, data []byte, expireSeconds int64) error
    Delete(ctx context.Context, userID int64) error
}
```

**特点：**
- 抽象了外部依赖（数据库、缓存、消息队列等）
- 由 Infrastructure 层实现
- 保持了 Domain 的纯净性

---

## 🔧 如何实现 Port 和 Adapter

### 步骤 1：定义 Port（接口）

```go
// application/ports/idgen/generator.go
package ports

// Generator ID 生成器端口
type Generator interface {
    Generate() (int64, error)
}
```

### 步骤 2：实现 Driven Adapter

```go
// infrastructure/platform/snowflake/node.go
package snowflake

import (
    "github.com/bwmarrin/snowflake"
    ports_idgen "github.com/shenfay/go-ddd-scaffold/internal/application/ports/idgen"
)

// Node Snowflake ID 生成器实现
type Node struct {
    node *snowflake.Node
}

// NewNode 创建 Snowflake 节点
func NewNode(nodeID int64) (*Node, error) {
    node, err := snowflake.NewNode(nodeID)
    if err != nil {
        return nil, err
    }
    
    return &Node{node: node}, nil
}

// Generate 生成唯一 ID
func (n *Node) Generate() (int64, error) {
    return n.node.Generate().Int64(), nil
}

// 编译期检查：确保实现了 Port 接口
var _ ports_idgen.Generator = (*Node)(nil)
```

### 步骤 3：在 Module 中组装

```go
// module/user.go
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 1. 创建基础设施（具体实现）
    idGenerator, err := snowflake.NewNode(1)
    if err != nil {
        panic(err)
    }
    
    // 2. ⭐ 基础设施已经是 Port 的实现，不需要额外适配
    // Node 已经实现了 ports_idgen.Generator 接口
    
    // 3. 创建应用服务（使用 Port）
    userSvc := userApp.NewUserService(
        infra.Logger.Named("user"),
        idGenerator,  // ← 传入 Port 的实现
        ...
    )
    
    return &UserModule{
        userService: userSvc,
    }
}
```

---

## 🌟 TokenService 完整实现示例

### Port 定义

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

// TokenService 令牌服务端口
type TokenService interface {
    // 生成令牌对
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    
    // 解析访问令牌
    ParseAccessToken(token string) (*TokenClaims, error)
    
    // 验证令牌
    ValidateToken(token string) (*TokenClaims, error)
    
    // 将令牌加入黑名单
    BlacklistToken(token string, expiresAt time.Time) error
}
```

### 具体实现（Infrastructure）

```go
// infrastructure/platform/auth/jwt_service.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// JWTService JWT 令牌服务具体实现
type JWTService struct {
    secret       string
    accessExpire time.Duration
    refreshExpire time.Duration
    issuer       string
    redis        *redis.Client
}

func NewJWTService(secret string, accessExpire, refreshExpire time.Duration, issuer string) *JWTService {
    return &JWTService{
        secret:       secret,
        accessExpire: accessExpire,
        refreshExpire: refreshExpire,
        issuer:       issuer,
    }
}

func (s *JWTService) SetRedisClient(client *redis.Client) {
    s.redis = client
}

// GenerateTokenPair 生成令牌对
func (s *JWTService) GenerateTokenPair(userID int64, username, email string) (*TokenPair, error) {
    now := time.Now()
    
    // Access Token
    accessClaims := &jwt.RegisteredClaims{
        Subject:   fmt.Sprintf("%d", userID),
        Issuer:    s.issuer,
        IssuedAt:  jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpire)),
        ID:        uuid.New().String(),
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(s.secret))
    if err != nil {
        return nil, err
    }
    
    // Refresh Token
    refreshClaims := &jwt.RegisteredClaims{
        Subject:   fmt.Sprintf("%d", userID),
        Issuer:    s.issuer,
        IssuedAt:  jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpire)),
        ID:        uuid.New().String(),
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(s.secret))
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresAt:    now.Add(s.accessExpire).Unix(),
    }, nil
}

// 其他方法实现...
```

### Adapter（类型转换层）

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

// NewTokenServiceAdapter 创建 TokenService 适配器
func NewTokenServiceAdapter(service *JWTService) *TokenServiceAdapter {
    return &TokenServiceAdapter{
        service: service,
    }
}

// GenerateTokenPair 生成令牌对
func (a *TokenServiceAdapter) GenerateTokenPair(
    userID int64, username, email string,
) (*ports_auth.TokenPair, error) {
    // 1. 调用基础设施的具体实现
    pair, err := a.service.GenerateTokenPair(userID, username, email)
    if err != nil {
        return nil, err
    }
    
    // 2. ⭐ 类型转换：infra.TokenPair → ports.TokenPair
    // 注意：如果 JWTService 返回的类型与 Port 定义不同，需要在这里转换
    return &ports_auth.TokenPair{
        AccessToken:  pair.AccessToken,
        RefreshToken: pair.RefreshToken,
        ExpiresAt:    pair.ExpiresAt,
    }, nil
}

// ParseAccessToken 解析访问令牌
func (a *TokenServiceAdapter) ParseAccessToken(token string) (*ports_auth.TokenClaims, error) {
    claims, err := a.service.ParseAccessToken(token)
    if err != nil {
        return nil, err
    }
    
    // 类型转换
    return &ports_auth.TokenClaims{
        UserID:    claims.UserID,
        Username:  claims.Username,
        Email:     claims.Email,
        JTI:       claims.JTI,
        IssuedAt:  claims.IssuedAt,
        ExpiresAt: claims.ExpiresAt,
    }, nil
}

// ValidateToken 验证令牌
func (a *TokenServiceAdapter) ValidateToken(token string) (*ports_auth.TokenClaims, error) {
    claims, err := a.service.ValidateToken(token)
    if err != nil {
        return nil, err
    }
    
    return &ports_auth.TokenClaims{
        UserID:    claims.UserID,
        Username:  claims.Username,
        Email:     claims.Email,
        JTI:       claims.JTI,
        IssuedAt:  claims.IssuedAt,
        ExpiresAt: claims.ExpiresAt,
    }, nil
}

// BlacklistToken 将令牌加入黑名单
func (a *TokenServiceAdapter) BlacklistToken(token string, expiresAt time.Time) error {
    return a.service.BlacklistToken(token, expiresAt)
}

// 编译期检查：确保实现了 Port 接口
var _ ports_auth.TokenService = (*TokenServiceAdapter)(nil)
```

### 在 Module 中使用

```go
// module/auth.go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(
        infra.Config.JWT.Secret,
        infra.Config.JWT.AccessExpire,
        infra.Config.JWT.RefreshExpire,
        "go-ddd-scaffold",
    )
    jwtSvc.SetRedisClient(infra.Redis)
    
    // 2. ⭐ 创建适配器
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 创建应用服务（使用 Port）
    authSvc := authApp.NewAuthService(
        uow,
        passwordHasher,
        tokenServiceAdapter,  // ← 使用适配器（Port 的实现）
        infra.EventPublisher,
        idGeneratorAdapter,
        infra.Logger.Named("auth"),
    )
    
    return &AuthModule{
        infra:      infra,
        jwtService: jwtSvc,
        routes:     routes,
    }
}
```

---

## 💡 Port 设计的最佳实践

### 1. Port 应该反映业务需求

```go
// ✅ 正确：以业务为中心命名
type UserService interface {
    RegisterUser(cmd *RegisterUserCommand) (*User, error)
    AuthenticateUser(cmd *AuthenticateUserCommand) (*AuthResult, error)
    UpdateUserProfile(cmd *UpdateProfileCommand) (*User, error)
}

// ❌ 避免：以技术为中心命名
type UserDataAccess interface {
    Insert(user map[string]interface{}) error
    Select(id int64) (map[string]interface{}, error)
    Update(id int64, data map[string]interface{}) error
}
```

### 2. Port 应该使用领域类型

```go
// ✅ 正确：使用领域类型
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
}

// ❌ 避免：使用基础设施类型
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*gorm.Model, error)
    Save(ctx context.Context, data map[string]interface{}) error
}
```

### 3. Port 应该保持精简

```go
// ✅ 正确：单一职责
type TokenService interface {
    GenerateTokenPair(...) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
}

type EmailService interface {
    SendWelcomeEmail(email string) error
    SendPasswordResetEmail(email string, token string) error
}

// ❌ 避免：大杂烩接口
type AuthService interface {
    // 20+ 个方法，包含认证、邮件、短信、日志...
}
```

### 4. 使用适配器处理类型转换

```go
// ✅ 正确：适配器负责类型转换
type TokenServiceAdapter struct {
    service *JWTService
}

func (a *TokenServiceAdapter) GenerateTokenPair(...) (*ports.TokenPair, error) {
    // 调用基础设施
    infraPair, err := a.service.GenerateTokenPair(...)
    if err != nil {
        return nil, err
    }
    
    // 类型转换
    return &ports.TokenPair{
        AccessToken:  infraPair.AccessToken,
        RefreshToken: infraPair.RefreshToken,
    }, nil
}

// ❌ 避免：在应用服务中处理基础设施类型
```

---

## 🔄 依赖流向

### 正确的依赖方向

```
Interfaces (Handler)
    ↓
Application (Service + Ports)
    ↓
Domain (Aggregate + Value Objects)
    ↓
Infrastructure (通过适配器实现 Ports)
```

### 代码示例

```go
// 1. Domain 层定义 Repository Port
// domain/user/repository/user_repository.go
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
}

// 2. Application 层使用 Port
// application/user/service.go
type UserService struct {
    userRepo repository.UserRepository  // ← 使用 Domain 定义的 Port
}

// 3. Infrastructure 层实现 Port
// infrastructure/persistence/repository/user_repository.go
type userRepositoryImpl struct {
    db *gorm.DB
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    // 具体实现
}

// 4. Module 层组装
// module/user.go
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 创建基础设施
    repo := &repositoryImpl{db: infra.DB}
    
    // 创建应用服务（注入 Port 的实现）
    svc := userApp.NewUserService(repo)
    
    return &UserModule{service: svc}
}
```

---

## 📊 Port vs Interface

### Go 接口 vs Ports 模式

| 特性 | Go 接口 | Ports 模式 |
|------|--------|-----------|
| **目的** | 多态、解耦 | 架构分层、依赖倒置 |
| **范围** | 语言特性 | 架构模式 |
| **位置** | 任何地方 | Application/Domain 层 |
| **实现** | Infrastructure 层 | Driven Adapters |

### 实际应用

```go
// Go 接口（通用解耦）
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Ports 模式（架构分层）
type TokenService interface {  // ← Port
    GenerateTokenPair(...) (*TokenPair, error)
}

type TokenServiceAdapter struct {  // ← Adapter
    service *JWTService
}
```

---

## ✅ 检查清单

### Port 设计检查

- [ ] Port 定义在 Application 或 Domain 层
- [ ] Port 使用领域类型，而非技术类型
- [ ] Port 保持精简，职责单一
- [ ] Adapter 负责类型转换
- [ ] 编译期检查：`var _ Port = (*Adapter)(nil)`

### 依赖关系检查

- [ ] Application 不 import Infrastructure
- [ ] Domain 无任何外部依赖
- [ ] Infrastructure 通过适配器实现 Ports
- [ ] Module 负责组装所有组件

---

## 📚 参考资源

- [Clean Architecture](../design/clean-architecture-spec.md)
- [架构总览](../design/architecture-overview.md)
- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Go 依赖注入最佳实践](https://github.com/google/wire)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
