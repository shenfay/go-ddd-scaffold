# 环境变量配置指南

本文档介绍如何正确配置和使用 DDD Scaffold 的环境变量。

---

## 📁 配置文件位置

```
backend/configs/
├── .env                    # 环境变量参考文件（需手动创建）
├── .env.example            # 环境变量示例
└── development.yaml        # 实际使用的配置文件 ⭐
```

---

## 🔧 配置方式对比

### 方式一：直接修改 YAML（推荐 ⭐⭐⭐⭐⭐）

**最简单直接的方式**，适合本地开发。

**步骤**：

1. 编辑 `backend/configs/development.yaml`
2. 直接修改数据库和密码配置

```yaml
database:
  host: localhost
  port: 5432
  name: ddd_scaffold
  user: shenfay          # ✅ 修改为你的用户名
  password: your_password # ✅ 修改为你的密码
  ssl_mode: disable

redis:
  addr: localhost:6379
  password: ""           # 如果有密码，在这里填写
```

**优点**：
- ✅ 简单直观
- ✅ 不需要额外工具
- ✅ 配置集中管理

**缺点**：
- ⚠️ 密码明文写在 YAML 中
- ⚠️ 不适合生产环境

---

### 方式二：使用环境变量（推荐 ⭐⭐⭐⭐）

**更安全的方式**，适合开发和生产环境。

**步骤**：

1. **复制环境变量文件**
   ```bash
   cd backend/configs
   cp .env.example .env
   ```

2. **编辑 `.env` 文件**
   ```bash
   vim backend/configs/.env
   ```

3. **修改配置**
   ```bash
   # 服务器配置
   APP_SERVER_PORT=8080
   APP_SERVER_MODE=debug
   
   # 数据库配置
   APP_DATABASE_HOST=localhost
   APP_DATABASE_PORT=5432
   APP_DATABASE_NAME=ddd_scaffold
   APP_DATABASE_USER=shenfay
   APP_DATABASE_PASSWORD=your_password  # ✅ 修改为你的密码
   APP_DATABASE_SSL_MODE=disable
   
   # Redis 配置
   APP_REDIS_ADDR=localhost:6379
   APP_REDIS_PASSWORD=
   
   # JWT 配置
   APP_JWT_SECRET=smLuhX79IBihMgtVucmefmdP8Gt8hshA
   APP_JWT_ACCESS_EXPIRE=30m
   APP_JWT_REFRESH_EXPIRE=7d
   
   # 日志配置
   APP_LOGGING_LEVEL=debug
   APP_LOGGING_FORMAT=console
   ```

4. **在启动脚本中使用**
   
   启动脚本会自动读取 `.env` 文件中的配置。

**优点**：
- ✅ 密码与代码分离
- ✅ 更安全
- ✅ 便于 CI/CD 集成

**缺点**：
- ⚠️ 需要额外的配置步骤

---

### 方式三：终端环境变量（临时使用）

**适合快速测试**。

```bash
# 设置环境变量
export APP_DATABASE_PASSWORD=your_password
export APP_JWT_SECRET=your_secret

# 启动服务
make run
```

**优点**：
- ✅ 快速方便
- ✅ 不修改文件

**缺点**：
- ⚠️ 只在当前终端有效
- ⚠️ 每次都要设置

---

## 🔐 安全建议

### 开发环境

```yaml
# development.yaml
database:
  user: shenfay
  password: postgres  # 可以使用简单密码
  
jwt:
  secret: dev-secret-key-not-for-production-use-long-random-string-in-prod
```

### 生产环境

```yaml
# production.yaml
database:
  user: app_user
  password: ${DB_PASSWORD}  # 从环境变量读取
  
jwt:
  secret: ${JWT_SECRET}     # 必须使用强随机字符串
```

**生成安全的 JWT Secret**：

```bash
# 方法 1：使用 openssl
openssl rand -base64 32

# 方法 2：使用 head 和 sha256sum
head -c 32 /dev/urandom | sha256sum

# 方法 3：使用 pwgen
pwgen -s 64 1
```

---

## 🛠️ 自动加载环境变量

