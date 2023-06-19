package util

import "golang.org/x/crypto/bcrypt"

// GetBcryptPwd 给密码加密
func GetBcryptPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(Str2Byte(pwd), bcrypt.DefaultCost)
	return Byte2Str(hash), err
}

// ComparePwd 比对密码
func ComparePwd(after string, before string) bool {
	// Returns true on success, after is for the database.
	err := bcrypt.CompareHashAndPassword(Str2Byte(after), Str2Byte(before))
	if err != nil {
		return false
	} else {
		return true
	}
}
