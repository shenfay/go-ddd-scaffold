-- 创建租户成员表
CREATE TABLE tenant_members (
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE, -- 租户 ID: 关联租户表
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,     -- 用户 ID: 关联用户表
    role_id BIGINT,                             -- 角色 ID: 用户在租户中的角色，关联角色表（暂时不添加外键约束）
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 加入时间：用户加入租户的时间
    invited_by BIGINT REFERENCES users(id),   -- 邀请人 ID: 邀请该用户加入的用户 ID
    status SMALLINT DEFAULT 0 NOT NULL,       -- 成员状态：0-待邀请，1-活跃，2-已移除，3-已拒绝
    metadata JSONB DEFAULT '{}',              -- 扩展信息：部门、职位等自定义字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    deleted_at TIMESTAMP,                     -- 软删除时间
    PRIMARY KEY (tenant_id, user_id)          -- 复合主键：确保同一用户在租户中只有一条记录
);

-- 表备注
COMMENT ON TABLE tenant_members IS '租户成员表：多对多关联用户和租户，记录用户在租户中的角色和状态';

-- 字段备注
COMMENT ON COLUMN tenant_members.tenant_id IS '租户 ID: 关联 tenants 表，标识成员所属的租户';
COMMENT ON COLUMN tenant_members.user_id IS '用户 ID: 关联 users 表，标识具体的成员用户';
COMMENT ON COLUMN tenant_members.role_id IS '角色 ID: 关联 roles 表，定义用户在租户中的权限等级';
COMMENT ON COLUMN tenant_members.joined_at IS '加入时间：用户正式加入租户的时间戳';
COMMENT ON COLUMN tenant_members.invited_by IS '邀请人 ID: 发起邀请的用户 ID，关联 users 表';
COMMENT ON COLUMN tenant_members.status IS '成员状态：0-待邀请，1-活跃，2-已移除，3-已拒绝';
COMMENT ON COLUMN tenant_members.metadata IS '扩展信息：JSONB 格式存储部门、职位、入职日期等自定义字段';
COMMENT ON COLUMN tenant_members.created_at IS '创建时间：记录首次创建的时间戳';
COMMENT ON COLUMN tenant_members.updated_at IS '更新时间：记录最后一次更新的时间戳，自动维护';
COMMENT ON COLUMN tenant_members.deleted_at IS '软删除时间：标记删除的时间戳，NULL 表示未删除';

-- 索引设计
CREATE INDEX idx_tenant_members_user_id ON tenant_members(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenant_members_tenant_id_status ON tenant_members(tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenant_members_role_id ON tenant_members(role_id);
CREATE INDEX idx_tenant_members_deleted_at ON tenant_members(deleted_at);
CREATE INDEX idx_tenant_members_metadata ON tenant_members USING GIN (metadata);

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_tenant_members_updated_at 
    BEFORE UPDATE ON tenant_members 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 约束检查
ALTER TABLE tenant_members ADD CONSTRAINT chk_tenant_members_status 
CHECK (status IN (0, 1, 2, 3));
