# Go DDD Scaffold

企业级 DDD 单体应用脚手架，基于领域驱动设计和事件驱动架构构建的标准化企业应用开发平台。

## 项目概述

这是一个面向企业级应用场景的完整 DDD 架构单体应用模板，提供：

- ✅ 标准化的领域驱动设计实践
- ✅ 自动路由注册机制
- ✅ 完整的领域事件机制
- ✅ 企业级安全和合规特性
- ✅ 高性能和可扩展性设计

## 核心架构特色

### DDD + Composition Root 分层架构

```
┌──────────────────────────────────────────────┐
│         Presentation Layer                    │
│  (HTTP Handlers / gRPC / Middleware)          │
├──────────────────────────────────────────────┤
│         Application Layer                     │
│  (Use Cases / Services / DTOs)                │
├──────────────────────────────────────────────┤
│           Domain Layer                        │
│  (Entities / Value Objects / Aggregates)      │
├──────────────────────────────────────────────┤
│       Infrastructure Layer                    │
│  (Persistence / External Services / Cache)    │
└──────────────────────────────────────────────┘
            ↑
    Bootstrap (Composition Root)
    - 创建所有依赖
    - 组装依赖关系
    - 类型安全，零运行时错误
```

### Composition Root 模式

项目采用 **Composition Root** 模式进行依赖管理：

```mermaid
graph TB
    A[main.go] --> B[Bootstrap]
    B --> C[Container]
    C --> D[Infrastructure]
    B --> E[Domain Initialization]
    E --> F[User Domain]
    E --> G[Tenant Domain]
    B --> H[Application Services]
    H --> I[User Service]
    H --> J[Auth Service]
    
    style B fill:#e1f5ff
    style C fill:#fff4e1
    style E fill:#f0f0f0
```

- ✅ **Bootstrap** 负责创建和组装所有领域组件
- ✅ **Container** 仅管理基础设施和 HTTP 路由
- ✅ **类型安全**：消除运行时类型断言，编译期检查
- ✅ **职责清晰**：各组件职责边界明确

**详细说明**: [Container 重构说明](docs/architecture/container-refactoring.md) | [开发规范](docs/guides/development-guidelines.md)

### 事件驱动架构

```mermaid
sequenceDiagram
    participant App as Application Service
    participant Repo as Repository
    participant DB as Database
    participant ES as EventStore
    
    App->>Repo: Save(Aggregate)
    Repo->>DB: Persist Aggregate
    Repo->>ES: Store Domain Events
    ES-->>Repo: Events Saved
    Repo-->>App: Success
    
    Note over ES: Events trigger side effects<br/>(Notifications, Projections, etc.)
```

- ✅ 完整的领域事件机制
- ✅ 事件存储和回放能力  
- ✅ 异步事件处理（用于触发副作用）

## 主要流程

### 1. 用户注册流程

```mermaid
sequenceDiagram
    participant Client as Client
    participant API as HTTP Handler
    participant App as Application Service
    participant User as User Aggregate
    participant Repo as UserRepository
    participant ES as EventStore
    
    Client->>API: POST /users/register
    API->>App: CreateUserRequest (DTO)
    App->>User: New(username, email, password)
    User-->>App: User Created + Domain Events
    App->>Repo: Save(User)
    Repo->>ES: Store Domain Events
    Repo->>DB: INSERT users
    ES-->>Repo: Events Stored
    Repo-->>App: Success
    App-->>API: UserDTO
    API-->>Client: 201 Created
    
    Note right of App: 应用服务协调领域对象<br/>保存聚合根和领域事件
```

### 2. Repository-DAO 数据流转

```mermaid
graph LR
    A[Domain Entity] -->|fromDomain| B[DAO Model]
    B -->|WithContext.Create| C[(Database)]
    C -->|WithContext.Find| D[DAO Model]
    D -->|toDomain| E[Domain Entity]
    
    subgraph Repository Layer
        A
        E
    end
    
    subgraph DAO Layer
        B
        D
    end
    
    style A fill:#e8f5e9
    style E fill:#e8f5e9
    style B fill:#fff3e0
    style D fill:#fff3e0
```

**详细实现**: [Repository-DAO 使用指南](docs/guides/repository-dao-usage.md)

