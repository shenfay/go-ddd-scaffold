package wire

// Package wire 提供依赖注入配置

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/infrastructure/auth"
	infraEvent "go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/pkg/metrics"
	"go-ddd-scaffold/internal/pkg/ratelimit"
)

// InitializeConfig 从配置文件加载配置
func InitializeConfig(configPath string) (*config.Config, error) {
	return config.LoadConfig(configPath)
}

// InitializeDB 根据配置初始化数据库连接（带连接池优化）
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

	// 获取底层 sql.DB 进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	
	// 解析时间配置
	connMaxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
	if err != nil {
		connMaxLifetime = time.Hour // 默认 1 小时
	}
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	
	connMaxIdleTime, err := time.ParseDuration(cfg.Database.ConnMaxIdleTime)
	if err != nil {
		connMaxIdleTime = 5 * time.Minute // 默认 5 分钟
	}
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return db, nil
}

// InitializeRedis 初始化 Redis 客户端
func InitializeRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
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

// InitializeEventBus 初始化事件总线（Redis Stream 版本）
func InitializeEventBus(cfg *config.Config, rdb *redis.Client) (*infraEvent.RedisEventBus, error) {
	// 创建 Redis Event Bus
	bus := infraEvent.NewRedisEventBus(rdb, infraEvent.RedisEventBusConfig{
		MaxRetries:     cfg.Redis.EventBusConfig.MaxRetries,
		RetryBaseDelay: cfg.Redis.EventBusConfig.RetryBaseDelay,
		PollInterval:   cfg.Redis.EventBusConfig.PollInterval,
		BatchSize:      cfg.Redis.EventBusConfig.BatchSize,
	})

	return bus, nil
}

// InitializeEventHandlers 初始化事件处理器并注册到事件总线
func InitializeEventHandlers(bus *infraEvent.RedisEventBus) {
	// TODO: 根据业务需要注册事件处理器
	// 示例：
	// bus.RegisterHandler("UserCreated", userEvent.NewUserCreatedHandler())
}

// InitializeTokenBlacklistService 初始化 Token 黑名单服务（带监控和限流熔断）
func InitializeTokenBlacklistService(
	rdb *redis.Client,
	metrics *metrics.Metrics,
	rateLimiter *ratelimit.RateLimiter,
	circuitBreaker *ratelimit.CircuitBreaker,
) auth.TokenBlacklistService {
	return auth.NewRedisTokenBlacklistService(
		rdb,
		"token:blacklist:",
		rateLimiter,
		circuitBreaker,
		metrics,
	)
}
