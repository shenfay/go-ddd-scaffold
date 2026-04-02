# 本地开发环境指南（不使用 Docker）

本文档介绍如何在已安装 PostgreSQL 和 Redis 的本地环境中运行 DDD Scaffold。

---

## 📋 前置条件

### 必需工具

确保已安装以下工具：

```bash
# Go 1.25+
go version

# PostgreSQL 15+
psql --version

# Redis 7+
redis-cli --version

# 其他工具
curl --version
jq --version
```

如果未安装，可以使用 Homebrew 安装：

```bash
brew install go postgresql@15 redis curl jq
```

---

## 🔧 配置说明

### 配置文件位置

```
backend/
├── configs/
│   ├── .env                    # 环境变量参考（不自动加载）
│   └── development.yaml        # 实际使用的配置 ✅
```

### 当前配置信息

**数据库配置** (`development.yaml`):
```yaml
database:
  host: localhost
  port: 5432
  name: ddd_scaffold
  user: shenfay
  password: postgres
  ssl_mode: disable
```

**Redis 配置**:
```yaml
redis:
  addr: localhost:6379
  password: ""
  db: 0
```

**JWT 配置**:
```yaml
jwt:
  secret: smLuhX79IBihMgtVucmefmdP8Gt8hshA
  access_expire: 30m
  refresh_expire: 168h  # 7 天
```

---

## 🚀 启动方式

### 方式一：一键启动脚本（推荐 ⭐）

```bash
# 进入项目根目录
cd /Users/shenfay/Projects/ddd-scaffold

# 赋予执行权限
chmod +x scripts/dev/start.sh

# 启动所有服务（包括监控）
./scripts/dev/start.sh --all

# 或只启动 API 和 Worker
./scripts/dev/start.sh

# 只启动 API + Asynqmon
./scripts/dev/start.sh --monitor

# 停止所有服务
./scripts/dev/start.sh --clean
```

**脚本功能**：
- ✅ 自动检查 PostgreSQL 和 Redis 连接
- ✅ 启动 API 服务（后台运行）
- ✅ 启动 Worker 服务（后台运行）
- ✅ 可选启动 Asynqmon 监控
- ✅ 可选启动 Swagger UI
- ✅ 保存进程 ID 到文件
- ✅ 彩色输出和详细日志

---

### 方式二：手动启动

#### 1. 检查依赖服务

```bash
# 检查 PostgreSQL
psql -h localhost -U shenfay -d ddd_scaffold -c "SELECT 1"

# 检查 Redis
redis-cli ping
```

#### 2. 启动 API 服务

```bash
cd backend

# 方式 A：前台运行（可以看到实时日志）
make run

# 方式 B：后台运行
nohup make run > /tmp/api.log 2>&1 &

# 查看日志
tail -f /tmp/api.log
```

#### 3. 启动 Worker 服务

```bash
cd backend

# 方式 A：前台运行
make run-worker

# 方式 B：后台运行
nohup make run-worker > /tmp/worker.log 2>&1 &

# 查看日志
tail -f /tmp/worker.log
```

#### 4. 启动 Asynqmon（可选）

```bash
# 安装
make asynqmon-install

# 启动
make asynqmon

# 访问 http://localhost:8080
```

---

## 🧪 测试验证

### 健康检查

```bash
# 简单检查
curl http://localhost:8080/health

# 详细检查
curl http://localhost:8080/health | jq .

# 使用 Makefile
make health
```

**预期输出**:
```json
{
  "status": "healthy",
  "checks": {
    "database": {
      "status": "ok",
      "response_time_ms": "<10ms"
    },
    "redis": {
      "status": "ok",
      "response_time_ms": "<10ms",
      "ping_result": "PONG"
    }
  }
}
```

### 执行核心流程测试

```bash
# 进入项目根目录
cd /Users/shenfay/Projects/ddd-scaffold

# 执行测试
./scripts/dev/core-flow-test.sh

# 查看测试结果
# 应该看到：注册 → 登录 → 获取用户信息 → 刷新 Token → 登出 → Token 失效
```

---

## 📊 监控和查看

### 查看 Asynq 队列和事件

