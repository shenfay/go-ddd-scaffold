# Go DDD Scaffold 文档中心

欢迎使用 Go DDD Scaffold！本文档中心提供完整的技术文档，帮助你快速上手和深入理解项目架构。

## 🎯 快速开始

### 新手入门路径

1. **第 1 天** - 阅读 [快速开始指南](guides/quickstart.md)，5 分钟上手
2. **第 2-3 天** - 学习 [架构设计](design/architecture-overview.md)，理解核心概念
3. **第 4-5 天** - 遵循 [开发规范](specifications/development-spec.md)，开始编码实践
4. **持续学习** - 参考 [实战教程](tutorials/getting-started-tutorial.md) 提升技能

---

## 📚 文档分类

### 📜 规范文档 (Specifications)

**必须遵守的开发规范和标准**

| 文档 | 说明 | 必读 |
|------|------|------|
| [开发规范](specifications/development-spec.md) | 编码标准、命名规范、代码结构 | ✅ |
| [架构规范](specifications/architecture-spec.md) | Clean Architecture 分层规则、依赖管理 | ✅ |
| [API 设计规范](specifications/api-spec.md) | RESTful API 设计、请求响应格式 | ✅ |
| [数据库规范](specifications/database-spec.md) | 表设计、索引策略、迁移流程 | ✅ |
| [错误处理规范](specifications/error-handling-spec.md) | 错误码体系、异常处理最佳实践 | ✅ |
| [安全规范](specifications/security-spec.md) | 认证授权、数据加密、审计日志 | ✅ |

### 📘 使用指南 (Guides)

**操作指南和最佳实践**

| 文档 | 说明 | 类型 |
|------|------|------|
| [快速开始](guides/quickstart.md) | 环境搭建、Hello World | 🔥 |
| [CLI 工具指南](guides/cli-tool-guide.md) | 代码生成器、项目初始化 | 🛠️ |
| [Module 开发指南](guides/module-development-guide.md) | 如何开发新功能模块 | ✅ |
| [Repository 指南](guides/repository-guide.md) | Repository 模式、事务处理 | ✅ |
| [DAO 使用指南](guides/dao-guide.md) | GORM Gen 配置、CRUD 操作 | 🛠️ |
| [DTO 使用指南](guides/dto-guide.md) | DTO 设计、数据转换 | ✅ |
| [路由配置指南](guides/routing-guide.md) | 路由自动加载机制 | ✅ |
| [工具包指南](guides/util-packages-guide.md) | pkg/util 工具函数 | 🛠️ |
| [AsynqMon 指南](guides/asynqmon-guide.md) | 任务队列监控 | 🛠️ |

### 🎨 设计文档 (Design)

**架构设计和核心概念深度解析**

| 文档 | 说明 | 图表 |
|------|------|------|
| [架构总览](design/architecture-overview.md) | 整体架构、技术选型、核心特性 | 📊 |
| [Ports 模式详解](design/ports-pattern-design.md) | Ports & Adapters 完整说明 | 📊 |
| [Clean Architecture](design/clean-architecture-spec.md) | 分层规范、依赖规则 | 📊 |
| [领域模型设计](design/domain-model.md) | 聚合根、值对象、领域服务 | 📊 |
| [DDD 设计指南](design/ddd-design-guide.md) | DDD 核心概念、限界上下文 | 📊 |
| [组合根模式](design/bootstrap-module-architecture.md) | Bootstrap + Module 架构 | 📊 |
| [事件驱动架构](design/event-driven-architecture.md) | EventBus vs EventPublisher | 📊 |
| [业务流程图](design/business-flow-diagrams.md) | 注册、登录等时序图 | 📊 |
| [技术实现流程](design/implementation-flow.md) | 启动流程、依赖注入 | 📊 |

### 📖 参考文档 (Reference)

**查询手册和技术规格**

| 文档 | 说明 | 类型 |
|------|------|------|
| [API 文档](reference/api-reference.md) | Swagger 自动生成 | 🌐 |
| [数据库 Schema](reference/database-schema.md) | 表结构、字段说明、索引 | 📋 |
| [配置项参考](reference/configuration-reference.md) | 所有配置项说明、默认值 | 📋 |
| [领域事件目录](reference/domain-events-catalog.md) | 所有领域事件列表 | 📋 |
| [错误码字典](reference/error-code-dictionary.md) | 完整错误码列表 | 📋 |

