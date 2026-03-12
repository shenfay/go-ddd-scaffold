-- 删除唯一索引
DROP INDEX IF EXISTS idx_roles_system_code_unique;
DROP INDEX IF EXISTS idx_roles_tenant_code_unique;

-- 删除约束
ALTER TABLE roles DROP CONSTRAINT IF EXISTS chk_roles_tenant_system;

-- 删除触发器
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;

-- 删除索引
DROP INDEX IF EXISTS idx_roles_permissions;
DROP INDEX IF EXISTS idx_roles_deleted_at;
DROP INDEX IF EXISTS idx_roles_is_system;
DROP INDEX IF EXISTS idx_roles_code;
DROP INDEX IF EXISTS idx_roles_tenant_id;

-- 删除表
DROP TABLE IF EXISTS roles;
