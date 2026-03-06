# 分层架构详解

## 📊 架构总览

### 四层架构图

```
┌─────────────────────────────────────────┐
│         Interfaces Layer                │  ← HTTP/gRPC/CLI 适配器
│  - Handler (协议转换)                    │
│  - Middleware (横切关注点)               │
│  - Router (路由注册)                     │
├─────────────────────────────────────────┤
│       Application Layer                 │  ← 应用编排，不含业务逻辑
│  - Service (协调 Domain + Infra)         │
│  - DTO (数据传输对象)                    │
│  - Assembler (DTO ↔ Entity)             │
├─────────────────────────────────────────┤
│         Domain Layer                    │  ← 纯业务逻辑（核心）
│  - Entity (实体 + 聚合根)                │
│  - ValueObject (值对象)                  │
│  - Repository Interface (仓储接口)       │
│  - Domain Service (领域服务)             │
│  - Domain Event (领域事件)               │
├─────────────────────────────────────────┤
│      Infrastructure Layer               │  ← 技术实现
│  - Repository Implementation            │
│  - Model (数据模型)                      │
│  - DAO (数据访问对象)                    │
│  - Redis/DB/Queue 客户端                 │
└─────────────────────────────────────────┘
```

**依赖规则**:
- ✅ 上层可以依赖下层（Interfaces → Application → Domain）
- ✅ Infrastructure 依赖 Domain（依赖倒置）
- ❌ **禁止跨层调用**（如 Interfaces 直接调用 Domain）
- ❌ **禁止反向依赖**（如 Domain 依赖 Infrastructure）

---

## 1. Domain Layer（领域层）

### 职责定义

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

### 1.1 实体设计

#### ✅ 正确示例

```go
// internal/domain/user/entity/user.go

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
    
    events []DomainEvent  // 领域事件
}

// 业务方法：封装业务逻辑
func (u *User) UpdateProfile(nickname Nickname, phone *string, bio *string) error {
    u.Nickname = nickname
    u.Phone = phone
    u.Bio = bio
    u.addEvent(UserProfileUpdatedEvent{
        UserID:   u.ID,
        Nickname: nickname.String(),
    })
    return nil
}

func (u *User) Lock() error {
    if u.Status == StatusLocked {
        return ErrAlreadyLocked
    }
    u.Status = StatusLocked
    u.addEvent(UserLockedEvent{UserID: u.ID})
    return nil
}

func (u *User) Activate() error {
    u.Status = StatusActive
    u.addEvent(UserActivatedEvent{UserID: u.ID})
    return nil
}
```

**优点**:
- ✅ 有丰富的业务方法（而非 getter/setter）
- ✅ 使用值对象封装（Email, Nickname）
- ✅ 重要的状态变更发布领域事件
- ✅ 无基础设施依赖

---

#### ❌ 错误示例（当前问题）

```go
// ❌ 问题 1: 缺少业务方法
type User struct {
    ID       uuid.UUID
    Email    string  // ❌ 直接使用字符串
    Password string  // ❌ 明文密码
}

func (u *User) SetEmail(email string) {  // ❌ Setter 模式
    u.Email = email
}

// ❌ 问题 2: 包含基础设施标签
type User struct {
    ID    uuid.UUID `gorm:"type:uuid"`  // ❌ GORM 标签
    Email string    `json:"email"`      // ❌ JSON 标签
}

// ❌ 问题 3: 依赖加密库
import "golang.org/x/crypto/bcrypt"  // ❌ 基础设施依赖

type HashedPassword string
func NewHashedPassword(password string) (HashedPassword, error) {
    return bcrypt.GenerateFromPassword(...)  // ❌ 在 Domain 层调用基础设施
}
```

