-- Up migration
-- Domain Events 表：用于事件溯源和审计追踪（永久保存）
CREATE TABLE IF NOT EXISTS domain_events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    metadata JSONB,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引优化
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_type, aggregate_id);
CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_occurred ON domain_events(occurred_at DESC);

-- 表备注
COMMENT ON TABLE domain_events IS '领域事件表：用于事件溯源和审计追踪，永久保存';

-- 字段备注
COMMENT ON COLUMN domain_events.id IS '主键 ID: 自增主键';
COMMENT ON COLUMN domain_events.aggregate_id IS '聚合根 ID: 关联的聚合根标识';
COMMENT ON COLUMN domain_events.aggregate_type IS '聚合根类型：如 user, order 等';
COMMENT ON COLUMN domain_events.event_type IS '事件类型：如 user.registered, order.created 等';
COMMENT ON COLUMN domain_events.event_data IS '事件数据：序列化后的事件数据（JSONB 格式）';
COMMENT ON COLUMN domain_events.metadata IS '元数据：事件的附加信息（可选）';
COMMENT ON COLUMN domain_events.occurred_at IS '事件发生时间：领域事件实际发生的时间';
COMMENT ON COLUMN domain_events.created_at IS '创建时间：记录插入数据库的时间';
