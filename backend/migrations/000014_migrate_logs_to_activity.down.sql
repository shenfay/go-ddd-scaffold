-- 回滚迁移（如果需要）
-- 注意：这个脚本不会删除 activity_logs 的数据，只是重新启用旧表

-- 如果需要恢复旧表结构，执行以下 SQL：
-- 1. 从 activity_logs 恢复 login_logs
-- INSERT INTO login_logs (...) SELECT ... FROM activity_logs WHERE action = 'USER_LOGIN';

-- 2. 从 activity_logs 恢复 audit_logs  
-- INSERT INTO audit_logs (...) SELECT ... FROM activity_logs WHERE action != 'USER_LOGIN';

-- 警告：这是一个危险操作，请谨慎使用！
SELECT 'Manual rollback required. Check migration script for details.' as warning;