### 🎓 教程文档 (Tutorials)

**循序渐进的学习教程**

| 文档 | 说明 | 难度 |
|------|------|------|
| [入门教程](tutorials/getting-started-tutorial.md) | 从零开始构建第一个功能 | ⭐ |
| [实战案例](tutorials/practical-examples.md) | 完整业务场景实现 | ⭐⭐ |
| [最佳实践集](tutorials/best-practices.md) | 常见问题解决方案 | ⭐⭐ |

### ⚙️ 运维文档 (Operations)

**部署、监控和故障排查**

| 文档 | 说明 | 类型 |
|------|------|------|
| [部署指南](operations/deployment-guide.md) | Docker、K8s、云原生部署 | 🚀 |
| [监控告警](operations/monitoring-alerting.md) | Prometheus、Grafana 配置 | 📈 |
| [性能优化](operations/performance-tuning.md) | 数据库优化、缓存策略 | ⚡ |
| [故障排查](operations/troubleshooting.md) | 常见问题诊断流程 | 🔧 |
| [备份恢复](operations/backup-recovery.md) | 数据库备份、灾难恢复 | 💾 |

### 📝 变更日志 (Changelog)

**版本记录和迁移指南**

| 文档 | 说明 |
|------|------|
| [更新日志](changelog/CHANGELOG.md) | 每个版本的变更说明 |
| [迁移指南](changelog/migration-guide.md) | 版本升级的破坏性变更 |

---

## 🎯 按角色查看文档

### 👶 新手开发者

```
1. 快速开始 → guides/quickstart.md
2. 开发规范 → specifications/development-spec.md
3. 入门教程 → tutorials/getting-started-tutorial.md
4. CLI 工具 → guides/cli-tool-guide.md
```

### 💻 应用开发者

```
1. 架构总览 → design/architecture-overview.md
2. Module 开发 → guides/module-development-guide.md
3. Repository 指南 → guides/repository-guide.md
4. 领域模型 → design/domain-model.md
```

### 🏗️ 架构师

```
1. Ports 模式 → design/ports-pattern-design.md
2. Clean Architecture → design/clean-architecture-spec.md
3. DDD 设计 → design/ddd-design-guide.md
4. 事件驱动 → design/event-driven-architecture.md
```

### 🔧 运维工程师

```
1. 部署指南 → operations/deployment-guide.md
2. 监控告警 → operations/monitoring-alerting.md
3. 故障排查 → operations/troubleshooting.md
4. 性能优化 → operations/performance-tuning.md
```

---

## 📊 文档统计

| 类别 | 文档数量 | 完成度 |
|------|---------|--------|
| 规范文档 | 6 篇 | 🟢 100% |
| 使用指南 | 9 篇 | 🟢 100% |
| 设计文档 | 9 篇 | 🟢 100% |
| 参考文档 | 5 篇 | 🟢 100% |
| 教程文档 | 3 篇 | 🟡 70% |
| 运维文档 | 5 篇 | 🟡 60% |
| 变更日志 | 2 篇 | 🔴 30% |
| **总计** | **39 篇** | **🟢 85%** |

---

## 🔧 开发环境要求

### 必需工具

- **Go** 1.21+
- **PostgreSQL** 14+
- **Redis** 7+
- **Git**

### 推荐工具

- **VS Code** + Go 插件
- **Postman** / Insomnia（API 调试）
- **TablePlus** / DBeaver（数据库管理）
- **Draw.io** / Mermaid（绘制流程图）

---

## 📞 支持与反馈

### 获取帮助

1. **查阅文档** - 使用上方索引查找相关文档
2. **检查示例** - 参考 `examples/` 目录下的示例代码
3. **查看 Issue** - GitHub Issues 中的常见问题
4. **创建 Issue** - 提交问题或建议

### 贡献文档

欢迎贡献文档！请遵循以下步骤：

1. Fork 项目
2. 创建分支 (`git checkout -b feature/docs-improvement`)
3. 提交更改 (`git commit -m 'Add some docs'`)
4. 推送到分支 (`git push origin feature/docs-improvement`)
5. 创建 Pull Request

---

## 📜 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](../../LICENSE) 文件

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team  
**版本：** v2.0.0
