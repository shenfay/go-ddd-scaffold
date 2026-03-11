# Go DDD Scaffold 领域模型设计文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的领域模型设计，包括核心领域概念、实体关系、聚合设计以及业务规则的领域表达。

## 领域概述

### 核心业务场景
项目定位为通用企业服务平台，主要解决以下业务场景：
- 多租户环境下的用户管理
- 基于RBAC的权限控制系统
- 安全可靠的认证授权机制
- 完整的操作审计追踪

### 领域边界定义

```
核心领域 (Core Domain):
├── User Management    (用户管理)
├── Tenant Management  (租户管理)
├── Authentication     (认证服务)
└── Authorization      (授权服务)

支撑领域 (Supporting Domain):
├── Audit Service      (审计服务)
└── Notification       (通知服务-预留)

通用领域 (Generic Domain):
├── Configuration      (配置管理)
├── Logging            (日志服务)
└── Validation         (验证服务)
```

## 核心领域模型设计

### 1. 用户领域 (User Domain)

#### 领域概念模型
```
User Aggregate (用户聚合)
├── User Entity (用户实体)
├── UserProfile Value Object (用户档案)
├── UserCredentials Value Object (用户凭证)
├── UserStatus Value Object (用户状态)
└── UserTenant Association (用户租户关联)
```

#### 实体设计

**User 实体**
```go
// User 用户聚合根
type User struct {
    ddd.BaseEntity

    username       *UserName
    email          *Email
    password       *HashedPassword
    status         UserStatus
    displayName    string
    firstName      string
    lastName       string
    gender         UserGender
    phoneNumber    string
    avatarURL      string
    lastLoginAt    *time.Time
    loginCount     int
    lockedUntil    *time.Time
    failedAttempts int
    createdAt      time.Time
    updatedAt      time.Time
}

// 用户行为方法 - 所有业务方法都会发布领域事件

// Activate 激活用户
func (u *User) Activate() error {
    if u.status != UserStatusPending {
        return ddd.NewBusinessError("USER_NOT_PENDING", "user is not in pending status")
    }
    u.status = UserStatusActive
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserActivatedEvent(u.ID().(UserID))
    u.ApplyEvent(event)
    
    return nil
}

// Deactivate 禁用用户
func (u *User) Deactivate(reason string) error {
    if u.status == UserStatusInactive {
        return ddd.NewBusinessError("USER_ALREADY_INACTIVE", "user is already inactive")
    }
    u.status = UserStatusInactive
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserDeactivatedEvent(u.ID().(UserID), reason)
    u.ApplyEvent(event)
    
    return nil
}

// Lock 锁定用户
func (u *User) Lock(duration time.Duration, reason string) error {
    if u.status == UserStatusLocked {
        return ddd.NewBusinessError("USER_ALREADY_LOCKED", "user is already locked")
    }
    lockUntil := time.Now().Add(duration)
    u.status = UserStatusLocked
    u.lockedUntil = &lockUntil
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserLockedEvent(u.ID().(UserID), reason, lockUntil)
    u.ApplyEvent(event)
    
    return nil
}

// Unlock 解锁用户
func (u *User) Unlock() error {
    if u.status != UserStatusLocked {
        return ddd.NewBusinessError("USER_NOT_LOCKED", "user is not locked")
    }
    u.status = UserStatusActive
    u.lockedUntil = nil
    u.failedAttempts = 0
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserUnlockedEvent(u.ID().(UserID))
    u.ApplyEvent(event)
    
    return nil
}

// RecordLogin 记录登录
func (u *User) RecordLogin(ipAddress, userAgent string) {
    now := time.Now()
    u.lastLoginAt = &now
    u.loginCount++
    u.failedAttempts = 0
    u.updatedAt = now
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserLoggedInEvent(u.ID().(UserID), ipAddress, userAgent)
    u.ApplyEvent(event)
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string, ipAddress string) error {
    if !u.password.Matches(oldPassword) {
        return ddd.NewBusinessError("INVALID_OLD_PASSWORD", "invalid old password")
    }
    u.password = NewHashedPassword(newPassword)
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserPasswordChangedEvent(u.ID().(UserID), ipAddress)
    u.ApplyEvent(event)
    
    return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
    oldEmail := u.email.Value()
    email, err := NewEmail(newEmail)
    if err != nil {
        return err
    }
    u.email = email
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserEmailChangedEvent(u.ID().(UserID), oldEmail, newEmail)
    u.ApplyEvent(event)
    
    return nil
}
```

