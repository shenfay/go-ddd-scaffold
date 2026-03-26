package module

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/cmd/shared/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/application"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
	httpMiddleware "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	v1 "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/rest/v1"
	authHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/rest/v1/auth"
)

// AuthModule 认证模块
// 实现 bootstrap.Module 和 bootstrap.HTTPModule 接口
type AuthModule struct {
	infra      *bootstrap.Infrastructure
	jwtService auth.TokenService
}

// NewAuthModule 创建认证模块
// 内部自行构建完整依赖链
func NewAuthModule(infra *bootstrap.Infrastructure) *AuthModule {
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
	respHandler := handler.NewHandler(m.infra.ErrorMapper)
	router := v1.NewRouter()

	// 创建认证中间件（用于需要登录的路由）
	authMiddleware := httpMiddleware.AuthMiddleware(m.jwtService)

	// 提供 Handler 的工厂函数（避免循环依赖）
	handlerProvider := func() (login, register, refresh, logout, getCurrentUser gin.HandlerFunc) {
		// 重新创建 AuthService（因为需要依赖注入）
		daoQuery := dao.Use(m.infra.DB)
		uow := application.NewUnitOfWork(m.infra.DB, daoQuery)

		passwordHasher := service.NewBcryptPasswordHasher(
			m.infra.Config.Security.PasswordHasher.Cost,
		)

		tokenServiceAdapter := auth.NewTokenServiceAdapter(m.jwtService)
		idGeneratorAdapter := idgen.NewGeneratorAdapter()

		authSvc := authApp.NewAuthService(
			uow,
			passwordHasher,
			tokenServiceAdapter,
			m.infra.EventPublisher,
			idGeneratorAdapter,
			m.infra.Logger.Named("auth"),
		)

		// 创建所有 Handler
		loginHandler := authHTTP.NewLoginHandler(authSvc, respHandler)
		registerHandler := authHTTP.NewRegisterHandler(authSvc, respHandler)
		refreshHandler := authHTTP.NewRefreshTokenHandler(authSvc, respHandler)
		logoutHandler := authHTTP.NewLogoutHandler(authSvc, respHandler)
		getCurrentUserHandler := authHTTP.NewGetCurrentUserHandler(authSvc, respHandler)

		return loginHandler.ServeHTTP, registerHandler.ServeHTTP,
			refreshHandler.ServeHTTP, logoutHandler.ServeHTTP,
			getCurrentUserHandler.ServeHTTP
	}

	// 注册认证路由（传入认证中间件）
	router.RegisterAuthRoutes(group, handlerProvider, authMiddleware)
}

// JWTService 返回 JWT 服务供中间件使用
func (m *AuthModule) JWTService() auth.TokenService {
	return m.jwtService
}
