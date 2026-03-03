#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
WebSocket Integration - WebSocket 快速集成工具
生成 WebSocket Manager、连接池、房间管理和心跳保活机制
"""

import os
import sys
import argparse
from pathlib import Path


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description='WebSocket 快速集成工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 初始化 WebSocket 项目
  %(prog)s init --project mathfun
  
  # 添加房间
  %(prog)s add room --name study-room --max-users 30
  
  # 配置心跳
  %(prog)s config heartbeat --interval 30s --timeout 90s
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # init 子命令
    init_parser = subparsers.add_parser('init', help='初始化 WebSocket 项目')
    init_parser.add_argument('--project', type=str, required=True, help='项目名称')
    init_parser.add_argument('--port', type=int, default=8080, help='WebSocket 端口')
    init_parser.add_argument('--max-connections', type=int, default=10000, help='最大连接数')
    init_parser.add_argument('--output', type=str, default='./backend', help='输出目录')
    
    # add 子命令
    add_parser = subparsers.add_parser('add', help='添加新元素')
    add_subparsers = add_parser.add_subparsers(dest='action', help='添加类型')
    
    # add room
    room_parser = add_subparsers.add_parser('room', help='添加房间')
    room_parser.add_argument('--name', type=str, required=True, help='房间名称')
    room_parser.add_argument('--max-users', type=int, default=100, help='最大用户数')
    room_parser.add_argument('--type', type=str, default='default', help='房间类型')
    
    return parser.parse_args()


def create_websocket_manager(project_name, output_dir):
    """生成 WebSocket Manager"""
    base_path = Path(output_dir) / 'internal' / 'infrastructure' / 'websocket'
    base_path.mkdir(parents=True, exist_ok=True)
    
    # manager.go
    manager_code = """// WebSocket Manager - 统一管理所有 WebSocket 连接
package websocket

import (
\t"sync"
\t"time"

\t"github.com/gorilla/websocket"
)

// ManagerConfig Manager 配置
type ManagerConfig struct {
\tMaxConnections    int           // 最大连接数
\tHeartbeatInterval time.Duration // 心跳间隔
\tWriteTimeout      time.Duration // 写入超时
\tReadTimeout       time.Duration // 读取超时
}

// Manager WebSocket 管理器
type Manager struct {
\tconfig ManagerConfig

\t// 所有活跃连接
\tconnections map[string]*Connection

\t// 房间管理
\troomManager *RoomManager

\t// 连接池
\tpool *ConnectionPool

\tmu sync.RWMutex
}

// NewManager 创建新的 Manager
func NewManager(config ManagerConfig) *Manager {
\treturn &Manager{
\t\tconfig:      config,
\t\tconnections: make(map[string]*Connection),
\t\troomManager: NewRoomManager(),
\t\tpool:        NewConnectionPool(config.MaxConnections),
\t}
}

// Start 启动 Manager
func (m *Manager) Start() {
\t// 启动心跳检查
\tgo m.startHeartbeatChecker()

\t// 启动连接池监控
\tgo m.pool.StartMonitoring()
}

// Shutdown 关闭 Manager
func (m *Manager) Shutdown() {
\tm.mu.Lock()
\tdefer m.mu.Unlock()

\t// 关闭所有连接
\tfor _, conn := range m.connections {
\t\tconn.Close()
\t}

\t// 关闭连接池
\tm.pool.Shutdown()
}

// AddConnection 添加新连接
func (m *Manager) AddConnection(connID string, wsConn *websocket.Conn) error {
\tm.mu.Lock()
\tdefer m.mu.Unlock()

\tif len(m.connections) >= m.config.MaxConnections {
\t\treturn ErrMaxConnectionsReached
\t}

\t// 从连接池获取 Connection 对象
\tconnection := m.pool.Get()
\tif connection == nil {
\t\tconnection = NewConnection(wsConn, m.config.WriteTimeout, m.config.ReadTimeout)
\t} else {
\t\tconnection.Reset(wsConn)
\t}

\tm.connections[connID] = connection

\t// 设置心跳
\tconnection.SetHeartbeat(
\t\tm.config.HeartbeatInterval,
\t\tm.config.HeartbeatInterval*3,
\t\tfunc() {
\t\t\tm.RemoveConnection(connID)
\t\t},
\t)

\treturn nil
}

// RemoveConnection 移除连接
func (m *Manager) RemoveConnection(connID string) {
\tm.mu.Lock()
\tdefer m.mu.Unlock()

\tif conn, exists := m.connections[connID]; exists {
\t\tconn.Close()
\t\tm.pool.Put(conn) // 回收到连接池
\t\tdelete(m.connections, connID)
\t}
}

// GetConnection 获取连接
func (m *Manager) GetConnection(connID string) (*Connection, bool) {
\tm.mu.RLock()
\tdefer m.mu.RUnlock()

\tconn, exists := m.connections[connID]
\treturn conn, exists
}

