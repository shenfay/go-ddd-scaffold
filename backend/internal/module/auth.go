package module

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/support/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	authHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/auth"
)

// AuthModule 认证模块
// 实现 bootstrap.Module 和 bootstrap.HTTPModule 接口
type AuthModule struct {
	infra      *bootstrap.Infra
	jwtService auth.TokenService
	routes     *authHTTP.Routes
}

// NewAuthModule 创建认证模块
// 内部自行构建完整依赖链
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
	// 1. 创建 DAO Query
	daoQuery := dao.Use(infra.DB)

	// 2. 创建 UnitOfWork
	uow := application.NewUnitOfWork(infra.DB, daoQuery)

	// 3. 创建 JWTService
	jwtSvc := auth.NewJWTService(
		infra.Config.JWT.Secret,
		infra.Config.JWT.AccessExpire,
		infra.Config.JWT.RefreshExpire,
		"go-ddd-scaffold",
	)
	// 4. 注入 Redis 客户端
	jwtSvc.SetRedisClient(infra.Redis)

	// 5. 创建 PasswordHasher（独立创建，不共享）
	passwordHasher := service.NewBcryptPasswordHasher(
		infra.Config.Security.PasswordHasher.Cost,
	)

	// 6. 创建 AuthService
	authSvc := authApp.NewAuthService(
		uow,
		passwordHasher,
		jwtSvc,
		infra.EventPublisher,
		infra.Snowflake,
		infra.Logger.Named("auth"),
	)

	// 7. 创建 respHandler（使用 Infra 中的 ErrorMapper）
	respHandler := httpShared.NewHandler(infra.ErrorMapper)

	// 8. 创建 Auth HTTP Handler 和 Routes
	handler := authHTTP.NewHandler(authSvc, respHandler)
	routes := authHTTP.NewRoutes(handler, jwtSvc)

	return &AuthModule{
		infra:      infra,
		jwtService: jwtSvc,
		routes:     routes,
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
	m.routes.Register(group)
}

// JWTService 返回 JWT 服务供中间件使用
func (m *AuthModule) JWTService() auth.TokenService {
	return m.jwtService
}
