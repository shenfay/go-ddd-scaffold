-- 删除索引
DROP INDEX IF EXISTS idx_audit_logs_metadata;
DROP INDEX IF EXISTS idx_audit_logs_status;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_occurred_at;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_tenant_id;

-- 删除表
DROP TABLE IF EXISTS audit_logs;
