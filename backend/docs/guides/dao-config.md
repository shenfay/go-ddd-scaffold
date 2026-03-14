# DAO 数据库生成器配置说明

## 概述

`go-ddd-scaffold generate dao` 命令支持多种数据库配置加载方式，优先级从高到低如下：

1. **命令行参数** (`--dsn`)
2. **环境变量**
3. **配置文件** (未来支持)
4. **默认值**

## 配置加载方式

### 方式一：命令行参数（最高优先级）

直接通过 `--dsn` 参数指定数据库连接字符串：

```bash
go run cmd/cli/main.go generate dao \
  --dsn "host=localhost user=postgres password=secret dbname=go_ddd_scaffold port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

**适用场景**：
- CI/CD 流水线
- 临时连接测试
- 脚本自动化

### 方式二：环境变量（推荐）

通过环境变量配置数据库连接信息：

```bash
# 方式 A：使用 APP_DATABASE_* 前缀（项目标准）
export APP_DATABASE_HOST=localhost
export APP_DATABASE_PORT=5432
export APP_DATABASE_USER=postgres
export APP_DATABASE_PASSWORD=your_password
export APP_DATABASE_NAME=go_ddd_scaffold
export APP_DATABASE_SSL_MODE=disable

# 方式 B：使用 DATABASE_* 前缀（通用标准）
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=your_password
export DATABASE_NAME=go_ddd_scaffold
export DATABASE_SSL_MODE=disable
```

然后执行生成命令：

```bash
go run cmd/cli/main.go generate dao
```

**环境变量优先级**：
1. `APP_DATABASE_*` (项目标准)
2. `DATABASE_*` (通用标准)
3. 默认值

**适用场景**：
- 本地开发环境
- Docker 容器部署
- 生产环境配置

### 方式三：配置文件（未来支持）

计划支持从 `backend/configs/config.yaml` 或 `backend/configs/config_development.yaml` 读取配置。

### 方式四：默认值

如果未指定任何配置，将使用以下默认值：

```
Host:     localhost
Port:     5432
User:     postgres
Password: postgres
Name:     go_ddd_scaffold
SSLMode:  disable
```

## 环境变量完整列表

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `APP_DATABASE_HOST` | 数据库主机地址 | `localhost` | `localhost`, `127.0.0.1`, `db.example.com` |
| `APP_DATABASE_PORT` | 数据库端口 | `5432` | `5432`, `5433` |
| `APP_DATABASE_USER` | 数据库用户名 | `postgres` | `postgres`, `admin` |
| `APP_DATABASE_PASSWORD` | 数据库密码 | `postgres` | `your_secure_password` |
| `APP_DATABASE_NAME` | 数据库名称 | `go_ddd_scaffold` | `go_ddd_scaffold`, `myapp_db` |
| `APP_DATABASE_SSL_MODE` | SSL 模式 | `disable` | `disable`, `require`, `verify-full` |

也支持不带 `APP_` 前缀的通用格式：
- `DATABASE_HOST`
- `DATABASE_PORT`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `DATABASE_NAME`
- `DATABASE_SSL_MODE`

## 使用示例

### 示例 1：本地开发环境

```bash
# 设置环境变量
export APP_DATABASE_HOST=localhost
export APP_DATABASE_PORT=5432
export APP_DATABASE_USER=postgres
export APP_DATABASE_PASSWORD=postgres
export APP_DATABASE_NAME=go_ddd_scaffold
export APP_DATABASE_SSL_MODE=disable

# 生成所有核心表
go run cmd/cli/main.go generate dao
```

### 示例 2：Docker 环境

```bash
# docker-compose.yml 中配置
services:
  app:
    image: go-ddd-scaffold
    environment:
      - APP_DATABASE_HOST=db
      - APP_DATABASE_PORT=5432
      - APP_DATABASE_USER=postgres
      - APP_DATABASE_PASSWORD=secret
      - APP_DATABASE_NAME=go_ddd_scaffold
      - APP_DATABASE_SSL_MODE=disable
    depends_on:
      - db
  
  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=go_ddd_scaffold
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=secret
    ports:
      - "5432:5432"
```

### 示例 3：生产环境

```bash
# 生产环境配置（使用 SSL）
export APP_DATABASE_HOST=prod-db.example.com
export APP_DATABASE_PORT=5432
export APP_DATABASE_USER=app_user
export APP_DATABASE_PASSWORD=very_secure_password
export APP_DATABASE_NAME=go_ddd_scaffold_prod
export APP_DATABASE_SSL_MODE=require

# 生成特定表
go run cmd/cli/main.go generate dao \
  -t users,tenants,tenant_members
```

### 示例 4：使用 .env 文件

创建 `.env` 文件：

```bash
# .env
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=postgres
APP_DATABASE_NAME=go_ddd_scaffold
APP_DATABASE_SSL_MODE=disable
```

使用 [direnv](https://direnv.net/) 自动加载：

```bash
# 安装 direnv
brew install direnv

# 配置 .envrc
echo "source .env" > .envrc
direnv allow

# 现在可以直接运行命令
go run cmd/cli/main.go generate dao
```

## 输出目录配置

可以通过 `-o` 参数指定生成代码的输出目录：

```bash
# 默认输出到 internal/infrastructure/persistence/dao
go run cmd/cli/main.go generate dao

