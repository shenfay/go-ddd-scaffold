package factory

import (
	"github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/provider"
	"go.uber.org/zap"
)

// DomainServices 领域服务集合
type DomainServices struct {
	UserService           *userApp.UserServiceImpl
	AuthJWTService        interface{} // *authInfra.JWTService
	AuthService           auth.AuthService
	UserSideEffectHandler interface{} // userEvent.SideEffectHandler
}

// ServiceDependencies 服务依赖
type ServiceDependencies struct {
	Container            container.Container
	DomainInfrastructure *provider.DomainInfrastructureProvider
	Logger               *zap.Logger
}

// CreateUserDomainServices 创建用户领域服务
// 系统级基础设施从容器获取，领域特定基础设施从 DomainInfrastructureProvider 获取
func CreateUserDomainServices(deps *ServiceDependencies) (*DomainServices, error) {
	// 从容器获取系统级基础设施
	uow := deps.Container.GetUnitOfWork()
	eventPublisher := deps.Container.GetEventPublisher()
	snowflakeNode := deps.Container.GetSnowflake()

	// 从领域基础设施提供者获取领域特定组件（懒加载 + 缓存）
	passwordHasher := deps.DomainInfrastructure.GetPasswordHasher()
	passwordPolicy := deps.DomainInfrastructure.GetPasswordPolicy()
	jwtService := deps.DomainInfrastructure.GetJWTService(nil) // nil 表示使用内部缓存的 Redis 客户端

	// 创建应用服务
	userService := userApp.NewUserService(
		uow,
		eventPublisher,
		passwordHasher,
		passwordPolicy,
		jwtService,
		snowflakeNode,
	)

	return &DomainServices{
		UserService: userService,
	}, nil
}

// CreateAuthDomainServices 创建认证领域服务
func CreateAuthDomainServices(deps *ServiceDependencies) (*DomainServices, error) {
	// 从容器获取系统级基础设施
	uow := deps.Container.GetUnitOfWork()
	eventPublisher := deps.Container.GetEventPublisher()
	snowflakeNode := deps.Container.GetSnowflake()
	logger := deps.Logger.Named("auth")

	// 从领域基础设施提供者获取领域特定组件（懒加载 + 缓存）
	passwordHasher := deps.DomainInfrastructure.GetPasswordHasher()
	jwtService := deps.DomainInfrastructure.GetJWTService(nil) // nil 表示使用内部缓存的 Redis 客户端

	// 创建应用服务
	authService := auth.NewAuthService(
		uow,
		passwordHasher,
		jwtService,
		eventPublisher,
		snowflakeNode,
		logger,
	)

	return &DomainServices{
		AuthJWTService: jwtService,
		AuthService:    authService,
	}, nil
}
