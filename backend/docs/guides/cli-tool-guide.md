# Go DDD Scaffold CLI工具使用文档

## 文档概述

本文档详细介绍了 go-ddd-scaffold 项目的CLI工具使用方法，包括安装配置、项目初始化、代码生成、文档生成等功能。

## CLI工具安装与配置

### 安装方式

#### 1. Go Install方式（推荐）
```bash
# 安装最新版本
go install github.com/your-org/go-ddd-scaffold@latest

# 或安装指定版本
go install github.com/your-org/go-ddd-scaffold@v1.0.0
```

#### 2. 源码编译方式
```bash
# 克隆源码
git clone https://github.com/your-org/go-ddd-scaffold.git
cd go-ddd-scaffold

# 编译安装
go build -o go-ddd-scaffold cmd/scaffold/main.go
sudo cp go-ddd-scaffold /usr/local/bin/
```

#### 3. 下载预编译二进制文件
```bash
# 从GitHub Releases下载对应平台的二进制文件
wget https://github.com/your-org/go-ddd-scaffold/releases/download/v1.0.0/go-ddd-scaffold-linux-amd64
chmod +x go-ddd-scaffold-linux-amd64
sudo mv go-ddd-scaffold-linux-amd64 /usr/local/bin/go-ddd-scaffold
```

### 环境配置

#### 配置文件位置
CLI工具会在以下位置查找配置文件：
1. `$HOME/.go-ddd-scaffold/config.yaml`
2. 当前工作目录的 `.scaffold.yaml`
3. 环境变量

#### 基本配置文件示例
```yaml
# ~/.go-ddd-scaffold/config.yaml
defaults:
  author: "Your Name"
  email: "your.email@example.com"
  license: "MIT"
  go_module: "github.com/your-org"

templates:
  path: "~/.go-ddd-scaffold/templates"
  default: "clean-architecture"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "${DB_PASSWORD}"

features:
  enable_frontend: true
  enable_docker: true
  enable_kubernetes: false
```

## 核心命令详解

### 1. 项目初始化命令

#### 基本用法
```bash
# 交互式初始化
go-ddd-scaffold init my-project

# 跳过前端初始化
go-ddd-scaffold init my-project --skip-frontend

# 指定模板
go-ddd-scaffold init my-project --template clean-architecture
```

#### 交互式初始化流程
```bash
$ go-ddd-scaffold init my-new-project

Welcome to Go DDD Scaffold!
? Project name: my-new-project
? Project description: My awesome enterprise application
? Author name: John Doe
? Author email: john@example.com
? Go module path: github.com/johndoe/my-new-project
? Database type: PostgreSQL
? Enable Redis cache: Yes
? Enable frontend (React): Yes
? Enable Docker support: Yes
? License type: MIT

Creating project structure...
✓ Created backend directory
✓ Created frontend directory
✓ Generated configuration files
✓ Initialized git repository
✓ Created initial commit

Project created successfully!
Next steps:
  cd my-new-project
  go mod tidy
  docker-compose up -d  # if Docker enabled
  go run cmd/server/main.go
```

#### 命令行参数详解
```bash
Usage:
  go-ddd-scaffold init [flags] PROJECT_NAME

Flags:
  -t, --template string     Template to use (default "clean-architecture")
  -s, --skip-frontend       Skip frontend initialization
  -d, --database string     Database type (postgresql/mysql) (default "postgresql")
  -r, --with-redis          Include Redis cache support
  -c, --config string       Path to config file
  -v, --verbose             Verbose output
  -h, --help                Help for init command
```

### 2. 代码生成命令

#### 实体代码生成
```bash
# 生成用户实体
go-ddd-scaffold generate entity user \
  --fields="username:string,email:string,password:string,status:int" \
  --methods="Validate,HashPassword" \
  --table-name="users"

# 生成租户实体
go-ddd-scaffold generate entity tenant \
  --fields="code:string,name:string,description:string,status:int" \
  --relationships="users:many" \
  --table-name="tenants"
```

#### CRUD代码生成
```bash
# 为用户实体生成完整CRUD
go-ddd-scaffold generate crud user \
  --with-validation \
  --with-pagination \
  --with-search

# 生成API接口
go-ddd-scaffold generate api user \
  --endpoints="list,get,create,update,delete" \
  --auth-required
```

