# 数据库表结构设计决策

## 概述

本文档说明 `ddd-scaffold` 项目的数据库表结构设计决策，包括**为什么有些表不需要**以及**推荐的存储方案**。

---

## 一、现有表结构

### ✅ 已实现的表

```sql
-- 1. users 表（用户核心表）
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,      -- 密码哈希
    email_verified BOOLEAN DEFAULT FALSE,
    locked BOOLEAN DEFAULT FALSE,
    failed_attempts INT DEFAULT 0,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 2. audit_logs 表（审计日志）
CREATE TABLE audit_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    ip VARCHAR(45),
    user_agent VARCHAR(500),
    device VARCHAR(100),
    browser VARCHAR(50),
    os VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. activity_logs 表（活动日志 - 简化版）
CREATE TABLE activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 4. email_verification_tokens 表（邮箱验证 Token）
CREATE TABLE email_verification_tokens (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. password_reset_tokens 表（密码重置 Token）
CREATE TABLE password_reset_tokens (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## 二、不需要的表及原因

### ❌ 不需要的表清单

```sql
-- 1. credentials 表（凭证表）- 不需要
-- 2. tokens 表（访问令牌表）- 不需要  
-- 3. devices 表（设备信息表）- 不需要
```

---

### 1. credentials 表（不需要）

#### ❌ 为什么不推荐？

```sql
-- ❌ 不推荐：独立的 credentials 表
CREATE TABLE credentials (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) UNIQUE REFERENCES users(id),
    password_hash VARCHAR(255) NOT NULL,
    failed_attempts INT DEFAULT 0,
    locked BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**问题**：
1. **违反 DDD 原则**：Credentials 应该是 User 的值对象，不是独立实体
2. **性能下降**：每次登录都需要 JOIN 查询
3. **增加复杂度**：需要维护两个表的事务一致性
4. **无实际收益**：没有带来额外的业务价值

#### ✅ 推荐方案

**将认证信息作为 User 实体的属性**：

```go
// domain/user/entity.go
package user

type User struct {
    ID             string     // 用户 ID
    Email          string     // 邮箱
    Password       string     // 密码哈希（值对象）
    EmailVerified  bool       // 邮箱是否验证
    Locked         bool       // 账户是否锁定
    FailedAttempts int        // 登录失败次数
    LastLoginAt    *time.Time // 最后登录时间
    CreatedAt      time.Time  // 创建时间
    UpdatedAt      time.Time  // 更新时间
}

// 业务方法
func (u *User) VerifyPassword(password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) IsLocked() bool {
    return u.Locked
}

func (u *User) IncrementFailedAttempts(max int) {
    u.FailedAttempts++
    if u.FailedAttempts >= max {
        u.Locked = true
    }
}
```

**数据库映射**：

```sql
-- ✅ users 表包含所有认证信息
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,      -- 密码哈希
    email_verified BOOLEAN DEFAULT FALSE,
    locked BOOLEAN DEFAULT FALSE,
    failed_attempts INT DEFAULT 0,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**优势**：
- ✅ **符合 DDD**：User 是聚合根，包含所有相关属性
- ✅ **高性能**：单次查询获取所有信息
- ✅ **简洁**：减少不必要的表
- ✅ **事务一致性**：在同一事务中更新

---

### 2. tokens 表（不需要）

#### ❌ 为什么不推荐？

```sql
-- ❌ 不推荐：数据库表存储 Token
CREATE TABLE tokens (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) REFERENCES users(id),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    device_info JSONB,
    created_at TIMESTAMP
);
```

**问题**：
1. **数据量大**：用户每次登录都插入新记录
2. **查询慢**：需要索引维护和清理
3. **过期处理复杂**：需要定时任务删除过期 Token
4. **并发性能低**：高频写入影响数据库性能

#### ✅ 推荐方案

**使用 Redis 存储 Token**：

```go
// infra/redis/token_store.go
package redis

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
)

type TokenStore struct {
    client *redis.Client
}

