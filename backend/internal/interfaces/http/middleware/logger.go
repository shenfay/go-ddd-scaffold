package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerWithTrace 带 TraceID 的日志中间件
// 在 Gin 默认日志基础上添加 trace_id 字段
// 依赖 TraceIDMiddleware 提供的 trace_id 进行日志关联
func LoggerWithTrace(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 TraceID
		traceID := GetTraceID(c)

		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 记录带 TraceID 的请求日志
		logger.Info("Request completed",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}
