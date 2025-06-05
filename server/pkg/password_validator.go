package pkg

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passwd", passwordValidator)
	}
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasMinLen := len(password) >= 8
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_.,<>?/{}\-]`).MatchString(password)

	return hasMinLen && hasUpper && hasNumber && hasSpecial
}
