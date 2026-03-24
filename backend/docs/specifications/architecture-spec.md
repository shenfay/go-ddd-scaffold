# 架构规范

本文档定义了 Go DDD Scaffold 项目的架构规范和设计原则。

## 🏗️ 架构分层

### 标准分层结构

```
┌─────────────────────────────────────────┐
│         Interfaces Layer                │  ← HTTP/gRPC/CLI
├─────────────────────────────────────────┤
│         Application Layer               │  ← Use Cases + Ports
├─────────────────────────────────────────┤
│           Domain Layer                  │  ← Business Logic
├─────────────────────────────────────────┤
│       Infrastructure Layer              │  ← Adapters
└─────────────────────────────────────────┘
            ↑
     Bootstrap (Composition Root)
```

### 各层职责

#### 1. Domain Layer（领域层）

**职责：** 封装业务逻辑和状态变更规则

**包含：**
- Aggregates（聚合根）
- Value Objects（值对象）
- Domain Services（领域服务）
- Repository Interfaces（仓储接口 - Domain 定义）

**依赖：** 无（最内层）

**示例：**
```go
// domain/user/aggregate/user.go
type User struct {
    *kernel.Entity
    username vo.Username
    email    vo.Email
    password vo.Password
}

func (u *User) CanLogin() bool {
    return u.Status() == vo.UserStatusActive
}
```

#### 2. Application Layer（应用层）

**职责：** 编排业务流程，定义 Ports 接口

**包含：**
- Application Services（应用服务）
- Ports（外部依赖接口）⭐
- DTOs（数据传输对象）

**依赖：** → Domain

**示例：**
```go
// application/ports/auth/token_service.go
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
}

// application/auth/service.go
type AuthServiceImpl struct {
    tokenService ports.TokenService  // ← 依赖 Port 接口
}
```

#### 3. Infrastructure Layer（基础设施层）

**职责：** 实现 Ports 接口的适配器

**包含：**
- Repository Implementations（仓储实现）
- Service Adapters（服务适配器）
- Persistence（持久化）
- External Services（外部服务）

**依赖：** → Application (通过适配器)

**示例：**
```go
// infrastructure/platform/auth/token_service_adapter.go
type TokenServiceAdapter struct {
    service *JWTService  // 具体实现
}

func (a *TokenServiceAdapter) GenerateTokenPair(...) (*ports.TokenPair, error) {
    // 类型转换
}
```

#### 4. Interfaces Layer（接口层）

**职责：** 协议适配，将外部请求转换为应用层调用

**包含：**
- HTTP Handlers
- gRPC Services
- CLI Commands

**依赖：** → Application

**示例：**
```go
// interfaces/http/auth/handler.go
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    c.ShouldBindJSON(&req)
    
    result, err := h.authService.AuthenticateUser(ctx, cmd)
    // 返回响应
}
```

---

## 🔑 核心设计原则

### 1. 依赖倒置原则 (DIP)

**规则：** 高层模块不应该依赖低层模块，两者都应该依赖于抽象

**实现方式：**
```go
// ✅ 正确：Application 定义 Port，Infrastructure 实现
package ports  // Application 层

type TokenService interface {
    GenerateTokenPair(...) (*TokenPair, error)
}

package auth  // Infrastructure 层

type TokenServiceAdapter struct {
    service TokenService  // 实现 Port
}
```

**禁止：**
```go
// ❌ 错误：Application 直接依赖 Infrastructure
package auth

import "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"

type AuthServiceImpl struct {
    tokenService *auth.JWTService  // ❌ 依赖具体实现
}
```

### 2. 单一职责原则 (SRP)

**规则：** 每个类应该只有一个引起变化的原因

**示例：**
```go
// ✅ 正确：职责分离
type UserService struct { ... }      // 用户管理
type AuthService struct { ... }      // 认证管理
type EmailService struct { ... }     // 邮件发送

// ❌ 错误：职责混杂
type UserService struct {
    // 包含用户管理、认证、邮件发送等所有功能
}
```

### 3. 开闭原则 (OCP)

**规则：** 软件实体应该对扩展开放，对修改关闭

**实现方式：** 使用接口和策略模式
```go
type PasswordHasher interface {
    Hash(password string) (string, error)
    Verify(password, hash string) bool
}

// 可以轻松扩展不同的实现
type BcryptPasswordHasher struct { ... }
type Argon2PasswordHasher struct { ... }
```

### 4. 接口隔离原则 (ISP)

**规则：** 客户端不应该被迫依赖它不使用的接口

**示例：**
```go
// ✅ 正确：接口细分
type Reader interface { Read() []byte }
type Writer interface { Write([]byte) error }

// ❌ 错误：大接口
type DataProcessor interface {
    Read() []byte
    Write([]byte) error
    Process() error
    Validate() bool
    // ... 太多方法
}
```

---

## 📦 Module 组装规范

### Composition Root 模式

**位置：** `internal/module/`

**职责：**
1. 创建基础设施组件
2. 创建适配器
3. 组装依赖
4. 注册路由和事件处理器

### AuthModule 示例

