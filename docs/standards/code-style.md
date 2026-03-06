# Go DDD Scaffold 代码规范

## 1. 命名规范

### 1.1 包命名

```go
// ✅ 正确：使用小写，单数形式
package user
package repository
package valueobject

// ❌ 错误：避免复数和大写
package Users      // 应该用 user
package Repository // 应该用 repository
```

**规则**:
- 包名使用小写
- 使用单数形式（`user` 而非 `users`）
- 避免下划线（`user_service` → `service`）

---

### 1.2 类型命名

```go
// ✅ 正确：PascalCase，名词为主
type User struct {}
type UserRepository interface {}
type Email struct {}

// ✅ 接口命名：使用 -er/-or 后缀
type Repository interface {}
type Publisher interface {}
type Handler interface {}

// ❌ 错误：避免 IInterface 风格
type IUser interface {}     // 应该用 User
type IUserService struct {} // 应该用 UserService
```

**规则**:
- 导出类型使用 PascalCase
- 接口名称反映行为（`Repository`, `Service`, `Handler`）
- 不使用 `I` 前缀

---

### 1.3 变量和函数命名

```go
// ✅ 正确：camelCase，语义清晰
var userID uuid.UUID
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
func (u *User) UpdateEmail(email Email) error

// ❌ 错误：避免缩写和模糊命名
var uid uuid.UUID           // 应该用 userID
func createUserReq(...)     // 应该用 CreateUserRequest
func (u *User) UpdEmail(...) // 应该用 UpdateEmail
```

**常用词完整拼写**:
| 缩写 | 完整形式 |
|------|----------|
| ID   | identifier |
| Req  | request |
| Resp | response |
| Err  | error |
| Msg  | message |

---

### 1.4 常量命名

```go
// ✅ 正确：PascalCase，语义清晰
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    MaxNameLength  = 100
)

// 枚举类型
type UserStatus string
const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
)
```

---

## 2. 注释规范

### 2.1 包注释

```go
// ✅ 正确：每个包必须有包注释
// Package user 提供用户领域的实体、值对象和服务
package user

// Package repository 定义用户领域仓储接口
package repository
```

---

### 2.2 导出元素注释

```go
// ✅ 正确：完整的句子，说明用途和行为
// User 用户聚合根
// 
// 包含用户的基础信息和行为方法。
// 用户的所有状态变更都应该通过其方法进行。
type User struct {
    ID    uuid.UUID
    Email Email
}

// UpdateEmail 更新用户邮箱
//
// 验证新邮箱的有效性，并在变更后发布 UserEmailChangedEvent 事件。
// 如果新旧邮箱相同，返回 nil。
func (u *User) UpdateEmail(email Email) error {
    // ...
}

// UserRepository 用户仓储接口
//
// 提供用户的持久化操作，隐藏底层数据存储细节。
type UserRepository interface {
    // GetByID 根据 ID 获取用户
    // 如果用户不存在，返回 ErrUserNotFound
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    
    // Create 创建新用户
    // 如果邮箱已存在，返回 ErrUserExists
    Create(ctx context.Context, user *User) error
}
```

**规则**:
- 所有导出元素必须有注释
- 使用完整句子，首字母大写
- 说明行为、边界条件、返回值
- 接口方法要说明错误情况

---

### 2.3 实现注释

```go
// ✅ 正确：解释"为什么"，而非"是什么"
// 使用 bcrypt 加密密码，成本因子为 10
// 这是在安全性和性能之间的平衡选择
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,
)

// ❌ 错误：重复代码本身
// 生成哈希密码
hashedPassword, err := bcrypt.GenerateFromPassword(...)
```

---

## 3. 错误处理规范

### 3.1 统一使用 AppError

```go
// ✅ 正确：使用 pkg/errors
import errPkg "go-ddd-scaffold/internal/pkg/errors"

func GetUser(id uuid.UUID) (*User, error) {
    if id == uuid.Nil {
        return nil, errPkg.ErrInvalidParameter.WithDetails("user ID is required")
    }
    
    user, err := repo.GetByID(id)
    if err != nil {
        return nil, errPkg.Wrap(err, "GET_USER_FAILED", "failed to get user from database")
    }
    
    return user, nil
}

// ❌ 错误：裸 error
func GetUser(id uuid.UUID) (*User, error) {
    if id == uuid.Nil {
        return nil, fmt.Errorf("user ID is required") // ❌
    }
    return nil, errors.New("not found") // ❌
}
```

