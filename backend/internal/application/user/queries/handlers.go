package queries

import (
	"context"

	domainUser "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	UserID domainUser.UserID
}

// GetUserQueryHandler 获取用户查询处理器
type GetUserQueryHandler struct {
	userRepo domainUser.UserRepository
}

// NewGetUserQueryHandler 创建获取用户查询处理器
func NewGetUserQueryHandler(userRepo domainUser.UserRepository) *GetUserQueryHandler {
	return &GetUserQueryHandler{
		userRepo: userRepo,
	}
}

// Handle 处理获取用户查询
func (h *GetUserQueryHandler) Handle(ctx context.Context, query *GetUserQuery) (*domainUser.User, error) {
	u, err := h.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}
	return u, nil
}

// ListUsersQuery 列出用户查询
type ListUsersQuery struct {
	Page       int
	PageSize   int
	Status     *domainUser.UserStatus
	Keyword    string
	SortBy     string
	SortDesc   bool
}

// ListUsersQueryHandler 列出用户查询处理器
type ListUsersQueryHandler struct {
	userRepo domainUser.UserRepository
}

// NewListUsersQueryHandler 创建列出用户查询处理器
func NewListUsersQueryHandler(userRepo domainUser.UserRepository) *ListUsersQueryHandler {
	return &ListUsersQueryHandler{
		userRepo: userRepo,
	}
}

// Handle 处理列出用户查询
func (h *ListUsersQueryHandler) Handle(ctx context.Context, query *ListUsersQuery) (*ddd.PaginatedResult[*domainUser.User], error) {
	pagination := ddd.NewPagination(query.Page, query.PageSize)

	criteria := domainUser.UserSearchCriteria{
		Keyword: query.Keyword,
		Status:  query.Status,
	}

	result, err := h.userRepo.FindByCriteria(ctx, criteria, pagination)
	if err != nil {
		return nil, err
	}

	return result, nil
}
