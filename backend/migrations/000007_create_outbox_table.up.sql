-- Up migration
-- Outbox Pattern 表：用于保证事件发布的原子性
CREATE TABLE IF NOT EXISTS outbox (
    id BIGINT PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed BOOLEAN NOT NULL DEFAULT false,
    processed_at TIMESTAMP NULL,
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引优化
CREATE INDEX idx_outbox_unprocessed ON outbox(occurred_at) WHERE processed = false;
CREATE INDEX idx_outbox_type ON outbox(event_type);
CREATE INDEX idx_outbox_aggregate ON outbox(aggregate_type, aggregate_id);
CREATE INDEX idx_outbox_retry ON outbox(retry_count) WHERE processed = false;

-- 表备注
COMMENT ON TABLE outbox IS 'Outbox Pattern 表：用于保证事件发布的原子性，发布后标记为已处理并可定期清理';

-- 字段备注
COMMENT ON COLUMN outbox.id IS '主键 ID: 自增主键';
COMMENT ON COLUMN outbox.event_type IS '事件类型：如 user.registered, order.created 等';
COMMENT ON COLUMN outbox.aggregate_type IS '聚合根类型：如 user, order 等';
COMMENT ON COLUMN outbox.aggregate_id IS '聚合根 ID: 关联的聚合根标识';
COMMENT ON COLUMN outbox.payload IS '事件载荷：序列化后的事件数据（JSONB 格式）';
COMMENT ON COLUMN outbox.metadata IS '元数据：事件的附加信息（可选）';
COMMENT ON COLUMN outbox.occurred_at IS '事件发生时间：领域事件实际发生的时间';
COMMENT ON COLUMN outbox.processed IS '是否已处理：false-待处理，true-已发布到消息队列';
COMMENT ON COLUMN outbox.processed_at IS '处理时间：事件被发布到消息队列的时间';
COMMENT ON COLUMN outbox.error_message IS '错误信息：最后一次处理失败的错误消息';
COMMENT ON COLUMN outbox.retry_count IS '重试次数：事件处理失败的重试次数，最多 10 次';
COMMENT ON COLUMN outbox.created_at IS '创建时间：记录首次插入数据库的时间';
COMMENT ON COLUMN outbox.updated_at IS '更新时间：记录最后一次更新的时间，自动维护';
