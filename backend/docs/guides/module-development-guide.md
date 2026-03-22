# Module 开发指南

## 概述

Module 是业务功能的封装单元，负责将领域层、应用层和接口层的组件组装为一个完整的业务模块。每个 Module：

- 实现 `bootstrap.Module` 基础接口
- 可选实现 `HTTPModule`（HTTP 路由）、`EventModule`（事件订阅）等能力接口
- 在构造函数中自行组装完整依赖链
- 向 Composition Root 暴露注册入口

## 快速开始：创建一个新领域模块

以创建 `Order` 订单模块为例，按以下步骤进行：

### Step 1: 定义领域层

```
internal/domain/order/
├── aggregate/
│   └── order.go           # 聚合根
├── repository/
│   └── order_repository.go # 仓储接口
├── valueobject/
│   └── order_status.go    # 值对象
└── event/
    └── order_created.go   # 领域事件
```

**仓储接口示例**（放在领域层）：

```go
// internal/domain/order/repository/order_repository.go
package repository

type OrderRepository interface {
    Save(ctx context.Context, order *aggregate.Order) error
    FindByID(ctx context.Context, id int64) (*aggregate.Order, error)
}
```

### Step 2: 实现基础设施层

```
internal/infrastructure/persistence/
├── dao/
│   └── order.gen.go       # GORM Gen 生成的 DAO
└── repository/
    └── order_repository.go # 仓储实现
```

**仓储实现示例**：

```go
// internal/infrastructure/persistence/repository/order_repository.go
package repository

type OrderRepository struct {
    query *dao.Query
}

func NewOrderRepository(query *dao.Query) *OrderRepository {
    return &OrderRepository{query: query}
}

func (r *OrderRepository) Save(ctx context.Context, order *aggregate.Order) error {
    // 使用 dao.Query 操作数据库
}
```

### Step 3: 实现应用层

```
internal/application/order/
├── dto.go       # 数据传输对象
└── service.go   # 应用服务
```

**应用服务示例**：

```go
// internal/application/order/service.go
package order

type OrderService struct {
    uow            application.UnitOfWork
    eventPublisher kernel.EventPublisher
    snowflake      *snowflake.Node
}

func NewOrderService(
    uow application.UnitOfWork,
    eventPublisher kernel.EventPublisher,
    snowflake *snowflake.Node,
) *OrderService {
    return &OrderService{
        uow:            uow,
        eventPublisher: eventPublisher,
        snowflake:      snowflake,
    }
}
```

### Step 4: 实现接口层

```
internal/interfaces/http/order/
├── handler.go   # HTTP 处理器
└── routes.go    # 路由定义
```

**Routes 示例**：

```go
// internal/interfaces/http/order/routes.go
package order

type Routes struct {
    handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
    return &Routes{handler: handler}
}

func (r *Routes) Register(group *gin.RouterGroup) {
    orders := group.Group("/orders")
    {
        orders.POST("", r.handler.Create)
        orders.GET("/:id", r.handler.GetByID)
    }
}
```

### Step 5: 创建 Module 文件

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
type OrderModule struct {
    infra  *bootstrap.Infra
    routes *orderHTTP.Routes
}

// NewOrderModule 创建订单模块
func NewOrderModule(infra *bootstrap.Infra) *OrderModule {
    // 1. 创建 DAO Query
    daoQuery := dao.Use(infra.DB)

    // 2. 创建 UnitOfWork
    uow := application.NewUnitOfWork(infra.DB, daoQuery)

    // 3. 创建 OrderService
    orderSvc := orderApp.NewOrderService(
        uow,
        infra.EventPublisher,
        infra.Snowflake,
    )

    // 4. 创建 HTTP Handler 和 Routes
    respHandler := httpShared.NewHandler(infra.ErrorMapper)
    handler := orderHTTP.NewHandler(orderSvc, respHandler)
    routes := orderHTTP.NewRoutes(handler)

    return &OrderModule{
        infra:  infra,
        routes: routes,
    }
}

// Name 返回模块名称
func (m *OrderModule) Name() string {
    return "order"
}

