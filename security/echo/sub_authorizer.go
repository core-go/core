package echo

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
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

func (h *SubAuthorizer) Authorize(privilegeId string, sub string, action int32) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			userId := FromContext(r, h.Authorization, h.Key)
			if len(userId) == 0 {
				ctx.JSON(http.StatusForbidden, "invalid User Id in http request")
				return errors.New("invalid User Id in http request")
			}
			p := h.Privilege(r.Context(), userId, privilegeId, sub)
			if p == ActionNone {
				ctx.JSON(http.StatusForbidden, "no permission for this user")
				return errors.New("no permission for this user")
			}
			if action == ActionNone || action == ActionAll {
				return next(ctx)
			}
			sum := action & p
			if h.Exact {
				if sum == action {
					return next(ctx)
				}
			} else if sum >= action {
				return next(ctx)
			}
			ctx.JSON(http.StatusForbidden, "no permission")
			return errors.New("no permission")
		}
	}
}
