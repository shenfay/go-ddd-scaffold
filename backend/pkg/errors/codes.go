package errors

// ==========================================
// 错误码定义（按领域分组）
// 格式：{DOMAIN}.{CATEGORY}.{ERROR}
// ==========================================

// 系统级错误码 (SYSTEM.*)
const (
	ErrCodeSystemInternal        = "SYSTEM.INTERNAL_ERROR"
	ErrCodeSystemInvalidRequest  = "SYSTEM.INVALID_REQUEST"
	ErrCodeSystemUnauthorized    = "SYSTEM.UNAUTHORIZED"
	ErrCodeSystemForbidden       = "SYSTEM.FORBIDDEN"
	ErrCodeSystemNotFound        = "SYSTEM.NOT_FOUND"
	ErrCodeSystemConflict        = "SYSTEM.CONFLICT"
	ErrCodeSystemTooManyRequests = "SYSTEM.TOO_MANY_REQUESTS"
)

// 认证域错误码 (AUTH.*)
const (
	ErrCodeAuthInvalidCredentials = "AUTH.INVALID_CREDENTIALS"
	ErrCodeAuthInvalidToken       = "AUTH.INVALID_TOKEN"
	ErrCodeAuthTokenExpired       = "AUTH.TOKEN_EXPIRED"
	ErrCodeAuthTokenRevoked       = "AUTH.TOKEN_REVOKED"
	ErrCodeAuthAccountLocked      = "AUTH.ACCOUNT_LOCKED"
	ErrCodeAuthInvalidResetToken  = "AUTH.INVALID_RESET_TOKEN" // 无效的重置令牌
	ErrCodeAuthResetTokenExpired  = "AUTH.RESET_TOKEN_EXPIRED" // 重置令牌已过期
)

// 用户域错误码 (USER.*)
const (
	ErrCodeUserNotFound           = "USER.NOT_FOUND"
	ErrCodeUserEmailAlreadyExists = "USER.EMAIL_ALREADY_EXISTS"
	ErrCodeUserEmailNotVerified   = "USER.EMAIL_NOT_VERIFIED"
)

// 安全域错误码 (SECURITY.*)
const (
	ErrCodeSecurityTooManyAttempts = "SECURITY.TOO_MANY_ATTEMPTS"
)
