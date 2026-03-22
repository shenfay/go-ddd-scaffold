# 路由注册指南

本文档介绍 go-ddd-scaffold 项目的模块化路由注册机制和使用方法。

## 概述

项目采用 **Go 惯用的"构造函数注入 + Infra 结构体 + Module 模式"**，通过 `HTTPModule` 接口实现模块化的路由注册。每个领域模块实现 `RegisterHTTP` 方法，在 `cmd/api/main.go`（Composition Root）中显式注册。

## 核心设计原则

1. **显式注册** - 所有模块和路由在 main.go 中显式创建和注册，无隐式依赖
2. **模块自包含** - 每个模块内部自行构建完整依赖链，通过构造函数注入
3. **接口隔离** - HTTPModule、EventModule 等接口按需实现，职责分离
4. **Go 惯用风格** - 无 DI 框架，无全局单例，无 init() 魔法

## 架构设计

### 文件组织

```
cmd/api/
└── main.go              # Composition Root（模块创建和路由注册）

internal/bootstrap/
├── infra.go             # Infra 结构体（聚合基础设施依赖）
└── module.go            # Module/HTTPModule/EventModule 接口定义

internal/module/
├── auth.go              # 认证模块（实现 HTTPModule）
└── user.go              # 用户模块（实现 HTTPModule + EventModule）

internal/interfaces/http/
├── handler.go           # HTTP 响应处理器
├── user/                # 用户领域 HTTP
│   ├── routes.go        # 路由定义
│   └── handler.go       # HTTP 处理器
└── auth/                # 认证领域 HTTP
    ├── routes.go        # 路由定义
    └── handler.go       # HTTP 处理器
```

### 模块接口定义

```go
// internal/bootstrap/module.go

// Module 基础接口：所有模块必须实现
type Module interface {
    Name() string
}

// HTTPModule 可选能力：支持 HTTP 路由注册
type HTTPModule interface {
    Module
    RegisterHTTP(group *gin.RouterGroup)
}

// EventModule 可选能力：支持事件订阅注册
type EventModule interface {
    Module
    RegisterSubscriptions(bus kernel.EventBus)
}
```

### 执行流程

```
1. 程序启动，main.go 执行
   ↓
2. 加载配置，创建 Logger
   ↓
3. 创建 Infra（聚合基础设施）
   ↓
4. 创建各领域模块：NewUserModule(infra), NewAuthModule(infra)
   ↓
5. 创建 Gin Engine 和中间件
   ↓
6. 遍历模块，调用 RegisterHTTP 注册路由
   ↓
7. 遍历模块，调用 RegisterSubscriptions 注册事件订阅
   ↓
8. 启动 HTTP 服务器
```

## 使用方法

### 启动应用

#### 方式 1：使用 Makefile（推荐）

```bash
# 查看可用命令
make help

# 启动开发服务器
make run
```

#### 方式 2：直接运行

```bash
# 运行 API 服务
go run ./cmd/api/

# 或者先编译再运行
go build -o bin/api ./cmd/api
./bin/api
```

### 新增领域路由

假设要添加 **订单领域（order）**：

#### Step 1: 创建 HTTP Handler 和 Routes

```go
// internal/interfaces/http/order/handler.go
package order

import (
    "github.com/gin-gonic/gin"
    orderApp "github.com/shenfay/go-ddd-scaffold/internal/application/order"
    httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

type Handler struct {
    orderService *orderApp.OrderService
    respHandler  *httpShared.Handler
}

func NewHandler(orderSvc *orderApp.OrderService, respHandler *httpShared.Handler) *Handler {
    return &Handler{
        orderService: orderSvc,
        respHandler:  respHandler,
    }
}

func (h *Handler) CreateOrder(c *gin.Context) {
    // 处理逻辑...
}

func (h *Handler) GetOrderByID(c *gin.Context) {
    // 处理逻辑...
}
```

```go
// internal/interfaces/http/order/routes.go
package order

import "github.com/gin-gonic/gin"

type Routes struct {
    handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
    return &Routes{handler: handler}
}

// Register 注册订单路由到指定的路由组
func (r *Routes) Register(group *gin.RouterGroup) {
    orders := group.Group("/orders")
    {
        orders.POST("", r.handler.CreateOrder)
        orders.GET("/:id", r.handler.GetOrderByID)
    }
}
```

#### Step 2: 创建 Module 文件

