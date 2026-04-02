# 活动日志调试指南

## 🔍 问题排查步骤

### 1️⃣ 检查中间件是否注册

**位置**: `backend/cmd/api/main.go`

```go
// 确保有以下代码
router.Use(activitylog.Middleware(activityLogService))
```

✅ **已修复**: 在最新提交中已经添加

---

### 2️⃣ 启动服务并观察调试日志

```bash
cd backend
go run ./cmd/api/main.go
```

**预期输出示例：**
```
[DEBUG] Activity log queued: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Action=REGISTER
[DEBUG] Middleware capturing activity: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Path=/api/v1/auth/register, Status=201
[DEBUG] Flushing batch of 1 activity logs
[DEBUG] Recording activity log: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Action=REGISTER, Status=SUCCESS
```

---

### 3️⃣ 测试触发活动日志

**方式 1：运行核心流程测试脚本**
```bash
scripts/dev/core-flow-test.sh
```

**方式 2：手动测试**
```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }'

# 2. 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }'

# 3. 获取当前用户信息（带 Token）
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 4️⃣ 查看数据库中的日志

**连接数据库：**
```bash
psql -h localhost -U postgres -d go_ddd_scaffold
```

**查询活动日志：**
```sql
-- 查看所有日志
SELECT id, user_id, email, action, status, created_at 
FROM activity_logs 
ORDER BY created_at DESC 
LIMIT 20;

-- 按用户 ID 查询
SELECT * FROM activity_logs 
WHERE user_id = 'YOUR_USER_ID' 
ORDER BY created_at DESC;

-- 统计每日日志数量
SELECT DATE(created_at) as date, COUNT(*) as count 
FROM activity_logs 
GROUP BY DATE(created_at) 
ORDER BY date DESC;

-- 查看不同动作的分布
SELECT action, status, COUNT(*) as count 
FROM activity_logs 
GROUP BY action, status 
ORDER BY count DESC;
```

---

### 5️⃣ 常见排查命令

```bash
# 检查表是否存在
psql -h localhost -U postgres -d go_ddd_scaffold -c "\dt activity_logs"

# 查看表结构
psql -h localhost -U postgres -d go_ddd_scaffold -c "\d activity_logs"

# 查看日志总数
psql -h localhost -U postgres -d go_ddd_scaffold -c "SELECT COUNT(*) FROM activity_logs;"

# 清空日志（测试用）
psql -h localhost -U postgres -d go_ddd_scaffold -c "TRUNCATE TABLE activity_logs;"
```

---

## 🐛 常见问题诊断

### 问题 1：没有调试日志输出

**可能原因：**
- 中间件未注册 ✅ 已修复
- 请求未经过认证（user_id 为空）
- 路径被过滤（健康检查等路径）

**解决方法：**
```bash
# 重启服务，观察启动日志
go run ./cmd/api/main.go

# 检查是否有以下日志
Starting server on port :8080
Swagger UI available at http://localhost:8080/swagger/index.html
```

---

### 问题 2：有队列日志但没有数据库写入

**可能原因：**
- 批量刷新间隔未到（默认 2 秒）
- 队列未满（默认 10 条）
- 数据库连接问题

**解决方法：**
```bash
# 等待 2-3 秒观察是否有批量刷新日志
# 或者发送更多请求填满队列（10 条）

# 检查数据库连接
psql -h localhost -U postgres -d go_ddd_scaffold -c "SELECT 1;"
```

---

### 问题 3：中间件没有捕获到用户信息

**可能原因：**
- JWT 认证失败
- Token 过期或无效
- authMiddleware 未正确设置 user_id

**调试方法：**
```bash
# 检查登录响应是否返回 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123!"}' | jq

# 确认响应中有 access_token 字段
# 使用正确的 token 格式：Bearer {token}
```

---

### 问题 4：数据库报错

**常见错误：**
```
ERROR: relation "activity_logs" does not exist
```

**解决方法：**
```bash
# 执行迁移
go run ./cmd/cli migrate up

# 验证表已创建
psql -h localhost -U postgres -d go_ddd_scaffold -c "\dt"
```

---

## 📊 完整的调试日志流程示例

**正常情况下的完整日志流：**

```
[GIN] 2026/04/02 - 18:30:00 | 201 |   15.234ms |         ::1 | POST     "/api/v1/auth/register"

[DEBUG] Middleware capturing activity: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Path=/api/v1/auth/register, Status=201
[DEBUG] Activity log queued: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Action=REGISTER

[DEBUG] Flushing batch of 1 activity logs
[DEBUG] Recording activity log: UserID=01KN6M7XKW01B9P6CAXF7Z0M4N, Action=REGISTER, Status=SUCCESS

✓ Database record created successfully
```

---

## 🎯 性能测试

**高并发场景下的日志表现：**

```bash
# 使用 ab 进行压力测试
ab -n 100 -c 10 \
   -H "Authorization: Bearer YOUR_TOKEN" \
   http://localhost:8080/api/v1/auth/me

# 观察日志输出
# 应该看到多条日志排队和批量刷新的过程
```

**预期行为：**
- 请求快速响应（异步写入不阻塞）
- 日志批量处理（每 10 条或 2 秒）
- 无数据库性能瓶颈

---

## 🔧 调整调试级别

**如果想减少调试日志输出：**

编辑 `backend/internal/activitylog/service.go`，注释掉以下行：

```go
// fmt.Printf("[DEBUG] Recording activity log: ...\n")
// fmt.Printf("[DEBUG] Activity log queued: ...\n")
// fmt.Printf("[DEBUG] Flushing batch of %d activity logs\n", ...)
// fmt.Printf("[DEBUG] Middleware capturing activity: ...\n")
```

**生产环境建议：**
- 使用日志库（如 logrus、zap）代替 fmt.Println
- 设置不同的日志级别（DEBUG/INFO/WARN/ERROR）
- 将调试日志输出到文件而非控制台

---

## 📝 下一步

调试完成后，记得：
1. 移除或注释掉调试日志（fmt.Printf）
2. 使用正式的日志库
3. 添加日志轮转和归档
4. 集成到监控系统（Prometheus + Grafana）

需要我帮你移除调试日志并优化吗？
