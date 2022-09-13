package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AESEncrypt
// @Description: AES加密
// @param        data 待加密数据
// @param        key  加密密钥
// @return       string 加密后数据
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

// AESDecrypt
// @Description: AES解密
// @param        cypted 待解密数据
// @param        key    解密密钥
// @return       string 解密后数据
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

// _PKCS7Padding
// @Description: PKCS7填充
// @param        ciphertext 待填充数据
// @param        blockSize  块大小
// @return       []byte     填充后数据
func _PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// _PKCS7UnPadding
// @Description: PKCS7去填充
// @param        origData 待去填充数据
// @return       []byte   去填充后数据
func _PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	} else {
		padding := int(origData[length-1])
		return origData[:(length - padding)]
	}
}
