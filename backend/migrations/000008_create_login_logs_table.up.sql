-- 创建登录日志表
CREATE TABLE login_logs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id), -- 用户 ID: 登录的用户
    tenant_id BIGINT REFERENCES tenants(id),  -- 租户 ID: 登录到的租户，NULL 表示全局登录
    
    -- 登录信息
    login_type VARCHAR(50) DEFAULT 'password',  -- 登录方式：password/sso/oauth/magic_link
    login_status VARCHAR(50) NOT NULL,          -- 登录状态：success/failed/blocked
    
    -- 设备与环境
    ip_address INET NOT NULL,                 -- IP 地址：客户端网络地址
    user_agent TEXT,                          -- 用户代理：浏览器或应用信息
    device_type VARCHAR(50),                  -- 设备类型：desktop/mobile/tablet
    os_info VARCHAR(100),                     -- 操作系统信息：Windows/macOS/iOS/Android
    browser_info VARCHAR(100),                -- 浏览器信息：Chrome/Safari/Firefox
    
    -- 地理位置（可选）
    country VARCHAR(50),                      -- 国家：IP 归属地
    city VARCHAR(100),                        -- 城市：IP 归属地
    
    -- 安全信息
    failure_reason VARCHAR(200),              -- 失败原因：登录失败的具体原因
    is_suspicious BOOLEAN DEFAULT FALSE,      -- 是否可疑：标记异常登录行为
    risk_score INTEGER DEFAULT 0,             -- 风险评分：0-100，分数越高风险越大
    
    -- 会话信息
    session_id VARCHAR(100),                  -- 会话 ID: 登录后生成的会话标识符
    access_token_id VARCHAR(100),             -- 访问令牌 ID: JWT Token 的唯一标识
    
    occurred_at TIMESTAMP NOT NULL,           -- 发生时间：登录尝试的时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 表备注
COMMENT ON TABLE login_logs IS '登录日志表：记录所有用户登录尝试，用于安全审计、风控分析和异常检测';

-- 字段备注
COMMENT ON COLUMN login_logs.id IS 'Snowflake ID: 登录日志唯一标识';
COMMENT ON COLUMN login_logs.user_id IS '用户 ID: 尝试登录的用户，关联 users 表';
COMMENT ON COLUMN login_logs.tenant_id IS '租户 ID: 用户登录到的具体租户，NULL 表示全局登录或未选择租户';
COMMENT ON COLUMN login_logs.login_type IS '登录方式：password(密码登录)/sso(单点登录)/oauth(第三方授权)/magic_link(魔术链接免密登录)';
COMMENT ON COLUMN login_logs.login_status IS '登录状态：success(成功)/failed(失败)/blocked(被阻止)';
COMMENT ON COLUMN login_logs.ip_address IS 'IP 地址：客户端的网络 IP 地址，INET 类型支持 IPv4 和 IPv6';
COMMENT ON COLUMN login_logs.user_agent IS '用户代理：客户端的完整 User-Agent 字符串，包含浏览器、设备等信息';
COMMENT ON COLUMN login_logs.device_type IS '设备类型：desktop(桌面端)/mobile(移动端)/tablet(平板)';
COMMENT ON COLUMN login_logs.os_info IS '操作系统信息：如 Windows 11、macOS 13.2、iOS 16.3、Android 13 等';
COMMENT ON COLUMN login_logs.browser_info IS '浏览器信息：如 Chrome 120.0、Safari 17.0、Firefox 121.0 等';
COMMENT ON COLUMN login_logs.country IS '国家：基于 IP 地址解析的地理位置，如"中国"、"美国"等';
COMMENT ON COLUMN login_logs.city IS '城市：基于 IP 地址解析的地理位置，如"北京"、"上海"等';
COMMENT ON COLUMN login_logs.failure_reason IS '失败原因：登录失败时的具体原因描述，如"密码错误"、"账户已锁定"等';
COMMENT ON COLUMN login_logs.is_suspicious IS '是否可疑：标记异常的登录行为，如异地登录、频繁失败等';
COMMENT ON COLUMN login_logs.risk_score IS '风险评分：0-100 分，综合评估本次登录的风险程度，分数越高风险越大';
COMMENT ON COLUMN login_logs.session_id IS '会话 ID: 登录成功后生成的会话标识符，用于后续请求的身份验证';
COMMENT ON COLUMN login_logs.access_token_id IS '访问令牌 ID: JWT Token 的唯一标识符，用于 Token 管理和撤销';
COMMENT ON COLUMN login_logs.occurred_at IS '发生时间：登录尝试实际发生的时间戳';
COMMENT ON COLUMN login_logs.created_at IS '创建时间：日志记录插入数据库的时间戳';

-- 索引设计（高频查询场景）
CREATE INDEX idx_login_logs_user_id ON login_logs(user_id, occurred_at DESC);
CREATE INDEX idx_login_logs_tenant_id ON login_logs(tenant_id, occurred_at DESC);
CREATE INDEX idx_login_logs_status ON login_logs(login_status);
CREATE INDEX idx_login_logs_occurred_at ON login_logs(occurred_at DESC);
CREATE INDEX idx_login_logs_ip_address ON login_logs(ip_address);
CREATE INDEX idx_login_logs_suspicious ON login_logs(is_suspicious) WHERE is_suspicious = TRUE;
