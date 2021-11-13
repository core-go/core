package authorization

import (
	"context"
	"net"
	"net/http"
)

type Handler struct {
	GetAndVerifyToken func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Secret            string
	Ip                string
	Authorization     string
}

func NewHandler(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, options ...string) *Handler {
	return NewHandlerWithIp(verifyToken, secret, "", options...)
}

func NewHandlerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, options ...string) *Handler {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &Handler{Authorization: authorization, GetAndVerifyToken: verifyToken, Secret: secret, Ip: ip}
}

func (c *Handler) HandleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au := r.Header["Authorization"]
		if au == nil || len(au) == 0 {
			next.ServeHTTP(w, r)
		} else {
			authorization := au[0]
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

func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
