# 架构优化方案

## 📋 现状分析

### 当前架构分层

```
backend/internal/
├── application/          # 应用层
│   ├── tenant/
│   │   ├── dto/         ✅ 已规范
│   │   └── service/     ⚠️ 职责混合
│   └── user/
│       ├── assembler/   ⚠️ 使用率低
│       ├── dto/         ✅ 已规范
│       ├── repository/  ❌ 空目录，应删除
│       └── service/     ⚠️ CQRS 已实现但命名混乱
├── domain/              # 领域层
│   ├── tenant/
│   │   ├── aggregate/   ⚠️ 与 entity 职责重复
│   │   ├── entity/      ⚠️ TenantMember 定位不清
│   │   ├── event/       ✅ 已实现
│   │   ├── factory/     ✅ 已实现
│   │   ├── repository/  ✅ 接口清晰
│   │   ├── service/     ⚠️ 领域服务职责不明
│   │   └── valueobject/ ❌ 空目录，应删除或补充
│   ├── user/
│   │   ├── aggregate/   ❌ 缺失，User 应该是聚合根
│   │   ├── entity/      ⚠️ User 是实体不是聚合根
│   │   ├── event/       ✅ 已实现
│   │   ├── repository/  ✅ 接口清晰
│   │   ├── service/     ✅ 已实现
│   │   └── valueobject/ ✅ 已实现
│   └── shared/          ✅ 共享领域概念
└── infrastructure/      # 基础设施层
    ├── app/             ⚠️ 命名模糊
    ├── auth/            ✅ Casbin + JWT
    ├── cache/           ❌ 空目录
    ├── errors/          ❌ 空目录
    ├── event/           ⚠️ 事件总线实现分散
    ├── external/        ❌ 空目录
    ├── middleware/      ✅ HTTP 中间件
    ├── notification/    ❌ 空目录
    ├── persistence/     ✅ GORM 仓储实现
    ├── queue/           ❌ 空目录
    ├── server/          ⚠️ 手动依赖注入
    └── wire/            ⚠️ Wire 配置不完整
```

---

## 🔴 核心问题清单

### P0: 领域层设计问题

#### 1. TenantMember 的聚合边界模糊

**问题描述**：
- `TenantMember` 被定义为 `Tenant` 的子实体（第 20 行注释）
- 但有独立的 Repository 接口和实现
- 有独立的生命周期（可独立创建、删除）
- 数据库中有独立的表 `tenant_members`

**违反原则**：
- DDD 聚合根边界不清晰
- 一致性边界定义错误

**解决方案**：

✅ **方案 A（推荐）：明确为独立聚合根**
```go
// domain/tenant/entity/tenant_member.go
// TenantMember 作为独立聚合根
type TenantMember struct {
    ID        uuid.UUID
   TenantID  uuid.UUID  // 外键引用
    UserID    uuid.UUID  // 外键引用
    Role      UserRole
    Status    MemberStatus
    // ... 其他字段
}

// 在 aggregate/ 目录下定义聚合行为
// domain/tenant/aggregate/tenant_member_aggregate.go
type TenantMemberAggregate struct {
    *entity.TenantMember
    // 聚合业务逻辑
}
```

❌ **方案 B（不推荐）：真正作为子实体**
- 需要重构为值对象集合存储在 Tenant 中
- 不符合实际业务场景，增加复杂度

---

#### 2. User 缺少聚合根定义

**问题描述**：
- User 只有 entity 目录，没有 aggregate 目录
- 但 User 有明显的不变量需要维护（Email 唯一、密码加密等）

**解决方案**：
```go
// domain/user/aggregate/user_aggregate.go
type UserAggregate struct {
    *entity.User
    // 聚合业务逻辑
}

func (a *UserAggregate) Register(email, password, nickname string) error {
    // 1. 验证 Email 格式
    // 2. 加密密码
    // 3. 设置初始状态
    // 4. 记录领域事件
}
```

---

#### 3. 空目录未清理

