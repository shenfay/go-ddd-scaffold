# Go DDD CQRS 设计指南

## 文档概述

本文档详细阐述了在go-ddd-scaffold项目中如何实现标准的DDD+CQRS架构模式，包括领域设计、边界划分、事件驱动等核心概念的具体实现方案。

## CQRS架构模式详解

### 核心设计理念

CQRS（Command Query Responsibility Segregation）将系统的读写操作完全分离，这种模式特别适合复杂的业务场景：

```
┌─────────────────┐    ┌─────────────────┐
│   Command Side  │    │   Query Side    │
│  (Write Model)  │    │  (Read Model)   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│  Command Bus    │    │  Query Service  │
│  & Handlers     │    │  & Projections  │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│  Domain Model   │◄──►│  Read Stores    │
│  (Aggregates)   │    │  (Optimized)    │
└─────────────────┘    └─────────────────┘
```

### 读写模型分离策略

#### 写模型（Command Model）
- **职责**：处理业务逻辑、维护数据一致性
- **特点**：规范化设计、强一致性、复杂业务规则
- **存储**：事务性数据库（PostgreSQL）

#### 读模型（Query Model）
- **职责**：优化查询性能、支持复杂展示需求
- **特点**：非规范化设计、最终一致性、高性能读取
- **存储**：可选用不同存储（Redis、Elasticsearch、专用读库）

## 领域设计核心要素

### 1. 聚合根设计原则

#### 聚合边界的确定
```go
// 用户聚合根 - 正确的边界设计
type UserAggregate struct {
    // 聚合根必须包含的元素
    baseAggregate BaseAggregate  // 基础聚合属性
    
    // 核心业务属性
    id       UserID
    username string
    email    Email
    password HashedPassword
    status   UserStatus
    
    // 相关值对象和实体
    profile  UserProfile      // 用户档案
    settings UserSettings     // 用户设置
    roles    []UserRole       // 用户角色列表
}

// 错误示例：聚合过大
type BadUserAggregate struct {
    // ... 用户基本信息
    
    orders   []Order         // 订单不应该在用户聚合内
    comments []Comment       // 评论不应该在用户聚合内
    logs     []OperationLog  // 操作日志不应该在用户聚合内
}
```

#### 聚合根的业务方法
```go
// 正确的聚合根方法设计 - 包含领域事件发布
func (u *User) ChangeEmail(newEmail string) error {
    // 1. 验证业务规则
    oldEmail := u.email.Value()
    if oldEmail == newEmail {
        return ddd.NewBusinessError("EMAIL_UNCHANGED", "email unchanged")
    }
    
    email, err := NewEmail(newEmail)
    if err != nil {
        return err
    }
    
    // 2. 执行业务逻辑
    u.email = email
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 3. 发布领域事件
    event := NewUserEmailChangedEvent(u.ID().(UserID), oldEmail, newEmail)
    u.ApplyEvent(event)
    
    return nil
}

// Activate 激活用户 - 状态变更示例
func (u *User) Activate() error {
    // 1. 验证业务规则
    if u.status != UserStatusPending {
        return ddd.NewBusinessError("USER_NOT_PENDING", "user is not in pending status")
    }

    // 2. 执行业务逻辑
    u.status = UserStatusActive
    u.updatedAt = time.Now()
    u.IncrementVersion()

    // 3. 发布领域事件
    event := NewUserActivatedEvent(u.ID().(UserID))
    u.ApplyEvent(event)

    return nil
}
```

### 2. 领域事件设计

