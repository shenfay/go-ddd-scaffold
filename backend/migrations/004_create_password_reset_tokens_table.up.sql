-- 创建密码重置令牌表
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires ON password_reset_tokens(expires_at);

-- 添加注释
COMMENT ON TABLE password_reset_tokens IS '密码重置令牌表';
COMMENT ON COLUMN password_reset_tokens.id IS '令牌 ID (ULID 格式)';
COMMENT ON COLUMN password_reset_tokens.user_id IS '关联用户 ID';
COMMENT ON COLUMN password_reset_tokens.token IS '加密安全的随机令牌(64字符十六进制)';
COMMENT ON COLUMN password_reset_tokens.expires_at IS '令牌过期时间(1小时有效期)';
COMMENT ON COLUMN password_reset_tokens.used IS '令牌是否已使用';
COMMENT ON COLUMN password_reset_tokens.created_at IS '创建时间';