```go
// internal/module/order.go
package module

import (
    "github.com/gin-gonic/gin"

    "github.com/shenfay/go-ddd-scaffold/internal/application"
    orderApp "github.com/shenfay/go-ddd-scaffold/internal/application/order"
    "github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
    httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
    orderHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/order"
)

// OrderModule 订单模块
// 实现 bootstrap.Module 和 bootstrap.HTTPModule 接口
type OrderModule struct {
    infra  *bootstrap.Infra
    routes *orderHTTP.Routes
}

// NewOrderModule 创建订单模块
// 内部自行构建完整依赖链
func NewOrderModule(infra *bootstrap.Infra) *OrderModule {
    // 1. 创建 DAO Query
    daoQuery := dao.Use(infra.DB)

    // 2. 创建 UnitOfWork
    uow := application.NewUnitOfWork(infra.DB, daoQuery)

    // 3. 创建 OrderService
    orderSvc := orderApp.NewOrderService(uow, infra.EventPublisher, infra.Snowflake)

    // 4. 创建 respHandler
    respHandler := httpShared.NewHandler(infra.ErrorMapper)

    // 5. 创建 HTTP Handler 和 Routes
    handler := orderHTTP.NewHandler(orderSvc, respHandler)
    routes := orderHTTP.NewRoutes(handler)

    return &OrderModule{
        infra:  infra,
        routes: routes,
    }
}

// Name 返回模块名称
// 实现 bootstrap.Module 接口
func (m *OrderModule) Name() string {
    return "order"
}

// RegisterHTTP 注册 HTTP 路由
// 实现 bootstrap.HTTPModule 接口
func (m *OrderModule) RegisterHTTP(group *gin.RouterGroup) {
    m.routes.Register(group)
}
```

#### Step 3: 在 main.go 中注册模块

```go
// cmd/api/main.go（关键代码片段）

func main() {
    // ... 加载配置、创建 Infra ...

    // 4. 创建模块
    userMod := module.NewUserModule(infra)
    authMod := module.NewAuthModule(infra)
    orderMod := module.NewOrderModule(infra)  // 新增

    modules := []bootstrap.Module{authMod, userMod, orderMod}

    // ... 创建路由和中间件 ...

    // 5.4 创建 API 路由组并注册模块路由
    api := router.Group("/api/v1")

    for _, m := range modules {
        if h, ok := m.(bootstrap.HTTPModule); ok {
            h.RegisterHTTP(api)
            logger.Info("HTTP routes registered", zap.String("module", m.Name()))
        }
    }

    // ... 启动服务器 ...
}
```

#### Step 4: 重启应用

```bash
make run
```

**完成！** ✅

## 实际代码示例

### UserModule 实现

```go
// internal/module/user.go

// UserModule 用户模块
// 实现 bootstrap.Module、bootstrap.HTTPModule 和 bootstrap.EventModule 接口
type UserModule struct {
    infra   *bootstrap.Infra
    routes  *userHTTP.Routes
    handler *userHTTP.Handler
    // 事件订阅器
    sideEffectHandler  *userEvent.SideEffectHandler
    auditSubscriber    *eventHandler.AuditSubscriber
    loginLogSubscriber *eventHandler.LoginLogSubscriber
}

// NewUserModule 创建用户模块
// 内部自行构建完整依赖链
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 1. 创建 DAO Query
    daoQuery := dao.Use(infra.DB)

    // 2. 创建 UnitOfWork
    uow := application.NewUnitOfWork(infra.DB, daoQuery)

    // 3. 创建 PasswordHasher
    passwordHasher := service.NewBcryptPasswordHasher(
        infra.Config.Security.PasswordHasher.Cost,
    )

    // ... 更多依赖构建 ...

    // 8. 创建 HTTP Handler 和 Routes
    handler := userHTTP.NewHandler(userSvc, respHandler)
    routes := userHTTP.NewRoutes(handler)

    return &UserModule{
        infra:              infra,
        routes:             routes,
        handler:            handler,
        // ...
    }
}

// Name 返回模块名称
func (m *UserModule) Name() string {
    return "user"
}

// RegisterHTTP 注册 HTTP 路由
func (m *UserModule) RegisterHTTP(group *gin.RouterGroup) {
    m.routes.Register(group)
}

// RegisterSubscriptions 注册事件订阅
func (m *UserModule) RegisterSubscriptions(bus sharedAggregate.EventBus) {
    subscriber := eventHandler.NewSubscriber(bus)
    subscriber.SubscribeAll(&eventHandler.Dependencies{
        AuditSubscriber:       m.auditSubscriber,
        LoginLogSubscriber:    m.loginLogSubscriber,
        UserSideEffectHandler: m.sideEffectHandler,
    })
}
```

