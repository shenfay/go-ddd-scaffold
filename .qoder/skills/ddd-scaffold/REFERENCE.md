# DDD Scaffold 参考文档

## 完整参数说明

### 必需参数

#### `--project-name` (string)
项目名称，用于生成项目目录和配置文件。

```bash
/ddd-scaffold --project-name myapp
```

**要求**:
- 只能包含小写字母、数字和连字符
- 长度 3-50 个字符
- 不能使用保留字（如 `go`, `main`, `test`）

---

#### `--domains` (array)
领域列表，逗号分隔。默认值：`user`

```bash
/ddd-scaffold --domains user,order,product,inventory
```

**支持的领域类型**:
- `user` - 用户管理
- `order` - 订单管理
- `product` - 商品管理
- `inventory` - 库存管理
- `payment` - 支付管理
- `notification` - 通知管理
- `analytics` - 分析管理
- 自定义领域名称

**最佳实践**:
- 按业务边界划分领域
- 保持领域间松耦合
- 避免循环依赖
- 建议不超过 8 个核心领域

---

### 可选参数

#### `--style` (string)
架构风格。默认值：`standard`

**选项**:
- `minimal` - 最小化架构（仅 Domain + Application）
- `standard` - 标准架构（完整四层 + 基础集成）
- `full` - 完整架构（四层 + 全部集成 + Docker + 监控）

**选择建议**:
- **minimal**: 快速原型验证，PoC 项目
- **standard**: 生产项目推荐，平衡复杂度和功能
- **full**: 企业级应用，需要完整工具链

---

#### `--output` (string)
输出目录。默认值：`./generated`

```bash
/ddd-scaffold --project-name myapp --output ./my-project
```

**注意**:
- 目录必须为空或不存在
- 需要有写入权限
- 支持相对路径和绝对路径

---

#### `--with-examples` (flag)
包含示例代码。默认关闭

```bash
/ddd-scaffold --with-examples
```

**包含内容**:
- 完整的用户领域示例
- CRUD 操作实现
- 单元测试示例
- API 调用示例

**推荐使用**: 新手建议开启，可快速理解项目结构

---

#### `--with-tests` (flag)
包含测试文件。默认关闭

```bash
/ddd-scaffold --with-tests
```

**测试覆盖**:
- 领域层单元测试（覆盖率目标：90%+）
- 应用层集成测试
- Repository 层 Mock 测试
- Handler 层 HTTP 测试

---

#### `--with-docker` (flag)
包含 Docker 配置。默认关闭

```bash
/ddd-scaffold --with-docker
```

**生成文件**:
- `Dockerfile` - 应用镜像构建
- `docker-compose.yml` - 多容器编排
- `.dockerignore` - 构建优化

**包含服务**:
- PostgreSQL 数据库
- Redis 缓存
- NATS 消息队列
- 应用容器

---

#### `--interactive` (flag)
交互模式。默认关闭

```bash
/ddd-scaffold --interactive
```

**交互流程**:
1. 输入项目名称
2. 选择领域列表
3. 选择架构风格
4. 确认生成选项
5. 预览项目结构
6. 开始生成

**推荐使用**: 首次使用时建议选择交互模式

---

## 配置详解

### config.yaml 完整配置

```yaml
# 基础配置
scaffold:
  version: "1.0.0"
  project_name: "myapp"
  output_dir: "./generated"
  style: "standard"
  
# 领域详细配置
domains:
  - name: "user"
    enabled: true
    entities:
      - name: "User"
        fields:
          - name: "ID"
            type: "string"
            validation: "required,uuid"
          - name: "Name"
            type: "string"
            validation: "required,min=2,max=50"
          - name: "Email"
            type: "string"
            validation: "required,email"
    aggregates:
      - name: "UserAggregate"
        root_entity: "User"
        business_methods:
          - HashPassword
          - VerifyPassword
          - UpdateProfile
    repository_methods:
      - FindByEmail
      - FindByRole
      - ExistsByEmail
      
# 基础设施集成
infrastructure:
  database:
    type: "postgresql"
    host: "localhost"
    port: 5432
    name: "myapp"
    user: "myapp"
    password: "secret"
    
  cache:
    type: "redis"
    host: "localhost"
    port: 6379
    
  messaging:
    type: "nats"
    host: "localhost"
    port: 4222
    subject_prefix: "myapp.events"
    
# 生成选项
generation:
  create_examples: true
  create_tests: true
  create_docker: true
  create_readme: true
  format_code: true
  generate_swagger: true
```

---

## 生成的代码结构详解

### 1. 领域层 (Domain Layer)

#### 实体 (`entity/user.go`)

