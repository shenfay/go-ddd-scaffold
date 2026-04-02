# 快速启动指南

本指南将帮助你快速启动和测试 DDD Scaffold 项目。

## 📋 前置条件

确保已安装以下工具：

- **Go** 1.25+
- **Docker** 和 **Docker Compose**（可选，用于容器化部署）
- **PostgreSQL** 15+（本地运行需要）
- **Redis** 7+（本地运行需要）

---

## 🚀 启动方式

### 方式一：使用 Docker Compose（推荐 ⭐）

最简单的方式，一键启动所有服务。

```bash
# 1. 启动所有服务（PostgreSQL, Redis, API, Worker）
docker-compose up -d

# 2. 查看日志
docker-compose logs -f

# 3. 检查服务状态
docker-compose ps

# 4. 停止所有服务
docker-compose down
```

**访问服务**：
- **API**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

---

### 方式二：本地运行（开发模式）

需要先手动启动 PostgreSQL 和 Redis。

#### 1. 启动基础设施

**使用 Docker 快速启动**：
```bash
# 只启动数据库和 Redis
docker-compose up -d postgres redis
```

**或使用本地安装**：
```bash
# macOS (Homebrew)
brew install postgresql@15 redis
brew services start postgresql@15
brew services start redis
```

#### 2. 设置环境变量

```bash
# 创建 .env 文件（或导出环境变量）
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=go_ddd_scaffold
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_SSLMODE=disable

export REDIS_ADDR=localhost:6379

export JWT_SECRET=dev-secret-key-not-for-production-use-long-random-string-in-prod
```

#### 3. 安装依赖

```bash
cd backend
make setup
make install-deps
```

#### 4. 运行数据库迁移

```bash
make migrate-up
```

#### 5. 启动 API 服务

```bash
# 终端 1
make run
```

#### 6. 启动 Worker 服务（可选）

```bash
# 终端 2
make run-worker
```

#### 7. 启动 asynqmon 监控（可选）

```bash
# 安装
make asynqmon-install

# 启动监控 UI
make asynqmon

# 访问 http://localhost:8080
```

---

### 方式三：二进制文件运行（生产模式）

```bash
# 1. 编译
make build
make build-worker

# 2. 运行
./bin/go-ddd-scaffold
./bin/go-ddd-scaffold-worker
```

---

## 🧪 测试验证

### 1. 健康检查

```bash
# 简单检查
curl http://localhost:8080/health

# 详细检查（包含 DB 和 Redis 状态）
curl http://localhost:8080/health | jq

# 预期输出：
# {
#   "status": "ok",
#   "checks": {
#     "database": {"status": "ok", "response_time_ms": "<10ms"},
#     "redis": {"status": "ok", "response_time_ms": "<10ms", "ping_result": "PONG"}
#   }
# }
```

### 2. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }' | jq

# 预期输出：
# {
#   "user": {...},
#   "access_token": "eyJhbGc...",
#   "refresh_token": "eyJhbGc...",
#   "expires_in": 1800
# }
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }' | jq
```

### 4. 刷新 Token

```bash
# 使用登录响应中的 refresh_token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGc..."
  }' | jq
```

### 5. 用户退出

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGc..."
  }' | jq
```

---

## 📊 监控和可观测性

### Prometheus 指标

```bash
# 访问 metrics 端点
curl http://localhost:8080/metrics

# 示例指标：
# http_requests_total
# http_request_duration_seconds_bucket
# auth_attempts_total
# user_registrations_total
```

### Grafana 仪表盘（预留）

配置 Prometheus 数据源后，导入以下 Dashboard ID：
- HTTP 请求监控
- 业务指标监控
- 系统资源监控

---

## 🔧 常用 Make 命令

```bash
# 开发
make run              # 启动 API
make run-worker       # 启动 Worker
make setup            # 配置开发环境
make install-deps     # 安装依赖

# 构建
make build            # 构建 API
make build-worker     # 构建 Worker
make build-linux      # 构建 Linux 版本
make clean            # 清理构建产物

# 测试
make test             # 运行所有测试
make test-short       # 运行快速测试
make coverage         # 生成测试覆盖率报告

# 代码质量
make fmt              # 格式化代码
make vet              # 运行 go vet
make lint             # 运行 linter

# 数据库
make migrate-up       # 执行数据库迁移
make migrate-down     # 回滚迁移

# 监控
make asynqmon-install    # 安装 asynqmon
make asynqmon            # 启动任务监控 UI
make health              # 检查应用健康状态

# 文档
make swagger-gen         # 生成 Swagger 文档
make swagger-serve       # 启动 Swagger UI 服务
```

---

## 🐛 故障排查

### 问题 1: 数据库连接失败

```bash
# 检查 PostgreSQL 是否运行
docker-compose ps postgres
# 或
pg_isready -h localhost -p 5432

# 重启 PostgreSQL
docker-compose restart postgres
```

### 问题 2: Redis 连接失败

```bash
# 检查 Redis 是否运行
docker-compose ps redis
# 或
redis-cli ping

# 重启 Redis
docker-compose restart redis
```

### 问题 3: 端口被占用

```bash
# 检查端口占用
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 修改 docker-compose.yml 中的端口映射
# 或停止占用端口的服务
```

### 问题 4: 迁移失败

```bash
# 查看详细错误
make migrate-up

# 重置数据库（⚠️ 会丢失数据）
docker-compose down -v
docker-compose up -d postgres
make migrate-up
```

---

## 📝 下一步

启动成功后，可以：

1. **阅读 API 文档**: http://localhost:8080/swagger/index.html
2. **测试认证流程**: 注册 → 登录 → 刷新 → 退出
3. **监控任务队列**: http://localhost:8080 (asynqmon)
4. **查看指标**: http://localhost:8080/metrics
5. **开始业务开发**: 参考 `internal/auth/` 模块的实现模式

---

## 🎯 成功标志

✅ API 服务启动在 http://localhost:8080  
✅ Worker 服务正常处理异步任务  
✅ PostgreSQL 数据库连接正常  
✅ Redis 缓存和消息队列正常  
✅ 健康检查返回 `{"status": "ok"}`  
✅ Swagger UI 可访问  

**恭喜！你的 DDD Scaffold 已经成功启动！** 🎉
