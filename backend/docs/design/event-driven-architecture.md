# 事件驱动架构

本文档详细介绍 Go DDD Scaffold 中的事件驱动架构设计与实现。

## 📋 事件驱动架构概述

### 什么是事件驱动？

**事件驱动架构（EDA）**是一种软件架构模式，通过事件的产生、传播和处理来组织系统组件之间的交互。

### 核心优势

1. **解耦** - 生产者不依赖消费者
2. **扩展性** - 轻松添加新的事件处理器
3. **可追溯性** - 完整的事件日志
4. **最终一致性** - 支持分布式事务

---

## 🎯 事件类型

### 领域事件（Domain Events）

**定义：** 领域中发生的重要事情，使用过去时态命名。

**示例：**
```go
UserRegistered      // 用户已注册
UserLoggedIn        // 用户已登录
TenantCreated       // 租户已创建
OrderPlaced         // 订单已提交
PaymentCompleted    // 支付已完成
```

**特点：**
- 业务语义明确
- 不可变
- 携带相关数据
- 用于解耦聚合

### 集成事件（Integration Events）

**定义：** 跨限界上下文传播的事件。

**示例：**
```go
// Auth Context → User Context
UserAuthenticated   // 用户已认证

// Order Context → Inventory Context
InventoryReserved   // 库存已预留
```

---

## 🏗️ 架构设计

### 完整的事件流

```
┌─────────────────────────────────────────┐
│     Domain Layer                        │
│                                         │
│  Aggregate                              │
│    ↓ RecordEvent()                      │
│  Domain Event                           │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│     Infrastructure Layer                │
│                                         │
│  Repository.Save()                      │
│    ↓                                    │
│  保存聚合 + 保存到 domain_events 表      │
│  (同一事务)                             │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│     Async Processor                     │
│                                         │
│  轮询 domain_events 表                  │
│    ↓                                    │
│  发布到 Redis Stream / Kafka            │
│    ↓                                    │
│  标记为已处理                           │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│     Event Handlers                      │
│                                         │
│  订阅特定事件类型                       │
│    ↓                                    │
│  执行业务逻辑                           │
│    ↓                                    │
│  可能触发新事件                         │
└─────────────────────────────────────────┘
```

---

## 💻 领域事件实现

### Domain Event 基类

```go
// domain/shared/kernel/domain_event.go
package kernel

import "time"

// DomainEvent 领域事件接口
type DomainEvent interface {
    Type() string           // 事件类型，如 "user.registered"
    AggregateID() int64     // 聚合根 ID
    AggregateType() string  // 聚合根类型，如 "User"
    Timestamp() time.Time   // 发生时间
}

// BaseDomainEvent 基础实现
type BaseDomainEvent struct {
    aggregateID   int64
    aggregateType string
    timestamp     time.Time
}

func (e *BaseDomainEvent) AggregateID() int64 {
    return e.aggregateID
}

func (e *BaseDomainEvent) AggregateType() string {
    return e.aggregateType
}

func (e *BaseDomainEvent) Timestamp() time.Time {
    return e.timestamp
}
```

### 事件示例

#### UserRegistered

```go
// domain/user/event/user_registered.go
package event

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// UserRegistered 用户已注册
type UserRegistered struct {
    kernel.BaseDomainEvent
    UserID  int64
    Email   string
    Created time.Time
}

func (e *UserRegistered) Type() string {
    return "user.registered"
}
```

#### UserLoggedIn

```go
// domain/user/event/user_logged_in.go
package event

import "time"

// UserLoggedIn 用户已登录
type UserLoggedIn struct {
    kernel.BaseDomainEvent
    UserID    int64
    Email     string
    IP        string
    UserAgent string
    Time      time.Time
}

func (e *UserLoggedIn) Type() string {
    return "user.logged_in"
}
```

---

## 🔄 Outbox Pattern 实现

### 为什么需要 Outbox Pattern？

**问题：** 如何保证业务操作和事件发布的原子性？

```
❌ 错误做法：
tx.Begin()
saveUser()
tx.Commit()
publishEvent()  // ← 如果失败怎么办？

✅ Outbox Pattern：
tx.Begin()
saveUser()
saveEventToOutbox()  // ← 同一事务
tx.Commit()
asyncPublishFromOutbox()  // ← 异步重试
```

### 数据库表设计

```sql
CREATE TABLE domain_events (
    id BIGINT PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    metadata JSONB,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed BOOLEAN NOT NULL DEFAULT false,
    processed_at TIMESTAMP NULL,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    
    INDEX idx_unprocessed (occurred_at) WHERE processed = false,
    INDEX idx_type (event_type)
);
```

