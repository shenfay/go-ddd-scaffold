# API基础设施完善行动计划

## 📋 当前状态评估

### ✅ 已完成 (Phase 1)
- [x] 统一错误码体系 (语义化英文 + 层级前缀)
- [x] 统一响应格式 (code/message/data/error/requestId/timestamp)
- [x] 结构化日志中间件 (Zap + Request-ID追踪)
- [x] 参数校验中间件 (go-playground/validator + 自定义规则)
- [x] 请求ID生成与传递机制
- [x] CORS跨域配置
- [x] DDD分层错误处理规范 (Repository/Service/Handler各层统一使用预定义错误码)

### ✅ 已完成 (Phase 2)
- [x] 中间件响应格式统一 (middleware/response 包)
- [x] 循环依赖问题修复 (errors/context.go)
- [x] Request DTO 规范化 (统一使用 query.QueryParams)

### ✅ 已完成 (Phase 3 - 当前实现)
- [x] 自定义指标收集器 (内存版本)
- [x] QPS/延迟/错误率指标采集
- [x] 按端点统计 (method + path)
- [x] 监控数据接口 (/api/metrics)

### ✅ 已完成 (Phase 4)
- [x] HTTP缓存中间件 (middleware/cache.go)
- [x] ETag生成与验证
- [x] Last-Modified支持
- [x] Cache-Control配置
- [x] 知识图谱只读接口缓存启用

### ✅ 已完成 (Phase 5)
- [x] 业务校验器 (pkg/validator/)
- [x] 邮箱唯一性校验
- [x] 密码强度校验
- [x] 昵称格式校验
- [x] 租户用户数限制校验
- [x] 按模块拆分校验文件 (validation.go/user_validator.go)
- [x] 在 User Service Register 方法中集成密码强度校验

### 📖 业务校验器使用说明

#### 位置
`backend/internal/pkg/validator/`

#### 文件结构
```
validator/
├── validator.go        # 参数校验封装（go-playground/validator）
├── validation.go       # 通用校验类型和函数
├── user_validator.go   # User 模块业务校验
└── knowledge_validator.go  # Knowledge 模块业务校验（待添加）
```

#### 已集成的校验
- User Service Register 方法：密码强度校验

#### 如何扩展新的校验规则

1. **在对应模块的 validator 文件中添加方法**：
```go
// user_validator.go
func (v *UserValidator) ValidateXxx(ctx context.Context, req *XxxRequest) BusinessValidationErrors {
    var errs BusinessValidationErrors
    // 添加校验逻辑
    return errs
}
```

2. **在 Service 层调用**：
```go
func (s *Service) Xxx(ctx context.Context, req *dto.XxxRequest) (*dto.Xxx, error) {
    // 调用业务校验
    if errs := s.userValidator.ValidateXxx(ctx, req); errs.HasError() {
        return nil, errPkg.InvalidParams.WithDetails(errs.Error())
    }
    // 继续业务逻辑
}
```

3. **添加新的校验函数**（无状态）：
```go
// validation.go
func ValidateXxxFormat(value string) error {
    // 校验逻辑
    return nil
}
```

### ⏳ 待处理问题 (按优先级排序)

#### 🔴 高优先级
1. ~~幂等性处理~~ - 防止重复操作 ✅ 已完成
2. **CI/CD流水线** - 自动化部署和测试

#### 🟡 中优先级
3. **负载与性能优化** - 分页/压缩/批量处理
4. ~~资源标识规范化~~ - UUID v4 统一使用 ✅ 已完成

#### 🟢 低优先级
5. ~~国际化支持~~ - 多语言错误信息 ✅ 已完成
6. **数据一致性保障** - 分布式事务处理
7. **可扩展性设计** - 预留字段/灵活结构

---

## 🚀 Phase 1.5: 幂等性处理 ✅ 已完成

### 目标
实现接口幂等性，防止重复提交造成数据不一致

### 技术方案
- **幂等Key**: 支持自动生成（基于请求内容MD5）或客户端提供
- **存储**: 内存存储（开发环境）+ Redis支持（生产环境）
- **适用接口**: POST, PUT, DELETE, PATCH

### 实施步骤
1. [完成] 创建幂等性中间件 `middleware/idempotency.go`
2. [完成] 定义幂等Key生成规则（自动+手动）
3. [完成] 内存存储实现
4. [待定] Redis存储支持
5. [完成] 在关键接口中应用中间件

### 关键文件
```
backend/internal/infrastructure/middleware/idempotency.go
```

---

## 🚀 Phase 1.5: CI/CD基础流水线 (预计 3-5天)

### 目标
建立完整的持续集成/持续部署流水线，实现前后端一体化项目的自动化测试、构建和部署

