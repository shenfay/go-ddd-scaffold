# Go DDD Scaffold 项目深度分析与优化建议

**分析时间：** 2024-03-26  
**项目路径：** `/Users/shenfay/Projects/go-ddd-scaffold`  
**分析者：** 小沐

---

## 📊 目录

1. [整体架构思路](#整体架构思路)
2. [现有架构评估](#现有架构评估)
3. [架构优化建议](#架构优化建议)
4. [领域事件设计与实现](#领域事件设计与实现)
5. [完整代码示例](#完整代码示例)
6. [最佳实践总结](#最佳实践总结)

---

## 一、整体架构思路

### 1.1 架构概览

该项目采用 **Clean Architecture + DDD** 的组合模式，分为四层：

```
┌─────────────────────────────────────────────────────────────┐
│ Interfaces Layer (接口层)                                   │
│ - HTTP REST API (Gin)                                       │
│ - Handler + Middleware                                      │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Application Layer (应用层)                                  │
│ - Use Cases (业务流程编排)                                  │
│ - Ports (外部依赖接口)                                      │
│ - DTO (数据传输对象)                                        │
│ - Unit of Work (工作单元)                                   │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Domain Layer (领域层)                                       │
│ - Aggregate Root (聚合根)                                   │
│ - Value Objects (值对象)                                    │
│ - Domain Services (领域服务)                                │
│ - Domain Events (领域事件)                                  │
│ - Repository Interfaces (仓储接口)                          │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Infrastructure Layer (基础设施层)                           │
│ - Repository Implementations (仓储实现)                     │
│ - Adapters (适配器)                                         │
│ - External Services (外部服务)                              │
│ - Database, Cache, Message Queue                           │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心设计原则

#### ✅ 已正确实现的原则

1. **依赖反转（Dependency Inversion）**
   - Domain 层定义 Repository 接口
   - Infrastructure 层实现接口
   - Application 层通过接口依赖

2. **聚合根设计**
   - User 聚合根封装用户业务逻辑
   - LoginStats 独立聚合根（解决高频更新问题）
   - 值对象（UserName, Email, UserProfile）

3. **事件驱动架构**
   - 领域事件（UserRegistered, UserLoggedIn 等）
   - Outbox Pattern 保证事务一致性
   - Asynq 异步处理

4. **模块化设计**
   - AuthModule, UserModule 独立模块
   - 清晰的职责划分
   - 易于扩展新模块

---

## 二、现有架构评估

### 2.1 优势分析

| 方面 | 评价 | 说明 |
|------|------|------|
| **分层清晰** | ⭐⭐⭐⭐⭐ | 四层架构划分明确，依赖方向正确 |
| **DDD 实践** | ⭐⭐⭐⭐ | 聚合根、值对象、领域服务设计合理 |
| **事件驱动** | ⭐⭐⭐⭐ | Outbox Pattern + Asynq 实现完整 |
| **可测试性** | ⭐⭐⭐⭐ | 依赖注入、接口设计支持单元测试 |
| **文档完整** | ⭐⭐⭐⭐ | 设计文档详细，易于理解 |

### 2.2 存在的问题

#### 问题 1：事件处理流程不完整

**现状：**
```
Domain Event → Outbox → Asynq → ❌ 缺少事件处理器注册机制
```

**问题描述：**
- `SideEffectHandler` 在 `UserModule` 中创建但未使用
- Worker 中没有事件处理器的注册和分发逻辑
- 事件处理器与事件类型的映射不清晰

**影响：**
- 事件发布后无法被正确处理
- 副作用（如发送邮件）无法执行

#### 问题 2：事件处理器的幂等性设计缺失

**现状：**
- 没有事件处理记录表
- 没有幂等性检查机制
- 重复处理同一事件会导致重复副作用

**影响：**
- 邮件可能被发送多次
- 统计数据可能重复计算

#### 问题 3：事件版本管理不完善

**现状：**
```go
// 所有事件版本都是 1
BaseEvent{
    version: 1,  // ❌ 硬编码
}
```

**问题：**
- 无法处理事件结构演变
- 无法支持事件升级

#### 问题 4：事件处理错误处理不足

**现状：**
- 没有死信队列（DLQ）机制
- 没有重试策略配置
- 没有事件处理失败告警

#### 问题 5：事件数据完整性问题

**现状：**
```go
// UserPasswordChangedEvent 中只有 UserID
// 但处理器需要用户邮箱来发送邮件
type UserPasswordChangedEvent struct {
    UserID    vo.UserID
    ChangedAt time.Time
    IPAddress string
    // ❌ 缺少 Email 字段
}
```

**影响：**
- 事件处理器需要额外查询数据库
- 增加延迟和数据库压力

#### 问题 6：事件订阅机制不灵活

**现状：**
- 只有同步的 SimpleEventBus
- 没有异步事件订阅
- 没有事件优先级机制

#### 问题 7：事件溯源（Event Sourcing）未实现

**现状：**
- 只有事件日志记录
- 无法从事件重建聚合根状态
- 无法支持完整的事件溯源

---

## 三、架构优化建议

### 3.1 目录结构优化

#### 当前结构
```
internal/
├── domain/
│   ├── shared/kernel/
│   ├── user/
│   │   ├── aggregate/
│   │   ├── event/
│   │   ├── service/
│   │   ├── repository/
│   │   └── valueobject/
│   └── tenant/
├── application/
│   ├── user/usecase/
│   ├── auth/
│   ├── ports/
│   └── shared/dto/
├── infrastructure/
│   ├── persistence/
│   ├── messaging/
│   ├── platform/
│   └── worker/
└── interfaces/
    └── http/
```

#### 优化建议

**1. 增强事件处理层**

```
internal/
├── domain/
│   ├── shared/
│   │   ├── kernel/
│   │   │   ├── event.go              # 事件基类
│   │   │   ├── event_bus.go          # 事件总线
│   │   │   ├── event_handler.go      # ✨ 新增：事件处理器接口
│   │   │   └── event_store.go        # ✨ 新增：事件存储接口
│   │   └── event/
│   │       ├── registry.go           # ✨ 新增：事件类型注册表
│   │       └── versioning.go         # ✨ 新增：事件版本管理
│   └── user/
│       ├── event/
│       │   ├── events.go
│       │   ├── handlers.go           # ✨ 改进：事件处理器
│       │   └── event_factory.go      # ✨ 新增：事件工厂
│       └── ...
├── application/
│   ├── event/                        # ✨ 新增：事件应用层
│   │   ├── event_handler_registry.go
│   │   ├── event_processor.go
│   │   └── event_idempotency.go
│   └── ...
└── infrastructure/
    ├── eventstore/
    │   ├── outbox_processor.go
    │   ├── event_store_impl.go       # ✨ 新增：事件存储实现
    │   ├── event_handler_registry.go # ✨ 新增：处理器注册
    │   └── dead_letter_queue.go      # ✨ 新增：死信队列
    ├── messaging/
    │   └── asynq/
    │       ├── event_publisher.go
    │       ├── event_handler_adapter.go  # ✨ 新增：处理器适配器
    │       └── event_processor.go    # ✨ 新增：事件处理器
    └── ...
```

**2. 改进 Repository 层**

```
domain/
├── user/
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── event_repository.go       # ✨ 新增：事件仓储
│   │   └── event_snapshot_repository.go  # ✨ 新增：快照仓储
│   └── ...
└── ...
```

**3. 增强 Ports 定义**

```
application/
├── ports/
│   ├── event/
│   │   ├── event_publisher.go        # ✨ 改进：发布器接口
│   │   ├── event_handler.go          # ✨ 新增：处理器接口
│   │   └── event_store.go            # ✨ 新增：存储接口
│   ├── idempotency/
│   │   └── idempotency_service.go    # ✨ 新增：幂等性服务
│   └── ...
└── ...
```

### 3.2 各层职责明确划分

#### Domain Layer（领域层）

**职责：**
- ✅ 定义聚合根和值对象
- ✅ 实现业务规则和状态转换
- ✅ 定义领域事件
- ✅ 定义 Repository 接口
- ✅ 定义 Domain Service 接口
- ✨ **新增：定义事件处理器接口**
- ✨ **新增：定义事件存储接口**

**示例：**
```go
// domain/shared/kernel/event_handler.go
package kernel

import "context"

// EventHandler 事件处理器接口（Domain 层定义）
type EventHandler interface {
    // 处理事件
    Handle(ctx context.Context, event DomainEvent) error
    
    // 检查是否可以处理该事件
    CanHandle(eventType string) bool
    
    // 获取处理器名称
    Name() string
}

// EventStore 事件存储接口（Domain 层定义）
type EventStore interface {
    // 保存事件
    SaveEvent(ctx context.Context, event DomainEvent) error
    
    // 获取聚合根的所有事件
    GetEventsByAggregateID(ctx context.Context, aggregateID interface{}) ([]DomainEvent, error)
    
    // 获取特定版本之后的事件
    GetEventsSince(ctx context.Context, aggregateID interface{}, version int) ([]DomainEvent, error)
}
```

#### Application Layer（应用层）

**职责：**
- ✅ 编排业务流程（Use Cases）
- ✅ 定义 Ports 接口
- ✅ 定义 DTO
- ✅ 管理事务（Unit of Work）
- ✨ **新增：事件处理器注册和分发**
- ✨ **新增：幂等性检查**
- ✨ **新增：事件版本管理**

**示例：**
```go
// application/event/event_handler_registry.go
package event

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// EventHandlerRegistry 事件处理器注册表
type EventHandlerRegistry interface {
    // 注册处理器
    Register(eventType string, handler kernel.EventHandler) error
    
    // 获取处理器
    GetHandlers(eventType string) []kernel.EventHandler
    
    // 分发事件
    Dispatch(ctx context.Context, event kernel.DomainEvent) error
}

// eventHandlerRegistryImpl 实现
type eventHandlerRegistryImpl struct {
    handlers map[string][]kernel.EventHandler
}

func (r *eventHandlerRegistryImpl) Dispatch(ctx context.Context, event kernel.DomainEvent) error {
    handlers := r.GetHandlers(event.EventName())
    
    for _, handler := range handlers {
        if err := handler.Handle(ctx, event); err != nil {
            // 错误处理：记录、重试、DLQ
            return err
        }
    }
    
    return nil
}
```

#### Infrastructure Layer（基础设施层）

**职责：**
- ✅ 实现 Repository 接口
- ✅ 实现 Ports 接口（适配器）
- ✅ 数据持久化
- ✅ 外部服务集成
- ✨ **新增：实现事件处理器注册**
- ✨ **新增：实现死信队列**
- ✨ **新增：实现幂等性检查**
- ✨ **新增：实现事件版本升级**

#### Interfaces Layer（接口层）

**职责：**
- ✅ HTTP/gRPC 协议适配
- ✅ 请求验证
- ✅ 响应格式化
- ✨ **新增：事件处理结果反馈**

### 3.3 模块间依赖关系优化

#### 当前依赖关系

```
AuthModule
    ↓
UserModule
    ↓
Infrastructure (DB, Redis, Asynq)
    ↓
Domain (User, Tenant)
```

#### 优化后的依赖关系

```
┌─────────────────────────────────────────┐
│ EventModule (新增)                      │
│ - EventHandlerRegistry                  │
│ - EventProcessor                        │
│ - EventIdempotency                      │
└─────────────────────────────────────────┘
        ↑                           ↑
        │                           │
┌───────┴──────────┐        ┌──────┴──────────┐
│ AuthModule       │        │ UserModule      │
│ - Auth Service   │        │ - User Service  │
│ - Auth Handler   │        │ - User Handler  │
└───────┬──────────┘        └──────┬──────────┘
        │                           │
        └───────────┬───────────────┘
                    ↓
        ┌───────────────────────────┐
        │ Infrastructure            │
        │ - Repository Impl         │
        │ - Event Store Impl        │
        │ - Adapters                │
        └───────────┬───────────────┘
                    ↓
        ┌───────────────────────────┐
        │ Domain                    │
        │ - Aggregates              │
        │ - Events                  │
        │ - Interfaces              │
        └───────────────────────────┘
```

### 3.4 代码复用和可维护性提升

#### 1. 事件处理器基类

```go
// domain/shared/kernel/base_event_handler.go
package kernel

import (
    "context"
    "go.uber.org/zap"
)

// BaseEventHandler 事件处理器基类
type BaseEventHandler struct {
    logger *zap.Logger
    name   string
}

func NewBaseEventHandler(logger *zap.Logger, name string) *BaseEventHandler {
    return &BaseEventHandler{
        logger: logger,
        name:   name,
    }
}

func (h *BaseEventHandler) Name() string {
    return h.name
}

func (h *BaseEventHandler) LogInfo(msg string, fields ...zap.Field) {
    h.logger.Info(msg, fields...)
}

func (h *BaseEventHandler) LogError(msg string, err error, fields ...zap.Field) {
    h.logger.Error(msg, append(fields, zap.Error(err))...)
}
```

#### 2. 事件处理器工厂

```go
// infrastructure/eventstore/event_handler_factory.go
package eventstore

import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
    "go.uber.org/zap"
)

// EventHandlerFactory 事件处理器工厂
type EventHandlerFactory struct {
    logger *zap.Logger
    // 其他依赖
}

func NewEventHandlerFactory(logger *zap.Logger) *EventHandlerFactory {
    return &EventHandlerFactory{
        logger: logger,
    }
}

// CreateHandlers 创建所有事件处理器
func (f *EventHandlerFactory) CreateHandlers() []kernel.EventHandler {
    return []kernel.EventHandler{
        userEvent.NewUserRegisteredHandler(f.logger),
        userEvent.NewUserLoggedInHandler(f.logger),
        userEvent.NewUserPasswordChangedHandler(f.logger),
        // ... 其他处理器
    }
}
```

#### 3. 统一的错误处理

```go
// application/event/event_error_handler.go
package event

import (
    "context"
    "fmt"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "go.uber.org/zap"
)

// EventErrorHandler 事件错误处理器
type EventErrorHandler struct {
    logger *zap.Logger
    dlq    DeadLetterQueue
}

func (h *EventErrorHandler) HandleError(
    ctx context.Context,
    event kernel.DomainEvent,
    handler kernel.EventHandler,
    err error,
) error {
    h.logger.Error("event handler failed",
        zap.String("event_type", event.EventName()),
        zap.String("handler", handler.Name()),
        zap.Error(err),
    )
    
    // 发送到死信队列
    return h.dlq.Send(ctx, &DeadLetterMessage{
        EventType:   event.EventName(),
        Handler:     handler.Name(),
        Event:       event,
        Error:       err.Error(),
        Timestamp:   time.Now(),
    })
}
```

---

## 四、领域事件设计与实现

### 4.1 领域事件的标准结构

#### 基础事件接口

```go
// domain/shared/kernel/event.go
package kernel

import (
    "time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
    // 事件名称（如 "UserRegistered"）
    EventName() string
    
    // 事件发生时间
    OccurredOn() time.Time
    
    // 聚合根 ID
    AggregateID() interface{}
    
    // 聚合根类型（如 "User"）
    AggregateType() string
    
    // 事件版本（用于事件升级）
    Version() int
    
    // 事件元数据
    Metadata() map[string]interface{}
}

// BaseEvent 基础事件实现
type BaseEvent struct {
    eventName     string
    aggregateID   interface{}
    aggregateType string
    version       int
    occurredOn    time.Time
    metadata      map[string]interface{}
}

func NewBaseEvent(
    eventName string,
    aggregateID interface{},
    aggregateType string,
    version int,
) *BaseEvent {
    return &BaseEvent{
        eventName:     eventName,
        aggregateID:   aggregateID,
        aggregateType: aggregateType,
        version:       version,
        occurredOn:    time.Now(),
        metadata:      make(map[string]interface{}),
    }
}

// 实现 DomainEvent 接口
func (e *BaseEvent) EventName() string       { return e.eventName }
func (e *BaseEvent) OccurredOn() time.Time   { return e.occurredOn }
func (e *BaseEvent) AggregateID() interface{} { return e.aggregateID }
func (e *BaseEvent) AggregateType() string   { return e.aggregateType }
func (e *BaseEvent) Version() int            { return e.version }
func (e *BaseEvent) Metadata() map[string]interface{} {
    if e.metadata == nil {
        e.metadata = make(map[string]interface{})
    }
    return e.metadata
}

func (e *BaseEvent) SetMetadata(key string, value interface{}) {
    if e.metadata == nil {
        e.metadata = make(map[string]interface{})
    }
    e.metadata[key] = value
}
```

#### 具体事件示例

```go
// domain/user/event/events.go
package event

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// ============================================================================
// UserRegisteredEvent - 用户注册事件
// ============================================================================

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
    *kernel.BaseEvent
    
    // 事件数据（包含处理所需的所有信息）
    UserID         vo.UserID `json:"user_id"`
    Username       string    `json:"username"`
    Email          string    `json:"email"`
    Status         string    `json:"status"`
    DisplayName    string    `json:"display_name"`
    RegistrationIP string    `json:"registration_ip"`
    TenantID       int64     `json:"tenant_id"`
    RegisteredAt   time.Time `json:"registered_at"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(
    userID vo.UserID,
    username, email, status, displayName, registrationIP string,
    tenantID int64,
) *UserRegisteredEvent {
    event := &UserRegisteredEvent{
        BaseEvent:      kernel.NewBaseEvent("UserRegistered", userID, "User", 1),
        UserID:         userID,
        Username:       username,
        Email:          email,
        Status:         status,
        DisplayName:    displayName,
        RegistrationIP: registrationIP,
        TenantID:       tenantID,
        RegisteredAt:   time.Now(),
    }
    
    // 设置元数据
    event.SetMetadata("event_type", "domain_event")
    event.SetMetadata("aggregate_type", "user")
    event.SetMetadata("security_event", false)
    
    return event
}

// ============================================================================
// UserPasswordChangedEvent - 用户密码修改事件（改进版）
// ============================================================================

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
    *kernel.BaseEvent
    
    // ✨ 改进：包含处理所需的所有信息
    UserID    vo.UserID `json:"user_id"`
    Email     string    `json:"email"`           // ✨ 新增：邮箱
    Username  string    `json:"username"`        // ✨ 新增：用户名
    ChangedAt time.Time `json:"changed_at"`
    IPAddress string    `json:"ip_address"`
    ChangedBy string    `json:"changed_by"`      // ✨ 新增：谁修改的
}

// NewUserPasswordChangedEvent 创建用户密码修改事件
func NewUserPasswordChangedEvent(
    userID vo.UserID,
    email, username, ipAddress, changedBy string,
) *UserPasswordChangedEvent {
    event := &UserPasswordChangedEvent{
        BaseEvent:  kernel.NewBaseEvent("UserPasswordChanged", userID, "User", 1),
        UserID:     userID,
        Email:      email,
        Username:   username,
        ChangedAt:  time.Now(),
        IPAddress:  ipAddress,
        ChangedBy:  changedBy,
    }
    
    event.SetMetadata("event_type", "domain_event")
    event.SetMetadata("aggregate_type", "user")
    event.SetMetadata("security_event", true)
    
    return event
}

// ============================================================================
// UserLoggedInEvent - 用户登录事件
// ============================================================================

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
    *kernel.BaseEvent
    
    UserID            vo.UserID `json:"user_id"`
    Email             string    `json:"email"`
    Username          string    `json:"username"`
    LoginAt           time.Time `json:"login_at"`
    IPAddress         string    `json:"ip_address"`
    UserAgent         string    `json:"user_agent"`
    Location          string    `json:"location"`
    DeviceType        string    `json:"device_type"`
    DeviceFingerprint string    `json:"device_fingerprint"`
    LoginMethod       string    `json:"login_method"`
    Success           bool      `json:"success"`
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(
    userID vo.UserID,
    email, username, ipAddress, userAgent, location, deviceType, deviceFingerprint, loginMethod string,
    success bool,
) *UserLoggedInEvent {
    event := &UserLoggedInEvent{
        BaseEvent:         kernel.NewBaseEvent("UserLoggedIn", userID, "User", 1),
        UserID:            userID,
        Email:             email,
        Username:          username,
        LoginAt:           time.Now(),
        IPAddress:         ipAddress,
        UserAgent:         userAgent,
        Location:          location,
        DeviceType:        deviceType,
        DeviceFingerprint: deviceFingerprint,
        LoginMethod:       loginMethod,
        Success:           success,
    }
    
    event.SetMetadata("event_type", "domain_event")
    event.SetMetadata("aggregate_type", "user")
    event.SetMetadata("security_event", true)
    
    return event
}
```

### 4.2 事件发布、订阅、处理的完整流程

#### 流程图

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. 聚合根产生事件                                                │
│    User.NewUser() → RecordEvent(UserRegisteredEvent)             │
└──────────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────────┐
│ 2. 应用服务保存聚合根和事件（同一事务）                          │
│    UnitOfWork.Transaction() {                                    │
│        userRepo.Save(user)                                       │
│        eventPublisher.Publish(event)  // 保存到 Outbox           │
│    }                                                              │
└──────────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────────┐
│ 3. Outbox 处理器轮询未发布事件                                   │
│    OutboxProcessor.Start() {                                     │
│        for event in unpublished_events {                         │
│            publisher.PublishDomainEvent(event)                   │
│            mark_as_processed(event)                              │
│        }                                                          │
│    }                                                              │
└──────────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────────┐
│ 4. 事件发布到消息队列（Asynq）                                   │
│    Publisher.PublishDomainEvent(event) {                         │
│        task = NewDomainEventTask(event)                          │
│        asynqClient.Enqueue(task)                                 │
│    }                                                              │
└──────────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────────┐
│ 5. Worker 处理事件                                               │
│    EventProcessor.ProcessDomainEvent(task) {                     │
│        event = ExtractEvent(task)                                │
│        handlers = registry.GetHandlers(event.Type)               │
│        for handler in handlers {                                 │
│            handler.Handle(event)                                 │
│        }                                                          │
│    }                                                              │
└──────────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────────┐
│ 6. 事件处理器执行副作用                                          │
│    UserRegisteredHandler.Handle(event) {                         │
│        emailService.SendWelcomeEmail(event.Email)                │
│        tenantService.InitializeTenant(event.UserID)              │
│        analyticsService.Track("user_registered", event)          │
│    }                                                              │
└──────────────────────────────────────────────────────────────────┘
```

#### 完整代码实现

**第 1 步：聚合根产生事件**

```go
// domain/user/aggregate/user.go
package aggregate

import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// User 聚合根
type User struct {
    kernel.BaseEntity
    username  *vo.UserName
    email     *vo.Email
    password  *vo.HashedPassword
    status    vo.UserStatus
    profile   *vo.UserProfile
}

// NewUser 创建新用户
func NewUser(username, email, hashedPassword string, idGenerator func() int64) (*User, error) {
    // ... 验证逻辑
    
    user := &User{
        status:  vo.UserStatusActive,
        profile: prof,
    }
    
    user.SetID(vo.NewUserID(idGenerator()))
    user.username = un
    user.email = em
    user.password = vo.NewHashedPassword(hashedPassword)
    
    // ✨ 产生领域事件
    registeredEvent := event.NewUserRegisteredEvent(
        user.ID().(vo.UserID),
        username,
        email,
        user.status.String(),
        user.profile.DisplayName(),
        "",  // registrationIP
        0,   // tenantID
    )
    user.ApplyEvent(registeredEvent)  // 记录到 uncommittedEvents
    
    return user, nil
}
```

**第 2 步：应用服务保存聚合根和事件**

```go
// internal/application/user/usecase/register_user.go
package usecase

import (
    "context"
    "fmt"
    "github.com/shenfay/go-ddd-scaffold/internal/application"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
)

// RegisterUserUseCase 注册用户用例
type RegisterUserUseCase struct {
    uow             application.UnitOfWork
    registrationSvc *service.RegistrationService
    eventPublisher  kernel.EventPublisher
}

// Execute 执行注册用户用例
func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
    var newUser *aggregate.User
    
    // ✨ 在事务中执行
    err := uc.uow.Transaction(ctx, func(ctx context.Context) error {
        // 1. 调用领域服务执行注册
        var err error
        newUser, err = uc.registrationSvc.Register(ctx, service.RegisterRequest{
            Username: cmd.Username,
            Email:    cmd.Email,
            Password: cmd.Password,
        })
        if err != nil {
            return err
        }
        
        // 2. 保存用户
        userRepo := uc.uow.UserRepository()
        if err := userRepo.Save(ctx, newUser); err != nil {
            return err
        }
        
        // 3. 保存登录统计
        loginStatsRepo := uc.uow.LoginStatsRepository()
        loginStats := aggregate.NewLoginStats(newUser.ID().(vo.UserID))
        if err := loginStatsRepo.Save(ctx, loginStats); err != nil {
            return err
        }
        
        // 4. ✨ 发布领域事件（在同一事务中）
        events := newUser.GetUncommittedEvents()
        for _, event := range events {
            if err := uc.eventPublisher.Publish(ctx, event); err != nil {
                return fmt.Errorf("failed to publish event: %w", err)
            }
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return &RegisterUserResult{
        UserID:   newUser.ID().(vo.UserID).Int64(),
        Username: newUser.Username().Value(),
        Email:    newUser.Email().Value(),
    }, nil
}
```

**第 3 步：Outbox 处理器轮询未发布事件**

```go
// internal/infrastructure/eventstore/outbox_processor.go
package eventstore

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/asynq"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

// OutboxProcessor Outbox 处理器
type OutboxProcessor struct {
    db           *gorm.DB
    publisher    *asynq.Publisher
    logger       *zap.Logger
    pollInterval time.Duration
    batchSize    int
}

// Start 启动轮询
func (p *OutboxProcessor) Start(ctx context.Context) error {
    ticker := time.NewTicker(p.pollInterval)
    defer ticker.Stop()
    
    p.logger.Info("Outbox processor started",
        zap.Duration("poll_interval", p.pollInterval),
        zap.Int("batch_size", p.batchSize),
    )
    
    for {
        select {
        case <-ctx.Done():
            p.logger.Info("Outbox processor stopped")
            return ctx.Err()
        case <-ticker.C:
            if err := p.processUnpublishedEvents(ctx); err != nil {
                p.logger.Error("Failed to process outbox", zap.Error(err))
            }
        }
    }
}

// processUnpublishedEvents 处理未发布的事件
func (p *OutboxProcessor) processUnpublishedEvents(ctx context.Context) error {
    // 1. 查询未处理的事件
    var events []*model.Outbox
    err := p.db.WithContext(ctx).
        Where("processed = ?", false).
        Order("occurred_at ASC").
        Limit(p.batchSize).
        Find(&events).Error
    
    if err != nil {
        return fmt.Errorf("failed to query unpublished events: %w", err)
    }
    
    if len(events) == 0 {
        return nil
    }
    
    p.logger.Debug("Found unpublished events", zap.Int("count", len(events)))
    
    // 2. 逐个发布
    successCount := 0
    failedCount := 0
    
    for _, event := range events {
        if err := p.publishEvent(ctx, event); err != nil {
            p.logger.Error("Failed to publish event",
                zap.Int64("id", event.ID),
                zap.String("type", event.EventType),
                zap.Error(err),
            )
            failedCount++
            
            // 增加重试计数
            if retryErr := p.incrementRetry(ctx, event.ID, err); retryErr != nil {
                p.logger.Error("Failed to increment retry count",
                    zap.Int64("id", event.ID),
                    zap.Error(retryErr),
                )
            }
        } else {
            successCount++
        }
    }
    
    p.logger.Info("Processed outbox events",
        zap.Int("success", successCount),
        zap.Int("failed", failedCount),
    )
    
    return nil
}

// publishEvent 发布单个事件
func (p *OutboxProcessor) publishEvent(ctx context.Context, event *model.Outbox) error {
    // 1. 反序列化事件数据
    var payload asynq.DomainEventPayload
    if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal event: %w", err)
    }
    
    // 2. 发布到 Asynq
    if err := p.publisher.PublishDomainEvent(ctx, payload, "default"); err != nil {
        return fmt.Errorf("failed to publish to asynq: %w", err)
    }
    
    // 3. 标记为已处理
    now := time.Now()
    return p.db.WithContext(ctx).Model(event).Updates(map[string]interface{}{
        "processed":    true,
        "processed_at": now,
        "updated_at":   now,
    }).Error
}
```

**第 4 步：事件发布到消息队列**

```go
// internal/infrastructure/messaging/asynq/event_publisher.go
package asynq

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
    idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
    "go.uber.org/zap"
)

// EventPublisherAdapter 事件发布器适配器
type EventPublisherAdapter struct {
    query         *dao.Query
    taskPublisher *Publisher
    logger        *zap.Logger
}

// Publish 发布领域事件
func (a *EventPublisherAdapter) Publish(ctx context.Context, event kernel.DomainEvent) error {
    a.logger.Debug("Publishing event",
        zap.String("event_type", event.EventName()),
        zap.Any("aggregate_id", event.AggregateID()),
    )
    
    // 1. 记录活动日志
    if err := a.saveActivityLog(ctx, event); err != nil {
        a.logger.Error("Failed to save activity log", zap.Error(err))
        return err
    }
    
    // 2. 记录事件日志
    if err := a.saveEventLog(ctx, event); err != nil {
        a.logger.Error("Failed to save event log", zap.Error(err))
        return err
    }
    
    // 3. ✨ 记录到 Outbox 表（Outbox Pattern）
    if err := a.saveToOutbox(ctx, event); err != nil {
        a.logger.Error("Failed to save to outbox", zap.Error(err))
        return err
    }
    
    return nil
}

// saveToOutbox 保存事件到 Outbox 表
func (a *EventPublisherAdapter) saveToOutbox(ctx context.Context, event kernel.DomainEvent) error {
    // 序列化事件数据
    payload, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    now := time.Now()
    daoModel := &model.Outbox{
        ID:            idgen.Generate(),
        EventType:     event.EventName(),
        AggregateType: event.AggregateType(),
        AggregateID:   fmt.Sprintf("%v", event.AggregateID()),
        Payload:       string(payload),
        OccurredAt:    &now,
        Processed:     false,
        RetryCount:    0,
        CreatedAt:     &now,
        UpdatedAt:     &now,
    }
    
    return a.query.Outbox.WithContext(ctx).Create(daoModel)
}
```

**第 5 步：Worker 处理事件**

```go
// internal/infrastructure/worker/processor.go
package worker

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/internal/application/event"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/asynq"
    "go.uber.org/zap"
)

// EventProcessor 事件处理器
type EventProcessor struct {
    logger           *zap.Logger
    handlerRegistry  event.EventHandlerRegistry
    idempotencyCheck event.IdempotencyService
}

// ProcessDomainEvent 处理领域事件
func (p *EventProcessor) ProcessDomainEvent(ctx context.Context, task *asynq.Task) error {
    // 1. 提取事件负载
    payload, err := asynq.ExtractDomainEventPayload(task)
    if err != nil {
        p.logger.Error("Failed to extract event payload", zap.Error(err))
        return err
    }
    
    p.logger.Info("Processing domain event",
        zap.String("event_type", payload.EventType),
        zap.String("aggregate_id", payload.AggregateID),
    )
    
    // 2. ✨ 检查幂等性（防止重复处理）
    processed, err := p.idempotencyCheck.IsProcessed(ctx, payload.EventType, payload.AggregateID)
    if err != nil {
        p.logger.Error("Failed to check idempotency", zap.Error(err))
        return err
    }
    
    if processed {
        p.logger.Info("Event already processed, skipping",
            zap.String("event_type", payload.EventType),
            zap.String("aggregate_id", payload.AggregateID),
        )
        return nil
    }
    
    // 3. 分发事件到处理器
    if err := p.handlerRegistry.Dispatch(ctx, payload); err != nil {
        p.logger.Error("Failed to dispatch event",
            zap.String("event_type", payload.EventType),
            zap.Error(err),
        )
        return err
    }
    
    // 4. ✨ 标记为已处理
    if err := p.idempotencyCheck.MarkAsProcessed(ctx, payload.EventType, payload.AggregateID); err != nil {
        p.logger.Error("Failed to mark as processed", zap.Error(err))
        return err
    }
    
    p.logger.Info("Event processed successfully",
        zap.String("event_type", payload.EventType),
        zap.String("aggregate_id", payload.AggregateID),
    )
    
    return nil
}
```

**第 6 步：事件处理器执行副作用**

```go
// internal/domain/user/event/handlers.go
package event

import (
    "context"
    "fmt"
    
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "go.uber.org/zap"
)

// UserRegisteredHandler 用户注册事件处理器
type UserRegisteredHandler struct {
    logger       *zap.Logger
    emailService EmailService
    tenantSvc    TenantService
}

func NewUserRegisteredHandler(
    logger *zap.Logger,
    emailService EmailService,
    tenantSvc TenantService,
) *UserRegisteredHandler {
    return &UserRegisteredHandler{
        logger:       logger,
        emailService: emailService,
        tenantSvc:    tenantSvc,
    }
}

// Handle 处理用户注册事件
func (h *UserRegisteredHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
    registeredEvent, ok := event.(*UserRegisteredEvent)
    if !ok {
        return fmt.Errorf("invalid event type")
    }
    
    h.logger.Info("Handling UserRegistered event",
        zap.String("user_id", registeredEvent.UserID.String()),
        zap.String("email", registeredEvent.Email),
    )
    
    // 1. 发送欢迎邮件
    if err := h.emailService.SendWelcomeEmail(registeredEvent.Email, registeredEvent.Username); err != nil {
        h.logger.Error("Failed to send welcome email",
            zap.String("email", registeredEvent.Email),
            zap.Error(err),
        )
        // 邮件发送失败不阻断主流程
    }
    
    // 2. 初始化租户配置
    if err := h.tenantSvc.InitializeTenant(ctx, registeredEvent.UserID.Int64()); err != nil {
        h.logger.Error("Failed to initialize tenant",
            zap.String("user_id", registeredEvent.UserID.String()),
            zap.Error(err),
        )
        // 租户初始化失败需要重试
        return err
    }
    
    return nil
}

// CanHandle 检查是否可以处理该事件
func (h *UserRegisteredHandler) CanHandle(eventType string) bool {
    return eventType == "UserRegistered"
}

// Name 返回处理器名称
func (h *UserRegisteredHandler) Name() string {
    return "UserRegisteredHandler"
}

// ============================================================================
// UserPasswordChangedHandler - 用户密码修改事件处理器
// ============================================================================

type UserPasswordChangedHandler struct {
    logger       *zap.Logger
    emailService EmailService
}

func NewUserPasswordChangedHandler(
    logger *zap.Logger,
    emailService EmailService,
) *UserPasswordChangedHandler {
    return &UserPasswordChangedHandler{
        logger:       logger,
        emailService: emailService,
    }
}

// Handle 处理用户密码修改事件
func (h *UserPasswordChangedHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
    passwordEvent, ok := event.(*UserPasswordChangedEvent)
    if !ok {
        return fmt.Errorf("invalid event type")
    }
    
    h.logger.Info("Handling UserPasswordChanged event",
        zap.String("user_id", passwordEvent.UserID.String()),
        zap.String("email", passwordEvent.Email),
    )
    
    // ✨ 改进：事件中已包含邮箱，无需额外查询
    if err := h.emailService.SendPasswordChangedEmail(
        passwordEvent.Email,
        passwordEvent.Username,
        passwordEvent.IPAddress,
    ); err != nil {
        h.logger.Error("Failed to send password changed email",
            zap.String("email", passwordEvent.Email),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}

func (h *UserPasswordChangedHandler) CanHandle(eventType string) bool {
    return eventType == "UserPasswordChanged"
}

func (h *UserPasswordChangedHandler) Name() string {
    return "UserPasswordChangedHandler"
}
```

### 4.3 事件持久化和异步处理的最佳实践

#### 1. 事件持久化策略

```go
// internal/infrastructure/persistence/repository/event_repository.go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
    idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
    "gorm.io/gorm"
)

// EventRepository 事件仓储接口
type EventRepository interface {
    // 保存事件
    SaveEvent(ctx context.Context, event kernel.DomainEvent) error
    
    // 获取聚合根的所有事件
    GetEventsByAggregateID(ctx context.Context, aggregateID interface{}) ([]kernel.DomainEvent, error)
    
    // 获取特定版本之后的事件
    GetEventsSince(ctx context.Context, aggregateID interface{}, version int) ([]kernel.DomainEvent, error)
    
    // 获取未处理的事件
    GetUnprocessedEvents(ctx context.Context, limit int) ([]kernel.DomainEvent, error)
    
    // 标记事件为已处理
    MarkAsProcessed(ctx context.Context, eventID int64) error
}

// eventRepositoryImpl 事件仓储实现
type eventRepositoryImpl struct {
    db    *gorm.DB
    query *dao.Query
}

func NewEventRepository(db *gorm.DB, query *dao.Query) EventRepository {
    return &eventRepositoryImpl{
        db:    db,
        query: query,
    }
}

// SaveEvent 保存事件
func (r *eventRepositoryImpl) SaveEvent(ctx context.Context, event kernel.DomainEvent) error {
    eventData, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    now := time.Now()
    daoEvent := &model.DomainEvent{
        ID:            idgen.Generate(),
        AggregateID:   fmt.Sprintf("%v", event.AggregateID()),
        AggregateType: event.AggregateType(),
        EventType:     event.EventName(),
        EventData:     string(eventData),
        Version:       int32(event.Version()),
        OccurredAt:    &now,
        CreatedAt:     &now,
    }
    
    return r.query.DomainEvent.WithContext(ctx).Create(daoEvent)
}

// GetEventsByAggregateID 获取聚合根的所有事件
func (r *eventRepositoryImpl) GetEventsByAggregateID(
    ctx context.Context,
    aggregateID interface{},
) ([]kernel.DomainEvent, error) {
    var models []*model.DomainEvent
    err := r.query.DomainEvent.WithContext(ctx).
        Where("aggregate_id = ?", fmt.Sprintf("%v", aggregateID)).
        Order("occurred_at ASC").
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    // 反序列化事件
    events := make([]kernel.DomainEvent, len(models))
    for i, m := range models {
        event, err := r.deserializeEvent(m)
        if err != nil {
            return nil, err
        }
        events[i] = event
    }
    
    return events, nil
}

// deserializeEvent 反序列化事件
func (r *eventRepositoryImpl) deserializeEvent(m *model.DomainEvent) (kernel.DomainEvent, error) {
    // 根据事件类型反序列化
    // 这里需要一个事件工厂来创建具体的事件对象
    // ...
    return nil, nil
}
```

#### 2. 异步处理配置

```go
// internal/infrastructure/messaging/asynq/config.go
package asynq

import (
    "time"
)

// AsyncConfig 异步处理配置
type AsyncConfig struct {
    // 重试配置
    MaxRetries      int           // 最大重试次数
    RetryBackoff    time.Duration // 重试退避时间
    
    // 队列配置
    Queues map[string]int // 队列名称 → 优先级
    
    // 处理器配置
    Concurrency int // 并发处理数
    
    // 超时配置
    TaskTimeout time.Duration // 任务超时时间
    
    // 死信队列配置
    DeadLetterQueueEnabled bool
    DeadLetterQueueName    string
}

// DefaultAsyncConfig 默认异步处理配置
func DefaultAsyncConfig() *AsyncConfig {
    return &AsyncConfig{
        MaxRetries:   10,
        RetryBackoff: 5 * time.Second,
        Queues: map[string]int{
            "critical": 6,
            "default":  3,
            "low":      1,
        },
        Concurrency:            10,
        TaskTimeout:            30 * time.Minute,
        DeadLetterQueueEnabled: true,
        DeadLetterQueueName:    "dead_letter",
    }
}
```

#### 3. 死信队列处理

```go
// internal/infrastructure/eventstore/dead_letter_queue.go
package eventstore

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
    idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

// DeadLetterMessage 死信消息
type DeadLetterMessage struct {
    ID            int64       `json:"id"`
    EventType     string      `json:"event_type"`
    AggregateID   string      `json:"aggregate_id"`
    Handler       string      `json:"handler"`
    EventData     string      `json:"event_data"`
    Error         string      `json:"error"`
    RetryCount    int         `json:"retry_count"`
    LastErrorAt   time.Time   `json:"last_error_at"`
    CreatedAt     time.Time   `json:"created_at"`
}

// DeadLetterQueue 死信队列
type DeadLetterQueue struct {
    db     *gorm.DB
    logger *zap.Logger
}

func NewDeadLetterQueue(db *gorm.DB, logger *zap.Logger) *DeadLetterQueue {
    return &DeadLetterQueue{
        db:     db,
        logger: logger,
    }
}

// Send 发送消息到死信队列
func (q *DeadLetterQueue) Send(ctx context.Context, msg *DeadLetterMessage) error {
    msg.ID = idgen.Generate()
    msg.CreatedAt = time.Now()
    msg.LastErrorAt = time.Now()
    
    // 保存到数据库
    daoMsg := &model.DeadLetterMessage{
        ID:          msg.ID,
        EventType:   msg.EventType,
        AggregateID: msg.AggregateID,
        Handler:     msg.Handler,
        EventData:   msg.EventData,
        Error:       msg.Error,
        RetryCount:  int32(msg.RetryCount),
        LastErrorAt: &msg.LastErrorAt,
        CreatedAt:   &msg.CreatedAt,
    }
    
    if err := q.db.WithContext(ctx).Create(daoMsg).Error; err != nil {
        q.logger.Error("Failed to save dead letter message",
            zap.String("event_type", msg.EventType),
            zap.Error(err),
        )
        return err
    }
    
    q.logger.Warn("Message sent to dead letter queue",
        zap.String("event_type", msg.EventType),
        zap.String("aggregate_id", msg.AggregateID),
        zap.String("error", msg.Error),
    )
    
    return nil
}

// Retry 重试死信消息
func (q *DeadLetterQueue) Retry(ctx context.Context, messageID int64) error {
    // 从死信队列中取出消息
    // 重新发布到消息队列
    // 如果成功，删除死信消息
    // 如果失败，增加重试计数
    
    q.logger.Info("Retrying dead letter message", zap.Int64("message_id", messageID))
    
    return nil
}
```

### 4.4 事件溯源和最终一致性的实现方案

#### 1. 事件溯源（Event Sourcing）

```go
// internal/infrastructure/eventstore/event_sourcing.go
package eventstore

import (
    "context"
    "fmt"
    
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// EventSourcingRepository 事件溯源仓储
type EventSourcingRepository struct {
    eventRepo EventRepository
}

func NewEventSourcingRepository(eventRepo EventRepository) *EventSourcingRepository {
    return &EventSourcingRepository{
        eventRepo: eventRepo,
    }
}

// ReconstructAggregate 从事件重建聚合根
func (r *EventSourcingRepository) ReconstructAggregate(
    ctx context.Context,
    aggregateID vo.UserID,
) (*aggregate.User, error) {
    // 1. 获取所有事件
    events, err := r.eventRepo.GetEventsByAggregateID(ctx, aggregateID)
    if err != nil {
        return nil, fmt.Errorf("failed to get events: %w", err)
    }
    
    if len(events) == 0 {
        return nil, kernel.ErrAggregateNotFound
    }
    
    // 2. 从事件重建聚合根
    user := &aggregate.User{}
    user.SetID(aggregateID)
    
    for _, event := range events {
        // 根据事件类型应用状态变更
        r.applyEvent(user, event)
    }
    
    return user, nil
}

// applyEvent 应用事件到聚合根
func (r *EventSourcingRepository) applyEvent(user *aggregate.User, event kernel.DomainEvent) {
    switch e := event.(type) {
    case *UserRegisteredEvent:
        // 应用用户注册事件
        // user.username = ...
        // user.email = ...
        // user.status = ...
    case *UserPasswordChangedEvent:
        // 应用密码修改事件
        // user.password = ...
    case *UserLoggedInEvent:
        // 应用登录事件
        // user.lastLoginAt = ...
    }
}

// GetEventHistory 获取事件历史
func (r *EventSourcingRepository) GetEventHistory(
    ctx context.Context,
    aggregateID vo.UserID,
) ([]kernel.DomainEvent, error) {
    return r.eventRepo.GetEventsByAggregateID(ctx, aggregateID)
}
```

#### 2. 最终一致性实现

```go
// internal/application/event/eventual_consistency.go
package event

import (
    "context"
    "fmt"
    "time"
    
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "go.uber.org/zap"
)

// EventualConsistencyManager 最终一致性管理器
type EventualConsistencyManager struct {
    logger              *zap.Logger
    handlerRegistry     EventHandlerRegistry
    idempotencyService  IdempotencyService
    deadLetterQueue     DeadLetterQueue
    maxRetries          int
    retryBackoff        time.Duration
}

func NewEventualConsistencyManager(
    logger *zap.Logger,
    handlerRegistry EventHandlerRegistry,
    idempotencyService IdempotencyService,
    deadLetterQueue DeadLetterQueue,
) *EventualConsistencyManager {
    return &EventualConsistencyManager{
        logger:             logger,
        handlerRegistry:    handlerRegistry,
        idempotencyService: idempotencyService,
        deadLetterQueue:    deadLetterQueue,
        maxRetries:         10,
        retryBackoff:       5 * time.Second,
    }
}

// ProcessEventWithRetry 处理事件并支持重试
func (m *EventualConsistencyManager) ProcessEventWithRetry(
    ctx context.Context,
    event kernel.DomainEvent,
) error {
    // 1. 检查幂等性
    processed, err := m.idempotencyService.IsProcessed(ctx, event.EventName(), event.AggregateID())
    if err != nil {
        m.logger.Error("Failed to check idempotency", zap.Error(err))
        return err
    }
    
    if processed {
        m.logger.Info("Event already processed, skipping",
            zap.String("event_type", event.EventName()),
        )
        return nil
    }
    
    // 2. 处理事件（带重试）
    var lastErr error
    for attempt := 0; attempt < m.maxRetries; attempt++ {
        err := m.handlerRegistry.Dispatch(ctx, event)
        if err == nil {
            // 处理成功，标记为已处理
            if err := m.idempotencyService.MarkAsProcessed(ctx, event.EventName(), event.AggregateID()); err != nil {
                m.logger.Error("Failed to mark as processed", zap.Error(err))
                return err
            }
            
            m.logger.Info("Event processed successfully",
                zap.String("event_type", event.EventName()),
                zap.Int("attempt", attempt+1),
            )
            return nil
        }
        
        lastErr = err
        m.logger.Warn("Event processing failed, retrying",
            zap.String("event_type", event.EventName()),
            zap.Int("attempt", attempt+1),
            zap.Int("max_retries", m.maxRetries),
            zap.Error(err),
        )
        
        // 指数退避
        if attempt < m.maxRetries-1 {
            time.Sleep(m.retryBackoff * time.Duration(attempt+1))
        }
    }
    
    // 3. 所有重试都失败，发送到死信队列
    m.logger.Error("Event processing failed after all retries",
        zap.String("event_type", event.EventName()),
        zap.Int("max_retries", m.maxRetries),
        zap.Error(lastErr),
    )
    
    if err := m.deadLetterQueue.Send(ctx, &DeadLetterMessage{
        EventType:   event.EventName(),
        AggregateID: fmt.Sprintf("%v", event.AggregateID()),
        EventData:   fmt.Sprintf("%v", event),
        Error:       lastErr.Error(),
        RetryCount:  m.maxRetries,
    }); err != nil {
        m.logger.Error("Failed to send to dead letter queue", zap.Error(err))
    }
    
    return lastErr
}

// ProcessEventWithTimeout 处理事件并支持超时
func (m *EventualConsistencyManager) ProcessEventWithTimeout(
    ctx context.Context,
    event kernel.DomainEvent,
    timeout time.Duration,
) error {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    return m.ProcessEventWithRetry(ctx, event)
}
```

---

## 五、完整代码示例

### 5.1 完整的用户注册流程

```go
// 1. 聚合根产生事件
user, err := aggregate.NewUser(username, email, hashedPassword, idGenerator)
if err != nil {
    return nil, err
}

// 2. 应用服务保存聚合根和事件
err = uc.uow.Transaction(ctx, func(ctx context.Context) error {
    // 保存用户
    if err := userRepo.Save(ctx, user); err != nil {
        return err
    }
    
    // 发布事件（保存到 Outbox）
    events := user.GetUncommittedEvents()
    for _, event := range events {
        if err := eventPublisher.Publish(ctx, event); err != nil {
            return err
        }
    }
    
    return nil
})

// 3. Outbox 处理器轮询并发布事件
outboxProcessor.Start(ctx)

// 4. Worker 处理事件
eventProcessor.ProcessDomainEvent(ctx, task)

// 5. 事件处理器执行副作用
userRegisteredHandler.Handle(ctx, event)
```

### 5.2 事件处理器注册

```go
// internal/cmd/worker/main.go
package main

import (
    "context"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/cmd/shared"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/eventstore"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/worker"
    userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
)

func main() {
    infra, logger, cleanup := shared.Initialize("worker")
    defer cleanup()
    
    // 1. 创建事件处理器注册表
    registry := eventstore.NewEventHandlerRegistry()
    
    // 2. 创建事件处理器工厂
    factory := eventstore.NewEventHandlerFactory(logger)
    
    // 3. 注册所有处理器
    handlers := factory.CreateHandlers()
    for _, handler := range handlers {
        for _, eventType := range handler.SupportedEventTypes() {
            registry.Register(eventType, handler)
        }
    }
    
    // 4. 创建事件处理器
    eventProcessor := worker.NewEventProcessor(logger, registry)
    
    // 5. 启动 Asynq Worker
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: infra.Config.Redis.Addr},
        asynq.Config{
            Concurrency: 10,
            Queues: map[string]int{
                "critical": 6,
                "default":  3,
                "low":      1,
            },
        },
    )
    
    mux := asynq.NewServeMux()
    mux.HandleFunc("domain:event", eventProcessor.ProcessDomainEvent)
    
    if err := srv.Run(mux); err != nil {
        logger.Fatal("failed to run worker", zap.Error(err))
    }
}
```

---

## 六、最佳实践总结

### 6.1 事件设计最佳实践

| 原则 | 说明 | 示例 |
|------|------|------|
| **过去时态** | 事件名称使用过去时态 | ✅ UserRegistered ❌ RegisterUser |
| **完整数据** | 事件包含处理所需的所有数据 | ✅ 包含 Email ❌ 只有 UserID |
| **不可变** | 事件创建后不可修改 | 使用值对象，避免 setter |
| **版本管理** | 支持事件结构演变 | 使用 Version 字段 |
| **元数据** | 记录事件上下文信息 | 事件类型、安全标记等 |

### 6.2 事件处理最佳实践

| 原则 | 说明 | 实现 |
|------|------|------|
| **幂等性** | 重复处理同一事件不产生副作用 | 使用 IdempotencyService |
| **异步处理** | 事件处理不阻塞主流程 | 使用 Asynq 消息队列 |
| **错误处理** | 处理失败进入死信队列 | 使用 DeadLetterQueue |
| **重试机制** | 支持指数退避重试 | 配置 MaxRetries 和 RetryBackoff |
| **超时控制** | 防止处理器无限等待 | 使用 context.WithTimeout |

### 6.3 架构最佳实践

| 方面 | 建议 | 理由 |
|------|------|------|
| **分层清晰** | 严格遵循四层架构 | 便于维护和扩展 |
| **依赖反转** | Domain 定义接口，Infrastructure 实现 | 解耦，易于测试 |
| **事务一致性** | 使用 Outbox Pattern | 保证业务操作和事件发布的原子性 |
| **模块化** | 按限界上下文划分模块 | 便于团队协作 |
| **文档完整** | 记录设计决策和架构图 | 便于新人理解 |

---

## 七、后续优化路线图

### Phase 1（立即实施）
- [ ] 完善事件处理器注册机制
- [ ] 实现幂等性检查
- [ ] 添加死信队列处理
- [ ] 改进事件数据完整性

### Phase 2（1-2 周）
- [ ] 实现事件版本管理
- [ ] 添加事件溯源支持
- [ ] 完善错误处理和重试机制
- [ ] 添加事件处理监控

### Phase 3（2-4 周）
- [ ] 实现 Saga 模式（分布式事务）
- [ ] 添加事件快照机制
- [ ] 支持事件投影（CQRS）
- [ ] 性能优化和压力测试

---

**分析完成时间：** 2024-03-26 14:30  
**下一步行动：** 根据优化建议逐步实施改进
