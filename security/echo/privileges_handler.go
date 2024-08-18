package echo

import (
	"context"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
)

type PrivilegesHandler struct {
	Privilege  func(ctx context.Context, userId string, privilegeId string) int32
	Privileges func(ctx context.Context, userId string) []string
}

func NewPrivilegeHandler(privilegeLoader func(ctx context.Context, userId string, privilegeId string) int32, privilegesLoader func(ctx context.Context, userId string) []string) *PrivilegesHandler {
	return &PrivilegesHandler{Privilege: privilegeLoader, Privileges: privilegesLoader}
}

func (h *PrivilegesHandler) GetPrivileges(ctx echo.Context) error {
	r := ctx.Request()
	id := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := io.ReadAll(r.Body)
		if er1 != nil {
			return ctx.String(http.StatusBadRequest, "Require id")
		}
		id = strings.Trim(string(b), " ")
	}
	if len(id) == 0 {
		return ctx.String(http.StatusBadRequest, "Require id")
	}
	result := h.Privileges(r.Context(), id)
	return ctx.JSON(http.StatusOK, result)
}
func (h *PrivilegesHandler) GetPrivilege(ctx echo.Context) error {
	r := ctx.Request()
	s := strings.Split(r.RequestURI, "/")
	if len(s) < 3 {
		return ctx.String(http.StatusBadRequest, "URL is not valid")
	}

	if r.Method != "GET" {
		return ctx.String(http.StatusBadRequest, "Must use GET method")
	}
	userId := s[len(s)-2]
	privilegeId := s[len(s)-1]
	if len(userId) == 0 || len(privilegeId) == 0 {
		return ctx.String(http.StatusBadRequest, "parameters cannot be empty")

	}
	result := h.Privilege(r.Context(), userId, privilegeId)
	return ctx.JSON(http.StatusOK, result)
}
