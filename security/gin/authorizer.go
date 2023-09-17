package gin

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (h *Authorizer) Authorize(privilegeId string, action int32) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		userId := FromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "invalid User Id in http request")
			return
		}
		p := h.Privilege(r.Context(), userId, privilegeId)
		if p == ActionNone {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission for this user")
			return
		}
		if action == ActionNone || action == ActionAll {
			ctx.Next()
			return
		}
		sum := action & p
		if h.Exact {
			if sum == action {
				ctx.Next()
				return
			}
		} else if sum >= action {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission")
	}
}
