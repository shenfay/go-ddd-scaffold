# 架构设计总结

## 📁 最终目录结构

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go                    # API 服务入口
│   ├── worker/
│   │   └── main.go                    # Worker 服务入口（待实现）
│   └── cli/
│       └── main.go                    # CLI 工具入口（待实现）
│
├── internal/
│   ├── auth/                          # 认证模块（垂直切片）
│   │   ├── domain.go                  # 领域模型（User 聚合根）
│   │   ├── repository.go              # 仓储接口
│   │   ├── repository_gorm.go         # GORM 仓储实现
│   │   ├── service.go                 # 应用服务
│   │   ├── token_service.go           # Token 服务
│   │   └── handler.go                 # HTTP Handler + 路由
│   │
│   └── infrastructure/                # 基础设施层
│       └── config/
│           └── config.go              # 配置管理
│
├── pkg/                               # 公共包
│   ├── constants/
│   │   └── constants.go               # 常量定义
│   ├── errors/
│   │   └── errors.go                  # 统一错误处理
│   └── utils/
│       └── ulid/
│           └── ulid.go                # ULID 生成工具
│
├── migrations/                        # 数据库迁移
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
│
├── configs/
│   ├── .env.example
│   ├── development.yaml               # 开发环境配置
│   └── production.yaml                # 生产环境配置（待创建）
│
├── test/                              # 测试（待实现）
│   ├── integration/
│   └── fixtures/
│
├── deployments/                       # 部署配置（待实现）
├── docs/                              # 文档（待实现）
├── Makefile
├── go.mod
└── go.sum
```

---

## 🎯 核心设计决策

### 1. **架构风格：简洁实用主义（方案 C）**
- ✅ 按业务模块垂直切片组织代码
- ✅ 每个模块包含完整的分层（Domain/Service/Handler）
- ✅ 平台层使用 `infrastructure` 命名

### 2. **ID 生成策略：ULID (oklog/ulid)**
- ✅ 有序的 UUID，利于数据库索引性能
- ✅ 格式：`user_{ulid}`、`tok_{ulid}`
- ✅ 线程安全、高并发友好

### 3. **密码存储：bcrypt**
- ✅ 默认 cost = 10
- ✅ 支持密码强度验证
- ✅ 预留密码历史检查扩展点

### 4. **Token 机制：JWT + Redis**
- ✅ Access Token: JWT 格式，30 分钟有效期
- ✅ Refresh Token: 随机 UUID，存储 Redis，7 天有效期
- ✅ 支持主动撤销（退出登录）

### 5. **错误处理：统一 AppError**
- ✅ 标准化错误响应格式
- ✅ 错误码 + 消息 + 元数据
- ✅ HTTP 状态码映射

### 6. **数据验证：三层模型**
- ✅ Handler 层：格式验证（validator）
- ✅ Domain 层：业务规则验证
- ✅ Repository 层：唯一性检查

---

## 🔑 核心功能实现

### **已实现的功能**
1. ✅ 用户注册（Register）
   - 邮箱唯一性检查
   - 密码 bcrypt 哈希
   - 自动生成 Token 对
   
2. ✅ 用户登录（Login）
   - 密码验证
   - 账户锁定保护（5 次失败）
   - 失败次数计数
   
3. ✅ 用户退出（Logout）
   - 撤销 Refresh Token
   - 清理会话
   
4. ✅ 刷新 Token（Refresh）
   - 验证 Refresh Token
   - 生成新的 Token 对
   - 更新最后登录时间

### **API 端点**
```
POST /api/v1/auth/register      # 注册
POST /api/v1/auth/login         # 登录
POST /api/v1/auth/logout        # 退出（需认证）
POST /api/v1/auth/refresh       # 刷新 Token
GET  /health                    # 健康检查
```

---

## 🛡️ 安全特性

### **已实现的安全措施**
1. ✅ **密码安全**
   - bcrypt 哈希
   - 最小长度 8 位
   - 最大长度 72 位（bcrypt 限制）

2. ✅ **账户保护**
   - 登录失败次数限制（5 次）
   - 失败后自动锁定
   - 重置机制

3. ✅ **Token 安全**
   - JWT 签名验证
   - Refresh Token 存储 Redis
   - 支持主动撤销

4. ✅ **请求验证**
   - 邮箱格式验证
   - 密码强度要求
   - 防止 SQL 注入（GORM 参数化）

### **待实现的安全加固**
- ⏳ 速率限制中间件
- ⏳ CORS 配置
- ⏳ 安全头（HSTS、XSS 防护）
- ⏳ 请求日志脱敏

---

## 📊 数据库设计

### **users 表**
```sql
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,          -- ULID 格式
    email VARCHAR(255) UNIQUE NOT NULL,  -- 邮箱（唯一索引）
    password VARCHAR(255) NOT NULL,      -- bcrypt 哈希
    email_verified BOOLEAN DEFAULT FALSE,-- 邮箱验证状态
    locked BOOLEAN DEFAULT FALSE,        -- 账户锁定状态
    failed_attempts INTEGER DEFAULT 0,   -- 失败次数
    last_login_at TIMESTAMP,             -- 最后登录时间
    created_at TIMESTAMP,                -- 创建时间
    updated_at TIMESTAMP                 -- 更新时间
);
```

### **索引**
- `idx_users_email`: 加速邮箱查询
- `idx_users_created_at`: 加速按创建时间排序

---

## 🔧 技术栈版本

```json
{
  "Go": "1.25.6",
  "Gin": "v1.12.0",
  "GORM": "v1.31.1",
  "PostgreSQL Driver": "v1.6.0",
  "Redis Client": "v9.18.0",
  "JWT": "v5.3.1",
  "bcrypt": "v0.49.0",
  "ULID": "v2.1.1",
  "Viper": "v1.21.0"
}
```

---

## 🚀 下一步建议

### **Phase 1: 完善核心功能** 🔥
- [ ] Worker 服务入口（Asynq）
- [ ] 领域事件发布/订阅
- [ ] 异步任务（发送邮件、审计日志）
- [ ] 单元测试（Domain 层）

### **Phase 2: 生产级加固** 🛡️
- [ ] 速率限制中间件
- [ ] CORS 配置
- [ ] Prometheus 指标
- [ ] 结构化日志（Zap）
- [ ] 集成测试

### **Phase 3: 功能扩展** 🚀
- [ ] 邮箱验证流程
- [ ] 密码找回流程
- [ ] 多设备登录支持
- [ ] RBAC 权限系统

---

## 📝 关键文件说明

### **1. `internal/auth/domain.go`**
- User 聚合根
- 密码哈希/验证逻辑
- 账户锁定逻辑
- 工厂方法 `NewUser()`

### **2. `internal/auth/service.go`**
- 应用服务（用例编排）
- 事务控制
- 领域事件发布
- DTO 定义（Commands/Responses）

### **3. `internal/auth/handler.go`**
- HTTP 处理器
- 请求验证
- 错误响应格式化
- 路由注册

### **4. `pkg/errors/errors.go`**
- 统一错误基类 `AppError`
- 预定义错误码
- HTTP 状态码映射

### **5. `pkg/utils/ulid/ulid.go`**
- ULID 生成工具
- 线程安全的熵源
- ID 格式化工具

---

## 💡 设计亮点

1. ✅ **垂直切片架构**：每个模块自包含，易于理解和维护
2. ✅ **依赖注入**：通过构造函数传递依赖，易于测试
3. ✅ **领域驱动设计**：清晰的领域模型，业务逻辑内聚
4. ✅ **务实的选择**：不过度设计，平衡优雅与实用
5. ✅ **生产级考虑**：错误处理、配置管理、优雅关闭

---

## 🎓 Go 设计哲学体现

1. **简洁优于复杂**
   - 没有过度抽象的接口
   - 直接的结构体组合

2. **显式优于隐式**
   - 错误显式返回
   - 依赖显式注入

3. **组合优于继承**
   - 使用 struct 组合
   - 没有复杂的继承层次

4. **接口小而精**
   - `UserRepository` 只有 5 个方法
   - 按需定义接口

5. **测试友好**
   - 依赖接口而非实现
   - 纯函数优先

---

## 📚 参考资源

- [DDD 快速入门](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservices-ddd-bounded-contexts/)
- [Go 最佳实践](https://github.com/golang-standards/project-layout)
- [Asynq 文档](https://github.com/hibiken/asynq)
- [Gin 框架](https://gin-gonic.com/docs/)
