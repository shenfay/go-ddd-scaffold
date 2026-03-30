package modules

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/cmd/app/factory"
	"github.com/shenfay/go-ddd-scaffold/internal/application"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
	v1 "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers/v1"
	authHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers/v1/auth"
	httpMiddleware "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
)

// AuthModule 认证模块
// 实现 bootstrap.Module 和 bootstrap.HTTPModule 接口
type AuthModule struct {
	infra      *factory.Infrastructure
	jwtService auth.TokenService
}

// NewAuthModule 创建认证模块
// 内部自行构建完整依赖链
func NewAuthModule(infra *factory.Infrastructure) *AuthModule {
	// 1. 创建 JWTService
	jwtSvc := auth.NewJWTService(
		infra.Config.JWT.Secret,
		infra.Config.JWT.AccessExpire,
		infra.Config.JWT.RefreshExpire,
		"go-ddd-scaffold",
	)
	// 2. 注入 Redis 客户端
	jwtSvc.SetRedisClient(infra.Redis)

	return &AuthModule{
		infra:      infra,
		jwtService: jwtSvc,
	}
}

// Name 返回模块名称
// 实现 bootstrap.Module 接口
func (m *AuthModule) Name() string {
	return "auth"
}

// RegisterHTTP 注册 HTTP 路由
// 实现 bootstrap.HTTPModule 接口
func (m *AuthModule) RegisterHTTP(group *gin.RouterGroup) {
	// 创建响应处理器（使用 handler 包）
	respHandler := handlers.NewHandler(m.infra.ErrorMapper)
	router := v1.NewRouter()

	// 创建认证中间件（用于需要登录的路由）
	authMiddleware := httpMiddleware.AuthMiddleware(m.jwtService)

	// 提前创建 AuthService（只创建一次，避免重复）
	authService := m.createAuthService()

	// 创建所有 Handler
	loginHandler := authHTTP.NewLoginHandler(authService, respHandler)
	registerHandler := authHTTP.NewRegisterHandler(authService, respHandler)
	refreshHandler := authHTTP.NewRefreshTokenHandler(authService, respHandler)
	logoutHandler := authHTTP.NewLogoutHandler(authService, respHandler)
	getCurrentUserHandler := authHTTP.NewGetCurrentUserHandler(authService, respHandler)

	// 注册认证路由（传入 Handler 集合）
	router.RegisterAuthRoutes(group, v1.AuthHandlers{
		Login:          loginHandler.Handle,
		Register:       registerHandler.Handle,
		Refresh:        refreshHandler.Handle,
		Logout:         logoutHandler.Handle,
		GetCurrentUser: getCurrentUserHandler.Handle,
	}, authMiddleware)
}

// createAuthService 创建认证服务（独立方法，避免重复代码）
func (m *AuthModule) createAuthService() authApp.AuthService {
	daoQuery := dao.Use(m.infra.DB)
	uow := application.NewUnitOfWork(m.infra.DB, daoQuery)

	passwordHasher := service.NewBcryptPasswordHasher(
		m.infra.Config.Security.PasswordHasher.Cost,
	)

	tokenServiceAdapter := auth.NewTokenServiceAdapter(m.jwtService)
	idGeneratorAdapter := idgen.NewGeneratorAdapter()

	return authApp.NewAuthService(
		uow,
		passwordHasher,
		tokenServiceAdapter,
		m.infra.EventPublisher,
		idGeneratorAdapter,
		m.infra.ActivityLogRepository(),
		m.infra.Logger.Named("auth"),
	)
}
