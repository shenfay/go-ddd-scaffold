-- 创建角色权限关联表
CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,   -- 角色 ID: 关联角色表
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE, -- 权限 ID: 关联权限表
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间：授权时间
    PRIMARY KEY (role_id, permission_id)          -- 复合主键：确保同一角色的同一权限只有一条记录
);

-- 表备注
COMMENT ON TABLE role_permissions IS '角色权限关联表：RBAC 模型中的多对多关系表，定义角色包含的权限';

-- 字段备注
COMMENT ON COLUMN role_permissions.role_id IS '角色 ID: 关联 roles 表，标识具体的角色';
COMMENT ON COLUMN role_permissions.permission_id IS '权限 ID: 关联 permissions 表，标识具体的权限项';
COMMENT ON COLUMN role_permissions.created_at IS '创建时间：角色被授予该权限的时间戳';

-- 索引设计
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
