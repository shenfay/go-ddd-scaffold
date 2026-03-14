package bootstrap

import (
	"context"
	"time"

	authCommands "github.com/shenfay/go-ddd-scaffold/internal/application/auth/commands"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	repositoryPkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"go.uber.org/zap"
)

// initAuthDomain 初始化认证领域
func (b *Bootstrap) initAuthDomain(ctx context.Context) error {
	b.logger.Info("Initializing authentication domain...")

	// === 1. 创建 JWT 服务 ===
	jwtConfig := b.config.JWT
	b.auth.jwtService = auth.NewJWTService(
		jwtConfig.Secret,
		jwtConfig.AccessExpire,
		jwtConfig.RefreshExpire,
		"go-ddd-scaffold", // issuer
	)

	// === 2. 创建基础设施服务 ===
	db := b.container.GetGormDB()

	// 初始化 DAO（必须在使用 repository 之前）
	dao.SetDefault(db)

	userRepo := repositoryPkg.NewUserRepository(db)
	passwordHasher := user.NewBcryptPasswordHasher(12)
	eventPublisher := NewInMemoryEventPublisher(b.logger.Named("auth-events"))
	snowflakeNode := b.container.GetSnowflake()

	// 包装 ID 生成器函数（忽略错误，因为 Snowflake 几乎不会失败）
	idGenerator := func() int64 {
		id, err := snowflakeNode.Generate()
		if err != nil {
			// 理论上不会发生，如果发生则回退到时间戳
			b.logger.Error("failed to generate snowflake id", zap.Error(err))
			return time.Now().UnixNano()
		}
		return id
	}

	// === 3. 创建 CQRS Handlers ===
	b.auth.authenticateHandler = authCommands.NewAuthenticateHandler(
		userRepo,
		passwordHasher,
		b.auth.jwtService,
		eventPublisher,
	)

	b.auth.registerHandler = authCommands.NewRegisterHandler(
		userRepo,
		passwordHasher,
		eventPublisher,
		idGenerator,
	)

	b.auth.refreshTokenHandler = authCommands.NewRefreshTokenHandler(
		userRepo,
		b.auth.jwtService,
	)

	b.auth.logoutHandler = authCommands.NewLogoutHandler(
		userRepo,
		b.auth.jwtService,
	)

	b.logger.Info("Authentication domain initialized successfully")
	return nil
}