**改进方案**: 见 [code-review-report.md#L28](code-review-report.md#L28)

---

### 1.2 值对象设计

#### ✅ 正确示例

```go
// internal/domain/user/valueobject/user_values.go

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

**优点**:
- ✅ 不可变（只有 getter，没有 setter）
- ✅ 构造函数强制验证
- ✅ 实现 Equals 方法
- ✅ 实现 String() 方法

---

#### ❌ 错误示例

```go
// ❌ 问题 1: 直接使用类型别名
type Email string  // ❌ 没有验证逻辑

// ❌ 问题 2: 可变值对象
type Email struct {
    Value string
}

func (e *Email) SetValue(v string) {  // ❌ Setter
    e.Value = v  // 没有验证
}

// ❌ 问题 3: 重复定义
// shared/valueobject/email.go
type Email struct{ value string }

// user/valueobject/user_values.go  
type Email string  // ❌ 两个不同的定义
```

**改进方案**: 统一到 `internal/domain/shared/valueobject/` 中定义

---

### 1.3 聚合根设计

#### ✅ 正确示例（Tenant 聚合根）

```go
// internal/domain/tenant/entity/tenant.go

// Tenant 租户聚合根
type Tenant struct {
    ID          uuid.UUID
    Name        TenantName     // 值对象
    Description string
    MaxMembers  int
    ExpiredAt   time.Time
    Status      TenantStatus
    members     []*TenantMember  // 内部管理实体
}

// AddMember 添加成员（聚合根方法）
func (t *Tenant) AddMember(userID uuid.UUID, role UserRole) (*TenantMember, error) {
    // 1. 检查限制
    if len(t.members) >= t.MaxMembers {
        return nil, ErrTenantLimitExceed
    }
    
    // 2. 检查是否已存在
    for _, m := range t.members {
        if m.UserID == userID {
            return nil, ErrMemberExists
        }
    }
    
    // 3. 创建成员关系
    member := NewTenantMember(t.ID, userID, role)
    t.members = append(t.members, member)
    
    // 4. 发布领域事件
    t.addEvent(TenantMemberAddedEvent{
        TenantID: t.ID,
        MemberID: member.ID,
        Role:     role,
    })
    
    return member, nil
}

// RemoveMember 移除成员
func (t *Tenant) RemoveMember(memberID uuid.UUID) error {
    for i, m := range t.members {
        if m.ID == memberID {
            t.members = append(t.members[:i], t.members[i+1:]...)
            t.addEvent(TenantMemberRemovedEvent{...})
            return nil
        }
    }
    return ErrMemberNotFound
}
```

**优点**:
- ✅ 聚合根管理内部实体的完整性
- ✅ 业务规则封装在聚合根内
- ✅ 发布领域事件记录重要变更

---

#### ❌ 错误示例（当前问题）

```go
// ❌ Tenant 只是数据容器，没有行为
type Tenant struct {
    ID         uuid.UUID
    Name       string
    MaxMembers int
    // ... 只有字段，没有方法
}

// ❌ 成员管理在外部进行（破坏聚合边界）
func CreateTenantMember(tenantID, userID uuid.UUID) *TenantMember {
    // 这里可以访问 Tenant 的内部成员列表吗？
    // 如何保证 MaxMembers 限制？
}
```

**改进方案**: 见 [code-review-report.md#L167](code-review-report.md#L167)

---

### 1.4 领域服务

#### ✅ 正确示例

```go
// internal/domain/user/service/user_registration_service.go

// UserRegistrationService 用户注册领域服务
type UserRegistrationService struct {
    userRepo UserRepository  // 依赖仓储接口
}

// NewUserRegistrationService 构造函数
func NewUserRegistrationService(userRepo UserRepository) *UserRegistrationService {
    return &UserRegistrationService{userRepo: userRepo}
}

// RegisterUser 注册用户（包含复杂业务逻辑）
func (s *UserRegistrationService) RegisterUser(
    ctx context.Context,
    email string,
    password string,
    nickname string,
) (*User, error) {
    // 1. 业务验证
    if err := s.validateEmailUnique(email); err != nil {
        return nil, err
    }
    
    if err := s.validatePasswordStrength(password); err != nil {
        return nil, err
    }
    
    // 2. 创建用户实体
    emailVO, _ := valueobject.NewEmail(email)
    hashedPwd, _ := entity.NewHashedPassword(password)
    nicknameVO, _ := valueobject.NewNickname(nickname)
    
    user := entity.NewUser(emailVO, hashedPwd, nicknameVO)
    
    // 3. 返回（由 Application Service 持久化）
    return user, nil
}

func (s *UserRegistrationService) validateEmailUnique(email string) error {
    existing, _ := s.userRepo.GetByEmail(context.Background(), email)
    if existing != nil {
        return ErrEmailExists
    }
    return nil
}

func (s *UserRegistrationService) validatePasswordStrength(password string) error {
    // 密码强度验证逻辑
    if len(password) < 6 {
        return ErrWeakPassword
    }
    // ... 更多验证
    return nil
}
```

**何时使用领域服务**:
1. 操作涉及多个聚合根
2. 需要访问外部资源（通过 Repository 接口）
3. 复杂的业务验证逻辑

---

### 1.5 领域事件

#### ✅ 正确示例

```go
// internal/domain/user/event/user_events.go

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
    UserID    uuid.UUID
    Email     string
    CreatedAt time.Time
    EventID   string
    EventType string
    AggregateID uuid.UUID
    OccurredAt  time.Time
    Version     int
}