**应删除的目录**：
- `domain/tenant/valueobject/` - 空目录
- `application/user/repository/` - 空目录，Repository 应该在 Domain 层
- `infrastructure/cache/` - 空目录
- `infrastructure/errors/` - 空目录
- `infrastructure/external/` - 空目录
- `infrastructure/notification/` - 空目录
- `infrastructure/queue/` - 空目录

**应补充的值对象**：
```go
// domain/tenant/valueobject/tenant_name.go
type TenantName struct {
    value string
}

// domain/tenant/valueobject/description.go
type Description struct {
    value string
}
```

---

### P1: 应用层职责问题

#### 4. CQRS 结构不清晰

**当前状态**：
```
application/user/service/
├── authentication_service.go   # 认证服务
├── user_command_service.go     # 命令服务
├── user_query_service.go       # 查询服务
├── interfaces.go               # 接口定义
└── transactional_auth_service_example.go  # 示例代码
```

**问题**：
- 文件名包含 `service` 后缀，不符合 Go 惯例
- CQRS 已实现但没有明确的目录结构
- `assembler/` 目录使用率低

**优化方案**：

✅ **推荐方案：按功能模块组织**
```
application/user/
├── auth/
│   ├── service.go             # AuthenticationService
│   ├── register.go            # 注册实现
│   ├── login.go               # 登录实现
│   └── logout.go              # 登出实现
├── command/
│   ├── service.go             # UserCommandService
│   ├── update_user.go         # 更新用户
│   └── update_profile.go      # 更新资料
├── query/
│   ├── service.go             # UserQueryService
│   ├── get_user.go            # 获取用户
│   └── list_users.go          # 列表查询
└── dto/                        # DTO 定义
```

**优势**：
- 职责清晰，易于导航
- 符合企业级应用规范
- 不过度设计，保持简洁

---

#### 5. Application Service 调用 Domain Service

**问题代码**：
```go
// application/tenant/service/tenant_service.go
func (s *tenantService) CreateTenant(...) {
    // 直接调用 Factory
    tenant, _ := factory.NewTenantBuilder(...).Build()
    
    // 直接操作 Repository
    tenantRepo.Create(ctx, tenant)
    memberRepo.Create(ctx, member)
}
```

**违反原则**：
- Application Service 不应该直接调用 Domain Factory
- 应该通过 Domain Service 或 Aggregate 方法

**正确做法**：
```go
// domain/tenant/service/tenant_domain_service.go
type TenantDomainService interface {
    CreateTenantWithOwner(ctx context.Context, name string, ownerID uuid.UUID) (*entity.Tenant, error)
}

// application/tenant/service/tenant_app_service.go
func (s *tenantAppService) CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.TenantResponse, error) {
    // 调用领域服务
    tenant, err := s.domainService.CreateTenantWithOwner(ctx, req.Name, ownerID)
    if err != nil {
        return nil, err
    }
    
    // 转换为 DTO
    return dto.ToTenantDTO(tenant, 0), nil
}
```

---

### P2: 基础设施层问题

#### 6. 依赖注入混乱

**当前状态**：
- `infrastructure/wire/` 有部分 Wire 配置
- `infrastructure/server/service.go` 手动依赖注入
- 两种方式混用

**优化方案**：

✅ **统一使用 Wire**
```go
// infrastructure/wire/wire.go
//go:build wireinject
// +build wireinject

func InitializeApp() (*App, error) {
    panic(wire.Build(
        // 基础设施
        NewDB,
        NewLogger,
        NewRedis,
        
        // 领域模块
       userModuleSet,
        tenantModuleSet,
        
        // HTTP 服务器
        NewServer,
    ))
}

// infrastructure/wire/modules.go
var userModuleSet = wire.NewSet(
   user_repo.NewUserRepository,
   user_service.NewAuthService,
   user_service.NewCommandService,
   user_service.NewQueryService,
)

var tenantModuleSet = wire.NewSet(
    tenant_repo.NewTenantRepository,
    tenant_repo.NewTenantMemberRepository,
    tenant_service.NewDomainService,
    tenant_app_service.NewApplicationService,
)
```

---

#### 7. 事件总线架构不完善

**当前状态**：
- 领域事件定义在 `domain/*/event/`
- 事件总线实现在 `infrastructure/event/`
- 缺少统一的事件处理器注册机制

