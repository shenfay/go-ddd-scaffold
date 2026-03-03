package middleware

import (
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator 校验实例
var validate *validator.Validate

// InitValidator 初始化校验器
func InitValidator() {
	validate = validator.New()

	// 自定义时间格式校验
	registerCustomValidators()
}

// registerCustomValidators 注册自定义校验器
func registerCustomValidators() {
	// UUID v4 校验
	validate.RegisterValidation("uuid4", func(fl validator.FieldLevel) bool {
		uuid := fl.Field().String()
		return len(uuid) == 36 && uuid[8] == '-' && uuid[13] == '-' && uuid[18] == '-' && uuid[23] == '-'
	})

	// 节点类型校验 (C/S/T/P)
	validate.RegisterValidation("nodeType", func(fl validator.FieldLevel) bool {
		nodeType := fl.Field().String()
		return nodeType == "C" || nodeType == "S" || nodeType == "T" || nodeType == "P"
	})

	// 关系类型校验
	validate.RegisterValidation("relType", func(fl validator.FieldLevel) bool {
		relType := fl.Field().String()
		return relType == "PREREQ" || relType == "SUP_SKILL" || relType == "THINK_PAT"
	})

	// 能力等级校验 (1-5)
	validate.RegisterValidation("competencyLevel", func(fl validator.FieldLevel) bool {
		level := fl.Field().String()
		return len(level) > 0 && level[0] >= '1' && level[0] <= '5'
	})
}

// ValidationMiddleware 请求参数校验中间件（需要配合 ValidateJSON 使用）
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if validate == nil {
			InitValidator()
		}

		// 获取结构体名称
		var obj interface{}
		if s, exists := c.Get("validatorObject"); exists {
			obj = s
		}

		if obj == nil {
			c.Next()
			return
		}

		// 执行校验
		if err := validate.Struct(obj); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				field := err.Field()
				tag := err.Tag()

				switch tag {
				case "required":
					errors[field] = field + " 是必填字段"
				case "min":
					errors[field] = field + " 长度不能小于 " + err.Param()
				case "max":
					errors[field] = field + " 长度不能大于 " + err.Param()
				case "uuid4":
					errors[field] = field + " 必须是有效的 UUID v4 格式"
				case "nodeType":
					errors[field] = field + " 必须是 C/S/T/P 之一"
				case "relType":
					errors[field] = field + " 必须是 PREREQ/SUP_SKILL/THINK_PAT 之一"
				default:
					errors[field] = field + " 校验失败 (" + tag + ")"
				}
			}

			c.JSON(400, response.ValidateErr(c.Request.Context(), errors))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateQuery 校验查询参数
func ValidateQuery(s interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validate == nil {
			InitValidator()
		}

		if err := c.ShouldBindQuery(s); err != nil {
			c.JSON(400, response.ValidateErr(c.Request.Context(), err.Error()))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateJSON 校验 JSON 请求体
func ValidateJSON(s interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validate == nil {
			InitValidator()
		}

		if err := c.ShouldBindJSON(s); err != nil {
			c.JSON(400, response.ValidateErr(c.Request.Context(), err.Error()))
			c.Abort()
			return
		}

		// 存储校验对象供后续中间件使用
		c.Set("validatorObject", s)

		c.Next()
	}
}

// GetValidator 获取校验器实例
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}
