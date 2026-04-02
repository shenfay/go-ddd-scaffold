# Docker Compose 开发环境指南

**日期**: 2026-04-02  
**阶段**: Phase 9（高优先级）  
**状态**: ✅ 完成  

---

## 📋 概述

本项目提供了完整的 **Docker Compose 开发环境配置**，支持一键启动所有服务（PostgreSQL、Redis、API、Worker）。

### **核心特性**

1. ✅ **docker-compose.yml** - 服务编排配置
   - PostgreSQL 15 数据库
   - Redis 7 缓存/消息队列
   - API 服务（热重载）
   - Worker 服务（热重载）

2. ✅ **backend/Dockerfile** - 多阶段构建
   - 开发环境（源代码运行）
   - 生产环境（编译优化）

3. ✅ **backend/.dockerignore** - 构建优化
   - 排除不必要的文件
   - 加速镜像构建

---

## 🚀 快速启动

### **前置条件**

确保已安装：
- Docker Desktop（Mac/Windows）或 Docker Engine + Docker Compose（Linux）
- Git

### **一键启动**

```bash
# 1. 克隆项目
git clone https://github.com/shenfay/go-ddd-scaffold.git
cd go-ddd-scaffold

# 2. 启动所有服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f

# 4. 访问服务
# API: http://localhost:8080
# PostgreSQL: localhost:5432
# Redis: localhost:6379
```

---

## 🔧 常用命令

### **服务管理**

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启所有服务
docker-compose restart

# 重启单个服务
docker-compose restart api
docker-compose restart worker
docker-compose restart postgres
docker-compose restart redis

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f           # 所有服务
docker-compose logs -f api       # 只看 API
docker-compose logs -f worker    # 只看 Worker

# 进入容器
docker-compose exec api sh
docker-compose exec worker sh
docker-compose exec postgres psql -U postgres
docker-compose exec redis redis-cli
```

### **构建管理**

```bash
# 重新构建镜像
docker-compose build

# 强制重新构建（不使用缓存）
docker-compose build --no-cache

# 构建并启动
docker-compose up -d --build
```

### **清理数据**

```bash
# 停止服务并删除容器（保留数据卷）
docker-compose down

# 停止服务并删除容器和数据卷（⚠️ 会丢失数据）
docker-compose down -v

# 只删除数据卷
docker volume rm go-ddd-scaffold_postgres_data
docker volume rm go-ddd-scaffold_redis_data
```

---

## 📊 服务详情

### **1. PostgreSQL 数据库**

```yaml
service: postgres
image: postgres:15-alpine
port: 5432
database: go_ddd_scaffold
user: postgres
password: postgres
```

**连接信息**:
- **Host**: `localhost` (宿主机) / `postgres` (容器内)
- **Port**: `5432`
- **Database**: `go_ddd_scaffold`
- **Username**: `postgres`
- **Password**: `postgres`

**健康检查**:
```bash
# 检查数据库是否就绪
docker-compose exec postgres pg_isready -U postgres

# 查看数据库
docker-compose exec postgres psql -U postgres -c "\l"
```

---

### **2. Redis 缓存/消息队列**

```yaml
service: redis
image: redis:7-alpine
port: 6379
```

**连接信息**:
- **Host**: `localhost` (宿主机) / `redis` (容器内)
- **Port**: `6379`
- **Password**: 无（开发环境）

**健康检查**:
```bash
# 检查 Redis 是否就绪
docker-compose exec redis redis-cli ping

# 查看 Redis 信息
docker-compose exec redis redis-cli INFO
```

---

### **3. API 服务**

```yaml
service: api
port: 8080
health: /health/live
hot-reload: ✅
```

**访问地址**:
- **API**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Liveness**: http://localhost:8080/health/live
- **Readiness**: http://localhost:8080/health/ready

**环境变量**:
```bash
ENV=development
DB_HOST=postgres
DB_PORT=5432
DB_NAME=go_ddd_scaffold
DB_USER=postgres
DB_PASSWORD=postgres
REDIS_ADDR=redis:6379
JWT_SECRET=dev-secret-key-not-for-production-use-long-random-string-in-prod
SERVER_PORT=8080
SERVER_MODE=debug
```

**热重载**:
```bash
# 修改代码后自动重载（Go 1.25+ 支持）
# 或者手动重启
docker-compose restart api
```

---

### **4. Worker 服务**

```yaml
service: worker
health: 基于 Redis 连接
hot-reload: ✅
```

**环境变量**:
```bash
ENV=development
REDIS_ADDR=redis:6379
ASYNQ_CONCURRENCY=10
ASYNQ_QUEUES=critical:6,default:3,low:1
```

**查看任务队列**:
```bash
# 进入 Worker 容器
docker-compose exec worker sh

# 使用 asynq CLI 查看队列状态（需要安装）
asynq ls
```

---

## 💡 开发工作流

### **1. 本地开发（推荐）**

```bash
# 在宿主机运行 Go，使用 Docker 的 DB 和 Redis
export DB_HOST=localhost
export DB_PORT=5432
export REDIS_ADDR=localhost:6379

