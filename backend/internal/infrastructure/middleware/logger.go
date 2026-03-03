package middleware

import (
	"bytes"
	"io"
	"time"

	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	SkipPaths    []string
	Format       string
	TimeFormat   string
	UTC          bool
	RequestIDKey string
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:    []string{"/health", "/metrics"},
		Format:       "json",
		TimeFormat:   "2006-01-02T15:04:05.000Z07:00",
		UTC:          false,
		RequestIDKey: "X-Request-ID",
	}
}

func Logger(cfg *LoggerConfig) gin.HandlerFunc {
	if cfg == nil {
		cfg = DefaultLoggerConfig()
	}

	logger := newZapLogger(cfg)

	return func(c *gin.Context) {
		for _, path := range cfg.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		latency := time.Since(start)

		requestID := GetReqIDFromGin(c)

		fields := []zapcore.Field{
			zap.String("requestId", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("clientIp", c.ClientIP()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.Int("bodySize", len(bodyBytes)),
		}

		if len(bodyBytes) > 0 && len(bodyBytes) < 10000 {
			fields = append(fields, zap.String("body", string(bodyBytes)))
		}

		if len(c.Errors) > 0 {
			errors := make([]string, 0, len(c.Errors))
			for _, e := range c.Errors {
				errors = append(errors, e.Error())
			}
			fields = append(fields, zap.Strings("errors", errors))
			logger.Error("request completed with errors", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("request completed with warning", fields...)
		} else {
			logger.Info("request completed successfully", fields...)
		}
	}
}

func newZapLogger(cfg *LoggerConfig) *zap.Logger {
	var config zap.Config

	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		if cfg.UTC {
			t = t.UTC()
		}
		enc.AppendString(t.Format(cfg.TimeFormat))
	}

	logger, _ := config.Build()
	return logger
}

func RecoveryWithLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.String("requestId", getRequestIDFromContext(c)),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Any("error", err),
					zap.Stack("stack"),
				)

				c.AbortWithStatusJSON(500, response.ServerErr(c.Request.Context()))
			}
		}()

		c.Next()
	}
}

func getRequestIDFromContext(c *gin.Context) string {
	return GetReqIDFromGin(c)
}
