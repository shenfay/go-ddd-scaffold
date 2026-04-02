# 活动日志调试进展报告

## 📊 当前状态

### ✅ 已确认正常的功能

1. **中间件注册成功** ✅
   - `activitylog.Middleware()` 已在 main.go 中注册
   - 中间件能正确捕获认证用户的请求

2. **异步队列工作正常** ✅
   - 日志成功进入队列：`Activity log queued`
   - 批量处理机制正常工作：`Flushing batch of 3 activity logs`

3. **部分日志能成功写入** ✅
   - PROFILE_UPDATE 类型的日志可以成功插入（rows:1）
   - 数据库连接和表结构正常

---

### 🐛 发现的问题

#### 问题 1：REGISTER 类型日志插入失败 ❌

**现象：**
```
[DEBUG] Recording activity log: UserID=01KN7CD38VXYPWNVTMT7GMC6E7, Action=REGISTER, Status=SUCCESS
[DEBUG] Inserting activity log: ID=ses_01KN7CD38YDWJT7C37A6CKNF0T, UserID=01KN7CD38VXYPWNVTMT7GMC6E7, Action=REGISTER
[DEBUG] Log details: Email=testuser_311518@test.com, IP=::1, Device=desktop, Browser=Other, OS=Other
[1.502ms] [rows:0] INSERT INTO "activity_logs" ...
```

**分析：**
- GORM 执行了 INSERT 语句但返回 `[rows:0]`
- 没有报错，但实际没有插入数据
- 同样的代码，PROFILE_UPDATE 类型能成功插入

**可能的原因：**
1. **主键冲突**：生成的 session ID 可能与现有数据冲突
2. **字段验证问题**：某个字段的值不符合数据库约束
3. **GORM Hook 问题**：可能存在 BeforeCreate 钩子阻止了插入
4. **并发问题**：多个 goroutine 同时写入导致的问题

---

#### 问题 2：登录/登出时中间件未捕获到用户

**现象：**
```
[DEBUG] No user found in context, skipping activity log
```

**原因：**
- 登录/注册接口在认证之前执行
- 此时 user_id 还未被设置到 Gin 上下文中
- 这是**预期行为**，需要在业务逻辑中手动记录这些活动

**解决方案：**
在 Login/Register Handler 中手动调用 `activityLogService.RecordAsync()`

---

### 📈 测试数据统计

**成功的日志（3 条）：**
```sql
SELECT action, description, created_at 
FROM activity_logs 
ORDER BY created_at DESC;

action          | description         | created_at
----------------+---------------------+----------------------
PROFILE_UPDATE  | GET /api/v1/auth/me | 2026-04-02 23:18:49
PROFILE_UPDATE  | GET /api/v1/users/:id| 2026-04-02 23:18:49
PROFILE_UPDATE  | GET /api/v1/auth/me | 2026-04-02 23:18:49
```

**失败的日志（REGISTER x 1）：**
- 有详细的 DEBUG 日志
- GORM 执行了 INSERT
- 返回 rows:0，实际未插入

---

## 🔍 下一步调试计划

### 方案 1：检查主键生成

```bash
# 查看数据库中是否有重复的 ID
psql -h localhost -U shenfay -d ddd_scaffold -c \
  "SELECT id, COUNT(*) FROM activity_logs GROUP BY id HAVING COUNT(*) > 1;"
```

### 方案 2：简化测试

创建一个最小化的测试用例：
```go
func TestActivityLogInsert(t *testing.T) {
    log := &activitylog.ActivityLog{
        ID:       "test_123",
        UserID:   "user_123",
        Email:    "test@test.com",
        Action:   activitylog.ActivityRegister,
        Status:   activitylog.ActivitySuccess,
        IP:       "127.0.0.1",
        UserAgent: "test",
        Device:   "desktop",
        Browser:  "Chrome",
        OS:       "Windows",
        Description: "Test log",
        CreatedAt: time.Now(),
    }
    
    err := repo.Create(context.Background(), log)
    if err != nil {
        t.Fatal(err)
    }
}
```

### 方案 3：使用原生 SQL 测试

```sql
-- 直接插入测试
INSERT INTO activity_logs (
    id, user_id, email, action, status, ip, user_agent, 
    device, browser, os, description, metadata, created_at
) VALUES (
    'test_manual_001',
    'user_test',
    'test@example.com',
    'REGISTER',
    'SUCCESS',
    '127.0.0.1',
    'Test Agent',
    'desktop',
    'Chrome',
    'Windows',
    'Manual test insert',
    '{}',
    NOW()
);
```

---

## 💡 临时解决方案

### 方案 A：在业务层手动记录（绕过中间件）

修改 `auth/handler.go`，在 Register/Login 成功后手动记录：

```go
// 注册成功后
if h.activityLog != nil {
    _ = h.activityLog.Record(c.Request.Context(), activitylog.LogParams{
        UserID:      resp.User.ID,
        Email:       resp.User.Email,
        Action:      activitylog.ActivityRegister,
        Status:      activitylog.ActivitySuccess,
        IP:          c.ClientIP(),
        UserAgent:   c.GetHeader("User-Agent"),
        Description: "用户注册成功",
    })
}
```

### 方案 B：移除 Metadata 字段的 JSON tag

修改 `domain.go`：
```go
Metadata string `gorm:"type:json" json:"metadata,omitempty"`
// 改为
Metadata string `gorm:"type:jsonb"`
```

---

## 🎯 当前优先级

1. **高优先级** 🔴
   - 解决 REGISTER 类型日志插入失败的问题
   - 在 Login/Register Handler 中添加手动记录

2. **中优先级** 🟡
   - 添加集成测试覆盖所有活动类型
   - 优化错误处理和日志记录

3. **低优先级** 🟢
   - 移除调试日志
   - 性能优化和生产环境配置

---

## 📝 相关文件

- [`service.go`](file:///Users/shenfay/Projects/ddd-scaffold/backend/internal/activitylog/service.go) - 服务层和中间件
- [`repository_gorm.go`](file:///Users/shenfay/Projects/ddd-scaffold/backend/internal/activitylog/repository_gorm.go) - GORM 仓储实现
- [`domain.go`](file:///Users/shenfay/Projects/ddd-scaffold/backend/internal/activitylog/domain.go) - 领域模型
- [`main.go`](file:///Users/shenfay/Projects/ddd-scaffold/backend/cmd/api/main.go) - 应用入口和路由配置

---

## 📊 调试命令速查

```bash
# 查看最新日志
tail -100 /tmp/api.log | grep -E "(DEBUG|INSERT)"

# 查看数据库中的日志
PGPASSWORD=postgres psql -h localhost -U shenfay -d ddd_scaffold \
  -c "SELECT * FROM activity_logs ORDER BY created_at DESC LIMIT 10;"

# 统计各类型日志数量
PGPASSWORD=postgres psql -h localhost -U shenfay -d ddd_scaffold \
  -c "SELECT action, COUNT(*) FROM activity_logs GROUP BY action;"

# 清空日志（重新测试）
PGPASSWORD=postgres psql -h localhost -U shenfay -d ddd_scaffold \
  -c "TRUNCATE TABLE activity_logs;"

# 运行测试脚本
bash scripts/dev/core-flow-test.sh
```
