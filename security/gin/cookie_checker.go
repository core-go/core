package gin

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (h *CookieChecker) Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		tokenCookie, err := r.Cookie(h.Token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Cookie error: " + err.Error())
			return
		}
		if tokenCookie == nil {
			ctx.JSON(http.StatusUnauthorized, h.Token + " is required in cookies")
			return
		}
		authorization := tokenCookie.Value
		if len(authorization) == 0 {
			ctx.JSON(http.StatusUnauthorized, h.Token + " is required in cookies")
			return
		}
		isToken, token, data, issuedAt, _, err := h.GetAndVerifyToken(authorization, h.Secret)
		if !isToken || err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid Authorization token")
			return
		}
		if data == nil {
			data = make(map[string]interface{})
		}
		iat := time.Unix(issuedAt, 0)
		data["token"] = token
		data["issuedAt"] = iat
		var ctx2 context.Context
		ctx2 = r.Context()
		if len(h.Ip) > 0 {
			ip := getRemoteIp(r)
			ctx2 = context.WithValue(ctx2, h.Ip, ip)
		}
		if h.CheckBlacklist != nil {
			user := ValueFromMap(h.Key, data)
			reason := h.CheckBlacklist(user, token, iat)
			if len(reason) > 0 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "token is not valid anymore")
			} else {
				if h.CheckWhitelist != nil {
					valid := h.CheckWhitelist(user, token)
					if !valid {
						ctx.AbortWithStatusJSON(http.StatusUnauthorized, "token is not valid anymore")
						return
					}
				}
				if len(h.Authorization) > 0 {
					ctx2 = context.WithValue(ctx2, h.Authorization, data)
					ctx.Request = r.WithContext(ctx2)
					ctx.Next()
				} else {
					for k, e := range data {
						if len(k) > 0 {
							ctx2 = context.WithValue(ctx2, k, e)
						}
					}
					ctx.Request = r.WithContext(ctx2)
					ctx.Next()
				}
			}
		} else {
			if h.CheckWhitelist != nil {
				user := ValueFromMap(h.Key, data)
				valid := h.CheckWhitelist(user, token)
				if !valid {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, "token is not valid anymore")
					return
				}
			}
			if len(h.Authorization) > 0 {
				ctx2 = context.WithValue(ctx2, h.Authorization, data)
				ctx.Request = r.WithContext(ctx2)
				ctx.Next()
			} else {
				for k, e := range data {
					if len(k) > 0 {
						ctx2 = context.WithValue(ctx2, k, e)
					}
				}
				ctx.Request = r.WithContext(ctx2)
				ctx.Next()
			}
		}
	}
}
