# Go DDD Scaffold 数据库设计文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的数据库设计，包括表结构设计、索引策略、Snowflake ID实现规范以及数据库迁移管理方案。

## 数据库选型与配置

### 技术选型
- **数据库系统**：PostgreSQL 15+
- **迁移工具**：golang-migrate/migrate
- **连接池**：GORM内置连接池
- **主键策略**：Snowflake ID算法

### 连接配置
```yaml
database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: ${DB_NAME:scaffold_db}
  user: ${DB_USER:postgres}
  password: ${DB_PASSWORD:}
  ssl_mode: ${DB_SSL_MODE:disable}
  max_idle_conns: ${DB_MAX_IDLE_CONNS:10}
  max_open_conns: ${DB_MAX_OPEN_CONNS:100}
  conn_max_lifetime: ${DB_CONN_MAX_LIFETIME:1h}
```

## 核心表结构设计

### 1. 用户相关表

#### users 表（用户主表）
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID
    username VARCHAR(50) NOT NULL UNIQUE,     -- 用户名
    email VARCHAR(255) NOT NULL UNIQUE,       -- 邮箱
    password_hash VARCHAR(255) NOT NULL,      -- 密码哈希
    status SMALLINT DEFAULT 0 NOT NULL,       -- 状态：0-待激活，1-激活，2-禁用，3-锁定
    display_name VARCHAR(100),                -- 显示名称
    first_name VARCHAR(50),                   -- 名字
    last_name VARCHAR(50),                    -- 姓氏
    gender SMALLINT DEFAULT 0,                -- 性别：0-未知，1-男，2-女，3-其他
    phone_number VARCHAR(20),                 -- 电话号码
    avatar_url VARCHAR(500),                  -- 头像URL
    last_login_at TIMESTAMP,                  -- 最后登录时间
    login_count INTEGER DEFAULT 0,            -- 登录次数
    failed_attempts INTEGER DEFAULT 0,        -- 连续失败登录次数
    locked_until TIMESTAMP,                   -- 账户锁定截止时间
    version INTEGER DEFAULT 0,                -- 乐观锁版本号
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_login_at ON users(last_login_at);
```

#### users 表索引优化
```sql
-- 性能优化索引
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_login_at ON users(last_login_at);
```

### 2. 租户相关表

#### tenants 表（租户主表）
```sql
CREATE TABLE tenants (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID
    code VARCHAR(20) NOT NULL UNIQUE,         -- 租户编码
    name VARCHAR(100) NOT NULL,               -- 租户名称
    description TEXT,                         -- 租户描述
    status SMALLINT DEFAULT 0 NOT NULL,       -- 状态：0-活跃，1-停用，2-暂停
    owner_id BIGINT NOT NULL,                 -- 所有者ID
    max_members INTEGER DEFAULT 100,          -- 最大成员数
    config JSONB,                             -- 租户配置（JSON格式）
    version INTEGER DEFAULT 0,                -- 乐观锁版本号
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_tenants_code ON tenants(code);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_owner_id ON tenants(owner_id);
```

#### tenant_members 表（租户成员关联）
```sql
CREATE TABLE tenant_members (
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role SMALLINT DEFAULT 2 NOT NULL,         -- 角色：0-所有者，1-管理员，2-成员，3-访客
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, user_id)
);

-- 索引设计
CREATE INDEX idx_tenant_members_user_id ON tenant_members(user_id);
CREATE INDEX idx_tenant_members_role ON tenant_members(role);
```

#### tenant_read_model 表（租户读模型 - CQRS）
```sql
CREATE TABLE tenant_read_model (
    tenant_id BIGINT PRIMARY KEY,             -- 租户ID
    code VARCHAR(20) NOT NULL,                -- 租户编码
    name VARCHAR(100) NOT NULL,               -- 租户名称
    description TEXT,                         -- 租户描述
    status SMALLINT DEFAULT 0,                -- 状态
    owner_id BIGINT NOT NULL,                 -- 所有者ID
    max_members INTEGER DEFAULT 100,          -- 最大成员数
    member_count INTEGER DEFAULT 0,           -- 当前成员数
    config JSONB,                             -- 租户配置
    created_at TIMESTAMP,                     -- 创建时间
    updated_at TIMESTAMP                      -- 更新时间
);

