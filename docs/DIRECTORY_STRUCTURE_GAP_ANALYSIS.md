# 目录结构对比分析

## 📊 理想 vs 实际

### **理想目录结构**（来自 GOALS_AND_ACCEPTANCE_CRITERIA.md）

```
backend/internal/
├── domain/                    # ✅ 已实现
│   ├── user/                  # ✅ 已实现
│   │   └── events.go          # ✅ 领域事件
│   └── shared/                # ✅ 已存在
│
├── app/                       # ❌ 缺失 - 应用层
│   ├── user/                  # ❌ 缺失
│   │   ├── service.go         # ❌ UserService
│   │   ├── dto.go             # ❌ DTO
│   │   └── mapper.go          # ❌ Mapper
│   └── authentication/        # ❌ 缺失
│       ├── service.go         # ❌ AuthService
│       ├── dto.go             # ❌ DTO
│       └── mapper.go          # ❌ Mapper
│
├── listener/                  # ✅ 已实现
│   ├── audit_log_listener.go  # ✅ 审计日志监听器
│   └── activity_log_listener.go # ❌ 缺失 - 活动日志监听器
│
├── transport/                 # ⚠️ 部分实现
│   ├── http/                  # ❌ 缺失 - HTTP 传输层
│   │   ├── router.go          # ❌ 路由配置
│   │   ├── middleware.go      # ❌ 中间件
│   │   └── handlers/          # ❌ HTTP Handlers
│   │       ├── auth.go        # ❌ AuthHandler
│   │       ├── user.go        # ❌ UserHandler
│   │       └── audit_log.go   # ❌ AuditLogHandler
│   └── worker/                # ✅ 已实现
│       └── handlers/
│           ├── audit_log_handler.go    # ✅ 审计日志处理器
│           └── activity_log_handler.go # ❌ 缺失
│
└── infra/                     # ✅ 已实现
    ├── persistence/           # ❌ 缺失
    │   ├── db.go              # ❌ GORM 连接
    │   └── redis.go           # ❌ Redis 连接
    ├── repository/            # ✅ 已实现
    │   ├── user.go            # ❌ 缺失 - UserRepository
    │   ├── audit_log.go       # ✅ AuditLogRepository
    │   └── activity_log.go    # ❌ ActivityLogPersistence
    ├── messaging/             # ✅ 已实现
    │   └── event_bus.go       # ✅ EventBus
    └── redis/                 # ✅ 已实现
        ├── token_store.go     # ✅ Token Store
        └── device_store.go    # ❌ 缺失 - Device Store
```

---

### **当前实际目录结构**

```
backend/internal/
├── activitylog/               # ⚠️ 旧的活动日志系统（应该迁移或删除）
│   ├── domain.go
│   ├── handler.go
│   ├── repository.go
│   ├── repository_gorm.go
│   └── service.go
│
├── asynq/                     # ⚠️ 旧的 Asynq 任务系统（应该迁移或删除）
│   ├── handlers/
│   │   └── activity_log_handler.go
│   └── tasks/
│       └── activity_log.go
│
├── auth/                      # ⚠️ 混合了应用层和传输层
│   ├── event_bus_adapter.go   # ✅ 适配器
│   ├── handler.go             # ⚠️ HTTP Handler（应该在 transport/http/handlers/）
│   ├── repository.go          # ⚠️ Repository 接口（应该在 infra/repository/）
│   ├── repository_gorm.go     # ⚠️ Repository 实现（应该在 infra/repository/）
│   ├── service.go             # ⚠️ Service（应该在 app/authentication/）
│   ├── token_service.go       # ⚠️ TokenService（应该在 infra/redis/ 或 app/）
│   └── tasks.go               # ⚠️ Asynq 任务（应该在 asynq/tasks/ 或移除）
│
├── domain/                    # ✅ 符合规范
│   ├── shared/                # ✅ 空目录
│   └── user/
│       └── events.go          # ✅ 领域事件
│
├── infra/                     # ✅ 符合规范
│   ├── messaging/
│   │   └── event_bus.go       # ✅ EventBus
│   ├── redis/
│   │   └── token_store.go     # ✅ Token Store
│   └── repository/
│       └── audit_log.go       # ✅ AuditLogRepository
│
├── infrastructure/            # ⚠️ 与 infra/ 重复，应该合并
│   └── config/
│       └── config.go
│
├── listener/                  # ✅ 符合规范
│   ├── audit_log_listener.go  # ✅ 审计日志监听器
│   └── dto.go                 # ✅ DTO
│
├── middleware/                # ⚠️ 应该在 transport/http/middleware/
│   ├── auth.go
│   ├── cors.go
│   ├── prometheus.go
│   └── ratelimit.go
│
└── transport/                 # ⚠️ 只有 worker，缺少 http
    └── worker/
        └── handlers/
            └── audit_log_handler.go  # ✅ Worker Handler
```

---

## 🗑️ **需要删除的文件/目录**

### **1. activitylog/** （旧的活动日志系统）

