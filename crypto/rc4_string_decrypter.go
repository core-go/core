package crypto

type RC4Decrypter struct {
}

func NewRC4Decrypter() *RC4Decrypter{
	return &RC4Decrypter{}
}

func (b *RC4Decrypter) Decrypt(cipherText string, secretKey string) (string, error) {
	bytes, err := DecryptRC4([]byte(cipherText), []byte(secretKey))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
