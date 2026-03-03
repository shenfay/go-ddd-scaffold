---
name: ddd-scaffold
description: DDD Clean Architecture 项目脚手架生成器。一键生成标准的四层架构目录结构、领域实体模板、聚合根、Repository 接口与实现。适用于快速启动 Go 后端 DDD 项目。
version: "1.0.0"
author: MathFun Team
tags: [ddd, scaffold, architecture, golang, clean-architecture, code-generation]
---

# DDD Scaffold Skill

## 功能概述

这是一个智能化的 DDD 项目脚手架生成工具，基于 MathFun 项目的最佳实践设计。它能够一键生成符合 Clean Architecture 标准的完整项目结构，包括领域层、应用层、基础设施层和接口层。

## 核心能力

### 1. 智能架构生成
- **四层分层架构** - 自动生成 Domain/Application/Infrastructure/Interface 目录
- **标准目录结构** - 遵循 DDD 和 Clean Architecture 最佳实践
- **模块化组织** - 按领域划分模块，支持多领域协同
- **依赖注入配置** - 预配置 Google Wire 依赖注入

### 2. 领域建模支持
- **实体模板** - 生成包含状态和行为的领域实体
- **值对象模板** - 生成不可变的值对象
- **聚合根模板** - 生成聚合根及其内部结构
- **领域服务模板** - 生成跨实体的业务逻辑
- **领域事件模板** - 生成事件定义和发布机制

### 3. 基础设施集成
- **Repository 实现** - 基于 GORM 的仓储实现
- **事件驱动** - NATS 消息队列集成
- **任务调度** - Asynq 分布式任务队列
- **数据库迁移** - golang-migrate 集成

### 4. 快速启动支持
- **编译验证** - 生成的代码可直接编译运行
- **示例代码** - 包含完整的用户领域示例
- **配置文件** - 预配置 Viper、日志、监控等
- **Docker 支持** - Docker Compose 一键启动

## 使用场景

### 适用情况
- 从零开始构建 Go DDD 项目
- 现有项目重构为 DDD 架构
- 需要标准化项目结构
- 团队 DDD 实践推广

### 不适用情况
- 简单的 CRUD 应用
- 不需要领域建模的技术组件
- 已有成熟架构的遗留系统改造

## 基本使用

### 快速开始
```bash
# 生成标准 DDD 项目结构
/ddd-scaffold --project-name myapp --domains user

# 生成多领域项目
/ddd-scaffold --project-name ecommerce --domains user,order,product,inventory

# 交互式生成（推荐）
/ddd-scaffold --interactive
```

### 高级用法
```bash
# 指定输出目录
/ddd-scaffold --project-name myapp --output ./custom/path

# 选择生成内容
/ddd-scaffold --project-name myapp --with-examples --with-tests --with-docker

# 自定义架构风格
/ddd-scaffold --project-name myapp --style minimal  # minimal | standard | full
```

