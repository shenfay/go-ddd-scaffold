// Package event 提供事务性事件发布支持
package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransactionalEventPublisher 事务性事件发布器
// 确保事件与业务操作在同一事务中提交
type TransactionalEventPublisher struct {
	outboxRepo      OutboxRepository
	eventBus        *EventBus
	logger          *zap.Logger
	publisherWorker *EventPublisherWorker
}

// TransactionalEventPublisherConfig 配置
type TransactionalEventPublisherConfig struct {
	PollInterval int  // 轮询间隔（秒）
	BatchSize    int  // 每批次处理的事件数
	MaxWorkers   int  // 并发工作协程数
	EnableWorker bool // 是否启动后台工作协程
}

// NewTransactionalEventPublisher 创建事务性事件发布器
func NewTransactionalEventPublisher(
	outboxRepo OutboxRepository,
	eventBus *EventBus,
	logger *zap.Logger,
	config TransactionalEventPublisherConfig,
) *TransactionalEventPublisher {
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.MaxWorkers == 0 {
		config.MaxWorkers = 3
	}
	if config.PollInterval == 0 {
		config.PollInterval = 5 // 默认 5 秒
	}

	publisher := &TransactionalEventPublisher{
		outboxRepo: outboxRepo,
		eventBus:   eventBus,
		logger:     logger,
	}

	// 如果启用 worker，则创建并启动
	if config.EnableWorker {
		publisher.publisherWorker = NewEventPublisherWorker(
			outboxRepo,
			eventBus,
			logger,
			config.BatchSize,
			config.MaxWorkers,
		)
	}

	return publisher
}

// PublishWithinTransaction 在事务内发布事件
// 该方法应该在事务函数内部调用，事件会与业务操作一起提交
func (p *TransactionalEventPublisher) PublishWithinTransaction(tx *gorm.DB, ctx context.Context, event DomainEvent) error {
	if event == nil {
		return ErrNilEvent
	}

	// 将事件保存到发件箱（在同一事务中）
	return p.outboxRepo.Save(ctx, tx, event)
}

// StartWorker 启动后台事件发布协程
func (p *TransactionalEventPublisher) StartWorker(ctx context.Context) error {
	if p.publisherWorker == nil {
		return fmt.Errorf("worker 未启用")
	}

	go func() {
		p.logger.Info("启动事件发布 Worker")
		if err := p.publisherWorker.Start(ctx); err != nil {
			p.logger.Error("事件发布 Worker 异常退出", zap.Error(err))
		}
	}()

	return nil
}

// StopWorker 停止后台事件发布协程
func (p *TransactionalEventPublisher) StopWorker(ctx context.Context) error {
	if p.publisherWorker == nil {
		return nil
	}

	return p.publisherWorker.Stop(ctx)
}

// EventPublisherWorker 事件发布工作协程
type EventPublisherWorker struct {
	outboxRepo OutboxRepository
	eventBus   *EventBus
	logger     *zap.Logger
	batchSize  int
	maxWorkers int
	stopChan   chan struct{}
}

// NewEventPublisherWorker 创建事件发布工作协程
func NewEventPublisherWorker(
	outboxRepo OutboxRepository,
	eventBus *EventBus,
	logger *zap.Logger,
	batchSize int,
	maxWorkers int,
) *EventPublisherWorker {
	return &EventPublisherWorker{
		outboxRepo: outboxRepo,
		eventBus:   eventBus,
		logger:     logger,
		batchSize:  batchSize,
		maxWorkers: maxWorkers,
		stopChan:   make(chan struct{}),
	}
}

// Start 启动工作协程
func (w *EventPublisherWorker) Start(ctx context.Context) error {
	ticker := newTicker(ctx, w.stopChan)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("停止事件发布 Worker")
			return nil
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.logger.Error("处理事件批次失败", zap.Error(err))
			}
		}
	}
}

// Stop 停止工作协程
func (w *EventPublisherWorker) Stop(ctx context.Context) error {
	close(w.stopChan)
	return nil
}

