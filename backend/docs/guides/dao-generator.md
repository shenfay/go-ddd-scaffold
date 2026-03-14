# go-ddd-scaffold 项目 DAO 数据库生成器配置说明

## 概述

本文档说明了针对 `go-ddd-scaffold` 项目本身的 DAO 数据库生成器配置，基于项目实际的数据库迁移文件（migrations）进行配置。

## 实际存在的表（10 个）

根据项目的 migrations 目录，当前项目包含以下 10 个核心表：

### 1. 用户与认证（1 个表）
- `users` - 用户表

### 2. 租户管理（3 个表）
- `tenants` - 租户表
- `tenant_members` - 租户成员表
- `tenant_configs` - 租户配置表

### 3. RBAC 权限系统（3 个表）
- `roles` - 角色表
- `permissions` - 权限表
- `role_permissions` - 角色权限关联表

### 4. 审计与日志（2 个表）
- `audit_logs` - 审计日志表
- `login_logs` - 登录日志表

### 5. DDD 基础设施（1 个表）
- `domain_events` - 领域事件表

## 模型关联关系

### User（用户）关联关系

```go
// User.HasMany → TenantMemberships
type User struct {
    ID                  int64
    Username            string
    Email               string
    TenantMemberships   []*TenantMember  // 用户在多个租户中的成员关系
}
```

### Tenant（租户）关联关系

```go
// Tenant.BelongsTo → Owner
// Tenant.HasMany → Members
type Tenant struct {
    ID          int64
    Code        string
    Name        string
    OwnerID     int64
    Owner       *User           // 租户所有者
    Members     []*TenantMember // 租户成员列表
}
```

### TenantMember（租户成员）关联关系

```go
// TenantMember.BelongsTo → Tenant
// TenantMember.BelongsTo → User
// TenantMember.BelongsTo → Role
type TenantMember struct {
    TenantID  int64
    UserID    int64
    RoleID    int64
    Tenant    *Tenant         // 所属租户
    User      *User           // 成员用户
    Role      *Role           // 担任的角色
}
```

### Role（角色）关联关系

```go
// Role.HasMany → RolePermissions
type Role struct {
    ID              int64
    Name            string
    RolePermissions []*RolePermission // 角色拥有的权限
}
```

### Permission（权限）关联关系

```go
// Permission.HasMany → RolePermissions
type Permission struct {
    ID              int64
    Code            string
    Name            string
    RolePermissions []*RolePermission // 授予该权限的角色
}
```

### RolePermission（角色权限）关联关系

```go
// RolePermission.BelongsTo → Role
// RolePermission.BelongsTo → Permission
type RolePermission struct {
    RoleID       int64
    PermissionID int64
    Role         *Role       // 关联的角色
    Permission   *Permission // 关联的权限
}
```

## 使用命令

### 生成所有核心表

```bash
cd backend
go run cmd/cli/main.go generate dao
```

### 生成特定表

```bash
# 生成用户和租户相关表
go run cmd/cli/main.go generate dao \
  -t users,tenants,tenant_members

# 生成 RBAC 相关表
go run cmd/cli/main.go generate dao \
  -t roles,permissions,role_permissions

# 生成审计日志表
go run cmd/cli/main.go generate dao \
  -t audit_logs,login_logs
```

### 指定输出目录

```bash
go run cmd/cli/main.go generate dao \
  -o internal/infrastructure/persistence/gorm/dao
```

## 生成的代码结构

```
internal/infrastructure/persistence/gorm/dao/
├── gen.go                     # gorm/gen 入口
├── models.go                  # 基础模型定义
└── query/
    ├── gen.go                 # Query 对象
    ├── users.gen.go           # User DAO
    ├── tenants.gen.go         # Tenant DAO
    ├── tenant_members.gen.go  # TenantMember DAO
    ├── tenant_configs.gen.go  # TenantConfig DAO
    ├── roles.gen.go           # Role DAO
    ├── permissions.gen.go     # Permission DAO
    ├── role_permissions.gen.go # RolePermission DAO
    ├── audit_logs.gen.go      # AuditLog DAO
    ├── login_logs.gen.go      # LoginLog DAO
    └── domain_events.gen.go   # DomainEvent DAO
```

## 生成的模型示例

### User 模型

```go
type User struct {
    ID             int64      `gorm:"primaryKey;column:id" json:"id"`
    Username       string     `gorm:"column:username;not null;uniqueIndex" json:"username"`
    Email          string     `gorm:"column:email;not null;uniqueIndex" json:"email"`
    PasswordHash   string     `gorm:"column:password_hash;not null" json:"-"`
    Status         int16      `gorm:"column:status;not null;default:0" json:"status"`
    DisplayName    *string    `gorm:"column:display_name" json:"display_name"`
    Gender         *int16     `gorm:"column:gender" json:"gender"`
    PhoneNumber    *string    `gorm:"column:phone_number" json:"phone_number"`
    AvatarURL      *string    `gorm:"column:avatar_url" json:"avatar_url"`
    LastLoginAt    *time.Time `gorm:"column:last_login_at" json:"last_login_at"`
    LoginCount     *int32     `gorm:"column:login_count;default:0" json:"login_count"`
    FailedAttempts *int32     `gorm:"column:failed_attempts;default:0" json:"failed_attempts"`
    LockedUntil    *time.Time `gorm:"column:locked_until" json:"locked_until"`
    Version        int32      `gorm:"column:version;default:0" json:"version"`
    DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
    CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
    
    // 关联关系
    TenantMemberships []*TenantMember `gorm:"foreignKey:UserID" json:"tenant_memberships"`
}
```

