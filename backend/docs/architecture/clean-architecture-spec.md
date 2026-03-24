# Clean Architecture & Ports 模式规范

## 🎯 设计理念

### 核心原则
1. **依赖规则（Dependency Rule）** - 依赖只能指向内层，不能反向依赖
2. **Ports & Adapters** - Application 层定义 Ports，Infrastructure 层实现适配器
3. **职责分离** - 每层都有明确的单一职责
4. **Go 风格** - 简洁实用，避免过度抽象
5. **易于测试** - 各层独立可测，支持 Mock

### 架构分层与依赖方向

```
Domain (最内层 - 业务核心)
    ↑
Application (用例编排 + Ports 定义)
    ↑  (通过 Ports 接口)
Infrastructure (适配器实现)
    ↑
Interfaces (协议适配)
    ↑
Bootstrap (组合根 - 组装所有依赖)
```

### 关键特性

✅ **Application 层不依赖 Infrastructure**  
✅ **所有外部依赖通过 Port 接口抽象**  
✅ **Infrastructure 通过适配器实现 Ports**  
✅ **Module 层作为组合根负责组装**

---

## 📁 完整目录结构

```
backend/
├── cmd/                           # 【应用入口】程序启动入口
│   ├── api/                       # API 服务（HTTP）
│   │   └── main.go
│   ├── worker/                    # Worker 服务（异步任务）
│   │   └── main.go
│   └── cli/                       # CLI 工具（数据迁移、生成器等）
│       └── main.go
│
├── internal/                      # 【内部代码】私有实现，不对外暴露
│   │
│   ├── domain/                    # 【领域层】业务核心（最纯净）
│   │   ├── shared/                # 共享领域模型
│   │   │   ├── kernel/            # 核心概念
│   │   │   │   ├── entity.go      # 实体基类
│   │   │   │   ├── valueobject.go # 值对象基类
│   │   │   │   ├── aggregate.go   # 聚合根基类
│   │   │   │   ├── repository.go  # 仓储接口基类
│   │   │   │   ├── event.go       # 领域事件接口
│   │   │   │   └── error.go       # 领域错误定义
│   │   │   └── event/             # 共享事件
│   │   │
│   │   ├── user/                  # 用户上下文（限界上下文）
│   │   │   ├── aggregate/         # 聚合根
│   │   │   │   ├── user.go        # User 聚合根
│   │   │   │   └── user_test.go   # 领域模型测试
│   │   │   ├── valueobject/       # 值对象
│   │   │   │   ├── username.go    # 用户名值对象
│   │   │   │   ├── email.go       # 邮箱值对象
│   │   │   │   └── password.go    # 密码值对象
│   │   │   ├── event/             # 领域事件
│   │   │   │   ├── user_registered.go  # 用户注册事件
│   │   │   │   └── handlers.go    # 领域事件处理器
│   │   │   ├── service/           # 领域服务
│   │   │   │   └── password_policy.go  # 密码策略
│   │   │   └── repository/        # 仓储接口（由 Domain 定义）
│   │   │       └── user_repository.go
│   │   │
│   │   └── tenant/                # 租户上下文
│   │       ├── aggregate/
│   │       ├── valueobject/
│   │       ├── event/
│   │       ├── service/
│   │       └── repository/
│   │
│   ├── application/               # 【应用层】业务流程编排 + ⭐Ports 定义
│   │   ├── ports/                 # ⭐ Ports 接口（由 App 层定义）
│   │   │   ├── auth/              # 认证相关 Ports
│   │   │   │   └── token_service.go   # TokenService 端口
│   │   │   ├── idgen/             # ID 生成器 Ports
│   │   │   │   └── generator.go
│   │   │   ├── cache/             # 缓存 Ports
│   │   │   │   └── user_cache.go
│   │   │   └── email/             # 邮件 Ports
│   │   │       └── service.go
│   │   │
│   │   ├── user/                  # 用户应用服务
│   │   │   ├── service.go         # 应用服务（编排流程）
│   │   │   └── dto.go             # 数据传输对象
│   │   │
│   │   ├── auth/                  # 认证应用服务
│   │   │   ├── service.go         # ✅ 使用 Ports，不使用 Infrastructure
│   │   │   └── dto.go
│   │   │
│   │   └── unit_of_work.go        # 工作单元接口
│   │
│   ├── infrastructure/            # 【基础设施层】Adapters 实现
│   │   │
│   │   ├── config/                # 配置管理
│   │   │   ├── loader.go
│   │   │   └── model.go
│   │   │
│   │   ├── persistence/           # 数据持久化
│   │   │   ├── dao/               # GORM 生成的 DAO
│   │   │   └── repository/        # Repository 实现（Ports 的适配器）
│   │   │       └── user_repository.go
│   │   │
│   │   ├── platform/         # 平台服务实现
│   │   │   ├── auth/              # 认证基础设施
│   │   │   │   ├── jwt_service.go     # JWT 具体实现
│   │   │   │   ├── token.go           # Token 类型定义
│   │   │   │   └── token_service_adapter.go  # ⭐ 适配器
│   │   │   ├── snowflake/       # Snowflake ID 生成器
│   │   │   │   └── node.go          # ✅ 直接实现 Generator Port
│   │   │   └── email/           # 邮件服务实现
│   │   │       └── smtp_service.go
│   │   │
│   │   ├── cache/                 # 缓存实现
│   │   │   └── redis/
│   │   │       └── user_cache.go
│   │   │
│   │   ├── eventstore/            # 事件存储
│   │   ├── messaging/             # 消息队列
│   │   ├── logging/               # 日志服务
│   │   └── taskqueue/             # 任务队列（Asynq）
│   │
│   ├── interfaces/                # 【接口层】协议适配
│   │   ├── http/                  # HTTP 接口
│   │   │   ├── auth/
│   │   │   │   ├── handler.go
│   │   │   │   └── routes.go
│   │   │   └── user/
│   │   │       ├── handler.go
│   │   │       └── routes.go
│   │   ├── grpc/                  # gRPC 接口
│   │   └── messaging/             # 消息接口
│   │
│   ├── module/                    # 【组合根】模块组装
│   │   ├── auth.go                # ✅ 创建适配器并注入
│   │   └── user.go
│   │
│   └── bootstrap/                 # 【启动引导】应用初始化
│       ├── infra.go               # Infra 容器（组合所有基础设施）
│       └── module.go              # Module 注册与管理
│
├── pkg/                           # 【公共库】可复用的工具包
│   ├── response/
│   ├── util/
│   └── useragent/
│
└── configs/                       # 【配置文件】
    ├── config.yaml
    └── .env
```

