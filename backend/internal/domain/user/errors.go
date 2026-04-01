package user

// 用户模块错误码 (20000-29999)
const (
	CodeUserNotFound    = 20001 // 用户不存在
	CodeUserExists      = 20002 // 用户已存在
	CodeInvalidPassword = 21001 // 密码错误
	CodeAccountLocked   = 21004 // 账户已锁定
	CodeInvalidEmail    = 22002 // 邮箱格式无效
	CodeEmailExists     = 22004 // 邮箱已存在

	// 用户状态相关
	CodeUserNotPending      = 20010 // 用户不在待激活状态
	CodeUserAlreadyInactive = 20011 // 用户已禁用
	CodeUserAlreadyLocked   = 20012 // 用户已锁定
	CodeUserNotLocked       = 20013 // 用户未锁定
	CodeUserCannotLogin     = 21003 // 用户无法登录
	CodeInvalidOldPassword  = 21005 // 旧密码错误
	CodeUsernameExists      = 22003 // 用户名已存在
)
