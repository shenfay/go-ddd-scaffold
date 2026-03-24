# 核心业务流程图

## 📝 用户注册流程

### 完整业务流程图

```mermaid
sequenceDiagram
    participant Client
    participant Handler as HTTP Handler
    participant Service as AuthService
    participant UoW as UnitOfWork
    participant Repo as UserRepository
    participant User as User Aggregate
    participant Event as EventPublisher

    Client->>Handler: POST /api/auth/register<br/>{username, email, password}
    activate Handler
    Handler->>Service: RegisterUser(cmd)
    activate Service
    
    Service->>UoW: Transaction(func())
    activate UoW
    
    UoW->>Repo: FindByUsername(username)
    activate Repo
    Repo-->>UoW: nil (不存在)
    deactivate Repo
    
    UoW->>Repo: FindByEmail(email)
    activate Repo
    Repo-->>UoW: nil (不存在)
    deactivate Repo
    
    UoW->>User: NewUser(username, email, password)
    activate User
    Note over User: 1. 验证用户名格式<br/>2. 验证邮箱格式<br/>3. 验证密码强度<br/>4. 哈希密码<br/>5. 生成 Snowflake ID
    User-->>UoW: *User (创建成功)
    deactivate User
    
    UoW->>Repo: Save(user)
    activate Repo
    Note over Repo: INSERT INTO users (...)
    Repo-->>UoW: nil (成功)
    deactivate Repo
    
    UoW-->>Service: COMMIT (事务提交)
    deactivate UoW
    
    Service->>Event: Publish(UserRegisteredEvent)
    activate Event
    Note over Event: 异步发布到 Asynq 队列<br/>不阻塞主流程
    Event-->>Service: 成功入队
    deactivate Event
    
    Service-->>Handler: RegisterResult{UserID, Username, Email}
    deactivate Service
    
    Handler-->>Client: 200 OK<br/>{code: 0, data: {...}}
    deactivate Handler
    
    Note right of Event: Worker 后台处理:<br/>1. 保存到 domain_events 表<br/>2. 发送欢迎邮件<br/>3. 其他后续处理
```

### 注册流程关键决策点

```mermaid
graph TD
    A[开始注册] --> B{检查用户名};
    B -->|已存在 | C[返回错误：<br/>用户名已存在];
    B -->|不存在 | D{检查邮箱};
    D -->|已存在 | E[返回错误：<br/>邮箱已被注册];
    D -->|不存在 | F{验证密码强度};
    F -->|太弱 | G[返回错误：<br/>密码强度不足];
    F -->|符合 | H{验证邮箱格式};
    H -->|无效 | I[返回错误：<br/>邮箱格式错误];
    H -->|有效 | J{验证用户名格式};
    J -->|无效 | K[返回错误：<br/>用户名格式错误];
    J -->|有效 | L[哈希密码 bcrypt];
    L --> M[生成 Snowflake ID];
    M --> N[创建 User 聚合根];
    N --> O[保存到数据库];
    O --> P{保存成功？};
    P -->|失败 | Q[返回错误：<br/>数据库异常];
    P -->|成功 | R[发布领域事件];
    R --> S[返回成功响应];
    
    C --> T[结束];
    E --> T;
    G --> T;
    I --> T;
    K --> T;
    Q --> T;
    S --> T;
    
    style L fill:#90EE90
    style M fill:#90EE90
    style N fill:#90EE90
    style O fill:#90EE90
    style R fill:#87CEEB
    style S fill:#90EE90
    style C fill:#FFB6C1
    style E fill:#FFB6C1
    style G fill:#FFB6C1
    style I fill:#FFB6C1
    style K fill:#FFB6C1
    style Q fill:#FFB6C1
```

---

## 🔑 用户登录流程

### 完整登录时序图