### main.go 路由注册流程

```go
// cmd/api/main.go（关键代码片段）

// 4. 创建模块（替代 Factory）
userMod := module.NewUserModule(infra)
authMod := module.NewAuthModule(infra)

modules := []bootstrap.Module{authMod, userMod}

// 4.3 注册事件订阅
for _, m := range modules {
    if em, ok := m.(bootstrap.EventModule); ok {
        em.RegisterSubscriptions(infra.EventBus)
        logger.Info("Event subscriptions registered", zap.String("module", m.Name()))
    }
}

// 5. 构建路由和中间件
router := gin.New()
// ... 应用中间件 ...

// 5.4 创建 API 路由组并注册模块路由
api := router.Group("/api/v1")

for _, m := range modules {
    if h, ok := m.(bootstrap.HTTPModule); ok {
        h.RegisterHTTP(api)
        logger.Info("HTTP routes registered", zap.String("module", m.Name()))
    }
}
```

## 新架构的优势

### 与旧架构对比

| 特性 | 旧架构（Router 单例 + init()） | 新架构（HTTPModule 接口） |
|------|------------------------------|--------------------------|
| 注册方式 | 隐式（init() 魔法） | 显式（main.go 注册） |
| 依赖管理 | 全局单例、延迟初始化 | 构造函数注入、显式传递 |
| 可测试性 | 困难（全局状态） | 容易（依赖可 mock） |
| 启动顺序 | 依赖 Go 包初始化顺序 | 完全可控 |
| 代码可读性 | 需要理解 init() 机制 | 直观清晰 |
| IDE 支持 | 较弱（隐式调用） | 强（可追踪引用） |

### 核心优势

1. **显式注册** - 所有模块在 main.go 中显式创建，依赖关系一目了然
2. **无隐式依赖** - 不依赖 init() 函数和全局单例，消除"魔法"
3. **模块自包含** - 每个模块内部完成依赖构建，职责清晰
4. **Go 惯用风格** - 符合 Go 社区推荐的依赖注入方式
5. **易于测试** - 依赖通过参数传入，方便 mock
6. **IDE 友好** - 可以追踪所有引用和调用链

## 常见问题

### Q1: 如何验证路由是否注册成功？

**A**: 查看启动日志，会显示所有已注册的模块：

```
INFO    HTTP routes registered    {"module": "auth"}
INFO    HTTP routes registered    {"module": "user"}
```

或者直接访问健康检查端点：

```bash
curl http://localhost:8080/health
```

### Q2: 模块之间如何共享服务？

**A**: 如果模块 A 需要使用模块 B 的服务，有两种方式：

1. **通过 Infra 共享**：将共享服务放入 Infra 结构体
2. **模块方法暴露**：如 AuthModule 暴露 `JWTService()` 方法供其他模块使用

```go
// AuthModule 暴露 JWTService
func (m *AuthModule) JWTService() auth.TokenService {
    return m.jwtService
}

// 在 main.go 中使用
authMod := module.NewAuthModule(infra)
jwtSvc := authMod.JWTService()  // 供中间件使用
```

### Q3: 如何添加需要事件订阅的模块？

**A**: 让模块同时实现 `HTTPModule` 和 `EventModule` 接口：

```go
type MyModule struct {
    // ...
}

func (m *MyModule) RegisterHTTP(group *gin.RouterGroup) {
    // 注册 HTTP 路由
}

func (m *MyModule) RegisterSubscriptions(bus kernel.EventBus) {
    // 注册事件订阅
}
```

main.go 中的循环会自动识别并调用相应的注册方法。

## 最佳实践

1. **一个领域一个 Module 文件** - 保持职责清晰，易于维护
2. **Module 内部构建依赖链** - 不要在 main.go 中手动创建各层对象
3. **显式传递 Infra** - 所有基础设施通过 Infra 结构体传递
4. **按需实现接口** - 只实现需要的能力接口（HTTPModule、EventModule 等）
5. **在 main.go 中控制注册顺序** - 如有依赖顺序要求，调整 modules 切片顺序

## 相关文档

- [开发规范指南](development-guidelines.md) - 编码标准和代码结构
- [CLI 工具指南](cli-tool-guide.md) - 项目初始化和代码生成
- [架构设计文档](../architecture/architecture-design.md) - 整体技术架构