# 指定输出目录
go run cmd/cli/main.go generate dao \
  -o internal/infrastructure/persistence/dao
```

**注意**：用户已将默认输出目录从 `internal/infrastructure/persistence/gorm/dao` 修改为 `internal/infrastructure/persistence/dao`。

## 生成的表配置

默认情况下，生成以下 10 个核心业务表：

### 1. 用户与认证
- `users`

### 2. 租户管理
- `tenants`
- `tenant_members`
- `tenant_configs`

### 3. RBAC 权限系统
- `roles`
- `permissions`
- `role_permissions`

### 4. 审计与日志
- `audit_logs`
- `login_logs`

### 5. DDD 基础设施
- `domain_events`

## 故障排除

### 问题 1: 数据库连接失败

**错误信息**：
```
数据库连接失败：dial tcp [::1]:5432: connect: connection refused
```

**解决方案**：
1. 确认 PostgreSQL 服务正在运行
2. 检查数据库配置是否正确
3. 验证数据库是否存在

```bash
# 检查 PostgreSQL 服务状态
brew services list | grep postgres  # macOS
systemctl status postgresql        # Linux

# 启动 PostgreSQL 服务
brew services start postgresql     # macOS
sudo systemctl start postgresql    # Linux

# 验证数据库存在
psql -U postgres -l | grep go_ddd_scaffold
```

### 问题 2: 认证失败

**错误信息**：
```
FATAL: password authentication failed for user "postgres"
```

**解决方案**：
1. 检查密码是否正确
2. 确认用户是否存在
3. 检查 pg_hba.conf 配置

```bash
# 重置密码
psql -U postgres
ALTER USER postgres WITH PASSWORD 'new_password';
```

### 问题 3: 表不存在

**错误信息**：
```
relation "users" does not exist
```

**解决方案**：
1. 确认数据库中存在这些表
2. 运行数据库迁移

```bash
# 查看当前数据库中的表
psql -U postgres -d go_ddd_scaffold -c "\dt"

# 运行迁移
make migrate-up
# 或手动执行
migrate -path migrations -database "postgres://..." up
```

### 问题 4: 环境变量未生效

**调试方法**：

```bash
# 打印当前环境变量
env | grep DATABASE

# 验证环境变量值
echo $APP_DATABASE_HOST
echo $DATABASE_HOST

# 在 CLI 中添加调试输出
go run cmd/cli/main.go generate dao --verbose
```

## 最佳实践

### 1. 使用 .env 文件管理本地配置

```bash
# 创建 .env 文件（不提交到 Git）
cat > .env << EOF
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=postgres
APP_DATABASE_NAME=go_ddd_scaffold
APP_DATABASE_SSL_MODE=disable
EOF

# 添加到 .gitignore
echo ".env" >> .gitignore
```

### 2. 提供配置示例文件

```bash
# 创建 .env.example 文件（提交到 Git）
cat > .env.example << EOF
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=your_password_here
APP_DATABASE_NAME=go_ddd_scaffold
APP_DATABASE_SSL_MODE=disable
EOF
```

### 3. 使用 Makefile 简化命令

```makefile
.PHONY: gen-dao

gen-dao:
	@echo "Generating DAO layer from database..."
	go run cmd/cli/main.go generate dao
```

然后执行：
```bash
make gen-dao
```

### 4. 区分不同环境

```bash
# development.env
APP_DATABASE_HOST=localhost
APP_DATABASE_NAME=go_ddd_scaffold_dev

# staging.env
APP_DATABASE_HOST=staging-db.example.com
APP_DATABASE_NAME=go_ddd_scaffold_staging

# production.env
APP_DATABASE_HOST=prod-db.example.com
APP_DATABASE_NAME=go_ddd_scaffold_prod
```

使用时指定环境文件：
```bash
set -a
source development.env
set +a
go run cmd/cli/main.go generate dao
```

## 安全建议

### 1. 不要在代码中硬编码密码

❌ **错误做法**：
```go
password := "hardcoded_password"
```

✅ **正确做法**：
```go
password := os.Getenv("APP_DATABASE_PASSWORD")
```

### 2. 使用强密码

```bash
# 生成随机强密码
openssl rand -base64 32
# 或
pwgen -s 32 1
```

### 3. 限制数据库权限

为应用创建专用的数据库用户，仅授予必要的权限：

```sql
-- 创建应用用户
CREATE USER app_user WITH PASSWORD 'secure_password';

-- 授予必要权限
GRANT CONNECT ON DATABASE go_ddd_scaffold TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
```

### 4. 生产环境启用 SSL

```bash
export APP_DATABASE_SSL_MODE=require
# 或更严格的验证模式
export APP_DATABASE_SSL_MODE=verify-full
```

## 总结

✅ **三种配置方式**：
1. 命令行参数（适合临时使用）
2. 环境变量（推荐，适合所有环境）
3. 配置文件（未来支持）

✅ **两个环境变量前缀**：
1. `APP_DATABASE_*` - 项目标准
2. `DATABASE_*` - 通用标准

✅ **一个默认值**：
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Name: go_ddd_scaffold
- SSLMode: disable

选择最适合你工作流的配置方式即可！
