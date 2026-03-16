# Go DDD Scaffold 开发规范文档

## 文档概述

本文档定义了 go-ddd-scaffold 项目的开发规范和编码标准，旨在确保代码质量、提高开发效率、便于团队协作和后期维护。

## 代码结构规范

### 项目目录结构
```
go-ddd-scaffold/
├── backend/              # 后端代码目录
│   ├── cmd/             # 应用程序入口
│   │   ├── api/         # API服务启动入口
│   │   ├── cli/         # CLI工具入口
│   │   ├── demo/        # 演示程序入口
│   │   └── worker/      # 后台工作进程入口
│   ├── internal/        # 内部包（不可被外部引用）
│   │   ├── domain/      # 领域层
│   │   │   ├── user/    # 用户领域
│   │   │   │   ├── domain.go        # 包入口（类型别名导出）
│   │   │   │   ├── model/           # 聚合根和值对象
│   │   │   │   │   ├── user.go          # User 聚合根（实体）
│   │   │   │   │   ├── valueobjects.go # 值对象集合（UserName、Email、HashedPassword 等）
│   │   │   │   │   └── builder.go       # Builder 模式实现（可选）
│   │   │   │   ├── event/           # 领域事件
│   │   │   │   │   └── events.go        # 事件定义
│   │   │   │   ├── repository/      # 仓储接口
│   │   │   │   │   └── user_repository.go
│   │   │   │   └── service/         # 领域服务
│   │   │   │       └── password_hasher.go
│   │   │   ├── tenant/  # 租户领域
│   │   │   ├── order/   # 订单领域（预留）
│   │   │   └── product/ # 产品领域（预留）
│   │   ├── application/ # 应用层
│   │   │   ├── user/    # 用户应用服务
│   │   │   │   ├── service.go     # Application Service 接口和实现
│   │   │   │   ├── dtos.go        # DTOs（Commands + Results）
│   │   │   │   └── event_handlers.go # 领域事件处理器（可选）
│   │   │   └── auth/    # 认证应用服务
│   │   ├── infrastructure/ # 基础设施层
│   │   │   ├── persistence/ # 数据持久化实现
│   │   │   │   └── user_repository_impl.go
│   │   │   ├── eventstore/  # 事件存储
│   │   │   │   └── event_store.go
│   │   │   ├── cache/       # 缓存实现
│   │   │   └── config/      # 配置实现
│   │   └── interfaces/      # 接口层
│   │       ├── http/        # HTTP 接口实现
│   │       │   ├── user/    # 用户领域 HTTP
│   │       │   │   ├── handler.go     # HTTP Handler
│   │       │   │   ├── request.go     # Request DTO
│   │       │   │   ├── response.go    # Response DTO
│   │       │   │   ├── mapper.go      # DTO 转换器
│   │       │   │   └── provider.go    # 路由提供者
│   │       │   └── auth/    # 认证领域 HTTP
│   │       └── middleware/  # HTTP 中间件
│   ├── shared/          # 共享领域内核
│   │   ├── ddd/         # DDD 基础组件
│   │   │   ├── entity.go        # 聚合根基类
│   │   │   ├── value_object.go  # 值对象基类
│   │   │   ├── event.go         # 领域事件基类
│   │   │   ├── repository.go    # 仓储接口
│   │   │   └── errors.go        # 领域错误
│   │   └── kernel/      # 共享内核
│   ├── configs/         # 配置文件
│   ├── migrations/      # 数据库迁移文件
│   └── tools/           # 开发工具
│       ├── generator/   # 代码生成器
│       └── migrator/    # 迁移工具
├── docs/                # 项目文档
├── deployments/         # 部署配置
│   ├── docker/          # Docker配置
│   └── kubernetes/      # K8s配置
└── scripts/             # 脚本工具
```

### 包命名规范
- 使用小写字母，避免驼峰命名
- 包名应简洁且具有描述性
- 避免使用复数形式（如用`user`而非`users`）
- 领域包名与业务概念保持一致
- 共享包放在 `shared/` 目录下，如 `shared/ddd`、`shared/kernel`

### 领域事件命名规范
```go
// 领域事件命名格式：[领域名][动作][Event]
// 动作使用过去式，表示已发生的事件

// 用户领域事件
type UserRegisteredEvent struct { ... }  // 用户已注册
type UserActivatedEvent struct { ... }    // 用户已激活
type UserEmailChangedEvent struct { ... } // 用户邮箱已变更

// 租户领域事件
type TenantCreatedEvent struct { ... }    // 租户已创建
type TenantMemberAddedEvent struct { ... }// 租户成员已添加
```

