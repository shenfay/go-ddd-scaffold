package entity

import (
	"errors"
)

// Tenant 相关错误定义
var (
	ErrTenantInvalid              = errors.New("tenant is invalid")
	ErrTenantMemberLimitExceeded  = errors.New("tenant member limit exceeded")
	ErrTenantMemberAlreadyExists  = errors.New("tenant member already exists")
	ErrTenantMemberNotFound       = errors.New("tenant member not found")
	ErrTenantNotFound             = errors.New("tenant not found")
)
