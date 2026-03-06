# 快速开始 - 安装指南

## 📋 前置要求

### 必需环境
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Git

### 推荐工具
- VS Code + Go 插件
- Docker Desktop（可选，用于容器化开发）

---

## 🚀 安装步骤

### 1. 克隆项目

```bash
git clone https://github.com/your-org/ddd-scaffold.git
cd ddd-scaffold/backend
```

---

### 2. 安装依赖

```bash
go mod download
```

---

### 3. 配置数据库

#### 方式 A: 使用 Docker（推荐）

```bash
docker run --name postgres \
  -e POSTGRES_DB=ddd_scaffold \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15
```

#### 方式 B: 本地安装

```sql
CREATE DATABASE ddd_scaffold;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE ddd_scaffold TO postgres;
```

---

### 4. 配置 Redis

#### 使用 Docker

```bash
docker run --name redis \
  -p 6379:6379 \
  -d redis:7-alpine
```

---

### 5. 配置文件

复制配置模板：

```bash
cp config/config.yaml.example config/config.yaml
```

编辑 `config/config.yaml`：

```yaml
app:
  name: "ddd-scaffold"
  port: 8080
  env: "development"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "ddd_scaffold"
  sslmode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret_key: "your-secret-key-change-in-production"
  expire_in: 24h
```

---

### 6. 运行数据库迁移

```bash
make migrate-up
```

---

### 7. 启动服务

```bash
make run
```

访问 http://localhost:8080/swagger/index.html 查看 API 文档。

---

## ✅ 验证安装

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 注册新用户
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123",
    "nickname": "TestUser"
  }'
```

---

## 🐛 常见问题

### Q: 数据库连接失败？
A: 检查 PostgreSQL 是否运行，端口是否正确

### Q: Redis 连接失败？
A: 检查 Redis 是否运行，默认端口 6379

### Q: 编译错误？
A: 确保 Go 版本 ≥ 1.21，运行 `go mod tidy`

---

**下一步**: [5 分钟快速体验](quickstart.md)
