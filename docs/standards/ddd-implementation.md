# DDD 实现规范

## 1. 分层架构职责

### 架构总览

```
┌─────────────────────────────────────┐
│         Interfaces Layer            │  ← HTTP/gRPC/CLI 适配器
├─────────────────────────────────────┤
│       Application Layer             │  ← 应用编排，不含业务逻辑
├─────────────────────────────────────┤
│         Domain Layer                │  ← 纯业务逻辑（核心）
├─────────────────────────────────────┤
│      Infrastructure Layer           │  ← 技术实现（数据库、缓存等）
└─────────────────────────────────────┘
```

**依赖规则**: 
- 上层可以依赖下层（Interfaces → Application → Domain）
- Infrastructure 依赖 Domain（依赖倒置）
- **禁止跨层调用**（如 Interfaces 直接调用 Domain）

---

## 2. Domain Layer（领域层）

### 2.1 职责定义

**只包含**:
- ✅ 实体（Entity）和聚合根（Aggregate Root）
- ✅ 值对象（Value Object）
- ✅ 领域服务（Domain Service）
- ✅ 领域事件（Domain Event）
- ✅ 仓储接口（Repository Interface）

**禁止包含**:
- ❌ 基础设施代码（数据库、HTTP、缓存）
- ❌ DTO 转换逻辑
- ❌ 应用编排逻辑

---

### 2.2 实体设计

```go
package user

// User 用户聚合根
type User struct {
    ID        uuid.UUID
    Email     Email              // 值对象
    Password  HashedPassword     // 值对象
    Nickname  Nickname           // 值对象
    Avatar    *string
    Phone     *string
    Bio       *string
    Status    UserStatus
    CreatedAt time.Time
    UpdatedAt time.Time
    
    // 领域事件（临时存储，发布后清空）
    events []DomainEvent
}

// 业务方法：封装业务逻辑
func (u *User) UpdateEmail(email Email) error {
    if u.Email.Equals(email) {
        return nil
    }
    
    u.Email = email
    u.addEvent(UserEmailChangedEvent{
        UserID:  u.ID,
        NewEmail: email.String(),
    })
    return nil
}

// Events 获取并清除已发生的事件
func (u *User) Events() []DomainEvent {
    events := u.events
    u.events = nil
    return events
}

func (u *User) addEvent(event DomainEvent) {
    u.events = append(u.events, event)
}
```

**关键规则**:
1. 实体必须有明确的业务方法（而非 getter/setter）
2. 状态变更必须通过方法进行
3. 重要的状态变更要发布领域事件

---

### 2.3 值对象设计

```go
package valueobject

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 构造函数（带验证）
func NewEmail(email string) (Email, error) {
    e := Email{value: strings.TrimSpace(strings.ToLower(email))}
    if !e.IsValid() {
        return Email{}, ErrInvalidEmail
    }
    return e, nil
}

// IsValid 验证邮箱格式
func (e Email) IsValid() bool {
    return emailRegex.MatchString(e.value)
}

// String 获取字符串表示
func (e Email) String() string {
    return e.value
}

// Equals 判断相等性
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

**关键规则**:
1. 值对象必须不可变（只有 getter，没有 setter）
2. 构造函数必须进行验证
3. 实现 `Equals` 方法用于比较
4. 实现 `String()` 方法用于显示

---

### 2.4 领域服务

```go
package service

// TenantMemberService 租户成员领域服务
// 当操作涉及多个聚合根时使用领域服务
type TenantMemberService struct {
    tenantRepo  TenantRepository
    memberRepo  TenantMemberRepository
    userRepo    UserRepository
}

// NewTenantMemberService 构造函数
func NewTenantMemberService(
    tenantRepo TenantRepository,
    memberRepo TenantMemberRepository,
    userRepo UserRepository,
) *TenantMemberService {
    return &TenantMemberService{
        tenantRepo: tenantRepo,
        memberRepo: memberRepo,
        userRepo:   userRepo,
    }
}

