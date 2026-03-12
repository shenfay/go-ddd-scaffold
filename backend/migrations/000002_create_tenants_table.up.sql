-- 创建租户表
CREATE TABLE tenants (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID: 租户唯一标识
    code VARCHAR(50) NOT NULL UNIQUE,         -- 租户编码：唯一标识，用于 URL 路由和数据隔离
    name VARCHAR(100) NOT NULL,               -- 租户名称：显示名称
    description TEXT,                         -- 租户描述：详细说明
    owner_id BIGINT NOT NULL REFERENCES users(id), -- 所有者 ID: 租户创建者的用户 ID
    
    -- 订阅相关字段
    subscription_plan VARCHAR(50) DEFAULT 'free',  -- 订阅计划：free/basic/pro/enterprise
    subscription_status VARCHAR(50) DEFAULT 'active', -- 订阅状态：active/canceled/expired/trial
    trial_ends_at TIMESTAMP,                    -- 试用截止时间：试用期结束时间
    subscription_starts_at TIMESTAMP,           -- 订阅开始时间：付费周期开始时间
    subscription_expires_at TIMESTAMP,          -- 订阅过期时间：当前订阅周期结束时间
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- 计费周期：monthly/yearly
    
    -- 状态与限制
    status SMALLINT DEFAULT 0 NOT NULL,       -- 状态：0-待激活，1-激活，2-停用，3-暂停
    max_members INTEGER DEFAULT 100,          -- 最大成员数：租户允许的最大用户数
    storage_limit BIGINT DEFAULT 10737418240, -- 存储限制（字节）：默认 10GB
    api_call_limit INTEGER DEFAULT 10000,     -- API 调用限制：每月 API 请求次数上限
    
    version INTEGER DEFAULT 0,                -- 乐观锁版本号：并发控制
    deleted_at TIMESTAMP,                     -- 软删除时间：标记删除时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

-- 表备注
COMMENT ON TABLE tenants IS '租户表：SaaS 多租户架构中的租户信息，支持订阅管理和资源限制';

-- 字段备注
COMMENT ON COLUMN tenants.id IS 'Snowflake ID: 租户唯一标识，使用雪花算法生成';
COMMENT ON COLUMN tenants.code IS '租户编码：全局唯一的租户标识符，用于 URL 路由（如 code.saas.com）和数据隔离';
COMMENT ON COLUMN tenants.name IS '租户名称：对外显示的组织或公司名称';
COMMENT ON COLUMN tenants.description IS '租户描述：租户的详细介绍或说明';
COMMENT ON COLUMN tenants.owner_id IS '所有者 ID: 租户创建者的用户 ID，关联 users 表';
COMMENT ON COLUMN tenants.subscription_plan IS '订阅计划：free(免费)/basic(基础)/pro(专业)/enterprise(企业)';
COMMENT ON COLUMN tenants.subscription_status IS '订阅状态：active(有效)/canceled(已取消)/expired(已过期)/trial(试用中)';
COMMENT ON COLUMN tenants.trial_ends_at IS '试用截止时间：试用期的结束时间戳，正式用户为 NULL';
COMMENT ON COLUMN tenants.subscription_starts_at IS '订阅开始时间：当前付费周期的开始时间戳';
COMMENT ON COLUMN tenants.subscription_expires_at IS '订阅过期时间：当前订阅周期的结束时间戳，永久有效为 NULL';
COMMENT ON COLUMN tenants.billing_cycle IS '计费周期：monthly(月付)/yearly(年付)，影响扣费频率';
COMMENT ON COLUMN tenants.status IS '租户状态：0-待激活，1-激活，2-停用，3-暂停';
COMMENT ON COLUMN tenants.max_members IS '最大成员数：租户允许添加的最大用户数量';
COMMENT ON COLUMN tenants.storage_limit IS '存储限制（字节）：租户可使用的总存储空间，默认 10GB';
COMMENT ON COLUMN tenants.api_call_limit IS 'API 调用限制：每月允许的 API 请求次数上限';
COMMENT ON COLUMN tenants.version IS '乐观锁版本号：用于并发更新控制，每次更新递增';
COMMENT ON COLUMN tenants.deleted_at IS '软删除时间：标记删除的时间戳，NULL 表示未删除';
COMMENT ON COLUMN tenants.created_at IS '创建时间：租户首次创建的时间戳';
COMMENT ON COLUMN tenants.updated_at IS '更新时间：租户记录最后一次更新的时间戳，自动维护';

-- 索引设计
CREATE INDEX idx_tenants_code ON tenants(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenants_owner_id ON tenants(owner_id);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_deleted_at ON tenants(deleted_at);
CREATE INDEX idx_tenants_subscription_status ON tenants(subscription_status);
CREATE INDEX idx_tenants_subscription_expires_at ON tenants(subscription_expires_at) WHERE subscription_expires_at IS NOT NULL;

-- 触发器：自动更新 updated_at
CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON tenants 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 约束检查
ALTER TABLE tenants ADD CONSTRAINT chk_tenants_status 
CHECK (status IN (0, 1, 2, 3));

ALTER TABLE tenants ADD CONSTRAINT chk_tenants_billing_cycle 
CHECK (billing_cycle IN ('monthly', 'yearly'));

ALTER TABLE tenants ADD CONSTRAINT chk_tenants_max_members 
CHECK (max_members > 0);
