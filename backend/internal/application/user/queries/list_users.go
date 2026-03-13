package queries

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserListItemDTO 用户列表项 DTO
type UserListItemDTO struct {
	UserID      string  `json:"user_id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName string  `json:"display_name"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	LastLoginAt *string `json:"last_login_at,omitempty"`
}

// ListUsersQuery 列出用户查询
type ListUsersQuery struct {
	Keyword  string
	Status   *user.UserStatus
	Page     int
	PageSize int
}

// UserListResult 用户列表结果
type UserListResult struct {
	Items      []*UserListItemDTO `json:"items"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// ListUsersHandler 列出用户查询处理器
type ListUsersHandler struct {
	userRepo user.UserRepository
}

// NewListUsersHandler 创建查询处理器
func NewListUsersHandler(userRepo user.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{
		userRepo: userRepo,
	}
}

// Handle 处理列出用户查询
func (h *ListUsersHandler) Handle(ctx context.Context, query *ListUsersQuery) (*UserListResult, error) {
	// 检查仓储是否初始化
	if h.userRepo == nil {
		// 返回空结果（Mock 模式）
		return &UserListResult{
			Items:      []*UserListItemDTO{},
			TotalCount: 0,
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalPages: 0,
		}, nil
	}

	criteria := user.UserSearchCriteria{
		Keyword: query.Keyword,
		Status:  query.Status,
	}

	pagination := ddd.NewPagination(query.Page, query.PageSize)

	result, err := h.userRepo.FindByCriteria(ctx, criteria, pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*UserListItemDTO, len(result.Items))
	for i, u := range result.Items {
		items[i] = toUserListItemDTO(u)
	}

	return &UserListResult{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// toUserListItemDTO 将领域用户转换为列表项 DTO
func toUserListItemDTO(u *user.User) *UserListItemDTO {
	// 简化实现：只访问公开方法，避免访问未初始化的私有字段
	dto := &UserListItemDTO{
		UserID:      u.ID().(user.UserID).String(),
		Username:    "",
		Email:       "",
		DisplayName: u.DisplayName(),
		Status:      u.Status().String(),
		CreatedAt:   u.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	// 安全访问可能为 nil 的字段
	if username := u.Username(); username != nil {
		dto.Username = username.Value()
	}
	if email := u.Email(); email != nil {
		dto.Email = email.Value()
	}
	if lastLogin := u.LastLoginAt(); lastLogin != nil {
		formatted := lastLogin.Format("2006-01-02T15:04:05Z07:00")
		dto.LastLoginAt = &formatted
	}

	return dto
}
