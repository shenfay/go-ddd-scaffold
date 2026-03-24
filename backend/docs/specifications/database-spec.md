# 数据库规范

本文档定义了 Go DDD Scaffold 项目的数据库设计和使用规范。

## 📋 设计原则

### 核心原则

1. **领域驱动** - 表结构反映领域模型
2. **规范化** - 遵循第三范式（3NF）
3. **性能优先** - 合理的索引和分区策略
4. **可维护性** - 清晰的命名和注释

---

## 🎯 命名规范

### 表命名

#### 基本规则

```sql
✅ 正确：
users              -- 使用复数形式
tenants            -- 使用复数形式
tenant_members     -- 使用下划线分隔
role_permissions   -- 关联表使用完整名称

❌ 错误：
user               -- 不应该用单数
Users              -- 不应该大写
tbl_users          -- 不需要前缀
```

#### 特殊表命名

```sql
-- 关联表（多对多）
user_roles         -- 用户角色关联
role_permissions   -- 角色权限关联
tenant_members     -- 租户成员关联

-- 审计表
audit_logs         -- 审计日志
login_logs         -- 登录日志
domain_events      -- 领域事件（Outbox）

-- 配置表
tenant_configs     -- 租户配置
system_configs     -- 系统配置
```

### 字段命名

#### 通用字段

```sql
-- 主键
id                 -- BIGINT PRIMARY KEY

-- 时间戳
created_at         -- TIMESTAMP NOT NULL DEFAULT NOW()
updated_at         -- TIMESTAMP NOT NULL DEFAULT NOW()
deleted_at         -- TIMESTAMP NULL (软删除)

-- 外键
user_id            -- 引用 users.id
tenant_id          -- 引用 tenants.id
role_id            -- 引用 roles.id
```

#### 业务字段

```sql
-- 用户相关
username           -- VARCHAR(50) NOT NULL UNIQUE
email              -- VARCHAR(255) NOT NULL UNIQUE
password_hash      -- VARCHAR(255) NOT NULL
status             -- VARCHAR(20) NOT NULL DEFAULT 'active'
display_name       -- VARCHAR(100)

-- 租户相关
name               -- VARCHAR(100) NOT NULL
slug               -- VARCHAR(50) NOT NULL UNIQUE
description        -- TEXT

-- 通用布尔值
is_active          -- BOOLEAN NOT NULL DEFAULT true
is_verified        -- BOOLEAN NOT NULL DEFAULT false
```

### 索引命名

```sql
-- 普通索引
idx_{table}_{column}
idx_users_email
idx_tenants_status

-- 唯一索引
uk_{table}_{column}
uk_users_username
uk_users_email

-- 复合索引
idx_{table}_{col1}_{col2}
idx_tenant_members_tenant_user

-- 外键索引
fk_{table}_{column}
fk_users_tenant_id
```

### 约束命名

```sql
-- 主键约束
pk_{table}
pk_users

-- 外键约束
fk_{table}_{reference}
fk_users_tenant

-- 检查约束
chk_{table}_{column}
chk_users_status

-- 默认值约束
df_{table}_{column}
df_users_created_at
```

---

## 🏗️ 表设计规范

### 基础表结构

#### 用户表

```sql
CREATE TABLE users (
    -- 主键
    id BIGINT PRIMARY KEY,
    
    -- 基本信息
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    
    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_verified BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    
    -- 约束
    CONSTRAINT uk_users_username UNIQUE (username),
    CONSTRAINT uk_users_email UNIQUE (email),
    CONSTRAINT chk_users_status CHECK (status IN ('active', 'inactive', 'locked', 'pending'))
);

-- 索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
```

#### 租户表

```sql
CREATE TABLE tenants (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    owner_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uk_tenants_slug UNIQUE (slug),
    CONSTRAINT fk_tenants_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    CONSTRAINT chk_tenants_status CHECK (status IN ('active', 'suspended', 'archived'))
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_owner_id ON tenants(owner_id);
```

#### 租户成员表（关联表）

```sql
CREATE TABLE tenant_members (
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (tenant_id, user_id),
    CONSTRAINT fk_tenant_members_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_tenant_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_tenant_members_role CHECK (role IN ('owner', 'admin', 'member', 'viewer'))
);

CREATE INDEX idx_tenant_members_user_id ON tenant_members(user_id);
```

#### 角色表

```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    tenant_id BIGINT NULL,  -- NULL 表示全局角色
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uk_roles_name UNIQUE (name),
    CONSTRAINT fk_roles_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
```

#### 权限表

