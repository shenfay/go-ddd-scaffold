package validator

import (
	"context"
)

// UserValidator User模块业务校验器
type UserValidator struct {
	userRepo UserRepository
}

// NewUserValidator 创建User校验器
func NewUserValidator(userRepo UserRepository) *UserValidator {
	return &UserValidator{userRepo: userRepo}
}

// ValidateRegisterRequest 校验注册请求
func (v *UserValidator) ValidateRegisterRequest(ctx context.Context, email, nickname, role string) ValidationErrors {
	var errs ValidationErrors

	// 邮箱唯一性校验
	if v.userRepo != nil && email != "" {
		exists, err := v.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			// TODO: 记录日志
		} else if exists {
			errs = append(errs, ValidationError{
				Field:   "email",
				Message: "邮箱已被注册",
			})
		}
	}

	// 昵称唯一性校验
	if nickname != "" && v.userRepo != nil {
		exists, err := v.userRepo.ExistsByNickname(ctx, nickname)
		if err != nil {
			// TODO: 记录日志
		} else if exists {
			errs = append(errs, ValidationError{
				Field:   "nickname",
				Message: "昵称已被使用",
			})
		}
	}

	// 角色校验
	validRoles := map[string]bool{
		"PARENT":  true,
		"CHILD":   true,
		"TEACHER": true,
		"ADMIN":   true,
	}
	if !validRoles[role] {
		errs = append(errs, ValidationError{
			Field:   "role",
			Message: "无效的用户角色",
		})
	}

	return errs
}
