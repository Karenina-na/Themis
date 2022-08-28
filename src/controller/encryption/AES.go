package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AESEncrypt(data interface{}, key []byte) string {
	origData := []byte(data.(string))
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData = _PKCS7Padding(origData, blockSize)
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blocMode.CryptBlocks(crypted, origData)
	return string(crypted)
}

func AESDecrypt(cypted string, key []byte) string {
	cryptedByte := []byte(cypted)
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	blockMode.CryptBlocks(origData, cryptedByte)
	origData = _PKCS7UnPadding(origData)
	return string(origData)
}

func _PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func _PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	} else {
		padding := int(origData[length-1])
		return origData[:(length - padding)]
	}
}
