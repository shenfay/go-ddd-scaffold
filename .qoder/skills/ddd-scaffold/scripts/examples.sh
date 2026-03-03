#!/bin/bash
# DDD Scaffold 使用示例脚本

set -e

echo "======================================"
echo "DDD Scaffold 使用示例"
echo "======================================"
echo ""

# 示例 1: 快速生成（最小配置）
echo "📦 示例 1: 快速生成博客系统"
echo "--------------------------------------"
cat << 'EOF'
命令:
/ddd-scaffold --project-name blog-system --domains user,post,comment

生成的项目:
blog-system/
├── cmd/server/main.go
├── internal/domain/
│   ├── user/
│   ├── post/
│   └── comment/
├── internal/application/
├── internal/infrastructure/
└── internal/interfaces/
EOF
echo ""
echo ""

# 示例 2: 标准电商项目
echo "📦 示例 2: 标准电商项目"
echo "--------------------------------------"
cat << 'EOF'
命令:
/ddd-scaffold \\
  --project-name ecommerce \\
  --domains user,product,order,inventory,payment \\
  --style standard \\
  --with-examples \\
  --output ./ecommerce-app

特点:
- 5 个核心领域
- 标准四层架构
- 包含完整示例代码
- NATS 事件驱动
- GORM Repository 实现
EOF
echo ""
echo ""

# 示例 3: 完整 SaaS 平台
echo "📦 示例 3: 完整 SaaS 平台（多租户）"
echo "--------------------------------------"
cat << 'EOF'
命令:
/ddd-scaffold \\
  --project-name saas-platform \\
  --domains tenant,user,subscription,billing \\
  --style full \\
  --with-examples \\
  --with-tests \\
  --with-docker \\
  --with-monitoring

特点:
- 多租户架构支持
- Casbin RBAC 权限管理
- 完整的 Docker Compose 配置
- Prometheus + Grafana 监控
- 单元测试覆盖率 > 80%
EOF
echo ""
echo ""

# 示例 4: 交互模式
echo "📦 示例 4: 交互模式（推荐新手）"
echo "--------------------------------------"
cat << 'EOF'
命令:
/ddd-scaffold --interactive

交互流程:
1. 输入项目名称
2. 选择领域列表
3. 选择架构风格 (minimal/standard/full)
4. 确认是否包含示例代码
5. 确认是否包含测试文件
6. 确认是否包含 Docker 配置
7. 预览项目结构
8. 开始生成
EOF
echo ""
echo ""

# 后续步骤
echo "🚀 生成完成后的操作"
echo "--------------------------------------"
cat << 'EOF'
# 1. 进入项目目录
cd your-project-name

# 2. 下载依赖
go mod tidy

# 3. 编译项目
make build

# 4. 启动开发环境
make docker-up

# 5. 运行数据库迁移
make migrate-up

# 6. 启动应用
make run

# 7. 访问 API
curl http://localhost:8080/api/v1/users
EOF
echo ""
echo ""

echo "======================================"
echo "💡 提示"
echo "======================================"
cat << 'EOF'
- 首次使用建议从交互模式开始
- 查看 EXAMPLES.md 了解更多使用场景
- 参考 REFERENCE.md 获取详细文档
- 遇到问题可咨询 DDD Architect Agent
EOF
echo ""
