# 数据库表结构总览

本文档提供 go-ddd-scaffold 项目所有数据库表的快速参考。

## 表列表

| 序号 | 表名 | 说明 | 软删除 |
|------|------|------|--------|
| 1 | `users` | 用户表 | ✅ |
| 2 | `tenants` | 租户表 | ✅ |
| 3 | `tenant_members` | 租户成员关联表 | ✅ |
| 4 | `roles` | 角色表 | ✅ |
| 5 | `permissions` | 权限定义表 | ❌ |
| 6 | `role_permissions` | 角色权限关联表 | ❌ |
| 7 | `audit_logs` | 审计日志表 | ❌ |
| 8 | `login_logs` | 登录日志表 | ❌ |
| 9 | `domain_events` | 领域事件表 | ❌ |
| 10 | `tenant_configs` | 租户配置表 | ❌ |

## 核心字段说明

### 1. users (用户表)

```sql
主键：id (BIGINT, Snowflake)
唯一索引：username, email
状态字段：status (0-待激活，1-激活，2-禁用，3-锁定)
软删除：deleted_at
```

**核心字段**：
- `username`: 用户名（唯一）
- `email`: 邮箱（唯一）
- `password_hash`: 密码哈希
- `display_name`: 显示名称
- `status`: 用户状态
- `failed_attempts`: 失败尝试次数
- `locked_until`: 锁定截止时间

### 2. tenants (租户表)

```sql
主键：id (BIGINT, Snowflake)
唯一索引：code
外键：owner_id → users(id)
状态字段：status (0-待激活，1-激活，2-停用，3-暂停)
软删除：deleted_at
```

**订阅相关字段**：
- `subscription_plan`: free/basic/pro/enterprise
- `subscription_status`: active/canceled/expired/trial
- `trial_ends_at`: 试用截止
- `subscription_expires_at`: 订阅过期
- `billing_cycle`: monthly/yearly

**资源限制**：
- `max_members`: 最大成员数
- `storage_limit`: 存储限制
- `api_call_limit`: API 调用限制

### 3. tenant_members (租户成员表)

```sql
复合主键：(tenant_id, user_id)
外键：tenant_id → tenants(id), user_id → users(id), role_id → roles(id)
状态字段：status (0-待邀请，1-活跃，2-已移除，3-已拒绝)
软删除：deleted_at
JSONB: metadata (部门、职位等扩展信息)
```

### 4. roles (角色表)

```sql
主键：id (BIGINT, Snowflake)
索引：(tenant_id, code) 唯一
外键：tenant_id → NULL (系统角色)
约束：CHECK (is_system = TRUE AND tenant_id IS NULL) OR (is_system = FALSE AND tenant_id IS NOT NULL)
软删除：deleted_at
JSONB: permissions (权限列表)
```

**设计要点**：
- `tenant_id IS NULL` + `is_system = TRUE`: 系统预定义角色
- `tenant_id = 123` + `is_system = FALSE`: 租户自定义角色

### 5. permissions (权限定义表)

```sql
主键：id (BIGINT, Snowflake)
唯一索引：code
作用域：scope (system/tenant)
JSONB: metadata
```

**权限编码示例**：
- `user:create`: 创建用户
- `user:read`: 查看用户
- `billing:manage`: 管理账单

### 6. role_permissions (角色权限关联表)

```sql
复合主键：(role_id, permission_id)
外键：role_id → roles(id), permission_id → permissions(id)
```

### 7. audit_logs (审计日志表)

```sql
主键：id (BIGINT, Snowflake)
外键：tenant_id → tenants(id), user_id → users(id)
IP 类型：ip_address (INET)
JSONB: metadata
```

**审计字段**：
- `action`: 操作类型
- `resource_type`: 资源类型
- `request_id`: 链路追踪 ID
- `status`: 成功/失败

### 8. login_logs (登录日志表)

```sql
主键：id (BIGINT, Snowflake)
外键：user_id → users(id), tenant_id → tenants(id)
IP 类型：ip_address (INET)
状态：login_status (success/failed/blocked)
```

**安全字段**：
- `login_type`: password/sso/oauth/magic_link
- `device_type`: desktop/mobile/tablet
- `is_suspicious`: 是否可疑
- `risk_score`: 风险评分 0-100

### 9. domain_events (领域事件表)

```sql
主键：id (BIGSERIAL 自增)
索引：(aggregate_id, aggregate_type), event_type, occurred_on
JSONB: event_data, metadata
状态：processed (是否已处理)
```

**事件溯源**：
- `aggregate_id`: 聚合根 ID
- `aggregate_type`: User/Tenant/Order
- `event_type`: UserRegistered/UserLoggedIn
- `event_version`: 事件版本

### 10. tenant_configs (租户配置表)

```sql
主键：tenant_id → tenants(id) ON DELETE CASCADE
JSONB: feature_flags, branding, custom_settings
```

**配置内容**：
- `feature_flags`: 功能开关
- `branding`: 品牌配置（logo、主题色）
- `custom_settings`: 自定义配置