---

### 3.2 错误码规范

```go
// 错误码格式：DOMAIN_ERROR_TYPE
var (
    // 用户相关
    ErrUserNotFound     = errors.New("USER_NOT_FOUND", "用户不存在")
    ErrUserExists       = errors.New("USER_ALREADY_EXISTS", "用户已存在")
    ErrInvalidEmail     = errors.New("INVALID_EMAIL", "邮箱格式不正确")
    ErrInvalidPassword  = errors.New("INVALID_PASSWORD", "密码不正确")
    
    // 租户相关
    ErrTenantNotFound   = errors.New("TENANT_NOT_FOUND", "租户不存在")
    ErrTenantLimitExceed = errors.New("TENANT_LIMIT_EXCEEDED", "超过租户限制")
)
```

**规则**:
- 错误码使用全大写下划线
- 错误消息简洁明确
- 按领域分组组织

---

### 3.3 错误包装

```go
// ✅ 正确：三层包装
// 1. 底层错误
if err := db.Query(); err != nil {
    return errPkg.Wrap(err, "DATABASE_QUERY_FAILED", "查询用户失败")
}

// 2. 业务错误
if err := validateEmail(req.Email); err != nil {
    return errPkg.ErrInvalidEmail.WithDetails(req.Email)
}

// 3. 应用错误
if err := service.CreateUser(req); err != nil {
    logger.Error("创建用户失败", zap.Error(err))
    return errPkg.Wrap(err, "CREATE_USER_FAILED", "无法创建新用户")
}
```

---

## 4. 测试规范

### 4.1 测试文件组织

```
internal/domain/user/
├── entity/
│   ├── user.go
│   └── user_test.go      # 实体测试
├── valueobject/
│   ├── email.go
│   └── email_test.go     # 值对象测试
└── repository/
    └── repository.go
```

---

### 4.2 测试命名

```go
// ✅ 正确：TestType_Method_Scenario_Result
func TestUser_UpdateEmail_ValidEmail_Success(t *testing.T)
func TestUser_UpdateEmail_InvalidEmail_Error(t *testing.T)
func TestUser_UpdateEmail_SameEmail_NoChange(t *testing.T)

// 简化版本（也接受）
func TestUser_UpdateEmail(t *testing.T)
```

---

### 4.3 测试结构

```go
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

**规则**:
- 使用表驱动测试（多个场景）
- Given/When/Then 三段式
- 测试独立，不依赖外部状态

---

### 4.4 测试覆盖率要求

| 层级 | 覆盖率要求 | 说明 |
|------|-----------|------|
| Domain | ≥ 90% | 核心业务逻辑 |
| Application | ≥ 80% | 应用编排逻辑 |
| Infrastructure | ≥ 70% | 基础设施实现 |
| Interfaces | ≥ 60% | HTTP Handler |

---

## 5. 代码组织规范

### 5.1 文件结构

```go
// 标准 Go 文件结构
package user

// 1. import（分组排序）
import (
    // 标准库
    "context"
    "time"
    
    // 第三方库
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    
    // 内部包
    "go-ddd-scaffold/internal/domain/user/valueobject"
    errPkg "go-ddd-scaffold/internal/pkg/errors"
)

// 2. 常量定义
const (
    StatusActive = "active"
)

// 3. 类型定义
type User struct {
    ID    uuid.UUID
    Email Email
}

// 4. 构造函数
func NewUser(email Email) *User {
    return &User{
        ID: uuid.New(),
        Email: email,
    }
}

// 5. 方法和接口实现
func (u *User) UpdateEmail(email Email) error {
    // ...
}
```

---

### 5.2 导入规范

```go
// ✅ 正确：分组清晰，有别名
import (
    // Standard library
    "context"
    "encoding/json"
    
    // Third-party packages
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // Internal packages
    "go-ddd-scaffold/internal/domain/user/entity"
    errPkg "go-ddd-scaffold/internal/pkg/errors"
)

