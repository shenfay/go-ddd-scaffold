package valueobject

import (
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// TenantMember 租户成员值对象
type TenantMember struct {
	UserID   vo.UserID
	TenantID TenantID
	Role     TenantRole
	JoinedAt string
}
