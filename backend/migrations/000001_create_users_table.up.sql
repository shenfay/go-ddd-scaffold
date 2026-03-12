-- 创建用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY,                    -- Snowflake ID: 用户唯一标识
    username VARCHAR(50) NOT NULL UNIQUE,     -- 用户名：用户登录名，全局唯一
    email VARCHAR(255) NOT NULL UNIQUE,       -- 邮箱：用户邮箱地址，全局唯一
    password_hash VARCHAR(255) NOT NULL,      -- 密码哈希：bcrypt 加密后的密码
    status SMALLINT DEFAULT 0 NOT NULL,       -- 状态：0-待激活，1-激活，2-禁用，3-锁定
    display_name VARCHAR(100),                -- 显示名称：用户对外显示的名称
    gender SMALLINT DEFAULT 0,                -- 性别：0-未知，1-男，2-女，3-其他
    phone_number VARCHAR(20),                 -- 电话号码：用户手机号码
    avatar_url VARCHAR(500),                  -- 头像 URL: 用户头像图片地址
    last_login_at TIMESTAMP,                  -- 最后登录时间：用户最后一次成功登录时间
    login_count INTEGER DEFAULT 0,            -- 登录次数：用户累计登录次数
    failed_attempts INTEGER DEFAULT 0,        -- 连续失败登录次数：用于账户锁定策略
    locked_until TIMESTAMP,                   -- 账户锁定截止时间：临时锁定的结束时间
    version INTEGER DEFAULT 0,                -- 乐观锁版本号：并发控制
    deleted_at TIMESTAMP,                     -- 软删除时间：标记删除时间，NULL 表示未删除
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间：记录创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 更新时间：记录最后更新时间
);

-- 表备注
COMMENT ON TABLE users IS '用户表：存储系统用户的基本信息和认证信息';

-- 字段备注
COMMENT ON COLUMN users.id IS 'Snowflake ID: 用户唯一标识，使用雪花算法生成';
COMMENT ON COLUMN users.username IS '用户名：用户登录名，全局唯一，长度 3-50 字符';
COMMENT ON COLUMN users.email IS '邮箱：用户邮箱地址，用于登录和接收通知，全局唯一';
COMMENT ON COLUMN users.password_hash IS '密码哈希：使用 bcrypt 算法加密后的密码';
COMMENT ON COLUMN users.status IS '用户状态：0-待激活，1-激活，2-禁用，3-锁定';
COMMENT ON COLUMN users.display_name IS '显示名称：用户对外显示的名称，可为空';
COMMENT ON COLUMN users.gender IS '性别：0-未知，1-男，2-女，3-其他';
COMMENT ON COLUMN users.phone_number IS '电话号码：用户手机号码，用于短信验证等';
COMMENT ON COLUMN users.avatar_url IS '头像 URL: 用户头像图片的存储地址';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间：用户最后一次成功登录的时间戳';
COMMENT ON COLUMN users.login_count IS '登录次数：用户累计成功登录的次数';
COMMENT ON COLUMN users.failed_attempts IS '连续失败登录次数：用于账户安全锁定策略';
COMMENT ON COLUMN users.locked_until IS '账户锁定截止时间：临时锁定的结束时间，永久锁定为 NULL';
COMMENT ON COLUMN users.version IS '乐观锁版本号：用于并发更新控制，每次更新递增';
COMMENT ON COLUMN users.deleted_at IS '软删除时间：标记删除的时间戳，NULL 表示未删除';
COMMENT ON COLUMN users.created_at IS '创建时间：记录首次创建的时间戳';
COMMENT ON COLUMN users.updated_at IS '更新时间：记录最后一次更新的时间戳，自动维护';

-- 索引设计
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_status_deleted_at ON users(status, deleted_at);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- 触发器：自动更新 updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 约束检查
ALTER TABLE users ADD CONSTRAINT chk_users_status 
CHECK (status IN (0, 1, 2, 3));

ALTER TABLE users ADD CONSTRAINT chk_users_gender 
CHECK (gender IN (0, 1, 2, 3));
