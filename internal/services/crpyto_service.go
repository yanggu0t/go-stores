package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type CryptoService struct {
	key []byte
}

func NewCryptoService(encodedKey string) *CryptoService {
	key, _ := base64.StdEncoding.DecodeString(encodedKey)
	return &CryptoService{key: key}
}

func (s *CryptoService) Decrypt(encryptedString string) (string, error) {
	// 步驟 1: Base64 解碼
	decodedData, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	// 檢查解碼後的數據長度
	if len(decodedData) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// 步驟 2: 提取 IV（前 16 字節）
	iv := decodedData[:aes.BlockSize]

	// 步驟 3: 提取密文（剩餘的字節）
	ciphertext := decodedData[aes.BlockSize:]

	// 創建 cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// 步驟 4: 使用 AES-CBC 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 步驟 5: 去除 PKCS7 填充
	unpaddedPlaintext, err := pkcs7Unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(unpaddedPlaintext), nil
}

// PKCS7 去除填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid padding")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}
