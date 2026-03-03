-- +goose Up
-- SQL section 'Up' is executed when this migration is applied

-- ============================================
-- 1. 创建租户表 (tenants) - 多租户 SaaS 基础表
-- ============================================
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    
    -- 租户配置
    max_members INTEGER NOT NULL DEFAULT 10 CHECK (max_members >= 1 AND max_members <= 100),  -- 最大成员数
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,                                              -- 过期时间
    
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加租户表注释
COMMENT ON TABLE tenants IS '租户表（多租户 SaaS 场景）';
COMMENT ON COLUMN tenants.id IS '租户唯一标识';
COMMENT ON COLUMN tenants.name IS '租户名称';
COMMENT ON COLUMN tenants.description IS '租户描述';
COMMENT ON COLUMN tenants.max_members IS '最大成员数限制';
COMMENT ON COLUMN tenants.expired_at IS '租户过期时间';
COMMENT ON COLUMN tenants.status IS '租户状态：active 活跃，expired 过期，suspended 暂停';
COMMENT ON COLUMN tenants.created_at IS '创建时间';
COMMENT ON COLUMN tenants.updated_at IS '更新时间';

-- 创建租户索引
CREATE INDEX idx_tenants_expired_at ON tenants(expired_at);
CREATE INDEX idx_tenants_status ON tenants(status);

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
    bio VARCHAR(500),
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
COMMENT ON COLUMN users.bio IS '用户个人简介';
COMMENT ON COLUMN users.status IS '用户状态：active/inactive/locked';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

-- 创建用户索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- ============================================
-- 3. 创建租户成员表 (tenant_members) - 关联用户、租户、角色
-- ============================================
CREATE TABLE tenant_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('member', 'guest', 'admin', 'owner')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'removed')),
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- 确保每个用户在同一租户中只有一个角色
    UNIQUE(tenant_id, user_id)
);

-- 添加租户成员表注释
COMMENT ON TABLE tenant_members IS '租户成员关系表';
COMMENT ON COLUMN tenant_members.id IS '成员关系唯一标识';
COMMENT ON COLUMN tenant_members.tenant_id IS '关联的租户 ID';
COMMENT ON COLUMN tenant_members.user_id IS '关联的用户 ID';
COMMENT ON COLUMN tenant_members.role IS '成员角色：member 普通成员，guest 访客，admin 管理员，owner 所有者';
COMMENT ON COLUMN tenant_members.status IS '成员状态：active 活跃，inactive 非活跃，removed 已移除';
COMMENT ON COLUMN tenant_members.invited_by IS '邀请人 ID';
COMMENT ON COLUMN tenant_members.joined_at IS '加入时间';
COMMENT ON COLUMN tenant_members.left_at IS '离开时间';
COMMENT ON COLUMN tenant_members.created_at IS '创建时间';
COMMENT ON COLUMN tenant_members.updated_at IS '更新时间';

-- 创建租户成员索引
CREATE INDEX idx_tenant_members_tenant_id ON tenant_members(tenant_id);
CREATE INDEX idx_tenant_members_user_id ON tenant_members(user_id);
CREATE INDEX idx_tenant_members_role ON tenant_members(role);
CREATE INDEX idx_tenant_members_status ON tenant_members(status);

-- ============================================
-- 4. 创建租户邀请表 (tenant_invitations)
-- ============================================
CREATE TABLE tenant_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('member', 'guest', 'admin', 'owner')),
    token VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'expired')),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    creator_id UUID REFERENCES users(id) ON DELETE SET NULL,
    accepted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 添加租户邀请表注释
COMMENT ON TABLE tenant_invitations IS '租户邀请表';
COMMENT ON COLUMN tenant_invitations.id IS '邀请唯一标识';
COMMENT ON COLUMN tenant_invitations.tenant_id IS '关联的租户 ID';
COMMENT ON COLUMN tenant_invitations.email IS '被邀请人邮箱';
COMMENT ON COLUMN tenant_invitations.role IS '邀请的角色';
COMMENT ON COLUMN tenant_invitations.token IS '邀请令牌';
COMMENT ON COLUMN tenant_invitations.status IS '邀请状态：pending 待处理，accepted 已接受，rejected 已拒绝，expired 已过期';
COMMENT ON COLUMN tenant_invitations.expires_at IS '邀请过期时间';
COMMENT ON COLUMN tenant_invitations.creator_id IS '创建人 ID';
COMMENT ON COLUMN tenant_invitations.accepted_by IS '接受人 ID';
COMMENT ON COLUMN tenant_invitations.accepted_at IS '接受时间';
COMMENT ON COLUMN tenant_invitations.rejected_at IS '拒绝时间';
COMMENT ON COLUMN tenant_invitations.created_at IS '创建时间';
COMMENT ON COLUMN tenant_invitations.updated_at IS '更新时间';

-- 创建租户邀请索引
CREATE INDEX idx_tenant_invitations_tenant_id ON tenant_invitations(tenant_id);
CREATE INDEX idx_tenant_invitations_email ON tenant_invitations(email);
CREATE INDEX idx_tenant_invitations_token ON tenant_invitations(token);
CREATE INDEX idx_tenant_invitations_status ON tenant_invitations(status);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- 删除表（注意顺序，先删除有外键依赖的表）
DROP TABLE IF EXISTS tenant_invitations;
DROP TABLE IF EXISTS tenant_members;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
