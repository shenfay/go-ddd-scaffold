# DDD 设计指南

本文档提供领域驱动设计（DDD）在 Go DDD Scaffold 中的实践指南。

## 📋 DDD 核心概念

### 战略设计（Strategic Design）

```
┌─────────────────────────────────────┐
│     战略设计                        │
│                                     │
│  • 限界上下文（Bounded Context）    │
│  • 通用语言（Ubiquitous Language）  │
│  • 上下文映射（Context Mapping）    │
└─────────────────────────────────────┘
```

### 战术设计（Tactical Design）

```
┌─────────────────────────────────────┐
│     战术设计                        │
│                                     │
│  • 实体（Entity）                   │
│  • 值对象（Value Object）           │
│  • 聚合根（Aggregate Root）         │
│  • 领域服务（Domain Service）       │
│  • 领域事件（Domain Event）         │
│  • 仓储（Repository）               │
└─────────────────────────────────────┘
```

---

## 🎯 限界上下文（Bounded Context）

### 什么是限界上下文？

**限界上下文**是领域模型的边界，定义了通用语言的适用范围。

### 项目中的限界上下文

```
Go DDD Scaffold
├── Auth Context（认证上下文）
│   ├── 职责：用户认证、令牌管理
│   ├── 核心概念：Token, JWT, Credential
│   └── 模块：auth
│
├── User Context（用户上下文）
│   ├── 职责：用户管理、资料维护
│   ├── 核心概念：User, Profile, Status
│   └── 模块：user
│
├── Tenant Context（租户上下文）
│   ├── 职责：租户管理、成员管理
│   ├── 核心概念：Tenant, Member, Role
│   └── 模块：tenant
│
└── RBAC Context（权限上下文）
    ├── 职责：角色权限管理
    ├── 核心概念：Role, Permission
    └── 模块：rbac
```

### 上下文映射关系

```
┌──────────────┐
│  Auth Context│
└──────┬───────┘
       │ 使用
       ↓
┌──────────────┐
│  User Context│
└──────┬───────┘
       │ 拥有
       ↓
┌──────────────┐
│ Tenant Context│
└──────┬───────┘
       │ 包含
       ↓
┌──────────────┐
│  RBAC Context │
└──────────────┘
```

---

## 💬 通用语言（Ubiquitous Language）

### 什么是通用语言？

**通用语言**是开发人员和领域专家共同使用的语言，体现在代码命名中。

### 项目中的通用语言

| 术语 | 含义 | 代码体现 |
|------|------|----------|
| **Aggregate（聚合根）** | 业务对象的集合 | `type User struct` |
| **Value Object（值对象）** | 不可变的度量概念 | `type Email struct` |
| **Entity（实体）** | 有唯一标识的对象 | `type Entity struct` |
| **Domain Event（领域事件）** | 领域中发生的事 | `type UserRegistered struct` |
| **Repository（仓储）** | 持久化抽象层 | `type UserRepository interface` |
| **Domain Service（领域服务）** | 跨聚合的业务逻辑 | `type PasswordHasher interface` |

### 命名一致性

```go
// ✅ 正确：使用通用语言
type UserRepository interface { ... }
type UserAggregate struct { ... }
type Email valueobject.Email

// ❌ 错误：混用术语
type UserDataAccess { ... }  // 不是 Repository
type UserDO struct { ... }    // 不是 Aggregate
type Mail string              // 不是 Value Object
```

---

## 🌳 聚合设计原则

### 1. 单一职责原则

```go
// ✅ 正确：User 聚合只负责用户相关业务
type User struct {
    *kernel.Entity
    username vo.Username
    email    vo.Email
    password vo.Password
    status   vo.UserStatus
}

func (u *User) Login(password string) error { ... }
func (u *User) Activate() error { ... }
func (u *User) ChangePassword(oldPwd, newPwd string) error { ... }

// ❌ 错误：User 聚合包含租户逻辑
type User struct {
    // ... 用户字段
    tenants []*Tenant  // ❌ 不应该在这里
}

func (u *User) CreateTenant(name string) (*Tenant, error) { ... }  // ❌ 职责混乱
```

### 2. 封装业务逻辑

