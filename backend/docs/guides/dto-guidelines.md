# DTO 设计与使用规范

## 文档概述

本文档定义了 go-ddd-scaffold 项目中数据传输对象（DTO）的设计规范和使用原则，确保各层之间的数据传递清晰、一致。

## 核心原则

### 1. 分层职责分离

- **Input DTO (Commands)**: 用于接收外部请求参数
- **Output DTO (Results)**: 用于返回操作结果
- **Domain Objects**: 领域对象，不跨层传递

### 2. DTO 命名规范

#### Input DTO - Command 格式
```go
// 格式：<Action><Entity>Command
type RegisterUserCommand struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type AuthenticateUserCommand struct {
    Username  string `json:"username" validate:"required"`
    Password  string `json:"password" validate:"required"`
    IPAddress string `json:"ip_address,omitempty"`
    UserAgent string `json:"user_agent,omitempty"`
}
```

#### Output DTO - Result 格式
```go
// 格式：<Action>Result
type RegisterUserResult struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

type GetUserResult struct {
    ID          int64     `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    DisplayName string    `json:"display_name"`
    Status      int32     `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

type AuthenticateUserResult struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}
```

### 3. Application Service 返回类型规范

**重要规则**：Application Service 的所有方法必须返回 Result DTO，而不是领域对象。

```go
// ✅ 正确：返回 Result DTO
type UserService interface {
    RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error)
    AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error)
    GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error)
    UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
}

// ❌ 错误：返回领域对象
type UserService interface {
    RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error) // 不允许
    GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error)        // 不允许
}
```

## DTO 组织方式

### Application 层文件结构

```
internal/application/user/
├── service.go           # Application Service 接口和实现
├── dtos.go              # 所有 DTOs（按功能分组）
└── event_handlers.go    # 领域事件处理器
```

### dtos.go 内容组织

```go
package user

import (
    "time"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// ============================================================================
// Input DTOs (Commands)
// ============================================================================

// RegisterUserCommand 用户注册命令
type RegisterUserCommand struct { ... }

// AuthenticateUserCommand 用户认证命令
type AuthenticateUserCommand struct { ... }

// UpdateUserProfileCommand 更新用户资料命令
type UpdateUserProfileCommand struct { ... }

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct { ... }

// ============================================================================
// Output DTOs (Results)
// ============================================================================

// RegisterUserResult 用户注册结果
type RegisterUserResult struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// AuthenticateUserResult 认证结果
type AuthenticateUserResult struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

// GetUserResult 获取用户结果
type GetUserResult struct {
    ID          int64     `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    DisplayName string    `json:"display_name"`
    Status      int32     `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

// ============================================================================
// Auxiliary DTOs (可选)
// ============================================================================

// UserProfileUpdate 用户资料更新数据（内部传输用）
type UserProfileUpdate struct {
    DisplayName *string
    FirstName   *string
    LastName    *string
}
```

## HTTP 层 DTO 转换

### HTTP Handler 职责

HTTP Handler 负责：
1. 接收 HTTP 请求并绑定到 Request DTO
2. 将 Request DTO 转换为 Application Command
3. 调用 Application Service
4. 将 Result DTO 转换为 Response DTO

### 转换流程

```
HTTP Request 
    ↓ (Bind)
Request DTO (interfaces/http/user/request.go)
    ↓ (Mapper.ToCommand)
Command DTO (application/user/dtos.go)
    ↓ (Service.Execute)
Result DTO (application/user/dtos.go)
    ↓ (Mapper.ToResponse)
Response DTO (interfaces/http/user/response.go)
    ↓ (JSON)
HTTP Response
```

### 示例代码

```go
// interfaces/http/user/handler.go
func (h *Handler) GetUser(c *gin.Context) {
    var req GetUserRequest
    if !h.respHandler.BindUri(c, &req) {
        return
    }

    // Request DTO → Command
    userID, err := h.mapper.ParseUserID(req.UserID)
    if err != nil {
        h.respHandler.BadRequest(c, "invalid user id")
        return
    }

    // Call Application Service
    result, err := h.userService.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        h.respHandler.Error(c, err)
        return
    }

    // Result DTO → Response DTO
    h.respHandler.Success(c, toUserResponse(result))
}

// interfaces/http/user/mapper.go
func toUserResponse(result *userApp.GetUserResult) *UserResponse {
    return &UserResponse{
        ID:          result.ID,
        Username:    result.Username,
        Email:       result.Email,
        DisplayName: result.DisplayName,
        Status:      result.Status,
        CreatedAt:   result.CreatedAt.Format(time.RFC3339),
    }
}
```

