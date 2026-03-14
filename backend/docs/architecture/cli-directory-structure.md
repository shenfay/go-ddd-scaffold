# CLI 工具目录结构说明

## 整体架构

```
backend/cmd/cli/
├── main.go                          # CLI 程序入口
└── internal/                        # 内部包（不对外暴露）
    ├── command/                     # 命令定义层
    │   ├── root.go                  # 根命令和统一注册
    │   ├── init.go                  # 项目初始化命令
    │   ├── generate.go              # 代码生成根命令
    │   └── migrate_docs.go          # 迁移和文档命令
    └── generators/                  # 代码生成层
        ├── types.go                 # 所有生成器的选项定义
        ├── dao_generator.go         # DAO 生成器（从数据库反向工程）⭐
        ├── entity_generator.go      # 实体生成器（待实现）
        ├── repository_generator.go  # Repository 生成器（待实现）
        ├── service_generator.go     # Service 生成器（待实现）
        ├── handler_generator.go     # Handler 生成器（待实现）
        ├── dto_generator.go         # DTO 生成器（待实现）
        ├── init_generator.go        # 项目初始化生成器
        └── stubs.go                 # 其他生成器存根
```

## 核心文件说明

### 1. 入口文件

#### `main.go` (35 行)
**职责**: CLI 程序启动入口
- 创建根命令
- 设置版本信息
- 执行命令

```go
func main() {
    rootCmd := newRootCmd()
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### 2. 命令层 (`internal/command/`)

#### `root.go` (~50 行)
**职责**: 注册所有子命令到根命令
- 统一命令注册入口
- 全局标志配置
- 版本信息显示

```go
func RegisterAll(rootCmd *cobra.Command) {
    rootCmd.AddCommand(initCmd())
    rootCmd.AddCommand(generateCmd())
    rootCmd.AddCommand(migrateCmd())
    rootCmd.AddCommand(docsCmd())
    rootCmd.AddCommand(cleanCmd())
    rootCmd.AddCommand(versionCmd())
}
```

#### `generate.go` (~160 行)
**职责**: 代码生成根命令及其子命令
- `generate entity` - 领域实体生成
- `generate dao` - 从数据库反向工程生成 DAO ⭐
- `generate repository` - Repository 层生成
- `generate service` - Service 层生成
- `generate handler` - Handler 层生成
- `generate dto` - DTO 层生成

```go
func generateCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "generate",
        Short: "Generate code for DDD scaffold",
        Aliases: []string{"gen", "g"},
    }
    
    cmd.AddCommand(generateEntityCmd())
    cmd.AddCommand(generateDAOCmd())
    // ... 其他命令
    
    return cmd
}
```

#### `init.go` (~45 行)
**职责**: 项目初始化命令
- 创建新项目结构
- 配置 Go module
- 选择项目模板

```bash
go-ddd-scaffold init my-project \
  --author "Your Name" \
  --email "your.email@example.com"
```

#### `migrate_docs.go` (~65 行)
**职责**: 数据库迁移和文档生成命令
- `migrate up` - 运行迁移
- `migrate down` - 回滚迁移
- `migrate create` - 创建迁移
- `docs swagger` - 生成 Swagger 文档

### 3. 生成器层 (`internal/generators/`)

#### `types.go` (~75 行)
**职责**: 定义所有生成器的选项结构
- 统一的类型定义
- 便于在 command 层和 generators 层之间传递参数

```go
// DAOOptions DAO 生成选项
type DAOOptions struct {
    Name          string
    Fields        string
    TableName     string
    OutputDir     string
    WithInterface bool
    WithCondition bool
}

// DAODBOptions DAO from database generation options
type DAODBOptions struct {
    OutputPath    string
    DSN           string
    ConfigFile    string
    WithUnitTest  bool
    FieldNullable bool
    Tables        []string
    AllTables     bool
}

// EntityOptions, RepositoryOptions, etc.
```

#### `dao_generator.go` (~400 行) ⭐
**职责**: 从数据库反向工程生成 DAO
- 连接 PostgreSQL 数据库
- 读取表结构
- 使用 gorm/gen 生成模型和 DAO
- 配置模型关联关系
- 支持环境变量配置加载

**命令**:
```bash
go-ddd-scaffold generate dao [flags]
```

**参数**:
- `-o, --output` - 输出目录
- `-d, --dsn` - 数据库连接字符串（可选，优先使用环境变量）
- `-t, --tables` - 指定要生成的表（逗号分隔）
- `--field-nullable` - 为可空字段生成指针
- `--with-test` - 生成单元测试

**示例**:
```bash
# 生成所有核心表（默认）
go-ddd-scaffold generate dao

# 生成特定表
go-ddd-scaffold generate dao -t users,tenants

