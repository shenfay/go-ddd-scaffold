# 数据库迁移指南

## 📋 概述

项目使用 SQL 文件 + GORM AutoMigrate 混合方式进行数据库迁移。

- **SQL 文件**：用于创建表结构、索引、注释等（推荐方式）
- **GORM AutoMigrate**：用于快速原型开发和测试环境

## 🚀 执行迁移

### 方式 1：使用 SQL 迁移文件（推荐）

```bash
# 执行正向迁移（创建表）
make migrate-up

# 回滚迁移（删除表）
make migrate-down
```

**迁移文件位置：**
```
backend/migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_activity_logs_table.up.sql
└── 002_create_activity_logs_table.down.sql
```

### 方式 2：使用 GORM AutoMigrate（开发环境）

```bash
# 使用 CLI 工具自动迁移
go run ./cmd/cli migrate up

# 回滚（需要手动执行 SQL）
go run ./cmd/cli migrate down
```

## 📊 迁移历史

### 001 - 用户表
- **文件**: `001_create_users_table.*`
- **描述**: 创建用户表及基础索引
- **包含字段**: id, email, password, email_verified, locked, failed_attempts, last_login_at

### 002 - 活动日志表 ✨ NEW
- **文件**: `002_create_activity_logs_table.*`
- **描述**: 创建活动日志表及查询索引
- **包含字段**: 
  - 基础信息：id, user_id, email, action, status
  - 设备信息：ip, user_agent, device, browser, os
  - 元数据：description, metadata (JSON), created_at

## 🔧 表结构详情

### activity_logs 表

```sql
CREATE TABLE activity_logs (
    id VARCHAR(50) PRIMARY KEY,           -- ULID 格式日志 ID
    user_id VARCHAR(50) NOT NULL,         -- 用户 ID（索引）
    email VARCHAR(255),                   -- 用户邮箱
    action VARCHAR(50) NOT NULL,          -- 活动类型（索引）
    status VARCHAR(20) NOT NULL,          -- 状态（索引）
    ip VARCHAR(45),                       -- IP 地址（支持 IPv6）
    user_agent VARCHAR(500),              -- User-Agent
    device VARCHAR(100),                  -- 设备类型
    browser VARCHAR(50),                  -- 浏览器名称
    os VARCHAR(50),                       -- 操作系统
    description TEXT,                     -- 描述
    metadata JSON,                        -- 元数据
    created_at TIMESTAMP WITH TIME ZONE   -- 创建时间（索引）
);
```

**索引：**
- `idx_activity_logs_user_id` - 按用户 ID 查询
- `idx_activity_logs_created_at` - 按时间排序
- `idx_activity_logs_action` - 按活动类型统计
- `idx_activity_logs_status` - 按状态筛选
- `idx_activity_logs_user_created` - 复合索引（用户 + 时间倒序）

## 💡 最佳实践

### 1. 创建新迁移

```bash
# 1. 创建 SQL 文件
touch backend/migrations/003_your_migration_name.up.sql
touch backend/migrations/003_your_migration_name.down.sql

# 2. 编写迁移 SQL
# 在 .up.sql 中写入 CREATE/ALTER 语句
# 在 .down.sql 中写入对应的 DROP/ALTER 语句

# 3. 更新 GORM AutoMigrate（可选）
# 在 cmd/cli/main.go 中添加新的 AutoMigrate 调用
```

### 2. 回滚策略

```sql
-- 好的回滚示例（可逆操作）
-- UP
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- DOWN
ALTER TABLE users DROP COLUMN phone;

-- 避免不可逆操作
-- ❌ 不推荐：删除字段且不留备份
-- ❌ 不推荐：修改字段类型导致数据丢失
```

### 3. 测试迁移

```bash
# 本地测试
make migrate-down
make migrate-up

# 验证表结构
psql -h localhost -U postgres -d go_ddd_scaffold -c "\dt"
psql -h localhost -U postgres -d go_ddd_scaffold -c "\d activity_logs"
```

## ⚠️ 注意事项

1. **生产环境谨慎使用 AutoMigrate**
   - AutoMigrate 不会删除字段
   - 可能产生冗余索引
   - 建议在生产环境使用明确的 SQL 迁移

2. **迁移顺序很重要**
   - 迁移文件按数字顺序执行（001 → 002 → 003）
   - 确保外键依赖的表先创建

3. **备份数据**
   - 执行迁移前后备份数据库
   - 特别是执行 DROP 或 ALTER 操作时

4. **向下兼容**
   - 尽量保证迁移可回滚
   - 避免破坏性变更（Breaking Changes）

## 🔍 常用命令

```bash
# 查看当前表
psql -h localhost -U postgres -d go_ddd_scaffold -c "\dt"

# 查看表结构
psql -h localhost -U postgres -d go_ddd_scaffold -c "\d activity_logs"

# 查看索引
psql -h localhost -U postgres -d go_ddd_scaffold -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'activity_logs';"

# 查看迁移历史（如果有 migration_history 表）
psql -h localhost -U postgres -d go_ddd_scaffold -c "SELECT * FROM migration_history ORDER BY executed_at DESC;"
```

## 📝 示例：添加新字段

假设要给用户表添加 `phone` 字段：

**步骤 1：创建迁移文件**
```bash
touch backend/migrations/003_add_phone_to_users.up.sql
touch backend/migrations/003_add_phone_to_users.down.sql
```

**步骤 2：编写 SQL**
```sql
-- 003_add_phone_to_users.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_phone ON users(phone);
COMMENT ON COLUMN users.phone IS '手机号码';

-- 003_add_phone_to_users.down.sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN phone;
```

**步骤 3：更新 GORM 模型（可选）**
```go
// internal/auth/user.go
type UserPO struct {
    // ... existing fields ...
    Phone string `gorm:"type:varchar(20)" json:"phone"`
}
```

**步骤 4：执行迁移**
```bash
make migrate-up
```

## 🎯 下一步

- [ ] 实现迁移版本管理（如 golang-migrate）
- [ ] 添加迁移前后数据验证
- [ ] 实现自动化迁移测试
- [ ] 集成到 CI/CD 流程