#### 领域服务生成
```bash
# 生成认证服务
go-ddd-scaffold generate service auth \
  --methods="Login,Logout,RefreshToken,ValidateToken" \
  --dependencies="userRepository,jwtService"

# 生成权限服务
go-ddd-scaffold generate service authorization \
  --methods="CheckPermission,GrantPermission,RevokePermission" \
  --dependencies="roleRepository,permissionRepository"
```

### 3. 文档生成命令

#### API文档生成
```bash
# 从代码注释生成Swagger文档
go-ddd-scaffold generate docs api \
  --format swagger \
  --output docs/swagger.yaml

# 生成Markdown API文档
go-ddd-scaffold generate docs api \
  --format markdown \
  --output docs/api-docs.md
```

#### 数据库文档生成
```bash
# 生成数据库ER图
go-ddd-scaffold generate docs database \
  --format mermaid \
  --output docs/er-diagram.mmd

# 生成数据库表结构文档
go-ddd-scaffold generate docs database \
  --format markdown \
  --output docs/database-schema.md
```

#### 项目架构文档生成
```bash
# 生成架构设计文档
go-ddd-scaffold generate docs architecture \
  --include-domains \
  --include-services \
  --output docs/architecture.md
```

### 4. 迁移管理命令

#### 数据库迁移操作
```bash
# 创建新的迁移文件
go-ddd-scaffold migrate create add_user_profiles \
  --type sql

# 应用所有待处理迁移
go-ddd-scaffold migrate up

# 回滚最近一次迁移
go-ddd-scaffold migrate down

# 查看迁移状态
go-ddd-scaffold migrate status

# 重置数据库（谨慎使用）
go-ddd-scaffold migrate reset
```

#### 迁移文件管理
```bash
# 生成迁移文件模板
go-ddd-scaffold migrate template create_users_table \
  --fields="id:bigint,username:string,email:string,password:string"

# 验证迁移文件语法
go-ddd-scaffold migrate validate

# 打包迁移文件
go-ddd-scaffold migrate pack \
  --output migrations.tar.gz
```

## Skills插件系统

### 内置Skills列表

#### 核心Skills（6个）
1. **init** - 项目初始化Skill
2. **entity** - 实体生成Skill
3. **crud** - CRUD代码生成Skill
4. **service** - 服务层代码生成Skill
5. **migration** - 数据库迁移Skill
6. **docs** - 文档生成Skill

#### 业务Skills（5个）
1. **auth** - 认证授权Skill
2. **rbac** - RBAC权限Skill
3. **tenant** - 多租户Skill
4. **audit** - 审计日志Skill
5. **notification** - 通知服务Skill（预留）

### 自定义Skill开发

#### Skill接口定义
```go
type Skill interface {
    Name() string
    Description() string
    Run(ctx *SkillContext) error
    Configure(config *viper.Viper) error
}

type SkillContext struct {
    Args    []string
    Flags   map[string]string
    Config  *viper.Viper
    Writer  io.Writer
    WorkingDir string
}
```

#### 开发自定义Skill示例
```go
// skills/custom/hello_skill.go
package custom

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/your-org/go-ddd-scaffold/pkg/skills"
)

type HelloSkill struct{}

func (s *HelloSkill) Name() string {
    return "hello"
}

func (s *HelloSkill) Description() string {
    return "Say hello with customization"
}

func (s *HelloSkill) Run(ctx *skills.SkillContext) error {
    name := ctx.Flags["name"]
    if name == "" {
        name = "World"
    }
    
    fmt.Fprintf(ctx.Writer, "Hello, %s!\n", name)
    return nil
}

func (s *HelloSkill) Configure(v *viper.Viper) error {
    // 配置Skill参数
    return nil
}

// 注册Skill
func init() {
    skills.Register(&HelloSkill{})
}
```

#### 使用自定义Skill
```bash
# 使用自定义Skill
go-ddd-scaffold hello --name="Go Developer"

# 输出: Hello, Go Developer!
```

## 配置管理

