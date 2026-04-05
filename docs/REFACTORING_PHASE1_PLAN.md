# 重构执行计划 - 阶段 1

## ⚠️ **重要发现**

在执行文件移动后，发现 `auth/` 目录的文件有复杂的相互依赖关系：

### **依赖关系图**

```
auth/service.go
├── 依赖: auth/repository.go (UserRepository)
├── 依赖: auth/token_service.go (TokenService)
├── 依赖: domain/user/events.go
└── 被依赖: auth/handler.go

auth/handler.go
├── 依赖: auth/service.go (Service)
├── 依赖: auth/token_service.go (TokenService)
├── 依赖: activitylog.Service
├── 依赖: middleware.LoginRateLimit()
├── 依赖: apperrors.*
└── 被依赖: cmd/api/main.go

auth/token_service.go
├── 依赖: Redis Client
└── 被依赖: auth/service.go, auth/handler.go
```

---

## 🎯 **问题分析**

### **问题 1：TokenService 的归属**

**当前状态**：
- `TokenService` 在 `auth/` 包中
- 被 `Service` 和 `Handler` 同时使用
- 包含 JWT Token 生成逻辑

**应该放在哪里？**
- 方案 A：`infra/redis/token_service.go` - 作为基础设施层
- 方案 B：`app/authentication/token_service.go` - 作为应用服务的一部分
- 方案 C：保持现状，不移动

**建议**：**方案 B**（与应用服务一起）

**理由**：
- Token 生成是认证业务逻辑的一部分
- 不仅仅是 Redis 操作，还包含 JWT 签名
- 符合 DDD 应用层职责

---

### **问题 2：Handler 的依赖太多**

**当前状态**：
```go
import (
    "github.com/shenfay/go-ddd-scaffold/internal/activitylog"
    "github.com/shenfay/go-ddd-scaffold/internal/middleware"
    apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
)
```

**问题**：
- Handler 依赖了多个包的类型
- 移动后需要更新所有 import
- 可能引入编译错误

---

### **问题 3：Service 中的未导出字段**

**当前状态**：
```go
type Service struct {
    repo         UserRepository
    tokenService *TokenService  // 未导出字段
    eventBus     event.EventBus
}
```

**问题**：
- `tokenService` 是未导出字段
- Handler 无法直接访问 `service.tokenService`
- 需要改为导出字段或提供 getter 方法

---

## 💡 **修正后的执行方案**

### **方案 A：保守方案**（推荐）⭐⭐⭐

**策略**：只移动必要的文件，保持最小变更

**步骤**：

#### **步骤 1：只移动 Service 到 app/**

```bash
mv internal/auth/service.go internal/app/authentication/service.go
```

**修改内容**：
1. 更新 package 声明：`package auth` → `package authentication`
2. 更新 import 路径（如果有）
3. 将 `tokenService` 改为导出字段：`TokenService`

**影响范围**：
- ✅ `cmd/api/main.go` - 需要更新 import
- ✅ `auth/handler.go` - 需要更新 import

---

#### **步骤 2：保持其他文件不动**

**保留位置**：
- `auth/handler.go` - 保持在 `auth/`（暂时作为 transport 层）
- `auth/token_service.go` - 保持在 `auth/`（与 Service 配套）
- `auth/repository*.go` - 移动到 `infra/repository/`

**理由**：
- 减少一次性变更的范围
- 降低编译错误的风险
- 可以逐步验证每一步

---

#### **步骤 3：更新 import 路径**

**需要更新的文件**：
1. `cmd/api/main.go`
2. `cmd/worker/main.go`
3. `auth/handler.go`
4. `test/integration/*.go`

---

### **方案 B：激进方案**（不推荐）⭐

**策略**：一次性完成所有文件移动

**步骤**：
1. 移动所有文件到目标位置
2. 批量更新所有 import 路径
3. 修复所有编译错误

**风险**：
- ❌ 可能引入大量编译错误
- ❌ 调试困难
- ❌ 回滚成本高

---

## 📋 **推荐的详细执行步骤**

