package security

import (
	"strconv"
	"strings"
	"time"
)

const joinChar = "-"

type DefaultTokenBlacklistTokenService struct {
	CacheService CacheService
	TokenPrefix  string
	TokenExpires int64
}

func NewTokenBlacklistTokenService(cacheService CacheService, keyPrefix string, tokenExpires int64) *DefaultTokenBlacklistTokenService {
	return &DefaultTokenBlacklistTokenService{CacheService: cacheService, TokenPrefix: keyPrefix, TokenExpires: tokenExpires}
}

func (s *DefaultTokenBlacklistTokenService) generateKey(token string) string {
	return s.TokenPrefix + token
}

func (s *DefaultTokenBlacklistTokenService) generateKeyForId(id string) string {
	return s.TokenPrefix + id
}

func (s *DefaultTokenBlacklistTokenService) Revoke(token string, reason string, expiredDate time.Time) error {
	key := s.generateKey(token)
	var value string
	if len(reason) > 0 {
		value = reason
	} else {
		value = ""
	}

	today := time.Now()
	expiresInSecond := expiredDate.Sub(today)
	if expiresInSecond <= 0 {
		return nil // Token already expires, don't need add to cache
	} else {
		return s.CacheService.Put(key, value, expiresInSecond*time.Second)
	}
}

func (s *DefaultTokenBlacklistTokenService) RevokeAllTokens(id string, reason string) error {
	key := s.generateKeyForId(id)
	today := time.Now()
	value := reason + joinChar + strconv.Itoa(int(today.Unix()))
	return s.CacheService.Put(key, value, time.Duration(s.TokenExpires)*time.Second)
}

func (s *DefaultTokenBlacklistTokenService) Check(id string, token string, createAt time.Time) string {
	idKey := s.generateKeyForId(id)
	tokenKey := s.generateKey(token)

	keys := []string{idKey, tokenKey}
	value, _, err := s.CacheService.GetManyStrings(keys)
	if err != nil {
		return ""
	}
	if len(value[idKey]) > 0 {
		index := strings.Index(value[idKey], joinChar)
		reason := value[idKey][0:index]
		strDate := value[idKey][index+1:]
		i, err := strconv.ParseInt(strDate, 10, 64)
		if err == nil {
			tmDate := time.Unix(i, 0)
			if tmDate.Sub(createAt) > 0 {
				return reason
			}
		}
	}
	if len(value[tokenKey]) > 0 {
		return value[tokenKey]
	}
	return ""
}