```go
// module/auth.go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 步骤 1: 创建基础设施
    daoQuery := dao.Use(infra.DB)
    uow := application.NewUnitOfWork(infra.DB, daoQuery)
    
    jwtSvc := auth.NewJWTService(
        infra.Config.JWT.Secret,
        infra.Config.JWT.AccessExpire,
        infra.Config.JWT.RefreshExpire,
        "go-ddd-scaffold",
    )
    jwtSvc.SetRedisClient(infra.Redis)
    
    // 步骤 2: ⭐ 创建适配器（关键！）
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    idGeneratorAdapter := infra.Snowflake
    
    // 步骤 3: 创建 PasswordHasher
    passwordHasher := service.NewBcryptPasswordHasher(
        infra.Config.Security.PasswordHasher.Cost,
    )
    
    // 步骤 4: 创建应用服务（使用 Ports）
    authSvc := authApp.NewAuthService(
        uow,
        passwordHasher,
        tokenServiceAdapter,  // ← 使用适配器，而非直接使用 jwtSvc
        infra.EventPublisher,
        idGeneratorAdapter,
        infra.Logger.Named("auth"),
    )
    
    // 步骤 5: 创建 Handler 和 Routes
    respHandler := httpShared.NewHandler(infra.ErrorMapper)
    handler := authHTTP.NewHandler(authSvc, respHandler)
    routes := authHTTP.NewRoutes(handler, jwtSvc)
    
    return &AuthModule{
        infra:      infra,
        jwtService: jwtSvc,
        routes:     routes,
    }
}
```

### 关键点说明

1. **适配器创建是核心** ⭐
   - Module 层知道所有具体实现
   - 负责将具体实现转换为 Port 接口
   - Application Service 只依赖 Port

2. **依赖注入方式**
   - 使用构造函数注入
   - 所有依赖在创建时传入
   - 清晰明确，易于测试

3. **职责分离**
   - Module 层：组装依赖
   - Application 层：业务编排
   - Infrastructure 层：技术实现

---

## 🎯 Ports 设计规范

### 什么是 Port？

Port 是 Application 层定义的接口，用于抽象外部依赖。

### Port 分类

#### 1. Repository Ports

**定义位置：** `domain/*/repository/`

```go
// domain/user/repository/user_repository.go
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
}
```

**实现位置：** `infrastructure/persistence/repository/`

#### 2. Service Ports

**定义位置：** `application/ports/*/`

```go
// application/ports/auth/token_service.go
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
    BlacklistToken(token string, expiresAt time.Time) error
}
```

**实现位置：** `infrastructure/platform/auth/`（通过适配器）

#### 3. Infrastructure Ports

**定义位置：** `application/ports/*/`

```go
// application/ports/idgen/generator.go
type Generator interface {
    Generate() (int64, error)
}

// application/ports/cache/user_cache.go
type UserCache interface {
    Get(ctx context.Context, userID int64) ([]byte, error)
    Set(ctx context.Context, userID int64, data []byte, expireSeconds int64) error
    Delete(ctx context.Context, userID int64) error
}
```

---

## 🔄 适配器模式规范

### 为什么需要适配器？

1. **类型隔离** - Port 层和 Infra 层的类型不同
2. **解耦** - Application 不知道 Infrastructure 的存在
3. **可替换** - 可以轻松更换实现

### TokenServiceAdapter 完整示例

```go
// infrastructure/platform/auth/token_service_adapter.go
package auth

import (
    "time"
    ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
)

// TokenServiceAdapter TokenService 端口适配器
type TokenServiceAdapter struct {
    service TokenService  // JWTService 具体实现
}

// NewTokenServiceAdapter 创建 TokenService 适配器
func NewTokenServiceAdapter(service TokenService) *TokenServiceAdapter {
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

// 其他方法类似...

// 编译期检查：确保实现了 Port 接口
var _ ports_auth.TokenService = (*TokenServiceAdapter)(nil)
```

---

## 🚫 禁止事项

### 依赖关系禁忌

```go
// ❌ 禁止：Application 导入 Infrastructure
import "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/..."

// ❌ 禁止：Domain 导入 Application
import "github.com/shenfay/go-ddd-scaffold/internal/application/..."

// ❌ 禁止：Infrastructure 导入 Application
import "github.com/shenfay/go-ddd-scaffold/internal/application/..."
```

### 代码组织禁忌

```go
// ❌ 禁止：在 Domain 层处理 HTTP 请求
func (u *User) HandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
    // Domain 层不应该知道 HTTP
}

// ❌ 禁止：在 Application 层直接操作数据库
func (s *UserService) queryFromDB(sql string) {
    // Application 层不应该直接访问数据库
}

// ❌ 禁止：在 Infrastructure 层编写业务逻辑
if user.Password().Value() != inputPassword {
    // 密码验证应该在 Domain 层的 Password 值对象中
}
```

---

## 📊 架构验证清单

### 代码审查检查项

- [ ] Application 层没有 import Infrastructure 包
- [ ] 所有外部依赖都有对应的 Port 接口
- [ ] Infrastructure 层通过适配器实现 Ports
- [ ] Module 层负责创建适配器并组装
- [ ] Domain 层无任何外部依赖
- [ ] 依赖方向符合规范（指向内层）
- [ ] 每层职责清晰，无职责混杂

### 编译验证

```bash
# 验证编译
cd backend && go build ./cmd/api

# 检查循环依赖
go get golang.org/x/tools/cmd/importcycle
importcycle ./...

# 查看依赖图
go get github.com/loov/goda
goda graph ./... | dot -Tpng -o deps.png
```

---

## 📚 参考资源

- [Clean Architecture](../../design/clean-architecture-spec.md)
- [Ports 模式详解](../../design/ports-pattern-design.md)
- [架构总览](../../design/architecture-overview.md)
- [开发规范](./development-spec.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
