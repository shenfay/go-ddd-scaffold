# WebSocket Integration 快速开始指南

## 5 分钟集成 WebSocket

### 第一步：安装 Skill

```bash
npx skills install websocket-integration
```

### 第二步：生成 WebSocket 基础设施

```bash
websocket-integration init \
  --project mathfun \
  --port 8080 \
  --max-connections 10000
```

这会生成完整的 WebSocket 架构：
- ✅ WebSocket Manager（统一管理）
- ✅ Connection 封装
- ✅ Connection Pool（连接池）
- ✅ Room Manager（房间管理）
- ✅ Heartbeat（心跳保活）

### 第三步：安装依赖

```bash
cd backend
go get github.com/gorilla/websocket
```

### 第四步：启动 WebSocket 服务

生成的代码已经可以使用，在 `cmd/server/main.go` 中添加：

```go
import (
    "your-project/internal/infrastructure/websocket"
)

func main() {
    // 创建 WebSocket Manager
    manager := websocket.NewManager(websocket.ManagerConfig{
        MaxConnections:    10000,
        HeartbeatInterval: 30 * time.Second,
        WriteTimeout:      10 * time.Second,
        ReadTimeout:       60 * time.Second,
    })
    
    // 启动 Manager
    go manager.Start()
    defer manager.Shutdown()
    
    // 注册 WebSocket 路由
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handleWebSocket(manager, w, r)
    })
    
    // 启动 HTTP 服务器
    http.ListenAndServe(":8080", nil)
}

func handleWebSocket(manager *websocket.Manager, w http.ResponseWriter, r *http.Request) {
    // 升级 HTTP 连接到 WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error("WebSocket upgrade failed:", err)
        return
    }
    
    // 生成连接 ID
    connID := generateConnectionID()
    
    // 添加到 Manager
    err = manager.AddConnection(connID, conn)
    if err != nil {
        logger.Error("Add connection failed:", err)
        return
    }
    
    // 处理消息循环
    for {
        message, err := manager.GetConnection(connID).ReadMessage()
        if err != nil {
            manager.RemoveConnection(connID)
            break
        }
        
        // 处理接收到的消息
        handleMessage(manager, connID, message)
    }
}
```

---

## 使用场景示例

### 场景 1：学生学习状态同步

```go
// 学生加入学习房间
manager.JoinRoom(studentConnID, "study-room-"+classID)

// 教师发送通知到房间
notification := map[string]interface{}{
    "type": "student_joined",
    "student_id": studentID,
    "timestamp": time.Now(),
}
manager.SendToRoom("study-room-"+classID, toJSON(notification))
```

### 场景 2：实时答题竞赛

```go
// 创建竞赛房间
room := roomManager.CreateRoom("competition-001", "数学竞赛", 50)

// 学生加入竞赛
manager.JoinRoom(studentConnID, "competition-001")

// 广播题目
question := map[string]interface{}{
    "question_id": "q001",
    "content": "1+1=?",
    "options": []string{"A. 1", "B. 2", "C. 3"},
}
manager.SendToRoom("competition-001", toJSON(question))

// 接收答案并统计
answers := make(map[string]string)
// ... 收集答案逻辑
```

### 场景 3：即时通知推送

```go
// 全局广播
announcement := map[string]interface{}{
    "type": "announcement",
    "title": "系统维护通知",
    "content": "今晚 22:00-23:00 系统维护",
}
manager.BroadcastToAll(toJSON(announcement))

// 发送给特定用户
manager.SendToConnection(userConnID, toJSON(message))
```

---

## 房间管理

### 创建房间

```bash
# 添加学习房间
websocket-integration add room \
  --name study-room-001 \
  --max-users 30

# 添加竞赛房间
websocket-integration add room \
  --name competition-room \
  --type game \
  --max-users 100
```

### 手动创建房间

```go
// 创建房间
room := roomManager.CreateRoom(
    "study-room-001",
    "数学学习室",
    30, // 最大人数
)

// 加入房间
err := manager.JoinRoom(room.ID, connection)

// 离开房间
err = manager.LeaveRoom(room.ID, connection)

// 获取房间信息
room, exists := roomManager.GetRoom("study-room-001")
if exists {
    fmt.Printf("房间人数：%d\n", room.UserCount())
}
```

