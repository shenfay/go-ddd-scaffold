#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DB Migrator - Goose 数据库迁移与 DAO 生成工具
基于 Goose 的数据库迁移管理、GORM Model 和 DAO 代码自动生成
"""

import os
import sys
import argparse
import subprocess
from pathlib import Path
from datetime import datetime


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description='Goose 数据库迁移与 DAO 生成工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 创建迁移文件
  %(prog)s create --table users --description "create users table"
  
  # 执行迁移
  %(prog)s migrate up --env dev
  
  # 生成 GORM Model 和 DAO
  %(prog)s generate --tables users,products --output ./internal/infrastructure/persistence/gorm
  
  # 完整工作流
  %(prog)s full-workflow --table orders --description "create orders table" --env dev --generate-dao
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # create 子命令
    create_parser = subparsers.add_parser('create', help='创建迁移文件')
    create_parser.add_argument('--table', type=str, required=True, help='表名')
    create_parser.add_argument('--description', type=str, required=True, help='迁移描述')
    create_parser.add_argument('--output', type=str, default='./migrations/sql', help='输出目录')
    
    # migrate 子命令
    migrate_parser = subparsers.add_parser('migrate', help='执行迁移')
    migrate_subparsers = migrate_parser.add_subparsers(dest='action', help='迁移操作')
    
    migrate_up = migrate_subparsers.add_parser('up', help='执行所有待执行的迁移')
    migrate_up.add_argument('--env', type=str, default='dev', help='环境 (dev/staging/prod)')
    
    migrate_down = migrate_subparsers.add_parser('down', help='回滚最后一个迁移')
    migrate_down.add_argument('--env', type=str, default='dev', help='环境')
    
    migrate_status = migrate_subparsers.add_parser('status', help='查看迁移状态')
    migrate_status.add_argument('--env', type=str, default='dev', help='环境')
    
    # generate 子命令
    generate_parser = subparsers.add_parser('generate', help='生成 GORM Model 和 DAO')
    generate_parser.add_argument('--tables', type=str, required=True, help='表名列表，逗号分隔')
    generate_parser.add_argument('--output', type=str, default='./internal/infrastructure/persistence/gorm', help='输出目录')
    generate_parser.add_argument('--with-comments', action='store_true', default=True, help='包含字段备注')
    generate_parser.add_argument('--soft-delete', action='store_true', default=True, help='启用软删除')
    generate_parser.add_argument('--with-tests', action='store_true', help='生成测试文件')
    
    # full-workflow 子命令
    workflow_parser = subparsers.add_parser('full-workflow', help='完整工作流')
    workflow_parser.add_argument('--table', type=str, required=True, help='表名')
    workflow_parser.add_argument('--description', type=str, required=True, help='迁移描述')
    workflow_parser.add_argument('--env', type=str, default='dev', help='环境')
    workflow_parser.add_argument('--generate-dao', action='store_true', help='生成 DAO 代码')
    workflow_parser.add_argument('--with-tests', action='store_true', help='生成测试文件')
    
    return parser.parse_args()


def get_goose_version_number():
    """生成 Goose 版本号（时间戳格式）"""
    return datetime.now().strftime('%Y%m%d%H%M%S')


def create_migration(table_name, description, output_dir):
    """创建 Goose 迁移文件"""
    version = get_goose_version_number()
    safe_description = description.lower().replace(' ', '_')
    filename = f"{version}_{safe_description}.sql"
    
    migration_path = Path(output_dir) / filename
    migration_path.parent.mkdir(parents=True, exist_ok=True)
    
    # 生成迁移模板
    content = f"""-- +goose Up
