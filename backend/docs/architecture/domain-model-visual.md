# 领域模型可视化

## 👤 User 聚合根完整模型

### 聚合根结构图

```
┌─────────────────────────────────────────────────────────────┐
│                      User Aggregate                         │
│                    (用户聚合根)                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │              User Entity (实体)                        │ │
│  │                                                       │ │
│  │  - ID: UserID (值对象)                                │ │
│  │  - Username: Username (值对象)                        │ │
│  │  - Email: Email (值对象)                              │ │
│  │  - Password: Password (值对象)                        │ │
│  │  - DisplayName: string                                │ │
│  │  - Avatar: string                                     │ │
│  │  - Status: UserStatus (枚举)                          │ │
│  │  - LastLoginAt: *time.Time                            │ │
│  │  - CreatedAt: time.Time                               │ │
│  │  - UpdatedAt: time.Time                               │ │
│  │                                                       │ │
│  │  + Register(username, email, password) error          │ │
│  │  + Login(password, ip, userAgent) error               │ │
│  │  + Logout()                                           │ │
│  │  + UpdateProfile(dispName, avatar) error              │ │
│  │  + ChangePassword(oldPwd, newPwd) error               │ │
│  │  + Activate() error                                   │ │
│  │  + Deactivate() error                                 │ │
│  │  + Lock(reason string) error                          │ │
│  │  + Unlock() error                                     │ │
│  │                                                       │ │
│  │  + CanLogin() bool                                    │ │
│  │  + IsActive() bool                                    │ │
│  │  + IsLocked() bool                                    │ │
│  │                                                       │ │
│  │  // 领域事件管理                                       │ │
│  │  + GetUncommittedEvents() []DomainEvent               │ │
│  │  + ClearUncommittedEvents()                           │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │           Value Objects (值对象集合)                   │ │
│  │                                                       │ │
│  │  ┌──────────────┐  ┌──────────────┐                  │ │
│  │  │   UserID     │  │   Username   │                  │ │
│  │  │  (int64)     │  │  (string)    │                  │ │
│  │  │              │  │              │                  │ │
│  │  │ + Validate() │  │ + Validate() │                  │ │
│  │  │ + Int64()    │  │ + Value()    │                  │ │
│  │  └──────────────┘  └──────────────┘                  │ │
│  │                                                       │ │
│  │  ┌──────────────┐  ┌──────────────┐                  │ │
│  │  │    Email     │  │   Password   │                  │ │
│  │  │  (string)    │  │  (string)    │                  │ │
│  │  │              │  │              │                  │ │
│  │  │ + Validate() │  │ + Validate() │                  │ │
│  │  │ + Value()    │  │ + Hash()     │                  │ │
│  │  │ + Matches()  │  │ + Verify()   │                  │ │
│  │  └──────────────┘  └──────────────┘                  │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         Domain Events (领域事件)                       │ │
│  │                                                       │ │
│  │  • UserRegisteredEvent (用户注册)                      │ │
│  │  • UserLoggedInEvent (用户登录)                        │ │
│  │  • UserLoggedOutEvent (用户登出)                       │ │
│  │  • UserProfileUpdatedEvent (资料更新)                  │ │
│  │  • UserPasswordChangedEvent (密码修改)                 │ │
│  │  • UserActivatedEvent (账户激活)                       │ │
│  │  • UserDeactivatedEvent (账户停用)                     │ │
│  │  • UserLockedEvent (账户锁定)                          │ │
│  │  • UserUnlockedEvent (账户解锁)                        │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏢 Tenant 聚合根模型

```
┌─────────────────────────────────────────────────────────────┐
│                    Tenant Aggregate                         │
│                   (租户聚合根)                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │              Tenant Entity (实体)                      │ │
│  │                                                       │ │
│  │  - ID: TenantID (值对象)                              │ │
│  │  - Name: string                                       │ │
│  │  - Slug: string (唯一标识)                            │ │
│  │  - Status: TenantStatus                               │ │
│  │  - Config: *TenantConfig (关联对象)                   │ │
│  │  - CreatedAt: time.Time                               │ │
│  │  - UpdatedAt: time.Time                               │ │
│  │                                                       │ │
│  │  + Create(name, slug) error                           │ │
│  │  + UpdateName(name) error                             │ │
│  │  + Activate() error                                   │ │
│  │  + Deactivate() error                                 │ │
│  │  + AddMember(userID, role) error                      │ │
│  │  + RemoveMember(userID) error                         │ │
│  │  + UpdateMemberRole(userID, role) error               │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │        TenantConfig (关联对象 - 非聚合根)              │ │
│  │                                                       │ │
│  │  - Theme: string                                      │ │
│  │  - Language: string                                   │ │
│  │  - Timezone: string                                   │ │
│  │  - MaxMembers: int                                    │ │
│  │  - Features: map[string]bool                          │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         Domain Events                                  │ │
│  │                                                       │ │
│  │  • TenantCreatedEvent                                  │ │
│  │  • TenantUpdatedEvent                                  │ │
│  │  • TenantActivatedEvent                                │ │
│  │  • TenantDeactivatedEvent                              │ │
│  │  • TenantMemberAddedEvent                              │ │
│  │  • TenantMemberRemovedEvent                            │ │
│  │  • TenantMemberRoleUpdatedEvent                        │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔐 Role & Permission 模型

