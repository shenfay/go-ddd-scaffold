# Go DDD Scaffold - 企业级 DDD 架构脚手架

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://github.com/your-org/ddd-scaffold/actions/workflows/ci.yml/badge.svg)](https://github.com/your-org/ddd-scaffold/actions)
[![Coverage Status](https://coveralls.io/repos/github/your-org/ddd-scaffold/badge.svg)](https://coveralls.io/github/your-org/ddd-scaffold)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/ddd-scaffold)](https://goreportcard.com/report/github.com/your-org/ddd-scaffold)

**Go DDD Scaffold** 是一个基于领域驱动设计（DDD）的企业级脚手架项目，采用 Clean Architecture + DDD 模式，提供完整的微服务基础设施。

---

## ✨ 核心特性

### 🏗️ 架构设计

- **DDD 领域驱动设计**: 严格的四层架构（Interfaces / Application / Domain / Infrastructure）
- **Clean Architecture**: 依赖倒置、边界清晰、易于测试
- **CQRS 模式**: 命令查询职责分离，支持读写优化
- **UnitOfWork 模式**: 事务一致性保证
- **Repository 模式**: 数据访问抽象

### 🔐 认证授权

- **JWT Token 认证**: 支持 Access Token + Refresh Token
- **多租户支持**: Tenant-based SaaS 架构
- **RBAC 权限控制**: 基于 Casbin 的细粒度权限管理
- **Token 黑名单**: Redis 缓存、限流熔断保护

### 📦 领域建模

- **聚合根设计**: 清晰的聚合边界
- **值对象**: 类型安全的值对象封装
- **领域服务**: 业务逻辑封装
- **领域事件**: 事件驱动架构

### 🚀 基础设施

- **数据库连接池**: 生产环境优化配置
- **Redis 集成**: 缓存、黑名单、EventBus
- **EventBus**: Redis Stream 持久化、重试机制
- **监控埋点**: Prometheus Metrics（15+ 核心指标）
- **限流熔断**: 令牌桶限流 + 三态熔断器

### 📊 可观测性

- **Prometheus**: 完整的监控指标体系
- **Grafana**: 可视化仪表盘
- **链路追踪**: Jaeger/Zipkin 集成
- **日志聚合**: ELK Stack 支持

### 🧪 测试支持

- **单元测试**: 38+ 核心测试用例
- **集成测试**: 端到端场景测试
- **基准测试**: 性能基准对比
- **Mock 框架**: gomock 支持

---

## 📁 项目结构

```
ddd-scaffold/
├── backend/                          # 后端服务
│   ├── cmd/                         # 应用入口
│   │   ├── server/                  # HTTP 服务器
│   │   └── migrate/                 # 数据库迁移
│   ├── internal/                    # 内部代码（不对外暴露）
│   │   ├── application/             # 应用层（CQRS Service）
│   │   │   └── user/
│   │   │       ├── dto/             # 数据传输对象
│   │   │       └── service/         # 应用服务实现
│   │   ├── domain/                  # 领域层（核心业务逻辑）
│   │   │   ├── user/                # 用户领域
│   │   │   │   ├── entity/          # 实体
│   │   │   │   ├── valueobject/     # 值对象
│   │   │   │   ├── repository/      # 仓储接口
│   │   │   │   └── event/           # 领域事件
│   │   │   └── tenant/              # 租户领域
│   │   ├── infrastructure/          # 基础设施层
│   │   │   ├── persistence/         # 持久化实现
│   │   │   ├── redis/               # Redis 客户端
│   │   │   ├── auth/                # 认证服务
│   │   │   ├── middleware/          # 中间件
│   │   │   └── wire/                # 依赖注入
│   │   ├── interfaces/              # 接口层（适配器）
│   │   │   └── http/                # HTTP Handler
│   │   └── pkg/                     # 通用工具包
│   │       ├── metrics/             # Prometheus 指标
│   │       └── ratelimit/           # 限流熔断
│   ├── config/                      # 配置文件
│   ├── migrations/                  # 数据库迁移脚本
│   ├── tests/                       # 测试代码
│   │   ├── unit/                    # 单元测试
│   │   └── integration/             # 集成测试
│   └── docs/                        # API 文档
├── frontend/                        # 前端服务（可选）
│   ├── src/
│   │   ├── business/                # 业务层
│   │   ├── data/                    # 数据层
│   │   ├── interaction/             # 交互层
│   │   └── presentation/            # 表现层
│   └── public/
└── docs/                            # 项目文档
    ├── deployment_guide.md          # 部署指南
    ├── monitoring_and_ratelimit.md  # 监控与限流
    └── api_reference.md             # API 参考
```

---

## 🚀 快速开始

### 前置要求

- Go >= 1.21
- PostgreSQL >= 13
- Redis >= 6.0
- Node.js >= 18 (可选，仅前端)

### 1. 克隆项目

```bash
git clone https://github.com/your-org/ddd-scaffold.git
cd ddd-scaffold
```

### 2. 启动基础设施

```bash
# 使用 Docker Compose
cd deployments/docker
docker-compose up -d postgres redis
```

### 3. 初始化数据库

```bash
cd backend
go run cmd/migrate/main.go up
```

### 4. 配置服务

编辑 `backend/config/config.yaml`，然后启动：

```bash
cd backend
go run cmd/server/main.go
```

访问：http://localhost:8080

### 5. 查看 API 文档

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go -o ./docs

# 访问 Swagger UI
http://localhost:8080/swagger/index.html
```

---

## 📖 核心功能示例

### 用户注册登录

```bash
# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "StrongPass123",
    "nickname": "Test User"
  }'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "StrongPass123"
  }'

