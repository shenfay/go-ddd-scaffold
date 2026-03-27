package modules

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/cmd/app/factory"
	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
	v1 "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers/v1"
	userHTTP "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers/v1/user"
)

// UserModule 用户模块
// 实现 bootstrap.Module、bootstrap.HTTPModule 和 bootstrap.EventModule 接口
type UserModule struct {
	infra         *factory.Infrastructure
	uow           application.UnitOfWork
	uowWithEvents application.UnitOfWorkWithEvents
	logWriter     *application.ActivityLogWriter
}

// NewUserModule 创建用户模块
// 内部自行构建完整依赖链
func NewUserModule(infra *factory.Infrastructure) *UserModule {
	// 1. 创建 DAO Query
	daoQuery := dao.Use(infra.DB)

	// 2. 创建 UnitOfWork（用于普通查询）
	uow := application.NewUnitOfWork(infra.DB, daoQuery)

	// 3. 创建 UnitOfWorkWithEvents（使用选项模式，用于需要发布事件的场景）
	uowWithEvents := application.NewUnitOfWork(
		infra.DB,
		daoQuery,
		application.WithEventPublisher(infra.EventPublisher),
	).(application.UnitOfWorkWithEvents)

	// 4. 创建 ActivityLogWriter（使用 UnitOfWorkWithEvents，确保在事务中写入）
	logWriter := application.NewActivityLogWriter(uowWithEvents, infra.Logger)

	// 5. 创建 JWTService（仅用于 Token 刷新等场景）
	jwtSvc := auth.NewJWTService(
		infra.Config.JWT.Secret,
		infra.Config.JWT.AccessExpire,
		infra.Config.JWT.RefreshExpire,
		"go-ddd-scaffold",
	)
	jwtSvc.SetRedisClient(infra.Redis)

	return &UserModule{
		infra:         infra,
		uow:           uow,
		uowWithEvents: uowWithEvents,
		logWriter:     logWriter,
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
	// 创建响应处理器
	respHandler := handlers.NewHandler(m.infra.ErrorMapper)
	router := v1.NewRouter()

	// 提前创建所有 Handler（只创建一次，避免重复）
	getUserHandler := userHTTP.NewGetUserHandler(
		usecase.NewGetUserUseCase(m.uow),
		respHandler,
	)
	updateProfileHandler := userHTTP.NewUpdateProfileHandler(
		usecase.NewUpdateProfileUseCase(m.uowWithEvents, m.logWriter),
		respHandler,
	)
	changePasswordHandler := userHTTP.NewChangePasswordHandler(
		usecase.NewChangePasswordUseCase(m.uowWithEvents,
			service.NewBcryptPasswordHasher(m.infra.Config.Security.PasswordHasher.Cost),
			auth.NewDefaultPasswordPolicy(service.PasswordPolicyConfig{
				MinLength:           m.infra.Config.Security.PasswordPolicy.MinLength,
				MaxLength:           m.infra.Config.Security.PasswordPolicy.MaxLength,
				RequireUppercase:    m.infra.Config.Security.PasswordPolicy.RequireUppercase,
				RequireLowercase:    m.infra.Config.Security.PasswordPolicy.RequireLowercase,
				RequireDigits:       m.infra.Config.Security.PasswordPolicy.RequireDigits,
				RequireSpecialChars: m.infra.Config.Security.PasswordPolicy.RequireSpecialChars,
				SpecialChars:        m.infra.Config.Security.PasswordPolicy.SpecialChars,
				DisallowCommon:      m.infra.Config.Security.PasswordPolicy.DisallowCommon,
			}),
			m.logWriter,
		),
		respHandler,
	)

	// 注册用户路由（传入 Handler 集合）
	router.RegisterUserRoutes(group, v1.UserHandlers{
		GetUser:        getUserHandler.Handle,
		UpdateProfile:  updateProfileHandler.Handle,
		ChangePassword: changePasswordHandler.Handle,
	})
}
