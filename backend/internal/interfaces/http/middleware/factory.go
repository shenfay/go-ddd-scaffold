package middleware

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Logger      *zap.Logger
	ErrorMapper *kernel.ErrorMapper
}

// DefaultMiddlewareConfig 创建默认中间件配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		Logger:      zap.NewExample(),
		ErrorMapper: kernel.NewErrorMapper(),
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

// CreateLogger 创建日志器（双输出模式：控制台 + 文件）
// 如果 logFile 为空，则只输出到控制台
func CreateLogger(logFile string) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 如果指定了日志文件，启用双输出模式
	if logFile != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}

		// 双输出：stdout + 日志文件
		appLogPath := logFile
		errorLogPath := filepath.Join(logDir, "error.log")

		config.OutputPaths = []string{"stdout", appLogPath}
		config.ErrorOutputPaths = []string{"stderr", errorLogPath}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
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