---

## 心跳配置

### 命令行配置

```bash
websocket-integration config heartbeat \
  --interval 30s \
  --timeout 90s \
  --max-retries 3
```

### 代码配置

```go
manager := websocket.NewManager(websocket.ManagerConfig{
    HeartbeatInterval: 30 * time.Second,  // 每 30 秒发送 ping
})

// 为单个连接设置心跳
connection.SetHeartbeat(
    30*time.Second,   // 心跳间隔
    90*time.Second,   // 超时时间
    func() {          // 超时回调
        logger.Warn("连接超时，断开连接")
    },
)
```

---

## 连接池使用

```go
// 创建连接池
pool := websocket.NewConnectionPool(1000)

// 获取连接
conn := pool.Get()
if conn == nil {
    // 池为空，创建新连接
    conn = websocket.NewConnection(wsConn, writeTimeout, readTimeout)
} else {
    // 复用连接
    conn.Reset(wsConn)
}

// 归还连接
defer pool.Put(conn)
```

---

## 监控和统计

### 获取统计信息

```go
stats := manager.GetStats()
fmt.Printf("总连接数：%d\n", stats.TotalConnections)
fmt.Printf("总房间数：%d\n", stats.TotalRooms)
fmt.Printf("连接池大小：%d\n", stats.PoolSize)
```

### Prometheus 指标

生成的代码包含 Prometheus 指标导出：

```go
// 在 /metrics 端点暴露指标
prometheus.MustRegister(websocketConnectionsTotal)
prometheus.MustRegister(websocketMessagesTotal)
prometheus.MustRegister(websocketErrorsTotal)
```

---

## 错误处理

```go
// 添加连接
err := manager.AddConnection(connID, conn)
if err != nil {
    if errors.Is(err, websocket.ErrMaxConnectionsReached) {
        // 连接数已达上限
        http.Error(w, "Server full", http.StatusServiceUnavailable)
    }
    return
}

// 发送消息
err = manager.SendToRoom(roomID, message)
if err != nil {
    if errors.Is(err, websocket.ErrRoomNotFound) {
        // 房间不存在
        logger.Error("Room not found:", roomID)
    } else if errors.Is(err, websocket.ErrConnectionClosed) {
        // 连接已关闭
        manager.RemoveConnection(connID)
    }
}
```

---

## 安全配置

### JWT 认证

```go
// WebSocket 连接时的 Token 验证
func authenticate(r *http.Request) (string, error) {
    token := r.URL.Query().Get("token")
    if token == "" {
        token = r.Header.Get("Authorization")
    }
    
    claims, err := parseJWT(token)
    if err != nil {
        return "", err
    }
    
    return claims.UserID, nil
}
```

### IP 限流

```go
// 限制每个 IP 的连接数
ipConnections := make(map[string]int)

func checkIPLimit(ip string) bool {
    return ipConnections[ip] < 10
}
```

---

## 下一步

### 1. 集成业务逻辑

在 `handleMessage` 函数中实现业务逻辑：

```go
func handleMessage(manager *websocket.Manager, connID string, message []byte) {
    var msg Message
    json.Unmarshal(message, &msg)
    
    switch msg.Type {
    case "join_room":
        manager.JoinRoom(connID, msg.RoomID)
    case "leave_room":
        manager.LeaveRoom(connID, msg.RoomID)
    case "send_message":
        manager.SendToRoom(msg.RoomID, message)
    }
}
```

### 2. 添加持久化

将重要消息持久化到数据库：

```go
func persistMessage(roomID string, message []byte) {
    // 保存到数据库
    db.Exec("INSERT INTO messages (room_id, content, created_at) VALUES (?, ?, ?)",
        roomID, message, time.Now())
}
```

### 3. 集群支持

使用 Redis Pub/Sub 实现跨节点通信：

```go
// 使用 Redis 广播
redisClient.Publish(ctx, "websocket:"+roomID, message)

// 订阅其他节点的消息
pubsub := redisClient.Subscribe(ctx, "websocket:*")
```

---

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：咨询 DDD Architect Agent

祝你开发顺利！🚀
