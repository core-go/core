package security

import (
	"net/http"
	"sort"
)

type RoleAuthorizer struct {
	Authorization string
	Key           string
	sortedRoles   bool
}

func NewRoleAuthorizer(sortedRoles bool, options ...string) *RoleAuthorizer {
	authorization := ""
	key := "roleId"
	if len(options) >= 2 {
		authorization = options[1]
	}
	if len(options) >= 1 {
		key = options[0]
	}
	return &RoleAuthorizer{sortedRoles: sortedRoles, Authorization: authorization, Key: key}
}

func (h *RoleAuthorizer) Authorize(next http.Handler, roles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRoles := ValuesFromContext(r, h.Authorization, h.Key)
		if userRoles == nil || len(*userRoles) == 0 {
			http.Error(w, "no permission: Require roles for this user", http.StatusForbidden)
			return
		}
		if h.sortedRoles {
			if HasSortedRole(roles, *userRoles) {
				next.ServeHTTP(w, r)
				return
			}
		}
		if HasRole(roles, *userRoles) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "no permission", http.StatusForbidden)
	})
}

func HasRole(roles []string, userRoles []string) bool {
	for _, role := range roles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
}

func HasSortedRole(roles []string, userRoles []string) bool {
	for _, role := range roles {
		i := sort.SearchStrings(userRoles, role)
		if i >= 0 && userRoles[i] == role {
			return true
		}
	}
	return false
}