```go
package entity

import (
    "time"
    "mathfun/pkg/errors"
)

// User 用户实体
type User struct {
    ID        string
    Name      string
    Email     string
    PasswordHash string
    Role      UserRole
    Status    UserStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewUser 创建新用户（工厂方法）
func NewUser(name, email string) (*User, error) {
    if err := validateUser(name, email); err != nil {
        return nil, err
    }
    
    return &User{
        ID:        uuid.New().String(),
        Name:      name,
        Email:     email,
        Role:      UserRoleMember,
        Status:    UserStatusActive,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}

// HashPassword 密码加密
func (u *User) HashPassword(password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hash)
    return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.PasswordHash), 
        []byte(password),
    )
    return err == nil
}

// UpdateProfile 更新资料
func (u *User) UpdateProfile(name, email string) error {
    if name != "" {
        u.Name = name
    }
    if email != "" {
        u.Email = email
    }
    u.UpdatedAt = time.Now()
    return nil
}

// 辅助方法
func validateUser(name, email string) error {
    if name == "" {
        return errors.BadRequest("name is required")
    }
    if len(name) < 2 || len(name) > 50 {
        return errors.BadRequest("name length must be between 2 and 50")
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`).MatchString(email) {
        return errors.BadRequest("invalid email format")
    }
    return nil
}
```

#### 值对象 (`valueobject/user_role.go`)

```go
package valueobject

type UserRole string

const (
    UserRoleAdmin   UserRole = "admin"
    UserRoleMember  UserRole = "member"
    UserRoleGuest   UserRole = "guest"
)

// IsValid 验证角色是否有效
func (r UserRole) IsValid() bool {
    switch r {
    case UserRoleAdmin, UserRoleMember, UserRoleGuest:
        return true
    default:
        return false
    }
}

// String 实现 fmt.Stringer 接口
func (r UserRole) String() string {
    return string(r)
}
```

#### 聚合根 (`aggregate/user_aggregate.go`)

```go
package aggregate

import (
    "mathfun/domain/user/entity"
    "mathfun/domain/user/event"
)

// UserAggregate 用户聚合根
type UserAggregate struct {
    User  *entity.User
    events []domain.Event  // 待发布的领域事件
}

// NewUserAggregate 创建聚合根
func NewUserAggregate(name, email string) (*UserAggregate, error) {
    user, err := entity.NewUser(name, email)
    if err != nil {
        return nil, err
    }
    
    agg := &UserAggregate{
        User:   user,
        events: make([]domain.Event, 0),
    }
    
    // 记录领域事件
    agg.recordEvent(event.UserCreated{
        UserID:  user.ID,
        Email:   user.Email,
        CreatedAt: time.Now(),
    })
    
    return agg, nil
}

// Activate 激活用户
func (a *UserAggregate) Activate() error {
    if a.User.Status == entity.UserStatusActive {
        return errors.Conflict("user already active")
    }
    
    a.User.Status = entity.UserStatusActive
    a.recordEvent(event.UserActivated{
        UserID:    a.User.ID,
        ActivatedAt: time.Now(),
    })
    
    return nil
}

// Deactivate 停用用户
func (a *UserAggregate) Deactivate(reason string) error {
    if a.User.Status == entity.UserStatusInactive {
        return errors.Conflict("user already inactive")
    }
    
    a.User.Status = entity.UserStatusInactive
    a.recordEvent(event.UserDeactivated{
        UserID:      a.User.ID,
        Reason:      reason,
        DeactivatedAt: time.Now(),
    })
    
    return nil
}

// DomainEvents 获取待发布的领域事件
func (a *UserAggregate) DomainEvents() []domain.Event {
    events := a.events
    a.events = nil  // 清空已读取的事件
    return events
}

// recordEvent 记录领域事件（内部方法）
func (a *UserAggregate) recordEvent(evt domain.Event) {
    a.events = append(a.events, evt)
}
```

#### 仓储接口 (`repository/user_repository.go`)

```go
package repository

import (
    "context"
    "mathfun/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
    // Save 保存用户（新增或更新）
    Save(ctx context.Context, user *entity.User) error
    
    // FindByID 根据 ID 查找
    FindByID(ctx context.Context, id string) (*entity.User, error)
    
    // FindByEmail 根据邮箱查找
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    
    // FindAll 查找所有（支持分页）
    FindAll(ctx context.Context, offset, limit int) ([]*entity.User, int, error)
    
    // Delete 删除用户
    Delete(ctx context.Context, id string) error
    
    // Count 统计用户数
    Count(ctx context.Context) (int, error)
}
```

---

### 2. 应用层 (Application Layer)

#### DTO 定义 (`dto/user_dto.go`)

```go
package dto

import (
    "time"
    "mathfun/domain/user/valueobject"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
    Name   string `json:"name" validate:"omitempty,min=2,max=50"`
    Email  string `json:"email" validate:"omitempty,email"`
}

