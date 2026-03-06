# 添加 API 端点 - 完整流程

## 🎯 目标

从零开始创建一个完整的用户查询 API，包含：
- ✅ Domain 层（Entity, ValueObject）
- ✅ Application 层（Service, DTO）
- ✅ Infrastructure 层（Repository 实现）
- ✅ Interfaces 层（HTTP Handler）
- ✅ Swagger 文档注解

---

## 📋 前置条件

在开始前，请确保已了解：
- [分层架构详解](../architecture/layers.md)
- [代码规范](../standards/code-style.md)
- [DDD 实现规范](../standards/ddd-implementation.md)

---

## 🚀 Step-by-Step 教程

### Step 1: 定义领域模型（Domain Layer）

#### 1.1 创建实体

```go
// internal/domain/user/entity/user.go

package entity

import (
    "time"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/domain/user/valueobject"
)

// User 用户聚合根
type User struct {
    ID        uuid.UUID
    Email     valueobject.Email
    Password  HashedPassword
    Nickname  valueobject.Nickname
    Avatar    *string
    Phone     *string
    Bio       *string
    Status    UserStatus
    CreatedAt time.Time
    UpdatedAt time.Time
    
    events []DomainEvent
}

// NewUser 构造函数
func NewUser(email valueobject.Email, password HashedPassword, nickname valueobject.Nickname) *User {
    return &User{
        ID:        uuid.New(),
        Email:     email,
        Password:  password,
        Nickname:  nickname,
        Status:    StatusActive,
        CreatedAt: time.Now(),
    }
}

// UpdateProfile 更新资料（业务方法）
func (u *User) UpdateProfile(nickname valueobject.Nickname, phone *string, bio *string) error {
    u.Nickname = nickname
    u.Phone = phone
    u.Bio = bio
    u.addEvent(UserProfileUpdatedEvent{
        UserID:   u.ID,
        Nickname: nickname.String(),
    })
    return nil
}
```

**关键点**:
- ✅ 使用值对象（Email, Nickname）
- ✅ 有业务方法（UpdateProfile）
- ✅ 无基础设施标签（GORM, JSON）

---

#### 1.2 创建值对象

```go
// internal/domain/user/valueobject/user_values.go

package valueobject

import (
    "errors"
    "regexp"
    "strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 构造函数（带验证）
func NewEmail(email string) (Email, error) {
    e := Email{value: strings.TrimSpace(strings.ToLower(email))}
    if !e.IsValid() {
        return Email{}, errors.New("invalid email format")
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

// Nickname 昵称值对象
type Nickname struct {
    value string
}

// NewNickname 构造函数
func NewNickname(nickname string) (Nickname, error) {
    n := Nickname{value: strings.TrimSpace(nickname)}
    if len(n.value) < 2 || len(n.value) > 50 {
        return Nickname{}, errors.New("nickname must be 2-50 characters")
    }
    return n, nil
}

func (n Nickname) String() string {
    return n.value
}
```

**关键点**:
- ✅ 不可变（只有 getter）
- ✅ 构造函数强制验证
- ✅ 实现 Equals 和 String 方法

---

### Step 2: 定义仓储接口（Domain Layer）

```go
// internal/domain/user/repository/repository.go

package repository

import (
    "context"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
    // GetByID 根据 ID 获取用户
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    
    // GetByEmail 根据邮箱获取用户
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    
    // Create 创建用户
    Create(ctx context.Context, user *entity.User) error
    
    // Update 更新用户
    Update(ctx context.Context, user *entity.User) error
    
    // Delete 删除用户
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**关键点**:
- ✅ 接口定义在 Domain 层
- ✅ 使用 Context 管理生命周期
- ✅ 返回 Entity 而非 Model

---

### Step 3: 创建应用服务（Application Layer）

#### 3.1 定义 DTO

```go
// internal/application/user/dto/user_dto.go

package dto

import "time"

// UserResponse 用户响应 DTO（扁平结构）
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Nickname  string    `json:"nickname"`
    Phone     *string   `json:"phone,omitempty"`
    Bio       *string   `json:"bio,omitempty"`
    Avatar    *string   `json:"avatar,omitempty"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUserRequest 创建用户请求 DTO
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Nickname string `json:"nickname" validate:"required,min=2,max=20"`
}

// UpdateUserRequest 更新用户请求 DTO
type UpdateUserRequest struct {
    Nickname string  `json:"nickname" validate:"omitempty,min=2,max=20"`
    Phone    *string `json:"phone"`
    Bio      *string `json:"bio"`
}
```

**关键点**:
- ✅ 扁平结构
- ✅ 基本类型（string, int, bool）
- ✅ 带验证标签

---

#### 3.2 创建 Assembler

```go
// internal/application/user/assembler/user_assembler.go

package assembler

import (
    "go-ddd-scaffold/internal/application/user/dto"
    "go-ddd-scaffold/internal/domain/user/entity"
    "go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserAssembler 用户 DTO 转换器
type UserAssembler struct{}

// NewUserAssembler 构造函数
func NewUserAssembler() *UserAssembler {
    return &UserAssembler{}
}

// ToResponse Entity → Response DTO
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

// FromCreateRequest Request → Entity
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
    
    return entity.NewUser(email, hashedPassword, nickname), nil
}
```

