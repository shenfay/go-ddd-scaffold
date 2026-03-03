# DB Migrator 参考手册

## 目录

1. [Goose 迁移工具](#goose-迁移工具)
2. [GORM Model 规范](#gorm-model-规范)
3. [Repository 模式实现](#repository-模式实现)
4. [多环境配置](#多环境配置)
5. [最佳实践](#最佳实践)
6. [故障排除](#故障排除)

---

## Goose 迁移工具

### 什么是 Goose？

Goose 是一个数据库迁移管理工具，用于版本化数据库 schema。它支持：

- ✅ UP/DOWN 迁移（正向/反向）
- ✅ 事务支持
- ✅ 迁移状态追踪
- ✅ PostgreSQL、MySQL、SQLite 等多种数据库
- ✅ Go 语言编写，易于集成

### 迁移文件结构

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
COMMENT ON COLUMN users.id IS '主键 ID';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

### 关键指令

| 指令 | 说明 |
|------|------|
| `-- +goose Up` | UP 迁移开始 |
| `-- +goose Down` | DOWN 迁移开始 |
| `-- +goose StatementBegin` | SQL 语句块开始 |
| `-- +goose StatementEnd` | SQL 语句块结束 |

### 常用命令

```bash
# 执行所有待执行的迁移
goose postgres "connection-string" up

# 回滚最后一个迁移
goose postgres "connection-string" down

# 查看迁移状态
goose postgres "connection-string" status

# 创建新的迁移文件
goose create create_users_table sql
```

### 连接字符串格式

```
postgresql://用户名:密码@主机:端口/数据库名?sslmode=disable
```

示例：
```
postgresql://mathfun:math111@localhost:5432/mathfun_dev?sslmode=disable
```

---

## GORM Model 规范

### 基础结构

```go
package model

import "time"

type User struct {
    ID        int64      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Name      string     `gorm:"column:name;size:50;notNull" json:"name"`
    Email     string     `gorm:"column:email;size:100;uniqueIndex;notNull" json:"email"`
    CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
    return "users"
}
```

### 字段标签详解

#### 基础标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `column` | 指定列名 | `gorm:"column:user_name"` |
| `size` | 字符串最大长度 | `gorm:"size:100"` |
| `primaryKey` | 主键 | `gorm:"primaryKey"` |
| `autoIncrement` | 自增 | `gorm:"autoIncrement"` |
| `default` | 默认值 | `gorm:"default:'active'"` |
| `notNull` | 非空约束 | `gorm:"notNull"` |
| `uniqueIndex` | 唯一索引 | `gorm:"uniqueIndex"` |
| `index` | 普通索引 | `gorm:"index"` |

#### 关联标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `foreignKey` | 外键字段名 | `gorm:"foreignKey:UserID"` |
| `references` | 引用的字段 | `gorm:"references:ID"` |
| `Many2Many` | 多对多关联表 | `gorm:"Many2Many:user_roles"` |
| `constraint` | 级联规则 | `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` |

### 关联关系定义

#### 一对一

```go
type Profile struct {
    ID     int64
    UserID int64
    Bio    string
    User   User `gorm:"foreignKey:UserID"`
}

type User struct {
    ID      int64
    Profile *Profile `gorm:"foreignKey:UserID"`
}
```

#### 一对多

```go
type Order struct {
    ID     int64
    UserID int64
    User   User       `gorm:"foreignKey:UserID"`
    Items  []OrderItem `gorm:"foreignKey:OrderID"`
}

type User struct {
    ID     int64
    Orders []Order `gorm:"foreignKey:UserID"`
}
```

#### 多对多

```go
type User struct {
    ID    int64
    Roles []Role `gorm:"Many2Many:user_roles;"`
}

type Role struct {
    ID    int64
    Users []User `gorm:"Many2Many:user_roles;"`
}
```

### 软删除实现

```go
type User struct {
    ID        int64
    Name      string
    DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// 查询时自动过滤已删除的记录
db.First(&user, 1) // WHERE deleted_at IS NULL

// 真正删除（物理删除）
db.Unscoped().Delete(&user)
```

---

## Repository 模式实现

### 接口定义

```go
package dao

import (
    "context"
    "your-project/internal/infrastructure/persistence/gorm/model"
)

type UserRepository interface {
    Create(ctx context.Context, entity *model.User) error
    GetByID(ctx context.Context, id int64) (*model.User, error)
    Update(ctx context.Context, entity *model.User) error
    Delete(ctx context.Context, id int64) error
    FindAll(ctx context.Context, offset, limit int) ([]model.User, int64, error)
    
    // 自定义查询方法
    FindByEmail(ctx context.Context, email string) (*model.User, error)
    FindByStatus(ctx context.Context, status string, offset, limit int) ([]model.User, int64, error)
}
```

### 实现类

```go
package dao

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "your-project/internal/infrastructure/persistence/gorm/model"
)

type UserRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &UserRepositoryImpl{db: db}
}

// Create 创建记录
func (r *UserRepositoryImpl) Create(ctx context.Context, entity *model.User) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据 ID 获取
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
    var entity model.User
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &entity, nil
}

// Update 更新记录
func (r *UserRepositoryImpl) Update(ctx context.Context, entity *model.User) error {
    return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除记录（软删除）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).Delete(&model.User{ID: id}).Error
}

// FindAll 分页查询
func (r *UserRepositoryImpl) FindAll(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
    var entities []model.User
    var total int64
    
    tx := r.db.WithContext(ctx)
    
    // 统计总数
    if err := tx.Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    err := tx.Offset(offset).Limit(limit).Find(&entities).Error
    return entities, total, err
}

// FindByEmail 自定义查询：根据邮箱查找
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var entity model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&entity).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &entity, nil
}
```

### 事务处理

```go
func (s *UserService) CreateUserWithProfile(ctx context.Context, req CreateRequest) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建用户
        user := &model.User{Name: req.Name, Email: req.Email}
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // 创建档案
        profile := &model.Profile{UserID: user.ID, Bio: req.Bio}
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

---

## 多环境配置

### 配置文件结构

```yaml
database:
  host: localhost
  port: 5432
  user: mathfun
  password: math111
  name: mathfun_dev
  ssl_mode: disable

environments:
  dev:
    host: localhost
    port: 5432
    user: dev_user
    password: dev_pass
    name: mathfun_dev
    
  staging:
    host: staging-db.example.com
    port: 5432
    user: staging_user
    password: staging_pass
    name: mathfun_staging
    
  prod:
    host: prod-db.example.com
    port: 5432
    user: prod_user
    password: prod_pass
    name: mathfun_prod
    ssl_mode: require
```

### 环境变量覆盖

可以通过环境变量覆盖配置：

```bash
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure_password
export DB_NAME=mathfun_prod
```

### 环境切换

```bash
# 开发环境
db-migrator migrate up --env dev

# 预发布环境
db-migrator migrate up --env staging

# 生产环境
db-migrator migrate up --env prod
```

---

## 最佳实践

### 命名规范

#### 表命名
- ✅ 使用蛇形命名法（snake_case）：`user_profiles`
- ✅ 表名使用单数形式：`user` 而非 `users`
- ✅ 关联表使用复数：`user_roles`
- ❌ 避免保留字：不要使用 `order`, `group` 等

#### 字段命名
- ✅ 主键统一为 `id`
- ✅ 时间戳：`created_at`, `updated_at`, `deleted_at`
- ✅ 外键：`{table}_id`，如 `user_id`
- ✅ 布尔值：`is_`, `has_`, `can_` 前缀

### 迁移文件编写

1. **单一职责** - 每个迁移只做一个逻辑变更
2. **可回滚** - 必须提供对应的 DOWN 迁移
3. **添加备注** - 所有字段都要有 COMMENT
4. **使用事务** - 确保迁移原子性
5. **测试充分** - 在测试环境验证后再应用到生产

### 索引设计

```sql
-- 高频查询字段添加索引
CREATE INDEX idx_users_email ON users(email);

-- 复合索引（注意顺序）
CREATE INDEX idx_orders_user_status ON orders(user_id, status);

-- 部分索引（只索引满足条件的数据）
CREATE INDEX idx_active_users ON users(created_at) WHERE status = 'active';
```

### 性能优化

1. **分页查询** - 使用 OFFSET/LIMIT 或基于游标的分页
2. **批量操作** - 使用 Batch Insert/Update
3. **预加载** - 使用 GORM 的 Preload 避免 N+1 查询
4. **选择字段** - 只查询需要的字段

```go
// 预加载关联
db.Preload("Orders").Find(&users)

// 选择特定字段
db.Select("id", "name", "email").Find(&users)

// 批量插入
db.CreateInBatches(&users, 100)
```

---

## 故障排除

### 常见问题

#### 1. 数据库连接失败

**症状：**
```
dial tcp [::1]:5432: connect: connection refused
```

**解决方案：**
```bash
# 检查数据库服务是否运行
docker-compose ps

# 重启数据库
docker-compose restart postgres

# 检查端口占用
lsof -i :5432
```

#### 2. 迁移执行失败

**症状：**
```
ERROR: relation "users" already exists
```

**解决方案：**
```sql
-- 检查 goose 迁移记录表
SELECT * FROM goose_db_version ORDER BY version_id DESC;

-- 手动标记迁移为已执行
INSERT INTO goose_db_version (version_id, is_applied) VALUES (20260225120000, true);
```

#### 3. GORM Model 类型不匹配

**症状：**
```
sql: Scan error on column "created_at": unsupported Scan
```

**解决方案：**
```go
// 确保 Go 类型与 PostgreSQL 类型匹配
type User struct {
    CreatedAt time.Time  // TIMESTAMP
    Balance   float64    // DECIMAL
    Active    bool       // BOOLEAN
}
```

#### 4. 外键约束冲突

**症状：**
```
ERROR: insert or update on table "orders" violates foreign key constraint
```

**解决方案：**
```sql
-- 检查引用的记录是否存在
SELECT * FROM users WHERE id = 123;

-- 或者先删除外键约束再重建
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_user_id_fkey;
ALTER TABLE orders ADD CONSTRAINT orders_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
```

#### 5. 迁移回滚失败

**症状：**
```
ERROR: cannot drop table because other objects depend on it
```

**解决方案：**
```sql
-- 使用 CASCADE 级联删除
DROP TABLE IF EXISTS users CASCADE;

-- 或者先删除依赖对象
DROP INDEX IF EXISTS idx_orders_user;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_user_id_fkey;
DROP TABLE IF EXISTS users;
```

### 调试技巧

#### 启用 GORM 日志

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

#### 查看生成的 SQL

```go
// 开启 SQL 日志
db.Debug().Find(&users)

// 输出：
// SELECT * FROM "users" ORDER BY "users"."id" LIMIT 100
```

#### 使用 EXPLAIN 分析查询

```go
var results []model.User
db.Debug().Exec("EXPLAIN ANALYZE SELECT * FROM users WHERE email = ?", "test@example.com").Scan(&results)
```

---

## 附录

### PostgreSQL 类型映射

| PostgreSQL | Go | GORM Tag |
|------------|----|----------|
| BIGSERIAL | int64 | `autoIncrement` |
| VARCHAR(n) | string | `size:n` |
| TEXT | string | `type:text` |
| TIMESTAMP | time.Time | - |
| BOOLEAN | bool | - |
| DECIMAL(p,s) | float64 | `type:decimal(p,s)` |
| JSONB | datatypes.JSON | `type:jsonb` |
| BYTEA | []byte | `type:bytea` |

### 有用的资源

- 📚 [Goose 官方文档](https://github.com/pressly/goose)
- 📚 [GORM 官方文档](https://gorm.io/docs/)
- 📚 [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- 🔧 [在线 SQL 格式化工具](https://sqlformat.org/)

---

*本参考手册由 db-migrator Skill 生成，专为 MathFun 项目优化*