### 角色聚合根

```
┌─────────────────────────────────────────────────────────────┐
│                     Role Aggregate                          │
│                    (角色聚合根)                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │               Role Entity (实体)                       │ │
│  │                                                       │ │
│  │  - ID: RoleID                                         │ │
│  │  - Name: string (如：admin, editor, viewer)            │ │
│  │  - DisplayName: string                                │ │
│  │  - Description: string                                │ │
│  │  - TenantID: TenantID (所属租户)                       │ │
│  │  - IsSystem: bool (系统角色不可删除)                   │ │
│  │  - Permissions: []PermissionID                        │ │
│  │                                                       │ │
│  │  + Create(name, displayName) error                    │ │
│  │  + AddPermission(permissionID) error                  │ │
│  │  + RemovePermission(permissionID) error               │ │
│  │  + HasPermission(permissionID) bool                   │ │
│  │  + HasAnyPermissions(permissionIDs) bool              │ │
│  │  + HasAllPermissions(permissionIDs) bool              │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         Domain Events                                  │ │
│  │                                                       │ │
│  │  • RoleCreatedEvent                                    │ │
│  │  • RoleUpdatedEvent                                    │ │
│  │  • RoleDeletedEvent                                    │ │
│  │  • PermissionAddedToRoleEvent                          │ │
│  │  • PermissionRemovedFromRoleEvent                      │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 权限值对象

```
┌─────────────────────────────────────────────────────────────┐
│              Permission Value Object                        │
│                  (权限值对象)                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Structure:                                                 │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Permission {                                         │ │
│  │    Resource: string  // 资源 (如：user, article)       │ │
│  │    Action: string    // 操作 (如：create, read)        │ │
│  │    Scope: string     // 范围 (如：own, all)            │ │
│  │  }                                                    │ │
│  │                                                       │ │
│  │  Examples:                                            │ │
│  │  • user:create:all   (创建任意用户)                    │ │
│  │  • user:read:own     (查看自己的用户信息)              │ │
│  │  • article:edit:own  (编辑自己的文章)                  │ │
│  │  • settings:manage   (管理系统设置)                    │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎭 聚合根之间的关系

### 关系图

