# Go DDD Scaffold 技术架构文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的技术架构设计，包括整体架构模式、分层设计、技术选型理由以及各组件间的交互关系。

## 整体架构设计

### 架构模式选择

项目采用 **Clean Architecture + DDD** 的混合架构模式，结合了两种架构思想的优势：

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│  (HTTP Controllers, CLI Commands, gRPC Services)        │
├─────────────────────────────────────────────────────────┤
│                    Application Layer                     │
│  (Use Cases, Application Services, DTOs)                │
├─────────────────────────────────────────────────────────┤
│                      Domain Layer                        │
│  (Entities, Value Objects, Aggregates, Domain Services) │
├─────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                     │
│  (Repositories, External Services, Persistence)         │
└─────────────────────────────────────────────────────────┘
```

### 核心设计原则

1. **依赖倒置原则**：高层模块不依赖低层模块，都依赖抽象
2. **单一职责原则**：每层都有明确的职责边界
3. **开闭原则**：对扩展开放，对修改封闭
4. **里氏替换原则**：子类型必须能够替换它们的基类型

## 分层详细设计

### 1. 领域层 (Domain Layer)

#### 核心职责
- 业务逻辑的核心表达
- 领域概念的精确建模
- 业务规则的封装实现

#### 主要组件

**实体 (Entities)**
```go
// 用户实体 - 具有唯一标识的领域对象
type User struct {
    id       UserID        // 值对象作为标识
    username string
    email    Email         // 值对象
    password Password      // 值对象
    status   UserStatus    // 值对象
    tenants  []TenantID    // 关联的租户
}

func (u *User) ChangePassword(oldPass, newPass string) error {
    if !u.password.Matches(oldPass) {
        return errors.New("invalid old password")
    }
    u.password = NewPassword(newPass)
    return nil
}
```

**值对象 (Value Objects)**
```go
// 用户ID值对象 - 不可变，基于值相等
type UserID int64

func NewUserID(id int64) UserID {
    if id <= 0 {
        panic("invalid user id")
    }
    return UserID(id)
}

func (uid UserID) Equals(other UserID) bool {
    return uid == other
}
```

**聚合根 (Aggregate Roots)**
```go
// 用户聚合根 - 管理用户相关的业务一致性
type UserAggregate struct {
    user     User
    profiles []UserProfile
    settings UserSettings
}

func (ua *UserAggregate) AddProfile(profile UserProfile) error {
    // 业务规则验证
    if len(ua.profiles) >= MaxProfilesPerUser {
        return errors.New("maximum profiles reached")
    }
    ua.profiles = append(ua.profiles, profile)
    return nil
}
```

**领域服务 (Domain Services)**
```go
// 密码策略领域服务
type PasswordPolicyService struct{}

func (ps *PasswordPolicyService) Validate(password string) error {
    if len(password) < 8 {
        return errors.New("password too short")
    }
    // 更多复杂验证规则...
    return nil
}
```

### 2. 应用层 (Application Layer)

#### 核心职责
- 协调领域对象完成业务用例
- 处理跨聚合的业务逻辑
- 定义应用服务接口

#### 主要用例实现
```go
// 用户管理用例
type UserUseCase struct {
    userRepo     UserRepository
    tenantRepo   TenantRepository
    eventBus     EventBus
    passwordSvc  PasswordPolicyService
}

func (uc *UserUseCase) CreateUser(req CreateUserRequest) (*User, error) {
    // 1. 验证业务规则
    if err := uc.passwordSvc.Validate(req.Password); err != nil {
        return nil, err
    }
    
    // 2. 创建领域对象
    user := NewUser(req.Username, req.Email, req.Password)
    
    // 3. 持久化
    if err := uc.userRepo.Save(user); err != nil {
        return nil, err
    }
    
    // 4. 发布领域事件
    uc.eventBus.Publish(UserCreatedEvent{UserID: user.ID()})
    
    return user, nil
}
```

### 3. 接口层 (Interfaces Layer)

#### 核心职责
- 处理外部请求和响应
- 协议转换（HTTP/gRPC/CLI）
- 输入验证和错误处理

#### HTTP控制器示例
```go
type UserController struct {
    userUseCase UserUseCase
    logger      *zap.Logger
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ctrl.handleValidationError(c, err)
        return
    }
    
    user, err := ctrl.userUseCase.CreateUser(req)
    if err != nil {
        ctrl.handleBusinessError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, UserResponseFromDomain(user))
}
```

### 4. 基础设施层 (Infrastructure Layer)

#### 核心职责
- 实现领域层定义的接口
- 处理技术细节（数据库、缓存、外部服务）
- 提供基础设施服务

#### 仓储实现示例
```go
type UserDatabaseRepository struct {
    db *gorm.DB
}

func (repo *UserDatabaseRepository) Save(user *User) error {
    return repo.db.Transaction(func(tx *gorm.DB) error {
        userModel := UserModelFromDomain(user)
        if err := tx.Save(userModel).Error; err != nil {
            return err
        }
        
        // 保存关联关系
        return repo.saveUserTenants(tx, user)
    })
}

