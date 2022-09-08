package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"hash"
)

//
// RSACreateKey
// @Description: Create RSA key
// @return       *rsa.PrivateKey Private key
// @return       *rsa.PublicKey  Public key
//
func RSACreateKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey
}

//
// RSAEncrypt
// @Description: Encrypt data
// @param        publicKey *rsa.PublicKey Public key
// @param        data      string         Data   to encrypt
// @param        al        crypto.Hash    Algorithm
// @return       string    Encrypted data
//
func RSAEncrypt(publicKey *rsa.PublicKey, data interface{}, al hash.Hash) string {
	oaep, _ := rsa.EncryptOAEP(al, rand.Reader, publicKey, []byte(fmt.Sprintf("%v", data)), nil)
	return string(oaep)
}

//
// RSADecrypt
// @Description: Decrypt data
// @param        privateKey     *rsa.PrivateKey Private   key
// @param        encryptedBytes []byte          Encrypted data
// @param        al             crypto.Hash     Algorithm
// @return       string         Decrypted data
//
func RSADecrypt(privateKey *rsa.PrivateKey, encryptedBytes string, al crypto.Hash) string {
	decryptedBytes, _ := privateKey.Decrypt(nil, []byte(fmt.Sprintf("%v", encryptedBytes)), &rsa.OAEPOptions{Hash: al})
	return string(decryptedBytes)
}