```bash
# 启动 Asynqmon
make asynqmon

# 浏览器访问
open http://localhost:8080

# 可以查看：
# - 队列状态（Critical, Default, Low）
# - 任务列表和处理情况
# - 事件流（入队、处理中、完成）
```

### 查看 Prometheus 指标

```bash
# 访问 metrics 端点
curl http://localhost:8080/metrics

# 实时监控
watch -n 1 'curl -s http://localhost:8080/metrics | head -50'
```

**关键指标**:
```promql
# HTTP 请求数
http_requests_total

# 响应时间
http_request_duration_seconds

# 认证尝试
auth_attempts_total{result="success"}

# 用户注册
user_registrations_total
```

### 查看应用日志

```bash
# API 日志
tail -f /tmp/api.log

# Worker 日志
tail -f /tmp/worker.log

# 搜索错误
grep -i error /tmp/api.log
```

---

## 🛠️ Makefile 常用命令

```bash
# 开发
make run              # 启动 API（前台）
make run-worker       # 启动 Worker（前台）

# 构建
make build            # 编译 API
make build-worker     # 编译 Worker
make clean            # 清理构建产物

# 测试
make test             # 运行所有测试
make test-short       # 运行快速测试
make coverage         # 生成覆盖率报告

# 代码质量
make fmt              # 格式化代码
make vet              # 运行 go vet
make lint             # 运行 linter

# 数据库
make migrate-up       # 执行迁移
make migrate-down     # 回滚迁移

# 监控
make asynqmon-install    # 安装 asynqmon
make asynqmon            # 启动监控 UI
make health              # 健康检查

# 文档
make swagger-gen         # 生成 Swagger 文档
make swagger-serve       # 启动 Swagger UI
```

---

## 🐛 故障排查

### 问题 1: PostgreSQL 连接失败

**症状**:
```
failed to connect to `user=shenfay database=ddd_scaffold`
```

**解决方案**:

```bash
# 1. 检查 PostgreSQL 是否运行
brew services list | grep postgresql

# 2. 启动 PostgreSQL
brew services start postgresql@15

# 3. 检查数据库是否存在
psql -h localhost -U postgres -c "\l" | grep ddd_scaffold

# 4. 如果不存在，创建数据库
psql -h localhost -U postgres << EOF
CREATE DATABASE ddd_scaffold OWNER shenfay;
EOF

# 5. 检查用户权限
psql -h localhost -U postgres -d ddd_scaffold -c "\du"
```

---

### 问题 2: Redis 连接失败

**症状**:
```
dial tcp [::1]:6379: connect: connection refused
```

**解决方案**:

```bash
# 1. 检查 Redis 是否运行
brew services list | grep redis

# 2. 启动 Redis
brew services start redis

# 3. 测试连接
redis-cli ping
# 应该返回 PONG
```

---

### 问题 3: 端口被占用

**症状**:
```
bind: address already in use
```

**解决方案**:

```bash
# 1. 查找占用端口的进程
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 2. 停止占用进程
kill -9 <PID>

# 3. 或者修改配置文件中的端口
vim backend/configs/development.yaml
```

---

### 问题 4: 迁移失败

**症状**:
```
relation "user_pos" does not exist
```

**解决方案**:

```bash
# 执行数据库迁移
make migrate-up

# 或手动执行
go run ./cmd/cli migrate up

# 验证表已创建
psql -h localhost -U shenfay -d ddd_scaffold -c "\dt"
```

---

### 问题 5: API 启动后立即退出

**解决方案**:

```bash
# 1. 查看详细错误
tail -50 /tmp/api.log

# 2. 检查配置是否正确
cat backend/configs/development.yaml

# 3. 检查 JWT Secret 是否设置
# 必须至少 32 个字符

# 4. 重新执行迁移
make migrate-down
make migrate-up
```

---

## 📝 管理后台进程

### 查看运行中的进程

```bash
# 查看所有相关进程
ps aux | grep "go run ./cmd"

# 或使用保存的 PID 文件
cat /tmp/ddd-scaffold-pids.txt
```

### 停止服务

