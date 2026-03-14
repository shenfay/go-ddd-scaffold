# Go DDD Scaffold CLI 工具使用文档

## 简介

`go-ddd-scaffold` 是一个专为 DDD（领域驱动设计）项目设计的代码生成 CLI 工具。它遵循 Clean Architecture 和 CQRS 模式，可以快速生成符合最佳实践的代码结构。

## 安装方式

### 方式一：从源码构建（推荐开发使用）

```bash
cd backend
make cli          # 构建当前系统的 CLI
make cli-install  # 安装到 GOPATH/bin
```

### 方式二：Go Install

```bash
go install github.com/shenfay/go-ddd-scaffold/cmd/cli@latest
```

### 方式三：下载预编译二进制

从 GitHub Releases 下载对应平台的二进制文件，添加到 PATH 即可。

## 快速开始

### 1. 查看帮助

```bash
# 查看主帮助
go-ddd-scaffold --help

# 查看所有命令
go-ddd-scaffold version

# 查看 generate 命令帮助
go-ddd-scaffold generate --help
```

### 2. 生成 DAO 层代码

```bash
# 基本用法
go-ddd-scaffold generate dao user \
  -f "username:string,email:string,password:string,status:int"

# 指定表名
go-ddd-scaffold generate dao user \
  -f "username:string,email:string" \
  -t users

# 指定输出目录
go-ddd-scaffold generate dao user \
  -f "username:string,email:string" \
  -o internal/infrastructure/dao

# 完整示例
go-ddd-scaffold generate dao user \
  -f "username:string,email:string,password:string,status:int,created_at:time" \
  -t users \
  -o internal/infrastructure/dao
```

**生成的文件**:
- `user_dao.go` - DAO 接口定义
- `user_dao_impl.go` - DAO 实现
- `user_model.go` - 数据模型和查询条件

### 3. 其他生成命令（待实现）

```bash
# 生成领域实体
go-ddd-scaffold generate entity user \
  -f "username:string,email:string" \
  --with-vo \
  --with-aggregate

# 生成 Repository
go-ddd-scaffold generate repository user \
  --domain user

# 生成应用服务
go-ddd-scaffold generate service user \
  --methods "Create,Get,Update,Delete"

# 生成 CQRS Handler
go-ddd-scaffold generate handler CreateUser \
  --type command \
  --domain user

# 生成 DTO
go-ddd-scaffold generate dto CreateUserRequest \
  --type request \
  -f "username:string,email:string,password:string"
```

## 命令详解

### 全局标志

```bash
-c, --config      # 配置文件路径 (默认：$HOME/.go-ddd-scaffold.yaml)
-v, --verbose     # 详细输出
-n, --dry-run     # 预览执行，不实际写入文件
```

### init - 初始化新项目

```bash
go-ddd-scaffold init my-project \
  --module-path github.com/username/my-project \
  --author "Your Name" \
  --email "your.email@example.com" \
  --license MIT \
  --template clean-architecture \
  --with-docker \
  --skip-frontend
```

**参数说明**:
- `--module-path`: Go module 路径
- `--author`: 作者名
- `--email`: 作者邮箱
- `--license`: 许可证类型 (默认：MIT)
- `--template`: 项目模板 (默认：clean-architecture)
- `--with-docker`: 包含 Docker 配置
- `--with-k8s`: 包含 Kubernetes manifests
- `--skip-frontend`: 跳过前端初始化

### generate dao - 生成 DAO 层

**字段格式**: `name:type,name:type`

支持的类型映射:
- `string` → `string`
- `int` / `integer` → `int64`
- `bool` / `boolean` → `bool`
- `time` / `datetime` / `timestamp` → `time.Time`
- `float` / `decimal` → `float64`
- `json` → `[]byte`

**示例**:

```bash
# 用户实体
go-ddd-scaffold generate dao user \
  -f "username:string,email:string,password:string,status:int"

# 订单实体（包含时间字段）
go-ddd-scaffold generate dao order \
  -f "order_no:string,user_id:int,total_amount:decimal,status:int,paid_at:time"

# 产品实体
go-ddd-scaffold generate dao product \
  -f "name:string,description:text,price:decimal,stock:int,is_active:bool"
```

### migrate - 数据库迁移管理

