package monitoring

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// EventMetricsCollector 事件指标收集器
type EventMetricsCollector struct {
	// 内置指标
	eventsPublished     map[string]*atomic.Uint64 // eventType -> count
	eventsProcessed     map[string]*atomic.Uint64 // eventType -> count
	eventsFailed        map[string]*atomic.Uint64 // eventType -> count
	eventProcessingTime map[string]*atomic.Int64  // eventType -> avg nanoseconds
	queueDepth          map[string]*atomic.Int64  // queue -> depth
	retryCount          map[string]*atomic.Uint64 // eventType -> count

	// 内部状态
	redisClient *redis.Client
	logger      *zap.Logger
	config      MetricsConfig

	// 告警相关
	alertRules []AlertRule
	alertCh    chan Alert
	mu         sync.RWMutex
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Namespace      string
	Subsystem      string
	RedisAddr      string
	ScrapeInterval time.Duration
	EnableAlerts   bool
}

// AlertRule 告警规则
type AlertRule struct {
	Name        string
	Description string
	Metric      string  // 指标名称
	Threshold   float64 // 阈值
	Duration    time.Duration
	Severity    AlertSeverity
	Enabled     bool
}

// AlertSeverity 告警严重程度
type AlertSeverity string

const (
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityInfo     AlertSeverity = "info"
)

// Alert 告警信息
type Alert struct {
	Name        string
	Description string
	Severity    AlertSeverity
	Value       float64
	Threshold   float64
	Timestamp   time.Time
}

// EventMetrics 事件指标数据
type EventMetrics struct {
	EventType     string
	QueueName     string
	Published     uint64
	Processed     uint64
	Failed        uint64
	AverageTime   time.Duration
	QueueDepth    int64
	RetryCount    uint64
	LastProcessed time.Time
}

// NewEventMetricsCollector 创建事件指标收集器
func NewEventMetricsCollector(
	redisClient *redis.Client,
	logger *zap.Logger,
	config MetricsConfig,
) *EventMetricsCollector {
	if logger == nil {
		logger = zap.L().Named("event_metrics")
	}

	// 设置默认配置
	if config.Namespace == "" {
		config.Namespace = "ddd_scaffold"
	}
	if config.Subsystem == "" {
		config.Subsystem = "events"
	}
	if config.ScrapeInterval == 0 {
		config.ScrapeInterval = 30 * time.Second
	}

	collector := &EventMetricsCollector{
		eventsPublished:     make(map[string]*atomic.Uint64),
		eventsProcessed:     make(map[string]*atomic.Uint64),
		eventsFailed:        make(map[string]*atomic.Uint64),
		eventProcessingTime: make(map[string]*atomic.Int64),
		queueDepth:          make(map[string]*atomic.Int64),
		retryCount:          make(map[string]*atomic.Uint64),
		redisClient:         redisClient,
		logger:              logger,
		config:              config,
		alertCh:             make(chan Alert, 100),
	}

	// 初始化默认告警规则
	collector.initDefaultAlertRules()

	return collector
}

// 初始化默认告警规则
func (emc *EventMetricsCollector) initDefaultAlertRules() {
	emc.alertRules = []AlertRule{
		{
			Name:        "high_failure_rate",
			Description: "Event failure rate exceeds 10%",
			Metric:      "failure_rate",
			Threshold:   0.1,
			Duration:    5 * time.Minute,
			Severity:    AlertSeverityWarning,
			Enabled:     true,
		},
		{
			Name:        "queue_backlog",
			Description: "Queue depth exceeds 1000",
			Metric:      "queue_depth",
			Threshold:   1000,
			Duration:    10 * time.Minute,
			Severity:    AlertSeverityWarning,
			Enabled:     true,
		},
		{
			Name:        "processing_slow",
			Description: "Average processing time exceeds 5 seconds",
			Metric:      "processing_time",
			Threshold:   5.0,
			Duration:    5 * time.Minute,
			Severity:    AlertSeverityCritical,
			Enabled:     true,
		},
	}
}

// RecordEventPublished 记录事件发布
func (emc *EventMetricsCollector) RecordEventPublished(eventType, queue string) {
	key := fmt.Sprintf("%s_%s", eventType, queue)

	if _, exists := emc.eventsPublished[key]; !exists {
		emc.mu.Lock()
		if _, exists := emc.eventsPublished[key]; !exists {
			emc.eventsPublished[key] = &atomic.Uint64{}
		}
		emc.mu.Unlock()
	}

	emc.eventsPublished[key].Add(1)
}

