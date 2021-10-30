// Package cipher implement configs system using AES Encryption and AES Decryption. Support Generating random 32 bytes key
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"math/big"
	"strings"
)

const defaultKey = "xbmcZMpQoGiRXlTSbHSuYPXynluuNyYh"

func Random() (string, error) {
	var output strings.Builder
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 32; i++ {
		bigIndex, err := rand.Int(rand.Reader, big.NewInt(52))
		if err != nil {
			return "", err
		}
		index := bigIndex.Int64()
		output.WriteString(string(letters[index]))
	}
	return output.String(), nil
}
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	cipherText := make([]byte, aes.BlockSize+len(b))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(b))
	return cipherText, nil
}
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
func Read(filePath string, outer map[string]string, field string, key string) error {
	if key == "" {
		key = defaultKey
	}
	var s strings.Builder
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(in, outer)
	if err != nil {
		return err
	}
	plain, err := Decrypt([]byte(key), []byte(outer[field]))
	if err != nil {
		return err
	}
	s.Write(plain)
	outer[field] = s.String()
	return err
}
func Write(filePath string, inter map[string]string, field string, key string) error {
	if key == "" {
		key = defaultKey
	}
	var s strings.Builder
	ciphered, err := Encrypt([]byte(key), []byte(inter[field]))
	s.Write(ciphered)
	if err != nil {
		return err
	}
	inter[field] = s.String()
	data, err1 := yaml.Marshal(inter)
	if err1 != nil {
		return err1
	}
	err = ioutil.WriteFile(filePath, data, 0666)
	return err
}
