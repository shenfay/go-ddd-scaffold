# 监控与可观测性指南

本文档介绍如何查看 DDD Scaffold 项目的队列、事件和各种监控信息。

---

## 📋 目录

1. [启动 API 并测试核心流程](#1-启动-api-并测试核心流程)
2. [查看 Asynq 队列和事件](#2-查看-asynq-队列和事件)
3. [查看 Prometheus 指标](#3-查看-prometheus-指标)
4. [查看应用日志](#4-查看应用日志)
5. [健康检查](#5-健康检查)
6. [Swagger API 文档](#6-swagger-api-文档)

---

## 1. 启动 API 并测试核心流程

### 1.1 启动服务

**方式 A：使用 Docker Compose（推荐）**
```bash
# 一键启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api
docker-compose logs -f worker

# 检查状态
docker-compose ps
```

**方式 B：本地运行**
```bash
# 终端 1 - 启动 API
cd backend
make run

# 终端 2 - 启动 Worker
cd backend
make run-worker
```

### 1.2 执行核心流程测试

```bash
# 赋予执行权限
chmod +x scripts/dev/core-flow-test.sh

# 执行测试（使用默认配置）
./scripts/dev/core-flow-test.sh

# 或者使用自定义参数
./scripts/dev/core-flow-test.sh \
  -u testuser \
  -e testuser@example.com \
  -p MyPassword123

# 交互模式
./scripts/dev/core-flow-test.sh --interactive
```

### 1.3 测试脚本功能清单

✅ **测试流程**：
1. 用户注册（随机用户名和邮箱）
2. 用户登录
3. 获取当前用户信息
4. 获取指定用户信息
5. 刷新 Token（验证令牌轮换）
6. 用户登出
7. 验证 Token 失效（黑名单机制）

✅ **输出内容**：
- 每个步骤的详细响应
- JSON 格式化的结果
- 成功/失败提示
- 关键指标验证

---

## 2. 查看 Asynq 队列和事件

### 2.1 Asynqmon UI（推荐 ⭐）

Asynqmon 是官方的 Web UI，可以直观地查看队列、任务、事件等信息。

#### 安装 Asynqmon

```bash
# 方式 A：通过 Go 安装
make asynqmon-install

# 方式 B：手动安装
go install github.com/hibiken/asynqmon@latest
```

#### 启动 Asynqmon

```bash
# 方式 A：使用 Makefile（默认端口 8080）
make asynqmon

# 方式 B：自定义端口
make asynqmon-port PORT=8081

# 方式 C：直接命令
$HOME/go/bin/asynqmon --redis-addr=localhost:6379

# 方式 D：Docker 运行
docker run --rm -p 8080:8080 hibiken/asynqmon --redis-addr=host.docker.internal:6379
```

#### 访问界面

打开浏览器访问：**http://localhost:8080**

#### 功能特性

✅ **Dashboard**：
- 实时显示各队列的任务数量
- 处理中的任务数
- 错误率统计

✅ **Queues**：
- Critical（高优先级队列）
- Default（默认队列）
- Low（低优先级队列）

✅ **Task Details**：
- 任务 ID
- 任务类型（如 `UserRegistered`、`SendWelcomeEmail`）
- 重试次数
- 错误信息
- 计划执行时间

✅ **Events（事件流）**：
- 任务入队事件
- 任务开始处理
- 任务完成/失败
- 重试事件

### 2.2 命令行查看

```bash
# 安装 asynq CLI
go install github.com/hibiken/asynq/cmd/asynq@latest

# 查看队列统计
asynq stats --redis-addr=localhost:6379

# 查看任务列表
asynq tasks --redis-addr=localhost:6379

# 查看特定队列
asynq tasks --queue=critical --redis-addr=localhost:6379

# 查看任务详情
asynq tasks --id=<TASK_ID> --redis-addr=localhost:6379
```

### 2.3 代码中查看事件处理器

项目中已实现的事件处理器：

```bash
# 查看事件定义
cat backend/internal/auth/events.go

# 查看事件处理器
cat backend/internal/auth/handler.go

# 查看 Worker 端处理器
cat backend/internal/auth/tasks.go
```

**已实现的事件**：
- `UserRegistered` - 用户注册事件
- `UserLoggedIn` - 用户登录事件
- `UserLoggedOut` - 用户登出事件
- `TokenRefreshed` - Token 刷新事件

---

## 3. 查看 Prometheus 指标

### 3.1 访问 Metrics 端点

```bash
# 直接访问
curl http://localhost:8080/metrics

# 保存到文件
curl http://localhost:8080/metrics > metrics.txt

# 实时监控（每秒刷新）
watch -n 1 'curl -s http://localhost:8080/metrics | head -50'
```

### 3.2 关键指标说明

#### HTTP 层指标

```promql
# HTTP 请求总数
http_requests_total{method="POST",path="/api/v1/auth/login",status="200"}

# HTTP 请求耗时（直方图）
http_request_duration_seconds_bucket{le="0.1"}  # < 100ms 的请求数
http_request_duration_seconds_bucket{le="1.0"}  # < 1s 的请求数

# 当前并发请求数
http_requests_in_flight
```

#### 数据库指标

```promql
# DB 查询总数
db_queries_total{query_type="SELECT"}

# DB 查询耗时
db_query_duration_seconds

# DB 连接池状态
db_connections_open
db_connections_idle
```

#### Redis 指标

```promql
# Redis 命令总数
redis_commands_total{command="GET"}

# Redis 操作耗时
redis_command_duration_seconds
```

#### 业务指标

```promql
# 认证尝试次数
auth_attempts_total{result="success"}
auth_attempts_total{result="failure"}

# 活跃用户数
active_users_count

# Token 刷新次数
token_refreshes_total

# 用户注册数
user_registrations_total

# 邮件发送数
emails_sent_total
```

### 3.3 使用 Prometheus 查询

启动 Prometheus 后，可以使用 PromQL 查询：

```promql
# 过去 5 分钟的请求速率
rate(http_requests_total[5m])

# 平均响应时间
avg(http_request_duration_seconds)

# 错误率
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))

# 认证成功率
sum(auth_attempts_total{result="success"}) / sum(auth_attempts_total)
```

### 3.4 Grafana 仪表盘（可选）

如果需要可视化展示，可以配置 Grafana：

1. 添加 Prometheus 数据源
2. 导入或创建 Dashboard
3. 配置告警规则

---

## 4. 查看应用日志

### 4.1 开发环境日志（Console 格式）

```bash
# API 日志（本地运行会自动输出到终端）
make run

# Worker 日志
make run-worker
```

### 4.2 Docker 容器日志

```bash
# 查看 API 日志
docker-compose logs -f api

# 查看 Worker 日志
docker-compose logs -f worker

# 查看最近 100 行
docker-compose logs --tail=100 api

# 导出日志到文件
docker-compose logs api > api.log
```

### 4.3 日志级别控制

```bash
# 修改配置文件
backend/configs/development.yaml

logger:
  level: debug  # debug, info, warn, error
  format: console  # json, console
```

### 4.4 结构化日志（JSON 格式）

生产环境建议使用 JSON 格式：

```yaml
logger:
  level: info
  format: json
  output_path: /var/log/app.log
```

查看日志：

```bash
# 使用 jq 格式化查看
tail -f /var/log/app.log | jq .

# 搜索特定级别
grep '"level":"error"' /var/log/app.log
```

---

## 5. 健康检查

### 5.1 快速健康检查

```bash
# 简单检查
curl http://localhost:8080/health

# 使用 Makefile
make health
```

### 5.2 详细健康检查

```bash
# 完整健康检查（包含 DB 和 Redis）
curl http://localhost:8080/health | jq .

# 预期输出：
# {
#   "status": "healthy",
#   "checks": {
#     "database": {
#       "status": "ok",
#       "response_time_ms": "<10ms"
#     },
#     "redis": {
#       "status": "ok",
#       "response_time_ms": "<10ms",
#       "ping_result": "PONG"
#     }
#   }
# }
```

### 5.3 Liveness 探针（Kubernetes）

```bash
# 只检查进程是否存活
curl http://localhost:8080/health/live
```

### 5.4 Readiness 探针（Kubernetes）

```bash
# 检查依赖是否就绪
curl http://localhost:8080/health/ready
```

---

## 6. Swagger API 文档

### 6.1 生成 Swagger 文档

```bash
# 生成文档
make swagger-gen

# 生成并启动服务
make swagger-serve
```

### 6.2 访问 Swagger UI

打开浏览器访问：**http://localhost:8080/swagger/index.html**

### 6.3 Swagger 功能

✅ **API 文档**：
- 查看所有端点的详细说明
- 请求/响应示例
- 参数说明

✅ **在线测试**：
- 直接在浏览器中测试 API
- 自动填充 JWT Token
- 查看实际响应

✅ **模型定义**：
- 查看数据结构定义
- 了解字段类型和约束

---

## 🔧 综合监控示例

### 场景 1：调试用户注册流程

```bash
# 1. 启动 Asynqmon 查看队列
make asynqmon

# 2. 在新终端执行注册测试
./scripts/dev/core-flow-test.sh

# 3. 观察 Asynqmon 中的任务变化
# - 应该看到 UserRegistered 任务入队
# - 然后被 Worker 处理
# - 最终完成

# 4. 同时查看 API 日志
docker-compose logs -f api | grep -i register

# 5. 查看 Worker 日志
docker-compose logs -f worker | grep -i registered
```

### 场景 2：性能问题排查

```bash
# 1. 实时监控指标
watch -n 1 'curl -s http://localhost:8080/metrics | grep http_request_duration'

# 2. 压测（使用 ab 或 wrk）
ab -n 1000 -c 10 http://localhost:8080/api/v1/auth/login

# 3. 观察 Prometheus 指标变化
# - 请求速率
# - 响应时间分布
# - 错误率

# 4. 查看慢查询日志
docker-compose logs api | grep "slow query"
```

### 场景 3：错误追踪

```bash
# 1. 查看错误日志
docker-compose logs api | grep -i error
docker-compose logs worker | grep -i error

# 2. 查看错误指标
curl http://localhost:8080/metrics | grep errors

# 3. 在 Asynqmon 中查看失败任务
# - 访问 Failed 队列
# - 查看错误堆栈
# - 分析重试历史
```

---

## 📊 监控工具对比

| 工具 | 用途 | 启动命令 | 访问地址 |
|------|------|----------|----------|
| **Asynqmon** | 队列和任务监控 | `make asynqmon` | http://localhost:8080 |
| **Prometheus** | 指标收集 | Docker Compose | http://localhost:9090 |
| **Grafana** | 可视化仪表盘 | Docker Compose | http://localhost:3000 |
| **Swagger UI** | API 文档和测试 | `make swagger-serve` | http://localhost:8080/swagger |
| **健康检查** | 服务状态 | - | http://localhost:8080/health |

---

## 🎯 快速开始

**推荐的监控组合**：

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 启动 Asynqmon（新终端）
make asynqmon

# 3. 执行测试流程
./scripts/dev/core-flow-test.sh

# 4. 浏览器访问
# - Asynqmon: http://localhost:8080
# - Swagger: http://localhost:8080/swagger
# - Health: http://localhost:8080/health
```

---

## ✅ 总结

### 查看队列和事件
- **工具**: Asynqmon UI（最直观）
- **命令**: `make asynqmon`
- **地址**: http://localhost:8080

### 查看指标
- **端点**: `/metrics`
- **命令**: `curl http://localhost:8080/metrics`
- **工具**: Prometheus + Grafana（可选）

### 查看日志
- **开发**: Console 输出
- **生产**: JSON 格式 + ELK/Loki

### 测试流程
- **脚本**: `scripts/dev/core-flow-test.sh`
- **覆盖**: 注册、登录、刷新、登出全流程

**所有监控工具都已就绪！** 🎊
