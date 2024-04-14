package jwt

import "strings"

type TokenAdapter struct {
	Prefix string
}
func NewTokenService(opts...string) *TokenAdapter {
	return NewTokenAdapter(opts...)
}
func NewTokenAdapter(opts...string) *TokenAdapter {
	prefix := "Bearer "
	if len(opts) > 0 && len(opts[0]) > 0 {
		prefix = opts[0]
	}
	return &TokenAdapter{Prefix: prefix}
}
func NewCookieTokenService() *TokenAdapter {
	return NewCookieTokenAdapter()
}
func NewCookieTokenAdapter() *TokenAdapter {
	return &TokenAdapter{Prefix: ""}
}
func (t *TokenAdapter) GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error) {
	return GenerateToken(payload, secret, expiresIn)
}

func (t *TokenAdapter) VerifyToken(token string, secret string) (map[string]interface{}, int64, int64, error) {
	payload, c, err := VerifyToken(token, secret)
	return payload, c.IssuedAt, c.ExpiresAt, err
}

func (t *TokenAdapter) GetAndVerifyToken(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error) {
	if len(t.Prefix) > 0 {
		if strings.HasPrefix(authorization, t.Prefix) == false {
			return false, "", nil, 0, 0, nil
		}
	}
	token := authorization[len(t.Prefix):]
	payload, c, err := VerifyToken(token, secret)
	return true, token, payload, c.IssuedAt, c.ExpiresAt, err
}
