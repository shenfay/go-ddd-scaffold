# 领域事件架构设计

## 架构概述

本项目采用 **事件溯源 + 独立任务队列** 的混合架构：

- **EventStore**: 纯事件溯源，仅记录历史事件用于审计和回放
- **asynq**: 独立的任务队列系统，负责任务调度和执行
- **asynqmon**: 独立的任务监控 UI，类似 Python 的 Flower

## 架构图

```
┌─────────────────────────────────────────────────┐
│  业务层 (Domain/Application)                     │
│     ↓ 发布事件                                   │
├─────────────────────────────────────────────────┤
│  EventStore (事件溯源)                           │
│  - 保存到 domain_events 表                        │
│  - 仅用于审计、历史追溯                          │
│  - 不包含状态字段                                │
├─────────────────────────────────────────────────┤
│  asynq (任务队列)                                │
│  - 事件处理器将事件推送到 Redis                  │
│  - asynq 负责任务调度和执行                      │
│  - asynqmon 负责监控                            │
└─────────────────────────────────────────────────┘
```

## 数据流

### 事件发布流程

```
用户注册 → UserCreated 事件
    ↓
EventStore.SaveEvents() → domain_events 表（纯历史记录）
    ↓
AsynqEventPublisher.Publish() → Redis (asynq 队列)
    ↓
asynq Worker → 执行任务（发送邮件等）
    ↓
asynqmon → 监控任务状态（独立服务）
```

## 核心组件

### 1. EventStore（纯事件溯源）

**职责**：
- 保存聚合根的所有历史事件
- 支持事件回放（重建聚合根状态）
- 提供审计和历史查询能力

**接口定义**：
```go
type EventStore interface {
    // SaveEvents 保存聚合根的未提交事件（仅用于溯源，不跟踪状态）
    SaveEvents(ctx context.Context, aggregateID string, aggregateType string, events []kernel.DomainEvent) error
    
    // GetEvents 获取聚合根的所有历史事件（用于事件回放）
    GetEvents(ctx context.Context, aggregateID string) ([]*EventRecord, error)
    
    // GetEventsByType 按类型获取事件（用于分析和审计）
    GetEventsByType(ctx context.Context, eventType string, limit int) ([]*EventRecord, error)
}
```

**domain_events 表结构**：
```sql
CREATE TABLE domain_events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_version INTEGER NOT NULL,
    event_data JSONB NOT NULL,
    occurred_on TIMESTAMP NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_id, aggregate_type);
CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_occurred ON domain_events(occurred_on DESC);
```

### 2. asynq 任务队列

**职责**：
- 接收事件发布器推送的任务
- 负责任务的调度、重试、执行
- 通过 Redis 进行任务存储和状态管理

**核心文件**：
- `asynq_client.go`: Client 和 Server 配置
- `asynq_task.go`: 任务定义和序列化
- `asynq_processor.go`: 任务处理器
- `asynq_publisher.go`: 事件发布器

**任务优先级队列**：
```go
Queues: map[string]int{
    "critical": 6, // 高优先级（如用户注册、支付）
    "default":  3, // 默认队列
    "low":      1, // 低优先级（如日志、通知）
}
```

### 3. asynqmon 监控

**功能**：
- 实时任务监控仪表盘
- 任务详情查看
- 失败任务重试
- 任务归档管理
- 统计图表展示

**部署方式**：
```bash
# 使用 Docker Compose 启动
docker-compose -f deployments/docker/docker-compose.asynq.yml up -d
```

**访问地址**：
- asynqmon UI: http://localhost:8081
- Redis: localhost:6379
- PostgreSQL: localhost:5432

## 使用示例

### 发布领域事件

```go
// 在应用服务中
event := user.NewUserCreatedEvent(user.ID, user.Email)
err := eventPublisher.Publish(ctx, event)
```

### 创建事件处理器

```go
// 实现 DomainEventHandler 接口
type WelcomeEmailHandler struct {
    emailService *email.Service
}

func (h *WelcomeEmailHandler) CanHandle(eventType string) bool {
    return eventType == "UserCreated"
}

func (h *WelcomeEmailHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
    // 发送邮件逻辑
    return nil
}
```

### 配置 asynq Worker

```go
// 在 bootstrap 中配置
asynqServer := messaging.NewAsynqServer(messaging.AsynqConfig{
    RedisAddr:     cfg.Redis.Addr,
    RedisPassword: cfg.Redis.Password,
    RedisDB:       cfg.Redis.DB,
})

processor := messaging.NewAsynqProcessor(logger, welcomeEmailHandler)
asynqServer.RegisterHandler(processor.ProcessTask)

// 启动 worker
go asynqServer.Run()
```

## 优势对比

### 原架构（使用 processed 字段）

❌ 缺点：
- 事件存储与队列耦合
- 需要手动管理事件状态
- 无法利用专业的任务调度系统
- 监控能力弱

### 新架构（EventStore + asynq）

✅ 优点：
- **职责分离**：EventStore 专注溯源，asynq 专注调度
- **零维护成本**：asynq 自动处理重试、超时、归档
- **强大监控**：asynqmon 提供完整的可视化界面
- **可扩展性**：支持分布式部署（虽然当前是单机）
- **灵活性**：可以轻松切换队列后端或调整策略

## 运维指南

### 数据库迁移

```bash
# 执行迁移（移除 processed 字段）
./tools/migrator/migrate.sh up 000011

# 回滚（恢复 processed 字段）
./tools/migrator/migrate.sh down 000011
```

### 启动 asynqmon

```bash
# 启动 Redis 和 asynqmon
docker-compose -f deployments/docker/docker-compose.asynq.yml up -d

# 查看日志
docker-compose logs -f asynqmon
```

### 监控指标

通过 asynqmon UI 可以查看：
- 各队列的任务数量
- 任务处理成功率和平均时间
- 失败任务和错误信息
- 历史趋势图表

## 注意事项

1. **不要直接修改 domain_events 表的状态**：该表现在是纯历史记录
2. **任务状态由 asynq 管理**：通过 asynqmon 或 API 进行操作
3. **事件回放**：使用 `GetEvents` 方法重建聚合根状态
4. **Redis 持久化**：建议开启 Redis AOF 保证任务不丢失

## 参考资料

- [asynq 官方文档](https://github.com/hibiken/asynq)
- [asynqmon 使用指南](https://github.com/hibiken/asynqmon)
- [事件溯源模式](https://martinfowler.com/eaaDev/EventSourcing.html)
