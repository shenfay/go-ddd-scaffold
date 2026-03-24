# 领域模型设计

本文档详细介绍 Go DDD Scaffold 项目的领域模型设计，包括聚合根、值对象、领域事件等核心概念。

## 📋 领域驱动设计核心概念

### 战略设计 vs 战术设计

```
战略设计（Strategic Design）
├── 限界上下文（Bounded Context）
├── 通用语言（Ubiquitous Language）
└── 上下文映射（Context Mapping）

战术设计（Tactical Design）⭐ 本文档重点
├── 实体（Entity）
├── 值对象（Value Object）
├── 聚合根（Aggregate Root）
├── 领域服务（Domain Service）
├── 领域事件（Domain Event）
└── 仓储（Repository）
```

---

## 🎯 实体（Entity）

### 什么是实体？

**实体**是具有唯一标识和生命周期的领域对象。

**特点：**
- 有唯一标识（ID）
- 有生命周期（创建、更新、删除）
- 可变状态
- 封装业务逻辑

### Base Entity 实现

```go
// domain/shared/kernel/entity.go
package kernel

import "time"

// Entity 基础实体
type Entity struct {
    id        ID
    createdAt time.Time
    updatedAt time.Time
}

// NewEntity 创建新实体
func NewEntity() *Entity {
    return &Entity{
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
}

// ReconstructEntity 从持久化层重建实体
func ReconstructEntity(id int64, createdAt, updatedAt time.Time) *Entity {
    return &Entity{
        id:        ID{id},
        createdAt: createdAt,
        updatedAt: updatedAt,
    }
}

// ID 获取实体 ID
func (e *Entity) ID() ID {
    return e.id
}

// CreatedAt 获取创建时间
func (e *Entity) CreatedAt() time.Time {
    return e.createdAt
}

// UpdatedAt 获取更新时间
func (e *Entity) UpdatedAt() time.Time {
    return e.updatedAt
}

// SetID 设置 ID（仅在创建时使用）
func (e *Entity) SetID(id int64) {
    e.id = ID{id}
}

// touch 更新更新时间
func (e *Entity) touch() {
    e.updatedAt = time.Now()
}
```

### ID 值对象

```go
// domain/shared/kernel/id.go
package kernel

// ID 实体 ID 值对象
type ID struct {
    value int64
}

// NewID 创建新 ID
func NewID(value int64) ID {
    if value <= 0 {
        panic("id must be positive")
    }
    return ID{value: value}
}

// Value 获取 ID 值
func (id ID) Value() int64 {
    return id.value
}

// Equals 比较两个 ID
func (id ID) Equals(other ID) bool {
    return id.value == other.value
}

// String 字符串表示
func (id ID) String() string {
    return fmt.Sprintf("%d", id.value)
}
```

---

## 💎 值对象（Value Object）

### 什么是值对象？

**值对象**是只描述事物属性，没有唯一标识的领域对象。

**特点：**
- 无唯一标识
- 不可变（Immutable）
- 可替换
- 封装验证逻辑

### 值对象示例：Email

```go
// domain/user/valueobject/email.go
package valueobject

import (
    "regexp"
    "strings"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 创建邮箱（工厂方法）
func NewEmail(value string) (Email, error) {
    // 1. 清理输入
    value = strings.TrimSpace(strings.ToLower(value))
    
    // 2. 验证格式
    if !isValidEmail(value) {
        return Email{}, kernel.FieldError(
            "email",
            "无效的邮箱格式",
            value,
        )
    }
    
    return Email{value: value}, nil
}

// String 字符串表示
func (e Email) String() string {
    return e.value
}

// Equals 比较两个邮箱
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
    pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
    return regexp.MustCompile(pattern).MatchString(email)
}
```

### 值对象示例：Username

```go
// domain/user/valueobject/username.go
package valueobject

import (
    "regexp"
    "strings"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// Username 用户名值对象
type Username struct {
    value string
}

// NewUsername 创建用户名
func NewUsername(value string) (Username, error) {
    value = strings.TrimSpace(strings.ToLower(value))
    
    // 验证长度
    if len(value) < 3 || len(value) > 50 {
        return Username{}, kernel.FieldError(
            "username",
            "用户名长度必须在 3-50 个字符之间",
            value,
        )
    }
    
    // 验证格式
    if !isValidUsername(value) {
        return Username{}, kernel.FieldError(
            "username",
            "用户名只能包含小写字母、数字和下划线",
            value,
        )
    }
    
    return Username{value: value}, nil
}

func (u Username) String() string {
    return u.value
}

func isValidUsername(username string) bool {
    pattern := `^[a-z][a-z0-9_]{2,49}$`
    return regexp.MustCompile(pattern).MatchString(username)
}
```