**优化方案**：
```go
// domain/shared/event/dispatcher.go
type EventDispatcher interface {
    Publish(events ...DomainEvent) error
    Subscribe(eventType string, handler EventHandler)
}

// infrastructure/event/in_memory_dispatcher.go
type InMemoryDispatcher struct {
    handlers map[string][]EventHandler
}

func (d *InMemoryDispatcher) Publish(events ...DomainEvent) error {
    for _, event := range events {
        if handlers, ok := d.handlers[event.GetEventType()]; ok {
            for _, handler := range handlers {
               go handler.Handle(event)
            }
        }
    }
    return nil
}

// application/event/handlers/tenant_created_handler.go
type TenantCreatedHandler struct {
    // 初始化默认配置
    // 发送欢迎邮件
}

func (h *TenantCreatedHandler) Handle(event DomainEvent) {
    // 处理租户创建事件
}
```

---

### P3: 代码规范问题

#### 8. 命名不一致

**问题**：
- Service 文件：`authentication_service.go` vs `user_command_service.go`
- 有的用下划线，有的用中横线
- 接口命名：`TenantService` vs `UserQueryService`

**规范**：
```
✅ 推荐：小写字母 + 下划线
- authentication_service.go
- user_command_service.go
- tenant_member_repository.go

✅ 或者：纯小写（更 Go 风格）
- authservice/service.go
- commandservice/user.go
- repository/tenant.go
```

---

#### 9. 示例代码位置不当

**问题**：
- `transactional_auth_service_example.go` 在生产代码目录

**优化**：
```
docs/examples/
└── transactional-auth-example.go  # 示例代码移到文档目录

// 或者在测试目录
internal/tests/examples/
└── transactional_auth_test.go
```

---

## 📊 优化优先级

### 第一阶段（P0 - 核心架构）
1. ✅ 明确 TenantMember 为独立聚合根
2. ✅ 补充 User 聚合根定义
3. ✅ 清理空目录

**工作量**：2-3 天  
**风险**：低  
**价值**：⭐⭐⭐⭐⭐

---

### 第二阶段（P1 - 应用层规范）
4. ✅ 重构 Application 层为清晰 CQRS 结构
5. ✅ 分离 Domain Service 和 Application Service

**工作量**：3-4 天  
**风险**：中  
**价值**：⭐⭐⭐⭐

---

### 第三阶段（P2 - 基础设施）
6. ✅ 统一使用 Wire 进行依赖注入
7. ✅ 完善事件驱动架构

**工作量**：2-3 天  
**风险**：低  
**价值**：⭐⭐⭐

---

### 第四阶段（P3 - 代码规范）
8. ✅ 统一命名规范
9. ✅ 迁移示例代码

**工作量**：1 天  
**风险**：低  
**价值**：⭐⭐

---

## 🎯 实施建议

### 立即执行（今天）
1. 清理所有空目录
2. 确认 TenantMember 聚合根定位
3. 创建 User 聚合根目录结构

### 本周完成
4. 重构 Application 层目录结构
5. 统一 Service 命名
6. 迁移示例代码

### 下周完成
7. 统一依赖注入方案
8. 完善事件总线
9. 编写架构文档

---

## 📝 决策点

需要您确认的设计决策：

1. **TenantMember 是否作为独立聚合根？**
   - ✅ 推荐：是（符合实际业务）
   - 否：需要大幅重构为值对象

2. **Application 层采用哪种组织方式？**
   - ✅ 推荐：按功能模块（auth/, command/, query/）
   - 备选：保持现状，仅优化命名

3. **依赖注入统一方案？**
   - ✅ 推荐：全部使用 Wire
   - 备选：保留手动注入（简单场景）

4. **是否需要完整的值对象层？**
   - ✅ 推荐：为核心概念添加值对象（TenantName, Email 等）
   - 备选：仅在必要时使用

---

## ✅ 下一步行动

请确认以上方案，我将：
1. 先执行第一阶段的 P0 优化
2. 每完成一个阶段提交一次 Git
3. 确保编译和测试通过

**您希望我现在开始执行第一阶段（P0）的优化吗？**
