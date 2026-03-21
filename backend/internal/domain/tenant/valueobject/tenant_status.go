package valueobject

// TenantStatus 租户状态枚举
type TenantStatus int

const (
	TenantStatusActive    TenantStatus = iota // 活跃
	TenantStatusInactive                      // 停用
	TenantStatusSuspended                     // 暂停
)

// String 返回状态字符串表示
func (ts TenantStatus) String() string {
	switch ts {
	case TenantStatusActive:
		return "active"
	case TenantStatusInactive:
		return "inactive"
	case TenantStatusSuspended:
		return "suspended"
	default:
		return "unknown"
	}
}