**原因**：
- ❌ 使用旧的同步/Asynq 直接发布方式
- ❌ 与新架构（EventBus → Listener → Worker）冲突
- ❌ 功能已被新的 `listener/` + `transport/worker/` 替代

**文件列表**：
```
internal/activitylog/
├── domain.go              ❌ 删除
├── handler.go             ❌ 删除
├── repository.go          ❌ 删除
├── repository_gorm.go     ❌ 删除
└── service.go             ❌ 删除
```

**依赖检查**：
- ✅ `auth/handler.go` 使用了 `activitylog.Service`（5 处引用）
- ✅ `cmd/api/main.go` 初始化了 activitylog service
- ✅ `cmd/worker/main.go` 注册了 activity log handler

**结论**：**暂时保留**，需要先迁移依赖后再删除

---

### **2. asynq/** （旧的 Asynq 任务系统）

**原因**：
- ❌ 与新的事件驱动架构重复
- ❌ 应该统一使用 EventBus

**文件列表**：
```
internal/asynq/
├── handlers/
│   └── activity_log_handler.go  ❌ 删除
└── tasks/
    └── activity_log.go          ❌ 删除
```

**依赖检查**：
- ✅ `cmd/worker/main.go` 注册了 activity:record handler

**结论**：**暂时保留**，等 activitylog 迁移完成后再删除

---

### **3. infrastructure/config/** （重复的目录）

**原因**：
- ❌ 与 `infra/` 命名不一致
- ❌ 应该合并到 `infra/persistence/` 或单独的 `config/` 包

**建议**：
- 移动到 `pkg/config/` 或保持现状但重命名为 `infra/config/`

---

## ➕ **缺少的文件/目录**

### **1. app/ 应用层**（核心缺失）⭐⭐⭐

**应该创建**：
```
internal/app/
├── user/
│   ├── service.go       # UserService
│   ├── dto.go           # UserDTO
│   └── mapper.go        # Entity ↔ DTO 转换
└── authentication/
    ├── service.go       # AuthService（从 auth/service.go 迁移）
    ├── dto.go           # AuthDTO
    └── mapper.go        # Entity ↔ DTO 转换
```

**优先级**：🔴 **高**

---

### **2. transport/http/** （HTTP 传输层）⭐⭐⭐

**应该创建**：
```
internal/transport/http/
├── router.go            # Gin 路由配置
├── middleware.go        # 中间件
└── handlers/
    ├── auth.go          # AuthHandler（从 auth/handler.go 迁移）
    ├── user.go          # UserHandler
    └── audit_log.go     # AuditLogQueryHandler
```

**优先级**：🔴 **高**

---

### **3. infra/persistence/** （持久化基础设施）⭐⭐

**应该创建**：
```
internal/infra/persistence/
├── db.go                # GORM 数据库连接
└── redis.go             # Redis 连接
```

**优先级**：🟡 **中**

---

### **4. infra/repository/user.go** （用户仓储）⭐⭐

**应该创建**：
```
internal/infra/repository/
└── user.go              # UserRepository GORM 实现
```

**当前状态**：
- `auth/repository.go` 定义了接口
- `auth/repository_gorm.go` 实现了 GORM
- 应该迁移到 `infra/repository/`

**优先级**：🟡 **中**

---

### **5. listener/activity_log_listener.go** （活动日志监听器）⭐

**应该创建**：
```
internal/listener/
└── activity_log_listener.go  # 订阅 USER.*, FEATURE.* 事件
```

**优先级**：🟢 **低**（可选，取决于是否需要活动日志）

---

### **6. infra/redis/device_store.go** （设备存储）⭐

**应该创建**：
```
internal/infra/redis/
└── device_store.go     # Device Store（Redis）
```

**当前状态**：
- Token 存储在 `token_store.go` 中
- 设备信息可能也在其中，但未分离

**优先级**：🟢 **低**

---

## 📋 **迁移计划**

### **阶段 1：整理现有代码**（1-2 天）

#### **步骤 1.1：移动 auth/ 到正确位置**

```bash
# 移动 Service 到 app/authentication/
mv internal/auth/service.go internal/app/authentication/service.go

# 移动 Handler 到 transport/http/handlers/
mkdir -p internal/transport/http/handlers
mv internal/auth/handler.go internal/transport/http/handlers/auth.go

# 移动 Repository 到 infra/repository/
mv internal/auth/repository*.go internal/infra/repository/

# 移动 TokenService 到 infra/redis/
mv internal/auth/token_service.go internal/infra/redis/
```

#### **步骤 1.2：更新 import 路径**

在所有文件中更新 import：
```go
// 旧
import "github.com/shenfay/go-ddd-scaffold/internal/auth"

// 新
import (
    "github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
    "github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
)
```

---

### **阶段 2：删除旧系统**（1 天）

#### **步骤 2.1：迁移 activitylog 依赖**

1. 将 `activitylog.Service` 改为通过 EventBus 发布
2. 创建 `ActivityLogListener`
3. 更新 `auth/handler.go` 中的调用

#### **步骤 2.2：删除旧文件**

```bash
rm -rf internal/activitylog/
rm -rf internal/asynq/
```

---

### **阶段 3：补充缺失组件**（2-3 天）

