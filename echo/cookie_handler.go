package echo

import (
	"context"
	"github.com/labstack/echo/v4"
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

func (c *CookieHandler) HandleAuthorization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			tokenCookie, err := r.Cookie(c.Token)
			if err != nil || tokenCookie == nil {
				return next(ctx)
			} else {
				authorization := tokenCookie.Value
				isToken, _, data, _, _, err := c.GetAndVerifyToken(authorization, c.Secret)
				var ctx2 context.Context
				ctx2 = r.Context()
				if len(c.Ip) > 0 {
					ip := GetRemoteIp(r)
					ctx2 = context.WithValue(ctx2, c.Ip, ip)
				}
				if !isToken {
					if len(c.Ip) == 0 {
						return next(ctx)
					} else {
						ctx.SetRequest(r.WithContext(ctx2))
						return next(ctx)
					}
				} else {
					if err != nil {
						if len(c.Ip) == 0 {
							return next(ctx)
						} else {
							ctx.SetRequest(r.WithContext(ctx2))
							return next(ctx)
						}
					} else {
						if len(c.Authorization) > 0 {
							ctx2 := context.WithValue(ctx2, c.Authorization, data)
							ctx.SetRequest(r.WithContext(ctx2))
							return next(ctx)
						} else {
							for k, e := range data {
								if len(k) > 0 {
									ctx2 = context.WithValue(ctx2, k, e)
								}
							}
							ctx.SetRequest(r.WithContext(ctx2))
							return next(ctx)
						}
					}
				}
			}
		}
	}
}
