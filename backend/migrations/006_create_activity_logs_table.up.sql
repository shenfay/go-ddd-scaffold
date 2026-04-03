-- 创建活动日志表（业务活动）
-- 用途：记录用户业务操作行为，用于产品分析、用户体验优化
-- 保存期限：建议 30-90 天（根据存储成本和分析需求）
-- 特点：轻量级、便于统计分析
CREATE TABLE IF NOT EXISTS activity_logs (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,           -- 标准化命名：USER.*, FEATURE.*, PAGE.*
    metadata JSONB DEFAULT '{}'::jsonb,    -- 结构化元数据
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 不需要外键约束（活动日志独立存在）
    -- 可以添加 deleted_at 支持软删除（可选）
    -- deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- 创建索引（优化查询性能）
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at_desc ON activity_logs(created_at DESC);  -- 降序，常用查询
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_created ON activity_logs(user_id, created_at DESC);  -- 复合索引：用户时间线
CREATE INDEX IF NOT EXISTS idx_activity_logs_action_created ON activity_logs(action, created_at DESC);  -- 复合索引：按动作类型统计

-- 添加注释
COMMENT ON TABLE activity_logs IS '活动日志表 - 记录用户业务操作行为（产品分析、用户体验优化）';
COMMENT ON COLUMN activity_logs.id IS '日志 ID（ULID 格式）';
COMMENT ON COLUMN activity_logs.user_id IS '用户 ID';
COMMENT ON COLUMN activity_logs.action IS '活动类型（如：USER.PROFILE.UPDATED, FEATURE.EXPORT_USED）';
COMMENT ON COLUMN activity_logs.metadata IS '元数据（JSONB 格式，结构化存储）';
COMMENT ON COLUMN activity_logs.created_at IS '创建时间（带时区）';
