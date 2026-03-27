# Go DDD Scaffold - CLI & Makefile 使用指南

本文档详细介绍 Go DDD Scaffold 项目的命令行工具（CLI）和 Makefile 的使用方法。

---

## 📋 目录

- [快速开始](#快速开始)
- [CLI 工具命令](#cli-工具命令)
- [Makefile 命令](#makefile-命令)
- [使用场景示例](#使用场景示例)
- [最佳实践](#最佳实践)

---

## 🚀 快速开始

### 安装 CLI 工具

```bash
# 构建并安装 CLI 到 GOPATH/bin
make cli-install

# 验证安装
go-ddd-scaffold version
```

### 设置开发环境

```bash
# 首次使用，设置开发环境
make setup
```

---

## 🔧 CLI 工具命令

CLI 工具提供项目初始化、代码生成、数据库迁移等高级功能。

### 命令总览

```bash
go-ddd-scaffold <command> [subcommand] [flags]
```

#### 可用命令

| 命令 | 说明 | 别名 |
|------|------|------|
| `init` | 项目初始化 | - |
| `generate` | 代码生成 | `gen`, `g` |
| `migrate` | 数据库迁移管理 | - |
| `docs` | 文档生成 | - |
| `clean` | 清理生成的文件 | - |
| `version` | 显示版本信息 | - |

---

### 1. generate - 代码生成

生成符合 DDD 和 Clean Architecture 模式的代码。

#### 1.1 生成领域实体

```bash
# 生成基础实体
go-ddd-scaffold generate entity User

# 生成实体及其值对象
go-ddd-scaffold generate entity User --with-vo

# 生成完整聚合根
go-ddd-scaffold generate entity User --with-aggregate --with-vo

# 指定字段和方法
go-ddd-scaffold generate entity User \
  -f "username:string,email:string,age:int" \
  -m "Validate,Activate,Deactivate" \
  -p "user"
```

**参数说明：**
- `-f, --fields`: 字段定义（格式：name:type,name:type）
- `-m, --methods`: 业务方法列表
- `-p, --package`: 领域包名
- `--with-vo`: 生成值对象
- `--with-aggregate`: 生成聚合根

#### 1.2 生成 DAO 层（从数据库）

```bash
# 从数据库生成所有核心表的 DAO
go-ddd-scaffold generate dao

# 从特定表生成 DAO
go-ddd-scaffold generate dao -t users,tenants,roles

# 自定义输出路径和配置
go-ddd-scaffold generate dao \
  -o internal/infrastructure/persistence/dao \
  --field-nullable \
  --with-test

# 使用自定义数据库连接
go-ddd-scaffold generate dao \
  -d "host=localhost user=postgres dbname=mydb"
```

**参数说明：**
- `-o, --output`: 输出目录（默认：internal/infrastructure/persistence/dao）
- `-d, --dsn`: 数据库连接 DSN
- `-c, --config`: 配置文件路径
- `--with-test`: 生成单元测试
- `--field-nullable`: 为可空字段生成指针类型（默认：true）
- `-t, --tables`: 指定表（逗号分隔）

#### 1.3 生成 Repository 层

```bash
# 生成用户仓储
go-ddd-scaffold generate repository UserRepository

# 指定领域和输出目录
go-ddd-scaffold generate repository UserRepository \
  -d user \
  -o internal/infrastructure/persistence/repository
```

**参数说明：**
- `-d, --domain`: 领域名称
- `-o, --output`: 输出目录

#### 1.4 生成应用服务

```bash
# 生成应用服务
go-ddd-scaffold generate service UserService

# 指定服务类型和方法
go-ddd-scaffold generate service UserService \
  -t application \
  -m "CreateUser,GetUser,UpdateUser,DeleteUser" \
  -d "UserRepository,EventPublisher"
```

**参数说明：**
- `-t, --type`: 服务类型（application/domain，默认：application）
- `-m, --methods`: 服务方法列表
- `-d, --deps`: 依赖注入列表（逗号分隔）

#### 1.5 生成 CQRS 处理器

```bash
# 生成命令处理器
go-ddd-scaffold generate handler CreateUserHandler -t command -d user

# 生成查询处理器
go-ddd-scaffold generate handler GetUserQueryHandler -t query -d user
```

**参数说明：**
- `-t, --type`: 处理器类型（command/query）
- `-d, --domain`: 领域名称

#### 1.6 生成 DTO

```bash
# 生成请求 DTO
go-ddd-scaffold generate dto CreateUserRequest -t request

# 生成响应 DTO
go-ddd-scaffold generate dto UserResponse -t response

# 指定字段和验证
go-ddd-scaffold generate dto CreateUserRequest \
  -f "Username:string,Email:string,Password:string" \
  --with-validation
```

**参数说明：**
- `-t, --type`: DTO 类型（request/response）
- `-f, --fields`: 字段定义
- `--with-validation`: 添加验证标签（默认：true）

---

### 2. migrate - 数据库迁移管理

使用 golang-migrate 管理数据库迁移。

#### 2.1 执行迁移

```bash
# 应用所有待处理的迁移
go-ddd-scaffold migrate up

# 应用指定数量的迁移
go-ddd-scaffold migrate up --steps 5
```

#### 2.2 回滚迁移

```bash
# 回滚 1 个迁移（默认）
go-ddd-scaffold migrate down

# 回滚多个迁移
go-ddd-scaffold migrate down --steps 3
```

#### 2.3 创建新迁移

```bash
# 创建新迁移文件
go-ddd-scaffold migrate create add_user_profile

# 会生成两个文件：
# - migrations/000010_add_user_profile.up.sql
# - migrations/000010_add_user_profile.down.sql
```

#### 2.4 查看迁移状态

```bash
# 查看当前迁移状态和待处理的迁移
go-ddd-scaffold migrate status

# 查看当前数据库版本
go-ddd-scaffold migrate version
```

#### 2.5 高级选项

```bash
# 使用自定义迁移目录
go-ddd-scaffold migrate up --path ./custom/migrations

# 使用自定义数据库连接
go-ddd-scaffold migrate up --dsn "postgres://user:pass@localhost/dbname?sslmode=disable"
```

---

### 3. docs - 文档生成

#### 3.1 生成 Swagger API 文档

```bash
# 生成 Swagger 文档（默认输出到 ./api/swagger）
go-ddd-scaffold docs swagger

# 自定义输出目录
go-ddd-scaffold docs swagger -o ./docs/api

# 指定项目目录
go-ddd-scaffold docs swagger -d /path/to/project
```

**参数说明：**
- `-d, --dir`: 项目根目录（默认：.）
- `-o, --output`: 输出目录（默认：./api/swagger）

---

### 4. clean - 清理生成的文件

```bash
# 清理当前目录的生成文件
go-ddd-scaffold clean

# 清理指定目录
go-ddd-scaffold clean ./internal

# 预览将删除的文件
go-ddd-scaffold clean --dry-run
```

---

### 5. version - 版本信息

```bash
# 显示版本信息
go-ddd-scaffold version
```

---

## 📝 Makefile 命令

Makefile 提供日常开发的快捷命令。

### 开发相关

```bash
make run              # 启动 API 服务（开发模式）
make run-worker       # 启动 Worker（开发模式）
make setup            # 设置开发环境
make install-deps     # 安装依赖
```

### CLI 工具构建

```bash
make cli              # 构建当前系统的 CLI 工具
make cli-linux        # 构建 Linux 版 CLI
make cli-darwin       # 构建 macOS 版 CLI
make cli-windows      # 构建 Windows 版 CLI
make cli-install      # 安装 CLI 到 GOPATH/bin
make cli-test         # 测试 CLI 工具
make cli-clean        # 清理 CLI 二进制文件
```

### 构建编译

```bash
make build            # 构建当前系统的应用
make build-worker     # 构建 Worker
make build-linux      # 构建 Linux 生产版本
make clean            # 清理构建产物
```

### 测试

```bash
make test             # 运行所有测试
make test-short       # 运行简短测试（跳过集成测试）
make coverage         # 生成测试覆盖率报告
```

### 代码质量

```bash
make fmt              # 格式化代码
make vet              # 运行 go vet
make lint             # 运行代码检查
```

### 数据库

```bash
make migrate-up       # 执行数据库迁移
make migrate-down     # 回滚数据库迁移
```

### asynqmon 监控

```bash
make asynqmon-install   # 安装 asynqmon 工具
make asynqmon           # 启动 asynqmon UI (端口 8080)
make asynqmon-port PORT=8081  # 自定义端口启动
make asynqmon-ui        # 在浏览器打开 asynqmon (仅 macOS)
```

### 工具

```bash
make health           # 检查应用健康状态
```

### 文档

```bash
make swagger-gen      # 生成 Swagger 文档
make swagger-serve    # 生成并启动 Swagger UI 服务
```

---

## 💡 使用场景示例

### 场景 1：新项目初始化

```bash
# 1. 设置开发环境
make setup

# 2. 安装 CLI 工具
make cli-install

# 3. 执行数据库迁移
make migrate-up

# 4. 生成 Swagger 文档
make swagger-gen

# 5. 启动开发服务器
make run
```

### 场景 2：添加新功能模块

```bash
# 1. 创建数据库迁移
go-ddd-scaffold migrate create add_product_catalog

# 2. 编辑迁移文件并执行
make migrate-up

# 3. 从数据库生成 DAO
go-ddd-scaffold generate dao -t products,product_categories --field-nullable

# 4. 手动编写领域模型（DDD 实践的核心）
# 编辑 domain/product/aggregate/product.go

# 5. 生成 Repository 层
go-ddd-scaffold generate repository ProductRepository -d product

# 6. 生成应用服务
go-ddd-scaffold generate service ProductService \
  -m "CreateProduct,GetProduct,UpdateProduct,ListProducts"

# 7. 生成 HTTP Handler
go-ddd-scaffold generate handler ProductHandler -d product

# 8. 生成 DTO
go-ddd-scaffold generate dto CreateProductRequest -t request
go-ddd-scaffold generate dto ProductResponse -t response

# 9. 更新 API 文档
make swagger-gen

# 10. 运行测试
make test
```

### 场景 3：日常开发流程

```bash
# 早上开始工作
make run              # 启动开发服务器

# 开发过程中
make fmt              # 格式化代码
make lint             # 代码检查
make test-short       # 运行快速测试

# 提交前检查
make vet
make lint
make test

# 生成覆盖率报告
make coverage
```

### 场景 4：生产部署

```bash
# 1. 构建生产版本
make build-linux

# 2. 构建 Worker
make build-worker-linux

# 3. 执行数据库迁移
./bin/go-ddd-scaffold migrate up

# 4. 健康检查
make health
```

---

## 🎯 最佳实践

### 1. CLI vs Makefile 的选择

| 场景 | 推荐方式 | 原因 |
|------|---------|------|
| 日常开发 | **Makefile** | 简单、直观、无需记忆参数 |
| 代码生成 | **CLI** | 需要灵活配置和交互 |
| 数据库迁移 | **两者结合** | Makefile 常规操作，CLI 高级功能 |
| 运维监控 | **CLI** | 提供详细状态和诊断 |

### 2. 代码生成建议

✅ **推荐做法：**
```bash
# 只生成基础设施代码
go-ddd-scaffold generate dao -t users

# 手动编写业务逻辑
# 编辑 domain/user/aggregate/user.go
```

❌ **不推荐：**
```bash
# 过度依赖自动化生成
go-ddd-scaffold generate entity User --with-everything
```

**理由：** DDD 的核心是领域建模，应该手动编写领域逻辑，而不是自动生成。

### 3. 数据库迁移规范

```bash
# ✅ 好的迁移命名
go-ddd-scaffold migrate create add_user_email
go-ddd-scaffold migrate create create_products_table

# ❌ 避免模糊的命名
go-ddd-scaffold migrate create update
go-ddd-scaffold migrate create fix
```

### 4. 开发工作流

```bash
# 功能开发流程
1. make migrate-create    # 创建迁移
2. 编辑迁移文件
3. make migrate-up        # 执行迁移
4. go-ddd-scaffold generate dao  # 生成基础设施
5. 手动编写领域逻辑      # DDD 核心
6. make test             # 测试验证
7. make swagger-gen      # 更新文档
```

### 5. 环境变量配置

建议在 `.env` 文件中配置数据库连接：

```bash
# .env
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=go_ddd_scaffold
DATABASE_SSL_MODE=disable
```

CLI 工具会自动读取这些环境变量。

---

## 🔗 相关资源

- [项目 README](../../README.md) - 项目介绍和快速开始
- [API 文档](../api/README.md) - API 接口文档
- [开发指南](../docs/DEVELOPMENT.md) - 详细开发规范

---

## ❓ 常见问题

### Q: CLI 工具无法找到？
```bash
# 确保已安装到 GOPATH/bin
make cli-install

# 检查 GOPATH 是否在 PATH 中
echo $GOPATH/bin
```

### Q: DAO 生成失败？
```bash
# 检查数据库连接
go-ddd-scaffold migrate status

# 检查配置文件
cat configs/config.yaml
```

### Q: Swagger 文档未生成？
```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 重新生成
make swagger-gen
```

---

## 📞 获取帮助

```bash
# 查看 CLI 帮助
go-ddd-scaffold --help
go-ddd-scaffold generate --help
go-ddd-scaffold migrate --help

# 查看 Makefile 帮助
make help
```