-- 索引设计
CREATE INDEX idx_tenant_read_model_code ON tenant_read_model(code);
CREATE INDEX idx_tenant_read_model_status ON tenant_read_model(status);
CREATE INDEX idx_tenant_read_model_owner_id ON tenant_read_model(owner_id);
```

#### domain_events 表（领域事件存储）
```sql
CREATE TABLE domain_events (
    id BIGSERIAL PRIMARY KEY,                 -- 事件ID
    aggregate_id VARCHAR(50) NOT NULL,        -- 聚合根ID
    aggregate_type VARCHAR(100) NOT NULL,     -- 聚合类型
    event_type VARCHAR(100) NOT NULL,         -- 事件类型
    event_version INTEGER NOT NULL,           -- 事件版本
    event_data JSONB NOT NULL,                -- 事件数据
    occurred_on TIMESTAMP NOT NULL,           -- 事件发生时间
    processed BOOLEAN DEFAULT FALSE,          -- 是否已处理
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_id, aggregate_type);
CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_occurred ON domain_events(occurred_on);
CREATE INDEX idx_domain_events_processed ON domain_events(processed) WHERE processed = FALSE;
```

#### tenant_configs 表（租户配置）
```sql
CREATE TABLE tenant_configs (
    tenant_id BIGINT PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    max_members INTEGER DEFAULT 100,          -- 最大成员数
    storage_limit BIGINT DEFAULT 10737418240, -- 存储限制（字节）
    feature_flags JSONB,                      -- 功能开关
    branding JSONB,                           -- 品牌配置
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. 用户租户关联表

#### user_tenants 表（用户租户多对多关系）
```sql
CREATE TABLE user_tenants (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role_id BIGINT,                           -- 在该租户下的角色ID
    status SMALLINT DEFAULT 1 NOT NULL,       -- 关系状态：1-正常，2-禁用
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,                        -- 离开时间（软删除）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, tenant_id)
);

-- 索引设计
CREATE INDEX idx_user_tenants_tenant_id ON user_tenants(tenant_id);
CREATE INDEX idx_user_tenants_role_id ON user_tenants(role_id);
CREATE INDEX idx_user_tenants_status ON user_tenants(status);
CREATE INDEX idx_user_tenants_joined_at ON user_tenants(joined_at);
```

### 4. 权限相关表

#### permissions 表（权限定义）
```sql
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID
    resource VARCHAR(50) NOT NULL,            -- 资源名称
    action VARCHAR(20) NOT NULL,              -- 操作名称
    description TEXT,                         -- 权限描述
    scope VARCHAR(20) DEFAULT 'system',       -- 权限作用域：system/global/tenant
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action, scope)
);

-- 索引设计
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_scope ON permissions(scope);
```

#### roles 表（角色定义）
```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID
    name VARCHAR(50) NOT NULL,                -- 角色名称
    description TEXT,                         -- 角色描述
    tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE, -- 所属租户（NULL表示系统角色）
    is_system BOOLEAN DEFAULT FALSE,          -- 是否系统角色
    status SMALLINT DEFAULT 1 NOT NULL,       -- 状态：1-正常，2-禁用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, tenant_id)
);

-- 索引设计
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_status ON roles(status);
```

#### role_permissions 表（角色权限关联）
```sql
CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- 索引设计
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

### 5. 审计日志表

#### audit_logs 表（操作审计）
```sql
CREATE TABLE audit_logs (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID
    user_id BIGINT REFERENCES users(id),      -- 操作用户
    tenant_id BIGINT REFERENCES tenants(id),  -- 操作租户
    action VARCHAR(50) NOT NULL,              -- 操作类型
    resource_type VARCHAR(50),                -- 资源类型
    resource_id BIGINT,                       -- 资源ID
    details JSONB,                            -- 操作详情
    ip_address INET,                          -- IP地址
    user_agent TEXT,                          -- 用户代理
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
```

## Snowflake ID 实现规范

### 算法实现
```go
package snowflake

import (
    "sync"
    "time"
)

const (
    epoch        = int64(1640995200000) // 2022-01-01 00:00:00 UTC
    nodeBits     = uint(10)
    stepBits     = uint(12)
    nodeMax      = int64(-1 ^ (-1 << nodeBits))
    stepMask     = int64(-1 ^ (-1 << stepBits))
    timeShift    = uint(nodeBits + stepBits)
    nodeShift    = uint(stepBits)
)

type Node struct {
    mu        sync.Mutex
    timestamp int64
    node      int64
    step      int64
}

func NewNode(node int64) (*Node, error) {
    if node < 0 || node > nodeMax {
        return nil, fmt.Errorf("node number must be between 0 and %d", nodeMax)
    }
    
    return &Node{
        timestamp: 0,
        node:      node,
        step:      0,
    }, nil
}

func (n *Node) Generate() int64 {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    now := time.Now().UnixMilli()
    
    if now == n.timestamp {
        n.step = (n.step + 1) & stepMask
        if n.step == 0 {
            for now <= n.timestamp {
                now = time.Now().UnixMilli()
            }
        }
    } else {
        n.step = 0
    }
    
    n.timestamp = now
    
    result := (now-epoch)<<timeShift |
        (n.node << nodeShift) |
        n.step
    
    return result
}
```

### 使用示例
```go
// 初始化节点
node, err := snowflake.NewNode(1)
if err != nil {
    log.Fatal(err)
}

// 生成ID
userID := node.Generate()
fmt.Printf("Generated Snowflake ID: %d\n", userID)

// 解析ID信息
func ParseSnowflakeID(id int64) (timestamp time.Time, nodeID, sequence int64) {
    timestamp = time.UnixMilli((id >> timeShift) + epoch)
    nodeID = (id >> nodeShift) & nodeMax
    sequence = id & stepMask
    return
}
```

### 性能基准测试
```go
func BenchmarkSnowflakeGenerate(b *testing.B) {
    node, _ := NewNode(1)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = node.Generate()
    }
}

// Benchmark结果：约1000万次/秒
```

## 索引策略优化

### 核心索引原则
1. **选择性原则**：高选择性的字段优先建立索引
2. **查询频率原则**：高频查询字段建立索引
3. **复合索引原则**：根据查询条件的使用频率排序

### 复合索引示例
```
-- 用户查询优化索引
CREATE INDEX idx_users_status_created_at ON users(status, created_at DESC);

-- 租户成员查询优化索引
CREATE INDEX idx_user_tenants_tenant_status_joined 
ON user_tenants(tenant_id, status, joined_at DESC);

-- 审计日志查询优化索引
CREATE INDEX idx_audit_logs_tenant_action_date 
ON audit_logs(tenant_id, action, created_at DESC);
```

### 部分索引
```sql
-- 只为活跃用户建立索引
CREATE INDEX idx_users_active_username 
ON users(username) WHERE status = 1;

-- 只为系统权限建立索引
CREATE INDEX idx_permissions_system_resource 
ON permissions(resource, action) WHERE scope = 'system';
```

## 数据库迁移管理

### 迁移文件命名规范
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_tenants_table.up.sql
├── 000002_create_tenants_table.down.sql
├── 000003_add_user_profiles.up.sql
└── 000003_add_user_profiles.down.sql
```

### 核心迁移脚本示例

**000001_create_users_table.up.sql**
```sql
-- 创建用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(32) NOT NULL,
    status SMALLINT DEFAULT 1 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- 创建触发器自动更新updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

**000001_create_users_table.down.sql**
```sql
-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除索引
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_status;

-- 删除表
DROP TABLE IF EXISTS users;
```

### 迁移执行策略

**开发环境迁移：**
```bash
# 自动应用所有迁移
migrate -database postgres://localhost:5432/scaffold_dev?sslmode=disable -path ./migrations up

# 回滚最近一次迁移
migrate -database postgres://localhost:5432/scaffold_dev?sslmode=disable -path ./migrations down 1
```

**生产环境迁移：**
```bash
#!/bin/bash
# production-migrate.sh

set -e

echo "Starting database migration..."

# 备份数据库
pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql

# 执行迁移
migrate -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSL_MODE" \
        -path ./migrations up

echo "Migration completed successfully"
```

## 数据一致性保障

### 事务隔离级别
```sql
-- 设置适当的隔离级别
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;

-- 对于关键业务操作使用SERIALIZABLE
BEGIN ISOLATION LEVEL SERIALIZABLE;
-- 关键业务逻辑
COMMIT;
```

### 外键约束
```sql
-- 强制外键约束
ALTER TABLE user_tenants 
ADD CONSTRAINT fk_user_tenants_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_tenants 
ADD CONSTRAINT fk_user_tenants_tenant 
FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
```

### 检查约束
```sql
-- 状态字段约束
ALTER TABLE users 
ADD CONSTRAINT chk_users_status 
CHECK (status IN (1, 2, 3));

-- 邮箱格式约束
ALTER TABLE users 
ADD CONSTRAINT chk_users_email 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
```

## 性能监控视图

### 数据库性能监控
```sql
-- 慢查询监控视图
CREATE VIEW slow_queries AS
SELECT 
    userid,
    dbid,
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
WHERE mean_time > 100  -- 平均执行时间超过100ms
ORDER BY mean_time DESC;

-- 表大小监控
CREATE VIEW table_sizes AS
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_stat_get_live_tuples(schemaname||'.'||tablename) AS live_rows,
    pg_stat_get_dead_tuples(schemaname||'.'||tablename) AS dead_rows
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

这个数据库设计文档为项目提供了完整的数据层设计规范和实现指南。