// AddMember 添加租户成员（涉及 Tenant 和 User 两个聚合根）
func (s *TenantMemberService) AddMember(
    ctx context.Context,
    tenantID uuid.UUID,
    userID uuid.UUID,
    role entity.UserRole,
) (*entity.TenantMember, error) {
    // 1. 获取租户聚合根
    tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
    if err != nil {
        return nil, err
    }
    
    // 2. 获取用户聚合根
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 3. 检查租户限制
    members, _ := s.memberRepo.ListByTenant(ctx, tenantID)
    if len(members) >= tenant.MaxMembers {
        return nil, ErrTenantLimitExceed
    }
    
    // 4. 创建成员关系
    member := entity.NewTenantMember(tenantID, userID, role)
    
    // 5. 持久化
    if err := s.memberRepo.Create(ctx, member); err != nil {
        return nil, err
    }
    
    return member, nil
}
```

**何时使用领域服务**:
1. 操作涉及多个聚合根
2. 需要访问外部资源（通过 Repository 接口）
3. 无状态的领域逻辑

---

### 2.5 领域事件

```go
package event

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
    UserID    uuid.UUID
    Email     string
    Role      string
    TenantID  *uuid.UUID
    EventID   string
    EventType string
    AggregateID uuid.UUID
    OccurredAt  time.Time
    Version     int
}

// NewUserRegisteredEvent 构造函数
func NewUserRegisteredEvent(
    userID uuid.UUID,
    email string,
    role string,
    tenantID *uuid.UUID,
) *UserRegisteredEvent {
    return &UserRegisteredEvent{
        UserID:      userID,
        Email:       email,
        Role:        role,
        TenantID:    tenantID,
        EventID:     uuid.New().String(),
        EventType:   "UserRegistered",
        AggregateID: userID,
        OccurredAt:  time.Now(),
        Version:     1,
    }
}

// 实现 DomainEvent 接口
func (e *UserRegisteredEvent) GetEventType() string   { return e.EventType }
func (e *UserRegisteredEvent) GetEventID() string     { return e.EventID }
func (e *UserRegisteredEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserRegisteredEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e *UserRegisteredEvent) GetVersion() int        { return e.Version }
```

**事件命名规范**:
- ✅ 使用过去式：`UserRegistered`, `EmailChanged`
- ✅ 反映业务事实：`PaymentCompleted` 而非 `UpdateOrder`
- ❌ 避免命令式：不要用 `RegisterUser`, `ChangeEmail`

---

## 3. Application Layer（应用层）

### 3.1 职责定义

**只包含**:
- ✅ 应用服务（协调领域对象完成任务）
- ✅ DTO（数据传输对象）
- ✅ Assembler（DTO ↔ Entity 转换器）
- ✅ 命令查询接口（CQRS）

**禁止包含**:
- ❌ 业务逻辑判断（应该在 Domain）
- ❌ 基础设施代码（应该在 Infrastructure）

---

### 3.2 应用服务

```go
package service

// UserService 用户应用服务
type UserService struct {
    userRepo     UserRepository
    assembler    *UserAssembler
    eventBus     EventBus
}

// NewUserService 构造函数
func NewUserService(
    userRepo UserRepository,
    eventBus EventBus,
) *UserService {
    return &UserService{
        userRepo:  userRepo,
        assembler: NewUserAssembler(),
        eventBus:  eventBus,
    }
}

// ChangeEmail 修改邮箱（应用服务方法）
func (s *UserService) ChangeEmail(
    ctx context.Context,
    userID uuid.UUID,
    newEmail string,
) error {
    // 1. 获取聚合根
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. 调用领域方法（业务逻辑在 Domain）
    emailVO, err := valueobject.NewEmail(newEmail)
    if err != nil {
        return err
    }
    
    if err := user.UpdateEmail(emailVO); err != nil {
        return err
    }
    
    // 3. 持久化
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // 4. 发布事件
    for _, event := range user.Events() {
        s.eventBus.Publish(ctx, event)
    }
    
    return nil
}
```

**应用服务 vs 领域服务**:
| 特性 | 应用服务 | 领域服务 |
|------|---------|---------|
| 职责 | 应用编排 | 领域逻辑 |
| 事务 | 管理事务边界 | 不参与事务 |
| 依赖 | 依赖 Repository 接口 | 依赖 Repository 接口 |
| 返回 | DTO 或无返回值 | 领域对象 |

---

### 3.3 DTO 设计

```go
package dto

