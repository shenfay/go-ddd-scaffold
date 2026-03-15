package bootstrap

import (
	"context"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	userDomain "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	repositoryPkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
	"go.uber.org/zap"
)

// initUserDomain 初始化用户领域
func (b *Bootstrap) initUserDomain(ctx context.Context) error {
	b.logger.Info("Initializing user domain...")

	baseLogger := b.logger.Named("user")

	// === 1. 创建基础设施服务 ===
	eventPublisher := NewInMemoryEventPublisher(baseLogger.Named("events"))

	// === 2. 创建仓储层 ===
	db := b.container.GetGormDB()

	// 初始化 DAO（必须在使用 repository 之前）
	dao.SetDefault(db)

	userRepo := repositoryPkg.NewUserRepository(db)

	// === 3. 创建应用服务（统一入口）===
	passwordHasher := userDomain.NewBcryptPasswordHasher(12)
	b.user.service = userApp.NewUserService(userRepo, eventPublisher, passwordHasher)

	// === 4. 创建领域事件处理器 ===
	b.user.eventHandler = userApp.NewUserEventHandler(baseLogger.Named("events"))

	b.logger.Info("User domain initialized successfully")
	return nil
}

// InMemoryEventPublisher 内存事件发布器（临时实现）
type InMemoryEventPublisher struct {
	logger *zap.Logger
}

// NewInMemoryEventPublisher 创建事件发布器
func NewInMemoryEventPublisher(logger *zap.Logger) *InMemoryEventPublisher {
	return &InMemoryEventPublisher{logger: logger}
}

// Publish 发布事件
func (p *InMemoryEventPublisher) Publish(ctx context.Context, event ddd.DomainEvent) error {
	// TODO: 实际应该发布到 Redis 或消息队列
	p.logger.Info("Domain event published", zap.Any("event", event))
	return nil
}

// 确保 InMemoryEventPublisher 实现 ddd.EventPublisher 接口
var _ ddd.EventPublisher = (*InMemoryEventPublisher)(nil)
