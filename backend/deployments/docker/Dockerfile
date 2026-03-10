# Backend CI/CD 配置

FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    ~/go/bin/swag init -g cmd/server/main.go -o docs

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 最终镜像
FROM alpine:latest

WORKDIR /root/

# 安装 CA 证书（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

# 从构建镜像复制二进制文件
COPY --from=builder/app/main .
COPY --from=builder/app/config/ ./config/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
