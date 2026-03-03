package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 参数校验器封装
type Validator struct {
	validate *validator.Validate
}

// New 创建校验器实例
func New() *Validator {
	v := validator.New()
	// 使用 tag name 作为字段名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validate: v}
}

// Validate 校验结构体
func (v *Validator) ValidateStruct(s interface{}) ValidationErrors {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs ValidationErrors
	// 尝试转换为 validator.ValidationErrors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrs {
			errs = append(errs, ValidationError{
				Field:   ve.Field(),
				Tag:     ve.Tag(),
				Value:   ve.Value(),
				Message: formatError(ve),
			})
		}
	}
	return errs
}

// ValidationError 单个字段校验错误
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// ValidationErrors 校验错误集合
type ValidationErrors []ValidationError

// Error 实现 error 接口
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Message)
	}
	return strings.Join(msgs, "; ")
}

// ToMap 转换为 map 用于响应
func (e ValidationErrors) ToMap() []ValidationError {
	return e
}

func formatError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s 是必填字段", err.Field())
	case "email":
		return fmt.Sprintf("%s 不是有效的邮箱地址", err.Field())
	case "min":
		return fmt.Sprintf("%s 长度不能小于 %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s 长度不能大于 %s", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", err.Field(), err.Param())
	case "eqfield":
		return fmt.Sprintf("%s 与 %s 不匹配", err.Field(), err.Param())
	case "nefield":
		return fmt.Sprintf("%s 不能与 %s 相同", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s 必须大于等于 %s", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于等于 %s", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s 必须是其中的值: %s", err.Field(), err.Param())
	case "uuid":
		return fmt.Sprintf("%s 不是有效的 UUID 格式", err.Field())
	case "uuid3":
		return fmt.Sprintf("%s 不是有效的 UUIDv3 格式", err.Field())
	case "uuid4":
		return fmt.Sprintf("%s 不是有效的 UUIDv4 格式", err.Field())
	case "uuid5":
		return fmt.Sprintf("%s 不是有效的 UUIDv5 格式", err.Field())
	case "url":
		return fmt.Sprintf("%s 不是有效的 URL", err.Field())
	case "uri":
		return fmt.Sprintf("%s 不是有效的 URI", err.Field())
	case "alpha":
		return fmt.Sprintf("%s 只能包含字母", err.Field())
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", err.Field())
	case "numeric":
		return fmt.Sprintf("%s 必须是数字", err.Field())
	case "hexadecimal":
		return fmt.Sprintf("%s 必须是十六进制", err.Field())
	case "hexcolor":
		return fmt.Sprintf("%s 不是有效的十六进制颜色", err.Field())
	case "rgb":
		return fmt.Sprintf("%s 不是有效的 RGB 颜色", err.Field())
	case "rgba":
		return fmt.Sprintf("%s 不是有效的 RGBA 颜色", err.Field())
	case "hsl":
		return fmt.Sprintf("%s 不是有效的 HSL 颜色", err.Field())
	case "hsla":
		return fmt.Sprintf("%s 不是有效的 HSLA 颜色", err.Field())
	default:
		return fmt.Sprintf("%s 校验失败: %s", err.Field(), err.Tag())
	}
}

// 全局校验器实例
var defaultValidator = New()

// Validate 全局校验函数
func Validate(s interface{}) ValidationErrors {
	return defaultValidator.ValidateStruct(s)
}