#### 值对象设计

**UserID 值对象**
```go
// UserID 用户标识
type UserID struct {
    ddd.Int64Identity
}

// NewUserID 创建用户标识
func NewUserID(value int64) UserID {
    return UserID{Int64Identity: ddd.NewInt64Identity(value)}
}

// String 返回用户标识字符串
func (uid UserID) String() string {
    return fmt.Sprintf("user-%d", uid.Int64())
}
```

**UserName 值对象**
```go
// UserName 用户名值对象
type UserName struct {
    value string
}

// NewUserName 创建用户名
func NewUserName(value string) (*UserName, error) {
    un := &UserName{value: strings.TrimSpace(value)}
    if err := un.Validate(); err != nil {
        return nil, err
    }
    return un, nil
}

// Value 返回用户名值
func (un *UserName) Value() string {
    return un.value
}

// Validate 验证用户名
func (un *UserName) Validate() error {
    if un.value == "" {
        return &ddd.ValidationError{
            Field:   "username",
            Message: "username cannot be empty",
        }
    }
    if len(un.value) < 3 {
        return &ddd.ValidationError{
            Field:   "username",
            Message: "username must be at least 3 characters long",
        }
    }
    if len(un.value) > 50 {
        return &ddd.ValidationError{
            Field:   "username",
            Message: "username cannot exceed 50 characters",
        }
    }
    // 只允许字母、数字、下划线和连字符
    validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    if !validPattern.MatchString(un.value) {
        return &ddd.ValidationError{
            Field:   "username",
            Message: "username can only contain letters, numbers, underscores and hyphens",
        }
    }
    return nil
}

// Equals 比较用户名是否相等
func (un *UserName) Equals(other *UserName) bool {
    if other == nil {
        return false
    }
    return strings.EqualFold(un.value, other.value)
}
```

**Email 值对象**
```go
// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 创建邮箱
func NewEmail(value string) (*Email, error) {
    email := &Email{value: strings.TrimSpace(strings.ToLower(value))}
    if err := email.Validate(); err != nil {
        return nil, err
    }
    return email, nil
}

// Value 返回邮箱值
func (e *Email) Value() string {
    return e.value
}

// Validate 验证邮箱格式
func (e *Email) Validate() error {
    if e.value == "" {
        return &ddd.ValidationError{
            Field:   "email",
            Message: "email cannot be empty",
        }
    }
    emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailPattern.MatchString(e.value) {
        return &ddd.ValidationError{
            Field:   "email",
            Message: "invalid email format",
        }
    }
    return nil
}

// Equals 比较邮箱是否相等
func (e *Email) Equals(other *Email) bool {
    if other == nil {
        return false
    }
    return strings.EqualFold(e.value, other.value)
}
```

**HashedPassword 值对象**
```go
// HashedPassword 加密密码值对象
type HashedPassword struct {
    value string
}

// NewHashedPassword 创建加密密码
func NewHashedPassword(hashedValue string) *HashedPassword {
    return &HashedPassword{value: hashedValue}
}

// Value 返回加密密码值
func (hp *HashedPassword) Value() string {
    return hp.value
}

// Matches 检查密码是否匹配
func (hp *HashedPassword) Matches(plainPassword string) bool {
    // 实际应用中应该使用 bcrypt 等安全哈希算法
    return hp.value == plainPassword
}
```

### 2. 租户领域 (Tenant Domain)

#### 领域概念模型
```
Tenant Aggregate (租户聚合)
├── Tenant Entity (租户实体)
├── TenantConfig Value Object (租户配置)
├── TenantStatus Value Object (租户状态)
└── TenantMember Association (租户成员)
```

#### 实体设计

