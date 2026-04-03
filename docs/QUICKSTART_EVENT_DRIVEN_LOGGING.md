# 事件驱动日志系统 - 快速启动指南

## 📋 概述

本文档提供事件驱动日志系统的完整启动和测试流程，预计耗时 **10 分钟**。

---

## 🎯 前置条件

确保以下服务已安装并运行：

- ✅ PostgreSQL 15+
- ✅ Redis 7.x
- ✅ Go 1.21+
- ✅ migrate 工具（可选）

---

## 🚀 快速启动（3 步）

### **步骤 1：执行数据库迁移**

#### 方式 A：使用迁移脚本（推荐）

```bash
cd /Users/shenfay/Projects/ddd-scaffold

# 添加执行权限
chmod +x scripts/dev/migrate-logging.sh

# 执行迁移
./scripts/dev/migrate-logging.sh
```

#### 方式 B：手动执行 SQL

```bash
cd backend/migrations

# 执行审计日志表迁移
psql $DATABASE_URL < 005_create_audit_logs_table.up.sql

# 执行活动日志表迁移
psql $DATABASE_URL < 006_create_activity_logs_table.up.sql
```

#### 验证迁移成功

```sql
-- 检查表是否存在
\dt audit_logs
\dt activity_logs

-- 查看表结构
\d audit_logs
```

---

### **步骤 2：启动 API 服务**

```bash
cd backend/cmd/api

# 启动服务
go run main.go
```

**预期输出**：
```
✓ Asynq client initialized
✓ Event Bus initialized          ← 关键：EventBus 启动成功
✓ Audit Log Listener registered  ← 关键：监听器注册成功
[GIN-debug] Listening and serving HTTP on :8080
```

**如果没有看到 "Event Bus initialized"，请检查**：
1. `cmd/api/main.go` 是否包含 EventBus 初始化代码
2. Redis 服务是否正常运行

---

### **步骤 3：启动 Worker 服务**

```bash
cd backend/cmd/worker

# 启动服务
go run main.go
```

**预期输出**：
```
✓ Database connection established
✓ AuditLogRepository and AuditLogHandler initialized  ← 关键
✓ Registered audit log handler for type: audit.log.task
✓ Asynq server created with concurrency=10
```

**如果没有看到相关日志，请检查**：
1. `cmd/worker/main.go` 是否注册了 Handler
2. 数据库连接是否正常

---

## 🧪 功能测试

### **测试 1：登录接口**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}'
```

**预期响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": "01H...",
      "email": "test@example.com"
    },
    "access_token": "eyJ...",
    "refresh_token": "dXNlcjoxMjM0NTY3ODkw",
    "expires_in": 3600
  }
}
```

**关键点**：
- ✅ 响应时间应该 < 50ms（P99）
- ✅ 返回 access_token 和 refresh_token

---

### **测试 2：检查审计日志**

#### 方式 A：使用测试脚本（推荐）

```bash
chmod +x scripts/dev/test-logging.sh
./scripts/dev/test-logging.sh
```

#### 方式 B：手动查询数据库

```sql
-- 查看最新的审计日志
SELECT 
    id,
    user_id,
    email,
    action,
    status,
    ip,
    created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT 10;
```

**预期结果**：
```
 id | user_id | email | action | status | ip | created_at
----|---------|-------|--------|--------|----|------------
 01H... | 01H... | test@example.com | AUTH.LOGIN.SUCCESS | SUCCESS | ::1 | 2026-04-03 10:00:00+00
```

**关键字段说明**：
- `action`: 应该是 `AUTH.LOGIN.SUCCESS`
- `status`: 应该是 `SUCCESS`
- `metadata`: 包含 IP、User-Agent、设备信息等

---

### **测试 3：Worker 处理验证**

查看 Worker 日志：

```bash
tail -f backend/logs/worker.log
```

**预期日志**：
```
[INFO] Processing audit.log.task task
[INFO] Audit log saved successfully: action=AUTH.LOGIN.SUCCESS user_id=01H...
```

如果看到错误日志，请检查：
1. 数据库连接是否正常
2. audit_logs 表是否存在
3. 字段类型是否匹配

---

## 📊 性能测试（可选）

### 使用 wrk 进行压测

```bash
# 安装 wrk
brew install wrk  # macOS
# 或
apt-get install wrk  # Linux

# 压测登录接口
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}'
```

**预期结果**：
```
Running 30s test @ http://localhost:8080/api/v1/auth/login
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    28.50ms   12.34ms  89.12ms   75.2%
    Req/Sec     1.23k    123.45     1.50k    85.0%
  37020 requests in 30.10s, 5.23MB read
Requests/sec: 1230.12
Transfer/sec: 177.45KB
```

