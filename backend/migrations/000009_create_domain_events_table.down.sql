-- 删除索引
DROP INDEX IF EXISTS idx_domain_events_metadata;
DROP INDEX IF EXISTS idx_domain_events_processed;
DROP INDEX IF EXISTS idx_domain_events_occurred;
DROP INDEX IF EXISTS idx_domain_events_type;
DROP INDEX IF EXISTS idx_domain_events_aggregate;

-- 删除表
DROP TABLE IF EXISTS domain_events;