#### 事件分类和命名规范
```go
// 领域事件接口定义
type DomainEvent interface {
    EventName() string           // 事件名称
    OccurredOn() time.Time       // 发生时间
    AggregateID() interface{}    // 聚合标识
    Version() int               // 事件版本
}

// 用户领域事件体系
type (
    // 用户生命周期事件
    UserCreatedEvent struct {
        UserID    UserID    `json:"user_id"`
        Username  string    `json:"username"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
    }
    
    UserActivatedEvent struct {
        UserID     UserID    `json:"user_id"`
        ActivatedAt time.Time `json:"activated_at"`
    }
    
    UserDeactivatedEvent struct {
        UserID       UserID    `json:"user_id"`
        Reason       string    `json:"reason"`
        DeactivatedAt time.Time `json:"deactivated_at"`
    }
    
    // 用户属性变更事件
    UserEmailChangedEvent struct {
        UserID    UserID `json:"user_id"`
        OldEmail  string `json:"old_email"`
        NewEmail  string `json:"new_email"`
        ChangedAt time.Time `json:"changed_at"`
    }
    
    UserPasswordChangedEvent struct {
        UserID    UserID    `json:"user_id"`
        ChangedAt time.Time `json:"changed_at"`
    }
    
    // 用户关联关系事件
    UserRoleAssignedEvent struct {
        UserID   UserID `json:"user_id"`
        RoleID   int64  `json:"role_id"`
        AssignedAt time.Time `json:"assigned_at"`
    }
)
```

#### 事件发布的时机和原则
```go
// 聚合根中的事件发布 - 使用 ApplyEvent 记录事件
func (u *User) RecordLogin(ipAddress, userAgent string) {
    now := time.Now()
    u.lastLoginAt = &now
    u.loginCount++
    u.failedAttempts = 0
    u.updatedAt = now
    u.IncrementVersion()

    // 发布领域事件 - 在业务逻辑成功执行后立即发布
    event := NewUserLoggedInEvent(u.ID().(UserID), ipAddress, userAgent)
    u.ApplyEvent(event)
}

// 重要原则：
// 1. 事件应该在业务逻辑成功执行后立即发布
// 2. 使用 ApplyEvent 将事件添加到未提交事件列表
// 3. 仓储层负责将未提交事件持久化到事件存储
// 4. 事件总线负责将事件分发给订阅者
```

### 3. 值对象设计模式

#### 不可变性保证
```go
// 正确的值对象设计
type Email struct {
    value string
}

// 工厂方法确保有效性
func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: strings.ToLower(email)}, nil
}

// 只提供访问方法，不提供修改方法
func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return strings.EqualFold(e.value, other.value)
}

func (e Email) IsValid() bool {
    return isValidEmail(e.value)
}

// 错误示例：可变的值对象
type MutableEmail struct {
    value string
}

func (me *MutableEmail) SetValue(newValue string) {
    me.value = newValue // 违反了值对象的不可变性原则
}
```

#### 复杂值对象示例
```go
// 地址值对象
type Address struct {
    street     string
    city       string
    state      string
    postalCode string
    country    string
}

func NewAddress(street, city, state, postalCode, country string) (Address, error) {
    addr := Address{
        street:     street,
        city:       city,
        state:      state,
        postalCode: postalCode,
        country:    country,
    }
    
    if err := addr.validate(); err != nil {
        return Address{}, err
    }
    
    return addr, nil
}

func (a Address) validate() error {
    if strings.TrimSpace(a.street) == "" {
        return errors.New("street is required")
    }
    if strings.TrimSpace(a.city) == "" {
        return errors.New("city is required")
    }
    if strings.TrimSpace(a.postalCode) == "" {
        return errors.New("postal code is required")
    }
    return nil
}

func (a Address) Formatted() string {
    return fmt.Sprintf("%s, %s, %s %s, %s", 
        a.street, a.city, a.state, a.postalCode, a.country)
}

// 值对象组合
type ContactInfo struct {
    email   Email
    phone   PhoneNumber
    address Address
}

func NewContactInfo(email Email, phone PhoneNumber, address Address) ContactInfo {
    return ContactInfo{
        email:   email,
        phone:   phone,
        address: address,
    }
}
```

## 边界划分策略

### 1. 限界上下文识别

#### 业务能力分析法
```go
// 用户管理限界上下文
type UserManagementContext struct {
    // 用户相关的所有聚合根
    userAggregate *UserAggregate
    tenantAggregate *TenantAggregate
    
    // 上下文特有的领域服务
    userService *UserService
    tenantService *TenantService
    
    // 上下文特有的仓储
    userRepo UserRepository
    tenantRepo TenantRepository
}

// 订单管理限界上下文
type OrderManagementContext struct {
    orderAggregate *OrderAggregate
    productAggregate *ProductAggregate
    
    orderService *OrderService
    productService *ProductService
    
    orderRepo OrderRepository
    productRepo ProductRepository
}
```

#### 上下文映射关系
```go
// 上下文间通信方式选择

// 1. 共享内核（Shared Kernel）
// 适用于紧密耦合的上下文
type SharedKernel struct {
    commonEntities []interface{}  // 共享的实体定义
    commonVOs      []interface{}  // 共享的值对象
    commonEvents   []DomainEvent  // 共享的领域事件
}

