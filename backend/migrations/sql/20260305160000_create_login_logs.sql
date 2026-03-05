-- +goose Up
-- +goose StatementBegin

-- 创建登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    
    -- IP 地址信息
    ip_address      VARCHAR(45) NOT NULL,
    ip_country      VARCHAR(50),
    ip_region       VARCHAR(50),
    ip_city         VARCHAR(50),
    
    -- 设备信息
    device_type     VARCHAR(20) DEFAULT 'desktop',
    os_info         VARCHAR(100),
    browser_info    VARCHAR(100),
    
    -- 登录状态
    login_status    VARCHAR(20) NOT NULL DEFAULT 'success',
    failure_reason  VARCHAR(200),
    
    -- 时间戳
    logged_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束
    CONSTRAINT fk_login_logs_user 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建索引

-- 复合索引：用户查询最近登录（最高频查询）
CREATE INDEX idx_login_logs_user_time ON login_logs(user_id, logged_at DESC);

-- 安全审计索引
CREATE INDEX idx_login_logs_ip_address ON login_logs(ip_address);
CREATE INDEX idx_login_logs_status_time ON login_logs(login_status, logged_at);
CREATE INDEX idx_login_logs_device_type ON login_logs(device_type);

-- 时间范围查询索引
CREATE INDEX idx_login_logs_logged_at ON login_logs(logged_at);

-- 注释说明
COMMENT ON TABLE login_logs IS '用户登录日志表 - 记录所有登录尝试（成功/失败）';
COMMENT ON COLUMN login_logs.user_id IS '用户 ID，关联 users 表';
COMMENT ON COLUMN login_logs.ip_address IS '登录 IP 地址（支持 IPv6）';
COMMENT ON COLUMN login_logs.ip_country IS 'IP 归属国家（通过 IP 库解析）';
COMMENT ON COLUMN login_logs.ip_region IS 'IP 归属省份/州';
COMMENT ON COLUMN login_logs.ip_city IS 'IP 归属城市';
COMMENT ON COLUMN login_logs.device_type IS '设备类型：mobile/desktop/tablet';
COMMENT ON COLUMN login_logs.os_info IS '操作系统信息（如：Windows 11, macOS 14.2）';
COMMENT ON COLUMN login_logs.browser_info IS '浏览器信息（如：Chrome 120.0, Safari 17.2）';
COMMENT ON COLUMN login_logs.login_status IS '登录状态：success/failed';
COMMENT ON COLUMN login_logs.failure_reason IS '失败原因：wrong_password/account_locked/too_many_attempts 等';
COMMENT ON COLUMN login_logs.logged_at IS '登录发生时间';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 删除表和索引
DROP INDEX IF EXISTS idx_login_logs_logged_at;
DROP INDEX IF EXISTS idx_login_logs_device_type;
DROP INDEX IF EXISTS idx_login_logs_status_time;
DROP INDEX IF EXISTS idx_login_logs_ip_address;
DROP INDEX IF EXISTS idx_login_logs_user_time;

DROP TABLE IF EXISTS login_logs;

-- +goose StatementEnd
