-- 创建角色表
CREATE TABLE roles (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID: 角色唯一标识
    tenant_id BIGINT,                         -- 租户 ID: NULL 表示系统角色，有值表示租户自定义角色
    name VARCHAR(50) NOT NULL,                -- 角色名称：显示名称
    code VARCHAR(50) NOT NULL,                -- 角色编码：唯一标识符
    description TEXT,                         -- 角色描述：详细说明
    is_system BOOLEAN DEFAULT FALSE,          -- 是否系统角色：TRUE=系统预定义，FALSE=租户自定义
    permissions JSONB DEFAULT '[]',           -- 权限列表：JSONB 数组存储权限码
    version INTEGER DEFAULT 0,                -- 乐观锁版本号：并发控制
    deleted_at TIMESTAMP,                     -- 软删除时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

-- 表备注
COMMENT ON TABLE roles IS '角色表：RBAC 权限模型中的角色定义，支持系统角色和租户自定义角色';

-- 字段备注
COMMENT ON COLUMN roles.id IS 'Snowflake ID: 角色唯一标识，使用雪花算法生成';
COMMENT ON COLUMN roles.tenant_id IS '租户 ID: NULL 表示系统预定义角色（全局可用），有值表示租户自定义角色（仅该租户可用）';
COMMENT ON COLUMN roles.name IS '角色名称：角色的显示名称，如"管理员"、"普通成员"等';
COMMENT ON COLUMN roles.code IS '角色编码：角色的唯一标识符，如"admin"、"member"等，同一租户下唯一';
COMMENT ON COLUMN roles.description IS '角色描述：角色的详细介绍或职责说明';
COMMENT ON COLUMN roles.is_system IS '是否系统角色：TRUE=系统预定义角色（tenant_id 必须为 NULL），FALSE=租户自定义角色（tenant_id 必须有值）';
COMMENT ON COLUMN roles.permissions IS '权限列表：JSONB 数组格式存储权限码，如 ["user:create", "user:read"]';
COMMENT ON COLUMN roles.version IS '乐观锁版本号：用于并发更新控制，每次更新递增';
COMMENT ON COLUMN roles.deleted_at IS '软删除时间：标记删除的时间戳，NULL 表示未删除';
COMMENT ON COLUMN roles.created_at IS '创建时间：角色首次创建的时间戳';
COMMENT ON COLUMN roles.updated_at IS '更新时间：角色记录最后一次更新的时间戳，自动维护';

-- 索引设计
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_roles_code ON roles(code, tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_is_system ON roles(is_system);
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX idx_roles_permissions ON roles USING GIN (permissions);

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_roles_updated_at 
    BEFORE UPDATE ON roles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 约束检查：确保逻辑一致性
ALTER TABLE roles ADD CONSTRAINT chk_roles_tenant_system 
CHECK (
    (is_system = TRUE AND tenant_id IS NULL) OR  -- 系统角色必须 tenant_id 为 NULL
    (is_system = FALSE AND tenant_id IS NOT NULL) -- 租户角色必须有 tenant_id
);

-- 唯一约束：同一租户下角色编码唯一
CREATE UNIQUE INDEX idx_roles_tenant_code_unique ON roles(tenant_id, code) WHERE deleted_at IS NULL AND tenant_id IS NOT NULL;
CREATE UNIQUE INDEX idx_roles_system_code_unique ON roles(code) WHERE deleted_at IS NULL AND is_system = TRUE;
