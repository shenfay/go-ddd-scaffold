package bootstrap

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/commands"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/queries"
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

	// === 3. 创建 CQRS Handlers（直接赋值给 Bootstrap 字段）===
	b.user.updateHandler = commands.NewUpdateUserHandler(userRepo, eventPublisher)

	b.user.getHandler = queries.NewGetUserHandler(userRepo)
	b.user.listHandler = queries.NewListUsersHandler(userRepo)

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

// 确保 InMemoryEventPublisher 实现 commands.EventPublisher 接口
var _ commands.EventPublisher = (*InMemoryEventPublisher)(nil)
