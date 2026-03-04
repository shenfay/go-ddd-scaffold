-- +goose Up
-- 为用户表添加 role 和 tenant_id 字段

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'member' 
CHECK (role IN ('member', 'guest', 'super_admin', 'content_admin', 'ops_admin'));

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL;

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);

-- 添加注释
COMMENT ON COLUMN users.role IS '用户角色：member/guest/super_admin/content_admin/ops_admin';
COMMENT ON COLUMN users.tenant_id IS '所属租户 ID（多租户场景）';

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_tenant_id;

ALTER TABLE users DROP COLUMN IF EXISTS role;
ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