```sql
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uk_permissions_name UNIQUE (name),
    CONSTRAINT uk_permissions_resource_action UNIQUE (resource, action)
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
```

#### 角色权限关联表

```sql
CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

#### 用户角色关联表

```sql
CREATE TABLE user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    tenant_id BIGINT NULL,  -- NULL 表示全局角色
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (user_id, role_id, tenant_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_tenant_id ON user_roles(tenant_id);
```

#### 审计日志表

```sql
CREATE TABLE audit_logs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(50) NOT NULL,
    resource_id BIGINT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

#### 领域事件表（Outbox Pattern）

```sql
CREATE TABLE domain_events (
    id BIGINT PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    metadata JSONB,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed BOOLEAN NOT NULL DEFAULT false,
    processed_at TIMESTAMP NULL,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT chk_domain_events_processed CHECK (
        (processed = true AND processed_at IS NOT NULL) OR
        (processed = false)
    )
);

CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_type, aggregate_id);
CREATE INDEX idx_domain_events_occurred_at ON domain_events(occurred_at);
CREATE INDEX idx_domain_events_unprocessed ON domain_events(occurred_at) WHERE processed = false;
```

---

## 🔍 索引策略

### 索引设计原则

1. **高频查询优先** - 为频繁查询的字段创建索引
2. **选择性原则** - 选择性高的字段更适合索引
3. **覆盖索引** - 尽可能包含查询所需的所有字段
4. **避免过度索引** - 平衡读写性能

### 索引类型选择

#### B-Tree 索引（默认）

```sql
-- 适用于等值查询、范围查询、排序
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
```

#### Hash 索引

```sql
-- 仅适用于等值查询
CREATE INDEX idx_sessions_token ON sessions USING HASH(token);
```

#### GIN 索引（JSONB）

```sql
-- 适用于 JSONB 字段查询
CREATE INDEX idx_audit_logs_old_values ON audit_logs USING GIN(old_values);
CREATE INDEX idx_audit_logs_new_values ON audit_logs USING GIN(new_values);
```

#### 部分索引

```sql
-- 仅索引满足条件的行
CREATE INDEX idx_users_active ON users(email) WHERE status = 'active';
CREATE INDEX idx_domain_events_unprocessed ON domain_events(occurred_at) WHERE processed = false;
```

#### 复合索引

```sql
-- 最左前缀原则
CREATE INDEX idx_tenant_members_tenant_user ON tenant_members(tenant_id, user_id);
-- 可用于：
-- WHERE tenant_id = ?
-- WHERE tenant_id = ? AND user_id = ?
-- ORDER BY tenant_id, user_id
```

---

## 🔄 数据库迁移

### 迁移文件命名

```bash
# 格式：{version}_{description}.up.sql / {version}_{description}.down.sql
000001_create_users_table.up.sql
000001_create_users_table.down.sql

000002_create_tenants_table.up.sql
000002_create_tenants_table.down.sql

000003_add_email_index_to_users.up.sql
000003_add_email_index_to_users.down.sql
```

### 迁移文件模板

```sql
-- 00000X_create_table_name.up.sql
-- 创建 {table_name} 表
-- 版本：{version}
-- 日期：{date}

BEGIN;

CREATE TABLE IF NOT EXISTS {table_name} (
    id BIGINT PRIMARY KEY,
    -- ... 其他字段
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_{table_name}_... ON {table_name}(...);

-- 注释
COMMENT ON TABLE {table_name} IS '{描述}';
COMMENT ON COLUMN {table_name}.id IS '主键 ID';
-- ... 其他注释

COMMIT;
```

```sql
-- 00000X_create_table_name.down.sql
-- 回滚：删除 {table_name} 表

BEGIN;

DROP TABLE IF EXISTS {table_name} CASCADE;

COMMIT;
```

### 迁移命令

```bash
# 安装 migrate 工具
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建新迁移
migrate create -ext sql -dir migrations -seq create_{table_name}_table

# 执行迁移
migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=disable" up

# 回滚迁移
migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=disable" down 1

# 查看迁移状态
migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=disable" version
```

---

## 💾 数据质量

### 约束使用

```sql
-- NOT NULL 约束
username VARCHAR(50) NOT NULL,

-- UNIQUE 约束
CONSTRAINT uk_users_email UNIQUE (email),

-- CHECK 约束
CONSTRAINT chk_users_status CHECK (status IN ('active', 'inactive', 'locked')),

-- FOREIGN KEY 约束
CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
```

### 触发器

```sql
-- 自动更新 updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### 数据验证

```sql
-- 邮箱格式验证
ALTER TABLE users ADD CONSTRAINT chk_users_email_format 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- 用户名格式验证
ALTER TABLE users ADD CONSTRAINT chk_users_username_format 
CHECK (username ~* '^[a-z][a-z0-9_]{2,49}$');
```

---

## 🚀 性能优化

### 查询优化

```sql
-- ✅ 正确：使用 EXPLAIN 分析
EXPLAIN ANALYZE
SELECT * FROM users
WHERE email = 'test@example.com';

-- ✅ 正确：避免 SELECT *
SELECT id, username, email FROM users
WHERE status = 'active';

-- ❌ 错误：函数导致索引失效
SELECT * FROM users WHERE LOWER(email) = 'test@example.com';

-- ✅ 正确：使用 CITEXT 或函数索引
CREATE INDEX idx_users_email_lower ON users(LOWER(email));
```

### 分区表

```sql
-- 按时间分区（审计日志）
CREATE TABLE audit_logs (
    id BIGINT,
    created_at TIMESTAMP NOT NULL,
    -- ... 其他字段
) PARTITION BY RANGE (created_at);

CREATE TABLE audit_logs_2024_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE audit_logs_2024_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

### 连接池配置

```yaml
# config.yaml
database:
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_lifetime: 1h   # 连接最大生命周期
  conn_max_idle_time: 5m  # 连接最大空闲时间
```

---

## 📊 数据库安全

### 权限管理

```sql
-- 创建应用用户
CREATE USER app_user WITH PASSWORD 'strong_password';

-- 授予最小权限
GRANT CONNECT ON DATABASE go_ddd_scaffold TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- 撤销危险权限
REVOKE CREATE ON SCHEMA public FROM app_user;
REVOKE ALL ON TABLE pg_user FROM app_user;
```

### 敏感数据加密

```sql
-- 使用 pgcrypto 扩展
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 加密存储
INSERT INTO users (password_hash) VALUES (crypt('password', gen_salt('bf')));

-- 验证密码
SELECT * FROM users WHERE password_hash = crypt('password', password_hash);
```

### 审计追踪

```sql
-- 记录所有 DDL 操作
CREATE EVENT TRIGGER audit_ddl_changes
ON ddl_command_end
EXECUTE FUNCTION log_ddl_changes();
```

---

## 📈 监控和维护

### 性能监控

```sql
-- 慢查询日志
ALTER SYSTEM SET log_min_duration_statement = 1000;  -- 1 秒

-- 查看慢查询
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### 日常维护

```sql
-- 定期 VACUUM
VACUUM ANALYZE users;

-- 重建索引
REINDEX TABLE users;

-- 更新统计信息
ANALYZE users;
```

### 备份恢复

```bash
# 备份数据库
pg_dump -U postgres go_ddd_scaffold > backup.sql

# 恢复数据库
psql -U postgres go_ddd_scaffold < backup.sql

# 仅备份结构
pg_dump -U postgres --schema-only go_ddd_scaffold > schema.sql

# 仅备份数据
pg_dump -U postgres --data-only go_ddd_scaffold > data.sql
```

---

## 📚 最佳实践

### 1. 表设计

✅ 推荐：
- 使用 BIGINT 作为主键
- 添加 created_at 和 updated_at
- 使用状态字段而不是多个布尔字段
- 适当使用 JSONB 存储灵活数据

❌ 避免：
- 使用 UUID 作为主键（除非必要）
- 使用业务字段作为主键
- 过度规范化导致查询复杂
- 在关系表中存储额外数据

### 2. 查询优化

✅ 推荐：
- 使用 EXPLAIN 分析查询
- 为常用查询创建覆盖索引
- 使用预编译语句
- 批量操作代替循环

❌ 避免：
- N+1 查询问题
- 在 WHERE 子句中使用函数
- LIKE '%pattern%' 前缀通配符
- 大表不加限制的 JOIN

### 3. 事务处理

✅ 推荐：
```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer func() {
    if p := recover(); p != nil {
        tx.Rollback()
        panic(p)
    } else if err != nil {
        tx.Rollback()
    }
}()

// 执行多个操作
err = tx.ExecContext(ctx, query1, args1...)
err = tx.ExecContext(ctx, query2, args2...)

err = tx.Commit()
```

❌ 避免：
- 长事务
- 事务中包含外部 API 调用
- 忘记提交或回滚
- 嵌套事务

---

## 📖 参考资源

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [SQL 反模式](https://pragprog.com/titles/bksqla/sql-antipatterns/)
- [数据库设计最佳实践](https://aws.amazon.com/cn/blogs/database/category/database/)
- [GORM 文档](https://gorm.io/docs/index.html)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
