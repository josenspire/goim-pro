package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const (
	iv = "0102030405060708"
)

// NewSHA256 - create
func NewSHA256(data string, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

func CreateAESSecretKey(size int) string {
	return GenerateRandString(16)
}

// AES加密的具体算法为: AES-128-CBC，输出格式为 base64
// AES加密时需要指定 iv：0102030405060708
// AES加密时需要 padding
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
// https://github.com/darknessomi/musicbox/wiki/%E7%BD%91%E6%98%93%E4%BA%91%E9%9F%B3%E4%B9%90%E6%96%B0%E7%99%BB%E5%BD%95API%E5%88%86%E6%9E%90
func AESEncrypt(encodeStr string, secretKeyStr string) (string, error) {
	secretKey := []byte(secretKeyStr)
	encodeBytes := []byte(encodeStr)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		logger.Errorf("[aes] create cipher error: %v", err)
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = pKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	cipherText := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(cipherText, encodeBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func AESDecrypt(decodeStr string, secretKeyStr string) (string, error) {
	// decode base64
	decodeBytes, _ := base64.StdEncoding.DecodeString(decodeStr)

	secretKey := []byte(secretKeyStr)
	block, _ := aes.NewCipher(secretKey)

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(originData, decodeBytes)
	originData = pKCS5UnPadding(originData)
	return string(originData[:]), nil
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize // 16, 32, 48 etc..
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

func pKCS5UnPadding(originData []byte) []byte {
	length := len(originData)
	unPadding := int(originData[length-1])
	return originData[:(length - unPadding)]
}
