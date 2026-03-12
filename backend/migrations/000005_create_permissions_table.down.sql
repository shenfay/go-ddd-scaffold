-- 删除约束
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS chk_permissions_scope;

-- 删除触发器
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;

-- 删除索引
DROP INDEX IF EXISTS idx_permissions_metadata;
DROP INDEX IF EXISTS idx_permissions_scope;
DROP INDEX IF EXISTS idx_permissions_resource_action;
DROP INDEX IF EXISTS idx_permissions_code;

-- 删除表
DROP TABLE IF EXISTS permissions;