```
┌──────────────────────────────────────────────────────────────┐
│                   聚合根关系总览                              │
└──────────────────────────────────────────────────────────────┘

┌─────────────┐         ┌─────────────┐
│    User     │         │   Tenant    │
│  Aggregate  │         │  Aggregate  │
└──────┬──────┘         └──────┬──────┘
       │                      │
       │ belongs to           │ has many
       │ (多对一)              │ (一对多)
       ↓                      ↓
┌─────────────────────────────────────────────────────────┐
│              TenantMember (关联表)                       │
│                                                         │
│  - TenantID: TenantID                                  │
│  - UserID: UserID                                      │
│  - Role: RoleID                                        │
│  - JoinedAt: time.Time                                 │
│                                                         │
│  职责：维护用户在租户中的成员关系和角色                   │
└─────────────────────────────────────────────────────────┘
       
       
┌─────────────┐         ┌─────────────┐
│   Tenant    │────────→│    Role     │
│  Aggregate  │  owns   │  Aggregate  │
└─────────────┘         └──────┬──────┘
                               │
                               │ contains
                               ↓
                    ┌──────────────────┐
                    │   Permission     │
                    │  (Value Object)  │
                    └──────────────────┘


详细关系说明：

1. User ←→ Tenant (通过 TenantMember)
   - 一个用户可以属于多个租户
   - 一个租户可以有多个成员
   - 关系存储在 TenantMember 表中
   
2. Tenant → Role
   - 租户拥有多个角色
   - 角色属于特定租户（隔离）
   
3. Role → Permission (值对象集合)
   - 角色包含多个权限
   - 权限是值对象，没有独立生命周期
   
4. User + Role → Access Control
   - 用户的权限 = 角色的权限并集
   - 通过角色间接授权给用户

```

---

## 📦 值对象详细设计

### 1. Email 值对象

```
┌─────────────────────────────────────────────────────────┐
│                  Email Value Object                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  type Email struct {                                    │
│    value string  // 存储邮箱地址                        │
│  }                                                      │
│                                                         │
│  Methods:                                               │
│  • NewEmail(email string) (*Email, error)              │
│    - 验证邮箱格式                                       │
│    - 转换为小写                                         │
│    - 检查长度限制                                       │
│                                                         │
│  • Value() string                                      │
│    - 返回邮箱字符串                                     │
│                                                         │
│  • Equals(other Email) bool                            │
│    - 比较两个邮箱是否相等                               │
│                                                         │
│  • Domain() string                                     │
│    - 提取域名部分                                       │
│                                                         │
│  Validation Rules:                                      │
│  ✓ 必须符合 RFC 5322 格式                               │
│  ✓ 长度不超过 255 字符                                   │
│  ✓ 必须包含 @ 符号                                       │
│  ✓ 域名必须有效                                         │
│  ✓ 自动转换为小写                                       │
└─────────────────────────────────────────────────────────┘
```

### 2. Password 值对象

```
┌─────────────────────────────────────────────────────────┐
│                Password Value Object                    │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  type Password struct {                                 │
│    hash string  // 存储哈希后的密码                      │
│  }                                                      │
│                                                         │
│  Methods:                                               │
│  • NewPassword(plain string, hasher PasswordHasher)    │
│    - 验证密码强度                                       │
│    - 使用 hasher 进行哈希                               │
│                                                         │
│  • Verify(plain string) bool                           │
│    - 验证明文密码是否匹配哈希                           │
│                                                         │
│  • Hash() string                                       │
│    - 返回密码哈希（用于存储）                           │
│                                                         │
│  Validation Rules:                                      │
│  ✓ 最少 8 个字符                                          │
│  ✓ 包含大小写字母                                       │
│  ✓ 包含数字或特殊字符                                   │
│  ✓ 不能是常见弱密码                                     │
│  ✓ 使用 bcrypt 加密 (cost=10)                           │
└─────────────────────────────────────────────────────────┘
```

### 3. UserID 值对象