**关键指标**：
- ✅ **P99 Latency**: < 50ms
- ✅ **Requests/sec**: > 1000 QPS
- ✅ **Error Rate**: < 0.1%

---

## 🔍 故障排查

### 问题 1：登录接口响应慢（> 200ms）

**可能原因**：
1. EventBus 未正确初始化
2. Redis 连接超时
3. 队列配置不当

**解决方案**：
```bash
# 检查 Redis 连接
redis-cli ping

# 检查 API 日志
tail -f backend/logs/api.log | grep "Event Bus"

# 确认 EventBus 已初始化
grep "✓ Event Bus initialized" backend/logs/api.log
```

---

### 问题 2：审计日志未记录

**可能原因**：
1. Listener 未订阅事件
2. Worker 未消费队列
3. 数据库写入失败

**排查步骤**：

```bash
# 1. 检查 Listener 是否注册
grep "Audit Log Listener registered" backend/logs/api.log

# 2. 检查 Worker 是否启动
ps aux | grep worker

# 3. 查看 Worker 日志
tail -f backend/logs/worker.log | grep "audit.log.task"

# 4. 检查队列中是否有任务
redis-cli
> KEYS asynq:*
> LLEN asynq:queue:default
```

---

### 问题 3：Redis 连接失败

**错误信息**：
```
Failed to connect to Redis: dial tcp [::1]:6379: connect: connection refused
```

**解决方案**：

```bash
# 1. 启动 Redis
redis-server

# 2. 或使用 Docker
docker run -d -p 6379:6379 redis:7

# 3. 验证连接
redis-cli ping  # 应返回 PONG

# 4. 检查配置文件
cat backend/configs/development.yaml | grep redis
```

---

## 📈 监控指标

### Prometheus Metrics

在 Grafana 中添加以下指标：

```promql
# 登录接口 P99 延迟
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{endpoint="/api/v1/auth/login"}[5m]))

# Asynq 队列大小
asynq_queue_size{queue="critical"}

# 审计日志保存耗时
audit_log_save_duration_seconds
```

---

## 🎯 验收标准

完成以下检查清单：

- [ ] **数据库迁移完成**
  - [ ] audit_logs 表创建成功
  - [ ] activity_logs 表创建成功
  - [ ] 所有索引已创建

- [ ] **API 服务正常**
  - [ ] EventBus 初始化成功
  - [ ] AuditLogListener 注册成功
  - [ ] 登录接口响应 < 50ms（P99）

- [ ] **Worker 服务正常**
  - [ ] AuditLogHandler 注册成功
  - [ ] 能够消费 audit.log.task 任务
  - [ ] 日志正常写入数据库

- [ ] **功能验证通过**
  - [ ] 登录成功触发审计日志
  - [ ] 登录失败触发审计日志
  - [ ] 账户锁定触发审计日志

- [ ] **性能达标**
  - [ ] P99 < 50ms
  - [ ] 吞吐量 > 1000 QPS
  - [ ] 错误率 < 0.1%

---

## 📝 下一步

### 立即可做

1. **查看完整的审计日志**
```sql
SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 20;
```

2. **查看 Worker 处理日志**
```bash
tail -n 100 backend/logs/worker.log
```

3. **测试其他场景**
   - 登录失败（错误密码）
   - 账户被锁定
   - Token 刷新

### 后续优化

1. **实现 ActivityLogListener**（可选）
   - 记录用户行为日志
   - 用于产品分析

2. **配置监控告警**
   - Prometheus + Grafana
   - 关键指标告警

3. **补充单元测试**
   - EventBus 测试
   - Listener 测试
   - Worker Handler 测试

---

## 📚 参考文档

- [`EVENT_DRIVEN_INTEGRATION_GUIDE.md`](./EVENT_DRIVEN_INTEGRATION_GUIDE.md) - 详细集成指南
- [`IMPLEMENTATION_SUMMARY.md`](./IMPLEMENTATION_SUMMARY.md) - 实施总结
- [`REFACTORING_IMPLEMENTATION_PLAN.md`](./REFACTORING_IMPLEMENTATION_PLAN.md) - 完整方案

---

## 🆘 获取帮助

如果遇到其他问题：

1. **查看日志文件**
   - `backend/logs/api.log`
   - `backend/logs/worker.log`

2. **检查 Git 提交历史**
```bash
git log --oneline -10
```

3. **回滚到旧版本**
```bash
git checkout feature/refactor-rebuild~1
```

---

**状态**: ✅ 核心功能已完成  
**最后更新**: 2026-04-03  
**维护者**: 后端团队
