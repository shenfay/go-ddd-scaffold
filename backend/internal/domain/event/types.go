package event

// ActivityType 活动类型
type ActivityType string

const (
	// 用户相关
	ActivityUserRegistered      ActivityType = "USER_REGISTERED"
	ActivityUserLoggedIn        ActivityType = "USER_LOGIN"
	ActivityUserLoggedOut       ActivityType = "USER_LOGOUT"
	ActivityUserActivated       ActivityType = "USER_ACTIVATED"
	ActivityUserDeactivated     ActivityType = "USER_DEACTIVATED"
	ActivityUserLocked          ActivityType = "USER_LOCKED"
	ActivityUserUnlocked        ActivityType = "USER_UNLOCKED"
	ActivityUserPasswordChanged ActivityType = "USER_PASSWORD_CHANGED"
	ActivityUserEmailChanged    ActivityType = "USER_EMAIL_CHANGED"
	ActivityUserProfileUpdated  ActivityType = "USER_PROFILE_UPDATED"

	// 订单相关（示例）
	ActivityOrderCreated   ActivityType = "ORDER_CREATED"
	ActivityOrderPaid      ActivityType = "ORDER_PAID"
	ActivityOrderShipped   ActivityType = "ORDER_SHIPPED"
	ActivityOrderCancelled ActivityType = "ORDER_CANCELLED"

	// 系统相关
	ActivitySystemError   ActivityType = "SYSTEM_ERROR"
	ActivitySecurityAlert ActivityType = "SECURITY_ALERT"
)

// ActivityStatus 活动状态
type ActivityStatus int16

const (
	// ActivityStatusSuccess 成功
	ActivityStatusSuccess ActivityStatus = 0
	// ActivityStatusFailed 失败
	ActivityStatusFailed ActivityStatus = 1
)