**关键点**:
- ✅ 集中管理转换逻辑
- ✅ 处理错误
- ✅ 双向转换

---

#### 3.3 实现应用服务

```go
// internal/application/user/service/user_service.go

package service

import (
    "context"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/application/user/assembler"
    "go-ddd-scaffold/internal/application/user/dto"
    "go-ddd-scaffold/internal/domain/user/repository"
    errPkg "go-ddd-scaffold/internal/pkg/errors"
)

// UserService 用户应用服务
type UserService struct {
    userRepo  repository.UserRepository
    assembler *assembler.UserAssembler
}

// NewUserService 构造函数
func NewUserService(
    userRepo repository.UserRepository,
    assembler *assembler.UserAssembler,
) *UserService {
    return &UserService{
        userRepo:  userRepo,
        assembler: assembler,
    }
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
    // 1. 获取聚合根
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Entity → DTO
    return s.assembler.ToResponse(user), nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. 哈希密码
    hashedPwd, err := entity.NewHashedPassword(req.Password)
    if err != nil {
        return nil, errPkg.Wrap(err, "HASH_PASSWORD_FAILED", "密码加密失败")
    }
    
    // 2. 创建实体
    user, err := s.assembler.FromCreateRequest(req, hashedPwd)
    if err != nil {
        return nil, err
    }
    
    // 3. 持久化
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 4. 返回 DTO
    return s.assembler.ToResponse(user), nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(
    ctx context.Context,
    userID uuid.UUID,
    req *dto.UpdateUserRequest,
) (*dto.UserResponse, error) {
    // 1. 获取聚合根
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. 调用领域方法（业务逻辑在 Domain）
    if req.Nickname != "" {
        nickname, _ := valueobject.NewNickname(req.Nickname)
        user.UpdateProfile(nickname, req.Phone, req.Bio)
    } else {
        user.UpdateProfile(user.Nickname, req.Phone, req.Bio)
    }
    
    // 3. 持久化
    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    // 4. 返回 DTO
    return s.assembler.ToResponse(user), nil
}
```

**关键点**:
- ✅ 只负责编排（不含业务逻辑）
- ✅ 返回 DTO 而非 Entity
- ✅ 统一错误处理（使用 AppError）

---

### Step 4: 实现仓储（Infrastructure Layer）

#### 4.1 定义数据模型

```go
// internal/infrastructure/persistence/gorm/model/user.go

package model

import (
    "time"
)

// User 用户数据模型
type User struct {
    ID        *string   `gorm:"type:uuid;primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"-"`
    Nickname  string    `gorm:"size:100;not null" json:"nickname"`
    Avatar    *string   `gorm:"size:500" json:"avatar,omitempty"`
    Phone     *string   `gorm:"size:20" json:"phone,omitempty"`
    Bio       *string   `gorm:"size:500" json:"bio,omitempty"`
    Status    string    `gorm:"size:20;not null;default:'active'" json:"status"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

**关键点**:
- ✅ 指针字段表示可选值
- ✅ 指定类型、长度、索引
- ✅ 敏感字段不导出（`json:"-"`）

---

#### 4.2 实现仓储

```go
// internal/infrastructure/persistence/gorm/repo/user_repository.go

package repo

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "go-ddd-scaffold/internal/domain/user/entity"
    "go-ddd-scaffold/internal/domain/user/repository"
    "go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
    "go-ddd-scaffold/internal/application/user/assembler"
    errPkg "go-ddd-scaffold/internal/pkg/errors"
)

// userRepository 用户仓储实现
type userRepository struct {
    db        *gorm.DB
    assembler *assembler.UserEntityAssembler
}

// NewUserRepository 构造函数
func NewUserRepository(db *gorm.DB, assembler *assembler.UserEntityAssembler) repository.UserRepository {
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

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    var model model.User
    if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    
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

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
    dbModel, err := r.assembler.ToModel(user)
    if err != nil {
        return err
    }
    
    return r.db.Save(dbModel).Error
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.Delete(&model.User{ID: ptr(id.String())}).Error
}

// 辅助函数
func ptr[T any](v T) *T {
    return &v
}
```

**关键点**:
- ✅ 实现 Domain 层定义的接口
- ✅ 使用 Assembler 转换
- ✅ 错误处理完整

---

### Step 5: 创建 HTTP Handler（Interfaces Layer）

```go
// internal/interfaces/http/user/handler.go

package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/application/user/dto"
    "go-ddd-scaffold/internal/application/user/service"
    "go-ddd-scaffold/internal/pkg/response"
)

// UserHandler 用户 HTTP Handler
type UserHandler struct {
    userService service.UserService
}

// NewUserHandler 构造函数
func NewUserHandler(userService service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户 ID 获取详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID" format(uuid)
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // 1. 绑定参数
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest("无效的用户 ID"))
        return
    }
    
    // 2. 调用应用服务
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(http.StatusOK, response.Success(user))
}

