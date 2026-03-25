# 数据库主键雪花 ID 统一规范

## 📋 背景

在 DDD 架构中，**所有聚合根的主键 ID 都应该使用雪花算法生成**，而不是依赖数据库的 `AUTO_INCREMENT`（自增）。

### 为什么？

| 特性 | 数据库自增 | 雪花 ID |
|------|-----------|---------|
| **分布式友好** | ❌ 单点限制 | ✅ 天然支持 |
| **ID 可见性** | ❌ 插入后才能获取 | ✅ 提前生成 |
| **并发安全** | ❌ 锁竞争 | ✅ 无锁设计 |
| **数据迁移** | ❌ ID 冲突风险 | ✅ 全局唯一 |
| **分库分表** | ❌ 不支持 | ✅ 原生支持 |
| **性能** | ⚠️ 高并发瓶颈 | ✅ 高性能 |

---

## ✅ 当前状态

### 已完成的表

| 表名 | ID 类型 | 状态 |
|------|--------|------|
| `users` | BIGINT | ✅ 雪花 ID |
| `tenants` | BIGINT | ✅ 雪花 ID |
| `tenant_members` | (tenant_id, user_id) | ✅ 复合主键 |
| `roles` | BIGINT | ✅ 雪花 ID |
| `permissions` | BIGINT | ✅ 雪花 ID |
| `role_permissions` | (role_id, permission_id) | ✅ 复合主键 |
| `audit_logs` | BIGINT | ✅ 雪花 ID |
| `login_logs` | BIGINT | ✅ 雪花 ID |
| `outbox` | BIGINT | ✅ 雪花 ID |
| `activity_logs` | BIGINT | ✅ 雪花 ID |
| `domain_events` | BIGINT | ✅ 雪花 ID（本次修复） |
| `tenant_configs` | tenant_id (BIGINT) | ✅ 外键主键 |

---

## 🔧 本次修改内容

### 1. 修改数据库迁移文件

**文件：** `migrations/000008_create_domain_events_table.up.sql`

```sql
-- Before
CREATE TABLE IF NOT EXISTS domain_events (
    id BIGSERIAL PRIMARY KEY,  -- ❌ 数据库自增
    ...
);

-- After
CREATE TABLE IF NOT EXISTS domain_events (
    id BIGINT PRIMARY KEY,     -- ✅ 雪花 ID
    ...
);
```

---

### 2. 更新 DAO 模型

**文件：** `internal/infrastructure/persistence/model/domain_events.gen.go`

```go
// Before
type DomainEvent struct {
    ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键 ID: 自增主键" json:"id"`
    // ...
}

