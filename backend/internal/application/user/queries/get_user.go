package queries

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	UserID user.UserID
}

// UserDetailDTO 用户详情DTO
type UserDetailDTO struct {
	UserID      string     `json:"user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Gender      string     `json:"gender"`
	PhoneNumber string     `json:"phone_number"`
	AvatarURL   string     `json:"avatar_url"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GetUserHandler 获取用户查询处理器
type GetUserHandler struct {
	userRepo user.UserRepository
}

// NewGetUserHandler 创建查询处理器
func NewGetUserHandler(userRepo user.UserRepository) *GetUserHandler {
	return &GetUserHandler{
		userRepo: userRepo,
	}
}

// Handle 处理获取用户查询
func (h *GetUserHandler) Handle(ctx context.Context, query *GetUserQuery) (*UserDetailDTO, error) {
	// 检查仓储是否初始化
	if h.userRepo == nil {
		return nil, ddd.ErrAggregateNotFound
	}

	u, err := h.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	return toUserDetailDTO(u), nil
}

// toUserDetailDTO 将领域用户转换为详情DTO
func toUserDetailDTO(u *user.User) *UserDetailDTO {
	return &UserDetailDTO{
		UserID:      u.ID().(user.UserID).String(),
		Username:    u.Username().Value(),
		Email:       u.Email().Value(),
		DisplayName: u.DisplayName(),
		FirstName:   u.FirstName(),
		LastName:    u.LastName(),
		Gender:      u.Gender().String(),
		PhoneNumber: u.PhoneNumber(),
		AvatarURL:   u.AvatarURL(),
		Status:      u.Status().String(),
		LastLoginAt: u.LastLoginAt(),
		LoginCount:  u.LoginCount(),
		CreatedAt:   u.CreatedAt(),
		UpdatedAt:   u.UpdatedAt(),
	}
}