### 值对象示例：Password

```go
// domain/user/valueobject/password.go
package valueobject

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// Password 密码值对象
type Password struct {
    hash string
}

// NewPassword 创建密码
func NewPassword(plainText string) (Password, error) {
    // 验证密码强度
    if err := validatePasswordStrength(plainText); err != nil {
        return Password{}, kernel.FieldError("password", err.Error(), plainText)
    }
    
    // 哈希密码
    hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
    if err != nil {
        return Password{}, err
    }
    
    return Password{hash: string(hash)}, nil
}

// Verify 验证密码
func (p Password) Verify(plainText string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainText))
    return err == nil
}

func validatePasswordStrength(password string) error {
    if len(password) < 8 {
        return errors.New("密码长度至少为 8 个字符")
    }
    
    hasUpper := false
    hasLower := false
    hasNumber := false
    hasSpecial := false
    
    for _, ch := range password {
        switch {
        case 'A' <= ch && ch <= 'Z':
            hasUpper = true
        case 'a' <= ch && ch <= 'z':
            hasLower = true
        case '0' <= ch && ch <= '9':
            hasNumber = true
        default:
            hasSpecial = true
        }
    }
    
    if !hasUpper {
        return errors.New("密码必须包含大写字母")
    }
    if !hasLower {
        return errors.New("密码必须包含小写字母")
    }
    if !hasNumber {
        return errors.New("密码必须包含数字")
    }
    if !hasSpecial {
        return errors.New("密码必须包含特殊字符")
    }
    
    return nil
}
```

### 值对象示例：UserStatus

```go
// domain/user/valueobject/status.go
package valueobject

// UserStatus 用户状态值对象
type UserStatus string

const (
    UserStatusPending   UserStatus = "pending"
    UserStatusActive    UserStatus = "active"
    UserStatusInactive  UserStatus = "inactive"
    UserStatusLocked    UserStatus = "locked"
    UserStatusDeleted   UserStatus = "deleted"
)

// IsValid 检查状态是否有效
func (s UserStatus) IsValid() bool {
    switch s {
    case UserStatusPending, UserStatusActive, 
         UserStatusInactive, UserStatusLocked, 
         UserStatusDeleted:
        return true
    default:
        return false
    }
}

// CanLogin 检查是否可以登录
func (s UserStatus) CanLogin() bool {
    return s == UserStatusActive
}

// IsActive 检查是否激活
func (s UserStatus) IsActive() bool {
    return s == UserStatusActive
}
```

---

## 🌳 聚合根（Aggregate Root）

### 什么是聚合根？

**聚合根**是聚合（Aggregate）的入口点，负责维护聚合内的一致性。

**特点：**
- 聚合的全局唯一标识
- 聚合的唯一访问入口
- 维护聚合内的一致性
- 发布领域事件

### User 聚合根完整实现

