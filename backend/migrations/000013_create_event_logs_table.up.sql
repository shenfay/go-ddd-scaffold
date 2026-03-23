-- 创建事件日志表（轻量级，用于事件溯源和审计）
CREATE TABLE event_logs (
    id BIGSERIAL PRIMARY KEY,                          -- 自增 ID
    aggregate_id VARCHAR(100) NOT NULL,                -- 聚合根 ID: 触发事件的领域对象 ID
    aggregate_type VARCHAR(50) NOT NULL,               -- 聚合类型：User, Tenant, Order 等
    event_type VARCHAR(100) NOT NULL,                  -- 事件类型：UserCreated, UserLoggedIn 等
    event_data JSONB NOT NULL,                         -- 事件数据：JSONB 格式存储事件的完整数据
    occurred_at TIMESTAMP NOT NULL,                    -- 事件发生时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP     -- 创建时间
);

-- 表备注
COMMENT ON TABLE event_logs IS '事件日志表：轻量级事件存储，用于审计追踪和事件回放，比 domain_events 更简洁';

-- 字段备注
COMMENT ON COLUMN event_logs.id IS '自增 ID: 事件记录的唯一标识';
COMMENT ON COLUMN event_logs.aggregate_id IS '聚合根 ID: 触发事件的领域对象的唯一标识，如用户 ID、订单 ID 等';
COMMENT ON COLUMN event_logs.aggregate_type IS '聚合类型：触发事件的领域对象类型，如 User, Tenant, Order 等';
COMMENT ON COLUMN event_logs.event_type IS '事件类型：具体的事件名称，如 UserCreated, UserLoggedIn 等';
COMMENT ON COLUMN event_logs.event_data IS '事件数据：JSONB 格式存储事件的完整数据，包含事件发生时所有的状态信息';
COMMENT ON COLUMN event_logs.occurred_at IS '事件发生时间：领域事件在业务逻辑中实际发生的时间戳';
COMMENT ON COLUMN event_logs.created_at IS '创建时间：事件记录插入数据库的时间戳';

-- 索引设计
CREATE INDEX idx_event_logs_aggregate ON event_logs(aggregate_id, aggregate_type);
CREATE INDEX idx_event_logs_type ON event_logs(event_type);
CREATE INDEX idx_event_logs_occurred ON event_logs(occurred_at DESC);

-- 复合索引：按聚合类型和时间查询（常用场景）
CREATE INDEX idx_event_logs_type_occurred ON event_logs(event_type, occurred_at DESC);
