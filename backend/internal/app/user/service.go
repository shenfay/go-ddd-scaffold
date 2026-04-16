package user

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Email    string
	Password string
}

// UpdateProfileCommand 更新资料命令
type UpdateProfileCommand struct {
	UserID string
	Email  string
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	Locked        bool       `json:"locked"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Service 用户应用服务
type Service struct {
	userRepo user.UserRepository
	eventBus messaging.EventBus
}

// NewService 创建用户应用服务
func NewService(userRepo user.UserRepository, eventBus messaging.EventBus) *Service {
	return &Service{
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// CreateUser 创建用户
func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (*UserDTO, error) {
	// 1. 创建领域实体
	u, err := user.NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	// 2. 持久化
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 3. 发布领域事件
	if s.eventBus != nil {
		evt := &user.UserRegistered{
			UserID:    u.ID,
			Email:     u.Email,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			// 事件发布失败不影响主流程，显式忽略错误
			_ = err
		}
	}

	// 4. 返回DTO
	return toUserDTO(u), nil
}

// GetUserByID 根据ID获取用户
func (s *Service) GetUserByID(ctx context.Context, userID string) (*UserDTO, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toUserDTO(u), nil
}

// UpdateProfile 更新用户资料
func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (*UserDTO, error) {
	// 1. 加载聚合根
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 2. 更新资料
	if err := u.UpdateEmail(cmd.Email); err != nil {
		return nil, err
	}

	// 3. 持久化
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 4. 发布领域事件
	if s.eventBus != nil {
		evt := &UserProfileUpdated{
			UserID:    u.ID,
			Email:     u.Email,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			// 事件发布失败不影响主流程，显式忽略错误
			_ = err
		}
	}

	return toUserDTO(u), nil
}

// toUserDTO 实体转DTO
func toUserDTO(u *user.User) *UserDTO {
	return &UserDTO{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Locked:        u.Locked,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
