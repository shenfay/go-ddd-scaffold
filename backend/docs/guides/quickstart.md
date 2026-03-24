# 快速开始

5 分钟快速上手 Go DDD Scaffold！

## 🚀 环境要求

### 必需工具

- **Go** 1.21+
- **PostgreSQL** 14+
- **Redis** 7+
- **Git**

### 可选工具

- **Docker** & Docker Compose（推荐）
- **TablePlus** / DBeaver（数据库管理）
- **Postman** / Insomnia（API 调试）

---

## 📦 快速安装

### 方式 1：克隆项目（推荐）

```bash
# 克隆项目
git clone https://github.com/shenfay/go-ddd-scaffold.git
cd go-ddd-scaffold/backend

# 安装依赖
go mod download

# 复制配置文件
cp configs/.env.example configs/.env
cp configs/config.yaml.example configs/config.yaml
```

### 方式 2：使用 CLI 工具

```bash
# 安装 CLI 工具
go install github.com/shenfay/go-ddd-scaffold/backend/cmd/cli@latest

# 创建新项目
ddd-scaffold new my-project

# 进入项目目录
cd my-project
```

---

## 🗄️ 数据库设置

### 使用 Docker（推荐）

```bash
# 启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis

# 验证服务
docker ps
# 应该看到 postgres 和 redis 容器运行中
```

### 手动安装

```bash
# 创建数据库
psql -U postgres
CREATE DATABASE go_ddd_scaffold;
CREATE USER scaffold WITH PASSWORD 'scaffold';
GRANT ALL PRIVILEGES ON DATABASE go_ddd_scaffold TO scaffold;
\q

# 修改 .env 文件
DB_HOST=localhost
DB_PORT=5432
DB_USER=scaffold
DB_PASSWORD=scaffold
DB_NAME=go_ddd_scaffold
```

---

## ⚙️ 配置说明

### 环境变量 (.env)

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=scaffold
DB_PASSWORD=scaffold
DB_NAME=go_ddd_scaffold

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRE=7200
JWT_REFRESH_EXPIRE=604800

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug  # debug, release, test
```

### 配置文件 (config.yaml)

```yaml
server:
  port: 8080
  mode: debug

jwt:
  secret: "your-secret-key"
  access_expire: 2h
  refresh_expire: 168h

database:
  driver: postgres
  host: localhost
  port: 5432
  user: scaffold
  password: scaffold
  database: go_ddd_scaffold
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

---

## 🔧 运行项目

### 1. 数据库迁移

```bash
# 安装 migrate 工具
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 执行迁移
migrate -path migrations -database "postgres://scaffold:scaffold@localhost:5432/go_ddd_scaffold?sslmode=disable" up

# 或使用 Makefile
make migrate-up
```

### 2. 启动 API 服务

```bash
# 开发模式
go run cmd/api/main.go

# 或使用 Makefile
make run

# 生产模式
go build -o bin/api cmd/api/main.go
./bin/api
```

### 3. 启动 Worker（可选）

```bash
# Worker 处理异步任务
go run cmd/worker/main.go
```

---

## ✅ 验证安装

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/api/health

# 预期响应
{"status":"ok","timestamp":"2024-03-23T10:00:00Z"}
```

### 测试注册功能

```bash
# 创建用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123!"
  }'

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "123456789",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### 测试登录功能

```bash
# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "password": "Password123!"
  }'

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

---

## 🎯 下一步

### 学习路径

1. **理解架构** - 阅读 [架构总览](../design/architecture-overview.md)
2. **开发规范** - 学习 [开发规范](../specifications/development-spec.md)
3. **实战教程** - 完成 [入门教程](../tutorials/getting-started-tutorial.md)
4. **深入 DDD** - 研究 [DDD 设计指南](../design/ddd-design-guide.md)

### 常用命令

```bash
# 查看所有可用命令
make help

# 常用命令
make run          # 运行 API
make worker       # 运行 Worker
make test         # 运行测试
make lint         # 代码检查
make migrate-up   # 数据库迁移
make migrate-down # 回滚迁移
```

---

## ❓ 常见问题

### Q1: 无法连接数据库

**解决方案：**
```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 或本地安装
pg_isready -h localhost -p 5432

# 检查防火墙
sudo ufw allow 5432
```

### Q2: 端口被占用

**解决方案：**
```bash
# 查看端口占用
lsof -i :8080

# 修改配置
# 在 config.yaml 中更改 server.port
```

### Q3: 依赖下载失败

**解决方案：**
```bash
# 使用国内镜像
export GOPROXY=https://goproxy.cn,direct

# 重新下载
go mod download
```

---

## 📚 参考资源

- [完整文档索引](../README.md)
- [CLI 工具指南](../guides/cli-tool-guide.md)
- [Module 开发指南](../guides/module-development-guide.md)
- [故障排查](../operations/troubleshooting.md)

---

**恭喜！** 🎉 你已经成功安装并运行了 Go DDD Scaffold！

继续学习 [架构总览](../design/architecture-overview.md) 深入了解项目架构。
