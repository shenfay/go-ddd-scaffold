# Go DDD Scaffold 技术架构总览

## 🎯 架构模式

项目采用 **Clean Architecture + Ports & Adapters + DDD + 构造函数注入** 的混合架构模式。

### 核心分层

```
┌─────────────────────────────────────────────────────────┐
│              Interfaces (HTTP/gRPC/CLI)                 │
│                    ↓ depends on                         │
├─────────────────────────────────────────────────────────┤
│              Application (Use Cases + ⭐Ports)          │
│                    ↓ depends on                         │
├─────────────────────────────────────────────────────────┤
│                  Domain (Entities, Aggregates)          │
│                    ↑ implements                         │
├─────────────────────────────────────────────────────────┤
│           Infrastructure (Adapters, Persistence)        │
└─────────────────────────────────────────────────────────┘
                         ↑
              Bootstrap (Composition Root)
```

### 关键特性

✅ **Ports & Adapters** - Application 层定义接口，Infrastructure 层实现适配器  
✅ **依赖倒置** - Application 不依赖 Infrastructure 具体实现  
✅ **组合根模式** - Module 层负责组装所有依赖  
✅ **构造函数注入** - 通过构造函数传递依赖，清晰明确  

---

## 📦 核心组件

### 1. Domain Layer（领域层）

**职责：** 封装业务逻辑和状态变更规则

**包含：**
- **Aggregates** - 聚合根（如 User、Tenant）
- **Value Objects** - 值对象（如 Email、Username）
- **Domain Services** - 领域服务（如 PasswordHasher）
- **Repository Interfaces** - 仓储接口（由 Domain 定义）

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

### 2. Application Layer（应用层）

**职责：** 编排业务流程，定义 Ports 接口

**包含：**
- **Application Services** - 应用服务（如 AuthService、UserService）
- **Ports** - 外部依赖接口（如 TokenService、Generator）
- **DTOs** - 数据传输对象

**Ports 示例：**
```go
// application/ports/auth/token_service.go
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
    ParseAccessToken(token string) (*TokenClaims, error)
    ValidateToken(token string) (*TokenClaims, error)
    BlacklistToken(token string, expiresAt time.Time) error
}
```

**应用服务示例：**
```go
// application/auth/service.go
type AuthServiceImpl struct {
    tokenService   ports_auth.TokenService  // ← 依赖 Port 接口
    idGenerator    ports_idgen.Generator    // ← 依赖 Port 接口
    eventPublisher kernel.EventPublisher
}
```

### 3. Infrastructure Layer（基础设施层）

**职责：** 实现 Ports 接口的适配器

**包含：**
- **Repository Implementations** - 仓储实现（使用 GORM DAO）
- **Service Adapters** - 服务适配器（如 TokenServiceAdapter）
- **Persistence** - 持久化（PostgreSQL + Redis）
- **External Services** - 外部服务（Email、SMS 等）

**适配器示例：**
```go
// infrastructure/platform/auth/token_service_adapter.go
type TokenServiceAdapter struct {
    service TokenService  // JWTService 具体实现
}

func (a *TokenServiceAdapter) GenerateTokenPair(
    userID int64, username, email string,
) (*ports_auth.TokenPair, error) {
    pair, err := a.service.GenerateTokenPair(userID, username, email)
    // 类型转换：infra.TokenPair → ports.TokenPair
    return &ports_auth.TokenPair{
        AccessToken:  pair.AccessToken,
        RefreshToken: pair.RefreshToken,
        ExpiresAt:    pair.ExpiresAt,
    }, nil
}
```

### 4. Interfaces Layer（接口层）

**职责：** 协议适配，将外部请求转换为应用层调用

**包含：**
- **HTTP Handlers** - HTTP 处理器
- **Routes** - 路由定义
- **gRPC Services** - gRPC 服务
- **CLI Commands** - 命令行工具

**示例：**
```go
// interfaces/http/auth/handler.go
type Handler struct {
    authService authApp.AuthService
}

func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 处理错误
        return
    }
    
    result, err := h.authService.AuthenticateUser(ctx, cmd)
    // 返回响应
}
```

### 5. Bootstrap（启动引导）

**职责：** 组合根，组装所有依赖

**包含：**
- **Infra Container** - 基础设施容器
- **Modules** - 功能模块（AuthModule、UserModule）

**示例：**
```go
// module/auth.go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 注入到应用服务
    authSvc := authApp.NewAuthService(
        tokenServiceAdapter,  // ← 使用适配器，而非直接使用 jwtSvc
        infra.Snowflake,      // ← Snowflake 直接实现 Generator Port
        ...
    )
    
    return &AuthModule{...}
}
```

---

## 🗂️ 完整目录结构

