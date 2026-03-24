# Module 开发指南

本文档介绍如何在 Go DDD Scaffold 中开发新的功能模块。

## 📋 Module 概述

### 什么是 Module？

Module（模块）是**组合根（Composition Root）**，负责：
1. 创建基础设施组件
2. 创建适配器
3. 组装依赖
4. 注册路由和事件处理器

### Module 的位置

```
backend/internal/
├── module/
│   ├── auth.go      # AuthModule
│   ├── user.go      # UserModule
│   └── tenant.go    # TenantModule（待开发）
```

### Module 的职责

```go
// ✅ Module 做什么
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 组装应用服务
    authSvc := authApp.NewAuthService(tokenServiceAdapter, ...)
    
    // 4. 注册路由
    routes := authHTTP.NewRoutes(handler, jwtSvc)
    
    return &AuthModule{routes: routes}
}

// ❌ Module 不做什么
// - 不包含业务逻辑
// - 不直接处理 HTTP 请求
// - 不直接访问数据库
```

---

## 🚀 开发新模块的完整流程

### 场景：开发一个"通知模块"（Notification）

#### 步骤 1：创建目录结构

```bash
# 创建领域层目录
mkdir -p backend/internal/domain/notification/{aggregate,valueobject,event,service,repository}

# 创建应用层目录
mkdir -p backend/internal/application/notification
mkdir -p backend/internal/application/ports/notification

# 创建基础设施目录
mkdir -p backend/internal/infrastructure/platform/notification

# 创建接口层目录
mkdir -p backend/internal/interfaces/http/notification

# 创建 Module 文件
touch backend/internal/module/notification.go
```

#### 步骤 2：定义领域模型

**2.1 创建聚合根**

```go
// domain/notification/aggregate/notification.go
package aggregate

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/valueobject"
)

// Notification 通知聚合根
type Notification struct {
    *kernel.Entity
    userID       vo.UserID
    messageType  vo.MessageType
    content      vo.Content
    status       vo.NotificationStatus
    sentAt       *time.Time
    readAt       *time.Time
}

// NewNotification 创建新通知
func NewNotification(
    userID vo.UserID,
    messageType vo.MessageType,
    content vo.Content,
) (*Notification, error) {
    n := &Notification{
        Entity:      kernel.NewEntity(),
        userID:      userID,
        messageType: messageType,
        content:     content,
        status:      vo.NotificationStatusPending,
    }
    
    // 验证
    if err := n.validate(); err != nil {
        return nil, err
    }
    
    // 发布领域事件
    n.RecordEvent(&event.NotificationCreated{
        NotificationID: n.ID().Value(),
        UserID:         userID.Value(),
        MessageType:    messageType,
    })
    
    return n, nil
}

// Send 发送通知
func (n *Notification) Send() error {
    if n.status != vo.NotificationStatusPending {
        return kernel.ErrDomainRuleViolated
    }
    
    n.status = vo.NotificationStatusSent
    now := time.Now()
    n.sentAt = &now
    
    n.RecordEvent(&event.NotificationSent{
        NotificationID: n.ID().Value(),
        SentAt:         *n.sentAt,
    })
    
    return nil
}

// MarkAsRead 标记为已读
func (n *Notification) MarkAsRead() error {
    if n.status != vo.NotificationStatusSent {
        return kernel.ErrDomainRuleViolated
    }
    
    n.status = vo.NotificationStatusRead
    now := time.Now()
    n.readAt = &now
    
    return nil
}

func (n *Notification) validate() error {
    // 验证逻辑
    return nil
}
```

**2.2 创建值对象**

```go
// domain/notification/valueobject/type.go
package valueobject

// MessageType 消息类型
type MessageType string

const (
    MessageTypeEmail      MessageType = "email"
    MessageTypeSMS        MessageType = "sms"
    MessageTypePush       MessageType = "push"
    MessageType站内信     MessageType = "internal"
)

// Content 消息内容
type Content struct {
    Subject string
    Body    string
}

// NotificationStatus 通知状态
type NotificationStatus string

const (
    NotificationStatusPending NotificationStatus = "pending"
    NotificationStatusSent    NotificationStatus = "sent"
    NotificationStatusRead    NotificationStatus = "read"
    NotificationStatusFailed  NotificationStatus = "failed"
)

// UserID 用户 ID
type UserID struct {
    value int64
}

func NewUserID(value int64) (UserID, error) {
    if value <= 0 {
        return UserID{}, kernel.FieldError("userID", "invalid user id", value)
    }
    return UserID{value: value}, nil
}

func (u UserID) Value() int64 { return u.value }
```

