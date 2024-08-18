package security

import (
	"context"
	"encoding/json"
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

func (h *PrivilegesHandler) GetPrivileges(w http.ResponseWriter, r *http.Request) {
	id := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			id = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := io.ReadAll(r.Body)
		if er1 != nil {
			http.Error(w, "Require id", http.StatusBadRequest)
			return
		}
		id = strings.Trim(string(b), " ")
	}
	if len(id) == 0 {
		http.Error(w, "Require id", http.StatusBadRequest)
		return
	}
	result := h.Privileges(r.Context(), id)
	respond(w, r, http.StatusOK, result)
}
func (h *PrivilegesHandler) GetPrivilege(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.RequestURI, "/")
	if len(s) < 3 {
		http.Error(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Must use GET method", http.StatusBadRequest)
		return
	}
	userId := s[len(s)-2]
	privilegeId := s[len(s)-1]
	if len(userId) == 0 || len(privilegeId) == 0 {
		http.Error(w, "parameters cannot be empty", http.StatusBadRequest)
		return
	}
	result := h.Privilege(r.Context(), userId, privilegeId)
	respond(w, r, http.StatusOK, result)
}
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	return err
}
