-- 恢复 processed 字段（回滚迁移）

-- 添加字段
ALTER TABLE domain_events 
ADD COLUMN processed BOOLEAN DEFAULT FALSE;

-- 添加注释
COMMENT ON COLUMN domain_events.processed IS '是否已处理：标记事件是否已被事件处理器或消息队列消费者处理';

-- 重建索引
CREATE INDEX idx_domain_events_processed ON domain_events(processed) WHERE processed = FALSE;
