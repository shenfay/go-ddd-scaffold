# 组合根模式（Composition Root）

本文档详细介绍组合根模式在 Go DDD Scaffold 中的应用，包括 Bootstrap 和 Module 的设计与实现。

## 📋 什么是组合根？

### 核心概念

**组合根（Composition Root）**是一个集中位置，负责组装应用程序的所有依赖，创建完整的对象图。

**位置：** `internal/module/`

**职责：**
- 创建基础设施组件
- 创建适配器
- 组装依赖（依赖注入）
- 注册路由和事件处理器

### 为什么需要组合根？

```
❌ 问题：依赖管理混乱

main.go
  ↓
创建各种服务...
  ↓
依赖关系分散在各处
  ↓
难以测试、难以维护
```

```
✅ 解决方案：组合根模式

main.go
  ↓
调用 NewAuthModule(infra)
  ↓
Module 集中组装所有依赖
  ↓
清晰的依赖关系，易于测试
```

---

## 🏗️ 架构设计

### 整体结构

```
┌─────────────────────────────────────┐
│          main.go                    │
│                                     │
│  func main() {                      │
│      infra := bootstrap.NewInfra()  │
│      modules := LoadModules(infra)  │
│      // ...                         │
│  }                                  │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│       bootstrap/infra.go            │
│                                     │
│  type Infra struct {                │
│      DB *gorm.DB                    │
│      Redis *redis.Client            │
│      Config *config.Config          │
│      Logger *zap.Logger             │
│      // ... 基础设施                │
│  }                                  │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│       module/auth.go                │
│                                     │
│  func NewAuthModule(infra) *Module {│
│      // 1. 创建基础设施             │
│      jwtSvc := auth.NewJWTService() │
│                                     │
│      // 2. 创建适配器 ⭐            │
│      adapter := auth.NewAdapter()   │
│                                     │
│      // 3. 组装应用服务             │
│      svc := app.NewAuthService()    │
│                                     │
│      // 4. 创建路由                 │
│      routes := http.NewRoutes()     │
│                                     │
│      return &Module{...}            │
│  }                                  │
└─────────────────────────────────────┘
```

---

## 💻 Infra 结构体详解

### 定义

```go
// bootstrap/infra.go
package bootstrap

import (
    "gorm.io/gorm"
    "github.com/redis/go-redis/v9"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

// Infra 基础设施容器
type Infra struct {
    // 核心基础设施
    DB        *gorm.DB
    Redis     *redis.Client
    Config    *config.Config
    Logger    *zap.Logger
    
    // DAO Query（GORM Gen 生成）
    DAOQuery  *dao_query.Query
    
    // 中间件
    JWTMiddleware *http_middleware.JWTMiddleware
    ErrorMapper   *response.ErrorMapper
    
    // 发布器
    EventPublisher eventstore.EventPublisher
}

// NewInfra 创建基础设施容器
func NewInfra() *Infra {
    // 1. 加载配置
    cfg := config.LoadConfig()
    
    // 2. 初始化日志
    logger := logging.NewZapLogger(cfg.Server.Mode)
    
    // 3. 初始化数据库
    db := initDatabase(cfg, logger)
    
    // 4. 初始化 Redis
    redisClient := initRedis(cfg)
    
    // 5. 初始化 DAO Query
    daoQuery := dao_query.Use(db)
    
    // 6. 初始化 JWT 中间件
    jwtMiddleware := http_middleware.NewJWTMiddleware(cfg.JWT.Secret)
    
    // 7. 初始化错误映射器
    errorMapper := response.NewErrorMapper()
    
    // 8. 初始化事件发布器
    eventPublisher := eventstore.NewEventPublisher(logger, redisClient)
    
    return &Infra{
        DB:             db,
        Redis:          redisClient,
        Config:         cfg,
        Logger:         logger,
        DAOQuery:       daoQuery,
        JWTMiddleware:  jwtMiddleware,
        ErrorMapper:    errorMapper,
        EventPublisher: eventPublisher,
    }
}

// Close 关闭所有资源
func (i *Infra) Close() {
    sqlDB, err := i.DB.DB()
    if err == nil {
        sqlDB.Close()
    }
    i.Redis.Close()
    i.Logger.Sync()
}
```