```go
// ✅ 正确：封装业务规则
func (u *User) Login(password string, ip string) error {
    // 检查锁定
    if u.isLocked() {
        return ErrUserLocked
    }
    
    // 验证密码
    if !u.password.Verify(password) {
        u.failedLoginAttempts++
        if u.failedLoginAttempts >= MaxLoginAttempts {
            u.Lock()
            return ErrUserLocked
        }
        return ErrInvalidCredentials
    }
    
    // 登录成功
    u.lastLoginAt = time.Now()
    u.lastLoginIP = ip
    return nil
}

// ❌ 错误：业务逻辑泄露到应用层
// application/auth/service.go
func (s *AuthService) Login(password string) error {
    if user.FailedAttempts >= 5 {  // ❌ 业务逻辑泄露
        return ErrLocked
    }
    if user.Password != password {  // ❌ 直接比较密码
        user.FailedAttempts++
    }
}
```

### 3. 通过值对象封装验证

```go
// ✅ 正确：验证逻辑封装在值对象中
username, err := vo.NewUsername("john")  // 内部验证
email, err := vo.NewEmail("test@example.com")  // 内部验证
password, err := vo.NewPassword("weak")  // 内部验证强度

// ❌ 错误：验证逻辑分散
if len(username) < 3 { ... }
if !regexp.MatchString(emailPattern, email) { ... }
if len(password) < 8 { ... }
```

### 4. 发布领域事件

```go
// ✅ 正确：记录重要业务操作
func (u *User) Activate() error {
    u.status = vo.UserStatusActive
    
    // 发布事件
    u.RecordEvent(&event.UserActivated{
        UserID: u.ID().Value(),
        Time:   time.Now(),
    })
    
    return nil
}

// ❌ 错误：无事件记录
func (u *User) Activate() error {
    u.status = vo.UserStatusActive
    // 没有发布事件，其他模块不知道用户已激活
}
```

---

## 🏗️ 分层架构中的 DDD

### 完整流程示例：用户注册

```
┌─────────────────────────────────────────┐
│  Interfaces Layer (HTTP Handler)        │
│                                         │
│  POST /api/users                        │
│  Body: {username, email, password}      │
│                                         │
│  handler.CreateUser(c)                  │
│    ↓                                    │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│  Application Layer (Service)            │
│                                         │
│  func (s *UserService) CreateUser(cmd)  │
│  1. 验证命令                            │
│  2. 调用领域方法                        │
│  3. 保存聚合                            │
│  4. 发布事件                            │
│                                         │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│  Domain Layer (Aggregate)               │
│                                         │
│  aggregate.NewUser(...)                 │
│    - 创建值对象                         │
│    - 验证业务规则                       │
│    - 记录领域事件                       │
│                                         │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│  Infrastructure Layer (Repository)      │
│                                         │
│  repo.Save(user)                        │
│    - 保存到数据库                       │
│    - 保存领域事件                       │
│    - 事务管理                           │
│                                         │
└─────────────────────────────────────────┘
```

### 代码实现

#### 1. HTTP Handler（Interfaces 层）

```go
// interfaces/http/user/handler.go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, http.StatusBadRequest, err)
        return
    }
    
    cmd := &app_user.CreateUserCommand{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }
    
    result, err := h.userService.CreateUser(c.Request.Context(), cmd)
    if err != nil {
        h.respHandler.Error(c, http.StatusInternalServerError, err)
        return
    }
    
    h.respHandler.Success(c, result)
}
```

#### 2. Application Service（应用层）

```go
// application/user/service.go
func (s *UserService) CreateUser(ctx context.Context, cmd *CreateUserCommand) (*UserResponse, error) {
    // 1. 创建领域对象（调用领域方法）
    user, err := aggregate.NewUser(cmd.Username, cmd.Email, cmd.Password)
    if err != nil {
        return nil, err
    }
    
    // 2. 保存聚合（使用 Repository Port）
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("save user failed: %w", err)
    }
    
    // 3. 发布领域事件（自动由 Repository 完成）
    // events := user.ReleaseEvents()
    // s.eventPublisher.Publish(events)
    
    // 4. 返回响应
    return &UserResponse{
        ID:       user.ID().Value(),
        Username: user.Username().String(),
        Email:    user.Email().String(),
    }, nil
}
```

#### 3. Domain Aggregate（领域层）

