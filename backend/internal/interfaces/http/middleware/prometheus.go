package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 请求指标
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP 请求总数",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration= promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP 请求延迟时间（秒）",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP 请求大小（字节）",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP 响应大小（字节）",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
		},
		[]string{"method", "endpoint"},
	)

	// 业务指标
	businessErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_errors_total",
			Help: "业务错误总数",
		},
		[]string{"error_type", "endpoint"},
	)

	dbOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "数据库操作总数",
		},
		[]string{"operation", "table", "status"},
	)

	cacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "缓存操作总数",
		},
		[]string{"operation", "status"},
	)

	cacheHitRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "缓存命中率",
		},
		[]string{"cache_type"},
	)
)

// MetricsMiddleware Prometheus 监控指标中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 记录请求大小
		requestSize := c.Request.ContentLength
		if requestSize > 0 {
			httpRequestSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(requestSize))
		}
		
		// 处理请求
		c.Next()
		
		// 计算延迟
		duration := time.Since(start).Seconds()
		
		// 记录延迟指标
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
		
		// 记录响应大小
		responseSize := c.Writer.Size()
		if responseSize > 0 {
			httpResponseSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(responseSize))
		}
		
		// 记录请求总数
		statusCode := c.Writer.Status()
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), string(rune(statusCode))).Inc()
		
		// 记录业务错误
		if statusCode >= 400 && statusCode < 500 {
			businessErrorsTotal.WithLabelValues("client_error", c.FullPath()).Inc()
		} else if statusCode >= 500 {
			businessErrorsTotal.WithLabelValues("server_error", c.FullPath()).Inc()
		}
	}
}

// RecordDBOperation 记录数据库操作
func RecordDBOperation(operation, table string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	dbOperationsTotal.WithLabelValues(operation, table, status).Inc()
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(cacheType string, hit bool) {
	if hit {
		cacheHitRatio.WithLabelValues(cacheType).Set(1)
	} else {
		cacheHitRatio.WithLabelValues(cacheType).Set(0)
	}
	cacheOperationsTotal.WithLabelValues("get", "hit").Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(cacheType string) {
	cacheHitRatio.WithLabelValues(cacheType).Set(0)
	cacheOperationsTotal.WithLabelValues("get", "miss").Inc()
}
