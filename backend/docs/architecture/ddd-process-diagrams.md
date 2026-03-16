## 📊 核心业务流程图集

### 用户注册流程（时序图）

```mermaid
sequenceDiagram
    participant Client as Client
    participant HTTP as HTTP Handler
    participant Service as UserService
    participant Domain as User Aggregate
    participant Repo as UserRepository
    participant Event as EventPublisher
    participant JWT as JWTService
    
    Client->>HTTP: POST /auth/register
    Note over Client,HTTP: {username, email, password}
    
    HTTP->>Service: RegisterUser(cmd)
    
    Service->>Repo: FindByUsername(username)
    Repo-->>Service: nil (不存在)
    
    Service->>Repo: FindByEmail(email)
    Repo-->>Service: nil (不存在)
    
    Service->>Domain: NewUser(username, email, hash)
    Note over Domain: 创建聚合根<br/>发布 UserRegisteredEvent
    
    Domain-->>Service: User
    
    Service->>Repo: Save(user)
    Note over Repo: GORM 持久化
    
    Repo-->>Service: OK
    
    Service->>Event: Publish(UserRegisteredEvent)
    Note over Event: 异步处理：<br/>1. 发送欢迎邮件<br/>2. 初始化统计<br/>3. 记录审计日志
    
    Service->>JWT: GenerateTokenPair(userID)
    JWT-->>Service: {access, refresh}
    
    Service-->>HTTP: User + Token
    
    HTTP-->>Client: 200 OK + UserDTO + Tokens
```

**关键点：**
1. ✅ 唯一性检查（用户名、邮箱）
2. ✅ 密码哈希（Bcrypt）
3. ✅ 领域事件发布
4. ✅ JWT 令牌生成（注册即登录）

---

### 用户登录流程（时序图）

```mermaid
sequenceDiagram
    participant Client as Client
    participant HTTP as HTTP Handler
    participant Service as UserService
    participant Domain as User Aggregate
    participant Repo as UserRepository
    participant Event as EventPublisher
    participant JWT as JWTService
    
    Client->>HTTP: POST /auth/login
    Note over Client,HTTP: {username, password}
    
    HTTP->>Service: AuthenticateUser(cmd)
    
    Service->>Repo: FindByUsername(username)
    Repo-->>Service: User
    
    Service->>Domain: Verify(password, hash)
    Domain-->>Service: true/false
    
    alt 密码错误
        Domain->>Domain: RecordFailedLogin()
        Domain-->>Service: Error
        Service-->>HTTP: INVALID_PASSWORD
        HTTP-->>Client: 401 Unauthorized
    else 密码正确
        Service->>Domain: CanLogin()
        Domain-->>Service: true
        
        Service->>Domain: RecordLogin(ip, ua)
        Note over Domain: 发布 UserLoggedInEvent
        
        Service->>Repo: Save(user)
        Note over Repo: 更新最后登录时间
        
        Service->>JWT: GenerateTokenPair(userID)
        JWT-->>Service: {access, refresh}
        
        Service-->>HTTP: AuthenticationResult
        
        HTTP-->>Client: 200 OK + Tokens
    end
```

**关键点：**
1. ✅ 密码验证（Bcrypt compare）
2. ✅ 登录失败计数
3. ✅ 领域事件发布
4. ✅ 登录日志记录

---

### 获取用户信息流程（流程图）

```mermaid
sequenceDiagram
    participant Client as Client
    participant HTTP as HTTP Handler
    participant Service as UserService
    participant Repo as UserRepository
    participant Domain as User Aggregate
    
    Client->>HTTP: GET /users/:id
    Note over Client,HTTP: Authorization: Bearer token
    
    HTTP->>HTTP: Parse JWT Token
    Note over HTTP: 解析出当前用户 ID
    
    HTTP->>Service: GetUserByID(targetUserID)
    
    Service->>Repo: FindByID(id)
    Repo-->>Service: User
    
    Service-->>HTTP: User
    
    HTTP->>HTTP: ToUserDetailDTO(user)
    Note over HTTP: 领域对象 → DTO
    
    HTTP-->>Client: 200 OK + UserDetailDTO
```

**关键点：**
1. ✅ JWT Token 解析
2. ✅ 权限验证（只能查看自己的信息）
3. ✅ 领域对象转换为 DTO

---

### Command 侧数据流（写操作）

```mermaid
flowchart TD
    A[HTTP Request] --> B[HTTP Handler]
    B --> C[Parse Request]
    C --> D[Create Command]
    D --> E[Application Service]
    E --> F{Find Aggregate}
    F -->|Not Found| G[Return Error]
    F -->|Found| H[Invoke Domain Method]
    H --> I[Aggregate State Changed]
    I --> J[Publish Domain Event]
    J --> K[Save Aggregate]
    K --> L[Return Result]
    L --> M[HTTP Response]
    
    style H fill:#ffe1e1
    style I fill:#ffe1e1
    style J fill:#fff4e1
```

---

### Query 侧数据流（读操作）

```mermaid
flowchart TD
    A[HTTP Request] --> B[HTTP Handler]
    B --> C[Parse JWT Token]
    C --> D[Extract UserID]
    D --> E[Application Service]
    E --> F[Repository.FindByID]
    F --> G{User Found?}
    G -->|No| H[Return Error]
    G -->|Yes| I[Map to DTO]
    I --> J[HTTP Response]
    
    style F fill:#e1f5ff
    style I fill:#e8f5e9
```

---

### Bootstrap 依赖注入组装流程

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Boot as Bootstrap
    participant Infra as Infrastructure
    participant App as Application
    participant Domain as Domain
    
    Main->>Boot: NewBootstrap(config)
    
    Boot->>Infra: Create DB Connection
    Boot->>Infra: Create Redis Client
    Boot->>Infra: Create Logger
    Boot->>Infra: Create JWT Service
    
    Boot->>Domain: Create Repository Interface
    Note over Boot,Domain: 实际是 RepositoryImpl
    
    Boot->>App: Create Application Service
    Note over Boot,App: 注入 Repository<br/>EventPublisher<br/>PasswordHasher
    
    Boot->>Boot: Create Event Handlers
    Boot->>Boot: Create HTTP Handlers
    
    Boot->>Boot: Setup Routes
    Note over Boot,Boot: 绑定 Handler 到路由
    
    Boot-->>Main: Ready to Serve
```
