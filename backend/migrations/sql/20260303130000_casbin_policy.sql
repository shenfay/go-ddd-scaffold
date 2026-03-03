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
COMMENT ON TABLE casbin_rule IS 'Casbin 权限策略存储表';
COMMENT ON COLUMN casbin_rule.ptype IS '策略类型：p(权限策略), g(角色关系)';
COMMENT ON COLUMN casbin_rule.v0 IS '策略字段 0(通常是角色或用户)';
COMMENT ON COLUMN casbin_rule.v1 IS '策略字段 1(通常是资源路径)';
COMMENT ON COLUMN casbin_rule.v2 IS '策略字段 2(通常是 HTTP 方法)';
COMMENT ON COLUMN casbin_rule.v3 IS '策略字段 3(通常是租户 ID)';

-- 创建索引
CREATE INDEX idx_casbin_rule_ptype ON casbin_rule(ptype);
CREATE INDEX idx_casbin_rule_v0 ON casbin_rule(v0);
CREATE INDEX idx_casbin_rule_v1 ON casbin_rule(v1);
CREATE INDEX idx_casbin_rule_v3 ON casbin_rule(v3);

-- 初始化默认策略数据
-- 注意：这些是示例数据，实际使用时应根据业务需求配置
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    -- 成员角色权限（在租户内）
    ('p', 'member', ':tenant', 'users', 'read'),
    ('p', 'member', ':tenant', 'users', 'write'),
    ('p', 'member', ':tenant', 'self', 'read'),
    ('p', 'member', ':tenant', 'self', 'write'),
    ('p', 'member', ':tenant', 'invitation', 'manage'),
    
    -- 访客角色权限
    ('p', 'guest', ':tenant', 'self', 'read'),
    ('p', 'guest', ':tenant', 'self', 'write'),
    
    -- 管理员角色权限
    ('p', 'admin', ':tenant', 'users', 'read'),
    ('p', 'admin', ':tenant', 'users', 'write'),
    ('p', 'admin', ':tenant', 'users', 'delete'),
    ('p', 'admin', ':tenant', 'settings', 'read'),
    ('p', 'admin', ':tenant', 'settings', 'write'),
    
    -- 所有者角色权限（完全控制）
    ('p', 'owner', ':tenant', '*', '*');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS casbin_rule;
