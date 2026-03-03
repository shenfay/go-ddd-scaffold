---
name: db-migrator
description: 数据库迁移与 DAO 生成工具。基于 Goose 的数据库迁移管理、GORM Model 和 DAO 代码自动生成。支持多环境配置、字段备注自动添加。适用于标准化数据库变更管理和快速数据访问层开发。
version: "2.0.0"
author: MathFun Team
tags: [database, migration, dao, postgresql, gorm, goose, code-generation, multi-environment]
---

# DB Migrator - 数据库迁移与 DAO 生成工具

## 功能概述

这是一个智能化的数据库迁移管理和 DAO 代码生成工具，基于 MathFun 项目的最佳实践设计。它使用 **Goose** 作为数据库迁移工具，提供标准化的数据库模式变更管理和 GORM Model/DAO 代码自动生成。

## 核心能力

### 1. 数据库迁移管理（Goose）
- **自动时间戳** - 遵循 Goose 版本命名规范（YYYYMMDDHHMMSS）
- **SQL 模板库** - 提供 CREATE、ALTER、DROP、INDEX 等常用操作模板
- **字段备注** - 自动添加 COMMENT ON COLUMN 备注说明
- **回滚支持** - 每个迁移都包含对应的 DOWN 迁移
- **多环境支持** - dev/staging/prod 环境配置管理
- **迁移状态检查** - 自动追踪已执行和待执行的迁移

### 2. GORM Model 生成
- **智能类型映射** - PostgreSQL 类型 → Go 类型自动转换
- **标签自动生成** - json/gorm/validate 标签一键生成
- **关联关系支持** - HasOne/HasMany/ManyToMany 自动识别
- **软删除集成** - 自动添加 deleted_at 字段支持
- **索引定义** - 根据数据库索引自动生成 GORM 索引标签

### 3. DAO 代码生成
- **Repository 模式** - 标准的 CRUD 接口定义
- **实现分离** - 接口与实现分离（dao/interface + dao/impl）
- **事务支持** - 自动注入事务管理
- **错误处理** - 统一的错误转换和包装
- **分页查询** - 内置分页和排序方法
- **批量操作** - 支持批量插入、更新、删除

## 使用场景

### 适用情况
- 新项目数据库初始化
- 现有表结构变更
- 快速生成数据访问层代码
- 需要字段备注的数据库项目
- 多环境部署的数据库管理

### 不适用情况
- 简单的数据查询操作
- 临时性的数据库调试
- NoSQL 数据库项目

## 基本使用

### 快速开始（5 分钟）

```bash
# 1. 创建迁移文件
db-migrator create --table users --description "create users table"

# 2. 编辑生成的 SQL 文件（自动打开编辑器）
# 添加字段定义和备注

# 3. 执行迁移
db-migrator migrate up --env dev

# 4. 生成 GORM Model 和 DAO
db-migrator generate --tables users --output ./internal/infrastructure/persistence/gorm
```

### 完整工作流

```bash
# 一键完成所有步骤
db-migrator full-workflow \
  --table products \
  --description "create products table" \
  --env dev \
  --generate-dao \
  --with-tests
```

## 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--table` | string | 是 | - | 表名 |
| `--tables` | array | 否 | - | 多个表名列表 |
| `--description` | string | 是 | - | 迁移描述 |
| `--env` | string | 否 | dev | 环境 (dev/staging/prod) |
| `--output` | string | 否 | ./generated | 代码输出目录 |
| `--generate-dao` | flag | 否 | false | 生成 DAO 代码 |
| `--with-tests` | flag | 否 | false | 生成测试文件 |
| `--with-comments` | flag | 否 | true | 包含字段备注 |
| `--soft-delete` | flag | 否 | true | 启用软删除支持 |

## 配置说明

技能使用 `.qoder/skills/db-migrator/config.yaml` 进行配置：

```yaml
# 数据库配置
database:
  host: localhost
  port: 5432
  user: mathfun
  password: math111
  name: mathfun
  ssl_mode: disable

# 多环境配置
environments:
  dev:
    host: localhost
    port: 5432
    user: dev_user
    password: dev_pass
    name: mathfun_dev
    
  staging:
    host: staging-db.example.com
    port: 5432
    user: staging_user
    password: staging_pass
    name: mathfun_staging
    
  prod:
    host: prod-db.example.com
    port: 5432
    user: prod_user
    password: prod_pass
    name: mathfun_prod

# 标准配置
standards:
  naming:
    table_format: snake_case
    field_format: snake_case
  timestamps:
    created: created_at
    updated: updated_at
    deleted: deleted_at
  comments:
    required: true
    format: "COMMENT ON COLUMN {{table}}.{{column}} IS '{{comment}}'"

# 路径配置
paths:
  migrations: backend/migrations/sql
  dao_output: backend/internal/infrastructure/persistence/gorm
  model_output: backend/internal/infrastructure/persistence/gorm/model
```

## 最佳实践

### 命名规范
- ✅ 表名使用蛇形命名法（snake_case）：`user_profiles`
- ✅ 字段名使用蛇形命名法：`created_at`, `updated_at`
- ✅ 迁移文件包含清晰的描述：`20260225100000_create_users_table.sql`
- ❌ 避免使用复数形式：用 `user` 而非 `users`（除非确实表示集合）

### 迁移文件编写
- 每个迁移文件只处理一个逻辑变更
- 提供清晰的 UP 和 DOWN 迁移
- 必须添加字段备注（COMMENT ON COLUMN）
- 使用事务确保原子性
- 遵循 PostgreSQL 最佳实践

### GORM Model 设计
- 使用指针类型表示可空字段
- 正确设置 json 标签
- 添加 validate 验证规则
- 明确指定主键和外键关系

### 团队协作
- 在开始工作前检查迁移状态
- 及时提交和推送迁移文件
- 协调团队成员的迁移操作
- 定期验证生成代码的正确性

## 故障排除

### 常见问题

**数据库连接失败**
- 检查数据库服务是否运行
- 验证连接配置是否正确
- 确认网络连通性
- 检查防火墙设置

**迁移执行失败**
- 检查 SQL 语法是否正确
- 验证表结构依赖关系
- 查看详细的错误信息
- 确认 Goose 版本兼容性

**代码生成失败**
- 确认数据库表已正确创建
- 检查 GORM 配置
- 验证输出路径权限
- 确保 go.mod 包含所需依赖

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

## 版本历史

- v2.0.0 (2026-02-25): 重大更新
  - 确定使用 Goose 作为迁移工具
  - 新增多环境配置支持
  - 增强字段备注自动生成
  - 改进 GORM Model 类型映射
  - 新增 Repository 模式 DAO 生成
  
- v1.0.0 (2026-01-26): 初始版本
  - 基础迁移文件生成功能
  - DAO 代码自动生成
  - 完整工作流支持

---
*本技能遵循 Qoder Skills 规范，专为 MathFun 项目优化设计*
