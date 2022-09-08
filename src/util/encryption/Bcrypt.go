package encryption

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//
// Bcrypt
// @Description: Bcrypt加密
// @param        data 待加密数据
// @return       string 加密后数据
//
func Bcrypt(data interface{}) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%v", data)), 14)
	return string(bytes)
}
