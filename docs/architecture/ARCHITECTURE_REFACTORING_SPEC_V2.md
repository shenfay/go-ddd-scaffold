# DDD-Scaffold 架构重构规范（2026 版）

## 概述

本文档描述 `ddd-scaffold` 项目从简洁实用主义垂直切片架构演进为**事件驱动 + 异步处理**的完整规范。所有重构工作必须严格遵循本文档。

---

## 一、架构演进的驱动力

### 1.1 当前架构问题

```
internal/auth/
├── domain.go           (User 实体包含认证信息)
├── handler.go          (HTTP Handler - 668 行，过于臃肿)
├── service.go          (应用服务)
├── token_service.go    (Token 管理)
├── repository.go       (仓储接口)
├── repository_gorm.go  (GORM 实现)
├── events.go           (领域事件)
└── tasks.go            (Asynq 任务)
```

**存在的问题**：
1. **职责混乱**：所有功能都在 `auth/` 下，违反单一职责原则
2. **同步阻塞**：日志记录直接写入数据库，影响主流程性能
3. **耦合严重**：业务代码依赖日志服务，难以测试和维护
4. **边界模糊**：用户管理、认证、Token、设备管理混在一起

---

### 1.2 演进目标

| 维度 | 当前状态 | 目标状态 | 收益 |
|------|---------|---------|------|
| **日志记录** | 同步写入 DB | ✅ 事件驱动 + 异步 | 性能提升 75% |
| **职责划分** | 大杂烩 | ✅ 清晰的 DDD 边界 | 可维护性↑ |
| **解耦程度** | 紧耦合 | ✅ 完全解耦 | 可测试性↑ |
| **扩展性** | 垂直扩展 | ✅ 水平扩展 | 并发能力↑ |

---

## 二、目标架构

### 2.1 架构分层

```
Domain Layer     → 纯业务逻辑（Entity, Value Object, Domain Event）
     ↓
Application Layer → 用例编排（发布领域事件）
     ↓
Listener Layer   → 事件监听（转换为日志任务）
     ↓
Transport Layer  → HTTP Handler + Worker Processor
     ↓
Infrastructure  → 数据库、Redis、消息队列
```

**核心原则**：
- ✅ **依赖倒置**：Domain 定义接口，Infra 实现
- ✅ **事件驱动**：通过 EventBus 异步通信
- ✅ **职责分离**：每层有明确的职责边界
- ✅ **简洁实用**：避免过度设计

---

### 2.2 目录结构

