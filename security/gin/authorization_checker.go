package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	Authorization = "authorization"
	Uid           = "uid"
	UserId        = "userId"
	UserName      = "userName"
	Username      = "username"
	UserType      = "userType"
	Roles         = "roles"
	Privileges    = "privileges"
	Permission    = "permission"
	Permissions   = "permissions"
	Ip            = "ip"
)

type AuthorizationChecker struct {
	GetAndVerifyToken func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Secret            string
	Ip                string
	CheckBlacklist    func(id string, token string, createAt time.Time) string
	Authorization     string
	Key               string
	CheckWhitelist    func(id string, token string) bool
}

func NewDefaultAuthorizationChecker(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, key string, options ...string) *AuthorizationChecker {
	return NewAuthorizationCheckerWithIp(verifyToken, secret, "", nil, nil, key, options...)
}
func NewAuthorizationChecker(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, checkToken func(string, string, time.Time) string, key string, options ...string) *AuthorizationChecker {
	return NewAuthorizationCheckerWithIp(verifyToken, secret, "", checkToken, nil, key, options...)
}
func NewAuthorizationCheckerWithWhitelist(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, checkToken func(string, string, time.Time) string, checkWhitelist func(string, string) bool, key string, options ...string) *AuthorizationChecker {
	return NewAuthorizationCheckerWithIp(verifyToken, secret, "", checkToken, checkWhitelist, key, options...)
}
func NewAuthorizationCheckerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, checkToken func(string, string, time.Time) string, checkWhitelist func(string, string) bool, key string, options ...string) *AuthorizationChecker {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &AuthorizationChecker{Authorization: authorization, Key: key, CheckBlacklist: checkToken, GetAndVerifyToken: verifyToken, Secret: secret, Ip: ip, CheckWhitelist: checkWhitelist}
}

func (h *AuthorizationChecker) Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		au := r.Header["Authorization"]
		if len(au) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "'Authorization' is required in http request header.")
			return
		}
		authorization := au[0]
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
