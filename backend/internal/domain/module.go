package domain

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/shenfay/go-ddd-scaffold/internal/container"
)

// Module 领域模块接口
type Module interface {
	// Name 返回领域名称（用于日志和配置）
	Name() string

	// Initialize 初始化领域模块
	// 在此方法中创建仓储、服务、注册路由
	Initialize(ctx context.Context, c container.Container) error
}

// LifecycleModule 可选的生命周期接口
type LifecycleModule interface {
	Module
	// Start 领域启动时的钩子（可选）
	Start(ctx context.Context, c container.Container) error
	// Stop 领域停止时的钩子（可选）
	Stop(ctx context.Context, c container.Container) error
}

// moduleRegistry 领域模块注册表
type moduleRegistry struct {
	modules map[string]Module
}

var globalRegistry = &moduleRegistry{
	modules: make(map[string]Module),
}

// Register 注册领域模块（通常在 init 中调用）
func Register(module Module) {
	if module == nil {
		panic("domain: register module is nil")
	}
	name := module.Name()
	if name == "" {
		panic("domain: module name is empty")
	}
	if _, exists := globalRegistry.modules[name]; exists {
		panic(fmt.Sprintf("domain: module already registered: %s", name))
	}
	globalRegistry.modules[name] = module
}

// GetModules 获取所有已注册的模块
func GetModules() map[string]Module {
	return globalRegistry.modules
}

// InitializeAll 初始化所有已注册的领域模块
func InitializeAll(ctx context.Context, c container.Container) error {
	logger := c.GetLogger("system")

	for name, module := range globalRegistry.modules {
		logger.Info("Initializing domain module",
			zap.String("module", name))

		if err := module.Initialize(ctx, c); err != nil {
			return fmt.Errorf("failed to initialize module %s: %w", name, err)
		}

		// 如果模块实现了生命周期接口，注册启动/停止钩子
		if lifecycleModule, ok := module.(LifecycleModule); ok {
			c.(container.ContainerInternal).OnStart(func(ctx context.Context) error {
				logger.Info("starting domain module", zap.String("module", name))
				return lifecycleModule.Start(ctx, c)
			})

			c.(container.ContainerInternal).OnStop(func(ctx context.Context) error {
				logger.Info("stopping domain module", zap.String("module", name))
				return lifecycleModule.Stop(ctx, c)
			})
		}
	}

	logger.Info("all domain modules initialized successfully",
		zap.Int("count", len(globalRegistry.modules)))
	return nil
}
