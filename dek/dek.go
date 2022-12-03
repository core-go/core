package dek

import (
	"crypto/aes"
	"crypto/cipher"
)

func Decrypt(firstKey, secondKey, ct, iv string, decode func(rq []byte) ([]byte, error)) (string, error) {
	key := firstKey + secondKey
	keyBytes, err := decode([]byte(key))
	if err != nil {
		return "", err
	}
	text, err := decode([]byte(ct))
	if err != nil {
		return "", err
	}
	no, err := decode([]byte(iv))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, no, text, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