# 响应示例
{
  "code": "Success",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "nickname": "Test User"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 创建租户

```bash
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "name": "My Company",
    "maxMembers": 50
  }'
```

### 受保护的 API 调用

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}"
```

---

## 🛠️ 开发指南

### 常用命令

```bash
# 查看所有可用命令
make help

# 运行开发服务器（热重载）
make dev

# 运行测试
make test

# 生成 Swagger 文档
make swagger

# 构建 Docker 镜像
make docker-build

# 代码格式化
make fmt

# 代码检查
make lint
```

### 添加新领域

1. **创建领域目录**
```bash
mkdir -p internal/domain/product/{entity,valueobject,repository,event,service}
```

2. **定义实体**
```go
// internal/domain/product/entity/product.go
type Product struct {
    ID          uuid.UUID
    Name        string
    Price       Money
    Description string
}
```

3. **定义 Repository 接口**
```go
// internal/domain/product/repository/product_repository.go
type ProductRepository interface {
    Create(ctx context.Context, product *Product) error
    GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
}
```

4. **实现 Infrastructure**
```go
// internal/infrastructure/persistence/product_repository_impl.go
type productRepositoryImpl struct {
    db *gorm.DB
}
```

5. **创建应用服务**
```go
// internal/application/product/service/product_service.go
type ProductService interface {
    CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.Product, error)
}
```

6. **添加 HTTP Handler**
```go
// internal/interfaces/http/product/handler.go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // 处理请求
}
```

7. **注册路由**
```go
// cmd/server/main.go
api.GET("/products", productHandler.ListProducts)
api.POST("/products", productHandler.CreateProduct)
```

---

## 📊 性能指标

### 基准测试

| 操作 | QPS | P99 延迟 | 说明 |
|------|-----|---------|------|
| **JWT 验证** | ~5000 | < 2ms | 单次验证 |
| **黑名单检查** | ~3000 | < 5ms | Redis EXISTS |
| **批量黑名单** | ~10000 | < 10ms | Pipeline 100 个 |
| **数据库查询** | ~2000 | < 20ms | 带索引查询 |

### 资源占用

| 组件 | 内存 | CPU | 磁盘 |
|------|------|-----|------|
| **Backend** | ~100MB | ~0.5 Core | ~50MB |
| **PostgreSQL** | ~500MB | ~1 Core | ~1GB |
| **Redis** | ~50MB | ~0.2 Core | ~100MB |

---

## 🔒 安全特性

### JWT Token 黑名单机制

```go
// 登出时自动加入黑名单
POST /api/auth/logout
Authorization: Bearer {token}

// 中间件自动检查黑名单
if isBlacklisted {
    return 401 Unauthorized
}
```

### 限流熔断保护

```yaml
# 配置示例
ratelimit:
  rate: 100        # 每秒 100 次请求
  burst: 200       # 突发容量 200
  
circuitbreaker:
  max_failures: 5           # 5 次失败触发熔断
  reset_timeout: 30s        # 30 秒后尝试恢复
  half_open_calls: 3        # 半开状态允许 3 次调用
```

### RBAC 权限控制

```conf
# Casbin 策略
p, admin, tenant1, users, read
p, admin, tenant1, users, write
p, member, tenant1, users, read

g, user1, admin, tenant1
```

---

## 📈 监控告警

### Prometheus 指标

```prometheus
# QPS
rate(http_requests_total[1m])

# P99 延迟
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
rate(http_requests_total{status=~"5.."}[5m])

# 熔断器状态
circuit_breaker_state{resource="redis"}
```

### Grafana 仪表盘

导入 Dashboard:
- **Backend Performance**: 自定义
- **Redis Overview**: 763
- **PostgreSQL Overview**: 9628

---

## 🚢 部署方案

### Docker 部署

```bash
# 构建镜像
docker build -t ddd-scaffold-backend:latest .

# 启动服务
docker-compose up -d
```

### Kubernetes 部署

```bash
cd deployments/k8s
kubectl apply -f .
```

详见 [部署指南](docs/deployment_guide.md)

---

## 📚 学习资源

### 文档

- [API 文档](docs/api_reference.md)
- [部署指南](docs/deployment_guide.md)
- [监控与限流](docs/monitoring_and_ratelimit.md)
- [DDD 架构设计](docs/ddd_architecture.md)

### 外部资源

- [DDD 实战指南](https://example.com/ddd-guide)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go 最佳实践](https://github.com/golang-standards/project-layout)

---

## 🤝 贡献指南

### 提交流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 保持 80% 以上的测试覆盖率
- 编写清晰的注释和文档
- 遵循 RESTful API 设计规范

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 👥 维护者

- **Your Name** - [@yourhandle](https://github.com/yourhandle)

---

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [GORM](https://github.com/go-gorm/gorm) - ORM 框架
- [Casbin](https://github.com/casbin/casbin) - 权限框架
- [Redis](https://github.com/redis/go-redis) - Redis 客户端
- [Prometheus](https://github.com/prometheus/client_golang) - 监控指标
- [Swag](https://github.com/swaggo/swag) - API 文档生成

---

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- GitHub Issues: [提交 Issue](https://github.com/your-org/ddd-scaffold/issues)
- Email: support@example.com
- Discord: [加入社区](https://discord.gg/example)

---

**🎉 Happy Coding!**
