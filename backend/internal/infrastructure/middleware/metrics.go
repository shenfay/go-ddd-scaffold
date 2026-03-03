package middleware

import (
	"time"

	"go-ddd-scaffold/internal/pkg/metrics"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// Metrics 监控中间件
// 收集API调用的QPS、延迟、错误率等指标
func Metrics(metricsCollector *metrics.APIMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 继续处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 记录指标
		metricsCollector.RecordRequest(
			c.Request.Method,
			c.Request.URL.Path,
			latency,
			c.Writer.Status(),
		)
	}
}

// MetricsHandler 监控数据查询接口
func MetricsHandler(metricsCollector *metrics.APIMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取指标快照
		snapshot := metricsCollector.GetMetrics()

		// 计算比率
		report := snapshot.CalculateRate()

		// 返回监控数据
		c.JSON(200, response.OKWithMsg(c.Request.Context(), report, "监控数据获取成功"))
	}
}

// ResetMetricsHandler 重置监控数据接口（谨慎使用）
func ResetMetricsHandler(metricsCollector *metrics.APIMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector.Reset()

		c.JSON(200, response.OKWithMsg(c.Request.Context(), nil, "监控数据已重置"))
	}
}