---

## 🔑 关键设计决策

### 1. Ports 的位置：为什么在 Application 层？

**传统做法：** Ports 放在 Domain 层  
**本项目做法：** Ports 放在 Application 层

**理由：**
- Repository Port 是应用层的需求，不是领域层的核心概念
- 避免 Domain 层被技术细节污染
- 符合"谁需要谁定义"的原则

### 2. 适配器的作用

**TokenServiceAdapter 示例：**

```go
// Application 层定义的 Port
type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
}

// Infrastructure 层的具体实现
type JWTService struct { ... }

// 适配器：将 JWTService 转换为 TokenService Port
type TokenServiceAdapter struct {
    service TokenService  // JWTService
}
```

**为什么需要适配器？**
1. **类型隔离**：Port 层和 Infra 层的类型不同
2. **解耦**：Application 不知道 Infrastructure 的存在
3. **可替换**：可以轻松更换 JWT 库而不影响业务逻辑

### 3. Module 层作为组合根

**职责：**
- 知道所有层的具体实现
- 负责创建适配器
- 组装依赖并注入到 Application Service

**示例：**
```go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 注入到应用服务
    authSvc := authApp.NewAuthService(
        tokenServiceAdapter,  // ← 使用适配器
        ...
    )
}
```

---

## ✅ 验证清单

### 代码检查

- [x] Application 层不 import `internal/infrastructure` 包
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

## 📚 相关文档

- [Ports 模式架构设计](./ports-pattern-design.md) - Ports 模式的详细说明
- [DDD 设计指南](./ddd-design-guide.md) - DDD 核心概念
- [Bootstrap 模块架构](./bootstrap-module-architecture.md) - 组合根模式
- [事件驱动架构](./event-driven-architecture.md) - 事件系统设计

---

## 🔄 架构演进

**当前状态：** 完整的 Ports 模式 + 组合根

**最近重构：**
- ✅ auth/service.go 完全使用 Ports 接口（方案 B）
- ✅ 移除对 Infrastructure 的直接依赖
- ✅ 通过 TokenServiceAdapter 进行类型转换
- ✅ 完全符合 Clean Architecture 规范