// SendToConnection 发送消息到指定连接
func (m *Manager) SendToConnection(connID string, message []byte) error {
\tconn, exists := m.GetConnection(connID)
\tif !exists {
\t\treturn ErrConnectionNotFound
\t}

\treturn conn.WriteMessage(message)
}

// BroadcastToAll 广播消息到所有连接
func (m *Manager) BroadcastToAll(message []byte) error {
\tm.mu.RLock()
\tdefer m.mu.RUnlock()

\tvar errors []error
\tfor _, conn := range m.connections {
\t\tif err := conn.WriteMessage(message); err != nil {
\t\t\terrors = append(errors, err)
\t\t}
\t}

\tif len(errors) > 0 {
\t\treturn errors[0]
\t}
\treturn nil
}

// JoinRoom 加入房间
func (m *Manager) JoinRoom(connID, roomID string) error {
\tconn, exists := m.GetConnection(connID)
\tif !exists {
\t\treturn ErrConnectionNotFound
\t}

\treturn m.roomManager.JoinRoom(roomID, conn)
}

// LeaveRoom 离开房间
func (m *Manager) LeaveRoom(connID, roomID string) error {
\tconn, exists := m.GetConnection(connID)
\tif !exists {
\t\treturn ErrConnectionNotFound
\t}

\treturn m.roomManager.LeaveRoom(roomID, conn)
}

// SendToRoom 发送消息到房间
func (m *Manager) SendToRoom(roomID string, message []byte) error {
\treturn m.roomManager.Broadcast(roomID, message)
}

// startHeartbeatChecker 启动心跳检查
func (m *Manager) startHeartbeatChecker() {
\tticker := time.NewTicker(m.config.HeartbeatInterval)
\tdefer ticker.Stop()

\tfor range ticker.C {
\t\tm.mu.RLock()
\t\tfor connID, conn := range m.connections {
\t\t\tif !conn.IsAlive() {
\t\t\t\tgo m.RemoveConnection(connID)
\t\t\t}
\t\t}
\t\tm.mu.RUnlock()
\t}
}

// GetStats 获取统计信息
func (m *Manager) GetStats() Stats {
\tm.mu.RLock()
\tdefer m.mu.RUnlock()

\treturn Stats{
\t\tTotalConnections: len(m.connections),
\t\tTotalRooms:       m.roomManager.RoomCount(),
\t\tPoolSize:         m.pool.Size(),
\t}
}

// Stats 统计信息
type Stats struct {
\tTotalConnections int `json:"total_connections"`
\tTotalRooms       int `json:"total_rooms"`
\tPoolSize         int `json:"pool_size"`
}
"""
    
    with open(base_path / 'manager.go', 'w', encoding='utf-8') as f:
        f.write(manager_code)
    
    print(f"✓ 生成 manager.go")
    
    # connection.go
    connection_code = """// Connection - WebSocket 连接封装
package websocket

import (
\t"sync"
\t"time"

\t"github.com/gorilla/websocket"
)

// Connection WebSocket 连接封装
type Connection struct {
\tconn *websocket.Conn

\twriteTimeout time.Duration
\treadTimeout  time.Duration

\t// 心跳相关
\theartbeatInterval time.Duration
\theartbeatTimeout  time.Duration
\tlastActiveTime    time.Time
\theartbeatStop     chan struct{}

\tmu sync.RWMutex
}

// NewConnection 创建新的 Connection
func NewConnection(wsConn *websocket.Conn, writeTimeout, readTimeout time.Duration) *Connection {
\treturn &Connection{
\t\tconn:         wsConn,
\t\twriteTimeout: writeTimeout,
\t\treadTimeout:  readTimeout,
\t\tlastActiveTime: time.Now(),
\t\theartbeatStop:  make(chan struct{}),
\t}
}

// Reset 重置连接（用于连接池复用）
func (c *Connection) Reset(wsConn *websocket.Conn) {
\tc.mu.Lock()
\tdefer c.mu.Unlock()

\tc.conn = wsConn
\tc.lastActiveTime = time.Now()
\tc.heartbeatStop = make(chan struct{})
}

// Close 关闭连接
func (c *Connection) Close() error {
\t// 停止心跳
\tclose(c.heartbeatStop)

\treturn c.conn.Close()
}

// WriteMessage 写入消息
func (c *Connection) WriteMessage(message []byte) error {
\tc.mu.Lock()
\tdefer c.mu.Unlock()

\tc.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
\treturn c.conn.WriteMessage(websocket.TextMessage, message)
}

// ReadMessage 读取消息
func (c *Connection) ReadMessage() ([]byte, error) {
\tc.mu.Lock()
\tdefer c.mu.Unlock()

\tc.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
\t_, message, err := c.conn.ReadMessage()
\tif err == nil {
\t\tc.lastActiveTime = time.Now()
\t}
\treturn message, err
}

