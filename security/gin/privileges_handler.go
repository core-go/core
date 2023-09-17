package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

func (h *PrivilegesHandler) GetPrivileges(ctx *gin.Context) {
	r := ctx.Request
	id := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := ioutil.ReadAll(r.Body)
		if er1 != nil {
			ctx.String(http.StatusBadRequest, "Require id")
			return
		}
		id = strings.Trim(string(b), " ")
	}
	if len(id) == 0 {
		ctx.String(http.StatusBadRequest, "Require id")
		return
	}
	result := h.Privileges(r.Context(), id)
	ctx.JSON(http.StatusOK, result)
}
func (h *PrivilegesHandler) GetPrivilege(ctx *gin.Context) {
	r := ctx.Request
	s := strings.Split(r.RequestURI, "/")
	if len(s) < 3 {
		ctx.String(http.StatusBadRequest, "URL is not valid")
		return
	}

	if r.Method != "GET" {
		ctx.String(http.StatusBadRequest, "Must use GET method")
		return
	}
	userId := s[len(s)-2]
	privilegeId := s[len(s)-1]
	if len(userId) == 0 || len(privilegeId) == 0 {
		ctx.String(http.StatusBadRequest, "parameters cannot be empty")
		return
	}
	result := h.Privilege(r.Context(), userId, privilegeId)
	ctx.JSON(http.StatusOK, result)
}