```mermaid
sequenceDiagram
    participant Client
    participant Handler as HTTP Handler
    participant Service as AuthService
    participant UoW as UnitOfWork
    participant Repo as UserRepository
    participant User as User Aggregate
    participant Token as TokenService
    participant Event as EventPublisher

    Client->>Handler: POST /api/auth/login<br/>{identifier, password}
    activate Handler
    Handler->>Service: AuthenticateUser(cmd)
    activate Service
    
    Service->>UoW: Transaction(func())
    activate UoW
    
    UoW->>Repo: FindByEmail/Username(identifier)
    activate Repo
    alt 找到用户
        Repo-->>UoW: *User
    else 未找到
        Repo-->>UoW: error
        Note over UoW: 返回"用户名或密码错误"<br/>不暴露用户是否存在
    end
    deactivate Repo
    
    UoW->>User: CanLogin()
    activate User
    alt 可以登录
        User-->>UoW: true
    else 被禁用
        User-->>UoW: false (Status=Inactive)
        Note over UoW: 返回"账户已被禁用"
    else 被锁定
        User-->>UoW: false (Status=Locked)
        Note over UoW: 返回"账户已被锁定"
    end
    deactivate User
    
    UoW->>Service: VerifyPassword(cmd.password, user.Password)
    alt 密码正确
        Service-->>UoW: true
    else 密码错误
        Service-->>UoW: false
        Note over UoW: 返回"用户名或密码错误"<br/>并记录失败尝试
    end
    
    UoW->>Token: GenerateTokenPair(userID, username, email)
    activate Token
    Note over Token: 1. 生成 JWT Access Token<br/>2. 生成 Refresh Token<br/>3. 设置过期时间
    Token-->>UoW: {AccessToken, RefreshToken, ExpiresAt}
    deactivate Token
    
    UoW->>User: Login(password, ip, userAgent)
    activate User
    Note over User: 1. 更新 LastLoginAt<br/>2. 发布 UserLoggedInEvent
    User-->>UoW: nil
    deactivate User
    
    UoW-->>Service: COMMIT
    deactivate UoW
    
    Service->>Event: PublishAsync(UserLoggedInEvent)
    activate Event
    Note over Event: 异步发布到队列<br/>用于审计日志、登录统计等
    Event-->>Service: 成功入队
    deactivate Event
    
    Service-->>Handler: AuthenticateResult<br/>{AccessToken, RefreshToken, ExpiresIn}
    deactivate Service
    
    Handler-->>Client: 200 OK<br/>{access_token, refresh_token, expires_in}
    deactivate Handler
    
    Note right of Event: Worker 后台处理:<br/>1. 保存到 domain_events 表<br/>2. 记录审计日志<br/>3. 发送登录通知（可选）
```

### 登录流程决策树

```mermaid
graph TD
    A[开始登录] --> B{查找用户};
    B -->|未找到 | C[返回：用户名或密码错误];
    B -->|找到 | D{CanLogin?};
    
    D -->|Status=Inactive| E[返回：账户已被禁用];
    D -->|Status=Locked| F[返回：账户已被锁定];
    D -->|Status=Active| G{验证密码};
    
    G -->|错误 | H[返回：用户名或密码错误<br/>记录失败次数];
    H --> I{失败>=5 次？};
    I -->|是 | J[锁定账户];
    I -->|否 | K[结束];
    
    G -->|正确 | L[生成 Token 对];
    L --> M[更新 LastLoginAt];
    M --> N[重置失败计数];
    N --> O[发布登录事件];
    O --> P[返回成功和 Token];
    P --> K;
    
    J --> K;
    
    style L fill:#90EE90
    style M fill:#90EE90
    style N fill:#90EE90
    style O fill:#87CEEB
    style P fill:#90EE90
    style C fill:#FFB6C1
    style E fill:#FFB6C1
    style F fill:#FFB6C1
    style H fill:#FFB6C1
    style J fill:#FF6B6B
```

---

## 🔄 Token 刷新流程

### Token 刷新时序图

```mermaid
sequenceDiagram
    participant Client
    participant Handler as HTTP Handler
    participant Service as AuthService
    participant Token as TokenService
    participant User as User Aggregate
    participant Repo as UserRepository

    Client->>Handler: POST /api/auth/refresh<br/>{refresh_token, current_token?}
    activate Handler
    Handler->>Service: RefreshToken(cmd)
    activate Service
    
    Service->>Token: ValidateToken(refresh_token)
    activate Token
    alt 令牌有效
        Token-->>Service: TokenClaims{UserID, ...}
    else 令牌无效/过期
        Token-->>Service: error
        Note over Service: 返回"无效的刷新令牌"
    end
    deactivate Token
    
    alt 令牌验证通过
        Service->>Repo: FindByID(UserID)
        activate Repo
        alt 找到用户
            Repo-->>Service: *User
        else 用户不存在
            Repo-->>Service: error
            Note over Service: 返回"用户不存在"
        end
        deactivate Repo
        
        Service->>User: CanLogin()
        activate User
        alt 可以登录
            User-->>Service: true
        else 不能登录
            User-->>Service: false
            Note over Service: 返回相应错误<br/>(禁用/锁定等)
        end
        deactivate User
        
        alt 用户状态正常
            opt 严格模式：提供 current_token
                Service->>Token: ParseAccessToken(current_token)
                activate Token
                Token-->>Service: oldClaims
                deactivate Token
                
                Service->>Token: BlacklistToken(current_token, expiresAt)
                activate Token
                Note over Token: 将旧 token 加入黑名单<br/>防止并发使用
                Token-->>Service: nil
                deactivate Token
            end
            
            Service->>Token: GenerateTokenPair(userID, username, email)
            activate Token
            Note over Token: 令牌轮换策略:<br/>生成新的 AccessToken 和 RefreshToken
            Token-->>Service: newTokenPair
            deactivate Token
            
            Service-->>Handler: RefreshTokenResult<br/>{new_access, new_refresh, expires_in}
        end
    end
    
    deactivate Service
    Handler-->>Client: 200 OK<br/>{access_token, refresh_token, expires_in}
    deactivate Handler
    
    Note right of Service: 安全特性:<br/>1. 验证 RefreshToken 有效性<br/>2. 检查用户状态<br/>3. Token 轮换（可选）<br/>4. 黑名单机制
```

