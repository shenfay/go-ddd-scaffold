# 活动日志功能使用指南

## 📋 功能概述

活动日志功能用于记录用户在系统中的所有关键操作，包括登录、登出、注册、Token 刷新等认证相关操作，以及个人资料修改等业务操作。

## 🎯 核心特性

### 1. 自动记录
- ✅ User-Agent 自动解析（设备、浏览器、操作系统识别）
- ✅ IP 地址自动获取（支持 IPv6）
- ✅ 请求路径和方法自动捕获
- ✅ HTTP 状态码自动记录

### 2. 手动记录
```go
// 在业务代码中手动记录活动
err := activityLogService.Record(ctx, activitylog.LogParams{
    UserID:      "user-123",
    Email:       "user@example.com",
    Action:      activitylog.ActivityLogin,
    Status:      activitylog.ActivitySuccess,
    IP:          "192.168.1.100",
    UserAgent:   c.GetHeader("User-Agent"),
    Description: "用户登录成功",
    Metadata: map[string]interface{}{
        "login_method": "password",
        "device_type":  "mobile",
    },
})
```

### 3. 中间件自动记录
```go
// 在路由中使用中间件自动记录
router.Use(activitylog.Middleware(activityLogService))
```

## 📊 支持的活动类型

| 类型 | 常量 | 说明 |
|------|------|------|
| 登录 | `ActivityLogin` | 用户登录操作 |
| 登出 | `ActivityLogout` | 用户登出操作 |
| 注册 | `ActivityRegister` | 用户注册操作 |
| 刷新 Token | `ActivityRefreshToken` | 刷新访问令牌 |
| 修改密码 | `ActivityPasswordChange` | 更改密码 |
| 更新资料 | `ActivityProfileUpdate` | 更新个人信息 |
| 验证邮箱 | `ActivityEmailVerify` | 邮箱验证 |
| 账户锁定 | `ActivityAccountLock` | 系统锁定账户 |
| 账户解锁 | `ActivityAccountUnlock` | 解除锁定 |

## 🔧 数据库表结构

```sql
CREATE TABLE activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL INDEX,     -- 用户 ID（索引）
    email VARCHAR(255),                      -- 用户邮箱
    action VARCHAR(50) NOT NULL,             -- 活动类型
    status VARCHAR(20),                      -- 状态（SUCCESS/FAILED）
    ip VARCHAR(45),                          -- IP 地址（支持 IPv6）
    user_agent VARCHAR(500),                 -- User-Agent
    device VARCHAR(100),                     -- 设备类型
    browser VARCHAR(50),                     -- 浏览器名称
    os VARCHAR(50),                          -- 操作系统
    description TEXT,                        -- 活动描述
    metadata JSON,                           -- 元数据（JSON）
    created_at TIMESTAMP NOT NULL INDEX      -- 创建时间（索引）
);
```

## 🌐 API 接口

### 1. 获取当前用户的活动日志
```http
GET /api/v1/users/me/activity-logs
Authorization: Bearer {token}

Query Parameters:
- limit: 返回数量限制（默认 10，最大 50）
```

**响应示例：**
```json
[
  {
    "id": "ses_01KN6M7XKW01B9P6CAXF7Z0M4N",
    "user_id": "01KN6M7XKW01B9P6CAXF7Z0M4N",
    "email": "user@example.com",
    "action": "LOGIN",
    "status": "SUCCESS",
    "ip": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "device": "desktop",
    "browser": "Chrome",
    "os": "Windows",
    "description": "用户登录成功",
    "metadata": null,
    "created_at": "2026-04-02T16:30:00+08:00"
  }
]
```

### 2. 获取指定用户的活动日志
```http
GET /api/v1/users/{id}/activity-logs
Authorization: Bearer {token}

Query Parameters:
- limit: 返回数量限制（默认 20，最大 100）
- offset: 偏移量
```

## 💡 使用场景

### 1. 安全审计
- 记录所有登录尝试（成功/失败）
- 检测异常登录行为（频繁失败、异地登录）
- 账户锁定/解锁记录

### 2. 用户行为分析
- 统计用户活跃度
- 分析用户使用习惯
- 优化产品功能

### 3. 问题排查
- 追踪用户操作流程
- 定位问题发生时间点
- 还原问题现场

## 🚀 集成示例

### 在 main.go 中初始化
```go
// 1. 创建活动日志仓储
activityLogRepo := activitylog.NewActivityLogRepository(db)

// 2. 创建活动日志服务
activityLogService := activitylog.NewService(activityLogRepo)

// 3. 创建活动日志 Handler
activityLogHandler := activitylog.NewHandler(activityLogService)

// 4. 创建认证 Handler（传入活动日志服务）
authHandler := auth.NewHandler(authService, activityLogService)

// 5. 注册路由
authHandler.RegisterRoutes(v1)
activityLogHandler.RegisterRoutes(v1)
```

### 在业务逻辑中使用
```go
// 登录成功时记录
if loginSuccess {
    _ = activityLogService.Record(ctx, activitylog.LogParams{
        UserID:      userID,
        Email:       email,
        Action:      activitylog.ActivityLogin,
        Status:      activitylog.ActivitySuccess,
        IP:          clientIP,
        UserAgent:   userAgent,
        Description: "用户登录成功",
    })
}

// 登录失败时也记录
if loginFailed {
    _ = activityLogService.Record(ctx, activitylog.LogParams{
        UserID:      userID,
        Email:       email,
        Action:      activitylog.ActivityLogin,
        Status:      activitylog.ActivityFailed,
        IP:          clientIP,
        UserAgent:   userAgent,
        Description: "用户登录失败：密码错误",
    })
}
```

## 📈 性能优化建议

1. **异步写入**：对于高频操作，可以使用消息队列异步写入
2. **定期清理**：设置定时任务定期删除过期日志（如保留 90 天）
3. **索引优化**：根据查询频率添加合适的索引
4. **分区表**：数据量大时可以按时间分区

## 🔒 隐私保护

1. **敏感信息脱敏**：不要在 metadata 中存储密码等敏感信息
2. **访问控制**：严格控制活动日志的访问权限
3. **合规性**：遵循 GDPR 等数据保护法规

## 📝 下一步计划

- [ ] 在登录 Handler 中添加活动日志记录
- [ ] 在登出 Handler 中添加活动日志记录
- [ ] 在 Token 刷新 Handler 中添加活动日志记录
- [ ] 实现活动日志清理任务
- [ ] 添加异常活动检测和告警
- [ ] 实现活动日志导出功能
