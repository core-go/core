package echo

import (
	"errors"
	"github.com/labstack/echo/v4"
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

func (h *ArrayAuthorizer) Authorize(array []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			v := FromContext(r, h.Authorization, h.Key)
			if len(v) == 0 {
				ctx.JSON(http.StatusForbidden, "cannot get '" + h.Key + "' from http request")
				return errors.New("cannot get '" + h.Key + "' from http request")
			}
			if len(array) == 0 {
				ctx.JSON(http.StatusForbidden, "no permission")
				return errors.New("no permission")
			}
			if h.sortedUsers {
				if IncludeOfSort(array, v) {
					return next(ctx)
				}
			}
			if Include(array, v) {
				return next(ctx)
			}
			ctx.JSON(http.StatusForbidden, "no permission")
			return errors.New("no permission")
		}
	}
}