// User 用户 DTO（扁平结构）
type User struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Nickname  string     `json:"nickname"`
    Phone     *string    `json:"phone,omitempty"`
    Bio       *string    `json:"bio,omitempty"`
    Avatar    *string    `json:"avatar,omitempty"`
    Status    string     `json:"status"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
}

// CreateUserRequest 创建用户请求 DTO
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Nickname string `json:"nickname" validate:"required,min=2,max=20"`
}

// UpdateProfileRequest 更新个人资料请求 DTO
type UpdateProfileRequest struct {
    Nickname *string `json:"nickname,omitempty" validate:"omitempty,min=2,max=20"`
    Phone    *string `json:"phone,omitempty" validate:"omitempty,e164"`
    Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}
```

**DTO 设计原则**:
1. 扁平结构（避免嵌套过深）
2. 使用基本类型（string, int, bool）
3. 不使用领域对象（Entity, ValueObject）
4. 带 JSON 标签和验证标签

---

### 3.4 Assembler 模式

```go
package assembler

// UserAssembler 用户 DTO 转换器
type UserAssembler struct{}

// NewUserAssembler 构造函数
func NewUserAssembler() *UserAssembler {
    return &UserAssembler{}
}

// ToDTO 领域实体 → DTO
func (a *UserAssembler) ToDTO(user *entity.User) *dto.User {
    if user == nil {
        return nil
    }
    
    return &dto.User{
        ID:        user.ID.String(),
        Email:     user.Email.String(),
        Nickname:  user.Nickname.String(),
        Phone:     user.Phone,
        Bio:       user.Bio,
        Avatar:    user.Avatar,
        Status:    string(user.Status),
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}

// FromCreateRequest 创建请求 → 领域实体
func (a *UserAssembler) FromCreateRequest(
    req *dto.CreateUserRequest,
    hashedPassword entity.HashedPassword,
) (*entity.User, error) {
    email, err := valueobject.NewEmail(req.Email)
    if err != nil {
        return nil, err
    }
    
    nickname, err := valueobject.NewNickname(req.Nickname)
    if err != nil {
        return nil, err
    }
    
    return &entity.User{
        Email:    email,
        Password: hashedPassword,
        Nickname: nickname,
        Status:   entity.StatusActive,
    }, nil
}
```

---

## 4. Infrastructure Layer（基础设施层）

### 4.1 职责定义

**只包含**:
- ✅ 仓储实现（Repository Implementation）
- ✅ 数据模型（Data Model）
- ✅ DAO（数据访问对象）
- ✅ 第三方服务集成（Redis, Kafka, etc.）

**禁止包含**:
- ❌ 业务逻辑
- ❌ 直接返回 Model 给应用层

---

### 4.2 仓储实现

```go
package repo

// userRepository 用户仓储实现
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository 构造函数
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    // 1. 查询数据库模型
    var model model.User
    if err := r.db.First(&model, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errPkg.ErrUserNotFound
        }
        return nil, err
    }
    
    // 2. Model → Entity 转换
    return r.toEntity(&model), nil
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    // 1. Entity → Model 转换
    model := r.toModel(user)
    
    // 2. 保存到数据库
    return r.db.Create(model).Error
}

// toEntity Model → Entity 转换
func (r *userRepository) toEntity(m *model.User) *entity.User {
    email, _ := valueobject.NewEmail(m.Email)
    nickname, _ := valueobject.NewNickname(m.Nickname)
    
    return &entity.User{
        ID:        uuid.MustParse(*m.ID),
        Email:     email,
        Password:  entity.HashedPassword(m.Password),
        Nickname:  nickname,
        Avatar:    m.Avatar,
        Phone:     m.Phone,
        Bio:       m.Bio,
        Status:    entity.UserStatus(m.Status),
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

// toModel Entity → Model 转换
func (r *userRepository) toModel(u *entity.User) *model.User {
    id := u.ID.String()
    return &model.User{
        ID:       &id,
        Email:    u.Email.String(),
        Password: u.Password.String(),
        Nickname: u.Nickname.String(),
        Avatar:   u.Avatar,
        Phone:    u.Phone,
        Bio:      u.Bio,
        Status:   string(u.Status),
    }
}
```

**关键规则**:
1. 仓储接口在 Domain 层定义，实现在 Infrastructure 层
2. Model ↔ Entity 转换必须在仓储内部完成
3. 不能直接返回 Model 给应用层

---

### 4.3 数据模型

