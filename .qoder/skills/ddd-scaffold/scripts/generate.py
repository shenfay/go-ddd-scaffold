#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DDD Scaffold - 主脚本
生成标准化的 DDD Clean Architecture 项目结构
"""

import os
import sys
import argparse
import yaml
from pathlib import Path
from datetime import datetime


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description='DDD Clean Architecture 项目脚手架生成器',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 快速生成
  %(prog)s --project-name myapp --domains user,order,product
  
  # 交互模式
  %(prog)s --interactive
  
  # 完整配置
  %(prog)s --project-name ecommerce --domains user,order,product,inventory \\
           --style full --with-examples --with-tests --with-docker
        """
    )
    
    parser.add_argument(
        '--project-name',
        type=str,
        help='项目名称 (必需)'
    )
    
    parser.add_argument(
        '--domains',
        type=str,
        default='user',
        help='领域列表，逗号分隔 (默认：user)'
    )
    
    parser.add_argument(
        '--output',
        type=str,
        default='./generated',
        help='输出目录 (默认：./generated)'
    )
    
    parser.add_argument(
        '--style',
        type=str,
        choices=['minimal', 'standard', 'full'],
        default='standard',
        help='架构风格 (默认：standard)'
    )
    
    parser.add_argument(
        '--with-examples',
        action='store_true',
        help='包含示例代码'
    )
    
    parser.add_argument(
        '--with-tests',
        action='store_true',
        help='包含测试文件'
    )
    
    parser.add_argument(
        '--with-docker',
        action='store_true',
        help='包含 Docker 配置'
    )
    
    parser.add_argument(
        '--interactive',
        action='store_true',
        help='交互模式'
    )
    
    return parser.parse_args()


def interactive_mode():
    """交互模式收集用户输入"""
    print("\n🏗️  DDD Scaffold - 交互式项目生成\n")
    
    project_name = input("请输入项目名称：").strip()
    if not project_name:
        print("❌ 项目名称不能为空")
        sys.exit(1)
    
    domains_input = input(
        "请输入领域列表 (逗号分隔，默认 user): "
    ).strip()
    domains = [d.strip() for d in domains_input.split(',')] if domains_input else ['user']
    
    print("\n选择架构风格:")
    print("  1. minimal   - 最小化架构 (仅 Domain + Application)")
    print("  2. standard  - 标准架构 (完整四层 + 基础集成)")
    print("  3. full      - 完整架构 (四层 + 全部集成 + Docker + 监控)")
    
    style_choice = input("请选择 (1/2/3, 默认 2): ").strip()
    style_map = {'1': 'minimal', '2': 'standard', '3': 'full'}
    style = style_map.get(style_choice, 'standard')
    
    with_examples = input("是否包含示例代码？(y/N): ").strip().lower() == 'y'
    with_tests = input("是否包含测试文件？(y/N): ").strip().lower() == 'y'
    with_docker = input("是否包含 Docker 配置？(y/N): ").strip().lower() == 'y'
    
    output = input("输出目录 (默认 ./generated): ").strip() or './generated'
    
    return {
        'project_name': project_name,
        'domains': domains,
        'style': style,
        'with_examples': with_examples,
        'with_tests': with_tests,
        'with_docker': with_docker,
        'output': output
    }


def validate_project_name(name):
    """验证项目名称"""
    import re
    if not re.match(r'^[a-z][a-z0-9-]{1,49}$', name):
        raise ValueError(
            f"无效的项目名称：{name}\n"
            "项目名称必须以小写字母开头，只能包含小写字母、数字和连字符"
        )
    return True


