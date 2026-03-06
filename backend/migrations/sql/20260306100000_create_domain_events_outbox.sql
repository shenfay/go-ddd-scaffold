-- +goose Up
-- 创建领域事件发件箱表
CREATE TABLE IF NOT EXISTS domain_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload TEXT NOT NULL,
    
    -- 状态管理
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'published', 'failed')),
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    last_error TEXT,
    
    -- 审计字段
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 添加表注释
COMMENT ON TABLE domain_events IS '领域事件发件箱表（事务性事件存储）';
COMMENT ON COLUMN domain_events.id IS '事件唯一标识';
COMMENT ON COLUMN domain_events.event_type IS '事件类型';
COMMENT ON COLUMN domain_events.aggregate_id IS '聚合根 ID';
COMMENT ON COLUMN domain_events.payload IS '事件数据（JSON 格式）';
COMMENT ON COLUMN domain_events.status IS '事件状态：pending 待发布，processing 处理中，published 已发布，failed 失败';
COMMENT ON COLUMN domain_events.retry_count IS '重试次数';
COMMENT ON COLUMN domain_events.max_retries IS '最大重试次数';
COMMENT ON COLUMN domain_events.next_retry_at IS '下次重试时间';
COMMENT ON COLUMN domain_events.last_error IS '最后一次错误信息';
COMMENT ON COLUMN domain_events.processed_at IS '处理完成时间';
COMMENT ON COLUMN domain_events.created_at IS '创建时间';
COMMENT ON COLUMN domain_events.updated_at IS '更新时间';

-- 创建索引优化查询性能
CREATE INDEX idx_domain_events_status_type ON domain_events(status, event_type);
CREATE INDEX idx_domain_events_aggregate_id ON domain_events(aggregate_id);
CREATE INDEX idx_domain_events_created_at ON domain_events(created_at);
CREATE INDEX idx_domain_events_next_retry ON domain_events(next_retry_at) WHERE status = 'pending';

-- 创建清理函数（可选，用于定期清理旧数据）
CREATE OR REPLACE FUNCTION cleanup_old_domain_events()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- 删除 30 天前已发布的事件
    DELETE FROM domain_events
    WHERE status = 'published'
      AND created_at < NOW() - INTERVAL '30 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- +goose Down
DROP FUNCTION IF EXISTS cleanup_old_domain_events();
DROP TABLE IF EXISTS domain_events;