**Tenant 实体**
```go
// 租户实体 - 聚合根
type Tenant struct {
    id          TenantID
    code        TenantCode      // 租户编码
    name        string          // 租户名称
    config      TenantConfig    // 租户配置
    status      TenantStatus    // 租户状态
    members     []TenantMember  // 租户成员
    createdAt   time.Time
    updatedAt   time.Time
}

func (t *Tenant) AddMember(userID UserID, role Role) error {
    // 业务规则验证
    if t.isMember(userID) {
        return ErrUserAlreadyMember
    }
    
    if len(t.members) >= t.config.MaxMembers {
        return ErrMaxMembersReached
    }
    
    member := TenantMember{
        UserID: userID,
        Role:   role,
        Joined: time.Now(),
    }
    
    t.members = append(t.members, member)
    t.updatedAt = time.Now()
    return nil
}

func (t *Tenant) RemoveMember(userID UserID) error {
    for i, member := range t.members {
        if member.UserID.Equals(userID) {
            t.members = append(t.members[:i], t.members[i+1:]...)
            t.updatedAt = time.Now()
            return nil
        }
    }
    return ErrUserNotMember
}
```

#### 值对象设计

**TenantID 值对象**
```go
// 租户ID - 全局唯一标识
type TenantID int64

func NewTenantID(id int64) TenantID {
    if id <= 0 {
        panic("tenant id must be positive")
    }
    return TenantID(id)
}

func (id TenantID) Equals(other TenantID) bool {
    return id == other
}
```

**TenantCode 值对象**
```go
// 租户编码 - 业务唯一标识
type TenantCode string

func NewTenantCode(code string) (TenantCode, error) {
    code = strings.ToUpper(strings.TrimSpace(code))
    
    if len(code) < 3 {
        return "", errors.New("tenant code too short")
    }
    
    if len(code) > 10 {
        return "", errors.New("tenant code too long")
    }
    
    if !tenantCodePattern.MatchString(code) {
        return "", errors.New("invalid tenant code format")
    }
    
    return TenantCode(code), nil
}

func (tc TenantCode) String() string {
    return string(tc)
}
```

### 3. 认证领域 (Authentication Domain)

#### 领域概念模型
```
Authentication Service (认证服务)
├── JWT Token Value Object
├── RefreshToken Value Object
├── AuthenticationResult Value Object
└── TokenPair Value Object
```

#### 核心服务设计

**AuthenticationService**
```go
// 认证服务 - 处理用户认证逻辑
type AuthenticationService struct {
    userRepo     UserRepository
    tokenService TokenService
    logger       *zap.Logger
}

func (as *AuthenticationService) Authenticate(username Username, password Password) (*AuthenticationResult, error) {
    // 1. 查找用户
    user, err := as.userRepo.FindByUsername(username)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // 2. 验证密码
    if !user.Credentials().VerifyPassword(password) {
        as.logger.Warn("invalid password attempt", 
            zap.String("username", username.String()))
        return nil, ErrInvalidCredentials
    }
    
    // 3. 检查账户状态
    if !user.IsActive() {
        return nil, ErrAccountDisabled
    }
    
    // 4. 生成令牌
    tokenPair, err := as.tokenService.GenerateTokens(user.ID())
    if err != nil {
        return nil, err
    }
    
    // 5. 更新最后登录时间
    user.RecordLogin()
    if err := as.userRepo.Update(user); err != nil {
        as.logger.Error("failed to update user login info", zap.Error(err))
    }
    
    return &AuthenticationResult{
        UserID:    user.ID(),
        Username:  user.Username(),
        TokenPair: tokenPair,
    }, nil
}

func (as *AuthenticationService) RefreshToken(refreshToken RefreshToken) (*TokenPair, error) {
    // 1. 验证刷新令牌
    claims, err := as.tokenService.ParseRefreshToken(refreshToken)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    // 2. 检查用户是否存在且激活
    user, err := as.userRepo.FindByID(claims.UserID)
    if err != nil || !user.IsActive() {
        return nil, ErrInvalidUser
    }
    
    // 3. 生成新的令牌对
    return as.tokenService.GenerateTokens(claims.UserID)
}
```

### 4. 授权领域 (Authorization Domain)

#### 领域概念模型
```
Authorization Service (授权服务)
├── Permission Entity (权限实体)
├── Role Entity (角色实体)
├── Policy Value Object (策略)
└── AccessControlList Value Object (访问控制列表)
```

#### RBAC模型实现

