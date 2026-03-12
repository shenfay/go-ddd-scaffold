#!/bin/bash

# 数据库迁移执行脚本
# 用于在开发/生产环境执行数据库迁移
# 使用方法：./tools/migrator/migrate.sh [up|down|version|force]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 从环境变量读取配置（优先使用 APP_ 前缀的变量）
DB_HOST="${APP_DATABASE_HOST:-${DB_HOST:-localhost}}"
DB_PORT="${APP_DATABASE_PORT:-${DB_PORT:-5432}}"
DB_NAME="${APP_DATABASE_NAME:-${DB_NAME:-go_ddd_scaffold}}"
DB_USER="${APP_DATABASE_USER:-${DB_USER:-shenfay}}"
DB_PASSWORD="${APP_DATABASE_PASSWORD:-${DB_PASSWORD:-postgres}}"
DB_SSL_MODE="${APP_DATABASE_SSL_MODE:-${DB_SSL_MODE:-disable}}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="${SCRIPT_DIR}/../../migrations"

# 显示帮助信息
show_help() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}数据库迁移工具${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    echo "用法：$0 [命令]"
    echo ""
    echo "命令:"
    echo "  up        应用所有待处理的迁移 (默认)"
    echo "  down      回滚最近一次迁移"
    echo "  version   查看当前版本"
    echo "  force N   强制设置版本为 N"
    echo "  help      显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 应用所有迁移"
    echo "  $0 up           # 应用所有迁移"
    echo "  $0 down         # 回滚一次"
    echo "  $0 version      # 查看版本"
    echo "  $0 force 10     # 强制设置为版本 10"
    echo ""
    echo -e "${YELLOW}配置信息：${NC}"
    echo "  主机：${DB_HOST}:${DB_PORT}"
    echo "  数据库：${DB_NAME}"
    echo "  用户：${DB_USER}"
    echo ""
}

# 检查 migrate 工具是否安装
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}错误：未找到 migrate 工具${NC}"
        echo "请先安装："
        echo "  macOS: brew install golang-migrate"
        echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && sudo mv migrate /usr/bin/migrate"
        echo "  Go: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
    echo -e "${GREEN}✓ migrate 工具已安装 ($(migrate -version))${NC}"
}

# 检查数据库连接
check_connection() {
    echo -e "${YELLOW}检查数据库连接...${NC}"
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 可以连接到 PostgreSQL${NC}"
    else
        echo -e "${RED}✗ 无法连接到 PostgreSQL${NC}"
        echo "请检查："
        echo "  1. PostgreSQL 服务是否运行"
        echo "  2. 数据库连接配置是否正确"
        echo "  3. 用户权限是否正确"
        exit 1
    fi
}

# 应用所有迁移
migrate_up() {
    echo -e "${YELLOW}开始应用迁移...${NC}"
    cd "$MIGRATIONS_DIR"
    if migrate -database "$DATABASE_URL" -path . up 2>&1; then
        echo -e "${GREEN}✓ 迁移应用成功${NC}"
    else
        echo -e "${RED}✗ 迁移失败${NC}"
        exit 1
    fi
}

# 回滚一次迁移
migrate_down() {
    echo -e "${YELLOW}回滚最近一次迁移...${NC}"
    cd "$MIGRATIONS_DIR"
    if migrate -database "$DATABASE_URL" -path . down 1 2>&1; then
        echo -e "${GREEN}✓ 回滚成功${NC}"
    else
        echo -e "${RED}✗ 回滚失败${NC}"
        exit 1
    fi
}

# 查看当前版本
migrate_version() {
    cd "$MIGRATIONS_DIR"
    CURRENT_VERSION=$(migrate -database "$DATABASE_URL" -path . version 2>&1)
    echo -e "${BLUE}当前数据库版本：${CURRENT_VERSION}${NC}"
}

# 强制设置版本
migrate_force() {
    local version=$1
    echo -e "${YELLOW}强制设置数据库版本为 ${version}...${NC}"
    cd "$MIGRATIONS_DIR"
    if migrate -database "$DATABASE_URL" -path . force $version 2>&1; then
        echo -e "${GREEN}✓ 版本已强制设置为 ${version}${NC}"
    else
        echo -e "${RED}✗ 设置版本失败${NC}"
        exit 1
    fi
}

# 主函数
main() {
    check_migrate
    check_connection
    
    echo ""
    
    case "${1:-up}" in
        up)
            migrate_up
            ;;
        down)
            migrate_down
            ;;
        version)
            migrate_version
            ;;
        force)
            if [ -z "$2" ]; then
                echo -e "${RED}错误：force 命令需要指定版本号${NC}"
                echo "用法：$0 force <version>"
                exit 1
            fi
            migrate_force "$2"
            ;;
        help|-h|--help)
            show_help
            ;;
        *)
            echo -e "${RED}未知命令：$1${NC}"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
