---
name: tenant-builder
description: 多租户 SaaS 架构快速搭建工具。提供租户管理、家庭角色、订阅计划、数据隔离等完整功能。基于 Casbin RBAC 权限控制和 Goose 数据库迁移。适用于教育科技 SaaS 平台。
version: "1.0.0"
author: MathFun Team
tags: [tenant, multi-tenancy, saas, rbac, casbin, subscription, family, data-isolation]
---

# Tenant Builder - 多租户 SaaS 架构搭建工具

## 功能概述

这是一个智能化的多租户 SaaS 架构生成工具，专为 MathFun 项目设计。它提供完整的租户管理体系、家庭角色系统、订阅管理和数据隔离机制，基于 **Casbin RBAC** 权限控制和 **Goose** 数据库迁移。

## 核心能力

### 1. 租户管理系统
- **租户生命周期** - 创建、激活、暂停、注销全周期管理
- **子域名识别** - 基于子域名的租户自动识别和路由
- **套餐计划** - free/basic/premium 多级套餐配置
- **过期管理** - 自动检测和处理过期租户

### 2. 家庭角色系统
- **家庭组管理** - 家长创建和管理家庭组
- **角色定义** - parent/child/educator 三种预定义角色
- **邀请机制** - 基于邮件或邀请码的家庭成员邀请
- **权限继承** - 家长对子女的监督和管理权限

### 3. Casbin RBAC 集成
- **领域隔离** - 租户级别的资源访问控制
- **角色权限** - 细粒度的操作权限定义
- **策略管理** - 动态添加和修改权限策略
- **审计日志** - 完整的权限使用记录

### 4. 数据隔离机制
- **行级隔离** - 基于 tenant_id 的数据过滤
- **中间件集成** - 自动注入租户上下文
- **查询拦截** - GORM 插件自动添加租户条件
- **跨租户禁止** - 严格防止数据越界访问

### 5. 订阅管理
- **套餐切换** - 支持升级降级操作
- **用量统计** - 用户数、存储空间等指标监控
- **账单生成** - 按月/年自动生成账单
- **支付集成** - 预留微信支付接口

## 使用场景

### 适用情况
- 教育科技 SaaS 平台搭建
- 需要多租户隔离的 B2B2C 应用
- 家庭学习管理系统
- 学校班级管理系统
- 订阅制知识服务平台

### 不适用情况
- 单租户内部系统
- 无需数据隔离的 C 端应用
- 非订阅制的纯免费产品

## 基本使用

### 快速开始（10 分钟）

```bash
# 1. 创建租户相关表
tenant-builder create --module tenant --description "create tenants table"
tenant-builder create --module family --description "create families table"
tenant-builder create --module subscription --description "create subscriptions table"

# 2. 执行迁移
tenant-builder migrate up --env dev

# 3. 生成 DAO 代码
tenant-builder generate \
  --tables tenants,families,family_members,subscriptions \
  --output ./internal/infrastructure/persistence/gorm

# 4. 初始化 Casbin 策略
tenant-builder init-casbin --default-role admin
```

### 完整工作流

```bash
# 一键完成所有步骤
tenant-builder full-workflow \
  --modules tenant,family,subscription \
  --env dev \
  --generate-dao \
  --init-casbin \
  --with-middleware
```

## 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--module` | string | 是 | - | 模块名 (tenant/family/subscription) |
| `--modules` | array | 否 | - | 多个模块列表 |
| `--description` | string | 是 | - | 迁移描述 |
| `--env` | string | 否 | dev | 环境 (dev/staging/prod) |
| `--output` | string | 否 | ./generated | 代码输出目录 |
| `--generate-dao` | flag | 否 | false | 生成 DAO 代码 |
| `--init-casbin` | flag | 否 | false | 初始化 Casbin 策略 |
| `--with-middleware` | flag | 否 | false | 生成租户中间件 |
| `--with-tests` | flag | 否 | false | 生成测试文件 |

## 生成的代码结构

```
backend/internal/
├── domain/
│   ├── tenant/
│   │   ├── entity/           # 租户实体
│   │   ├── repository/       # 仓储接口
│   │   └── service/          # 领域服务
│   └── family/
│       ├── entity/           # 家庭实体
│       └── repository/       # 家庭仓储
├── infrastructure/
│   ├── persistence/
│   │   └── gorm/
│   │       ├── model/        # GORM Model
│   │       └── dao/          # DAO 实现
│   └── auth/
│   │       └── casbin/       # Casbin 策略
│   └── middleware/
│   │       └── tenant.go     # 租户中间件
└── interfaces/
    └── http/
        └── tenant/           # HTTP 处理器
```

## 最佳实践

### 租户隔离策略

#### 数据库层面
- ✅ 所有表添加 `tenant_id` 字段（系统表除外）
- ✅ 创建复合索引：`(tenant_id, id)`
- ✅ 使用外键约束确保引用完整性

#### 应用层面
- ✅ 中间件自动提取和注入 tenant_id
- ✅ Repository 层统一添加租户过滤条件
- ✅ 单元测试验证隔离有效性

### Casbin 策略设计

```conf
# 租户管理员策略
p, tenant_admin, tenants, read, allow
p, tenant_admin, tenants, update, allow
p, tenant_admin, users, manage, allow

# 普通家长策略
p, parent, children, read, allow
p, parent, learning_progress, read, allow
p, parent, subscriptions, view, allow

# 孩子策略
p, child, learning_resources, read, allow
p, child, games, play, allow
```

### 订阅管理

```go
type SubscriptionPlan string

const (
    PlanFree     SubscriptionPlan = "free"
    PlanBasic    SubscriptionPlan = "basic"
    PlanPremium  SubscriptionPlan = "premium"
)

// 套餐限制
var PlanLimits = map[SubscriptionPlan]PlanLimit{
    PlanFree: {
        MaxUsers: 3,
        MaxStorage: 100 * MB,
        Features: []string{"basic_games"},
    },
    PlanBasic: {
        MaxUsers: 10,
        MaxStorage: 1 * GB,
        Features: []string{"basic_games", "progress_tracking"},
    },
}
```

## 故障排除

### 常见问题

**租户数据未隔离**
- 检查中间件是否正确注册
- 验证 tenant_id 是否自动注入
- 查看 Repository 是否添加租户过滤

**Casbin 策略不生效**
- 确认策略已加载到适配器
- 检查 enforcer 配置
- 验证角色继承关系

**订阅过期未处理**
- 检查定时任务是否运行
- 验证过期检测逻辑
- 查看通知发送状态

### 获取帮助
- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- 🚀 快速开始：查看 [QUICKSTART.md](./QUICKSTART.md)

## 版本历史

- v1.0.0 (2026-02-25): 初始版本发布
  - 完整的租户管理功能
  - 家庭角色系统
  - Casbin RBAC 集成
  - 订阅管理框架
  - 数据隔离中间件

---
*本技能遵循 Qoder Skills 规范，专为 MathFun 项目优化设计*
