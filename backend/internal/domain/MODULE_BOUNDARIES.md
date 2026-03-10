# Domain Module Boundaries

## 领域模块边界定义

### 1. User Domain (用户领域)

**职责**: 管理用户身份、认证、个人资料

**聚合根**:
- `User` - 用户聚合根

**实体**:
- 无（User 是唯一实体）

**值对象**:
- `UserRole` - 用户角色
- `UserStatus` - 用户状态
- `Email` - 邮箱地址
- `Password` - 密码

**领域事件**:
- `UserRegisteredEvent` - 用户注册事件
- `UserLoggedInEvent` - 用户登录事件
- `UserLockedEvent` - 用户锁定事件
- `UserActivatedEvent` - 用户激活事件
- `UserProfileUpdatedEvent` - 用户资料更新事件
- `UserEmailChangedEvent` - 用户邮箱变更事件

**Repository**:
- `UserRepository` - 用户仓储

**Domain Service**:
- `UserDomainService` - 用户领域服务

---

### 2. Tenant Domain (租户领域)

**职责**: 管理租户组织、成员关系、角色权限

**聚合根**:
- `Tenant` - 租户聚合根（包含 Members 集合）

**实体**:
- `TenantMember` - 租户成员（属于 Tenant 聚合）

**值对象**:
- `MemberRole` - 成员角色（owner, admin, member, guest）
- `MemberStatus` - 成员状态（active, inactive, removed）

**领域事件**:
- `TenantCreatedEvent` - 租户创建事件
- `MemberJoinedEvent` - 成员加入事件
- `MemberLeftEvent` - 成员离开事件
- `RoleChangedEvent` - 角色变更事件

**Repository**:
- `TenantRepository` - 租户仓储
- `TenantMemberRepository` - 租户成员仓储

**Domain Service**:
- `MembershipDomainService` - 成员资格领域服务

**Factory**:
- `TenantBuilder` - 租户构建器

---

### 3. Shared Domain (共享领域)

**职责**: 跨领域的通用领域概念

**模块**:
- `specification/` - Specification 模式实现
  - `Specification[T]` - 泛型规格接口
  - `AndSpec`, `OrSpec`, `NotSpec` - 组合规格
  - 业务规则封装（ActiveMemberSpec, OwnerRoleSpec 等）

---

## 模块依赖关系

```
┌─────────────────┐
│  User Domain    │
│                 │
│  - User        │
│  - UserRole     │
│  - UserStatus   │
└────────┬────────┘
         │
         │ 引用
         │
┌────────▼────────┐
│  Tenant Domain  │
│                 │
│  - Tenant       │
│  - TenantMember │
│  - MemberRole   │
│  - MemberStatus │
└─────────────────┘
```

**依赖方向**: Tenant Domain → User Domain

**解释**:
- TenantMember 引用 UserID（租户成员关联到用户）
- User 不感知 Tenant 的存在（用户可以是多个租户的成员）
- 这是正确的依赖方向，符合单一职责原则

---

## 不变量保证

### User Domain
1. Email 全局唯一
2. Password 必须加密存储
3. Status 必须是有效枚举值
4. Role 必须在允许范围内

### Tenant Domain
1. Name 在租户内唯一
2. MaxMembers > 0
3. ExpiredAt > Now
4. Members 数量 <= MaxMembers
5. 每个用户在每个租户只能有一个成员记录
6. Owner 角色唯一

---

## 跨模块操作

### 场景 1: 用户加入租户
```go
// Application Service
func (s *tenantService) AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
    return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 1. 获取租户聚合根
        tenant, err := tenantRepo.GetByID(ctx, tenantID)
        if err != nil {
            return err
        }
        
        // 2. 调用聚合根方法（检查不变量）
        member, err := tenant.AddMember(userID, role, currentUserID)
        if err != nil {
            return err
        }
        
        // 3. 持久化
        if err := tenantRepo.Save(ctx, tenant); err != nil {
            return err
        }
        
        if err := memberRepo.Create(ctx, member); err != nil {
            return err
        }
        
        // 4. 发布领域事件
        eventBus.Publish(NewMemberJoinedEvent(tenantID, userID, role))
        
        return nil
    })
}
```

### 场景 2: 租户创建（自动添加创始成员）
```go
func (s *tenantService) CreateTenant(ctx context.Context, name, description string, ownerID uuid.UUID) (*entity.Tenant, error) {
    return s.uow.WithTransaction(ctx, func(ctx context.Context) error {
        // 使用 Factory 构建租户
        tenant, err := NewTenantBuilder(ownerID, name).
            WithDescription(description).
            Build()
        if err != nil {
            return err
        }
        
        // 保存到数据库
        if err := tenantRepo.Create(ctx, tenant); err != nil {
            return err
        }
        
        // 自动添加创始成员（owner 角色）
        member := &entity.TenantMember{
            ID:      uuid.New(),
           UserID:   ownerID,
           TenantID: tenant.ID,
            Role:     entity.RoleOwner,
            Status:   entity.MemberStatusActive,
            JoinedAt: time.Now(),
        }
        
        if err := memberRepo.Create(ctx, member); err != nil {
            return err
        }
        
        return nil
    })
}
```

---

## Specification 模式应用

### 示例 1: 验证活跃的所有者
```go
spec := And(&ActiveMemberSpec{}, &OwnerRoleSpec{})
if spec.IsSatisfiedBy(member) {
    // 执行所有者专属操作
}
```

### 示例 2: 租户容量检查
```go
capacitySpec := NewTenantHasCapacitySpec(currentMemberCount)
if !capacitySpec.IsSatisfiedBy(tenant) {
    return ErrTenantCapacityExceeded
}
```

### 示例 3: 复杂规则组合
```go
// 可以加入的管理员：活跃状态 + (管理员或所有者角色)
eligibleSpec := And(
    &ActiveMemberSpec{},
    Or(&AdminRoleSpec{}, &OwnerRoleSpec{}),
)
```

---

## CQRS DTO 映射

### Command Side (写模型)
- `CreateTenantRequest` - 创建租户请求
- `UpdateTenantRequest` - 更新租户请求
- `AddMemberRequest` - 添加成员请求
- `RemoveMemberRequest` - 移除成员请求

### Query Side (读模型)
- `TenantResponse` - 租户基础信息
- `TenantDetailResponse` - 租户详情（含成员列表）
- `TenantMemberResponse` - 成员信息（嵌套 User DTO）

---

## 测试策略

### Unit Test
- 测试聚合根的不变量保证
- 测试 Specification 的规则逻辑
- 测试 Factory 的构建逻辑

### Integration Test
- 测试 Repository 的持久化
- 测试 Domain Service 的跨聚合操作
- 测试领域事件的发布

---

## 未来扩展

### 可能的子领域
1. **Permission Domain** - 权限管理（RBAC/ABAC）
2. **Billing Domain** - 计费订阅
3. **Audit Domain** - 审计日志

### 扩展点
- Specification 模式支持更多业务规则
- 领域事件支持 Event Sourcing
- CQRS 读写分离支持独立扩展