```
backend/
├── cmd/                           # 应用入口
│   ├── api/main.go               # HTTP API
│   ├── worker/main.go            # Worker 服务
│   └── cli/main.go               # CLI 工具
│
├── internal/
│   ├── domain/                   # 领域层
│   │   ├── shared/kernel/        # 核心概念
│   │   ├── user/                 # 用户限界上下文
│   │   │   ├── aggregate/
│   │   │   ├── valueobject/
│   │   │   ├── event/
│   │   │   ├── service/
│   │   │   └── repository/       # ← 仓储接口（Domain 定义）
│   │   └── tenant/
│   │
│   ├── application/              # 应用层 + ⭐Ports
│   │   ├── ports/                # ⭐ Ports 接口定义
│   │   │   ├── auth/
│   │   │   │   └── token_service.go
│   │   │   ├── idgen/
│   │   │   │   └── generator.go
│   │   │   ├── cache/
│   │   │   └── email/
│   │   ├── user/
│   │   │   ├── service.go        # ✅ 使用 Ports
│   │   │   └── dto.go
│   │   ├── auth/
│   │   │   ├── service.go        # ✅ 使用 Ports
│   │   │   └── dto.go
│   │   └── unit_of_work.go
│   │
│   ├── infrastructure/           # 基础设施层（Adapters）
│   │   ├── persistence/
│   │   │   ├── dao/              # GORM 生成的 DAO
│   │   │   └── repository/       # ← Repository 适配器
│   │   ├── platform/
│   │   │   ├── auth/
│   │   │   │   ├── jwt_service.go
│   │   │   │   └── token_service_adapter.go  # ⭐ 适配器
│   │   │   └── snowflake/
│   │   │       └── node.go       # ✅ 直接实现 Generator
│   │   ├── cache/redis/
│   │   └── email/
│   │
│   ├── interfaces/               # 接口层
│   │   ├── http/
│   │   │   ├── auth/
│   │   │   └── user/
│   │   └── grpc/
│   │
│   ├── module/                   # ⭐ 组合根
│   │   ├── auth.go               # 创建适配器并注入
│   │   └── user.go
│   │
│   └── bootstrap/                # 启动引导
│       ├── infra.go              # Infra 容器
│       └── module.go             # Module 管理
│
├── pkg/                          # 公共库
├── configs/                      # 配置文件
├── migrations/                   # 数据库迁移
└── docs/                         # 文档
```

---

## 🎖️ 架构优势

### 1. 易于测试

```go
// 可以轻松创建 Mock 实现
type MockTokenService struct{}

func (m *MockTokenService) GenerateTokenPair(...) (*TokenPair, error) {
    return &TokenPair{AccessToken: "mock_token"}, nil
}
```

### 2. 技术无关性

- Application 逻辑不依赖具体技术（JWT/Redis/MySQL）
- 可以随时替换基础设施实现
- 便于技术升级和迁移

### 3. 清晰的职责分离

| 层级 | 职责 | 依赖 |
|------|------|------|
| Domain | 业务规则 | 无 |
| Application | 用例编排 | → Domain |
| Infrastructure | 技术实现 | → Application (通过适配器) |
| Interfaces | 协议适配 | → Application |
| Bootstrap | 依赖组装 | 知道所有层 |

### 4. 高度可维护

- 每层都有明确的单一职责
- 修改一层不影响其他层
- 新成员容易理解代码结构

---

## 🔧 技术栈

### 核心技术

- **语言：** Go 1.21+
- **Web 框架：** Gin
- **ORM:** GORM (Gen 生成 DAO)
- **数据库：** PostgreSQL 14+
- **缓存：** Redis 7+
- **任务队列：** Asynq
- **JWT:** golang-jwt/jwt/v5

### 架构模式

- Clean Architecture
- Ports & Adapters
- Domain-Driven Design
- Composition Root
- Dependency Injection (构造函数)

---

## 📚 相关文档

- [Ports 模式架构设计](./ports-pattern-design.md) - Ports 模式的详细说明
- [Clean Architecture 规范](./clean-architecture-spec.md) - 架构分层规范
- [DDD 设计指南](./ddd-design-guide.md) - DDD 核心概念
- [Bootstrap 模块架构](./bootstrap-module-architecture.md) - 组合根模式
- [事件驱动架构](./event-driven-architecture.md) - 事件系统设计

---

## ✅ 验证清单

### 代码检查

- [x] Application 层不 import `internal/infrastructure` 包
- [x] 所有外部依赖都有对应的 Port 接口定义
- [x] Infrastructure 层实现 Ports 接口
- [x] Module 层负责组装所有依赖
- [x] Domain 层无任何外部依赖

### 编译验证

```bash
cd backend && go build ./cmd/api
# ✅ 编译成功，无错误
```

---

## 🔄 架构演进

**当前状态：** 完整的 Ports 模式 + 组合根

**最近重构（2024）：**
- ✅ auth/service.go 完全使用 Ports 接口（方案 B）
- ✅ 移除对 Infrastructure 的直接依赖
- ✅ 通过 TokenServiceAdapter 进行类型转换
- ✅ 完全符合 Clean Architecture 规范

**下一步计划：**
- [ ] 为 Application Services 添加单元测试
- [ ] 完善文档和示例代码
- [ ] 优化错误处理机制
