# DB Migrator - 数据库迁移与 DAO 生成工具 v2.0

## 📦 完整文件列表

```
.qoder/skills/db-migrator/
├── SKILL.md              # 技能主文档（6.4KB）
├── config.yaml           # 配置文件（3.8KB）
├── QUICKSTART.md         # 快速开始指南（7.9KB）
├── EXAMPLES.md           # 使用示例（12KB）
├── REFERENCE.md          # 参考手册（13KB）
├── README.md             # 本文件
└── scripts/
    ├── generate.py       # Python 生成脚本（14KB，可执行）
    └── helper.sh         # Bash 辅助脚本（7KB）
```

**总计**: 约 52KB 文档 + 21KB 代码

---

## ✅ 完成的功能

### 1. 核心功能 (100%)

- ✅ **Goose 迁移管理**
  - 自动生成时间戳版本号（YYYYMMDDHHMMSS）
  - UP/DOWN 迁移支持
  - 字段备注自动生成（COMMENT ON COLUMN）
  - 事务支持
  
- ✅ **GORM Model 生成**
  - PostgreSQL → Go 类型自动映射
  - json/gorm/validate 标签生成
  - 关联关系支持（HasOne/HasMany/ManyToMany）
  - 软删除集成
  
- ✅ **Repository 模式 DAO**
  - 接口与实现分离
  - CRUD 方法自动生成
  - 分页查询支持
  - 事务处理集成

### 2. 多环境支持 (100%)

- ✅ dev/staging/prod 三套环境配置
- ✅ 环境变量覆盖支持
- ✅ SSL 模式配置
- ✅ 密码加密存储（通过环境变量）

### 3. 文档完整性 (100%)

- ✅ SKILL.md - 技能说明和功能介绍
- ✅ QUICKSTART.md - 5 分钟快速上手
- ✅ EXAMPLES.md - 4 个完整场景示例
- ✅ REFERENCE.md - Goose/GORM 最佳实践
- ✅ config.yaml - 多环境配置

---

## 🚀 使用方式

### 基本命令

```bash
# 创建迁移
db-migrator create --table users --description "create users table"

# 执行迁移
db-migrator migrate up --env dev

# 生成 DAO
db-migrator generate --tables users --output ./internal/persistence/gorm

# 完整工作流
db-migrator full-workflow --table products --description "create products" --generate-dao
```

### 典型工作流

1. 使用 `db-migrator create` 创建迁移
2. 编辑 SQL 添加业务字段和备注
3. 使用 `db-migrator migrate up` 执行迁移
4. 使用 `db-migrator generate` 生成 GORM Model 和 DAO
5. 在应用服务中注入并使用 DAO

---

## 📊 生成的代码结构

```
backend/internal/infrastructure/persistence/gorm/
├── model/
│   ├── user.go              # User 模型
│   ├── product.go           # Product 模型
│   └── order.go             # Order 模型
└── dao/
    ├── user_repository.go          # User 仓储接口
    ├── user_repository_impl.go     # User 仓储实现
    ├── product_repository.go       # Product 仓储接口
    └── ...
```

---

## 🎯 特色功能

### 1. 字段备注自动化

所有字段都会自动生成 COMMENT：

```sql
COMMENT ON COLUMN users.id IS '主键 ID';
COMMENT ON COLUMN users.email IS '邮箱地址';
COMMENT ON COLUMN users.created_at IS '创建时间';
```

### 2. 智能类型映射

| PostgreSQL | Go | GORM Tag |
|------------|----|----------|
| BIGSERIAL | int64 | `autoIncrement` |
| VARCHAR(50) | string | `size:50` |
| TIMESTAMP | time.Time | `autoCreateTime` |
| JSONB | datatypes.JSON | `type:jsonb` |

### 3. Repository 模式

```go
// 接口定义
type UserRepository interface {
    Create(ctx context.Context, entity *model.User) error
    GetByID(ctx context.Context, id int64) (*model.User, error)
    Update(ctx context.Context, entity *model.User) error
    Delete(ctx context.Context, id int64) error
    FindAll(ctx context.Context, offset, limit int) ([]model.User, int64, error)
}

// 依赖注入
func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        userRepo: dao.NewUserRepository(db),
    }
}
```

---

## 🔧 配置说明

### 开发环境

```yaml
database:
  host: localhost
  port: 5432
  user: mathfun
  password: math111
  name: mathfun_dev
  ssl_mode: disable
```

### 生产环境

```yaml
environments:
  prod:
    host: prod-db.mathfun.local
    port: 5432
    user: mathfun_prod
    password: ${PROD_DB_PASSWORD}  # 环境变量
    name: mathfun_prod
    ssl_mode: require
```

---

## 📚 学习路径

1. **快速开始** (5 分钟)
   - 阅读 [QUICKSTART.md](./QUICKSTART.md)
   - 完成第一个表创建
   
2. **深入理解** (30 分钟)
   - 阅读 [REFERENCE.md](./REFERENCE.md)
   - 学习 Goose 和 GORM 最佳实践
   
3. **实战演练** (1 小时)
   - 参考 [EXAMPLES.md](./EXAMPLES.md)
   - 完成 4 个场景练习
   
4. **高级应用** (按需)
   - 自定义 Repository 方法
   - 复杂关联关系处理
   - 性能优化技巧

---

## 🐛 常见问题

### 数据库连接失败

```bash
# 检查数据库服务
docker-compose ps

# 重启数据库
docker-compose restart postgres
```

### 迁移执行失败

```bash
# 查看迁移状态
db-migrator migrate status --env dev

# 回滚最后一个迁移
db-migrator migrate down --env dev
```

### 代码生成失败

确保：
1. 数据库表已正确创建
2. go.mod 包含所需依赖
3. 输出目录有写权限

---

## 📈 版本历史

### v2.0.0 (2026-02-25) - 重大更新

- ✅ 确定使用 **Goose** 作为迁移工具
- ✅ 新增 **多环境配置** 支持（dev/staging/prod）
- ✅ 增强 **字段备注** 自动生成
- ✅ 改进 **GORM Model** 类型映射
- ✅ 新增 **Repository 模式** DAO 生成
- ✅ 完善文档体系（QUICKSTART/EXAMPLES/REFERENCE）

### v1.0.0 (2026-01-26) - 初始版本

- 基础迁移文件生成
- DAO 代码自动生成
- 完整工作流支持

---

## 🤝 与其他 Skills 协同

### 上游依赖
- `ddd-modeling-assistant` - 领域模型设计

### 下游依赖
- `api-endpoint-generator` - API 端点生成
- `ddd-development-workflow` - DDD 开发流程

### 推荐工作流

```
1. ddd-modeling-assistant (设计领域模型)
   ↓
2. db-migrator (生成数据库和 DAO)
   ↓
3. api-endpoint-generator (生成 API 端点)
   ↓
4. 手动完善业务逻辑
```

---

## 📖 相关资源

- 📚 [Goose 官方文档](https://github.com/pressly/goose)
- 📚 [GORM 官方文档](https://gorm.io/docs/)
- 📚 [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- 💡 [MathFun 技术栈规范](../../../docs/system_design/服务端技术栈规范文档.md)

---

*本 Skill 专为 MathFun 项目优化设计，遵循 Qoder Skills 规范*
