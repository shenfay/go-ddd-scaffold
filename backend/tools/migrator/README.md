# 数据库迁移工具

## 目录结构

```
tools/migrator/
├── migrate.sh      # 迁移执行脚本
└── README.md       # 使用说明
```

## migrate.sh - 迁移执行脚本

### 功能特点

✅ **支持多种环境**  
- 开发环境、测试环境、生产环境均可使用
- 自动读取环境变量配置
- 支持 SSL 连接配置

✅ **完整的迁移命令**
- `up` - 应用所有待处理的迁移
- `down` - 回滚最近一次迁移
- `version` - 查看当前版本
- `force` - 强制设置版本（用于修复脏状态）

✅ **友好的交互**
- 彩色终端输出
- 详细的帮助信息
- 完整的错误提示

### 使用方法

```bash
# 从项目根目录执行
cd backend

# 应用所有迁移（默认命令）
./tools/migrator/migrate.sh

# 或者指定命令
./tools/migrator/migrate.sh up       # 应用所有迁移
./tools/migrator/migrate.sh down     # 回滚一次
./tools/migrator/migrate.sh version  # 查看版本
./tools/migrator/migrate.sh force 10 # 强制设置为版本 10
./tools/migrator/migrate.sh help     # 显示帮助
```

### 环境变量配置

脚本会按以下优先级读取配置：

1. `APP_DATABASE_*` (推荐，与项目配置一致)
2. `DB_*` (兼容标准命名)
3. 默认值

**完整的环境变量列表：**

```bash
# 数据库连接配置（推荐）
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=go_ddd_scaffold
APP_DATABASE_USER=shenfay
APP_DATABASE_PASSWORD=postgres
APP_DATABASE_SSL_MODE=disable

# 或者使用简化的 DB_ 前缀（兼容）
DB_HOST=localhost
DB_PORT=5432
DB_NAME=go_ddd_scaffold
DB_USER=shenfay
DB_PASSWORD=postgres
DB_SSL_MODE=disable
```

### 命令说明

### 命令说明

#### up - 应用所有迁移

```bash
./tools/migrator/migrate.sh up
```

执行所有待处理的 `.up.sql` 文件，按版本号顺序应用。

**输出示例：**
```
检查数据库连接...
✓ 可以连接到 PostgreSQL

开始应用迁移...
1/u create_users_table (19.67ms)
2/u create_tenants_table (34.00ms)
...
✓ 迁移应用成功
```

#### down - 回滚迁移

```bash
./tools/migrator/migrate.sh down
```

回滚最近一次执行的迁移（执行对应的 `.down.sql` 文件）。

#### version - 查看版本

```bash
./tools/migrator/migrate.sh version
```

显示当前数据库的迁移版本号。

**输出示例：**
```
当前数据库版本：10
```

#### force - 强制设置版本

```bash
./tools/migrator/migrate.sh force 10
```

强制将数据库版本设置为指定值，不执行任何 SQL 文件。用于修复脏状态或版本不一致问题。

#### help - 显示帮助

```bash
./tools/migrator/migrate.sh help
```

显示完整的帮助信息，包括可用命令和配置说明。

## 手动操作命令

### 应用迁移

```bash
migrate -database "postgres://user:password@host:5432/dbname?sslmode=disable" \
  -path ./migrations up
```

### 回滚迁移

```bash
# 回滚最近一次
migrate -database "$DATABASE_URL" -path ./migrations down 1

# 回滚所有
migrate -database "$DATABASE_URL" -path ./migrations down -all
```

### 查看版本

```bash
migrate -database "$DATABASE_URL" -path ./migrations version
```

### 强制设置版本

```bash
# 用于修复脏状态
migrate -database "$DATABASE_URL" -path ./migrations force 10
```

## 故障排查

### 问题：Dirty database version

**原因**：迁移过程中断，数据库处于不一致状态

**解决方案**：
```bash
# 1. 强制设置到干净状态
migrate -database "$DATABASE_URL" -path ./migrations force 10

# 2. 重新运行测试脚本
./tools/migrator/migrate_test.sh
```

### 问题：无法连接数据库

**检查清单**：
- [ ] PostgreSQL 服务是否运行
- [ ] 主机名和端口是否正确
- [ ] 用户名和密码是否正确
- [ ] 数据库是否存在
- [ ] 用户是否有权限访问该数据库

### 问题：权限不足

**解决方案**：
```sql
-- 授予用户创建数据库的权限
ALTER USER your_user CREATEDB;

-- 或者授予超级用户权限
ALTER USER your_user SUPERUSER;
```

## 相关文档

- [迁移文件 README](../migrations/README.md)
- [数据库设计文档](../docs/reference/database-design.md)
- [表结构总览](../docs/reference/database-schema-overview.md)
