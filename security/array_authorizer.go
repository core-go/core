package security

import "net/http"

type ArrayAuthorizer struct {
	Authorization string
	Key           string
	sortedUsers   bool
}

func NewArrayAuthorizer(sortedUsers bool, options ...string) *ArrayAuthorizer {
	authorization := ""
	key := "userId"
	if len(options) >= 2 {
		authorization = options[1]
	}
	if len(options) >= 1 && len(options[0]) > 0 {
		key = options[0]
	}
	return &ArrayAuthorizer{sortedUsers: sortedUsers, Authorization: authorization, Key: key}
}

func (h *ArrayAuthorizer) Authorize(next http.Handler, array []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := FromContext(r, h.Authorization, h.Key)
		if len(v) == 0 {
			http.Error(w, "cannot get '" + h.Key + "' from http request", http.StatusForbidden)
			return
		}
		if len(array) == 0 {
			http.Error(w, "no permission", http.StatusForbidden)
			return
		}
		if h.sortedUsers {
			if IncludeOfSort(array, v) {
				next.ServeHTTP(w, r)
				return
			}
		}
		if Include(array, v) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "no permission", http.StatusForbidden)
	})
}
