package bootstrap

import (
	"context"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	userDomain "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
)

// initUserDomain 初始化用户领域
func (b *Bootstrap) initUserDomain(ctx context.Context) error {
	b.logger.Info("Initializing user domain...")

	baseLogger := b.logger.Named("user")

	// === 1. 创建基础设施服务 ===
	// 使用 Bootstrap 中创建的事件总线
	eventPublisher := b.eventBus

	// === 2. 从容器获取仓储层 ===
	userRepo := b.container.GetUserRepo()

	// === 3. 创建应用服务（统一入口）===
	// 从配置获取密码策略和哈希配置
	securityConfig := b.config.Security

	passwordHasher := userDomain.NewBcryptPasswordHasher(securityConfig.PasswordHasher.Cost)
	passwordPolicy := auth.NewDefaultPasswordPolicy(userDomain.PasswordPolicyConfig{
		MinLength:           securityConfig.PasswordPolicy.MinLength,
		MaxLength:           securityConfig.PasswordPolicy.MaxLength,
		RequireUppercase:    securityConfig.PasswordPolicy.RequireUppercase,
		RequireLowercase:    securityConfig.PasswordPolicy.RequireLowercase,
		RequireDigits:       securityConfig.PasswordPolicy.RequireDigits,
		RequireSpecialChars: securityConfig.PasswordPolicy.RequireSpecialChars,
		SpecialChars:        securityConfig.PasswordPolicy.SpecialChars,
		DisallowCommon:      securityConfig.PasswordPolicy.DisallowCommon,
	})
	b.user.service = userApp.NewUserService(
		userRepo,
		eventPublisher,
		passwordHasher,
		passwordPolicy,
		b.auth.jwtService,
		b.container.GetSnowflake(),
	)

	// === 4. 创建领域副作用处理器 ===
	b.user.sideEffectHandler = userDomain.NewSideEffectHandler(baseLogger.Named("events"))

	b.logger.Info("User domain initialized successfully")
	return nil
}
