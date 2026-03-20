package provider

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	authInfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
)

// DomainInfrastructureProvider 领域基础设施提供者
// 负责提供各领域特定的基础设施组件（非全局共享）
// 使用懒加载和缓存机制，确保每个组件只创建一次
type DomainInfrastructureProvider struct {
	config *config.AppConfig

	// === 实例缓存（懒加载 + 单例）===
	passwordHasher service.PasswordHasher
	passwordPolicy service.PasswordPolicy
	jwtService     *authInfra.JWTService
	redisClient    *redis.Client // 用于注入 JWT 服务

	mu sync.Mutex // 并发锁，保护缓存
}

// NewDomainInfrastructureProvider 创建领域基础设施提供者
func NewDomainInfrastructureProvider(cfg *config.AppConfig) *DomainInfrastructureProvider {
	return &DomainInfrastructureProvider{
		config: cfg,
	}
}

// SetRedisClient 设置 Redis 客户端（由外部在创建后注入）
func (p *DomainInfrastructureProvider) SetRedisClient(client *redis.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.redisClient = client
}

// GetPasswordHasher 获取密码哈希器（用户/认证领域使用）
// 懒加载：首次调用时创建，后续调用返回缓存实例
func (p *DomainInfrastructureProvider) GetPasswordHasher() service.PasswordHasher {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.passwordHasher == nil {
		securityConfig := p.config.Security
		p.passwordHasher = service.NewBcryptPasswordHasher(
			securityConfig.PasswordHasher.Cost,
		)
	}
	return p.passwordHasher
}

// GetPasswordPolicy 获取密码策略（用户/认证领域使用）
// 懒加载：首次调用时创建，后续调用返回缓存实例
func (p *DomainInfrastructureProvider) GetPasswordPolicy() service.PasswordPolicy {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.passwordPolicy == nil {
		securityConfig := p.config.Security
		p.passwordPolicy = authInfra.NewDefaultPasswordPolicy(service.PasswordPolicyConfig{
			MinLength:           securityConfig.PasswordPolicy.MinLength,
			MaxLength:           securityConfig.PasswordPolicy.MaxLength,
			RequireUppercase:    securityConfig.PasswordPolicy.RequireUppercase,
			RequireLowercase:    securityConfig.PasswordPolicy.RequireLowercase,
			RequireDigits:       securityConfig.PasswordPolicy.RequireDigits,
			RequireSpecialChars: securityConfig.PasswordPolicy.RequireSpecialChars,
			SpecialChars:        securityConfig.PasswordPolicy.SpecialChars,
			DisallowCommon:      securityConfig.PasswordPolicy.DisallowCommon,
		})
	}
	return p.passwordPolicy
}

// GetJWTService 获取 JWT 服务（认证领域使用）
// 懒加载：首次调用时创建，后续调用返回缓存实例
// redisClient 参数：如果提供了新的 Redis 客户端，会更新到 JWT 服务中
func (p *DomainInfrastructureProvider) GetJWTService(redisClient *redis.Client) *authInfra.JWTService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.jwtService == nil {
		jwtConfig := p.config.JWT
		p.jwtService = authInfra.NewJWTService(
			jwtConfig.Secret,
			jwtConfig.AccessExpire,
			jwtConfig.RefreshExpire,
			"go-ddd-scaffold", // issuer
		)
	}

	// 如果提供了新的 Redis 客户端或者还没有注入，则更新
	if redisClient != nil {
		p.jwtService.SetRedisClient(redisClient)
	} else if p.redisClient != nil {
		// 使用内部缓存的 Redis 客户端
		p.jwtService.SetRedisClient(p.redisClient)
	}

	return p.jwtService
}
