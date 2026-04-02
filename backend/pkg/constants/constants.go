package constants

// 系统常量
const (
	// 项目名称
	ProjectName = "Go DDD Scaffold"

	// API 版本
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion

	// 默认分页参数
	DefaultPageSize   = 20
	MaxPageSize       = 100
	DefaultPageNumber = 1

	// 时间格式
	TimeLayoutRFC3339 = "2006-01-02T15:04:05Z07:00"
	TimeLayoutDate    = "2006-01-02"
	TimeLayoutTime    = "15:04:05"
)

// Redis Key 前缀
const (
	RedisKeyPrefix            = "go_ddd_scaffold:"
	RedisKeyRefreshToken      = RedisKeyPrefix + "refresh_token:"
	RedisKeyUserSession       = RedisKeyPrefix + "user_session:"
	RedisKeyLoginAttempts     = RedisKeyPrefix + "login_attempts:"
	RedisKeyEmailVerification = RedisKeyPrefix + "email_verification:"
	RedisKeyPasswordReset     = RedisKeyPrefix + "password_reset:"
)

// Asynq 任务类型
const (
	AsynqTaskSendVerificationEmail = "auth:send_verification_email"
	AsynqTaskSendWelcomeEmail      = "auth:send_welcome_email"
	AsynqTaskLogUserRegistration   = "auth:log_user_registration"
	AsynqTaskLogLoginAttempt       = "auth:log_login_attempt"
	AsynqTaskCleanupExpiredTokens  = "auth:cleanup_expired_tokens"
)

// 队列名称
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)