// ❌ 错误：混在一起，无分组
import (
    "context"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/domain/user/entity"
    "encoding/json"
)
```

---

## 6. DDD 分层规范

### 6.1 Domain Layer（领域层）

**职责**: 纯业务逻辑，无基础设施依赖

```go
package user

// ✅ 正确：只有业务逻辑
type User struct {
    ID       uuid.UUID
    Email    Email      // 值对象
    Password HashedPassword
}

func (u *User) UpdateEmail(email Email) error {
    // 业务验证
    if u.Email.Equals(email) {
        return nil
    }
    
    u.Email = email
    u.addEvent(UserEmailChangedEvent{...})
    return nil
}

// ❌ 错误：包含基础设施
type User struct {
    ID    uuid.UUID `gorm:"type:uuid"` // ❌ GORM 标签
    Email string    `json:"email"`     // ❌ JSON 标签
}
```

**禁止**:
- ❌ import infrastructure 包
- ❌ GORM/JSON 标签
- ❌ 数据库操作
- ❌ HTTP 相关代码

---

### 6.2 Application Layer（应用层）

**职责**: 应用编排，不含业务逻辑

```go
package service

// ✅ 正确：只负责协调
type UserService struct {
    userRepo  UserRepository
    eventBus  EventBus
}

func (s *UserService) ChangeEmail(ctx context.Context, userID uuid.UUID, email string) error {
    // 1. 获取聚合根
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. 调用领域方法（业务逻辑在领域层）
    if err := user.UpdateEmail(valueobject.NewEmail(email)); err != nil {
        return err
    }
    
    // 3. 持久化
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // 4. 发布事件
    s.eventBus.Publish(user.Events())
    
    return nil
}
```

**禁止**:
- ❌ 业务判断逻辑（应该在 Domain）
- ❌ 直接访问数据库（通过 Repository）
- ❌ 直接返回 Entity（转换为 DTO）

---

### 6.3 Infrastructure Layer（基础设施层）

**职责**: 技术实现，依赖倒置

```go
package repo

// ✅ 正确：实现领域层定义的接口
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    // 1. 查询数据库模型
    var model model.User
    if err := r.db.First(&model, id).Error; err != nil {
        return nil, err
    }
    
    // 2. 转换为领域实体
    return model.ToEntity(), nil
}
```

**禁止**:
- ❌ 直接返回 Model 给应用层
- ❌ 包含业务逻辑
- ❌ 依赖上层包

---

### 6.4 Interfaces Layer（接口层）

**职责**: 协议转换，适配不同接口

```go
package http

// ✅ 正确：只做协议转换
type UserHandler struct {
    userService application.UserService
}

func (h *UserHandler) UpdateEmail(c *gin.Context) {
    // 1. 绑定请求参数
    var req UpdateEmailRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, response.BadRequest(err))
        return
    }
    
    // 2. 调用应用服务
    if err := h.userService.ChangeEmail(c.Request.Context(), req.UserID, req.Email); err != nil {
        c.JSON(500, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(200, response.OK())
}
```

**禁止**:
- ❌ 业务逻辑
- ❌ 直接访问仓储
- ❌ 返回领域实体

---

## 7. 值对象规范

### 7.1 值对象定义

```go
// ✅ 正确：不可变，有验证
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    e := Email{value: strings.TrimSpace(email)}
    if !e.IsValid() {
        return Email{}, ErrInvalidEmail
    }
    return e, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) IsValid() bool {
    return emailRegex.MatchString(e.value)
}

// ❌ 错误：可变，无验证
type Email string  // 直接使用类型别名
```

---

### 7.2 值对象使用

```go
// ✅ 正确：在实体中使用值对象
type User struct {
    Email    Email         // 值对象
    Password HashedPassword // 值对象
}

func (u *User) UpdateEmail(emailStr string) error {
    // 创建值对象（自动验证）
    email, err := valueobject.NewEmail(emailStr)
    if err != nil {
        return err
    }
    
    u.Email = email
    return nil
}

