package skm

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
)

var iv = "VqEVUW_7MVc3)yyV"

func AES256Encrypt(s string, key string) (string, error) {
	keySHA256 := sha256.Sum256([]byte(key))
	cipher, err := aes.NewCipher(keySHA256[0:32])
	if err != nil {
		return "", err
	}

	padding := PKCS7Padding([]byte(s), aes.BlockSize)

	encryptData := make([]byte, len(padding))
	cipher.Encrypt(encryptData, padding)

	return base64.StdEncoding.EncodeToString(encryptData), nil

}

func AES256Decrypt(s string, key string) (string, error) {

	encryptData, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	keySHA256 := sha256.Sum256([]byte(key))
	cipher, err := aes.NewCipher(keySHA256[0:32])
	if err != nil {
		return "", err
	}

	decryptData := make([]byte, len(encryptData))
	cipher.Decrypt(decryptData, encryptData)

	unPadding := PKCS7UnPadding(decryptData)

	return string(unPadding), nil

}

func PKCS7Padding(data []byte, blockSize int) []byte {
	paddingSize := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(data, paddingText...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
