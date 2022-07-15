package ldap

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const internalServerError = "Internal Server Error"

type LdapInfoHandler struct {
	Load     func(ctx context.Context, id string) (map[string]interface{}, error)
	Resource string
	Action   string
	Error    func(context.Context, string)
	Log      func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewLdapInfoHandler(load func(context.Context, string) (map[string]interface{}, error), logError func(context.Context, string), options ...func(context.Context, string, string, bool, string) error) *LdapInfoHandler {
	var writeLog func(context.Context, string, string, bool, string) error
	if len(options) > 0 && options[0] != nil {
		writeLog = options[0]
	}
	return NewLdapInfoHandlerWithLog(load, logError, writeLog)
}
func NewLdapInfoHandlerWithLog(load func(context.Context, string) (map[string]interface{}, error), logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, options ...string) *LdapInfoHandler {
	var resource, action string
	if len(options) > 0 && len(options[0]) > 0 {
		action = options[0]
	} else {
		action = "load"
	}
	if len(options) > 1 && len(options[1]) > 0 {
		resource = options[1]
	} else {
		resource = "ldap"
	}
	h := LdapInfoHandler{Load: load, Resource: resource, Action: action, Error: logError, Log: writeLog}
	return &h
}
func (h *LdapInfoHandler) GetLdapInfo(w http.ResponseWriter, r *http.Request) {
	uid := ""
	if r.Method == "GET" {
		s := r.URL.Path
		i := strings.LastIndex(s, "/")
		if i >= 0 {
			uid = s[i+1:]
		}
	} else {
		b, er1 := ioutil.ReadAll(r.Body)
		if er1 != nil {
			http.Error(w, "Body cannot is empty", http.StatusBadRequest)
			return
		}
		uid = strings.Trim(string(b), " ")
	}
	result, err := h.Load(r.Context(), uid)
	if err != nil {
		if h.Error != nil {
			h.Error(r.Context(), err.Error())
		}
		respond(w, r, http.StatusInternalServerError, internalServerError, h.Log, h.Resource, h.Action, false, err.Error())
	} else {
		if result == nil {
			respond(w, r, http.StatusNotFound, result, h.Log, h.Resource, h.Action, false, "Not Found")
		} else {
			respond(w, r, http.StatusOK, result, h.Log, h.Resource, h.Action, true, "")
		}
	}
}
func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}, writeLog func(context.Context, string, string, bool, string) error, resource string, action string, success bool, desc string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if writeLog != nil {
		writeLog(r.Context(), resource, action, success, desc)
	}
	return err
}
