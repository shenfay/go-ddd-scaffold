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

// EventName 领域事件名称类型
type EventName string

// QueueName 消息队列名称类型
type QueueName string

// AsynqTaskType Asynq 任务类型
type AsynqTaskType string

// 领域事件名称常量
const (
	EventUserRegistered     EventName = "user.registered"
	EventUserLoggedIn       EventName = "user.logged_in"
	EventUserLoginFailed    EventName = "user.login_failed"
	EventUserAccountLocked  EventName = "user.account_locked"
	EventUserLoggedOut      EventName = "user.logged_out"
	EventUserTokenRefreshed EventName = "user.token_refreshed"
	EventUserProfileUpdated EventName = "user.profile_updated"
)

// Asynq 任务类型
const (
	AsynqTaskSendVerificationEmail  = "auth:send_verification_email"
	AsynqTaskSendPasswordResetEmail = "auth:send_password_reset_email"
	AsynqTaskSendWelcomeEmail       = "auth:send_welcome_email"
	AsynqTaskLogUserRegistration    = "auth:log_user_registration"
	AsynqTaskLogLoginAttempt        = "auth:log_login_attempt"
	AsynqTaskCleanupExpiredTokens   = "auth:cleanup_expired_tokens"
)

// 队列名称
const (
	QueueCritical QueueName = "critical"
	QueueDefault  QueueName = "default"
	QueueLow      QueueName = "low"
)
