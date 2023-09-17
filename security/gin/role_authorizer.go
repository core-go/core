package gin

import (
	"github.com/gin-gonic/gin"
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

func (h *RoleAuthorizer) Authorize(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		userRoles := ValuesFromContext(r, h.Authorization, h.Key)
		if userRoles == nil || len(*userRoles) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission: Require roles for this user")
			return
		}
		if h.sortedRoles {
			if HasSortedRole(roles, *userRoles) {
				ctx.Next()
				return
			}
		}
		if HasRole(roles, *userRoles) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission")
	}
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
