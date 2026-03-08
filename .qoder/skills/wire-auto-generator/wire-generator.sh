#!/bin/bash

# Wire Auto Generator Script
# Wire 依赖注入自动化生成工具

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
DEFAULT_WIRE_DIR="./backend/internal/infrastructure/wire"
VERBOSE=false
DRY_RUN=false
AUTO_FIX=true

# 打印函数
print_info() {
    echo -e "${BLUE}💡 $1${NC}"
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

# 使用帮助
show_help() {
    cat << EOF
Wire 依赖注入自动化生成工具

用法: $0 [选项] [目录路径]

选项:
  -h, --help        显示帮助信息
  -v, --verbose     详细输出模式
  -d, --dry-run     仅检查，不生成代码
  -f, --auto-fix    自动修复可修复的问题
  -o, --output-dir  指定 Wire 配置目录（默认：$DEFAULT_WIRE_DIR）

示例:
  $0                                    # 使用默认配置
  $0 -v                                 # 详细模式
  $0 -d                                 # 仅检查
  $0 ./backend/internal/infrastructure/wire
  $0 --output-dir ./backend/internal/infrastructure/wire

EOF
}

# 检查 Go 环境
check_go_environment() {
    print_info "检查 Go 环境..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    
    GO_VERSION=$(go version)
    print_success "Go 环境：$GO_VERSION"
}

# 检查 Wire 是否安装
check_wire_installed() {
    print_info "检查 Wire 工具..."
    
    if ! command -v wire &> /dev/null; then
        print_warning "Wire 未安装，将使用 go run 方式"
        WIRE_CMD="go run github.com/google/wire/cmd/wire@latest"
    else
        WIRE_VERSION=$(wire version 2>&1 | head -1)
        print_success "Wire 已安装：$WIRE_VERSION"
        WIRE_CMD="wire"
    fi
}

# 检查项目结构
check_project_structure() {
    print_info "检查项目结构..."
    
    if [ ! -f "go.mod" ]; then
        print_error "未找到 go.mod，请在项目根目录运行"
        exit 1
    fi
    
    print_success "项目结构正常"
}

# 分析 Wire 配置
analyze_wire_config() {
    local wire_dir=$1
    
    print_info "分析 Wire 配置：$wire_dir"
    
    if [ ! -d "$wire_dir" ]; then
        print_error "目录不存在：$wire_dir"
        exit 1
    fi
    
    # 检查 injector.go
    if [ ! -f "$wire_dir/injector.go" ]; then
        print_warning "未找到 injector.go"
        print_info "建议创建 injector.go 定义 Provider Set"
        
        if [ "$AUTO_FIX" = true ]; then
            create_injector_template "$wire_dir"
        fi
        return 1
    fi
    
    # 检查 Provider Set
    PROVIDER_SET_COUNT=$(grep -c "var.*Set = wire.NewSet" "$wire_dir/injector.go" || echo "0")
    print_success "发现 $PROVIDER_SET_COUNT 个 Provider Set"
    
    # 检查 Injector 函数
    INJECTOR_COUNT=$(grep -c "^func Initialize" "$wire_dir/injector.go" || echo "0")
    print_success "发现 $INJECTOR_COUNT 个 Injector 函数"
    
    # 检查类型导出
    if grep -q "wire.Bind" "$wire_dir/injector.go"; then
        print_info "检测到 wire.Bind 使用"
        
        # 检查是否有未导出的类型
        UNEXPORTED_TYPES=$(grep "wire.Bind.*new(\*" "$wire_dir/injector.go" | grep -v "new([A-Z]" || echo "")
        if [ -n "$UNEXPORTED_TYPES" ]; then
            print_error "发现未导出的类型用于 wire.Bind:"
            echo "$UNEXPORTED_TYPES"
            
            if [ "$AUTO_FIX" = true ]; then
                print_info "建议：移除 wire.Bind 或重命名为导出类型"
            fi
        fi
    fi
    
    return 0
}

# 创建 injector.go 模板
create_injector_template() {
    local wire_dir=$1
    local template_file="$wire_dir/injector.go.template"
    
    print_info "创建 injector.go 模板..."
    
    cat > "$template_file" << 'TEMPLATE'
//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	
	// 导入你的服务包
	// "your-project/internal/application/service"
	// "your-project/internal/infrastructure/persistence/gorm/repo"
)

// ==================== 核心 Provider Set ====================

// RepositorySet 仓储集合
var RepositorySet = wire.NewSet(
	// repo.NewUserDAORepository,
	// repo.NewTenantDAORepository,
)

// TransactionSet 事务管理集合
var TransactionSet = wire.NewSet(
	// transaction.NewGormUnitOfWork,
)

// ==================== 应用服务初始化 ====================

// InitializeUserService 初始化用户服务
func InitializeUserService(
	db *gorm.DB,
	// 其他依赖...
) interface{} { // 替换为实际的服务接口
	wire.Build(
		RepositorySet,
		// service.NewUserService,
	)
	return nil
}
TEMPLATE

    print_success "模板已创建：$template_file"
    print_info "请参考模板创建 injector.go"
}

