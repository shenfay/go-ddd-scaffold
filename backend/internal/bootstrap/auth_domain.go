package bootstrap

import (
	"context"

	authCommands "github.com/shenfay/go-ddd-scaffold/internal/application/auth/commands"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence"
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
	db := b.container.GetDB()
	userRepo := persistence.NewUserRepository(db)
	passwordHasher := user.NewBcryptPasswordHasher(12)
	eventPublisher := NewInMemoryEventPublisher(b.logger.Named("auth-events"))

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
