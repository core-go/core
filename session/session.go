package session

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type SecureSession struct {
	Secret string
}

func New(secret string) SecureSession {
	return SecureSession{Secret: secret}
}

func (s SecureSession) EncodeSessionID(sid string) string {
	b := base64.StdEncoding.EncodeToString([]byte(sid))
	sf := fmt.Sprintf("%s.%s", b, signature(sid, s.Secret))
	return url.QueryEscape(sf)
}

func (s SecureSession) DecodeSessionID(value string) (string, error) {
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", err
	}

	values := strings.Split(value, ".")
	if len(values) != 2 {
		return "", errors.New("invalid session id")
	}

	bsid, err := base64.StdEncoding.DecodeString(values[0])
	if err != nil {
		return "", err
	}
	sid := string(bsid)

	signData := signature(sid, s.Secret)
	if signData != values[1] {
		return "", errors.New("invalid session id")
	}
	return sid, nil
}

func signature(sid string, sign string) string {
	h := hmac.New(sha1.New, []byte(sign))
	h.Write([]byte(sid))
	return fmt.Sprintf("%x", h.Sum(nil))
}