### 环境变量支持
```bash
# 数据库配置
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=secret

# Redis配置
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=redis_secret

# JWT配置
export JWT_SECRET=my_jwt_secret_key
export JWT_ACCESS_EXPIRE=30m
export JWT_REFRESH_EXPIRE=7d
```

### 多环境配置
```yaml
# config/development.yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  name: scaffold_dev

# config/production.yaml
server:
  port: 80
  mode: release

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
```

### 配置覆盖优先级
1. 命令行标志（最高优先级）
2. 环境变量
3. 配置文件
4. 默认值（最低优先级）

## 模板系统

### 内置模板
```
templates/
├── clean-architecture/     # 标准Clean Architecture模板
├── microservice/          # 微服务模板
├── minimal/               # 最小化模板
└── enterprise/            # 企业级完整模板
```

### 自定义模板开发
```go
// 创建自定义模板结构
my-template/
├── scaffold.yaml          # 模板配置文件
├── backend/               # 后端模板
│   ├── cmd/
│   ├── internal/
│   └── configs/
├── frontend/              # 前端模板
│   ├── src/
│   └── public/
└── docker/                # Docker模板
    ├── Dockerfile
    └── docker-compose.yml
```

#### 模板配置文件示例
```yaml
# scaffold.yaml
name: "my-custom-template"
version: "1.0.0"
description: "My custom project template"
author: "Your Name"

variables:
  - name: project_name
    prompt: "Project name"
    default: "my-project"
  - name: module_path
    prompt: "Go module path"
    default: "github.com/username/project"
  - name: with_docker
    prompt: "Include Docker support?"
    type: bool
    default: true

file_templates:
  - source: "backend/cmd/server/main.go.tmpl"
    target: "backend/cmd/server/main.go"
  - source: "backend/internal/domain/{{.project_name}}/entity.go.tmpl"
    target: "backend/internal/domain/{{.project_name}}/entity.go"

post_hooks:
  - command: "go mod init {{.module_path}}"
    working_dir: "backend"
  - command: "go mod tidy"
    working_dir: "backend"
```

## 高级用法

### 批量生成
```bash
# 批量生成多个实体
go-ddd-scaffold batch generate \
  --entities="user,tenant,role,permission" \
  --template-file=entities.yaml

# entities.yaml内容示例
entities:
  - name: user
    fields:
      - name: username
        type: string
      - name: email
        type: string
    methods:
      - Validate
      - HashPassword
  
  - name: tenant
    fields:
      - name: code
        type: string
      - name: name
        type: string
    relationships:
      - users:many
```

### 项目模板导出/导入
```bash
# 导出当前项目为模板
go-ddd-scaffold template export my-project-template \
  --exclude="*.log,*.tmp,node_modules" \
  --output=my-template.tar.gz

# 从模板创建项目
go-ddd-scaffold template import my-template.tar.gz my-new-project
```

### 代码质量检查集成
```bash
# 生成代码后自动运行检查
go-ddd-scaffold generate entity user --run-lint --run-test

# 配置文件中启用自动检查
code_quality:
  run_lint: true
  run_tests: true
  check_coverage: 80
```

## 故障排除

### 常见问题及解决方案

#### 1. 权限问题
```bash
# 问题：Permission denied when creating files
# 解决：检查目录权限
sudo chown -R $USER:$USER /path/to/project
```

#### 2. 依赖安装失败
```bash
# 问题：Go modules download failed
# 解决：清理缓存并重试
go clean -modcache
go mod tidy
```

#### 3. 数据库连接失败
```bash
# 问题：Cannot connect to database
# 解决：检查数据库服务状态
docker-compose up -d postgres
# 或检查连接参数
go-ddd-scaffold migrate status --debug
```

#### 4. 模板渲染错误
```bash
# 问题：Template parsing error
# 解决：验证模板语法
go-ddd-scaffold template validate /path/to/template
```

### 调试模式
```bash
# 启用详细输出
go-ddd-scaffold init my-project --verbose --debug

# 查看执行日志
go-ddd-scaffold --log-level=debug generate entity user

# 生成诊断报告
go-ddd-scaffold diagnose > diagnosis-report.txt
```

这个CLI工具使用文档为用户提供了完整的工具使用指南和最佳实践。