---

## 🚪 用户登出流程

### 登出流程图

```mermaid
graph TD
    A[收到登出请求] --> B{是否提供<br/>access_token?};
    
    B -->|是 | C[解析 access_token];
    C --> D{解析成功？};
    D -->|失败 | E[记录警告日志];
    D -->|成功 | F[获取过期时间];
    F --> G[加入黑名单];
    G --> H{加入成功？};
    H -->|失败 | I[记录错误日志<br/>但不影响登出];
    H -->|成功 | J[继续];
    
    B -->|否 | J;
    E --> J;
    I --> J;
    
    J --> K[返回成功];
    K --> L[结束];
    
    style G fill:#90EE90
    style K fill:#90EE90
    style I fill:#FFD700
    style E fill:#FFD700
```

### 登出时序图

```mermaid
sequenceDiagram
    participant Client
    participant Handler as HTTP Handler
    participant Service as AuthService
    participant Token as TokenService

    Client->>Handler: POST /api/auth/logout<br/>{access_token?}
    activate Handler
    Handler->>Service: Logout(cmd)
    activate Service
    
    opt 提供了 access_token
        Service->>Token: ParseAccessToken(token)
        activate Token
        alt 解析成功
            Token-->>Service: claims{ExpiresAt}
        else 解析失败
            Token-->>Service: error
            Note over Service: 记录警告<br/>但继续执行
        end
        deactivate Token
        
        Service->>Token: BlacklistToken(token, expiresAt)
        activate Token
        Note over Token: Redis SETEX<br/>key: "token_blacklist:{jti}"<br/>ttl: expiresAt - now
        Token-->>Service: nil
        deactivate Token
    end
    
    Service-->>Handler: LogoutResult{Success: true}
    deactivate Service
    
    Handler-->>Client: 200 OK<br/>{code: 0, message: "登出成功"}
    deactivate Handler
    
    Note right of Service: 登出特点:<br/>1. 幂等性 - 重复登出也成功<br/>2. 可选 token - 不提供也能登出<br/>3. 黑名单防止重用<br/>4. 不修改用户状态
```

---

## 👤 用户资料更新流程

### 更新资料时序图

```mermaid
sequenceDiagram
    participant Client
    participant Middleware as Auth Middleware
    participant Handler as HTTP Handler
    participant Service as UserService
    participant UoW as UnitOfWork
    participant Repo as UserRepository
    participant User as User Aggregate
    participant Cache as UserCache

    Client->>Middleware: PUT /api/user/profile<br/>Authorization: Bearer {token}<br/>{display_name, avatar}
    activate Middleware
    Middleware->>Middleware: ValidateToken()
    alt 令牌有效
        Middleware-->>Handler: ctx with userID
    else 令牌无效
        Middleware-->>Client: 401 Unauthorized
        destroy Middleware
    end
    
    Handler->>Service: UpdateProfile(ctx, userID, req)
    activate Service
    
    Service->>UoW: Transaction(func())
    activate UoW
    
    UoW->>Repo: FindByID(userID)
    activate Repo
    Repo-->>UoW: *User
    deactivate Repo
    
    UoW->>User: UpdateProfile(displayName, avatar)
    activate User
    Note over User: 1. 验证 displayName 长度<br/>2. 验证 avatar URL 格式<br/>3. 更新属性<br/>4. 发布 UserProfileUpdatedEvent
    User-->>UoW: nil
    deactivate User
    
    UoW->>Repo: Save(user)
    activate Repo
    Note over Repo: UPDATE users SET ... WHERE id = ?
    Repo-->>UoW: nil
    deactivate Repo
    
    UoW-->>Service: COMMIT
    deactivate UoW
    
    Service->>Cache: Delete(userID)
    activate Cache
    Note over Cache: 清除缓存<br/>保证下次读取最新数据
    Cache-->>Service: nil
    deactivate Cache
    
    Service-->>Handler: nil (成功)
    deactivate Service
    
    Handler-->>Client: 200 OK<br/>{code: 0, message: "更新成功"}
    deactivate Handler
    
    Note right of Service: 关键点:<br/>1. 需要认证<br/>2. 事务保证一致性<br/>3. 清除缓存<br/>4. 发布领域事件
```

