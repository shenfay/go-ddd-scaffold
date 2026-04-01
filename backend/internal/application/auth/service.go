package auth

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
	ports_idgen "github.com/shenfay/go-ddd-scaffold/internal/application/ports/idgen"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"go.uber.org/zap"
)

// 错误码常量
const (
	CodeAccountLocked  = user.CodeAccountLocked
	CodeUsernameExists = user.CodeUsernameExists
	CodeEmailExists    = user.CodeEmailExists
	CodeUserNotFound   = user.CodeUserNotFound
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
	uow            application.UnitOfWork
	passwordHasher service.PasswordHasher
	tokenService   ports_auth.TokenService
	eventPublisher common.EventPublisher
	idGenerator    ports_idgen.Generator
	logger         *zap.Logger
	activityRepo   model.ActivityLogRepository
}

// NewAuthService 创建认证应用服务
// 如果 logger 为 nil，则使用默认的全局 logger
func NewAuthService(
	uow application.UnitOfWork,
	passwordHasher service.PasswordHasher,
	tokenService ports_auth.TokenService,
	eventPublisher common.EventPublisher,
	idGenerator ports_idgen.Generator,
	activityRepo model.ActivityLogRepository,
	logger *zap.Logger,
) *AuthServiceImpl {
	if logger == nil {
		logger = zap.L().Named("auth")
	}
	return &AuthServiceImpl{
		uow:            uow,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
		idGenerator:    idGenerator,
		activityRepo:   activityRepo,
		logger:         logger,
	}
}

