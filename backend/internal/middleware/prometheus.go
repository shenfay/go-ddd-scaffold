package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
)

// PrometheusMiddleware Prometheus 监控中间件
func PrometheusMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 增加请求计数
		m.IncHTTPRequests()
		
		// 增加正在处理的请求数
		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()
		
		// 继续处理请求
		c.Next()
		
		// 计算耗时
		duration := time.Since(start).Seconds()
		
		// 记录请求耗时
		status := strconv.Itoa(c.Writer.Status())
		m.ObserveHTTPDuration(c.Request.Method, c.FullPath(), status, duration)
	}
}