### 技术方案
- **CI工具**: GitHub Actions
- **项目结构**: 前后端一体化单体项目
- **测试覆盖**: 前端单元测试 + 后端单元测试 + 集成测试
- **代码质量**: 静态检查 + 代码覆盖率
- **部署方式**: Docker容器化部署

### 实施步骤
1. 创建GitHub Actions工作流配置（支持前后端）
2. 配置前端构建和测试（React + Vite/webpack）
3. 配置后端单元测试和集成测试
4. 设置代码质量检查（ESLint + golangci-lint）
5. 实现一体化Docker镜像构建和推送
6. 配置自动化部署到测试环境

### 关键文件
```
.github/workflows/ci.yml
.github/workflows/cd.yml
Dockerfile              # 前后端一体化镜像
docker-compose.test.yml
frontend/vite.config.js # 或 webpack配置
backend/Makefile
```

### 单体项目特殊考虑

**构建策略**：
```
一体化构建流程：
1. 前端构建 → 生成静态文件
2. 后端构建 → 编译Go二进制
3. 打包 → 前端静态文件嵌入后端
4. 镜像构建 → 单一可执行镜像
```

**测试策略**：
```
并行测试执行：
├── 前端测试：Jest + React Testing Library
├── 后端测试：Go test + testify
└── 集成测试：API接口测试 + E2E测试
```

**部署策略**：
```
单一部署单元：
- 一个Docker镜像包含前后端
- 通过环境变量控制运行模式
- 统一的健康检查和监控
```

---

## 🚀 Phase 2: 用户管理基础模块 (预计 4-5天)

### 目标
实现用户管理核心功能，验证技术栈完整性和DDD架构

### 功能范围
- 用户注册/登录
- 家庭租户管理
- 权限角色分配
- 基础信息维护

### 技术验证点
- DDD分层架构实践
- JWT认证机制
- 数据库事务处理
- API文档生成

### 关键文件
```
backend/internal/domain/user/
backend/internal/application/user/service/
backend/internal/application/user/dto/
backend/internal/infrastructure/persistence/gorm/repo/user_repository.go
backend/internal/interfaces/http/user_handler.go
```

---

### 目标
实现接口幂等性，防止重复提交造成数据不一致

### 技术方案
- **幂等Key生成**: 基于业务主键 + 时间戳 + 操作类型
- **存储介质**: Redis (过期时间 24小时)
- **适用接口**: PUT(更新), DELETE(删除), POST(创建订单等)

### 实施步骤
1. 创建幂等性中间件 `middleware/idempotency.go`
2. 定义幂等Key生成规则
3. Redis存储幂等记录
4. 在关键接口中应用中间件
5. 编写单元测试

### 关键文件
```
backend/internal/infrastructure/middleware/idempotency.go
backend/internal/pkg/idempotency/key_generator.go
backend/test/integration/idempotency_test.go
```

---

## 🚀 Phase 3: 监控与日志增强 (预计 2天)

### 目标
建立完整的API监控体系，支持Prometheus指标采集

### 技术方案
- **指标类型**: QPS, 平均延迟, P95/P99延迟, 错误率
- **标签维度**: 接口路径, HTTP方法, 状态码, 客户端类型
- **日志增强**: 结构化字段标准化

### 实施步骤
1. 集成 Prometheus client_golang
2. 创建指标收集器 `pkg/metrics/api_metrics.go`
3. 在Logger中间件中增加指标上报
4. 配置Grafana仪表板
5. 设置告警规则

### 关键文件
```
backend/internal/pkg/metrics/api_metrics.go
backend/internal/infrastructure/middleware/metrics.go
docker-compose.monitoring.yml
docs/monitoring/grafana-dashboard.json
```

---

## 🚀 Phase 4: 缓存策略 ✅ 已完成

### 目标
通过HTTP缓存头部和服务端缓存提升性能

### 技术方案
- **HTTP缓存**: ETag + Last-Modified + Cache-Control
- **服务端缓存**: Redis缓存热点数据
- **缓存策略**: LRU + TTL

### 实施步骤
1. [完成] 创建缓存中间件 `middleware/cache.go`
2. [完成] 实现ETag生成器 (基于内容哈希)
3. [完成] 添加Last-Modified支持
4. [待定] 配置Redis缓存客户端
5. [完成] 在只读接口中启用缓存

### 关键文件
```
backend/internal/infrastructure/middleware/cache.go
backend/internal/pkg/cache/redis_client.go
backend/internal/pkg/cache/etag_generator.go
```

---

## 🚀 Phase 5: 数据校验强化 ✅ 已完成

### 目标
加强业务层面的数据校验，防止非法数据入库

### 技术方案
- **校验层级**: 参数校验 → 业务校验 → 数据库约束
- **校验内容**: 唯一性检查, 关联性验证, 业务规则校验
- **错误反馈**: 精准定位错误字段和原因

