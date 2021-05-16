package crypto

import "encoding/base64"

func EncodeBase64(plainText []byte) string {
	return base64.StdEncoding.EncodeToString(plainText)
}

func EncodeBase64FromString(plainText string) string {
	return EncodeBase64([]byte(plainText))
}

func DecodeBase64(cipherText []byte) []byte {
	plainText := make([]byte, len(cipherText))
	base64.StdEncoding.Decode(plainText, cipherText)
	return plainText
}

func DecodeBase64FromString(cipherText string) []byte {
	return DecodeBase64([]byte(cipherText))
}
