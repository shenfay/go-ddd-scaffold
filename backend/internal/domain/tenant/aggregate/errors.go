package aggregate

// 租户模块错误码 (30000-39999)
const (
	CodeTenantNotFound          = 30001 // 租户不存在
	CodeTenantExists            = 30002 // 租户已存在
	CodeTenantCodeExists        = 30003 // 租户编码已存在
	CodeNotTenantMember         = 31001 // 不是租户成员
	CodeNotTenantOwner          = 32003 // 不是租户所有者
	CodeTenantAlreadyActive     = 30010 // 租户已激活
	CodeTenantAlreadyInactive   = 30011 // 租户已禁用
	CodeTenantAlreadySuspended  = 30012 // 租户已暂停
	CodeTenantNotActive         = 30013 // 租户未激活
	CodeTenantMaxMembersReached = 30020 // 租户成员数已达上限
	CodeInvalidMaxMembers       = 30021 // 无效的最大成员数
	CodeAlreadyMember           = 30022 // 用户已是成员
	CodeOperatorNotMember       = 30023 // 操作者不是成员
	CodeInsufficientPermissions = 30024 // 权限不足
	CodeCannotRemoveOwner       = 30025 // 不能移除所有者
	CodeCannotRemoveAdmin       = 30026 // 不能移除管理员
	CodeCannotChangeOwnRole     = 30027 // 不能修改自己的角色
	CodeInvalidOperation        = 30028 // 无效操作
)
