package crypto

import "crypto/rc4"

var rc4Cache = make(map[string]*rc4.Cipher, 0)

func GetRC4CipherFromBytes(secretKey []byte) (*rc4.Cipher, error) {
	s := string(secretKey)
	return GetRC4Cipher(s)
}
func GetRC4Cipher(secretKey string) (*rc4.Cipher, error) {
	c := rc4Cache[secretKey]
	if c != nil {
		return c, nil
	}
	b := []byte(secretKey)
	return rc4.NewCipher(b)
}
func DecryptRC4ByCipher(cipherText []byte, cipher *rc4.Cipher) ([]byte, error) {
	dst := make([]byte, len(cipherText))
	cipher.XORKeyStream(dst, cipherText)
	return dst, nil
}
func EncryptRC4ByCipher(plaintText []byte, cipher *rc4.Cipher) ([]byte, error) {
	dst := make([]byte, len(plaintText))
	cipher.XORKeyStream(dst, plaintText)
	return dst, nil
}
func DecryptRC4(cipherText []byte, secretKey []byte) ([]byte, error) {
	if cipher, err := GetRC4CipherFromBytes(secretKey); err != nil {
		return nil, err
	} else {
		return DecryptRC4ByCipher(cipherText, cipher)
	}
}
func EncryptRC4(plaintText []byte, secretKey []byte) ([]byte, error) {
	if cipher, err := GetRC4CipherFromBytes(secretKey); err != nil {
		return nil, err
	} else {
		return EncryptRC4ByCipher(plaintText, cipher)
	}
}