// After
type DomainEvent struct {
    ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false;comment:Snowflake ID：事件唯一标识，使用雪花算法生成" json:"id"`
    // ...
}
```

**关键点：**
- ✅ `autoIncrement: false` - 禁用自增
- ✅ 更新注释说明使用雪花 ID

---

### 3. 添加 ID 生成逻辑

**文件：** `internal/infrastructure/messaging/asynq/event_publisher.go`

```go
func (a *EventPublisherAdapter) saveEventLog(ctx context.Context, event kernel.DomainEvent) error {
    // ...
    
    now := time.Now()
    daoModel := &model.DomainEvent{
        ID:            idgen.Generate(), // ✅ 生成雪花 ID
        AggregateID:   a.aggregateIDToString(event.AggregateID()),
        AggregateType: aggregateType,
        EventType:     event.EventName(),
        EventData:     string(eventDataJSON),
        OccurredAt:    &now,
        CreatedAt:     &now,
    }
    
    return a.query.DomainEvent.WithContext(ctx).Create(daoModel)
}
```

**导入包：**
```go
import (
    // ...
    idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
)
```

---

## 🎯 完整实施清单

### Phase 1: 数据库迁移 ✅

- [x] 检查所有表的 ID 生成策略
- [x] 修复 `domain_events` 表的 `BIGSERIAL` 为 `BIGINT`
- [x] 更新字段注释说明

### Phase 2: DAO 模型更新 ✅

- [x] 更新 `DomainEvent` 模型的 gorm tag
- [x] 设置 `autoIncrement: false`
- [x] 更新注释

### Phase 3: 业务代码更新 ✅

- [x] 在 `saveEventLog()` 中添加雪花 ID 生成
- [x] 导入 `idgen` 包
- [x] 编译验证

### Phase 4: 验证测试 ⏳

- [ ] 运行集成测试
- [ ] 验证实体创建流程
- [ ] 验证事件保存流程

---

## 📊 使用场景对比

### ✅ 使用雪花 ID 的场景

1. **聚合根主键**
   ```go
   type User struct {
       id int64 // 雪花 ID
   }
   
   func NewUser(...) (*User, error) {
       id, _ := idGenerator.Generate()
       return &User{id: id, ...}, nil
   }
   ```

2. **领域事件日志**
   ```go
   daoModel := &model.DomainEvent{
       ID: idgen.Generate(), // ✅
   }
   ```

3. **活动日志**
   ```go
   daoModel := &model.ActivityLog{
       ID: idgen.Generate(), // ✅
   }
   ```

4. **Outbox 消息**
   ```go
   daoModel := &model.Outbox{
       ID: idgen.Generate(), // ✅
   }
   ```

---

### ⚠️ 特殊情况

#### 1. 复合主键

某些关联表使用复合主键，不需要额外的 ID：

```sql
-- tenant_members 表
CREATE TABLE tenant_members (
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(50),
    PRIMARY KEY (tenant_id, user_id),  -- ✅ 复合主键
    ...
);
```

**原因：**
- ✅ 自然键足够唯一
- ✅ 避免冗余索引
- ✅ 符合业务语义

---

#### 2. 外键作为主键

一对一关系时，直接使用外键作为主键：

```sql
-- tenant_configs 表
CREATE TABLE tenant_configs (
    tenant_id BIGINT PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    config JSONB,
    ...
);
```

**原因：**
- ✅ 强一致性约束
- ✅ 避免额外索引
- ✅ 明确的业务关系

---

## 🚀 最佳实践

### 1. 实体创建时生成 ID

```go
// ✅ 正确：在聚合根构造函数中生成
func NewUser(username, email, password string, idGenerator ports_idgen.Generator) (*User, error) {
    id, err := idGenerator.Generate()
    if err != nil {
        return nil, fmt.Errorf("failed to generate user id: %w", err)
    }
    
    user := &User{
        id:       id,
        username: username,
        // ...
    }
    return user, nil
}

// ❌ 错误：保存到数据库时才生成
func (r *UserRepository) Save(user *User) error {
    if user.id == 0 {
        user.id = idgen.Generate() // 不应该在这里生成
    }
    // ...
}
```

---

### 2. 事务外生成 ID

```go
// ✅ 正确：事务外生成
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // 1. 生成 ID（事务外）
    userID, _ := s.idGenerator.Generate()
    
    // 2. 创建聚合根
    user, _ := aggregate.NewUser(req.Username, req.Email, req.Password, userID)
    
    // 3. 保存到数据库（事务内）
    err := s.uow.Transaction(ctx, func(ctx context.Context) error {
        return s.userRepo.Save(ctx, user)
    })
    
    return user, err
}
```

---

### 3. 批量插入时预先生成 ID

```go
// ✅ 正确：预先生成所有 ID
func BatchCreateUsers(users []UserData) error {
    // 1. 预先生成所有 ID
    for i := range users {
        users[i].ID = idgen.Generate()
    }
    
    // 2. 批量插入
    return db.Model(&User{}).Create(users).Error
}