// ❌ 错误：直接使用原始类型
type User struct {
    Email string  // ❌ 缺少验证和语义
}
```

---

## 8. 数据转移对象（DTO）规范

### 8.1 DTO 定义

```go
// ✅ 正确：扁平结构，带 JSON 标签
type User struct {
    ID        string     `json:"id"`
    Email     string     `json:"email"`
    Nickname  string     `json:"nickname"`
    Phone     *string    `json:"phone,omitempty"`
    Bio       *string    `json:"bio,omitempty"`
    CreatedAt time.Time  `json:"createdAt"`
}

// ❌ 错误：嵌套过深，缺少标签
type User struct {
    ID    uuid.UUID  // ❌ 没有 JSON 标签
    Email Email      // ❌ 使用领域值对象
}
```

---

### 8.2 Assembler 模式

```go
// ✅ 正确：专门的转换层
type UserAssembler struct{}

func (a *UserAssembler) ToDTO(user *entity.User) *dto.User {
    return &dto.User{
        ID:        user.ID.String(),
        Email:     user.Email.String(),
        Nickname:  user.Nickname.String(),
        CreatedAt: user.CreatedAt,
    }
}

func (a *UserAssembler) FromDTO(dto *dto.User) *entity.User {
    email, _ := valueobject.NewEmail(dto.Email)
    nickname, _ := valueobject.NewNickname(dto.Nickname)
    
    return &entity.User{
        ID:        uuid.MustParse(dto.ID),
        Email:     email,
        Nickname:  nickname,
    }
}
```

---

## 9. 性能优化规范

### 9.1 数据库查询

```go
// ✅ 正确：预加载关联，避免 N+1
func (r *userRepository) ListWithTenant(ctx context.Context, tenantID uuid.UUID) ([]*User, error) {
    var users []*model.User
    // 使用 Preload 加载关联数据
    err := r.db.Preload("TenantMembers").
        Where("tenant_id = ?", tenantID).
        Find(&users).Error
    return users, err
}

// ❌ 错误：N+1 查询
for _, user := range users {
    members, _ := memberRepo.ListByUser(user.ID) // ❌ 循环查询
}
```

---

### 9.2 连接池配置

```go
// ✅ 正确：生产环境配置
sqlDB.SetMaxIdleConns(10)        // 最大空闲连接
sqlDB.SetMaxOpenConns(100)       // 最大打开连接
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

// ❌ 错误：默认配置（性能差）
// 无配置，使用 GORM 默认值
```

---

## 10. 安全规范

### 10.1 输入验证

```go
// ✅ 正确：所有输入必须验证
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        return err
    }
    
    // 使用 validator 验证
    if err := validator.ValidateEmail(req.Email); err != nil {
        return errPkg.ErrInvalidEmail
    }
    
    if err := validator.ValidatePasswordStrength(req.Password); err != nil {
        return errPkg.ErrInvalidPassword
    }
}

// ❌ 错误：直接使用未验证输入
var req CreateUserRequest
c.ShouldBindJSON(&req)
service.CreateUser(req) // ❌ 没有验证
```

---

### 10.2 敏感信息处理

```go
// ✅ 正确：密码加密存储
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost,
)

// ✅ 正确：日志中脱敏
logger.Info("用户登录", 
    zap.String("user_id", userID.String()),
    // zap.String("password", password), // ❌ 不能记录密码
)

// ❌ 错误：明文密码
type User struct {
    Password string // ❌ 明文存储
}
```

---

## 遵守检查清单

在提交代码前，请确认：

- [ ] 命名符合规范（包、类型、变量）
- [ ] 所有导出元素有注释
- [ ] 错误使用 AppError 处理
- [ ] 测试覆盖率达到要求
- [ ] DDD 分层清晰（Domain/Application/Infrastructure）
- [ ] 值对象封装合理
- [ ] DTO 使用正确
- [ ] 输入已验证
- [ ] 敏感信息已脱敏

---

**版本**: v1.0  
**生效日期**: 2026-03-06  
**下次回顾**: 2026-03-20
