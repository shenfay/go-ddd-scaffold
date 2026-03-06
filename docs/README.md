# 技术文档索引

## 📚 规范文档（Standards）

**必读** - 所有团队成员必须遵守的开发规范

- **[代码规范](standards/code-style.md)** - 命名、注释、错误处理、测试规范
- **[DDD 实现规范](standards/ddd-implementation.md)** - 分层架构职责、领域建模规范

---

## 🚀 快速开始（Getting Started）

- [安装指南](getting-started/installation.md) - 环境配置和项目启动
- [5 分钟快速体验](getting-started/quickstart.md) - 创建第一个 API
- [配置说明](getting-started/configuration.md) - 配置文件和环境变量

---

## 📖 开发指南（Guides）

### 基础教程
- [如何创建领域模块](guides/create-domain-module.md) - 从零开始创建完整领域
- [数据库迁移](guides/database-migration.md) - Goose 使用指南
- [添加 API 端点](guides/add-api-endpoint.md) - HTTP Handler 编写
- [实现业务逻辑](guides/implement-business-logic.md) - Domain Service 编写

### 进阶主题
- CQRS 模式实现 - 命令查询分离
- 领域事件驱动 - EventBus 使用
- 多租户架构 - Tenant-based SaaS
- 认证授权 - JWT + Casbin

---

## 🏗️ 架构文档（Architecture）

- [架构总览](architecture/overview.md) - 整体架构图和技术栈
- [分层架构](architecture/layers.md) - 四层架构详解
- [依赖关系](architecture/dependencies.md) - 包依赖图
- [架构决策记录 (ADR)](architecture/adr/) - 重要技术决策

---

## 📡 API 文档（API Reference）

- [认证 API](api-reference/authentication.md) - 登录、注册、Token
- [用户管理 API](api-reference/user-management.md) - CRUD 接口
- [租户管理 API](api-reference/tenant-management.md) - 多租户相关
- [Swagger UI](http://localhost:8080/swagger/index.html) - 在线文档

---

## 🚢 部署文档（Deployment）

- [本地开发](deployment/local-development.md) - Docker Compose 开发环境
- [Docker 部署](deployment/docker-deployment.md) - 容器化部署
- [Kubernetes 部署](deployment/kubernetes-deployment.md) - K8s 配置
- [生产环境配置](deployment/production-setup.md) - 性能优化和安全加固

---

## 🧪 测试文档（Testing）

- [单元测试](testing/unit-testing.md) - 领域层测试
- [集成测试](testing/integration-testing.md) - 应用层测试
- [E2E 测试](testing/e2e-testing.md) - 端到端测试
- [基准测试](testing/benchmark-testing.md) - 性能测试

---

## 🛠️ 工具文档（Tools）

- [代码生成器](tools/code-generator.md) - DAO/DTO 自动生成
- [项目脚手架](tools/scaffold.md) - 一键创建领域模块
- [Makefile 命令](tools/makefile-commands.md) - 常用操作命令

---

## 📋 检查清单（Checklists）

审核和验证标准：

- ✅ **[代码审查清单](checklists/code-review.md)** - PR 审核标准
- ⏳ [发布前检查](checklists/release-checklist.md) - 上线前验证
- ⏳ [安全检查清单](checklists/security-checklist.md) - 安全性验证

---

## 📝 版本历史

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| v1.0 | 2026-03-06 | 初始版本，包含基础规范和 DDD 实现规范 |

---

## 🔗 外部资源

- [Go 官方文档](https://golang.org/doc/)
- [GORM 文档](https://gorm.io/docs/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [DDD 参考](https://martinfowler.com/tags/domain_driven_design.html)

---

**维护人员**: Development Team  
**最后更新**: 2026-03-06
