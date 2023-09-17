package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (h *TokenAuthorizer) Authorize(privilegeId string, action int32) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		privileges := ValuesFromContext(r, h.Authorization, h.Key)
		if privileges == nil || len(*privileges) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission: Require privileges for this user")
			return
		}

		privilegeAction := GetAction(*privileges, privilegeId, h.sortedPrivilege)
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
