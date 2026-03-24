# EventBus vs EventPublisher 使用指南

**日期:** 2026-03-24  
**状态:** ✅ 已实施  
**目标:** 明确两种事件处理模式的使用场景，避免重复执行和混淆

---

## 📋 架构概览

### 当前事件系统架构

```
┌─────────────────────────────────────────────────────────┐
│                   领域事件产生点                          │
│              (AggregateRoot.Publish)                    │
└────────────────────┬────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌──────────────────┐
│    EventBus     │    │ EventPublisher   │
│   (同步总线)    │    │  (异步发布器)     │
└────────┬────────┘    └─────────┬────────┘
         │                       │
         │                       ├──────────────┐
         │                       │              │
         ▼                       ▼              ▼
   ┌──────────┐          ┌──────────┐   ┌──────────┐
   │同步处理器 │          │保存日志  │   │Asynq队列 │
   │立即执行  │          │Activity  │   │Worker处理│
   │          │          │EventLog  │   │          │
   └──────────┘          └──────────┘   └──────────┘
```

---

## 🎯 两种模式对比

### EventBus（同步事件总线）

**特点：**
- ✅ **同步执行** - 在调用线程中立即执行
- ✅ **事务内触发** - 可以在数据库事务内完成
- ✅ **可回滚** - 失败可以回滚事务
- ⚠️ **阻塞主流程** - 影响响应时间
- ⚠️ **不能跨服务** - 仅限进程内通信

**使用场景：**
1. **需要立即看到效果的操作**
   ```go
   // 例：用户注册后立即发送欢迎邮件
   UserRegistered → 立即发送邮件 → 返回成功
   
   // 如果改为异步：
   UserRegistered → 加入队列 → Worker 稍后处理 → 用户可能已离开
   ```

2. **需要在事务内完成的操作**
   ```go
   tx.Begin()
   user.Save()
   EventBus.Publish()  // ← 立即触发，可回滚
   tx.Commit()
   ```

3. **失败需要回滚的操作**
   ```go
   // 例：激活用户时需要同时更新相关数据
   UserActivated → 更新用户状态 + 更新权限
   // 任何一步失败都需要回滚
   ```

4. **实时性要求高的业务逻辑**
   ```go
   // 例：用户余额变更通知
   BalanceChanged → 立即检查是否透支 → 发送警告
   ```

---

### EventPublisher（异步事件发布器）

**特点：**
- ✅ **异步执行** - 不阻塞主流程
- ✅ **高性能** - 快速响应
- ✅ **可重试** - 失败可以自动重试
- ✅ **可削峰填谷** - 缓冲突发流量
- ⚠️ **最终一致性** - 不是立即执行
- ⚠️ **调试复杂** - 分布式追踪困难

**使用场景：**
1. **可以延迟处理的操作**
   ```go
   // 例：审计日志记录
   UserLoggedIn → 记录登录日志 → 可以延迟几秒
   
   // 例：统计分析
   OrderPaid → 更新销售统计 → 不需要实时
   ```

2. **耗时较长的操作**
   ```go
   // 例：生成报表
   ReportRequested → 加入队列 → Worker 后台生成
   
   // 例：批量邮件发送
   NewsletterSent → 分批处理 → 避免超时
   ```

3. **失败可以重试的操作**
   ```go
   // 例：调用外部 API
   UserCreated → 同步到 CRM 系统 → 失败可重试
   
   // 例：发送短信
   VerificationCodeSent → 第三方服务 → 可重试
   ```

4. **解耦系统组件**
   ```go
   // 例：微服务间通信
   OrderCreated → 消息队列 → 库存服务/物流服务/财务服务
   ```

---

## 📊 决策矩阵

| 场景特征 | EventBus（同步） | EventPublisher（异步） |
|---------|-----------------|---------------------|
| **实时性要求** | 高（毫秒级） | 低（秒级/分钟级） |
| **执行时机** | 立即 | 稍后 |
| **事务需求** | 需要同步 | 可以异步 |
| **失败处理** | 需要回滚 | 可以重试 |
| **性能影响** | 可接受 | 需要优化 |
| **复杂度** | 简单 | 较复杂 |

---

## ✅ 本项目实施策略

### 当前分工

#### EventBus 负责（同步处理器）
```go
// 仅处理需要同步执行的业务逻辑
UserRegistered    → UserSideEffectHandler (业务副作用处理)
UserActivated     → UserSideEffectHandler
UserDeactivated   → UserSideEffectHandler
UserLoggedIn      → UserSideEffectHandler (更新最后登录时间等)
UserPasswordChanged → UserSideEffectHandler
UserEmailChanged  → UserSideEffectHandler
UserLocked        → UserSideEffectHandler
UserUnlocked      → UserSideEffectHandler
UserProfileUpdated → UserSideEffectHandler
```