#### **步骤 3.1：创建 app/ 层**

```bash
mkdir -p internal/app/{user,authentication}
```

#### **步骤 3.2：创建 transport/http/ 层**

```bash
mkdir -p internal/transport/http/handlers
```

#### **步骤 3.3：完善 infra/ 层**

```bash
mkdir -p internal/infra/persistence
touch internal/infra/persistence/{db,redis}.go
```

---

## 🎯 **最终目标目录结构**

```
backend/internal/
├── domain/                    # ✅ 领域层
│   ├── user/
│   │   └── events.go
│   └── shared/
│
├── app/                       # ✅ 应用层（待创建）
│   ├── user/
│   │   ├── service.go
│   │   ├── dto.go
│   │   └── mapper.go
│   └── authentication/
│       ├── service.go
│       ├── dto.go
│       └── mapper.go
│
├── listener/                  # ✅ 监听器层
│   ├── audit_log_listener.go
│   └── activity_log_listener.go  # 待创建
│
├── transport/                 # ✅ 传输层（待完善）
│   ├── http/
│   │   ├── router.go
│   │   ├── middleware.go
│   │   └── handlers/
│   │       ├── auth.go
│   │       ├── user.go
│   │       └── audit_log.go
│   └── worker/
│       └── handlers/
│           ├── audit_log_handler.go
│           └── activity_log_handler.go  # 待创建
│
└── infra/                     # ✅ 基础设施层（待完善）
    ├── persistence/
    │   ├── db.go
    │   └── redis.go
    ├── repository/
    │   ├── user.go
    │   ├── audit_log.go
    │   └── activity_log.go
    ├── messaging/
    │   └── event_bus.go
    └── redis/
        ├── token_store.go
        └── device_store.go
```

---

## 📊 **完成度评估**

| 层级 | 理想状态 | 当前状态 | 完成度 |
|------|----------|----------|--------|
| **domain/** | ✅ 完整 | ✅ 完整 | **100%** |
| **app/** | ✅ 完整 | ❌ 缺失 | **0%** |
| **listener/** | ✅ 2 个监听器 | ⚠️ 1 个监听器 | **50%** |
| **transport/http/** | ✅ 完整 | ❌ 缺失 | **0%** |
| **transport/worker/** | ✅ 完整 | ⚠️ 部分 | **50%** |
| **infra/messaging/** | ✅ 完整 | ✅ 完整 | **100%** |
| **infra/redis/** | ⚠️ 2 个 Store | ⚠️ 1 个 Store | **50%** |
| **infra/repository/** | ⚠️ 3 个 Repo | ⚠️ 1 个 Repo | **33%** |
| **infra/persistence/** | ✅ 完整 | ❌ 缺失 | **0%** |

**总体完成度**: **~40%**

---

## 💡 **建议行动方案**

### **方案 A：渐进式重构**（推荐）⭐⭐⭐

**优点**：
- ✅ 风险低，每次只改一点
- ✅ 可以边改边测试
- ✅ 不影响现有功能

**步骤**：
1. 先整理 `auth/` 目录（1-2 天）
2. 再删除 `activitylog/` 和 `asynq/`（1 天）
3. 最后补充缺失的组件（2-3 天）

**总时间**: 4-6 天

---

### **方案 B：保持现状**（保守）⭐⭐

**优点**：
- ✅ 无需大量重构
- ✅ 当前系统可正常工作

**缺点**：
- ❌ 不符合 DDD 规范
- ❌ 目录结构混乱
- ❌ 技术债务累积

---

### **方案 C：彻底重构**（激进）⭐

**优点**：
- ✅ 完全符合规范
- ✅ 清晰的架构分层

**缺点**：
- ❌ 工作量大（7-10 天）
- ❌ 风险高，可能引入 bug
- ❌ 需要全面回归测试

---

## 🏁 **结论**

### **需要删除的文件**

| 文件/目录 | 状态 | 原因 |
|-----------|------|------|
| `internal/activitylog/` | ⚠️ 暂保留 | 有依赖，需先迁移 |
| `internal/asynq/` | ⚠️ 暂保留 | 有依赖，需先迁移 |
| `internal/infrastructure/` | ⚠️ 待重命名 | 与 `infra/` 重复 |

### **需要创建的文件**

| 文件/目录 | 优先级 | 说明 |
|-----------|--------|------|
| `internal/app/` | 🔴 高 | 应用层（核心缺失） |
| `internal/transport/http/` | 🔴 高 | HTTP 传输层（核心缺失） |
| `internal/infra/persistence/` | 🟡 中 | 持久化基础设施 |
| `internal/infra/repository/user.go` | 🟡 中 | 用户仓储 |
| `internal/listener/activity_log_listener.go` | 🟢 低 | 活动日志监听器 |
| `internal/infra/redis/device_store.go` | 🟢 低 | 设备存储 |

---

**下一步建议**：
1. 先执行**方案 A 的步骤 1**：整理 `auth/` 目录
2. 验证编译通过
3. 再决定是否继续删除旧文件或补充新组件

您希望我执行哪个方案？
