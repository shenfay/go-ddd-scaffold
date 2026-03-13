package user

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=100"`
	FirstName   *string `json:"first_name,omitempty" binding:"omitempty,max=50"`
	LastName    *string `json:"last_name,omitempty" binding:"omitempty,max=50"`
	Gender      *string `json:"gender,omitempty" binding:"omitempty,oneof=male female other unknown"`
	PhoneNumber *string `json:"phone_number,omitempty" binding:"omitempty,max=20"`
	AvatarURL   *string `json:"avatar_url,omitempty" binding:"omitempty,max=500,url"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// DeactivateUserRequest 禁用用户请求
type DeactivateUserRequest struct {
	Reason string `json:"reason,omitempty" binding:"omitempty,max=500"`
}

// ListUsersRequest 列出用户请求
type ListUsersRequest struct {
	Keyword  string `form:"keyword"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive pending locked"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// GetUserRequest 获取用户请求（路径参数）
type GetUserRequest struct {
	UserID string `uri:"id" binding:"required"`
}

// ActivateUserRequest 激活用户请求（路径参数）
type ActivateUserRequest struct {
	UserID string `uri:"id" binding:"required"`
}
