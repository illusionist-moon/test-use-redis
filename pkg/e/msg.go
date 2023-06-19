package e

var MsgFlags = map[int]string{
	Success:       "success",
	Error:         "fail",
	InvalidParams: "参数错误",

	ErrorExistUser:    "用户名已存在",
	ErrorNotExistUser: "该用户不存在",
	ErrorIncorrectPwd: "用户存在但密码错误",

	ErrorAuthCheckTokenFail:    "Token鉴权失败",
	ErrorAuthCheckTokenTimeout: "Token已超时",
	ErrorAuthToken:             "Token生成失败",

	PageNotFound: "Page not found",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