**2.3 创建领域事件**

```go
// domain/notification/event/notification_created.go
package event

import "time"

// NotificationCreated 通知已创建
type NotificationCreated struct {
    NotificationID int64
    UserID         int64
    MessageType    string
    OccurredAt     time.Time
}

func (e *NotificationCreated) Type() string {
    return "notification.created"
}

func (e *NotificationCreated) AggregateID() int64 {
    return e.NotificationID
}

func (e *NotificationCreated) AggregateType() string {
    return "Notification"
}

func (e *NotificationCreated) Timestamp() time.Time {
    return e.OccurredAt
}
```

**2.4 定义 Repository 接口**

```go
// domain/notification/repository/notification_repository.go
package repository

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/aggregate"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/valueobject"
)

// NotificationRepository 通知仓储接口
type NotificationRepository interface {
    // 基本 CRUD
    FindByID(ctx context.Context, id int64) (*aggregate.Notification, error)
    FindByUserID(ctx context.Context, userID vo.UserID, limit int) ([]*aggregate.Notification, error)
    Save(ctx context.Context, notification *aggregate.Notification) error
    Delete(ctx context.Context, id int64) error
    
    // 查询方法
    CountUnread(ctx context.Context, userID vo.UserID) (int, error)
}
```

#### 步骤 3：实现应用层

**3.1 定义 Ports**

```go
// application/ports/notification/email_service.go
package ports

// EmailService 邮件服务端口
type EmailService interface {
    SendEmail(to, subject, body string) error
    SendTemplateEmail(to, templateName string, data map[string]interface{}) error
}
```

**3.2 创建 DTO**

```go
// application/notification/dto.go
package notification

import vo "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/valueobject"

// CreateNotificationCommand 创建通知命令
type CreateNotificationCommand struct {
    UserID      int64
    MessageType vo.MessageType
    Subject     string
    Body        string
}

// NotificationResponse 通知响应
type NotificationResponse struct {
    ID          int64  `json:"id"`
    UserID      int64  `json:"user_id"`
    MessageType string `json:"message_type"`
    Subject     string `json:"subject"`
    Status      string `json:"status"`
    CreatedAt   string `json:"created_at"`
}
```

**3.3 实现应用服务**

```go
// application/notification/service.go
package notification

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/repository"
    ports_notification "github.com/shenfay/go-ddd-scaffold/internal/application/ports/notification"
    "go.uber.org/zap"
)

// NotificationService 通知应用服务
type NotificationService struct {
    logger        *zap.Logger
    repo          repository.NotificationRepository
    emailService  ports_notification.EmailService
    eventPublisher EventPublisher
}

func NewNotificationService(
    logger *zap.Logger,
    repo repository.NotificationRepository,
    emailService ports_notification.EmailService,
    eventPublisher EventPublisher,
) *NotificationService {
    return &NotificationService{
        logger:        logger,
        repo:          repo,
        emailService:  emailService,
        eventPublisher: eventPublisher,
    }
}

// CreateNotification 创建通知
func (s *NotificationService) CreateNotification(
    ctx context.Context, 
    cmd *CreateNotificationCommand,
) (*NotificationResponse, error) {
    
    // 1. 创建领域对象
    userID, err := vo.NewUserID(cmd.UserID)
    if err != nil {
        return nil, err
    }
    
    content := vo.Content{
        Subject: cmd.Subject,
        Body:    cmd.Body,
    }
    
    notification, err := aggregate.NewNotification(userID, cmd.MessageType, content)
    if err != nil {
        return nil, err
    }
    
    // 2. 保存通知
    err = s.repo.Save(ctx, notification)
    if err != nil {
        return nil, err
    }
    
    // 3. 发送邮件（如果适用）
    if cmd.MessageType == vo.MessageTypeEmail {
        err = s.emailService.SendEmail(
            cmd.UserID,  // 需要转换为用户邮箱
            cmd.Subject,
            cmd.Body,
        )
        if err != nil {
            s.logger.Error("send email failed", zap.Error(err))
            // 不返回错误，继续执行
        }
    }
    
    // 4. 发布领域事件
    events := notification.ReleaseEvents()
    for _, event := range events {
        s.eventPublisher.Publish(event)
    }
    
    // 5. 返回响应
    return &NotificationResponse{
        ID:          notification.ID().Value(),
        UserID:      cmd.UserID,
        MessageType: string(cmd.MessageType),
        Subject:     cmd.Subject,
        Status:      string(notification.Status()),
        CreatedAt:   notification.CreatedAt().Format(time.RFC3339),
    }, nil
}
```

