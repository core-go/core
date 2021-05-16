package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type Decrypter interface {
	Decrypt(cipherText []byte, iv []byte) ([]byte, error)
}

type DefaultDecrypter struct {
	Block cipher.Block
}

func NewDecrypter(secretKey []byte) (*DefaultDecrypter, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	return &DefaultDecrypter{Block: block}, nil
}

func (d *DefaultDecrypter) Decrypt(cipherText []byte, iv []byte) ([]byte, error) {
	return Decrypt(cipherText, iv, d.Block), nil
}
