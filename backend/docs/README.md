# Backend 文档索引

本目录包含 go-ddd-scaffold 项目后端技术文档，按类别组织。

## 文档目录

### 开发指南 (guides/)

| 文档 | 说明 |
|------|------|
| [development-guidelines.md](guides/development-guidelines.md) | 开发规范、编码标准、代码结构 |
| [cli-tool-guide.md](guides/cli-tool-guide.md) | CLI 工具安装配置、项目初始化、代码生成 |
| [error-handling.md](guides/error-handling.md) | 统一错误处理机制、错误码体系 |
| [route-auto-loading.md](guides/route-auto-loading.md) | 路由自动加载机制和使用方法 |

### 架构设计 (architecture/)

| 文档 | 说明 |
|------|------|
| [architecture-design.md](architecture/architecture-design.md) | 整体技术架构、分层设计、核心设计原则 |
| [ddd-cqrs-design-guide.md](architecture/ddd-cqrs-design-guide.md) | DDD 实践指南、CQRS 模式应用 |
| [domain-model.md](architecture/domain-model.md) | 领域模型设计、实体关系、聚合设计 |

### 参考文档 (reference/)

| 文档 | 说明 |
|------|------|
| [api-specification.md](reference/api-specification.md) | API 接口规范、请求响应格式 |
| [database-design.md](reference/database-design.md) | 数据库设计、表结构、索引策略 |
| [security-compliance.md](reference/security-compliance.md) | 安全设计、等保合规、权限模型 |

### 运维文档 (operations/)

| 文档 | 说明 |
|------|------|
| [deployment-operations.md](operations/deployment-operations.md) | 部署方案、运维配置、监控告警 |

## 阅读顺序建议

### 新手入门
1. 先阅读项目主文档 `/docs/project-overview.md` 了解项目整体
2. 阅读 [development-guidelines.md](guides/development-guidelines.md) 熟悉开发规范
3. 阅读 [architecture-design.md](architecture/architecture-design.md) 了解技术架构

### 业务开发
1. 参考 [domain-model.md](architecture/domain-model.md) 了解领域模型
2. 参考 [ddd-cqrs-design-guide.md](architecture/ddd-cqrs-design-guide.md) 掌握设计模式
3. 查阅 [api-specification.md](reference/api-specification.md) 了解接口规范
4. 参考 [error-handling.md](guides/error-handling.md) 处理错误

### 部署运维
1. 参考 [database-design.md](reference/database-design.md) 了解数据模型
2. 阅读 [deployment-operations.md](operations/deployment-operations.md) 了解部署方案
3. 参考 [security-compliance.md](reference/security-compliance.md) 确保安全合规
