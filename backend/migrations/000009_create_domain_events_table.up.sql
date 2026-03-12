-- 创建领域事件表（事件溯源）
CREATE TABLE domain_events (
    id BIGSERIAL PRIMARY KEY,                 -- 自增 ID: 事件记录的唯一标识
    aggregate_id VARCHAR(100) NOT NULL,       -- 聚合根 ID: 触发事件的领域对象 ID
    aggregate_type VARCHAR(100) NOT NULL,     -- 聚合类型：User/Tenant/Order 等
    event_type VARCHAR(100) NOT NULL,         -- 事件类型：UserRegistered/UserLoggedIn 等
    event_version INTEGER NOT NULL,           -- 事件版本：事件的版本号，用于并发控制
    event_data JSONB NOT NULL,                -- 事件数据：JSONB 格式存储事件的完整数据
    occurred_on TIMESTAMP NOT NULL,           -- 事件发生时间：领域事件实际发生的时间
    processed BOOLEAN DEFAULT FALSE,          -- 是否已处理：标记事件是否已被消费者处理
    metadata JSONB DEFAULT '{}',              -- 事件元数据：JSONB 格式存储额外上下文
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 表备注
COMMENT ON TABLE domain_events IS '领域事件表：事件溯源模式的核心表，存储所有领域对象的状态变更事件';

-- 字段备注
COMMENT ON COLUMN domain_events.id IS '自增 ID: 事件记录的唯一标识，使用 PostgreSQL 的 BIGSERIAL 类型';
COMMENT ON COLUMN domain_events.aggregate_id IS '聚合根 ID: 触发事件的领域对象的唯一标识，如用户 ID、租户 ID 等';
COMMENT ON COLUMN domain_events.aggregate_type IS '聚合类型：触发事件的领域对象类型，如 User、Tenant、Order 等';
COMMENT ON COLUMN domain_events.event_type IS '事件类型：具体的事件名称，如 UserRegisteredEvent、UserLoggedInEvent 等';
COMMENT ON COLUMN domain_events.event_version IS '事件版本：事件的版本号，用于乐观锁和事件升级管理';
COMMENT ON COLUMN domain_events.event_data IS '事件数据：JSONB 格式存储事件的完整数据，包含事件发生时所有的状态信息';
COMMENT ON COLUMN domain_events.occurred_on IS '事件发生时间：领域事件在业务逻辑中实际发生的时间戳';
COMMENT ON COLUMN domain_events.processed IS '是否已处理：标记事件是否已被事件处理器或消息队列消费者处理';
COMMENT ON COLUMN domain_events.metadata IS '事件元数据：JSONB 格式存储事件的额外上下文，如 trace_id、user_id、correlation_id 等';
COMMENT ON COLUMN domain_events.created_at IS '创建时间：事件记录插入数据库的时间戳';

-- 索引设计
CREATE INDEX idx_domain_events_aggregate ON domain_events(aggregate_id, aggregate_type);
CREATE INDEX idx_domain_events_type ON domain_events(event_type);
CREATE INDEX idx_domain_events_occurred ON domain_events(occurred_on DESC);
CREATE INDEX idx_domain_events_processed ON domain_events(processed) WHERE processed = FALSE;
CREATE INDEX idx_domain_events_metadata ON domain_events USING GIN (metadata);