// processBatch 处理一批事件
func (w *EventPublisherWorker) processBatch(ctx context.Context) error {
	// 获取待处理的事件
	events, err := w.outboxRepo.GetPendingEvents(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("获取待处理事件失败：%w", err)
	}

	if len(events) == 0 {
		return nil
	}

	w.logger.Debug("获取到待处理事件", zap.Int("count", len(events)))

	// 并发处理事件（限制并发数）
	semaphore := make(chan struct{}, w.maxWorkers)
	errChan := make(chan error, len(events))

	for _, outboxEvent := range events {
		semaphore <- struct{}{} // 获取信号量

		go func(event *OutboxEvent) {
			defer func() { <-semaphore }() // 释放信号量

			if err := w.processEvent(ctx, event); err != nil {
				errChan <- err
			}
		}(outboxEvent)
	}

	// 等待所有协程完成
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}

	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("处理事件时发生 %d 个错误", len(errors))
	}

	return nil
}

// processEvent 处理单个事件
func (w *EventPublisherWorker) processEvent(ctx context.Context, outboxEvent *OutboxEvent) error {
	// 标记为处理中
	if err := w.outboxRepo.MarkAsProcessing(ctx, outboxEvent.ID); err != nil {
		return fmt.Errorf("标记事件处理中失败：%w", err)
	}

	// 反序列化事件
	domainEvent, err := w.deserializeEvent(outboxEvent)
	if err != nil {
		w.logger.Error("反序列化事件失败",
			zap.String("eventId", outboxEvent.ID.String()),
			zap.Error(err),
		)
		return w.outboxRepo.MarkAsFailed(ctx, outboxEvent.ID, err.Error())
	}

	// 发布到事件总线
	if err := w.eventBus.PublishSync(ctx, domainEvent); err != nil {
		w.logger.Error("发布事件失败",
			zap.String("eventId", outboxEvent.ID.String()),
			zap.String("eventType", outboxEvent.EventType),
			zap.Error(err),
		)
		return w.outboxRepo.MarkAsFailed(ctx, outboxEvent.ID, err.Error())
	}

	// 标记为已发布
	if err := w.outboxRepo.MarkAsPublished(ctx, outboxEvent.ID); err != nil {
		w.logger.Error("标记事件已发布失败",
			zap.String("eventId", outboxEvent.ID.String()),
			zap.Error(err),
		)
		return err
	}

	w.logger.Info("事件发布成功",
		zap.String("eventId", outboxEvent.ID.String()),
		zap.String("eventType", outboxEvent.EventType),
	)

	return nil
}

// deserializeEvent 反序列化事件
func (w *EventPublisherWorker) deserializeEvent(outboxEvent *OutboxEvent) (DomainEvent, error) {
	// 根据 EventType 动态创建具体的事件类型
	// 这里需要一个事件注册表来映射 EventType -> 构造函数
	// 简化版本：直接反序列化为 GenericDomainEvent

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(outboxEvent.Payload), &data); err != nil {
		return nil, ErrDeserializationFailed.WithCause(err)
	}

	// TODO: 实现事件类型注册表，根据 EventType 创建具体事件类型
	// 目前返回通用事件包装器
	return &GenericDomainEvent{
		ID:          outboxEvent.ID.String(),
		EventType:   outboxEvent.EventType,
		AggregateID: outboxEvent.AggregateID,
		Data:        data,
		Timestamp:   outboxEvent.CreatedAt,
		Version:     1,
	}, nil
}

// ticker 封装 time.Ticker 以支持 context 取消
type tickerImpl struct {
	ticker   *time.Ticker
	stopChan <-chan struct{}
	done     chan struct{}
	C        <-chan time.Time
}

func newTicker(ctx context.Context, stopChan <-chan struct{}) *tickerImpl {
	interval := time.Second * 5 // 默认 5 秒轮询一次
	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	c := make(chan time.Time)

	go func() {
		defer close(done)
		defer ticker.Stop()

		for {
			select {
			case tickTime := <-ticker.C:
				select {
				case c <- tickTime:
				case <-stopChan:
					return
				case <-ctx.Done():
					return
				}
			case <-stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return &tickerImpl{
		ticker:   ticker,
		stopChan: stopChan,
		done:     done,
		C:        c,
	}
}

func (t *tickerImpl) Stop() {
	<-t.done
}
