# go-ddd-scaffold 项目实施完整计划

## 项目背景
当前项目已完成9个规范性文档体系建设和Git初始化，后端目录结构已建立但缺少核心实现，需要按阶段完善基础设施和核心功能。

## 实施目标
在4周内完成最小可行的产品，包含后端核心领域模型和前端基础框架，实现用户注册登录、租户管理等核心功能。

## 阶段划分

### 阶段一：基础设施完善（第1周，P0优先级）

#### 1.1 配置文件系统
**产出文件**：
- `backend/configs/config.yaml` - 主配置文件
- `backend/configs/.env.example` - 环境变量示例
- `backend/configs/config_dev.yaml` - 开发环境配置
- `backend/configs/config_prod.yaml` - 生产环境配置

**配置内容**：
```yaml
server:
  port: ${SERVER_PORT:8080}
  mode: ${ENV_MODE:debug}

database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: ${DB_NAME:scaffold}
  user: ${DB_USER:postgres}
  password: ${DB_PASSWORD:}

redis:
  addr: ${REDIS_ADDR:localhost:6379}
  password: ${REDIS_PASSWORD:}

jwt:
  secret: ${JWT_SECRET:changeme}
  access_expire: 30m
  refresh_expire: 7d
```

#### 1.2 数据库迁移脚本
**产出文件**：
```
backend/migrations/
├── 000001_create_users_table.up.sql
├── 000002_create_tenants_table.up.sql
├── 000003_create_user_tenants_table.up.sql
├── 000004_create_roles_table.up.sql
├── 000005_create_permissions_table.up.sql
├── 000006_create_role_permissions_table.up.sql
└── 000007_create_audit_logs_table.up.sql
```

### 阶段二：后端核心领域模型（第2周，P0优先级）

#### 2.1 用户领域完善
**检查现有文件**：
- `backend/internal/domain/user/entity.go` - User实体
- `backend/internal/domain/user/service.go` - 用户服务
- `backend/internal/domain/user/repository.go` - 仓储接口

**新增文件**：
- `backend/internal/domain/user/handler.go` - HTTP处理函数
- `backend/internal/domain/user/dto.go` - 请求/响应DTO

#### 2.2 租户领域完善
**检查现有文件**：
- `backend/internal/domain/tenant/` 目录下的实体和服务

**新增文件**：
- `backend/internal/domain/tenant/handler.go`
- `backend/internal/domain/tenant/dto.go`

#### 2.3 认证授权模块（新增）
**新建目录和文件**：
```
backend/internal/auth/
├── service.go      # 认证服务
├── middleware.go   # JWT中间件
├── token.go       # Token生成/验证
└── dto.go         # 认证DTO

backend/internal/permission/
├── service.go      # 权限服务
├── middleware.go   # 权限中间件
└── dto.go         # 权限DTO
```

#### 2.4 基础设施层完善
**新增目录和文件**：
```
backend/internal/infrastructure/
├── persistence/
│   ├── postgres.go    # 数据库连接
│   └── user_repo.go   # 用户仓储实现
├── cache/
│   └── redis.go       # Redis连接
└── config/
    └── loader.go      # 配置加载器
```

#### 2.5 核心API接口实现
**需要实现的接口**：
- POST /api/v1/users - 用户注册
- POST /api/v1/auth/login - 用户登录
- POST /api/v1/auth/refresh - 刷新Token
- GET /api/v1/users/:id - 获取用户信息
- POST /api/v1/tenants - 创建租户
- GET /api/v1/tenants - 租户列表

### 阶段三：前端基础框架（第3周，P1优先级）

#### 3.1 项目初始化
**创建目录结构**：
```
frontend/
├── src/
│   ├── business/
│   │   ├── auth/          # 认证模块
│   │   ├── user/          # 用户模块
│   │   └── tenant/        # 租户模块
│   ├── data/
│   │   ├── api/           # API客户端
│   │   ├── hooks/         # 自定义hooks
│   │   └── store/         # 状态管理
│   ├── interaction/
│   │   ├── components/    # 通用组件
│   │   ├── forms/         # 表单组件
│   │   └── layouts/       # 布局组件
│   ├── presentation/
│   │   ├── pages/         # 页面组件
│   │   └── styles/        # 样式文件
│   ├── i18n/              # 国际化
│   ├── utils/             # 工具函数
│   ├── App.tsx
│   └── main.tsx
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── vite.config.ts
```

#### 3.2 核心页面开发
**优先级页面**：
1. 登录页面 `/login` （P0）
2. 注册页面 `/register` （P0）
3. 仪表板 `/dashboard` （P1）
4. 用户管理 `/users` （P1）
5. 租户管理 `/tenants` （P1）

### 阶段四：部署配置（第4周，P1-P2优先级）

#### 4.1 Docker配置
**产出文件**：
```
backend/deployments/docker/
├── docker-compose.yml      # 开发环境
├── docker-compose.prod.yml # 生产环境
└── Dockerfile             # 应用镜像
```

#### 4.2 Kubernetes配置（预留）
**产出文件**：
```
backend/deployments/kubernetes/
├── deployment.yaml
├── service.yaml
├── ingress.yaml
├── configmap.yaml
└── secrets.yaml
```

## 时间安排

| 周次 | 阶段 | 主要任务 | 交付物 |
|------|------|----------|--------|
| 第1周 | 阶段一 | 配置文件 + 数据库迁移 | config.yaml + migrations |
| 第2周 | 阶段二 | 后端核心领域模型 | User/Tenant/Auth API |
| 第3周 | 阶段三 | 前端基础框架 | React项目 + 登录页 |
| 第4周 | 阶段四 | 部署配置 + 功能完善 | Docker + 核心页面 |

## 质量保证
- 每个阶段完成后进行单元测试
- 核心功能需通过端到端测试验证
- 代码遵循既定的开发规范
- 文档与代码同步更新

## 风险控制
- 严格按照优先级实施，避免范围蔓延
- 每日检查进度，及时调整计划
- 保留完整的工作记录和版本历史
- 确保可随时回滚到稳定状态