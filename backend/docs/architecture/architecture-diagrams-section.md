## 📊 架构图集

### Clean Architecture 分层架构

```mermaid
graph TB
    subgraph Interfaces["接口层 (Interfaces)"]
        HTTP[HTTP Handler<br/>Gin Framework]
        Middleware[Middleware<br/>Auth/CORS/Logger]
        DTOs[Request/Response DTOs]
    end
    
    subgraph Application["应用层 (Application)"]
        UserService[UserService<br/>Register/Auth/GetUser]
        Commands[Commands<br/>RegisterUserCmd<br/>AuthenticateUserCmd]
        Events[Event Handlers<br/>UserRegistered<br/>UserLoggedIn]
    end
    
    subgraph Domain["领域层 (Domain) ⭐"]
        UserAggregate[User Aggregate<br/>聚合根]
        VOs[值对象<br/>UserName<br/>Email<br/>HashedPassword]
        DomainEvents[领域事件<br/>UserRegisteredEvent<br/>UserLoggedInEvent]
        RepoInterface[Repository Interface]
    end
    
    subgraph Infrastructure["基础设施层 (Infrastructure)"]
        RepoImpl[Repository Impl<br/>GORM + PostgreSQL]
        JWT[JWT Service]
        Bcrypt[Bcrypt Hasher]
        EventBus[Event Publisher]
    end
    
    Interfaces --> Application
    Application --> Domain
    Infrastructure -.-> Domain
    
    style Domain fill:#e1f5ff
    style UserAggregate fill:#ffe1e1
```

**说明：**
- **依赖方向**：外层依赖内层，内层不依赖外层
- **领域层**：核心业务逻辑，不依赖任何框架
- **接口层**：处理 HTTP 请求，转换为应用层命令
- **基础设施层**：实现技术细节（数据库、缓存等）

---

### Composition Root 设计

```mermaid
graph TB
    Bootstrap[Bootstrap<br/>Composition Root]
    
    subgraph DomainComponents[领域组件]
        UserHandlers[user.*Handler]
        TenantHandlers[tenant.*Handler]
    end
    
    subgraph InfraComponents[基础设施组件]
        DB[PostgreSQL DB]
        Redis[Redis Cache]
        Logger[Zap Logger]
    end
    
    subgraph HttpComponents[接口层组件]
        HttpHandlers[HTTP Handlers]
        Router[Gin Router]
    end
    
    Bootstrap --> DomainComponents
    Bootstrap --> InfraComponents
    Bootstrap --> HttpComponents
    
    style Bootstrap fill:#ffe1e1
    style DomainComponents fill:#e1f5ff