// RegisterHTTP 注册 HTTP 路由
func (m *OrderModule) RegisterHTTP(group *gin.RouterGroup) {
    m.routes.Register(group)
}
```

### Step 6: 在 main.go 中注册模块

```go
// cmd/api/main.go
func main() {
    // ...
    
    // 4. 创建模块
    userMod := module.NewUserModule(infra)
    authMod := module.NewAuthModule(infra)
    orderMod := module.NewOrderModule(infra)  // 新增

    modules := []bootstrap.Module{authMod, userMod, orderMod}

    // 后续注册逻辑自动处理...
}
```

## Module 接口实现详解

### 基本结构

以 `UserModule` 为参考，一个完整的 Module 结构如下：

```go
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
```

### 实现 HTTPModule

```go
// RegisterHTTP 注册 HTTP 路由
// 实现 bootstrap.HTTPModule 接口
func (m *UserModule) RegisterHTTP(group *gin.RouterGroup) {
    m.routes.Register(group)
}
```

### 实现 EventModule

```go
// RegisterSubscriptions 注册事件订阅
// 实现 bootstrap.EventModule 接口
func (m *UserModule) RegisterSubscriptions(bus sharedAggregate.EventBus) {
    subscriber := eventHandler.NewSubscriber(bus)
    subscriber.SubscribeAll(&eventHandler.Dependencies{
        AuditSubscriber:       m.auditSubscriber,
        LoginLogSubscriber:    m.loginLogSubscriber,
        UserSideEffectHandler: m.sideEffectHandler,
    })
}
```

### 暴露公共依赖

某些模块需要向外暴露组件供其他模块或中间件使用。例如 `AuthModule` 暴露 `JWTService`：

```go
// JWTService 返回 JWT 服务供中间件使用
func (m *AuthModule) JWTService() auth.TokenService {
    return m.jwtService
}
```

## 依赖链组装模式

### DAO 初始化

每个模块独立调用 `dao.Use()` 获取 Query 对象：

```go
// 1. 创建 DAO Query
daoQuery := dao.Use(infra.DB)
```

### UnitOfWork 创建

每个模块创建自己的 UnitOfWork 实例（不共享）：

```go
// 2. 创建 UnitOfWork
uow := application.NewUnitOfWork(infra.DB, daoQuery)
```

### 完整依赖链示例（UserModule）

```go
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 1. 创建 DAO Query
    daoQuery := dao.Use(infra.DB)

    // 2. 创建 UnitOfWork
    uow := application.NewUnitOfWork(infra.DB, daoQuery)

    // 3. 创建 PasswordHasher
    passwordHasher := service.NewBcryptPasswordHasher(
        infra.Config.Security.PasswordHasher.Cost,
    )

    // 4. 创建 PasswordPolicy
    policyConfig := service.PasswordPolicyConfig{
        MinLength:           infra.Config.Security.PasswordPolicy.MinLength,
        // ...其他配置
    }
    passwordPolicy := auth.NewDefaultPasswordPolicy(policyConfig)

    // 5. 创建 JWTService
    jwtSvc := auth.NewJWTService(
        infra.Config.JWT.Secret,
        infra.Config.JWT.AccessExpire,
        infra.Config.JWT.RefreshExpire,
        "go-ddd-scaffold",
    )
    jwtSvc.SetRedisClient(infra.Redis)

    // 6. 创建 UserService
    userSvc := userApp.NewUserService(
        uow,
        infra.EventPublisher,
        passwordHasher,
        passwordPolicy,
        jwtSvc,
        infra.Snowflake,
    )

    // 7. 创建 respHandler
    respHandler := httpShared.NewHandler(infra.ErrorMapper)

    // 8. 创建 HTTP Handler 和 Routes
    handler := userHTTP.NewHandler(userSvc, respHandler)
    routes := userHTTP.NewRoutes(handler)

    // 9. 创建事件订阅器
    // ...

    return &UserModule{...}
}
```

### 跨模块依赖处理

当模块 A 需要模块 B 的服务时，有两种方式：

**方式一：模块暴露公共依赖**

```go
// AuthModule 暴露 JWTService
func (m *AuthModule) JWTService() auth.TokenService {
    return m.jwtService
}

