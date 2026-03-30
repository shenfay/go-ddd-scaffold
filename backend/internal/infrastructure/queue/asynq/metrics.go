package asynq

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// TaskEnqueuedTotal 任务入队总数
	TaskEnqueuedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "asynq_tasks_enqueued_total",
			Help: "Total number of tasks enqueued",
		},
		[]string{"task_type", "queue"},
	)

	// TaskProcessedTotal 任务处理总数
	TaskProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "asynq_tasks_processed_total",
			Help: "Total number of tasks processed",
		},
		[]string{"task_type", "queue", "status"},
	)

	// TaskProcessingDuration 任务处理时长
	TaskProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "asynq_task_processing_duration_seconds",
			Help:    "Task processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"task_type", "queue"},
	)

	// QueueSize 队列大小
	QueueSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "asynq_queue_size",
			Help: "Current size of the queue",
		},
		[]string{"queue"},
	)

	// WorkerConcurrency 工作并发度
	WorkerConcurrency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "asynq_worker_concurrency",
			Help: "Current worker concurrency",
		},
	)
)

// RegisterMetrics 注册 Prometheus 指标
func RegisterMetrics() {
	prometheus.MustRegister(TaskEnqueuedTotal)
	prometheus.MustRegister(TaskProcessedTotal)
	prometheus.MustRegister(TaskProcessingDuration)
	prometheus.MustRegister(QueueSize)
	prometheus.MustRegister(WorkerConcurrency)
}