def generate_project_structure(config):
    """生成项目目录结构"""
    print(f"\n📦 开始生成项目：{config['project_name']}")
    print(f"   领域：{', '.join(config['domains'])}")
    print(f"   风格：{config['style']}")
    print(f"   输出：{config['output']}\n")
    
    base_path = Path(config['output']) / config['project_name']
    
    # 创建基础目录
    directories = [
        'cmd/server',
        'cmd/worker',
        'internal/config',
        'internal/domain',
        'internal/application',
        'internal/infrastructure/wire',
        'internal/infrastructure/persistence/gorm',
        'internal/infrastructure/cache',
        'internal/infrastructure/queue',
        'internal/infrastructure/web/middleware',
        'internal/interfaces/http',
        'pkg/errors',
        'pkg/response',
        'pkg/validator',
        'migrations/sql',
        'configs',
        'scripts',
        'tests',
    ]
    
    # 根据架构风格添加额外目录
    if config['style'] in ['standard', 'full']:
        directories.extend([
            'internal/infrastructure/message/nats',
            'internal/infrastructure/task/asynq',
        ])
    
    if config['style'] == 'full':
        directories.extend([
            'monitoring/prometheus',
            'monitoring/grafana/dashboards',
            'docs/api',
        ])
    
    # 为每个领域创建子目录
    for domain in config['domains']:
        directories.extend([
            f'internal/domain/{domain}/entity',
            f'internal/domain/{domain}/valueobject',
            f'internal/domain/{domain}/aggregate',
            f'internal/domain/{domain}/repository',
            f'internal/domain/{domain}/service',
            f'internal/domain/{domain}/event',
            f'internal/application/{domain}/service',
            f'internal/application/{domain}/dto',
            f'internal/application/{domain}/event',
            f'internal/interfaces/http/{domain}',
        ])
        
        if config['with_tests']:
            directories.extend([
                f'tests/domain/{domain}',
                f'tests/application/{domain}',
                f'tests/interfaces/{domain}',
            ])
    
    # 创建所有目录
    for directory in directories:
        dir_path = base_path / directory
        dir_path.mkdir(parents=True, exist_ok=True)
        print(f"✓ 创建目录：{directory}")
    
    return base_path


