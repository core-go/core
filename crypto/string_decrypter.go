package crypto

type StringDecrypter interface {
	Decrypt(cipherText string, secretKey string) (string, error)
}

func NewStringDecrypter(decrypter Decrypter, toBytes func(string) []byte, toString func([]byte) string) StringDecrypter {
	return &DefaultStringDecrypter{Decrypter: decrypter, ToBytes: toBytes, ToString: toString}
}

type DefaultStringDecrypter struct {
	Decrypter Decrypter
	ToBytes   func(string) []byte
	ToString  func([]byte) string
}

func (d *DefaultStringDecrypter) Decrypt(cipherText string, secretKey string) (string, error) {
	var b1, b2 []byte
	if d.ToBytes == nil {
		b1 = []byte(cipherText)
		b2 = []byte(secretKey)
	} else {
		b1 = d.ToBytes(cipherText)
		b2 = d.ToBytes(secretKey)
	}
	bs, err := d.Decrypter.Decrypt(b1, b2)
	if err != nil {
		return "", err
	}
	if d.ToString != nil {
		return d.ToString(bs), nil
	}
	return string(bs), nil
}
