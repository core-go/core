package authorization

import (
	"context"
	"net/http"
)

type CookieHandler struct {
	GetAndVerifyToken func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Token             string
	Secret            string
	Ip                string
	Authorization     string
}

func NewCookieHandler(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, options ...string) *CookieHandler {
	return NewCookieHandlerWithIp(verifyToken, secret, "", options...)
}

func NewCookieHandlerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, options ...string) *CookieHandler {
	var authorization string
	token := "token"
	if len(options) > 0 {
		authorization = options[0]
	}
	if len(options) > 1 {
		token = options[1]
	}
	return &CookieHandler{Authorization: authorization, GetAndVerifyToken: verifyToken, Secret: secret, Token: token, Ip: ip}
}

func (c *CookieHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(c.Token)
		if err != nil || tokenCookie == nil {
			next.ServeHTTP(w, r)
		} else {
			authorization := tokenCookie.Value
			isToken, _, data, _, _, err := c.GetAndVerifyToken(authorization, c.Secret)
			var ctx context.Context
			ctx = r.Context()
			if len(c.Ip) > 0 {
				ip := GetRemoteIp(r)
				ctx = context.WithValue(ctx, c.Ip, ip)
			}
			if !isToken {
				if len(c.Ip) == 0 {
					next.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				if err != nil {
					if len(c.Ip) == 0 {
						next.ServeHTTP(w, r)
					} else {
						next.ServeHTTP(w, r.WithContext(ctx))
					}
				} else {
					if len(c.Authorization) > 0 {
						ctx := context.WithValue(ctx, c.Authorization, data)
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
			}
		}
	})
}
