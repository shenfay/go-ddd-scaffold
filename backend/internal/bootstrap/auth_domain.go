package bootstrap

import (
	"context"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	repositoryPkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
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
	// 使用 Bootstrap 中创建的事件总线
	eventPublisher := b.eventBus

	// === 3. 创建应用服务（统一入口）===
	b.auth.authService = authApp.NewAuthService(
		userRepo,
		passwordHasher,
		b.auth.jwtService,
		eventPublisher,
	)

	b.logger.Info("Authentication domain initialized successfully")
	return nil
}