### Repository 实现

```go
// infrastructure/persistence/repository/user_repository.go
func (r *userRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 1. 保存聚合根
    err := r.saveUser(tx, user)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 2. ⭐ 保存领域事件到 Outbox（同一事务）
    events := user.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}

func (r *userRepositoryImpl) saveEvent(tx *gorm.DB, event kernel.DomainEvent) error {
    eventData, _ := json.Marshal(event)
    
    daoEvent := &dao.DomainEvent{
        ID:            snowflake.Generate(),
        EventType:     event.Type(),
        AggregateType: event.AggregateType(),
        AggregateID:   fmt.Sprintf("%d", event.AggregateID()),
        EventData:     eventData,
        OccurredAt:    event.Timestamp(),
        Processed:     false,
    }
    
    return tx.Create(daoEvent).Error
}
```

---

## 📨 异步事件发布

### Event Publisher

```go
// infrastructure/eventstore/event_publisher.go
package eventstore

import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "go.uber.org/zap"
)

// EventPublisher 事件发布器
type EventPublisher interface {
    Publish(event kernel.DomainEvent)
}

type eventPublisherImpl struct {
    logger *zap.Logger
    redis  *redis.Client
}

func NewEventPublisher(logger *zap.Logger, redis *redis.Client) EventPublisher {
    return &eventPublisherImpl{
        logger: logger,
        redis:  redis,
    }
}

func (p *eventPublisherImpl) Publish(event kernel.DomainEvent) {
    ctx := context.Background()
    
    // 序列化事件
    data, err := json.Marshal(event)
    if err != nil {
        p.logger.Error("serialize event failed", zap.Error(err))
        return
    }
    
    // 发布到 Redis Stream
    streamKey := "stream:domain_events"
    err = p.redis.XAdd(ctx, &redis.XAddArgs{
        Stream: streamKey,
        Values: map[string]interface{}{
            "event_type":      event.Type(),
            "aggregate_type":  event.AggregateType(),
            "aggregate_id":    event.AggregateID(),
            "event_data":      string(data),
            "timestamp":       event.Timestamp().Unix(),
        },
    }).Err()
    
    if err != nil {
        p.logger.Error("publish event to redis failed", zap.Error(err))
        // 可以加入本地队列稍后重试
    }
}
```

### Domain Event Worker

```go
// cmd/worker/main.go
package main

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
    "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
)

func main() {
    cfg := config.LoadConfig()
    logger := logging.NewZapLogger(cfg.Server.Mode)
    
    // 创建 Asynq 客户端
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: cfg.Redis.Addr},
        asynq.Config{
            Concurrency: 10,
            Queues: map[string]int{
                "critical": 6,
                "default":  3,
                "low":      1,
            },
        },
    )
    
    // 创建事件处理器
    mux := asynq.NewServeMux()
    
    // 注册领域事件处理器
    mux.HandleFunc(event.TypeDomainEvent, event.HandleDomainEvent(logger))
    
    // 启动 Worker
    if err := srv.Run(mux); err != nil {
        logger.Fatal("failed to run worker", zap.Error(err))
    }
}
```

### 事件处理器

```go
// interfaces/event/domain_event_handler.go
package event

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "go.uber.org/zap"
)

const (
    TypeDomainEvent = "domain_event"
)

// HandleDomainEvent 处理领域事件
func HandleDomainEvent(logger *zap.Logger) func(ctx context.Context, t *asynq.Task) error {
    return func(ctx context.Context, t *asynq.Task) error {
        var payload struct {
            EventType     string          `json:"event_type"`
            AggregateType string          `json:"aggregate_type"`
            AggregateID   string          `json:"aggregate_id"`
            EventData     json.RawMessage `json:"event_data"`
        }
        
        if err := json.Unmarshal(t.Payload(), &payload); err != nil {
            return err
        }
        
        logger.Info("processing domain event",
            zap.String("event_type", payload.EventType),
            zap.String("aggregate_id", payload.AggregateID),
        )
        
        // 根据事件类型分发到不同的处理器
        switch payload.EventType {
        case "user.registered":
            return handleUserRegistered(ctx, payload.EventData, logger)
        case "user.logged_in":
            return handleUserLoggedIn(ctx, payload.EventData, logger)
        default:
            logger.Warn("unknown event type",
                zap.String("event_type", payload.EventType),
            )
        }
        
        return nil
    }
}

func handleUserRegistered(ctx context.Context, data json.RawMessage, logger *zap.Logger) error {
    var evt event.UserRegistered
    if err := json.Unmarshal(data, &evt); err != nil {
        return err
    }
    
    logger.Info("handling user registered",
        zap.Int64("user_id", evt.UserID),
        zap.String("email", evt.Email),
    )
    
    // 1. 发送欢迎邮件
    // emailService.SendWelcomeEmail(evt.Email)
    
    // 2. 初始化租户配置
    // tenantService.InitializeTenant(evt.UserID)
    
    // 3. 记录统计信息
    // analyticsService.Track("user_registered", ...)
    
    return nil
}
```