**Permission 实体**
```go
// 权限实体
type Permission struct {
    id          PermissionID
    resource    ResourceName    // 资源名称
    action      ActionName      // 操作名称
    description string          // 权限描述
    createdAt   time.Time
}

func NewPermission(resource ResourceName, action ActionName, desc string) *Permission {
    return &Permission{
        id:          NewPermissionID(),
        resource:    resource,
        action:      action,
        description: desc,
        createdAt:   time.Now(),
    }
}

func (p *Permission) Identifier() string {
    return fmt.Sprintf("%s:%s", p.resource, p.action)
}
```

**Role 实体**
```go
// 角色实体
type Role struct {
    id          RoleID
    name        RoleName
    description string
    permissions []PermissionID  // 关联的权限
    tenantID    TenantID        // 所属租户
    createdAt   time.Time
}

func (r *Role) AddPermission(permissionID PermissionID) {
    if !r.hasPermission(permissionID) {
        r.permissions = append(r.permissions, permissionID)
    }
}

func (r *Role) RemovePermission(permissionID PermissionID) {
    for i, pid := range r.permissions {
        if pid.Equals(permissionID) {
            r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
            break
        }
    }
}

func (r *Role) HasPermission(resource ResourceName, action ActionName) bool {
    // 检查角色是否拥有指定权限
    for _, permissionID := range r.permissions {
        // 通过仓储获取权限详情进行比较
        permission := getPermissionByID(permissionID)
        if permission.Resource().Equals(resource) && 
           permission.Action().Equals(action) {
            return true
        }
    }
    return false
}
```

## 领域关系设计

### 用户-租户多对多关系

```
User (1) ──── (*) UserTenant (*) ──── (1) Tenant
              │
              └── Role (在特定租户下的角色)
```

**UserTenant 关联实体**
```go
// 用户租户关联 - 值对象
type UserTenant struct {
    userID   UserID
    tenantID TenantID
    roleID   RoleID      // 在该租户下的角色
    joinedAt time.Time   // 加入时间
}

func NewUserTenant(userID UserID, tenantID TenantID, roleID RoleID) UserTenant {
    return UserTenant{
        userID:   userID,
        tenantID: tenantID,
        roleID:   roleID,
        joinedAt: time.Now(),
    }
}

func (ut UserTenant) UserID() UserID     { return ut.userID }
func (ut UserTenant) TenantID() TenantID { return ut.tenantID }
func (ut UserTenant) RoleID() RoleID     { return ut.roleID }
```

### 聚合边界设计

```
┌─────────────────────────────────────────────────┐
│                User Aggregate                    │
│  User + UserProfile + UserCredentials           │
│  (通过 UserID 聚合根标识)                        │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│               Tenant Aggregate                   │
│  Tenant + TenantConfig + TenantMembers          │
│  (通过 TenantID 聚合根标识)                      │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│            Permission Aggregate                  │
│  Permission (独立聚合)                           │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│               Role Aggregate                     │
│  Role + RolePermissions (通过 RoleID 标识)       │
└─────────────────────────────────────────────────┘
```

## 领域事件设计

### 领域事件基类

```go
// DomainEvent 领域事件接口
type DomainEvent interface {
    EventName() string
    OccurredOn() time.Time
    AggregateID() interface{}
    Version() int
    Metadata() map[string]interface{}
}

// BaseEvent 领域事件基础结构
type BaseEvent struct {
    eventName   string
    aggregateID interface{}
    version     int
    occurredOn  time.Time
    metadata    map[string]interface{}
}

// EventName 返回事件名称
func (e *BaseEvent) EventName() string { return e.eventName }

// OccurredOn 返回事件发生时间
func (e *BaseEvent) OccurredOn() time.Time { return e.occurredOn }

// AggregateID 返回聚合根ID
func (e *BaseEvent) AggregateID() interface{} { return e.aggregateID }

// Version 返回事件版本
func (e *BaseEvent) Version() int { return e.version }

// Metadata 返回事件元数据
func (e *BaseEvent) Metadata() map[string]interface{} { return e.metadata }

// SetMetadata 设置事件元数据
func (e *BaseEvent) SetMetadata(key string, value interface{}) {
    if e.metadata == nil {
        e.metadata = make(map[string]interface{})
    }
    e.metadata[key] = value
}
```

### 用户领域事件