### 实施步骤
1. [完成] 创建业务校验器 `pkg/validator/`
2. [完成] 实现常用校验规则 (唯一性, 关联检查)
3. [完成] 在Service层集成业务校验
4. [完成] 统一校验错误返回格式

### 关键文件
```
backend/internal/pkg/validator/
├── validator.go        # 参数校验封装
├── validation.go       # 通用校验类型和函数
├── user_validator.go  # User模块业务校验
```

---

## 🚀 Phase 6: 负载与性能优化 (预计 2天)

### 目标
优化大数据量接口性能，支持分页和压缩

### 技术方案
- **分页策略**: 基于游标的分页 (Cursor-based Pagination)
- **数据压缩**: gzip压缩响应体 (>1KB数据)
- **批量操作**: 支持批量创建/更新

### 实施步骤
1. 实现游标分页工具 `pkg/pagination/cursor.go`
2. 添加gzip压缩中间件
3. 设计批量操作接口规范
4. 性能压测对比优化效果

### 关键文件
```
backend/internal/pkg/pagination/cursor.go
backend/internal/infrastructure/middleware/compression.go
```

---

## 🚀 Phase 7: 国际化支持 ✅ 已完成

### 目标
支持多语言错误信息和业务文案

### 技术方案
- **语言文件**: YAML配置文件（可独立编辑）
- **语言检测**: URL参数 > Header > Cookie > 默认
- **支持语言**: zh-CN, en-US, ja-JP
- **fallback**: YAML加载失败时使用内置map

### 实施步骤
1. [完成] 创建i18n核心包 `pkg/i18n/`
2. [完成] 实现语言检测和上下文传递
3. [完成] 创建多语言资源文件 (YAML格式)
4. [完成] DTO层集成国际化翻译
5. [完成] 响应字段name/description支持多语言

### 关键文件
```
backend/internal/pkg/i18n/
├── i18n.go           # 核心函数
├── loader.go         # YAML加载器
└── locale/
    ├── zh-CN.yaml    # 中文翻译
    ├── en-US.yaml    # 英文翻译
    └── ja-JP.yaml    # 日文翻译
```

### 使用示例
```go
// DTO转换时自动翻译
func ToDomainResponse(ctx context.Context, domain *entity.Domain) *DomainResponse {
    return &DomainResponse{
        Name:        i18n.GetMessageByContext(ctx, domain.NameKey),
        Description: i18n.GetMessageByContext(ctx, domain.DescriptionKey),
    }
}
```

---

## 🚀 Phase 8: 数据一致性保障 (预计 3-4天)

### 目标
在分布式环境下保证数据最终一致性

### 技术方案
- **事务模式**: Saga模式 + 补偿事务
- **消息队列**: 使用RabbitMQ/RocketMQ
- **幂等消费**: 消费端幂等性保证

### 实施步骤
1. 搭建消息队列环境
2. 实现Saga事务协调器
3. 创建补偿事务处理器
4. 设计死信队列和重试机制
5. 编写集成测试

### 关键文件
```
backend/internal/pkg/messaging/saga_coordinator.go
backend/internal/pkg/messaging/consumer.go
docker-compose.mq.yml
```

---

## 🚀 Phase 9: 可扩展性设计 (持续进行)

### 目标
接口设计具备良好的向前兼容性

### 设计原则
- **预留字段**: response中预留ext字段
- **版本管理**: URI版本(/v1/) + Accept头版本
- **渐进式增强**: 新功能通过可选字段提供

### 实施要点
1. 制定API版本管理规范
2. 在Response结构中添加Ext字段
3. 建立API变更评审流程
4. 维护API兼容性矩阵

---

## 📊 资源需求

### 人力投入
- 后端开发: 1人全程负责
- 前端配合: 0.3人(国际化、监控面板)
- 测试支持: 0.5人(接口测试、性能测试)

### 技术依赖
- Redis (幂等性、缓存)
- Prometheus + Grafana (监控)
- RabbitMQ (消息队列，可选)

### 时间估算
- **总计**: 14-19个工作日
- **建议节奏**: 每个Phase完成后进行集成测试

---

## ✅ 验收标准

每个Phase完成后需满足:
1. 通过单元测试 (覆盖率≥80%)
2. 通过集成测试
3. 性能指标达标 (如QPS提升、延迟降低)
4. 文档更新完整
5. 代码评审通过

---

## 📝 风险管控

### 技术风险
- **缓存穿透**: 增加空值缓存策略
- **Redis单点**: 配置主从+哨兵
- **幂等性失效**: 增加重试机制

### 进度风险
- **依赖阻塞**: 提前准备替代方案
- **需求变更**: 保持架构灵活性
- **人员变动**: 完善文档和注释

---

**负责人**: 后端开发团队  
**启动时间**: 待确认  
**预计完成**: 3-4周