// CreateUser godoc
// @Summary 创建用户
// @Description 注册新用户
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "用户信息"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. 绑定参数
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
        return
    }
    
    // 2. 调用应用服务
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(http.StatusCreated, response.Success(user))
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 更新用户资料
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID" format(uuid)
// @Param request body dto.UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    // 1. 绑定参数
    userID, _ := uuid.Parse(c.Param("id"))
    var req dto.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
        return
    }
    
    // 2. 调用应用服务
    user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(http.StatusOK, response.Success(user))
}
```

**关键点**:
- ✅ 只负责协议转换
- ✅ 参数验证
- ✅ Swagger 注释完整

---

### Step 6: 注册路由

```go
// internal/interfaces/http/router.go

package http

import (
    "github.com/gin-gonic/gin"
    "go-ddd-scaffold/internal/interfaces/http/middleware"
)

// SetupRouter 设置路由
func SetupRouter(
    userHandler *http.UserHandler,
) *gin.Engine {
    router := gin.Default()
    
    // 全局中间件
    router.Use(middleware.CORS())
    router.Use(middleware.RateLimit())
    
    // 健康检查
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API 路由
    api := router.Group("/api")
    {
        // 用户路由
        users := api.Group("/users")
        {
            users.POST("", userHandler.CreateUser)
            users.GET("/:id", userHandler.GetUser)
            users.PUT("/:id", userHandler.UpdateUser)
        }
    }
    
    return router
}
```

---

### Step 7: Wire 依赖注入

```go
// internal/infrastructure/wire/user.go

package wire

import (
    "github.com/google/wire"
    "go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
    "go-ddd-scaffold/internal/application/user/service"
    "go-ddd-scaffold/internal/application/user/assembler"
    "go-ddd-scaffold/internal/interfaces/http/user"
)

// UserModuleSet 用户模块 Provider
var UserModuleSet = wire.NewSet(
    assembler.NewUserAssembler,
    repo.NewUserRepository,
    service.NewUserService,
    http.NewUserHandler,
)
```

```go
// internal/infrastructure/wire/wire.go

//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "go-ddd-scaffold/internal/config"
    "gorm.io/gorm"
)

// InitializeApp 初始化应用
func InitializeApp(cfg *config.Config) (*gin.Engine, error) {
    panic(wire.Build(
        // 基础设施
        InitializeDB,
        
        // 用户模块
        UserModuleSet,
    ))
}
```

运行 Wire:

```bash
cd backend/internal/infrastructure/wire
wire gen .
```

---

### Step 8: 生成 Swagger 文档

```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go -o docs

# 启动服务后访问
# http://localhost:8080/swagger/index.html
```

---

## ✅ 完成！

你已经创建了一个完整的 API 端点，包含：

- ✅ **Domain Layer**: Entity, ValueObject, Repository Interface
- ✅ **Application Layer**: Service, DTO, Assembler
- ✅ **Infrastructure Layer**: Repository Implementation, Model
- ✅ **Interfaces Layer**: HTTP Handler
- ✅ **Swagger Documentation**: API 文档自动生成

---

## 🐛 常见错误

### ❌ 错误 1: Handler 中包含业务逻辑

```go
// ❌ 错误
func (h *UserHandler) UpdateUser(c *gin.Context) {
    var req UpdateUserRequest
    c.ShouldBindJSON(&req)
    
    // ❌ 业务判断在这里
    if req.Email != "" && !isValidEmail(req.Email) {
        c.JSON(400, gin.H{"error": "无效邮箱"})
        return
    }
    
    // ❌ 直接访问仓储
    user, _ := userRepo.GetByID(ctx, userID)
}

// ✅ 正确
func (h *UserHandler) UpdateUser(c *gin.Context) {
    var req UpdateUserRequest
    c.ShouldBindJSON(&req)  // 参数验证
    
    // 调用应用服务
    user, err := userService.UpdateUser(ctx, userID, &req)
    c.JSON(200, response.Success(user))
}
```

---

### ❌ 错误 2: Application Service 返回 Entity

```go
// ❌ 错误
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    return s.userRepo.GetByID(ctx, id)  // 返回 Entity
}

// ✅ 正确
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
    user, _ := s.userRepo.GetByID(ctx, id)
    return s.assembler.ToResponse(user)  // 返回 DTO
}
```

---

### ❌ 错误 3: Repository 直接返回 Model

```go
// ❌ 错误
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    var model model.User
    r.db.First(&model, id)
    return &model, nil  // 返回 Model
}

// ✅ 正确
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    var model model.User
    r.db.First(&model, id)
    return r.assembler.ToEntity(&model)  // 转换为 Entity
}
```

---

## 📚 相关文档

- [分层架构详解](../architecture/layers.md)
- [代码规范](../standards/code-style.md)
- [DDD 实现规范](../standards/ddd-implementation.md)
- [Code Review 报告](./code-review-report.md)

---

**版本**: v1.0  
**最后更新**: 2026-03-06
