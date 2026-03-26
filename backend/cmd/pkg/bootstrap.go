package pkg

import (
	"log"
	"os"

	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	logging "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/log"
	"go.uber.org/zap"
)

// InitInfrastructure 初始化基础设施
// 返回 Infra、Logger 和 cleanup 函数
func InitInfrastructure(appName string) (*bootstrap.Infra, *zap.Logger, func()) {
	// 1. 加载配置
	env := os.Getenv("ENV_MODE")
	if env == "" {
		env = "development"
	}

	configLoader := config.NewConfigLoader(nil)
	appConfig, err := configLoader.Load(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 创建 Logger
	logConfig := &config.LoggingConfig{
		Level:      appConfig.Logging.Level,
		Format:     appConfig.Logging.Format,
		File:       appConfig.Logging.File,
		MaxSize:    appConfig.Logging.MaxSize,
		MaxBackups: appConfig.Logging.MaxBackups,
		MaxAge:     appConfig.Logging.MaxAge,
	}
	appLogger, err := logging.New(logConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	logger := appLogger.Logger.Named(appName)
	logger.Info("Starting "+appName+"...", zap.String("env", env))

	// 3. 创建基础设施（调用现有的 bootstrap.NewInfra）
	infra, cleanup, err := bootstrap.NewInfra(appConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create infrastructure", zap.Error(err))
	}

	return infra, logger, cleanup
}
