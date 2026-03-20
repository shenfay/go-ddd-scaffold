package factory

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	authInfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	authHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/auth"
	userHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
)

// HTTPInterfaces HTTP 接口集合
type HTTPInterfaces struct {
	Router *gin.Engine
	Deps   *http.Dependencies
}

// BuildHTTPInterfaces 构建 HTTP 接口
func BuildHTTPInterfaces(
	cfg *config.AppConfig,
	logger *zap.Logger,
	userService *userApp.UserServiceImpl,
	authService auth.AuthService,
	jwtService *authInfra.JWTService,
) (*HTTPInterfaces, error) {
	// 创建路由配置
	routerConfig := &http.RouterConfig{
		APIPrefix: "/api/v1",
		Port:      cfg.Server.Port,
	}

	// 获取全局路由总线
	router := http.GetRouter(routerConfig)

	// 创建错误映射器和 Handler
	errorMapper := kernel.NewErrorMapper()
	handler := http.NewHandler(errorMapper)
	deps := http.NewDependencies(handler)

	// 注册领域服务到依赖容器
	deps.RegisterProvider("userService", userService)
	deps.RegisterProvider("authService", authService)
	deps.RegisterProvider("jwtService", jwtService)

	// 创建并注册认证领域路由
	authHandler := authHttp.NewHandler(authService, handler)
	authRoutes := authHttp.NewRoutes(authHandler, jwtService)
	router.Register(func(routerGroup *gin.RouterGroup, h *http.Handler, d *http.Dependencies) {
		authRoutes.Register(routerGroup)
	})

	// 创建并注册用户领域路由
	userHandler := userHttp.NewHandler(userService)
	userRoutes := userHttp.NewRoutes(userHandler)
	router.Register(func(routerGroup *gin.RouterGroup, h *http.Handler, d *http.Dependencies) {
		userRoutes.Register(routerGroup, d)
	})

	// 构建完整路由（应用中间件并触发所有已注册的领域路由）
	engine := router.Build(deps, logger)

	return &HTTPInterfaces{
		Router: engine,
		Deps:   deps,
	}, nil
}

// RegisterEventHandlers 注册事件处理器
func RegisterEventHandlers(
	eventBus kernel.EventBus,
	userSideEffectHandler *userEvent.SideEffectHandler,
	logger *zap.Logger,
) error {
	// 注册用户领域事件的副作用处理器
	if userSideEffectHandler != nil {
		// 这里可以注册具体的事件处理器
		// eventBus.Subscribe("UserCreated", userSideEffectHandler.HandleUserCreated)
		// eventBus.Subscribe("UserUpdated", userSideEffectHandler.HandleUserUpdated)
		logger.Info("Event handlers registered")
	}

	return nil
}
