package auth

import (
	"context"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
	"go.uber.org/zap"
)

// AuthService 认证应用服务接口
type AuthService interface {
	// AuthenticateUser 用户认证（登录）
	AuthenticateUser(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error)
	// RegisterUser 用户注册
	RegisterUser(ctx context.Context, cmd *RegisterCommand) (*RegisterResult, error)
	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResult, error)
	// Logout 用户登出
	Logout(ctx context.Context, cmd *LogoutCommand) (*LogoutResult, error)
	// GetUserByID 根据 ID 获取用户信息
	GetUserByID(ctx context.Context, userID int64) (*UserInfoResult, error)
}

// AuthServiceImpl 认证应用服务实现
type AuthServiceImpl struct {
	userRepo       user.UserRepository
	passwordHasher user.PasswordHasher
	tokenService   auth.TokenService
	eventPublisher kernel.EventPublisher
	idGenerator    *snowflake.Node
	logger         *zap.Logger
}

// NewAuthService 创建认证应用服务
// 如果 logger 为 nil，则使用默认的全局 logger
func NewAuthService(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	tokenService auth.TokenService,
	eventPublisher kernel.EventPublisher,
	idGenerator *snowflake.Node,
	logger *zap.Logger,
) *AuthServiceImpl {
	if logger == nil {
		logger = zap.L().Named("auth")
	}
	return &AuthServiceImpl{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
		idGenerator:    idGenerator,
		logger:         logger,
	}
}

