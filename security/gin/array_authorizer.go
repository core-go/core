package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (h *ArrayAuthorizer) Authorize(array []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		v := FromContext(r, h.Authorization, h.Key)
		if len(v) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "cannot get '" + h.Key + "' from http request")
			return
		}
		if len(array) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission")
			return
		}
		if h.sortedUsers {
			if IncludeOfSort(array, v) {
				ctx.Next()
				return
			}
		}
		if Include(array, v) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, "no permission")
	}
}
