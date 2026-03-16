package bootstrap

import (
	"context"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	userDomain "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	repositoryPkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
)

// initUserDomain 初始化用户领域
func (b *Bootstrap) initUserDomain(ctx context.Context) error {
	b.logger.Info("Initializing user domain...")

	baseLogger := b.logger.Named("user")

	// === 1. 创建基础设施服务 ===
	// 使用 Bootstrap 中创建的事件总线
	eventPublisher := b.eventBus

	// === 2. 创建仓储层 ===
	db := b.container.GetGormDB()

	// 初始化 DAO（必须在使用 repository 之前）
	dao.SetDefault(db)

	userRepo := repositoryPkg.NewUserRepository(db)

	// === 3. 创建应用服务（统一入口）===
	passwordHasher := userDomain.NewBcryptPasswordHasher(12)
	b.user.service = userApp.NewUserService(userRepo, eventPublisher, passwordHasher, b.auth.jwtService)

	// === 4. 创建领域事件处理器 ===
	b.user.eventHandler = userApp.NewUserEventHandler(baseLogger.Named("events"))

	b.logger.Info("User domain initialized successfully")
	return nil
}