```
┌─────────────────────────────────────────────────────────┐
│                 UserID Value Object                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  type UserID struct {                                   │
│    value int64  // Snowflake ID                         │
│  }                                                      │
│                                                         │
│  Methods:                                               │
│  • NewUserID(id int64) UserID                          │
│  • Int64() int64                                       │
│  • String() string                                     │
│  • Equals(other UserID) bool                           │
│                                                         │
│  Characteristics:                                       │
│  • 不可变                                               │
│  • 全局唯一 (Snowflake)                                 │
│  • 大致有序 (时间相关)                                  │
│  • 高性能比较                                           │
└─────────────────────────────────────────────────────────┘
```

---

## 🎯 聚合根设计原则

### 边界划分

```
┌─────────────────────────────────────────────────────────┐
│            聚合根边界划分原则                            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. User 聚合根                                         │
│     ✓ 负责用户身份认证                                  │
│     ✓ 管理用户状态（激活/锁定）                         │
│     ✓ 记录登录行为                                      │
│     ✗ 不负责：租户成员关系（由 Tenant 管理）             │
│                                                         │
│  2. Tenant 聚合根                                       │
│     ✓ 负责租户信息管理                                  │
│     ✓ 管理成员资格（添加/移除成员）                     │
│     ✓ 管理成员角色                                      │
│     ✗ 不负责：用户认证（由 User 管理）                   │
│                                                         │
│  3. Role 聚合根                                         │
│     ✓ 负责角色定义                                      │
│     ✓ 管理权限集合                                      │
│     ✗ 不负责：用户分配（由 Tenant 管理）                 │
│                                                         │
│  关键原则：                                              │
│  • 每个聚合根有明确的业务边界                           │
│  • 聚合根之间通过 ID 引用，不直接持有对象               │
│  • 跨聚合的操作通过领域事件协调                         │
│  • 保持聚合的小而专注                                   │
└─────────────────────────────────────────────────────────┘
```

### 一致性边界

```
┌─────────────────────────────────────────────────────────┐
│           事务一致性边界示意图                           │
└─────────────────────────────────────────────────────────┘

单个聚合内的事务：
┌──────────────────────────────────────────────────────┐
│  User 聚合                                            │
│                                                       │
│  [开启事务]                                          │
│  ├─ 更新密码                                         │
│  ├─ 记录密码修改时间                                 │
│  ├─ 发布 UserPasswordChangedEvent                    │
│  └─ [提交事务]                                       │
│                                                       │
│  ✅ 保证原子性和一致性                                │
└──────────────────────────────────────────────────────┘

跨聚合的最终一致性：
┌──────────────┐      领域事件       ┌──────────────┐
│    User      │ ──────────────────→ │   Tenant     │
│  Aggregate   │  UserActivatedEvent │  Aggregate   │
│              │                     │              │
│  [激活用户]  │                     │ [更新状态]   │
│              │                     │              │
│  ✅ 立即一致  │                     │ ⏱️ 最终一致  │
└──────────────┘                     └──────────────┘

说明：
• 同一聚合内：强一致性（事务保证）
• 跨聚合：最终一致性（领域事件）
• 避免大聚合（性能考虑）

```

---

## 🔄 生命周期管理

### User 聚合根的生命周期

```
┌─────────────────────────────────────────────────────────┐
│            User 聚合根生命周期状态机                     │
└─────────────────────────────────────────────────────────┘

                    [User.Register()]
                         ↓
              ┌──────────────────┐
              │   Registered     │ ← 初始状态
              │  (已注册未激活)   │
              └────────┬─────────┘
                       │
        ┌──────────────┼──────────────┐
        │ [Activate()] │              │ [Deactivate()]
        ↓              │              ↓
  ┌──────────┐    ┌──────────┐  ┌──────────┐
  │  Active  │    │ Inactive │  │ Suspended│
  │ (激活)   │←──→│ (未激活)  │  │ (暂停)   │
  └────┬─────┘    └──────────┘  └────┬─────┘
       │                             │
       │ [Lock()]                    │ [Unlock()]
       ↓                             ↓
  ┌──────────┐                 ┌──────────┐
  │  Locked  │ ←──────────────→│ Suspended│
  │ (锁定)   │   [Admin Action]│ (暂停)   │
  └──────────┘                 └──────────┘

状态转换规则：
✓ Registered → Active (邮件验证后)
✓ Active → Inactive (用户主动停用)
✓ Active → Locked (多次登录失败)
✓ Locked → Suspended (人工审核)
✓ Suspended → Active (申诉成功)
✗ Locked 不能直接转为 Active (需要审核)

每个状态转换都会发布领域事件！

```

