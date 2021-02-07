package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
)

type CodeHandler struct {
	Loader         CodeLoader
	Resource       string
	Action         string
	RequiredMaster bool
	LogError       func(context.Context, string)
	WriteLog       func(ctx context.Context, resource string, action string, success bool, desc string) error
}

func NewDefaultCodeHandler(loader CodeLoader, resource string, action string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error) *CodeHandler {
	return NewCodeHandler(loader, resource, action, true, logError, writeLog)
}
func NewCodeHandler(loader CodeLoader, resource string, action string, requiredMaster bool, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error) *CodeHandler {
	if len(resource) == 0 {
		resource = "code"
	}
	if len(action) == 0 {
		action = "load"
	}
	h := CodeHandler{Loader: loader, Resource: resource, Action: action, RequiredMaster: requiredMaster, WriteLog: writeLog, LogError: logError}
	return &h
}
func (c *CodeHandler) Load(w http.ResponseWriter, r *http.Request) {
	code := ""
	if c.RequiredMaster {
		if r.Method == "GET" {
			i := strings.LastIndex(r.RequestURI, "/")
			if i >= 0 {
				code = r.RequestURI[i+1:]
			}
		} else {
			b, er1 := ioutil.ReadAll(r.Body)
			if er1 != nil {
				http.Error(w, "Body cannot is empty", http.StatusBadRequest)
				return
			}
			code = strings.Trim(string(b), " ")
		}
	}
	result, er4 := c.Loader.Load(r.Context(), code)
	if er4 != nil {
		Error(w, r, http.StatusInternalServerError, InternalServerError, c.LogError, c.Resource, c.Action, er4, c.WriteLog)
	} else {
		Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, c.Action)
	}
}
