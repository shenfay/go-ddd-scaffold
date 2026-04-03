# 重构清理计划

## 📊 当前状态评估

### ✅ 已完成的核心重构

| 组件 | 状态 | 说明 |
|------|------|------|
| **EventBus** | ✅ 完成 | `internal/infra/messaging/event_bus.go` |
| **领域事件** | ✅ 完成 | `internal/domain/user/events.go` |
| **Listener** | ⚠️ 待修复 | `internal/listener/audit_log_listener.go`（编译错误） |
| **Redis TokenStore** | ✅ 完成 | `internal/infra/redis/token_store.go` |
| **Repository** | ✅ 完成 | `internal/infra/repository/audit_log.go` |
| **Worker Handler** | ✅ 完成 | `internal/transport/worker/handlers/audit_log_handler.go` |
| **适配器** | ✅ 完成 | `internal/auth/event_bus_adapter.go` |

---

## 🗑️ 需要清理的文件

### **1. 旧的事件定义**（可删除）⚠️

**文件**：`backend/internal/auth/events.go` (143 行)

**问题**：
- ❌ 与新架构重复（`domain/user/events.go`）
- ❌ 不符合 DDD 规范（在 auth 层定义领域事件）
- ❌ 仍在被 `service.go` 使用

**依赖关系**：
```go
// backend/internal/auth/service.go 第 96 行
event := NewUserRegisteredEvent(user.ID, user.Email, "", "")

// backend/internal/auth/service.go 第 163 行
event := NewUserLoggedInEvent(user.ID, user.Email, cmd.IP, cmd.UserAgent, true)
```

**处理方案**：
1. ✅ **暂时保留** - 需要先迁移 service.go 中的引用
2. ⬜ **迁移到新的领域事件**
3. ⬜ **删除旧文件**

---

### **2. Listener 目录**（暂时保留）

**文件**：`backend/internal/listener/audit_log_listener.go` (85 行)
**DTO**：`backend/internal/listener/dto.go` (30 行)

**状态**：
- ⚠️ 有编译错误（import 路径问题）
- ⚠️ 在 main.go 中被注释掉

**处理方案**：
1. ✅ **暂时保留** - 需要修复编译错误
2. ⬜ **修复后启用**
3. ⬜ **或删除**（如果确定不需要自动审计日志）

---

### **3. 文档目录**（建议整理）

**当前文档**：
```
docs/
├── GOALS_AND_ACCEPTANCE_CRITERIA.md          ✅ 保留
├── ARCHITECTURE_REFACTORING_SPEC_V2.md       ✅ 保留
├── DATABASE_SCHEMA_DESIGN.md                 ✅ 保留
├── REFACTORING_IMPLEMENTATION_PLAN.md        ✅ 保留
├── EVENT_DRIVEN_INTEGRATION_GUIDE.md         ✅ 保留
├── IMPLEMENTATION_SUMMARY.md                 ✅ 保留
├── QUICKSTART_EVENT_DRIVEN_LOGGING.md        ✅ 保留
├── COMPILATION_FIX_RECORD.md                 ✅ 保留（临时）
└── [其他旧文档...]                           ⚠️ 待检查
```

**建议**：
- ✅ 保留本次重构相关的 7 份核心文档
- ⬜ 删除或归档过时的文档

---

## 📋 清理步骤

### **阶段 1：迁移服务层事件**（优先级：高）⭐⭐⭐

#### **步骤 1.1：更新 auth/service.go**

**目标**：使用新的领域事件替代旧的事件

**修改前**：
```go
import "github.com/shenfay/go-ddd-scaffold/internal/auth"

// 创建注册事件
event := auth.NewUserRegisteredEvent(user.ID, user.Email, "", "")

// 创建登录事件
event := auth.NewUserLoggedInEvent(user.ID, user.Email, cmd.IP, cmd.UserAgent, true)
```

**修改后**：
```go
import "github.com/shenfay/go-ddd-scaffold/internal/domain/user/events"

// 创建注册事件
event := &events.UserRegistered{
    UserID:    user.ID,
    Email:     user.Email,
    Timestamp: time.Now(),
}

// 创建登录事件
event := &events.UserLoggedIn{
    UserID:    user.ID,
    Email:     user.Email,
    IP:        cmd.IP,
    UserAgent: cmd.UserAgent,
    Timestamp: time.Now(),
}
```

**注意**：需要确保新事件的字段名与旧事件兼容

---