// 2. 客户-供应商（Customer-Supplier）
// 上游供应方，下游消费方
type UpstreamSupplier interface {
    // 提供标准化的API
    GetUserProfile(userID UserID) (*UserProfile, error)
    GetUserOrders(userID UserID) ([]OrderSummary, error)
}

type DownstreamConsumer struct {
    supplierClient UpstreamSupplier
}

// 3. 防腐层（Anti-Corruption Layer）
// 保护核心领域不受外部影响
type ACLayer struct {
    externalService ExternalUserService
    userTranslator  UserTranslator
}

func (acl *ACLayer) GetExternalUserInfo(externalID string) (*User, error) {
    externalUser, err := acl.externalService.GetUser(externalID)
    if err != nil {
        return nil, err
    }
    
    // 转换为内部领域模型
    return acl.userTranslator.ToDomain(externalUser), nil
}
```

### 2. 聚合间协作模式

#### 最终一致性协作
```go
// 通过领域事件实现聚合间协作
type UserRegistrationSaga struct {
    eventBus EventBus
    userRepo UserRepository
    tenantRepo TenantRepository
}

func (s *UserRegistrationSaga) HandleUserCreated(event UserCreatedEvent) error {
    // 1. 创建默认用户档案
    profile := NewUserProfile(event.UserID, ProfileDefaults{})
    
    // 2. 分配默认租户
    defaultTenant, err := s.tenantRepo.GetDefaultTenant()
    if err != nil {
        return err
    }
    
    userTenant := UserTenant{
        UserID:   event.UserID,
        TenantID: defaultTenant.ID(),
        RoleID:   DefaultUserRoleID,
    }
    
    // 3. 发布后续事件
    s.eventBus.Publish(UserProfileCreatedEvent{
        UserID: event.UserID,
        Profile: profile,
    })
    
    s.eventBus.Publish(UserTenantAssignedEvent{
        UserID:   event.UserID,
        TenantID: defaultTenant.ID(),
    })
    
    return nil
}
```

#### 直接引用vs ID引用
```go
// 正确的方式：通过ID引用其他聚合
type Order struct {
    id         OrderID
    userID     UserID        // 引用用户聚合的ID
    productIDs []ProductID   // 引用产品聚合的ID
    amount     Money
}

// 错误的方式：直接引用其他聚合根
type BadOrder struct {
    id     OrderID
    user   *User         // 直接引用用户聚合根 - 违反聚合边界
    products []*Product  // 直接引用产品聚合根 - 违反聚合边界
}
```

## CQRS具体实现方案

### 1. 命令侧实现

#### 命令对象设计
```go
// 命令接口定义
type Command interface {
    CommandName() string
    Validate() error
}

// 具体命令实现
type CreateUserCommand struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    TenantID *int64 `json:"tenant_id,omitempty"`
}

func (cmd CreateUserCommand) CommandName() string {
    return "CreateUser"
}

func (cmd CreateUserCommand) Validate() error {
    // 基础验证
    if len(cmd.Username) < 3 {
        return errors.New("username too short")
    }
    
    if !isValidEmail(cmd.Email) {
        return errors.New("invalid email format")
    }
    
    if len(cmd.Password) < 8 {
        return errors.New("password too short")
    }
    
    return nil
}
```

#### 命令处理器实现
```go
// 命令处理器接口
type CommandHandler interface {
    Handle(command Command) (interface{}, error)
}

// 用户命令处理器
type UserCommandHandler struct {
    userRepo      UserRepository
    tenantRepo    TenantRepository
    eventBus      EventBus
    passwordSvc   PasswordService
    validator     *validator.Validate
}

