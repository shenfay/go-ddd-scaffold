# 快速启动指南

## 🚀 5 分钟快速开始

### **前置条件**
- Go 1.25+
- Docker（用于运行 PostgreSQL 和 Redis）

### **1. 启动基础设施**

```bash
# 启动 PostgreSQL
docker run --name go-ddd-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15

# 启动 Redis
docker run --name go-ddd-redis \
  -p 6379:6379 \
  -d redis:7-alpine

# 验证容器运行
docker ps
```

### **2. 配置环境变量**

```bash
# 开发环境（临时）
export DB_PASSWORD=postgres
export JWT_SECRET=dev-secret-key-not-for-production

# 或创建 .env 文件（推荐）
cat > backend/configs/.env << EOF
DB_PASSWORD=postgres
JWT_SECRET=dev-secret-key-not-for-production
EOF
```

### **3. 启动 API 服务**

```bash
cd backend

# 方式 1: 直接运行
go run ./cmd/api

# 方式 2: 使用 Makefile
make run
```

看到以下日志表示启动成功：
```
Database connection established and tables migrated
Redis connection established
Starting server on port :8080
```

### **4. 测试健康检查**

```bash
curl http://localhost:8080/health
# 返回：{"status":"healthy"}
```

---

## 📦 完整功能测试

### **场景 1: 用户注册**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "Password123!"
  }'
```

**成功响应**:
```json
{
  "user": {
    "id": "user_01H8K3M9N2P4Q5R6S7T8U9V0W1",
    "email": "john@example.com",
    "email_verified": false,
    "created_at": "2026-04-02T10:30:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "tok_01H8K3M9N2P4Q5R6S7T8U9V0W2",
  "expires_in": 1800
}
```

### **场景 2: 用户登录**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "Password123!"
  }'
```

### **场景 3: 测试速率限制**

连续快速请求登录接口 6 次，第 6 次会触发限制：

```bash
for i in {1..6}; do
  echo "Request $i:"
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrong"}' \
    -w "\nHTTP Status: %{http_code}\n\n"
done
```

**预期结果**:
- 前 5 次：返回 `401 Unauthorized`
- 第 6 次：返回 `429 Too Many Requests`

---

## 🔧 Worker 服务（可选）

如果需要处理异步任务（如发送邮件），启动 Worker：

```bash
# 另开终端
cd backend
go run ./cmd/worker
```

看到以下日志表示 Worker 启动成功：
```
Redis connection established
Starting Asynq Worker with concurrency=10
```

---

## 🛠️ 常用命令

### **Makefile 命令**

```bash
# 启动 API
make run

# 启动 Worker
make run-worker

# 构建
make build

# 测试
make test

# 代码格式化
make fmt

# 代码检查
make vet

# Lint 检查
make lint

# Swagger 文档
make swagger-gen
make swagger-serve
```

### **Docker 命令**

```bash
# 查看容器状态
docker ps

# 停止容器
docker stop go-ddd-postgres go-ddd-redis

# 删除容器
docker rm go-ddd-postgres go-ddd-redis

# 查看日志
docker logs go-ddd-postgres
docker logs go-ddd-redis
```

---

## 🐛 故障排查

### **问题 1: 数据库连接失败**

**错误信息**:
```
Failed to initialize database: dial tcp [::1]:5432: connect: connection refused
```

**解决方案**:
1. 检查 PostgreSQL 是否运行：`docker ps`
2. 检查端口是否被占用：`lsof -i :5432`
3. 重启 PostgreSQL 容器

### **问题 2: Redis 连接失败**

**错误信息**:
```
Failed to connect to Redis: dial tcp [::1]:6379: connect: connection refused
```

**解决方案**:
1. 检查 Redis 是否运行：`docker ps`
2. 检查端口是否被占用：`lsof -i :6379`
3. 重启 Redis 容器

### **问题 3: 端口被占用**

**错误信息**:
```
listen tcp :8080: bind: address already in use
```

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或修改配置文件中的端口
vim configs/development.yaml
```

### **问题 4: 编译错误**

```bash
# 清理并重新下载依赖
cd backend
go clean -modcache
go mod tidy

# 重新编译
go build ./...
```

---

## 📝 配置说明

### **开发环境配置** (`configs/development.yaml`)

```yaml
server:
  port: 8080          # API 端口
  mode: debug         # debug/release/test

database:
  host: localhost
  port: 5432
  name: go_ddd_scaffold_dev
  user: postgres
  password: ${DB_PASSWORD}  # 从环境变量读取

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  secret: ${JWT_SECRET}     # 从环境变量读取
  access_expire: 30m        # Access Token 有效期
  refresh_expire: 168h      # Refresh Token 有效期 (7 天)
```

---

## 🎯 下一步

完成快速启动后，建议按以下顺序学习：

1. **阅读架构文档**: [`docs/ARCHITECTURE_SUMMARY.md`](ARCHITECTURE_SUMMARY.md)
2. **了解 Phase 2 实现**: [`docs/PHASE2_IMPLEMENTATION.md`](PHASE2_IMPLEMENTATION.md)
3. **运行完整测试**: 参考上文的测试场景
4. **开始业务开发**: 基于此脚手架实现你的业务逻辑

---

## 💡 提示

- **开发模式**: 使用 `debug` 模式可以看到详细的日志
- **热重载**: 安装 `air` 实现代码热重载
  ```bash
  go install github.com/cosmtrek/air@latest
  air
  ```
- **监控 Worker**: 使用 `asynqmon` 查看任务队列
  ```bash
  make asynqmon-install
  make asynqmon
  # 访问 http://localhost:8080
  ```

---

**祝你开发愉快！** 🎉
