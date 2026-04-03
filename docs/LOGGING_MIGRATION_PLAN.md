# 日志系统事件驱动架构迁移计划

## 概述

本文档描述如何从当前的同步日志记录架构迁移到**事件驱动 + 异步处理**的架构。

---

## 一、架构对比

### 当前架构（同步）

```
HTTP Handler → App Service → auditLogService.Record() 
                                      ↓
                              直接写入数据库（阻塞）
                                      ↓
                              返回响应
```

**问题**：
- ❌ 同步阻塞，影响主流程性能
- ❌ 耦合严重，业务代码依赖日志服务
- ❌ 难以扩展，无法应对高并发场景

---

### 目标架构（事件驱动 + 异步）

```
HTTP Handler → App Service → eventBus.Publish(UserLoggedIn) ← 立即返回
                                      ↓
                              Asynq 队列（异步）
                                      ↓
                              Listener 监听并转换
                                      ↓
                              Worker 消费并写入数据库
```

**优势**：
- ✅ 异步非阻塞，主流程性能提升
- ✅ 完全解耦，业务代码不依赖日志实现
- ✅ 水平扩展，Worker 可独立扩容

---

## 二、迁移步骤

### 阶段 1：准备工作（预计 1 小时）

#### 1.1 创建 Git 分支

```bash
git checkout -b feature/event-driven-logging
```

#### 1.2 备份当前代码

确保当前功能正常运行，记录测试用例通过率。

#### 1.3 安装依赖

```bash
cd backend
go get github.com/hibiken/asynq@latest
go mod tidy
```

---

### 阶段 2：创建新的目录结构（预计 30 分钟）

#### 2.1 创建目录

```bash
cd backend/internal

# 创建 Listener 层
mkdir listener

# 创建 Worker handlers 目录
mkdir -p worker/handlers
```

#### 2.2 调整文件位置

```bash
# 移动现有的 activitylog handler 到 worker
mv internal/activitylog/handler.go worker/handlers/activity_log_handler.go

# 更新包名
sed -i '' 's/package activitylog/package handlers/' worker/handlers/activity_log_handler.go
```

---

### 阶段 3：实现 EventBus（预计 2 小时）

#### 3.1 定义 EventBus 接口

**文件**：`infra/messaging/event_bus.go`

```go
package messaging

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// EventBus 事件总线接口
type EventBus interface {
    Publish(ctx context.Context, evt event.Event) error
    Subscribe(eventType string, handler event.EventHandler)
}

// NewEventBus 工厂方法
func NewEventBus(redisAddr string, queueConfig QueueConfig) EventBus {
    return NewAsynqEventBus(redisAddr, queueConfig)
}
```

#### 3.2 实现 AsynqEventBus

**文件**：`infra/messaging/asynq_event_bus.go`

```go
package messaging

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/shenfay/go-ddd-scaffold/pkg/event"
)

type asynqEventBus struct {
    client *asynq.Client
}

func NewAsynqEventBus(redisAddr string, queueConfig QueueConfig) EventBus {
    client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
    
    return &asynqEventBus{client: client}
}

func (b *asynqEventBus) Publish(ctx context.Context, evt event.Event) error {
    payload, _ := json.Marshal(evt.GetPayload())
    
    _, err := b.client.EnqueueContext(ctx,
        asynq.NewTask(evt.GetType(), payload),
        asynq.Queue(b.getQueueForEvent(evt.GetType())),
    )
    
    return err
}

func (b *asynqEventBus) getQueueForEvent(eventType string) string {
    // 审计日志 → critical 队列
    if strings.HasPrefix(eventType, "AUTH.") || strings.HasPrefix(eventType, "SECURITY.") {
        return "critical"
    }
    // 活动日志 → default 队列
    return "default"
}
```

---

### 阶段 4：创建 Listener（预计 3 小时）

#### 4.1 审计日志监听器

**文件**：`listener/audit_log_listener.go`

```go
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
    // 类似处理...
}
```

#### 4.2 活动日志监听器

**文件**：`listener/activity_log_listener.go`