### **阶段 1.1：准备阶段**

#### **步骤 1：备份当前状态**

```bash
git add .
git commit -m "backup: before refactoring phase 1"
```

---

#### **步骤 2：创建目标目录**

```bash
mkdir -p internal/app/authentication
mkdir -p internal/transport/http/handlers
mkdir -p internal/infra/repository
```

---

### **阶段 1.2：移动 Service**

#### **步骤 3：移动 service.go**

```bash
mv internal/auth/service.go internal/app/authentication/service.go
```

#### **步骤 4：更新 package 声明**

文件：`internal/app/authentication/service.go`

```go
// 修改前
package auth

// 修改后
package authentication
```

#### **步骤 5：修改未导出字段**

文件：`internal/app/authentication/service.go`

```go
// 修改前
type Service struct {
    repo         UserRepository
    tokenService *TokenService  // 未导出
    eventBus     event.EventBus
}

// 修改后
type Service struct {
    Repo         UserRepository
    TokenService *TokenService  // 导出
    EventBus     event.EventBus
}
```

---

### **阶段 1.3：移动 Repository**

#### **步骤 6：移动 repository 文件**

```bash
mv internal/auth/repository.go internal/infra/repository/user_repository.go
mv internal/auth/repository_gorm.go internal/infra/repository/user_repository_gorm.go
```

#### **步骤 7：更新 package 声明**

文件：`internal/infra/repository/user_repository.go`

```go
// 修改前
package auth

// 修改后
package repository
```

文件：`internal/infra/repository/user_repository_gorm.go`

```go
// 修改前
package auth

// 修改后
package repository
```

---

### **阶段 1.4：更新 Import 路径**

#### **步骤 8：更新 cmd/api/main.go**

```go
// 修改前
import (
    "github.com/shenfay/go-ddd-scaffold/internal/auth"
)

authService := auth.NewService(userRepo, tokenService)

// 修改后
import (
    "github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
)

userRepo := repository.NewUserRepository(db)
authService := authentication.NewService(userRepo, tokenService)
```

---

#### **步骤 9：更新 auth/handler.go**

```go
// 修改前
import (
    "github.com/shenfay/go-ddd-scaffold/internal/auth"
)

type Handler struct {
    service *auth.Service
}

func NewHandler(service *auth.Service, ...) *Handler {
    return &Handler{
        service: service,
    }
}

// 修改后
import (
    "github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
)

type Handler struct {
    service *authentication.Service
}

func NewHandler(service *authentication.Service, ...) *Handler {
    return &Handler{
        service: service,
    }
}
```

---

#### **步骤 10：更新测试文件**

文件：`test/integration/auth_integration_test.go`

```go
// 修改前
import (
    "github.com/shenfay/go-ddd-scaffold/internal/auth"
)

// 修改后
import (
    "github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
    "github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
)
```

---

### **阶段 1.5：验证编译**

#### **步骤 11：编译测试**

```bash
cd backend
go build ./cmd/api
go build ./cmd/worker
go test ./test/integration/...
```

#### **步骤 12：功能测试**

```bash
# 启动 API 服务
cd cmd/api && go run main.go

# 测试登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}'
```

---

## ⚠️ **风险评估**

| 风险项 | 概率 | 影响 | 缓解措施 |
|--------|------|------|----------|
| 编译错误 | 中 | 中 | 分步执行，每步验证 |
| 运行时错误 | 低 | 高 | 充分的功能测试 |
| 测试失败 | 中 | 中 | 更新测试代码 |
| 回滚困难 | 低 | 高 | Git 备份 |

---

## 🎯 **最终决策**

**建议采用方案 A：保守方案**

**理由**：
1. ✅ 风险可控，每次只改一点
2. ✅ 易于调试和回滚
3. ✅ 可以边改边测试
4. ✅ 不影响现有功能

**预计时间**：1-2 小时

**下一步**：
- 等待用户确认
- 开始执行阶段 1.1（准备阶段）

---

**您是否同意按照此方案执行？**

如果同意，我将立即开始执行**阶段 1.1：准备阶段**。