---

## 📊 数据库映射关系

### ORM 映射图

```
┌─────────────────────────────────────────────────────────┐
│          领域模型 → 数据库表映射                          │
└─────────────────────────────────────────────────────────┘

Domain Layer (Go)              Database (PostgreSQL)
═══════════════════════════════════════════════════════

User Aggregate                 → users 表
  ├─ ID: UserID                → id (BIGINT, PK)
  ├─ Username: Username        → username (VARCHAR, UNIQUE)
  ├─ Email: Email              → email (VARCHAR, UNIQUE)
  ├─ Password: Password        → password_hash (VARCHAR)
  ├─ DisplayName               → display_name (VARCHAR)
  ├─ Avatar                    → avatar (VARCHAR)
  ├─ Status: UserStatus        → status (SMALLINT, ENUM)
  ├─ LastLoginAt               → last_login_at (TIMESTAMP)
  ├─ CreatedAt                 → created_at (TIMESTAMP)
  └─ UpdatedAt                 → updated_at (TIMESTAMP)

Tenant Aggregate               → tenants 表
  ├─ ID: TenantID              → id (BIGINT, PK)
  ├─ Name                      → name (VARCHAR)
  ├─ Slug                      → slug (VARCHAR, UNIQUE)
  ├─ Status: TenantStatus      → status (SMALLINT)
  └─ Config                    → config (JSONB)

TenantMember (关联)            → tenant_members 表
  ├─ TenantID                  → tenant_id (BIGINT, FK)
  ├─ UserID                    → user_id (BIGINT, FK)
  ├─ Role                      → role_id (BIGINT, FK)
  └─ JoinedAt                  → joined_at (TIMESTAMP)

Role Aggregate                 → roles 表
  ├─ ID: RoleID                → id (BIGINT, PK)
  ├─ Name                      → name (VARCHAR)
  ├─ DisplayName               → display_name (VARCHAR)
  ├─ TenantID                  → tenant_id (BIGINT, FK)
  └─ IsSystem                  → is_system (BOOLEAN)

Permissions (嵌入)             → role_permissions 表
  ├─ RoleID                    → role_id (BIGINT, FK)
  ├─ Resource                  → resource (VARCHAR)
  ├─ Action                    → action (VARCHAR)
  └─ Scope                     → scope (VARCHAR)

Domain Events                  → domain_events 表
  ├─ ID                        → id (BIGINT, PK)
  ├─ EventType                 → event_type (VARCHAR)
  ├─ AggregateType             → aggregate_type (VARCHAR)
  ├─ AggregateID               → aggregate_id (BIGINT)
  ├─ EventData                 → event_data (JSONB)
  ├─ OccurredAt                → occurred_at (TIMESTAMP)
  └─ Processed                 → processed (BOOLEAN)

索引策略：
✅ 主键索引：所有表的 id 字段
✅ 唯一索引：username, email, slug
✅ 外键索引：tenant_id, user_id, role_id
✅ 组合索引：(tenant_id, user_id) for TenantMember
✅ GIN 索引：config (JSONB), event_data (JSONB)

```

---

## 📚 参考文档

- [领域模型](./domain-model.md) - 详细领域模型说明
- [DDD 设计指南](./ddd-design-guide.md) - DDD 核心概念
- [架构总览](./architecture-overview.md) - 整体架构