// AuthenticateUser 认证用户
func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error) {
	var authResult *AuthenticateResult
	var loginEvent common.DomainEvent
	var savedUserID int64
	var savedUsername string

	// 在事务中执行认证
	err := s.uow.Transaction(ctx, func(ctx context.Context) error {
		// 1. 查找用户
		userRepo := s.uow.UserRepository()
		foundUser, err := userRepo.FindByEmail(ctx, cmd.Identifier)
		if err != nil {
			// 再尝试作为用户名查找
			foundUser, err = userRepo.FindByUsername(ctx, cmd.Identifier)
			if err != nil {
				return common.NewBusinessError(common.CodeInvalidCredentials, "用户名或密码错误")
			}
		}

		// 2. 检查账户状态
		if err := s.checkAccountStatus(foundUser); err != nil {
			return err
		}

		// 3. 验证密码
		if !s.passwordHasher.Verify(cmd.Password, foundUser.Password().Value()) {
			return common.NewBusinessError(common.CodeInvalidCredentials, "用户名或密码错误")
		}

		// 4. 生成令牌对
		tokenPair, err := s.tokenService.GenerateTokenPair(
			foundUser.ID().(vo.UserID).Int64(),
			foundUser.Username().Value(),
			foundUser.Email().Value(),
		)
		if err != nil {
			return common.NewBusinessError(common.CodeTokenGenerationFailed, "令牌生成失败")
		}

		// 5. 创建登录成功事件
		loginEvent = userEvent.NewUserLoggedInEvent(
			foundUser.ID().(vo.UserID),
			cmd.IPAddress,
			cmd.UserAgent,
			"", // location - 可通过 IP 查询
			"desktop",
			"",
			"password",
			true,
		)

		// 6. 保存用户信息用于后续活动日志
		savedUserID = foundUser.ID().(vo.UserID).Int64()
		savedUsername = foundUser.Username().Value()

		// 7. 返回结果
		expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())
		authResult = &AuthenticateResult{
			UserID:       foundUser.ID().(vo.UserID).String(),
			Username:     foundUser.Username().Value(),
			Email:        foundUser.Email().Value(),
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    expiresIn,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 6. 在事务内同步写入活动日志（保证审计可靠性）
	activityLog := &model.ActivityLog{
		ID:     idgen.Generate(),
		UserID: savedUserID,
		Action: model.ActivityUserLoggedIn,
		Status: model.ActivityStatusSuccess,
		Metadata: map[string]interface{}{
			"username":   savedUsername,
			"ip_address": cmd.IPAddress,
			"user_agent": cmd.UserAgent,
		},
		OccurredAt: time.Now(),
		CreatedAt:  time.Now(),
	}
	if err := s.activityRepo.Save(ctx, activityLog); err != nil {
		s.logger.Warn("活动日志写入失败", zap.Error(err))
		// 不因活动日志失败影响主流程
	}

	// 7. 事务成功后，异步发布登录事件到 Redis 队列
	if loginEvent != nil {
		if pubErr := s.eventPublisher.Publish(ctx, loginEvent); pubErr != nil {
			s.logger.Warn("登录事件发布失败", zap.Error(pubErr))
		}
	}

	return authResult, nil
}

// checkAccountStatus 检查账户状态
func (s *AuthServiceImpl) checkAccountStatus(user *user.User) error {
	if !user.CanLogin() {
		switch user.Status() {
		case vo.UserStatusInactive:
			return common.NewBusinessError(common.CodeAccountDisabled, "账户已被禁用")
		case vo.UserStatusLocked:
			return common.NewBusinessError(CodeAccountLocked, "账户已被锁定")
		default:
			return common.NewBusinessError(common.CodeAccountCannotLogin, "账户无法登录")
		}
	}
	return nil
}

// RegisterUser 注册用户 - 简化版
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterCommand) (*RegisterResult, error) {
	var newUser *user.User

	err := s.uow.Transaction(ctx, func(ctx context.Context) error {
		userRepo := s.uow.UserRepository()

		// 1. 检查用户名唯一性
		if existingUser, _ := userRepo.FindByUsername(ctx, cmd.Username); existingUser != nil {
			return common.NewBusinessError(CodeUsernameExists, "用户名已存在")
		}

		// 2. 检查邮箱唯一性
		if existingUser, _ := userRepo.FindByEmail(ctx, cmd.Email); existingUser != nil {
			return common.NewBusinessError(CodeEmailExists, "邮箱已被注册")
		}

		// 3. 哈希密码
		hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
		if err != nil {
			return err
		}

		// 4. 生成 ID 并创建用户实体
		userID, err := s.idGenerator.Generate()
		if err != nil {
			return err
		}

		newUser, err = user.NewUser(cmd.Username, cmd.Email, hashedPassword, func() int64 {
			return userID
		})
		if err != nil {
			return err
		}

		// 5. 保存用户
		return userRepo.Save(ctx, newUser)
	})

	if err != nil {
		return nil, err
	}

	// 6. 发布领域事件（事务外）
	s.publishUserEvents(ctx, newUser)

	// 7. 返回结果
	return &RegisterResult{
		UserID:   newUser.ID().(vo.UserID).String(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}

// publishUserEvents 发布用户的领域事件
func (s *AuthServiceImpl) publishUserEvents(ctx context.Context, user *user.User) {
	events := user.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			s.logger.Error("Failed to publish event",
				zap.String("event_type", event.EventName()),
				zap.Error(err),
			)
		}
	}
	user.ClearUncommittedEvents()
}

// RefreshToken 刷新令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResult, error) {
	// 1. 验证刷新令牌并提取用户 ID
	claims, err := s.tokenService.ValidateToken(cmd.RefreshToken)
	if err != nil {
		return nil, common.NewBusinessError(common.CodeInvalidToken, "无效的刷新令牌")
	}

	// 2. 查找用户（只读操作，不需要事务）
	userRepo := s.uow.UserRepository()
	foundUser, err := userRepo.FindByID(ctx, vo.NewUserID(claims.UserID))
	if err != nil {
		return nil, common.NewBusinessError(CodeUserNotFound, "用户不存在")
	}

	// 3. 检查账户状态
	if err := s.checkAccountStatus(foundUser); err != nil {
		return nil, err
	}

	// 4. ⭐ 令牌轮换策略
	s.handleTokenRotation(cmd.CurrentToken, claims)

	// 5. 生成新的令牌对
	tokenPair, err := s.tokenService.GenerateTokenPair(foundUser.ID().(vo.UserID).Int64(), foundUser.Username().Value(), foundUser.Email().Value())
	if err != nil {
		return nil, common.NewBusinessError(common.CodeTokenGenerationFailed, "令牌生成失败")
	}

	// 6. 返回结果
	expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())

	return &RefreshTokenResult{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// handleTokenRotation 处理令牌轮换策略
func (s *AuthServiceImpl) handleTokenRotation(currentToken string, claims *ports_auth.TokenClaims) {
	if currentToken != "" {
		// 严格模式：将旧 access_token 加入黑名单
		oldClaims, parseErr := s.tokenService.ParseAccessToken(currentToken)
		if parseErr == nil && oldClaims != nil {
			blacklistErr := s.tokenService.BlacklistToken(currentToken, oldClaims.ExpiresAt)
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
}

// Logout 用户登出
func (s *AuthServiceImpl) Logout(ctx context.Context, cmd *LogoutCommand) (*LogoutResult, error) {
	// 1. 查找用户（验证用户存在）- 只读操作，不需要事务
	userRepo := s.uow.UserRepository()
	_, err := userRepo.FindByID(ctx, vo.NewUserID(cmd.UserID))
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
	// 查找用户（只读操作，不需要事务）
	userRepo := s.uow.UserRepository()
	foundUser, err := userRepo.FindByID(ctx, vo.NewUserID(userID))
	if err != nil {
		return nil, common.NewBusinessError(CodeUserNotFound, "用户不存在")
	}

	return &UserInfoResult{
		ID:          foundUser.ID().(vo.UserID).Int64(),
		Username:    foundUser.Username().Value(),
		Email:       foundUser.Email().Value(),
		DisplayName: foundUser.DisplayName(),
		Status:      foundUser.Status().String(),
	}, nil
}
