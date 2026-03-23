package module

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	sharedAggregate "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	eventHandler "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	userHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/useragent"
)

// UserModule 用户模块
// 实现 bootstrap.Module、bootstrap.HTTPModule 和 bootstrap.EventModule 接口
type UserModule struct {
	infra   *bootstrap.Infra
	routes  *userHTTP.Routes
	handler *userHTTP.Handler
	// 事件订阅器
	sideEffectHandler  *userEvent.SideEffectHandler
	auditSubscriber    *eventHandler.AuditSubscriber
	loginLogSubscriber *eventHandler.LoginLogSubscriber
}

// NewUserModule 创建用户模块
// 内部自行构建完整依赖链
func NewUserModule(infra *bootstrap.Infra) *UserModule {
	// 1. 创建 DAO Query
	daoQuery := dao.Use(infra.DB)

	// 2. 创建 UnitOfWork
	uow := application.NewUnitOfWork(infra.DB, daoQuery)

	// 3. 创建 PasswordHasher
	passwordHasher := service.NewBcryptPasswordHasher(
		infra.Config.Security.PasswordHasher.Cost,
	)

	// 4. 创建 PasswordPolicy
	policyConfig := service.PasswordPolicyConfig{
		MinLength:           infra.Config.Security.PasswordPolicy.MinLength,
		MaxLength:           infra.Config.Security.PasswordPolicy.MaxLength,
		RequireUppercase:    infra.Config.Security.PasswordPolicy.RequireUppercase,
		RequireLowercase:    infra.Config.Security.PasswordPolicy.RequireLowercase,
		RequireDigits:       infra.Config.Security.PasswordPolicy.RequireDigits,
		RequireSpecialChars: infra.Config.Security.PasswordPolicy.RequireSpecialChars,
		SpecialChars:        infra.Config.Security.PasswordPolicy.SpecialChars,
		DisallowCommon:      infra.Config.Security.PasswordPolicy.DisallowCommon,
	}
	passwordPolicy := auth.NewDefaultPasswordPolicy(policyConfig)

	// 5. 创建 JWTService (UserService 需要 TokenService)
	jwtSvc := auth.NewJWTService(
		infra.Config.JWT.Secret,
		infra.Config.JWT.AccessExpire,
		infra.Config.JWT.RefreshExpire,
		"go-ddd-scaffold",
	)
	jwtSvc.SetRedisClient(infra.Redis)

	// 6. 创建 UserService
	userSvc := userApp.NewUserService(
		uow,
		infra.EventPublisher,
		passwordHasher,
		passwordPolicy,
		jwtSvc,
		infra.Snowflake,
	)

	// 7. 创建 respHandler（使用 Infra 中的 ErrorMapper）
	respHandler := httpShared.NewHandler(infra.ErrorMapper)

	// 8. 创建 HTTP Handler 和 Routes
	handler := userHTTP.NewHandler(userSvc, respHandler)
	routes := userHTTP.NewRoutes(handler)

	// 9. 创建事件订阅器（SideEffectHandler 在 Worker 中使用，这里创建空实现）
	sideEffectHandler := userEvent.NewSideEffectHandler(infra.Logger, nil)

	// 创建审计日志订阅器
	auditLogRepo := repository.NewAuditLogRepository(daoQuery)
	auditSubscriber := eventHandler.NewAuditSubscriber(auditLogRepo, infra.Snowflake)

	// 创建登录日志订阅器
	loginLogRepo := repository.NewLoginLogRepository(daoQuery)
	uaParser := &userAgentParserAdapter{parser: useragent.NewParser()}
	loginLogSubscriber := eventHandler.NewLoginLogSubscriber(loginLogRepo, infra.Snowflake, uaParser)

	return &UserModule{
		infra:              infra,
		routes:             routes,
		handler:            handler,
		sideEffectHandler:  sideEffectHandler,
		auditSubscriber:    auditSubscriber,
		loginLogSubscriber: loginLogSubscriber,
	}
}

// Name 返回模块名称
// 实现 bootstrap.Module 接口
func (m *UserModule) Name() string {
	return "user"
}

// RegisterHTTP 注册 HTTP 路由
// 实现 bootstrap.HTTPModule 接口
func (m *UserModule) RegisterHTTP(group *gin.RouterGroup) {
	m.routes.Register(group)
}

// RegisterSubscriptions 注册事件订阅
// 实现 bootstrap.EventModule 接口
func (m *UserModule) RegisterSubscriptions(bus sharedAggregate.EventBus) {
	subscriber := eventHandler.NewSubscriber(bus)
	subscriber.SubscribeAll(&eventHandler.Dependencies{
		AuditSubscriber:       m.auditSubscriber,
		LoginLogSubscriber:    m.loginLogSubscriber,
		UserSideEffectHandler: m.sideEffectHandler,
	})
}

// userAgentParserAdapter 适配 useragent.Parser 到 event.UserAgentParser 接口
type userAgentParserAdapter struct {
	parser *useragent.Parser
}

func (a *userAgentParserAdapter) Parse(ua string) eventHandler.DeviceInfo {
	info := a.parser.Parse(ua)
	return eventHandler.DeviceInfo{
		DeviceType: info.DeviceType,
		OS:         info.OS,
		Browser:    info.Browser,
	}
}
