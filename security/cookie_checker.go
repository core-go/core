package security

import (
	"context"
	"net/http"
	"time"
)

type CookieChecker struct {
	GetAndVerifyToken func(token string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Secret            string
	Ip                string
	CheckBlacklist    func(id string, token string, createAt time.Time) string
	Token             string
	Authorization     string
	Key               string
	CheckWhitelist    func(id string, token string) bool
}

func NewDefaultCookieChecker(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, key string, options ...string) *CookieChecker {
	return NewCookieCheckerWithIp(verifyToken, secret, "", nil, nil, key, options...)
}
func NewCookieChecker(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, checkToken func(string, string, time.Time) string, key string, options ...string) *CookieChecker {
	return NewCookieCheckerWithIp(verifyToken, secret, "", checkToken, nil, key, options...)
}
func NewCookieCheckerWithWhitelist(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, checkToken func(string, string, time.Time) string, checkWhitelist func(string, string) bool, key string, options ...string) *CookieChecker {
	return NewCookieCheckerWithIp(verifyToken, secret, "", checkToken, checkWhitelist, key, options...)
}
func NewCookieCheckerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, checkToken func(string, string, time.Time) string, checkWhitelist func(string, string) bool, key string, options ...string) *CookieChecker {
	var authorization string
	if len(options) > 0 {
		authorization = options[0]
	}
	token := "token"
	if len(options) > 1 {
		token = options[1]
	}
	return &CookieChecker{Token: token, Authorization: authorization, Key: key, CheckBlacklist: checkToken, GetAndVerifyToken: verifyToken, Secret: secret, Ip: ip, CheckWhitelist: checkWhitelist}
}

func (h *CookieChecker) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(h.Token)
		if err != nil {
			http.Error(w, "Cookie error: " + err.Error(), http.StatusUnauthorized)
			return
		}
		if tokenCookie == nil {
			http.Error(w, h.Token + " is required in cookies", http.StatusUnauthorized)
			return
		}
		authorization := tokenCookie.Value
		if len(authorization) == 0 {
			http.Error(w, h.Token + " is required in cookies", http.StatusUnauthorized)
			return
		}
		isToken, token, data, issuedAt, _, err := h.GetAndVerifyToken(authorization, h.Secret)
		if !isToken || err != nil {
			http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
			return
		}
		if data == nil {
			data = make(map[string]interface{})
		}
		iat := time.Unix(issuedAt, 0)
		data["token"] = token
		data["issuedAt"] = iat
		var ctx context.Context
		ctx = r.Context()
		if len(h.Ip) > 0 {
			ip := getRemoteIp(r)
			ctx = context.WithValue(ctx, h.Ip, ip)
		}
		if h.CheckBlacklist != nil {
			user := ValueFromMap(h.Key, data)
			reason := h.CheckBlacklist(user, token, iat)
			if len(reason) > 0 {
				http.Error(w, "token is not valid anymore", http.StatusUnauthorized)
			} else {
				if h.CheckWhitelist != nil {
					valid := h.CheckWhitelist(user, token)
					if !valid {
						http.Error(w, "token is not valid anymore", http.StatusUnauthorized)
						return
					}
				}
				if len(h.Authorization) > 0 {
					ctx := context.WithValue(ctx, h.Authorization, data)
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					for k, e := range data {
						if len(k) > 0 {
							ctx = context.WithValue(ctx, k, e)
						}
					}
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			if h.CheckWhitelist != nil {
				user := ValueFromMap(h.Key, data)
				valid := h.CheckWhitelist(user, token)
				if !valid {
					http.Error(w, "token is not valid anymore", http.StatusUnauthorized)
					return
				}
			}
			if len(h.Authorization) > 0 {
				ctx := context.WithValue(ctx, h.Authorization, data)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				for k, e := range data {
					if len(k) > 0 {
						ctx = context.WithValue(ctx, k, e)
					}
				}
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
