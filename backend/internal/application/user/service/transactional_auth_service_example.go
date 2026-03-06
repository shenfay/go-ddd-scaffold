// Package service 认证应用服务实现 - 事务性事件发布示例
//
// ⚠️ 注意：这是一个示例文件，展示了如何使用事务性事件发件箱模式
// 实际使用时需要确保仓储层支持 WithTx 方法
//
// TODO: 在仓储层添加 WithTx 支持后，将此文件重命名为 transactional_auth_service.go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	user_assembler "go-ddd-scaffold/internal/application/user/assembler"
	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	user_event "go-ddd-scaffold/internal/domain/user/event"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/service"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/auth"
	"go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/infrastructure/transaction"
	errPkg "go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/validator"
	cast "go-ddd-scaffold/pkg/uitl"
)

// TransactionalAuthenticationService 支持事务性事件发布的认证服务
type TransactionalAuthenticationService struct {
	userRepo         repository.UserRepository
	tenantRepo       repository.TenantRepository
	tenantMemberRepo repository.TenantMemberRepository
	jwtService       entity.JWTService
	eventPublisher   *event.TransactionalEventPublisher
	userAssembler    user_assembler.UserAssembler
	userValidator    *validator.UserValidator
	tokenBlacklist   auth.TokenBlacklistService
	uow              transaction.UnitOfWork
	passwordHasher   service.PasswordHasher // 密码哈希器
}

// NewTransactionalAuthenticationService 创建支持事务性事件发布的认证服务
func NewTransactionalAuthenticationService(
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
	tenantMemberRepo repository.TenantMemberRepository,
	jwtService entity.JWTService,
	eventPublisher *event.TransactionalEventPublisher,
	tokenBlacklist auth.TokenBlacklistService,
	uow transaction.UnitOfWork,
	passwordHasher service.PasswordHasher, // 新增参数
) *TransactionalAuthenticationService {
	svc := &TransactionalAuthenticationService{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		tenantMemberRepo: tenantMemberRepo,
		jwtService:       jwtService,
		eventPublisher:   eventPublisher,
		userAssembler:    user_assembler.NewUserAssembler(),
		tokenBlacklist:   tokenBlacklist,
		uow:              uow,
		passwordHasher:   passwordHasher,
	}
	svc.userValidator = validator.NewUserValidator(nil)
	return svc
}

