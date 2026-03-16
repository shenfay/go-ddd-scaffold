## 📊 领域模型图集

### User 聚合根结构

```mermaid
classDiagram
    class User {
        -ID: UserID
        -username: UserName
        -email: Email
        -password: HashedPassword
        -status: UserStatus
        -displayName: string
        -createdAt: time.Time
        -updatedAt: time.Time
        +Register(username, email, pwd): User
        +Activate(): void
        +Deactivate(): void
        +ChangePassword(old, new): void
        +RecordLogin(ip, ua): void
        +CanLogin(): bool
    }
    
    class UserName {
        -value: string
        +NewUserName(value): UserName
        +Value(): string
    }
    
    class Email {
        -value: string
        +NewEmail(value): Email
        +Value(): string
    }
    
    class HashedPassword {
        -value: string
        +NewHashedPassword(hash): HashedPassword
        +Value(): string
    }
    
    class UserStatus {
        <<enumeration>>
        Pending
        Active
        Inactive
        Locked
    }
    
    User *-- UserName : contains
    User *-- Email : contains
    User *-- HashedPassword : contains
    User *-- UserStatus : has
    
    note for User "聚合根\n维护业务一致性"
    note for UserName "值对象\n不可变"
    note for Email "值对象\n不可变"
    note for HashedPassword "值对象\n不可变"
```

**说明：**
- **聚合根**：User 是聚合的根，所有外部访问都通过 User
- **值对象**：UserName、Email、HashedPassword 都是不可变的值对象
- **封装性**：外部不能直接修改内部状态，必须通过行为方法

---

### 领域事件关系图

```mermaid
graph LR
    subgraph UserActions[用户行为]
        Register[注册]
        Login[登录]
        Update[更新资料]
        ChangePwd[修改密码]
    end
    
    subgraph Events[领域事件]
        Registered[UserRegisteredEvent]
        LoggedIn[UserLoggedInEvent]
        Updated[UserUpdatedEvent]
        PasswordChanged[UserPasswordChangedEvent]
    end
    
    subgraph SideEffects[副作用]
        WelcomeEmail[发送欢迎邮件]
        InitStats[初始化统计]
        AuditLog[记录审计日志]
        LoginLog[记录登录日志]
        Notify[发送通知]
    end
    
    Register --> Registered
    Login --> LoggedIn
    Update --> Updated
    ChangePwd --> PasswordChanged
    
    Registered --> WelcomeEmail
    Registered --> InitStats
    Registered --> AuditLog
    
    LoggedIn --> LoginLog
    LoggedIn --> AuditLog
    
    Updated --> AuditLog
    Updated --> Notify
    
    PasswordChanged --> AuditLog
    PasswordChanged --> Notify
    
    style Events fill:#e1f5ff
    style SideEffects fill:#fff4e1
```

**说明：**
- 领域事件由聚合根的行为触发
- 事件处理器监听事件并执行副作用
- 副作用不影响主流程，异步执行