// SetHeartbeat 设置心跳
func (c *Connection) SetHeartbeat(interval, timeout time.Duration, onTimeout func()) {
\tc.mu.Lock()
\tdefer c.mu.Unlock()

\tc.heartbeatInterval = interval
\tc.heartbeatTimeout = timeout

\t// 启动心跳协程
\tgo func() {
\t\tticker := time.NewTicker(interval)
\t\tdefer ticker.Stop()

\t\tfor {
\t\t\tselect {
\t\t\tcase <-ticker.C:
\t\t\t\tif time.Since(c.lastActiveTime) > timeout {
\t\t\t\t\tonTimeout()
\t\t\t\t\treturn
\t\t\t\t}
\t\t\t\t// 发送 ping
\t\t\t\tc.sendPing()
\t\t\tcase <-c.heartbeatStop:
\t\t\t\treturn
\t\t\t}
\t\t}
\t}()
}

// sendPing 发送 Ping
func (c *Connection) sendPing() {
\tc.mu.Lock()
\tdefer c.mu.Unlock()

\tc.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
\tc.conn.WriteMessage(websocket.PingMessage, []byte{})
}

// IsAlive 检查连接是否存活
func (c *Connection) IsAlive() bool {
\tc.mu.RLock()
\tdefer c.mu.RUnlock()

\treturn time.Since(c.lastActiveTime) <= c.heartbeatTimeout
}

// ConnID 获取连接 ID（实现 stub）
func (c *Connection) ConnID() string {
\t// TODO: 实现连接 ID 生成逻辑
\treturn "unknown"
}
"""
    
    with open(base_path / 'connection.go', 'w', encoding='utf-8') as f:
        f.write(connection_code)
    
    print(f"✓ 生成 connection.go")
    
    # connection_pool.go
    pool_code = """// ConnectionPool - 连接对象池
package websocket

import (
\t"sync"
\t"time"
)

// ConnectionPool 连接对象池
type ConnectionPool struct {
\tmaxSize int
\tpool    chan *Connection

\tmu sync.Mutex
}

// NewConnectionPool 创建连接池
func NewConnectionPool(maxSize int) *ConnectionPool {
\treturn &ConnectionPool{
\t\tmaxSize: maxSize,
\t\tpool:    make(chan *Connection, maxSize),
\t}
}

// Get 从池中获取连接
func (p *ConnectionPool) Get() *Connection {
\tselect {
\tcase conn := <-p.pool:
\t\treturn conn
\tdefault:
\t\treturn nil
\t}
}

// Put 归还连接到池
func (p *ConnectionPool) Put(conn *Connection) {
\tselect {
\tcase p.pool <- conn:
\t\t// 成功归还
\tdefault:
\t\t// 池已满，丢弃
\t\tconn.Close()
\t}
}

// Size 获取池大小
func (p *ConnectionPool) Size() int {
\treturn len(p.pool)
}

// StartMonitoring 启动监控
func (p *ConnectionPool) StartMonitoring() {
\tticker := time.NewTicker(1 * time.Minute)
\tdefer ticker.Stop()

\tfor range ticker.C {
\t\tp.mu.Lock()
\t\tsize := p.Size()
\t\tp.mu.Unlock()

\t\t// 这里可以添加监控指标上报
\t\t_ = size
\t}
}

// Shutdown 关闭连接池
func (p *ConnectionPool) Shutdown() {
\tclose(p.pool)
}
"""
    
    with open(base_path / 'connection_pool.go', 'w', encoding='utf-8') as f:
        f.write(pool_code)
    
    print(f"✓ 生成 connection_pool.go")
    
    # room.go
    room_code = """// Room - WebSocket 房间
package websocket

import (
\t"sync"
)

// Room 房间定义
type Room struct {
\tID       string
\tName     string
\tMaxUsers int

\tconnections map[string]*Connection
\tmu          sync.RWMutex
}

// NewRoom 创建新房间
func NewRoom(id, name string, maxUsers int) *Room {
\treturn &Room{
\t\tID:          id,
\t\tName:        name,
\t\tMaxUsers:    maxUsers,
\t\tconnections: make(map[string]*Connection),
\t}
}

// Join 加入房间
func (r *Room) Join(conn *Connection) error {
\tr.mu.Lock()
\tdefer r.mu.Unlock()

\tif len(r.connections) >= r.MaxUsers {
\t\treturn ErrRoomFull
\t}

\tr.connections[conn.ConnID()] = conn
\treturn nil
}

// Leave 离开房间
func (r *Room) Leave(conn *Connection) {
\tr.mu.Lock()
\tdefer r.mu.Unlock()

\tdelete(r.connections, conn.ConnID())
}