### 在 Go 代码中加载 .env

如果你希望在 Go 代码中自动加载 `.env` 文件，可以添加以下依赖：

```bash
go get github.com/joho/godotenv
```

然后在 `main.go` 中导入：

```go
import _ "github.com/joho/godotenv"
```

这样程序启动时会自动加载 `.env` 文件。

---

## 📝 当前配置状态

### 默认配置（development.yaml）

```yaml
database:
  host: localhost
  port: 5432
  name: ddd_scaffold
  user: shenfay
  password: postgres      # ⚠️ 请根据实际情况修改
  ssl_mode: disable

redis:
  addr: localhost:6379
  password: ""            # 如果有密码请修改
```

### 环境变量示例（.env）

```bash
# 数据库配置
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=ddd_scaffold
APP_DATABASE_USER=shenfay
APP_DATABASE_PASSWORD=postgres  # ⚠️ 请修改为你的密码

# Redis 配置
APP_REDIS_ADDR=localhost:6379
APP_REDIS_PASSWORD=

# JWT 配置（生产环境必须修改）
APP_JWT_SECRET=smLuhX79IBihMgtVucmefmdP8Gt8hshA
```

---

## 🚀 使用启动脚本

启动脚本会自动读取 `.env` 文件中的配置：

```bash
# 启动服务（自动使用 .env 中的配置）
./scripts/dev/start-local.sh --all

# 脚本会：
# 1. 读取 backend/configs/.env
# 2. 提取数据库和 Redis 配置
# 3. 使用 PGPASSWORD 传递密码
# 4. 检查连接是否正常
```

---

## ✅ 验证配置

### 检查 PostgreSQL 配置

```bash
# 从 .env 读取配置
source backend/configs/.env

# 测试连接（使用 PGPASSWORD 避免密码提示）
export PGPASSWORD="$APP_DATABASE_PASSWORD"
psql -h "$APP_DATABASE_HOST" -p "$APP_DATABASE_PORT" \
     -U "$APP_DATABASE_USER" -d "$APP_DATABASE_NAME" \
     -c "SELECT current_database()"

# 清理
unset PGPASSWORD
```

### 检查 Redis 配置

```bash
# 从 .env 读取配置
source backend/configs/.env

# 测试连接
redis-cli -h "${APP_REDIS_ADDR%:*}" -p "${APP_REDIS_ADDR#*:}" ping
```

---

## 🎯 最佳实践

### 1. 使用 .gitignore 保护敏感信息

```bash
# .gitignore
backend/configs/.env          # ✅ 忽略实际配置
!backend/configs/.env.example # ✅ 保留示例文件
```

### 2. 多环境配置

```
backend/configs/
├── .env.example              # 示例模板
├── development.env           # 开发环境
├── test.env                  # 测试环境
├── production.env            # 生产环境
├── development.yaml          # 开发环境 YAML
├── test.yaml                 # 测试环境 YAML
└── production.yaml           # 生产环境 YAML
```

### 3. CI/CD 中使用环境变量

```yaml
# .github/workflows/ci.yml
env:
  DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
  JWT_SECRET: ${{ secrets.JWT_SECRET }}
```

---

## 📖 相关文档

- [本地开发指南](LOCAL_DEVELOPMENT_GUIDE.md) - 完整的本地开发流程
- [监控指南](MONITORING_GUIDE.md) - 监控和可观测性
- [快速参考](QUICK_REFERENCE.md) - 常用命令速查

---

## ✅ 总结

**推荐配置方式**：

| 场景 | 推荐方式 | 说明 |
|------|----------|------|
| **本地开发** | 直接修改 YAML | 简单快捷 |
| **团队协作** | 使用 .env 文件 | 配置统一管理 |
| **生产部署** | 环境变量注入 | 最安全 |
| **CI/CD** | Secrets 管理 | 自动化集成 |

**立即开始**：

```bash
# 1. 复制环境变量文件
cp backend/configs/.env.example backend/configs/.env

# 2. 修改配置
vim backend/configs/.env

# 3. 启动服务
./scripts/dev/start-local.sh --all
```

**配置完成！** 🎉