### 领域事件实现规范
```go
// 1. 嵌入 BaseEvent
type UserActivatedEvent struct {
    *ddd.BaseEvent
    UserID      UserID    `json:"user_id"`
    ActivatedAt time.Time `json:"activated_at"`
}

// 2. 提供构造函数
func NewUserActivatedEvent(userID UserID) *UserActivatedEvent {
    event := &UserActivatedEvent{
        BaseEvent:   ddd.NewBaseEvent("UserActivated", userID, 1),
        UserID:      userID,
        ActivatedAt: time.Now(),
    }
    event.SetMetadata("event_type", "domain_event")
    event.SetMetadata("aggregate_type", "user")
    return event
}
```

### 聚合根实现规范
```go
// 1. 嵌入 BaseEntity
type User struct {
    ddd.BaseEntity
    // ... 业务字段
}

// 2. 业务方法中发布事件
func (u *User) Activate() error {
    // 验证业务规则
    if u.status != UserStatusPending {
        return ddd.NewBusinessError("USER_NOT_PENDING", "...")
    }
    
    // 执行业务逻辑
    u.status = UserStatusActive
    u.updatedAt = time.Now()
    u.IncrementVersion()
    
    // 发布领域事件
    event := NewUserActivatedEvent(u.ID().(UserID))
    u.ApplyEvent(event)
    
    return nil
}
```

### 仓储实现规范
```go
// 1. 接口定义在领域层
type UserRepository interface {
    ddd.Repository
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    // ... 其他方法
}

// 2. 实现放在基础设施层
type UserRepositoryImpl struct {
    db DB
}

// 3. Save 方法需要处理领域事件
func (r *UserRepositoryImpl) Save(ctx context.Context, u *User) error {
    // ... 保存聚合
    
    // 保存未提交的事件
    events := u.GetUncommittedEvents()
    if len(events) > 0 {
        if err := r.eventStore.AppendEvents(ctx, u.ID(), events); err != nil {
            return err
        }
        u.ClearUncommittedEvents()
    }
    
    return nil
}
```

## 命名约定

### 变量命名
```go
// 好的命名示例
var userName string
var isActive bool
var userList []User
var configMap map[string]interface{}

// 避免的命名
var usr string          // 过于简短
var user_list []User    // 混合命名风格
var UserData []User     // 驼峰与下划线混合
```

### 函数命名
```go
// 动词开头，描述函数行为
func CreateUser(user User) error
func GetUserByID(id int64) (User, error)
func IsValidEmail(email string) bool
func GenerateAccessToken(userID int64) (string, error)

// 避免模糊的命名
func Process(u User) error     // 不清楚具体处理什么
func Handle(data interface{})  // 参数类型不明确
```

### 结构体命名
```go
// 使用名词，首字母大写
type User struct {
    ID        int64  `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 接口命名
```go
// 通常以"er"结尾
type UserRepository interface {
    Save(user User) error
    FindByID(id int64) (User, error)
    Delete(id int64) error
}

type PasswordEncoder interface {
    Encode(password string) (string, error)
    Matches(rawPassword, encodedPassword string) bool
}
```

## 错误处理标准

### 错误定义规范
```go
// 统一错误码定义
const (
    SuccessCode           = 0
    InvalidParamCode      = 1001
    UnauthorizedCode      = 1002
    ForbiddenCode         = 1003
    NotFoundCode          = 1004
    InternalServerErrorCode = 5000
)

// 错误响应结构
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### 错误处理实践
```go
// 好的做法
func GetUserByID(id int64) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid user id: %d", id)
    }
    
    user, err := repo.FindByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

// 避免的做法
func BadExample(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, err // 直接返回原始错误
    }
    return user, nil
}
```

### 错误包装和日志记录
```go
import (
    "fmt"
    "errors"
    "go.uber.org/zap"
)

func ProcessUser(user User) error {
    err := validateUser(user)
    if err != nil {
        logger.Error("user validation failed", 
            zap.Error(err),
            zap.Int64("user_id", user.ID))
        return fmt.Errorf("validation error: %w", err)
    }
    
    err = saveUser(user)
    if err != nil {
        logger.Error("failed to save user",
            zap.Error(err),
            zap.Int64("user_id", user.ID))
        return fmt.Errorf("save failed: %w", err)
    }
    
    return nil
}
```

## 日志记录规范

### 日志级别使用
```go
// 使用 Uber Zap 日志库
logger.Debug("debug message", zap.String("key", "value"))
logger.Info("info message", zap.Int("count", 42))
logger.Warn("warning message", zap.Error(err))
logger.Error("error message", zap.Error(err))
logger.Fatal("fatal message", zap.Error(err)) // 程序退出
```

### 日志内容规范
```go
// 好的日志记录
logger.Info("user login successful",
    zap.Int64("user_id", userID),
    zap.String("ip_address", clientIP),
    zap.String("user_agent", userAgent))

