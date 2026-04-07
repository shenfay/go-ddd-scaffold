package integration_test

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
)

// AuthIntegrationSuite 认证集成测试套件
type AuthIntegrationSuite struct {
	suite.Suite
	db           *gorm.DB
	redis        *redis.Client
	asynqClient  *asynq.Client
	router       *gin.Engine
	authService  *authentication.Service
	tokenService authentication.TokenService
}

// SetupSuite 在测试套件开始前执行一次
func (s *AuthIntegrationSuite) SetupSuite() {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 加载配置
	cfg, err := config.Load("test")
	if err != nil {
		// 如果 test 配置不存在，使用默认配置
		cfg = &config.Config{
			Database: config.DatabaseConfig{
				Host:            getEnv("TEST_DB_HOST", "localhost"),
				Port:            5432,
				Name:            getEnv("TEST_DB_NAME", "go_ddd_scaffold_test"),
				User:            getEnv("TEST_DB_USER", "postgres"),
				Password:        getEnv("TEST_DB_PASSWORD", "postgres"),
				SSLMode:         "disable",
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
			},
			Redis: config.RedisConfig{
				Addr:     getEnv("TEST_REDIS_ADDR", "localhost:6379"),
				Password: "",
				DB:       15, // 使用 DB 15 用于测试，避免污染生产数据
			},
			JWT: config.JWTConfig{
				Secret:        "test-secret-key-not-for-production",
				AccessExpire:  30 * time.Minute,
				RefreshExpire: 7 * 24 * time.Hour,
				Issuer:        "go-ddd-scaffold-test",
			},
		}
	}

	// 初始化数据库
	s.initDatabase(cfg.Database)

	// 初始化 Redis
	s.initRedis(cfg.Redis)

	// 初始化 Asynq Client
	s.initAsynqClient(cfg.Redis)

	// 初始化服务
	s.initServices(cfg)

	// 初始化路由
	s.initRouter()
}

// SetupTest 在每个测试开始前执行
func (s *AuthIntegrationSuite) SetupTest() {
	// 清理测试数据
	s.cleanupTestData()
}

// TearDownSuite 在测试套件结束后执行
func (s *AuthIntegrationSuite) TearDownSuite() {
	// 关闭数据库连接
	sqlDB, _ := s.db.DB()
	sqlDB.Close()

	// 关闭 Redis 连接
	s.redis.Close()

	// 关闭 Asynq Client
	if s.asynqClient != nil {
		s.asynqClient.Close()
	}
}

// initDatabase 初始化数据库
func (s *AuthIntegrationSuite) initDatabase(cfg config.DatabaseConfig) {
	dsn := cfg.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默日志
	})
	s.Require().NoError(err, "Failed to connect to database")

	// 自动迁移表结构
	err = db.AutoMigrate(&repository.UserPO{})
	s.Require().NoError(err, "Failed to migrate database")

	s.db = db
}

// initRedis 初始化 Redis
func (s *AuthIntegrationSuite) initRedis(cfg config.RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 10,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	s.Require().NoError(err, "Failed to connect to Redis")

	s.redis = client
}

// initAsynqClient 初始化 Asynq Client
func (s *AuthIntegrationSuite) initAsynqClient(cfg config.RedisConfig) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	s.asynqClient = client
}

// initServices 初始化服务
func (s *AuthIntegrationSuite) initServices(cfg *config.Config) {
	userRepo := repository.NewUserRepository(s.db)
	tokenService := authentication.NewTokenServiceImpl(
		s.redis,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshExpire,
	)
	s.tokenService = tokenService
	s.authService = authentication.NewService(userRepo, tokenService)
}

// initRouter 初始化路由
func (s *AuthIntegrationSuite) initRouter() {
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// 注册认证路由
	v1 := router.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(s.authService, s.tokenService, s.asynqClient)
		authHandler.RegisterRoutes(v1)
	}

	s.router = router
}

// cleanupTestData 清理测试数据
func (s *AuthIntegrationSuite) cleanupTestData() {
	// 删除所有用户数据
	s.db.Exec("DELETE FROM users")

	// 清理 Redis 测试数据
	s.redis.FlushDB(context.Background())
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
