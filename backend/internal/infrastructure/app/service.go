// Package app 提供应用级别的服务协调
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/infrastructure/server"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Application 应用程序主服务
type Application struct {
	config       *config.Config
	logger       *zap.Logger
	db           *gorm.DB
	server       *server.ServerService
	shutdownChan chan os.Signal
}

// NewApplication 创建新的应用实例
func NewApplication() (*Application, error) {

	// 1. 加载配置
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 2. 初始化日志器
	logger, err := initLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化日志失败: %w", err)
	}

	// 3. 初始化数据库
	db, err := initDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 4. 创建服务器服务
	serverService := server.NewServerService(cfg, db, logger)

	return &Application{
		config:       cfg,
		logger:       logger,
		db:           db,
		server:       serverService,
		shutdownChan: make(chan os.Signal, 1),
	}, nil
}

// Run 运行应用程序
func (app *Application) Run() error {
	app.logger.Info("开始启动 Go DDD Scaffold 应用",
		zap.String("name", app.config.App.Name),
		zap.String("env", app.config.App.Env),
		zap.Int("port", app.config.App.Port))

	// 初始化服务器服务
	if err := app.server.Initialize(); err != nil {
		return fmt.Errorf("服务器初始化失败: %w", err)
	}

	// 启动服务器（在goroutine中）
	go func() {
		if err := app.server.Start(); err != nil {
			app.logger.Error("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待关闭信号
	return app.waitForShutdown()
}

// waitForShutdown 等待关闭信号并优雅关闭
func (app *Application) waitForShutdown() error {
	// 注册关闭信号
	signal.Notify(app.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-app.shutdownChan
	app.logger.Info("收到关闭信号", zap.String("signal", sig.String()))

	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("服务器关闭失败", zap.Error(err))
		return err
	}

	// 关闭数据库连接
	sqlDB, err := app.db.DB()
	if err != nil {
		app.logger.Error("获取数据库连接失败", zap.Error(err))
	} else {
		if err := sqlDB.Close(); err != nil {
			app.logger.Error("关闭数据库连接失败", zap.Error(err))
		}
	}

	// 同步日志
	app.logger.Sync()

	app.logger.Info("应用已优雅关闭")
	return nil
}

// GetServer 获取服务器服务实例
func (app *Application) GetServer() *server.ServerService {
	return app.server
}

// GetLogger 获取日志器实例
func (app *Application) GetLogger() *zap.Logger {
	return app.logger
}

// GetDatabase 获取数据库实例
func (app *Application) GetDatabase() *gorm.DB {
	return app.db
}

// Close 关闭应用资源
func (app *Application) Close() {
	// 关闭数据库连接
	if app.db != nil {
		if sqlDB, err := app.db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// 同步日志
	if app.logger != nil {
		app.logger.Sync()
	}
}

// initLogger 初始化日志器
func initLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, fmt.Errorf("创建日志器失败: %w", err)
	}

	return logger, nil
}

// initDatabase 初始化数据库连接
func initDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)

	logger.Info("连接数据库", zap.String("host", cfg.Database.Host))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return db, nil
}
