-- 删除触发器
DROP TRIGGER IF EXISTS update_tenant_configs_updated_at ON tenant_configs;

-- 删除索引
DROP INDEX IF EXISTS idx_tenant_configs_custom_settings;
DROP INDEX IF EXISTS idx_tenant_configs_branding;
DROP INDEX IF EXISTS idx_tenant_configs_feature_flags;

-- 删除表
DROP TABLE IF EXISTS tenant_configs;