logger.Error("database operation failed",
    zap.Error(err),
    zap.String("operation", "create_user"),
    zap.Int64("user_id", userID),
    zap.Duration("duration", time.Since(startTime)))

// 避免记录敏感信息
logger.Info("user authenticated", 
    zap.String("username", username))
    // 不要记录密码等敏感信息
```

### 结构化日志格式
```go
type LogFields struct {
    UserID      int64  `json:"user_id,omitempty"`
    TenantID    int64  `json:"tenant_id,omitempty"`
    RequestID   string `json:"request_id,omitempty"`
    Operation   string `json:"operation"`
    Duration    string `json:"duration,omitempty"`
    StatusCode  int    `json:"status_code,omitempty"`
    ErrorMessage string `json:"error_message,omitempty"`
}
```

## 数据库操作规范

### 实体定义
```go
// 使用 Snowflake ID 作为主键
type User struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"size:50;uniqueIndex" json:"username" validate:"required,min=3,max=20"`
    Email     string    `gorm:"size:100;uniqueIndex" json:"email" validate:"required,email"`
    Password  string    `gorm:"size:255" json:"-" validate:"required,min=8"`
    Status    int       `gorm:"default:1" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 多对多关系定义
type UserTenant struct {
    UserID   int64 `gorm:"primaryKey"`
    TenantID int64 `gorm:"primaryKey"`
    RoleID   int64 `gorm:"index"` // 在该租户下的角色
}
```

### 查询规范
```go
// 使用预加载避免 N+1 查询
func GetUsersWithTenants() ([]User, error) {
    var users []User
    err := db.Preload("Tenants").Find(&users).Error
    return users, err
}

// 使用事务处理关联操作
func CreateUserWithTenant(user User, tenant Tenant) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        userTenant := UserTenant{
            UserID:   user.ID,
            TenantID: tenant.ID,
            RoleID:   defaultRoleID,
        }
        
        if err := tx.Create(&userTenant).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

## API 接口规范

### 请求/响应结构
```go
// 请求结构体
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// 响应结构体
type UserResponse struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Status    int       `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 控制器实现
```go
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, APIResponse{
            Code:    InvalidParamCode,
            Message: "Invalid request parameters",
            Data:    err.Error(),
        })
        return
    }
    
    // 业务逻辑处理
    user, err := userService.CreateUser(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, APIResponse{
            Code:    InternalServerErrorCode,
            Message: "Failed to create user",
        })
        return
    }
    
    c.JSON(http.StatusCreated, APIResponse{
        Code:    SuccessCode,
        Message: "User created successfully",
        Data:    UserResponseFromDomain(user),
    })
}
```

## 测试规范

### 单元测试结构
```go
// 文件命名：xxx_test.go
func TestUserService_CreateUser(t *testing.T) {
    // 准备测试数据
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // 定义测试用例
    tests := []struct {
        name    string
        request CreateUserRequest
        wantErr bool
    }{
        {
            name: "valid user creation",
            request: CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            request: CreateUserRequest{
                Username: "testuser",
                Email:    "invalid-email",
                Password: "password123",
            },
            wantErr: true,
        },
    }
    
    // 执行测试
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := service.CreateUser(tt.request)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 测试覆盖率要求
- 核心业务逻辑：≥ 90%
- API接口层：≥ 80%
- 数据访问层：≥ 85%
- 工具函数：≥ 95%

## 代码审查清单

### 必须检查项
- [ ] 代码符合命名规范
- [ ] 错误处理完整且合理
- [ ] 日志记录适当且不包含敏感信息
- [ ] 数据库操作使用事务保护
- [ ] API响应格式统一
- [ ] 单元测试覆盖率达标
- [ ] 无明显的性能问题
- [ ] 遵循DDD分层架构原则

### 建议检查项
- [ ] 代码可读性和维护性
- [ ] 是否存在重复代码
- [ ] 异常情况处理是否完善
- [ ] 配置项是否合理抽象
- [ ] 是否遵循SOLID原则
- [ ] 注释是否清晰必要

## 持续集成规范

### Git提交规范
```
feat: 添加新功能
fix: 修复bug
docs: 更新文档
style: 代码格式调整
refactor: 代码重构
test: 添加测试
chore: 构建过程或辅助工具变动
```

### 分支管理策略
- `main`：生产环境分支
- `develop`：开发主分支
- `feature/*`：功能开发分支
- `hotfix/*`：紧急修复分支
- `release/*`：发布准备分支

这个开发规范文档为项目提供了统一的编码标准和最佳实践指导，有助于提高代码质量和团队协作效率。