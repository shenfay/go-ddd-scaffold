-- 迁移旧的 audit_logs 和 login_logs 到 activity_logs

-- 1. 迁移 audit_logs 数据
INSERT INTO activity_logs (id, tenant_id, user_id, action, status, ip_address, user_agent, metadata, occurred_at, created_at)
SELECT 
    id,
    tenant_id,
    user_id,
    action,
    status,
    ip_address,
    user_agent,
    metadata,
    occurred_at,
    COALESCE(created_at, CURRENT_TIMESTAMP)
FROM audit_logs;

-- 2. 迁移 login_logs 数据到 activity_logs
INSERT INTO activity_logs (id, tenant_id, user_id, action, status, ip_address, user_agent, metadata, occurred_at, created_at)
SELECT 
    l.id,
    l.tenant_id,
    l.user_id,
    'USER_LOGIN' as action,  -- 统一操作类型
    CASE 
        WHEN l.login_status = 'success' THEN 0
        ELSE 1
    END as status,
    l.ip_address,
    l.user_agent,
    jsonb_build_object(
        'login_type', l.login_type,
        'device_type', l.device_type,
        'os_info', l.os_info,
        'browser_info', l.browser_info,
        'country', l.country,
        'city', l.city,
        'failure_reason', l.failure_reason,
        'is_suspicious', l.is_suspicious,
        'risk_score', l.risk_score,
        'session_id', l.session_id,
        'access_token_id', l.access_token_id
    ) as metadata,
    l.occurred_at,
    COALESCE(l.created_at, CURRENT_TIMESTAMP)
FROM login_logs l;

-- 3. 验证迁移结果
DO $$
DECLARE
    old_audit_count INTEGER;
    old_login_count INTEGER;
    new_activity_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO old_audit_count FROM audit_logs;
    SELECT COUNT(*) INTO old_login_count FROM login_logs;
    SELECT COUNT(*) INTO new_activity_count FROM activity_logs;
    
    RAISE NOTICE '迁移完成:';
    RAISE NOTICE '  - 旧 audit_logs: % 条', old_audit_count;
    RAISE NOTICE '  - 旧 login_logs: % 条', old_login_count;
    RAISE NOTICE '  - 新 activity_logs: % 条', new_activity_count;
    RAISE NOTICE '  - 总计应迁移：% 条', old_audit_count + old_login_count;
    
    IF new_activity_count != (old_audit_count + old_login_count) THEN
        RAISE EXCEPTION '迁移数据量不匹配！请检查迁移脚本。';
    END IF;
END $$;