### 使用示例

```go
// main.go
func main() {
    // 创建基础设施容器
    infra := bootstrap.NewInfra()
    defer infra.Close()
    
    // 加载所有模块
    modules := LoadModules(infra)
    
    // 注册路由
    router := gin.Default()
    for _, module := range modules {
        module.RegisterRoutes(router)
    }
    
    // 启动服务器
    router.Run(":8080")
}
```

---

## 🎯 Module 模式详解

### Module 接口定义

```go
// bootstrap/module.go
package bootstrap

import (
    "github.com/gin-gonic/gin"
)

// Module 模块接口
type Module interface {
    // RegisterRoutes 注册 HTTP 路由
    RegisterRoutes(router *gin.Engine)
    
    // SubscribeEvents 订阅领域事件
    SubscribeEvents()
}

// LoadModules 加载所有模块
func LoadModules(infra *Infra) []Module {
    modules := []Module{
        NewUserModule(infra),
        NewAuthModule(infra),
        NewTenantModule(infra),
        // 添加新模块...
    }
    
    return modules
}
```

### AuthModule 完整实现

```go
// module/auth.go
package module

import (
    "github.com/gin-gonic/gin"
    "github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
    app_auth "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
    infra_auth "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
    http_auth "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/auth"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
    repository "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
)

// AuthModule 认证模块
type AuthModule struct {
    infra      *bootstrap.Infra
    jwtService *infra_auth.JWTService
    routes     *http_auth.Routes
}

// NewAuthModule 创建认证模块 ⭐ 组合根
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // === 步骤 1: 创建基础设施组件 ===
    
    // 1.1 获取 DAO Query
    daoQuery := dao_query.Use(infra.DB)
    
    // 1.2 创建 JWT Service
    jwtSvc := infra_auth.NewJWTService(
        infra.Config.JWT.Secret,
        infra.Config.JWT.AccessExpire,
        infra.Config.JWT.RefreshExpire,
        "go-ddd-scaffold",
    )
    jwtSvc.SetRedisClient(infra.Redis)
    
    // === 步骤 2: 创建适配器（关键！）⭐ ===
    
    // 2.1 TokenService 适配器
    tokenServiceAdapter := infra_auth.NewTokenServiceAdapter(jwtSvc)
    
    // 2.2 ID Generator 适配器（Snowflake 已实现 Port）
    idGeneratorAdapter := infra.Snowflake
    
    // === 步骤 3: 创建 Repository ===
    
    userRepo := repository.NewUserRepository(infra.DB, daoQuery)
    
    // === 步骤 4: 创建应用服务 ===
    
    // 4.1 创建响应 Handler
    respHandler := http_shared.NewHandler(infra.ErrorMapper)
    
    // 4.2 创建 AuthService（依赖 Ports）
    authSvc := app_auth.NewAuthService(
        infra.Logger.Named("auth"),
        userRepo,              // ← Domain 定义的 Port
        tokenServiceAdapter,   // ← Application 定义的 Port
        idGeneratorAdapter,    // ← Application 定义的 Port
        infra.EventPublisher,  // ← Application 定义的 Port
    )
    
    // === 步骤 5: 创建 HTTP 层 ===
    
    // 5.1 创建 Handler
    handler := http_auth.NewHandler(authSvc, respHandler)
    
    // 5.2 创建 Routes
    routes := http_auth.NewRoutes(handler, jwtSvc)
    
    // === 返回 Module ===
    
    return &AuthModule{
        infra:      infra,
        jwtService: jwtSvc,
        routes:     routes,
    }
}

// RegisterRoutes 注册路由
func (m *AuthModule) RegisterRoutes(router *gin.Engine) {
    m.routes.Register(&router.Group("/api/v1"))
}

// SubscribeEvents 订阅领域事件
func (m *AuthModule) SubscribeEvents() {
    // 如果需要订阅事件，在这里注册
    // m.infra.EventPublisher.Subscribe("user.created", handler)
}
```

