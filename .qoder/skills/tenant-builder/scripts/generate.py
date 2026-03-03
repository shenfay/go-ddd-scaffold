#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Tenant Builder - 多租户 SaaS 架构生成工具
提供租户管理、家庭角色、订阅管理等完整功能
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
        description='多租户 SaaS 架构生成工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 创建租户表
  %(prog)s create --module tenant --description "create tenants table"
  
  # 创建家庭系统表
  %(prog)s create --module family --description "create families table"
  
  # 执行迁移
  %(prog)s migrate up --env dev
  
  # 生成 DAO 代码
  %(prog)s generate --tables tenants,families --output ./internal/infrastructure/persistence/gorm
  
  # 完整工作流
  %(prog)s full-workflow --modules tenant,family,subscription --env dev --generate-dao --init-casbin
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # create 子命令
    create_parser = subparsers.add_parser('create', help='创建迁移文件')
    create_parser.add_argument('--module', type=str, required=True, 
                              choices=['tenant', 'family', 'subscription', 'invitation'],
                              help='模块名')
    create_parser.add_argument('--description', type=str, required=True, help='迁移描述')
    create_parser.add_argument('--output', type=str, default='./backend/migrations/sql', help='输出目录')
    
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
    
    # init-casbin 子命令
    casbin_parser = subparsers.add_parser('init-casbin', help='初始化 Casbin 策略')
    casbin_parser.add_argument('--default-role', type=str, default='admin', help='默认角色')
    casbin_parser.add_argument('--output', type=str, default='./backend/config/auth/rbac.conf', help='输出文件')
    
    # full-workflow 子命令
    workflow_parser = subparsers.add_parser('full-workflow', help='完整工作流')
    workflow_parser.add_argument('--modules', type=str, required=True, help='模块列表，逗号分隔')
    workflow_parser.add_argument('--env', type=str, default='dev', help='环境')
    workflow_parser.add_argument('--generate-dao', action='store_true', help='生成 DAO 代码')
    workflow_parser.add_argument('--init-casbin', action='store_true', help='初始化 Casbin')
    workflow_parser.add_argument('--with-middleware', action='store_true', help='生成中间件')
    workflow_parser.add_argument('--with-tests', action='store_true', help='生成测试文件')
    
    return parser.parse_args()


def get_goose_version_number():
    """生成 Goose 版本号（时间戳格式）"""
    return datetime.now().strftime('%Y%m%d%H%M%S')


