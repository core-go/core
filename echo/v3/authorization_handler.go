package echo

import (
	"context"
	"github.com/labstack/echo"
	"net"
	"net/http"
)

type AuthorizationHandler struct {
	GetAndVerifyToken func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Secret            string
	Ip                string
	Authorization     string
}

func NewAuthorizationHandler(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, options ...string) *AuthorizationHandler {
	return NewAuthorizationHandlerWithIp(verifyToken, secret, "", options...)
}

func NewAuthorizationHandlerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, options ...string) *AuthorizationHandler {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &AuthorizationHandler{Authorization: authorization, GetAndVerifyToken: verifyToken, Secret: secret, Ip: ip}
}

func (c *AuthorizationHandler) HandleAuthorization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()

			au := r.Header["Authorization"]
			if au == nil || len(au) == 0 {
				return next(ctx)
			} else {
				authorization := au[0]
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
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
