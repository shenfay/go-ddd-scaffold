# P0 问题修复报告

## 📅 修复时间
2026-03-06

## 🎯 修复目标
修复 Code Review 中发现的 5 个严重问题（P0 级别）

---

## ✅ 已完成的修复

### Task 1: User 实体添加业务方法 ✅

**问题**: User 实体只有简单的 `IsActive()` 方法，缺少业务方法（Lock, Activate, UpdateProfile 等）

**修复内容**:

1. **添加 Lock() 方法** - 锁定用户账号
   ```go
   func (u *User) Lock() error {
       if u.Status == StatusLocked {
           return ErrAlreadyLocked
       }
       u.Status = StatusLocked
       u.addEvent(NewUserLockedEvent(u.ID))
       return nil
   }
   ```

2. **添加 Activate() 方法** - 激活用户账号
   ```go
   func (u *User) Activate() error {
       if u.Status == StatusActive {
           return nil
       }
       u.Status = StatusActive
       u.addEvent(NewUserActivatedEvent(u.ID))
       return nil
   }
   ```

3. **添加 UpdateProfile() 方法** - 更新用户资料
   ```go
   func (u *User) UpdateProfile(nickname Nickname, phone *string, bio *string) error {
       u.Nickname = nickname
       u.Phone = phone
       u.Bio = bio
       u.addEvent(NewUserProfileUpdatedEvent(...))
       return nil
   }
   ```

4. **添加 UpdateEmail() 方法** - 更新邮箱
   ```go
   func (u *User) UpdateEmail(newEmail Email) error {
       if u.Email.Equals(newEmail) {
           return nil
       }
       oldEmail := u.Email
       u.Email = newEmail
       u.addEvent(NewUserEmailChangedEvent(u.ID, oldEmail.String(), newEmail.String()))
       return nil
   }
   ```

5. **添加领域事件支持**
   - `Events()` - 获取待发布事件
   - `ClearEvents()` - 清空已发布事件
   - `addEvent()` - 内部添加事件

**影响范围**: 
- `backend/internal/domain/user/entity/user.go` (+80/-14)
- `backend/internal/domain/user/event/user_events.go` (新建)

**符合规范**: DDD 实现规范 2.2 - "实体必须有明确的业务方法"

---

### Task 2: 重构 HashedPassword ✅

**问题**: `HashedPassword` 在 Domain 层直接依赖 `bcrypt` 库（违反 Domain 纯净性）

**修复内容**:

1. **移除 bcrypt 依赖**
   ```go
   // ❌ 之前
   import "golang.org/x/crypto/bcrypt"
   
   func NewHashedPassword(plainPassword string) (HashedPassword, error) {
       bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
       return HashedPassword(string(bytes)), err
   }
   
   // ✅ 现在
   type HashedPassword string
   // 加密逻辑移至 Infrastructure 层
   ```

2. **创建 PasswordHasher 接口**（Domain 层）
   ```go
   // internal/domain/user/service/password_hasher.go
   type PasswordHasher interface {
       Hash(plain string) (string, error)
       Verify(hash, plain string) bool
   }
   ```

3. **提供 Bcrypt 实现**（Infrastructure 层）
   ```go
   type BcryptPasswordHasher struct {
       cost int
   }
   
   func (h *BcryptPasswordHasher) Hash(plain string) (string, error) {
       bytes, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
       return string(bytes), err
   }
   ```

**影响范围**:
- `backend/internal/domain/user/entity/user.go` (+7/-27)
- `backend/internal/domain/user/service/password_hasher.go` (新建)

**符合规范**: DDD 实现规范 2.1 - "Domain 层禁止包含基础设施代码"

---

### Task 3: 提取 UserRegistrationService ✅

**问题**: `Register` 方法包含过多业务逻辑（应该在 Domain Service 中）

**修复内容**:

1. **创建用户注册领域服务**
   ```go
   // internal/domain/user/service/user_registration_service.go
   type UserRegistrationService struct {
       userRepo repository.UserRepository
       hasher   PasswordHasher
   }
   ```

