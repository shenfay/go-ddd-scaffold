# Backend 文档索引

本文档索引包含 go-ddd-scaffold 项目后端所有技术文档。

## 🎯 快速开始

**新手入门推荐阅读顺序：**

1. [架构总览](architecture/architecture-overview.md) - 快速了解整体架构
2. [Ports 模式设计](architecture/ports-pattern-design.md) - 理解核心架构模式
3. [开发规范](guides/development-guidelines.md) - 熟悉编码标准
4. [错误处理](guides/error-handling.md) - 掌握统一错误处理

---

## 📚 文档目录

### 🏗️ 架构设计 (architecture/) - 13 篇

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [architecture-overview.md](architecture/architecture-overview.md) | 🆕 技术架构总览、Ports 模式、完整目录结构 | 🔥 必读 |
| [ports-pattern-design.md](architecture/ports-pattern-design.md) | 🆕 Ports & Adapters 模式详细说明 | 🔥 必读 |
| [clean-architecture-spec.md](architecture/clean-architecture-spec.md) | 🆕 Clean Architecture 分层规范 | 🔥 必读 |
| [architecture-diagrams-detailed.md](architecture/architecture-diagrams-detailed.md) | 🆕 详细架构分层图、依赖关系图 | 📊 图表 |
| [domain-model-visual.md](architecture/domain-model-visual.md) | 🆕 领域模型可视化、聚合根结构图 | 📊 图表 |
| [business-flow-diagrams.md](architecture/business-flow-diagrams.md) | 🆕 核心业务流程时序图、流程图 | 📊 图表 |
| [implementation-flow-diagrams.md](architecture/implementation-flow-diagrams.md) | 🆕 技术实现流程图、启动流程 | 📊 图表 |
| [ddd-design-guide.md](architecture/ddd-design-guide.md) | DDD 实践指南、聚合根、值对象 | ✅ 核心 |
| [domain-model.md](architecture/domain-model.md) | 领域模型设计、实体关系 | ✅ 核心 |
| [bootstrap-module-architecture.md](architecture/bootstrap-module-architecture.md) | 组合根模式、Module 架构 | ✅ 核心 |
| [event-driven-architecture.md](architecture/event-driven-architecture.md) | 事件驱动架构设计 | ✅ 核心 |
| [architecture-diagrams-section.md](architecture/architecture-diagrams-section.md) | Clean Architecture 分层图集 | 📊 图表 |
| [domain-model-diagrams.md](architecture/domain-model-diagrams.md) | 领域模型结构图集 | 📊 图表 |
| [ddd-process-diagrams.md](architecture/ddd-process-diagrams.md) | 业务流程图（注册/登录） | 📊 图表 |

### 📖 开发指南 (guides/) - 12 篇

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [development-guidelines.md](guides/development-guidelines.md) | 开发规范、编码标准 | 🔥 必读 |
| [error-handling.md](guides/error-handling.md) | 统一错误处理机制 | ✅ 核心 |
| [route-auto-loading.md](guides/route-auto-loading.md) | 路由自动加载机制 | ✅ 核心 |
| [cli-tool-guide.md](guides/cli-tool-guide.md) | CLI 工具使用指南 | 🛠️ 工具 |
| [dao-generator.md](guides/dao-generator.md) | DAO 代码生成器 | 🛠️ 工具 |
| [repository-dao-usage.md](guides/repository-dao-usage.md) | Repository 和 DAO 使用 | ✅ 核心 |
| [dto-guidelines.md](guides/dto-guidelines.md) | DTO 设计规范 | ✅ 核心 |
| [module-development-guide.md](guides/module-development-guide.md) | Module 开发指南 | ✅ 核心 |
| [util-packages-guide.md](guides/util-packages-guide.md) | 工具包使用指南 | 🛠️ 工具 |
| [asynqmon-usage-guide.md](guides/asynqmon-usage-guide.md) | AsynqMon 监控使用 | 🛠️ 工具 |

### 💻 实现文档 (implementation/)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [core-features-summary.md](implementation/core-features-summary.md) | 核心功能实现总结 | ✅ 核心 |

### 📋 参考文档 (reference/)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [api-specification.md](reference/api-specification.md) | API 接口规范 | ✅ 核心 |
| [database-design.md](reference/database-design.md) | 数据库设计、表结构 | ✅ 核心 |
| [database-schema-overview.md](reference/database-schema-overview.md) | 数据库 Schema 总览 | ℹ️ 参考 |
| [security-compliance.md](reference/security-compliance.md) | 安全设计、等保合规 | ✅ 核心 |

