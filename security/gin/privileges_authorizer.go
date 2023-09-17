package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PrivilegesAuthorizer struct {
	Privileges      func(ctx context.Context, userId string) []string
	Authorization   string
	Key             string
	sortedPrivilege bool
	exact           bool
}

func NewPrivilegesAuthorizer(loadPrivileges func(ctx context.Context, userId string) []string, sortedPrivilege bool, exact bool, key string, options ...string) *PrivilegesAuthorizer {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &PrivilegesAuthorizer{Privileges: loadPrivileges, Authorization: authorization, Key: key, sortedPrivilege: sortedPrivilege, exact: exact}
}

func (h *PrivilegesAuthorizer) Authorize(privilegeId string, action int32) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		userId := FromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "Invalid User Id")
			return
		}
		privileges := h.Privileges(r.Context(), userId)
		if privileges == nil || len(privileges) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission: Require privileges for this user")
			return
		}

		privilegeAction := GetAction(privileges, privilegeId, h.sortedPrivilege)
		if privilegeAction == ActionNone {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission for this user")
			return
		}
		if action == ActionNone || action == ActionAll {
			ctx.Next()
			return
		}
		sum := action & privilegeAction
		if h.exact {
			if sum == action {
				ctx.Next()
				return
			}
		} else {
			if sum >= action {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission")
	}
}