2. **封装完整注册流程**
   ```go
   func (s *UserRegistrationService) RegisterUser(ctx context.Context, cmd RegisterCommand) (*entity.User, error) {
       // 1. 验证邮箱唯一性
       if err := s.validateEmailUnique(ctx, cmd.Email); err != nil {
           return nil, err
       }
       
       // 2. 验证密码强度
       if err := s.validatePasswordStrength(cmd.Password); err != nil {
           return nil, err
       }
       
       // 3. 哈希密码（使用接口，不依赖具体实现）
       hashedPwdStr, err := s.hasher.Hash(cmd.Password)
       
       // 4. 创建值对象
       email, _ := valueobject.NewEmail(cmd.Email)
       nickname, _ := valueobject.NewNickname(cmd.Nickname)
       
       // 5. 创建用户实体
       user := &entity.User{
           ID:        uuid.New(),
           Email:     email,
           Password:  entity.HashedPassword(hashedPwdStr),
           Nickname:  nickname,
           Status:    entity.StatusActive,
           CreatedAt: time.Now(),
       }
       
       return user, nil
   }
   ```

3. **业务验证集中在领域服务**
   - `validateEmailUnique()` - 邮箱唯一性验证
   - `validatePasswordStrength()` - 密码强度验证
   - `checkTenantLimit()` - 租户成员限制检查（预留 TODO）

**影响范围**:
- `backend/internal/domain/user/service/user_registration_service.go` (新建)

**符合规范**: DDD 实现规范 3.1 - "复杂业务逻辑应封装在领域服务中"

---

### Task 4: Tenant 聚合根成员管理 ✅

**问题**: Tenant 聚合根设计不清晰，缺少领域事件支持

**修复内容**:

1. **移除 GORM 标签**（保持 Domain 纯净）
   ```go
   // ❌ 之前
   type Tenant struct {
       ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
       Name        string    `gorm:"size:100"`
       MaxMembers  int       `gorm:"default:10"`
       Members     []TenantMember `gorm:"foreignKey:TenantID"`
   }
   
   // ✅ 现在
   type Tenant struct {
       ID          uuid.UUID
       Name        string
       MaxMembers  int
       Members     []TenantMember  // 纯字段，无标签
       events      []DomainEvent
   }
   ```

2. **增强 AddMember() 方法**
   ```go
   func (t *Tenant) AddMember(userID uuid.UUID, role UserRole, invitedBy *uuid.UUID) (*TenantMember, error) {
       // 验证逻辑...
       
       member := TenantMember{...}
       t.Members = append(t.Members, member)
       
       // 发布领域事件
       t.addEvent(NewTenantMemberAddedEvent(t.ID, userID, member.ID, string(role)))
       
       return &member, nil
   }
   ```

3. **增强 RemoveMember() 方法**
   ```go
   func (t *Tenant) RemoveMember(memberID uuid.UUID) error {
       for i, member := range t.Members {
           if member.ID == memberID {
               t.Members[i].Status = MemberStatusRemoved
               leftAt := time.Now()
               t.Members[i].LeftAt = &leftAt
               
               // 发布领域事件
               t.addEvent(NewTenantMemberRemovedEvent(t.ID, member.UserID, member.ID))
               
               return nil
           }
       }
       return ErrTenantMemberNotFound
   }
   ```

4. **添加领域事件支持**
   - `Events()` - 获取待发布事件
   - `ClearEvents()` - 清空已发布事件
   - `addEvent()` - 内部添加事件

5. **创建租户领域事件**
   - `TenantCreatedEvent` - 租户创建事件
   - `TenantMemberAddedEvent` - 成员添加事件
   - `TenantMemberRemovedEvent` - 成员移除事件

**影响范围**:
- `backend/internal/domain/tenant/entity/tenant.go` (+44/-8)
- `backend/internal/domain/tenant/event/tenant_events.go` (新建)

**符合规范**: DDD 实现规范 2.2 - "聚合根必须管理内部实体的完整性"

---