```go
// domain/user/aggregate/user.go
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
        Entity:   kernel.NewEntity(),
        username: uName,
        email:    uEmail,
        password: uPassword,
        status:   vo.UserStatusPending,
    }
    
    // 记录领域事件
    user.RecordEvent(&event.UserRegistered{
        UserID:  user.ID().Value(),
        Email:   user.email.String(),
        Created: time.Now(),
    })
    
    return user, nil
}
```

#### 4. Repository Implementation（基础设施层）

```go
// infrastructure/persistence/repository/user_repository.go
func (r *userRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 1. 保存聚合根
    err := r.saveUser(tx, user)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 2. 保存领域事件（Outbox Pattern）
    events := user.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}
```

---

## 🔄 领域事件驱动架构

### 事件发布订阅模式

```
┌─────────────────────────────────────────┐
│  User Registered (Event)                │
│                                         │
│  Publisher: UserAggregate               │
│  Subscribers:                           │
│    1. EmailService → 发送欢迎邮件       │
│    2. TenantService → 初始化租户配置    │
│    3. AnalyticsService → 记录统计       │
└─────────────────────────────────────────┘
```

### 实现方式

```go
// 1. 定义事件
type UserRegistered struct {
    UserID  int64
    Email   string
    Created time.Time
}

// 2. 聚合根发布事件
func (u *User) RecordEvent(event DomainEvent) {
    u.events = append(u.events, event)
}

// 3. Repository 保存事件
func (r *userRepositoryImpl) Save(user *User) error {
    // 保存用户...
    
    // 保存事件到 domain_events 表
    events := user.ReleaseEvents()
    for _, event := range events {
        r.saveEvent(tx, event)
    }
    
    return nil
}

// 4. 异步处理器发布到消息队列
func (w *DomainEventHandler) Handle() {
    events := r.getUnprocessedEvents()
    for _, event := range events {
        // 发布到 Redis Stream / Kafka
        publisher.Publish(event)
    }
}

// 5. 订阅者处理事件
func (s *EmailService) OnUserRegistered(event UserRegistered) {
    s.SendWelcomeEmail(event.Email)
}

func (t *TenantService) OnUserRegistered(event UserRegistered) {
    t.InitializeTenant(event.UserID)
}
```

---

## 📊 CQRS 模式（可选演进）

### 当前架构（简单 CQRS）

```
写模型（Command）
└── 聚合根 + Repository
    └── 处理创建、更新、删除

读模型（Query）
└── DAO 直接查询
    └── 处理列表、详情查询
```

### 未来演进（完整 CQRS）

```
┌─────────────────────────────────────┐
│          Command Side               │
│                                     │
│  UserAggregate                      │
│  UserRepository                     │
│  UserService (Write)                │
└──────────────┬──────────────────────┘
               │
               │ 发布事件
               ↓
┌─────────────────────────────────────┐
│          Query Side                 │
│                                     │
│  ReadModel (扁平化数据)             │
│  UserQueryDAO (直接查询)            │
│  UserQueryService (Read)            │
└─────────────────────────────────────┘
```

---

## ✅ DDD 实施检查清单

### 战略设计

- [ ] 识别了限界上下文
- [ ] 建立了通用语言
- [ ] 明确了上下文映射关系

### 战术设计

- [ ] 识别了聚合根
- [ ] 创建了值对象
- [ ] 封装了业务逻辑
- [ ] 定义了领域服务
- [ ] 发布了领域事件
- [ ] 实现了 Repository 模式

### 代码质量

- [ ] 聚合根有丰富的业务方法
- [ ] 值对象是不可变的
- [ ] 实体有唯一标识
- [ ] 领域事件使用过去时态
- [ ] Repository 接口在 Domain 层
- [ ] 应用服务编排业务流程

---

## 📚 参考资源

- [Domain-Driven Design](https://martinfowler.com/books/domainDrivenDesign.html) - Eric Evans
- [Implementing Domain-Driven Design](https://msdn.microsoft.com/en-us/library/jj973677.aspx)
- [DDD Quick Reference](https://www.infoq.com/minibooks/domain-driven-design-quick-reference)
- [微服务架构中的 DDD](https://azure.microsoft.com/zh-cn/resources/cloud-computing-dictionary/what-is-domain-driven-design/)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