#### EventPublisher 负责（异步处理）
```go
// 所有领域事件都会通过 EventPublisher 异步处理
// 包括：
// 1. 保存 ActivityLog（活动日志）
// 2. 保存 EventLog（事件日志）
// 3. 发布到 Asynq 队列（Worker 处理）

// Worker 中处理：
AuditSubscriber   → 保存审计日志（异步）
LoginLogSubscriber → 保存登录日志（异步）
UserSideEffectHandler → 其他副作用处理（异步备份）
```

---

## 🚫 避免的问题

### 问题 1: 重复执行 ❌

**错误示例（已修复）：**
```go
// ❌ 同时在 EventBus 和 EventPublisher 中注册
// EventBus.Subscribe("UserLoggedIn", AuditSubscriber.Handle)
// EventPublisher.Publish() → Asynq → AuditSubscriber.Handle

// 结果：同一段代码被执行了两次！
```

**正确做法：**
```go
// ✅ 只在 EventPublisher 中处理（异步）
// EventBus 不再注册 AuditSubscriber

// 或者只在 EventBus 中处理（同步）
// EventPublisher 不调用 AuditSubscriber
```

---

### 问题 2: 混淆使用场景 ❌

**错误示例：**
```go
// ❌ 在 EventBus 中处理耗时操作
EventBus.Subscribe("OrderCreated", func() {
    // 耗时 10 秒生成报表
    generateReport()
})

// 结果：API 响应极慢
```

**正确做法：**
```go
// ✅ 使用 EventPublisher 异步处理
EventPublisher.Publish("OrderCreated")
// Asynq → Worker → generateReport()
// API 立即响应
```

---

### 问题 3: 事务外触发 ❌

**错误示例：**
```go
// ❌ 在事务提交后才触发事件
tx.Begin()
user.Save()
tx.Commit()
EventBus.Publish()  // ← 已经无法回滚

// 如果事件处理失败，数据已经提交
```

**正确做法：**
```go
// ✅ 在事务内触发（使用 EventBus）
tx.Begin()
user.Save()
EventBus.Publish()  // ← 立即触发，可回滚
tx.Commit()

// 或者使用本地消息表 + EventPublisher
tx.Begin()
user.Save()
saveDomainEventToDB()  // 保存到 domain_events 表
tx.Commit()
publishFromDB()  // 从数据库读取并发布到队列
```

---

## 📝 代码示例

### 示例 1: 使用 EventBus（同步）

```go
// 1. 定义事件处理器
type UserActivationHandler struct {
    db *gorm.DB
}

func (h *UserActivationHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
    activatedEvent := event.(*event.UserActivatedEvent)
    
    // 立即更新相关数据（在同一事务中）
    err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 更新用户状态
        if err := tx.Model(&User{}).Where("id = ?", activatedEvent.UserID).
            Update("status", "active").Error; err != nil {
            return err
        }
        
        // 更新权限
        if err := tx.Model(&Permission{}).Where("user_id = ?", activatedEvent.UserID).
            Update("is_active", true).Error; err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return err  // 失败会回滚整个事务
    }
    
    return nil
}

// 2. 注册到 EventBus
bus.Subscribe("UserActivated", handler.Handle)

// 3. 在聚合根中触发
func (u *User) Activate() error {
    u.Status = "active"
    u.AddEvent(event.NewUserActivatedEvent(u.ID))
    return nil
}

// 4. 在应用服务中发布
func (s *UserService) ActivateUser(ctx context.Context, userID int64) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        user, err := s.userRepo.FindByID(ctx, userID)
        if err != nil {
            return err
        }
        
        if err := user.Activate(); err != nil {
            return err
        }
        
        if err := s.userRepo.Save(ctx, user); err != nil {
            return err
        }
        
        // 同步触发事件（在事务内）
        for _, evt := range user.Events() {
            if err := s.eventBus.Publish(ctx, evt); err != nil {
                return err  // 会回滚事务
            }
        }
        
        return nil
    })
}
```

---

### 示例 2: 使用 EventPublisher（异步）

