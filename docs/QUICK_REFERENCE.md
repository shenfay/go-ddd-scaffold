# 快速参考卡片 🚀

## ⚡ 一键启动

```bash
# 启动所有服务（推荐）
./scripts/dev/start.sh --all

# 或手动启动
docker-compose up -d
```

---

## 📊 监控工具

| 工具 | 命令 | 地址 | 用途 |
|------|------|------|------|
| **Asynqmon** | `make asynqmon` | http://localhost:8080 | 队列和任务监控 |
| **Swagger** | `make swagger-serve` | http://localhost:8080/swagger | API 文档 |
| **健康检查** | `curl localhost:8080/health` | http://localhost:8080/health | 服务状态 |
| **Metrics** | `curl localhost:8080/metrics` | http://localhost:8080/metrics | Prometheus 指标 |

---

## 🧪 测试流程

```bash
# 执行核心流程测试
./scripts/dev/core-flow-test.sh

# 测试内容：
# 1. 用户注册
# 2. 登录
# 3. 获取用户信息
# 4. 刷新 Token
# 5. 登出
# 6. Token 失效验证
```

---

## 🔍 查看队列和事件

### Asynqmon UI（最直观）

```bash
# 1. 安装
make asynqmon-install

# 2. 启动
make asynqmon

# 3. 访问
open http://localhost:8080
```

### 查看内容

✅ **Queues（队列）**
- Critical（高优先级）
- Default（默认）
- Low（低优先级）

✅ **Tasks（任务）**
- UserRegistered
- UserLoggedIn
- UserLoggedOut
- TokenRefreshed

✅ **Events（事件流）**
- 入队 → 处理中 → 完成/失败

---

## 📈 查看指标

```bash
# HTTP 指标
curl localhost:8080/metrics | grep http_requests

# 数据库指标
curl localhost:8080/metrics | grep db_queries

# Redis 指标
curl localhost:8080/metrics | grep redis_commands

# 业务指标
curl localhost:8080/metrics | grep auth_attempts
```

---

## 📝 查看日志

```bash
# API 日志
docker-compose logs -f api

# Worker 日志
docker-compose logs -f worker

# 错误日志
docker-compose logs api | grep -i error
```

---

## 🛠️ Makefile 常用命令

```bash
# 开发
make run              # 启动 API
make run-worker       # 启动 Worker

# 构建
make build            # 编译 API
make build-worker     # 编译 Worker

# 测试
make test             # 运行测试
make coverage         # 覆盖率报告

# 数据库
make migrate-up       # 执行迁移
make migrate-down     # 回滚迁移

# 监控
make asynqmon         # 启动任务监控
make health           # 健康检查

# 文档
make swagger-gen      # 生成 Swagger
make swagger-serve    # 启动 Swagger UI
```

---

## 🎯 典型工作流

### 场景 1：日常开发

```bash
# 1. 启动基础设施
docker-compose up -d postgres redis

# 2. 启动 API（终端 1）
make run

# 3. 启动 Worker（终端 2）
make run-worker

# 4. 测试功能
./scripts/dev/core-flow-test.sh
```

### 场景 2：调试问题

```bash
# 1. 启动 Asynqmon 查看队列
make asynqmon

# 2. 查看日志
docker-compose logs -f api worker

# 3. 查看指标
watch 'curl localhost:8080/metrics | head -50'
```

### 场景 3：性能测试

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 启动监控
make asynqmon

# 3. 压测
ab -n 1000 -c 10 http://localhost:8080/api/v1/auth/login

# 4. 观察指标
curl localhost:8080/metrics | grep duration
```

---

## 🔑 关键路径

```
backend/
├── cmd/
│   ├── api/          # API 入口
│   ├── worker/       # Worker 入口
│   └── docs/         # Swagger 文档入口
├── internal/
│   ├── auth/         # 认证模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── tasks.go   # 队列任务
│   │   └── events.go  # 领域事件
│   └── middleware/   # 中间件
└── pkg/
    ├── metrics/      # Prometheus 指标
    └── health/       # 健康检查
```

---

## 🐛 故障排查

### API 无法启动

```bash
# 1. 检查配置
cat backend/configs/development.yaml

# 2. 检查数据库连接
psql -h localhost -U postgres -d ddd_scaffold

# 3. 查看日志
tail -f /tmp/api.log
```

### Worker 不处理任务

```bash
# 1. 检查 Redis 连接
redis-cli ping

# 2. 查看队列
asynq stats --redis-addr=localhost:6379

# 3. 重启 Worker
make run-worker
```

### 测试失败

```bash
# 1. 查看详细错误
./scripts/dev/core-flow-test.sh 2>&1 | tail -50

# 2. 检查服务状态
make health

# 3. 重置数据库
docker-compose down -v
make migrate-up
```

---

## 📞 快速链接

- **API**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Swagger**: http://localhost:8080/swagger
- **Asynqmon**: http://localhost:8080
- **Metrics**: http://localhost:8080/metrics

---

## 🎓 学习路径

1. ✅ **启动服务** - `docker-compose up -d`
2. ✅ **执行测试** - `./scripts/dev/core-flow-test.sh`
3. ✅ **查看队列** - `make asynqmon`
4. ✅ **查看指标** - `curl localhost:8080/metrics`
5. ✅ **查看日志** - `docker-compose logs -f`

---

**祝开发愉快！** 🎉
