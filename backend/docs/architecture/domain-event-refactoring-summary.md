# 领域事件架构重构总结

## 重构概述

本次重构将原有的 `messaging` 包按领域概念拆分为两个独立的目录：
- **domain_event** - 领域事件相关（EventStore、EventBus、Publisher）
- **task_queue** - 任务队列相关（asynq Client、Server、Processor）

## 新的目录结构

```
backend/internal/infrastructure/
├── domain_event/              # 领域事件基础设施
│   ├── store.go              # EventStore 接口定义
│   ├── repository.go         # DomainEventRepository 实现
│   ├── publisher.go          # AsynqPublisher 实现
│   └── bus.go                # EventBus 实现（Simple/Async）
│
├── task_queue/               # 任务队列基础设施
│   ├── client.go             # asynq Client/Server 配置
│   ├── publisher.go          # Publisher 封装 + 任务定义
│   └── processor.go          # Processor 处理器 + Handler 接口
│
└── messaging/                # 保留的旧目录（仅 event_bus.go）
    └── event_bus.go          # 临时保留，供其他模块迁移使用
```

## 核心变更

### 1. domain_event 包

**职责**：领域事件的存储、发布和总线通信

#### store.go
- `EventStore` 接口 - 纯事件溯源，不包含状态管理
- `EventRecord` 结构 - 事件记录 DTO

#### repository.go  
- `DomainEventRepository` - EventStore 的实现
- `SaveEvents()` - 保存事件到数据库
- `GetEvents()` - 按聚合根查询历史事件
- `GetEventsByType()` - 按类型查询事件

#### publisher.go
- `Publisher` 接口 - 事件发布器抽象
- `AsynqPublisher` - 基于 asynq 的实现
- 支持优先级队列（critical/default/low）

#### bus.go
- `EventBus` 接口 - 事件总线
- `SimpleEventBus` - 同步实现
- `AsyncEventBus` - 异步实现（带 worker 池）
- `EventHandler` 接口 - 事件处理器

### 2. task_queue 包

**职责**：asynq 任务队列的封装和管理

#### client.go
- `Config` - asynq 配置
- `NewClient()` - 创建 asynq Client
- `NewServer()` - 创建 asynq Server
- 支持多优先级队列配置

#### publisher.go
- `Publisher` - asynq 发布器封装
- `DomainEventPayload` - 领域事件任务负载
- `TaskTypeDomainEvent` - 任务类型常量
- `PublishDomainEvent()` - 发布领域事件任务

#### processor.go
- `Processor` - asynq 任务处理器
- `Handler` 接口 - 领域事件处理器
- `ProcessTask()` - 处理 asynq 任务
- `domainEventAdapter` - 适配器模式重建领域事件

## 依赖关系

```
应用层
  ↓
domain_event.Publisher (发布事件)
  ↓
task_queue.Publisher (推送到 Redis)
  ↓
asynq Worker (消费任务)
  ↓
task_queue.Processor (调用 Handler)
  ↓
domain_event.EventHandler (业务逻辑)
```

## 使用示例

### 1. 保存并发布领域事件

```go
// 在 Repository 中
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
    // 保存到数据库
    err := r.saveToDB(ctx, u)
    if err != nil {
        return err
    }
    
    // 发布领域事件（通过 asynq）
    events := u.GetUncommittedEvents()
    for _, event := range events {
        err = r.eventPublisher.Publish(ctx, event)
        if err != nil {
            logger.Error("Failed to publish event", zap.Error(err))
        }
    }
    
    return nil
}
```

### 2. 实现事件处理器

```go
// 实现 task_queue.Handler 接口
type WelcomeEmailHandler struct {
    emailService *email.Service
}

func (h *WelcomeEmailHandler) CanHandle(eventType string) bool {
    return eventType == "UserCreated"
}

func (h *WelcomeEmailHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
    // 处理事件：发送欢迎邮件
    payload := event.(*domainEventAdapter)
    var data UserCreatedData
    json.Unmarshal(payload.eventData, &data)
    
    return h.emailService.SendWelcomeEmail(ctx, data.Email)
}
```

### 3. 注册处理器

```go
// 在 bootstrap 中
processor := task_queue.NewProcessor(logger, welcomeEmailHandler, auditHandler)
asynqServer := task_queue.NewServer(cfg.Asynq)
asynqServer.RegisterHandler(processor.ProcessTask)

// 启动 worker
go asynqServer.Run()
```

## 迁移清单

### 已完成的迁移
- ✅ EventStore 接口 → `domain_event/store.go`
- ✅ DomainEventRepository → `domain_event/repository.go`
- ✅ EventBus → `domain_event/bus.go`
- ✅ AsynqPublisher → `domain_event/publisher.go`
- ✅ asynq Client/Server → `task_queue/client.go`
- ✅ asynq Task/Payload → `task_queue/publisher.go`
- ✅ asynq Processor → `task_queue/processor.go`

### 引用更新
- ✅ `domain_event_repository.go` - 更新为 `domain_event` 包
- ✅ `user_repository.go` - 更新为 `domain_event.EventStore`
- ✅ 其他文件保持使用 `kernel.EventBus`（临时兼容）

## 优势对比

### 原架构（messaging 包）
❌ 职责不清 - EventStore 和 asynq 混在一起  
❌ 命名混乱 - 所有文件都带 `asynq_` 前缀  
❌ 难以理解 - 不符合领域驱动设计理念  

### 新架构（domain_event + task_queue）
✅ 职责清晰 - 领域事件 vs 任务队列分离  
✅ 命名简洁 - 文件名直接反映职责  
✅ 符合 DDD - 按领域概念组织代码  
✅ 易于扩展 - 可以独立替换某个组件  

## 注意事项

1. **messaging/event_bus.go** 暂时保留，因为部分代码仍在使用 `kernel.EventBus`
2. **后续迁移**：可以将 `kernel.EventBus` 统一迁移到 `domain_event.EventBus`
3. **asynq 依赖**：确保已运行 `go mod tidy` 下载 asynq 包
4. **测试验证**：需要测试事件发布和处理的完整流程

## 相关文件

- [架构设计文档](../../docs/architecture/event-driven-architecture.md)
- [README](../../README.md)
- [go.mod](../../go.mod)

---

**重构完成时间**: 2026-03-19  
**重构方案**: 方案 2（按领域概念组织）