#### 步骤 4：实现基础设施层

**4.1 实现 Repository**

```go
// infrastructure/persistence/repository/notification_repository.go
package repository

import (
    "context"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/notification/repository"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

type notificationRepository struct {
    db       *gorm.DB
    daoQuery *dao_query.Query
}

func NewNotificationRepository(db *gorm.DB, daoQuery *dao_query.Query) repository.NotificationRepository {
    return &notificationRepository{
        db:       db,
        daoQuery: daoQuery,
    }
}

func (r *notificationRepository) FindByID(ctx context.Context, id int64) (*aggregate.Notification, error) {
    dao, err := r.daoQuery.Notification.WithContext(ctx).Where(r.daoQuery.Notification.ID.Eq(id)).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, kernel.ErrAggregateNotFound
        }
        return nil, err
    }
    
    return r.toDomain(dao)
}

func (r *notificationRepository) Save(ctx context.Context, notification *aggregate.Notification) error {
    // 开始事务
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 保存或更新
    err := r.saveNotification(tx, notification)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 保存领域事件
    events := notification.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}

func (r *notificationRepository) toDomain(dao *dao.Notification) (*aggregate.Notification, error) {
    // DAO → Domain 转换
}

func (r *notificationRepository) saveNotification(tx *gorm.DB, notification *aggregate.Notification) error {
    // Domain → DAO 转换并保存
}
```

**4.2 实现 EmailService**

```go
// infrastructure/platform/notification/email_service_impl.go
package notification

import (
    "github.com/shenfay/go-ddd-scaffold/internal/application/ports/notification"
    "gopkg.in/gomail.v2"
)

type emailServiceImpl struct {
    smtpHost string
    smtpPort int
    username string
    password string
}

func NewEmailServiceImpl(smtpHost string, smtpPort int, username, password string) ports.EmailService {
    return &emailServiceImpl{
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        username: username,
        password: password,
    }
}

func (s *emailServiceImpl) SendEmail(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.username)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)
    
    d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.username, s.password)
    
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("send email failed: %w", err)
    }
    
    return nil
}

func (s *emailServiceImpl) SendTemplateEmail(to, templateName string, data map[string]interface{}) error {
    // TODO: 实现模板邮件
    return nil
}

// 编译期检查
var _ ports.EmailService = (*emailServiceImpl)(nil)
```

#### 步骤 5：实现接口层

**5.1 创建 Handler**

```go
// interfaces/http/notification/handler.go
package notification

import (
    "github.com/gin-gonic/gin"
    app_notification "github.com/shenfay/go-ddd-scaffold/internal/application/notification"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
)

type Handler struct {
    service     *app_notification.NotificationService
    respHandler *http_shared.ResponseHandler
}

func NewHandler(
    service *app_notification.NotificationService,
    respHandler *http_shared.ResponseHandler,
) *Handler {
    return &Handler{
        service:     service,
        respHandler: respHandler,
    }
}

// CreateNotification 创建通知
func (h *Handler) CreateNotification(c *gin.Context) {
    var req CreateNotificationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, http.StatusBadRequest, err)
        return
    }
    
    cmd := &app_notification.CreateNotificationCommand{
        UserID:      GetUserIDFromContext(c),
        MessageType: vo.MessageType(req.MessageType),
        Subject:     req.Subject,
        Body:        req.Body,
    }
    
    result, err := h.service.CreateNotification(c.Request.Context(), cmd)
    if err != nil {
        h.respHandler.Error(c, http.StatusInternalServerError, err)
        return
    }
    
    h.respHandler.Success(c, result)
}

// ListNotifications 获取通知列表
func (h *Handler) ListNotifications(c *gin.Context) {
    // 实现类似...
}
```

**5.2 创建 Routes**

