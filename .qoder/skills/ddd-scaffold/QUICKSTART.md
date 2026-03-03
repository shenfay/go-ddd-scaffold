# DDD Scaffold 快速开始指南

## 5 分钟快速上手

### 第一步：安装 Skill

```bash
npx skills install ddd-scaffold
```

### 第二步：生成项目结构

#### 方式 1: 使用默认配置（推荐新手）

```bash
/ddd-scaffold --project-name myapp --interactive
```

按照提示回答以下问题：
1. 项目名称：`myapp`
2. 需要哪些领域：`user,order,product`（逗号分隔）
3. 架构风格：`standard`（minimal/standard/full）
4. 是否需要示例代码：`y`
5. 是否需要测试文件：`y`
6. 是否需要 Docker 配置：`y`

#### 方式 2: 命令行一键生成

```bash
# 生成标准电商项目（包含用户、订单、商品、库存四个领域）
/ddd-scaffold \
  --project-name ecommerce \
  --domains user,order,product,inventory \
  --style standard \
  --with-examples \
  --with-tests \
  --with-docker \
  --output ./my-ecommerce-app
```

### 第三步：初始化项目

```bash
# 进入项目目录
cd my-ecommerce-app

# 下载依赖
go mod tidy

# 编译项目
make build

# 启动开发环境
make docker-up

# 运行数据库迁移
make migrate-up

# 启动应用
make run
```

访问 `http://localhost:8080/api/v1/users` 查看 API

## 生成的项目结构

```
my-ecommerce-app/
├── cmd/
│   ├── server/           # API 服务入口
│   │   └── main.go
│   └── worker/           # Worker 入口
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user/         # 用户领域
│   │   ├── order/        # 订单领域
│   │   ├── product/      # 商品领域
│   │   └── inventory/    # 库存领域
│   ├── application/      # 应用层
│   ├── infrastructure/   # 基础设施层
│   └── interfaces/       # 接口层
├── pkg/                  # 公共包
├── migrations/           # 数据库迁移
├── configs/              # 配置文件
├── go.mod
├── Makefile
└── README.md
```

## 下一步

### 1. 修改领域模型

编辑 `internal/domain/user/entity/user.go`：

```go
type User struct {
    ID        string
    Name      string
    Email     string
    // 添加你的自定义字段
}
```

### 2. 添加业务逻辑

编辑 `internal/domain/user/service/user_service.go`：

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // 添加你的业务逻辑
    return nil
}
```

### 3. 调用 API

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'

# 获取用户列表
curl http://localhost:8080/api/v1/users
```

## 常用命令

```bash
# 开发模式运行
make run

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 生成依赖注入代码
make generate

# 查看日志
docker-compose logs -f app
```

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：查看故障排除章节或咨询 DDD Architect Agent

## 推荐学习路径

1. ✅ 完成本快速开始（5 分钟）
2. 📚 阅读 [REFERENCE.md](./REFERENCE.md) 了解详细配置（15 分钟）
3. 🔍 研究生成的代码结构（30 分钟）
4. 🎯 参考 [EXAMPLES.md](./EXAMPLES.md) 添加新功能（1 小时）
5. 🤖 咨询 DDD Architect Agent 优化架构设计（按需）

祝你开发顺利！🚀
