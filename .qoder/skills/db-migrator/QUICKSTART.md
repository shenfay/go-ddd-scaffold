# DB Migrator 快速开始指南

## 5 分钟快速上手

### 第一步：安装 Skill

```bash
npx skills install db-migrator
```

### 第二步：创建迁移文件

```bash
db-migrator create \
  --table users \
  --description "create users table"
```

这会生成 Goose 迁移文件：
```
migrations/sql/20260225120000_create_users_table.sql
```

### 第三步：编辑迁移文件

自动打开编辑器，添加业务字段：

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    
    -- 业务字段
    name            VARCHAR(50) NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(20) DEFAULT 'member',
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

-- 添加字段备注
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '主键 ID';
COMMENT ON COLUMN users.name IS '用户姓名';
COMMENT ON COLUMN users.email IS '邮箱地址';
COMMENT ON COLUMN users.password_hash IS '密码哈希';
COMMENT ON COLUMN users.role IS '用户角色';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';
COMMENT ON COLUMN users.deleted_at IS '删除时间（软删除）';

-- 添加索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- +goose Down
DROP TABLE IF EXISTS users;
```

### 第四步：执行数据库迁移

```bash
# 确保数据库已启动
docker-compose up -d postgres

# 执行迁移
db-migrator migrate up --env dev
```

查看迁移状态：
```bash
db-migrator migrate status --env dev
```

### 第五步：生成 GORM Model 和 DAO

```bash
db-migrator generate \
  --tables users \
  --output ./internal/infrastructure/persistence/gorm \
  --soft-delete
```

生成的文件：
```
internal/infrastructure/persistence/gorm/
├── model/
│   └── user.go              # GORM Model
└── dao/
    ├── user_repository.go          # Repository 接口
    └── user_repository_impl.go     # Repository 实现
```

### 第六步：使用 DAO

在应用服务中注入并使用：

```go
import (
    "your-project/internal/infrastructure/persistence/gorm/dao"
)

type UserService struct {
    userRepo dao.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        userRepo: dao.NewUserRepository(db),
    }
}

func (s *UserService) CreateUser(ctx context.Context, req CreateRequest) error {
    entity := &model.User{
        Name:  req.Name,
        Email: req.Email,
    }
    return s.userRepo.Create(ctx, entity)
}
```

## 完整工作流（一键完成）

```bash
db-migrator full-workflow \
  --table products \
  --description "create products table" \
  --env dev \
  --generate-dao
```

这会引导你完成：
1. ✅ 创建迁移文件
2. ✅ 编辑迁移（手动）
3. ✅ 执行迁移
4. ✅ 生成 DAO 代码

## 多环境配置

### 开发环境
```bash
db-migrator migrate up --env dev
```

### 预发布环境
```bash
db-migrator migrate up --env staging
```

### 生产环境
```bash
db-migrator migrate up --env prod
```

## 常用命令

### 创建迁移
```bash
# 基础用法
db-migrator create --table orders --description "create orders table"

# 批量创建
db-migrator create --table order_items --description "create order items table"
```

### 执行迁移
```bash
# 执行所有待执行的迁移
db-migrator migrate up --env dev

# 回滚最后一个迁移
db-migrator migrate down --env dev

# 查看迁移状态
db-migrator migrate status --env dev
```

### 生成代码
```bash
# 单个表
db-migrator generate --tables users

# 多个表
db-migrator generate --tables users,products,orders

# 包含测试文件
db-migrator generate --tables users --with-tests
```

## 生成的代码结构

### GORM Model (`model/user.go`)

```go
package model

import "time"

// User 用户表模型
type User struct {
    ID        int64      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Name      string     `gorm:"column:name;size:50;notNull" json:"name"`
    Email     string     `gorm:"column:email;size:100;uniqueIndex;notNull" json:"email"`
    Role      string     `gorm:"column:role;size:20;default:'member'" json:"role"`
    CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
    return "users"
}
```

### Repository 接口 (`dao/user_repository.go`)

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
}
```

### Repository 实现 (`dao/user_repository_impl.go`)

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

func (r *UserRepositoryImpl) Create(ctx context.Context, entity *model.User) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

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

// ... 其他方法
```

## 配置说明

Skill 使用 `.qoder/skills/db-migrator/config.yaml` 进行配置：

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
```

## 下一步

### 1. 添加关联关系

编辑生成的 Model 添加外键关联：

```go
type Order struct {
    ID        int64
    UserID    int64      `gorm:"column:user_id;index" json:"user_id"`
    User      *User      `gorm:"foreignKey:UserID" json:"user"`
    // ...
}
```

### 2. 添加自定义查询方法

在 Repository 实现中添加：

```go
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var entity model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&entity).Error
    return &entity, err
}
```

### 3. 编写单元测试

如果使用了 `--with-tests` 选项：

```go
func TestUserRepository_Create(t *testing.T) {
    // TODO: 实现测试逻辑
}
```

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：咨询 DDD Architect Agent

## 推荐学习路径

1. ✅ 完成本快速开始（5 分钟）
2. 📚 阅读 [REFERENCE.md](./REFERENCE.md) 了解 Goose 最佳实践（15 分钟）
3. 🔍 研究生成的代码结构（30 分钟）
4. 🎯 参考 [EXAMPLES.md](./EXAMPLES.md) 添加复杂表结构（1 小时）
5. 🤖 咨询 DDD Architect Agent 优化领域模型（按需）

祝你开发顺利！🚀
