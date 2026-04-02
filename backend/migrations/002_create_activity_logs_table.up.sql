-- 创建活动日志表
CREATE TABLE IF NOT EXISTS activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    ip VARCHAR(45),
    user_agent VARCHAR(500),
    device VARCHAR(100),
    browser VARCHAR(50),
    os VARCHAR(50),
    description TEXT,
    metadata JSON,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引（优化查询性能）
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_status ON activity_logs(status);
CREATE INDEX idx_activity_logs_user_created ON activity_logs(user_id, created_at DESC);

-- 添加注释
COMMENT ON TABLE activity_logs IS '用户活动日志表';
COMMENT ON COLUMN activity_logs.id IS '日志 ID（ULID 格式）';
COMMENT ON COLUMN activity_logs.user_id IS '用户 ID（关联 users 表）';
COMMENT ON COLUMN activity_logs.email IS '用户邮箱（冗余字段，便于查询）';
COMMENT ON COLUMN activity_logs.action IS '活动类型（LOGIN/LOGOUT/REGISTER 等）';
COMMENT ON COLUMN activity_logs.status IS '活动状态（SUCCESS/FAILED）';
COMMENT ON COLUMN activity_logs.ip IS 'IP 地址（支持 IPv6）';
COMMENT ON COLUMN activity_logs.user_agent IS 'User-Agent 原始字符串';
COMMENT ON COLUMN activity_logs.device IS '设备类型（mobile/tablet/desktop）';
COMMENT ON COLUMN activity_logs.browser IS '浏览器名称';
COMMENT ON COLUMN activity_logs.os IS '操作系统';
COMMENT ON COLUMN activity_logs.description IS '活动描述';
COMMENT ON COLUMN activity_logs.metadata IS '元数据（JSON 格式）';
COMMENT ON COLUMN activity_logs.created_at IS '创建时间';
