# DDD 脚手架优化规划 - 回归本质

## 🎯 核心定位

**打造企业级 Go 语言 DDD 脚手架，让标准架构快速落地**

### 三个核心原则

1. **规范化** - 统一的代码结构和开发规范
2. **流程化** - 清晰的开发流程和最佳实践  
3. **工具化** - 自动化的代码生成和辅助工具

### 三个避免

1. ❌ **避免过度设计** - 不引入不必要的复杂性
2. ❌ **避免过早优化** - 不为假设的场景做设计
3. ❌ **避免重复造轮子** - 使用成熟的开源方案

---

## 📊 当前问题分析

### 已做得好的地方 ✅

1. **DDD 四层架构清晰**
   - Interfaces / Application / Domain / Infrastructure
   - 依赖方向正确，符合 Clean Architecture

2. **领域建模规范**
   - Entity（实体）、ValueObject（值对象）
   - Aggregate Root（聚合根）、Repository（仓储）
   - Domain Service（领域服务）、Domain Event（领域事件）

3. **基础设施完善**
   - JWT + Casbin 认证授权
   - GORM + PostgreSQL 数据持久化
   - Redis 缓存支持
   - Prometheus 监控指标

### 需要改进的地方 ⚠️

1. **规范性不足**
   - [ ] 缺少统一的命名规范文档
   - [ ] 代码注释风格不一致
   - [ ] 错误处理不规范（有的用 AppError，有的直接 fmt.Errorf）
   - [ ] 测试覆盖率要求不明确

2. **流程化缺失**
   - [ ] 新增领域模块的流程不清晰
   - [ ] 数据库迁移流程不规范
   - [ ] 代码审查 checklist 缺失
   - [ ] 版本发布流程未定义

3. **工具化不够**
   - [ ] 缺少代码生成工具（DAO、DTO、Assembler）
   - [ ] 缺少项目脚手架工具（一键创建新模块）
   - [ ] 缺少自动化测试工具
   - [ ] Makefile 命令不完善

4. **文档不完善**
   - [ ] docs 目录大部分是空的
   - [ ] 缺少"如何创建一个新领域"的教程
   - [ ] 缺少架构决策记录（ADR）
   - [ ] API 文档不完整

---

## 🚀 优化方向（按优先级）

### 优先级 1: 规范化建设 ⭐⭐⭐⭐⭐

#### 1.1 代码规范文档

**目标**: 统一代码风格，降低协作成本

**内容**:
```
docs/standards/
├── code-style.md           # 代码风格指南
├── naming-conventions.md   # 命名规范
├── error-handling.md       # 错误处理规范
├── comment-guidelines.md   # 注释编写规范
└── testing-standards.md    # 测试编写规范
```

**具体规范**:
- 包命名：使用小写，复数形式（`users`, `repositories`）
- 接口命名：`-er` 后缀（`UserRepository`, `EventPublisher`）
- 错误处理：统一使用 `pkg/errors`，禁止裸 error
- 注释规范：所有导出元素必须有注释，使用完整句子
- 测试规范：单元测试覆盖率 ≥ 80%，表驱动测试

**工作量**: 1-2 天

---

#### 1.2 DDD 实现规范

**目标**: 明确 DDD 各层的职责和边界

**内容**:
```
docs/standards/ddd-implementation.md
```

**核心规则**:
```go
// Domain Layer - 纯业务逻辑，无基础设施依赖
type User struct {              // 实体
    ID        uuid.UUID
    Email     Email         // 值对象
    Password  HashedPassword // 值对象
}

func (u *User) UpdateEmail(email Email) error {  // 业务方法
    if u.Email.Equals(email) {
        return nil
    }
    u.Email = email
    u.RecordEvent(UserEmailChangedEvent{...})
    return nil
}

// Application Layer - 应用编排，不含业务逻辑
type UserService struct {      // 应用服务
    userRepo UserRepository
    eventBus EventBus
}

func (s *UserService) ChangeEmail(ctx context.Context, userID uuid.UUID, email string) error {
    // 1. 获取聚合根
    user, _ := s.userRepo.GetByID(ctx, userID)
    
    // 2. 调用领域方法
    user.UpdateEmail(valueobject.NewEmail(email))
    
    // 3. 持久化
    s.userRepo.Update(ctx, user)
    
    // 4. 发布事件
    s.eventBus.Publish(user.Events())
    
    return nil
}

// Infrastructure Layer - 基础设施实现
type userRepository struct {   // 仓储实现
    db *gorm.DB
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
    var model UserModel
    r.db.First(&model, id)
    return model.ToEntity(), nil  // Model ↔ Entity 转换
}
```

