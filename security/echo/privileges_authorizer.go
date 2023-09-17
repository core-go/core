package echo

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
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

func (h *PrivilegesAuthorizer) Authorize(privilegeId string, action int32) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			userId := FromContext(r, h.Authorization, h.Key)
			if len(userId) == 0 {
				ctx.JSON(http.StatusForbidden, "invalid User Id")
				return errors.New("invalid User Id")
			}
			privileges := h.Privileges(r.Context(), userId)
			if privileges == nil || len(privileges) == 0 {
				ctx.JSON(http.StatusForbidden, "no permission: Require privileges for this user")
				return errors.New("no permission: Require privileges for this user")
			}

			privilegeAction := GetAction(privileges, privilegeId, h.sortedPrivilege)
			if privilegeAction == ActionNone {
				ctx.JSON(http.StatusForbidden, "no permission for this user")
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
