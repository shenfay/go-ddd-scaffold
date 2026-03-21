package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
	uservo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// TenantRepository 租户仓储接口
type TenantRepository interface {
	// 基础仓储操作
	Save(ctx context.Context, tenant *aggregate.Tenant) error
	FindByID(ctx context.Context, id valueobject.TenantID) (*aggregate.Tenant, error)
	Delete(ctx context.Context, id valueobject.TenantID) error
	Exists(ctx context.Context, id valueobject.TenantID) (bool, error)

	// 查询操作
	FindByCode(ctx context.Context, code string) (*aggregate.Tenant, error)
	FindByOwnerID(ctx context.Context, ownerID uservo.UserID) ([]*aggregate.Tenant, error)
	FindByStatus(ctx context.Context, status valueobject.TenantStatus) ([]*aggregate.Tenant, error)

	// 分页查询
	FindAll(ctx context.Context, pagination kernel.Pagination) (*kernel.PaginatedResult[*aggregate.Tenant], error)
	FindByCriteria(ctx context.Context, criteria TenantSearchCriteria, pagination kernel.Pagination) (*kernel.PaginatedResult[*aggregate.Tenant], error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status valueobject.TenantStatus) (int64, error)

	// 成员操作
	AddMember(ctx context.Context, tenantID valueobject.TenantID, member *valueobject.TenantMember) error
	RemoveMember(ctx context.Context, tenantID valueobject.TenantID, userID uservo.UserID) error
	FindMembers(ctx context.Context, tenantID valueobject.TenantID) ([]*valueobject.TenantMember, error)
	FindMemberByUserID(ctx context.Context, tenantID valueobject.TenantID, userID uservo.UserID) (*valueobject.TenantMember, error)

	// 乐观锁支持
	SaveWithVersion(ctx context.Context, tenant *aggregate.Tenant, expectedVersion int) error
}

// TenantSearchCriteria 租户搜索条件
type TenantSearchCriteria struct {
	Keyword    string                    `json:"keyword,omitempty"`
	Status     *valueobject.TenantStatus `json:"status,omitempty"`
	OwnerID    *uservo.UserID            `json:"owner_id,omitempty"`
	CodePrefix string                    `json:"code_prefix,omitempty"`
}

// TenantReadModel 租户读模型接口
type TenantReadModel interface {
	GetTenantProfile(ctx context.Context, tenantID valueobject.TenantID) (*TenantProfileDTO, error)
	ListTenants(ctx context.Context, criteria TenantListCriteria, pagination kernel.Pagination) (*kernel.PaginatedResult[*TenantListItemDTO], error)
	SearchTenants(ctx context.Context, keyword string, pagination kernel.Pagination) (*kernel.PaginatedResult[*TenantListItemDTO], error)
	GetTenantStatistics(ctx context.Context) (*TenantStatisticsDTO, error)
	GetTenantMembers(ctx context.Context, tenantID valueobject.TenantID) ([]*TenantMemberDTO, error)
}

// TenantProfileDTO 租户资料DTO
type TenantProfileDTO struct {
	TenantID    valueobject.TenantID      `json:"tenant_id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Status      valueobject.TenantStatus  `json:"status"`
	OwnerID     uservo.UserID             `json:"owner_id"`
	MaxMembers  int                       `json:"max_members"`
	Config      *valueobject.TenantConfig `json:"config"`
	CreatedAt   string                    `json:"created_at"`
	UpdatedAt   string                    `json:"updated_at"`
	MemberCount int                       `json:"member_count"`
}

// TenantListItemDTO 租户列表项DTO
type TenantListItemDTO struct {
	TenantID    valueobject.TenantID     `json:"tenant_id"`
	Code        string                   `json:"code"`
	Name        string                   `json:"name"`
	Status      valueobject.TenantStatus `json:"status"`
	OwnerID     uservo.UserID            `json:"owner_id"`
	MemberCount int                      `json:"member_count"`
	CreatedAt   string                   `json:"created_at"`
}

// TenantListCriteria 租户列表查询条件
type TenantListCriteria struct {
	Status   *valueobject.TenantStatus `json:"status,omitempty"`
	Keyword  string                    `json:"keyword,omitempty"`
	OwnerID  *uservo.UserID            `json:"owner_id,omitempty"`
	SortBy   string                    `json:"sort_by,omitempty"`
	SortDesc bool                      `json:"sort_desc,omitempty"`
}

// TenantStatisticsDTO 租户统计DTO
type TenantStatisticsDTO struct {
	TotalTenants        int64   `json:"total_tenants"`
	ActiveTenants       int64   `json:"active_tenants"`
	InactiveTenants     int64   `json:"inactive_tenants"`
	SuspendedTenants    int64   `json:"suspended_tenants"`
	TotalMembers        int64   `json:"total_members"`
	AvgMembersPerTenant float64 `json:"avg_members_per_tenant"`
}

// TenantMemberDTO 租户成员DTO
type TenantMemberDTO struct {
	UserID   uservo.UserID        `json:"user_id"`
	TenantID valueobject.TenantID `json:"tenant_id"`
	Username string               `json:"username"`
	Email    string               `json:"email"`
	Role     string               `json:"role"`
	JoinedAt string               `json:"joined_at"`
}