**禁止事项**:
- ❌ Domain 层不能 import Infrastructure 包
- ❌ Application 层不能有业务判断逻辑
- ❌ Infrastructure 层不能直接返回 Entity（必须通过 Repository 接口）

**工作量**: 1 天

---

#### 1.3 项目结构规范

**目标**: 统一项目目录结构

**标准结构**:
```
backend/
├── cmd/                    # 应用入口
│   ├── server/            # HTTP 服务
│   ├── worker/            # 后台任务
│   └── migrate/           # 数据库迁移
│
├── internal/              # 内部代码（不对外暴露）
│   ├── domain/           # 领域层
│   │   ├── user/         # 用户领域
│   │   │   ├── entity/          # 实体
│   │   │   ├── valueobject/     # 值对象
│   │   │   ├── repository/      # 仓储接口
│   │   │   ├── service/         # 领域服务
│   │   │   └── event/           # 领域事件
│   │   └── tenant/       # 租户领域
│   │
│   ├── application/      # 应用层
│   │   ├── user/
│   │   │   ├── service/       # 应用服务
│   │   │   ├── dto/           # 数据传输对象
│   │   │   └── assembler/     # DTO↔Entity 转换器
│   │   └── tenant/
│   │
│   ├── infrastructure/   # 基础设施层
│   │   ├── persistence/  # 持久化实现
│   │   │   ├── gorm/
│   │   │   │   ├── model/     # 数据模型
│   │   │   │   ├── dao/       # DAO 层（代码生成）
│   │   │   │   └── repo/      # 仓储实现
│   │   │   └── redis/
│   │   ├── auth/             # 认证服务
│   │   ├── event/            # 事件总线
│   │   ├── cache/            # 缓存服务
│   │   └── middleware/       # 中间件
│   │
│   ├── interfaces/       # 接口层（适配器）
│   │   ├── http/
│   │   │   ├── handler/     # HTTP Handler
│   │   │   ├── middleware/  # HTTP 中间件
│   │   │   └── router/      # 路由注册
│   │   └── grpc/
│   │
│   └── pkg/              # 通用工具包
│       ├── errors/       # 错误处理
│       ├── validator/    # 参数验证
│       └── response/     # 统一响应
│
├── migrations/            # 数据库迁移
│   ├── sql/              # SQL 迁移脚本
│   └── goose/            # Goose 配置
│
├── tools/                # 开发工具
│   ├── dao-gen/         # DAO 代码生成
│   └── scaffold/        # 项目脚手架
│
├── tests/               # 测试代码
│   ├── unit/            # 单元测试
│   ├── integration/     # 集成测试
│   └── e2e/             # 端到端测试
│
├── config/              # 配置文件
├── scripts/             # 脚本文件
├── Makefile             # 构建脚本
└── go.mod               # Go 模块定义
```

**关键规则**:
- `internal/` 下的代码只能被本项目引用
- `pkg/` 下的代码可以被外部项目引用
- 按领域划分目录，而非按技术类型（如 `user/` 而非 `services/`）

**工作量**: 0.5 天（整理现有结构）

---

### 优先级 2: 流程化建设 ⭐⭐⭐⭐⭐

#### 2.1 新增领域模块流程

**目标**: 10 分钟内创建一个完整的领域模块

**标准流程**:
```bash
# 1. 运行脚手架工具
make create-domain name=user

# 自动生成以下目录结构:
# backend/internal/domain/user/
#   ├── entity/
#   │   └── user.go
#   ├── valueobject/
#   │   └── email.go
#   ├── repository/
#   │   └── repository.go
#   └── event/
#       └── user_events.go

# 2. 运行应用服务脚手架
make create-application name=user

# 生成:
# backend/internal/application/user/
#   ├── service/
#   │   └── service.go
#   ├── dto/
#   │   └── user_dto.go
#   └── assembler/
#       └── user_assembler.go

# 3. 运行 HTTP Handler 脚手架
make create-handler name=user

# 生成:
# backend/internal/interfaces/http/user/
#   └── handler.go

# 4. 运行数据库迁移生成
make migration-create name=create_users_table

# 生成:
# migrations/sql/TIMESTAMP_create_users_table.sql
```