```go
// domain/user/aggregate/user.go
package aggregate

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
)

const (
    MaxLoginAttempts = 5
    LockDuration     = 30 * time.Minute
)

// User 用户聚合根
type User struct {
    *kernel.Entity
    username              vo.Username
    email                 vo.Email
    password              vo.Password
    status                vo.UserStatus
    failedLoginAttempts   int
    lockedUntil           *time.Time
    lastLoginAt           *time.Time
    lastLoginIP           string
}

// NewUser 创建新用户（构造函数）
func NewUser(username, email, password string) (*User, error) {
    // 验证并创建值对象
    uName, err := vo.NewUsername(username)
    if err != nil {
        return nil, err
    }
    
    uEmail, err := vo.NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    uPassword, err := vo.NewPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Entity:              kernel.NewEntity(),
        username:            uName,
        email:               uEmail,
        password:            uPassword,
        status:              vo.UserStatusPending,
        failedLoginAttempts: 0,
    }
    
    // 发布领域事件
    user.RecordEvent(&event.UserRegistered{
        UserID:  user.ID().Value(),
        Email:   user.email.String(),
        Created: time.Now(),
    })
    
    return user, nil
}

// Login 登录业务逻辑
func (u *User) Login(password string, ip string, userAgent string) error {
    // 1. 检查账户锁定
    if u.isLocked() {
        return kernel.NewBusinessError(
            CodeUserLocked,
            "账户已被锁定，请稍后再试",
        )
    }
    
    // 2. 检查状态
    if u.status != vo.UserStatusActive {
        return kernel.NewBusinessError(
            CodeUserInactive,
            "账户未激活",
        )
    }
    
    // 3. 验证密码
    if !u.password.Verify(password) {
        u.failedLoginAttempts++
        
        // 检查是否达到最大失败次数
        if u.failedLoginAttempts >= MaxLoginAttempts {
            u.Lock()
            return kernel.NewBusinessError(
                CodeUserLocked,
                "多次密码错误，账户已被锁定",
            )
        }
        
        return kernel.NewBusinessError(
            CodeInvalidCredentials,
            "用户名或密码错误",
        )
    }
    
    // 4. 登录成功
    u.failedLoginAttempts = 0
    now := time.Now()
    u.lastLoginAt = &now
    u.lastLoginIP = ip
    
    u.RecordEvent(&event.UserLoggedIn{
        UserID:    u.ID().Value(),
        Email:     u.email.String(),
        IP:        ip,
        UserAgent: userAgent,
        Time:      now,
    })
    
    return nil
}

// Activate 激活账户
func (u *User) Activate() error {
    if u.status != vo.UserStatusPending {
        return kernel.ErrDomainRuleViolated
    }
    
    u.status = vo.UserStatusActive
    
    u.RecordEvent(&event.UserActivated{
        UserID: u.ID().Value(),
        Time:   time.Now(),
    })
    
    return nil
}

// Lock 锁定账户
func (u *User) Lock() {
    u.status = vo.UserStatusLocked
    until := time.Now().Add(LockDuration)
    u.lockedUntil = &until
    
    u.RecordEvent(&event.UserLocked{
        UserID:    u.ID().Value(),
        Until:     until,
        Reason:    "multiple_failed_attempts",
        Timestamp: time.Now(),
    })
}

// Unlock 解锁账户
func (u *User) Unlock() error {
    if u.status != vo.UserStatusLocked {
        return kernel.ErrDomainRuleViolated
    }
    
    u.status = vo.UserStatusActive
    u.failedLoginAttempts = 0
    u.lockedUntil = nil
    
    u.RecordEvent(&event.UserUnlocked{
        UserID: u.ID().Value(),
        Time:   time.Now(),
    })
    
    return nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string) error {
    // 验证旧密码
    if !u.password.Verify(oldPassword) {
        return kernel.NewBusinessError(
            CodeInvalidCredentials,
            "原密码错误",
        )
    }
    
    // 创建新密码
    newPwd, err := vo.NewPassword(newPassword)
    if err != nil {
        return err
    }
    
    u.password = newPwd
    
    u.RecordEvent(&event.UserPasswordChanged{
        UserID: u.ID().Value(),
        Time:   time.Now(),
    })
    
    return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
    email, err := vo.NewEmail(newEmail)
    if err != nil {
        return err
    }
    
    oldEmail := u.email
    u.email = email
    
    u.RecordEvent(&event.UserEmailChanged{
        UserID:   u.ID().Value(),
        OldEmail: oldEmail.String(),
        NewEmail: email.String(),
        Time:     time.Now(),
    })
    
    return nil
}

// Getters - 只读访问
func (u *User) Username() vo.Username { return u.username }
func (u *User) Email() vo.Email       { return u.email }
func (u *User) Status() vo.UserStatus { return u.status }
func (u *User) LastLoginAt() *time.Time { return u.lastLoginAt }

// 私有方法
func (u *User) isLocked() bool {
    if u.lockedUntil == nil {
        return false
    }
    
    // 检查锁是否已过期
    if time.Now().After(*u.lockedUntil) {
        u.Unlock()
        return false
    }
    
    return true
}
```

---

## 📢 领域事件（Domain Event）

### 什么是领域事件？

**领域事件**是领域中发生的重要事情的记录。

**特点：**
- 过去时态（已发生）
- 不可变
- 携带相关数据
- 用于解耦

### Domain Event 基类

```go
// domain/shared/kernel/domain_event.go
package kernel

import "time"

// DomainEvent 领域事件接口
type DomainEvent interface {
    Type() string           // 事件类型
    AggregateID() int64     // 聚合根 ID
    AggregateType() string  // 聚合根类型
    Timestamp() time.Time   // 发生时间
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
    aggregateID   int64
    aggregateType string
    timestamp     time.Time
}

// AggregateID 获取聚合根 ID
func (e *BaseDomainEvent) AggregateID() int64 {
    return e.aggregateID
}

// AggregateType 获取聚合根类型
func (e *BaseDomainEvent) AggregateType() string {
    return e.aggregateType
}

// Timestamp 获取时间戳
func (e *BaseDomainEvent) Timestamp() time.Time {
    return e.timestamp
}
```

### 领域事件示例

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

#### UserLocked

```go
// domain/user/event/user_locked.go
package event

import "time"

// UserLocked 用户已锁定
type UserLocked struct {
    kernel.BaseDomainEvent
    UserID    int64
    Until     time.Time
    Reason    string
    Timestamp time.Time
}

func (e *UserLocked) Type() string {
    return "user.locked"
}
```

