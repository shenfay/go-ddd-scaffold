-- 创建审计日志表
-- 用途：记录关键安全事件（登录、认证、权限变更），满足合规性和安全审计需求
-- 保存期限：建议至少 1-7 年（根据行业合规要求）
-- 特点：不可篡改、详细完整、长期保存
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    email VARCHAR(255),                    -- 冗余字段，便于查询和展示
    action VARCHAR(50) NOT NULL,           -- 标准化命名：AUTH.*, SECURITY.*
    status VARCHAR(20) NOT NULL,           -- SUCCESS / FAILED
    ip VARCHAR(45),                        -- IPv6 最大长度
    user_agent VARCHAR(500),               -- 原始 User-Agent
    device VARCHAR(100),                   -- mobile/tablet/desktop
    browser VARCHAR(50),                   -- Chrome/Firefox/Safari
    os VARCHAR(50),                        -- Windows/macOS/Linux
    metadata JSONB DEFAULT '{}'::jsonb,    -- 结构化元数据（使用 JSONB 提升性能）
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 不添加外键约束：审计日志需独立于用户存在（合规要求）
    -- 不添加 deleted_at：审计日志不允许软删除（防篡改）
);

-- 创建索引（优化查询性能）
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at_desc ON audit_logs(created_at DESC);  -- 降序，常用查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs(status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);  -- 复合索引：用户审计轨迹
CREATE INDEX IF NOT EXISTS idx_audit_logs_action_created ON audit_logs(action, created_at DESC);  -- 复合索引：按事件类型查询
CREATE INDEX IF NOT EXISTS idx_audit_logs_failed ON audit_logs(status, created_at DESC) WHERE status = 'FAILED';  -- 部分索引：只索引失败记录

-- 添加注释
COMMENT ON TABLE audit_logs IS '审计日志表 - 记录关键安全事件（登录、认证、权限变更）';
COMMENT ON COLUMN audit_logs.id IS '日志 ID（ULID 格式）';
COMMENT ON COLUMN audit_logs.user_id IS '用户 ID（关联 users 表，但不设外键）';
COMMENT ON COLUMN audit_logs.email IS '用户邮箱（冗余字段，便于查询）';
COMMENT ON COLUMN audit_logs.action IS '审计事件类型（如：AUTH.LOGIN.SUCCESS, SECURITY.ACCOUNT.LOCKED）';
COMMENT ON COLUMN audit_logs.status IS '事件状态（SUCCESS/FAILED）';
COMMENT ON COLUMN audit_logs.ip IS 'IP 地址（支持 IPv6）';
COMMENT ON COLUMN audit_logs.user_agent IS 'User-Agent 原始字符串';
COMMENT ON COLUMN audit_logs.device IS '设备类型（mobile/tablet/desktop）';
COMMENT ON COLUMN audit_logs.browser IS '浏览器名称';
COMMENT ON COLUMN audit_logs.os IS '操作系统';
COMMENT ON COLUMN audit_logs.metadata IS '元数据（JSONB 格式，结构化存储）';
COMMENT ON COLUMN audit_logs.created_at IS '创建时间（带时区）';
