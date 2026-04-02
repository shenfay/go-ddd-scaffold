#!/bin/bash

# 一键启动开发环境和监控工具
# 
# 用法:
#   ./start-dev.sh [选项]
#
# 选项:
#   -m, --monitor      同时启动监控工具（Asynqmon）
#   -s, --swagger      同时启动 Swagger UI
#   -a, --all          启动所有服务（默认）
#   -l, --local        使用本地模式（不使用 Docker）
#   -h, --help         显示帮助信息

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

show_help() {
    echo "用法：$0 [选项]"
    echo ""
    echo "启动开发环境和监控工具"
    echo ""
    echo "选项:"
    echo "  -m, --monitor      同时启动监控工具（Asynqmon）"
    echo "  -s, --swagger      同时启动 Swagger UI"
    echo "  -a, --all          启动所有服务（默认）"
    echo "  -l, --local        使用本地模式（不使用 Docker）"
    echo "  -h, --help         显示帮助信息"
    echo ""
    echo "示例:"
    echo "  ./start-dev.sh                    # 使用 Docker 启动基础服务"
    echo "  ./start-dev.sh --local            # 使用本地 PostgreSQL 和 Redis"
    echo "  ./start-dev.sh --local --monitor  # 本地模式 + Asynqmon"
    echo "  ./start-dev.sh --all              # 启动所有服务"
    echo ""
}

# 检查是否在 backend 目录
check_backend_dir() {
    if [ ! -f "backend/cmd/api/main.go" ]; then
        print_error "请在项目根目录运行此脚本"
        exit 1
    fi
}

# 检查 Docker Compose
check_docker_compose() {
    if command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE_CMD="docker-compose"
    elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
        DOCKER_COMPOSE_CMD="docker compose"
    else
        print_error "Docker Compose 未安装"
        exit 1
    fi
}

# 启动基础设施
start_infrastructure() {
    print_info "启动基础设施（PostgreSQL, Redis）..."
    $DOCKER_COMPOSE_CMD up -d postgres redis
    
    print_info "等待服务就绪..."
    sleep 5
    
    # 检查 PostgreSQL
    if docker-compose ps postgres | grep -q "healthy"; then
        print_success "PostgreSQL 已就绪"
    else
        print_warning "PostgreSQL 可能还未完全就绪"
    fi
    
    # 检查 Redis
    if docker-compose ps redis | grep -q "healthy"; then
        print_success "Redis 已就绪"
    else
        print_warning "Redis 可能还未完全就绪"
    fi
}

# 启动 API 和 Worker
start_services() {
    print_info "启动 API 和 Worker 服务..."
    
    # 在后台启动 API
    cd backend
    make run > /tmp/api.log 2>&1 &
    API_PID=$!
    cd ..
    
    print_info "API 进程已启动 (PID: $API_PID)"
    
    # 等待 API 启动
    sleep 3
    
    # 检查 API 是否成功启动
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "API 服务已启动"
    else
        print_warning "API 服务可能还在启动中..."
    fi
}

# 启动 Asynqmon
start_asynqmon() {
    print_info "启动 Asynqmon 监控..."
    
    # 检查是否已安装
    if ! command -v asynqmon &> /dev/null; then
        print_warning "asynqmon 未安装，正在安装..."
        make asynqmon-install
    fi
    
    # 在后台启动
    asynqmon --redis-addr=localhost:6379 > /tmp/asynqmon.log 2>&1 &
    ASYNQMON_PID=$!
    
    print_success "Asynqmon 已启动 (PID: $ASYNQMON_PID)"
    print_info "访问地址：http://localhost:8080"
}

# 启动 Swagger
start_swagger() {
    print_info "启动 Swagger UI..."
    
    cd backend
    
    # 生成文档
    print_info "生成 Swagger 文档..."
    make swagger-gen > /dev/null 2>&1
    
    # 在后台启动 docs server
    go run ./cmd/docs/main.go > /tmp/swagger.log 2>&1 &
    SWAGGER_PID=$!
    
    cd ..
    
    print_success "Swagger UI 已启动 (PID: $SWAGGER_PID)"
    print_info "访问地址：http://localhost:8080/swagger/index.html"
}

# 主函数
main() {
    START_MONITOR=false
    START_SWAGGER=false
    USE_LOCAL=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--monitor)
                START_MONITOR=true
                shift
                ;;
            -s|--swagger)
                START_SWAGGER=true
                shift
                ;;
            -a|--all)
                START_MONITOR=true
                START_SWAGGER=true
                shift
                ;;
            -l|--local)
                USE_LOCAL=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知选项：$1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 根据模式启动
    if [ "$USE_LOCAL" = true ]; then
        print_info "使用本地模式（不使用 Docker）"
        echo ""
        # 调用本地启动脚本
        SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
        LOCAL_SCRIPT="$SCRIPT_DIR/start-local.sh"
        
        if [ ! -f "$LOCAL_SCRIPT" ]; then
            print_error "找不到本地启动脚本：$LOCAL_SCRIPT"
            exit 1
        fi
        
        # 构建参数
        local args=()
        [ "$START_MONITOR" = true ] && args+=("--monitor")
        [ "$START_SWAGGER" = true ] && args+=("--swagger")
        [ "$START_MONITOR" = true ] && [ "$START_SWAGGER" = true ] && args+=("--all")
        
        exec "$LOCAL_SCRIPT" "${args[@]}"
    else
        # Docker 模式（原有逻辑）
        main_docker "$START_MONITOR" "$START_SWAGGER"
    fi
}
    
# Docker 模式主函数
main_docker() {
    local start_monitor=$1
    local start_swagger=$2
    
    # 检查
    check_backend_dir
    check_docker_compose
    
    # 启动步骤
    start_infrastructure
    start_services
    
    if [ "$START_MONITOR" = true ]; then
        start_asynqmon
    fi
    
    if [ "$START_SWAGGER" = true ]; then
        start_swagger
    fi
    
    echo ""
    echo "========================================"
    echo -e "${GREEN}✅ Docker 环境已启动${NC}"
    echo "========================================"
    echo ""
    print_info "服务状态:"
    echo ""
    echo "  📡 API 服务:"
    echo "     - 健康检查：http://localhost:8080/health"
    echo "     - Swagger:   http://localhost:8080/swagger/index.html"
    echo ""
    
    if [ "$START_MONITOR" = true ]; then
        echo "  📊 Asynqmon:"
        echo "     - 监控面板：http://localhost:8080"
        echo ""
    fi
    
    if [ "$START_SWAGGER" = true ]; then
        echo "  📖 Swagger Docs:"
        echo "     - API 文档：http://localhost:8080/swagger/index.html"
        echo ""
    fi
    
    echo "  💾 PostgreSQL:"
    echo "     - 连接：localhost:5432"
    echo ""
    
    echo "  🗄️  Redis:"
    echo "     - 连接：localhost:6379"
    echo ""
    
    echo "========================================"
    echo ""
    print_info "提示:"
    echo "  - 查看 API 日志：tail -f /tmp/api.log"
    echo "  - 查看 Asynqmon 日志：tail -f /tmp/asynqmon.log"
    echo "  - 执行测试流程：./scripts/dev/core-flow-test.sh"
    echo "  - 停止所有服务：docker-compose down"
    echo ""
    print_success "开发环境准备就绪！"
    echo ""
}

# Docker 模式主函数
main_docker() {
    local start_monitor=$1
    local start_swagger=$2
    
    echo ""
    echo "========================================"
    echo -e "${GREEN}🚀 启动 Docker 开发环境${NC}"
    echo "========================================"
    echo ""

# 执行主函数
main "$@"
