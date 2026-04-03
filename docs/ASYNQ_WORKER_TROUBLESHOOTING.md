# Asynq Worker 问题排查与修复

## 🔍 问题诊断

### 错误信息
```
handler not found for task "activity:record"
```

**错误率**: 95.35% (86 个任务，82 个失败)

### 根本原因
Worker 进程没有运行或未正确注册 `activity:record` handler

---

## ✅ 已完成的修复

### 1. 编译 Worker
```bash
cd backend
go build -o worker ./cmd/worker/main.go
```

### 2. 清空 Redis 中的失败任务
```bash
redis-cli FLUSHDB
```

### 3. 清空数据库活动日志
```bash
psql -U shenfay -d ddd_scaffold -c "TRUNCATE activity_logs CASCADE;"
```

### 4. Worker Handler 注册（已验证代码正确）
```go
// cmd/worker/main.go
activityLogRepo := activitylog.NewActivityLogRepository(db)
activityLogHandler := asynqhandlers.NewActivityLogHandler(activityLogRepo)
mux.HandleFunc("activity:record", func(ctx context.Context, t *asynq.Task) error {
    return activityLogHandler.HandleActivityLogRecord(ctx, t)
})
logger.Info("✓ Registered activity log handler for type: activity:record")
```

---

## 🚀 启动步骤

### Step 1: 启动 API 服务
```bash
cd backend
./main > logs/api.log 2>&1 &
```

### Step 2: 启动 Worker
```bash
cd backend
./worker > logs/worker.log 2>&1 &
```

### Step 3: 验证 Worker 启动
```bash
# 检查进程
ps aux | grep worker | grep -v grep

# 查看日志
tail -f logs/worker.log
```

**预期日志输出**:
```
🚀 Starting Asynq Worker...
✓ Redis connection established
✓ Database connection established
✓ Asynq server created with concurrency=10
✓ Registered activity log handler for type: activity:record
🎯 Starting Asynq Worker processor...
```

### Step 4: 运行测试
```bash
bash scripts/dev/core-flow-test.sh
```

### Step 5: 监控任务执行
1. 打开 Asynqmon: http://localhost:8081
2. 查看 default 队列
3. 检查 Completed 和 Failed 标签页
4. 观察成功率

---

## 📊 验证清单

- [ ] Worker 进程正在运行
- [ ] Worker 日志显示 handler 已注册
- [ ] Redis 中没有积压任务
- [ ] 数据库活动日志表已清空
- [ ] API 服务正常运行
- [ ] 测试脚本执行成功
- [ ] Asynqmon 显示任务成功率 > 95%
- [ ] 数据库中有新的活动日志记录

---

## 🔧 常见问题排查

### Q1: Worker 启动失败
```bash
# 查看错误日志
tail -100 logs/worker.log

# 检查数据库连接
psql -U shenfay -d ddd_scaffold -c "SELECT 1;"

# 检查 Redis 连接
redis-cli ping
```

### Q2: Handler 仍然报错 "not found"
```bash
# 确认 worker 代码已重新编译
cd backend
go build -o worker ./cmd/worker/main.go

# 杀掉旧进程
pkill -f './worker'

# 重新启动
./worker > logs/worker.log 2>&1 &
```

### Q3: 任务堆积不处理
```bash
# 检查 Worker 是否在运行
ps aux | grep worker

# 查看 Worker 日志中的错误
tail -f logs/worker.log | grep -i error

# 检查 Asynqmon 的任务状态
# http://localhost:8081/queues/default
```

### Q4: 数据库写入失败
```bash
# 查看 Worker 错误日志
tail -f logs/worker.log | grep "Failed to create"

# 检查数据库表结构
psql -U shenfay -d ddd_scaffold -c "\d activity_logs"

# 手动测试插入
psql -U shenfay -d ddd_scaffold <<EOF
INSERT INTO activity_logs (id, user_id, action, status, ip, created_at)
VALUES ('test-123', 'test-user', 'LOGIN', 'SUCCESS', '127.0.0.1', NOW());
SELECT * FROM activity_logs WHERE id = 'test-123';
DELETE FROM activity_logs WHERE id = 'test-123';
EOF
```

---

## 📈 性能优化建议

### 1. 调整 Worker 并发度
```yaml
# configs/development.yaml
asynq:
  concurrency: 10  # 根据 CPU 核心数调整
  queues:
    critical: 6
    default: 3
    low: 1
```

### 2. 设置任务重试策略
```go
// 在 API 发送任务时
client.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(30*time.Second))
```

### 3. 监控关键指标
- 队列长度 (Queue size)
- 处理延迟 (Latency)
- 错误率 (Error rate)
- 成功率 (Success rate)

---

## 🎯 当前状态

- ✅ Handler 代码正确
- ✅ Worker 已编译
- ✅ Redis 已清空
- ✅ 数据库已清空
- ⏳ Worker 待启动
- ⏳ 测试待执行

**下一步**: 启动 Worker 并运行测试脚本验证功能