```
backend/
├── internal/
│   ├── domain/                    # 领域层
│   │   ├── user/                  # 用户核心域
│   │   │   ├── entity.go          # User 聚合根
│   │   │   │                      # - ID, Email, Password
│   │   │   │                      # - Locked, FailedAttempts
│   │   │   ├── value_objects.go   # Email, Password 值对象
│   │   │   ├── repository.go      # UserRepository 接口
│   │   │   └── events.go          # 用户事件
│   │   │                          #   - UserLoggedIn
│   │   │                          #   - UserLoggedOut
│   │   │
│   │   └── shared/                # 共享领域概念
│   │       ├── id.go              # ID 生成器（ULID）
│   │       ├── base_event.go      # 领域事件基类
│   │       └── errors.go          # 通用错误
│   │
│   ├── app/                       # 应用层
│   │   ├── user/                  # 用户应用服务
│   │   │   ├── service.go         # UserService
│   │   │   ├── dto.go             # DTO
│   │   │   └── mapper.go          # Entity ↔ DTO
│   │   │
│   │   └── authentication/        # 认证应用服务
│   │       ├── service.go         # AuthService
│   │       │                      # - Login(), Logout()
│   │       │                      # - RefreshToken()
│   │       ├── dto.go
│   │       └── mapper.go
│   │
│   ├── listener/                  # 🆕 事件监听器层
│   │   ├── audit_log_listener.go  # ✅ 审计日志监听器
│   │   │                          # • 订阅 AUTH.*, SECURITY.* 事件
│   │   │                          # • 转换为审计日志任务
│   │   │                          # • 发布到 Worker 队列（critical）
│   │   │
│   │   └── activity_log_listener.go # ✅ 活动日志监听器
│   │                                # • 订阅 USER.*, FEATURE.* 事件
│   │                                # • 转换为活动日志任务
│   │                                # • 发布到 Worker 队列（default）
│   │
│   ├── transport/                 # 传输层
│   │   ├── http/
│   │   │   ├── router.go          # 路由配置
│   │   │   ├── middleware.go      # 中间件
│   │   │   └── handlers/
│   │   │       ├── auth.go        # AuthHandler
│   │   │       ├── user.go        # UserHandler
│   │   │       └── audit_log.go   # AuditLogHandler（查询）
│   │   │
│   │   └── worker/                # Worker 服务
│   │       └── handlers/
│   │           ├── audit_log_handler.go    # 审计日志处理器
│   │           └── activity_log_handler.go # 活动日志处理器
│   │
│   └── infra/                     # 基础设施层
│       ├── persistence/
│       │   ├── db.go              # GORM 数据库连接
│       │   └── redis.go           # Redis 连接
│       │
│       ├── repository/            # 仓储实现
│       │   ├── user.go            # UserRepository GORM
│       │   ├── audit_log.go       # AuditLogRepository GORM
│       │   └── activity_log.go    # ActivityLogPersistence
│       │
│       ├── messaging/             # 消息队列
│       │   └── event_bus.go       # ✅ EventBus 接口 + 工厂
│       │                          # • type EventBus interface { ... }
│       │                          # • func NewEventBus(...) EventBus
│       │                          # • 内部调用 asynqEventBus
│       │
│       └── redis/                 # Redis 存储
│           ├── token_store.go     # Token 存储（Redis）
│           └── device_store.go    # Device 存储（Redis）
│
├── pkg/
│   └── event/
│       ├── event.go               # Event 基础接口
│       └── asynq_event_bus.go     # （已废弃，迁移到 infra/messaging）
│
└── migrations/
    ├── 001_create_users_table.sql
    ├── 005_create_audit_logs_table.sql
    └── 006_create_activity_logs_table.sql
```

---

## 三、数据库表结构设计

### 3.1 现有表结构

#### ✅ users 表（用户核心表）

```sql
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,      -- 密码哈希
    email_verified BOOLEAN DEFAULT FALSE,
    locked BOOLEAN DEFAULT FALSE,
    failed_attempts INT DEFAULT 0,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_users_email ON users(email);
```

**说明**：
- ✅ **包含所有认证信息**：password, locked, failed_attempts
- ✅ **符合 DDD**：User 是聚合根，认证信息是值对象
- ✅ **高性能**：不需要 JOIN 查询 credentials

---

#### ✅ audit_logs 表（审计日志）

```sql
CREATE TABLE audit_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    action VARCHAR(50) NOT NULL,           -- AUTH.*, SECURITY.*
    status VARCHAR(20) NOT NULL,           -- SUCCESS / FAILED
    ip VARCHAR(45),
    user_agent VARCHAR(500),
    device VARCHAR(100),
    browser VARCHAR(50),
    os VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at_desc ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_status ON audit_logs(status);
CREATE INDEX idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_action_created ON audit_logs(action, created_at DESC);
CREATE INDEX idx_audit_logs_failed ON audit_logs(status, created_at DESC) 
    WHERE status = 'FAILED';
```

**用途**：
- ✅ 安全审计、合规检查
- ✅ 保存期限：1-7 年
- ✅ 不可篡改、详细完整

---

#### ✅ activity_logs 表（活动日志 - 简化版）

```sql
CREATE TABLE activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,           -- USER.*, FEATURE.*
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at_desc ON activity_logs(created_at DESC);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_user_created ON activity_logs(user_id, created_at DESC);
```

**用途**：
- ✅ 产品分析、用户体验优化
- ✅ 保存期限：30-90 天
- ✅ 轻量级、便于统计

---

#### ✅ email_verification_tokens 表

