# 快速开始指南

5 分钟内在本地运行 DDD Scaffold 项目。

## 📋 前置要求

### 必需软件

| 软件 | 版本 | 用途 | 安装方式 |
|------|------|------|---------|
| Go | 1.21+ | 后端运行时 | `brew install go` |
| PostgreSQL | 14+ | 主数据库 | `brew install postgresql` |
| Redis | 7+ | 缓存/会话存储 | `brew install redis` |

### 可选软件

| 软件 | 用途 | 安装方式 |
|------|------|---------|
| Docker & Docker Compose | 容器化部署 | [Docker Desktop](https://www.docker.com/products/docker-desktop/) |
| Swag | Swagger 文档生成 | `go install github.com/swaggo/swag/cmd/swag@latest` |

## 🚀 快速启动（推荐方式）

### 方式一：使用 Docker Compose（最简单）

```bash
# 1. 克隆项目
git clone <repository-url>
cd ddd-scaffold

# 2. 启动所有基础设施（PostgreSQL + Redis）
docker-compose up -d postgres redis

# 3. 等待服务就绪（约 10 秒）
sleep 10

# 4. 进入后端目录
cd backend

# 5. 运行数据库迁移
make migrate up

# 6. 启动 API 服务
make run api
```

**验证服务**：
```bash
# 健康检查
curl http://localhost:8080/health

# 访问 Swagger 文档
open http://localhost:8080/swagger/index.html
```

### 方式二：本地安装（开发推荐）

#### 步骤 1：启动基础设施

```bash
# 启动 PostgreSQL
brew services start postgresql

# 启动 Redis
brew services start redis

# 验证服务
psql -U postgres -c "SELECT version();"
redis-cli ping  # 应返回 PONG
```

#### 步骤 2：创建数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE ddd_scaffold;"
psql -U postgres -c "CREATE USER ddd_scaffold WITH PASSWORD 'ddd_scaffold';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ddd_scaffold TO ddd_scaffold;"
```

#### 步骤 3：配置环境变量

```bash
cd backend
cp configs/.env.example configs/.env

# 编辑配置文件（根据实际情况修改）
vim configs/.env
```

**关键配置项**：
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ddd_scaffold

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d
```

#### 步骤 4：安装依赖并运行

```bash
# 安装 Go 依赖
go mod download

# 运行数据库迁移
make migrate up

# 启动 API 服务
make run api
```

## 📝 测试 API

### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456!"
  }'
```

**预期响应**：
```json
{
  "code": "SUCCESS",
  "message": "注册成功",
  "data": {
    "user": {
      "id": "01JQMXYZ...",
      "email": "test@example.com",
      "email_verified": false
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456!"
  }'
```

### 3. 获取当前用户信息

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

## 🛠️ 常用操作

### 启动 Worker（异步任务处理）

```bash
# 新终端窗口
cd backend
make run worker
```

### 生成 Swagger 文档

```bash
make swagger gen

# 查看生成的文档
ls -la api/swagger/
# docs.go
# swagger.json
# swagger.yaml
```

### 查看数据库迁移状态

```bash
make db-status
```

### 回滚数据库迁移

```bash
make migrate down
```

### 运行测试

```bash
# 运行所有测试
make test

# 仅运行单元测试（跳过集成测试）
make test-short

# 生成覆盖率报告
make coverage
open coverage.html
```

## 🐛 常见问题

### 问题 1：数据库连接失败

**错误信息**：
```
failed to connect database: dial tcp [::1]:5432: connect: connection refused
```

**解决方案**：
```bash
# 检查 PostgreSQL 是否运行
brew services list | grep postgresql

# 启动 PostgreSQL
brew services start postgresql

# 验证连接
psql -U postgres -d ddd_scaffold -c "SELECT 1;"
```

### 问题 2：Redis 连接失败

**错误信息**：
```
redis: connection refused
```

**解决方案**：
```bash
# 检查 Redis 是否运行
brew services list | grep redis

# 启动 Redis
brew services start redis

# 验证连接
redis-cli ping
```

### 问题 3：端口已被占用

**错误信息**：
```
listen tcp :8080: bind: address already in use
```

**解决方案**：
```bash
# 查找占用端口的进程
lsof -i :8080

# 终止进程（替换 PID）
kill -9 <PID>

# 或者修改配置使用其他端口
vim configs/.env
# 修改 SERVER_PORT=8081
```

### 问题 4：Swag 命令未找到

**错误信息**：
```
zsh: command not found: swag
```

**解决方案**：
```bash
# 安装 Swag
go install github.com/swaggo/swag/cmd/swag@latest

# 添加到 PATH（如需要）
export PATH=$PATH:$(go env GOPATH)/bin
```

## 📚 下一步

- [📖 开发指南](DEVELOPMENT_GUIDE.md) - 了解开发规范和流程
- [🏗️ DDD 架构设计](../architecture/DDD_ARCHITECTURE.md) - 深入理解架构设计
- [📡 API 文档](http://localhost:8080/swagger/index.html) - 查看所有 API 接口
- [🗄️ 数据库设计](../database/SCHEMA_DESIGN.md) - 了解数据库结构

## 💡 开发建议

1. **使用 IDE**：推荐 GoLand 或 VS Code（安装 Go 插件）
2. **启用自动保存**：避免忘记保存文件
3. **使用 `make run api`**：而非 `go run`，确保环境变量正确加载
4. **定期运行 `make vet`**：及早发现潜在问题
5. **编写测试**：新功能必须包含单元测试

---

**遇到问题？** 查看 [故障排查指南](../operations/TROUBLESHOOTING.md) 或提交 [GitHub Issue](https://github.com/shenfay/go-ddd-scaffold/issues)
