-- 创建审计日志表
CREATE TABLE audit_logs (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT REFERENCES tenants(id),  -- 租户 ID: NULL 表示系统级操作
    user_id BIGINT REFERENCES users(id),      -- 用户 ID: 执行操作的用户
    action VARCHAR(100) NOT NULL,             -- 操作类型：如 CREATE_USER、UPDATE_TENANT
    resource_type VARCHAR(100),               -- 资源类型：如 User、Tenant、Order
    resource_id BIGINT,                       -- 资源 ID: 被操作资源的 ID
    request_id VARCHAR(100),                  -- 请求 ID: 链路追踪标识符
    ip_address INET,                          -- IP 地址：客户端 IP 地址
    user_agent TEXT,                          -- 用户代理：客户端浏览器或应用信息
    metadata JSONB DEFAULT '{}',              -- 扩展元数据：JSONB 格式存储额外上下文
    status SMALLINT DEFAULT 0,                -- 状态：0-成功，1-失败
    error_message TEXT,                       -- 错误信息：操作失败时的错误详情
    occurred_at TIMESTAMP NOT NULL,           -- 发生时间：操作实际发生的时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 表备注
COMMENT ON TABLE audit_logs IS '审计日志表：记录系统中所有重要操作，用于安全审计、问题排查和合规检查';

-- 字段备注
COMMENT ON COLUMN audit_logs.id IS 'Snowflake ID: 审计日志唯一标识';
COMMENT ON COLUMN audit_logs.tenant_id IS '租户 ID: 操作所属的租户，NULL 表示系统级操作（如管理员操作）';
COMMENT ON COLUMN audit_logs.user_id IS '用户 ID: 执行操作的用户，关联 users 表';
COMMENT ON COLUMN audit_logs.action IS '操作类型：具体的操作名称，如 CREATE_USER、UPDATE_TENANT_CONFIG、DELETE_ORDER 等';
COMMENT ON COLUMN audit_logs.resource_type IS '资源类型：被操作的目标资源类型，如 User、Tenant、Order 等';
COMMENT ON COLUMN audit_logs.resource_id IS '资源 ID: 被操作资源的具体 ID，用于定位具体记录';
COMMENT ON COLUMN audit_logs.request_id IS '请求 ID: 分布式链路追踪的唯一标识符，用于关联同一次请求的所有日志';
COMMENT ON COLUMN audit_logs.ip_address IS 'IP 地址：客户端的网络 IP 地址，INET 类型支持 IPv4 和 IPv6';
COMMENT ON COLUMN audit_logs.user_agent IS '用户代理：客户端的浏览器、移动设备或 API 客户端的完整信息';
COMMENT ON COLUMN audit_logs.metadata IS '扩展元数据：JSONB 格式存储操作的额外上下文，如变更前后对比、操作参数等';
COMMENT ON COLUMN audit_logs.status IS '操作状态：0-成功，1-失败，用于统计操作成功率';
COMMENT ON COLUMN audit_logs.error_message IS '错误信息：当操作失败时，记录详细的错误消息或堆栈跟踪';
COMMENT ON COLUMN audit_logs.occurred_at IS '发生时间：操作实际发生的时间戳，可能与 created_at 略有不同';
COMMENT ON COLUMN audit_logs.created_at IS '创建时间：日志记录插入数据库的时间戳';

-- 索引设计
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_occurred_at ON audit_logs(occurred_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_status ON audit_logs(status);
CREATE INDEX idx_audit_logs_metadata ON audit_logs USING GIN (metadata);
