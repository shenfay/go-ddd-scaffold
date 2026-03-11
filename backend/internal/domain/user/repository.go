package user

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	ddd.Repository

	// 基础仓储操作
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	Delete(ctx context.Context, id UserID) error
	Exists(ctx context.Context, id UserID) (bool, error)

	// 查询操作
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByStatus(ctx context.Context, status UserStatus) ([]*User, error)

	// 分页查询
	FindAll(ctx context.Context, pagination ddd.Pagination) (*ddd.PaginatedResult[*User], error)
	FindByCriteria(ctx context.Context, criteria UserSearchCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*User], error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status UserStatus) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, users []*User) error
	DeleteBatch(ctx context.Context, ids []UserID) error

	// 乐观锁支持
	SaveWithVersion(ctx context.Context, user *User, expectedVersion int) error
}

// UserSearchCriteria 用户搜索条件
type UserSearchCriteria struct {
	Keyword       string      `json:"keyword,omitempty"`
	Status        *UserStatus `json:"status,omitempty"`
	Gender        *UserGender `json:"gender,omitempty"`
	CreatedFrom   *string     `json:"created_from,omitempty"`
	CreatedTo     *string     `json:"created_to,omitempty"`
	LastLoginFrom *string     `json:"last_login_from,omitempty"`
	LastLoginTo   *string     `json:"last_login_to,omitempty"`
}

// UserReadModel 用户读模型接口
type UserReadModel interface {
	GetUserProfile(ctx context.Context, userID UserID) (*UserProfileDTO, error)
	ListUsers(ctx context.Context, criteria UserListCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*UserListItemDTO], error)
	SearchUsers(ctx context.Context, keyword string, pagination ddd.Pagination) (*ddd.PaginatedResult[*UserListItemDTO], error)
	GetUserStatistics(ctx context.Context) (*UserStatisticsDTO, error)
}

// UserProfileDTO 用户资料DTO
type UserProfileDTO struct {
	UserID      UserID     `json:"user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Gender      UserGender `json:"gender"`
	PhoneNumber string     `json:"phone_number"`
	AvatarURL   string     `json:"avatar_url"`
	Status      UserStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	LastLoginAt *string    `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count"`
}

// UserListItemDTO 用户列表项DTO
type UserListItemDTO struct {
	UserID      UserID     `json:"user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	Status      UserStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
	LastLoginAt *string    `json:"last_login_at,omitempty"`
}

// UserListCriteria 用户列表查询条件
type UserListCriteria struct {
	Status   *UserStatus `json:"status,omitempty"`
	Keyword  string      `json:"keyword,omitempty"`
	SortBy   string      `json:"sort_by,omitempty"`
	SortDesc bool        `json:"sort_desc,omitempty"`
}

// UserStatisticsDTO 用户统计DTO
type UserStatisticsDTO struct {
	TotalUsers             int64 `json:"total_users"`
	ActiveUsers            int64 `json:"active_users"`
	InactiveUsers          int64 `json:"inactive_users"`
	PendingUsers           int64 `json:"pending_users"`
	LockedUsers            int64 `json:"locked_users"`
	TodayLogins            int64 `json:"today_logins"`
	ThisWeekLogins         int64 `json:"this_week_logins"`
	ThisMonthLogins        int64 `json:"this_month_logins"`
	AverageSessionDuration int64 `json:"average_session_duration"` // 秒
}
