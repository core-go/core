package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

func GenerateKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	publicKey := &privateKey.PublicKey

	sPrivate, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	sPublic, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(sPrivate), base64.StdEncoding.EncodeToString(sPublic), nil
}
func Decrypt(v string, privateKey string) (string, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(privateKey)))
	_, err := base64.StdEncoding.Decode(dst, []byte(privateKey))
	if err != nil {
		return "", err
	}

	privateKeyPKCS, err := x509.ParsePKCS8PrivateKey(dst)
	if err != nil {
		return "", err
	}

	dst2, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}

	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKeyPKCS.(*rsa.PrivateKey), dst2)
	if err != nil {
		return "", err
	}
	plainPin := string(res)
	return plainPin, nil
}