```go
// interfaces/http/notification/routes.go
package notification

import (
    "github.com/gin-gonic/gin"
    "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware/auth"
)

type Routes struct {
    handler *Handler
    auth    *auth.JWTMiddleware
}

func NewRoutes(handler *Handler, auth *auth.JWTMiddleware) *Routes {
    return &Routes{
        handler: handler,
        auth:    auth,
    }
}

// Register 注册路由
func (r *Routes) Register(router *gin.RouterGroup) {
    notification := router.Group("/notifications")
    notification.Use(r.auth.Middleware())
    {
        notification.POST("", r.handler.CreateNotification)
        notification.GET("", r.handler.ListNotifications)
        notification.GET("/:id", r.handler.GetNotification)
        notification.DELETE("/:id", r.handler.DeleteNotification)
        notification.POST("/:id/read", r.handler.MarkAsRead)
    }
}
```

#### 步骤 6：创建 Module

```go
// module/notification.go
package module

import (
    "github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
    app_notification "github.com/shenfay/go-ddd-scaffold/internal/application/notification"
    infra_notification "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/notification"
    http_notification "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/notification"
    http_shared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

// NotificationModule 通知模块
type NotificationModule struct {
    infra  *bootstrap.Infra
    routes *http_notification.Routes
}

// NewNotificationModule 创建通知模块
func NewNotificationModule(infra *bootstrap.Infra) *NotificationModule {
    // 1. 创建基础设施
    daoQuery := dao_query.Use(infra.DB)
    
    emailService := infra_notification.NewEmailServiceImpl(
        infra.Config.SMTP.Host,
        infra.Config.SMTP.Port,
        infra.Config.SMTP.Username,
        infra.Config.SMTP.Password,
    )
    
    // 2. 创建 Repository
    notificationRepo := repository.NewNotificationRepository(infra.DB, daoQuery)
    
    // 3. 创建应用服务
    respHandler := http_shared.NewHandler(infra.ErrorMapper)
    handler := http_notification.NewHandler(notificationSvc, respHandler)
    routes := http_notification.NewRoutes(handler, infra.JWTMiddleware)
    
    return &NotificationModule{
        infra:  infra,
        routes: routes,
    }
}

// RegisterRoutes 注册路由
func (m *NotificationModule) RegisterRoutes(router *gin.Engine) {
    m.routes.Register(&router.Group("/api/v1"))
}

// SubscribeEvents 订阅领域事件
func (m *NotificationModule) SubscribeEvents() {
    // 订阅感兴趣的事件
    infra.EventPublisher.Subscribe("user.created", m.handleUserCreated)
}

func (m *NotificationModule) handleUserCreated(event event.UserCreated) error {
    // 用户注册后发送欢迎通知
    cmd := &app_notification.CreateNotificationCommand{
        UserID:      event.UserID,
        MessageType: "email",
        Subject:     "欢迎加入我们！",
        Body:        "感谢您的注册...",
    }
    
    _, err := m.notificationSvc.CreateNotification(context.Background(), cmd)
    return err
}
```

#### 步骤 7：注册 Module

```go
// bootstrap/module.go
func LoadModules(infra *Infra) []Module {
    modules := []Module{
        NewUserModule(infra),
        NewAuthModule(infra),
        NewNotificationModule(infra),  // ← 添加新模块
    }
    
    return modules
}

// main.go
func main() {
    // ...
    
    // 加载所有模块
    modules := bootstrap.LoadModules(infra)
    
    // 注册路由
    for _, module := range modules {
        module.RegisterRoutes(router)
        module.SubscribeEvents()
    }
    
    // ...
}
```

---

## ✅ Module 开发检查清单

### 领域层

- [ ] 创建聚合根（Aggregate）
- [ ] 创建值对象（Value Objects）
- [ ] 创建领域事件（Domain Events）
- [ ] 定义 Repository 接口

### 应用层

- [ ] 定义 Ports（外部依赖接口）
- [ ] 创建 DTO（命令和响应）
- [ ] 实现应用服务

### 基础设施层

- [ ] 实现 Repository
- [ ] 实现 Ports 适配器
- [ ] 实现 Unit of Work

### 接口层

- [ ] 创建 Handler
- [ ] 创建 Routes
- [ ] 创建请求/响应 DTO

### Module

- [ ] 创建 Module 结构体
- [ ] 实现构造函数
- [ ] 实现 RegisterRoutes
- [ ] 实现 SubscribeEvents
- [ ] 在 bootstrap 中注册

---

## 📚 参考资源

- [架构总览](../design/architecture-overview.md)
- [Ports 模式详解](../design/ports-pattern-design.md)
- [领域模型设计](../design/domain-model.md)
- [开发规范](../specifications/development-spec.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