// RecordEventProcessed 记录事件处理成功
func (emc *EventMetricsCollector) RecordEventProcessed(eventType, queue, handler string, duration time.Duration) {
	key := fmt.Sprintf("%s_%s", eventType, queue)

	// 增加处理计数
	if _, exists := emc.eventsProcessed[key]; !exists {
		emc.mu.Lock()
		if _, exists := emc.eventsProcessed[key]; !exists {
			emc.eventsProcessed[key] = &atomic.Uint64{}
		}
		emc.mu.Unlock()
	}
	emc.eventsProcessed[key].Add(1)

	// 更新平均处理时间（简单移动平均）
	if _, exists := emc.eventProcessingTime[key]; !exists {
		emc.mu.Lock()
		if _, exists := emc.eventProcessingTime[key]; !exists {
			emc.eventProcessingTime[key] = &atomic.Int64{}
		}
		emc.mu.Unlock()
	}

	currentAvg := emc.eventProcessingTime[key].Load()
	newDuration := duration.Nanoseconds()
	// 简单移动平均：new_avg = (old_avg + new_value) / 2
	newAvg := (currentAvg + newDuration) / 2
	emc.eventProcessingTime[key].Store(newAvg)
}

// RecordEventFailed 记录事件处理失败
func (emc *EventMetricsCollector) RecordEventFailed(eventType, queue, handler, errorType string) {
	key := fmt.Sprintf("%s_%s", eventType, queue)

	if _, exists := emc.eventsFailed[key]; !exists {
		emc.mu.Lock()
		if _, exists := emc.eventsFailed[key]; !exists {
			emc.eventsFailed[key] = &atomic.Uint64{}
		}
		emc.mu.Unlock()
	}

	emc.eventsFailed[key].Add(1)
}

// RecordRetry 记录重试次数
func (emc *EventMetricsCollector) RecordRetry(eventType, queue string) {
	key := fmt.Sprintf("%s_%s", eventType, queue)

	if _, exists := emc.retryCount[key]; !exists {
		emc.mu.Lock()
		if _, exists := emc.retryCount[key]; !exists {
			emc.retryCount[key] = &atomic.Uint64{}
		}
		emc.mu.Unlock()
	}

	emc.retryCount[key].Add(1)
}

// UpdateQueueDepth 更新队列深度
func (emc *EventMetricsCollector) UpdateQueueDepth(queue string, depth int64) {
	if _, exists := emc.queueDepth[queue]; !exists {
		emc.mu.Lock()
		if _, exists := emc.queueDepth[queue]; !exists {
			emc.queueDepth[queue] = &atomic.Int64{}
		}
		emc.mu.Unlock()
	}

	emc.queueDepth[queue].Store(depth)
}

// GetMetrics 获取当前指标数据
func (emc *EventMetricsCollector) GetMetrics(ctx context.Context) ([]EventMetrics, error) {
	emc.mu.RLock()
	defer emc.mu.RUnlock()

	var metrics []EventMetrics

	// 收集所有指标数据
	for key := range emc.eventsPublished {
		// 解析 key
		var eventType, queue string
		if n, err := fmt.Sscanf(key, "%s_%s", &eventType, &queue); err != nil || n != 2 {
			continue
		}

		metric := EventMetrics{
			EventType:     eventType,
			QueueName:     queue,
			Published:     emc.eventsPublished[key].Load(),
			Processed:     0,
			Failed:        0,
			AverageTime:   time.Duration(emc.eventProcessingTime[key].Load()) * time.Nanosecond,
			QueueDepth:    0,
			RetryCount:    0,
			LastProcessed: time.Now(),
		}

		// 获取其他指标
		if processed, exists := emc.eventsProcessed[key]; exists {
			metric.Processed = processed.Load()
		}
		if failed, exists := emc.eventsFailed[key]; exists {
			metric.Failed = failed.Load()
		}
		if retry, exists := emc.retryCount[key]; exists {
			metric.RetryCount = retry.Load()
		}

		metrics = append(metrics, metric)
	}

	// 添加队列深度信息
	for queue, depth := range emc.queueDepth {
		// 查找对应的事件类型（简化处理）
		found := false
		for i := range metrics {
			if metrics[i].QueueName == queue {
				metrics[i].QueueDepth = depth.Load()
				found = true
				break
			}
		}
		if !found {
			// 创建新的队列指标
			metrics = append(metrics, EventMetrics{
				QueueName:  queue,
				QueueDepth: depth.Load(),
			})
		}
	}

	return metrics, nil
}