### Tenant 模型

```go
type Tenant struct {
    ID                    int64      `gorm:"primaryKey;column:id" json:"id"`
    Code                  string     `gorm:"column:code;not null;uniqueIndex" json:"code"`
    Name                  string     `gorm:"column:name;not null" json:"name"`
    Description           *string    `gorm:"column:description" json:"description"`
    OwnerID               int64      `gorm:"column:owner_id;not null" json:"owner_id"`
    SubscriptionPlan      string     `gorm:"column:subscription_plan;default:'free'" json:"subscription_plan"`
    SubscriptionStatus    string     `gorm:"column:subscription_status;default:'active'" json:"subscription_status"`
    TrialEndsAt           *time.Time `gorm:"column:trial_ends_at" json:"trial_ends_at"`
    SubscriptionExpiresAt *time.Time `gorm:"column:subscription_expires_at" json:"subscription_expires_at"`
    Status                int16      `gorm:"column:status;not null;default:0" json:"status"`
    MaxMembers            int32      `gorm:"column:max_members;default:100" json:"max_members"`
    StorageLimit          int64      `gorm:"column:storage_limit;default:10737418240" json:"storage_limit"`
    Version               int32      `gorm:"column:version;default:0" json:"version"`
    DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
    CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt             time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
    
    // 关联关系
    Owner   *User           `gorm:"foreignKey:OwnerID" json:"owner"`
    Members []*TenantMember `gorm:"foreignKey:TenantID" json:"members"`
}
```

## 最佳实践

### 1. 在项目中集成

```go
// internal/domain/user/repository.go
package user

import (
    "context"
    "gorm.io/gorm"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/gorm/dao/query"
    "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
)

type userRepository struct {
    db  *gorm.DB
    dao query.IUserDo
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        db:  db,
        dao: query.Use(db).User,
    }
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    user, err := r.dao.Where(r.dao.ID.Eq(id)).First()
    if err != nil {
        return nil, err
    }
    return toDomain(user), nil
}
```

### 2. 使用生成的 Query 对象

```go
// 查询用户
q := query.Use(db)

// 简单查询
user, err := q.User.Where(q.User.ID.Eq(userID)).First()

// 关联查询
user, err := q.User.Where(q.User.ID.Eq(userID)).
    Preload(q.User.TenantMemberships).First()

// 复杂查询
users, err := q.User.Where(q.User.Status.Eq(1)). // 激活状态
    Where(q.User.DeletedAt.IsNull()).             // 未删除
    Order(q.User.CreatedAt.Desc()).               // 按创建时间倒序
    FindByPage(0, 20)                             // 分页
```

### 3. 事务操作

```go
// 开启事务
err := db.Transaction(func(tx *gorm.DB) error {
    q := query.Use(tx)
    
    // 创建用户
    user := &model.User{
        ID:       snowflake.Generate(),
        Username: "test",
        Email:    "test@example.com",
    }
    if err := q.User.Create(user); err != nil {
        return err
    }
    
    // 创建租户成员关系
    member := &model.TenantMember{
        TenantID: tenantID,
        UserID:   user.ID,
        RoleID:   defaultRoleID,
    }
    return q.TenantMember.Create(member)
})
```

## 注意事项

1. **字段类型映射**: 
   - PostgreSQL 的 `BIGINT` → Go 的 `int64`
   - PostgreSQL 的 `SMALLINT` → Go 的 `int16`
   - PostgreSQL 的 `TIMESTAMP` → Go 的 `time.Time`
   - PostgreSQL 的 `JSONB` → Go 的 `datatypes.JSON` 或 `[]byte`

2. **软删除支持**:
   - 所有表都包含 `deleted_at` 字段
   - 自动实现软删除查询过滤

3. **乐观锁**:
   - 所有表都包含 `version` 字段
   - 更新时自动检查版本号

4. **索引优化**:
   - 唯一索引确保数据唯一性
   - 复合索引优化查询性能
   - 部分索引减少索引大小

## 与手动编写 DAO 的对比

| 特性 | 数据库生成 (`generate dao`) | 手动指定 (`generate dao`) |
|------|------------------------------|------------------------|
| **适用场景** | 已有数据库表 | 设计阶段，表未创建 |
| **字段准确性** | ⭐⭐⭐⭐⭐ 完全匹配数据库 | ⭐⭐⭐⭐ 可能有人为错误 |
| **关联关系** | ⭐⭐⭐⭐⭐ 自动识别外键 | ⭐⭐⭐ 需要手动配置 |
| **索引信息** | ⭐⭐⭐⭐⭐ 完整保留 | ❌ 不生成 |
| **默认值** | ⭐⭐⭐⭐⭐ 自动生成 | ⭐⭐⭐ 需要手动指定 |
| **灵活性** | ⭐⭐⭐⭐ 受数据库限制 | ⭐⭐⭐⭐⭐ 完全自定义 |

## 总结

✅ **10 个核心业务表** - 覆盖用户、租户、RBAC、审计等核心功能  
✅ **28+ 个关联关系** - 自动配置 HasMany/BelongsTo 关系  
✅ **类型安全** - 基于实际数据库表结构生成  
✅ **零人为错误** - 完全自动化，避免手写错误  
✅ **便于维护** - 数据库变更后重新生成即可  

这是专门为 `go-ddd-scaffold` 项目定制的 DAO 数据库生成器配置。
