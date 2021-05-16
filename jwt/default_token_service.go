package jwt

import "strings"

type DefaultTokenService struct {
}

func NewTokenService() *DefaultTokenService {
	return &DefaultTokenService{}
}
func (t *DefaultTokenService) GenerateToken(payload interface{}, secret string, expiresIn int64) (string, error) {
	return GenerateToken(payload, secret, expiresIn)
}

func (t *DefaultTokenService) VerifyToken(token string, secret string) (map[string]interface{}, int64, int64, error) {
	payload, c, err := VerifyToken(token, secret)
	return payload, c.IssuedAt, c.ExpiresAt, err
}

func (t *DefaultTokenService) GetAndVerifyToken(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error) {
	if strings.HasPrefix(authorization, "Bearer ") == false {
		return false, "", nil, 0, 0, nil
	}
	token := authorization[7:]
	payload, c, err := VerifyToken(token, secret)
	return true, token, payload, c.IssuedAt, c.ExpiresAt, err
}