func (h *UserCommandHandler) HandleCreateUser(cmd CreateUserCommand) (UserID, error) {
    // 1. 命令验证
    if err := h.validator.Struct(cmd); err != nil {
        return 0, fmt.Errorf("validation failed: %w", err)
    }
    
    // 2. 业务规则验证
    if exists, _ := h.userRepo.ExistsByEmail(cmd.Email); exists {
        return 0, errors.New("email already registered")
    }
    
    // 3. 创建聚合根
    user, err := NewUser(cmd.Username, cmd.Email, cmd.Password)
    if err != nil {
        return 0, fmt.Errorf("failed to create user: %w", err)
    }
    
    // 4. 处理租户关联
    if cmd.TenantID != nil {
        tenant, err := h.tenantRepo.GetByID(TenantID(*cmd.TenantID))
        if err != nil {
            return 0, fmt.Errorf("tenant not found: %w", err)
        }
        
        if err := user.AssignToTenant(tenant.ID()); err != nil {
            return 0, fmt.Errorf("failed to assign tenant: %w", err)
        }
    }
    
    // 5. 持久化聚合
    if err := h.userRepo.Save(user); err != nil {
        return 0, fmt.Errorf("failed to save user: %w", err)
    }
    
    // 6. 发布领域事件
    events := user.GetUncommittedEvents()
    for _, event := range events {
        if err := h.eventBus.Publish(event); err != nil {
            // 记录错误但不中断主流程
            log.Printf("Failed to publish event: %v", err)
        }
    }
    
    // 7. 清除已提交事件
    user.ClearUncommittedEvents()
    
    return user.ID(), nil
}
```

### 2. 查询侧实现

#### 查询对象设计
```go
// 查询接口定义
type Query interface {
    QueryName() string
}

// 具体查询实现
type GetUserProfileQuery struct {
    UserID UserID `json:"user_id"`
}

func (q GetUserProfileQuery) QueryName() string {
    return "GetUserProfile"
}

type ListUsersQuery struct {
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    Status   *int   `json:"status,omitempty"`
    Keyword  string `json:"keyword,omitempty"`
}

func (q ListUsersQuery) QueryName() string {
    return "ListUsers"
}
```

#### 查询服务实现
```go
// 查询服务接口
type QueryService interface {
    Execute(query Query) (interface{}, error)
}

// 用户查询服务
type UserQueryService struct {
    db *gorm.DB
}

// 优化的读模型查询
func (qs *UserQueryService) GetUserProfile(query GetUserProfileQuery) (*UserProfileDTO, error) {
    var profile UserProfileDTO
    
    err := qs.db.Table("user_read_model").
        Select(`
            id,
            username,
            email,
            status,
            created_at,
            updated_at,
            tenant_count,
            last_login_at
        `).
        Where("id = ?", query.UserID).
        First(&profile).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, fmt.Errorf("query failed: %w", err)
    }
    
    return &profile, nil
}

// 复杂列表查询
func (qs *UserQueryService) ListUsers(query ListUsersQuery) (*PagedResult[UserListItemDTO], error) {
    if query.Page <= 0 {
        query.Page = 1
    }
    if query.PageSize <= 0 || query.PageSize > 100 {
        query.PageSize = 20
    }
    
    offset := (query.Page - 1) * query.PageSize
    
    db := qs.db.Table("user_read_model")
    
    // 动态查询条件
    if query.Status != nil {
        db = db.Where("status = ?", *query.Status)
    }
    
    if query.Keyword != "" {
        keyword := "%" + strings.ToLower(query.Keyword) + "%"
        db = db.Where(
            "LOWER(username) LIKE ? OR LOWER(email) LIKE ?",
            keyword, keyword,
        )
    }
    
    // 获取总数
    var totalCount int64
    if err := db.Count(&totalCount).Error; err != nil {
        return nil, fmt.Errorf("count query failed: %w", err)
    }
    
    // 获取数据
    var users []UserListItemDTO
    if err := db.
        Offset(offset).
        Limit(query.PageSize).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, fmt.Errorf("data query failed: %w", err)
    }
    
    return &PagedResult[UserListItemDTO]{
        Items:      users,
        TotalCount: totalCount,
        Page:       query.Page,
        PageSize:   query.PageSize,
    }, nil
}
```

### 3. 读模型同步策略

#### 事件驱动的读模型更新
```go
// 读模型投影器
type UserProjector struct {
    db DB
}

// Project 投影领域事件到读模型
func (p *UserProjector) Project(ctx context.Context, event ddd.DomainEvent) error {
    switch e := event.(type) {
    case *user.UserRegisteredEvent:
        return p.handleUserRegistered(ctx, e)
    case *user.UserActivatedEvent:
        return p.handleUserActivated(ctx, e)
    case *user.UserDeactivatedEvent:
        return p.handleUserDeactivated(ctx, e)
    case *user.UserLoggedInEvent:
        return p.handleUserLoggedIn(ctx, e)
    case *user.UserEmailChangedEvent:
        return p.handleUserEmailChanged(ctx, e)
    // ... 其他事件处理
    default:
        return nil
    }
}

