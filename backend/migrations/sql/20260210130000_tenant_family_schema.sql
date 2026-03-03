-- +goose Up
-- SQL section 'Up' is executed when this migration is applied

-- ============================================
-- 1. 创建租户表 (tenants) - 包含订阅信息
-- ============================================
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    subscription_type VARCHAR(20) NOT NULL DEFAULT 'trial' CHECK (subscription_type IN ('trial', 'monthly', 'yearly', 'lifetime')),
    subscription_started_at TIMESTAMP WITH TIME ZONE,
    subscription_expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    trial_days INTEGER DEFAULT 7,
    max_children INTEGER NOT NULL DEFAULT 3 CHECK (max_children >= 1 AND max_children <= 10),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加租户表注释
COMMENT ON TABLE tenants IS '家庭租户表（包含订阅信息）';
COMMENT ON COLUMN tenants.id IS '租户唯一标识';
COMMENT ON COLUMN tenants.name IS '租户名称';
COMMENT ON COLUMN tenants.description IS '租户描述';
COMMENT ON COLUMN tenants.subscription_type IS '订阅类型：trial试用, monthly月付, yearly年付, lifetime终身';
COMMENT ON COLUMN tenants.subscription_started_at IS '订阅开始时间';
COMMENT ON COLUMN tenants.subscription_expired_at IS '订阅过期时间';
COMMENT ON COLUMN tenants.trial_days IS '试用天数';
COMMENT ON COLUMN tenants.max_children IS '最大子女数量限制';
COMMENT ON COLUMN tenants.status IS '租户状态：active活跃, expired过期, suspended暂停';
COMMENT ON COLUMN tenants.created_at IS '创建时间';
COMMENT ON COLUMN tenants.updated_at IS '更新时间';

-- 创建租户索引
CREATE INDEX idx_tenants_subscription_expired_at ON tenants(subscription_expired_at);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_subscription_type ON tenants(subscription_type);

-- ============================================
-- 2. 创建用户表 (users) - 纯用户信息
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    avatar VARCHAR(500),
    phone VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'locked')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加用户表注释
COMMENT ON TABLE users IS '用户基础信息表';
COMMENT ON COLUMN users.id IS '用户唯一标识';
COMMENT ON COLUMN users.email IS '用户邮箱';
COMMENT ON COLUMN users.password IS '用户密码';
COMMENT ON COLUMN users.nickname IS '用户昵称';
COMMENT ON COLUMN users.avatar IS '用户头像';
COMMENT ON COLUMN users.phone IS '用户手机号';
COMMENT ON COLUMN users.status IS '用户状态：active/inactive/locked';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

-- 创建用户索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- ============================================
-- 3. 创建家庭成员表 (tenant_members) - 关联用户、租户、角色
-- ============================================
CREATE TABLE tenant_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('parent', 'child', 'admin')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'removed')),
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- 确保每个用户在同一租户中只有一个角色
    UNIQUE(tenant_id, user_id)
);

-- 添加家庭成员表注释
COMMENT ON TABLE tenant_members IS '租户成员关系表';
COMMENT ON COLUMN tenant_members.id IS '成员关系唯一标识';
COMMENT ON COLUMN tenant_members.tenant_id IS '关联的家庭租户ID';
COMMENT ON COLUMN tenant_members.user_id IS '关联的用户ID';
COMMENT ON COLUMN tenant_members.role IS '成员角色：parent父母, child子女, admin管理员';
COMMENT ON COLUMN tenant_members.status IS '成员状态：active活跃, inactive非活跃, removed已移除';
COMMENT ON COLUMN tenant_members.invited_by IS '邀请人ID';
COMMENT ON COLUMN tenant_members.joined_at IS '加入时间';
COMMENT ON COLUMN tenant_members.left_at IS '离开时间';
COMMENT ON COLUMN tenant_members.created_at IS '创建时间';
COMMENT ON COLUMN tenant_members.updated_at IS '更新时间';

-- 创建家庭成员索引
CREATE INDEX idx_tenant_members_tenant ON tenant_members(tenant_id);
CREATE INDEX idx_tenant_members_user ON tenant_members(user_id);
CREATE INDEX idx_tenant_members_role ON tenant_members(role);
CREATE INDEX idx_tenant_members_status ON tenant_members(status);
CREATE INDEX idx_tenant_members_invited_by ON tenant_members(invited_by);

-- 添加约束：确保每个租户的父母角色数量限制
CREATE UNIQUE INDEX idx_tenant_parent_unique 
ON tenant_members(tenant_id, role) 
WHERE role = 'parent' AND status = 'active';

-- ============================================
-- 4. 创建租户邀请表 (tenant_invitations)
-- ============================================
CREATE TABLE tenant_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invitee_email VARCHAR(255) NOT NULL,
    invitee_role VARCHAR(20) NOT NULL CHECK (invitee_role IN ('parent', 'child')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'expired', 'cancelled')),
    max_uses INTEGER DEFAULT 1,
    used_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 添加租户邀请表注释
