package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func Decrypt(base64EncryptedAESKey string, privateKeyStr string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	encryptedData, err := base64.StdEncoding.DecodeString(base64EncryptedAESKey)
	if err != nil {
		return nil, err
	}

	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey.(*rsa.PrivateKey), encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func GetPKCS1Signature(data string, privateKeyStr string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	hashed := sha1.Sum([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func Encrypt(key []byte, publicKeyStr string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func GenRandomArray(size int) ([]byte, error) {
	random := make([]byte, size)
	_, err := rand.Read(random)
	if err != nil {
		return nil, err
	}
	return random, nil
}

func ParseRsaPublicKeyFromPem(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, errors.New("key type is not RSA")
}
func ParseRsaPrivateKeyFromPem(privatePEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func CreateSignature(msg string, privateKey string) (string, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	pri, err := ParseRsaPrivateKeyFromPem(privateKey)
	if err != nil {
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, pri, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return "", err
	}
	return string(signature), nil
}

func Decode(data string, signature string, iv string, privateKey string, aes string, decode func(cipherText string, encKey string, iv string) (string, error)) (string, bool, error) {
	plainData, err := decode(data, aes, iv)
	if err != nil {
		return "", false, err
	}

	b64signature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", false, err
	}

	priKey, err := ParseRsaPrivateKeyFromPem(privateKey)
	if err != nil {
		return "", false, err
	}

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(plainData))
	if err != nil {
		return "", false, err
	}

	msgHashSum := msgHash.Sum(nil)
	hashData := hex.EncodeToString([]byte(msgHashSum))

	decryptedSignature, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, b64signature)
	if err != nil {
		return "", false, err
	}

	verify := string(decryptedSignature) == hashData
	if !verify {
		return "", false, err
	}

	return plainData, true, err
}