// UserResponse 用户响应
type UserResponse struct {
    ID        string                 `json:"id"`
    Name      string                 `json:"name"`
    Email     string                 `json:"email"`
    Role      valueobject.UserRole   `json:"role"`
    Status    string                 `json:"status"`
    CreatedAt time.Time              `json:"created_at"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
    Users      []UserResponse `json:"users"`
    Total      int            `json:"total"`
    Page       int            `json:"page"`
    PageSize   int            `json:"page_size"`
    TotalPages int            `json:"total_pages"`
}
```

#### 应用服务 (`service/user_app_service.go`)

```go
package service

import (
    "context"
    "mathfun/application/user/dto"
    "mathfun/domain/user/aggregate"
    "mathfun/domain/user/repository"
    "mathfun/pkg/errors"
)

// UserAppService 用户应用服务
type UserAppService struct {
    userRepo      repository.UserRepository
    eventPublisher EventPublisher
}

// NewUserAppService 构造函数
func NewUserAppService(
    userRepo repository.UserRepository,
    publisher EventPublisher,
) *UserAppService {
    return &UserAppService{
        userRepo:      userRepo,
        eventPublisher: publisher,
    }
}

// CreateUser 创建用户
func (s *UserAppService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 检查邮箱是否已存在
    existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, errors.Conflict("email already exists")
    }
    
    // 创建聚合根
    userAgg, err := aggregate.NewUserAggregate(req.Name, req.Email)
    if err != nil {
        return nil, err
    }
    
    // 密码加密
    err = userAgg.User.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // 保存用户
    err = s.userRepo.Save(ctx, userAgg.User)
    if err != nil {
        return nil, err
    }
    
    // 发布领域事件
    for _, event := range userAgg.DomainEvents() {
        s.eventPublisher.Publish(ctx, event)
    }
    
    // 返回响应
    return s.toResponse(userAgg.User), nil
}

// GetUserByID 获取用户详情
func (s *UserAppService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.NotFound("user not found")
    }
    
    return s.toResponse(user), nil
}

// UpdateUser 更新用户
func (s *UserAppService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.NotFound("user not found")
    }
    
    // 更新资料
    err = user.UpdateProfile(req.Name, req.Email)
    if err != nil {
        return nil, err
    }
    
    // 保存变更
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, err
    }
    
    return s.toResponse(user), nil
}

// DeleteUser 删除用户
func (s *UserAppService) DeleteUser(ctx context.Context, id string) error {
    return s.userRepo.Delete(ctx, id)
}

// ListUsers 用户列表
func (s *UserAppService) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error) {
    offset := (page - 1) * pageSize
    
    users, total, err := s.userRepo.FindAll(ctx, offset, pageSize)
    if err != nil {
        return nil, err
    }
    
    totalPages := (total + pageSize - 1) / pageSize
    
    return &dto.ListUsersResponse{
        Users:      s.toResponses(users),
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }, nil
}

// 辅助方法
func (s *UserAppService) toResponse(user *domain.User) *dto.UserResponse {
    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        Status:    string(user.Status),
        CreatedAt: user.CreatedAt,
    }
}

func (s *UserAppService) toResponses(users []*domain.User) []dto.UserResponse {
    responses := make([]dto.UserResponse, len(users))
    for i, user := range users {
        responses[i] = *s.toResponse(user)
    }
    return responses
}
```

---

## 最佳实践总结

### 1. 领域设计原则

✅ **DO**:
- 保持聚合内高内聚
- 通过 ID 引用其他聚合
- 使用值对象表示不变概念
- 通过领域事件实现最终一致性

❌ **DON'T**:
- 跨聚合直接对象引用
- 在领域层依赖技术细节
- 暴露聚合内部状态
- 忽略不变条件

### 2. 命名规范

**实体命名**:
- 使用名词：`User`, `Order`, `Product`
- 避免动词：`ProcessUser`, `HandleOrder`

**服务命名**:
- 领域服务：`UserService`, `OrderService`
- 应用服务：`UserAppService`, `OrderAppService`

**Repository 命名**:
- 接口：`UserRepository` (domain 层)
- 实现：`UserRepositoryImpl` (infrastructure 层)

### 3. 错误处理

```go
// 使用统一的错误类型
if err != nil {
    if errors.IsNotFound(err) {
        return nil, errors.NotFound("resource not found")
    }
    if errors.IsConflict(err) {
        return nil, errors.Conflict("resource conflict")
    }
    return nil, errors.Internal("internal error: %v", err)
}
```

### 4. 事务管理

```go
// 在应用服务层管理事务
func (s *OrderAppService) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 业务逻辑
    // ...
    
    if err := tx.Commit().Error; err != nil {
        return err
    }
    
    return nil
}
```

---

## 故障排除

### 常见问题速查

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 编译失败：import 找不到 | Go 模块路径不对 | 运行 `go mod tidy` |
| Wire 生成失败 | Provider 未注册 | 检查 `wire.go` 中的 Provider 列表 |
| 数据库连接失败 | 配置错误 | 检查 `config.yaml` 中的数据库配置 |
| 测试失败：Mock 不对 | Mock 实现与接口不匹配 | 重新生成 Mock 或手动修复 |

---

这是参考文档的核心部分。完整的 REFERENCE.md 将包含更多配置选项、代码示例和故障排除指南。