// AuthenticateUser 认证用户
func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error) {
	// 1. 查找用户
	var foundUser *user.User
	var err error

	// 先尝试作为邮箱查找
	foundUser, err = s.userRepo.FindByEmail(ctx, cmd.Identifier)
	if err != nil {
		// 再尝试作为用户名查找
		foundUser, err = s.userRepo.FindByUsername(ctx, cmd.Identifier)
		if err != nil {
			return nil, kernel.NewBusinessError(kernel.CodeInvalidCredentials, "用户名或密码错误")
		}
	}

	// 2. 检查账户状态
	if !foundUser.CanLogin() {
		switch foundUser.Status() {
		case user.UserStatusInactive:
			return nil, kernel.NewBusinessError(kernel.CodeAccountDisabled, "账户已被禁用")
		case user.UserStatusLocked:
			return nil, kernel.NewBusinessError(kernel.CodeAccountLocked, "账户已被锁定")
		default:
			return nil, kernel.NewBusinessError(kernel.CodeAccountCannotLogin, "账户无法登录")
		}
	}

	// 3. 验证密码
	if !s.passwordHasher.Verify(cmd.Password, foundUser.Password().Value()) {
		// 记录失败登录
		foundUser.RecordFailedLogin(cmd.IPAddress, cmd.UserAgent, "invalid_password")
		_ = s.userRepo.Save(ctx, foundUser)
		return nil, kernel.NewBusinessError(kernel.CodeInvalidCredentials, "用户名或密码错误")
	}

	// 4. 生成令牌对
	tokenPair, err := s.tokenService.GenerateTokenPair(
		foundUser.ID().(user.UserID).Int64(),
		foundUser.Username().Value(),
		foundUser.Email().Value(),
	)
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeTokenGenerationFailed, "令牌生成失败")
	}

	// 5. 记录成功登录
	foundUser.RecordLogin(cmd.IPAddress, cmd.UserAgent)

	// 6. 获取未提交的事件并先发布到 EventBus
	events := foundUser.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			s.logger.Error("Failed to publish domain event", zap.String("event_name", event.EventName()), zap.Error(err))
		}
	}

	// 7. 保存用户（会保存领域事件到 domain_events 表）
	if err := s.userRepo.Save(ctx, foundUser); err != nil {
		return nil, err
	}

	// 7. 返回结果
	expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())

	return &AuthenticateResult{
		UserID:       foundUser.ID().(user.UserID).String(),
		Username:     foundUser.Username().Value(),
		Email:        foundUser.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RegisterUser 注册用户
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterCommand) (*RegisterResult, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil && !errors.Is(err, kernel.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, kernel.NewBusinessError(kernel.CodeUsernameExists, "用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingUser, err = s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, kernel.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, kernel.NewBusinessError(kernel.CodeEmailExists, "邮箱已被注册")
	}

	// 3. 哈希密码
	hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	// 4. 使用 Snowflake 生成唯一 ID
	userID, err := s.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	// 5. 创建用户实体
	newUser, err := user.NewUser(cmd.Username, cmd.Email, hashedPassword, func() int64 {
		return userID
	})
	if err != nil {
		return nil, err
	}

	// 5. 保存用户
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	newUser.ClearUncommittedEvents()

	// 7. 返回结果
	return &RegisterResult{
		UserID:   newUser.ID().(user.UserID).String(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResult, error) {
	// 1. 验证刷新令牌并提取用户 ID
	claims, err := s.tokenService.ValidateToken(cmd.RefreshToken)
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeInvalidToken, "无效的刷新令牌")
	}

	// 2. 查找用户
	foundUser, err := s.userRepo.FindByID(ctx, user.NewUserID(claims.UserID))
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeUserNotFound, "用户不存在")
	}

	// 3. 检查账户状态
	if !foundUser.CanLogin() {
		switch foundUser.Status() {
		case user.UserStatusInactive:
			return nil, kernel.NewBusinessError(kernel.CodeAccountDisabled, "账户已被禁用")
		case user.UserStatusLocked:
			return nil, kernel.NewBusinessError(kernel.CodeAccountLocked, "账户已被锁定")
		default:
			return nil, kernel.NewBusinessError(kernel.CodeAccountCannotLogin, "账户无法登录")
		}
	}

	// 4. ⭐ 令牌轮换策略
	if cmd.CurrentToken != "" {
		// 严格模式：将旧 access_token 加入黑名单
		oldClaims, parseErr := s.tokenService.ParseAccessToken(cmd.CurrentToken)
		if parseErr == nil && oldClaims != nil {
			blacklistErr := s.tokenService.BlacklistToken(cmd.CurrentToken, oldClaims.ExpiresAt)
			if blacklistErr != nil {
				s.logger.Warn("failed to blacklist current access token",
					zap.Int64("user_id", claims.UserID),
					zap.Error(blacklistErr),
				)
			} else {
				s.logger.Debug("current access token blacklisted (strict mode)",
					zap.Int64("user_id", claims.UserID),
				)
			}
		}
	} else {
		// 宽松模式：记录日志
		s.logger.Debug("refreshing token without current_token (relaxed mode)",
			zap.Int64("user_id", claims.UserID),
		)
	}

	// 5. 生成新的令牌对
	tokenPair, err := s.tokenService.GenerateTokenPair(foundUser.ID().(user.UserID).Int64(), foundUser.Username().Value(), foundUser.Email().Value())
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeTokenGenerationFailed, "令牌生成失败")
	}

	// 5. 返回结果
	expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())

	return &RefreshTokenResult{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Logout 用户登出
func (s *AuthServiceImpl) Logout(ctx context.Context, cmd *LogoutCommand) (*LogoutResult, error) {
	// 1. 查找用户（验证用户存在）
	_, err := s.userRepo.FindByID(ctx, user.NewUserID(cmd.UserID))
	if err != nil {
		// 用户不存在也返回成功（幂等性）
		return &LogoutResult{Success: true}, nil
	}

	// 2. 将令牌加入黑名单（如果提供了 access_token）
	if cmd.AccessToken != "" {
		// 解析令牌获取过期时间
		claims, parseErr := s.tokenService.ParseAccessToken(cmd.AccessToken)
		if parseErr == nil && claims != nil {
			// 成功解析，加入黑名单
			blacklistErr := s.tokenService.BlacklistToken(cmd.AccessToken, claims.ExpiresAt)
			if blacklistErr != nil {
				// 记录错误但不影响主流程
				s.logger.Warn("failed to blacklist token on logout",
					zap.Int64("user_id", cmd.UserID),
					zap.String("ip_address", cmd.IPAddress),
					zap.Error(blacklistErr),
				)
			} else {
				s.logger.Debug("token blacklisted successfully",
					zap.Int64("user_id", cmd.UserID),
					zap.String("ip_address", cmd.IPAddress),
				)
			}
		}
	}

	// 3. 返回成功
	return &LogoutResult{Success: true}, nil
}

// GetUserByID 根据 ID 获取用户信息
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, userID int64) (*UserInfoResult, error) {
	foundUser, err := s.userRepo.FindByID(ctx, user.NewUserID(userID))
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeUserNotFound, "用户不存在")
	}

	return &UserInfoResult{
		ID:          foundUser.ID().(user.UserID).Int64(),
		Username:    foundUser.Username().Value(),
		Email:       foundUser.Email().Value(),
		DisplayName: foundUser.DisplayName(),
		Status:      foundUser.Status().String(),
	}, nil
}