```go
// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
    *BaseEvent
    UserID       UserID    `json:"user_id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    RegisteredAt time.Time `json:"registered_at"`
}

// UserActivatedEvent 用户激活事件
type UserActivatedEvent struct {
    *BaseEvent
    UserID      UserID    `json:"user_id"`
    ActivatedAt time.Time `json:"activated_at"`
}

// UserDeactivatedEvent 用户禁用事件
type UserDeactivatedEvent struct {
    *BaseEvent
    UserID        UserID    `json:"user_id"`
    Reason        string    `json:"reason"`
    DeactivatedAt time.Time `json:"deactivated_at"`
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
    *BaseEvent
    UserID    UserID    `json:"user_id"`
    LoginAt   time.Time `json:"login_at"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
    *BaseEvent
    UserID    UserID    `json:"user_id"`
    ChangedAt time.Time `json:"changed_at"`
    IPAddress string    `json:"ip_address"`
}

// UserEmailChangedEvent 用户邮箱修改事件
type UserEmailChangedEvent struct {
    *BaseEvent
    UserID    UserID    `json:"user_id"`
    OldEmail  string    `json:"old_email"`
    NewEmail  string    `json:"new_email"`
    ChangedAt time.Time `json:"changed_at"`
}

// UserLockedEvent 用户锁定事件
type UserLockedEvent struct {
    *BaseEvent
    UserID      UserID    `json:"user_id"`
    Reason      string    `json:"reason"`
    LockedUntil time.Time `json:"locked_until"`
    LockedAt    time.Time `json:"locked_at"`
}

// UserUnlockedEvent 用户解锁事件
type UserUnlockedEvent struct {
    *BaseEvent
    UserID     UserID    `json:"user_id"`
    UnlockedAt time.Time `json:"unlocked_at"`
}

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
    *BaseEvent
    UserID        UserID    `json:"user_id"`
    UpdatedFields []string  `json:"updated_fields"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// UserFailedLoginAttemptEvent 用户登录失败事件
type UserFailedLoginAttemptEvent struct {
    *BaseEvent
    UserID    UserID    `json:"user_id"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
    Reason    string    `json:"reason"`
    AttemptAt time.Time `json:"attempt_at"`
}
```

### 租户领域事件

```go
// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct {
    *BaseEvent
    TenantID  TenantID    `json:"tenant_id"`
    Code      string      `json:"code"`
    Name      string      `json:"name"`
    OwnerID   user.UserID `json:"owner_id"`
    CreatedAt time.Time   `json:"created_at"`
}

// TenantActivatedEvent 租户激活事件
type TenantActivatedEvent struct {
    *BaseEvent
    TenantID    TenantID  `json:"tenant_id"`
    ActivatedAt time.Time `json:"activated_at"`
}

// TenantDeactivatedEvent 租户停用事件
type TenantDeactivatedEvent struct {
    *BaseEvent
    TenantID      TenantID  `json:"tenant_id"`
    Reason        string    `json:"reason"`
    DeactivatedAt time.Time `json:"deactivated_at"`
}

// TenantSuspendedEvent 租户暂停事件
type TenantSuspendedEvent struct {
    *BaseEvent
    TenantID    TenantID  `json:"tenant_id"`
    Reason      string    `json:"reason"`
    SuspendedAt time.Time `json:"suspended_at"`
}

// TenantNameChangedEvent 租户名称变更事件
type TenantNameChangedEvent struct {
    *BaseEvent
    TenantID  TenantID  `json:"tenant_id"`
    OldName   string    `json:"old_name"`
    NewName   string    `json:"new_name"`
    ChangedAt time.Time `json:"changed_at"`
}

// TenantConfigChangedEvent 租户配置变更事件
type TenantConfigChangedEvent struct {
    *BaseEvent
    TenantID    TenantID    `json:"tenant_id"`
    ConfigKey   string      `json:"config_key"`
    ConfigValue interface{} `json:"config_value"`
    ChangedAt   time.Time   `json:"changed_at"`
}

// TenantMemberAddedEvent 租户成员添加事件
type TenantMemberAddedEvent struct {
    *BaseEvent
    TenantID TenantID    `json:"tenant_id"`
    UserID   user.UserID `json:"user_id"`
    Role     string      `json:"role"`
    AddedBy  user.UserID `json:"added_by"`
    AddedAt  time.Time   `json:"added_at"`
}

// TenantMemberRemovedEvent 租户成员移除事件
type TenantMemberRemovedEvent struct {
    *BaseEvent
    TenantID  TenantID    `json:"tenant_id"`
    UserID    user.UserID `json:"user_id"`
    RemovedBy user.UserID `json:"removed_by"`
    RemovedAt time.Time   `json:"removed_at"`
}

// TenantMemberRoleChangedEvent 租户成员角色变更事件
type TenantMemberRoleChangedEvent struct {
    *BaseEvent
    TenantID  TenantID    `json:"tenant_id"`
    UserID    user.UserID `json:"user_id"`
    OldRole   string      `json:"old_role"`
    NewRole   string      `json:"new_role"`
    ChangedBy user.UserID `json:"changed_by"`
    ChangedAt time.Time   `json:"changed_at"`
}
```

### 事件处理示例

```go
// 领域事件处理器接口
type EventHandler interface {
    Handle(ctx context.Context, event DomainEvent) error
}

// 用户事件处理器
type UserEventHandler struct {
    emailService EmailService
    auditService AuditService
}

func (h *UserEventHandler) Handle(ctx context.Context, event DomainEvent) error {
    switch e := event.(type) {
    case *UserRegisteredEvent:
        return h.handleUserRegistered(ctx, e)
    case *UserActivatedEvent:
        return h.handleUserActivated(ctx, e)
    // ... 其他事件处理
    }
    return nil
}

func (h *UserEventHandler) handleUserRegistered(ctx context.Context, event *UserRegisteredEvent) error {
    // 发送欢迎邮件
    if err := h.emailService.SendWelcomeEmail(event.Email, event.Username); err != nil {
        return err
    }
    
    // 记录审计日志
    return h.auditService.LogEvent(AuditLog{
        EventType: "USER_REGISTERED",
        UserID:    event.UserID,
        Timestamp: event.RegisteredAt,
        Details: map[string]interface{}{
            "username": event.Username,
            "email":    event.Email,
        },
    })
}
```

## 领域服务

### 认证服务

```go
// AuthenticationService 认证服务
type AuthenticationService struct {
    userRepo       UserRepository
    tokenService   TokenService
    passwordPolicy PasswordPolicyService
}

// Authenticate 用户认证
func (s *AuthenticationService) Authenticate(ctx context.Context, usernameOrEmail, password string, ipAddress, userAgent string) (*AuthenticateResult, error) {
    // 1. 查找用户
    var u *User
    var err error
    
    // 尝试作为邮箱查找
    u, err = s.userRepo.FindByEmail(ctx, usernameOrEmail)
    if err != nil {
        // 尝试作为用户名查找
        u, err = s.userRepo.FindByUsername(ctx, usernameOrEmail)
        if err != nil {
            return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
        }
    }

    // 2. 验证密码
    if !u.Password().Matches(password) {
        u.RecordFailedLogin(ipAddress, userAgent, "invalid_password")
        _ = s.userRepo.Save(ctx, u)
        return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
    }

    // 3. 检查账户状态
    if !u.CanLogin() {
        return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "account cannot login")
    }

    // 4. 记录成功登录
    u.RecordLogin(ipAddress, userAgent)
    if err := s.userRepo.Save(ctx, u); err != nil {
        return nil, err
    }

    // 5. 生成令牌
    tokenPair, err := s.tokenService.GenerateTokenPair(u.ID().(UserID))
    if err != nil {
        return nil, err
    }

    return &AuthenticateResult{
        UserID:       u.ID().(UserID),
        Username:     u.Username().Value(),
        Email:        u.Email().Value(),
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
        ExpiresAt:    tokenPair.ExpiresAt,
    }, nil
}
```

### 密码策略服务

```go
// PasswordPolicyService 密码策略服务接口
type PasswordPolicyService interface {
    Validate(password string) error
    GetPolicy() PasswordPolicy
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
    MinLength           int
    MaxLength           int
    RequireUppercase    bool
    RequireLowercase    bool
    RequireDigit        bool
    RequireSpecialChar  bool
    SpecialChars        string
    DisallowUsername    bool
    MaxRepeatedChars    int
    PasswordHistorySize int
}
```

### 租户服务

```go
// TenantService 租户领域服务
type TenantService struct {
    tenantRepo TenantRepository
    userRepo   UserRepository
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, code, name string, ownerID user.UserID) (*Tenant, error) {
    // 1. 检查租户编码是否已存在
    if _, err := s.tenantRepo.FindByCode(ctx, code); err == nil {
        return nil, ddd.NewBusinessError("TENANT_CODE_EXISTS", "tenant code already exists")
    }

    // 2. 检查所有者是否存在
    if _, err := s.userRepo.FindByID(ctx, ownerID); err != nil {
        return nil, ddd.NewBusinessError("OWNER_NOT_FOUND", "owner user not found")
    }

    // 3. 创建租户
    tenant, err := NewTenant(code, name, ownerID)
    if err != nil {
        return nil, err
    }

    // 4. 保存租户
    if err := s.tenantRepo.Save(ctx, tenant); err != nil {
        return nil, err
    }

    // 5. 添加所有者为成员
    member := &TenantMember{
        UserID:   ownerID,
        TenantID: tenant.ID().(TenantID),
        Role:     TenantRoleOwner,
        JoinedAt: time.Now().Format(time.RFC3339),
    }
    if err := s.tenantRepo.AddMember(ctx, tenant.ID().(TenantID), member); err != nil {
        return nil, err
    }

    return tenant, nil
}

// AddMember 添加成员到租户
func (s *TenantService) AddMember(ctx context.Context, tenantID TenantID, userID, addedBy user.UserID, role TenantRole) error {
    // 1. 检查租户是否存在
    tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
    if err != nil {
        return ddd.ErrAggregateNotFound
    }

    // 2. 检查租户状态
    if !tenant.IsActive() {
        return ddd.NewBusinessError("TENANT_NOT_ACTIVE", "tenant is not active")
    }

    // 3. 检查成员数量限制
    members, err := s.tenantRepo.FindMembers(ctx, tenantID)
    if err != nil {
        return err
    }

    if !tenant.CanAddMember(len(members)) {
        return ddd.NewBusinessError("MAX_MEMBERS_REACHED", "tenant has reached maximum member limit")
    }

    // 4. 检查操作权限
    addedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, addedBy)
    if err != nil {
        return ddd.NewBusinessError("OPERATOR_NOT_MEMBER", "operator is not a member")
    }

    if addedByMember.Role != TenantRoleOwner && addedByMember.Role != TenantRoleAdmin {
        return ddd.NewBusinessError("INSUFFICIENT_PERMISSIONS", "insufficient permissions")
    }

    // 5. 添加成员
    member := &TenantMember{
        UserID:   userID,
        TenantID: tenantID,
        Role:     role,
        JoinedAt: time.Now().Format(time.RFC3339),
    }

    if err := s.tenantRepo.AddMember(ctx, tenantID, member); err != nil {
        return err
    }

    // 6. 发布领域事件
    event := NewTenantMemberAddedEvent(tenantID, userID, addedBy, role)
    tenant.ApplyEvent(event)

    return s.tenantRepo.Save(ctx, tenant)
}
```

## 业务规则约束

### 用户领域规则
1. 用户名必须唯一且符合命名规范（3-50字符，字母数字下划线连字符）
2. 邮箱必须唯一且格式正确
3. 密码必须满足复杂度要求（8位以上，包含大小写字母、数字、特殊字符）
4. 用户状态变更需要发布领域事件
5. 删除用户需要软删除保护
6. 登录失败超过阈值自动锁定账户

### 租户领域规则
1. 租户编码必须全局唯一（3-20字符，大写字母数字下划线连字符）
2. 租户成员数量不能超过配置上限
3. 租户状态影响其下所有用户的访问权限
4. 跨租户数据必须严格隔离
5. 只有 Owner 和 Admin 可以添加/移除成员
6. 所有者不能被移除，所有权可以转移

### 认证授权规则
1. 登录失败次数超过阈值需要锁定账户
2. JWT令牌必须设置合理的过期时间
3. 权限检查必须在业务操作前完成
4. 敏感操作需要二次验证
5. 密码修改需要验证旧密码（管理员重置除外）

这个领域模型设计文档为项目提供了清晰的业务概念模型和实现指导。