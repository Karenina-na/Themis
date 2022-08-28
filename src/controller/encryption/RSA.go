package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash"
)

func RSACreateKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey
}

func RSAEncrypt(publicKey *rsa.PublicKey, data interface{}, al hash.Hash) string {
	oaep, _ := rsa.EncryptOAEP(al, rand.Reader, publicKey, []byte(fmt.Sprintf("%v", data)), nil)
	return string(oaep)
}
func RSADecrypt(privateKey *rsa.PrivateKey, encryptedBytes string, al crypto.Hash) string {
	decryptedBytes, _ := privateKey.Decrypt(nil, []byte(fmt.Sprintf("%v", encryptedBytes)), &rsa.OAEPOptions{Hash: al})
	return string(decryptedBytes)
}