# 使用自定义数据库连接
go-ddd-scaffold generate dao -d "postgres://user:pass@host:5432/dbname"
```

**生成的文件**:
- `internal/infrastructure/persistence/dao/user_dao_gen.go`
- `internal/infrastructure/persistence/dao/tenant_dao_gen.go`
- `internal/infrastructure/persistence/dao/role_dao_gen.go`
- 等...

#### 其他生成器（待实现）

- `entity_generator.go` - 领域实体生成器
- `repository_generator.go` - Repository 层生成器
- `service_generator.go` - Service 层生成器
- `handler_generator.go` - CQRS Handler 生成器
- `dto_generator.go` - DTO 生成器

#### `init_generator.go` (~27 行)
**职责**: 项目初始化生成器（存根）
- 创建项目目录结构
- 生成配置文件
- 初始化 Git 仓库

#### `stubs.go` (~30 行)
**职责**: 其他生成器的存根实现
- 快速占位，便于编译
- 后续逐步实现完整功能

## 两种 DAO 生成方式对比

| 特性 | `generate dao` | `generate dao` |
|------|---------------|------------------|
| **数据来源** | 手动指定字段 | 数据库表结构 |
| **适用场景** | 设计阶段，表未创建 | 已有数据库，需要快速生成 |
| **灵活性** | 高，完全自定义 | 中，依赖数据库结构 |
| **代码风格** | 手写风格，符合 DDD | gorm/gen 风格，类型安全 |
| **关联关系** | 手动配置 | 自动配置 |
| **生成位置** | `internal/infrastructure/dao` | `internal/infrastructure/persistence/gorm/dao` |

## 使用场景

### 场景 1: 新项目的领域建模阶段

```bash
# 1. 先设计领域实体
go-ddd-scaffold generate entity user \
  -f "username:string,email:string"

# 2. 生成 DAO 接口
go-ddd-scaffold generate dao user \
  -f "username:string,email:string,password:string"

# 3. 创建数据库迁移
go-ddd-scaffold migrate create create_users_table
```

### 场景 2: 已有数据库，需要生成代码

```bash
# 直接从数据库反向工程
go-ddd-scaffold generate dao \
  -t users,tenants,tenant_members

# 或生成所有表
go-ddd-scaffold generate dao --all-tables
```

### 场景 3: 混合使用

```bash
# 1. 核心业务表用 dao 生成
go-ddd-scaffold generate dao \
  -t users,tenants,ingredients,food_items

# 2. 特殊业务逻辑手动编写
# 手动编写复杂的业务查询和服务
```

## 生成的代码集成

### DAO 层集成示例

```go
// internal/domain/user/repository.go
package user

import (
    "context"
    "gorm.io/gorm"
    "foodknow/internal/infrastructure/persistence/gorm/dao/query"
    "foodknow/internal/infrastructure/persistence/gorm/model"
)

type userRepository struct {
    db   *gorm.DB
    dao  query.IUserDo
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
    
    // 转换为领域实体
    return toDomain(user), nil
}
```

## 最佳实践

### 1. 目录隔离

```
internal/infrastructure/
├── dao/                         # 手动编写的 DAO
│   ├── user_dao.go
│   └── user_dao_impl.go
└── persistence/
    └── gorm/
        └── dao/                 # gorm/gen 生成的 DAO
            ├── gen.go
            ├── models.go
            └── query/
                ├── users.gen.go
                └── ...
```

### 2. 命名规范

- **手动 DAO**: `{entity}_dao.go`, `{entity}_dao_impl.go`, `{entity}_model.go`
- **生成 DAO**: `{table}.gen.go`（由 gorm/gen 自动命名）

### 3. 版本控制

```bash
# 将生成的代码纳入版本控制
git add internal/infrastructure/persistence/gorm/dao/

# 但忽略临时文件
echo "*.tmp" >> .gitignore
```

### 4. 定期同步

数据库变更后及时重新生成：

```bash
# 应用数据库迁移
go-ddd-scaffold migrate up

# 重新生成 DAO
go-ddd-scaffold generate dao -t modified_table
```

## 扩展机制

### 添加新的生成器

1. **在 `generators/types.go` 中添加选项**:
```go
type ProjectorOptions struct {
    Name      string
    Domain    string
    OutputDir string
}
```

2. **创建生成器实现**:
```go
// generators/projector_generator.go
type ProjectorGenerator struct {
    opts ProjectorOptions
}

func NewProjectorGenerator(opts ProjectorOptions) *ProjectorGenerator {
    return &ProjectorGenerator{opts: opts}
}

func (g *ProjectorGenerator) Generate() error {
    // 实现生成逻辑
}
```

3. **在 `command/generate.go` 中添加命令**:
```go
func generateProjectorCmd() *cobra.Command {
    var opts generators.ProjectorOptions
    
    cmd := &cobra.Command{
        Use:   "projector [name]",
        Short: "Generate projector for read model",
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.Name = args[0]
            generator := generators.NewProjectorGenerator(opts)
            return generator.Generate()
        },
    }
    
    return cmd
}
```

4. **注册命令**:
```go
// generate.go
func generateCmd() *cobra.Command {
    cmd.AddCommand(generateProjectorCmd())
    return cmd
}
```

## 总结

CLI 工具的目录结构遵循以下原则：

✅ **清晰的三层架构**: main → command → generators  
✅ **职责分离**: 每层只做一件事  
✅ **易于扩展**: 四步即可添加新命令  
✅ **类型安全**: 强类型系统保证正确性  
✅ **文档完善**: 每个文件都有清晰说明  

两种 DAO 生成方式：
- ✅ **`generate dao`** - 手动指定字段，适合设计阶段
- ✅ **`generate dao`** - 从数据库反向工程，适合已有数据库

开发者可以根据项目阶段和需求选择合适的生成方式。
