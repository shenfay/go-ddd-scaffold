package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// PublishDomainEventJob 发布领域事件作业
// 职责：从 outbox_events 表读取事件并发布到 EventBus
type PublishDomainEventJob struct {
	logger *zap.Logger
}

// NewPublishDomainEventJob 创建发布领域事件作业
func NewPublishDomainEventJob(logger *zap.Logger) *PublishDomainEventJob {
	return &PublishDomainEventJob{
		logger: logger.Named("publish_domain_event"),
	}
}

// Execute 执行发布领域事件任务
func (j *PublishDomainEventJob) Execute(ctx context.Context, payload map[string]interface{}) error {
	j.logger.Debug("Processing domain event publication", zap.Any("payload", payload))

	// TODO: 实现 Outbox 模式的事件发布逻辑
	// 1. 从 outbox_events 表查询待发布的事件
	// 2. 反序列化事件数据
	// 3. 通过 EventBus 发布事件
	// 4. 删除或标记已处理的 outbox 记录

	j.logger.Info("Domain event processed (placeholder)",
		zap.String("event_type", payload["event_type"].(string)))

	return nil
}

// Queue 返回队列名称
func (j *PublishDomainEventJob) Queue() string {
	return "jobs_realtime"
}

// MaxRetry 返回最大重试次数
func (j *PublishDomainEventJob) MaxRetry() int {
	return 3
}

// Timeout 返回超时时间
func (j *PublishDomainEventJob) Timeout() time.Duration {
	return 30 * time.Second
}
