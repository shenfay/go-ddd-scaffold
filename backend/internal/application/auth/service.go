package auth

import (
	"context"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
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
	eventPublisher ddd.EventPublisher
	logger         *zap.Logger
}

// NewAuthService 创建认证应用服务
// 如果 logger 为 nil，则使用默认的全局 logger
func NewAuthService(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	tokenService auth.TokenService,
	eventPublisher ddd.EventPublisher,
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
			return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "用户名或密码错误")
		}
	}

	// 2. 检查账户状态
	if !foundUser.CanLogin() {
		switch foundUser.Status() {
		case user.UserStatusInactive:
			return nil, ddd.NewBusinessError("ACCOUNT_DISABLED", "账户已被禁用")
		case user.UserStatusLocked:
			return nil, ddd.NewBusinessError("ACCOUNT_LOCKED", "账户已被锁定")
		default:
			return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "账户无法登录")
		}
	}

	// 3. 验证密码
	if !s.passwordHasher.Verify(cmd.Password, foundUser.Password().Value()) {
		// 记录失败登录
		foundUser.RecordFailedLogin(cmd.IPAddress, cmd.UserAgent, "invalid_password")
		_ = s.userRepo.Save(ctx, foundUser)
		return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "用户名或密码错误")
	}

	// 4. 生成令牌对
	tokenPair, err := s.tokenService.GenerateTokenPair(
		foundUser.ID().(user.UserID).Int64(),
		foundUser.Username().Value(),
		foundUser.Email().Value(),
	)
	if err != nil {
		return nil, ddd.NewBusinessError("TOKEN_GENERATION_FAILED", "令牌生成失败")
	}

	// 5. 记录成功登录（在保存之前发布领域事件）
	foundUser.RecordLogin(cmd.IPAddress, cmd.UserAgent)
	// 先发布领域事件，因为 Save 后事件会被清除
	events := foundUser.GetUncommittedEvents()
	s.logger.Debug("Publishing domain events", zap.Int("event_count", len(events)), zap.String("username", foundUser.Username().Value()))
	for _, event := range events {
		s.logger.Debug("Publishing event", zap.String("event_name", event.EventName()))
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			s.logger.Error("Failed to publish event", zap.String("event_name", event.EventName()), zap.Error(err))
			// 记录错误但不中断主流程
		}
	}
	foundUser.ClearUncommittedEvents()

	// 6. 保存用户
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
	if err != nil && !errors.Is(err, ddd.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ddd.NewBusinessError("USERNAME_EXISTS", "用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingUser, err = s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, ddd.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ddd.NewBusinessError("EMAIL_EXISTS", "邮箱已被注册")
	}

	// 3. 哈希密码
	hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户实体
	newUser, err := user.NewUser(cmd.Username, cmd.Email, hashedPassword, func() int64 {
		return time.Now().UnixNano() // 临时实现，实际应该使用 Snowflake
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
		return nil, ddd.NewBusinessError("INVALID_TOKEN", "无效的刷新令牌")
	}

	// 2. 查找用户
	foundUser, err := s.userRepo.FindByID(ctx, user.NewUserID(claims.UserID))
	if err != nil {
		return nil, ddd.NewBusinessError("USER_NOT_FOUND", "用户不存在")
	}

	// 3. 检查账户状态
	if !foundUser.CanLogin() {
		switch foundUser.Status() {
		case user.UserStatusInactive:
			return nil, ddd.NewBusinessError("ACCOUNT_DISABLED", "账户已被禁用")
		case user.UserStatusLocked:
			return nil, ddd.NewBusinessError("ACCOUNT_LOCKED", "账户已被锁定")
		default:
			return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "账户无法登录")
		}
	}

	// 4. ⭐ 令牌轮换：生成新 token 后，旧 Refresh Token 自动失效
	// 注意：我们采用宽松策略 - 不主动将旧 AT 加入黑名单
	// 而是依赖短生命周期（5-15 分钟）自然过期
	// 如果提供了当前 Access Token，可以主动加入黑名单（严格模式）
	if cmd.CurrentToken != "" {
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
	}

	// 5. 生成新的令牌对
	tokenPair, err := s.tokenService.GenerateTokenPair(foundUser.ID().(user.UserID).Int64(), foundUser.Username().Value(), foundUser.Email().Value())
	if err != nil {
		return nil, ddd.NewBusinessError("TOKEN_GENERATION_FAILED", "令牌生成失败")
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
		return nil, ddd.NewBusinessError("USER_NOT_FOUND", "用户不存在")
	}

	return &UserInfoResult{
		ID:          foundUser.ID().(user.UserID).Int64(),
		Username:    foundUser.Username().Value(),
		Email:       foundUser.Email().Value(),
		DisplayName: foundUser.DisplayName(),
		Status:      foundUser.Status().String(),
	}, nil
}
