# 文件清理清单

## 📋 当前状态分析

### ✅ **需要保留的文件/目录**

#### **1. activitylog 包**（旧的活动日志系统）
```
backend/internal/activitylog/
├── domain.go              ✅ 保留 - 活动日志领域模型
├── handler.go             ✅ 保留 - HTTP Handler
├── repository.go          ✅ 保留 - Repository 接口
├── repository_gorm.go     ✅ 保留 - GORM 实现
└── service.go             ✅ 保留 - 服务层（被 auth/handler 使用）
```

**原因**：
- 用于记录业务活动日志（用户行为、功能使用等）
- 被 `auth/handler.go` 广泛使用（5 处引用）
- 与新的 audit_log 功能不同，互补存在

---

#### **2. asynq/handlers 和 asynq/tasks**（旧的 Worker 处理器）
```
backend/internal/asynq/
├── handlers/
│   └── activity_log_handler.go  ✅ 保留 - 处理 activity:record 任务
└── tasks/
    └── activity_log.go          ✅ 保留 - 定义 ActivityLogPayload
```

**原因**：
- 被 `cmd/worker/main.go` 使用
- 处理 `activity:record` 类型的任务
- 与新的 `audit.log.task` 任务类型不同

---

#### **3. listener 包**（新的审计日志监听器）
```
backend/internal/listener/
├── audit_log_listener.go  ✅ 保留 - 订阅领域事件
└── dto.go                 ✅ 保留 - AuditLogTask DTO
```

**原因**：
- 新架构的核心组件
- 将领域事件转换为 Worker 任务

---

#### **4. transport/worker/handlers**（新的 Worker 处理器）
```
backend/internal/transport/worker/handlers/
└── audit_log_handler.go   ✅ 保留 - 处理 audit.log.task 任务
```

**原因**：
- 新架构的 Worker 处理器
- 写入 audit_logs 表

---

### ❌ **可以删除的文件/目录**

经过检查，**目前没有可以安全删除的文件**。

所有现有的文件和目录都在被使用，且功能不重复：

| 包/目录 | 用途 | 是否在用 |
|---------|------|----------|
| `activitylog/` | 业务活动日志 | ✅ 是（auth/handler 使用） |
| `asynq/handlers/` | 活动日志 Worker | ✅ 是（cmd/worker 注册） |
| `asynq/tasks/` | 活动日志任务定义 | ✅ 是（handler 依赖） |
| `listener/` | 审计日志监听器 | ✅ 是（cmd/api 初始化） |
| `transport/worker/handlers/` | 审计日志 Worker | ✅ 是（cmd/worker 注册） |
| `infra/messaging/` | EventBus | ✅ 是（核心组件） |
| `infra/redis/` | Token Store | ✅ 是（token_service 使用） |
| `infra/repository/` | AuditLogRepository | ✅ 是（Worker handler 使用） |
| `domain/user/events.go` | 领域事件 | ✅ 是（service.go 使用） |

---

## 🎯 结论

### **不需要删除任何文件**

当前的架构设计是**双轨制**：

1. **活动日志系统**（旧）
   - 用途：记录用户行为、功能使用等业务活动
   - 技术栈：Asynq 直接发布任务
   - 数据表：`activity_logs`
   - 位置：`activitylog/` + `asynq/`

2. **审计日志系统**（新）
   - 用途：记录认证、安全相关事件
   - 技术栈：EventBus → Listener → Worker
   - 数据表：`audit_logs`
   - 位置：`listener/` + `transport/worker/` + `infra/`

这两个系统：
- ✅ **功能互补**：一个记录业务活动，一个记录安全审计
- ✅ **技术不同**：一个是直接 Asynq，一个是事件驱动
- ✅ **都必要**：满足不同场景的需求

---

## 💡 建议

### **方案 A：保持现状**（推荐）⭐⭐⭐

**优点**：
- ✅ 功能完整，满足多种需求
- ✅ 技术演进清晰（从简单到复杂）
- ✅ 无破坏性变更

**缺点**：
- ⚠️ 两套系统并存，可能增加维护成本
- ⚠️ 新人学习曲线稍陡

---

### **方案 B：统一为事件驱动**（长期目标）

**步骤**：
1. 将 `activitylog` 也改为通过 EventBus 发布
2. 创建 `ActivityLogListener`
3. 复用现有的 Worker Handler 机制
4. 逐步废弃旧的 `asynq/handlers`

**时间估计**：2-3 天

**风险**：中等（需要重构现有代码）

---

### **方案 C：完全替换**（不推荐）

**步骤**：
1. 删除 `activitylog/` 和 `asynq/`
2. 将所有日志记录改为通过 EventBus
3. 统一使用 `audit_logs` 表

**风险**：高（破坏性变更，影响现有功能）

---

## 📝 最终决策

**建议采用方案 A：保持现状**

理由：
1. 两个系统功能不同，都有存在的价值
2. 当前没有编译错误，系统运行正常
3. 未来可以根据实际需要决定是否统一
4. 避免不必要的重构风险

---

**清理状态**: ❌ **无需清理**  
**原因**: 所有文件都在使用中，功能不重复  
**下一步**: 验证系统功能，确保正常运行
