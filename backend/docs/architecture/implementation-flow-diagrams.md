# 技术实现流程图

## 🏗️ 应用启动流程

### main.go 启动流程图

```mermaid
graph TD
    A[程序启动] --> B[加载配置文件];
    B --> C{配置有效？};
    
    C -->|否 | D[记录错误并退出];
    C -->|是 | E[初始化 Logger];
    
    E --> F[初始化 Sentry<br/>错误追踪];
    F --> G[初始化数据库连接];
    
    G --> H{DB 连接成功？};
    H -->|否 | I[记录错误并退出];
    H -->|是 | J[运行自动迁移];
    
    J --> K[初始化 Redis 连接];
    K --> L{Redis 连接成功？};
    L -->|否 | M[记录错误并退出];
    L -->|是 | N[初始化 Asynq Client];
    
    N --> O[创建 Infra 容器];
    O --> P[注册 Modules<br/>AuthModule, UserModule...];
    
    P --> Q{遍历 Module};
    Q --> R[调用 module.RegisterHTTP<br/>注册路由];
    R --> S{还有更多 Module?};
    S -->|是 | Q;
    S -->|否 | T[设置全局错误处理中间件];
    
    T --> U[启动 HTTP Server];
    U --> V[启动 Asynq Worker<br/>后台协程];
    
    V --> W[等待中断信号];
    W --> X[优雅关闭];
    
    X --> Y[关闭 DB 连接];
    Y --> Z[关闭 Redis 连接];
    Z --> AA[关闭 HTTP Server];
    AA --> AB[关闭 Worker];
    AB --> AC[程序退出];
    
    style O fill:#87CEEB
    style P fill:#87CEEB
    style U fill:#90EE90
    style V fill:#90EE90
    style X fill:#FFD700
```

---

## 🔧 Module 注册流程

### Module 组装详细流程

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Infra as Infra Container
    participant Auth as AuthModule
    participant User as UserModule
    participant Router as Gin Router

    Main->>Infra: NewInfra(config)
    activate Infra
    Note over Infra: 1. 加载配置<br/>2. 创建 DB 连接<br/>3. 创建 Redis 连接<br/>4. 创建 Logger<br/>5. 创建 Snowflake Node
    Infra-->>Main: *Infra
    deactivate Infra
    
    Main->>Auth: NewAuthModule(infra)
    activate Auth
    
    Note over Auth: 开始组装认证模块
    
    Auth->>Auth: dao.Use(infra.DB)
    Note over Auth: 创建 DAO Query
    
    Auth->>Auth: NewUnitOfWork(infra.DB, daoQuery)
    Note over Auth: 创建工作单元
    
    Auth->>Auth: NewJWTService(...)<br/>jwtSvc.SetRedisClient(infra.Redis)
    Note over Auth: 创建 JWT 服务并注入 Redis
    
    Auth->>Auth: NewTokenServiceAdapter(jwtSvc)
    Note over Auth: ⭐ 创建适配器<br/>将 JWTService 转换为 Port
    
    Auth->>Auth: infra.Snowflake
    Note over Auth: Snowflake 直接实现 Generator Port
    
    Auth->>Auth: NewBcryptPasswordHasher(cost)
    Note over Auth: 创建密码哈希器
    
    Auth->>Auth: NewAuthService(uow, passwordHasher,<br/>tokenServiceAdapter, eventPublisher,<br/>idGenerator, logger)
    Note over Auth: 创建应用服务<br/>所有依赖都是 Port 接口
    
    Auth->>Auth: NewHandler(authSvc, respHandler)
    Note over Auth: 创建 HTTP Handler
    
    Auth->>Auth: NewRoutes(handler, jwtSvc)
    Note over Auth: 创建路由配置
    
    Auth-->>Main: *AuthModule
    deactivate Auth
    
    Main->>User: NewUserModule(infra)
    activate User
    Note over User: 类似流程组装用户模块
    User-->>Main: *UserModule
    deactivate User
    
    Main->>Router: gin.New()
    Main->>Router: 遍历 modules<br/>module.RegisterHTTP(group)
    loop 每个 Module
        Auth->>Router: 注册认证相关路由<br/>POST /auth/register<br/>POST /auth/login<br/>POST /auth/logout<br/>POST /auth/refresh
        User->>Router: 注册用户相关路由<br/>GET /user/profile<br/>PUT /user/profile<br/>PATCH /user/password
    end
    
    Note over Router: 所有路由注册完成<br/>准备启动服务