### UserModule 实现

```go
// module/user.go
package module

import (
    "github.com/gin-gonic/gin"
    "github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
    app_user "github.com/shenfay/go-ddd-scaffold/internal/application/user"
    repository "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
    http_user "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

// UserModule 用户模块
type UserModule struct {
    infra   *bootstrap.Infra
    routes  *http_user.Routes
}

// NewUserModule 创建用户模块
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 1. 创建基础设施
    daoQuery := dao_query.Use(infra.DB)
    
    // 2. 创建 Repository
    userRepo := repository.NewUserRepository(infra.DB, daoQuery)
    
    // 3. 创建应用服务
    respHandler := http_shared.NewHandler(infra.ErrorMapper)
    
    userSvc := app_user.NewUserService(
        infra.Logger.Named("user"),
        userRepo,
        infra.EventPublisher,
        infra.Snowflake,
    )
    
    // 4. 创建 HTTP 层
    handler := http_user.NewHandler(userSvc, respHandler)
    routes := http_user.NewRoutes(handler, infra.JWTMiddleware)
    
    return &UserModule{
        infra:  infra,
        routes: routes,
    }
}

func (m *UserModule) RegisterRoutes(router *gin.Engine) {
    m.routes.Register(&router.Group("/api/v1"))
}

func (m *UserModule) SubscribeEvents() {
    // 订阅感兴趣的事件
}
```

---

## 🔑 关键设计要点

### 1. 适配器是核心 ⭐

```go
// ❌ 错误：直接使用基础设施
authSvc := app_auth.NewAuthService(
    jwtSvc,  // ← 具体实现，不是 Port
)

// ✅ 正确：使用适配器
tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
authSvc := app_auth.NewAuthService(
    tokenServiceAdapter,  // ← Port 的实现
)
```

### 2. 依赖注入方式

```go
// ✅ 推荐：构造函数注入
func NewAuthService(
    logger *zap.Logger,
    userRepo repository.UserRepository,
    tokenService ports.TokenService,
    eventPublisher EventPublisher,
) *AuthServiceImpl {
    return &AuthServiceImpl{
        logger:         logger,
        userRepo:       userRepo,
        tokenService:   tokenService,
        eventPublisher: eventPublisher,
    }
}

// ❌ 避免：setter 注入或全局变量
```

### 3. Module 的职责边界

```go
// ✅ Module 做什么
func NewAuthModule(infra *Infra) *AuthModule {
    // 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 创建适配器
    adapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 组装应用服务
    svc := app.NewAuthService(adapter, ...)
    
    // 创建路由
    routes := http.NewRoutes(handler, jwtSvc)
    
    return &AuthModule{routes: routes}
}

// ❌ Module 不做什么
// - 不包含业务逻辑
// - 不直接处理 HTTP 请求
// - 不直接访问数据库
```

### 4. Unit of Work 的创建

```go
// module/auth.go 中
func NewAuthModule(infra *Infra) *AuthModule {
    // ...
    
    // 创建 Unit of Work
    uow := application.NewUnitOfWork(infra.DB, daoQuery)
    
    // 传递给需要事务的服务
    authSvc := app_auth.NewAuthService(uow, ...)
    
    // ...
}
```

---

## 📊 依赖组装流程图

```
main.go
  ↓
bootstrap.NewInfra()
  ├─→ 加载配置
  ├─→ 初始化数据库
  ├─→ 初始化 Redis
  ├─→ 初始化日志
  └─→ 创建 Infra 容器
  
LoadModules(infra)
  ↓
NewAuthModule(infra)
  ├─→ 创建 JWT Service
  ├─→ 创建 TokenService Adapter ⭐
  ├─→ 创建 UserRepository
  ├─→ 创建 AuthService
  ├─→ 创建 Handler
  └─→ 创建 Routes
  
NewUserModule(infra)
  ├─→ 创建 UserRepository
  ├─→ 创建 UserService
  ├─→ 创建 Handler
  └─→ 创建 Routes

返回 []Module
  ↓
注册所有路由
  ↓
启动服务器
```

