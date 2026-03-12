-- 删除触发器
DROP TRIGGER IF EXISTS update_tenant_members_updated_at ON tenant_members;

-- 删除约束
ALTER TABLE tenant_members DROP CONSTRAINT IF EXISTS chk_tenant_members_status;

-- 删除索引
DROP INDEX IF EXISTS idx_tenant_members_metadata;
DROP INDEX IF EXISTS idx_tenant_members_deleted_at;
DROP INDEX IF EXISTS idx_tenant_members_role_id;
DROP INDEX IF EXISTS idx_tenant_members_tenant_id_status;
DROP INDEX IF EXISTS idx_tenant_members_user_id;

-- 删除表
DROP TABLE IF EXISTS tenant_members;
