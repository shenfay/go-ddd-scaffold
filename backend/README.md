# Go DDD Scaffold

企业级 DDD 单体应用脚手架，基于领域驱动设计和 CQRS 架构模式构建的标准化企业应用开发平台。

## 项目概述

这是一个面向企业级应用场景的完整 DDD+CQRS 架构单体应用模板，提供：

- 标准化的领域驱动设计实践
- CQRS 读写分离架构
- 自动路由注册机制
- 完整的事件驱动机制
- 企业级安全和合规特性
- 高性能和可扩展性设计

## 核心架构特色

### DDD分层架构
```
┌─────────────────────────────────┐
│         Presentation Layer       │  (HTTP/gRPC/API)
├─────────────────────────────────┤
│        Application Layer         │  (Use Cases/Services)
├─────────────────────────────────┤
│          Domain Layer            │  (Business Logic)
├─────────────────────────────────┤
│      Infrastructure Layer        │  (Persistence/External)
└─────────────────────────────────┘
```

### CQRS模式实现
- **命令侧 (Write Model)**: 处理复杂的业务逻辑和数据一致性
- **查询侧 (Read Model)**: 优化读取性能，支持复杂的数据展示需求

### 事件驱动架构
- 完整的领域事件机制
- 事件存储和回放能力
- 异步事件处理支持

## 目录结构

```
go-ddd-scaffold/
├── cmd/                    # 应用程序入口
│   ├── api/               # REST API 服务（包含 main.go 和 domains.go）
│   ├── worker/            # 后台工作任务
│   └── cli/               # 命令行工具
├── internal/              # 内部包
│   ├── domain/            # 领域层
│   │   ├── user/          # 用户领域
│   │   ├── tenant/        # 租户领域
│   │   ├── order/         # 订单领域
│   │   ├── product/       # 产品领域
│   │   └── common/        # 通用领域概念
│   ├── application/       # 应用层
│   │   ├── commands/      # 命令处理
│   │   ├── queries/       # 查询处理
│   │   ├── services/      # 应用服务
│   │   └── dtos/          # 数据传输对象
│   ├── interfaces/        # 接口层
│   │   ├── http/          # HTTP 接口（路由自动注册）
│   │   ├── grpc/          # gRPC 接口
│   │   └── messaging/     # 消息接口
│   └── infrastructure/    # 基础设施层
│       ├── persistence/   # 数据持久化
│       ├── messaging/     # 消息传递
│       ├── eventstore/    # 事件存储
│       ├── cache/         # 缓存实现
│       └── config/        # 配置管理
├── configs/               # 配置文件
├── migrations/            # 数据库迁移
├── docs/                  # 文档
├── deployments/           # 部署配置
├── tools/                 # 开发工具
├── shared/                # 共享库
└── go.mod                 # Go 模块定义
```

## 快速开始

### 开发环境要求
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### 使用 Makefile（推荐）
```bash
# 查看可用命令
make help

# 安装依赖
make install-deps

# 启动开发服务器
make run

# 构建应用
make build

# 健康检查
make health
```

### 手动启动
```bash
# 安装依赖
go mod tidy

# 启动应用（重要：使用 ./cmd/api/ 而非 ./cmd/api/main.go，以包含 domains.go）
go run ./cmd/api/

# 或者先编译再运行
go build -o bin/api ./cmd/api
./bin/api
```

**注意**：必须使用 `go run ./cmd/api/` 而不是 `go run ./cmd/api/main.go`，因为 domains.go 包含了领域路由的自动注册导入。

## 企业级特性

### 安全合规
- 多因素认证支持
- 完整的审计日志
- RBAC权限控制系统
- 数据加密和隐私保护

### 性能优化
- 多级缓存架构
- 读写分离设计
- 查询优化器
- 连接池管理

### 可观测性
- 分布式追踪
- 结构化日志
- 性能监控
- 健康检查

## 开发指南

详细的技术文档请参考：
- [架构设计文档](docs/architecture/)
- [DDD模式指南](docs/ddd-patterns/)
- [CQRS实现手册](docs/cqrs-implementation/)

## 许可证

Apache 2.0 License