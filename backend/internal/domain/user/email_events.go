package user

import "time"

// EmailVerified 邮箱验证成功事件
type EmailVerified struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewEmailVerifiedEvent 创建邮箱验证成功事件
func NewEmailVerifiedEvent(userID, email string) *EmailVerified {
	return &EmailVerified{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
}

func (e *EmailVerified) GetPayload() interface{} {
	return e
}

func (e *EmailVerified) GetType() string {
	return "USER.EMAIL.VERIFIED"
}

// EmailVerificationRequested 邮箱验证请求事件
type EmailVerificationRequested struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Timestamp time.Time `json:"timestamp"`
}

// NewEmailVerificationRequestedEvent 创建邮箱验证请求事件
func NewEmailVerificationRequestedEvent(userID, email, token string) *EmailVerificationRequested {
	return &EmailVerificationRequested{
		UserID:    userID,
		Email:     email,
		Token:     token,
		Timestamp: time.Now(),
	}
}

func (e *EmailVerificationRequested) GetPayload() interface{} {
	return e
}

func (e *EmailVerificationRequested) GetType() string {
	return "USER.EMAIL.VERIFICATION.REQUESTED"
}