// NewUserRegisteredEvent 构造函数
func NewUserRegisteredEvent(userID uuid.UUID, email string) *UserRegisteredEvent {
    return &UserRegisteredEvent{
        UserID:      userID,
        Email:       email,
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

## 2. Application Layer（应用层）

### 职责定义

**只包含**:
- ✅ 应用服务（协调领域对象完成任务）
- ✅ DTO（数据传输对象）
- ✅ Assembler（DTO ↔ Entity 转换器）
- ✅ 命令查询接口（CQRS）

**禁止包含**:
- ❌ 业务逻辑判断（应该在 Domain）
- ❌ 基础设施代码（应该在 Infrastructure）

---

### 2.1 应用服务

#### ✅ 正确示例

```go
// internal/application/user/service/user_service.go

// UserService 用户应用服务
type UserService struct {
    userRepo     repository.UserRepository
    assembler    *UserAssembler
    eventBus     EventBus
}

// NewUserService 构造函数
func NewUserService(
    userRepo repository.UserRepository,
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

#### ❌ 错误示例（当前问题）

```go
// ❌ 问题：Application Service 包含业务判断
func (s *authenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
    // ❌ 业务校验在这里
    if err := validator.ValidatePasswordStrength(req.Password); err != nil {
        return nil, errPkg.ErrInvalidPassword
    }
    
    // ❌ 检查邮箱是否已存在（应该在 Domain）
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, errPkg.ErrUserExists
    }
    
    // ❌ 验证租户限制（复杂业务逻辑）
    count, _ := s.tenantMemberRepo.ListByTenant(ctx, *tenantID)
    tenant, _ := s.tenantRepo.GetByID(ctx, *tenantID)
    if len(count) >= tenant.MaxMembers {
        return nil, errPkg.ErrTenantLimitExceed
    }
    
    // ❌ 直接创建实体（应该用 Factory 或 Domain Service）
    newUser := &entity.User{...}
}
```

**改进方案**: 将业务逻辑移到 Domain Service，见 [code-review-report.md#L98](code-review-report.md#L98)

---

### 2.2 DTO 设计

#### ✅ 正确示例

```go
// internal/application/user/dto/user_dto.go

// UserResponse 用户响应 DTO（扁平结构）
type UserResponse struct {
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
```

**DTO 设计原则**:
1. 扁平结构（避免嵌套过深）
2. 使用基本类型（string, int, bool）
3. 不使用领域对象（Entity, ValueObject）
4. 带 JSON 标签和验证标签

---

### 2.3 Assembler 模式

#### ✅ 正确示例

```go
// internal/application/user/assembler/user_assembler.go

// UserAssembler 用户 DTO 转换器
type UserAssembler struct{}

// NewUserAssembler 构造函数
func NewUserAssembler() *UserAssembler {
    return &UserAssembler{}
}

// ToResponse 领域实体 → DTO
func (a *UserAssembler) ToResponse(user *entity.User) *dto.UserResponse {
    if user == nil {
        return nil
    }
    
    return &dto.UserResponse{
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

## 3. Infrastructure Layer（基础设施层）

### 职责定义

**只包含**:
- ✅ 仓储实现（Repository Implementation）
- ✅ 数据模型（Data Model）
- ✅ DAO（数据访问对象）
- ✅ 第三方服务集成（Redis, Kafka, etc.）

**禁止包含**:
- ❌ 业务逻辑
- ❌ 直接返回 Model 给应用层

---

### 3.1 仓储实现

#### ✅ 正确示例

```go
// internal/infrastructure/persistence/gorm/repo/user_repository.go

// userRepository 用户仓储实现
type userRepository struct {
    db        *gorm.DB
    assembler UserEntityAssembler  // 使用 Assembler 转换
}

// NewUserRepository 构造函数
func NewUserRepository(db *gorm.DB, assembler UserEntityAssembler) repository.UserRepository {
    return &userRepository{db: db, assembler: assembler}
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
    
    // 2. Model → Entity 转换（使用 Assembler，处理错误）
    return r.assembler.ToEntity(&model)
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    // 1. Entity → Model 转换
    dbModel, err := r.assembler.ToModel(user)
    if err != nil {
        return err
    }
    
    // 2. 保存到数据库
    return r.db.Create(dbModel).Error
}
```

**关键规则**:
1. 仓储接口在 Domain 层定义，实现在 Infrastructure 层
2. Model ↔ Entity 转换必须在仓储内部完成
3. 不能直接返回 Model 给应用层

---

#### ❌ 错误示例（当前问题）

```go
// ❌ 问题 1: 忽略转换错误
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    var model model.User
    r.db.First(&model, id)  // ❌ 忽略错误
    
    email, _ := valueobject.NewEmailFromString(model.Email)  // ❌ 忽略错误
    
    return &entity.User{
        ID:    uuid.MustParse(*model.ID),  // ❌ 可能 panic
        Email: email,
        // ...
    }, nil
}

// ❌ 问题 2: 手动转换代码冗长
func (r *userRepository) toEntity(m *model.User) *entity.User {
    // 几十行重复的转换代码
    // 容易出错且难以维护
}
```

**改进方案**: 使用 Assembler 模式，集中管理转换逻辑

---

### 3.2 数据模型

#### ✅ 正确示例

```go
// internal/infrastructure/persistence/gorm/model/user.go

// User 用户数据模型（GORM Model）
type User struct {
    ID        *string    `gorm:"type:uuid;primaryKey" json:"id"`
    Email     string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
    Password  string     `gorm:"size:255;not null" json:"-"`  // 敏感字段不导出
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

## 4. Interfaces Layer（接口层）

### 职责定义

**只包含**:
- ✅ HTTP Handler（协议转换）
- ✅ Middleware（横切关注点）
- ✅ Router（路由注册）

**禁止包含**:
- ❌ 业务逻辑
- ❌ 直接访问仓储
- ❌ 返回领域实体

---

### 4.1 HTTP Handler

#### ✅ 正确示例

```go
// internal/interfaces/http/user/handler.go

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
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // 1. 绑定参数（使用 validator 验证）
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest("无效的用户 ID"))
        return
    }
    
    // 2. 调用应用服务
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        // 统一错误处理（或使用中间件）
        h.logger.Error("获取用户失败", zap.Error(err))
        c.JSON(http.StatusInternalServerError, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(http.StatusOK, response.Success(user))
}
```

**Handler 职责**:
1. 参数绑定和验证
2. 调用应用服务
3. 格式化响应
4. 记录日志

---

#### ❌ 错误示例（当前问题）

```go
// ❌ 问题 1: Handler 包含业务逻辑
func (h *UserHandler) UpdateUser(c *gin.Context) {
    var req dto.UpdateUserRequest
    c.ShouldBindJSON(&req)
    
    // ❌ 业务判断在这里
    if req.Email != "" && !isValidEmail(req.Email) {
        c.JSON(400, gin.H{"error": "无效邮箱"})
        return
    }
    
    // ❌ 直接访问仓储
    user, _ := userRepo.GetByID(ctx, userID)
    
    // ❌ 修改实体
    user.Email = req.Email
    userRepo.Update(ctx, user)
}

// ❌ 问题 2: 错误处理重复
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ... 业务逻辑
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})  // 每个 Handler 都重复
        return
    }
}
```

**改进方案**: 
1. 业务逻辑移到 Domain/Application
2. 使用统一错误处理中间件

---

## 5. 依赖注入（Wire）

### 模块化组织

#### ✅ 正确示例

```go
// internal/infrastructure/wire/providers_database.go
package wire