# 启动基础设施
docker-compose up -d postgres redis

# 本地运行 API
cd backend
go run ./cmd/api/main.go

# 本地运行 Worker
cd backend
go run ./cmd/worker/main.go
```

**优势**:
- ✅ 最快的编译速度
- ✅ 完整的 IDE 支持
- ✅ 方便的调试
- ✅ 使用容器的数据存储服务

---

### **2. 容器内开发**

```bash
# 进入 API 容器
docker-compose exec api sh

# 在容器内运行
go run ./cmd/api/main.go

# 实时重载（需要 air 等工具）
air -c .air.toml
```

**优势**:
- ✅ 环境一致性
- ✅ 无需本地安装 Go
- ✅ 接近生产环境

---

### **3. 完整容器开发**

```bash
# 修改 docker-compose.yml 中的 volumes 挂载
# 代码变更会自动触发重载

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api worker
```

**优势**:
- ✅ 完全隔离
- ✅ 一键启动
- ✅ 团队协作方便

---

## 🔍 故障排查

### **问题 1: 服务启动失败**

```bash
# 查看详细日志
docker-compose logs postgres
docker-compose logs redis
docker-compose logs api
docker-compose logs worker

# 检查端口占用
lsof -i :5432
lsof -i :6379
lsof -i :8080

# 解决：停止占用端口的服务或修改 docker-compose.yml 端口映射
```

---

### **问题 2: 数据库连接失败**

```bash
# 检查数据库是否就绪
docker-compose exec postgres pg_isready -U postgres

# 如果未就绪，等待几秒后重试
sleep 5
docker-compose exec postgres pg_isready -U postgres

# 查看数据库日志
docker-compose logs postgres

# 重置数据库
docker-compose down -v
docker-compose up -d postgres
```

---

### **问题 3: API 无法连接 Redis**

```bash
# 检查 Redis 是否就绪
docker-compose exec redis redis-cli ping

# 应该返回 PONG
# 如果返回错误，重启 Redis
docker-compose restart redis

# 检查网络连接
docker-compose exec api ping redis
```

---

### **问题 4: 代码修改不生效**

```bash
# 检查 volumes 是否正确挂载
docker-compose config

# 应该看到:
# volumes:
#   - ..:/app

# 重启服务
docker-compose restart api
docker-compose restart worker

# 或者重新构建
docker-compose up -d --build
```

---

## 🎯 最佳实践

### **1. 环境变量管理**

```bash
# 创建 .env 文件（不要提交到 Git）
cp .env.example .env

# .env 内容
DB_PASSWORD=your-password
JWT_SECRET=your-secret-key
ENV=development
```

```yaml
# docker-compose.yml 中使用
services:
  api:
    environment:
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
```

---

### **2. 数据持久化**

```yaml
# 使用命名卷（推荐）
volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

services:
  postgres:
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    volumes:
      - redis_data:/data
```

**优势**:
- ✅ 数据不随容器删除而丢失
- ✅ 便于备份和迁移
- ✅ 性能优于 bind mount

---

### **3. 健康检查依赖**

```yaml
services:
  api:
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
```

**优势**:
- ✅ 确保依赖服务完全就绪
- ✅ 避免启动顺序问题
- ✅ 提高系统可靠性

---

### **4. 资源限制**

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

**用途**:
- 防止单个服务占用过多资源
- 保证系统稳定性
- 适合多服务部署

---

## 📝 生产环境配置（预留）

### **生产 Dockerfile**

```dockerfile
# 使用多阶段构建减小镜像大小
FROM golang:1.25-alpine AS builder
# ... 编译代码 ...

FROM alpine:latest
# ... 复制二进制文件 ...
CMD ["/app/api"]
```

**优势**:
- ✅ 镜像大小 < 20MB
- ✅ 不包含源代码
- ✅ 非 root 用户运行

---

### **生产 docker-compose.prod.yml**

```yaml
version: '3.8'

services:
  api:
    build:
      context: ..
      dockerfile: backend/Dockerfile
      target: production  # 使用生产阶段
    environment:
      ENV: production
      JWT_SECRET: ${JWT_SECRET}  # 从环境变量读取
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 1G
```

**使用方式**:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## 🎉 总结

通过 Docker Compose，我们实现了：

✅ **一键启动** - 所有服务一键启动  
✅ **环境一致** - 开发、测试、生产环境一致  
✅ **易于协作** - 新成员快速上手  
✅ **热重载** - 代码修改即时生效  
✅ **健康检查** - 自动检测服务状态  
✅ **数据持久化** - 数据不丢失  
✅ **资源隔离** - 服务互不干扰  

**这是提升开发效率的关键工具！** 🚀

---

## 📞 参考文档

- [Docker Compose 官方文档](https://docs.docker.com/compose/)
- [Docker 多阶段构建](https://docs.docker.com/build/building/multi-stage/)
- [QUICKSTART.md](QUICKSTART.md) - 快速启动指南
- [ARCHITECTURE_SUMMARY.md](ARCHITECTURE_SUMMARY.md) - 整体架构说明
