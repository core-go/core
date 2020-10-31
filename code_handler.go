package service

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type CodeHandler struct {
	Loader    CodeLoader
	Resource  string
	Action    string
	LogWriter LogWriter
}

func NewCodeHandler(loader CodeLoader, resource string, action string, logWriter LogWriter) *CodeHandler {
	if len(resource) == 0 {
		resource = "code"
	}
	if len(action) == 0 {
		action = "load"
	}
	h := CodeHandler{Loader: loader, Resource: resource, Action: action, LogWriter: logWriter}
	return &h
}
func (c *CodeHandler) Load(w http.ResponseWriter, r *http.Request) {
	code := ""
	if r.Method == "GET" {
		i := strings.LastIndex(r.RequestURI, "/")
		if i >= 0 {
			code = r.RequestURI[i+1:]
		}
	} else {
		b, er1 := ioutil.ReadAll(r.Body)
		if er1 != nil {
			RespondString(w, r, http.StatusBadRequest, "Cannot get the body of 'Forgot Password'")
			return
		}
		code = strings.Trim(string(b), " ")
	}
	result, er4 := c.Loader.Load(r.Context(), code)
	if er4 != nil {
		Respond(w, r, http.StatusInternalServerError, InternalServerError, c.LogWriter, c.Resource, c.Action, false, er4.Error())
	} else {
		Respond(w, r, http.StatusOK, result, c.LogWriter, c.Resource, c.Action, true, "")
	}
}
