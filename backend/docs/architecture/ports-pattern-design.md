# Ports 模式架构设计

## 📋 概述

本项目采用 **Ports & Adapters**（六边形架构）模式，结合 **Clean Architecture** 和 **DDD** 的最佳实践，构建了一个高度可维护、可测试的 Go DDD 脚手架。

---

## 🏗️ 核心架构原则

### 1. 依赖倒置（DIP）

**Application 层不依赖 Infrastructure 具体实现**，所有外部依赖都通过 Port 接口抽象。

```
Domain Layer (最内层)
    ↓ depends on
Application Layer (定义 Ports)
    ↓ depends on  
Infrastructure Layer (实现 Ports 的适配器)
```

### 2. 分层职责

| 层级 | 职责 | 依赖方向 |
|------|------|---------|
| **Domain** | 领域模型、业务规则、领域服务 | 无依赖 |
| **Application** | 应用服务、用例、Ports 接口定义 | → Domain |
| **Infrastructure** | 技术实现、适配器、持久化 | → Application (通过适配器) |
| **Interfaces** | HTTP/gRPC/CLI 等接口层 | → Application |
| **Bootstrap** | 组合根、依赖注入、启动引导 | 知道所有层 |

---

## 🎯 Ports 接口设计

### 已定义的 Ports

#### 1. Repository Ports（领域仓储）

**位置：** `internal/domain/*/repository/`

```go
// UserRepository 用户仓储端口
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
}
```

**Infrastructure 适配器：** `internal/infrastructure/persistence/repository/user_repository.go`

#### 2. TokenService Port（认证服务）

**位置：** `internal/application/ports/auth/token_service.go`

```go
// TokenService Token 服务端口
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ParseRefreshToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
    BlacklistToken(token string, expiresAt time.Time) error
    IsTokenBlacklisted(token string) (bool, error)
}
```

**Infrastructure 适配器：** `internal/infrastructure/platform/auth/token_service_adapter.go`
- 适配器将 JWTService 转换为 Port 接口
- 进行类型转换（infra.TokenPair → ports.TokenPair）

#### 3. ID Generator Port（ID 生成器）

**位置：** `internal/application/ports/idgen/generator.go`

```go
// Generator ID 生成器端口
type Generator interface {
    Generate() (int64, error)
}
```

**Infrastructure 实现：** `internal/infrastructure/platform/snowflake/node.go`
- Snowflake Node 直接实现了 Generate() 方法
- 无需适配器，天然符合 Port 接口

#### 4. Cache Ports（缓存服务）

**位置：** `internal/application/ports/cache/`

```go
// UserCache 用户缓存端口
type UserCache interface {
    Get(ctx context.Context, userID int64) ([]byte, error)
    Set(ctx context.Context, userID int64, data []byte, expireSeconds int64) error
    Delete(ctx context.Context, userID int64) error
    Exists(ctx context.Context, userID int64) (bool, error)
}
```

**Infrastructure 适配器：** `internal/infrastructure/cache/redis/user_cache.go`

#### 5. Email Port（邮件服务）

**位置：** `internal/application/ports/email/service.go`

```go
// EmailService 邮件服务端口
type EmailService interface {
    SendWelcomeEmail(toEmail string, username string) error
    // ... 其他邮件发送方法
}
```

**Infrastructure 适配器：** `internal/infrastructure/email/smtp_service.go`

---

## 🔧 适配器模式实施

### TokenService 适配器示例

**为什么需要适配器？**

1. **类型隔离**：Port 层和 Infra 层的 TokenPair/TokenClaims 是不同类型
2. **解耦**：Application 层不知道 Infrastructure 的具体实现
3. **可测试性**：可以轻松替换为 Mock 实现

**适配器实现：**

```go
// internal/infrastructure/platform/auth/token_service_adapter.go
type TokenServiceAdapter struct {
    service TokenService  // Infrastructure 的具体实现
}

func (a *TokenServiceAdapter) GenerateTokenPair(
    userID int64, username, email string,
) (*ports_auth.TokenPair, error) {
    pair, err := a.service.GenerateTokenPair(userID, username, email)
    if err != nil {
        return nil, err
    }
    // 类型转换：infra.TokenPair → ports.TokenPair
    return &ports_auth.TokenPair{
        AccessToken:  pair.AccessToken,
        RefreshToken: pair.RefreshToken,
        ExpiresAt:    pair.ExpiresAt,
    }, nil
}
```

---

## 📦 Module 组合根模式

### AuthModule 示例

**位置：** `internal/module/auth.go`

```go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施组件
    jwtSvc := auth.NewJWTService(...)
    jwtSvc.SetRedisClient(infra.Redis)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    idGeneratorAdapter := infra.Snowflake
    
    // 3. 创建应用服务（依赖 Ports）
    authSvc := authApp.NewAuthService(
        uow,
        passwordHasher,
        tokenServiceAdapter,  // ← 使用适配器，而非直接使用 jwtSvc
        infra.EventPublisher,
        idGeneratorAdapter,
        logger,
    )
    
    return &AuthModule{...}
}
```

**关键点：**
- Module 层作为**组合根**，知道所有具体实现
- 负责创建适配器并注入到 Application Service
- Application Service 只依赖 Ports 接口

---

## ✅ 架构验证清单

### 代码检查点

- [x] Application 层不 import Infrastructure 包
- [x] 所有外部依赖都有对应的 Port 接口定义
- [x] Infrastructure 层实现 Ports 接口（直接或通过适配器）
- [x] Module 层负责组装所有依赖
- [x] Domain 层无任何外部依赖

### 编译验证

```bash
cd backend && go build ./cmd/api
# ✅ 编译成功，无错误
```

---

## 🎖️ 架构优势

### 1. 易于测试

```go
// 可以轻松创建 Mock 实现
type MockTokenService struct {
    ports_auth.TokenService
}

func (m *MockTokenService) GenerateTokenPair(...) (*ports_auth.TokenPair, error) {
    return &ports_auth.TokenPair{AccessToken: "mock_token", ...}, nil
}
```

### 2. 技术无关性

- Application 逻辑不依赖具体技术（JWT/Redis/MySQL）
- 可以随时替换基础设施实现
- 便于技术升级和迁移

### 3. 清晰的职责分离

- Domain：业务规则
- Application：用例编排
- Infrastructure：技术实现
- Interfaces：协议适配

---

## 📚 相关文档

- [Clean Architecture 规范](./clean-architecture-spec.md)
- [DDD 设计指南](./ddd-design-guide.md)
- [事件驱动架构](./event-driven-architecture.md)
- [Bootstrap 模块架构](./bootstrap-module-architecture.md)

---

## 🔄 架构演进历史

1. **初始阶段**：传统三层架构
2. **引入 DDD**：分离 Domain 层
3. **Clean Architecture**：引入 Ports & Adapters
4. **当前状态**：完整的 Ports 模式 + 组合根

**最近重构：** auth/service.go 完全使用 Ports 接口（方案 B）
- 移除了对 `infrastructure/platform/auth` 的直接依赖
- 通过 TokenServiceAdapter 进行类型转换
- 完全符合 Clean Architecture 规范
