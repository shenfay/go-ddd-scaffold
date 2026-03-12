-- 删除索引
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;

-- 删除表
DROP TABLE IF EXISTS role_permissions;
