# Backend 文档索引

本目录包含 go-ddd-scaffold 项目后端技术文档，按类别组织。

## 📚 文档目录

### 架构设计 (architecture/) - 6 篇

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture-design.md](architecture/architecture-design.md) | 整体技术架构、分层设计、核心设计原则 | ✅ 核心 |
| [ddd-design-guide.md](architecture/ddd-design-guide.md) | DDD 实践指南、聚合根、值对象、领域事件 | ✅ 核心 |
| [domain-model.md](architecture/domain-model.md) | 领域模型设计、实体关系、聚合设计 | ✅ 核心 |
| [architecture-diagrams-section.md](architecture/architecture-diagrams-section.md) | Clean Architecture 分层架构图集 | ✨ 图表 |
| [domain-model-diagrams.md](architecture/domain-model-diagrams.md) | 领域模型结构图集（聚合根/值对象/事件） | ✨ 图表 |
| [ddd-process-diagrams.md](architecture/ddd-process-diagrams.md) | 核心业务流程图集（注册/登录时序图） | ✨ 图表 |

### 开发指南 (guides/) - 9 篇

| 文档 | 说明 | 状态 |
|------|------|------|
| [development-guidelines.md](guides/development-guidelines.md) | 开发规范、编码标准、代码结构 | ✅ 核心 |
| [cli-tool-guide.md](guides/cli-tool-guide.md) | CLI 工具安装配置、项目初始化、代码生成 | ✅ 工具 |
| [error-handling.md](guides/error-handling.md) | 统一错误处理机制、错误码体系 | ✅ 核心 |
| [route-auto-loading.md](guides/route-auto-loading.md) | 路由自动加载机制和使用方法 | ✅ 核心 |
| [dao-config.md](guides/dao-config.md) | DAO 配置和使用指南 | ✅ 工具 |
| [dao-generator.md](guides/dao-generator.md) | DAO 代码生成器使用指南 | ✅ 工具 |
| [repository-dao-usage.md](guides/repository-dao-usage.md) | Repository 和 DAO 使用指南 | ✅ 核心 |
| [util-packages-guide.md](guides/util-packages-guide.md) | 工具包使用指南（cast/time） | ✅ 工具 |

### 实现文档 (implementation/) - 1 篇

| 文档 | 说明 | 状态 |
|------|------|------|
| [core-features-summary.md](implementation/core-features-summary.md) | 核心功能实现总结（注册/登录/获取信息） | ✨ 新增 |

### 参考文档 (reference/) - 4 篇

| 文档 | 说明 | 状态 |
|------|------|------|
| [api-specification.md](reference/api-specification.md) | API 接口规范、请求响应格式 | ✅ 核心 |
| [database-design.md](reference/database-design.md) | 数据库设计、表结构、索引策略 | ✅ 核心 |
| [database-schema-overview.md](reference/database-schema-overview.md) | 数据库 Schema 总览 | ℹ️ 参考 |
| [security-compliance.md](reference/security-compliance.md) | 安全设计、等保合规、权限模型 | ✅ 核心 |

### 运维文档 (operations/) - 1 篇

| 文档 | 说明 | 状态 |
|------|------|------|
| [deployment-operations.md](operations/deployment-operations.md) | 部署方案、运维配置、监控告警 | ✅ 核心 |

### API 文档 (swagger/) - 自动生成

| 文档 | 说明 |
|------|------|
| [swagger/README.md](swagger/README.md) | Swagger API 文档说明 |
| swagger.json | OpenAPI 规范文件 |
| swagger.yaml | OpenAPI YAML 格式 |
| docs.go | Swagger 初始化代码 |

## 阅读顺序建议

### 新手入门
1. 先阅读项目主文档 `/docs/project-overview.md` 了解项目整体
2. 阅读 [development-guidelines.md](guides/development-guidelines.md) 熟悉开发规范
3. 阅读 [architecture-design.md](architecture/architecture-design.md) 了解技术架构

### 业务开发
1. 参考 [domain-model.md](architecture/domain-model.md) 了解领域模型
| [architecture-diagrams-section.md](architecture/architecture-diagrams-section.md) | Clean Architecture 分层架构图集 | ✨ 图表 |
| [domain-model-diagrams.md](architecture/domain-model-diagrams.md) | 领域模型结构图集（聚合根/值对象/事件） | ✨ 图表 |
| [ddd-process-diagrams.md](architecture/ddd-process-diagrams.md) | 核心业务流程图集（注册/登录时序图） | ✨ 图表 |
2. 参考 [ddd-design-guide.md](architecture/ddd-design-guide.md) 掌握设计模式
3. 查阅 [api-specification.md](reference/api-specification.md) 了解接口规范
4. 参考 [error-handling.md](guides/error-handling.md) 处理错误

### 部署运维
1. 参考 [database-design.md](reference/database-design.md) 了解数据模型
2. 阅读 [deployment-operations.md](operations/deployment-operations.md) 了解部署方案
3. 参考 [security-compliance.md](reference/security-compliance.md) 确保安全合规
