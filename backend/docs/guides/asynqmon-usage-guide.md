# asynqmon 使用指南

## 什么是 asynqmon？

asynqmon 是 asynq 任务队列的 Web UI 监控工具，类似 Python Celery 的 Flower。

**功能**：
- ✅ 实时查看任务队列状态
- ✅ 监控任务执行情况（成功/失败）
- ✅ 查看任务详情和错误信息
- ✅ 手动重试失败任务
- ✅ 统计图表和趋势分析

---

## 🚀 快速开始（不使用 Docker）

### 1️⃣ 安装 asynqmon

```bash
cd /Users/shenfay/Projects/go-ddd-scaffold/backend

# 方法 1：使用 Makefile（推荐）
make asynqmon-install

# 方法 2：直接使用 go install
go install github.com/hibiken/asynqmon/cmd/asynqmon@latest
```

**说明**：asynqmon 现已独立为单独的仓库 `github.com/hibiken/asynqmon`

### 2️⃣ 确保 Redis 正在运行

```bash
# 检查 Redis 是否运行
redis-cli ping
# 应该返回 PONG

# 如果没有运行（macOS）
brew services start redis

# 或者手动启动
redis-server
```

### 3️⃣ 启动 asynqmon

```bash
# 方法 1：使用 Makefile（推荐）
make asynqmon

# 方法 2：直接运行
asynqmon --redis-addr=localhost:6379

# 方法 3：指定端口
make asynqmon-port PORT=8081
# 或
asynqmon --redis-addr=localhost:6379 --port=8081
```

### 4️⃣ 访问 UI

打开浏览器访问：
- **http://localhost:8080** （默认端口）
- **http://localhost:8081** （如果你指定了 --port=8081）

---

## 📋 Makefile 命令

项目提供了以下 asynqmon 相关命令：

```bash
# 安装 asynqmon CLI 工具
make asynqmon-install

# 启动 asynqmon UI（默认端口 8080）
make asynqmon

# 启动 asynqmon UI（自定义端口）
make asynqmon-port PORT=8081

# 在浏览器中打开 asynqmon UI（仅 macOS）
make asynqmon-ui
```

---

## 🔧 高级配置

### 连接远程 Redis

```bash
asynqmon --redis-addr=redis.example.com:6379 \
         --redis-password=your_password \
         --redis-db=0
```

### 启用 TLS 连接

```bash
asynqmon --redis-addr=redis.example.com:6379 \
         --redis-tls
```

### 查看所有可用选项

```bash
asynqmon --help
```

---

## 🎯 UI 功能介绍

### Dashboard（仪表盘）

显示：
- 各队列的任务数量（Pending, Active, Completed, Failed）
- 任务处理速率图表
- 最近处理的任务列表

### Queues（队列）

查看每个队列的详细信息：
- **critical** - 高优先级队列（权重 6）
- **default** - 默认队列（权重 3）
- **low** - 低优先级队列（权重 1）

### Task Detail（任务详情）

点击任意任务查看详情：
- 任务类型和 Payload
- 执行时间和状态
- 错误信息（如果失败）
- 重试次数

### 操作

- **Archive** - 归档已完成的任务
- **Delete** - 删除任务
- **Retry** - 重试失败的任务

---

## 💡 常见问题

### Q1: asynqmon 找不到命令？

**A:** 确保 GOPATH/bin 在 PATH 环境变量中：

```bash
# macOS/Linux
export PATH=$PATH:$(go env GOPATH)/bin

# 添加到 ~/.zshrc 或 ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Q2: 连接不到 Redis？

**A:** 检查 Redis 是否运行：

```bash
# 检查 Redis 状态
redis-cli ping

# 如果不是 PONG，启动 Redis
brew services start redis  # macOS
# 或
redis-server
```

### Q3: 端口被占用？

**A:** 使用其他端口：

```bash
make asynqmon-port PORT=9000
```

### Q4: 能看到 UI 但没有任务数据？

**A:** 确保：
1. 应用已经发布了一些领域事件
2. asynq Worker 正在运行
3. Redis 地址正确（应用和 asynqmon 使用同一个 Redis）

---

## 🔗 项目集成

### 在应用中发布事件

```go
// 在 Repository 或 Service 中
event := NewUserCreatedEvent(userID, email)
err = eventPublisher.Publish(ctx, event)
```

### 启动 asynq Worker

需要在应用启动时同时启动 asynq Worker：

```go
// 在 bootstrap 或 main.go 中
processor := task_queue.NewProcessor(logger, handlers...)
asynqServer := task_queue.NewServer(cfg.Asynq)
asynqServer.RegisterHandler(processor.ProcessTask)

// 后台运行
go asynqServer.Run()
```

---

## 📊 架构图

```
┌─────────────┐
│   API App   │ 发布领域事件
└──────┬──────┘
       ↓
┌─────────────┐
│    Redis    │ asynq 队列后端
└──────┬──────┘
       ├──────────────┐
       ↓              ↓
┌─────────────┐ ┌────────────┐
│ asynq Worker│ │  asynqmon  │ 监控 UI
└─────────────┘ └────────────┘
```

---

## 📚 参考资料

- [asynq 官方文档](https://github.com/hibiken/asynq)
- [asynqmon GitHub](https://github.com/hibiken/asynqmon)
- [项目架构文档](./docs/architecture/event-driven-architecture.md)

---

**最后更新**: 2026-03-19
