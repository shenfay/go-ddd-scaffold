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
// 用户实体 - 聚合根
type User struct {
    id          UserID           // 用户唯一标识
    username    Username         // 用户名
    email       Email            // 邮箱
    credentials UserCredentials  // 凭证信息
    profile     UserProfile      // 用户档案
    status      UserStatus       // 用户状态
    createdAt   time.Time        // 创建时间
    updatedAt   time.Time        // 更新时间
}

// 用户行为方法
func (u *User) ChangePassword(current, new Password) error {
    if !u.credentials.VerifyPassword(current) {
        return ErrInvalidCurrentPassword
    }
    
    if err := u.credentials.SetPassword(new); err != nil {
        return err
    }
    
    u.updatedAt = time.Now()
    return nil
}

func (u *User) UpdateProfile(profile UserProfile) {
    u.profile = profile
    u.updatedAt = time.Now()
}

func (u *User) Activate() {
    u.status = UserStatusActive
    u.updatedAt = time.Now()
}

func (u *User) Deactivate() {
    u.status = UserStatusInactive
    u.updatedAt = time.Now()
}
```

#### 值对象设计

**UserID 值对象**
```go
// 用户ID - 使用Snowflake算法生成
type UserID int64

func NewUserID(id int64) UserID {
    if id <= 0 {
        panic("user id must be positive")
    }
    return UserID(id)
}

func (id UserID) String() string {
    return strconv.FormatInt(int64(id), 10)
}

func (id UserID) Equals(other UserID) bool {
    return id == other
}
```

**Username 值对象**
```go
// 用户名 - 业务规则封装
type Username string

func NewUsername(name string) (Username, error) {
    name = strings.TrimSpace(name)
    
    if len(name) < 3 {
        return "", errors.New("username too short, minimum 3 characters")
    }
    
    if len(name) > 20 {
        return "", errors.New("username too long, maximum 20 characters")
    }
    
    if !usernamePattern.MatchString(name) {
        return "", errors.New("username contains invalid characters")
    }
    
    return Username(name), nil
}

func (u Username) String() string {
    return string(u)
}
```

**Email 值对象**
```go
// 邮箱地址 - 格式验证和标准化
type Email string

func NewEmail(email string) (Email, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    
    if !emailPattern.MatchString(email) {
        return "", errors.New("invalid email format")
    }
    
    return Email(email), nil
}

func (e Email) String() string {
    return string(e)
}

func (e Email) GetDomain() string {
    parts := strings.Split(string(e), "@")
    if len(parts) == 2 {
        return parts[1]
    }
    return ""
}
```

**UserCredentials 值对象**
```go
// 用户凭证 - 密码和认证相关信息
type UserCredentials struct {
    passwordHash string
    salt         string
    lastLogin    *time.Time
    loginCount   int
}

func NewUserCredentials(password Password) (*UserCredentials, error) {
    salt := generateSalt()
    hash, err := hashPassword(password.String(), salt)
    if err != nil {
        return nil, err
    }
    
    return &UserCredentials{
        passwordHash: hash,
        salt:         salt,
        loginCount:   0,
    }, nil
}

func (c *UserCredentials) VerifyPassword(password Password) bool {
    hash, err := hashPassword(password.String(), c.salt)
    if err != nil {
        return false
    }
    return hash == c.passwordHash
}

func (c *UserCredentials) SetPassword(password Password) error {
    salt := generateSalt()
    hash, err := hashPassword(password.String(), salt)
    if err != nil {
        return err
    }
    
    c.passwordHash = hash
    c.salt = salt
    return nil
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

### 核心领域事件

```go
// 用户相关事件
type UserRegisteredEvent struct {
    UserID    UserID
    Username  Username
    Email     Email
    Timestamp time.Time
}

type UserActivatedEvent struct {
    UserID    UserID
    Timestamp time.Time
}

type UserDeactivatedEvent struct {
    UserID    UserID
    Reason    string
    Timestamp time.Time
}

// 租户相关事件
type TenantCreatedEvent struct {
    TenantID   TenantID
    TenantCode TenantCode
    CreatorID  UserID
    Timestamp  time.Time
}

type UserAddedToTenantEvent struct {
    UserID    UserID
    TenantID  TenantID
    RoleID    RoleID
    AddedBy   UserID
    Timestamp time.Time
}
```

### 事件处理示例

```go
// 领域事件处理器
type UserEventHandler struct {
    emailService EmailService
    auditService AuditService
}

func (h *UserEventHandler) HandleUserRegistered(event UserRegisteredEvent) {
    // 发送欢迎邮件
    h.emailService.SendWelcomeEmail(event.Email, event.Username)
    
    // 记录审计日志
    h.auditService.LogEvent(AuditLog{
        EventType: "USER_REGISTERED",
        UserID:    event.UserID,
        Timestamp: event.Timestamp,
        Details: map[string]interface{}{
            "username": event.Username.String(),
            "email":    event.Email.String(),
        },
    })
}
```

## 业务规则约束

### 用户领域规则
1. 用户名必须唯一且符合命名规范
2. 邮箱必须唯一且格式正确
3. 密码必须满足复杂度要求
4. 用户状态变更需要审计记录
5. 删除用户需要软删除保护

### 租户领域规则
1. 租户编码必须全局唯一
2. 租户成员数量不能超过配置上限
3. 租户状态影响其下所有用户的访问权限
4. 跨租户数据必须严格隔离

### 认证授权规则
1. 登录失败次数超过阈值需要锁定账户
2. JWT令牌必须设置合理的过期时间
3. 权限检查必须在业务操作前完成
4. 敏感操作需要二次验证

这个领域模型设计文档为项目提供了清晰的业务概念模型和实现指导。