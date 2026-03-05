package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// APIMetrics API指标收集器
type APIMetrics struct {
	// 原子计数器
	totalRequests uint64
	totalErrors   uint64
	totalLatency  uint64 // 纳秒
	totalSuccess  uint64

	// 按路径和方法的详细统计
	endpointStats map[string]*EndpointStat
	mu            sync.RWMutex

	// 启动时间
	startTime time.Time
}

// EndpointStat 端点统计信息
type EndpointStat struct {
	Path          string `json:"path"`
	Method        string `json:"method"`
	TotalRequests uint64 `json:"total_requests"`
	TotalErrors   uint64 `json:"total_errors"`
	TotalLatency  uint64 `json:"total_latency_ns"` // 纳秒
	MinLatency    uint64 `json:"min_latency_ns"`
	MaxLatency    uint64 `json:"max_latency_ns"`
}

// NewAPIMetrics 创建新的指标收集器
func NewAPIMetrics() *APIMetrics {
	return &APIMetrics{
		endpointStats: make(map[string]*EndpointStat),
		startTime:     time.Now(),
	}
}

// RecordRequest 记录请求指标
func (m *APIMetrics) RecordRequest(method, path string, latency time.Duration, statusCode int) {
	// 原子增加总请求数
	atomic.AddUint64(&m.totalRequests, 1)
	atomic.AddUint64(&m.totalLatency, uint64(latency.Nanoseconds()))

	// 统计成功/错误
	if statusCode >= 200 && statusCode < 400 {
		atomic.AddUint64(&m.totalSuccess, 1)
	} else {
		atomic.AddUint64(&m.totalErrors, 1)
	}

	// 更新端点详细统计
	key := method + ":" + path
	m.mu.Lock()
	stat, exists := m.endpointStats[key]
	if !exists {
		stat = &EndpointStat{
			Path:       path,
			Method:     method,
			MinLatency: uint64(latency.Nanoseconds()),
			MaxLatency: uint64(latency.Nanoseconds()),
		}
		m.endpointStats[key] = stat
	}

	// 更新端点统计
	atomic.AddUint64(&stat.TotalRequests, 1)
	atomic.AddUint64(&stat.TotalLatency, uint64(latency.Nanoseconds()))

	// 更新最大最小延迟
	latencyNs := uint64(latency.Nanoseconds())
	if latencyNs > atomic.LoadUint64(&stat.MaxLatency) {
		atomic.StoreUint64(&stat.MaxLatency, latencyNs)
	}
	if latencyNs < atomic.LoadUint64(&stat.MinLatency) {
		atomic.StoreUint64(&stat.MinLatency, latencyNs)
	}

	// 统计错误
	if statusCode >= 400 {
		atomic.AddUint64(&stat.TotalErrors, 1)
	}
	m.mu.Unlock()
}

// GetMetrics 获取当前指标快照
func (m *APIMetrics) GetMetrics() *MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := &MetricsSnapshot{
		StartTime:     m.startTime,
		Uptime:        time.Since(m.startTime),
		TotalRequests: atomic.LoadUint64(&m.totalRequests),
		TotalErrors:   atomic.LoadUint64(&m.totalErrors),
		TotalLatency:  atomic.LoadUint64(&m.totalLatency),
		TotalSuccess:  atomic.LoadUint64(&m.totalSuccess),
		Endpoints:     make([]*EndpointStat, 0, len(m.endpointStats)),
	}

	// 复制端点统计数据
	for _, stat := range m.endpointStats {
		// 创建副本避免竞态
		copyStat := &EndpointStat{
			Path:          stat.Path,
			Method:        stat.Method,
			TotalRequests: atomic.LoadUint64(&stat.TotalRequests),
			TotalErrors:   atomic.LoadUint64(&stat.TotalErrors),
			TotalLatency:  atomic.LoadUint64(&stat.TotalLatency),
			MinLatency:    atomic.LoadUint64(&stat.MinLatency),
			MaxLatency:    atomic.LoadUint64(&stat.MaxLatency),
		}
		snapshot.Endpoints = append(snapshot.Endpoints, copyStat)
	}

	return snapshot
}

// MetricsSnapshot 指标快照
type MetricsSnapshot struct {
	StartTime     time.Time       `json:"start_time"`
	Uptime        time.Duration   `json:"uptime"`
	TotalRequests uint64          `json:"total_requests"`
	TotalErrors   uint64          `json:"total_errors"`
	TotalSuccess  uint64          `json:"total_success"`
	TotalLatency  uint64          `json:"total_latency_ns"`
	Endpoints     []*EndpointStat `json:"endpoints"`
}

// CalculateRate 计算各种比率
func (s *MetricsSnapshot) CalculateRate() *MetricsReport {
	report := &MetricsReport{
		Uptime:           s.Uptime,
		TotalRequests:    s.TotalRequests,
		QPS:              float64(s.TotalRequests) / s.Uptime.Seconds(),
		ErrorRate:        0,
		SuccessRate:      0,
		AverageLatencyMs: 0,
		P95LatencyMs:     0,
		P99LatencyMs:     0,
	}

	if s.TotalRequests > 0 {
		report.ErrorRate = float64(s.TotalErrors) / float64(s.TotalRequests) * 100
		report.SuccessRate = float64(s.TotalSuccess) / float64(s.TotalRequests) * 100
		report.AverageLatencyMs = float64(s.TotalLatency) / float64(s.TotalRequests) / 1000000 // 转换为毫秒
	}

	// 计算P95和P99延迟（简化版）
	if len(s.Endpoints) > 0 {
		var latencies []float64
		for _, ep := range s.Endpoints {
			if ep.TotalRequests > 0 {
				avg := float64(ep.TotalLatency) / float64(ep.TotalRequests) / 1000000
				latencies = append(latencies, avg)
			}
		}

		if len(latencies) > 0 {
			// 简单排序计算百分位数
			// 实际项目中应使用更精确的算法
			report.P95LatencyMs = calculatePercentile(latencies, 95)
			report.P99LatencyMs = calculatePercentile(latencies, 99)
		}
	}

	return report
}

// MetricsReport 指标报告
type MetricsReport struct {
	Uptime           time.Duration `json:"uptime"`
	TotalRequests    uint64        `json:"total_requests"`
	QPS              float64       `json:"qps"`
	ErrorRate        float64       `json:"error_rate_percent"`
	SuccessRate      float64       `json:"success_rate_percent"`
	AverageLatencyMs float64       `json:"avg_latency_ms"`
	P95LatencyMs     float64       `json:"p95_latency_ms"`
	P99LatencyMs     float64       `json:"p99_latency_ms"`
}

// calculatePercentile 计算百分位数（简化版）
func calculatePercentile(values []float64, percentile int) float64 {
	if len(values) == 0 {
		return 0
	}

	// 简化实现：取最大值作为近似
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// Reset 重置所有统计（谨慎使用）
func (m *APIMetrics) Reset() {
	atomic.StoreUint64(&m.totalRequests, 0)
	atomic.StoreUint64(&m.totalErrors, 0)
	atomic.StoreUint64(&m.totalLatency, 0)
	atomic.StoreUint64(&m.totalSuccess, 0)

	m.mu.Lock()
	m.endpointStats = make(map[string]*EndpointStat)
	m.mu.Unlock()
}
