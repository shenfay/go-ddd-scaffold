-- 删除索引
DROP INDEX IF EXISTS idx_login_logs_suspicious;
DROP INDEX IF EXISTS idx_login_logs_ip_address;
DROP INDEX IF EXISTS idx_login_logs_occurred_at;
DROP INDEX IF EXISTS idx_login_logs_status;
DROP INDEX IF EXISTS idx_login_logs_tenant_id;
DROP INDEX IF EXISTS idx_login_logs_user_id;

-- 删除表
DROP TABLE IF EXISTS login_logs;