```go
package model

// User 用户数据模型（GORM Model）
type User struct {
    ID        *string    `gorm:"type:uuid;primaryKey" json:"id"`
    Email     string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
    Password  string     `gorm:"size:255;not null" json:"-"`
    Nickname  string     `gorm:"size:100;not null" json:"nickname"`
    Avatar    *string    `gorm:"size:500" json:"avatar,omitempty"`
    Phone     *string    `gorm:"size:20" json:"phone,omitempty"`
    Bio       *string    `gorm:"size:500" json:"bio,omitempty"`
    Status    string     `gorm:"size:20;not null;default:'active'" json:"status"`
    CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

**数据模型规范**:
1. 使用指针字段表示可选值（便于处理 NULL）
2. 必须指定表名、字段类型、长度
3. 敏感字段加 `json:"-"` 标签（防止泄露）

---

## 5. Interfaces Layer（接口层）

### 5.1 HTTP Handler

```go
package http

// UserHandler 用户 HTTP Handler
type UserHandler struct {
    userService application.UserService
    logger      *zap.Logger
}

// NewUserHandler 构造函数
func NewUserHandler(
    userService application.UserService,
    logger *zap.Logger,
) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户 ID 获取详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} dto.User
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // 1. 绑定参数
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(400, response.BadRequest("无效的用户 ID"))
        return
    }
    
    // 2. 调用应用服务
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.logger.Error("获取用户失败", zap.Error(err))
        c.JSON(500, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(200, response.Success(user))
}
```

**Handler 职责**:
1. 参数绑定和验证
2. 调用应用服务
3. 格式化响应
4. 记录日志

**禁止**:
- ❌ 业务逻辑
- ❌ 直接访问仓储
- ❌ 返回领域实体

---

## 6. 依赖注入规范

### 6.1 Wire 配置

```go
// providers.go
package wire

import (
    "github.com/google/wire"
    "gorm.io/gorm"
    "go-ddd-scaffold/internal/domain/user/repository"
    "go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
)

// UserProviderSet 用户模块依赖集合
var UserProviderSet = wire.NewSet(
    NewUserRepository,
    NewUserService,
    NewUserHandler,
)

// NewUserRepository 提供用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return repo.NewUserRepository(db)
}

// NewUserService 提供用户服务
func NewUserService(userRepo repository.UserRepository) *service.UserService {
    return service.NewUserService(userRepo)
}

// NewUserHandler 提供用户 Handler
func NewUserHandler(svc *service.UserService) *http.UserHandler {
    return http.NewUserHandler(svc)
}
```

---

## 7. 测试规范

### 7.1 领域层测试

```go
package user_test

func TestUser_UpdateEmail(t *testing.T) {
    t.Run("有效邮箱_更新成功", func(t *testing.T) {
        // Given
        user := createTestUser()
        newEmail, _ := valueobject.NewEmail("new@example.com")
        
        // When
        err := user.UpdateEmail(newEmail)
        
        // Assert
        assert.NoError(t, err)
        assert.Equal(t, newEmail, user.Email)
        
        // 验证事件发布
        events := user.Events()
        require.Len(t, events, 1)
        assert.IsType(t, &event.UserEmailChangedEvent{}, events[0])
    })
    
    t.Run("无效邮箱_返回错误", func(t *testing.T) {
        // Given
        user := createTestUser()
        invalidEmail, _ := valueobject.NewEmail("invalid")
        
        // When
        err := user.UpdateEmail(invalidEmail)
        
        // Assert
        assert.Error(t, err)
    })
}
```

---

## 8. 检查清单

在提交领域代码前，请确认：

### Domain Layer
- [ ] 是否只有纯业务逻辑？
- [ ] 是否有基础设施依赖？（不应该有）
- [ ] 实体是否有业务方法？（而非 getter/setter）
- [ ] 值对象是否不可变？
- [ ] 重要的状态变更是否发布事件？

### Application Layer
- [ ] 是否只负责编排？
- [ ] 是否有业务判断逻辑？（应该在 Domain）
- [ ] 是否返回 DTO？（而非 Entity）

### Infrastructure Layer
- [ ] 是否实现了 Repository 接口？
- [ ] Model ↔ Entity 转换是否在内部完成？
- [ ] 是否直接返回 Model？（不应该）

### Interfaces Layer
- [ ] 是否只做协议转换？
- [ ] 是否有业务逻辑？（不应该有）
- [ ] 参数是否验证？

---

**版本**: v1.0  
**生效日期**: 2026-03-06
