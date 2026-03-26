# 事件系统设计规范

## 📋 核心设计原则

### 1. 领域事件 vs 活动日志的本质区别

| 维度 | 领域事件 (Domain Events) | 活动日志 (Activity Logs) |
|------|------------------------|------------------------|
| **目的** | 通知其他模块业务状态变化 | 审计、安全、合规 |
| **性质** | 业务概念 | 技术概念 |
| **处理方式** | 异步、可重试 | 同步、必须可靠 |
| **失败影响** | 用户体验下降 | 合规风险、审计缺失 |
| **事务一致性** | 最终一致性 | 强一致性（事务内） |
| **示例** | UserRegistered, OrderPlaced | 登录日志、操作记录 |

### 2. 什么应该是领域事件？

**✅ 正确的领域事件特征：**
```go
// ✅ UserRegistered - 业务状态变化，需要触发多个副作用
UserRegistered → {
    - 发送欢迎邮件（异步）
    - 初始化用户统计（异步）
    - 发送注册通知给管理员（异步）
}

// ✅ UserPasswordChanged - 安全相关，需要通知用户
UserPasswordChanged → {
    - 发送安全通知邮件（异步）
}

// ✅ OrderPlaced - 跨模块协作
OrderPlaced → {
    - 扣减库存（异步）
    - 通知商家（异步）
    - 生成配送单（异步）
}
```

**❌ 不应该是领域事件：**
```go
// ❌ ActivityLogWritten - 这只是日志记录，不是业务事件
// ❌ UserLoginSuccess - 这应该直接记录到 activity_logs
// ❌ AuditTrailCreated - 这是审计，不是业务
```

---

## 🏗️ 系统架构

### 完整的事件处理流程

```
┌─────────────────────────────────────────────────────────┐
│ 应用层 (Application Layer)                              │
│                                                         │
│  UseCase:                                               │
│  1. 事务开始                                            │
│  2. user := domain.NewUser(...)  ← 产生领域事件        │
│  3. repo.Save(user)              ← 保存聚合根          │
│  4. activityLogRepo.Save(log)    ← ⚠️ 直接写日志！    │
│  5. eventPublisher.Publish(event)← 保存到 outbox       │
│  6. 事务提交                                            │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│ 基础设施层 (Infrastructure Layer)                       │
│                                                         │
│  ┌─────────────────┐  ┌──────────────────┐             │
│  │ ActivityLogRepo │  │ OutboxProcessor  │             │
│  │ - 直接写入      │  │ - 轮询 outbox    │             │
│  │ - 事务内完成    │  │ - 发布到 Asynq   │             │
│  └─────────────────┘  └────────┬─────────┘             │
│                                │                        │
│                                ↓                        │
│                     ┌──────────────────────┐            │
│                     │  Asynq Message Queue │            │
│                     └──────────┬───────────┘            │
│                                │                        │
│                                ↓                        │
│                     ┌──────────────────────┐            │
│                     │  Asynq Workers       │            │
│                     │  ┌────────────────┐  │            │
│                     │  │EmailSubscriber │  │            │
│                     │  │- 发送邮件      │  │            │
│                     │  └────────────────┘  │            │
│                     │  ┌────────────────┐  │            │
│                     │  │StatsSubscriber │  │            │
│                     │  │- 初始化统计    │  │            │
│                     │  └────────────────┘  │            │
│                     └──────────────────────┘            │
└─────────────────────────────────────────────────────────┘
```

---

## 💡 实现模式

### 1. UseCase 中直接写入 ActivityLog

```go
type RegisterUserUseCase struct {
    uow             UnitOfWorkWithEvents
    registrationSvc *RegistrationService
    logWriter       *ActivityLogWriter  // ← 新增
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) error {
    return uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
        // 1. 创建用户（产生领域事件）
        user, err := uc.registrationSvc.Register(ctx, cmd)
        if err != nil {
            return err
        }

        // 2. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
        if err := uc.logWriter.WriteSuccess(
            ctx,
            user.ID().Int64(),
            aggregate.ActivityUserRegistered,
            map[string]interface{}{
                "username": user.Username(),
                "email":    user.Email(),
            },
        ); err != nil {
            return err  // ← 失败会回滚整个事务
        }

        // 3. 注册聚合根以自动发布领域事件
        uc.uow.TrackAggregate(user)

        // 4. 保存用户
        return uc.uow.UserRepository().Save(ctx, user)
    })
}
```

### 2. ActivityLogWriter 工具类