// TokenData Token 数据结构
type TokenData struct {
    UserID     string `json:"user_id"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresAt  time.Time `json:"expires_at"`
    DeviceID   string `json:"device_id"`
}

// Store 存储 Token（7 天有效期）
func (s *TokenStore) Store(ctx context.Context, refreshToken string, data *TokenData) error {
    key := s.buildKey(refreshToken)
    value, _ := json.Marshal(data)
    
    return s.client.Set(ctx, key, value, 7*24*time.Hour).Err()
}

// Get 获取 Token 信息
func (s *TokenStore) Get(ctx context.Context, refreshToken string) (*TokenData, error) {
    key := s.buildKey(refreshToken)
    value, err := s.client.Get(ctx, key).Bytes()
    
    if err == redis.Nil {
        return nil, ErrTokenNotFound
    }
    
    var data TokenData
    json.Unmarshal(value, &data)
    return &data, nil
}

// Delete 删除 Token（登出时使用）
func (s *TokenStore) Delete(ctx context.Context, refreshToken string) error {
    key := s.buildKey(refreshToken)
    return s.client.Del(ctx, key).Err()
}

func (s *TokenStore) buildKey(refreshToken string) string {
    return "auth:token:" + refreshToken
}
```

**Redis Key 设计**：
```
Key 格式：auth:token:{refresh_token}
Value: {
  "user_id": "user123",
  "access_token": "eyJ...",
  "refresh_token": "dXNlcjoxMjM0NTY3ODkw",
  "expires_at": "2024-04-04T10:00:00Z",
  "device_id": "device456"
}
TTL: 7 天（自动过期）
```

**优势**：
- ✅ **高性能**：O(1) 查询速度
- ✅ **自动过期**：Redis TTL 机制，无需手动清理
- ✅ **黑名单机制**：登出时直接删除 Key
- ✅ **水平扩展**：Redis Cluster 支持高并发

---

### 3. devices 表（不需要）

#### ❌ 为什么不推荐？

```sql
-- ❌ 不推荐：独立设备表
CREATE TABLE devices (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) REFERENCES users(id),
    device_name VARCHAR(100),
    device_type VARCHAR(20),
    os VARCHAR(50),
    browser VARCHAR(50),
    ip VARCHAR(45),
    is_active BOOLEAN DEFAULT TRUE,
    last_active_at TIMESTAMP,
    created_at TIMESTAMP
);
```

**问题**：
1. **数据冗余**：每个设备多条记录
2. **查询复杂**：需要统计活跃设备数
3. **维护成本高**：需要定期清理不活跃设备
4. **联表查询**：增加数据库负担

#### ✅ 推荐方案

**方案 A：Redis 存储（推荐）**

```go
// infra/redis/device_store.go
package redis

type DeviceStore struct {
    client *redis.Client
}

// Device 设备信息
type Device struct {
    DeviceID     string    `json:"device_id"`
    DeviceName   string    `json:"device_name"`
    DeviceType   string    `json:"device_type"`  // mobile/tablet/desktop
    OS           string    `json:"os"`
    Browser      string    `json:"browser"`
    LastActiveAt time.Time `json:"last_active_at"`
    IsActive     bool      `json:"is_active"`
}

// AddDevice 添加设备
func (s *DeviceStore) AddDevice(ctx context.Context, userID string, device *Device) error {
    key := s.buildKey(userID)
    
    // 获取现有设备列表
    devices, _ := s.GetDevices(ctx, userID)
    
    // 检查是否已存在
    exists := false
    for i, d := range devices {
        if d.DeviceID == device.DeviceID {
            devices[i] = device  // 更新现有设备
            exists = true
            break
        }
    }
    
    if !exists {
        devices = append(devices, device)
    }
    
    // 限制最多 10 个设备
    if len(devices) > 10 {
        devices = devices[:10]
    }
    
    value, _ := json.Marshal(devices)
    return s.client.Set(ctx, key, value, 30*24*time.Hour).Err()
}

// GetDevices 获取用户的所有设备
func (s *DeviceStore) GetDevices(ctx context.Context, userID string) ([]*Device, error) {
    key := s.buildKey(userID)
    value, err := s.client.Get(ctx, key).Bytes()
    
    if err == redis.Nil {
        return []*Device{}, nil
    }
    
    var devices []*Device
    json.Unmarshal(value, &devices)
    return devices, nil
}

// RevokeDevice 撤销设备（踢出）
func (s *DeviceStore) RevokeDevice(ctx context.Context, userID, deviceID string) error {
    devices, _ := s.GetDevices(ctx, userID)
    
    filtered := make([]*Device, 0, len(devices))
    for _, d := range devices {
        if d.DeviceID != deviceID {
            filtered = append(filtered, d)
        }
    }
    
    value, _ := json.Marshal(filtered)
    key := s.buildKey(userID)
    return s.client.Set(ctx, key, value, 30*24*time.Hour).Err()
}

func (s *DeviceStore) buildKey(userID string) string {
    return "auth:devices:" + userID
}
```

**Redis Key 设计**：
```
Key 格式：auth:devices:{user_id}
Value: [
  {
    "device_id": "device123",
    "device_name": "Chrome on macOS",
    "device_type": "desktop",
    "os": "macOS 14.0",
    "browser": "Chrome 120",
    "last_active_at": "2024-04-03T10:00:00Z",
    "is_active": true
  }
]
TTL: 30 天（自动清理不活跃设备）
```

**方案 B：users 表 JSON 字段（备选）**

```sql
-- ✅ 备选方案：在 users 表中添加 JSON 字段
ALTER TABLE users ADD COLUMN devices JSONB DEFAULT '[]'::jsonb;

-- 创建 GIN 索引（可选，如果需要搜索设备）
CREATE INDEX idx_users_devices ON users USING GIN (devices);
```

**Go 代码映射**：

```go
type User struct {
    ID        string
    Email     string
    Password  string
    Devices   []Device `gorm:"type:jsonb;default:'[]'::jsonb"`
    // ... 其他字段
}

type Device struct {
    DeviceID     string    `json:"device_id"`
    DeviceName   string    `json:"device_name"`
    DeviceType   string    `json:"device_type"`
    LastActiveAt time.Time `json:"last_active_at"`
    IsActive     bool      `json:"is_active"`
}
```

**对比**：

| 维度 | Redis 方案 | JSON 字段方案 |
|------|-----------|-------------|
| **性能** | ⭐⭐⭐⭐⭐ O(1) | ⭐⭐⭐⭐ 需要解析 JSON |
| **扩展性** | ⭐⭐⭐⭐⭐ 独立 Redis Cluster | ⭐⭐⭐ users 表会变大 |
| **灵活性** | ⭐⭐⭐⭐⭐ 易于修改结构 | ⭐⭐⭐⭐ JSON 结构灵活 |
| **一致性** | ⭐⭐⭐ 最终一致性 | ⭐⭐⭐⭐⭐ 强一致性 |
| **推荐度** | ✅ **强烈推荐** | ⭐ 备选方案 |

**推荐**：使用 **Redis 方案**，因为：
- ✅ 设备信息是**高频读取、低频写入**
- ✅ 不需要强一致性（允许短暂不一致）
- ✅ 可以独立水平扩展

---

## 三、总结

### ✅ 推荐的表结构

```sql
-- 核心业务表
users                          # 用户表（包含认证信息）

-- 日志表
audit_logs                     # 审计日志表
activity_logs                  # 活动日志表（简化版）

-- Token 表（一次性 Token）
email_verification_tokens      # 邮箱验证 Token 表
password_reset_tokens          # 密码重置 Token 表

-- ❌ 不需要的表（使用 Redis 替代）
-- credentials                # → users.password
-- tokens                     # → Redis: auth:token:{refresh_token}
-- devices                    # → Redis: auth:devices:{user_id}
```

### 🎯 设计原则

1. **符合 DDD**：User 是聚合根，Credentials 是值对象
2. **高性能**：高频数据用 Redis，持久化数据用数据库
3. **简洁实用**：避免过度设计和不必要的表
4. **可扩展**：Redis Cluster 支持水平扩展

### 📊 存储方案对比

| 数据类型 | 推荐方案 | 理由 |
|---------|---------|------|
| **用户信息** | PostgreSQL users 表 | 需要强一致性、事务支持 |
| **密码哈希** | PostgreSQL users 表 | 用户实体的一部分 |
| **登录 Token** | Redis | 高频访问、自动过期 |
| **设备信息** | Redis | 高频读取、最终一致性可接受 |
| **审计日志** | PostgreSQL audit_logs 表 | 需要长期保存、合规要求 |
| **活动日志** | PostgreSQL activity_logs 表 | 中期保存、分析用途 |
| **一次性 Token** | PostgreSQL tokens 表 | 短期有效、需要关系约束 |

---

**文档版本**：v1.0  
**创建日期**：2026-04-03  
**状态**：已批准