**输出物**:
- ✅ 领域模型模板（含注释）
- ✅ 仓储接口和实现
- ✅ 应用服务模板
- ✅ DTO 和 Assembler
- ✅ HTTP Handler 模板
- ✅ 数据库迁移模板
- ✅ 单元测试模板

**工作量**: 2-3 天（开发脚手架工具）

---

#### 2.2 数据库迁移流程

**目标**: 规范数据库变更管理

**流程**:
```bash
# 1. 创建迁移
make migration-create name=add_user_phone_field

# 生成: migrations/sql/TIMESTAMP_add_user_phone_field.sql

# 2. 编辑迁移脚本
# -- +goose Up
# ALTER TABLE users ADD COLUMN phone VARCHAR(20);
# 
# -- +goose Down
# ALTER TABLE users DROP COLUMN phone;

# 3. 执行迁移
make migrate-up

# 4. 生成 DAO 代码
make dao-generate

# 5. 提交代码
git add migrations/ internal/infrastructure/persistence/
```

**规范**:
- 每个迁移文件必须是幂等的
- 必须包含 Up 和 Down 两部分
- 文件名格式：`TIMESTAMP_description.sql`
- 禁止在生产环境直接修改表结构

**工作量**: 0.5 天（完善 Makefile）

---

#### 2.3 代码审查 Checklist

**目标**: 保证代码质量

**Checklist**:
```markdown
## DDD 规范检查
- [ ] 领域层是否包含基础设施依赖？（不应该有）
- [ ] 应用服务是否只负责编排，不含业务逻辑？
- [ ] 仓储实现是否通过接口依赖？
- [ ] 实体是否包含业务方法？（应该有）

## 代码质量检查
- [ ] 所有导出元素是否有注释？
- [ ] 错误处理是否规范？（使用 AppError）
- [ ] 是否有单元测试？（覆盖率≥80%）
- [ ] 是否有竞态检测？（go test -race）

## 安全检査
- [ ] 用户输入是否验证？（使用 validator）
- [ ] SQL 注入防护？（使用参数化查询）
- [ ] 敏感信息是否加密？（密码、Token）
- [ ] 权限校验是否正确？（Casbin）

## 性能检查
- [ ] 是否有 N+1 查询问题？
- [ ] 数据库索引是否合理？
- [ ] 是否有缓存策略？
- [ ] 连接池配置是否优化？
```

**工作量**: 0.5 天

---

#### 2.4 测试流程

**目标**: 保证代码质量

**测试策略**:
```
测试金字塔:
        /\
       /  \
      / E2E \      端到端测试（10%）
     /______\
    /        \
   / Integration\  集成测试（20%）
  /______________\
 /                \
/    Unit Tests    \  单元测试（70%）
--------------------
```

**测试要求**:
```bash
# 单元测试（运行快，无外部依赖）
make test-unit
# 要求：覆盖率 ≥ 80%, 执行时间 < 30s

# 集成测试（需要数据库、Redis）
make test-integration  
# 要求：覆盖核心业务流程

# 端到端测试（完整系统）
make test-e2e
# 要求：覆盖关键用户路径
```

**测试模板**:
```go
package user_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "go-ddd-scaffold/internal/domain/user/entity"
    "go-ddd-scaffold/internal/domain/user/valueobject"
)

// 测试实体行为
func TestUser_UpdateEmail(t *testing.T) {
    // Arrange
    user := createUser()
    newEmail, _ := valueobject.NewEmail("new@example.com")
    
    // Act
    err := user.UpdateEmail(newEmail)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, newEmail, user.Email)
}

// 测试值对象验证
func TestEmail_InvalidFormat(t *testing.T) {
    // Arrange
    invalidEmails := []string{
        "",
        "invalid",
        "@example.com",
    }
    
    // Act & Assert
    for _, email := range invalidEmails {
        _, err := valueobject.NewEmail(email)
        assert.Error(t, err, "email: %s", email)
    }
}
```

