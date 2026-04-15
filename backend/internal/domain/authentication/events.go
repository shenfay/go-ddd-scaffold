package authentication

import "time"

// PasswordResetRequested 密码重置请求事件
//
// 当用户请求密码重置时发布此事件,触发邮件发送等后续操作
type PasswordResetRequested struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// GetPayload 获取事件数据
func (e *PasswordResetRequested) GetPayload() interface{} {
	return e
}

// GetType 获取事件类型
func (e *PasswordResetRequested) GetType() string {
	return "AUTH.PASSWORD_RESET.REQUESTED"
}

// PasswordResetCompleted 密码重置完成事件
//
// 当用户成功重置密码后发布此事件,用于审计日志和安全监控
type PasswordResetCompleted struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// GetPayload 获取事件数据
func (e *PasswordResetCompleted) GetPayload() interface{} {
	return e
}

// GetType 获取事件类型
func (e *PasswordResetCompleted) GetType() string {
	return "AUTH.PASSWORD_RESET.COMPLETED"
}
