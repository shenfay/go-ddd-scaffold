package valueobject

// TenantRole 租户角色枚举
type TenantRole int

const (
	TenantRoleOwner  TenantRole = iota // 所有者
	TenantRoleAdmin                    // 管理员
	TenantRoleMember                   // 成员
	TenantRoleGuest                    // 访客
)

// String 返回角色字符串表示
func (tr TenantRole) String() string {
	switch tr {
	case TenantRoleOwner:
		return "owner"
	case TenantRoleAdmin:
		return "admin"
	case TenantRoleMember:
		return "member"
	case TenantRoleGuest:
		return "guest"
	default:
		return "unknown"
	}
}
