//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/application/user/service"
	appService "go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
)

// UserModuleDeps User 模块的依赖声明
type UserModuleDeps struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	JWTService entity.JWTService
	EventBus   *event.EventBus
}

// InitializeUserModule 初始化 User 模块
// 返回值用于路由注册
func InitializeUserModule(
	db *gorm.DB,
	logger *zap.Logger,
	jwtService entity.JWTService,
	eventBus *event.EventBus,
) (
	*service.Service,
	error,
) {
	wire.Build(
		// Repositories
		repo.NewUserDAORepository,
		repo.NewTenantDAORepository,
		repo.NewTenantMemberDAORepository,

		// 事件总线适配器
		newUserEventBusAdapter,

		// Service
		service.NewService,
	)

	return nil, nil
}

// userEventBusAdapter User 模块的事件总线适配器
type userEventBusAdapter struct {
	bus *event.EventBus
}

func newUserEventBusAdapter(bus *event.EventBus) appService.EventBus {
	return &userEventBusAdapter{bus: bus}
}

func (a *userEventBusAdapter) Publish(event interface{}) error {
	if event == nil {
		return nil
	}
	return nil
}