import (
    "github.com/google/wire"
    "gorm.io/gorm"
)

// DatabaseProviderSet 数据库相关 Provider
var DatabaseProviderSet = wire.NewSet(
    InitializeDB,
    InitializeTransaction,
)

func InitializeDB(cfg *config.Config) (*gorm.DB, error) {
    // 数据库连接初始化
}

func InitializeTransaction(db *gorm.DB) transaction.UnitOfWork {
    return transaction.NewGormUnitOfWork(db)
}
```

```go
// internal/infrastructure/wire/user.go
package wire

import (
    "github.com/google/wire"
)

// UserModuleSet 用户模块 Provider
var UserModuleSet = wire.NewSet(
    NewUserRepository,
    NewUserService,
    NewUserHandler,
)
```

---

#### ❌ 错误示例（当前问题）

```go
// ❌ 所有 Provider 都在一个文件（135 行）
// internal/infrastructure/wire/providers.go

func InitializeDB(cfg *config.Config) (*gorm.DB, error)
func InitializeRedis(cfg *config.Config) (*redis.Client, error)
func InitializeJWTService(cfg *config.Config) entity.JWTService
func InitializeCasbinEnforcer() *casbin.Enforcer
// ... 几十个函数混在一起
```

**改进方案**: 按功能模块拆分到不同文件

---

## 📋 分层检查清单

在审查代码时，请确认：

### Domain Layer
- [ ] 是否只有纯业务逻辑？
- [ ] 是否有基础设施依赖？（❌不应该有）
- [ ] 实体是否有业务方法？（✅应该有）
- [ ] 值对象是否不可变且有验证？
- [ ] 重要的状态变更是否发布事件？

### Application Layer
- [ ] 是否只负责编排？
- [ ] 是否有业务判断逻辑？（❌应该在 Domain）
- [ ] 是否返回 DTO？（✅而非 Entity）

### Infrastructure Layer
- [ ] 是否实现了 Repository 接口？
- [ ] Model ↔ Entity 转换是否在内部完成？
- [ ] 是否直接返回 Model？（❌不应该）

### Interfaces Layer
- [ ] 是否只做协议转换？
- [ ] 是否有业务逻辑？（❌不应该有）
- [ ] 参数是否验证？

---

**版本**: v1.0  
**生效日期**: 2026-03-06  
**下次回顾**: 2026-03-20