COMMENT ON TABLE tenant_invitations IS '租户家庭成员邀请表';
COMMENT ON COLUMN tenant_invitations.id IS '邀请记录唯一标识';
COMMENT ON COLUMN tenant_invitations.code IS '邀请码';
COMMENT ON COLUMN tenant_invitations.tenant_id IS '关联的家庭租户ID';
COMMENT ON COLUMN tenant_invitations.creator_id IS '邀请创建者';
COMMENT ON COLUMN tenant_invitations.invitee_email IS '被邀请人邮箱';
COMMENT ON COLUMN tenant_invitations.invitee_role IS '被邀请人角色（parent/child）';
COMMENT ON COLUMN tenant_invitations.status IS '邀请状态：pending待接受, accepted已接受, expired已过期, cancelled已取消';
COMMENT ON COLUMN tenant_invitations.max_uses IS '最大使用次数';
COMMENT ON COLUMN tenant_invitations.used_count IS '已使用次数';
COMMENT ON COLUMN tenant_invitations.expires_at IS '邀请过期时间';
COMMENT ON COLUMN tenant_invitations.created_at IS '创建时间';
COMMENT ON COLUMN tenant_invitations.updated_at IS '更新时间';

-- 创建租户邀请索引
CREATE INDEX idx_tenant_invitations_tenant ON tenant_invitations(tenant_id);
CREATE INDEX idx_tenant_invitations_creator ON tenant_invitations(creator_id);
CREATE INDEX idx_tenant_invitations_code ON tenant_invitations(code);
CREATE INDEX idx_tenant_invitations_status ON tenant_invitations(status);
CREATE INDEX idx_tenant_invitations_email ON tenant_invitations(invitee_email);

-- ============================================
-- 5. 添加外键约束
-- ============================================
-- tenant_members 表的外键约束已在创建时定义

-- ============================================
-- 6. 创建视图和函数
-- ============================================

-- 创建租户成员统计视图
CREATE OR REPLACE VIEW tenant_member_stats AS
SELECT 
    t.id as tenant_id,
    t.name as tenant_name,
    t.subscription_type,
    t.subscription_expired_at,
    t.max_children,
    COALESCE(parent_stats.parent_count, 0) as parent_count,
    COALESCE(child_stats.child_count, 0) as child_count,
    t.status as tenant_status,
    CASE 
        WHEN t.subscription_expired_at > NOW() AND t.status = 'active' THEN 'active'
        WHEN t.subscription_expired_at <= NOW() THEN 'expired'
        ELSE t.status
    END as effective_status
FROM tenants t
LEFT JOIN (
    SELECT tenant_id, COUNT(*) as parent_count
    FROM tenant_members 
    WHERE role = 'parent' AND status = 'active'
    GROUP BY tenant_id
) parent_stats ON t.id = parent_stats.tenant_id
LEFT JOIN (
    SELECT tenant_id, COUNT(*) as child_count
    FROM tenant_members 
    WHERE role = 'child' AND status = 'active'
    GROUP BY tenant_id
) child_stats ON t.id = child_stats.tenant_id;

-- 创建函数：自动更新邀请使用次数
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_invitation_used_count()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'accepted' AND OLD.status != 'accepted' THEN
        UPDATE tenant_invitations 
        SET used_count = used_count + 1
        WHERE id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- 创建触发器：当邀请被接受时更新使用次数
CREATE TRIGGER trigger_update_invitation_count
    AFTER UPDATE OF status ON tenant_invitations
    FOR EACH ROW
    EXECUTE FUNCTION update_invitation_used_count();

-- ============================================
-- 7. 插入初始化数据
-- ============================================

-- 插入系统默认租户
INSERT INTO tenants (id, name, description, subscription_type, subscription_expired_at, max_children, status) VALUES
('00000000-0000-0000-0000-000000000001', '系统默认租户', '系统内置默认租户，用于演示和测试', 'lifetime', NOW() + INTERVAL '10 years', 5, 'active')
ON CONFLICT (id) DO NOTHING;

-- 插入系统管理员用户
INSERT INTO users (id, email, password, nickname, status) VALUES
('00000000-0000-0000-0000-000000000010', 'admin@example.com', '$2a$10$hashed_password_here', '系统管理员', 'active')
ON CONFLICT (id) DO NOTHING;

-- 将管理员用户关联到默认租户
INSERT INTO tenant_members (tenant_id, user_id, role, status, joined_at) VALUES
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000010', 'admin', 'active', NOW())
ON CONFLICT DO NOTHING;

-- 插入其他系统角色用户
INSERT INTO users (id, email, password, nickname, status) VALUES
('00000000-0000-0000-0000-000000000011', 'content.admin@example.com', '$2a$10$hashed_password_here', '内容管理员', 'active'),
('00000000-0000-0000-0000-000000000012', 'ops.admin@example.com', '$2a$10$hashed_password_here', '运营管理员', 'active')
ON CONFLICT (id) DO NOTHING;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- 删除初始化数据
DELETE FROM tenant_members WHERE user_id IN ('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000012');
DELETE FROM users WHERE id IN ('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000012');
DELETE FROM tenants WHERE id = '00000000-0000-0000-0000-000000000001';

-- 删除触发器和函数
DROP TRIGGER IF EXISTS trigger_update_invitation_count ON tenant_invitations;
DROP FUNCTION IF EXISTS update_invitation_used_count();
DROP VIEW IF EXISTS tenant_member_stats;

-- 删除表（按依赖顺序）
DROP TABLE IF EXISTS tenant_invitations;
DROP TABLE IF EXISTS tenant_members;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;