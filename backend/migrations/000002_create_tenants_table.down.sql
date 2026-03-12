-- 删除触发器
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;

-- 删除约束
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS chk_tenants_max_members;
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS chk_tenants_billing_cycle;
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS chk_tenants_status;

-- 删除索引
DROP INDEX IF EXISTS idx_tenants_subscription_expires_at;
DROP INDEX IF EXISTS idx_tenants_subscription_status;
DROP INDEX IF EXISTS idx_tenants_deleted_at;
DROP INDEX IF EXISTS idx_tenants_status;
DROP INDEX IF EXISTS idx_tenants_owner_id;
DROP INDEX IF EXISTS idx_tenants_code;

-- 删除表
DROP TABLE IF EXISTS tenants;
