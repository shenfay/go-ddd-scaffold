package valueobject

import (
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// TenantID 租户标识
type TenantID struct {
	kernel.Int64Identity
}

// NewTenantID 创建租户标识
func NewTenantID(value int64) TenantID {
	return TenantID{Int64Identity: kernel.NewInt64Identity(value)}
}

// String 返回租户标识字符串
func (tid TenantID) String() string {
	return fmt.Sprintf("tenant-%d", tid.Int64())
}