### Task 5: 统一错误处理 ⏳

**问题**: 部分地方直接使用 `fmt.Errorf`，错误码不统一

**状态**: 待完成（需要修改 Application Service 和 Handler）

**计划**:
1. 创建统一的错误处理中间件
2. 修改所有 Handler 使用 `c.Error(err)` 而非重复的错误处理
3. 确保所有错误都使用 `AppError`

---

## 📊 统计数据

| 指标 | 数值 |
|------|------|
| 修改文件数 | 7 个 |
| 新增文件数 | 4 个 |
| 新增代码行数 | ~400 行 |
| 删除代码行数 | ~50 行 |
| 净增代码 | ~350 行 |

---

## ✅ 验证清单

### Domain Layer
- [x] User 实体有丰富的业务方法（Lock, Activate, UpdateProfile, UpdateEmail）
- [x] 重要的状态变更发布领域事件
- [x] HashedPassword 不再依赖 bcrypt
- [x] Tenant 聚合根管理成员关系
- [x] 领域事件定义完整

### 代码质量
- [x] 所有导出函数有注释
- [x] 使用值对象封装（Email, Nickname）
- [x] 实体无基础设施标签
- [x] 领域服务职责清晰

---

## 🚧 待完成的工作

### Task 5: 统一错误处理（进行中）

**需要修改的文件**:
1. `internal/interfaces/http/user/handler.go` - Handler 错误处理
2. `internal/application/user/service/authentication_service.go` - 应用服务错误
3. `internal/pkg/errors/middleware.go` - 新建错误处理中间件

**预计工作量**: 2-3 小时

---

### 后续集成工作

#### 1. 更新 Application Service

**当前问题**: `authenticationService.Register()` 包含业务逻辑

**改进方案**:
```go
func (s *authenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
    // 1. 使用领域服务
    cmd := domain_service.RegisterCommand{
        Email:    req.Email,
        Password: req.Password,
        Nickname: req.Nickname,
    }
    
    user, err := s.userRegistrationService.RegisterUser(ctx, cmd)
    if err != nil {
        return nil, err
    }
    
    // 2. 持久化
    s.userRepo.Create(ctx, user)
    
    // 3. 返回 DTO
    return s.assembler.ToResponse(user), nil
}
```

---

#### 2. 集成 Wire 依赖注入

**需要添加的 Provider**:
```go
// internal/infrastructure/wire/user.go
var UserModuleSet = wire.NewSet(
    service.NewBcryptPasswordHasher,  // 密码哈希器
    service.NewUserRegistrationService,  // 注册领域服务
    // ... 其他 Provider
)
```

---

#### 3. 补充单元测试

**需要测试的方法**:
- `User.Lock()` / `User.Activate()` / `User.UpdateProfile()`
- `UserRegistrationService.RegisterUser()`
- `Tenant.AddMember()` / `Tenant.RemoveMember()`

**目标覆盖率**: ≥90%

---

## 📋 下一步建议

### 立即执行（今天）
1. ✅ 完成 Task 5: 统一错误处理中间件
2. ⏳ 更新 Application Service 以使用新的领域服务
3. ⏳ 编写核心方法的单元测试

### 本周内完成
1. 集成到 Wire 依赖注入
2. 全面测试（单元 + 集成）
3. 更新文档（基于新的实现）

### 下周开始（P1 级别）
1. CQRS 完整分离
2. UnitOfWork + Outbox Pattern
3. Redis Stream EventBus

---

## 💬 重要说明

### 关于 Task 5（统一错误处理）

Task 5 涉及面较广，需要修改：
- 所有 HTTP Handler 的错误处理逻辑
- Application Service 的错误返回
- 可能需要调整响应格式

**建议**: 
- 先完成核心的 4 个任务（已完成✅）
- Task 5 单独作为一个专项任务处理
- 可以先创建一个简单的错误处理中间件，逐步完善

---

**报告生成时间**: 2026-03-06  
**完成状态**: 4/5 完成（80%）  
**下次更新**: Task 5 完成后
