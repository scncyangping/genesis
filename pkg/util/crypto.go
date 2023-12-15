// @Author: YangPing
// @Create: 2023/10/23
// @Description: 加解密工具类

package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	AesIv     = "DEF_FAULT=CF(Aiv"
	AesKey    = "DEF_FAULT=CF(Key"
	DefaultIv = "REVGX0ZBVUxUPUNGKEFpdg==:"
)

func AESEncrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher([]byte(AesKey))
	if err != nil {
		return "", err
	}

	// 填充明文以适应块大小
	plaintext = PKCS7Padding(plaintext, aes.BlockSize)

	// 使用CBC模式加密
	mode := cipher.NewCBCEncrypter(block, []byte(AesIv))
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 将IV和密文组合成一个字符串返回
	return base64.StdEncoding.EncodeToString([]byte(AesIv)) + ":" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESDecrypt(ciphertext string) (string, error) {
	parts := strings.Split(ciphertext, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid ciphertext format")
	}

	iv, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	ciphertextBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(AesKey))
	if err != nil {
		return "", err
	}

	// 使用CBC模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	// 去除填充
	plaintext = PKCS7UnPadding(plaintext)

	return string(plaintext), nil
}

// PKCS7Padding 添加PKCS#7填充到明文
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去除PKCS#7填充
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