**工作量**: 1 天（完善测试基础设施）

---

### 优先级 3: 工具化建设 ⭐⭐⭐⭐⭐

#### 3.1 DAO 代码生成器

**目标**: 自动生成 CRUD 代码

**现状**: 已有 `tools/dao-gen`，但需要完善

**改进方向**:
```bash
# 当前
cd tools/dao-gen && go run main.go

# 改进后
make dao-generate  # 一键生成

# 生成的代码包含:
# - Model 结构体（带 GORM 标签）
# - DAO 接口和实现
# - 基础 CRUD 方法
# - 自定义查询方法模板
```

**功能增强**:
- ✅ 支持多表关联查询生成
- ✅ 支持软删除字段自动处理
- ✅ 支持审计字段（CreatedAt, UpdatedAt）
- ✅ 支持分页查询生成

**工作量**: 1-2 天

---

#### 3.2 领域模块脚手架

**目标**: 一键创建完整的领域模块

**设计**:
```bash
# 命令行工具
./scaffold create domain user --fields="email:string,nickname:string,password:string"

# 或使用 Makefile
make create-domain name=user fields="email,nickname,password"
```

**生成内容**:
```
backend/internal/domain/user/
├── entity/
│   └── user.go              # 实体（含业务方法模板）
├── valueobject/
│   ├── email.go             # 值对象（含验证逻辑）
│   └── nickname.go
├── repository/
│   └── repository.go        # 仓储接口
└── event/
    └── user_events.go       # 领域事件模板

backend/internal/application/user/
├── service/
│   └── user_service.go      # 应用服务（含 CQRS 分离）
├── dto/
│   ├── user_dto.go          # DTO 定义
│   └── requests.go          # 请求 DTO
└── assembler/
    └── user_assembler.go    # 转换器

backend/internal/interfaces/http/user/
└── handler.go               # HTTP Handler（含 Swagger 注释）

backend/tests/unit/domain/user/
└── user_test.go             # 单元测试模板

migrations/sql/
└── TIMESTAMP_create_users.sql  # 迁移脚本模板
```

**工作量**: 3-4 天

---

#### 3.3 Makefile 完善

**目标**: 常用操作命令化

**当前命令**:
```bash
make help              # 显示帮助
make build            # 编译
make run              # 运行
make test             # 测试
make migrate-up       # 迁移数据库
```

**需要增加**:
```bash
# 代码生成
make generate              # 运行所有代码生成
make generate-dao          # 生成 DAO 代码
make generate-dto          # 生成 DTO
make swagger              # 生成 Swagger 文档

# 脚手架
make create-domain         # 创建领域模块
make create-application    # 创建应用服务
make create-handler        # 创建 HTTP Handler

# 代码质量
make lint                 # 代码检查
make fmt                  # 格式化代码
make tidy                 # 清理依赖
make cover                # 查看测试覆盖率

# 开发辅助
make dev                  # 启动开发服务器（热重载）
make docker-build         # 构建 Docker 镜像
make docker-run           # Docker 运行
make clean                # 清理构建产物
```

**工作量**: 1 天

---

#### 3.4 开发环境容器化

**目标**: 一键启动开发环境

**Docker Compose**:
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    depends_on:
      - postgres
      - redis
  
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: ddd_scaffold
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
  
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

**使用**:
```bash
make dev-environment    # 启动所有服务
make logs              # 查看日志
make stop              # 停止所有服务
```

**工作量**: 1 天

---

### 优先级 4: 文档完善 ⭐⭐⭐⭐

#### 4.1 文档结构重建

**目标**: 建立清晰的文档体系

