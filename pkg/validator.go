package pkg

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 密码验证 只能包含字母、数字、下划线和连字符
func PasswordValidate(fl validator.FieldLevel) bool {
	pattern := `^[a-zA-Z0-9_-]+$`
	value := fl.Field().String()
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}
	return matched
}