```bash
# 运行所有待处理迁移
go-ddd-scaffold migrate up

# 回滚最后一次迁移
go-ddd-scaffold migrate down

# 创建新迁移
go-ddd-scaffold migrate create add_users_table

# 查看迁移状态
go-ddd-scaffold migrate status
```

### docs - 生成文档

```bash
# 生成 Swagger 文档
go-ddd-scaffold docs swagger

# 生成 API 文档
go-ddd-scaffold docs api
```

### clean - 清理生成的文件

```bash
# 清理所有生成的文件
go-ddd-scaffold clean

# 清理指定目录
go-ddd-scaffold clean internal/infrastructure/dao

# 预览删除（不实际执行）
go-ddd-scaffold clean --dry-run
```

### version - 显示版本信息

```bash
go-ddd-scaffold version
```

## 配置文件

创建 `~/.go-ddd-scaffold.yaml` 配置文件：

```yaml
defaults:
  author: "Your Name"
  email: "your.email@example.com"
  license: "MIT"
  
generator:
  dao:
    output_dir: "internal/infrastructure/dao"
    with_interface: true
    with_condition: true
    
  entity:
    with_vo: false
    with_aggregate: true
```

## 最佳实践

### 1. 命名约定

- **实体名**: 使用单数形式（`user`, `order`）
- **表名**: 使用复数形式（`users`, `orders`）
- **字段名**: 使用驼峰命名（`username`, `createdAt`）
- **列名**: 自动生成蛇形命名（`username`, `created_at`）

### 2. 字段定义技巧

```bash
# 常用字段组合
# 基础字段
-f "name:string,code:string,status:int"

# 审计字段（自动生成 created_at, updated_at）
-f "name:string"  # 无需手动指定时间字段

# 金额字段使用 decimal
-f "price:decimal,total_amount:decimal"

# JSON 字段
-f "metadata:json,config:json"
```

### 3. 分层生成顺序

推荐的代码生成顺序：

1. **领域层**: `generate entity` (定义领域模型)
2. **基础设施层**: `generate dao` (数据访问层)
3. **领域层**: `generate repository` (仓储层)
4. **应用层**: `generate service` (应用服务)
5. **应用层**: `generate handler` (CQRS 处理器)
6. **接口层**: `generate dto` (数据传输对象)

### 4. 与现有代码集成

生成的 DAO 可以直接注入到 Repository 中使用：

```go
// internal/domain/user/repository.go
type userRepository struct {
    dao *dao.UserDAOImpl
}

func NewUserRepository(db *sql.DB) *userRepository {
    return &userRepository{
        dao: dao.NewUserDAOImpl(db),
    }
}
```

## 架构设计

CLI 工具采用三层架构：

```
main.go (入口)
  ↓
command 层 (命令解析)
  ↓
generators 层 (代码生成)
```

详细架构设计请参考：[CLI 架构设计文档](docs/architecture/cli-design.md)

## 扩展开发

### 添加新的生成器

1. 在 `generators/types.go` 中定义选项
2. 创建生成器实现文件
3. 在 `command/generate.go` 中添加命令
4. 在 `command/root.go` 中注册

详见：[CLI 架构设计文档 - 扩展机制](docs/architecture/cli-design.md#扩展机制)

## 故障排除

### 常见问题

#### 1. 权限错误

```bash
# 问题：Permission denied when creating files
# 解决：检查目录权限或使用 sudo
sudo chown -R $USER:$USER /path/to/project
```

#### 2. 模板解析错误

```bash
# 问题：Template parsing error
# 解决：检查字段格式是否正确
# 错误：-f "username:string email:string"  (缺少逗号)
# 正确：-f "username:string,email:string"
```

#### 3. 导入错误

```bash
# 问题：package not found
# 解决：运行 go mod tidy
go mod tidy
```

### 调试模式

```bash
# 启用详细输出
go-ddd-scaffold generate dao user -v

# 预览执行
go-ddd-scaffold generate dao user -n
```

## Makefile 命令

```bash
# 构建 CLI
make cli

# 构建 Linux 版本
make cli-linux

# 安装到 GOPATH/bin
make cli-install

# 测试 CLI
make cli-test

# 清理构建产物
make cli-clean
```

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