```bash
# 方式 A：使用脚本
./scripts/dev/start.sh --clean

# 方式 B：手动停止
pkill -f "go run ./cmd/api"
pkill -f "go run ./cmd/worker"
pkill -f "asynqmon"

# 方式 C：使用 kill 命令
kill $(cat /tmp/ddd-scaffold-pids.txt | grep API_PID | cut -d= -f2)
```

---

## 🎯 典型工作流

### 日常开发流程

```bash
# 早上开始工作
cd /Users/shenfay/Projects/ddd-scaffold

# 1. 启动所有服务
./scripts/dev/start.sh --all

# 2. 验证启动成功
make health

# 3. 执行功能测试
./scripts/dev/core-flow-test.sh

# 4. 开始编码...
```

### 调试问题流程

```bash
# 1. 启动 Asynqmon 查看队列
make asynqmon

# 2. 在新终端查看日志
tail -f /tmp/api.log
tail -f /tmp/worker.log

# 3. 执行测试触发问题
./scripts/dev/core-flow-test.sh

# 4. 观察 Asynqmon 和日志
```

### 下班收工流程

```bash
# 停止所有服务
./scripts/dev/start.sh --clean

# 或直接关闭终端（进程会自动结束）
```

---

## 📖 配置文件详解

### development.yaml 完整结构

```yaml
server:
  port: 8080                    # API 端口
  mode: debug                   # debug, release, test
  read_timeout: 30s             # 读取超时
  write_timeout: 30s            # 写入超时
  idle_timeout: 60s             # 空闲超时

database:
  host: localhost               # 数据库主机
  port: 5432                    # 数据库端口
  name: ddd_scaffold            # 数据库名
  user: shenfay                 # 用户名
  password: postgres            # 密码 ⚠️
  ssl_mode: disable             # SSL 模式
  max_open_conns: 25            # 最大连接数
  max_idle_conns: 5             # 最大空闲连接
  conn_max_lifetime: 5m         # 连接最大生命周期

redis:
  addr: localhost:6379          # Redis 地址
  password: ""                  # Redis 密码
  db: 0                         # Redis 数据库编号
  pool_size: 10                 # 连接池大小

jwt:
  secret: smLuhX79IBihMgtVucmefmdP8Gt8hshA  # JWT 密钥 ⚠️
  access_expire: 30m            # Access Token 过期时间
  refresh_expire: 168h          # Refresh Token 过期时间 (7 天)
  issuer: go-ddd-scaffold       # 签发者

asynq:
  addr: localhost:6379          # Redis 地址（Asynq 使用）
  concurrency: 10               # 并发处理数
  queues:
    critical: 6                 # 高优先级队列 worker 数
    default: 3                  # 默认队列 worker 数
    low: 1                      # 低优先级队列 worker 数

logger:
  level: debug                  # 日志级别
  format: console               # json, console
  output_path: stdout           # 输出路径
```

---

## 🔒 安全提示

⚠️ **重要**：

1. **不要提交敏感配置**
   ```bash
   # .gitignore 已包含
   backend/configs/.env
   backend/configs/*.local.yaml
   ```

2. **生产环境必须修改**
   - JWT Secret（使用随机生成的长字符串）
   - 数据库密码
   - Redis 密码

3. **使用环境变量**
   ```bash
   # 可以通过环境变量覆盖配置
   export APP_DATABASE_PASSWORD=your_password
   export APP_JWT_SECRET=your_secret
   ```

---

## 📚 相关文档

- [监控指南](MONITORING_GUIDE.md) - 详细的监控和可观测性文档
- [快速参考](QUICK_REFERENCE.md) - 快速命令参考
- [Docker Compose 指南](DOCKER_COMPOSE_GUIDE.md) - 容器化部署指南

---

## ✅ 总结

**本地开发环境的优势**：

✅ 无需学习 Docker  
✅ 直接使用本地工具链  
✅ 调试更方便  
✅ 性能更好（无容器开销）  

**启动步骤**：

1. ✅ 确保 PostgreSQL 和 Redis 已安装并运行
2. ✅ 配置 `backend/configs/development.yaml`
3. ✅ 运行 `./scripts/dev/start.sh --all`
4. ✅ 访问 http://localhost:8080/health 验证
5. ✅ 开始开发！

**祝开发愉快！** 🎉