// 在 main.go 中使用
authMod := module.NewAuthModule(infra)
jwtService := authMod.JWTService()
// 将 jwtService 传给需要它的中间件
```

**方式二：各模块独立创建（推荐）**

如果服务是无状态的，各模块可以独立创建自己的实例：

```go
// UserModule 和 AuthModule 各自创建 JWTService
jwtSvc := auth.NewJWTService(...)
```

## 最佳实践

### 1. 每个模块一个文件

```
internal/module/
├── auth.go
├── user.go
├── order.go
└── ...
```

### 2. 模块内部创建 UoW（不共享）

```go
// 正确：每个模块创建自己的 UoW
uow := application.NewUnitOfWork(infra.DB, daoQuery)

// 错误：从 Infra 共享 UoW
// uow := infra.UoW  // 不要这样做
```

### 3. Infra 字段只在 Module 构造函数中使用

```go
func NewUserModule(infra *bootstrap.Infra) *UserModule {
    // 正确：在构造函数中提取需要的组件
    db := infra.DB
    config := infra.Config
    
    // 错误：将整个 Infra 下传给 Service
    // svc := userApp.NewUserService(infra)  // 不要这样做
}
```

### 4. 接口由消费者定义

```go
// 正确：接口定义在领域层（消费方）
// internal/domain/user/repository/user_repository.go
type UserRepository interface {
    Save(ctx context.Context, user *aggregate.User) error
}

// 实现在基础设施层
// internal/infrastructure/persistence/repository/user_repository.go
type UserRepository struct { ... }
```

## 常见问题

### Q1: 模块间依赖如何处理？

**A1**: 有三种策略：

1. **公共依赖暴露**：通过模块方法暴露（如 `AuthModule.JWTService()`）
2. **独立创建**：无状态服务各模块独立创建
3. **提升到 Infra**：真正的共享组件放入 Infra（如 EventPublisher）

### Q2: 如何添加 gRPC 支持？

**A2**: 

1. 取消 `bootstrap/module.go` 中 `GRPCModule` 的注释
2. 在模块中实现 `RegisterGRPC(srv *grpc.Server)` 方法
3. 在 `cmd/api/main.go` 中添加 gRPC Server 初始化和模块遍历注册

```go
// GRPCModule 可选能力：支持 gRPC 服务注册
type GRPCModule interface {
    Module
    RegisterGRPC(srv *grpc.Server)
}
```

### Q3: 如何为模块编写测试？

**A3**: 模块测试分为两类：

**单元测试**：测试模块内的各个组件

```go
func TestUserService_Create(t *testing.T) {
    // Mock 依赖
    mockUoW := &MockUnitOfWork{}
    mockPublisher := &MockEventPublisher{}
    
    // 创建被测对象
    svc := userApp.NewUserService(mockUoW, mockPublisher, ...)
    
    // 执行测试
    err := svc.Create(ctx, dto)
    assert.NoError(t, err)
}
```

**集成测试**：测试完整模块

```go
func TestUserModule_Integration(t *testing.T) {
    // 创建测试用 Infra（使用测试数据库）
    infra := testutil.NewTestInfra(t)
    defer infra.Cleanup()
    
    // 创建模块
    userMod := module.NewUserModule(infra)
    
    // 通过 HTTP 接口测试
    router := gin.New()
    api := router.Group("/api/v1")
    userMod.RegisterHTTP(api)
    
    // 发送请求并验证
    req := httptest.NewRequest("POST", "/api/v1/users", ...)
    // ...
}
```

### Q4: 如何控制模块加载顺序？

**A4**: 在 `main.go` 中显式控制模块数组顺序：

```go
// 依赖关系：userMod 可能需要 authMod 的某些服务
// 确保 authMod 先创建
authMod := module.NewAuthModule(infra)
userMod := module.NewUserModule(infra)

modules := []bootstrap.Module{authMod, userMod}  // 顺序重要
```
