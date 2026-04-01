package log

import (
	"os"
	"path/filepath"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DefaultConfig 返回默认配置
func DefaultConfig() *config.LoggingConfig {
	return &config.LoggingConfig{
		Level:      "debug",
		Format:     "console",
		File:       "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
	}
}

// Logger 应用日志器
type Logger struct {
	*zap.Logger
	config *config.LoggingConfig
}

// New 创建日志器
func New(config *config.LoggingConfig) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 确保日志目录存在
	logDir := filepath.Dir(config.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 根据格式配置编码器
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if config.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// 开发环境使用彩色文本格式
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 日志级别
	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, err
	}

	// 双输出：stdout + 日志文件
	appLogPath := config.File
	errorLogPath := filepath.Join(logDir, "error.log")

	// 创建写入器
	writeSyncer := getLogWriter(appLogPath, config.MaxSize, config.MaxBackups, config.MaxAge)
	errorWriteSyncer := getLogWriter(errorLogPath, config.MaxSize, config.MaxBackups, config.MaxAge)

	// 合并输出目标
	core := zapcore.NewTee(
		// stdout 输出
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		// 文件输出 - 所有级别
		zapcore.NewCore(encoder, writeSyncer, level),
		// 错误文件输出 - 仅 error 及以上级别
		zapcore.NewCore(encoder, errorWriteSyncer, zap.ErrorLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{
		Logger: logger,
		config: config,
	}, nil
}

// getLogWriter 创建日志写入器（带轮转功能）
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	// 使用 lumberjack 实现日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,    // 每个文件最大大小 (MB)
		MaxBackups: maxBackups, // 保留旧文件个数
		MaxAge:     maxAge,     // 文件保留天数
		Compress:   true,       // 压缩旧文件
	}

	return zapcore.AddSync(lumberJackLogger)
}

// Sync 刷新日志缓冲区
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// GetConfig 获取日志配置
func (l *Logger) GetConfig() *config.LoggingConfig {
	return l.config
}