**新结构**:
```
docs/
├── README.md                      # 文档索引
│
├── getting-started/              # 快速开始
│   ├── installation.md            # 安装指南
│   ├── quickstart.md              # 5 分钟快速体验
│   └── configuration.md           # 配置说明
│
├── standards/                    # 规范文档
│   ├── code-style.md              # 代码风格
│   ├── naming-conventions.md      # 命名规范
│   ├── ddd-implementation.md      # DDD 实现规范
│   ├── error-handling.md          # 错误处理
│   └── testing-standards.md       # 测试规范
│
├── guides/                       # 开发指南
│   ├── create-domain-module.md    # 创建领域模块
│   ├── database-migration.md      # 数据库迁移
│   ├── add-api-endpoint.md        # 添加 API 端点
│   └── implement-business-logic.md # 实现业务逻辑
│
├── architecture/                 # 架构文档
│   ├── overview.md                # 架构总览
│   ├── layers.md                  # 分层架构
│   ├── dependencies.md            # 依赖关系
│   └── adr/                       # 架构决策记录
│       └── 001-use-gorm.md
│
├── api-reference/                # API 文档
│   ├── authentication.md          # 认证 API
│   ├── user-management.md         # 用户管理 API
│   └── tenant-management.md       # 租户管理 API
│
└── deployment/                   # 部署文档
    ├── local-development.md       # 本地开发
    ├── docker-deployment.md       # Docker 部署
    └── kubernetes-deployment.md   # K8s 部署
```

**工作量**: 2-3 天

---

#### 4.2 示例领域完善

**目标**: 提供完整的参考实现

**选择**: `Book` 领域（简单、易懂、无复杂业务）

**完善内容**:
```
backend/internal/domain/book/
├── entity/
│   ├── book.go                  # 聚合根
│   └── author.go                # 实体
├── valueobject/
│   ├── isbn.go                  # ISBN 值对象
│   ├── title.go                 # 书名值对象
│   └── price.go                 # 价格值对象
├── repository/
│   └── repository.go            # 仓储接口
├── service/
│   └── book_service.go          # 领域服务（业务逻辑）
└── event/
    └── book_events.go           # 领域事件

backend/internal/application/book/
├── service/
│   ├── book_command_service.go  # 命令服务
│   └── book_query_service.go    # 查询服务
├── dto/
│   ├── book_dto.go
│   └── requests.go
└── assembler/
    └── book_assembler.go

backend/internal/interfaces/http/book/
└── handler.go                   # HTTP Handler

backend/tests/
├── unit/domain/book/           # 领域层测试
├── integration/application/book/ # 应用层测试
└── e2e/book_api_test.go        # API 测试
```

**特点**:
- ✅ 完整的 DDD 实现
- ✅ 包含所有最佳实践
- ✅ 高测试覆盖率（≥90%）
- ✅ 详细的代码注释

**工作量**: 2 天

---

## 📋 实施计划

### 第一阶段（1-2 周）: 基础规范
- [ ] 代码规范文档（1.1）
- [ ] DDD 实现规范（1.2）
- [ ] 项目结构整理（1.3）
- [ ] Makefile 完善（3.3）

### 第二阶段（2-3 周）: 工具开发
- [ ] DAO 生成器完善（3.1）
- [ ] 领域模块脚手架（3.2）
- [ ] 开发环境容器化（3.4）

### 第三阶段（3-4 周）: 流程建设
- [ ] 新增领域模块流程（2.1）
- [ ] 数据库迁移流程（2.2）
- [ ] 代码审查 Checklist（2.3）
- [ ] 测试流程（2.4）

### 第四阶段（4-5 周）: 文档完善
- [ ] 文档结构重建（4.1）
- [ ] 示例领域完善（4.2）
- [ ] API 文档完善

**总工作量**: 约 15-20 个工作日

---

## 🎯 成功标准

### 规范化 ✅
- [ ] 有完整的规范文档
- [ ] 团队成员能说出规范内容
- [ ] 代码审查有据可依

### 流程化 ✅
- [ ] 新人能在 1 小时内创建第一个领域模块
- [ ] 所有数据库变更都通过迁移脚本
- [ ] 每次 PR 都有 Checklist

### 工具化 ✅
- [ ] 80% 的重复代码通过工具生成
- [ ] Makefile 覆盖所有常用操作
- [ ] 开发环境一键启动

---

## 💡 核心理念

**少即是多**:
- 不引入不必要的技术栈
- 不为假设的场景做设计
- 保持代码简洁易读

**约定优于配置**:
- 明确的目录结构
- 统一的命名规范
- 标准的代码模板

**自动化优先**:
- 能自动化的绝不手动
- 工具能解决的绝不靠人脑
- 减少人为错误

**渐进式演进**:
- 不追求一步到位
- 根据实际需求调整
- 保持向后兼容

---

**制定日期**: 2026-03-06  
**审核状态**: Draft  
**下次回顾**: 2026-03-20