---

## 🔐 密码修改流程

### 密码修改流程图

```mermaid
graph TD
    A[收到密码修改请求] --> B[验证当前密码];
    B --> C{密码正确？};
    
    C -->|错误 | D[返回：当前密码错误];
    D --> E[结束];
    
    C -->|正确 | F{验证新密码强度};
    F -->|太弱 | G[返回：密码强度不足<br/>需包含大小写、数字、特殊字符];
    G --> E;
    
    F -->|符合 | H{新旧密码相同？};
    H -->|是 | I[返回：新密码不能与旧密码相同];
    I --> E;
    
    H -->|否 | J[哈希新密码];
    J --> K[更新密码];
    K --> L[记录修改时间];
    L --> M[发布密码修改事件];
    M --> N[使其他设备 token 失效<br/>可选];
    N --> O[返回成功];
    O --> E;
    
    style J fill:#90EE90
    style K fill:#90EE90
    style L fill:#90EE90
    style M fill:#87CEEB
    style N fill:#FFD700
    style O fill:#90EE90
    style D fill:#FFB6C1
    style G fill:#FFB6C1
    style I fill:#FFB6C1
```

---

## 🎯 错误处理流程

### 统一错误处理流程图

```mermaid
graph TD
    A[Controller 捕获异常] --> B{错误类型};
    
    B -->|BusinessError| C[提取错误码和消息];
    C --> D[返回标准响应<br/>{code, message}];
    
    B -->|ValidationError| E[收集验证错误];
    E --> F[返回 400<br/>{errors: [...]}];
    
    B -->|UnauthorizedError| G[返回 401<br/>Unauthorized];
    
    B -->|ForbiddenError| H[返回 403<br/>Forbidden];
    
    B -->|NotFoundError| I[返回 404<br/>Not Found];
    
    B -->|DBError| J[记录详细日志];
    J --> K[返回 500<br/>Internal Server Error];
    
    B -->|Panic| L[恢复 Panic];
    L --> M[记录堆栈跟踪];
    M --> N[返回 500<br/>服务异常];
    
    D --> O[结束];
    F --> O;
    G --> O;
    H --> O;
    I --> O;
    K --> O;
    N --> O;
    
    style C fill:#FFD700
    style D fill:#FF6B6B
    style J fill:#FF6B6B
    style K fill:#FF6B6B
    style M fill:#FF6B6B
    style N fill:#FF6B6B
```

---

## 📊 领域事件处理流程

### 事件发布与订阅流程

```mermaid
sequenceDiagram
    participant App as Application Service
    participant Pub as EventPublisher
    participant Queue as Asynq Queue
    participant Worker as Asynq Worker
    participant Store as EventStore
    participant Handler as DomainEventHandler

    Note over App: 业务操作完成<br/>产生领域事件
    
    App->>Pub: Publish(event)
    activate Pub
    
    pub Note over Pub: 1. 序列化事件为 JSON<br/>2. 创建 Asynq Task<br/>3. Push 到 Redis 队列
    
    Pub->>Queue: Enqueue(task)
    activate Queue
    Queue-->>Pub: task_id
    deactivate Queue
    
    Pub-->>App: nil (立即返回)
    deactivate Pub
    
    Note over App: 继续处理其他业务<br/>不等待事件处理
    
    par 异步并行处理
        Worker->>Queue: Dequeue(task)
        activate Worker
        
        Worker->>Store: Save(event)
        activate Store
        Note over Store: 保存到 domain_events 表<br/>用于事件溯源和审计
        Store-->>Worker: nil
        deactivate Store
        
        Worker->>Handler: Handle(event)
        activate Handler
        Note over Handler: 根据事件类型调用<br/>相应的处理器
        
        alt 有注册处理器
            Handler-->>Worker: nil
        else 无处理器
            Note over Worker: 记录警告日志<br/>但标记为成功
        end
        deactivate Handler
        
        Worker-->>Queue: Acknowledge(task_id)
        deactivate Worker
    end
    
    Note right of Worker: 重试机制:<br/>• 失败自动重试 (最多 3 次)<br/>• 指数退避<br/>• 最终失败记录到 dead letter queue
```

---

## 📚 参考文档

- [架构总览](./architecture-overview.md) - 整体架构介绍
- [Ports 模式设计](./ports-pattern-design.md) - Ports 模式详细说明
- [领域模型可视化](./domain-model-visual.md) - 领域模型图表
- [架构分层详解](./architecture-diagrams-detailed.md) - 分层架构图
