#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
API Doc Generator - API 文档生成器
基于 swaggo/swag 自动生成 Swagger/OpenAPI 文档
"""

import os
import sys
import subprocess
import shutil
from pathlib import Path


class ColorLogger:
    """彩色日志输出"""
    
    BLUE = '\033[0;34m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    RED = '\033[0;31m'
    NC = '\033[0m'  # No Color
    
    @classmethod
    def info(cls, msg):
        print(f"{cls.BLUE}[INFO]{cls.NC} {msg}")
    
    @classmethod
    def success(cls, msg):
        print(f"{cls.GREEN}[SUCCESS]{cls.NC} {msg}")
    
    @classmethod
    def warning(cls, msg):
        print(f"{cls.YELLOW}[WARNING]{cls.NC} {msg}")
    
    @classmethod
    def error(cls, msg):
        print(f"{cls.RED}[ERROR]{cls.NC} {msg}")


def get_project_root():
    """动态获取项目根目录"""
    script_dir = Path(__file__).parent
    return (script_dir / '../../../..').resolve()


def check_dependencies():
    """检查依赖工具"""
    ColorLogger.info("检查必要的工具...")
    
    missing_tools = []
    
    # 检查 go
    if not shutil.which('go'):
        missing_tools.append('go')
    
    # 检查 swag
    if not shutil.which('swag'):
        ColorLogger.warning("swag 工具未安装，正在安装...")
        subprocess.run(['go', 'install', 'github.com/swaggo/swag/cmd/swag@latest'], check=True)
    
    if missing_tools:
        ColorLogger.error(f"缺少必要的工具：{' '.join(missing_tools)}")
        ColorLogger.info("请安装缺失的工具后再继续")
        sys.exit(1)
    
    ColorLogger.success("所有依赖工具检查通过")


def validate_project_structure(project_root):
    """验证 Go 项目结构"""
    ColorLogger.info("验证项目结构...")
    
    required_paths = [
        project_root / "backend" / "go.mod",
        project_root / "backend" / "internal" / "interfaces" / "http",
    ]
    
    for path in required_paths:
        if not path.exists():
            ColorLogger.error(f"项目结构不完整：缺少 {path.relative_to(project_root)}")
            return False
    
    ColorLogger.success("项目结构验证通过")
    return True


def generate_docs(format='json', package='', output_dir=None, project_root=None):
    """生成 API 文档"""
    if project_root is None:
        project_root = get_project_root()
    
    if output_dir is None:
        output_dir = project_root / "backend" / "docs"
    
    ColorLogger.info(f"生成 API 文档 (格式：{format})...")
    
    # 创建输出目录
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # 构建 swag 命令
    cmd = [
        'swag', 'init',
        '-g', 'main.go',
        '--parseInternal',
        '--parseDependency',
        '-d', str(project_root / 'cmd' / 'server'),
        '-o', str(output_dir),
    ]
    
    if format:
        cmd.extend(['--ot', format])
    
    if package:
        cmd.extend(['--pd', package])
    
    ColorLogger.info(f"执行命令：{' '.join(cmd)}")
    
    try:
        # 切换到 backend 目录执行
        backend_dir = project_root / 'backend'
        result = subprocess.run(
            cmd,
            cwd=str(backend_dir),
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            ColorLogger.error(f"生成失败：{result.stderr}")
            return False
        
        ColorLogger.success(f"API 文档生成成功：{output_dir}")
        
        # 显示生成的文件
        generated_files = list(output_dir.glob('*.*'))
        if generated_files:
            ColorLogger.info("生成的文件:")
            for f in generated_files:
                print(f"  - {f.name}")
        
        return True
        
    except subprocess.CalledProcessError as e:
        ColorLogger.error(f"执行失败：{e}")
        return False
    except Exception as e:
        ColorLogger.error(f"意外错误：{e}")
        return False


def clean_docs(output_dir=None, project_root=None):
    """清理生成的文档"""
    if project_root is None:
        project_root = get_project_root()
    
    if output_dir is None:
        output_dir = project_root / "backend" / "docs"
    
    ColorLogger.info("清理生成的文档...")
    
    files_to_clean = ['docs.json', 'docs.yaml', 'swagger.json', 'swagger.yaml']
    
    for filename in files_to_clean:
        file_path = output_dir / filename
        if file_path.exists():
            file_path.unlink()
            ColorLogger.info(f"已删除：{filename}")
    
    ColorLogger.success("清理完成")


def preview_docs(output_dir=None, project_root=None):
    """预览文档内容"""
    if project_root is None:
        project_root = get_project_root()
    
    if output_dir is None:
        output_dir = project_root / "backend" / "docs"
    
    json_file = output_dir / "docs.json"
    yaml_file = output_dir / "docs.yaml"
    
    if not json_file.exists() and not yaml_file.exists():
        ColorLogger.error("文档不存在，请先生成文档")
        return
    
    ColorLogger.info("API 文档概览:")
    
    # 尝试读取 JSON 文件
    if json_file.exists():
        import json
        try:
            with open(json_file, 'r', encoding='utf-8') as f:
                data = json.load(f)
                print(f"\nAPI 版本：{data.get('info', {}).get('version', 'unknown')}")
                print(f"标题：{data.get('info', {}).get('title', 'unknown')}")
                print(f"路径数量：{len(data.get('paths', {}))}")
        except Exception as e:
            ColorLogger.warning(f"无法解析 JSON: {e}")
    
    # 显示文件列表
    ColorLogger.info("\n生成的文件:")
    for f in output_dir.iterdir():
        if f.is_file():
            size = f.stat().st_size
            print(f"  - {f.name} ({size:,} bytes)")


def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(description='API 文档生成器')
    parser.add_argument('--action', type=str, default='generate',
                       choices=['generate', 'clean', 'preview'],
                       help='操作类型')
    parser.add_argument('--format', type=str, default='json',
                       choices=['json', 'yaml'],
                       help='输出格式')
    parser.add_argument('--package', type=str, default='',
                       help='包名')
    parser.add_argument('--output', type=str, default='',
                       help='输出目录')
    
    args = parser.parse_args()
    
    project_root = get_project_root()
    output_dir = Path(args.output) if args.output else None
    
    if args.action == 'generate':
        # 检查依赖
        check_dependencies()
        
        # 验证项目结构
        if not validate_project_structure(project_root):
            sys.exit(1)
        
        # 生成文档
        success = generate_docs(
            format=args.format,
            package=args.package,
            output_dir=output_dir,
            project_root=project_root
        )
        
        sys.exit(0 if success else 1)
        
    elif args.action == 'clean':
        clean_docs(output_dir, project_root)
        
    elif args.action == 'preview':
        preview_docs(output_dir, project_root)


if __name__ == '__main__':
    main()