#### **步骤 1.2：验证编译**

```bash
cd backend
go build ./cmd/api
go build ./cmd/worker
```

---

#### **步骤 1.3：删除旧事件文件**

```bash
rm backend/internal/auth/events.go
```

---

### **阶段 2：修复 Listener 编译**（优先级：中）⭐⭐

#### **步骤 2.1：检查编译错误**

```bash
cd backend
go build ./internal/listener/...
```

**预期错误**：
```
no required module provides package github.com/shenfay/go-ddd-scaffold/internal/domain/user/events
```

#### **步骤 2.2：可能的修复方案**

**方案 A**：检查 import 路径
```go
// audit_log_listener.go
import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/events"
)
```

**方案 B**：简化事件结构（如果太复杂）

**方案 C**：如果 Listener 不是必需的，可以考虑删除

---

#### **步骤 2.3：恢复 Listener 功能**

取消注释 `cmd/api/main.go` 中的代码：
```go
// 恢复以下代码
import (
    "github.com/shenfay/go-ddd-scaffold/internal/listener"
)

// 创建审计日志监听器
auditLogListener := listener.NewAuditLogListener(eventBus)
_ = auditLogListener // 保持引用，防止被 GC
pkglogger.Info("✓ Audit Log Listener registered")
```

---

### **阶段 3：最终验证**（优先级：高）⭐⭐⭐

#### **步骤 3.1：完整编译测试**

```bash
cd backend
go build ./cmd/api
go build ./cmd/worker
```

#### **步骤 3.2：功能测试**

```bash
# 启动 API 服务
cd cmd/api && go run main.go

# 测试登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}'

# 检查审计日志
psql $DATABASE_URL -c "SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;"
```

---

## 🎯 建议的清理顺序

### **立即可做**（今天）

1. ✅ **评估哪些文件真正需要删除**
   - 检查每个文件的实际使用情况
   - 确认没有循环依赖

2. ⬜ **迁移 auth/events.go 到新位置**
   - 更新 service.go 中的引用
   - 删除旧文件

3. ⬜ **验证基本功能**
   - 确保登录接口正常工作
   - 验证 EventBus 工作

---

### **短期完成**（本周）

1. ⬜ **修复 Listener 编译错误**
   - 或者决定是否需要这个功能

2. ⬜ **整理文档**
   - 保留核心文档
   - 删除或归档过时文档

3. ⬜ **更新实施总结**
   - 记录最终的完成度
   - 列出所有变更

---

## ⚠️ 风险提示

### **高风险操作**

1. **直接删除 auth/events.go**
   - ❌ 会导致 service.go 编译失败
   - ✅ 必须先迁移引用

2. **直接删除 listener 目录**
   - ❌ 如果将来需要审计日志会很麻烦
   - ✅ 应该先修复或明确决定不需要

---

### **安全策略**

1. **渐进式清理**
   - 每次只修改一个文件
   - 每次修改后立即验证编译

2. **Git 备份**
   ```bash
   git add .
   git commit -m "backup: before cleanup"
   ```

3. **回滚计划**
   ```bash
   # 如果出现问题
   git checkout HEAD~1
   ```

---

## 📝 决策清单

### **需要用户确认的事项**

- [ ] **是否真的需要 Listener 功能？**
  - 是 → 修复编译错误并启用
  - 否 → 删除 listener 目录和相关代码

- [ ] **是否要统一事件定义？**
  - 是 → 迁移到 domain/user/events
  - 否 → 保留现状，但会有重复

- [ ] **文档如何处理？**
  - 保留所有重构文档（推荐）
  - 只保留最核心的 3-4 份

---

## 🏁 重构完成标准

### **代码层面**

- [x] EventBus 实现完成
- [x] 领域事件定义完成
- [x] Repository 实现完成
- [x] Worker Handler 实现完成
- [ ] **Listener 编译通过** ← 待完成
- [ ] **旧事件文件删除** ← 待完成

### **功能层面**

- [x] 数据库迁移完成
- [x] API 可以编译运行
- [x] EventBus 正常工作
- [ ] **审计日志自动记录** ← 待完成（依赖 Listener）

### **文档层面**

- [x] 实施方案文档
- [x] 集成指南文档
- [x] 快速启动指南
- [x] 编译错误修复记录

---

**当前状态**: ⚠️ **85% 完成**  
**剩余工作**: Listener 修复 + 旧文件清理  
**预计时间**: 2-3 小时
