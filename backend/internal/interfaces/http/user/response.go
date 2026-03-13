package user

// ==================== 用户模块 Swagger Response 类型 ====================

// CreateUserResponse 创建用户响应
// @Description 创建用户操作返回的数据结构
type CreateUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// UserResponse 用户详情响应
// @Description 用户详细信息数据结构
type UserResponse struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Status      int32   `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// UserItemResponse 用户列表项响应
// @Description 用户列表中的单个用户信息
type UserItemResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

// UserListResponse 用户列表响应
// @Description 分页用户列表
type UserListResponse struct {
	Items []UserItemResponse `json:"items"`
	Meta  PageMeta           `json:"meta"`
}

// PageMeta 分页元信息
// @Description 分页查询的元数据
type PageMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ChangePasswordResponse 修改密码响应
// @Description 修改密码操作返回结果
type ChangePasswordResponse struct {
	Message string `json:"message"`
}