// ❌ 错误：依赖数据库自增
func BatchCreateUsers(users []UserData) error {
    // 依赖数据库自增，无法控制 ID
    return db.Model(&User{}).Create(users).Error
}
```

---

## ⚠️ 注意事项

### 1. 不要在数据库中设置自增

```sql
-- ❌ 错误
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,  -- 数据库自增
);

-- ✅ 正确
CREATE TABLE users (
    id BIGINT PRIMARY KEY,     -- 雪花 ID
);
```

---

### 2. GORM Model 配置

```go
// ❌ 错误
type User struct {
    ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true"`
}

// ✅ 正确
type User struct {
    ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false"`
}
```

---

### 3. 空值检查

确保在保存前 ID 已生成：

```go
// ✅ 正确：验证 ID
func (r *UserRepository) Save(ctx context.Context, user *User) error {
    if user.ID() == 0 {
        return fmt.Errorf("user id is required")
    }
    // ...
}

// ❌ 错误：不检查 ID
func (r *UserRepository) Save(ctx context.Context, user *User) error {
    // 直接保存，可能导致 ID=0
    return r.db.Save(user).Error
}
```

---

## 📈 性能对比

### 基准测试（yitter/idgenerator-go）

```bash
BenchmarkNextId-8    10000000    112 ns/op    0 B/op    0 allocs/op
```

### 对比数据库自增

| 操作 | 雪花 ID | 数据库自增 |
|------|---------|------------|
| 单次插入 | ~112ns | ~1ms（含网络） |
| 批量插入（1000 条） | ~0.1ms | ~10ms（含锁竞争） |
| 并发插入（100 QPS） | ✅ 无锁 | ⚠️ 锁竞争 |
| 并发插入（10000 QPS） | ✅ 稳定 | ❌ 性能下降 |

---

## 🔍 故障排查

### 问题 1：ID 重复

**可能原因：**
1. WorkerId 配置重复
2. 时间回拨（已自动处理）

**解决方案：**
```bash
# 检查日志中的 WorkerId
grep "worker_id" /var/log/app.log

# 确保每个节点 WorkerId 唯一
# 开发环境：手动配置不同值
# 生产环境：使用 Redis 自动注册
```

---

### 问题 2：ID 为 0

**可能原因：**
- 忘记调用 `idgen.Generate()`
- 错误处理被忽略

**解决方案：**
```go
// ✅ 正确：始终检查错误和返回值
id, err := idGenerator.Generate()
if err != nil {
    return fmt.Errorf("failed to generate id: %w", err)
}
if id == 0 {
    return fmt.Errorf("generated id is zero")
}

// ❌ 错误：忽略错误
id, _ := idGenerator.Generate()
```

---

### 问题 3：数据库报错

**错误信息：**
```
pq: duplicate key value violates unique constraint "users_pkey"
```

**原因：** 表中已有相同 ID

**解决方案：**
1. 清空表（开发环境）
2. 调整 WorkerId（生产环境）
3. 使用更大的 WorkerIdBitLength

---

## ✅ 检查清单

在新项目中应用此规范时，请确认：

### 数据库层面

- [ ] 所有表使用 `BIGINT` 而非 `BIGSERIAL`
- [ ] 主键注释说明"Snowflake ID"
- [ ] 迁移文件已更新

### 代码层面

- [ ] DAO 模型设置 `autoIncrement: false`
- [ ] 实体创建时生成 ID
- [ ] 保存前验证 ID 不为 0
- [ ] 错误处理完整

### 配置层面

- [ ] WorkerId 配置唯一
- [ ] 初始化在应用启动时完成
- [ ] 单元测试覆盖 ID 生成

---

## 📚 相关资源

- [雪花 ID 使用规范](./id-generator-usage-guide.md)
- [迁移报告](./ID_GENERATOR_MIGRATION_REPORT.md)
- [yitter/idgenerator-go GitHub](https://github.com/yitter/idgenerator-go)

---

**文档版本：** v1.0  
**创建日期：** 2026-03-25  
**维护者：** 架构委员会  
**状态：** ✅ 已批准并实施
