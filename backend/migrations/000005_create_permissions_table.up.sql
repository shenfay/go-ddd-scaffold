-- 创建权限表
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID: 权限唯一标识
    name VARCHAR(100) NOT NULL,               -- 权限名称：显示名称
    code VARCHAR(100) NOT NULL UNIQUE,        -- 权限编码：全局唯一的权限标识符
    resource VARCHAR(100) NOT NULL,           -- 资源类型：user/tenant/billing 等
    action VARCHAR(50) NOT NULL,              -- 操作类型：create/read/update/delete/*
    scope VARCHAR(50) DEFAULT 'tenant',       -- 作用域：system-系统级，tenant-租户级
    description TEXT,                         -- 权限描述：详细说明
    metadata JSONB DEFAULT '{}',              -- 扩展元数据：JSONB 格式存储额外信息
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

-- 表备注
COMMENT ON TABLE permissions IS '权限定义表：RBAC 模型中的权限字典，定义系统中所有可授予的权限';

-- 字段备注
COMMENT ON COLUMN permissions.id IS 'Snowflake ID: 权限唯一标识，使用雪花算法生成';
COMMENT ON COLUMN permissions.name IS '权限名称：权限的显示名称，如"创建用户"、"查看账单"等';
COMMENT ON COLUMN permissions.code IS '权限编码：全局唯一的权限标识符，用于代码中判断权限，如"user:create"';
COMMENT ON COLUMN permissions.resource IS '资源类型：权限作用的目标资源，如 user/tenant/billing/order 等';
COMMENT ON COLUMN permissions.action IS '操作类型：对资源的操作，如 create/read/update/delete，*表示所有操作';
COMMENT ON COLUMN permissions.scope IS '作用域：system-系统级权限（仅管理员可用），tenant-租户级权限（普通用户可用）';
COMMENT ON COLUMN permissions.description IS '权限描述：权限的详细介绍或使用场景说明';
COMMENT ON COLUMN permissions.metadata IS '扩展元数据：JSONB 格式存储权限的额外配置或约束条件';
COMMENT ON COLUMN permissions.created_at IS '创建时间：权限首次创建的时间戳';
COMMENT ON COLUMN permissions.updated_at IS '更新时间：权限记录最后一次更新的时间戳，自动维护';

-- 索引设计
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_metadata ON permissions USING GIN (metadata);

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_permissions_updated_at 
    BEFORE UPDATE ON permissions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 约束检查
ALTER TABLE permissions ADD CONSTRAINT chk_permissions_scope 
CHECK (scope IN ('system', 'tenant'));