```go
// 创建活动日志写入器
logWriter := application.NewActivityLogWriter(activityLogRepo, logger)

// 使用方式 1：写入成功的活动
err := logWriter.WriteSuccess(
    ctx,
    userID,
    aggregate.ActivityUserLoggedIn,
    map[string]interface{}{
        "ip_address": "192.168.1.1",
        "device":     "iPhone",
    },
    application.WithIPAddress("192.168.1.1"),
    application.WithUserAgent("Mozilla/5.0..."),
)

// 使用方式 2：写入失败的活动
err := logWriter.WriteFailure(
    ctx,
    userID,
    aggregate.ActivityUserLocked,
    "多次登录失败",
    application.WithIPAddress("192.168.1.1"),
)
```

### 3. 领域事件处理器（仅处理副作用）

```go
type UserEventSubscriber struct {
    logger       *zap.Logger
    emailService EmailService
    statsRepo    StatisticsRepository
    // ❌ 不再有 auditLogger
}

func (s *UserEventSubscriber) handleUserRegistered(ctx context.Context, event kernel.DomainEvent) error {
    e := event.(*userEvent.UserRegisteredEvent)
    
    // ✅ 只处理副作用：发送邮件（异步）
    go func() {
        _ = s.emailService.SendWelcomeEmail(e.Email, e.Username)
    }()
    
    // ✅ 只处理副作用：初始化统计（异步）
    go func() {
        _ = s.statsRepo.InitializeUserStats(e.UserID.Int64())
    }()
    
    // ❌ 不再写入 ActivityLog
}
```

---

## 📊 事件分类矩阵

| 事件类型 | ActivityLog | DomainEvent | Outbox | 说明 |
|---------|-------------|-------------|--------|------|
| UserRegistered | ✅ 同步写入 | ✅ 永久 | ✅ | 业务事件 + 审计 |
| UserLoggedIn | ✅ 同步写入 | ⚠️ 可选 | ✅ | 主要是审计 |
| UserPasswordChanged | ✅ 同步写入 | ✅ 永久 | ✅ | 安全事件 |
| UserEmailChanged | ✅ 同步写入 | ✅ 永久 | ✅ | 安全事件 |
| UserLocked | ✅ 同步写入 | ✅ 永久 | ✅ | 安全事件 |
| UserProfileUpdated | ⚠️ 可选 | ⚠️ 可选 | ✅ | 普通业务 |
| StatsInitialized | ❌ 不记录 | ⚠️ 1 年 | ✅ | 内部事件 |

---

## 🔧 配置与最佳实践

### 1. 事件过期策略（待实现）

```sql
-- Outbox 表：7-30 天清理
ALTER TABLE outbox ADD COLUMN expires_at TIMESTAMP;
DELETE FROM outbox WHERE processed = true AND processed_at < NOW() - INTERVAL '30 days';

-- ActivityLog 表：6-12 个月归档
DELETE FROM activity_logs WHERE occurred_at < NOW() - INTERVAL '1 year';

-- DomainEvent 表：冷热分离
CREATE TABLE domain_events_hot AS 
SELECT * FROM domain_events WHERE occurred_at > NOW() - INTERVAL '6 months';
```

### 2. 监控指标

```go
// 待实现的监控
- activity_log_write_duration_seconds (P99 < 50ms)
- domain_event_publish_latency (P99 < 100ms)
- outbox_unprocessed_count (告警阈值 > 1000)
- asynq_queue_depth (告警阈值 > 10000)
```

---

## ✅ 检查清单

在实现新的事件时，请确认：

- [ ] **是否需要审计？** → 是 → 在 UseCase 内直接写入 ActivityLog
- [ ] **是否需要触发其他模块？** → 是 → 产生领域事件
- [ ] **是否是副作用？** → 是 → 通过 EventSubscriber 异步处理
- [ ] **是否可以失败？** → 否 → ActivityLog 必须在事务内
- [ ] **是否需要重试？** → 是 → 通过 Asynq 消息队列

---

## 📚 参考文档

- [DDD 领域驱动设计](docs/design/ddd-design-guide.md)
- [整洁架构规范](docs/design/clean-architecture-spec.md)
- [Outbox Pattern 实现](internal/infrastructure/messaging/asynq/event_publisher.go)
- [ActivityLogWriter 使用](internal/application/activity_log_writer.go)

---

**最后更新时间**: 2026-03-26  
**版本**: v2.0  
**主要变更**: ActivityLog 从事件异步处理改为 UseCase 同步写入
