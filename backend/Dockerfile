# 多阶段构建 - 开发环境
FROM golang:1.25-alpine AS development

# 安装必要的工具
RUN apk add --no-cache git wget ca-certificates

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY backend/go.mod backend/go.sum ./backend/

# 下载依赖
RUN cd backend && go mod download

# 复制源代码
COPY backend/ ./backend/

# 设置工作目录到 backend
WORKDIR /app/backend

# 暴露端口
EXPOSE 8080

# 默认命令（在 docker-compose 中会被覆盖）
CMD ["go", "run", "./cmd/api/main.go"]

# 生产环境构建（预留）
FROM golang:1.25-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY backend/go.mod backend/go.sum ./backend/

# 下载依赖
RUN cd backend && go mod download

# 复制源代码
COPY backend/ ./backend/

# 编译 API 服务
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api ./cmd/api/main.go

# 编译 Worker 服务
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/worker ./cmd/worker/main.go

# 最终镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# 从 builder 复制二进制文件
COPY --from=builder --chown=appuser:appgroup /app/api /app/api
COPY --from=builder --chown=appuser:appgroup /app/worker /app/worker

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 默认命令
CMD ["/app/api"]
