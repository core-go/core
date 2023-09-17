package security

import (
	"context"
	"net/http"
)

type Authorizer struct {
	Privilege     func(ctx context.Context, userId string, privilegeId string) int32
	Authorization string
	Key           string
	Exact         bool
}

func NewAuthorizer(loadPrivilege func(context.Context, string, string) int32, exact bool, options ...string) *Authorizer {
	authorization := ""
	key := "userId"
	if len(options) >= 2 {
		authorization = options[1]
	}
	if len(options) >= 1 {
		key = options[0]
	}
	return &Authorizer{Privilege: loadPrivilege, Exact: exact, Authorization: authorization, Key: key}
}

func (h *Authorizer) Authorize(next http.Handler, privilegeId string, action int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := FromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			http.Error(w, "invalid User Id in http request", http.StatusForbidden)
			return
		}
		p := h.Privilege(r.Context(), userId, privilegeId)
		if p == ActionNone {
			http.Error(w, "no permission for this user", http.StatusForbidden)
			return
		}
		if action == ActionNone || action == ActionAll {
			next.ServeHTTP(w, r)
			return
		}

		if h.Exact {
			sum := action & p
			if sum == action {
				next.ServeHTTP(w, r)
				return
			}
		} else if p >= action {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "no permission", http.StatusForbidden)
	})
}
