package echo

import (
	"errors"
	"github.com/labstack/echo/v4"
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

func (h *TokenAuthorizer) Authorize(privilegeId string, action int32) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			privileges := ValuesFromContext(r, h.Authorization, h.Key)
			if privileges == nil || len(*privileges) == 0 {
				ctx.String(http.StatusForbidden, "no permission: Require privileges for this user")
				return errors.New("no permission: Require privileges for this user")
			}

			privilegeAction := GetAction(*privileges, privilegeId, h.sortedPrivilege)
			if privilegeAction == ActionNone {
				ctx.String(http.StatusForbidden, "no permission for this user")
				return errors.New("no permission for this user")
			}
			if action == ActionNone || action == ActionAll {
				return next(ctx)
			}
			sum := action & privilegeAction
			if h.exact {
				if sum == action {
					return next(ctx)
				}
			} else {
				if sum >= action {
					return next(ctx)
				}
			}
			ctx.JSON(http.StatusForbidden, "no permission")
			return errors.New("no permission")
		}
	}
}
