package gin

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (h *SubAuthorizer) Authorize(privilegeId string, sub string, action int32) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		userId := FromContext(r, h.Authorization, h.Key)
		if len(userId) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "invalid User Id in http request")
			return
		}
		p := h.Privilege(r.Context(), userId, privilegeId, sub)
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