func (p *UserProjector) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
    _, err := p.db.Exec(ctx,
        `INSERT INTO user_read_model (user_id, username, email, status, created_at, updated_at, login_count) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
        event.UserID.Int64(),
        event.Username,
        event.Email,
        int(user.UserStatusPending),
        event.RegisteredAt,
        event.RegisteredAt,
        0,
    )
    return err
}

func (p *UserProjector) handleUserEmailChanged(ctx context.Context, event *user.UserEmailChangedEvent) error {
    _, err := p.db.Exec(ctx,
        "UPDATE user_read_model SET email = ?, updated_at = ? WHERE user_id = ?",
        event.NewEmail,
        event.ChangedAt,
        event.UserID.Int64(),
    )
    return err
}

func (p *UserProjector) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
    _, err := p.db.Exec(ctx,
        `UPDATE user_read_model SET last_login_at = ?, login_count = login_count + 1, updated_at = ? WHERE user_id = ?`,
        event.LoginAt,
        event.LoginAt,
        event.UserID.Int64(),
    )
    return err
}
```

## 事件驱动架构实现

### 1. 事件存储设计

#### 事件表结构
```sql
CREATE TABLE domain_events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(50) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_version INT NOT NULL,
    event_data JSONB NOT NULL,
    occurred_on TIMESTAMP NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_aggregate_lookup (aggregate_id, aggregate_type),
    INDEX idx_event_type (event_type),
    INDEX idx_occurred_on (occurred_on)
);
```

#### 事件存储实现
```go
type EventStore interface {
    AppendEvents(aggregateID interface{}, events []DomainEvent) error
    GetEventsForAggregate(aggregateID interface{}, afterVersion int) ([]DomainEvent, error)
    GetAllEvents(afterID int64, limit int) ([]StoredEvent, error)
}

type PostgresEventStore struct {
    db *gorm.DB
}

func (es *PostgresEventStore) AppendEvents(aggregateID interface{}, events []DomainEvent) error {
    return es.db.Transaction(func(tx *gorm.DB) error {
        for _, event := range events {
            eventData, err := json.Marshal(event)
            if err != nil {
                return fmt.Errorf("failed to marshal event: %w", err)
            }
            
            storedEvent := StoredEvent{
                AggregateID:   fmt.Sprintf("%v", aggregateID),
                AggregateType: getAggregateType(event),
                EventType:     event.EventName(),
                EventVersion:  event.Version(),
                EventData:     string(eventData),
                OccurredOn:    event.OccurredOn(),
            }
            
            if err := tx.Create(&storedEvent).Error; err != nil {
                return fmt.Errorf("failed to store event: %w", err)
            }
        }
        return nil
    })
}
```

### 2. 事件发布订阅机制

#### 事件总线实现
```go
type EventBus interface {
    Publish(event DomainEvent) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
}

type SimpleEventBus struct {
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

func (eb *SimpleEventBus) Publish(event DomainEvent) error {
    eb.mutex.RLock()
    defer eb.mutex.RUnlock()
    
    eventType := event.EventName()
    handlers, exists := eb.handlers[eventType]
    if !exists {
        return nil // 没有订阅者是正常的
    }
    
    // 异步处理事件
    wg := sync.WaitGroup{}
    for _, handler := range handlers {
        wg.Add(1)
        go func(h EventHandler) {
            defer wg.Done()
            if err := h(event); err != nil {
                log.Printf("Error handling event %s: %v", eventType, err)
            }
        }(handler)
    }
    
    wg.Wait()
    return nil
}

func (eb *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()
    
    if eb.handlers[eventType] == nil {
        eb.handlers[eventType] = make([]EventHandler, 0)
    }
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
    return nil
}
```

### 3. 事件处理器分类

#### 实时处理器
```go
// 发送欢迎邮件
type WelcomeEmailHandler struct {
    emailService EmailService
}

func (h *WelcomeEmailHandler) Handle(event DomainEvent) error {
    userCreated, ok := event.(UserCreatedEvent)
    if !ok {
        return nil // 不是关心的事件类型
    }
    
    return h.emailService.SendWelcomeEmail(userCreated.Email, userCreated.Username)
}

// 更新搜索索引
type SearchIndexHandler struct {
    searchService SearchService
}

func (h *SearchIndexHandler) Handle(event DomainEvent) error {
    switch evt := event.(type) {
    case UserCreatedEvent:
        return h.searchService.IndexUser(SearchUser{
            ID:       int64(evt.UserID),
            Username: evt.Username,
            Email:    evt.Email,
        })
    case UserEmailChangedEvent:
        return h.searchService.UpdateUserEmail(int64(evt.UserID), evt.NewEmail)
    }
    return nil
}
```

#### 批处理处理器
```go
// 数据仓库同步处理器
type DataWarehouseHandler struct {
    warehouseDB *gorm.DB
    batchSize   int
    buffer      []DomainEvent
    mutex       sync.Mutex
}

func (h *DataWarehouseHandler) Handle(event DomainEvent) error {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    
    h.buffer = append(h.buffer, event)
    
    if len(h.buffer) >= h.batchSize {
        return h.flushBuffer()
    }
    
    return nil
}

func (h *DataWarehouseHandler) flushBuffer() error {
    if len(h.buffer) == 0 {
        return nil
    }
    
    // 批量处理事件
    events := make([]interface{}, len(h.buffer))
    for i, event := range h.buffer {
        events[i] = convertToWarehouseFormat(event)
    }
    
    err := h.warehouseDB.Table("user_analytics").CreateInBatches(events, 100).Error
    
    // 清空缓冲区
    h.buffer = h.buffer[:0]
    
    return err
}
```

## 开发者友好性设计

### 1. 代码生成工具

#### CLI命令设计
```bash
# 生成新的聚合根
go run ./cmd/cli ddd generate aggregate user username:string email:string

# 生成领域事件
go run ./cmd/cli ddd generate event UserCreated username:string email:string

# 生成CQRS服务
go run ./cmd/cli ddd generate service user --cqrs

# 验证DDD架构合规性
go run ./cmd/cli ddd validate
```

#### 生成器实现示例
```go
type AggregateGenerator struct {
    templateEngine *template.Template
}

func (g *AggregateGenerator) Generate(name string, fields []Field) error {
    data := struct {
        Name   string
        Fields []Field
        Time   time.Time
    }{
        Name:   name,
        Fields: fields,
        Time:   time.Now(),
    }
    
    var buf bytes.Buffer
    if err := g.templateEngine.ExecuteTemplate(&buf, "aggregate.tmpl", data); err != nil {
        return err
    }
    
    filename := fmt.Sprintf("internal/domain/%s/%s.go", 
        strings.ToLower(name), strings.ToLower(name))
    
    return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
```

### 2. 架构验证工具

#### 静态分析规则
```go
type ArchitectureValidator struct {
    rules []ValidationRule
}

type ValidationRule struct {
    Name        string
    Description string
    Check       func(pkg *Package) error
}

var DefaultRules = []ValidationRule{
    {
        Name:        "AggregateRootMethods",
        Description: "聚合根方法应该返回错误而不是panic",
        Check:       validateAggregateRootMethods,
    },
    {
        Name:        "DomainEventNaming",
        Description: "领域事件应该以过去时态命名",
        Check:       validateDomainEventNaming,
    },
    {
        Name:        "LayerDependencies",
        Description: "验证分层架构依赖方向",
        Check:       validateLayerDependencies,
    },
}

func validateAggregateRootMethods(pkg *Package) error {
    for _, method := range pkg.Methods {
        if isAggregateRoot(method.Receiver) && method.ReturnsError {
            // 检查方法实现是否正确处理错误
            if hasPanicCall(method.Body) {
                return fmt.Errorf("aggregate root method %s panics instead of returning error", 
                    method.Name)
            }
        }
    }
    return nil
}
```

### 3. 调试和监控支持

#### 聚合状态查看器
```go
type AggregateDebugger struct {
    eventStore EventStore
    logger     *zap.Logger
}

func (ad *AggregateDebugger) GetAggregateHistory(aggregateID interface{}) (*AggregateHistory, error) {
    events, err := ad.eventStore.GetEventsForAggregate(aggregateID, 0)
    if err != nil {
        return nil, err
    }
    
    history := &AggregateHistory{
        AggregateID: aggregateID,
        Events:      make([]EventSnapshot, len(events)),
    }
    
    for i, event := range events {
        history.Events[i] = EventSnapshot{
            Sequence:    i + 1,
            Type:        event.EventName(),
            Data:        event,
            OccurredAt:  event.OccurredOn(),
        }
    }
    
    return history, nil
}

type EventSnapshot struct {
    Sequence   int
    Type       string
    Data       interface{}
    OccurredAt time.Time
}
```

这个设计指南为实现稳健的DDD+CQRS架构提供了完整的理论基础和实践指导，既保持了架构的严谨性，又考虑了开发者的实际使用需求。