func (repo *UserDatabaseRepository) FindByID(id UserID) (*User, error) {
    var userModel UserModel
    if err := repo.db.Where("id = ?", int64(id)).First(&userModel).Error; err != nil {
        return nil, err
    }
    
    return UserFromModel(userModel), nil
}
```

## 技术选型详解

### 核心技术栈

#### 后端技术选型

**Go 1.25.6**
- **选择理由**：高性能、静态类型、并发支持好、生态成熟
- **优势**：编译速度快、部署简单、内存占用低
- **适用场景**：高并发Web服务、微服务、API网关

**Gin Web框架**
- **选择理由**：轻量级、高性能、中间件生态丰富
- **核心特性**：路由性能优异、错误恢复机制、验证中间件
- **性能表现**：比标准库快40倍的路由匹配速度

**GORM ORM**
- **选择理由**：功能完整、社区活跃、支持多种数据库
- **关键特性**：关联预加载、事务支持、迁移工具集成
- **性能优化**：批量操作、连接池管理、查询缓存

**PostgreSQL**
- **选择理由**：ACID特性、JSON支持、扩展性好
- **优势**：数据完整性保障、复杂查询支持、地理信息系统
- **适用场景**：OLTP系统、数据分析、复杂业务逻辑

**Redis**
- **选择理由**：高性能、丰富的数据结构、持久化支持
- **使用场景**：会话存储、缓存、分布式锁、消息队列
- **集群支持**：主从复制、哨兵模式、集群模式

#### 主键生成方案

**Snowflake ID**
```
时间戳(41位) + 机器ID(10位) + 序列号(12位)
```
- **优势**：全局唯一、趋势递增、分布式友好
- **性能**：纳秒级生成速度，无网络依赖
- **存储**：64位整数，索引效率高

#### 认证授权方案

**JWT双Token机制**
```
Access Token (短期) + Refresh Token (长期)
```
- **安全性**：签名验证、过期控制、撤销机制
- **性能**：无状态认证，服务端无需存储会话
- **扩展性**：支持负载均衡和水平扩展

### 前端技术选型（预留）

**React 18**
- **选择理由**：组件化、虚拟DOM、生态丰富
- **核心特性**：Hooks、Suspense、并发渲染
- **开发体验**：热重载、TypeScript支持、调试工具

**Tailwind CSS**
- **选择理由**：实用优先、原子化CSS、开发效率高
- **优势**：减少CSS文件体积、响应式设计简单
- **定制性**：可通过配置文件深度定制

## 架构模式详解

### DDD战术模式应用

#### 聚合设计原则
```
用户聚合边界
├── User (聚合根)
├── UserProfile (实体)
├── UserSettings (值对象)
└── UserTenant (值对象)
```

**设计要点**：
- 聚合根负责维护一致性边界
- 跨聚合通过领域事件通信
- 聚合尽量保持小巧专注

#### 领域事件模式
```go
// 领域事件定义
type UserRegisteredEvent struct {
    UserID    UserID
    Username  string
    Timestamp time.Time
}

// 事件发布器
type EventBus interface {
    Publish(event DomainEvent) error
    Subscribe(eventType string, handler EventHandler)
}

// 异步事件处理
func (h *EmailNotificationHandler) Handle(event DomainEvent) {
    switch evt := event.(type) {
    case UserRegisteredEvent:
        h.sendWelcomeEmail(evt.UserID, evt.Username)
    }
}
```

### CQRS模式（预留）

**命令查询职责分离**：
- **命令端**：处理写操作，保证数据一致性
- **查询端**：优化读操作，支持复杂查询

```go
// 命令处理器
type CreateUserCommandHandler struct {
    userRepo UserRepository
}

func (h *CreateUserCommandHandler) Handle(cmd CreateUserCommand) error {
    user := NewUser(cmd.Username, cmd.Email, cmd.Password)
    return h.userRepo.Save(user)
}

// 查询服务
type UserQueryService struct {
    db *gorm.DB
}

func (qs *UserQueryService) GetUserList(filter UserFilter) ([]UserDTO, error) {
    var users []UserDTO
    // 优化的只读查询
    return qs.db.Select("id, username, email, created_at").
        Where("status = ?", Active).
        Find(&users).Error
}
```

## 部署架构

### 单体应用部署
```
┌─────────────────────────────────────────┐
│              Load Balancer               │
├─────────────────────────────────────────┤
│           Application Servers            │
│        [Go App Instance 1..N]           │
├─────────────┬───────────────────────────┤
│   Cache     │        Database           │
│  (Redis)    │     (PostgreSQL)          │
└─────────────┴───────────────────────────┘
```

### 配置管理策略

**环境变量配置**：
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=scaffold_db
DB_USER=postgres
DB_PASSWORD=${SECRET_DB_PASSWORD}

# Redis配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=${SECRET_REDIS_PASSWORD}

# JWT配置
JWT_SECRET=${SECRET_JWT_KEY}
JWT_ACCESS_EXPIRE=30m
JWT_REFRESH_EXPIRE=7d
```

**YAML配置文件**：
```yaml
server:
  port: ${SERVER_PORT:8080}
  mode: ${ENV_MODE:debug}

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  max_idle_conns: 10
  max_open_conns: 100
```

## 性能优化策略

### 数据库优化
- **连接池配置**：合理的最大连接数和空闲连接数
- **索引策略**：复合索引、部分索引、表达式索引
- **查询优化**：预加载、批处理、分页查询

### 缓存策略
- **多级缓存**：本地缓存 + Redis缓存
- **缓存穿透**：布隆过滤器、空值缓存
- **缓存雪崩**：随机过期时间、熔断机制

### 并发处理
- **Goroutine池**：限制并发数量，防止资源耗尽
- **限流降级**：令牌桶算法、漏桶算法
- **超时控制**：合理的超时设置和重试机制

这个技术架构文档为项目提供了完整的技术蓝图，确保团队成员对系统设计有一致的理解。