### ⚙️ 运维文档 (operations/)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [deployment-operations.md](operations/deployment-operations.md) | 部署方案、运维配置 | ✅ 核心 |

### 🔄 重构文档 (refactoring/)

| 文档 | 说明 | 优先级 |
|------|------|--------|
| [simplicity-design-assessment.md](refactoring/simplicity-design-assessment.md) | 简洁设计评估 | ℹ️ 参考 |

### 📡 API 文档 (swagger/)

Swagger 自动生成的 API 文档，包含：
- `swagger.json` - OpenAPI 规范文件
- `swagger.yaml` - OpenAPI YAML 格式
- `docs.go` - Swagger 初始化代码

详见：[swagger/README.md](swagger/README.md)

---

## 🎓 学习路径

### 阶段 1：入门（1-2 天）

1. ✅ 阅读 [architecture-overview.md](architecture/architecture-overview.md)
2. ✅ 阅读 [development-guidelines.md](guides/development-guidelines.md)
3. ✅ 运行项目并测试 API

### 阶段 2：理解架构（2-3 天）

1. ✅ 深入理解 [ports-pattern-design.md](architecture/ports-pattern-design.md)
2. ✅ 学习 [ddd-design-guide.md](architecture/ddd-design-guide.md)
3. ✅ 研究 [domain-model.md](architecture/domain-model.md)
4. ✅ 查看 [architecture-diagrams-detailed.md](architecture/architecture-diagrams-detailed.md) 架构图

### 阶段 3：实战开发（持续）

1. ✅ 参考 [module-development-guide.md](guides/module-development-guide.md) 开发新功能
2. ✅ 遵循 [error-handling.md](guides/error-handling.md) 处理错误
3. ✅ 使用 [cli-tool-guide.md](guides/cli-tool-guide.md) 生成代码
4. ✅ 查看 [business-flow-diagrams.md](architecture/business-flow-diagrams.md) 理解业务流程

### 阶段 4：深入优化（进阶）

1. ✅ 研究 [event-driven-architecture.md](architecture/event-driven-architecture.md)
2. ✅ 优化数据库设计 [database-design.md](reference/database-design.md)
3. ✅ 实施监控和告警 [deployment-operations.md](operations/deployment-operations.md)

---

## 📊 图表索引

### 架构图表

- [架构分层详解](architecture/architecture-diagrams-detailed.md)
  - Clean Architecture 完整架构图
  - 依赖方向详解
  - Ports & Adapters 模式详解
  - Module 组装流程图
  - 数据流图

- [领域模型可视化](architecture/domain-model-visual.md)
  - User 聚合根完整模型图
  - Tenant 聚合根模型图
  - Role & Permission 模型图
  - 聚合根关系图
  - 值对象详细设计图
  - 生命周期状态机图
  - 数据库映射关系图

### 业务流程图表

- [核心业务流程图](architecture/business-flow-diagrams.md)
  - 用户注册流程时序图
  - 注册流程决策树
  - 用户登录流程时序图
  - 登录流程决策树
  - Token 刷新流程时序图
  - 用户登出流程图
  - 用户资料更新时序图
  - 密码修改流程图
  - 统一错误处理流程图
  - 领域事件处理流程图

### 技术实现图表

- [技术实现流程图](architecture/implementation-flow-diagrams.md)
  - 应用启动流程图
  - Module 注册时序图
  - 依赖注入流程图
  - Repository 适配器模式流程图
  - TokenService 适配器转换流程图
  - HTTP 请求处理链路图
  - UnitOfWork 事务管理流程图
  - 领域事件异步处理流程图

---

## 🎖️ 架构特性

本项目采用现代化的架构设计：

✅ **Clean Architecture** - 清晰的依赖规则  
✅ **Ports & Adapters** - 高度解耦的接口设计  
✅ **Domain-Driven Design** - 业务为核心的领域模型  
✅ **Composition Root** - 明确的依赖组装  
✅ **Constructor Injection** - 清晰的依赖注入  

详见：[architecture-overview.md](architecture/architecture-overview.md)

---

## 🔧 开发环境

**必需工具：**
- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Git

**推荐工具：**
- VS Code + Go 插件
- Postman / Insomnia（API 调试）
- TablePlus / DBeaver（数据库管理）
- Draw.io / Mermaid（绘制流程图）

---

## 📞 支持与反馈

如有问题或建议，请：
1. 查阅相关文档
2. 检查现有 Issue
3. 创建新的 Issue

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
