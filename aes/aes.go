package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func Encrypt(plainData string, initVector []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(initVector) != aes.BlockSize {
		return "", fmt.Errorf("init vector size must be %d bytes", aes.BlockSize)
	}
	encrypter := cipher.NewCBCEncrypter(block, initVector)
	padData := PKCS5Pad([]byte(plainData), aes.BlockSize)
	ciphertext := make([]byte, len(padData))
	encrypter.CryptBlocks(ciphertext, padData)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
func Encode(plaintext string, key string, iv string, blockSize int) (string, error) {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Pad([]byte(plaintext), blockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	encrypter := cipher.NewCBCEncrypter(block, bIV)
	encrypter.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext), nil
}
func PKCS5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func Decrypt(base64EncryptedData string, initVector []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(initVector) != aes.BlockSize {
		return "", fmt.Errorf("init vector size must be %d bytes", aes.BlockSize)
	}

	cipherText, err := base64.StdEncoding.DecodeString(base64EncryptedData)
	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, initVector)
	data := make([]byte, len(cipherText))
	decrypter.CryptBlocks(data, cipherText)

	data, err = PKCS5Unpad(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func Decode(cipherText string, encKey string, iv string) (string, error) {
	if len(iv) != 16 {
		return "", errors.New("length IV != 16")
	}
	if cipherText == "" {
		return "", errors.New("invalid data")
	}
	key := []byte(encKey)
	bIV := []byte(iv)
	decodedCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, bIV)
	decrypter.CryptBlocks(decodedCipherText, decodedCipherText)
	decodedCipherText, err = PKCS5Unpad(decodedCipherText)
	if err != nil {
		return "", err
	}
	return string(decodedCipherText), nil
}
func PKCS5Unpad(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("data is empty")
	}
	unpad := int(src[len(src)-1])
	if unpad > len(src) || unpad == 0 {
		return nil, errors.New("invalid padding")
	}
	pad := bytes.Repeat([]byte{byte(unpad)}, unpad)
	if bytes.HasSuffix(src, pad) {
		return src[:len(src)-unpad], nil
	}
	return nil, errors.New("padding mismatch")
}