```go
// 1. 定义异步处理器（Worker 中）
type AuditLogHandler struct {
    repo *ActivityLogRepository
}

func (h *AuditLogHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
    var eventData map[string]interface{}
    if err := json.Unmarshal(task.Payload(), &eventData); err != nil {
        return err
    }
    
    // 异步保存审计日志（不影响主流程）
    log := &model.ActivityLog{
        UserID:     eventData["user_id"].(int64),
        Action:     eventData["action"].(string),
        OccurredAt: time.Now(),
    }
    
    return h.repo.Save(ctx, log)
    // 失败会由 Asynq 自动重试
}

// 2. 在 EventPublisher 中发布
func (p *EventPublisherAdapter) Publish(ctx context.Context, event kernel.DomainEvent) error {
    // 1. 保存 ActivityLog（同步）
    p.saveActivityLog(ctx, event)
    
    // 2. 保存 EventLog（同步）
    p.saveEventLog(ctx, event)
    
    // 3. 发布到 Asynq 队列（异步）
    p.publishToQueue(ctx, event)
    
    return nil
}

// 3. Worker 处理
func (w *Worker) Start() error {
    srv := asynq.NewServer(w.redisOpt, asynq.Config{
        Concurrency: 10,
    })
    
    mux := asynq.NewServeMux()
    mux.HandleFunc("domain:event", w.processDomainEvent)
    
    return srv.Run(mux)
}
```

---

## 🎯 最佳实践

### 1. 优先使用异步模式

```go
// ✅ 默认选择 EventPublisher
// 除非有明确的同步需求

// 理由：
// - 更好的性能
// - 更好的可扩展性
// - 更好的容错能力
```

### 2. 明确标注执行方式

```go
// @SyncEventHandler 同步事件处理器
// 在事务内执行，失败会回滚
type UserActivationHandler struct {
    // ...
}

// @AsyncEventHandler 异步事件处理器
// 在 Worker 中执行，失败会重试
type AuditLogHandler struct {
    // ...
}
```

### 3. 避免混合模式

```go
// ❌ 不要同时在两个地方注册
// EventBus.Subscribe("UserRegistered", handler)
// EventPublisher.Publish() → Asynq → handler

// ✅ 选择一个模式
// 要么只在 EventBus 中注册（同步）
// 要么只在 EventPublisher 中处理（异步）
```

### 4. 监控和日志

```go
// 记录事件处理方式和性能
func (b *SimpleEventBus) Publish(ctx context.Context, event DomainEvent) error {
    start := time.Now()
    defer func() {
        log.Debug("Sync event processed", 
            zap.String("event", event.EventName()),
            zap.Duration("duration", time.Since(start)))
    }()
    
    // ... 处理逻辑
}

func (a *EventPublisherAdapter) Publish(ctx context.Context, event DomainEvent) error {
    start := time.Now()
    defer func() {
        log.Debug("Async event published",
            zap.String("event", event.EventName()),
            zap.Duration("duration", time.Since(start)))
    }()
    
    // ... 处理逻辑
}
```

---

## 📊 总结对比表

| 特性 | EventBus | EventPublisher |
|------|----------|----------------|
| **执行方式** | 同步 | 异步 |
| **执行位置** | 调用线程 | Worker 线程池 |
| **响应时间** | 立即（ms 级） | 延迟（s 级） |
| **事务支持** | ✅ 可在事务内 | ❌ 事务外 |
| **回滚能力** | ✅ 可回滚 | ❌ 需补偿 |
| **失败处理** | 直接失败 | 自动重试 |
| **性能影响** | 阻塞主流程 | 不阻塞 |
| **可扩展性** | 低（单机） | 高（分布式） |
| **调试难度** | 简单 | 较复杂 |
| **适用场景** | 核心业务逻辑 | 辅助功能、日志、通知 |

---

## 🔧 维护建议

### 定期检查项

1. **检查重复注册**
   ```bash
   # 查找同一事件是否在两个地方都注册
   grep -r "Subscribe.*UserRegistered"
   grep -r "publishToQueue.*UserRegistered"
   ```

2. **监控处理时间**
   ```go
   // 同步处理器应该 < 100ms
   // 异步处理器应该 < 5s
   ```

3. **检查失败率**
   ```go
   // 同步失败率应该 < 0.1%
   // 异步失败率应该 < 1%（有重试机制）
   ```

### 未来优化方向

1. **考虑完全异步化**
   - 如果业务可以接受最终一致性
   - 如果有完善的监控告警
   - 可以考虑完全移除 EventBus

2. **引入 CQRS 模式**
   - 命令侧使用 EventBus（保证一致性）
   - 查询侧使用 EventPublisher（提高性能）

3. **事件溯源**
   - 使用 EventLog 作为唯一事实来源
   - 通过事件回放重建状态

---

**最后更新:** 2026-03-24  
**维护者:** Development Team  
**下次复审:** 根据业务发展需要决定
