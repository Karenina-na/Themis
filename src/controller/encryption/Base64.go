package encryption

import (
	"encoding/base64"
	"fmt"
)

func Base64(data interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", data)))
}

func Base64Decode(data interface{}) string {
	decodeString, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", data))
	return string(decodeString)
}
