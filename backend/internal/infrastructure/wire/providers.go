package wire

// Package wire 提供依赖注入配置

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/infrastructure/auth"
	infraEvent "go-ddd-scaffold/internal/infrastructure/event"
)

// InitializeConfig 从配置文件加载配置
func InitializeConfig(configPath string) (*config.Config, error) {
	return config.LoadConfig(configPath)
}

// InitializeDB 根据配置初始化数据库连接
func InitializeDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// InitializeJWTService 初始化 JWT 服务
func InitializeJWTService(cfg *config.Config) entity.JWTService {
	return auth.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.ExpireIn)
}

// InitializeCasbinService 初始化 Casbin 权限服务
func InitializeCasbinService(db *gorm.DB) (auth.CasbinService, error) {
	enforcer, err := auth.NewCasbinEnforcer(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Casbin: %w", err)
	}
	return auth.NewCasbinService(enforcer), nil
}

// InitializeEventBus 初始化事件总线
func InitializeEventBus() *infraEvent.EventBus {
	return infraEvent.NewEventBus()
}

// InitializeEventHandlers 初始化事件处理器并注册到事件总线
func InitializeEventHandlers(bus *infraEvent.EventBus) {
	// TODO: 根据业务需要注册事件处理器
	// 示例：
	// bus.RegisterHandler("UserCreated", userEvent.NewUserCreatedHandler())
}