// StartMonitoring 开始监控
func (emc *EventMetricsCollector) StartMonitoring(ctx context.Context) {
	emc.logger.Info("Starting event monitoring")

	ticker := time.NewTicker(emc.config.ScrapeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			emc.logger.Info("Stopping event monitoring")
			return
		case <-ticker.C:
			emc.scrapeMetrics(ctx)
			if emc.config.EnableAlerts {
				emc.evaluateAlerts(ctx)
			}
		}
	}
}

// scrapeMetrics 抓取指标数据
func (emc *EventMetricsCollector) scrapeMetrics(ctx context.Context) {
	metrics, err := emc.GetMetrics(ctx)
	if err != nil {
		emc.logger.Error("Failed to scrape metrics", zap.Error(err))
		return
	}

	emc.logger.Debug("Metrics scraped", zap.Int("count", len(metrics)))
}

// evaluateAlerts 评估告警规则
func (emc *EventMetricsCollector) evaluateAlerts(ctx context.Context) {
	metrics, err := emc.GetMetrics(ctx)
	if err != nil {
		emc.logger.Error("Failed to get metrics for alerts", zap.Error(err))
		return
	}

	for _, rule := range emc.alertRules {
		if !rule.Enabled {
			continue
		}

		value := emc.calculateMetricValue(rule.Metric, metrics)
		if value >= rule.Threshold {
			alert := Alert{
				Name:        rule.Name,
				Description: rule.Description,
				Severity:    rule.Severity,
				Value:       value,
				Threshold:   rule.Threshold,
				Timestamp:   time.Now(),
			}

			select {
			case emc.alertCh <- alert:
				emc.logger.Warn("Alert triggered",
					zap.String("alert_name", alert.Name),
					zap.Float64("value", alert.Value),
					zap.Float64("threshold", alert.Threshold))
			default:
				emc.logger.Warn("Alert channel full, dropping alert", zap.String("alert_name", rule.Name))
			}
		}
	}
}

// calculateMetricValue 计算指标值
func (emc *EventMetricsCollector) calculateMetricValue(metric string, metrics []EventMetrics) float64 {
	switch metric {
	case "failure_rate":
		var totalProcessed, totalFailed uint64
		for _, m := range metrics {
			totalProcessed += m.Processed
			totalFailed += m.Failed
		}
		if totalProcessed == 0 {
			return 0
		}
		return float64(totalFailed) / float64(totalProcessed)

	case "queue_depth":
		var maxDepth int64
		for _, m := range metrics {
			if m.QueueDepth > maxDepth {
				maxDepth = m.QueueDepth
			}
		}
		return float64(maxDepth)

	case "processing_time":
		var totalTime time.Duration
		var count int
		for _, m := range metrics {
			if m.AverageTime > 0 {
				totalTime += m.AverageTime
				count++
			}
		}
		if count == 0 {
			return 0
		}
		return totalTime.Seconds() / float64(count)

	default:
		return 0
	}
}

// Alerts 返回告警通道
func (emc *EventMetricsCollector) Alerts() <-chan Alert {
	return emc.alertCh
}

// AddAlertRule 添加告警规则
func (emc *EventMetricsCollector) AddAlertRule(rule AlertRule) {
	emc.mu.Lock()
	defer emc.mu.Unlock()

	emc.alertRules = append(emc.alertRules, rule)
	emc.logger.Info("Alert rule added", zap.String("name", rule.Name))
}

// RemoveAlertRule 删除告警规则
func (emc *EventMetricsCollector) RemoveAlertRule(ruleName string) {
	emc.mu.Lock()
	defer emc.mu.Unlock()

	for i, rule := range emc.alertRules {
		if rule.Name == ruleName {
			emc.alertRules = append(emc.alertRules[:i], emc.alertRules[i+1:]...)
			emc.logger.Info("Alert rule removed", zap.String("name", ruleName))
			return
		}
	}
}

// HealthCheck 健康检查
func (emc *EventMetricsCollector) HealthCheck(ctx context.Context) error {
	// 检查 Redis 连接
	if emc.redisClient != nil {
		if err := emc.redisClient.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("redis connection failed: %w", err)
		}
	}

	return nil
}

// GetAlertRules 获取所有告警规则
func (emc *EventMetricsCollector) GetAlertRules() []AlertRule {
	emc.mu.RLock()
	defer emc.mu.RUnlock()

	// 返回副本
	rules := make([]AlertRule, len(emc.alertRules))
	copy(rules, emc.alertRules)
	return rules
}