// Register 用户注册（使用事务性事件发布）
func (s *TransactionalAuthenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	// 1. 开启事务
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return nil, errPkg.Wrap(err, "TRANSACTION_BEGIN_FAILED", "开启事务失败")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	db := tx.GetDB()

	// 2. 业务校验：密码强度
	if err := validator.ValidatePasswordStrength(req.Password); err != nil {
		return nil, errPkg.ErrInvalidPassword
	}

	// 3. 检查邮箱是否已存在
	// TODO: 需要在 UserRepository 中实现 GetByEmailWithTx 方法
	// existingUser, _ := s.userRepo.GetByEmailWithTx(db, ctx, req.Email)
	_ = db // 避免编译错误，实际使用时删除
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errPkg.ErrUserExists
	}

	// 4. 将 string 的 TenantID 转换为 uuid.UUID
	var tenantID *uuid.UUID
	if req.TenantID != nil {
		var err error
		tenantID, err = cast.ToUUIDPtr(*req.TenantID)
		if err != nil {
			return nil, err
		}
	}

	// 5. 验证租户限制（如果是成员注册）
	if req.Role != nil && *req.Role == "member" {
		if tenantID == nil {
			return nil, errPkg.ErrUnauthorized // 成员必须指定租户
		}

		// 检查租户总用户数限制
		// TODO: 需要在 TenantMemberRepository 中实现 ListByTenantWithTx 方法
		count, err := s.tenantMemberRepo.ListByTenant(ctx, *tenantID)
		if err != nil {
			return nil, err
		}

		tenant, err := s.tenantRepo.GetByID(ctx, *tenantID)
		if err != nil {
			return nil, err
		}

		// 检查租户最大用户数限制
		if len(count) >= tenant.MaxMembers {
			return nil, errPkg.ErrTenantLimitExceed
		}
	}

	// 6. 验证密码强度并使用 PasswordHasher 加密
	plainPassword, err := valueobject.NewPlainPassword(req.Password)
	if err != nil {
		return nil, errPkg.ErrInvalidPassword
	}

	// ✅ 使用注入的 PasswordHasher
	hashedPasswordStr, err := s.passwordHasher.Hash(plainPassword.String())
	if err != nil {
		return nil, errPkg.Wrap(err, "HASH_PASSWORD_FAILED", "密码加密失败")
	}

	hashedPassword := entity.HashedPassword(hashedPasswordStr)

	// 7. 使用 Assembler 创建用户实体
	newUser, err := s.userAssembler.FromRegisterRequest(req, hashedPassword)
	if err != nil {
		return nil, errPkg.Wrap(err, "CREATE_USER_FAILED", "failed to create user from register request")
	}

	// 8. 在事务中保存用户
	// TODO: 需要在 UserRepository 中实现 CreateWithTx 方法
	// if err := s.userRepo.CreateWithTx(db, ctx, newUser); err != nil { ... }
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, errPkg.Wrap(err, "SAVE_USER_FAILED", "保存用户失败")
	}

	// 9. 如果指定了租户，则创建租户成员关系
	if tenantID != nil {
		tenantMember := &entity.TenantMember{
			ID:       uuid.New(),
			TenantID: *tenantID,
			UserID:   newUser.ID,
			Role:     "member", // 默认角色
			Status:   entity.MemberStatusActive,
			JoinedAt: time.Now(),
		}

		// TODO: 需要在 TenantMemberRepository 中实现 CreateWithTx 方法
		// if err := s.tenantMemberRepo.CreateWithTx(db, ctx, tenantMember); err != nil { ... }
		if err := s.tenantMemberRepo.Create(ctx, tenantMember); err != nil {
			return nil, errPkg.Wrap(err, "SAVE_TENANT_MEMBER_FAILED", "创建租户成员关系失败")
		}
	}

	// 10. ⭐ 关键：在事务中发布领域事件到发件箱
	// 事件会与用户数据在同一事务中原子性地提交
	registeredEvent := user_event.NewUserRegisteredEvent(
		newUser.ID,
		newUser.Email.String(),
	)

	if err := s.eventPublisher.PublishWithinTransaction(db, ctx, registeredEvent); err != nil {
		return nil, errPkg.Wrap(err, "SAVE_EVENT_FAILED", "保存领域事件失败")
	}

	// 11. 提交事务
	// 此时用户数据和事件都已持久化到数据库
	// 后台 Worker 会异步轮询发件箱并发布事件
	if err := tx.Commit(); err != nil {
		return nil, errPkg.Wrap(err, "COMMIT_TRANSACTION_FAILED", "提交事务失败")
	}

	// 12. 转换为 DTO 返回
	userDTO := s.userAssembler.ToDTO(newUser)
	if tenantID != nil {
		member, err := s.tenantMemberRepo.GetByUserAndTenant(ctx, newUser.ID, *tenantID)
		if err == nil {
			userDTO.Role = string(member.Role)
			userDTO.TenantID = cast.ToStringPtr(member.TenantID.String())
		}
	}
	return userDTO, nil
}

// Login 用户登录（简化版本，不使用事务性事件）
func (s *TransactionalAuthenticationService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 登录操作不涉及复杂的事务，可以直接发布事件

	// 1. 查找用户
	userEntity, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errPkg.ErrUserNotFound
	}

	// 2. 验证密码
	// TODO: 需要使用 PasswordHasher 验证
	if string(userEntity.Password) == req.Password {
		return nil, errPkg.ErrInvalidPassword
	}

	// 3. 检查用户状态
	if userEntity.Status != entity.StatusActive {
		return nil, errPkg.ErrUnauthorized
	}

	// 4. 生成 JWT 令牌
	token, err := s.jwtService.GenerateToken(userEntity.ID)
	if err != nil {
		return nil, err
	}

	// 5. 发布登录事件（示例：暂时不实际调用）
	_ = user_event.NewUserLoggedInEvent(
		userEntity.ID,
		"", // IP
		"", // UserAgent
		"", // DeviceType
		"success",
		nil,
	)

	// TODO: 需要在 TransactionalEventPublisher 中暴露 PublishSync 方法
	// 目前暂时直接使用 eventBus
	// s.eventPublisher.PublishSync(ctx, loginEvent)

	// 6. 返回登录响应
	return s.userAssembler.ToLoginResponseDTO(userEntity, token), nil
}

// Logout 用户登出
func (s *TransactionalAuthenticationService) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	// 如果配置了 Token 黑名单服务，将 token 加入黑名单
	if s.tokenBlacklist != nil && token != "" {
		_, err := s.jwtService.ValidateToken(token)
		if err == nil {
			expireAt := time.Now().Add(24 * time.Hour) // 默认 24 小时
			err = s.tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
			if err != nil {
				// 记录错误但不影响登出流程
			}
		}
	}

	return nil
}
