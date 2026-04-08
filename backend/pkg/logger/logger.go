package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// Init 初始化日志系统
func Init(level string) error {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(parseLevel(level))

	// 自定义编码格式
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	Logger = logger.Sugar()
	return nil
}

// parseLevel 解析日志级别字符串
func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 输出调试日志
func Debug(args ...interface{}) {
	if Logger != nil {
		Logger.Debug(args...)
	}
}

// Info 输出信息日志
func Info(args ...interface{}) {
	if Logger != nil {
		Logger.Info(args...)
	}
}

// Warn 输出警告日志
func Warn(args ...interface{}) {
	if Logger != nil {
		Logger.Warn(args...)
	}
}

// Error 输出错误日志
func Error(args ...interface{}) {
	if Logger != nil {
		Logger.Error(args...)
	}
}

// Sync 刷新日志缓冲区到磁盘
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

// SetLevel 动态设置日志级别（需要外部维护 AtomicLevel）
// 注意：此方法仅用于示例，实际使用中建议通过配置文件控制
