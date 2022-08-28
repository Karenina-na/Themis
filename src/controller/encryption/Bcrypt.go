package encryption

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Bcrypt(data interface{}) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%v", data)), 14)
	return string(bytes)
}