## 最佳实践

### 1. 避免贫血模型

DTO 应该只包含数据，不包含业务逻辑：

```go
// ✅ 正确：纯数据对象
type RegisterUserCommand struct {
    Username string `json:"username" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// ❌ 错误：包含业务逻辑
type RegisterUserCommand struct {
    Username string
    Email    string
    Password string
    
    func (c *RegisterUserCommand) Validate() error { // 不应该有业务逻辑
        // ...
    }
}
```

### 2. 使用指针表示可选字段

```go
type UpdateUserProfileCommand struct {
    UserID      user.UserID `json:"user_id"`
    DisplayName *string     `json:"display_name,omitempty"` // 可选
    FirstName   *string     `json:"first_name,omitempty"`   // 可选
    Gender      *user.UserGender `json:"gender,omitempty"`  // 可选
}
```

### 3. 时间格式化统一

```go
// 在 Response DTO 中使用 RFC3339 格式
type GetUserResult struct {
    CreatedAt time.Time `json:"created_at"` // Go 的 time.Time 自动支持 JSON 序列化
}

// 如果需要特定格式，在 HTTP 层转换
func toUserResponse(result *userApp.GetUserResult) *UserResponse {
    return &UserResponse{
        CreatedAt: result.CreatedAt.Format(time.RFC3339),
    }
}
```

### 4. ID 类型处理

```go
// Command 中使用基础类型（int64/string）
type GetUserRequest struct {
    UserID string `uri:"id" validate:"required"` // HTTP 层使用 string
}

type GetUserResult struct {
    ID int64 `json:"id"` // Application 层使用 int64
}

// 在 Mapper 中转换
func (m *Mapper) ParseUserID(id string) (user.UserID, error) {
    intID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return user.UserID{}, err
    }
    return user.NewUserID(intID), nil
}
```

## 常见错误

### 错误 1: 直接返回领域对象

```go
// ❌ 错误
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error) {
    return s.userRepo.FindByID(ctx, userID)
}

// ✅ 正确
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error) {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    return &GetUserResult{
        ID:          user.ID().(user.UserID).Int64(),
        Username:    user.Username().Value(),
        Email:       user.Email().Value(),
        DisplayName: user.DisplayName(),
    }, nil
}
```

### 错误 2: DTO 命名不一致

```go
// ❌ 错误：混合使用不同命名风格
type RegisterUserCommand struct { ... }
type UserResult struct { ... }           // 应该改为 GetUserResult
type AuthenticationResult struct { ... } // 应该改为 AuthenticateUserResult

// ✅ 正确：统一使用 <Action><Entity>Result 格式
type RegisterUserCommand struct { ... }
type RegisterUserResult struct { ... }
type GetUserResult struct { ... }
type AuthenticateUserResult struct { ... }
```

### 错误 3: 跨层使用 DTO

```go
// ❌ 错误：Domain 层依赖 Application DTO
package user

import "github.com/shenfay/go-ddd-scaffold/internal/application/user" // 不允许

func NewUser(cmd *user.RegisterUserCommand) (*User, error) { // 不允许
    // ...
}

// ✅ 正确：Domain 层只依赖领域概念
package user

func NewUser(username, email, hashedPassword string) (*User, error) {
    // ...
}
```

## 总结

1. **Input DTO** 统一使用 `Command` 后缀
2. **Output DTO** 统一使用 `Result` 后缀
3. **Application Service** 必须返回 Result DTO，不能返回领域对象
4. **HTTP Handler** 负责 DTO 之间的转换
5. **DTO 只包含数据**，不包含业务逻辑
6. **保持命名一致性**，遵循 `<Action><Entity><Type>` 格式

遵循这些规范可以确保代码清晰、可维护，并且符合 DDD 的分层架构原则。
