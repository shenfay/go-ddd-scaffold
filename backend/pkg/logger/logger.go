package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Config 日志配置
type Config struct {
	Level      string `yaml:"level" mapstructure:"level"`             // 日志级别：debug, info, warn, error
	Format     string `yaml:"format" mapstructure:"format"`           // 日志格式：json, console
	FilePath   string `yaml:"file_path" mapstructure:"file_path"`     // 日志文件路径（可选）
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`       // 单个文件最大大小 (MB)
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"` // 保留旧文件最大数量
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`         // 文件保留最大天数 (天)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "json",
		FilePath:   "", // 空表示只输出到 stdout
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
	}
}

// Init 初始化日志（单例）
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		log, err = newLogger(cfg)
	})
	return err
}

// newLogger 创建新的 logger
func newLogger(cfg *Config) (*zap.Logger, error) {
	// 解析日志级别
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level.SetLevel(zap.InfoLevel) // 默认 info
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建输出目标
	var cores []zapcore.Core

	// 标准输出
	stdoutCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, stdoutCore)

	// 文件输出（如果配置了文件路径）
	if cfg.FilePath != "" {
		fileWriter, _, err := zap.Open(cfg.FilePath)
		if err == nil {
			fileCore := zapcore.NewCore(
				encoder,
				zapcore.AddSync(fileWriter),
				level,
			)
			cores = append(cores, fileCore)
		}
	}

	// 创建核心（多输出）
	core := zapcore.NewTee(cores...)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}

// Get 获取全局 logger
func Get() *zap.Logger {
	if log == nil {
		// 如果未初始化，返回一个开发模式的 logger
		l, _ := zap.NewDevelopment()
		return l
	}
	return log
}

// Sync 同步日志到磁盘
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// 便捷的日志方法

// Debug 记录 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 记录 info 级别日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 记录 warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 记录 error 级别日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 记录 fatal 级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// With 添加字段
func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

// 常用字段辅助函数

// String 字符串字段
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int 整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 int64 字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Bool 布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Error 错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration 时长字段
func Duration(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Any 任意类型字段
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