---

## ✅ 最佳实践

### 1. 保持 Module 精简

```go
// ✅ 正确：Module 只负责组装
func NewAuthModule(infra *Infra) *AuthModule {
    adapter := createAdapter(infra)
    service := createService(infra, adapter)
    routes := createRoutes(infra, service)
    return &AuthModule{routes: routes}
}

// ❌ 错误：Module 包含业务逻辑
func NewAuthModule(infra *Infra) *AuthModule {
    // 不应该在这里写业务逻辑
    if someCondition {
        // ...
    }
}
```

### 2. 使用适配器转换类型

```go
// ✅ 正确：适配器负责类型转换
type TokenServiceAdapter struct {
    service *JWTService
}

func (a *TokenServiceAdapter) GenerateTokenPair(...) (*ports.TokenPair, error) {
    infraPair, err := a.service.GenerateTokenPair(...)
    return &ports.TokenPair{
        AccessToken:  infraPair.AccessToken,
        RefreshToken: infraPair.RefreshToken,
    }, nil
}
```

### 3. 清晰的依赖链

```go
Infra (基础设施容器)
  ↓
Module (组合根)
  ↓
Application Service (应用服务)
  ↓
Domain Aggregate (领域聚合)
  ↓
Repository (仓储实现)
```

### 4. 延迟初始化

```go
// ✅ 正确：按需创建
func NewAuthModule(infra *Infra) *AuthModule {
    // 只在需要时才创建
    if needFeature {
        expensiveService := createExpensiveService()
    }
}

// ❌ 错误：过度初始化
func NewAuthModule(infra *Infra) *AuthModule {
    // 不管用不用都创建
    unusedService := createUnusedService()
}
```

---

## 🔄 扩展新模块

### 添加 TenantModule

```go
// module/tenant.go
package module

import (
    "github.com/gin-gonic/gin"
    "github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
    app_tenant "github.com/shenfay/go-ddd-scaffold/internal/application/tenant"
    repository "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
    http_tenant "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/tenant"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

// TenantModule 租户模块
type TenantModule struct {
    infra  *bootstrap.Infra
    routes *http_tenant.Routes
}

// NewTenantModule 创建租户模块
func NewTenantModule(infra *bootstrap.Infra) *TenantModule {
    // 1. 创建基础设施
    daoQuery := dao_query.Use(infra.DB)
    
    // 2. 创建 Repository
    tenantRepo := repository.NewTenantRepository(infra.DB, daoQuery)
    
    // 3. 创建应用服务
    respHandler := http_shared.NewHandler(infra.ErrorMapper)
    
    tenantSvc := app_tenant.NewTenantService(
        infra.Logger.Named("tenant"),
        tenantRepo,
        infra.EventPublisher,
    )
    
    // 4. 创建 HTTP 层
    handler := http_tenant.NewHandler(tenantSvc, respHandler)
    routes := http_tenant.NewRoutes(handler, infra.JWTMiddleware)
    
    return &TenantModule{
        infra:  infra,
        routes: routes,
    }
}

func (m *TenantModule) RegisterRoutes(router *gin.Engine) {
    m.routes.Register(&router.Group("/api/v1"))
}

func (m *TenantModule) SubscribeEvents() {
    // 订阅事件
}
```

### 注册到 main.go

```go
// bootstrap/module.go
func LoadModules(infra *Infra) []Module {
    modules := []Module{
        NewUserModule(infra),
        NewAuthModule(infra),
        NewTenantModule(infra),  // ← 添加新模块
    }
    
    return modules
}
```

---

## 📚 参考资源

- [Composition Root Pattern](https://blog.ploeh.dk/2011/07/28/CompositionRoot/)
- [Dependency Injection](https://martinfowler.com/articles/injection.html)
- [Go 依赖注入最佳实践](https://github.com/google/wire)
- [Clean Architecture](../design/clean-architecture-spec.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
