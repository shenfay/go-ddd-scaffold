# DDD 架构设计说明

## 整体架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Presentation Layer                             │
│                                                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│  │ HTTP Controller│  gRPC Service │  │  Message Consumer        │  │
│  │   (Handler)   │  │             │  │                          │  │
│  └──────┬───────┘  └──────┬───────┘  └────────────┬─────────────┘  │
└─────────┼─────────────────┼───────────────────────┼─────────────────┘
          │                 │                       │
          ▼                 ▼                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Application Layer                               │
│                  (协调领域对象完成用例)                                │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    Command Handlers                             │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌────────────────┐ │ │
│  │  │ RegisterUser    │  │ ActivateUser    │  │ ChangePassword │ │ │
│  │  │ CommandHandler  │  │ CommandHandler  │  │ CommandHandler │ │ │
│  │  └─────────────────┘  └─────────────────┘  └────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                     Query Handlers                              │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌────────────────┐ │ │
│  │  │ GetUser         │  │ ListUsers       │  │ SearchUsers     │ │ │
│  │  │ QueryHandler    │  │ QueryHandler    │  │ QueryHandler    │ │ │
│  │  └─────────────────┘  └─────────────────┘  └────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Domain Layer                                  │
│                   (纯业务逻辑，无外部依赖)                             │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                     User Aggregate                              │ │
│  │                                                                 │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │ │
│  │  │   Entity     │  │ Value Object │  │   Domain Events      │ │ │
│  │  │   User       │  │ - UserID     │  │ - UserRegistered     │ │ │
│  │  │              │  │ - Email      │  │ - UserActivated      │ │ │
│  │  │  + Activate()│  │ - UserName   │  │ - UserDeactivated    │ │ │
│  │  │  + Lock()    │  │ - Password   │  │ - UserLoggedIn       │ │ │
│  │  └──────────────┘  └──────────────┘  └──────────────────────┘ │ │
│  │                                                                 │ │
│  │  ┌──────────────────────────────────────────────────────────┐  │ │
│  │  │              Repository Interface                        │  │ │
│  │  │  - Save(user *User) error                                │  │ │
│  │  │  - FindByID(id UserID) (*User, error)                    │  │ │
│  │  │  - FindByUsername(username string) (*User, error)        │  │ │
│  │  └──────────────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                              │
│                (实现 domain 层定义的接口)                              │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐ │
│  │  Persistence     │  │    Cache         │  │   Event Store    │ │
│  │  GORMRepository  │  │  RedisCache      │  │  Kafka/RabbitMQ  │ │
│  │                  │  │                  │  │                  │ │
│  │  func Save()     │  │  func Get()      │  │  func Publish()  │ │
│  │  func FindByID() │  │  func Set()      │  │  func Subscribe()│ │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

## CQRS 数据流

### 命令侧（写操作）流程

```
HTTP POST /api/users/register
        │
        ▼
┌──────────────────┐
│  HTTP Handler    │
│  (Presentation)  │
└────────┬─────────┘
         │ RegisterUserCommand
         ▼
┌──────────────────────────────┐
│  RegisterUserCommandHandler  │
│  (Application Layer)         │
│                              │
│  1. 验证命令参数              │
│  2. 检查业务规则              │
│  3. 创建聚合根                │
│  4. 调用仓储保存              │
│  5. 发布领域事件              │
└────────┬─────────────────────┘
         │
         ├──────────────┐
         │              │
         ▼              ▼
┌─────────────┐  ┌──────────────┐
│ UserRepository│  │ EventPublisher│
│ (Interface) │  │ (Interface)  │
└──────┬──────┘  └──────┬───────┘
       │                │
       ▼                ▼
┌─────────────┐  ┌──────────────┐
│GormRepository│  │KafkaPublisher│
│(Infrastruct.)│  │(Infrastructure)│
└─────────────┘  └──────────────┘
```

### 查询侧（读操作）流程

```
HTTP GET /api/users/{id}
        │
        ▼
┌──────────────────┐
│  HTTP Handler    │
│  (Presentation)  │
└────────┬─────────┘
         │ GetUserQuery
         ▼
┌──────────────────────────────┐
│   GetUserQueryHandler        │
│   (Application Layer)        │
│                              │
│   1. 执行查询                │
│   2. 返回 DTO/Entity         │
└────────┬─────────────────────┘
         │
         ▼
┌─────────────┐
│ UserRepository│
│ (Interface) │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│GormRepository│
│(Infrastructure)│
└─────────────┘
```

## 防腐层 (ACL) 设计

```go
// Domain Layer - 定义接口
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
}

// Infrastructure Layer - 实现接口
type UserGormRepository struct {
    db *gorm.DB
}

func (r *UserGormRepository) Save(ctx context.Context, user *User) error {
    // GORM 具体实现
    userModel := convertToModel(user)
    return r.db.Save(userModel).Error
}

func (r *UserGormRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    var userModel UserModel
    err := r.db.First(&userModel, id).Error
    if err != nil {
        return nil, err
    }
    return convertToDomain(userModel), nil
}
```

## 新的目录结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── user/
│   │       ├── entity.go        # 聚合根
│   │       ├── vo.go            # 值对象
│   │       ├── events.go        # 领域事件
│   │       └── repository.go    # 仓储接口
│   │
│   ├── application/
│   │   └── user/
│   │       ├── service.go       # 应用服务接口
│   │       ├── commands/
│   │       │   └── handlers.go  # 命令处理器
│   │       └── queries/
│   │           └── handlers.go  # 查询处理器
│   │
│   ├── interfaces/
│   │   ├── http/
│   │   │   └── handlers/
│   │   └── grpc/
│   │
│   └── infrastructure/
│       ├── persistence/
│       │   └── repositories/
│       ├── cache/
│       ├── eventstore/
│       └── messaging/
│
└── shared/
    ├── ddd/
    │   ├── entity.go
    │   ├── event.go
    │   ├── value_object.go
    │   └── repository.go
    └── cqrs/
        └── command.go
```

## DDD标准符合性检查

✅ **符合项**:
1. 清晰的四层架构（Presentation → Application → Domain → Infrastructure）
2. 聚合根包含完整的业务逻辑
3. 值对象不可变且自验证
4. 领域事件驱动设计
5. 仓储接口定义在 Domain 层
6. 基础设施层实现防腐层
7. CQRS 读写分离

⚠️ **待完善项**:
1. 需要实现具体的基础设施层（GORM、Redis 等）
2. 需要添加 HTTP/gRPC 控制器
3. 需要实现命令总线和事件总线
4. 需要添加事务管理