### 3. 智能指针转换流程

```mermaid
graph TD
    A[Domain Field] --> B{Field Type}
    B -->|String NOT NULL| C[util.String]
    B -->|String NULLABLE| D[util.StringPtrNilIfEmpty]
    B -->|Int NOT NULL| E[util.Int32]
    B -->|Int NULLABLE| F[util.Int32PtrNilIfZero]
    
    C --> G[*string]
    D --> G
    E --> H[*int32]
    F --> H
    
    G --> I[DAO Model]
    H --> I
    
    style C fill:#bbdefb
    style D fill:#e3f2fd
    style E fill:#bbdefb
    style F fill:#e3f2fd
```

**技术细节**: [util/cast.go 方法命名规范](docs/guides/repository-dao-usage.md#类型转换工具)

## 目录结构

```
go-ddd-scaffold/
├── cmd/                    # 应用程序入口
│   ├── api/               # REST API 服务（main.go + domains.go）
│   ├── worker/            # 后台工作任务
│   └── cli/               # 命令行工具（代码生成/迁移管理）
├── internal/              # 内部包
│   ├── domain/            # 领域层（核心业务逻辑）
│   │   ├── user/          # 用户领域
│   │   ├── tenant/        # 租户领域
│   │   └── common/        # 通用领域概念
│   ├── bootstrap/         # Composition Root（依赖组装）
│   │   ├── bootstrap.go   # 应用启动器
│   │   └── user_domain.go # 用户领域初始化
│   ├── container/         # 基础设施容器
│   │   └── container.go   # Container 实现（HTTP 路由）
│   ├── application/       # 应用层（用例编排）
│   │   ├── auth/          # 认证服务
│   │   ├── user/          # 用户应用服务
│   │   └── shared/        # 共享应用组件
│   ├── interfaces/        # 接口层（适配器）
│   │   ├── http/          # HTTP 接口（路由自动注册）
│   │   ├── grpc/          # gRPC 接口
│   │   └── messaging/     # 消息接口
│   └── infrastructure/    # 基础设施层
│       ├── persistence/   # 数据持久化（DAO + Repository）
│       ├── messaging/     # 消息传递
│       ├── eventstore/    # 事件存储
│       ├── cache/         # 缓存实现
│       └── config/        # 配置管理
├── configs/               # 配置文件（.env + .yaml）
├── migrations/            # 数据库迁移脚本
├── docs/                  # 文档
│   ├── architecture/      # 架构设计文档
│   ├── guides/            # 开发指南
│   ├── implementation/    # 实现细节
│   └── reference/         # 技术参考
├── deployments/           # 部署配置（Docker + K8s）
├── tools/                 # 开发工具（CLI/Generator）
├── shared/                # 共享库（DDD 基础/CQRS/Response）
└── go.mod                 # Go 模块定义
```

## 快速开始

### 开发环境要求
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 安装依赖（GORM/gen、swag 等）
make install-deps

# 启动开发服务器（热重载）
make run

# 构建应用
make build

# 运行测试
make test

# 健康检查
make health
```

### 手动启动

```bash
# 1. 安装依赖
go mod tidy

# 2. 启动应用
# 重要：使用 ./cmd/api/ 而非 ./cmd/api/main.go，以包含 domains.go
go run ./cmd/api/

# 或者先编译再运行
go build -o bin/api ./cmd/api
./bin/api
```

**注意**：必须使用 `go run ./cmd/api/` 而不是 `go run ./cmd/api/main.go`，因为 `domains.go` 包含了领域路由的自动注册导入。

### 验证启动成功

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy"
  }
}
```

## 企业级特性

### 🔐 安全合规
- ✅ 多因素认证支持（MFA）
- ✅ 完整的审计日志（所有关键操作可追溯）
- ✅ RBAC 权限控制系统（角色 - 权限矩阵）
- ✅ 数据加密和隐私保护（传输中 + 静态）

### ⚡ 性能优化
- ✅ 多级缓存架构（Redis + 内存）
- ✅ 读写分离设计（主从数据库）
- ✅ 查询优化器（执行计划分析）
- ✅ 连接池管理（高效资源复用）

### 📊 可观测性
- ✅ 分布式追踪（OpenTelemetry 兼容）
- ✅ 结构化日志（JSON 格式，ELK 集成）
- ✅ 性能监控（Prometheus + Grafana）
- ✅ 健康检查（实时状态反馈）

## 设计亮点

### 1. 智能指针转换系统

项目创新性地设计了四类指针转换函数，完美解决了 Go 语言中值类型与指针类型的转换难题：

```go
// ✅ 场景 1：可 NULL 字段 - 区分空值和零值
DisplayName: util.StringPtrNilIfEmpty(displayName)
// "" → nil (数据库存储为 NULL)
// "张三" → *string("张三") (数据库存储为 '张三')

// ✅ 场景 2：NOT NULL 字段 - 简洁直接
LoginType: util.String(log.LoginType)
// 直接创建指针，不检查空值

// ✅ 场景 3：读取数据 - 安全解引用
domainName := util.StringValue(daoModel.DisplayName)
// nil → "", *string("张三") → "张三"
```

**优势**：
- 🎯 **类型安全**：编译期检查，无运行时反射
- 🚀 **代码简洁**：从 ~50 行减少到 ~20 行（-60%）
- 📦 **高度复用**：全项目统一使用 `pkg/util`
- 🧪 **易于测试**：纯函数，无副作用

**详细文档**: [Repository-DAO 使用指南](docs/guides/repository-dao-usage.md#类型转换工具)

### 2. Composition Root 依赖管理

采用 **Composition Root** 模式，实现了清晰的依赖管理：

```
main.go
  └── Bootstrap (组合根)
      ├── Container (基础设施容器)
      │   ├── GORM DB
      │   ├── Redis Cache
      │   └── HTTP Router
      ├── Domain Initialization (领域初始化)
      │   ├── User Domain Services
      │   └── Tenant Domain Services
      └── Application Services (应用服务)
          ├── User Service
          └── AuthService
```

**核心优势**：
- 🔒 **类型安全**：零运行时类型断言
- 📋 **职责清晰**：Bootstrap 组装，Container 路由，Domain 业务
- 🧩 **易于替换**：依赖注入，方便 Mock 和测试
- ⚙️ **编译期检查**：所有依赖在编译期验证

### 3. Repository + DAO 分层设计

创新的 Repository + DAO 双层设计，兼顾了 DDD 的纯净性和工程实践：

```
┌─────────────────┐
│  Domain Entity  │ ← 领域对象（业务逻辑）
├─────────────────┤
│   Repository    │ ← 领域转换 + 业务规则
├─────────────────┤
│     DAO Model   │ ← 数据库模型（GORM/gen）
├─────────────────┤
│     Database    │ ← PostgreSQL
└─────────────────┘
```

**实现细节**：
- Repository 负责：领域对象转换、业务规则验证、领域事件保存
- DAO 负责：类型安全的 CRUD、基础查询封装
- 两者通过 `fromDomain()` / `toDomain()` 方法互相转换

**完整示例**: [Repository-DAO 使用指南](docs/guides/repository-dao-usage.md)

### 4. 事件驱动架构

基于领域事件的松耦合架构，支持最终一致性和异步处理：

```
用户注册
  ↓
[UserCreated] 事件发布
  ├─→ 发送欢迎邮件（Notification）
  ├─→ 更新读模型（Projection）
  ├─→ 触发风控检查（Policy）
  └─→ 同步第三方系统（Integration）
```

**特性**：
- 📦 事件存储在 `domain_events` 表，支持事件溯源
- 🔄 自动事件持久化，Repository 自动保存未提交事件
- ⚡ 异步处理器，后台 worker 消费事件
- 🔍 完整的事件历史，便于审计和问题排查

---

## 开发指南

详细的技术文档请参考：
- 📐 [架构设计文档](docs/architecture/) - 整体架构和设计原则
- 📖 [DDD 模式指南](docs/ddd-patterns/) - DDD 模式详解
- 🛠️ [开发规范](docs/guides/development-guidelines.md) - 编码规范和最佳实践
- 🗄️ [Repository-DAO 使用指南](docs/guides/repository-dao-usage.md) - Repository 层实现细节
- 🔄 [类型转换和时间处理](docs/guides/util-packages-guide.md) - util 包使用手册

---

## 许可证

Apache 2.0 License