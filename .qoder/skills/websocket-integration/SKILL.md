---
name: websocket-integration
description: WebSocket 快速集成工具。生成 WebSocket Manager、连接池管理、房间管理逻辑和心跳保活机制。适用于实时通信场景，如在线学习状态同步、即时通知、协作学习等。
version: "1.0.0"
author: MathFun Team
tags: [websocket, real-time, connection-pool, room-manager, heartbeat, gorilla-websocket]
---

# WebSocket Integration - WebSocket 快速集成工具

## 功能概述

这是一个智能化的 WebSocket 集成工具，专为 MathFun 项目设计。它提供完整的 WebSocket 连接管理、房间系统、心跳保活和消息广播机制，基于 **gorilla/websocket** 库实现，支持大规模并发连接。

## 核心能力

### 1. WebSocket Manager
- **连接管理** - 统一的 WebSocket 连接生命周期管理
- **并发控制** - 支持数千并发连接的优化处理
- **错误恢复** - 自动重连和异常处理机制
- **性能监控** - 连接数、消息量实时监控

### 2. 连接池管理
- **池化复用** - 连接对象池减少 GC 压力
- **资源限制** - 最大连接数控制和超限处理
- **健康检查** - 定期检测连接可用性
- **优雅关闭** - 连接清理和资源释放

### 3. 房间管理系统
- **动态分组** - 支持动态创建和加入房间
- **广播机制** - 房间内消息广播和单播
- **权限控制** - 房间访问权限验证
- **人数统计** - 实时房间人数监控

### 4. 心跳保活机制
- **定时心跳** - 可配置的心跳间隔
- **超时检测** - 连接超时自动断开
- **断线重连** - 客户端断线自动重连
- **状态同步** - 断线重连后状态恢复

## 使用场景

### 适用情况
- 学生在线学习状态实时同步
- 教师监控学生学习进度
- 即时通知和消息推送
- 协作学习和讨论室
- 实时答题和竞赛
- 在线客服系统

### 不适用情况
- 简单的请求响应场景
- 低频数据更新
- 纯静态内容分发

## 基本使用

### 快速开始（5 分钟）

```bash
# 1. 生成 WebSocket 基础架构
websocket-integration init --project mathfun

# 2. 安装依赖
go get github.com/gorilla/websocket

# 3. 启动 WebSocket 服务
go run cmd/server/main.go
```

### 添加房间管理

```bash
# 添加学习房间
websocket-integration add room --name study-room --max-users 30

# 添加竞赛房间
websocket-integration add room --name competition-room --type game
```

### 配置心跳

```bash
# 配置心跳参数
websocket-integration config heartbeat \
  --interval 30s \
  --timeout 90s \
  --max-retries 3
```

## 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--project` | string | 是 | - | 项目名称 |
| `--port` | int | 否 | 8080 | WebSocket 端口 |
| `--max-connections` | int | 否 | 10000 | 最大连接数 |
| `--heartbeat-interval` | duration | 否 | 30s | 心跳间隔 |
| `--room-type` | string | 否 | default | 房间类型 |
| `--broadcast` | flag | 否 | false | 启用全局广播 |

## 生成的代码结构

```
backend/
├── internal/
│   ├── infrastructure/
│   │   └── websocket/
│   │       ├── manager.go          # WebSocket Manager
│   │       ├── connection.go       # 连接封装
│   │       ├── connection_pool.go  # 连接池
│   │       ├── room.go             # 房间管理
│   │       ├── room_manager.go     # 房间管理器
│   │       ├── heartbeat.go        # 心跳管理
│   │       └── message.go          # 消息处理
│   │
│   └── domain/
│       └── websocket/
│           ├── entity/
│           │   ├── connection.go
│           │   └── room.go
│           └── service/
│               └── notification_service.go
│
└── pkg/
    └── websocket/
        ├── config.go              # 配置
        └── middleware/
            ├── auth.go            # 认证中间件
            └── ratelimit.go       # 限流中间件
```

## 最佳实践

### 连接管理

```go
// ✅ 正确：使用 Manager 统一管理
manager := websocket.NewManager(websocket.ManagerConfig{
    MaxConnections: 10000,
    HeartbeatInterval: 30 * time.Second,
})

go manager.Start()
defer manager.Shutdown()

// ❌ 错误：直接操作原始连接
conn, _ := upgrader.Upgrade(w, r, nil)
go handleConnection(conn) // 无管理、无清理
```

### 房间使用

```go
// ✅ 正确：使用房间管理器
roomManager := websocket.NewRoomManager()
room := roomManager.CreateRoom("study-room-001", 30)
room.Join(connection)
room.Broadcast(message)

// ❌ 错误：手动维护房间映射
rooms := make(map[string][]*websocket.Conn)
rooms["room1"] = append(rooms["room1"], conn)
```

### 心跳处理

```go
// ✅ 正确：启用心跳检测
connection.SetHeartbeat(
    30*time.Second,  // 心跳间隔
    90*time.Second,  // 超时时间
    func() {         // 超时回调
        logger.Warn("connection timeout")
    },
)

// ❌ 错误：不处理心跳
for {
    _, msg, _ := conn.ReadMessage()
    handleMessage(msg) // 无超时检测
}
```

### 错误处理

```go
// ✅ 正确：完善的错误处理
err := manager.SendToRoom("room-001", message)
if err != nil {
    if errors.Is(err, websocket.ErrRoomNotFound) {
        logger.Error("room not found")
    } else if errors.Is(err, websocket.ErrConnectionClosed) {
        logger.Warn("connection closed")
    }
}

// ❌ 错误：忽略错误
manager.SendToRoom("room-001", message) // 不检查错误
```

## 故障排除

### 常见问题

**连接频繁断开**
- 检查心跳间隔是否过长
- 验证防火墙设置
- 查看服务器日志确认错误原因

**内存占用过高**
- 检查连接池是否正确释放
- 验证最大连接数限制
- 使用 pprof 分析内存泄漏

**消息丢失**
- 确认房间广播逻辑
- 检查消息队列是否阻塞
- 验证连接状态

### 获取帮助
- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- 🚀 快速开始：查看 [QUICKSTART.md](./QUICKSTART.md)

## 版本历史

- v1.0.0 (2026-02-25): 初始版本发布
  - WebSocket Manager 完整实现
  - 连接池管理
  - 房间管理系统
  - 心跳保活机制
  - 消息广播功能

---
*本技能遵循 Qoder Skills 规范，专为 MathFun 项目优化设计*