# 生成 Wire 代码
generate_wire_code() {
    local wire_dir=$1
    
    if [ "$DRY_RUN" = true ]; then
        print_warning "跳过代码生成（--dry-run 模式）"
        return 0
    fi
    
    print_info "开始生成 Wire 代码..."
    
    cd "$wire_dir" || exit 1
    
    # 运行 Wire 生成
    if $WIRE_CMD gen . 2>&1 | tee /tmp/wire_output.log; then
        print_success "Wire 代码生成成功"
        
        # 检查生成的文件
        if [ -f "wire_gen.go" ]; then
            LINES=$(wc -l < wire_gen.go)
            print_success "生成 wire_gen.go ($LINES 行)"
        fi
    else
        print_error "Wire 生成失败"
        
        # 显示错误信息
        if [ -f /tmp/wire_output.log ]; then
            print_info "错误详情:"
            cat /tmp/wire_output.log
        fi
        
        return 1
    fi
    
    cd - > /dev/null || exit 1
    return 0
}

# 验证编译
verify_compilation() {
    local wire_dir=$1
    
    print_info "验证编译..."
    
    cd "$wire_dir" || exit 1
    
    if go build . 2>&1 | tee /tmp/build_output.log; then
        print_success "编译通过"
    else
        print_error "编译失败"
        
        # 诊断错误
        diagnose_compilation_errors "/tmp/build_output.log"
        
        if [ "$AUTO_FIX" = true ]; then
            attempt_auto_fix "/tmp/build_output.log"
        fi
        
        return 1
    fi
    
    cd - > /dev/null || exit 1
    return 0
}

# 诊断编译错误
diagnose_compilation_errors() {
    local error_file=$1
    
    print_info "诊断编译错误..."
    
    if [ ! -f "$error_file" ]; then
        return 0
    fi
    
    # 常见错误模式匹配
    if grep -q "undefined:" "$error_file"; then
        print_error "发现未定义的类型或函数"
        print_info "可能原因：缺少 import 语句或类型名错误"
    fi
    
    if grep -q "redeclared in this block" "$error_file"; then
        print_error "发现重复定义"
        print_info "解决方案：删除重复的定义，只保留一个"
    fi
    
    if grep -q "name.*not exported by package" "$error_file"; then
        print_error "发现未导出的类型"
        print_info "解决方案：将类型名首字母大写或使用接口类型"
    fi
    
    if grep -q "import cycle not allowed" "$error_file"; then
        print_error "发现循环依赖"
        print_info "解决方案：重构代码结构，引入中间接口层"
    fi
}

# 尝试自动修复
attempt_auto_fix() {
    local error_file=$1
    
    print_info "尝试自动修复..."
    
    # 这里可以添加更多的自动修复逻辑
    # 目前仅提供指导性建议
    
    if grep -q "undefined:.*transaction" "$error_file"; then
        print_info "建议：在 injector.go 中添加 import"
        echo '    "go-ddd-scaffold/internal/infrastructure/transaction"'
    fi
    
    if grep -q "redeclared.*userEventBusAdapter" "$error_file"; then
        print_info "建议：删除 user.go 中的 userEventBusAdapter 定义（第 63-77 行）"
    fi
}

# 主函数
main() {
    local wire_dir="$DEFAULT_WIRE_DIR"
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -f|--auto-fix)
                AUTO_FIX=true
                shift
                ;;
            -o|--output-dir)
                wire_dir="$2"
                shift 2
                ;;
            *)
                if [ -d "$1" ]; then
                    wire_dir="$1"
                else
                    print_error "未知参数：$1"
                    show_help
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # 执行流程
    print_info "🚀 Wire 自动化生成工具启动"
    echo ""
    
    check_go_environment
    check_wire_installed
    check_project_structure
    
    echo ""
    
    if analyze_wire_config "$wire_dir"; then
        echo ""
        
        if generate_wire_code "$wire_dir"; then
            echo ""
            
            if verify_compilation "$wire_dir"; then
                echo ""
                print_success "🎉 Wire 代码生成完成！"
                echo ""
                
                # 显示下一步建议
                print_info "下一步建议:"
                echo "  1. 在 main.go 中使用生成的 Injector 函数"
                echo "  2. 替换手动依赖初始化代码"
                echo "  3. 运行集成测试验证"
                echo ""
            fi
        fi
    fi
}

# 执行主函数
main "$@"
