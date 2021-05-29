package security

import "net/http"

type TokenAuthorizer struct {
	Authorization   string
	Key             string
	sortedPrivilege bool
	exact           bool
}

func NewTokenAuthorizer(sortedPrivilege bool, exact bool, key string, options ...string) *TokenAuthorizer {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &TokenAuthorizer{Authorization: authorization, Key: key, sortedPrivilege: sortedPrivilege, exact: exact}
}

func (h *TokenAuthorizer) Authorize(next http.Handler, privilegeId string, action int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		privileges := ValuesFromContext(r, h.Authorization, h.Key)
		if privileges == nil || len(*privileges) == 0 {
			http.Error(w, "no permission: Require privileges for this user", http.StatusForbidden)
			return
		}

		privilegeAction := GetAction(*privileges, privilegeId, h.sortedPrivilege)
		if privilegeAction == ActionNone {
			http.Error(w, "no permission for this user", http.StatusForbidden)
			return
		}
		if action == ActionNone || action == ActionAll {
			next.ServeHTTP(w, r)
			return
		}
		sum := action & privilegeAction
		if h.exact {
			if sum == action {
				next.ServeHTTP(w, r)
				return
			}
		} else {
			if sum >= action {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "no permission", http.StatusForbidden)
	})
}
