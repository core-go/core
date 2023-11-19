package security

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
)

type ICacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, timeToLive time.Duration) (bool, error)
}

type SessionHandler struct {
	PrefixSessionIndex string
	SecretKey          string
	CookieName         string
	VerifyToken        func(tokenString string, secret string) (map[string]interface{}, int64, int64, error)
	Cache              ICacheService
	sessionExpiredTime time.Duration
	LogError           func(ctx context.Context, format string, args ...interface{})
}

func NewSessionHandler(secretKey string, verifyToken func(tokenString string, secret string) (map[string]interface{}, int64, int64, error), cache ICacheService, sessionExpiredTime time.Duration, enableSession bool, logError func(ctx context.Context, format string, args ...interface{}), opts ...string) *SessionHandler {
	var prefixSessionIndex, cookieName string
	if len(opts) > 0 {
		prefixSessionIndex = opts[0]
	} else {
		prefixSessionIndex = "index:"
	}
	if len(opts) > 1 {
		cookieName = opts[1]
	} else {
		cookieName = "id"
	}
	newHandler := &SessionHandler{
		VerifyToken:        verifyToken,
		SecretKey:          secretKey,
		PrefixSessionIndex: prefixSessionIndex,
		CookieName:         cookieName,
		sessionExpiredTime: sessionExpiredTime,
		Cache:              cache,
		LogError:           logError,
	}
	return newHandler
}

func (h *SessionHandler) Handle(next http.Handler, skipRefreshTTL bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := ""
		// case if set sessionID in cookie, need get token from cookie
		cookie, err := r.Cookie(h.CookieName)
		if err != nil {
			http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
			return
		}

		if cookie == nil || cookie.Value == "" {
			http.Error(w, "invalid Authorization token", http.StatusUnauthorized)
			return
		}
		sessionId = cookie.Value
		ctx := r.Context()
		if h.Cache != nil {
			var sessionData map[string]string
			s, err := h.Cache.Get(r.Context(), sessionId)
			if err != nil {
				http.Error(w, "Session is expired", http.StatusUnauthorized)
				return
			}
			err2 := json.Unmarshal([]byte(s), &sessionData)
			if err2 != nil {
				if h.LogError != nil {
					h.LogError(r.Context(), "error unmarshal: %s ", err2.Error())
				}
				http.Error(w, "Session is expired", http.StatusUnauthorized)
				return
			}
			if id, ok := sessionData["id"]; ok {
				uData := map[string]interface{}{}
				s, err := h.Cache.Get(r.Context(), h.PrefixSessionIndex+id)
				if err != nil {
					http.Error(w, "Session is expired", http.StatusUnauthorized)
					return
				}
				err2 := json.Unmarshal([]byte(s), &uData)
				if err2 != nil {
					if h.LogError != nil {
						h.LogError(r.Context(), "error unmarshal: %s ", err2.Error())
					}
					http.Error(w, "Session is expired", http.StatusInternalServerError)
					return
				}
				ip := getForwardedRemoteIp(r)
				sid, ok := uData["sid"]
				if !ok || sid != sessionId ||
					getValue(uData, "userAgent") != r.UserAgent() ||
					getValue(uData, "ip") != ip {
					http.Error(w, "You cannot use multiple devices with a single account", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Session is expired", http.StatusUnauthorized)
				return
			}

			azureToken := getString(sessionData, "azure_token")
			ctx = context.WithValue(ctx, "azure_token", azureToken)

			authorizationToken := getString(sessionData, "token")
			ctx = context.WithValue(ctx, "token", authorizationToken)
		}
		h.Verify(next, skipRefreshTTL, sessionId).ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *SessionHandler) Verify(next http.Handler, skipRefreshTTL bool, sessionId string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authorizationToken, exists := ctx.Value("token").(string)
		if !exists || authorizationToken == "" {
			http.Error(writer, "invalid authorization token", http.StatusUnauthorized)
			return
		}
		payload, _, _, err := h.VerifyToken(authorizationToken, h.SecretKey)
		if err != nil {
			http.Error(writer, "invalid authorization token", http.StatusUnauthorized)
			return
		}
		ip := getForwardedRemoteIp(r)
		ctx = context.WithValue(ctx, "ip", ip)
		ctx = context.WithValue(ctx, "token", authorizationToken)
		for k, e := range payload {
			if len(k) > 0 {
				ctx = context.WithValue(ctx, k, e)
			}
		}
		if !skipRefreshTTL && sessionId != "" {
			_, err := h.Cache.Expire(ctx, sessionId, h.sessionExpiredTime)
			if err != nil {
				if h.LogError != nil {
					h.LogError(ctx, err.Error())
				}
				http.Error(writer, "error set expire sessionId", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(writer, r.WithContext(ctx))
	})
}

func getForwardedRemoteIp(r *http.Request) string {
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	return ""
}

func getValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		return value.(string)
	}
	return ""
}

func getString(data map[string]string, key string) string {
	if value, ok := data[key]; ok {
		return value
	}
	return ""
}