-- +goose StatementBegin
-- 创建 {table_name} 表
CREATE TABLE IF NOT EXISTS {table_name} (
    id              BIGSERIAL PRIMARY KEY,
    -- TODO: 添加其他字段
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

-- 添加字段备注
COMMENT ON TABLE {table_name} IS '{description}';
COMMENT ON COLUMN {table_name}.id IS '主键 ID';
COMMENT ON COLUMN {table_name}.created_at IS '创建时间';
COMMENT ON COLUMN {table_name}.updated_at IS '更新时间';
COMMENT ON COLUMN {table_name}.deleted_at IS '删除时间（软删除）';

-- TODO: 添加索引
-- CREATE INDEX idx_{table_name}_column ON {table_name}(column_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS {table_name};
-- +goose StatementEnd
"""
    
    with open(migration_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ 创建迁移文件：{migration_path}")
    print(f"\n下一步:")
    print(f"  1. 编辑文件添加业务字段")
    print(f"  2. 运行：db-migrator migrate up --env dev")
    print(f"  3. 生成代码：db-migrator generate --tables {table_name}\n")
    
    return migration_path


def run_migrate(action, env):
    """执行数据库迁移"""
    db_config = load_db_config(env)
    
    if action == 'up':
        cmd = ['goose', 'postgres', build_connection_string(db_config), 'up']
    elif action == 'down':
        cmd = ['goose', 'postgres', build_connection_string(db_config), 'down']
    elif action == 'status':
        cmd = ['goose', 'postgres', build_connection_string(db_config), 'status']
    else:
        print(f"❌ 未知的迁移操作：{action}")
        return
    
    print(f"🚀 执行迁移：{' '.join(cmd)}")
    try:
        subprocess.run(cmd, check=True)
        print("✅ 迁移执行成功")
    except subprocess.CalledProcessError as e:
        print(f"❌ 迁移失败：{e}")
        sys.exit(1)


def generate_dao(tables, output_dir, config):
    """生成 GORM Model 和 DAO 代码"""
    table_list = [t.strip() for t in tables.split(',')]
    
    base_path = Path(output_dir)
    model_dir = base_path / 'model'
    dao_dir = base_path / 'dao'
    
    model_dir.mkdir(parents=True, exist_ok=True)
    dao_dir.mkdir(parents=True, exist_ok=True)
    
    for table_name in table_list:
        # 生成 GORM Model
        model_code = generate_gorm_model(table_name, config)
        model_file = model_dir / f"{table_name.lower()}.go"
        with open(model_file, 'w', encoding='utf-8') as f:
            f.write(model_code)
        print(f"✓ 生成 Model: {model_file}")
        
        # 生成 DAO 接口
        dao_interface_code = generate_dao_interface(table_name, config)
        dao_interface_file = dao_dir / f"{table_name.lower()}_repository.go"
        with open(dao_interface_file, 'w', encoding='utf-8') as f:
            f.write(dao_interface_code)
        print(f"✓ 生成 DAO 接口：{dao_interface_file}")
        
        # 生成 DAO 实现
        dao_impl_code = generate_dao_impl(table_name, config)
        dao_impl_file = dao_dir / f"{table_name.lower()}_repository_impl.go"
        with open(dao_impl_file, 'w', encoding='utf-8') as f:
            f.write(dao_impl_code)
        print(f"✓ 生成 DAO 实现：{dao_impl_file}")
    
    print(f"\n✅ 代码生成完成!")
    print(f"下一步:")
    print(f"  1. 检查生成的代码")
    print(f"  2. 根据需要修改业务逻辑")
    print(f"  3. 运行 go mod tidy 下载依赖\n")


def generate_gorm_model(table_name, config):
    """生成 GORM Model 代码"""
    entity_name = table_name.title().replace('_', '')
    
    code = f"""// Code generated by db-migrator. DO NOT EDIT.
// source: {table_name}

package model

import (
\t"time"
)

// {entity_name} {table_name} 表模型
type {entity_name} struct {{
\tID        int64     `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
\t// TODO: 添加其他字段
	
\t// 标准时间戳字段
\tCreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
\tUpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
\tDeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}}

// TableName 指定表名
func ({entity_name}) TableName() string {{
\treturn "{table_name}"
}}
"""
    
    return code


def generate_dao_interface(table_name, config):
    """生成 DAO 接口"""
    entity_name = table_name.title().replace('_', '')
    
    code = f"""// Code generated by db-migrator. DO NOT EDIT.
// source: {table_name}

package dao

import (
\t"context"
\t"your-project/internal/infrastructure/persistence/gorm/model"
)

// {entity_name}Repository {table_name} 仓储接口
type {entity_name}Repository interface {{
\t// Create 创建记录
\tCreate(ctx context.Context, entity *model.{entity_name}) error
\t
\t// GetByID 根据 ID 获取
\tGetByID(ctx context.Context, id int64) (*model.{entity_name}, error)
\t
\t// Update 更新记录
\tUpdate(ctx context.Context, entity *model.{entity_name}) error
\t
\t// Delete 删除记录（软删除）
\tDelete(ctx context.Context, id int64) error
\t
\t// FindAll 获取所有记录
\tFindAll(ctx context.Context, offset, limit int) ([]model.{entity_name}, int64, error)
}}
"""
    
    return code


def generate_dao_impl(table_name, config):
    """生成 DAO 实现"""
    entity_name = table_name.title().replace('_', '')
    
    code = f"""// Code generated by db-migrator. DO NOT EDIT.
// source: {table_name}

package dao

import (
\t"context"
\t"errors"
\t"gorm.io/gorm"
\t"your-project/internal/infrastructure/persistence/gorm/model"
)

// {entity_name}RepositoryImpl {table_name} 仓储实现
type {entity_name}RepositoryImpl struct {{
\tdb *gorm.DB
}}

// New{entity_name}Repository 构造函数
func New{entity_name}Repository(db *gorm.DB) {entity_name}Repository {{
\treturn &{entity_name}RepositoryImpl{{
\t\tdb: db,
\t}}
}}

// Create 创建记录
func (r *{entity_name}RepositoryImpl) Create(ctx context.Context, entity *model.{entity_name}) error {{
\treturn r.db.WithContext(ctx).Create(entity).Error
}}

// GetByID 根据 ID 获取
func (r *{entity_name}RepositoryImpl) GetByID(ctx context.Context, id int64) (*model.{entity_name}, error) {{
\tvar entity model.{entity_name}
\terr := r.db.WithContext(ctx).First(&entity, id).Error
\tif err != nil {{
\t\tif errors.Is(err, gorm.ErrRecordNotFound) {{
\t\t\treturn nil, nil
\t\t}}
\t\treturn nil, err
\t}}
\treturn &entity, nil
}}

// Update 更新记录
func (r *{entity_name}RepositoryImpl) Update(ctx context.Context, entity *model.{entity_name}) error {{
\treturn r.db.WithContext(ctx).Save(entity).Error
}}

// Delete 删除记录（软删除）
func (r *{entity_name}RepositoryImpl) Delete(ctx context.Context, id int64) error {{
\treturn r.db.WithContext(ctx).Delete(&model.{entity_name}{{ID: id}}).Error
}}

// FindAll 获取所有记录
func (r *{entity_name}RepositoryImpl) FindAll(ctx context.Context, offset, limit int) ([]model.{entity_name}, int64, error) {{
\tvar entities []model.{entity_name}
\tvar total int64
\t
\tif err := r.db.WithContext(ctx).Model(&model.{entity_name}{{}}).Count(&total).Error; err != nil {{
\t\treturn nil, 0, err
\t}}
\t
\terr := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
\tif err != nil {{
\t\treturn nil, 0, err
\t}}
\t
\treturn entities, total, nil
}}
"""
    
    return code


def load_db_config(env):
    """加载数据库配置"""
    configs = {
        'dev': {
            'host': 'localhost',
            'port': '5432',
            'user': 'mathfun',
            'password': 'math111',
            'name': 'mathfun_dev',
            'ssl_mode': 'disable',
        },
        'staging': {
            'host': 'staging-db.example.com',
            'port': '5432',
            'user': 'staging_user',
            'password': 'staging_pass',
            'name': 'mathfun_staging',
            'ssl_mode': 'disable',
        },
        'prod': {
            'host': 'prod-db.example.com',
            'port': '5432',
            'user': 'prod_user',
            'password': 'prod_pass',
            'name': 'mathfun_prod',
            'ssl_mode': 'require',
        },
    }
    return configs.get(env, configs['dev'])


def build_connection_string(config):
    """构建 PostgreSQL 连接字符串"""
    return f"postgresql://{config['user']}:{config['password']}@{config['host']}:{config['port']}/{config['name']}?sslmode={config['ssl_mode']}"


def main():
    """主函数"""
    args = parse_args()
    
    if args.command == 'create':
        create_migration(args.table, args.description, args.output)
    
    elif args.command == 'migrate':
        if not args.action:
            print("❌ 错误：必须指定迁移操作 (up/down/status)")
            sys.exit(1)
        run_migrate(args.action, args.env)
    
    elif args.command == 'generate':
        config = {
            'with_comments': args.with_comments,
            'soft_delete': args.soft_delete,
        }
        generate_dao(args.tables, args.output, config)
    
    elif args.command == 'full-workflow':
        print(f"\n🚀 开始完整工作流：{args.table}\n")
        
        # Step 1: 创建迁移
        print("Step 1: 创建迁移文件...")
        migration_file = create_migration(args.table, args.description, './migrations/sql')
        
        # Step 2: 提示用户编辑
        print("\n⚠️  请编辑迁移文件添加业务字段，然后按回车继续...")
        input()
        
        # Step 3: 执行迁移
        print("\nStep 2: 执行迁移...")
        run_migrate('up', args.env)
        
        # Step 4: 生成 DAO
        if args.generate_dao:
            print("\nStep 3: 生成 DAO 代码...")
            generate_dao(args.table, './internal/infrastructure/persistence/gorm', {})
        
        print("\n✅ 完整工作流执行完成!\n")
    
    else:
        print("❌ 错误：未知命令。使用 --help 查看可用命令")
        sys.exit(1)


if __name__ == '__main__':
    main()
