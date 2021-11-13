package jwt

import "strings"

type TokenService struct {
	Prefix string
}

func NewTokenService(opts...string) *TokenService {
	prefix := "Bearer "
	if len(opts) > 0 && len(opts[0]) > 0 {
		prefix = opts[0]
	}
	return &TokenService{Prefix: prefix}
}
func NewCookieTokenService() *TokenService {
	return &TokenService{Prefix: ""}
}
func (t *TokenService) GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error) {
	return GenerateToken(payload, secret, expiresIn)
}

func (t *TokenService) VerifyToken(token string, secret string) (map[string]interface{}, int64, int64, error) {
	payload, c, err := VerifyToken(token, secret)
	return payload, c.IssuedAt, c.ExpiresAt, err
}

func (t *TokenService) GetAndVerifyToken(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error) {
	if len(t.Prefix) > 0 {
		if strings.HasPrefix(authorization, t.Prefix) == false {
			return false, "", nil, 0, 0, nil
		}
	}
	token := authorization[len(t.Prefix):]
	payload, c, err := VerifyToken(token, secret)
	return true, token, payload, c.IssuedAt, c.ExpiresAt, err
}
