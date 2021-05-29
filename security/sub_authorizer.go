package security

import (
	"context"
	"net/http"
)

type SubAuthorizer struct {
	Privilege     func(ctx context.Context, userId string, privilegeId string, sub string) int32
	Authorization string
	Key           string
	Exact         bool
}

func NewSubAuthorizer(loadPrivilege func(context.Context, string, string, string) int32, exact bool, options ...string) *SubAuthorizer {
	authorization := ""
	key := "userId"
	if len(options) >= 2 {
		authorization = options[1]
	}
	if len(options) >= 1 {
		key = options[0]
	}
	return &SubAuthorizer{Privilege: loadPrivilege, Exact: exact, Authorization: authorization, Key: key}
}

func (h *SubAuthorizer) Authorize(next http.Handler, privilegeId string, sub string, action int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := FromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			http.Error(w, "invalid User Id in http request", http.StatusForbidden)
			return
		}
		p := h.Privilege(r.Context(), userId, privilegeId, sub)
		if p == ActionNone {
			http.Error(w, "no permission for this user", http.StatusForbidden)
			return
		}
		if action == ActionNone || action == ActionAll {
			next.ServeHTTP(w, r)
			return
		}
		sum := action & p
		if h.Exact {
			if sum == action {
				next.ServeHTTP(w, r)
				return
			}
		} else if sum >= action {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "no permission", http.StatusForbidden)
	})
}
