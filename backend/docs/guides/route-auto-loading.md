# 路由自动加载指南

本文档介绍 go-ddd-scaffold 项目的路由自动加载机制和使用方法。

## 概述

项目采用**领域驱动的路由自动注册模式**，通过 Go 的 `init()` 函数实现零侵入的领域路由自动注册。新增领域时无需修改 `main.go`，只需在 `domains.go` 中添加导入即可。

## 核心设计原则

1. **配置分离** - 端口和 API 前缀在 main.go 中定义，router.go 不包含任何默认配置
2. **main.go 稳定性** - main.go 文件永久固定，无需因新增领域而修改
3. **依赖倒置** - 接口层主动依赖领域层，符合 DDD 原则
4. **自动发现** - 利用 init() 函数自动注册，无需手动调用
5. **延迟注册** - 使用 pendingRegs 暂存机制，确保 main.go 的配置优先生效
6. **防止重复** - 使用标志位确保路由只注册一次

## 架构设计

### 文件组织

```
cmd/api/
├── main.go          # API 启动入口（永不修改）
└── domains.go       # 领域导入清单（唯一修改点）

internal/interfaces/http/
├── router.go        # 路由总线核心
├── handler.go       # HTTP 响应处理器
└── user/            # 领域路由目录
    ├── routes.go    # 路由注册（含 init()）
    └── handler.go   # HTTP 处理器
```

### 执行流程

```
1. 程序启动
   ↓
2. 导入 domains.go → 触发 _ "user" 导入
   ↓
3. user/routes.go 的 init() 执行
   ↓
4. 调用 MustGetRouter() → 创建临时 Router，保存注册函数到 pendingRegs
   ↓
5. main.go 调用 GetRouter(config) → 创建真正的 Router（使用 main 中的配置）
   ↓
6. GetRouter 应用所有 pendingRegs 中暂存的注册函数
   ↓
7. main.go 调用 router.Build()
   ↓
8. Build() 遍历所有已注册的 registrar 并执行
   ↓
9. 完整路由构建完成
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

**注意**：`make run` 使用的是 `go run ./cmd/api/`，这会包含 `domains.go`。

#### 方式 2：手动启动

```bash
# ✅ 正确 - 运行整个包
go run ./cmd/api/

# ❌ 错误 - 只运行 main.go，会返回 404
go run ./cmd/api/main.go

# 或者先编译再运行
go build -o bin/api ./cmd/api
./bin/api
```

**为什么不能直接运行 main.go？**

因为 `domains.go` 包含了领域的 `_ "..."` 导入，这些导入会触发各领域的 `init()` 函数执行。如果只运行 `main.go`，这些导入不会生效，导致路由没有注册，最终返回 404 错误。

### 新增领域路由

假设要添加 **订单领域（order）**：

#### Step 1: 在 domains.go 添加导入

```go
// cmd/api/domains.go
package main

import (
    // 用户领域
    _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
    
    // 新增领域时，在这里添加导入：
    _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/order"
)
```

#### Step 2: 创建 order/routes.go

```go
// internal/interfaces/http/order/routes.go
package order

