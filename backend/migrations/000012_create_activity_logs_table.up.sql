-- 创建活动日志表
CREATE TABLE activity_logs (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT REFERENCES tenants(id),           -- 租户 ID: NULL 表示系统级操作
    user_id BIGINT NOT NULL REFERENCES users(id),      -- 用户 ID
    
    -- 核心字段
    action VARCHAR(100) NOT NULL,                      -- 操作类型：USER_LOGIN, USER_REGISTERED, ORDER_CREATED
    status SMALLINT DEFAULT 0,                         -- 状态：0-成功，1-失败
    
    -- 上下文信息
    ip_address INET,                                   -- IP 地址
    user_agent TEXT,                                   -- User-Agent
    metadata JSONB DEFAULT '{}',                       -- 元数据：存储设备信息、登录方式等额外数据
    
    -- 时间戳
    occurred_at TIMESTAMP NOT NULL,                    -- 发生时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP     -- 创建时间
);

-- 表备注
COMMENT ON TABLE activity_logs IS '活动日志表：统一记录所有用户活动和系统事件，用于审计、分析和安全监控';

-- 字段备注
COMMENT ON COLUMN activity_logs.id IS 'Snowflake ID: 活动日志唯一标识';
COMMENT ON COLUMN activity_logs.tenant_id IS '租户 ID: 操作所属的租户，NULL 表示系统级操作';
COMMENT ON COLUMN activity_logs.user_id IS '用户 ID: 执行操作的用户，关联 users 表';
COMMENT ON COLUMN activity_logs.action IS '操作类型：具体的活动名称，如 USER_LOGIN, USER_REGISTERED, USER_LOGOUT, ORDER_CREATED 等';
COMMENT ON COLUMN activity_logs.status IS '操作状态：0-成功，1-失败，用于统计成功率';
COMMENT ON COLUMN activity_logs.ip_address IS 'IP 地址：客户端的网络 IP 地址，INET 类型支持 IPv4 和 IPv6';
COMMENT ON COLUMN activity_logs.user_agent IS '用户代理：客户端的浏览器、移动设备或 API 客户端的完整信息';
COMMENT ON COLUMN activity_logs.metadata IS '扩展元数据：JSONB 格式存储操作的额外上下文，如登录方式、设备信息、地理位置等';
COMMENT ON COLUMN activity_logs.occurred_at IS '发生时间：操作实际发生的时间戳';
COMMENT ON COLUMN activity_logs.created_at IS '创建时间：日志记录插入数据库的时间戳';

-- 索引设计
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id, occurred_at DESC);
CREATE INDEX idx_activity_logs_tenant_id ON activity_logs(tenant_id, occurred_at DESC);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_status ON activity_logs(status);
CREATE INDEX idx_activity_logs_occurred_at ON activity_logs(occurred_at DESC);
CREATE INDEX idx_activity_logs_metadata ON activity_logs USING GIN (metadata);

-- 部分索引：只索引失败的记录（用于错误分析）
CREATE INDEX idx_activity_logs_failed ON activity_logs(status) WHERE status = 1;