```go
package listener

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
    "github.com/shenfay/go-ddd-scaffold/pkg/event"
)

type ActivityLogListener struct {
    eventBus messaging.EventBus
}

func NewActivityLogListener(eventBus messaging.EventBus) *ActivityLogListener {
    l := &ActivityLogListener{eventBus}
    
    // 订阅用户行为事件
    eventBus.Subscribe("USER.PROFILE.UPDATED", l.HandleProfileUpdated)
    eventBus.Subscribe("FEATURE.EXPORT_USED", l.HandleExportUsed)
    
    return l
}

func (l *ActivityLogListener) HandleProfileUpdated(ctx context.Context, evt event.Event) error {
    // 转换为活动日志任务并发布到 Worker 队列
    return l.eventBus.Publish(ctx, &ActivityLogTask{
        Action: "USER.PROFILE.UPDATED",
        Data:   evt.GetPayload(),
    })
}
```

---

### 阶段 5：重构 App 层服务（预计 2 小时）

#### 5.1 更新认证服务

**文件**：`app/authentication/service.go`

```go
package authentication

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/events"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
)

type authenticationService struct {
    userRepo UserRepository
    eventBus messaging.EventBus  // ✅ 依赖抽象接口
}

func NewAuthenticationService(userRepo UserRepository, eventBus messaging.EventBus) *authenticationService {
    return &authenticationService{
        userRepo: userRepo,
        eventBus: eventBus,
    }
}

func (s *authenticationService) Login(ctx context.Context, cmd LoginCommand) (*Token, error) {
    // 验证密码...
    if !user.VerifyPassword(cmd.Password) {
        // 发布登录失败事件
        s.eventBus.Publish(ctx, &events.LoginFailed{
            UserID: user.ID,
            Email:  cmd.Email,
            IP:     cmd.IP,
            Reason: "invalid_password",
        })
        return nil, ErrInvalidCredentials
    }
    
    // 生成 Token...
    token := GenerateToken(user.ID)
    
    // ✅ 发布登录成功事件（异步，不阻塞）
    s.eventBus.Publish(ctx, &events.UserLoggedIn{
        UserID:    user.ID,
        Email:     user.Email,
        IP:        cmd.IP,
        UserAgent: cmd.UA,
        Device:    parseDevice(cmd.UA),
    })
    
    // ✅ 立即返回，不等待日志写入
    return token, nil
}
```

#### 5.2 删除 app/activitylog/ 目录

```bash
# 删除不再需要的目录
rm -rf app/activitylog/
```

---

### 阶段 6：更新 Worker 处理器（预计 2 小时）

#### 6.1 审计日志处理器

**文件**：`worker/handlers/audit_log_handler.go`

```go
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

#### 6.2 注册 Handler 到 Worker

**文件**：`cmd/worker/main.go`

```go
func main() {
    // 初始化组件...
    eventBus := messaging.NewEventBus(redisAddr, queueConfig)
    
    // 注册 Handler
    srv := asynq.NewServer(redisOpt, asynq.Config{
        Queues: map[string]int{
            "critical": 6,
            "default":  3,
            "low":      1,
        },
    })
    
    mux := asynq.NewServeMux()
    
    // 注册审计日志处理器
    auditLogHandler := handlers.NewAuditLogHandler(auditLogRepo)
    mux.HandleFunc("AUTH.LOGIN.SUCCESS", auditLogHandler.ProcessTask)
    mux.HandleFunc("AUTH.LOGIN.FAILED", auditLogHandler.ProcessTask)
    mux.HandleFunc("SECURITY.ACCOUNT.LOCKED", auditLogHandler.ProcessTask)
    
    // 注册活动日志处理器
    activityLogHandler := handlers.NewActivityLogHandler(activityLogRepo)
    mux.HandleFunc("USER.PROFILE.UPDATED", activityLogHandler.ProcessTask)
    
    // 启动 Worker
    if err := srv.Run(mux); err != nil {
        log.Fatalf("Could not run server: %v", err)
    }
}
```

---

### 阶段 7：初始化 Listener（预计 1 小时）

**文件**：`cmd/api/main.go`

```go
func main() {
    // 初始化组件...
    eventBus := messaging.NewEventBus(redisAddr, queueConfig)
    
    // 初始化 Listener（订阅事件）
    auditLogListener := listener.NewAuditLogListener(eventBus)
    activityLogListener := listener.NewActivityLogListener(eventBus)
    
    // 不需要显式调用，Listener 已在构造函数中订阅事件
    
    // 初始化 App 服务
    authService := authentication.NewAuthService(userRepo, eventBus)
    
    // 初始化 HTTP Handler
    authHandler := transport.NewAuthHandler(authService)
    
    // 启动 API...
}
```

---

### 阶段 8：测试验证（预计 2 小时）

#### 8.1 单元测试

```bash
# 测试 EventBus
go test ./infra/messaging/... -v