## 索引策略

### 通用索引模式

1. **软删除部分索引**：
   ```sql
   CREATE INDEX idx_table_field ON table(field) WHERE deleted_at IS NULL;
   ```

2. **JSONB GIN 索引**：
   ```sql
   CREATE INDEX idx_table_jsonb ON table USING GIN (jsonb_field);
   ```

3. **复合索引**：
   ```sql
   CREATE INDEX idx_table_field1_field2 ON table(field1, field2 DESC);
   ```

### 各表核心索引

| 表名 | 核心索引 |
|------|----------|
| users | `idx_users_username`, `idx_users_email`, `idx_users_status_deleted_at` |
| tenants | `idx_tenants_code`, `idx_tenants_subscription_status` |
| tenant_members | `idx_tenant_members_tenant_id_status`, `idx_tenant_members_metadata` (GIN) |
| roles | `idx_roles_tenant_id_code`, `idx_roles_permissions` (GIN) |
| permissions | `idx_permissions_resource_action` |
| audit_logs | `idx_audit_logs_occurred_at DESC`, `idx_audit_logs_metadata` (GIN) |
| login_logs | `idx_login_logs_user_id_occurred_at DESC`, `idx_login_logs_suspicious` (部分索引) |
| domain_events | `idx_domain_events_aggregate`, `idx_domain_events_processed` (部分索引) |

## 约束检查

### CHECK 约束汇总

```sql
-- users 表
chk_users_status: status IN (0, 1, 2, 3)
chk_users_gender: gender IN (0, 1, 2, 3)

-- tenants 表
chk_tenants_status: status IN (0, 1, 2, 3)
chk_tenants_billing_cycle: billing_cycle IN ('monthly', 'yearly')
chk_tenants_max_members: max_members > 0

-- tenant_members 表
chk_tenant_members_status: status IN (0, 1, 2, 3)

-- roles 表
chk_roles_tenant_system: 
  (is_system = TRUE AND tenant_id IS NULL) OR 
  (is_system = FALSE AND tenant_id IS NOT NULL)

-- permissions 表
chk_permissions_scope: scope IN ('system', 'tenant')
```

## 触发器

所有包含 `updated_at` 字段的表都有自动更新触发器：

```sql
CREATE TRIGGER update_{table}_updated_at 
    BEFORE UPDATE ON {table} 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

适用表：
- users
- tenants
- tenant_members
- roles
- permissions
- tenant_configs

## 外键关系图

```
┌─────────────┐
│    users    │
└──────┬──────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
┌─────────────┐    ┌─────────────┐
│   tenants   │    │ tenant_members │
│ (owner_id)  │◀───┤ (user_id)   │
└──────┬──────┘    └──────┬──────┘
       │                  │
       │                  ├──────────┐
       │                  │          │
       ▼                  ▼          ▼
┌─────────────┐    ┌─────────────┐
│tenant_configs│   │    roles    │
│ (tenant_id) │   │ (tenant_id) │
└─────────────┘    └──────┬──────┘
                          │
                          ▼
                   ┌─────────────┐
                   │permissions  │
                   │             │
                   └──────▲──────┘
                          │
                   ┌──────┴──────┐
                   │role_permissions│
                   └─────────────┘
```

## 数据字典

### 状态码枚举

#### users.status
- `0`: 待激活 (pending)
- `1`: 激活 (active)
- `2`: 禁用 (inactive)
- `3`: 锁定 (locked)

#### tenants.status
- `0`: 待激活 (pending)
- `1`: 激活 (active)
- `2`: 停用 (inactive)
- `3`: 暂停 (suspended)

#### tenants.subscription_status
- `active`: 有效订阅
- `canceled`: 已取消
- `expired`: 已过期
- `trial`: 试用中

#### tenant_members.status
- `0`: 待邀请 (pending)
- `1`: 活跃 (active)
- `2`: 已移除 (removed)
- `3`: 已拒绝 (rejected)

#### permissions.scope
- `system`: 系统级权限
- `tenant`: 租户级权限

## PostgreSQL 特性应用

1. **JSONB 类型**：灵活的半结构化数据存储
2. **GIN 索引**：加速 JSONB 查询
3. **部分索引**：只为未删除的数据建立索引
4. **CHECK 约束**：保证数据一致性
5. **INET 类型**：原生支持 IP 地址
6. **触发器**：自动维护 updated_at
7. **CASCADE 删除**：自动清理关联数据

## 迁移文件

对应迁移文件位于 `../migrations/` 目录：

```bash
# 执行所有迁移
migrate -database "postgres://localhost:5432/go_ddd_scaffold_dev?sslmode=disable" -path ./migrations up

# 运行测试脚本
./migrations/test_migration.sh
```

## 相关文档

- [数据库设计详细文档](./reference/database-design.md)
- [迁移文件使用说明](../migrations/README.md)
- [领域模型文档](./architecture/domain-model.md)