---

## 🔄 聚合根与领域事件

### 事件记录机制

```go
// domain/shared/kernel/aggregate_root.go
package kernel

// AggregateRoot 聚合根接口
type AggregateRoot interface {
    Entity
    RecordEvent(event DomainEvent)
    ReleaseEvents() []DomainEvent
    ClearEvents()
}

// aggregateRootMixin 聚合根混入
type aggregateRootMixin struct {
    events []DomainEvent
}

// RecordEvent 记录领域事件
func (a *aggregateRootMixin) RecordEvent(event DomainEvent) {
    a.events = append(a.events, event)
}

// ReleaseEvents 释放所有事件
func (a *aggregateRootMixin) ReleaseEvents() []DomainEvent {
    events := a.events
    a.events = nil
    return events
}

// ClearEvents 清空事件
func (a *aggregateRootMixin) ClearEvents() {
    a.events = nil
}
```

### User 聚合根使用示例

```go
// 创建用户
user, err := aggregate.NewUser("john", "john@example.com", "Password123!")
if err != nil {
    return err
}

// 此时会记录 UserRegistered 事件
events := user.ReleaseEvents()
// events[0] => UserRegistered

// 用户登录
err = user.Login("Password123!", "192.168.1.1", "Mozilla/5.0...")
if err != nil {
    return err
}

// 如果登录成功，会记录 UserLoggedIn 事件
events = user.ReleaseEvents()
// events[0] => UserLoggedIn

// 保存聚合根和事件
err = userRepository.Save(ctx, user)
// Repository 会自动保存释放的事件到 domain_events 表
```

---

## 📊 领域模型关系图

```
┌─────────────────────────────────────┐
│          User (Aggregate Root)      │
│                                     │
│  - id: ID                           │
│  - username: Username (VO)          │
│  - email: Email (VO)                │
│  - password: Password (VO)          │
│  - status: UserStatus (VO)          │
│                                     │
│  + Login()                          │
│  + Activate()                       │
│  + Lock()                           │
│  + ChangePassword()                 │
│  + UpdateEmail()                    │
└─────────────────────────────────────┘
           ↓ (依赖)
┌─────────────────────────────────────┐
│      Value Objects (值对象)         │
├─────────────────────────────────────┤
│  • Username                         │
│  • Email                            │
│  • Password                         │
│  • UserStatus                       │
└─────────────────────────────────────┘
           ↓ (发布)
┌─────────────────────────────────────┐
│      Domain Events (领域事件)       │
├─────────────────────────────────────┤
│  • UserRegistered                   │
│  • UserLoggedIn                     │
│  • UserLocked                       │
│  • UserActivated                    │
│  • UserPasswordChanged              │
│  • UserEmailChanged                 │
└─────────────────────────────────────┘
```

---

## ✅ 最佳实践

### 1. 值对象应该不可变

```go
// ✅ 正确：值对象不可变
type Email struct {
    value string  // 私有字段
}

func (e Email) String() string {  // 只读方法
    return e.value
}

// ❌ 错误：值对象可变
type Email struct {
    value string
}

func (e *Email) SetValue(v string) {  // ❌ 不应该有 setter
    e.value = v
}
```

### 2. 聚合根应该封装业务逻辑

```go
// ✅ 正确：聚合根封装业务逻辑
func (u *User) Login(password string) error {
    if !u.password.Verify(password) {
        u.failedLoginAttempts++
        // ... 业务逻辑
    }
}

// ❌ 错误：贫血模型
type User struct {
    Password string
}
// 在应用层验证密码
if user.Password != inputPassword { ... }
```

### 3. 使用工厂方法创建值对象

```go
// ✅ 正确：工厂方法
email, err := vo.NewEmail("test@example.com")
if err != nil {
    // 处理验证错误
}

// ❌ 错误：直接构造
email := vo.Email{value: "test@example.com"}  // 绕过了验证
```

### 4. 聚合根通过事件与外部通信

```go
// ✅ 正确：通过事件
u.RecordEvent(&event.UserLoggedIn{...})

// ❌ 错误：直接调用外部服务
func (u *User) Login(password string) error {
    // 不应该在聚合根内直接发送邮件
    emailService.SendWelcomeEmail(u.email)  
}
```

---

## 📚 参考资源

- [Domain-Driven Design](https://martinfowler.com/books/domainDrivenDesign.html)
- [Implementing Domain-Driven Design](https://msdn.microsoft.com/en-us/library/jj973677.aspx)
- [值对象模式](https://martinfowler.com/bliki/ValueObject.html)
- [聚合模式](https://martinfowler.com/bliki/DDD_Aggregate.html)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team