// Broadcast 广播消息
func (r *Room) Broadcast(message []byte) error {
\tr.mu.RLock()
\tdefer r.mu.RUnlock()

\tvar errors []error
\tfor _, conn := range r.connections {
\t\tif err := conn.WriteMessage(message); err != nil {
\t\t\terrors = append(errors, err)
\t\t}
\t}

\tif len(errors) > 0 {
\t\treturn errors[0]
\t}
\treturn nil
}

// UserCount 获取用户数
func (r *Room) UserCount() int {
\tr.mu.RLock()
\tdefer r.mu.RUnlock()
\treturn len(r.connections)
}
"""
    
    with open(base_path / 'room.go', 'w', encoding='utf-8') as f:
        f.write(room_code)
    
    print(f"✓ 生成 room.go")
    
    # room_manager.go
    room_manager_code = """// RoomManager - 房间管理器
package websocket

import (
\t"sync"
)

// RoomManager 房间管理器
type RoomManager struct {
\trooms map[string]*Room
\tmu    sync.RWMutex
}

// NewRoomManager 创建房间管理器
func NewRoomManager() *RoomManager {
\treturn &RoomManager{
\t\trooms: make(map[string]*Room),
\t}
}

// CreateRoom 创建房间
func (m *RoomManager) CreateRoom(id, name string, maxUsers int) *Room {
\tm.mu.Lock()
\tdefer m.mu.Unlock()

\troom := NewRoom(id, name, maxUsers)
\tm.rooms[id] = room
\treturn room
}

// GetRoom 获取房间
func (m *RoomManager) GetRoom(id string) (*Room, bool) {
\tm.mu.RLock()
\tdefer m.mu.RUnlock()

\troom, exists := m.rooms[id]
\treturn room, exists
}

// DeleteRoom 删除房间
func (m *RoomManager) DeleteRoom(id string) {
\tm.mu.Lock()
\tdefer m.mu.Unlock()

\tif room, exists := m.rooms[id]; exists {
\t\t// 关闭房间内所有连接
\t\tfor _, conn := range room.connections {
\t\t\tconn.Close()
\t\t}
\t\tdelete(m.rooms, id)
\t}
}

// JoinRoom 加入房间
func (m *RoomManager) JoinRoom(id string, conn *Connection) error {
\troom, exists := m.GetRoom(id)
\tif !exists {
\t\treturn ErrRoomNotFound
\t}

\treturn room.Join(conn)
}

// LeaveRoom 离开房间
func (m *RoomManager) LeaveRoom(id string, conn *Connection) error {
\troom, exists := m.GetRoom(id)
\tif !exists {
\t\treturn ErrRoomNotFound
\t}

\troom.Leave(conn)
\treturn nil
}

// Broadcast 广播消息到房间
func (m *RoomManager) Broadcast(id string, message []byte) error {
\troom, exists := m.GetRoom(id)
\tif !exists {
\t\treturn ErrRoomNotFound
\t}

\treturn room.Broadcast(message)
}

// RoomCount 获取房间数量
func (m *RoomManager) RoomCount() int {
\tm.mu.RLock()
\tdefer m.mu.RUnlock()
\treturn len(m.rooms)
}
"""
    
    with open(base_path / 'room_manager.go', 'w', encoding='utf-8') as f:
        f.write(room_manager_code)
    
    print(f"✓ 生成 room_manager.go")
    
    # errors.go
    errors_code = """// Package websocket - WebSocket 错误定义
package websocket

import "errors"

// WebSocket 常见错误
var (
\tErrConnectionNotFound    = errors.New("connection not found")
\tErrMaxConnectionsReached = errors.New("max connections reached")
\tErrRoomNotFound          = errors.New("room not found")
\tErrRoomFull              = errors.New("room is full")
\tErrAlreadyInRoom         = errors.New("already in room")
)
"""
    
    with open(base_path / 'errors.go', 'w', encoding='utf-8') as f:
        f.write(errors_code)
    
    print(f"✓ 生成 errors.go")
    
    print(f"\n✅ WebSocket 基础设施生成完成!\n")
    print(f"下一步:")
    print(f"  go get github.com/gorilla/websocket")
    print(f"  查看生成的文件：{base_path}\n")


def main():
    """主函数"""
    args = parse_args()
    
    if args.command == 'init':
        create_websocket_manager(args.project, args.output)
    
    elif args.command == 'add':
        if args.action == 'room':
            print(f"🏠 添加房间：{args.name} (最大用户数：{args.max_users})")
            # TODO: 实现房间添加逻辑
        else:
            print("❌ 未知的添加类型")
            sys.exit(1)
    
    elif args.command == 'config':
        print("⚙️  配置功能开发中...")
    
    else:
        print("❌ 错误：未知命令。使用 --help 查看可用命令")
        sys.exit(1)


if __name__ == '__main__':
    main()
