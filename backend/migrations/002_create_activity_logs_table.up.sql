-- 创建活动日志表（审计日志）
-- 用途：记录用户关键操作，满足合规性和安全审计需求
-- 保存期限：建议至少 1 年（根据合规要求调整）
CREATE TABLE IF NOT EXISTS activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    email VARCHAR(255),                    -- 冗余字段，便于查询和展示
    action VARCHAR(50) NOT NULL,           -- 标准化命名：DOMAIN.CATEGORY.ACTION
    status VARCHAR(20) NOT NULL,           -- SUCCESS / FAILED
    ip VARCHAR(45),                        -- IPv6 最大长度
    user_agent VARCHAR(500),               -- 原始 User-Agent
    device VARCHAR(100),                   -- mobile/tablet/desktop
    browser VARCHAR(50),                   -- Chrome/Firefox/Safari
    os VARCHAR(50),                        -- Windows/macOS/Linux
    description TEXT,                      -- 人类可读的描述
    metadata JSONB DEFAULT '{}'::jsonb,    -- 结构化元数据（使用 JSONB 提升性能）
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 不添加外键约束：审计日志需独立于用户存在（合规要求）
    -- 不添加 deleted_at：审计日志不允许软删除（防篡改）
);

-- 创建索引（优化查询性能）
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at_desc ON activity_logs(created_at DESC);  -- 降序，常用查询
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_status ON activity_logs(status);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_created ON activity_logs(user_id, created_at DESC);  -- 复合索引：用户时间线查询
CREATE INDEX IF NOT EXISTS idx_activity_logs_action_created ON activity_logs(action, created_at DESC);  -- 复合索引：按动作类型查询

-- GIN 索引（可选）：如果需要搜索 metadata 中的字段
-- CREATE INDEX IF NOT EXISTS idx_activity_logs_metadata_gin ON activity_logs USING GIN (metadata);

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
