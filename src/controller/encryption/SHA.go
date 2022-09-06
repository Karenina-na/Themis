package encryption

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

//
// Sha1
// @Description: sha1加密
// @param        data 加密数据
// @return       string 加密后的数据
//
func Sha1(data interface{}) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//
// Sha256
// @Description: sha256加密
// @param        data 加密数据
// @return       string 加密后的数据
//
func Sha256(data interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha512(data interface{}) string {
	h := sha512.New()
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