```

---

## 🎯 依赖注入流程

### 构造函数注入详解

```mermaid
graph TB
    subgraph "Composition Root (main.go)"
        A[创建 Infra 容器]
        A --> B[DB *gorm.DB]
        A --> C[Redis *redis.Client]
        A --> D[Logger *zap.Logger]
        A --> E[Config *config.Config]
        A --> F[Snowflake *idgen.Node]
        
        B & C & D & E & F --> G[调用 NewAuthModule infra]
    end
    
    subgraph "AuthModule 内部"
        G --> H[创建基础设施组件]
        H --> I[JWTService]
        H --> J[DAO Query]
        
        I --> K[创建适配器]
        K --> L[TokenServiceAdapter]
        K --> M[ID Generator Adapter]
        
        L & M --> N[创建应用服务]
        N --> O[AuthService 依赖:<br/>• UnitOfWork<br/>• PasswordHasher<br/>• TokenService Port<br/>• EventPublisher<br/>• Generator Port<br/>• Logger]
    end
    
    subgraph "Application Service"
        O --> P[业务逻辑实现<br/>只依赖 Port 接口<br/>不知道具体实现]
    end
    
    style A fill:#87CEEB
    style G fill:#87CEEB
    style K fill:#FFD700
    style L fill:#90EE90
    style M fill:#90EE90
    style O fill:#90EE90
    style P fill:#90EE90
```

---

## 📦 Repository 适配器模式

### Repository 实现流程

```mermaid
sequenceDiagram
    participant App as Application Service
    participant Port as Repository Port
    participant Adapter as Repository Adapter
    participant DAO as GORM DAO
    participant DB as PostgreSQL

    Note over App: 需要查询用户
    
    App->>Port: FindByID(userID)
    Note over Port: 接口定义在<br/>domain/user/repository.go
    
    Port->>Adapter: FindByID(userID)
    Note over Adapter: 实现 Port 接口<br/>internal/infrastructure/persistence/repository/
    
    Adapter->>DAO: q.User.FindOne(id)
    Note over DAO: GORM 生成的查询方法<br/>SELECT * FROM users WHERE id = ? LIMIT 1
    
    DAO->>DB: SQL Query
    activate DB
    DB-->>DAO: sql.Rows
    deactivate DB
    
    DAO-->>Adapter: *dao.User (或 error)
    
    Adapter->>Adapter: 转换 DAO → Entity
    Note over Adapter: 1. 读取 DAO 字段<br/>2. 创建值对象<br/>3. 构建聚合根<br/>4. 加载关联数据（可选）
    
    Adapter-->>Port: *aggregate.User
    
    Port-->>App: *aggregate.User
    
    Note right of App: 拿到领域对象<br/>继续执行业务逻辑
    
    Note over Adapter: 反向操作同样适用:<br/>Save(user) 时<br/>Entity → DAO → DB
