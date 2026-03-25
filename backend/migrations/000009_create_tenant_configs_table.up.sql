-- 创建租户配置表
CREATE TABLE tenant_configs (
    tenant_id BIGINT PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE, -- 租户 ID: 关联租户表，一对一关系
    
    -- 资源限制
    max_members INTEGER DEFAULT 100,          -- 最大成员数：租户允许的最大用户数量
    storage_limit BIGINT DEFAULT 10737418240, -- 存储限制（字节）：默认 10GB
    
    -- JSONB 配置字段
    feature_flags JSONB DEFAULT '{}',         -- 功能开关：JSONB 格式存储各功能的启用状态
    branding JSONB DEFAULT '{}',              -- 品牌配置：logo、主题色等白标定制信息
    custom_settings JSONB DEFAULT '{}',       -- 自定义配置：租户特定的业务配置项
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

-- 表备注
COMMENT ON TABLE tenant_configs IS '租户配置表：存储租户的个性化配置和资源限制，支持白标定制和功能开关';

-- 字段备注
COMMENT ON COLUMN tenant_configs.tenant_id IS '租户 ID: 与 tenants 表一对一关联，作为主键同时是外键，级联删除';
COMMENT ON COLUMN tenant_configs.max_members IS '最大成员数：租户允许添加的用户数量上限，受订阅计划限制';
COMMENT ON COLUMN tenant_configs.storage_limit IS '存储限制（字节）：租户可使用的总存储空间，默认 10GB(10737418240 字节)';
COMMENT ON COLUMN tenant_configs.feature_flags IS '功能开关：JSONB 格式存储各功能的启用状态，如 {"dark_mode": true, "advanced_analytics": false}';
COMMENT ON COLUMN tenant_configs.branding IS '品牌配置：JSONB 格式存储白标定制信息，如 logo URL、主题色、自定义域名等';
COMMENT ON COLUMN tenant_configs.custom_settings IS '自定义配置：JSONB 格式存储租户特定的业务配置项，如审批流程、通知模板等';
COMMENT ON COLUMN tenant_configs.created_at IS '创建时间：租户配置首次创建的时间戳';
COMMENT ON COLUMN tenant_configs.updated_at IS '更新时间：租户配置最后一次更新的时间戳，自动维护';

-- JSONB 索引（提取常用字段）
CREATE INDEX idx_tenant_configs_feature_flags ON tenant_configs USING GIN (feature_flags);
CREATE INDEX idx_tenant_configs_branding ON tenant_configs USING GIN (branding);
CREATE INDEX idx_tenant_configs_custom_settings ON tenant_configs USING GIN (custom_settings);

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_tenant_configs_updated_at 
    BEFORE UPDATE ON tenant_configs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
