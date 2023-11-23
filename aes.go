package gotools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unfading := int(origData[length-1])
	return origData[:(length - unfading)]
}

// AesEncrypt AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AesDecrypt AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// EncryptString 加密 字符串
func EncryptString(key, originData string) (string, error) {
	encryptedData, err := AesEncrypt([]byte(originData), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// DecryptString 解密 字符串
func DecryptString(key, encryptedData string) (string, error) {
	encrypted, _ := base64.StdEncoding.DecodeString(encryptedData)
	originData, err := AesDecrypt([]byte(encrypted), []byte(key))
	if err != nil {
		return "", err
	}
	return string(originData), nil
}