### 参数说明
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--project-name` | string | 是 | - | 项目名称 |
| `--domains` | array | 否 | [user] | 领域列表，逗号分隔 |
| `--output` | string | 否 | ./generated | 输出目录 |
| `--style` | string | 否 | standard | 架构风格 (minimal/standard/full) |
| `--with-examples` | flag | 否 | false | 包含示例代码 |
| `--with-tests` | flag | 否 | false | 包含测试文件 |
| `--with-docker` | flag | 否 | false | 包含 Docker 配置 |
| `--interactive` | flag | 否 | false | 交互模式 |

## 生成的项目结构

### 标准架构（standard）
```
myapp/
├── cmd/
│   ├── server/
│   │   └── main.go           # API 服务入口
│   └── worker/
│       └── main.go           # Worker 入口
├── internal/
│   ├── config/
│   │   ├── config.go         # 配置结构体
│   │   └── viper.go          # Viper 配置
│   ├── domain/
│   │   └── {domain}/
│   │       ├── entity/       # 领域实体
│   │       ├── valueobject/  # 值对象
│   │       ├── aggregate/    # 聚合根
│   │       ├── repository/   # 仓储接口
│   │       ├── service/      # 领域服务
│   │       └── event/        # 领域事件
│   ├── application/
│   │   └── {domain}/
│   │       ├── service/      # 应用服务
│   │       ├── dto/          # DTO
│   │       └── event/        # 应用事件处理器
│   ├── infrastructure/
│   │   ├── wire/             # 依赖注入
│   │   ├── persistence/      # 持久化实现
│   │   ├── cache/            # 缓存实现
│   │   ├── queue/            # 消息队列
│   │   └── web/              # Web 框架
│   └── interfaces/
│       └── http/
│           └── {domain}/     # HTTP Handler
├── pkg/
│   ├── errors/               # 错误定义
│   ├── response/             # 统一响应
│   └── validator/            # 验证器
├── migrations/
│   └── sql/                  # SQL 迁移文件
├── configs/
│   └── config.yaml           # 配置文件
├── scripts/
├── tests/
├── go.mod
├── Makefile
└── README.md
```

## 配置说明

Skill 使用 `.qoder/skills/ddd-scaffold/config.yaml` 进行配置：

```yaml
# 基础配置
scaffold:
  project_name: "myapp"
  output_dir: "./generated"
  style: "standard"  # minimal | standard | full
  
# 领域配置
domains:
  - name: "user"
    enabled: true
    entities:
      - name: "User"
        fields:
          - name: "ID"
            type: "string"
          - name: "Name"
            type: "string"
          - name: "Email"
            type: "string"
    aggregates:
      - name: "UserAggregate"
        root_entity: "User"
        
# 生成选项
generation:
  create_examples: true
  create_tests: true
  create_docker: true
  create_readme: true
  
# 集成配置
integrations:
  database:
    type: "postgresql"  # postgresql | mysql | mongodb
    orm: "gorm"
    
  messaging:
    enabled: true
    type: "nats"  # nats | rabbitmq | kafka
    
  task_queue:
    enabled: true
    type: "asynq"  # asynq | cron
```

## 最佳实践

### 1. 领域划分原则
- **单一职责** - 每个领域有明确的业务边界
- **高内聚低耦合** - 领域内高度内聚，领域间松耦合
- **依赖倒置** - 依赖抽象不依赖具体实现
- **领域事件解耦** - 通过事件实现领域间最终一致性

### 2. 命名规范
- **领域名称** - 使用小写复数形式（user, order, product）
- **实体名称** - 使用大驼峰命名（User, OrderItem）
- **Repository 接口** - `{Entity}Repository`
- **Service** - `{Domain}Service`

### 3. 代码组织
- **领域层** - 只包含纯业务逻辑，无技术依赖
- **应用层** - 编排领域对象，不包含业务规则
- **基础设施层** - 实现技术细节，依赖抽象接口
- **接口层** - 处理协议转换，不包含业务逻辑

## 故障排除

### 常见问题

**生成的代码编译失败**
- 检查 Go 版本是否 >= 1.21
- 确认已运行 `go mod tidy`
- 验证导入路径是否正确

**领域依赖循环**
- 使用 `--analyze-dependencies` 检查循环依赖
- 通过领域事件解耦相互依赖的领域
- 引入共享内核（Shared Kernel）

**不知道如何划分领域**
- 使用 `--interactive` 模式获得引导
- 参考示例项目中的领域划分
- 咨询 DDD Architect Agent

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

## 版本历史

- v1.0.0 (2026-02-25): 初始版本
  - 基于 ddd-development-workflow 实践经验重构
  - 聚焦脚手架生成功能
  - 支持通用项目（排除教育行业特定内容）
  - 简化配置和使用方式

---
*本技能遵循 Qoder Skills 规范，专为快速启动 DDD 项目优化设计*
