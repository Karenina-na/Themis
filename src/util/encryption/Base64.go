package encryption

import (
	"encoding/base64"
	"fmt"
)

// Base64
// @Description: Base64加密
// @param        data 加密数据
// @return       string 加密后的数据
func Base64(data interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", data)))
}

// Base64Decode
// @Description: Base64解密
// @param        data 解密数据
// @return       string 解密后的数据
func Base64Decode(data interface{}) string {
	decodeString, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", data))
	return string(decodeString)
}
