-- +goose Up
-- SQL section 'Up' is executed when this migration is applied

-- ============================================
-- Casbin 策略存储表
-- 用于存储 RBAC + 多租户权限策略
-- ============================================

-- 创建 casbin_rule 表（如果使用 gorm-adapter，可自动创建）
CREATE TABLE IF NOT EXISTS casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL DEFAULT 'p',
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100),
    v6 VARCHAR(100),
    v7 VARCHAR(100),
    v8 VARCHAR(100),
    v9 VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 添加注释
COMMENT ON TABLE casbin_rule IS 'Casbin权限策略存储表';
COMMENT ON COLUMN casbin_rule.ptype IS '策略类型: p(权限策略), g(角色关系)';

-- 创建索引优化查询性能
CREATE INDEX idx_casbin_rule_ptype ON casbin_rule(ptype);
CREATE INDEX idx_casbin_rule_v0 ON casbin_rule(v0);
CREATE INDEX idx_casbin_rule_v0_v1 ON casbin_rule(v0, v1);
CREATE INDEX idx_casbin_rule_v1 ON casbin_rule(v1);

-- ============================================
-- 预置系统级权限策略
-- ============================================

-- 超级管理员：拥有所有权限
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'super_admin', 'all', 'all', 'all')
ON CONFLICT DO NOTHING;

-- 内容管理员：可以管理知识内容
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'content_admin', 'knowledge', 'read', 'allow'),
('p', 'content_admin', 'knowledge', 'write', 'allow'),
('p', 'content_admin', 'knowledge', 'delete', 'allow'),
('p', 'content_admin', 'kg_domains', 'read', 'allow'),
('p', 'content_admin', 'kg_domains', 'write', 'allow'),
('p', 'content_admin', 'kg_nodes', 'read', 'allow'),
('p', 'content_admin', 'kg_nodes', 'write', 'allow')
ON CONFLICT DO NOTHING;

-- 运营管理员：可以管理用户
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'ops_admin', 'users', 'read', 'allow'),
('p', 'ops_admin', 'users', 'write', 'allow'),
('p', 'ops_admin', 'tenants', 'read', 'allow'),
('p', 'ops_admin', 'tenants', 'write', 'allow')
ON CONFLICT DO NOTHING;

-- ============================================
-- 预置租户级权限策略（默认模板）
-- 格式: p, 角色, 租户ID, 资源, 操作
-- ============================================

-- 爸爸角色权限
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'father', ':tenant', 'children', 'read'),
('p', 'father', ':tenant', 'children', 'write'),
('p', 'father', ':tenant', 'children', 'delete'),
('p', 'father', ':tenant', 'progress', 'read'),
('p', 'father', ':tenant', 'self', 'read'),
('p', 'father', ':tenant', 'self', 'write'),
('p', 'father', ':tenant', 'invitation', 'manage'),
('p', 'father', ':tenant', 'settings', 'read'),
('p', 'father', ':tenant', 'settings', 'write')
ON CONFLICT DO NOTHING;

-- 妈妈角色权限
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'mother', ':tenant', 'children', 'read'),
('p', 'mother', ':tenant', 'children', 'write'),
('p', 'mother', ':tenant', 'progress', 'read'),
('p', 'mother', ':tenant', 'self', 'read'),
('p', 'mother', ':tenant', 'self', 'write'),
('p', 'mother', ':tenant', 'invitation', 'manage'),
('p', 'mother', ':tenant', 'settings', 'read'),
('p', 'mother', ':tenant', 'settings', 'write')
ON CONFLICT DO NOTHING;

-- 孩子角色权限
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'child', ':tenant', 'progress', 'read'),
('p', 'child', ':tenant', 'progress', 'write'),
('p', 'child', ':tenant', 'self', 'read'),
('p', 'child', ':tenant', 'self', 'write'),
('p', 'child', ':tenant', 'knowledge', 'read'),
('p', 'child', ':tenant', 'learning', 'access')
ON CONFLICT DO NOTHING;

-- 爷爷/奶奶角色权限（只读）
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
('p', 'grandpa', ':tenant', 'children', 'read'),
('p', 'grandpa', ':tenant', 'progress', 'read'),
('p', 'grandpa', ':tenant', 'self', 'read'),
('p', 'grandpa', ':tenant', 'self', 'write'),
('p', 'grandma', ':tenant', 'children', 'read'),
('p', 'grandma', ':tenant', 'progress', 'read'),
('p', 'grandma', ':tenant', 'self', 'read'),
('p', 'grandma', ':tenant', 'self', 'write')
ON CONFLICT DO NOTHING;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS casbin_rule;