def create_tenant_migration(description, output_dir):
    """创建租户表迁移"""
    version = get_goose_version_number()
    filename = f"{version}_{description.lower().replace(' ', '_')}.sql"
    
    migration_path = Path(output_dir) / filename
    migration_path.parent.mkdir(parents=True, exist_ok=True)
    
    content = """-- +goose Up
-- +goose StatementBegin
-- 租户表 - SaaS 多租户架构核心表
CREATE TABLE IF NOT EXISTS tenants (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    subdomain       VARCHAR(50) UNIQUE NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    plan            VARCHAR(20) DEFAULT 'free',
    expired_at      TIMESTAMP WITH TIME ZONE,
    max_users       INT DEFAULT 10,
    max_storage     BIGINT DEFAULT 1073741824,  -- 1GB
    current_users   INT DEFAULT 1,
    current_storage BIGINT DEFAULT 0,
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

-- 添加字段备注
COMMENT ON TABLE tenants IS '租户表 - SaaS 多租户架构核心';
COMMENT ON COLUMN tenants.id IS '主键 ID';
COMMENT ON COLUMN tenants.name IS '租户名称（学校/机构名）';
COMMENT ON COLUMN tenants.subdomain IS '子域名（用于租户识别）';
COMMENT ON COLUMN tenants.status IS '租户状态：active, suspended, expired, cancelled';
COMMENT ON COLUMN tenants.plan IS '套餐类型：free, basic, premium';
COMMENT ON COLUMN tenants.expired_at IS '过期时间';
COMMENT ON COLUMN tenants.max_users IS '最大用户数限制';
COMMENT ON COLUMN tenants.max_storage IS '最大存储空间（字节）';
COMMENT ON COLUMN tenants.current_users IS '当前用户数';
COMMENT ON COLUMN tenants.current_storage IS '已用存储空间';
COMMENT ON COLUMN tenants.created_at IS '创建时间';
COMMENT ON COLUMN tenants.updated_at IS '更新时间';
COMMENT ON COLUMN tenants.deleted_at IS '删除时间（软删除）';

-- 添加索引
CREATE INDEX idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_plan ON tenants(plan);
CREATE INDEX idx_tenants_expired ON tenants(expired_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants CASCADE;
-- +goose StatementEnd
"""
    
    with open(migration_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ 创建租户迁移文件：{migration_path}")
    return migration_path


def create_family_migration(description, output_dir):
    """创建家庭系统迁移"""
    version = get_goose_version_number()
    filename = f"{version}_{description.lower().replace(' ', '_')}.sql"
    
    migration_path = Path(output_dir) / filename
    migration_path.parent.mkdir(parents=True, exist_ok=True)
    
    content = """-- +goose Up
-- +goose StatementBegin
-- 家庭表 - 家庭学习组织单元
CREATE TABLE IF NOT EXISTS families (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    creator_id      BIGINT NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE families IS '家庭表 - 家庭学习组织单元';
COMMENT ON COLUMN families.id IS '主键 ID';
COMMENT ON COLUMN families.tenant_id IS '所属租户 ID（数据隔离）';
COMMENT ON COLUMN families.name IS '家庭名称';
COMMENT ON COLUMN families.creator_id IS '创建者 ID（家长）';
COMMENT ON COLUMN families.status IS '家庭状态：active, inactive';

CREATE INDEX idx_families_tenant ON families(tenant_id);
CREATE INDEX idx_families_creator ON families(creator_id);

-- 家庭成员表
CREATE TABLE IF NOT EXISTS family_members (
    id              BIGSERIAL PRIMARY KEY,
    family_id       BIGINT NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'child',
    joined_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    invited_by      BIGINT REFERENCES users(id),
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL,
    
    UNIQUE(family_id, user_id)
);

COMMENT ON TABLE family_members IS '家庭成员表';
COMMENT ON COLUMN family_members.id IS '主键 ID';
COMMENT ON COLUMN family_members.family_id IS '所属家庭 ID';
COMMENT ON COLUMN family_members.user_id IS '用户 ID';
COMMENT ON COLUMN family_members.role IS '家庭角色：parent, child, educator';
COMMENT ON COLUMN family_members.joined_at IS '加入时间';
COMMENT ON COLUMN family_members.invited_by IS '邀请人 ID';

CREATE INDEX idx_family_members_family ON family_members(family_id);
CREATE INDEX idx_family_members_user ON family_members(user_id);
CREATE INDEX idx_family_members_role ON family_members(role);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS family_members CASCADE;
DROP TABLE IF EXISTS families CASCADE;
-- +goose StatementEnd
"""
    
    with open(migration_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ 创建家庭系统迁移文件：{migration_path}")
    return migration_path


def create_subscription_migration(description, output_dir):
    """创建订阅管理迁移"""
    version = get_goose_version_number()
    filename = f"{version}_{description.lower().replace(' ', '_')}.sql"
    
    migration_path = Path(output_dir) / filename
    migration_path.parent.mkdir(parents=True, exist_ok=True)
    
    content = """-- +goose Up
-- +goose StatementBegin
-- 订阅计划表
CREATE TABLE IF NOT EXISTS subscription_plans (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(20) UNIQUE NOT NULL,
    name            VARCHAR(100) NOT NULL,
    price_monthly   DECIMAL(10,2) DEFAULT 0,
    price_yearly    DECIMAL(10,2) DEFAULT 0,
    max_users       INT DEFAULT 10,
    max_storage     BIGINT DEFAULT 1073741824,
    features        JSONB DEFAULT '[]'::jsonb,
    is_active       BOOLEAN DEFAULT true,
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE subscription_plans IS '订阅计划表';
COMMENT ON COLUMN subscription_plans.code IS '计划代码：free, basic, premium';
COMMENT ON COLUMN subscription_plans.features IS '功能特性列表（JSON 数组）';

-- 租户订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT UNIQUE NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id         BIGINT REFERENCES subscription_plans(id),
    status          VARCHAR(20) DEFAULT 'active',
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    auto_renew      BOOLEAN DEFAULT false,
    last_payment_at TIMESTAMP WITH TIME ZONE,
    next_billing_at TIMESTAMP WITH TIME ZONE,
    
    -- 标准时间戳字段
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE subscriptions IS '租户订阅表';
COMMENT ON COLUMN subscriptions.tenant_id IS '租户 ID';
COMMENT ON COLUMN subscriptions.plan_id IS '订阅计划 ID';
COMMENT ON COLUMN subscriptions.status IS '订阅状态：active, cancelled, expired';
COMMENT ON COLUMN subscriptions.start_date IS '开始日期';
COMMENT ON COLUMN subscriptions.end_date IS '结束日期';
COMMENT ON COLUMN subscriptions.auto_renew IS '是否自动续费';

CREATE INDEX idx_subscriptions_tenant ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_plan ON subscriptions(plan_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(end_date);

-- 订阅计划初始数据
INSERT INTO subscription_plans (code, name, price_monthly, price_yearly, max_users, max_storage, features) VALUES
('free', '免费版', 0, 0, 3, 104857600, '["基础游戏", "有限进度追踪"]'::jsonb),
('basic', '基础版', 29.9, 299, 10, 1073741824, '["基础游戏", "完整进度追踪", "家庭报告"]'::jsonb),
('premium', '高级版', 59.9, 599, 50, 5368709120, '["全部游戏", "高级分析", "优先支持", "API 访问"]'::jsonb);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions CASCADE;
DROP TABLE IF EXISTS subscription_plans CASCADE;
-- +goose StatementEnd
"""
    
    with open(migration_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ 创建订阅管理迁移文件：{migration_path}")
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
        entity_name = table_name.title().replace('_', '')
        
        # 生成 GORM Model
        model_code = generate_gorm_model(entity_name, table_name, config)
        model_file = model_dir / f"{table_name.lower()}.go"
        with open(model_file, 'w', encoding='utf-8') as f:
            f.write(model_code)
        print(f"✓ 生成 Model: {model_file}")
        
        # 生成 DAO 接口
        dao_interface_code = generate_dao_interface(entity_name, table_name, config)
        dao_interface_file = dao_dir / f"{table_name.lower()}_repository.go"
        with open(dao_interface_file, 'w', encoding='utf-8') as f:
            f.write(dao_interface_code)
        print(f"✓ 生成 DAO 接口：{dao_interface_file}")
        
        # 生成 DAO 实现
        dao_impl_code = generate_dao_impl(entity_name, table_name, config)
        dao_impl_file = dao_dir / f"{table_name.lower()}_repository_impl.go"
        with open(dao_impl_file, 'w', encoding='utf-8') as f:
            f.write(dao_impl_code)
        print(f"✓ 生成 DAO 实现：{dao_impl_file}")
    
    print(f"\n✅ 代码生成完成!")


def generate_gorm_model(entity_name, table_name, config):
    """生成 GORM Model 代码"""
    # 根据表名添加特定字段
    fields = ""
    if table_name == 'tenants':
        fields = """\tName           string     `gorm:"column:name;size:100;notNull" json:"name"`
\tSubdomain      string     `gorm:"column:subdomain;size:50;uniqueIndex;notNull" json:"subdomain"`
\tStatus         string     `gorm:"column:status;size:20;default:'active'" json:"status"`
\tPlan           string     `gorm:"column:plan;size:20;default:'free'" json:"plan"`
\tExpiredAt      *time.Time `gorm:"column:expired_at" json:"expired_at,omitempty"`
\tMaxUsers       int        `gorm:"column:max_users;default:10" json:"max_users"`
\tMaxStorage     int64      `gorm:"column:max_storage;default:1073741824" json:"max_storage"`
\tCurrentUsers   int        `gorm:"column:current_users;default:1" json:"current_users"`
\tCurrentStorage int64      `gorm:"column:current_storage;default:0" json:"current_storage"`"""
    elif table_name == 'families':
        fields = """\tTenantID   int64  `gorm:"column:tenant_id;notNull;index" json:"tenant_id"`
\tCreatorID  int64  `gorm:"column:creator_id;notNull" json:"creator_id"`
\tName       string `gorm:"column:name;size:100;notNull" json:"name"`
\tStatus     string `gorm:"column:status;size:20;default:'active'" json:"status"`"""
    elif table_name == 'subscriptions':
        fields = """\tTenantID      int64      `gorm:"column:tenant_id;uniqueIndex;notNull" json:"tenant_id"`
\tPlanID        *int64     `gorm:"column:plan_id;index" json:"plan_id,omitempty"`
\tStatus        string     `gorm:"column:status;size:20;default:'active'" json:"status"`
\tStartDate     time.Time  `gorm:"column:start_date;type:date;notNull" json:"start_date"`
\tEndDate       time.Time  `gorm:"column:end_date;type:date;notNull" json:"end_date"`
\tAutoRenew     bool       `gorm:"column:auto_renew;default:false" json:"auto_renew"`
\tLastPaymentAt *time.Time `gorm:"column:last_payment_at" json:"last_payment_at,omitempty"`
\tNextBillingAt *time.Time `gorm:"column:next_billing_at" json:"next_billing_at,omitempty"`"""
    
    code = f"""// Code generated by tenant-builder. DO NOT EDIT.
// source: {table_name}

package model

import (
\t"time"
)

// {entity_name} {table_name} 表模型
type {entity_name} struct {{
\tID        int64      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
{fields}
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


def generate_dao_interface(entity_name, table_name, config):
    """生成 DAO 接口"""
    code = f"""// Code generated by tenant-builder. DO NOT EDIT.
// source: {table_name}

package dao

import (
\t"context"
\t"your-project/internal/infrastructure/persistence/gorm/model"
)

// {entity_name}Repository {table_name} 仓储接口
type {entity_name}Repository interface {{
\tCreate(ctx context.Context, entity *model.{entity_name}) error
\tGetByID(ctx context.Context, id int64) (*model.{entity_name}, error)
\tUpdate(ctx context.Context, entity *model.{entity_name}) error
\tDelete(ctx context.Context, id int64) error
\tFindAll(ctx context.Context, tenantID int64, offset, limit int) ([]model.{entity_name}, int64, error)
}}
"""
    
    return code


def generate_dao_impl(entity_name, table_name, config):
    """生成 DAO 实现"""
    code = f"""// Code generated by tenant-builder. DO NOT EDIT.
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

// FindAll 分页查询（租户隔离）
func (r *{entity_name}RepositoryImpl) FindAll(ctx context.Context, tenantID int64, offset, limit int) ([]model.{entity_name}, int64, error) {{
\tvar entities []model.{entity_name}
\tvar total int64
\t
\ttx := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
\t
\tif err := tx.Model(&model.{entity_name}{{}}).Count(&total).Error; err != nil {{
\t\treturn nil, 0, err
\t}}
\t
\terr := tx.Offset(offset).Limit(limit).Find(&entities).Error
\tif err != nil {{
\t\treturn nil, 0, err
\t}}
\t
\treturn entities, total, nil
}}
"""
    
    return code


def init_casbin(default_role, output_file):
    """初始化 Casbin 策略文件"""
    output_path = Path(output_file)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    content = f"""[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "admin"

# ============================================
# 租户管理员策略
# ============================================
p, {default_role}, tenants, read, allow
p, {default_role}, tenants, update, allow
p, {default_role}, families, manage, allow
p, {default_role}, subscriptions, view, allow
p, {default_role}, users, invite, allow

# ============================================
# 家长角色策略
# ============================================
p, parent, children, read, allow
p, parent, children, manage, allow
p, parent, learning_progress, read, allow
p, parent, learning_progress, view_reports, allow
p, parent, subscriptions, view, allow
p, parent, families, join, allow

# ============================================
# 孩子角色策略
# ============================================
p, child, learning_resources, read, allow
p, child, games, play, allow
p, child, learning_progress, update, allow
p, child, families, view, allow

# ============================================
# 教育工作者角色策略
# ============================================
p, educator, students, read, allow
p, educator, students, assign_tasks, allow
p, educator, learning_progress, read, allow
p, educator, learning_progress, analyze, allow

# ============================================
# 角色继承关系
# ============================================
g, admin, {default_role}
"""
    
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✓ 初始化 Casbin 策略文件：{output_path}")


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
            'host': 'staging-db.mathfun.local',
            'port': '5432',
            'user': 'mathfun_staging',
            'password': '${STAGING_DB_PASSWORD}',
            'name': 'mathfun_staging',
            'ssl_mode': 'disable',
        },
        'prod': {
            'host': 'prod-db.mathfun.local',
            'port': '5432',
            'user': 'mathfun_prod',
            'password': '${PROD_DB_PASSWORD}',
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
        if args.module == 'tenant':
            create_tenant_migration(args.description, args.output)
        elif args.module == 'family':
            create_family_migration(args.description, args.output)
        elif args.module == 'subscription':
            create_subscription_migration(args.description, args.output)
        else:
            print(f"❌ 未知的模块：{args.module}")
            sys.exit(1)
    
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
    
    elif args.command == 'init-casbin':
        init_casbin(args.default_role, args.output)
    
    elif args.command == 'full-workflow':
        modules = [m.strip() for m in args.modules.split(',')]
        print(f"\n🚀 开始完整工作流：{', '.join(modules)}\n")
        
        # Step 1: 创建迁移
        print("Step 1: 创建迁移文件...")
        for module in modules:
            if module == 'tenant':
                create_tenant_migration(f"create {module}s table", './backend/migrations/sql')
            elif module == 'family':
                create_family_migration(f"create {module} system", './backend/migrations/sql')
            elif module == 'subscription':
                create_subscription_migration(f"create {module} management", './backend/migrations/sql')
        
        # Step 2: 提示用户编辑
        print("\n⚠️  请检查迁移文件，然后按回车继续...")
        input()
        
        # Step 3: 执行迁移
        print("\nStep 2: 执行迁移...")
        run_migrate('up', args.env)
        
        # Step 4: 生成 DAO
        if args.generate_dao:
            print("\nStep 3: 生成 DAO 代码...")
            tables = ','.join([f"{m}s" if m != 'family' else 'families,family_members' for m in modules])
            generate_dao(tables, './internal/infrastructure/persistence/gorm', {})
        
        # Step 5: 初始化 Casbin
        if args.init_casbin:
            print("\nStep 4: 初始化 Casbin 策略...")
            init_casbin('admin', './backend/config/auth/rbac.conf')
        
        print("\n✅ 完整工作流执行完成!\n")
    
    else:
        print("❌ 错误：未知命令。使用 --help 查看可用命令")
        sys.exit(1)


if __name__ == '__main__':
    main()
