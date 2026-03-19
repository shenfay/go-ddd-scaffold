-- 移除 domain_events 表的 processed 字段
-- 原因：事件存储仅用于溯源和审计，不跟踪处理状态
-- 状态跟踪由 asynq 队列独立管理

-- 删除索引
DROP INDEX IF EXISTS idx_domain_events_processed;

-- 删除字段
ALTER TABLE domain_events 
DROP COLUMN IF EXISTS processed;

-- 更新注释
COMMENT ON TABLE domain_events IS '领域事件表：纯事件溯源模式，仅记录历史事件用于审计和回放，不包含状态信息';
