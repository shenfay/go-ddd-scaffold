#!/bin/bash

# 服务启动脚本
# 用法：./start-services.sh [backend|frontend|all]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "检查前置条件..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装，请先安装 Node.js"
        exit 1
    fi
    
    # 检查 MySQL
    if ! command -v mysql &> /dev/null; then
        log_warning "MySQL 客户端未找到，请确保 MySQL 服务正在运行"
    fi
    
    # 检查 Redis
    if ! command -v redis-cli &> /dev/null; then
        log_warning "Redis 客户端未找到，请确保 Redis 服务正在运行"
    fi
    
    log_success "前置检查完成"
}

start_backend() {
    log_info "启动后端服务..."
    cd "$SCRIPT_DIR/backend"
    
    if [ -f ".env" ]; then
        log_info "加载环境变量..."
        export $(cat .env | grep -v '^#' | xargs)
    fi
    
    log_info "后端服务启动中..."
    log_info "API 地址：http://localhost:8080"
    log_info "健康检查：http://localhost:8080/health"
    log_info "监控指标：http://localhost:8080/metrics"
    
    go run cmd/server/main.go
}

start_frontend() {
    log_info "启动前端服务..."
    cd "$SCRIPT_DIR/frontend"
    
    if [ ! -d "node_modules" ]; then
        log_warning "依赖未安装，正在安装依赖..."
        npm install
    fi
    
    log_info "前端服务启动中..."
    log_info "访问地址：http://localhost:3000"
    
    npm start
}

start_all() {
    log_info "同时启动前后端服务..."
    
    # 启动后端（后台运行）
    log_info "启动后端服务（后台运行）..."
    cd "$SCRIPT_DIR/backend"
    go run cmd/server/main.go > /tmp/backend.log 2>&1 &
    BACKEND_PID=$!
    log_success "后端服务已启动 (PID: $BACKEND_PID)"
    
    # 等待后端完全启动
    log_info "等待后端服务启动..."
    sleep 5
    
    # 检查后端是否正常启动
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "后端服务运行正常"
    else
        log_warning "后端服务可能还未完全启动，请稍后检查"
    fi
    
    # 启动前端（后台运行）
    log_info "启动前端服务（后台运行）..."
    cd "$SCRIPT_DIR/frontend"
    npm start > /tmp/frontend.log 2>&1 &
    FRONTEND_PID=$!
    log_success "前端服务已启动 (PID: $FRONTEND_PID)"
    
    echo ""
    log_success "所有服务已启动！"
    echo ""
    log_info "访问地址:"
    log_info "  前端应用：http://localhost:3000"
    log_info "  后端 API: http://localhost:8080"
    log_info "  健康检查：http://localhost:8080/health"
    log_info "  监控指标：http://localhost:8080/metrics"
    echo ""
    log_info "查看日志:"
    log_info "  后端日志：tail-f /tmp/backend.log"
    log_info "  前端日志：tail-f /tmp/frontend.log"
    echo ""
    log_info "停止服务:"
    log_info "  kill $BACKEND_PID  # 停止后端"
    log_info "  kill $FRONTEND_PID  # 停止前端"
}

show_help() {
    echo "用法：$0 [command]"
    echo ""
    echo "Commands:"
    echo "  backend   仅启动后端服务"
    echo "  frontend  仅启动前端服务"
    echo "  all       同时启动前后端服务（默认）"
    echo "  help      显示帮助信息"
    echo ""
    echo "Examples:"
    echo "  $0 backend    # 只启动后端"
    echo "  $0 frontend   # 只启动前端"
    echo "  $0 all        # 启动所有服务"
    echo "  $0            # 默认启动所有服务"
    echo ""
}

# 主程序
main() {
    case "${1:-all}" in
        backend)
            check_prerequisites
            start_backend
            ;;
        frontend)
            check_prerequisites
            start_frontend
            ;;
        all)
            check_prerequisites
            start_all
            ;;
        help|-h|--help)
            show_help
            ;;
        *)
            log_error "未知命令：$1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
