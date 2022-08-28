package encryption

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func Sha1(data interface{}) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

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