import (
    "github.com/gin-gonic/gin"
    httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// RegisterRoutes 注册订单领域路由
func RegisterRoutes(router *gin.RouterGroup, handler *httpiface.Handler) {
    h := NewOrderHandler(nil, handler)
    
    // 订单资源路由
    v1 := router.Group("v1/orders")
    {
        v1.POST("", h.CreateOrder)
        v1.GET("/:id", h.GetOrderByID)
        v1.PUT("/:id", h.UpdateOrder)
        v1.DELETE("/:id", h.DeleteOrder)
    }
}

// init 自动注册到全局路由总线
// 注意：MustGetRouter() 会在 init 时暂存注册函数，等待 main.go 初始化真正的 Router
func init() {
    httpiface.MustGetRouter().Register(RegisterRoutes)
}
```

#### Step 3: 重启应用

```bash
make run
```

**完成！** ✅ `main.go` 无需任何修改。

## 技术细节

### Router 延迟注册机制

```go
var (
    globalRouter   *Router
    routerOnce     sync.Once
    pendingRegs    []func(*Router) // 存储 init 时注册的函数，延迟初始化时使用
)

// GetRouter 获取全局路由总线实例（单例）
// config 参数仅在首次调用时生效，必须由 main.go 提供配置
func GetRouter(config *RouterConfig) *Router {
    routerOnce.Do(func() {
        globalRouter = NewRouter(config)
        
        // 应用所有在 init 中注册的函数
        for _, regFunc := range pendingRegs {
            regFunc(globalRouter)
        }
        pendingRegs = nil // 清理内存
    })
    return globalRouter
}

// MustGetRouter 获取全局路由总线实例（用于模块注册）
// 如果尚未初始化，会将注册函数暂存到 pendingRegs
func MustGetRouter() *Router {
    // 如果已经初始化，直接返回
    if globalRouter != nil {
        return globalRouter
    }
    
    // 否则返回一个临时的 Router 用于收集注册函数
    tempRouter := &Router{
        registrars: make([]RouteRegistrar, 0),
    }
    
    // 包装 Register 方法，使其能延迟执行
    wrappedReg := func(r *Router) {
        for _, reg := range tempRouter.registrars {
            r.Register(reg)
        }
    }
    
    pendingRegs = append(pendingRegs, wrappedReg)
    return tempRouter
}
```

**工作原理：**
1. `init()` 调用 `MustGetRouter()` → 创建临时 Router，保存注册函数到 `pendingRegs`
2. `main.go` 调用 `GetRouter(config)` → 创建真正的 Router（使用 main 中的配置）
3. `GetRouter` 应用所有 `pendingRegs` 中暂存的注册函数
4. 确保 main.go 的配置优先级高于 init 的自动注册

### Router 防重复注册机制

```go
type Router struct {
    config      *RouterConfig
    ginEngine   *gin.Engine
    registrars  []RouteRegistrar
    handler     *Handler
    initialized bool  // 防止重复初始化
}

// Build 方法确保只执行一次
func (r *Router) Build(deps *Dependencies) *gin.Engine {
    if !r.initialized {
        // 注册所有路由
        r.initialized = true
    }
    return r.ginEngine
}
```

### init() 自动注册时机

Go 的包初始化顺序保证：
1. 首先初始化导入的包（`domains.go` 的 `_ "user"`）
2. 执行被导入包的 `init()` 函数 → 调用 `MustGetRouter().Register()`，暂存到 pendingRegs
3. 最后执行 `main()` 函数 → 调用 `GetRouter(config)`，应用 pendingRegs 中的注册函数

这确保了：
- main.go 的配置优先生效（端口、API 前缀）
- 所有领域路由正确注册
- router.go 不包含任何默认配置

## 常见问题

### Q1: 为什么会出现 404 错误？

**A**: 最常见的原因是使用了错误的启动命令：

```bash
# ❌ 错误
go run ./cmd/api/main.go

# ✅ 正确
go run ./cmd/api/
```

### Q2: 如何验证路由是否注册成功？

**A**: 查看启动日志，会显示所有已注册的路由：

```
[GIN-debug] POST   /api/v1/users             --> ...
[GIN-debug] GET    /api/v1/users/:id         --> ...
[GIN-debug] PUT    /api/v1/users/:id         --> ...
```

或者直接访问健康检查端点：

```bash
curl http://localhost:8080/health
```

### Q3: 能否同时支持多个领域？

**A**: 可以！只需在 `domains.go` 中添加多个导入：

```go
import (
    _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
    _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/order"
    _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/product"
)
```

### Q4: 路由注册的顺序重要吗？

**A**: 不重要。Gin 框架会自动处理路由冲突，按照最匹配的原则执行。

## 最佳实践

1. **始终使用 `go run ./cmd/api/`** - 确保包含所有包文件
2. **在 domains.go 中按字母顺序导入** - 保持代码整洁
3. **每个领域独立的 routes.go** - 职责清晰，易于维护
4. **使用 `initialized` 标志** - 防止重复注册
5. **避免在 Register() 中立即执行** - 统一在 Build() 时注册

## 相关文档

- [开发规范指南](development-guidelines.md) - 编码标准和代码结构
- [CLI 工具指南](cli-tool-guide.md) - 项目初始化和代码生成
- [架构设计文档](../architecture/architecture-design.md) - 整体技术架构