```

---

## 🔄 TokenService 适配器转换流程

### 类型转换详细流程

```mermaid
graph TD
    subgraph "Infrastructure Layer"
        A[JWTService.GenerateTokenPair]
        A --> B[返回 *infra.TokenPair]
        B --> C{infra.TokenPair:<br/>AccessToken string<br/>RefreshToken string<br/>ExpiresAt time.Time}
    end
    
    subgraph "Adapter Layer"
        C --> D[TokenServiceAdapter.GenerateTokenPair]
        D --> E[调用 jwtSvc.GenerateTokenPair]
        E --> F[收到 *infra.TokenPair]
        
        F --> G[类型转换]
        G --> H[创建 ports_auth.TokenPair]
        H --> I{ports_auth.TokenPair:<br/>AccessToken string<br/>RefreshToken string<br/>ExpiresAt time.Time}
    end
    
    subgraph "Application Layer"
        I --> J[返回 *ports_auth.TokenPair]
        J --> K[AuthService 使用 Token]
    end
    
    style B fill:#FFB6C1
    style C fill:#FFB6C1
    style F fill:#FFD700
    style G fill:#90EE90
    style H fill:#90EE90
    style I fill:#90EE90
    style J fill:#90EE90
```

### 详细转换代码流程

```mermaid
sequenceDiagram
    participant App as AuthService
    participant Adapter as TokenServiceAdapter
    participant JWT as JWTService

    App->>Adapter: GenerateTokenPair(userID, username, email)
    
    Adapter->>JWT: GenerateTokenPair(userID, username, email)
    activate JWT
    
    Note over JWT: 1. 创建 JWT Claims<br/>2. 签名 Access Token<br/>3. 生成 Refresh Token<br/>4. 计算过期时间
    
    JWT-->>Adapter: &infra.TokenPair{<br/>  AccessToken: "eyJhbGc...",<br/>  RefreshToken: "dGhpcyBp...",<br/>  ExpiresAt: time.Time{...}<br/>}
    deactivate JWT
    
    Note over Adapter: 类型转换开始
    
    Adapter->>Adapter: 创建新对象<br/>&ports_auth.TokenPair{<br/>  AccessToken: pair.AccessToken,<br/>  RefreshToken: pair.RefreshToken,<br/>  ExpiresAt: pair.ExpiresAt,<br/>}
    
    Adapter-->>Adapter: &ports_auth.TokenPair{...}
    
    Adapter-->>App: &ports_auth.TokenPair{...}
    
    Note over App: 使用 Port 类型的 TokenPair<br/>完全不知道 infra.TokenPair 的存在
    
    Note right of Adapter: ⭐ 关键优势:<br/>1. Application 层解耦<br/>2. 易于 Mock 测试<br/>3. 可以随时替换 JWT 库<br/>4. 类型安全
```

---

## 🌐 HTTP 请求处理流程

### 完整请求链路图

```mermaid
graph TB
    A[Client 请求] --> B[Nginx/负载均衡<br/>可选];
    B --> C[Gin Engine];
    
    C --> D[全局中间件链];
    D --> E[CORS 中间件];
    E --> F[Recovery 中间件<br/>Panic 恢复];
    F --> G[Logger 中间件<br/>请求日志];
    G --> H[RequestID 中间件<br/>追踪 ID];
    
    H --> I{路由匹配};
    I -->|未找到 | J[404 Not Found];
    I -->|找到 | K[认证中间件<br/>JWT 验证];
    
    K --> L{Token 有效？};
    L -->|无效 | M[401 Unauthorized];
    L -->|有效 | N[注入 UserID 到 Context];
    
    N --> O[业务中间件<br/>可选];
    O --> P[限流中间件<br/>可选];
    P --> Q[权限检查<br/>RBAC];
    
    Q --> R{有权限？};
    R -->|无 | S[403 Forbidden];
    R -->|有 | T[HTTP Handler];
    
    T --> U[解析请求体];
    U --> V[参数验证];
    V --> W{验证通过？};
    
    W -->|失败 | X[400 Bad Request<br/>返回验证错误];
    W -->|成功 | Y[调用 Application Service];
    
    Y --> Z[执行事务];
    Z --> AA{成功？};
    
    AA -->|失败 | AB[回滚事务];
    AB --> AC[记录错误日志];
    AC --> AD[返回错误响应];
    
    AA -->|成功 | AE[发布领域事件<br/>异步];
    AE --> AF[构建响应 DTO];
    AF --> AG[返回 JSON 响应];
    
    J --> AH[结束];
    M --> AH;
    S --> AH;
    X --> AH;
    AD --> AH;
    AG --> AH;
    
    style C fill:#87CEEB
    style D fill:#FFD700
    style K fill:#FF6B6B
    style Q fill:#FF6B6B
    style T fill:#90EE90
    style Y fill:#90EE90
    style Z fill:#90EE90
    style AG fill:#90EE90
