package validation

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type EmailValidator struct {
	Email string `json:"email" binding:"required,email" msg:"Email: Email必须为合法的邮箱"`
}

type UserNameValidator struct {
	UserName string `json:"username" binding:"required,alphanumunicode,min=1,max=20" msg:"UserName: 用户名长度在1-20之间，仅包含Unicode字符与数字"`
}

type PasswordValidator struct {
	Password string `json:"password" binding:"required,ascii,min=8,max=72,noCRLF" msg:"Password: 用户密码只能由ASCII字符组成，但不能包含换行符和回车，长度在8到72之间"`
}

type UpdatePasswordValidator struct {
	PasswordValidator
	RePassword string `json:"re-password" binding:"required,eqfield=PasswordValidator.Password" msg:"RePassword: 两次密码输入需要一致"`
}

type UserLoginValidator struct {
	EmailValidator
	PasswordValidator
}

type UserRegisterValidator struct {
	UserNameValidator
	PasswordValidator
	RePassword string `json:"re-password" binding:"required,eqfield=PasswordValidator.Password" msg:"RePassword: 两次密码输入需要一致"`
}

// user defined validator
// 定义密码验证
func noCRLF(fl validator.FieldLevel) bool {
	passwd := fl.Field().String()
	if strings.Contains(passwd, "\r") || strings.Contains(passwd, "\n") {
		return false
	}
	return true
}
