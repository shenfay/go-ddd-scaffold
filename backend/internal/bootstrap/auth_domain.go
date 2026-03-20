package bootstrap

import (
	"context"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	authInfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/domain_event"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/task_queue"
)

// initAuthDomain 初始化认证领域
func (b *Bootstrap) initAuthDomain(ctx context.Context) error {
	b.logger.Info("Initializing authentication domain...")

	// === 1. 创建 JWT 服务 ===
	jwtConfig := b.config.JWT
	b.auth.jwtService = authInfra.NewJWTService(
		jwtConfig.Secret,
		jwtConfig.AccessExpire,
		jwtConfig.RefreshExpire,
		"go-ddd-scaffold", // issuer
	)

	// === 1.5. 注入 Redis 客户端（用于令牌黑名单）===
	redisClient := b.container.GetRedis()
	b.auth.jwtService.SetRedisClient(redisClient)

	// === 2. 从容器获取基础设施服务 ===
	userRepo := b.container.GetUserRepo()
	passwordHasher := user.NewBcryptPasswordHasher(12)
	// 从容器获取 logger
	logger := b.container.GetLogger("auth")
	// 创建 asynq 事件发布器
	asynqClient := task_queue.NewClient(task_queue.Config{
		RedisAddr:     b.config.Redis.Addr,
		RedisPassword: b.config.Redis.Password,
		RedisDB:       b.config.Redis.DB,
	})
	asynqPublisher := task_queue.NewPublisher(asynqClient)
	eventPublisher := domain_event.NewAsynqPublisher(asynqPublisher, logger.Named("publisher"))

	// === 3. 创建应用服务（统一入口）===
	b.auth.authService = authApp.NewAuthService(
		userRepo,
		passwordHasher,
		b.auth.jwtService,
		eventPublisher,
		b.container.GetSnowflake(),
		logger,
	)

	b.logger.Info("Authentication domain initialized successfully")
	return nil
}
