package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Logger      *zap.Logger
	ErrorMapper *apperrors.ErrorMapper
}

// DefaultMiddlewareConfig 创建默认中间件配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		Logger:      zap.NewExample(),
		ErrorMapper: apperrors.NewErrorMapper(),
	}
}

// MiddlewareFactory 中间件工厂
// 用于统一创建和管理中间件，避免重复依赖
type MiddlewareFactory struct {
	config *MiddlewareConfig
}

// NewMiddlewareFactory 创建中间件工厂
func NewMiddlewareFactory(config *MiddlewareConfig) *MiddlewareFactory {
	if config == nil {
		config = DefaultMiddlewareConfig()
	}
	return &MiddlewareFactory{
		config: config,
	}
}

// Chain 返回完整的中间件链（按正确顺序）
// 顺序：TraceID → Gin Logger → Recovery → Error → LoggerWithTrace
func (f *MiddlewareFactory) Chain() []interface{} {
	return []interface{}{
		f.TraceID(),
		f.GinLogger(),
		f.Recovery(),
		f.Error(),
		f.LoggerWithTrace(),
	}
}

// TraceID 创建 TraceID 追踪中间件
func (f *MiddlewareFactory) TraceID() interface{} {
	return TraceIDMiddleware()
}

// GinLogger 创建 Gin 默认日志中间件（彩色文本格式）
func (f *MiddlewareFactory) GinLogger() interface{} {
	return gin.Logger()
}

// Recovery 创建 Panic 恢复中间件
func (f *MiddlewareFactory) Recovery() interface{} {
	return Recovery(f.config.Logger)
}

// Error 创建错误处理中间件
func (f *MiddlewareFactory) Error() interface{} {
	return Error(f.config.ErrorMapper, f.config.Logger)
}

// LoggerWithTrace 创建带 TraceID 的日志中间件
func (f *MiddlewareFactory) LoggerWithTrace() interface{} {
	return LoggerWithTrace(f.config.Logger)
}
