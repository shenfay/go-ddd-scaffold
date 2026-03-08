//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	authService "go-ddd-scaffold/internal/application/user/service"
	tenantService "go-ddd-scaffold/internal/application/tenant/service"
	"go-ddd-scaffold/internal/domain/user/entity"
	domainService "go-ddd-scaffold/internal/domain/user/service"
	"go-ddd-scaffold/internal/infrastructure/auth"
	infraEvent "go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
	"go-ddd-scaffold/internal/infrastructure/transaction"
)

// ==================== 核心 Provider Set ====================

// RepositorySet 仓储集合
var RepositorySet = wire.NewSet(
	repo.NewUserDAORepository,
	repo.NewTenantDAORepository,
	repo.NewTenantMemberDAORepository,
)

// TransactionSet 事务管理集合
var TransactionSet = wire.NewSet(
	transaction.NewGormUnitOfWork,
)

// AuthServiceSet 认证服务集合
var AuthServiceSet = wire.NewSet(
	auth.NewJWTService,
	
	InitializeCasbinService,
)

// ==================== 应用服务初始化 ====================

// InitializeUserCommandService 初始化用户命令服务
func InitializeUserCommandService(
	db *gorm.DB,
	passwordHasher domainService.PasswordHasher,
	uow transaction.UnitOfWork,
) authService.UserCommandService {
	wire.Build(
		RepositorySet,
		authService.NewUserCommandService,
	)
	return nil
}

// InitializeTenantService 初始化租户服务（带 UnitOfWork）
func InitializeTenantService(
	db *gorm.DB,
	casbinService auth.CasbinService,
	uow transaction.UnitOfWork,
) tenantService.TenantService {
	wire.Build(
		RepositorySet,
		tenantService.NewTenantService,
	)
	return nil
}

// InitializeAuthenticationService 初始化认证服务
func InitializeAuthenticationService(
	db *gorm.DB,
	logger interface{}, // *zap.Logger
	jwtService entity.JWTService,
	eventBus *infraEvent.EventBus,
	tokenBlacklistService auth.TokenBlacklistService,
	passwordHasher domainService.PasswordHasher,
) authService.AuthenticationService {
	wire.Build(
		RepositorySet,
		newUserEventBusAdapter,
		authService.NewAuthenticationService,
	)
	return nil
}

// ==================== 辅助函数 ====================

// userEventBusAdapter User 模块的事件总线适配器
type userEventBusAdapter struct {
	bus *infraEvent.EventBus
}

func newUserEventBusAdapter(bus *infraEvent.EventBus) authService.EventBus {
	return &userEventBusAdapter{bus: bus}
}

func (a *userEventBusAdapter) Publish(ctx context.Context, event infraEvent.DomainEvent) error {
	if event == nil {
		return nil
	}
	return a.bus.Publish(ctx, event)
}