def generate_go_mod(project_name, base_path):
    """生成 go.mod 文件"""
    content = f"""module {project_name}

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.14.0
	github.com/jinzhu/gorm v1.9.16
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/spf13/viper v1.16.0
	go.uber.org/zap v1.24.0
	github.com/prometheus/client_golang v1.15.2
	github.com/casbin/casbin/v2 v2.68.0
	github.com/nats-io/nats.go v1.28.0
	github.com/hibiken/asynq v0.18.0
	github.com/gorilla/websocket v1.5.0
	github.com/stretchr/testify v1.8.4
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/files v1.0.2
	github.com/google/uuid v1.3.0
	golang.org/x/crypto v0.9.0
	golang.org/x/text v0.10.0
)
"""
    
    with open(base_path / 'go.mod', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ 创建 go.mod")


def generate_makefile(project_name, base_path):
    """生成 Makefile"""
    content = f"""# {project_name} Makefile

.PHONY: build run test clean fmt lint generate migrate-up migrate-down docker-up docker-down help

# 变量定义
BINARY_NAME={project_name}
CMD_PATH=./cmd/server
MAIN_FILE=${{CMD_PATH}}/main.go
BUILD_PATH=./bin/${{BINARY_NAME}}

# 默认目标
all: build

# 编译
build:
\t@echo "正在编译..."
\t@go build -o ${{BUILD_PATH}} ${{MAIN_FILE}}
\t@echo "编译完成：${{BUILD_PATH}}"

# 运行
run:
\t@echo "启动应用..."
\t@go run ${{MAIN_FILE}}

# 测试
test:
\t@echo "运行测试..."
\t@go test -v ./... -cover

# 清理
clean:
\t@echo "清理编译文件..."
\t@rm -rf ./bin/*
\t@go clean
\t@echo "清理完成"

# 格式化代码
fmt:
\t@echo "格式化代码..."
\t@gofmt -s -w .
\t@goimports -w .
\t@echo "格式化完成"

# 代码检查
lint:
\t@echo "运行代码检查..."
\t@golangci-lint run
\t@echo "检查完成"

# 生成依赖注入代码
generate:
\t@echo "生成 Wire 代码..."
\t@wire gen ./...
\t@echo "生成完成"

# 数据库迁移 - 向上
migrate-up:
\t@echo "执行数据库迁移..."
\t@migrate -path migrations/sql -database "postgres://localhost:5432/${{PROJECT_NAME}}?sslmode=disable" up
\t@echo "迁移完成"

# 数据库迁移 - 向下
migrate-down:
\t@echo "回滚数据库迁移..."
\t@migrate -path migrations/sql -database "postgres://localhost:5432/${{PROJECT_NAME}}?sslmode=disable" down
\t@echo "回滚完成"

# Docker - 启动
docker-up:
\t@echo "启动 Docker 容器..."
\t@docker-compose up -d
\t@echo "Docker 容器已启动"

# Docker - 停止
docker-down:
\t@echo "停止 Docker 容器..."
\t@docker-compose down
\t@echo "Docker 容器已停止"

# 帮助信息
help:
\t@echo "{project_name} - DDD Clean Architecture 项目"
\t@echo ""
\t@echo "可用命令:"
\t@echo "  make build        编译项目"
\t@echo "  make run          运行应用"
\t@echo "  make test         运行测试"
\t@echo "  make clean        清理编译文件"
\t@echo "  make fmt          格式化代码"
\t@echo "  make lint         代码检查"
\t@echo "  make generate     生成 Wire 代码"
\t@echo "  make migrate-up   执行数据库迁移"
\t@echo "  make migrate-down 回滚数据库迁移"
\t@echo "  make docker-up    启动 Docker 容器"
\t@echo "  make docker-down  停止 Docker 容器"
\t@echo ""
"""
    
    with open(base_path / 'Makefile', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ 创建 Makefile")


def generate_readme(project_name, domains, base_path):
    """生成 README.md"""
    content = f"""# {project_name}

基于 DDD Clean Architecture 的 Go 后端项目

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 14+
- Redis 7.0+
- Docker & Docker Compose

### 本地开发

```bash
# 克隆项目
git clone https://github.com/your-org/{project_name}.git
cd {project_name}

# 安装依赖
go mod download

# 配置环境变量
cp configs/config.example.yaml configs/config.yaml
# 编辑配置文件

# 启动开发环境
make docker-up

# 运行数据库迁移
make migrate-up

# 启动应用
make run
```

访问 `http://localhost:8080/api/v1` 查看 API

## 项目结构

```
{project_name}/
├── cmd/                    # 应用入口
│   ├── server/            # API 服务
│   └── worker/            # 后台任务
├── internal/              # 内部代码
│   ├── domain/           # 领域层
│   │   ├── {'/'.join(domains)}
│   ├── application/      # 应用层
│   ├── infrastructure/   # 基础设施层
│   └── interfaces/       # 接口层
├── pkg/                  # 公共包
├── migrations/           # 数据库迁移
├── configs/              # 配置文件
└── tests/               # 测试文件
```

## 可用领域

{''.join([f"- **{d}** - {d.title()}管理模块\\n" for d in domains])}

## 常用命令

```bash
make build          # 编译
make run           # 运行
make test          # 测试
make fmt           # 格式化
make lint          # 检查
make generate      # 生成代码
make docker-up     # 启动 Docker
```

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **消息队列**: NATS
- **任务调度**: Asynq
- **认证**: JWT + Casbin
- **文档**: Swagger

## 开发规范

本项目遵循 DDD Clean Architecture 架构原则：

1. **领域层** - 核心业务逻辑，无技术依赖
2. **应用层** - 协调领域对象，编排业务流程
3. **基础设施层** - 实现技术细节
4. **接口层** - 处理协议转换

## License

MIT License
"""
    
    with open(base_path / 'README.md', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ 创建 README.md")


def main():
    """主函数"""
    args = parse_args()
    
    # 交互模式
    if args.interactive:
        config = interactive_mode()
    else:
        if not args.project_name:
            print("❌ 错误：必须指定项目名称 (--project-name)")
            sys.exit(1)
        
        config = {
            'project_name': args.project_name,
            'domains': [d.strip() for d in args.domains.split(',')],
            'style': args.style,
            'output': args.output,
            'with_examples': args.with_examples,
            'with_tests': args.with_tests,
            'with_docker': args.with_docker,
        }
    
    try:
        # 验证项目名称
        validate_project_name(config['project_name'])
        
        # 生成项目结构
        base_path = generate_project_structure(config)
        
        # 生成配置文件
        generate_go_mod(config['project_name'], base_path)
        generate_makefile(config['project_name'], base_path)
        generate_readme(config['project_name'], config['domains'], base_path)
        
        print("\n✅ 项目生成完成!\n")
        print(f"下一步:")
        print(f"  cd {base_path}")
        print(f"  go mod tidy")
        print(f"  make build")
        print(f"  make run\n")
        
    except ValueError as e:
        print(f"❌ 错误：{e}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ 意外错误：{e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()