# 测试 Listener
go test ./listener/... -v

# 测试 Worker Handler
go test ./worker/handlers/... -v
```

#### 8.2 集成测试

```bash
# 运行核心流程测试
bash scripts/dev/core-flow-test.sh
```

**验证点**：
- ✅ 用户登录后，立即返回 Token（不等待日志写入）
- ✅ Worker 异步写入 audit_logs 表
- ✅ 查询审计日志接口正常
- ✅ 性能提升（响应时间 < 100ms）

#### 8.3 压力测试

```bash
# 使用 ab 或 wrk 进行压测
ab -n 1000 -c 10 http://localhost:8080/api/v1/auth/login

# 观察指标
# - 响应时间从 200ms → 50ms
# - 吞吐量提升 3-4 倍
```

---

### 阶段 9：清理和文档（预计 1 小时）

#### 9.1 删除旧代码

```bash
# 删除旧的 activitylog 服务
rm -rf internal/app/activitylog/

# 删除 pkg/event/asynq_event_bus.go（如果不再需要）
rm pkg/event/asynq_event_bus.go
```

#### 9.2 更新文档

- ✅ 更新 README.md
- ✅ 更新 ARCHITECTURE_REFACTORING_SPEC.md（已完成）
- ✅ 创建本迁移计划文档

#### 9.3 提交代码

```bash
git add .
git commit -m "feat: implement event-driven logging architecture

- Add EventBus interface and Asynq implementation
- Add AuditLogListener and ActivityLogListener
- Move log writing to Worker (async)
- Remove app/activitylog service
- Improve login response time from 200ms to 50ms

BREAKING CHANGE: Logging is now asynchronous via event bus"

git push origin feature/event-driven-logging
```

---

## 三、回滚方案

如果出现问题，随时回滚：

```bash
# 1. 切换回主分支
git checkout main

# 2. 删除特性分支
git branch -D feature/event-driven-logging

# 3. 恢复工作环境
go mod tidy
go build ./...

# 4. 验证功能正常
bash scripts/dev/core-flow-test.sh
```

---

## 四、风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|---------|
| EventBus 实现错误 | 日志丢失 | 中 | 完整的单元测试 + 集成测试 |
| Worker 崩溃 | 日志积压 | 低 | Asynq 自动重试 + 监控告警 |
| Redis 故障 | 队列不可用 | 低 | Redis Cluster + 降级方案（同步写入） |
| 性能不如预期 | 架构价值降低 | 低 | 压测验证 + 调优队列配置 |

---

## 五、验收标准

### 功能验收

- [ ] 用户登录立即返回 Token（< 50ms）
- [ ] Worker 在 1 秒内写入 audit_logs
- [ ] 查询审计日志接口正常
- [ ] 所有现有测试通过

### 性能验收

- [ ] 登录接口 P99 < 100ms
- [ ] 吞吐量 > 1000 QPS
- [ ] Worker 消费延迟 < 1 秒

### 代码质量

- [ ] 无循环依赖
- [ ] 测试覆盖率 > 70%
- [ ] 通过 lint 检查
- [ ] 文档完整清晰

---

## 六、总结

### 迁移收益

- ✅ **性能提升**：登录响应时间从 200ms → 50ms（75% 提升）
- ✅ **架构解耦**：业务代码不再依赖日志实现
- ✅ **可扩展性**：Worker 可独立水平扩展
- ✅ **可靠性**：队列缓冲，避免数据库瞬时压力

### 工作量估算

- **总工时**：约 14 小时
- **建议安排**：2 个工作日完成
- **参与人员**：1-2 名后端工程师

### 后续优化

1. **监控告警**：添加队列长度、消费延迟监控
2. **死信队列**：处理失败的任务归档
3. **重试策略**：指数退避重试机制
4. **多队列隔离**：审计日志 vs 活动日志物理隔离

---

**文档版本**：v1.0  
**创建日期**：2026-04-03  
**最后更新**：2026-04-03  
**状态**：待执行
