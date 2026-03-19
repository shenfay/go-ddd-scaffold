package valueobject

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusPending  UserStatus = iota // 待激活
	UserStatusActive                     // 激活
	UserStatusInactive                   // 禁用
	UserStatusLocked                     // 锁定
)

// String 返回状态字符串表示
func (us UserStatus) String() string {
	switch us {
	case UserStatusPending:
		return "pending"
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}
