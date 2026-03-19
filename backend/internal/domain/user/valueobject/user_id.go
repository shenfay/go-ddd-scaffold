package valueobject

import (
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// ============================================================================
// Identity - 身份标识
// ============================================================================

// UserID 用户标识
type UserID struct {
	kernel.Int64Identity
}

// NewUserID 创建用户标识
func NewUserID(value int64) UserID {
	return UserID{Int64Identity: kernel.NewInt64Identity(value)}
}

// String 返回用户标识字符串
func (uid UserID) String() string {
	return fmt.Sprintf("%d", uid.Int64())
}
