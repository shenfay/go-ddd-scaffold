# 应用容器重构说明

## 文档概述

本文档记录了 `go-ddd-scaffold` 项目中 Container（应用容器）组件的职责重构过程，包括重构背景、设计方案、实施步骤和影响范围。

## 重构背景

### 原有设计问题

在早期的架构设计中，Container 承担了过多的职责：

```go
// ❌ 旧设计：Container 职责过重
type Container interface {
    // === 基础设施访问 ===
    GetDB() *sql.DB
    GetRedis() *redis.Client
    GetCache() CacheClient
    GetLogger(name string) *zap.Logger
    GetConfig() *config.AppConfig
    
    // === 路由访问 ===
    GetRouter() *gin.Engine
    
    // === 服务注册表（问题所在）===
    RegisterService(name string, service interface{})
    GetService(name string) interface{}
    
    // === 生命周期 ===
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

**存在的问题：**

1. **类型安全问题**
   ```go
   // 运行时才能发现错误
   handler := container.GetService("user.command.create").(*CreateUserHandler)
   // ↑ 如果类型断言失败，运行时 panic
   ```

2. **职责不清晰**
   - Container 既管理基础设施，又管理业务服务
   - 模糊了 Composition Root 和 Infrastructure 的边界

3. **依赖关系隐式**
   ```go
   // 无法从代码中直接看出依赖关系
   userHandler := NewHandler(
       container.GetService("user.command.create"),
       container.GetService("user.query.get"),
   )
   ```

4. **不符合 DDD 原则**
   - Container 侵入领域层
   - 领域组件的创建应该由 Composition Root 负责

---

## 新设计方案

### 核心设计理念

**Container 职责限定为：**
- ✅ HTTP 层路由管理
- ✅ 基础设施管理
- ✅ 生命周期管理
- ❌ **不负责**领域层的注册和初始化

```go
// ✅ 新设计：职责精简
type Container interface {
    // === 基础设施访问 ===
    GetDB() *sql.DB
    GetRedis() *redis.Client
    GetCache() CacheClient
    GetLogger(name string) *zap.Logger
    GetConfig() *config.AppConfig
    
    // === HTTP 路由访问 ===
    GetRouter() *gin.Engine
    
    // === 生命周期管理 ===
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

### 职责分离

```
┌─────────────────────────────────────────┐
│          Bootstrap                      │
│     (Composition Root)                  │
│  ┌─────────────────────────────────┐   │
│  │ ✅ 创建领域组件                 │   │
│  │ ✅ 组装依赖关系                 │   │
│  │ ✅ 初始化应用服务               │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
         ↓ 直接使用字段引用
┌─────────────────────────────────────────┐
│          Container                      │
│  ┌─────────────────────────────────┐   │
│  │ ✅ HTTP 路由管理                │   │
│  │    - RegisterRoutes()           │   │
│  │    - AutoLoadDomains()          │   │
│  │    - BuildRoutes()              │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │ ✅ 基础设施管理                  │   │
│  │    - GetDB()                    │   │
│  │    - GetRedis()                 │   │
│  │    - GetLogger()                │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │ ✅ 生命周期管理                  │   │
│  │    - Start(ctx)                 │   │
│  │    - Stop(ctx)                  │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

---

## 实施细节

### 1. Bootstrap 按领域分组

```go
// bootstrap/bootstrap.go
type Bootstrap struct {
    container container.Container
    config    *config.AppConfig
    logger    *zap.Logger
    httpDeps  *http.Dependencies
    
    // === 用户领域组件（按领域分组）===
    user struct {
        createHandler     *commands.CreateUserHandler
        updateHandler     *commands.UpdateUserHandler
        activateHandler   *commands.ActivateUserHandler
        deactivateHandler *commands.DeactivateUserHandler
        changePassHandler *commands.ChangePasswordHandler
        getHandler        *queries.GetUserHandler
        listHandler       *queries.ListUsersHandler
    }
}
```

**优势：**
- ✅ 结构清晰，一目了然
- ✅ 类型安全，无需类型断言
- ✅ IDE 友好，支持自动补全

---

### 2. 领域组件创建

```go
// bootstrap/user_domain.go
func (b *Bootstrap) initUserDomain(ctx context.Context) error {
    b.logger.Info("Initializing user domain...")
    
    // 1. 创建基础设施服务
    passwordHasher := &user.SimplePasswordHasher{}
    eventPublisher := NewInMemoryEventPublisher(baseLogger.Named("events"))
    
    // 2. 创建仓储层
    db := b.container.GetDB()
    userRepo := persistence.NewUserRepository(db)
    
    // 3. 创建 CQRS Handlers（直接赋值给 Bootstrap 字段）
    b.user.createHandler = commands.NewCreateUserHandler(
        userRepo, passwordHasher, eventPublisher,
    )
    b.user.updateHandler = commands.NewUpdateUserHandler(userRepo, eventPublisher)
    b.user.activateHandler = commands.NewActivateUserHandler(userRepo, eventPublisher)
    b.user.deactivateHandler = commands.NewDeactivateUserHandler(userRepo, eventPublisher)
    b.user.changePassHandler = commands.NewChangePasswordHandler(
        userRepo, passwordHasher, eventPublisher,
    )
    
    b.user.getHandler = queries.NewGetUserHandler(userRepo)
    b.user.listHandler = queries.NewListUsersHandler(userRepo)
    
    return nil
}
```

**关键点：**
- ✅ 不再调用 `container.RegisterService()`
- ✅ 直接赋值给 `b.user.*` 字段
- ✅ 类型明确，无需接口转换

---

### 3. 跨领域调用处理

```go
// bootstrap/initializeInterfaces.go
func (b *Bootstrap) initializeInterfaces(ctx context.Context) error {
    // 创建 HTTP Handler（响应处理）
    respHandler := http.NewHandler(apperrors.NewErrorMapper())
    
    // 创建用户领域 HTTP Handler（业务处理）
    // ✅ 直接使用 Bootstrap 中持有的领域组件，类型安全
    userHandler := userHttp.NewHandler(
        userHttp.WithCreateUserHandler(b.user.createHandler),
        userHttp.WithUpdateUserHandler(b.user.updateHandler),
        userHttp.WithActivateUserHandler(b.user.activateHandler),
        userHttp.WithDeactivateUserHandler(b.user.deactivateHandler),
        userHttp.WithChangePasswordHandler(b.user.changePassHandler),
        userHttp.WithGetUserHandler(b.user.getHandler),
        userHttp.WithListUsersHandler(b.user.listHandler),
        userHttp.WithResponseHandler(respHandler),
    )
    
    // 构建路由...
    return nil
}
```

**对比旧方案：**

```go
// ❌ 旧方案：类型不安全
userHandler := userHttp.NewHandler(
    userHttp.WithCreateUserHandler(
        b.container.GetService("user.command.create").(*commands.CreateUserHandler),
    ),
    // ... 其他 7 个类似的调用
)

// ✅ 新方案：类型安全
userHandler := userHttp.NewHandler(
    userHttp.WithCreateUserHandler(b.user.createHandler),
    // ... 其他 7 个直接的字段引用
)
```

---

### 4. 未来扩展：跨领域依赖

当 Order 领域需要调用 User 领域时：

```go
// bootstrap/bootstrap.go
type Bootstrap struct {
    // ... 用户领域组件
    user struct {
        getHandler *queries.GetUserHandler
        // ...
    }
    
    // ... 订单领域组件
    order struct {
        createHandler *commands.CreateOrderHandler
        // ...
    }
}

// initializeDomains 注意创建顺序
func (b *Bootstrap) initializeDomains(ctx context.Context) error {
    // 1. 先创建 User Handlers
    if err := b.initUserDomain(ctx); err != nil {
        return err
    }
    
    // 2. 再创建 Order Handlers，传入 User Handler 作为依赖
    if err := b.initOrderDomain(ctx); err != nil {
        return err
    }
    
    return nil
}

// bootstrap/order_domain.go
func (b *Bootstrap) initOrderDomain(ctx context.Context) error {
    // 创建 Order Repository
    orderRepo := persistence.NewOrderRepository(db)
    
    // 创建 Order Handlers，直接传入 User Handler
    b.order.createHandler = commands.NewCreateOrderHandler(
        orderRepo,
        b.user.getHandler,  // ← 直接传递，类型安全
    )
    
    return nil
}
```

---

## 重构效果

### 代码统计

| 指标 | 旧方案 | 新方案 | 改进 |
|------|--------|--------|------|
| **代码行数** | ~286 行 | ~271 行 | -15 行 |
| **类型断言** | 7 处 | 0 处 | ✅ 消除 |
| **运行时错误风险** | 有（panic） | 无 | ✅ 编译期检查 |
| **IDE 支持** | 弱 | 强 | ✅ 智能提示 |

### 架构优势

#### 1. 类型安全提升 100%

```go
// ❌ 旧方案：运行时可能 panic
handler := container.GetService("user.command.create").(*CreateUserHandler)

// ✅ 新方案：编译期检查
handler := b.user.createHandler
```

#### 2. 职责边界清晰

```
Container: "我负责提供基础设施和管理路由"
Bootstrap: "我负责创建和组装所有依赖"
Domain:    "我只关注业务逻辑"
HTTP:      "我只负责路由和响应"
```

#### 3. 依赖关系显式

```go
// ✅ 一眼就能看出依赖关系
userHandler := userHttp.NewHandler(
    WithCreateUserHandler(b.user.createHandler),
    WithGetUserHandler(b.user.getHandler),
    // ...
)
```

#### 4. 更符合 DDD 原则

- ✅ Container 不侵入领域层
- ✅ 领域组件的创建由 Composition Root 负责
- ✅ 依赖关系清晰，易于测试

---

## 迁移指南

### 从旧方案迁移到新方案

#### Step 1: 修改 Container 接口

```go
// internal/container/container.go

// ❌ 移除
type Container interface {
    RegisterService(name string, service interface{})
    GetService(name string) interface{}
    // ...
}

// ✅ 保留
type Container interface {
    GetDB() *sql.DB
    GetRedis() *redis.Client
    GetCache() CacheClient
    GetLogger(name string) *zap.Logger
    GetConfig() *config.AppConfig
    GetRouter() *gin.Engine
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

#### Step 2: 重构 Bootstrap

```go
// internal/bootstrap/bootstrap.go

type Bootstrap struct {
    container container.Container
    // ...
    
    // ✅ 新增：按领域分组
    user struct {
        createHandler     *commands.CreateUserHandler
        updateHandler     *commands.UpdateUserHandler
        // ...
    }
}
```

#### Step 3: 更新领域初始化

```go
// internal/bootstrap/user_domain.go

// ❌ 移除
b.container.RegisterService("user.command.create", createUserHandler)

// ✅ 改为
b.user.createHandler = commands.NewCreateUserHandler(...)
```

#### Step 4: 更新接口层

```go
// internal/bootstrap/initializeInterfaces.go

// ❌ 移除
userHandler := userHttp.NewHandler(
    userHttp.WithCreateUserHandler(
        b.container.GetService("user.command.create").(*commands.CreateUserHandler),
    ),
)

// ✅ 改为
userHandler := userHttp.NewHandler(
    userHttp.WithCreateUserHandler(b.user.createHandler),
)
```

---

## 常见问题

### Q1: Bootstrap 字段会很多怎么办？

**A:** 使用嵌套 struct 按领域分组：

```go
type Bootstrap struct {
    user struct {
        createHandler     *commands.CreateUserHandler
        updateHandler     *commands.UpdateUserHandler
        // ...
    }
    
    tenant struct {
        createHandler *commands.CreateTenantHandler
        // ...
    }
    
    order struct {
        createHandler *commands.CreateOrderHandler
        // ...
    }
}

// 使用时
b.user.createHandler.Handle(...)
b.tenant.createHandler.Handle(...)
```

---

### Q2: 如何处理跨领域调用？

**A:** 在 `initializeDomains()` 中注意创建顺序：

```go
func (b *Bootstrap) initializeDomains(ctx context.Context) error {
    // 1. 先创建基础领域（User、Tenant）
    if err := b.initUserDomain(ctx); err != nil {
        return err
    }
    if err := b.initTenantDomain(ctx); err != nil {
        return err
    }
    
    // 2. 再创建依赖其他领域的领域（Order）
    if err := b.initOrderDomain(ctx); err != nil {
        return err
    }
    
    return nil
}
```

---

### Q3: 是否需要引入 DI 容器？

**A:** 不需要。当前方案已经足够简洁：

- ✅ 类型安全
- ✅ 零额外依赖
- ✅ 易于理解和维护
- ✅ 符合 Composition Root 原则

引入 DI 容器（如 Google Wire）会增加复杂度，违背轻量级设计原则。

---

## 最佳实践

### 1. 领域组件命名规范

```go
// 字段命名
b.user.createHandler     // Command Handler
b.user.getHandler        // Query Handler

// 方法命名
b.initUserDomain()       // 领域初始化方法
b.initializeInterfaces() // 接口层初始化方法
```

### 2. 依赖传递原则

```go
// ✅ 推荐：直接传递具体类型
orderHandler := commands.NewCreateOrderHandler(
    orderRepo,
    b.user.getHandler,  // 具体类型
)

// ⚠️ 避免：使用接口解耦（除非必要）
// 只有当循环依赖时才考虑定义接口
```

### 3. 测试建议

```go
// 单元测试：可以直接 Mock Bootstrap
func TestCreateUser(t *testing.T) {
    boot := &Bootstrap{}
    boot.user.createHandler = mockCreateUserHandler
    
    // 直接调用，无需启动容器
    result, err := boot.user.createHandler.Handle(ctx, cmd)
    assert.NoError(t, err)
}
```

---

## 总结

### 重构收益

✅ **类型安全**：消除运行时类型断言  
✅ **职责清晰**：Container 和 Bootstrap 各司其职  
✅ **易于维护**：依赖关系显式，IDE 友好  
✅ **符合 DDD**：Composition Root 模式得到正确实施  

### 适用场景

本方案特别适合：
- ✅ 单体应用
- ✅ 领域数量适中（< 10 个）
- ✅ 追求简洁和类型安全
- ✅ 不需要动态模块加载

### 未来演进

如果项目发展到需要：
- 动态模块加载
- 插件化架构
- 微服务拆分

可以考虑引入更强大的 DI 容器或服务网格方案。

---

## 参考文档

- [DDD+CQRS 设计指南](./ddd-cqrs-design-guide.md)
- [路由自动加载机制](../guides/route-auto-loading.md)
- [Composition Root 模式](https://martinfowler.com/articles/injection.html)

---

**最后更新时间**: 2026-03-13  
**版本**: v1.0
