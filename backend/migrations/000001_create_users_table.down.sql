-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- 删除函数
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除索引
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_status_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- 删除约束
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_gender;
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_status;

-- 删除表
DROP TABLE IF EXISTS users;
