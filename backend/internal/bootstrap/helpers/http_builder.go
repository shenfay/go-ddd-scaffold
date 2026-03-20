package helpers

import (
	"github.com/gin-gonic/gin"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	authHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/auth"
	userHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
	"go.uber.org/zap"
)

// HTTPInterfaces HTTP 接口层依赖
type HTTPInterfaces struct {
	Router  *http.Router
	Deps    *http.Dependencies
	Handler *http.Handler
}

// BuildHTTPInterfaces 构建 HTTP 接口层
func BuildHTTPInterfaces(
	config *config.AppConfig,
	logger *zap.Logger,
	userService *userApp.UserServiceImpl,
	authService authApp.AuthService,
	jwtService *auth.JWTService,
) (*HTTPInterfaces, error) {
	logger.Info("Building HTTP interfaces...")

	// 创建错误映射器
	errorMapper := kernel.NewErrorMapper()

	// 创建 HTTP Handler（响应处理）
	respHandler := http.NewHandler(errorMapper)

	// 创建用户领域 HTTP Handler（业务处理）
	userHandler := userHttp.NewHandler(userService)

	// 获取 router 并构建路由
	router := http.GetRouter(&http.RouterConfig{
		APIPrefix: "/api/v1",
		Port:      config.Server.Port,
	})

	// 创建依赖容器
	deps := http.NewDependencies(respHandler)

	// 注册用户领域路由
	userRoutes := userHttp.NewRoutes(userHandler)
	userRoutes.RegisterTo(deps)

	// 手动注册用户路由
	router.Register(func(routerGroup *gin.RouterGroup, handler *http.Handler, deps *http.Dependencies) {
		userRoutes.Register(routerGroup, deps)
	})

	// 注册认证领域路由
	authHandler := authHttp.NewHandler(authService, respHandler)
	authRoutes := authHttp.NewRoutes(authHandler, jwtService)

	router.Register(func(routerGroup *gin.RouterGroup, handler *http.Handler, deps *http.Dependencies) {
		authRoutes.Register(routerGroup)
	})

	// 构建路由（触发所有领域的注册）
	engine := router.Build(deps, logger)
	if engine == nil {
		logger.Error("Failed to build HTTP routes")
		return nil, nil
	}

	logger.Info("HTTP interfaces built successfully")

	return &HTTPInterfaces{
		Router:  router,
		Deps:    deps,
		Handler: respHandler,
	}, nil
}