```sql
CREATE TABLE email_verification_tokens (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

#### ✅ password_reset_tokens 表

```sql
CREATE TABLE password_reset_tokens (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

### 3.2 不需要的表（使用 Redis）

#### ❌ credentials 表（不需要）

**原因**：
- ✅ User 实体已包含所有认证信息
- ✅ 符合 DDD：Credentials 是值对象，不是独立实体
- ✅ 高性能：不需要 JOIN

---

#### ❌ tokens 表（不需要）

**替代方案**：使用 Redis 存储

```go
// infra/redis/token_store.go
package redis

type TokenStore struct {
    client *redis.Client
}

// Key: auth:token:{refresh_token}
// Value: {
//   "user_id": "user123",
//   "access_token": "eyJ...",
//   "expires_at": "2024-04-04T10:00:00Z"
// }
// TTL: 7 days

func (s *TokenStore) Store(ctx context.Context, refreshToken string, data *TokenData) error {
    key := "auth:token:" + refreshToken
    value, _ := json.Marshal(data)
    return s.client.Set(ctx, key, value, 7*24*time.Hour).Err()
}

func (s *TokenStore) Get(ctx context.Context, refreshToken string) (*TokenData, error) {
    key := "auth:token:" + refreshToken
    value, err := s.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, ErrTokenNotFound
    }
    var data TokenData
    json.Unmarshal(value, &data)
    return &data, nil
}

func (s *TokenStore) Delete(ctx context.Context, refreshToken string) error {
    key := "auth:token:" + refreshToken
    return s.client.Del(ctx, key).Err()
}
```

**优势**：
- ✅ 高性能：O(1) 查询
- ✅ 自动过期：Redis TTL 机制
- ✅ 黑名单：登出时直接删除

---

#### ❌ devices 表（不需要）

**替代方案**：使用 Redis 存储

```go
// infra/redis/device_store.go
package redis

type DeviceStore struct {
    client *redis.Client
}

// Key: auth:devices:{user_id}
// Value: [
//   {
//     "device_id": "device123",
//     "device_name": "Chrome on macOS",
//     "device_type": "desktop",
//     "last_active_at": "2024-04-03T10:00:00Z"
//   }
// ]
// TTL: 30 days

func (s *DeviceStore) AddDevice(ctx context.Context, userID string, device *Device) error {
    key := "auth:devices:" + userID
    devices, _ := s.GetDevices(ctx, userID)
    
    // 更新或添加设备
    exists := false
    for i, d := range devices {
        if d.DeviceID == device.DeviceID {
            devices[i] = device
            exists = true
            break
        }
    }
    if !exists {
        devices = append(devices, device)
    }
    
    // 限制最多 10 个设备
    if len(devices) > 10 {
        devices = devices[:10]
    }
    
    value, _ := json.Marshal(devices)
    return s.client.Set(ctx, key, value, 30*24*time.Hour).Err()
}

func (s *DeviceStore) GetDevices(ctx context.Context, userID string) ([]*Device, error) {
    key := "auth:devices:" + userID
    value, err := s.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return []*Device{}, nil
    }
    var devices []*Device
    json.Unmarshal(value, &devices)
    return devices, nil
}
```

**优势**：
- ✅ 高频读取、低延迟
- ✅ 最终一致性可接受
- ✅ 自动清理不活跃设备

---

## 四、事件驱动架构

### 4.1 事件流转

```
HTTP Request
    ↓
Transport (Handler)
    ↓
App Service → eventBus.Publish(UserLoggedIn) ← 立即返回
    ↓
Domain Event
    ↓
EventBus (Asynq Client) → Redis Queue ← 异步
    ↓
Listener (订阅者)
    ↓
Worker (消费者)
    ↓
Repository (GORM)
    ↓
Database (audit_logs)
```

---

### 4.2 EventBus 实现

```go
// infra/messaging/event_bus.go
package messaging

import (
    "context"
    "encoding/json"
    "strings"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// EventBus 事件总线接口
type EventBus interface {
    Publish(ctx context.Context, evt event.Event) error
    Subscribe(eventType string, handler event.EventHandler)
}

// asynqEventBus EventBus 的 Asynq 实现
type asynqEventBus struct {
    client *asynq.Client
}

// NewEventBus 工厂方法（统一入口）
func NewEventBus(redisAddr string, queueConfig QueueConfig) EventBus {
    client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
    return &asynqEventBus{client: client}
}

// Publish 发布事件到 Asynq 队列
func (b *asynqEventBus) Publish(ctx context.Context, evt event.Event) error {
    payload, _ := json.Marshal(evt.GetPayload())
    
    _, err := b.client.EnqueueContext(ctx,
        asynq.NewTask(evt.GetType(), payload),
        asynq.Queue(b.getQueueForEvent(evt.GetType())),
    )
    
    return err
}

// getQueueForEvent 根据事件类型选择队列
func (b *asynqEventBus) getQueueForEvent(eventType string) string {
    // 审计日志 → critical 队列（高优先级）
    if strings.HasPrefix(eventType, "AUTH.") || strings.HasPrefix(eventType, "SECURITY.") {
        return "critical"
    }
    // 活动日志 → default 队列
    return "default"
}

// Subscribe 订阅事件（由 Listener 调用）
func (b *asynqEventBus) Subscribe(eventType string, handler event.EventHandler) {
    // Listener 负责注册到 Worker 的 ServeMux
    // 这里只是标记已订阅
}
```

---

### 4.3 Listener 实现

```go
// listener/audit_log_listener.go
package listener

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/events"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
    "github.com/shenfay/go-ddd-scaffold/pkg/event"
)

type AuditLogListener struct {
    eventBus messaging.EventBus
}

func NewAuditLogListener(eventBus messaging.EventBus) *AuditLogListener {
    l := &AuditLogListener{eventBus}
    
    // 订阅认证相关事件
    eventBus.Subscribe("AUTH.LOGIN.SUCCESS", l.HandleUserLoggedIn)
    eventBus.Subscribe("AUTH.LOGIN.FAILED", l.HandleLoginFailed)
    eventBus.Subscribe("SECURITY.ACCOUNT.LOCKED", l.HandleAccountLocked)
    
    return l
}

// HandleUserLoggedIn 处理用户登录事件
func (l *AuditLogListener) HandleUserLoggedIn(ctx context.Context, evt event.Event) error {
    e := evt.(*events.UserLoggedIn)
    
    // 转换为审计日志任务并发布到 Worker 队列
    return l.eventBus.Publish(ctx, &AuditLogTask{
        Action: "AUTH.LOGIN.SUCCESS",
        Status: "SUCCESS",
        Data: map[string]interface{}{
            "user_id":    e.UserID,
            "email":      e.Email,
            "ip":         e.IP,
            "user_agent": e.UserAgent,
            "device":     e.Device,
        },
    })
}

func (l *AuditLogListener) HandleLoginFailed(ctx context.Context, evt event.Event) error {
    e := evt.(*events.LoginFailed)
    
    return l.eventBus.Publish(ctx, &AuditLogTask{
        Action: "AUTH.LOGIN.FAILED",
        Status: "FAILED",
        Data: map[string]interface{}{
            "user_id": e.UserID,
            "email":   e.Email,
            "ip":      e.IP,
            "reason":  e.Reason,
        },
    })
}
```

---

### 4.4 Worker Handler 实现

```go
// worker/handlers/audit_log_handler.go
package handlers

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
)

type AuditLogHandler struct {
    repo *repository.AuditLogRepository
}

func NewAuditLogHandler(repo *repository.AuditLogRepository) *AuditLogHandler {
    return &AuditLogHandler{repo}
}

// ProcessTask 处理审计日志任务
func (h *AuditLogHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
    var data AuditLogTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        return err
    }
    
    log := &AuditLog{
        UserID:   data.Data["user_id"].(string),
        Action:   data.Action,
        Status:   data.Status,
        Metadata: data.Data,
    }
    
    return h.repo.Save(ctx, log)
}
```

---

## 五、命名规范

### 5.1 目录命名

| 层级 | 命名规则 | 示例 |
|------|---------|------|
| 顶层目录 | 小写，可缩写 | `domain`, `app`, `listener`, `transport`, `infra` |
| 业务域 | 小写，单数 | `user`, `authentication` |
| 共享目录 | 小写 | `shared`, `common` |
| 集合类 | 复数 | `handlers`, `repositories` |

### 5.2 文件命名

| 类型 | 规则 | 示例 |
|------|------|------|
| 聚合根/实体 | 单数 | `entity.go` |
| 值对象集合 | 复数 | `value_objects.go` |
| 事件集合 | 复数 | `events.go` |
| 服务 | 单数 | `service.go` |
| 仓储接口 | 单数 | `repository.go` |
| 仓储实现 | 简洁形式 | `user.go` (不用 `_gorm.go`) |
| EventBus | 简洁形式 | `event_bus.go` (不用 `_asynq.go`) |
| Listener | `{domain}_listener.go` | `audit_log_listener.go` |
| Handler | `{domain}_handler.go` | `audit_log_handler.go` |

### 5.3 事件命名

**格式**：`{DOMAIN}.{CATEGORY}.{ACTION}`

```go
// ✅ 正确
"AUTH.LOGIN.SUCCESS"
"AUTH.LOGIN.FAILED"
"SECURITY.ACCOUNT.LOCKED"
"USER.PROFILE.UPDATED"

// ❌ 错误
"user_registered"        # 不能用下划线
"userLoggedIn"           # 不能用驼峰
"UserRegistered"         # 不能大写
```

---

## 六、实施步骤

详见 [`LOGGING_MIGRATION_PLAN.md`](./LOGGING_MIGRATION_PLAN.md)

### 阶段概览

1. **准备工作**（1 小时）
2. **创建目录结构**（30 分钟）
3. **实现 EventBus**（2 小时）
4. **创建 Listener**（3 小时）
5. **重构 App 层**（2 小时）
6. **更新 Worker**（2 小时）
7. **初始化 Listener**（1 小时）
8. **测试验证**（2 小时）
9. **清理文档**（1 小时）

**总工时**：约 14 小时（建议 2 个工作日）

---

## 七、验收标准

### 7.1 功能验收

- [ ] 用户登录立即返回 Token（< 50ms）
- [ ] Worker 在 1 秒内写入 audit_logs
- [ ] 查询审计日志接口正常
- [ ] 所有现有测试通过

### 7.2 性能验收

- [ ] 登录接口 P99 < 100ms
- [ ] 吞吐量 > 1000 QPS
- [ ] Worker 消费延迟 < 1 秒

### 7.3 架构验收

- [ ] 业务代码不直接调用日志服务
- [ ] Listener 完全解耦领域事件和日志记录
- [ ] EventBus 可轻松替换实现（如切换到 Kafka）
- [ ] 无循环依赖
- [ ] 测试覆盖率 > 70%

---

## 八、关键设计决策

### 8.1 为什么不拆分 Credentials 表？

**决策**：User 聚合根包含所有认证信息

**理由**：
- ✅ 符合 DDD：Credentials 是值对象，不是独立实体
- ✅ 高性能：不需要 JOIN 查询
- ✅ 简洁：减少不必要的表
- ✅ 事务一致性：在同一事务中更新

---

### 8.2 为什么用 Redis 存储 Token？

**决策**：使用 Redis，不用数据库表

**理由**：
- ✅ 高性能：O(1) 查询
- ✅ 自动过期：Redis TTL 机制
- ✅ 黑名单：登出时直接删除 Key
- ✅ 水平扩展：Redis Cluster

---

### 8.3 为什么需要 Listener 层？

**决策**：新增 Listener 层，负责事件转换和路由

**理由**：
- ✅ 职责分离：App 层只发布事件，不关心如何处理
- ✅ 解耦：Listener 可以独立测试和修改
- ✅ 灵活：可以轻松添加新的事件处理器

---

### 8.4 为什么 EventBus 文件名要抽象？

**决策**：使用 `event_bus.go`，不用 `asynq_event_bus.go`

**理由**：
- ✅ 框架无关：如果换消息队列，文件名不用改
- ✅ 依赖抽象：符合 DDD 原则
- ✅ 易于 Mock：测试时可以使用 FakeEventBus

---

## 九、参考文档

- [`DATABASE_SCHEMA_DESIGN.md`](./DATABASE_SCHEMA_DESIGN.md) - 数据库表结构设计
- [`LOGGING_MIGRATION_PLAN.md`](./LOGGING_MIGRATION_PLAN.md) - 迁移计划
- [`ARCHITECTURE_REFACTORING_SPEC.md`](./ARCHITECTURE_REFACTORING_SPEC.md) - 原规范文档

---

**文档版本**：v2.0  
**创建日期**：2026-04-03  
**最后更新**：2026-04-03  
**状态**：已批准  
**作者**：DDD-Scaffold Team