```

---

## 🗄️ 数据库操作流程

### UnitOfWork 事务管理流程

```mermaid
sequenceDiagram
    participant App as Application Service
    participant UoW as UnitOfWork
    participant Tx as *gorm.DB (Tx)
    participant Repo as Repository
    participant DB as PostgreSQL

    App->>UoW: Transaction(func(ctx context.Context) error)
    activate UoW
    
    UoW->>Tx: db.Begin()
    activate Tx
    Tx-->>UoW: *gorm.DB (transaction)
    
    UoW->>UoW: 将 Tx 注入到 Repositories
    Note over UoW: 所有 Repository 操作<br/>都使用这个 Tx
    
    UoW->>App: 执行传入的函数
    
    App->>Repo: FindByID(id)
    activate Repo
    Repo->>Tx: SELECT ... FOR UPDATE
    Tx-->>Repo: result
    Repo-->>App: entity
    deactivate Repo
    
    App->>Repo: Save(entity)
    activate Repo
    Repo->>Tx: UPDATE ...
    Tx-->>Repo: rows affected
    Repo-->>App: nil
    deactivate Repo
    
    App-->>UoW: nil (成功)
    
    alt 执行成功
        UoW->>Tx: Commit()
        Tx-->>UoW: nil
        UoW-->>App: nil
    else 发生错误
        UoW->>Tx: Rollback()
        Tx-->>UoW: nil
        UoW-->>App: error
    end
    
    deactivate Tx
    deactivate UoW
    
    Note over UoW: 事务隔离级别:<br/>Read Committed (默认)<br/>可配置为 Repeatable Read 或 Serializable
```

---

## 📨 领域事件异步处理流程

### 事件队列处理流程

```mermaid
graph TB
    A[领域事件产生] --> B[EventPublisher.Publish];
    B --> C{序列化事件为 JSON};
    
    C --> D[创建 Asynq Task];
    D --> E[Push 到 Redis Queue];
    E --> F[立即返回 nil];
    
    F --> G[事件入队成功];
    G --> H[主流程继续];
    
    par Asynq Worker 后台处理
        I[Worker 轮询队列] --> J{有新任务？};
        J -->|是 | K[Dequeue Task];
        K --> L[反序列化 JSON];
        L --> M[查找对应 Handler];
        
        M --> N{Handler 存在？};
        N -->|否 | O[记录警告];
        N -->|是 | P[调用 Handler];
        
        P --> Q{处理成功？};
        Q -->|是 | R[Acknowledge Task];
        Q -->|否 | S{重试次数 < 3?};
        
        S -->|是 | T[重新入队<br/>指数退避];
        S -->|否 | U[移入 Dead Letter Queue];
        
        R --> V[保存事件到<br/>domain_events 表];
        V --> W[Task 完成];
        
        O --> W;
        U --> W;
        T --> W;
    end
    
    style E fill:#87CEEB
    style F fill:#90EE90
    style R fill:#90EE90
    style V fill:#90EE90
    style U fill:#FF6B6B
```

---

## 📚 参考文档

- [架构总览](./architecture-overview.md) - 整体架构介绍
- [Ports 模式设计](./ports-pattern-design.md) - Ports 模式详细说明
- [业务流程图](./business-flow-diagrams.md) - 业务流程时序图
- [架构分层详解](./architecture-diagrams-detailed.md) - 分层架构图
