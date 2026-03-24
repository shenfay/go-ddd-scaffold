# 事件发布器适配器使用指南

**日期:** 2026-03-23  
**状态:** ✅ 已完成

---

## 📋 背景

为了简化架构并统一删除旧代码，我们采用了**适配器模式**来平滑过渡：

### 问题
- 已删除 `AsynqPublisher` 中间层
- 应用服务中大量引用 `kernel.EventPublisher` 接口
- 直接重构会导致大量编译错误和工作量

### 解决方案
创建 `EventPublisherAdapter`，实现旧的接口但使用新的基础设施。

---

## 🎯 核心思想

```go
// 旧的接口（保持不变）
type EventPublisher interface {
    Publish(ctx context.Context, event kernel.DomainEvent) error
}

// 新的实现（适配器）
type EventPublisherAdapter struct {
    activityRepo  aggregate.ActivityLogRepository  // 新：活动日志
    eventRepo     aggregate.EventLogRepository     // 新：事件日志
    taskPublisher *task_queue.Publisher            // 旧：队列发布
}
```

**关键点：**
- ✅ 保持接口不变（兼容旧代码）
- ✅ 内部实现三重写逻辑（ActivityLog + EventLog + Queue）
- ✅ 应用服务无需修改或只需少量修改

---

## 📦 文件位置

```
backend/internal/infrastructure/eventstore/
├── publisher_adapter.go          ← 新增：适配器实现
├── bus.go                        ← 保留：事件总线
├── repository.go                 ← 保留：事件仓储
└── store.go                      ← 保留：事件存储
```

**已删除：**
- ❌ publisher.go (AsynqPublisher 实现)

---

## 🔧 使用方法

### 方式 1: 在 bootstrap 中配置（推荐）

```go
// backend/internal/bootstrap/module.go 或相关配置文件中

func provideEventPublisher(
    activityRepo aggregate.ActivityLogRepository,
    eventRepo aggregate.EventLogRepository,
    taskPublisher *task_queue.Publisher,
    logger *zap.Logger,
) kernel.EventPublisher {
    return domain_event.NewEventPublisherAdapter(
        activityRepo,
        eventRepo,
        taskPublisher,
        logger,
    )
}
```

### 方式 2: 直接在 main.go 中注入

```go
// backend/cmd/api/main.go

// 创建依赖
activityRepo := repository.NewActivityLogRepository(query)
eventRepo := repository.NewEventLogRepository(query)
taskPublisher := task_queue.NewPublisher(asynqClient)

// 创建适配器
eventPublisher := domain_event.NewEventPublisherAdapter(
    activityRepo,
    eventRepo,
    taskPublisher,
    logger,
)

// 注入到应用服务
authService := auth.NewAuthService(
    uow,
    passwordHasher,
    tokenService,
    eventPublisher,  // ← 使用适配器
    idGenerator,
    logger,
)
```

---

## 🎪 工作流程

当应用服务调用 `eventPublisher.Publish(event)` 时：

```
┌─────────────────────────────────────┐
│ Application Service                 │
│   eventPublisher.Publish(event)     │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│ EventPublisherAdapter               │
│                                     │
│  1. saveActivityLog()              │
│     → activity_logs 表              │
│                                     │
│  2. saveEventLog()                 │
│     → event_log 表                  │
│                                     │
│  3. publishToQueue()               │
│     → Asynq Redis 队列              │
└──────────────┬──────────────────────┘
               ↓
         三个操作并行/串行执行
```

---

## 📊 数据流示例

以用户注册为例：

```go
// 1. 应用服务创建事件
event := userEvent.NewUserRegisteredEvent(
    userID, username, email, status, displayName, 
    registrationIP, tenantID,
)

// 2. 发布事件（对应用服务透明）
err := s.eventPublisher.Publish(ctx, event)

// 3. 适配器内部执行：
//    a) 保存 ActivityLog
//       - action: "USER_REGISTERED"
//       - status: ActivityStatusSuccess
//       - metadata: {email, username, ...}
//       
//    b) 保存 EventLog
//       - aggregate_id: userID
//       - aggregate_type: "UserRegistered"
//       - event_data: {完整事件 JSON}
//       
//    c) 发布到 Asynq 队列
//       - queue: "critical"
//       - payload: {序列化后的事件}
```

---

## ✅ 优势

### 1. 最小化改动
- ✅ 应用服务**无需修改**
- ✅ 保持向后兼容
- ✅ 降低测试成本

### 2. 清晰的职责分离
- ✅ ActivityLog: 业务活动记录（审计）
- ✅ EventLog: 领域事件记录（溯源）
- ✅ Asynq: 异步任务处理

### 3. 易于测试
```go
// 可以单独测试每个组件
func TestActivityLog(t *testing.T) { ... }
func TestEventLog(t *testing.T) { ... }
func TestPublishToQueue(t *testing.T) { ... }

// 也可以测试整体
func TestAdapter(t *testing.T) { ... }
```

### 4. 渐进式迁移
- 第一阶段：使用适配器（当前）
- 第二阶段：逐步替换为直接调用（可选）
- 第三阶段：移除适配器（最终）

---

## 🔍 配置说明

### 数据库表
适配器会使用以下两个新表：

```sql
-- activity_logs: 业务活动记录
CREATE TABLE activity_logs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action VARCHAR(100) NOT NULL,
    status SMALLINT DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    occurred_at TIMESTAMP NOT NULL
);

-- event_log: 领域事件记录
CREATE TABLE event_log (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL
);
```

### 队列配置
适配器会根据事件类型自动选择队列：

| 队列 | 优先级 | 事件示例 |
|------|--------|---------|
| `critical` | 高 | UserRegistered, PaymentCompleted |
| `default` | 中 | 大部分业务事件 |
| `low` | 低 | UserLoggedIn（日志类） |

---

## ⚠️ 注意事项

### 1. 错误处理
适配器采用**宽松的错误处理策略**：
- ActivityLog 保存失败 → 记录日志，不阻断
- EventLog 保存失败 → 记录日志，不阻断
- Queue 发布失败 → **返回错误**（关键流程）

### 2. 性能考虑
三个操作是**顺序执行**的：
```go
saveActivityLog()  // ~5ms
saveEventLog()     // ~5ms
publishToQueue()   // ~10ms
// 总计约 20ms
```

如果需要优化，可以改为并发执行：
```go
// 并发保存日志（可选）
go a.saveActivityLog(ctx, event)
go a.saveEventLog(ctx, event)
// 同步发布队列
return a.publishToQueue(ctx, event)
```

### 3. 事务边界
- ActivityLog 和 EventLog 的保存**不在**业务事务中
- 如果需要在事务中保存，需要调整架构设计

---

## 🚀 下一步计划

### Phase 3A: 完成适配（当前）
- ✅ 创建 EventPublisherAdapter
- ⏳ 更新 DI 配置
- ⏳ 运行测试验证

### Phase 3B: 简化 Handler（后续）
- 将 Handler 接口改为函数类型
- 支持多个处理器注册

### Phase 3C: 清理旧代码（最后）
- 删除旧的 migration 文件
- 删除不再需要的代码

---

## 📖 相关文档

- [Phase 3 重构报告](./phase3-refactoring-report.md) - 详细分析
- [重构进度总结](./event-system-refactoring-summary.md) - 总体进度
- [ActivityLog 使用指南](../guides/activity-log-usage.md) - 待创建
- [EventLog 使用指南](../guides/event-log-usage.md) - 待创建

---

**最后更新:** 2026-03-23  
**状态:** ✅ 可以使用