---

## 🎯 事件订阅模式

### 应用内订阅

```go
// application/auth/service.go
func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthResult, error) {
    // ... 认证逻辑
    
    // 发布事件
    s.eventPublisher.Publish(&event.UserLoggedIn{
        UserID:    user.ID().Value(),
        Email:     user.Email().String(),
        IP:        cmd.IP,
        UserAgent: cmd.UserAgent,
    })
    
    return result, nil
}

// module/auth.go
func (m *AuthModule) SubscribeEvents() {
    // 订阅用户登录事件
    m.infra.EventPublisher.Subscribe("user.logged_in", 
        func(event kernel.DomainEvent) {
            // 更新最后登录时间
            // 记录登录日志
        },
    )
}
```

### 跨服务订阅（微服务场景）

```go
// 发布到消息队列（Kafka/RabbitMQ）
func (p *eventPublisherImpl) Publish(event kernel.DomainEvent) {
    // 发布到 Kafka
    msg := kafka.Message{
        Topic: "domain-events." + event.Type(),
        Value: toJSON(event),
    }
    
    p.kafkaProducer.Send(msg)
}

// 其他服务订阅
consumer.Subscribe("domain-events.user.registered", 
    func(msg kafka.Message) {
        // 处理事件
    },
)
```

---

## ✅ 最佳实践

### 1. 事件命名规范

```go
// ✅ 正确：使用过去时态
UserRegistered      // 用户已注册
OrderPlaced         // 订单已提交
PaymentCompleted    // 支付已完成

// ❌ 错误：使用现在时或将来时
RegisterUser        // 这是命令，不是事件
PlaceOrder          // 这是命令
```

### 2. 事件携带足够的数据

```go
// ✅ 正确：携带完整数据
type UserRegistered struct {
    UserID  int64
    Email   string
    Created time.Time
}

// ❌ 错误：数据不足
type UserRegistered struct {
    UserID int64
    // 还需要查询数据库才能获取邮箱
}
```

### 3. 事件处理器幂等性

```go
// ✅ 正确：幂等处理
func handleUserRegistered(data EventData) error {
    // 检查是否已处理
    if alreadyProcessed(data.ID) {
        return nil
    }
    
    // 处理事件
    processEvent(data)
    
    // 标记已处理
    markAsProcessed(data.ID)
    
    return nil
}
```

### 4. 错误处理和重试

```go
// ✅ 正确：带重试的错误处理
func processWithRetry(event Event, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := handler(event)
        if err == nil {
            return nil
        }
        
        if i == maxRetries-1 {
            // 移到死信队列
            sendToDeadLetterQueue(event)
            return err
        }
        
        // 指数退避
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return nil
}
```

---

## 📊 事件流程图

### 用户注册流程

```
┌──────────────┐
│ HTTP Handler │
│ POST /users  │
└──────┬───────┘
       │
       ↓
┌──────────────┐
│ App Service  │
│ CreateUser() │
└──────┬───────┘
       │
       ↓
┌──────────────┐
│ Domain       │
│ User.New()   │
│              │
│ RecordEvent: │
│ UserRegistered
└──────┬───────┘
       │
       ↓
┌──────────────┐
│ Repository   │
│ Save()       │
│              │
│ TX:          │
│ 1. Save User │
│ 2. Save Event│
│    (Outbox)  │
└──────┬───────┘
       │
       ↓
┌──────────────┐
│ Async Worker │
│ Poll Outbox  │
│              │
│ Publish to   │
│ Redis Stream │
└──────┬───────┘
       │
       ↓
┌─────────────────────────────────┐
│ Event Handlers (并行执行)       │
├─────────────────────────────────┤
│ 1. EmailService                 │
│    → Send Welcome Email         │
│                                 │
│ 2. TenantService                │
│    → Initialize Tenant Config   │
│                                 │
│ 3. AnalyticsService             │
│    → Track Registration         │
└─────────────────────────────────┘
```

---

## 📚 参考资源

- [Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html)
- [Event-Driven Architecture](https://www.enterpriseintegrationpatterns.com/patterns/messaging/toc.html)
- [Domain Events](https://martinfowler.com/li/DomainEvent.html)
- [Asynq 文档](https://github.com/hibiken/asynq)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
