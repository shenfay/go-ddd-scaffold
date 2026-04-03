# 事件驱动日志系统集成指南

## 概述

本文档说明如何在 API 和 Worker 服务中集成新的事件驱动日志系统。

---

## 架构变更

### 新增组件

1. **EventBus** (`internal/infra/messaging/event_bus.go`)
   - 基于 Asynq 的事件总线实现
   - 支持队列路由（critical/default）

2. **AuditLogListener** (`internal/listener/audit_log_listener.go`)
   - 订阅认证事件（AUTH.LOGIN.SUCCESS, AUTH.LOGIN.FAILED, SECURITY.ACCOUNT.LOCKED）
   - 转换为审计日志任务并发布到 Worker 队列

3. **Worker Handlers** (`internal/transport/worker/handlers/`)
   - AuditLogHandler: 处理审计日志任务并写入数据库

4. **Redis Token Store** (`internal/infra/redis/token_store.go`)
   - 使用 Redis 存储 Token 和设备信息
   - TTL: 7 天自动过期

---

## API 服务集成 (cmd/api/main.go)

### 步骤 1：导入新包

```go
import (
    // 现有导入...
    
    // 新增
    "github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
    "github.com/shenfay/go-ddd-scaffold/internal/listener"
)
```

### 步骤 2：创建 EventBus

在 `initRedis()` 之后添加：

```go
// 创建 EventBus
eventBus := messaging.NewEventBus(cfg.Redis.Addr, messaging.QueueConfig{
    Critical: "critical",
    Default:  "default",
})
pkglogger.Info("✓ Event Bus initialized")
```

### 步骤 3：创建 Listener

```go
// 创建审计日志监听器
auditLogListener := listener.NewAuditLogListener(eventBus)
_ = auditLogListener // 保持引用，防止被 GC
pkglogger.Info("✓ Audit Log Listener registered")
```

### 步骤 4：设置 AuthService 的 EventBus

找到创建 authService 的代码：

```go
// 原代码
authService := auth.NewService(userRepo, tokenService)

// 修改为
authService := auth.NewService(userRepo, tokenService)
authService.SetEventBus(auth.NewEventBusAdapter(eventBus))
```

---

## Worker 服务集成 (cmd/worker/main.go)

### 步骤 1：导入新包

```go
import (
    // 现有导入...
    
    // 新增
    "github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
    workerHandlers "github.com/shenfay/go-ddd-scaffold/internal/transport/worker/handlers"
)
```

### 步骤 2：创建 Repository 和 Handler

在创建 mux 之前：

```go
// 创建审计日志仓储
auditLogRepo := repository.NewAuditLogRepository(db)
logger.Info("✓ AuditLogRepository initialized")

// 创建审计日志处理器
auditLogHandler := workerHandlers.NewAuditLogHandler(auditLogRepo)
```

### 步骤 3：注册 Handler 到 ServeMux

```go
mux := asynq.NewServeMux()

// 注册审计日志处理器
mux.HandleFunc("audit.log.task", auditLogHandler.ProcessTask)

// 保留原有的活动日志处理器
mux.HandleFunc(constants.ActivityLogTaskType, activityLogHandler.ProcessActivityLogTask)
```

---

## 数据库迁移

运行 SQL 迁移创建新表：

```bash
cd backend

# 运行迁移
migrate -path migrations -database "$DATABASE_URL" up
```

或手动执行：

```bash
psql $DATABASE_URL < migrations/005_create_audit_logs_table.up.sql
psql $DATABASE_URL < migrations/006_create_activity_logs_table.up.sql
```

---

## 验证

### 1. 启动 API 服务

```bash
cd backend/cmd/api
go run main.go
```

应该看到日志：
```
✓ Event Bus initialized
✓ Audit Log Listener registered
```

### 2. 启动 Worker 服务

```bash
cd backend/cmd/worker
go run main.go
```

应该看到日志：
```
✓ AuditLogRepository initialized
```

### 3. 测试登录接口

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 4. 检查审计日志

```sql
SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;
```

---

## 性能监控

### Prometheus Metrics

新增指标：
- `asynq_jobs_processed_total` - 已处理的任务数
- `asynq_queue_size` - 队列大小
- `audit_log_save_duration_seconds` - 审计日志保存耗时

### Grafana 仪表盘

导入新的仪表盘模板（待添加）。

---

## 故障排查

### 问题 1：事件未发布

**症状**：登录成功但 audit_logs 表无记录

**检查**：
1. API 日志中是否有 "Failed to publish UserLoggedInEvent"
2. EventBus 是否正确初始化
3. Listener 是否订阅了事件

### 问题 2：Worker 未消费

**症状**：队列中有任务但 Worker 未处理

**检查**：
1. Worker 日志中是否有错误
2. Handler 是否正确注册到 ServeMux
3. 任务类型名称是否匹配（"audit.log.task"）

### 问题 3：Redis 连接失败

**症状**：启动时报 "Failed to connect to Redis"

**解决**：
```bash
# 检查 Redis 是否运行
redis-cli ping

# 检查配置文件中的 Redis 地址
cat configs/development.yaml | grep redis
```

---

## 回滚方案

如果需要回滚到旧版本：

### 1. 代码回滚

```bash
git checkout HEAD~1
```

### 2. 清理数据

```sql
-- 删除新表
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS activity_logs;

-- 清理 Redis
redis-cli FLUSHDB
```

---

## 下一步

1. 性能压测：验证 P99 < 50ms
2. 监控告警：配置关键指标告警
3. 文档完善：更新 API 文档和运维手册
